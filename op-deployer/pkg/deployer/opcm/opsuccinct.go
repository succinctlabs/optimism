package opcm

import (
	"github.com/ethereum-optimism/optimism/op-chain-ops/script"
	"github.com/ethereum/go-ethereum/common"
)

type DeployOPSuccinctOutput struct {
	SP1MockVerifier          common.Address
	OPSuccinctL2OutputOracle common.Address
}

type DeploySP1MockVerifierScript script.DeployScriptWithoutInput[common.Address]

func NewDeploySP1MockVerifierScript(host *script.Host) (DeploySP1MockVerifierScript, error) {
	return script.NewDeployScriptWithoutInputFromFile[common.Address](host, "DeployMockVerifier.s.sol", "DeployMockVerifier")
}

type OPSuccinctDeployerScript script.DeployScriptWithoutInput[common.Address]

func NewOPSuccinctDeployerScript(host *script.Host) (OPSuccinctDeployerScript, error) {
	return script.NewDeployScriptWithoutInputFromFile[common.Address](host, "OPSuccinctDeployer.s.sol", "OPSuccinctDeployer")
}

func NewDeployOPSuccinctScripts(host *script.Host) DeployOPSuccinctOutput {
	deploySP1MockVerifier, err := NewDeploySP1MockVerifierScript(host)
	if err != nil {
		panic("failed to load DeploySP1MockVerifier script: " + err.Error())
	}

	sp1MockVerifierAddr, err := deploySP1MockVerifier.Run()

	opSuccinctDeployer, err := NewOPSuccinctDeployerScript(host)
	if err != nil {
		panic("failed to load OPSuccinctDeployer script: " + err.Error())
	}

	opSuccinctL2OutputOracleAddr, err := opSuccinctDeployer.Run()

	return DeployOPSuccinctOutput{
		SP1MockVerifier:          sp1MockVerifierAddr,
		OPSuccinctL2OutputOracle: opSuccinctL2OutputOracleAddr,
	}
}
