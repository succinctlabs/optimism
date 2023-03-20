package actions

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum-optimism/optimism/op-batcher/batcher"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils"
	"github.com/ethereum-optimism/optimism/op-node/testlog"
	"github.com/ethereum-optimism/optimism/op-service/crypto"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
)

func TestBatcherL1Reorg(gt *testing.T) {
	t := NewDefaultTesting(gt)
	t.Log("Starting test")
	b, miner, _, sequencer := setupBatcherTest(t)
	l1Client := miner.EthClient()

	// Prep the sequencer
	sequencer.ActL1HeadSignal(t)
	sequencer.ActL2PipelineFull(t)
	origSafeHead := sequencer.L2Safe()

	// Create an L2 block to be submitted
	sequencer.ActL2StartBlock(t)
	sequencer.ActL2EndBlock(t)
	sequencer.ActL2PipelineFull(t)

	// Run the batcher to submit the transaction
	progressBatcher(t, b, l1Client, miner)

	// Check the transaction got included on L1
	blockNum, err := l1Client.BlockNumber(t.Ctx())
	require.NoError(t, err)
	block, err := l1Client.BlockByNumber(t.Ctx(), big.NewInt(int64(blockNum)))
	require.NoError(t, err)
	require.Equal(t, 1, len(block.Transactions()))

	// Progress the sequencer
	sequencer.ActL1HeadSignal(t)
	sequencer.ActL2PipelineFull(t)
	safeHead := sequencer.L2Safe()
	require.Greater(t, safeHead.Number, origSafeHead.Number, "Safe head did not progress")
}

func progressBatcher(t StatefulTesting, b *batcher.BatchSubmitter, l1Client *ethclient.Client, miner *L1Miner) {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		b.Step()
	}()

	for {
		select {
		case <-ch:
			t.Log("Batcher step complete")
			return
		case <-time.After(200 * time.Millisecond):
			t.Log("Checking for pending transactions")
			// Check for pending transactions
			count, err := l1Client.PendingTransactionCount(t.Ctx())
			require.NoError(t, err)
			if count == 0 {
				continue
			}

			t.Log("Producing L1 block")
			miner.ActL1StartBlock(4)(t)
			pendingTxs := miner.eth.TxPool().Pending(false)
			for _, transactions := range pendingTxs {
				for _, tx := range transactions {
					miner.IncludeTx(t, tx)
				}
			}
			miner.ActL1EndBlock(t)
		}
	}
}

func setupBatcherTest(t StatefulTesting) (*batcher.BatchSubmitter, *L1Miner, *L2Engine, *L2Sequencer) {
	p := &e2eutils.TestParams{
		MaxSequencerDrift:   20, // larger than L1 block time we simulate in this test (12)
		SequencerWindowSize: 24,
		ChannelTimeout:      20,
	}
	dp := e2eutils.MakeDeployParams(t, p)
	sd := e2eutils.Setup(t, dp, defaultAlloc)
	l := testlog.Logger(t, log.LvlDebug)
	miner, engine, sequencer := setupSequencerTest(t, sd, l)
	miner.ActL1SetFeeRecipient(common.Address{'A'})

	rollupNode := sequencer.RollupClient()
	sign := crypto.PrivateKeySignerFn(dp.Secrets.Batcher, big.NewInt(int64(dp.DeployConfig.L1ChainID)))
	batcherCfg := batcher.Config{
		L1Client:     miner.EthClient(),
		L2Client:     engine.EthClient(),
		RollupNode:   rollupNode,
		PollInterval: 15,
		TxManagerConfig: txmgr.Config{
			ResubmissionTimeout:       10 * time.Minute,
			ReceiptQueryInterval:      time.Second,
			NumConfirmations:          1,
			SafeAbortNonceTooLowCount: 3,
			From:                      dp.Addresses.Batcher,
			Signer: func(ctx context.Context, from common.Address, tx *types.Transaction) (*types.Transaction, error) {
				return sign(from, tx)
			},
		},
		From:   dp.Addresses.Batcher,
		Rollup: sd.RollupCfg,
		Channel: batcher.ChannelConfig{
			SeqWindowSize:      15,
			ChannelTimeout:     40,
			MaxChannelDuration: 1,
			SubSafetyMargin:    4,
			// Set the max frame size to 24 so that we can test sending transactions
			// The fixed overhead for frame size is 23, so we must be larger or else
			// the uint64 will underflow, causing the frame size to be essentially unbound
			MaxFrameSize:     200,
			TargetFrameSize:  1,
			TargetNumFrames:  1,
			ApproxComprRatio: 0.4,
		},
	}
	sequencer.L2Unsafe()

	b, err := batcher.NewBatchSubmitter(t.Ctx(), batcherCfg, l)
	require.NoError(t, err)
	return b, miner, engine, sequencer
}
