package presets

import (
	"github.com/ethereum-optimism/optimism/op-devstack/devtest"
	"github.com/ethereum-optimism/optimism/op-devstack/dsl"
	"github.com/ethereum-optimism/optimism/op-devstack/shim"
	"github.com/ethereum-optimism/optimism/op-devstack/stack/match"
)

type MinimalWithProposer struct {
	Minimal

	L2Proposer *dsl.L2Proposer
}

func NewMinimalWithProposer(t devtest.T) *MinimalWithProposer {
	system := shim.NewSystem(t)
	orch := Orchestrator()
	orch.Hydrate(system)
	minimal := MinimalFromSystem(t, system, orch)
	l2 := system.L2Network(match.Assume(t, match.L2ChainA))
	proposer := l2.L2Proposer(match.Assume(t, match.FirstL2Proposer))
	return &MinimalWithProposer{
		Minimal:    *minimal,
		L2Proposer: dsl.NewL2Proposer(proposer),
	}
}
