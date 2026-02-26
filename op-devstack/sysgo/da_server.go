package sysgo

import (
	"os"

	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	bss "github.com/ethereum-optimism/optimism/op-batcher/batcher"
	"github.com/ethereum-optimism/optimism/op-devstack/devtest"
	"github.com/ethereum-optimism/optimism/op-devstack/stack"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

// WithAltDA starts a DA server and configures the batcher and op-node to use AltDA mode
// with GenericCommitment (no on-chain challenge contract needed).
//
// The DA server is started and batcher/L2CL options are accumulated during the Deploy phase.
// The rollup config is overridden via L2CLConfig.RollupAltDAConfig, which is applied
// inside op-node's AfterDeploy hook (after the deployer registers the L2 network).
// This avoids ordering dependencies with the deployer's AfterDeploy.
func WithAltDA(l2ChainID eth.ChainID) stack.Option[*Orchestrator] {
	return stack.Deploy(func(orch *Orchestrator) {
		p := orch.P()
		logger := p.Logger()

		// Start DA server with in-memory storage
		store := altda.NewMemStore()
		server := altda.NewDAServer("127.0.0.1", 0, store, logger, true)
		p.Require().NoError(server.Start(), "failed to start DA server")
		p.Cleanup(func() {
			logger.Info("Stopping DA server")
			_ = server.Stop()
		})

		endpoint := server.HttpEndpoint()
		logger.Info("Started AltDA server", "endpoint", endpoint)

		// Export DA server URL for op-succinct proposer subprocess.
		// The Rust host reads ALTDA_SERVER_URL from env to fetch batch data.
		os.Setenv("ALTDA_SERVER_URL", endpoint)

		altDACLICfg := altda.CLIConfig{
			Enabled:      true,
			DAServerURL:  endpoint,
			VerifyOnRead: true,
			GenericDA:    true,
		}

		altDAConfig := &rollup.AltDAConfig{
			CommitmentType:    altda.GenericCommitmentString,
			DAChallengeWindow: 10,
			DAResolveWindow:   10,
		}

		// Append batcher option — consumed by WithBatcher's AfterDeploy
		orch.batcherOptions = append(orch.batcherOptions, func(id stack.L2BatcherID, cfg *bss.CLIConfig) {
			if id.ChainID() == l2ChainID {
				cfg.AltDA = altDACLICfg
			}
		})

		// Append L2CL option — consumed by WithOpNode's AfterDeploy.
		// Sets both the CLI config and the rollup config override.
		orch.l2CLOptions = append(orch.l2CLOptions, L2CLOptionFn(func(p devtest.P, id stack.L2CLNodeID, cfg *L2CLConfig) {
			if id.ChainID() == l2ChainID {
				cfg.AltDA = altDACLICfg
				cfg.RollupAltDAConfig = altDAConfig
			}
		}))
	})
}
