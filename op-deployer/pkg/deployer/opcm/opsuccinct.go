package opcm

import (
	"os"
	"strings"

	"github.com/ethereum-optimism/optimism/op-chain-ops/script"
	"github.com/ethereum/go-ethereum/common"
)

// SP1ProverMode represents the SP1 prover backend type.
type SP1ProverMode string

const (
	SP1ProverModePlonk   SP1ProverMode = "plonk"
	SP1ProverModeGroth16 SP1ProverMode = "groth16"
)

// GetSP1ProverMode returns the configured SP1 prover mode from environment.
// Defaults to Plonk if not set or invalid.
func GetSP1ProverMode() SP1ProverMode {
	mode := strings.ToLower(os.Getenv("SP1_PROVER_MODE"))
	if mode == "groth16" {
		return SP1ProverModeGroth16
	}
	return SP1ProverModePlonk
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

// NewDeploySP1VerifierScript loads the appropriate SP1Verifier deployment script based on SP1_PROVER_MODE.
// Use this for network proving; use NewDeploySP1MockVerifierScript for local testing.
func NewDeploySP1VerifierScript(host *script.Host) (DeploySP1VerifierScript, error) {
	if GetSP1ProverMode() == SP1ProverModeGroth16 {
		return NewDeploySP1VerifierGroth16Script(host)
	}
	return NewDeploySP1VerifierPlonkScript(host)
}

// NewDeployOPSuccinctScripts deploys OP Succinct contracts at genesis.
func NewDeployOPSuccinctScripts(host *script.Host) DeployOPSuccinctOutput {
	deploySP1Verifier, err := NewDeploySP1VerifierScript(host)
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
