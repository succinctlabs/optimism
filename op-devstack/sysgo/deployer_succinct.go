package sysgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum-optimism/optimism/op-chain-ops/devkeys"
	"github.com/ethereum-optimism/optimism/op-devstack/devtest"
	"github.com/ethereum-optimism/optimism/op-devstack/stack"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/geth"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/wait"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

// OPSuccinctGameType is the game type used by OPSuccinct fault dispute games.
// This must match the GAME_TYPE value used in the deployment scripts.
const OPSuccinctGameType uint32 = 42

// setStandardPortalRespectedGameType reads the AnchorStateRegistry from the standard OptimismPortal2
// and sets its respectedGameType to the specified game type. This is necessary because the
// StandardBridge DSL uses the portal from the rollup config (standard devstack deployment),
// which has its own AnchorStateRegistry separate from the OPSuccinct-deployed one.
func setStandardPortalRespectedGameType(o *Orchestrator, l1ELID stack.L1ELNodeID, l2ChainID eth.ChainID, gameType uint32) {
	p := o.P()
	require := p.Require()
	logger := p.Logger().New("component", "succinct-deployer")

	l1ChainID := l1ELID.ChainID()

	l1EL, ok := o.l1ELs.Get(l1ELID)
	require.True(ok, "l1 EL node required")

	rpcClient, err := rpc.DialContext(p.Ctx(), l1EL.UserRPC())
	require.NoError(err, "failed to dial L1 RPC")
	client := ethclient.NewClient(rpcClient)

	// Get the standard portal address from the L2 network's rollup config
	l2Net, ok := o.l2Nets.Get(l2ChainID)
	require.True(ok, "l2 network required")
	portalAddr := l2Net.rollupCfg.DepositContractAddress

	logger.Info("Reading AnchorStateRegistry from standard OptimismPortal2",
		"portal", portalAddr.Hex())

	// Read anchorStateRegistry() from the portal
	// Function selector: bytes4(keccak256("anchorStateRegistry()")) = 0x72d5fe21
	asrSelector := crypto.Keccak256([]byte("anchorStateRegistry()"))[:4]
	asrResult, err := client.CallContract(p.Ctx(), ethereum.CallMsg{
		To:   &portalAddr,
		Data: asrSelector,
	}, nil)
	require.NoError(err, "failed to read anchorStateRegistry from portal")
	require.Len(asrResult, 32, "unexpected anchorStateRegistry result length")

	anchorStateRegistryAddr := common.BytesToAddress(asrResult[12:32])
	logger.Info("Found standard AnchorStateRegistry",
		"anchorStateRegistry", anchorStateRegistryAddr.Hex())

	// Get the guardian key (SuperchainConfigGuardian - required for standard ASR)
	superOps := devkeys.SuperchainOperatorKeys(l1ChainID.ToBig())
	guardianKey, err := o.keys.Secret(superOps(devkeys.SuperchainConfigGuardianKey))
	require.NoError(err, "failed to get guardian key (SuperchainConfigGuardian)")

	logger.Info("Setting respectedGameType on standard AnchorStateRegistry",
		"anchorStateRegistry", anchorStateRegistryAddr.Hex(),
		"gameType", gameType)

	// AnchorStateRegistry.setRespectedGameType(uint32)
	selector := crypto.Keccak256([]byte("setRespectedGameType(uint32)"))[:4]
	gameTypeBytes := common.LeftPadBytes(big.NewInt(int64(gameType)).Bytes(), 32)
	data := append(selector, gameTypeBytes...)

	nonce, err := client.PendingNonceAt(p.Ctx(), crypto.PubkeyToAddress(guardianKey.PublicKey))
	require.NoError(err, "failed to get nonce")

	gasPrice, err := client.SuggestGasPrice(p.Ctx())
	require.NoError(err, "failed to get gas price")

	tx := types.NewTransaction(nonce, anchorStateRegistryAddr, big.NewInt(0), 100000, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(l1ChainID.ToBig()), guardianKey)
	require.NoError(err, "failed to sign tx")

	err = client.SendTransaction(p.Ctx(), signedTx)
	require.NoError(err, "failed to send setRespectedGameType tx")

	_, err = wait.ForReceiptOK(p.Ctx(), client, signedTx.Hash())
	require.NoError(err, "failed to wait for setRespectedGameType receipt")

	logger.Info("Successfully set respectedGameType on standard AnchorStateRegistry", "txHash", signedTx.Hash().Hex())
}

// setRespectedGameType updates the respectedGameType on AnchorStateRegistry to the specified game type.
// This is necessary for withdrawals to work correctly with the StandardBridge DSL,
// which filters games by the portal's respectedGameType.
// Note: In recent Optimism contract versions, setRespectedGameType is on AnchorStateRegistry,
// not OptimismPortal2. The OPSuccinct FDG deployment creates a MockSystemConfig with the deployer
// (L1ProxyAdminOwnerRole) as the guardian, so we must use that key to call setRespectedGameType.
func setRespectedGameType(o *Orchestrator, l1ELID stack.L1ELNodeID, anchorStateRegistryAddr common.Address, gameType uint32) {
	p := o.P()
	require := p.Require()
	logger := p.Logger().New("component", "succinct-deployer")

	l1ChainID := l1ELID.ChainID()

	l1EL, ok := o.l1ELs.Get(l1ELID)
	require.True(ok, "l1 EL node required")

	rpcClient, err := rpc.DialContext(p.Ctx(), l1EL.UserRPC())
	require.NoError(err, "failed to dial L1 RPC")
	client := ethclient.NewClient(rpcClient)

	// Get the guardian key - the OPSuccinct FDG deployment creates a MockSystemConfig with
	// the deployer (L1ProxyAdminOwnerRole) as the guardian. This is different from the standard
	// devstack which uses SuperchainConfigGuardianKey.
	chainOps := devkeys.ChainOperatorKeys(l1ChainID.ToBig())
	guardianKey, err := o.keys.Secret(chainOps(devkeys.L1ProxyAdminOwnerRole))
	require.NoError(err, "failed to get guardian key (L1ProxyAdminOwner)")

	logger.Info("Setting respectedGameType on AnchorStateRegistry",
		"anchorStateRegistry", anchorStateRegistryAddr.Hex(),
		"gameType", gameType)

	// AnchorStateRegistry.setRespectedGameType(uint32)
	// Selector: bytes4(keccak256("setRespectedGameType(uint32)")) = 0x7fc48504
	// We use a raw call since there's no AnchorStateRegistry binding available
	selector := crypto.Keccak256([]byte("setRespectedGameType(uint32)"))[:4]
	// Encode gameType as uint32 (padded to 32 bytes)
	gameTypeBytes := common.LeftPadBytes(big.NewInt(int64(gameType)).Bytes(), 32)
	data := append(selector, gameTypeBytes...)

	// Send raw transaction to AnchorStateRegistry
	nonce, err := client.PendingNonceAt(p.Ctx(), crypto.PubkeyToAddress(guardianKey.PublicKey))
	require.NoError(err, "failed to get nonce")

	gasPrice, err := client.SuggestGasPrice(p.Ctx())
	require.NoError(err, "failed to get gas price")

	tx := types.NewTransaction(nonce, anchorStateRegistryAddr, big.NewInt(0), 100000, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(l1ChainID.ToBig()), guardianKey)
	require.NoError(err, "failed to sign tx")

	err = client.SendTransaction(p.Ctx(), signedTx)
	require.NoError(err, "failed to send setRespectedGameType tx")

	_, err = wait.ForReceiptOK(p.Ctx(), client, signedTx.Hash())
	require.NoError(err, "failed to wait for setRespectedGameType receipt")

	logger.Info("Successfully set respectedGameType", "txHash", signedTx.Hash().Hex())
}

// =============================================================
// SP1MockVerifier Deployment
// =============================================================

// Deploys an SP1MockVerifier contract for the specified L2 chain, and updates
// the orchestrator's L2 network deployments accordingly.
func WithDeploySP1MockVerifier(
	l1ELID stack.L1ELNodeID,
	l2ChainID eth.ChainID,
) stack.Option[*Orchestrator] {
	return stack.AfterDeploy(func(o *Orchestrator) {
		WithDeploySP1MockVerifierPostDeploy(o, l1ELID, l2ChainID)
	})
}

// Super version of WithDeploySP1MockVerifier that runs in the Finally phase
func WithSuperDeploySP1MockVerifier(
	l1ELID stack.L1ELNodeID,
	l2ChainID eth.ChainID,
) stack.Option[*Orchestrator] {
	return stack.Finally(func(o *Orchestrator) {
		WithDeploySP1MockVerifierPostDeploy(o, l1ELID, l2ChainID)
	})
}

func WithDeploySP1MockVerifierPostDeploy(
	o *Orchestrator,
	l1ELID stack.L1ELNodeID,
	l2ChainID eth.ChainID,
) {
	require := o.P().Require()

	rootPrefix, err := findMonorepoRoot("Cargo.lock")
	require.NoError(err, "failed to locate monorepo root")

	repoRoot, err := filepath.Abs(rootPrefix)
	require.NoError(err, "failed to resolve monorepo root")

	addr, err := o.deploySP1MockVerifier(repoRoot, l1ELID, l2ChainID)
	require.NoError(err, "failed to deploy SP1MockVerifier")

	l2Net, ok := o.l2Nets.Get(l2ChainID)
	o.P().Require().True(ok, "l2 network required")
	l2Net.deployment.sp1MockVerifier = common.HexToAddress(addr)
}

// deploySP1MockVerifier deploys an SP1MockVerifier contract
func (o *Orchestrator) deploySP1MockVerifier(
	repoRoot string,
	l1ELID stack.L1ELNodeID,
	l2ChainID eth.ChainID,
) (string, error) {

	p := o.P()
	logger := p.Logger().New("component", "succinct-deployer")
	require := p.Require()

	l1ChainID := l1ELID.ChainID()

	l1EL, ok := o.GetL1EL(l1ELID)
	require.True(ok, "l1 EL node required")

	l1PAOKey, err := o.GetKeys().Secret(devkeys.L1ProxyAdminOwnerRole.Key(l1ChainID.ToBig()))
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
	if err = WriteEnvFile(envFile, envVars); err != nil {
		return "", fmt.Errorf("failed to write sp1-mock-verifier env: %w", err)
	}

	addr, err := execDeployMockVerifier(o.P(), repoRoot, envFile)
	if err != nil {
		return "", err
	}

	logger.Info("Deployed SP1MockVerifier", "address", addr)
	return addr, nil
}

// execDeployMockVerifier runs `just deploy-mock-verifier <envFile>` and parses the output
func execDeployMockVerifier(p devtest.P, repoRoot, envFile string) (string, error) {
	cmd := exec.CommandContext(p.Ctx(), "just", "deploy-mock-verifier", envFile)
	cmd.Dir = repoRoot

	logger := p.Logger().New("component", "succinct-deployer")

	logger.Info("Executing deploy-mock-verifier", "cmd", strings.Join(cmd.Args, " "))
	stdoutStr, runErr := execCommand(cmd, logger)
	p.Require().NoError(runErr, "failed to execute deploy-mock-verifier command")

	addrMap, err := parseNamedAddresses(stdoutStr, "0")
	return addrMap["0"], err
}

// ============================================================
// OPSuccinctL2OutputOracle Deployment
// ============================================================

// Deploys an OPSuccinctL2OutputOracle contract for each specified chain, and
// updates the orchestrator's L2 network deployments accordingly.
func WithDeployOpSuccinctL2OutputOracle(
	l1CLID stack.L1CLNodeID,
	l1ELID stack.L1ELNodeID,
	l2CLID stack.L2CLNodeID,
	l2ELID stack.L2ELNodeID,
	opts ...L2OOOption,
) stack.Option[*Orchestrator] {
	return stack.AfterDeploy(func(o *Orchestrator) {
		WithDeployOpSuccinctL2OutputOraclePostDeploy(o, l1CLID, l1ELID, l2CLID, l2ELID, opts...)
	})
}

// Super version of WithDeployOpSuccinctL2OutputOracle that runs in the Finally phase
func WithSuperDeployOpSuccinctL2OutputOracle(
	l1CLID stack.L1CLNodeID,
	l1ELID stack.L1ELNodeID,
	l2CLID stack.L2CLNodeID,
	l2ELID stack.L2ELNodeID,
	opts ...L2OOOption,
) stack.Option[*Orchestrator] {
	return stack.Finally(func(o *Orchestrator) {
		WithDeployOpSuccinctL2OutputOraclePostDeploy(o, l1CLID, l1ELID, l2CLID, l2ELID, opts...)
	})
}

func WithDeployOpSuccinctL2OutputOraclePostDeploy(o *Orchestrator,
	l1CLID stack.L1CLNodeID,
	l1ELID stack.L1ELNodeID,
	l2CLID stack.L2CLNodeID,
	l2ELID stack.L2ELNodeID,
	opts ...L2OOOption,
) {
	cfg := &L2OOConfigs{}
	for _, opt := range opts {
		opt(cfg)
	}

	require := o.P().Require()

	rootPrefix, err := findMonorepoRoot("Cargo.lock")
	require.NoError(err, "failed to locate monorepo root")

	repoRoot, err := filepath.Abs(rootPrefix)
	require.NoError(err, "failed to resolve monorepo root")

	addr, err := o.deployOpSuccinctL2OutputOracle(repoRoot, l1CLID, l1ELID, l2CLID, l2ELID, cfg)
	require.NoError(err, "failed to deploy OPSuccinctL2OutputOracle")

	l2Net, ok := o.l2Nets.Get(l2CLID.ChainID())
	o.P().Require().True(ok, "l2 network required")
	l2Net.deployment.opSuccinctL2OutputOracle = common.HexToAddress(addr)
}

// deployOpSuccinctL2OutputOracle deploys an OPSuccinctL2OutputOracle contract
func (o *Orchestrator) deployOpSuccinctL2OutputOracle(
	repoRoot string,
	l1CLID stack.L1CLNodeID,
	l1ELID stack.L1ELNodeID,
	l2CLID stack.L2CLNodeID,
	l2ELID stack.L2ELNodeID,
	cfgs *L2OOConfigs,
) (string, error) {

	p := o.P()
	l2ChainID := l2CLID.ChainID()
	logger := p.Logger().New("component", "succinct-deployer", "chain", l2ChainID.String())
	require := p.Require()

	l1Net, ok := o.l1Nets.Get(l1CLID.ChainID())
	require.True(ok, "l1 network required")

	l1CL, ok := o.GetL1CL(l1CLID)
	require.True(ok, "l1 CL node required")

	l1EL, ok := o.GetL1EL(l1ELID)
	require.True(ok, "l1 EL node required")

	l2Net, ok := o.l2Nets.Get(l2CLID.ChainID())
	require.True(ok, "l2 network required")

	l2CL, ok := o.GetL2CL(l2CLID)
	require.True(ok, "l2 CL node required")

	l2EL, ok := o.GetL2EL(l2ELID)
	require.True(ok, "l2 EL node required")

	l1ChainID := l1CLID.ChainID().ToBig()
	l1PAOKey, err := o.GetKeys().Secret(devkeys.L1ProxyAdminOwnerRole.Key(l1ChainID))
	require.NoError(err, "failed to get L1ProxyAdminOwnerRole key")
	l1PAOKeyStr := hexutil.Encode(crypto.FromECDSA(l1PAOKey))

	proposerKey, err := o.GetKeys().Secret(devkeys.ProposerRole.Key(l2ChainID.ToBig()))
	require.NoError(err, "failed to get ProposerRole key")
	proposerAddr := crypto.PubkeyToAddress(proposerKey.PublicKey)

	finalizationPeriodSecs := resolveFinalizationPeriodSecs(cfgs.FinalizationPeriodSecs)
	startingBlockNumber, err := resolveStartingBlockNumber(p, l2EL.UserRPC(), l2Net.rollupCfg.BlockTime, cfgs.StartingBlockNumber, finalizationPeriodSecs)
	o.P().Require().NoError(err, "failed to resolve starting block number")

	base := p.TempDir()

	l1CfgDir := l1ConfigDir(base)
	err = os.MkdirAll(l1CfgDir, 0o755)
	require.NoError(err, "mkdir l1 config dir")

	l2CfgDir := l2ConfigDir(base)
	err = os.MkdirAll(l2CfgDir, 0o755)
	require.NoError(err, "mkdir l2 config dir")

	// Enables test-specific config file path for parallel test isolation
	l2ooConfigPath := filepath.Join(base, "opsuccinctl2ooconfig.json")

	WithValidityConfigDirsOption(o, l1CfgDir, l2CfgDir)

	verifierAddr, err := l2Net.deployment.resolveSP1VerifierAddr()
	require.NoError(err, "failed to get verifier address")

	envVars := map[string]string{
		"L1_RPC":           l1EL.UserRPC(),
		"L1_BEACON_RPC":    l1CL.beaconHTTPAddr,
		"L2_RPC":           strings.ReplaceAll(l2EL.UserRPC(), "ws://", "http://"),
		"L2_NODE_RPC":      strings.ReplaceAll(l2CL.UserRPC(), "ws://", "http://"),
		"VERIFIER_ADDRESS": verifierAddr.Hex(),
		"PRIVATE_KEY":      l1PAOKeyStr,
		"PROPOSER":         proposerAddr.Hex(),
		"L1_CONFIG_DIR":    l1CfgDir,
		"L2_CONFIG_DIR":    l2CfgDir,
		"OP_SUCCINCT_L2_OUTPUT_ORACLE_CONFIG_PATH": l2ooConfigPath,
		"STARTING_BLOCK_NUMBER":                    fmt.Sprintf("%d", startingBlockNumber),
		"RUST_LOG":                                 "info",
	}

	setEnvIfNotNil(envVars, "SUBMISSION_INTERVAL", cfgs.SubmissionInterval)
	setEnvIfNotNil(envVars, "RANGE_PROOF_INTERVAL", cfgs.RangeProofInterval)

	envDir := p.TempDir()
	envFile := filepath.Join(envDir, fmt.Sprintf("op-succinct-l2oo-%s.env", strings.ReplaceAll(l2ChainID.String(), "-", "_")))
	if err = WriteEnvFile(envFile, envVars); err != nil {
		return "", fmt.Errorf("failed to write op-succinct-l2oo env: %w", err)
	}

	l1ChainConfig := l1Net.genesis.Config

	err = writeL1ChainConfig(l1ChainConfig, l1CLID.ChainID(), l1CfgDir, logger)
	if err != nil {
		return "", fmt.Errorf("failed to write L1 chain config: %w", err)
	}

	addr, err := execDeployOracle(o.P(), repoRoot, envFile)
	if err != nil {
		return "", err
	}

	logger.Info("Deployed OPSuccinctL2OutputOracle", "address", addr)
	return addr, nil
}

// execDeployOracle runs `just deploy-oracle <envFile>` and parses the output
func execDeployOracle(p devtest.P, repoRoot, envFile string) (string, error) {
	cmd := exec.CommandContext(p.Ctx(), "just", "deploy-oracle", envFile)
	cmd.Dir = repoRoot

	logger := p.Logger().New("component", "succinct-deployer")

	logger.Info("Executing deploy-oracle", "cmd", strings.Join(cmd.Args, " "))
	stdoutStr, runErr := execCommand(cmd, logger)
	p.Require().NoError(runErr, "failed to execute deploy-oracle command")

	addrMap, err := parseNamedAddresses(stdoutStr, "0")
	return addrMap["0"], err
}

// L2OOConfigs holds configuration for OPSuccinctL2OutputOracle contract deployment
type L2OOConfigs struct {
	StartingBlockNumber    *uint64
	SubmissionInterval     *uint64
	RangeProofInterval     *uint64
	FinalizationPeriodSecs *uint64
}

type L2OOOption func(*L2OOConfigs)

// WithL2OOStartingBlockNumber sets the starting block number for the L2OO deployment
func WithL2OOStartingBlockNumber(n uint64) L2OOOption {
	return func(cfg *L2OOConfigs) {
		cfg.StartingBlockNumber = &n
	}
}

// WithL2OOSubmissionInterval sets the submission interval for the L2OO contract
func WithL2OOSubmissionInterval(n uint64) L2OOOption {
	return func(cfg *L2OOConfigs) {
		cfg.SubmissionInterval = &n
	}
}

// WithL2OORangeProofInterval sets the range proof interval for the L2OO contract
func WithL2OORangeProofInterval(n uint64) L2OOOption {
	return func(cfg *L2OOConfigs) {
		cfg.RangeProofInterval = &n
	}
}

// WithL2OOFinalizationPeriodSecs sets the finalization period in seconds for the L2OO deployment.
// This determines how long to wait for L2 blocks to finalize before starting the oracle.
// Default is 3600 (1 hour). For e2e tests, a lower value speeds up deployment.
func WithL2OOFinalizationPeriodSecs(n uint64) L2OOOption {
	return func(cfg *L2OOConfigs) {
		cfg.FinalizationPeriodSecs = &n
	}
}

const defaultFinalizationPeriodSecs = 3600

// resolveFinalizationPeriodSecs returns the configured finalization period or the default.
func resolveFinalizationPeriodSecs(cfgFinalizationPeriodSecs *uint64) uint64 {
	if cfgFinalizationPeriodSecs != nil {
		return *cfgFinalizationPeriodSecs
	}
	return uint64(defaultFinalizationPeriodSecs)
}

// resolveStartingBlockNumber determines the starting block number for L2OO and FDG deployments
func resolveStartingBlockNumber(p devtest.P, l2Rpc string, l2BlockTime uint64, cfgStartingBlockNumber *uint64, finalizationPeriodSecs uint64) (uint64, error) {
	logger := p.Logger().New("component", "succinct-deployer")

	var v uint64
	if cfgStartingBlockNumber != nil {
		v = *cfgStartingBlockNumber
	} else {
		v = finalizationPeriodSecs/l2BlockTime + 1
	}
	target := new(big.Int).SetUint64(v)

	res, err := ethclient.DialContext(p.Ctx(), l2Rpc)
	if err != nil {
		return 0, err
	}
	defer res.Close()

	block, err := geth.WaitForBlockToBeFinalized(target, res, 90*time.Minute)
	if err != nil {
		logger.Warn("L2 chain did not reach finalized block within timeout", "err", err)
		return 0, err
	}
	blockNumber := block.Number().Uint64()
	logger.Info("Finalized L2 block reached; proceeding with deployment", "block", blockNumber)

	return blockNumber, nil
}

// ===========================================================
// OPSuccinctFaultDisputeGame Deployment
// ===========================================================

const (
	disputeGameFinalityDelaySecsDefault = 604800 // 7 days in seconds
	maxChallengeDurationDefault         = 604800 // 7 days in seconds
	maxProveDurationDefault             = 86400  // 1 day in seconds
)

// Deploys an OPSuccinctFaultDisputeGame contract for each specified chain, and
// updates the orchestrator's L2 network deployments accordingly.
func WithDeployOPSuccinctFaultDisputeGame(
	l1CLID stack.L1CLNodeID,
	l1ELID stack.L1ELNodeID,
	l2CLID stack.L2CLNodeID,
	l2ELID stack.L2ELNodeID,
	opts ...FdgOption,
) stack.Option[*Orchestrator] {
	return stack.AfterDeploy(func(o *Orchestrator) {
		WithDeployOPSuccinctFaultDisputeGamePostDeploy(o, l1CLID, l1ELID, l2CLID, l2ELID, opts...)
	})
}

// Super version of WithDeployOPSuccinctFaultDisputeGame that runs in the Finally phase
func WithSuperDeployOPSuccinctFaultDisputeGame(
	l1CLID stack.L1CLNodeID,
	l1ELID stack.L1ELNodeID,
	l2CLID stack.L2CLNodeID,
	l2ELID stack.L2ELNodeID,
	opts ...FdgOption,
) stack.Option[*Orchestrator] {
	return stack.Finally(func(o *Orchestrator) {
		WithDeployOPSuccinctFaultDisputeGamePostDeploy(o, l1CLID, l1ELID, l2CLID, l2ELID, opts...)
	})
}

func WithDeployOPSuccinctFaultDisputeGamePostDeploy(o *Orchestrator,
	l1CLID stack.L1CLNodeID,
	l1ELID stack.L1ELNodeID,
	l2CLID stack.L2CLNodeID,
	l2ELID stack.L2ELNodeID,
	opts ...FdgOption,
) {
	cfg := &FdgConfigs{}
	for _, opt := range opts {
		opt(cfg)
	}

	require := o.P().Require()

	rootPrefix, err := findMonorepoRoot("Cargo.lock")
	require.NoError(err, "failed to locate monorepo root")

	repoRoot, err := filepath.Abs(rootPrefix)
	require.NoError(err, "failed to resolve monorepo root")

	addrs, err := o.deployOpSuccinctFaultDisputeGame(repoRoot, l1CLID, l1ELID, l2CLID, l2ELID, cfg)
	require.NoError(err, "failed to deploy OPSuccinctFaultDisputeGame")

	l2Net, ok := o.l2Nets.Get(l2CLID.ChainID())
	o.P().Require().True(ok, "l2 network required")
	l2Net.deployment.sp1Verifier = addrs.Sp1Verifier
	l2Net.deployment.anchorStateRegistry = addrs.AnchorStateRegistry
	l2Net.deployment.disputeGameFactoryProxy = addrs.FactoryProxy

	// Set respectedGameType to 42 (OPSuccinct game type) on BOTH AnchorStateRegistries:
	// 1. The OPSuccinct-deployed AnchorStateRegistry (for the OPSuccinct DGF)
	// 2. The standard devstack's AnchorStateRegistry (for the StandardBridge DSL which uses the standard portal)
	setRespectedGameType(o, l1ELID, addrs.AnchorStateRegistry, OPSuccinctGameType)
	setStandardPortalRespectedGameType(o, l1ELID, l2CLID.ChainID(), OPSuccinctGameType)
}

// deployOpSuccinctFaultDisputeGame deploys an OPSuccinctFaultDisputeGame contract
func (o *Orchestrator) deployOpSuccinctFaultDisputeGame(
	repoRoot string,
	l1CLID stack.L1CLNodeID,
	l1ELID stack.L1ELNodeID,
	l2CLID stack.L2CLNodeID,
	l2ELID stack.L2ELNodeID,
	cfgs *FdgConfigs,
) (FdgAddresses, error) {

	p := o.P()
	l2ChainID := l2CLID.ChainID()
	logger := p.Logger().New("component", "succinct-deployer", "chain", l2ChainID.String())
	require := p.Require()

	l1Net, ok := o.l1Nets.Get(l1CLID.ChainID())
	require.True(ok, "l1 network required")

	l1CL, ok := o.GetL1CL(l1CLID)
	require.True(ok, "l1 CL node required")

	l1EL, ok := o.GetL1EL(l1ELID)
	require.True(ok, "l1 EL node required")

	l2Net, ok := o.l2Nets.Get(l2CLID.ChainID())
	require.True(ok, "l2 network required")

	l2CL, ok := o.GetL2CL(l2CLID)
	require.True(ok, "l2 CL node required")

	l2EL, ok := o.GetL2EL(l2ELID)
	require.True(ok, "l2 EL node required")

	l1ChainID := l1CLID.ChainID().ToBig()
	l1PAOKey, err := o.GetKeys().Secret(devkeys.L1ProxyAdminOwnerRole.Key(l1ChainID))
	if err != nil {
		return FdgAddresses{}, fmt.Errorf("failed to get L1ProxyAdminOwnerRole key: %w", err)
	}
	l1PAOKeyStr := hexutil.Encode(crypto.FromECDSA(l1PAOKey))

	disputeGameFinalityDelaySecs := resolveDisputeGameFinalityDelaySecs(cfgs.disputeGameFinalityDelaySecs)
	maxChallengeDuration := resolveMaxChallengeDuration(cfgs.maxChallengeDuration)
	maxProveDuration := resolveMaxProveDuration(cfgs.maxProveDuration)

	startingL2BlockNumber, err := resolveStartingBlockNumber(p, l2EL.UserRPC(), l2Net.rollupCfg.BlockTime, cfgs.startingL2BlockNumber, disputeGameFinalityDelaySecs)
	o.P().Require().NoError(err, "failed to resolve starting block number")

	base := p.TempDir()

	l1CfgDir := l1ConfigDir(base)
	err = os.MkdirAll(l1CfgDir, 0o755)
	require.NoError(err, "mkdir l1 config dir")

	l2CfgDir := l2ConfigDir(base)
	err = os.MkdirAll(l2CfgDir, 0o755)
	require.NoError(err, "mkdir l2 config dir")

	// Enables test-specific config file path for parallel test isolation
	fdgConfigPath := filepath.Join(base, "opsuccinctfdgconfig.json")

	WithFPConfigDirsOption(o, l1CfgDir, l2CfgDir)

	verifierAddr, err := l2Net.deployment.resolveSP1VerifierAddr()
	require.NoError(err, "failed to get verifier address")

	envVars := map[string]string{
		"L1_RPC":                              l1EL.UserRPC(),
		"L1_BEACON_RPC":                       l1CL.beaconHTTPAddr,
		"L2_RPC":                              strings.ReplaceAll(l2EL.UserRPC(), "ws://", "http://"),
		"L2_NODE_RPC":                         strings.ReplaceAll(l2CL.UserRPC(), "ws://", "http://"),
		"GAME_TYPE":                           "42",
		"DISPUTE_GAME_FINALITY_DELAY_SECONDS": fmt.Sprintf("%d", disputeGameFinalityDelaySecs),
		"MAX_CHALLENGE_DURATION":              fmt.Sprintf("%d", maxChallengeDuration),
		"MAX_PROVE_DURATION":                  fmt.Sprintf("%d", maxProveDuration),
		"VERIFIER_ADDRESS":                    verifierAddr.Hex(),
		"PRIVATE_KEY":                         l1PAOKeyStr,
		"STARTING_L2_BLOCK_NUMBER":            fmt.Sprintf("%d", startingL2BlockNumber),
		"L1_CONFIG_DIR":                       l1CfgDir,
		"L2_CONFIG_DIR":                       l2CfgDir,
		"OP_SUCCINCT_FAULT_DISPUTE_GAME_CONFIG_PATH": fdgConfigPath,
		"PERMISSIONLESS_MODE":                        "true",
		"OP_SUCCINCT_MOCK":                           strconv.FormatBool(os.Getenv("NETWORK_PRIVATE_KEY") == ""),
		"RUST_LOG":                                   "info",
	}

	envDir := p.TempDir()
	envFile := filepath.Join(envDir, fmt.Sprintf("op-succinct-fdg-%s.env", strings.ReplaceAll(l2ChainID.String(), "-", "_")))
	if err = WriteEnvFile(envFile, envVars); err != nil {
		return FdgAddresses{}, fmt.Errorf("failed to write op-succinct-fdg env: %w", err)
	}

	l1ChainConfig := l1Net.genesis.Config

	err = writeL1ChainConfig(l1ChainConfig, l1CLID.ChainID(), l1CfgDir, logger)
	if err != nil {
		return FdgAddresses{}, fmt.Errorf("failed to write L1 chain config: %w", err)
	}

	addrs, err := execDeployFdgContracts(o.P(), repoRoot, envFile)
	if err != nil {
		return FdgAddresses{}, err
	}

	logger.Info("Deployed OPSuccinctFaultDisputeGame", "address", addrs)
	return addrs, nil
}

// execDeployFdgContracts runs `just deploy-fdg-contracts <envFile>` and parses the output
func execDeployFdgContracts(p devtest.P, repoRoot, envFile string) (FdgAddresses, error) {
	cmd := exec.CommandContext(p.Ctx(), "just", "deploy-fdg-contracts", envFile)
	cmd.Dir = repoRoot

	logger := p.Logger().New("component", "succinct-deployer")

	logger.Info("Executing deploy-fdg-contracts", "cmd", strings.Join(cmd.Args, " "))
	stdoutStr, err := execCommand(cmd, logger)
	p.Require().NoError(err, "failed to execute deploy-fdg-contracts command")

	return parseDeploymentAddresses(stdoutStr)
}

type FdgAddresses struct {
	AnchorStateRegistry common.Address
	FactoryProxy        common.Address
	Sp1Verifier         common.Address
}

func parseDeploymentAddresses(stdoutStr string) (FdgAddresses, error) {
	m, err := parseNamedAddresses(stdoutStr,
		"anchorStateRegistry",
		"factoryProxy",
		"sp1Verifier",
	)
	if err != nil {
		return FdgAddresses{}, err
	}

	return FdgAddresses{
		AnchorStateRegistry: common.HexToAddress(m["anchorStateRegistry"]),
		FactoryProxy:        common.HexToAddress(m["factoryProxy"]),
		Sp1Verifier:         common.HexToAddress(m["sp1Verifier"]),
	}, nil
}

type FdgConfigs struct {
	startingL2BlockNumber        *uint64
	disputeGameFinalityDelaySecs *uint64
	maxChallengeDuration         *uint64
	maxProveDuration             *uint64
}

type FdgOption func(*FdgConfigs)

// WithFdgL2StartingBlockNumber sets the starting block number for the FDG deployment
func WithFdgL2StartingBlockNumber(n uint64) FdgOption {
	return func(cfg *FdgConfigs) {
		cfg.startingL2BlockNumber = &n
	}
}

// WithFdgDisputeGameFinalityDelaySecs sets the starting block number for the FDG deployment
func WithFdgDisputeGameFinalityDelaySecs(n uint64) FdgOption {
	return func(cfg *FdgConfigs) {
		cfg.disputeGameFinalityDelaySecs = &n
	}
}

// WithFdgMaxChallengeDuration sets the max challenge duration for the FDG deployment
func WithFdgMaxChallengeDuration(n uint64) FdgOption {
	return func(cfg *FdgConfigs) {
		cfg.maxChallengeDuration = &n
	}
}

// WithFdgMaxProveDuration sets the max prove duration for the FDG deployment
func WithFdgMaxProveDuration(n uint64) FdgOption {
	return func(cfg *FdgConfigs) {
		cfg.maxProveDuration = &n
	}
}

func resolveDisputeGameFinalityDelaySecs(cfgDisputeGameFinalityDelaySecs *uint64) uint64 {
	if cfgDisputeGameFinalityDelaySecs != nil {
		return *cfgDisputeGameFinalityDelaySecs
	}

	return disputeGameFinalityDelaySecsDefault
}

func resolveMaxChallengeDuration(cfgMaxChallengeDuration *uint64) uint64 {
	if cfgMaxChallengeDuration != nil {
		return *cfgMaxChallengeDuration
	}

	return maxChallengeDurationDefault
}

func resolveMaxProveDuration(cfgMaxProveDuration *uint64) uint64 {
	if cfgMaxProveDuration != nil {
		return *cfgMaxProveDuration
	}

	return maxProveDurationDefault
}

// ============================================================
// Succinct Deployment Helper Functions
// ============================================================

func writeL1ChainConfig(
	l1ChainConfig any,
	l1ChainID fmt.Stringer,
	dir string,
	logger log.Logger,
) error {
	path := filepath.Join(dir, l1ChainID.String()+".json")
	logger.Info("writing L1 chain config", "path", path)

	data, err := json.Marshal(l1ChainConfig)
	if err != nil {
		return fmt.Errorf("marshal L1 chain config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write L1 chain config: %w", err)
	}

	return nil
}

func l1ConfigDir(base string) string {
	return filepath.Join(base, "Configs", "L1")
}

func l2ConfigDir(base string) string {
	return filepath.Join(base, "Configs", "L2")
}

func execCommand(cmd *exec.Cmd, logger log.Logger) (stdoutStr string, err error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	stdoutStr = strings.TrimSpace(stdout.String())
	stderrStr := strings.TrimSpace(stderr.String())

	if err == nil {
		return stdoutStr, nil
	}

	// Foundry errored *but* it's the known transient indexing error AND we have result
	if strings.Contains(stderrStr, "transaction indexing is in progress") && stdoutStr != "" {
		logger.Warn("ignoring indexing error and using parsed addresses", "err", err)
		return stdoutStr, nil
	}

	return stdoutStr, fmt.Errorf(
		"failed to execute command: %w\nstdout:\n%s\nstderr:\n%s",
		err, stdoutStr, stderrStr,
	)
}

// name: address 0x....
var namedAddrRE = regexp.MustCompile(`(?m)^([A-Za-z0-9_]+):\s+address\s+(0x[0-9a-fA-F]{40})\b`)

// parseNamedAddresses scans stdoutStr for lines of the form:
//
//	<name>: address 0x...
//
// and returns a map[name]address only for the requested names.
// It errors if any requested name is missing.
func parseNamedAddresses(stdoutStr string, names ...string) (map[string]string, error) {
	needed := make(map[string]bool, len(names))
	for _, n := range names {
		needed[n] = true
	}

	result := make(map[string]string, len(names))

	matches := namedAddrRE.FindAllStringSubmatch(stdoutStr, -1)
	for _, m := range matches {
		if len(m) != 3 {
			continue
		}
		name := m[1]
		addr := m[2]

		if needed[name] {
			result[name] = addr
		}
	}

	for _, n := range names {
		if _, ok := result[n]; !ok {
			return nil, fmt.Errorf("missing expected address: %s", n)
		}
	}

	return result, nil
}

// WriteEnvFile writes key-value pairs to a file in .env format.
func WriteEnvFile(path string, kv map[string]string) error {
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
