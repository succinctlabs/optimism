package altda

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum-optimism/optimism/op-devstack/devtest"
	"github.com/ethereum-optimism/optimism/op-devstack/dsl"
	"github.com/ethereum-optimism/optimism/op-devstack/presets"
	"github.com/ethereum-optimism/optimism/op-devstack/stack/match"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-supervisor/supervisor/types"
)

// TestAltDA_SafeHeadProgresses verifies the full AltDA pipeline:
//  1. Batcher posts batch data to the DA server
//  2. Batcher posts only the commitment to L1
//  3. Op-node reads commitment from L1
//  4. Op-node fetches batch data from DA server using the commitment
//  5. Op-node derives the batch and advances the safe head
//
// If the safe head advances after an L2 transaction, the entire AltDA
// derivation pipeline is proven to work end-to-end.
func TestAltDA_SafeHeadProgresses(gt *testing.T) {
	t := devtest.SerialT(gt)
	sys := presets.NewMinimal(t)
	l := t.Logger()

	// Record initial safe head
	initialSafe := sys.L2CL.HeadBlockRef(types.LocalSafe)
	l.Info("Initial safe head", "number", initialSafe.Number, "hash", initialSafe.Hash)

	// Send an L2 transaction
	alice := sys.FunderL2.NewFundedEOA(eth.OneTenthEther)
	bob := sys.Wallet.NewEOA(sys.L2EL)
	alice.Transfer(bob.Address(), eth.OneHundredthEther)

	// Wait for safe head to advance — this is the critical assertion.
	// Safe head can only advance if the batcher submitted a commitment,
	// and op-node resolved it via the DA server.
	dsl.CheckAll(t, sys.L2CL.AdvancedFn(types.LocalSafe, 1, 30))

	// Log final sync status as evidence
	status := sys.L2CL.SyncStatus()
	l.Info("AltDA chain progressed successfully",
		"unsafeL2", status.UnsafeL2.Number,
		"safeL2", status.SafeL2.Number,
		"localSafeL2", status.LocalSafeL2.Number,
	)

	// Verify L1 batch transactions contain commitments (start with 0x01),
	// not raw batch data (which starts with 0x00).
	verifyL1CommitmentData(t, sys)
}

// verifyL1CommitmentData scans recent L1 blocks for batcher transactions
// and verifies they contain AltDA commitment data (version byte 0x01).
func verifyL1CommitmentData(t devtest.T, sys *presets.Minimal) {
	l := t.Logger()
	rollupCfg := sys.L2Chain.Escape().RollupConfig()
	batchInbox := rollupCfg.BatchInboxAddress

	l1EC := sys.L1Network.Escape().L1ELNode(match.FirstL1EL).EthClient()

	ctx, cancel := context.WithTimeout(t.Ctx(), dsl.DefaultTimeout)
	defer cancel()

	head, err := l1EC.BlockRefByLabel(ctx, eth.Unsafe)
	require.NoError(t, err)

	// Scan backwards through recent L1 blocks looking for batcher txs
	scanFloor := uint64(0)
	if head.Number > 20 {
		scanFloor = head.Number - 20
	}
	foundCommitment := false
	for blockNum := head.Number; blockNum > scanFloor; blockNum-- {
		_, txs, err := l1EC.InfoAndTxsByNumber(ctx, blockNum)
		require.NoError(t, err)

		for _, tx := range txs {
			if tx.To() == nil || *tx.To() != batchInbox {
				continue
			}
			data := tx.Data()
			if len(data) == 0 {
				continue
			}
			// In AltDA mode, batcher prefixes commitment data with DerivationVersion1 (0x01).
			// Regular calldata batches use DerivationVersion0 (0x00).
			version := data[0]
			l.Info("Found batcher tx on L1",
				"block", blockNum,
				"txHash", tx.Hash(),
				"dataLen", len(data),
				"versionByte", version,
			)
			require.Equal(t, byte(0x01), version,
				"batcher tx should contain AltDA commitment (version 0x01), got version 0x%02x", version)
			foundCommitment = true
		}
	}
	require.True(t, foundCommitment, "should find at least one batcher commitment tx on L1")
	l.Info("L1 commitment data verified — batcher posted commitments, not raw batch data")
}
