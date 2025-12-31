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
	"github.com/ethereum-optimism/optimism/op-service/logpipe"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

// L2SuccinctFaultProofChallenger wraps the OP Succinct fault-proof challenger binary as a subprocess.
type L2SuccinctFaultProofChallenger struct {
	mu                 sync.Mutex
	id                 stack.L2ChallengerID
	execPath           string
	args               []string
	p                  devtest.P
	logger             log.Logger
	sub                *SubProcess
	l2MetricsRegistrar L2MetricsRegistrar
	metricsPort        string
}

var _ L2ChallengerBackend = (*L2SuccinctFaultProofChallenger)(nil)

// FaultProofChallenger extends L2ChallengerBackend with faultproof-specific methods.
type FaultProofChallenger interface {
	L2ChallengerBackend
	Start()
	Stop()
}

var _ FaultProofChallenger = (*L2SuccinctFaultProofChallenger)(nil)

func (c *L2SuccinctFaultProofChallenger) hydrate(system stack.ExtensibleSystem) {
	bFrontend := shim.NewL2Challenger(shim.L2ChallengerConfig{
		CommonConfig: shim.NewCommonConfig(system.T()),
		ID:           c.id,
		Config:       nil, // Succinct challenger runs as subprocess, no op-challenger config
	})
	l2Net := system.L2Network(stack.L2NetworkID(c.id.ChainID()))
	l2Net.(stack.ExtensibleL2Network).AddL2Challenger(bFrontend)
}

// Start starts the fault-proof challenger subprocess.
func (c *L2SuccinctFaultProofChallenger) Start() {
	c.mu.Lock()
	if c.sub != nil {
		c.logger.Warn("Fault Proof Challenger already started")
		c.mu.Unlock()
		return
	}

	// We pipe sub-process logs to the test-logger.
	logOut := logpipe.ToLogger(c.logger.New("src", "stdout"))
	logErr := logpipe.ToLogger(c.logger.New("src", "stderr"))

	stdOutLogs := logpipe.LogProcessor(func(line []byte) {
		e := logpipe.ParseRustStructuredLogs(line)
		logOut(e)
	})
	stdErrLogs := logpipe.LogProcessor(func(line []byte) {
		e := logpipe.ParseRustStructuredLogs(line)
		logErr(e)
	})
	c.sub = NewSubProcess(c.p, stdOutLogs, stdErrLogs)
	c.mu.Unlock()

	c.sub.OnExit(func(err error) {
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

		c.p.Require().NoError(err, "fault-proof challenger exited unexpectedly")
	})

	err := c.sub.Start(c.execPath, c.args, []string{})
	c.p.Require().NoError(err, "Must start challenger")

	if c.metricsPort != "" && c.l2MetricsRegistrar != nil {
		metricsTarget := NewPrometheusMetricsTarget("localhost", c.metricsPort, false)
		c.l2MetricsRegistrar.RegisterL2MetricsTargets(c.id, metricsTarget)
		c.logger.Info("Registered fault-proof challenger metrics", "port", c.metricsPort)
	}
}

// Stop stops the fault-proof challenger subprocess.
func (c *L2SuccinctFaultProofChallenger) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sub == nil {
		c.logger.Warn("fault-proof challenger already stopped")
		return
	}

	err := c.sub.Stop(true)
	c.p.Require().NoError(err, "Must stop challenger")
	c.sub = nil
}

// WithSuccinctFaultProofChallenger creates a fault-proof challenger after deployment.
func WithSuccinctFaultProofChallenger(challengerID stack.L2ChallengerID, l1ELID stack.L1ELNodeID, l2ELID stack.L2ELNodeID, opts ...FaultProofChallengerOption) stack.Option[*Orchestrator] {
	return stack.AfterDeploy(func(orch *Orchestrator) {
		WithSuccinctFaultProofChallengerPostDeploy(orch, challengerID, l1ELID, l2ELID, opts...)
	})
}

// WithSuperSuccinctFaultProofChallenger creates a fault-proof challenger in the Finally phase.
func WithSuperSuccinctFaultProofChallenger(challengerID stack.L2ChallengerID,
	l1ELID stack.L1ELNodeID, l2ELID stack.L2ELNodeID, opts ...FaultProofChallengerOption) stack.Option[*Orchestrator] {
	return stack.Finally(func(orch *Orchestrator) {
		WithSuccinctFaultProofChallengerPostDeploy(orch, challengerID, l1ELID, l2ELID, opts...)
	})
}

// WithSuccinctFaultProofChallengerPostDeploy sets up and starts the OP Succinct fault-proof challenger.
func WithSuccinctFaultProofChallengerPostDeploy(orch *Orchestrator, challengerID stack.L2ChallengerID, l1ELID stack.L1ELNodeID, l2ELID stack.L2ELNodeID, opts ...FaultProofChallengerOption) {
	ctx := stack.ContextWithID(orch.P().Ctx(), challengerID)
	p := orch.P().WithCtx(ctx)
	logger := p.Logger().New("component", "succinct-faultproof-challenger")

	require := p.Require()
	require.False(orch.challengers.Has(challengerID), "challenger must not already exist")

	l2Net, ok := orch.l2Nets.Get(challengerID.ChainID())
	require.True(ok, "l2 network required")

	l1EL, ok := orch.GetL1EL(l1ELID)
	require.True(ok, "l1 EL node required")

	l2EL, ok := orch.GetL2EL(l2ELID)
	require.True(ok, "l2 EL node required")

	// Use ChallengerRole for the challenger key
	challengerKey, err := orch.GetKeys().Secret(devkeys.ChallengerRole.Key(challengerID.ChainID().ToBig()))
	require.NoError(err, "failed to get challenger key")
	challengerKeyStr := hexutil.Encode(crypto.FromECDSA(challengerKey))

	cfg := &FaultProofChallengerConfig{}
	for _, opt := range opts {
		opt(p, challengerID, cfg)
	}

	l1RPC := l1EL.UserRPC()
	l2RPC := strings.ReplaceAll(l2EL.UserRPC(), "ws://", "http://")
	anchorStateRegistryAddr := l2Net.deployment.anchorStateRegistry
	factoryAddr := l2Net.deployment.disputeGameFactoryProxy

	logger.Info("L1_RPC", "url", l1RPC)
	logger.Info("L2_RPC", "url", l2RPC)
	logger.Info("ANCHOR_STATE_REGISTRY_ADDRESS", "address", anchorStateRegistryAddr)
	logger.Info("FACTORY_ADDRESS", "address", factoryAddr)

	envVars := map[string]string{
		"L1_RPC":                        l1RPC,
		"L2_RPC":                        l2RPC,
		"ANCHOR_STATE_REGISTRY_ADDRESS": anchorStateRegistryAddr.String(),
		"FACTORY_ADDRESS":               factoryAddr.String(),
		"GAME_TYPE":                     "42",
		"PRIVATE_KEY":                   challengerKeyStr,
		"LOG_FORMAT":                    "json",
	}

	// Optional parameters (override defaults if set)
	setEnvIfNotNil(envVars, "FETCH_INTERVAL", cfg.fetchInterval)
	setEnvIfNotNil(envVars, "MALICIOUS_CHALLENGE_PERCENTAGE", cfg.maliciousChallengePercentage)
	setEnvIfNotNil(envVars, "RUST_LOG", cfg.rustLog)

	var metricsPort string
	if areMetricsEnabled() {
		metricsPort, err = getAvailableLocalPort()
		require.NoError(err, "failed to get available port for challenger metrics")
		envVars["CHALLENGER_METRICS_PORT"] = metricsPort
	}

	envDir := p.TempDir()
	envFile := filepath.Join(envDir, fmt.Sprintf("fault-proof-challenger-%s.env", challengerID.String()))
	err = WriteEnvFile(envFile, envVars)
	p.Require().NoError(err, "must write fault proof challenger env file")

	if cfg.envFilePath != nil {
		err = WriteEnvFile(*cfg.envFilePath, envVars)
		p.Require().NoError(err, "must write challenger env file")
		logger.Info("challenger env file written", "path", *cfg.envFilePath)
	}

	execPath := os.Getenv("FAULT_PROOF_CHALLENGER_EXEC_PATH")
	p.Require().NotEmpty(execPath, "FAULT_PROOF_CHALLENGER_EXEC_PATH environment variable must be set")
	_, err = os.Stat(execPath)
	p.Require().NotErrorIs(err, os.ErrNotExist, "challenger executable must exist")

	c := &L2SuccinctFaultProofChallenger{
		id:                 challengerID,
		execPath:           execPath,
		args:               []string{"--env-file", envFile},
		p:                  p,
		logger:             logger,
		l2MetricsRegistrar: orch,
		metricsPort:        metricsPort,
	}
	logger.Info("Starting fault-proof challenger")
	c.Start()
	p.Cleanup(func() {
		logger.Info("Stopping fault-proof challenger")
		c.Stop()
	})
	logger.Info("fault-proof challenger is running")

	// Store the challenger in the orchestrator's challengers map
	require.True(orch.challengers.SetIfMissing(challengerID, c), "challenger must not already exist")
}

// FaultProofChallengerConfig holds configuration for the OP Succinct fault-proof challenger.
type FaultProofChallengerConfig struct {
	fetchInterval                *uint64
	maliciousChallengePercentage *float64
	rustLog                      *string
	envFilePath                  *string
}

// FaultProofChallengerOption is a function that configures the FaultProofChallengerConfig.
type FaultProofChallengerOption func(p devtest.P, id stack.L2ChallengerID, cfg *FaultProofChallengerConfig)

// WithFPChallengerFetchInterval sets the polling interval in seconds.
func WithFPChallengerFetchInterval(n uint64) FaultProofChallengerOption {
	return FaultProofChallengerOption(func(p devtest.P, id stack.L2ChallengerID, cfg *FaultProofChallengerConfig) {
		cfg.fetchInterval = &n
	})
}

// WithFPChallengerMaliciousChallengePercentage sets the percentage of valid games to challenge maliciously (for testing).
func WithFPChallengerMaliciousChallengePercentage(pct float64) FaultProofChallengerOption {
	return FaultProofChallengerOption(func(p devtest.P, id stack.L2ChallengerID, cfg *FaultProofChallengerConfig) {
		cfg.maliciousChallengePercentage = &pct
	})
}

// WithFPChallengerRustLog sets the RUST_LOG environment variable.
func WithFPChallengerRustLog(level string) FaultProofChallengerOption {
	return FaultProofChallengerOption(func(p devtest.P, id stack.L2ChallengerID, cfg *FaultProofChallengerConfig) {
		cfg.rustLog = &level
	})
}

// WithFPChallengerWriteEnvFile enables writing environment variables to a file.
func WithFPChallengerWriteEnvFile(path string) FaultProofChallengerOption {
	return FaultProofChallengerOption(func(p devtest.P, id stack.L2ChallengerID, cfg *FaultProofChallengerConfig) {
		cfg.envFilePath = &path
	})
}
