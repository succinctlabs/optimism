package opcm

import (
	"fmt"

	"github.com/ethereum-optimism/optimism/op-chain-ops/script"
)

// Scripts contains all the deployment scripts for ease of passing them around
type Scripts struct {
	DeployAlphabetVM      DeployAlphabetVMScript
	DeployAltDA           DeployAltDAScript
	DeployAsterisc        DeployAsteriscScript
	DeployDisputeGame     DeployDisputeGameScript
	DeployImplementations DeployImplementationsScript
	DeployMIPS            DeployMIPSScript
	DeploySuperchain      DeploySuperchainScript
	DeployOPChain         DeployOPChainScript
	DeploySP1MockVerifier DeploySP1MockVerifierScript
	OPSuccinctDeployer    OPSuccinctDeployerScript
}

// NewScripts collects all the deployment scripts, raising exceptions if any of them
// are not found or if the Go types don't match the ABI
func NewScripts(host *script.Host) (*Scripts, error) {
	deployImplementations, err := NewDeployImplementationsScript(host)
	if err != nil {
		return nil, fmt.Errorf("failed to load DeployImplementations script: %w", err)
	}

	deploySuperchain, err := NewDeploySuperchainScript(host)
	if err != nil {
		return nil, fmt.Errorf("failed to load DeploySuperchain script: %w", err)
	}

	deployAlphabetVM, err := NewDeployAlphabetVMScript(host)
	if err != nil {
		return nil, fmt.Errorf("failed to load DeployAlphabetVM script: %w", err)
	}

	deployAltDA, err := NewDeployAltDAScript(host)
	if err != nil {
		return nil, fmt.Errorf("failed to load DeployAltDA script: %w", err)
	}

	deployAsterisc, err := NewDeployAsteriscScript(host)
	if err != nil {
		return nil, fmt.Errorf("failed to load DeployAsterisc script: %w", err)
	}

	deployDisputeGame, err := NewDeployDisputeGameScript(host)
	if err != nil {
		return nil, fmt.Errorf("failed to load DeployDisputeGame script: %w", err)
	}

	deployMIPSScript, err := NewDeployMIPSScript(host)
	if err != nil {
		return nil, fmt.Errorf("failed to load DeployMIPSScript script: %w", err)
	}

	deployOPChain, err := NewDeployOPChainScript(host)
	if err != nil {
		return nil, fmt.Errorf("failed to load DeployOPChain script: %w", err)
	}

	deploySP1MockVerifier, err := NewDeploySP1MockVerifierScript(host)
	if err != nil {
		return nil, fmt.Errorf("failed to load DeploySP1MockVerifier script: %w", err)
	}

	opSuccinctDeployer, err := NewOPSuccinctDeployerScript(host)
	if err != nil {
		return nil, fmt.Errorf("failed to load OPSuccinctDeployer script: %w", err)
	}

	return &Scripts{
		DeployAlphabetVM:      deployAlphabetVM,
		DeployAltDA:           deployAltDA,
		DeployAsterisc:        deployAsterisc,
		DeployDisputeGame:     deployDisputeGame,
		DeployMIPS:            deployMIPSScript,
		DeployImplementations: deployImplementations,
		DeploySuperchain:      deploySuperchain,
		DeployOPChain:         deployOPChain,
		DeploySP1MockVerifier: deploySP1MockVerifier,
		OPSuccinctDeployer:    opSuccinctDeployer,
	}, nil
}
