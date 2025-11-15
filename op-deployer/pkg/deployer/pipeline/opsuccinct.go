package pipeline

import (
	"fmt"

	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/state"
	"github.com/ethereum/go-ethereum/common"
)

func DeploySP1MockVerifier(env *Env, st *state.State, chainID common.Hash) error {
	lgr := env.Logger.New("stage", "deploy-sp1-mock-verifier")

	lgr.Info("deploying sp1 mock verifier")

	address, err := env.Scripts.DeploySP1MockVerifier.Run()
	if err != nil {
		lgr.Error("failed to deploy sp1 mock verifier", "error", err)
		return fmt.Errorf("failed to deploy sp1 mock verifier: %w", err)
	}

	cs := findChainState(st, chainID)
	if cs == nil {
		cs = &state.ChainState{
			ID: chainID,
		}
		st.Chains = append(st.Chains, cs)
	}

	cs.OpSuccinctContracts.SP1MockVerifier = address

	lgr.Info("deployed sp1 mock verifier", "address", address)

	return nil
}

func DeployOPSuccinctL2OutputOracle(env *Env, st *state.State, chainID common.Hash) error {
	lgr := env.Logger.New("stage", "deploy-op-succinct-l2-output-oracle")

	lgr.Info("deploying op-succinct l2 output oracle")

	address, err := env.Scripts.OPSuccinctDeployer.Run()
	if err != nil {
		lgr.Error("failed to deploy op-succinct l2 output oracle", "error", err)
		return fmt.Errorf("failed to deploy op-succinct l2 output oracle: %w", err)
	}

	cs := findChainState(st, chainID)
	if cs == nil {
		cs = &state.ChainState{
			ID: chainID,
		}
		st.Chains = append(st.Chains, cs)
	}

	cs.OpSuccinctContracts.OPSuccinctL2OutputOracle = address

	lgr.Info("deployed op-succinct l2 output oracle", "address", address)

	return nil
}

func findChainState(st *state.State, chainID common.Hash) *state.ChainState {
	for _, cs := range st.Chains {
		if cs.ID == chainID {
			return cs
		}
	}
	return nil
}
