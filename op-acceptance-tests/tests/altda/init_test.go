package altda

import (
	"testing"

	"github.com/ethereum-optimism/optimism/op-chain-ops/devkeys"
	"github.com/ethereum-optimism/optimism/op-devstack/presets"
	"github.com/ethereum-optimism/optimism/op-devstack/stack"
	"github.com/ethereum-optimism/optimism/op-devstack/sysgo"
)

// minimalSystemNoChallenger mirrors DefaultMinimalSystem but omits the L2 challenger
// (which requires the cannon binary) since AltDA tests don't need fault proofs.
func minimalSystemNoChallenger() stack.Option[*sysgo.Orchestrator] {
	ids := sysgo.NewDefaultMinimalSystemIDs(sysgo.DefaultL1ID, sysgo.DefaultL2AID)
	dest := &sysgo.DefaultMinimalSystemIDs{}

	opt := stack.Combine[*sysgo.Orchestrator]()

	opt.Add(stack.BeforeDeploy(func(o *sysgo.Orchestrator) {
		o.P().Logger().Info("Setting up (AltDA, no challenger)")
	}))

	opt.Add(sysgo.WithMnemonicKeys(devkeys.TestMnemonic))

	opt.Add(sysgo.WithDeployer(),
		sysgo.WithDeployerOptions(
			sysgo.WithLocalContractSources(),
			sysgo.WithCommons(ids.L1.ChainID()),
			sysgo.WithPrefundedL2(ids.L1.ChainID(), ids.L2.ChainID()),
		),
	)

	opt.Add(sysgo.WithL1Nodes(ids.L1EL, ids.L1CL))

	opt.Add(sysgo.WithL2ELNode(ids.L2EL))
	opt.Add(sysgo.WithL2CLNode(ids.L2CL, ids.L1CL, ids.L1EL, ids.L2EL, sysgo.L2CLSequencer()))

	opt.Add(sysgo.WithBatcher(ids.L2Batcher, ids.L1EL, ids.L2CL, ids.L2EL))
	opt.Add(sysgo.WithProposer(ids.L2Proposer, ids.L1EL, &ids.L2CL, nil))

	opt.Add(sysgo.WithFaucets([]stack.L1ELNodeID{ids.L1EL}, []stack.L2ELNodeID{ids.L2EL}))

	opt.Add(sysgo.WithTestSequencer(ids.TestSequencer, ids.L1CL, ids.L2CL, ids.L1EL, ids.L2EL))

	// Deliberately omit WithL2Challenger — AltDA tests don't need fault proofs.

	opt.Add(stack.Finally(func(orch *sysgo.Orchestrator) {
		*dest = ids
	}))

	return opt
}

func TestMain(m *testing.M) {
	presets.DoMain(m,
		stack.MakeCommon(sysgo.WithAltDA(sysgo.DefaultL2AID)),
		stack.MakeCommon[*sysgo.Orchestrator](minimalSystemNoChallenger()),
	)
}
