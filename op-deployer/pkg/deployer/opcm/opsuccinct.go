package opcm

import (
	"strings"

	"github.com/ethereum-optimism/optimism/op-chain-ops/script"
	"github.com/ethereum/go-ethereum/common"
)

// SP1ProofMode represents the SP1 prover backend type.
type SP1ProofMode string

const (
	SP1ProofModePlonk   SP1ProofMode = "plonk"
	SP1ProofModeGroth16 SP1ProofMode = "groth16"
)

// UseOPSuccinctOverrideKey is the global deploy override key that enables OP Succinct deployment.
const UseOPSuccinctOverrideKey = "useOPSuccinct"

// SP1ProofModeOverrideKey is the global deploy override key used to select the SP1 proof mode.
// This is consumed by the op-deployer OP Succinct pipeline stage.
const SP1ProofModeOverrideKey = "sp1ProofMode"

// ParseSP1ProofMode parses a user-provided proof mode string.
// Defaults to Plonk if unset or unrecognized.
func ParseSP1ProofMode(mode string) SP1ProofMode {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case string(SP1ProofModeGroth16):
		return SP1ProofModeGroth16
	default:
		return SP1ProofModePlonk
	}
}

type DeployOPSuccinctOutput struct {
	SP1Verifier     common.Address
	SP1MockVerifier common.Address
}

type DeploySP1MockVerifierScript script.DeployScriptWithoutInput[common.Address]

func NewDeploySP1MockVerifierScript(host *script.Host) (DeploySP1MockVerifierScript, error) {
	return script.NewDeployScriptWithoutInputFromFile[common.Address](host, "DeployMockVerifier.s.sol", "DeployMockVerifier")
}

type DeploySP1VerifierScript script.DeployScriptWithoutInput[common.Address]

// NewDeploySP1VerifierPlonkScript loads the SP1Verifier Plonk deployment script (v5.0.0).
func NewDeploySP1VerifierPlonkScript(host *script.Host) (DeploySP1VerifierScript, error) {
	return script.NewDeployScriptWithoutInputFromFile[common.Address](host, "DeployVerifier.s.sol", "DeployVerifierPlonk")
}

// NewDeploySP1VerifierGroth16Script loads the SP1Verifier Groth16 deployment script (v5.0.0).
func NewDeploySP1VerifierGroth16Script(host *script.Host) (DeploySP1VerifierScript, error) {
	return script.NewDeployScriptWithoutInputFromFile[common.Address](host, "DeployVerifier.s.sol", "DeployVerifierGroth16")
}

// NewDeploySP1VerifierScript loads the appropriate SP1Verifier deployment script based on mode.
// Use this for network proving; use NewDeploySP1MockVerifierScript for local testing.
func NewDeploySP1VerifierScript(host *script.Host, mode SP1ProofMode) (DeploySP1VerifierScript, error) {
	if mode == SP1ProofModeGroth16 {
		return NewDeploySP1VerifierGroth16Script(host)
	}
	return NewDeploySP1VerifierPlonkScript(host)
}

// NewDeployOPSuccinctScripts deploys OP Succinct contracts at genesis.
func NewDeployOPSuccinctScripts(host *script.Host, mode SP1ProofMode) DeployOPSuccinctOutput {
	deploySP1Verifier, err := NewDeploySP1VerifierScript(host, mode)
	if err != nil {
		panic("failed to load DeploySP1Verifier script: " + err.Error())
	}

	sp1VerifierAddr, err := deploySP1Verifier.Run()
	if err != nil {
		panic("failed to deploy SP1Verifier: " + err.Error())
	}

	deploySP1MockVerifier, err := NewDeploySP1MockVerifierScript(host)
	if err != nil {
		panic("failed to load DeploySP1MockVerifier script: " + err.Error())
	}

	sp1MockVerifierAddr, err := deploySP1MockVerifier.Run()
	if err != nil {
		panic("failed to deploy SP1MockVerifier: " + err.Error())
	}

	return DeployOPSuccinctOutput{
		SP1Verifier:     sp1VerifierAddr,
		SP1MockVerifier: sp1MockVerifierAddr,
	}
}
