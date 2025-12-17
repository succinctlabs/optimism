package pipeline

import (
	"fmt"

	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/opcm"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/state"
	"github.com/ethereum/go-ethereum/common"
)

func DeployOPSuccinct(env *Env, intent *state.Intent, st *state.State, chainID common.Hash) error {
	lgr := env.Logger.New("stage", "deploy-opsuccinct", "chain", chainID.Hex())
	lgr.Info("Deploying OP Succinct contracts")

	output := opcm.NewDeployOPSuccinctScripts(env.L1ScriptHost)

	chainState, err := st.Chain(chainID)
	if err != nil {
		return fmt.Errorf("failed to get chain state: %w", err)
	}

	chainState.SP1Verifier = output.SP1Verifier
	chainState.SP1MockVerifier = output.SP1MockVerifier

	return nil
}
