package sysgo

import (
	"os"

	"github.com/ethereum/go-ethereum/common"

	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	bss "github.com/ethereum-optimism/optimism/op-batcher/batcher"
	"github.com/ethereum-optimism/optimism/op-devstack/devtest"
	"github.com/ethereum-optimism/optimism/op-devstack/stack"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

// WithAltDA starts a DA server and configures the batcher and op-node to use AltDA mode
// with Keccak256 commitments. The DA server computes keccak256(data) as the commitment,
// which the Rust AltDA data source can verify inside the zkVM proof.
//
// The DA server is started and batcher/L2CL options are accumulated during the Deploy phase.
// The rollup config is overridden via L2CLConfig.RollupAltDAConfig, which is applied
// inside op-node's AfterDeploy hook (after the deployer registers the L2 network).
// This avoids ordering dependencies with the deployer's AfterDeploy.
func WithAltDA(l2ChainID eth.ChainID) stack.Option[*Orchestrator] {
	return stack.Deploy(func(orch *Orchestrator) {
		p := orch.P()
		logger := p.Logger()

		// Start DA server with in-memory storage (Keccak256 commitment mode)
		store := altda.NewMemStore()
		server := altda.NewDAServer("127.0.0.1", 0, store, logger, false)
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
			GenericDA:    false,
		}

		// Keccak256 mode requires a non-zero DAChallengeAddress to pass rollup config validation.
		// The address is unused — challenges are never triggered in e2e tests.
		altDAConfig := &rollup.AltDAConfig{
			CommitmentType:     altda.KeccakCommitmentString,
			DAChallengeAddress: common.HexToAddress("0x0000000000000000000000000000000000000001"),
			DAChallengeWindow:  10,
			DAResolveWindow:    10,
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
