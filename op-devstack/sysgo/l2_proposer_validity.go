package sysgo

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	"github.com/ethereum/go-ethereum/log"
	embeddedpg "github.com/fergusstrange/embedded-postgres"
)

type L2SuccinctValidityProposer struct {
	mu                 sync.Mutex
	id                 stack.L2ProposerID
	service            *ps.ProposerService
	userRPC            string
	execPath           string
	args               []string
	p                  devtest.P
	logger             log.Logger
	sub                *SubProcess
	databaseURL        string
	l2MetricsRegistrar L2MetricsRegistrar
}

var _ L2Prop = (*L2SuccinctValidityProposer)(nil)

// ValidityProposer extends L2Prop with validity-specific methods.
type ValidityProposer interface {
	L2Prop
	Start()
	Stop()
	DatabaseURL() string
}

var _ ValidityProposer = (*L2SuccinctValidityProposer)(nil)

func (p *L2SuccinctValidityProposer) hydrate(system stack.ExtensibleSystem) {
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

func (k *L2SuccinctValidityProposer) UserRPC() string {
	return k.userRPC
}

func (k *L2SuccinctValidityProposer) DatabaseURL() string {
	return k.databaseURL
}

type ValidityProposerConfig struct {
	l1ConfigDir                string
	l2ConfigDir                string
	submissionInterval         *uint64
	rangeProofInterval         *uint64
	rangeProofEvmGasLimit      *uint64
	maxConcurrentProofRequests *uint64
	maxConcurrentWitnessGen    *uint64
	loopInterval               *uint64
	provingTimeout             *uint64
	opSuccinctConfigName       *string
	mockMode                   *bool
	rustLog                    *string
	envFilePath                *string
}

type ValidityProposerOption = ProposerOption[ValidityProposerConfig]

func WithValidityConfigDirsOption(
	o *Orchestrator,
	l1Dir, l2Dir string,
) {
	AppendProposerOption(o, ValidityProposerOption(
		func(p devtest.P, id stack.L2ProposerID, cfg *ValidityProposerConfig) {
			cfg.l1ConfigDir = l1Dir
			cfg.l2ConfigDir = l2Dir
		},
	))
}

func WithVPSubmissionInterval(n uint64) ValidityProposerOption {
	return ValidityProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *ValidityProposerConfig) {
		cfg.submissionInterval = &n
	})
}

func WithVPRangeProofInterval(n uint64) ValidityProposerOption {
	return ValidityProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *ValidityProposerConfig) {
		cfg.rangeProofInterval = &n
	})
}

func WithVPRangeProofEvmGasLimit(n uint64) ValidityProposerOption {
	return ValidityProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *ValidityProposerConfig) {
		cfg.rangeProofEvmGasLimit = &n
	})
}

func WithVPMaxConcurrentProofRequests(n uint64) ValidityProposerOption {
	return ValidityProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *ValidityProposerConfig) {
		cfg.maxConcurrentProofRequests = &n
	})
}

func WithVPMaxConcurrentWitnessGen(n uint64) ValidityProposerOption {
	return ValidityProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *ValidityProposerConfig) {
		cfg.maxConcurrentWitnessGen = &n
	})
}

func WithVPLoopInterval(n uint64) ValidityProposerOption {
	return ValidityProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *ValidityProposerConfig) {
		cfg.loopInterval = &n
	})
}

func WithVPProvingTimeout(n uint64) ValidityProposerOption {
	return ValidityProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *ValidityProposerConfig) {
		cfg.provingTimeout = &n
	})
}

func WithVPOpSuccinctConfigName(name string) ValidityProposerOption {
	return ValidityProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *ValidityProposerConfig) {
		cfg.opSuccinctConfigName = &name
	})
}

func WithVPMockMode(enabled bool) ValidityProposerOption {
	return ValidityProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *ValidityProposerConfig) {
		cfg.mockMode = &enabled
	})
}

func WithVPRustLog(level string) ValidityProposerOption {
	return ValidityProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *ValidityProposerConfig) {
		cfg.rustLog = &level
	})
}

// WithVPWriteEnvFile enables writing environment variables to a file.
// When set, the proposer will write all env vars to the specified path at startup.
func WithVPWriteEnvFile(path string) ValidityProposerOption {
	return ValidityProposerOption(func(p devtest.P, id stack.L2ProposerID, cfg *ValidityProposerConfig) {
		cfg.envFilePath = &path
	})
}

func (k *L2SuccinctValidityProposer) Start() {
	k.mu.Lock()
	if k.sub != nil {
		k.logger.Warn("Validity Proposer already started")
		k.mu.Unlock()
		return
	}

	// We pipe sub-process logs to the test-logger.
	// And inspect them along the way, to get the RPC server address.
	logOut := logpipe.ToLogger(k.logger.New("src", "stdout"))
	logErr := logpipe.ToLogger(k.logger.New("src", "stderr"))

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
			k.logger.Info("validity proposer already reaped on shutdown", "err", err)
			return
		}

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				sig := ws.Signal()
				if sig == syscall.SIGINT || sig == syscall.SIGTERM {
					// Keeping postgres running for restart tests
					return
				}
			}
		}

		k.p.Require().NoError(err, "validity proposer exited unexpectedly")
	})

	err := k.sub.Start(k.execPath, k.args, []string{})
	k.p.Require().NoError(err, "Must start")
}

// Stops the validity proposer.
func (k *L2SuccinctValidityProposer) Stop() {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.sub == nil {
		k.logger.Warn("validity proposer already stopped")
		return
	}

	err := k.sub.Stop(true)
	k.p.Require().NoError(err, "Must stop")
	k.sub = nil
}

func WithSuccinctValidityProposer(proposerID stack.L2ProposerID, l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2CLID stack.L2CLNodeID, l2ELID stack.L2ELNodeID, opts ...ValidityProposerOption) stack.Option[*Orchestrator] {
	return stack.AfterDeploy(func(orch *Orchestrator) {
		WithSuccinctValidityProposerPostDeploy(orch, proposerID, l1CLID, l1ELID, l2CLID, l2ELID, opts...)
	})
}

func WithSuperSuccinctValidityProposer(proposerID stack.L2ProposerID,
	l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2CLID stack.L2CLNodeID, l2ELID stack.L2ELNodeID, opts ...ValidityProposerOption) stack.Option[*Orchestrator] {
	return stack.Finally(func(orch *Orchestrator) {
		WithSuccinctValidityProposerPostDeploy(orch, proposerID, l1CLID, l1ELID, l2CLID, l2ELID, opts...)
	})
}

func WithSuccinctValidityProposerPostDeploy(orch *Orchestrator, proposerID stack.L2ProposerID, l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2CLID stack.L2CLNodeID, l2ELID stack.L2ELNodeID, opts ...ValidityProposerOption) {
	ctx := stack.ContextWithID(orch.P().Ctx(), proposerID)
	p := orch.P().WithCtx(ctx)
	logger := p.Logger().New("component", "succinct-validity")

	require := p.Require()
	require.False(orch.proposers.Has(proposerID), "proposer must not already exist")

	l2Net, ok := orch.l2Nets.Get(proposerID.ChainID())
	require.True(ok, "l2 network required")

	l1EL, ok := orch.GetL1EL(l1ELID)
	require.True(ok, "l1 EL node required")

	l1CL, ok := orch.GetL1CL(l1CLID)
	require.True(ok, "l1 CL node required")

	l2EL, ok := orch.GetL2EL(l2ELID)
	require.True(ok, "l2 EL node required")

	l2CL, ok := orch.GetL2CL(l2CLID)
	require.True(ok, "l2 CL node required")

	// --- Embedded Postgres setup ---
	embeddedPG, err := startEmbeddedPostgres(p)
	require.NoError(err, "must start embedded postgres (or read DATABASE_URL)")
	logger.Info("Using embedded Postgres", "url", embeddedPG.URL)
	p.Cleanup(func() {
		if embeddedPG != nil {
			embeddedPG.stop()
			logger.Info("Stopped embedded Postgres and removed temp dir")
		}
	})

	proposerKey, err := orch.GetKeys().Secret(devkeys.ProposerRole.Key(proposerID.ChainID().ToBig()))
	require.NoError(err, "failed to get proposer key")
	proposerKeyStr := hexutil.Encode(crypto.FromECDSA(proposerKey))

	cfg := &ValidityProposerConfig{}
	orch.proposerOptions.Apply(p, proposerID, cfg)
	for _, opt := range opts {
		opt(p, proposerID, cfg)
	}

	require.NotEmpty(cfg.l1ConfigDir, "validity proposer L1 config dir must be set")
	require.NotEmpty(cfg.l2ConfigDir, "validity proposer L2 config dir must be set")

	l1RPC := l1EL.UserRPC()
	l1BeaconRPC := l1CL.beaconHTTPAddr
	l2RPC := strings.ReplaceAll(l2EL.UserRPC(), "ws://", "http://")
	l2NodeRPC := strings.ReplaceAll(l2CL.UserRPC(), "ws://", "http://")
	l2ooAddr := l2Net.deployment.opSuccinctL2OutputOracle

	verifierAddr, err := l2Net.deployment.resolveSP1VerifierAddr()
	require.NoError(err, "failed to get verifier address")

	logger.Info("L1_RPC", "url", l1RPC)
	logger.Info("L1_BEACON_RPC", "url", l1BeaconRPC)
	logger.Info("L2_RPC", "url", l2RPC)
	logger.Info("L2_NODE_RPC", "url", l2NodeRPC)
	logger.Info("VERIFIER_ADDRESS", "address", verifierAddr)
	logger.Info("L2OO_ADDRESS", "address", l2ooAddr)

	envVars := map[string]string{
		"L1_RPC":           l1RPC,
		"L1_BEACON_RPC":    l1BeaconRPC,
		"L2_RPC":           l2RPC,
		"L2_NODE_RPC":      l2NodeRPC,
		"VERIFIER_ADDRESS": verifierAddr.String(),
		"L2OO_ADDRESS":     l2ooAddr.String(),
		"DATABASE_URL":     embeddedPG.URL,
		"PRIVATE_KEY":      proposerKeyStr,
		"L1_CONFIG_DIR":    cfg.l1ConfigDir,
		"L2_CONFIG_DIR":    cfg.l2ConfigDir,
		"LOG_FORMAT":       "json",
	}

	setEnvFromEnvOrDefault(envVars, "NETWORK_PRIVATE_KEY", "")

	// Mock mode: default true, false if NETWORK_PRIVATE_KEY is set (real proving)
	if envVars["NETWORK_PRIVATE_KEY"] != "" {
		envVars["OP_SUCCINCT_MOCK"] = "false"
	} else {
		envVars["OP_SUCCINCT_MOCK"] = "true"
	}

	// Optional parameters (override defaults if set)
	setEnvIfNotNil(envVars, "SUBMISSION_INTERVAL", cfg.submissionInterval)
	setEnvIfNotNil(envVars, "RANGE_PROOF_INTERVAL", cfg.rangeProofInterval)
	setEnvIfNotNil(envVars, "RANGE_PROOF_EVM_GAS_LIMIT", cfg.rangeProofEvmGasLimit)
	setEnvIfNotNil(envVars, "MAX_CONCURRENT_PROOF_REQUESTS", cfg.maxConcurrentProofRequests)
	setEnvIfNotNil(envVars, "MAX_CONCURRENT_WITNESS_GEN", cfg.maxConcurrentWitnessGen)
	setEnvIfNotNil(envVars, "LOOP_INTERVAL", cfg.loopInterval)
	setEnvIfNotNil(envVars, "PROVING_TIMEOUT", cfg.provingTimeout)
	setEnvIfNotNil(envVars, "OP_SUCCINCT_CONFIG_NAME", cfg.opSuccinctConfigName)
	setEnvIfNotNil(envVars, "OP_SUCCINCT_MOCK", cfg.mockMode)
	setEnvIfNotNil(envVars, "RUST_LOG", cfg.rustLog)

	if areMetricsEnabled() {
		metricsPort, err := getAvailableLocalPort()
		require.NoError(err, "failed to get available port for metrics")
		envVars["METRICS_PORT"] = metricsPort
		metricsTarget := NewPrometheusMetricsTarget("localhost", metricsPort, false)
		orch.RegisterL2MetricsTargets(proposerID, metricsTarget)
		logger.Info("Registered validity proposer metrics", "port", metricsPort)
	}

	envDir := p.TempDir()
	envFile := filepath.Join(envDir, fmt.Sprintf("validity-proposer-%s.env", proposerID.String()))
	err = WriteEnvFile(envFile, envVars)
	p.Require().NoError(err, "must write validity proposer env file")

	if cfg.envFilePath != nil {
		err = WriteEnvFile(*cfg.envFilePath, envVars)
		p.Require().NoError(err, "must write env file")
		logger.Info("env file written", "path", *cfg.envFilePath)
	}

	execPath := os.Getenv("VALIDITY_PROPOSER_EXEC_PATH")
	p.Require().NotEmpty(execPath, "VALIDITY_PROPOSER_EXEC_PATH environment variable must be set")
	_, err = os.Stat(execPath)
	p.Require().NotErrorIs(err, os.ErrNotExist, "executable must exist")

	k := &L2SuccinctValidityProposer{
		id:                 proposerID,
		userRPC:            "", // retrieved from logs
		execPath:           execPath,
		args:               []string{"--env-file", envFile},
		p:                  p,
		logger:             logger,
		databaseURL:        embeddedPG.URL,
		l2MetricsRegistrar: orch,
	}
	logger.Info("Starting validity proposer")
	k.Start()
	p.Cleanup(func() {
		logger.Info("Stopping validity proposer")
		k.Stop()
	})
	logger.Info("validity proposer is running", "rpc", k.UserRPC())
	require.True(orch.proposers.SetIfMissing(proposerID, k), "must not already exist")
}

// EmbeddedPG holds the running instance and the URL we pass to children.
type EmbeddedPG struct {
	pg  *embeddedpg.EmbeddedPostgres
	URL string
}

// startEmbeddedPostgres starts a local postgres ONLY if we don't already have `DATABASE_URL`.
func startEmbeddedPostgres(p devtest.P) (*EmbeddedPG, error) {

	const (
		pgUser        = "op-succinct"
		pgDB          = "op-succinct"
		pgPass        = "postgres"
		pgRuntimePath = "runtime"
		pgDataPath    = "data"
	)

	// 1) Caller already provided a DB → just wrap it.
	if v := os.Getenv("DATABASE_URL"); v != "" {
		epg := &EmbeddedPG{pg: nil, URL: v}
		return epg, nil
	}

	// 2) We need to start our own.
	base := p.TempDir()
	pgRoot, err := os.MkdirTemp(base, "embedded-pg-*")
	runtimePath := filepath.Join(pgRoot, pgRuntimePath)
	dataPath := filepath.Join(pgRoot, pgDataPath)

	portStr, err := getAvailableLocalPort()
	if err != nil {
		return nil, fmt.Errorf("getAvailableLocalPort: %w", err)
	}

	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("get available port: %w", err)
	}

	port := uint32(portInt)

	cfg := embeddedpg.DefaultConfig().
		Port(port).
		Database(pgDB).
		Username(pgUser).
		Password(pgPass).
		RuntimePath(runtimePath).
		DataPath(dataPath)

	pg := embeddedpg.NewDatabase(cfg)
	if err := pg.Start(); err != nil {
		return nil, fmt.Errorf("start embedded postgres: %w", err)
	}

	url := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable", pgUser, pgPass, port, pgDB)
	epg := &EmbeddedPG{pg: pg, URL: url}
	return epg, nil
}

// stop stops PG if we actually started it.
func (e *EmbeddedPG) stop() {
	if e == nil || e.pg == nil {
		return
	}
	_ = e.pg.Stop()
}
