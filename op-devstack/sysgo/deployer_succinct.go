package sysgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ethereum-optimism/optimism/op-chain-ops/devkeys"
	"github.com/ethereum-optimism/optimism/op-devstack/stack"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

// Deploys an SP1MockVerifier contract for the specified L2 chain, and updates
// the orchestrator's L2 network deployments accordingly.
func WithDeploySP1MockVerifier(
	l1EL stack.L1ELNodeID,
	l2ChainID eth.ChainID,
) stack.Option[*Orchestrator] {
	return stack.AfterDeploy(func(o *Orchestrator) {
		rootPrefix, err := findMonorepoRoot("Cargo.lock")
		o.P().Require().NoError(err, "failed to locate monorepo root")

		repoRoot, err := filepath.Abs(rootPrefix)
		o.P().Require().NoError(err, "failed to resolve monorepo root")

		addr, err := o.deploySP1MockVerifier(repoRoot, l1EL, l2ChainID)
		o.P().Require().NoError(err, "failed to deploy SP1MockVerifier")

		l2Net, ok := o.l2Nets.Get(l2ChainID)
		o.P().Require().True(ok, "l2 network required")
		l2Net.deployment.sp1MockVerifier = common.HexToAddress(addr)
	})
}

// deploySP1MockVerifier deploys an SP1MockVerifier contract
func (o *Orchestrator) deploySP1MockVerifier(
	repoRoot string,
	l1ELID stack.L1ELNodeID,
	l2ChainID eth.ChainID,
) (string, error) {

	p := o.P()
	logger := p.Logger().New("chain", l2ChainID.String())
	require := p.Require()

	l1ChainID := l1ELID.ChainID()

	l1EL, ok := o.l1ELs.Get(l1ELID)
	require.True(ok, "l1 EL node required")

	l1PAOKey, err := o.keys.Secret(devkeys.L1ProxyAdminOwnerRole.Key(l1ChainID.ToBig()))
	if err != nil {
		return "", fmt.Errorf("failed to get L1ProxyAdminOwnerRole key: %w", err)
	}
	l1PAOKeyStr := hexutil.Encode(crypto.FromECDSA(l1PAOKey))

	envVars := map[string]string{
		"L1_RPC":      l1EL.UserRPC(),
		"PRIVATE_KEY": l1PAOKeyStr,
	}

	envDir := p.TempDir()
	envFile := filepath.Join(envDir, fmt.Sprintf("sp1-mock-verifier-%s.env", strings.ReplaceAll(l2ChainID.String(), "-", "_")))
	if err = writeEnvFile(envFile, envVars); err != nil {
		return "", fmt.Errorf("failed to write sp1-mock-verifier env: %w", err)
	}

	addr, err := execDeployMockVerifier(o.P().Ctx(), repoRoot, envFile, logger)
	if err != nil {
		return "", err
	}

	logger.Info("Deployed SP1MockVerifier", "address", addr)
	return addr, nil
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
	l2ChainID := l2CLID.ChainID()
	logger := p.Logger().New("chain", l2ChainID.String())
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

	base := p.TempDir()

	l1CfgDir := l1ConfigDir(base)
	err = os.MkdirAll(l1CfgDir, 0o755)
	require.NoError(err, "mkdir l1 config dir")

	l2CfgDir := l2ConfigDir(base)
	os.MkdirAll(l2CfgDir, 0o755)
	require.NoError(err, "mkdir l2 config dir")

	WithValidityConfigDirsOption(o, l1CfgDir, l2CfgDir)

	envVars := map[string]string{
		"L1_RPC":           l1EL.UserRPC(),
		"L1_NODE_RPC":      l1CL.beaconHTTPAddr,
		"L2_RPC":           strings.ReplaceAll(l2EL.UserRPC(), "ws://", "http://"),
		"L2_NODE_RPC":      strings.ReplaceAll(l2CL.UserRPC(), "ws://", "http://"),
		"VERIFIER_ADDRESS": l2Net.deployment.sp1MockVerifier.Hex(),
		"PRIVATE_KEY":      l1PAOKeyStr,
		"L1_CONFIG_DIR":    l1CfgDir,
		"L2_CONFIG_DIR":    l2CfgDir,
		"RUST_LOG":         "info",
	}

	envDir := p.TempDir()
	envFile := filepath.Join(envDir, fmt.Sprintf("op-succinct-l2oo-%s.env", strings.ReplaceAll(l2ChainID.String(), "-", "_")))
	if err = writeEnvFile(envFile, envVars); err != nil {
		return "", fmt.Errorf("failed to write op-succinct-l2oo env: %w", err)
	}

	l1ChainConfig := l1Net.genesis.Config

	err = writeL1ChainConfig(l1ChainConfig, l1CLID.ChainID(), l1CfgDir, logger)
	if err != nil {
		return "", fmt.Errorf("failed to write L1 chain config: %w", err)
	}

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
	dir string,
	logger log.Logger,
) error {
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

func l1ConfigDir(base string) string {
	return filepath.Join(base, "Configs", "L1")
}

func l2ConfigDir(base string) string {
	return filepath.Join(base, "Configs", "L2")
}

// execDeployMockVerifier runs `just deploy-mock-verifier <envFile>` and parses the output
func execDeployMockVerifier(ctx context.Context, repoRoot, envFile string, logger log.Logger) (string, error) {
	cmd := exec.CommandContext(ctx, "just", "deploy-mock-verifier", envFile)
	cmd.Dir = repoRoot

	logger.Info("Executing deploy-mock-verifier", "cmd", strings.Join(cmd.Args, " "))

	return execCommand(cmd, logger)
}

// execDeployOracle runs `just deploy-oracle <envFile>` and parses the output
func execDeployOracle(ctx context.Context, repoRoot, envFile string, logger log.Logger) (string, error) {
	cmd := exec.CommandContext(ctx, "just", "deploy-oracle", envFile)
	cmd.Dir = repoRoot

	logger.Info("Executing deploy-oracle", "cmd", strings.Join(cmd.Args, " "))

	return execCommand(cmd, logger)
}

func execCommand(cmd *exec.Cmd, logger log.Logger) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	stdoutStr := strings.TrimSpace(stdout.String())
	stderrStr := strings.TrimSpace(stderr.String())

	re := regexp.MustCompile(`(?m)^0:\s+address\s+(0x[0-9a-fA-F]{40})\b`)
	match := re.FindStringSubmatch(stdoutStr)
	addr := ""
	if len(match) == 2 {
		addr = match[1]
	}

	if err == nil && addr != "" {
		return addr, nil
	}

	// Foundry errored *but* it's the known transient indexing error AND we have an address
	if err != nil && addr != "" && strings.Contains(stderrStr, "transaction indexing is in progress") {
		logger.Warn("ignoring indexing error and using parsed address", "err", err, "address", addr)
		return addr, nil
	}

	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w\nstdout:\n%s\nstderr:\n%s",
			err, stdoutStr, stderrStr)
	}

	return "", fmt.Errorf("command succeeded but could not find the address.\nstdout:\n%s", stdoutStr)
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
