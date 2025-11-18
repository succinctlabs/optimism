package sysgo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	embeddedpg "github.com/fergusstrange/embedded-postgres"
)

type L2SVProposer struct {
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

var _ L2Prop = (*L2SVProposer)(nil)

func (p *L2SVProposer) hydrate(system stack.ExtensibleSystem) {
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

func (k *L2SVProposer) Start() {
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

	userRPCChan := make(chan string, 1)
	defer close(userRPCChan)
	metricsTargetChan := make(chan PrometheusMetricsTarget, 1)
	defer close(metricsTargetChan)

	onLogEntry := func(e logpipe.LogEntry) {
		msg := e.LogMessage()
		if msg == "RPC server bound to address" {
			userRPCChan <- "http://" + e.FieldValue("addr").(string)
		} else if metricsUrl, found := strings.CutPrefix(msg, "Serving metrics at: "); found {
			// Matching messages like "Serving metrics at: http://0.0.0.0:9091"
			if !strings.HasPrefix(metricsUrl, "http") {
				metricsUrl = fmt.Sprintf("http://%s", metricsUrl)
			}
			parsedUrl, err := url.Parse(metricsUrl)
			k.p.Require().NoError(err, "invalid metrics url output to logs", "log", msg)
			k.p.Require().NotEmpty(parsedUrl.Port(), "empty port in logged metrics url", "log", msg)
			metricsTargetChan <- NewPrometheusMetricsTarget(parsedUrl.Hostname(), parsedUrl.Port(), false)
		}
	}
	stdOutLogs := logpipe.LogProcessor(func(line []byte) {
		e := logpipe.ParseRustStructuredLogs(line)
		logOut(e)
		onLogEntry(e)
	})
	stdErrLogs := logpipe.LogProcessor(func(line []byte) {
		e := logpipe.ParseRustStructuredLogs(line)
		logErr(e)
	})
	k.sub = NewSubProcess(k.p, stdOutLogs, stdErrLogs)
	k.mu.Unlock()

	k.sub.OnExit(func(err error) {
		k.embeddedPG.stop()

		if errors.Is(err, syscall.ECHILD) {
			k.p.Logger().Info("validity proposer already reaped on shutdown", "err", err)
			return
		}

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				sig := ws.Signal()
				if sig == syscall.SIGINT || sig == syscall.SIGTERM {
					k.p.Logger().Info("validity proposer interrupted during shutdown", "signal", sig)
					return
				}
			}
		}

		k.p.Require().NoError(err, "validity proposer exited unexpectedly")
	})

	err := k.sub.Start(k.execPath, k.args, []string{})
	k.p.Require().NoError(err, "Must start")

	if areMetricsEnabled() {
		var metricsTarget PrometheusMetricsTarget
		k.p.Require().NoError(tasks.Await(k.p.Ctx(), metricsTargetChan, &metricsTarget), "need metrics endpoint")
		k.l2MetricsRegistrar.RegisterL2MetricsTargets(k.id, metricsTarget)
	}
}

// Stops the validity proposer.
// warning: no restarts supported yet, since the RPC port is not remembered.
func (k *L2SVProposer) Stop() {
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

func (k *L2SVProposer) UserRPC() string {
	return k.userRPC
}

func WithSVProposer(proposerID stack.L2ProposerID, l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2CLID stack.L2CLNodeID, l2ELID stack.L2ELNodeID) stack.Option[*Orchestrator] {
	return stack.AfterDeploy(func(orch *Orchestrator) {
		WithL2SVProposerPostDeploy(orch, proposerID, l1CLID, l1ELID, l2CLID, l2ELID)
	})
}

func WithSuperSVProposer(proposerID stack.L2ProposerID,
	l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2CLID stack.L2CLNodeID, l2ELID stack.L2ELNodeID) stack.Option[*Orchestrator] {
	return stack.Finally(func(orch *Orchestrator) {
		WithL2SVProposerPostDeploy(orch, proposerID, l1CLID, l1ELID, l2CLID, l2ELID)
	})
}

func WithL2SVProposerPostDeploy(orch *Orchestrator, proposerID stack.L2ProposerID, l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2CLID stack.L2CLNodeID, l2ELID stack.L2ELNodeID, opts ...L2CLOption) {
	ctx := orch.P().Ctx()
	ctx = stack.ContextWithID(ctx, proposerID)
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
	p.Require().NoError(err, "must start embedded postgres (or read DATABASE_URL)")
	logger.Info("Using embedded Postgres", "url", embeddedPG.URL)
	p.Cleanup(func() {
		if embeddedPG != nil {
			embeddedPG.stop()
			logger.Info("Stopped embedded Postgres and removed temp dir")
		}
	})

	dgf := l2Net.deployment.disputeGameFactoryProxy
	logger.Info("Using DisputeGameFactory", "address", dgf)

	mockVerifierAddr := l2Net.deployment.sp1MockVerifier
	logger.Info("Using mock verifier", "address", mockVerifierAddr)

	l2ooAddr := l2Net.deployment.opSuccinctL2OutputOracle
	logger.Info("Using L2OO", "address", l2ooAddr)

	proposerKey, err := orch.keys.Secret(devkeys.ProposerRole.Key(proposerID.ChainID().ToBig()))
	require.NoError(err)
	proposerKeyStr := hexutil.Encode(crypto.FromECDSA(proposerKey))

	envVars := []string{
		"L1_RPC=" + l1EL.UserRPC(),
		"L1_NODE_RPC=" + l1CL.beaconHTTPAddr,
		"L2_RPC=" + strings.ReplaceAll(l2EL.UserRPC(), "ws://", "http://"),
		"L2_NODE_RPC=" + strings.ReplaceAll(l2CL.UserRPC(), "ws://", "http://"),
		"VERIFIER_ADDRESS=" + mockVerifierAddr.String(),
		"L2OO_ADDRESS=" + l2ooAddr.String(),
		"DATABASE_URL=" + embeddedPG.URL,
		"PRIVATE_KEY=" + proposerKeyStr,
		propagateEnvVarOrDefault("NETWORK_PRIVATE_KEY", ""),
		"LOG_FORMAT=json",
	}

	envDir := p.TempDir()
	envFile := filepath.Join(envDir, fmt.Sprintf("l2-sv-proposer-%s.env", proposerID.String()))
	err = os.WriteFile(envFile, []byte(strings.Join(envVars, "\n")), 0o600)
	p.Require().NoError(err, "must write sv proposer env file")

	if areMetricsEnabled() {
		metricsPort, err := getAvailableLocalPort()
		p.Require().NoError(err, "WithL2SVProposer: getting metrics port")

		envVars = append(envVars, propagateEnvVarOrDefault("SV_PROPOSER_METRICS_PORT", metricsPort))
		envVars = append(envVars, "SV_PROPOSER_METRICS_ENABLED=true")
	}

	execPath := os.Getenv("SV_PROPOSER_EXEC_PATH")
	p.Require().NotEmpty(execPath, "SV_PROPOSER_EXEC_PATH environment variable must be set")
	_, err = os.Stat(execPath)
	p.Require().NotErrorIs(err, os.ErrNotExist, "executable must exist")

	k := &L2SVProposer{
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

// Deploys an OPSuccinctL2OutputOracle contract for each specified chain, and
// updates the orchestrator's L2 network deployments accordingly.
func WithDeployOpSuccinctL2OutputOracle(
	l1CL stack.L1CLNodeID,
	l1EL stack.L1ELNodeID,
	l2CL stack.L2CLNodeID,
	l2EL stack.L2ELNodeID,
) stack.Option[*Orchestrator] {
	return stack.AfterDeploy(func(o *Orchestrator) {
		rootPrefix, err := findMonorepoRoot("Cargo.lock")
		o.P().Require().NoError(err, "failed to locate monorepo root")

		repoRoot, err := filepath.Abs(rootPrefix)
		o.P().Require().NoError(err, "failed to resolve monorepo root")

		addr, err := o.deployOpSuccinctL2OutputOracle(repoRoot, l1CL, l1EL, l2CL, l2EL)
		o.P().Require().NoError(err, "failed to deploy OPSuccinctL2OutputOracle")

		l2Net, ok := o.l2Nets.Get(l2CL.ChainID())
		o.P().Require().True(ok, "l2 network required")
		l2Net.deployment.opSuccinctL2OutputOracle = common.HexToAddress(addr)
	})
}

// deployOpSuccinctL2OutputOracle deploys an OPSuccinctL2OutputOracle contract
func (o *Orchestrator) deployOpSuccinctL2OutputOracle(
	repoRoot string,
	l1CLID stack.L1CLNodeID,
	l1ELID stack.L1ELNodeID,
	l2CLID stack.L2CLNodeID,
	l2ELID stack.L2ELNodeID,
) (string, error) {

	p := o.P()
	logger := p.Logger().New("chain", l2CLID.ChainID().String())
	require := p.Require()

	l1Net, ok := o.l1Nets.Get(l1CLID.ChainID())
	require.True(ok, "l1 network required")

	l1CL, ok := o.l1CLs.Get(l1CLID)
	require.True(ok, "l1 CL node required")

	l1EL, ok := o.l1ELs.Get(l1ELID)
	require.True(ok, "l1 EL node required")

	l2Net, ok := o.l2Nets.Get(l2CLID.ChainID())
	require.True(ok, "l2 network required")

	l2CL, ok := o.l2CLs.Get(l2CLID)
	require.True(ok, "l2 CL node required")

	l2EL, ok := o.l2ELs.Get(l2ELID)
	require.True(ok, "l2 EL node required")

	l1ChainID := l1CLID.ChainID().ToBig()
	l1PAOKey, err := o.keys.Secret(devkeys.L1ProxyAdminOwnerRole.Key(l1ChainID))
	if err != nil {
		return "", fmt.Errorf("failed to get L1ProxyAdminOwnerRole key: %w", err)
	}
	l1PAOKeyStr := hexutil.Encode(crypto.FromECDSA(l1PAOKey))

	envVars := map[string]string{
		"L1_RPC":           l1EL.UserRPC(),
		"L1_NODE_RPC":      l1CL.beaconHTTPAddr,
		"L2_RPC":           strings.ReplaceAll(l2EL.UserRPC(), "ws://", "http://"),
		"L2_NODE_RPC":      strings.ReplaceAll(l2CL.UserRPC(), "ws://", "http://"),
		"VERIFIER_ADDRESS": l2Net.deployment.sp1MockVerifier.Hex(),
		"PRIVATE_KEY":      l1PAOKeyStr,
		"RUST_LOG":         "info",
	}

	envDir := p.TempDir()
	envFile := filepath.Join(envDir, fmt.Sprintf("opsuccinctl2oo-%s.env", strings.ReplaceAll(l2CLID.ChainID().String(), "-", "_")))
	if err = writeEnvFile(envFile, envVars); err != nil {
		return "", fmt.Errorf("failed to write opsuccinct env: %w", err)
	}

	l1ChainConfig := l1Net.genesis.Config

	err = writeL1ChainConfig(l1ChainConfig, l1CLID.ChainID(), logger)
	if err != nil {
		return "", fmt.Errorf("failed to write L1 chain config: %w", err)
	}

	logger.Info("Deploying OPSuccinctL2OutputOracle")
	addr, err := execDeployOracle(o.P().Ctx(), repoRoot, envFile, logger)
	if err != nil {
		return "", err
	}

	logger.Info("Deployed OPSuccinctL2OutputOracle", "address", addr)
	return addr, nil
}

func writeL1ChainConfig(
	l1ChainConfig any,
	l1ChainID fmt.Stringer,
	logger log.Logger,
) error {

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get cwd: %w", err)
	}

	dir := filepath.Join(cwd, "Configs", "L1")
	if err = os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir %q: %w", dir, err)
	}

	path := filepath.Join(dir, l1ChainID.String()+".json")
	logger.Info("writing L1 chain config for opsuccinct L2OO", "path", path)

	data, err := json.Marshal(l1ChainConfig)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// execDeployOracle runs `just deploy-oracle <envFile>` and parses the output
func execDeployOracle(ctx context.Context, repoRoot, envFile string, logger log.Logger) (string, error) {
	cmd := exec.CommandContext(ctx, "just", "deploy-oracle", envFile)
	cmd.Dir = repoRoot

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	logger.Info("Executing deploy-oracle", "cmd", strings.Join(cmd.Args, " "))

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("deploy-oracle failed: %w\nstdout:\n%s\nstderr:\n%s",
			err, stdout.String(), stderr.String())
	}

	stdoutStr := strings.TrimSpace(stdout.String())
	if stdoutStr != "" {
		logger.Info("deploy-oracle output", "stdout", stdoutStr)
	}

	// Try to parse the `== Return ==` section:
	// 0: address 0x123...
	reReturn := regexp.MustCompile(`(?m)^0:\s+address\s+(0x[0-9a-fA-F]{40})\b`)
	if m := reReturn.FindStringSubmatch(stdoutStr); len(m) == 2 {
		addr := m[1]
		return addr, nil
	}

	return "", fmt.Errorf("deploy-oracle succeeded but could not find the address.\nstdout:\n%s", stdoutStr)
}

func writeEnvFile(path string, kv map[string]string) error {
	var keys []string
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for _, k := range keys {
		val := kv[k]
		if val == "" {
			continue
		}
		fmt.Fprintf(&b, "%s=%s\n", k, val)
	}
	return os.WriteFile(path, []byte(b.String()), 0o600)
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
