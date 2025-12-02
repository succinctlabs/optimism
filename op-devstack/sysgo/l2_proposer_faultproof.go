package sysgo

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/ethereum-optimism/optimism/op-chain-ops/devkeys"
	"github.com/ethereum-optimism/optimism/op-devstack/devtest"
	"github.com/ethereum-optimism/optimism/op-devstack/shim"
	"github.com/ethereum-optimism/optimism/op-devstack/stack"
	ps "github.com/ethereum-optimism/optimism/op-proposer/proposer"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/logpipe"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type L2SuccinctFaultProofProposer struct {
	mu                 sync.Mutex
	id                 stack.L2ProposerID
	service            *ps.ProposerService
	userRPC            string
	execPath           string
	args               []string
	p                  devtest.P
	sub                *SubProcess
	l2MetricsRegistrar L2MetricsRegistrar
}

var _ L2Prop = (*L2SuccinctFaultProofProposer)(nil)

func (p *L2SuccinctFaultProofProposer) hydrate(system stack.ExtensibleSystem) {
	require := system.T().Require()
	rpcCl, err := client.NewRPC(system.T().Ctx(), system.Logger(), p.userRPC, client.WithLazyDial())
	require.NoError(err)
	system.T().Cleanup(rpcCl.Close)

	bFrontend := shim.NewL2Proposer(shim.L2ProposerConfig{
		CommonConfig: shim.NewCommonConfig(system.T()),
		ID:           p.id,
		Client:       rpcCl,
	})
	l2Net := system.L2Network(stack.L2NetworkID(p.id.ChainID()))
	l2Net.(stack.ExtensibleL2Network).AddL2Proposer(bFrontend)
}

func (k *L2SuccinctFaultProofProposer) UserRPC() string {
	return k.userRPC
}

func (k *L2SuccinctFaultProofProposer) Start() {
	k.mu.Lock()
	if k.sub != nil {
		k.p.Logger().Warn("Fault Proof Proposer already started")
		k.mu.Unlock()
		return
	}

	// We pipe sub-process logs to the test-logger.
	// And inspect them along the way, to get the RPC server address.
	logOut := logpipe.ToLogger(k.p.Logger().New("component", "fault-proof", "src", "stdout"))
	logErr := logpipe.ToLogger(k.p.Logger().New("component", "fault-proof", "src", "stderr"))

	stdOutLogs := logpipe.LogProcessor(func(line []byte) {
		e := logpipe.ParseRustStructuredLogs(line)
		logOut(e)
	})
	stdErrLogs := logpipe.LogProcessor(func(line []byte) {
		e := logpipe.ParseRustStructuredLogs(line)
		logErr(e)
	})
	k.sub = NewSubProcess(k.p, stdOutLogs, stdErrLogs)
	k.mu.Unlock()

	k.sub.OnExit(func(err error) {
		if errors.Is(err, syscall.ECHILD) {
			return
		}

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				sig := ws.Signal()
				if sig == syscall.SIGINT || sig == syscall.SIGTERM {
					return
				}
			}
		}

		k.p.Require().NoError(err, "fault-proof proposer exited unexpectedly")
	})

	err := k.sub.Start(k.execPath, k.args, []string{})
	k.p.Require().NoError(err, "Must start")
}

// Stops the fault-proof proposer.
func (k *L2SuccinctFaultProofProposer) Stop() {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.sub == nil {
		k.p.Logger().Warn("fault-proof proposer already stopped")
		return
	}

	err := k.sub.Stop(true)
	k.p.Require().NoError(err, "Must stop")
	k.sub = nil
}

func WithSuccinctFaultProofProposer(proposerID stack.L2ProposerID, l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2CLID stack.L2CLNodeID, l2ELID stack.L2ELNodeID, opts ...FaultProofProposerOption) stack.Option[*Orchestrator] {
	return stack.AfterDeploy(func(orch *Orchestrator) {
		WithSuccinctFaultProofProposerPostDeploy(orch, proposerID, l1CLID, l1ELID, l2CLID, l2ELID, opts...)
	})
}

func WithSuperSuccinctFaultProofProposer(proposerID stack.L2ProposerID,
	l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2CLID stack.L2CLNodeID, l2ELID stack.L2ELNodeID, opts ...FaultProofProposerOption) stack.Option[*Orchestrator] {
	return stack.Finally(func(orch *Orchestrator) {
		WithSuccinctFaultProofProposerPostDeploy(orch, proposerID, l1CLID, l1ELID, l2CLID, l2ELID, opts...)
	})
}

func WithSuccinctFaultProofProposerPostDeploy(orch *Orchestrator, proposerID stack.L2ProposerID, l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2CLID stack.L2CLNodeID, l2ELID stack.L2ELNodeID, opts ...FaultProofProposerOption) {
	ctx := stack.ContextWithID(orch.P().Ctx(), proposerID)
	p := orch.P().WithCtx(ctx)
	logger := p.Logger()

	require := p.Require()
	require.False(orch.proposers.Has(proposerID), "proposer must not already exist")

	l2Net, ok := orch.l2Nets.Get(proposerID.ChainID())
	require.True(ok, "l2 network required")

	l1EL, ok := orch.l1ELs.Get(l1ELID)
	require.True(ok, "l1 EL node required")

	l1CL, ok := orch.l1CLs.Get(l1CLID)
	require.True(ok, "l1 CL node required")

	l2EL, ok := orch.l2ELs.Get(l2ELID)
	require.True(ok, "l2 EL node required")

	l2CL, ok := orch.l2CLs.Get(l2CLID)
	require.True(ok, "l2 CL node required")

	proposerKey, err := orch.keys.Secret(devkeys.ProposerRole.Key(proposerID.ChainID().ToBig()))
	require.NoError(err)
	proposerKeyStr := hexutil.Encode(crypto.FromECDSA(proposerKey))

	cfg := &FaultProofProposerConfig{}
	orch.proposerOptions.Apply(p, proposerID, cfg)
	for _, opt := range opts {
		opt(p, proposerID, cfg)
	}

	require.NotEmpty(cfg.l1ConfigDir, "fault-proof proposer L1 config dir must be set")
	require.NotEmpty(cfg.l2ConfigDir, "fault-proof proposer L2 config dir must be set")

	l1RPC := l1EL.UserRPC()
	l1BeaconRPC := l1CL.beaconHTTPAddr
	l2RPC := strings.ReplaceAll(l2EL.UserRPC(), "ws://", "http://")
	l2NodeRPC := strings.ReplaceAll(l2CL.UserRPC(), "ws://", "http://")
	mockVerifierAddr := l2Net.deployment.sp1MockVerifier
	disputeGameFactoryProxy := l2Net.deployment.disputeGameFactoryProxy

	logger.Info("L1_RPC", "url", l1RPC)
	logger.Info("L1_BEACON_RPC", "url", l1BeaconRPC)
	logger.Info("L2_RPC", "url", l2RPC)
	logger.Info("L2_NODE_RPC", "url", l2NodeRPC)
	logger.Info("SP1MockVerifier", "address", mockVerifierAddr)
	logger.Info("DisputeGameFactory", "address", disputeGameFactoryProxy)

	envVars := map[string]string{
		"L1_RPC":           l1RPC,
		"L1_BEACON_RPC":    l1BeaconRPC,
		"L2_RPC":           l2RPC,
		"L2_NODE_RPC":      l2NodeRPC,
		"VERIFIER_ADDRESS": mockVerifierAddr.String(),
		"FACTORY_ADDRESS":  disputeGameFactoryProxy.String(),
		"GAME_TYPE":        "42",
		"MOCK_MODE":        "true",
		"PRIVATE_KEY":      proposerKeyStr,
		"L1_CONFIG_DIR":    cfg.l1ConfigDir,
		"L2_CONFIG_DIR":    cfg.l2ConfigDir,
		"LOG_FORMAT":       "json",
	}

	setEnvFromEnvOrDefault(envVars, "NETWORK_PRIVATE_KEY", "")

	// Optional parameters (override defaults if set)
	setEnvIfNotNil(envVars, "PROPOSAL_INTERVAL_IN_BLOCKS", cfg.proposalIntervalInBlocks)
	setEnvIfNotNil(envVars, "FETCH_INTERVAL", cfg.fetchInterval)
	setEnvIfNotNil(envVars, "FAST_FINALITY_MODE", cfg.fastFinalityMode)
	setEnvIfNotNil(envVars, "FAST_FINALITY_PROVING_LIMIT", cfg.fastFinalityProvingLimit)
	setEnvIfNotNil(envVars, "RANGE_SPLIT_COUNT", cfg.rangeSplitCount)
	setEnvIfNotNil(envVars, "MAX_CONCURRENT_RANGE_PROOFS", cfg.maxConcurrentRangeProofs)
	setEnvIfNotNil(envVars, "MOCK_MODE", cfg.mockMode)
	setEnvIfNotNil(envVars, "RUST_LOG", cfg.rustLog)

	if areMetricsEnabled() {
		metricsPort, err := getAvailableLocalPort()
		p.Require().NoError(err, "must get available port for metrics")
		setEnvFromEnvOrDefault(envVars, "FAULT_PROOF_PROPOSER_METRICS_PORT", metricsPort)
		envVars["FAULT_PROOF_PROPOSER_METRICS_ENABLED"] = "true"
	}

	envDir := p.TempDir()
	envFile := filepath.Join(envDir, fmt.Sprintf("fault-proof-proposer-%s.env", proposerID.String()))
	err = writeEnvFile(envFile, envVars)
	p.Require().NoError(err, "must write fault proof proposer env file")

	execPath := os.Getenv("FAULT_PROOF_PROPOSER_EXEC_PATH")
	p.Require().NotEmpty(execPath, "FAULT_PROOF_PROPOSER_EXEC_PATH environment variable must be set")
	_, err = os.Stat(execPath)
	p.Require().NotErrorIs(err, os.ErrNotExist, "executable must exist")

	k := &L2SuccinctFaultProofProposer{
		id:                 proposerID,
		userRPC:            "", // retrieved from logs
		execPath:           execPath,
		args:               []string{"--env-file", envFile},
		p:                  p,
		l2MetricsRegistrar: orch,
	}
	p.Logger().Info("Starting fault-proof proposer")
	k.Start()
	p.Cleanup(func() {
		logger.Info("Stopping fault-proof proposer")
		k.Stop()
	})
	p.Logger().Info("fault-proof proposer is running", "rpc", k.UserRPC())
	require.True(orch.proposers.SetIfMissing(proposerID, k), "must not already exist")
}

type FaultProofProposerConfig struct {
	l1ConfigDir              string
	l2ConfigDir              string
	proposalIntervalInBlocks *uint64
	fastFinalityMode         *bool
	fastFinalityProvingLimit *uint64
	rangeSplitCount          *uint64
	maxConcurrentRangeProofs *uint64
	fetchInterval            *uint64
	mockMode                 *bool
	rustLog                  *string
}

type FaultProofProposerOption = ProposerOption[FaultProofProposerConfig]

func WithFPConfigDirsOption(
	o *Orchestrator,
	l1Dir, l2Dir string,
) {
	AppendProposerOption(o, FaultProofProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *FaultProofProposerConfig) {
		cfg.l1ConfigDir = l1Dir
		cfg.l2ConfigDir = l2Dir
	},
	))
}

func WithFPProposalIntervalInBlocks(n uint64) FaultProofProposerOption {
	return FaultProofProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *FaultProofProposerConfig) {
		cfg.proposalIntervalInBlocks = &n
	},
	)
}

func WithFPFetchInterval(n uint64) FaultProofProposerOption {
	return FaultProofProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *FaultProofProposerConfig) {
		cfg.fetchInterval = &n
	},
	)
}

func WithFPFastFinalityMode(enabled bool) FaultProofProposerOption {
	return FaultProofProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *FaultProofProposerConfig) {
		cfg.fastFinalityMode = &enabled
	},
	)
}

func WithFPFastFinalityProvingLimit(n uint64) FaultProofProposerOption {
	return FaultProofProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *FaultProofProposerConfig) {
		cfg.fastFinalityProvingLimit = &n
	},
	)
}

func WithFPRangeSplitCount(n uint64) FaultProofProposerOption {
	return FaultProofProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *FaultProofProposerConfig) {
		cfg.rangeSplitCount = &n
	},
	)
}

func WithFPMaxConcurrentRangeProofs(n uint64) FaultProofProposerOption {
	return FaultProofProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *FaultProofProposerConfig) {
		cfg.maxConcurrentRangeProofs = &n
	},
	)
}

func WithFPMockMode(enabled bool) FaultProofProposerOption {
	return FaultProofProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *FaultProofProposerConfig) {
		cfg.mockMode = &enabled
	},
	)
}

func WithFPRustLog(level string) FaultProofProposerOption {
	return FaultProofProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *FaultProofProposerConfig) {
		cfg.rustLog = &level
	},
	)
}
