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
	"github.com/ethereum-optimism/optimism/op-service/tasks"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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
	sub                *SubProcess
	embeddedPG         *EmbeddedPG
	l2MetricsRegistrar L2MetricsRegistrar
}

var _ L2Prop = (*L2SuccinctValidityProposer)(nil)

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

type ValidityProposerConfig struct {
	l1ConfigDir string
	l2ConfigDir string
}

type ValidityProposerOption func(id stack.L2ProposerID, cfg *ValidityProposerConfig)

func WithValidityProposerOption(o *Orchestrator, opt ValidityProposerOption) {
	o.proposerOptions = append(o.proposerOptions, func(id stack.L2ProposerID, cfg any) {
		c, ok := cfg.(*ValidityProposerConfig)
		if !ok {
			return
		}
		opt(id, c)
	})
}

func WithValidityConfigDirsOption(
	o *Orchestrator,
	l1Dir, l2Dir string,
) {
	WithValidityProposerOption(o, func(id stack.L2ProposerID, cfg *ValidityProposerConfig) {
		cfg.l1ConfigDir = l1Dir
		cfg.l2ConfigDir = l2Dir
	})
}

func (k *L2SuccinctValidityProposer) Start() {
	k.mu.Lock()
	if k.sub != nil {
		k.p.Logger().Warn("Validity Proposer already started")
		k.mu.Unlock()
		return
	}

	// We pipe sub-process logs to the test-logger.
	// And inspect them along the way, to get the RPC server address.
	logOut := logpipe.ToLogger(k.p.Logger().New("component", "validity", "src", "stdout"))
	logErr := logpipe.ToLogger(k.p.Logger().New("component", "validity", "src", "stderr"))

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
		if k.embeddedPG != nil {
			k.embeddedPG.stop()
		}

		if errors.Is(err, syscall.ECHILD) {
			k.p.Logger().Info("validity proposer already reaped on shutdown", "err", err)
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

		k.p.Require().NoError(err, "validity proposer exited unexpectedly")
	})

	err := k.sub.Start(k.execPath, k.args, []string{})
	k.p.Require().NoError(err, "Must start")

	metricsTargetChan := make(chan PrometheusMetricsTarget, 1)
	if areMetricsEnabled() {
		var metricsTarget PrometheusMetricsTarget
		k.p.Require().NoError(tasks.Await(k.p.Ctx(), metricsTargetChan, &metricsTarget), "need metrics endpoint")
		k.l2MetricsRegistrar.RegisterL2MetricsTargets(k.id, metricsTarget)
	}
}

// Stops the validity proposer.
func (k *L2SuccinctValidityProposer) Stop() {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.sub == nil {
		k.p.Logger().Warn("validity proposer already stopped")
		return
	}

	err := k.sub.Stop(true)
	k.p.Require().NoError(err, "Must stop")
	k.sub = nil
}

func WithSuccinctValidityProposer(proposerID stack.L2ProposerID, l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2CLID stack.L2CLNodeID, l2ELID stack.L2ELNodeID) stack.Option[*Orchestrator] {
	return stack.AfterDeploy(func(orch *Orchestrator) {
		WithSuccinctValidityProposerPostDeploy(orch, proposerID, l1CLID, l1ELID, l2CLID, l2ELID)
	})
}

func WithSuperSuccinctValidityProposer(proposerID stack.L2ProposerID,
	l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2CLID stack.L2CLNodeID, l2ELID stack.L2ELNodeID) stack.Option[*Orchestrator] {
	return stack.Finally(func(orch *Orchestrator) {
		WithSuccinctValidityProposerPostDeploy(orch, proposerID, l1CLID, l1ELID, l2CLID, l2ELID)
	})
}

func WithSuccinctValidityProposerPostDeploy(orch *Orchestrator, proposerID stack.L2ProposerID, l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2CLID stack.L2CLNodeID, l2ELID stack.L2ELNodeID, opts ...L2CLOption) {
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

	cfg := DefaultL2CLConfig()
	orch.l2CLOptions.Apply(orch.P(), l2CLID, cfg)       // apply global options
	L2CLOptionBundle(opts).Apply(orch.P(), l2CLID, cfg) // apply specific options

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

	mockVerifierAddr := l2Net.deployment.sp1MockVerifier
	logger.Info("Using SP1MockVerifier", "address", mockVerifierAddr)

	l2ooAddr := l2Net.deployment.opSuccinctL2OutputOracle
	logger.Info("Using OPSuccinctL2OutputOracle", "address", l2ooAddr)

	proposerKey, err := orch.keys.Secret(devkeys.ProposerRole.Key(proposerID.ChainID().ToBig()))
	require.NoError(err)
	proposerKeyStr := hexutil.Encode(crypto.FromECDSA(proposerKey))

	vpCfg := &ValidityProposerConfig{}
	for _, opt := range orch.proposerOptions {
		opt(proposerID, vpCfg)
	}

	require.NotEmpty(vpCfg.l1ConfigDir, "validity proposer L1 config dir must be set")
	require.NotEmpty(vpCfg.l2ConfigDir, "validity proposer L2 config dir must be set")

	envVars := map[string]string{
		"L1_RPC":               l1EL.UserRPC(),
		"L1_BEACON_RPC":        l1CL.beaconHTTPAddr,
		"L2_RPC":               strings.ReplaceAll(l2EL.UserRPC(), "ws://", "http://"),
		"L2_NODE_RPC":          strings.ReplaceAll(l2CL.UserRPC(), "ws://", "http://"),
		"VERIFIER_ADDRESS":     mockVerifierAddr.String(),
		"L2OO_ADDRESS":         l2ooAddr.String(),
		"DATABASE_URL":         embeddedPG.URL,
		"PRIVATE_KEY":          proposerKeyStr,
		"SUBMISSION_INTERVAL":  "10",
		"RANGE_PROOF_INTERVAL": "10",
		"OP_SUCCINCT_MOCK":     "true",
		"L1_CONFIG_DIR":        vpCfg.l1ConfigDir,
		"L2_CONFIG_DIR":        vpCfg.l2ConfigDir,
		"LOG_FORMAT":           "json",
	}

	setEnvFromEnvOrDefault(envVars, "NETWORK_PRIVATE_KEY", "")

	if areMetricsEnabled() {
		metricsPort, err := getAvailableLocalPort()
		p.Require().NoError(err, "must get available port for metrics")
		setEnvFromEnvOrDefault(envVars, "VALIDITY_PROPOSER_METRICS_PORT", metricsPort)
		envVars["VALIDITY_PROPOSER_METRICS_ENABLED"] = "true"
	}

	envDir := p.TempDir()
	envFile := filepath.Join(envDir, fmt.Sprintf("validity-proposer-%s.env", proposerID.String()))
	err = writeEnvFile(envFile, envVars)
	p.Require().NoError(err, "must write validity proposer env file")

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
		embeddedPG:         embeddedPG,
		l2MetricsRegistrar: orch,
	}
	p.Logger().Info("Starting validity proposer")
	k.Start()
	p.Cleanup(func() {
		logger.Info("Stopping validity proposer")
		k.Stop()
	})
	p.Logger().Info("validity proposer is running", "rpc", k.UserRPC())
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
		pgPass        = "posgres"
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

	url := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s", pgUser, pgPass, port, pgDB)
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

func setEnvFromEnvOrDefault(env map[string]string, key, def string) {
	if v := os.Getenv(key); v != "" {
		env[key] = v
	} else if def != "" {
		env[key] = def
	}
}
