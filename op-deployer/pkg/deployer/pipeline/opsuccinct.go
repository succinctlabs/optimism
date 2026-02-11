package pipeline

import (
	"fmt"

	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/opcm"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/state"
	"github.com/ethereum/go-ethereum/common"
)

func DeployOPSuccinct(env *Env, intent *state.Intent, st *state.State, chainID common.Hash) error {
	lgr := env.Logger.New("stage", "deploy-opsuccinct", "chain", chainID.Hex())

	if !shouldDeployOPSuccinct(intent) {
		lgr.Info("op-succinct deployment not needed")
		return nil
	}

	lgr.Info("Deploying OP Succinct contracts")

	mode, err := sp1ProofModeFromIntent(intent)
	if err != nil {
		return err
	}
	output := opcm.NewDeployOPSuccinctScripts(env.L1ScriptHost, mode)

	chainState, err := st.Chain(chainID)
	if err != nil {
		return fmt.Errorf("failed to get chain state: %w", err)
	}

	chainState.SP1Verifier = output.SP1Verifier
	chainState.SP1MockVerifier = output.SP1MockVerifier

	return nil
}

func shouldDeployOPSuccinct(intent *state.Intent) bool {
	if intent == nil || len(intent.GlobalDeployOverrides) == 0 {
		return false
	}
	raw, ok := intent.GlobalDeployOverrides[opcm.UseOPSuccinctOverrideKey]
	if !ok || raw == nil {
		return false
	}
	use, ok := raw.(bool)
	return ok && use
}

func sp1ProofModeFromIntent(intent *state.Intent) (opcm.SP1ProofMode, error) {
	if intent == nil || len(intent.GlobalDeployOverrides) == 0 {
		return opcm.SP1ProofModePlonk, nil
	}
	raw, ok := intent.GlobalDeployOverrides[opcm.SP1ProofModeOverrideKey]
	if !ok || raw == nil {
		return opcm.SP1ProofModePlonk, nil
	}
	modeStr, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string (got %T)", opcm.SP1ProofModeOverrideKey, raw)
	}
	return opcm.ParseSP1ProofMode(modeStr), nil
}
