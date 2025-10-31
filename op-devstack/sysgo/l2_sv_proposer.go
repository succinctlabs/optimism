package sysgo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ethereum-optimism/optimism/op-devstack/devtest"
	"github.com/ethereum-optimism/optimism/op-devstack/shim"
	"github.com/ethereum-optimism/optimism/op-devstack/stack"
	ps "github.com/ethereum-optimism/optimism/op-proposer/proposer"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/logpipe"
	"github.com/ethereum-optimism/optimism/op-service/tasks"
	"github.com/ethereum-optimism/optimism/op-service/testutils/tcpproxy"
)

type L2SVProposer struct {
	mu                 sync.Mutex
	id                 stack.L2ProposerID
	service            *ps.ProposerService
	userRPC            string
	userProxy          *tcpproxy.Proxy
	execPath           string
	args               []string
	env                []string
	p                  devtest.P
	sub                *SubProcess
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
	defer k.mu.Unlock()
	if k.sub != nil {
		k.p.Logger().Warn("Validity Proposer already started")
		return
	}

	// Create a proxy for the user RPC,
	// so other services can connect, and stay connected, across restarts.
	if k.userProxy == nil {
		k.userProxy = tcpproxy.New(k.p.Logger())
		k.p.Require().NoError(k.userProxy.Start())
		k.p.Cleanup(func() {
			k.userProxy.Close()
		})
		k.userRPC = "http://" + k.userProxy.Addr()
	}

	// Create the sub-process.
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

	err := k.sub.Start(k.execPath, k.args, k.env)
	k.p.Require().NoError(err, "Must start")

	var userRPCAddr string
	k.p.Require().NoError(tasks.Await(k.p.Ctx(), userRPCChan, &userRPCAddr), "need user RPC")

	if areMetricsEnabled() {
		var metricsTarget PrometheusMetricsTarget
		k.p.Require().NoError(tasks.Await(k.p.Ctx(), metricsTargetChan, &metricsTarget), "need metrics endpoint")
		k.l2MetricsRegistrar.RegisterL2MetricsTargets(k.id, metricsTarget)
	}

	k.userProxy.SetUpstream(ProxyAddr(k.p.Require(), userRPCAddr))
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

func WithSVProposer(l2CLID stack.L2CLNodeID, l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2ELID stack.L2ELNodeID, proposerID stack.L2ProposerID) stack.Option[*Orchestrator] {
	return stack.AfterDeploy(func(orch *Orchestrator) {
		WithL2SVProposerPostDeploy(orch, l2CLID, l1CLID, l1ELID, l2ELID, proposerID)
	})
}

func WithSuperSVProposer(l2CLID stack.L2CLNodeID, l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2ELID stack.L2ELNodeID, proposerID stack.L2ProposerID) stack.Option[*Orchestrator] {
	return stack.Finally(func(orch *Orchestrator) {
		WithL2SVProposerPostDeploy(orch, l2CLID, l1CLID, l1ELID, l2ELID, proposerID)
	})
}

func WithL2SVProposerPostDeploy(orch *Orchestrator, l2CLID stack.L2CLNodeID, l1CLID stack.L1CLNodeID, l1ELID stack.L1ELNodeID, l2ELID stack.L2ELNodeID, proposerID stack.L2ProposerID, opts ...L2CLOption) {
	p := orch.P().WithCtx(stack.ContextWithID(orch.P().Ctx(), l2CLID))

	require := p.Require()

	l1Net, ok := orch.l1Nets.Get(l1CLID.ChainID())
	require.True(ok, "l1 network required")

	l2Net, ok := orch.l2Nets.Get(l2CLID.ChainID())
	require.True(ok, "l2 network required")

	l1ChainConfig := l1Net.genesis.Config

	l1EL, ok := orch.l1ELs.Get(l1ELID)
	require.True(ok, "l1 EL node required")

	l1CL, ok := orch.l1CLs.Get(l1CLID)
	require.True(ok, "l1 CL node required")

	l2EL, ok := orch.l2ELs.Get(l2ELID)
	require.True(ok, "l2 EL node required")

	l2CL, ok := orch.l2CLs.Get(l2CLID)

	cfg := DefaultL2CLConfig()
	orch.l2CLOptions.Apply(orch.P(), l2CLID, cfg)       // apply global options
	L2CLOptionBundle(opts).Apply(orch.P(), l2CLID, cfg) // apply specific options

	tempKonaDir := p.TempDir()

	tempRollupCfgPath := filepath.Join(tempKonaDir, "rollup.json")
	rollupCfgData, err := json.Marshal(l2Net.rollupCfg)
	p.Require().NoError(err, "must write rollup config")
	p.Require().NoError(err, os.WriteFile(tempRollupCfgPath, rollupCfgData, 0o644))

	tempL1CfgPath := filepath.Join(tempKonaDir, "l1-chain-config.json")
	l1CfgData, err := json.Marshal(l1ChainConfig)
	p.Require().NoError(err, "must write l1 chain config")
	p.Require().NoError(err, os.WriteFile(tempL1CfgPath, l1CfgData, 0o644))

	envVars := []string{
		"L1_RPC=" + l1EL.UserRPC(),
		"L1_NODE_RPC=" + l1CL.beaconHTTPAddr,
		"L2_RPC=" + strings.ReplaceAll(l2EL.EngineRPC(), "ws://", "http://"),
		"L2_NODE_RPC=" + l2CL.UserRPC(),
	}

	if areMetricsEnabled() {
		// NB: Instead of getAvailableLocalPort, we should pass "0" so the OS picks its
		// own port, but that is not currently logged properly so we cannot parse it.
		// See: https://github.com/op-rs/kona/issues/2987
		metricsPort, err := getAvailableLocalPort()
		p.Require().NoError(err, "WithL2SVProposer: getting metrics port")

		envVars = append(envVars, propagateEnvVarOrDefault("SV_PROPOSER_METRICS_PORT", metricsPort))
		envVars = append(envVars, "SV_PROPOSRE_METRICS_ENABLED=true")
	}

	execPath := os.Getenv("SV_PROPOSER_EXEC_PATH")
	p.Require().NotEmpty(execPath, "SV_PROPOSER_EXEC_PATH environment variable must be set")
	_, err = os.Stat(execPath)
	p.Require().NotErrorIs(err, os.ErrNotExist, "executable must exist")

	k := &L2SVProposer{
		id:                 proposerID,
		userRPC:            "", // retrieved from logs
		execPath:           execPath,
		args:               []string{},
		env:                envVars,
		p:                  p,
		l2MetricsRegistrar: orch,
	}
	p.Logger().Info("Starting validity proposer")
	k.Start()
	p.Cleanup(k.Stop)
	p.Logger().Info("validity proposer is up", "rpc", k.UserRPC())
	require.True(orch.proposers.SetIfMissing(proposerID, k), "must not already exist")
}
