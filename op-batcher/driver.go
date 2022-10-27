package op_batcher

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"io"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum-optimism/optimism/op-batcher/sequencer"
	"github.com/ethereum-optimism/optimism/op-node/eth"
	"github.com/ethereum-optimism/optimism/op-proposer/txmgr"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

// BatchSubmitter encapsulates a service responsible for submitting L2 tx
// batches to L1 for availability.
type BatchSubmitter struct {
	txMgr txmgr.TxManager
	addr  common.Address
	cfg   sequencer.Config
	wg    sync.WaitGroup
	done  chan struct{}
	log   log.Logger

	ctx    context.Context
	cancel context.CancelFunc

	lastSubmittedBlock eth.BlockID

	s channelManager
}

// NewBatchSubmitter initializes the BatchSubmitter, gathering any resources
// that will be needed during operation.
func NewBatchSubmitter(cfg Config, l log.Logger) (*BatchSubmitter, error) {
	ctx := context.Background()

	var err error
	var sequencerPrivKey *ecdsa.PrivateKey
	var addr common.Address

	if cfg.PrivateKey != "" && cfg.Mnemonic != "" {
		return nil, errors.New("cannot specify both a private key and a mnemonic")
	}

	if cfg.PrivateKey == "" {
		// Parse wallet private key that will be used to submit L2 txs to the batch
		// inbox address.
		wallet, err := hdwallet.NewFromMnemonic(cfg.Mnemonic)
		if err != nil {
			return nil, err
		}

		acc := accounts.Account{
			URL: accounts.URL{
				Path: cfg.SequencerHDPath,
			},
		}
		addr, err = wallet.Address(acc)
		if err != nil {
			return nil, err
		}

		sequencerPrivKey, err = wallet.PrivateKey(acc)
		if err != nil {
			return nil, err
		}
	} else {
		sequencerPrivKey, err = crypto.HexToECDSA(strings.TrimPrefix(cfg.PrivateKey, "0x"))
		if err != nil {
			return nil, err
		}

		addr = crypto.PubkeyToAddress(sequencerPrivKey.PublicKey)
	}

	batchInboxAddress, err := parseAddress(cfg.SequencerBatchInboxAddress)
	if err != nil {
		return nil, err
	}

	// Connect to L1 and L2 providers. Perform these last since they are the
	// most expensive.
	l1Client, err := dialEthClientWithTimeout(ctx, cfg.L1EthRpc)
	if err != nil {
		return nil, err
	}

	l2Client, err := dialEthClientWithTimeout(ctx, cfg.L2EthRpc)
	if err != nil {
		return nil, err
	}

	rollupClient, err := dialRollupClientWithTimeout(ctx, cfg.RollupRpc)
	if err != nil {
		return nil, err
	}

	chainID, err := l1Client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	sequencerBalance, err := l1Client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, err
	}

	log.Info("starting batch submitter", "submitter_addr", addr, "submitter_bal", sequencerBalance)

	txManagerConfig := txmgr.Config{
		Log:                       l,
		Name:                      "Batch Submitter",
		ResubmissionTimeout:       cfg.ResubmissionTimeout,
		ReceiptQueryInterval:      time.Second,
		NumConfirmations:          cfg.NumConfirmations,
		SafeAbortNonceTooLowCount: cfg.SafeAbortNonceTooLowCount,
	}

	batcherCfg := sequencer.Config{
		Log:               l,
		Name:              "Batch Submitter",
		L1Client:          l1Client,
		L2Client:          l2Client,
		RollupNode:        rollupClient,
		MinL1TxSize:       cfg.MinL1TxSize,
		MaxL1TxSize:       cfg.MaxL1TxSize,
		BatchInboxAddress: batchInboxAddress,
		ChannelTimeout:    cfg.ChannelTimeout,
		ChainID:           chainID,
		PrivKey:           sequencerPrivKey,
		PollInterval:      cfg.PollInterval,
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &BatchSubmitter{
		cfg:   batcherCfg,
		addr:  addr,
		txMgr: txmgr.NewSimpleTxManager("batcher", txManagerConfig, l1Client),
		done:  make(chan struct{}),
		log:   l,
		// TODO: this context only exists because the even loop doesn't reach done
		// if the tx manager is blocking forever due to e.g. insufficient balance.
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (l *BatchSubmitter) Start() error {
	l.wg.Add(1)
	go l.loop()
	return nil
}

func (l *BatchSubmitter) Stop() {
	l.cancel()
	close(l.done)
	l.wg.Wait()
}

// loadBlocksIntoState loads all blocks since the previous
func (l *BatchSubmitter) loadBlocksIntoState(ctx context.Context) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	defer cancel()
	syncStatus, err := l.cfg.RollupNode.SyncStatus(ctx)
	if err != nil {
		l.log.Warn("issue fetching L2 head", "err", err)
		return
	}
	if syncStatus.HeadL1 == (eth.L1BlockRef{}) {
		l.log.Info("Rollup node has no L1 head info yet")
		return
	}
	l.log.Info("Got new L2 sync status", "safe_head", syncStatus.SafeL2, "unsafe_head", syncStatus.UnsafeL2, "last_submitted", l.lastSubmittedBlock, "l1_head", syncStatus.HeadL1)
	if syncStatus.SafeL2.Number >= syncStatus.UnsafeL2.Number {
		l.log.Trace("No unsubmitted blocks from sequencer")
		return
	}

	// If we just started, start at safe-head
	if l.lastSubmittedBlock == (eth.BlockID{}) {
		l.log.Info("Starting batch-submitter work at safe-head", "safe", syncStatus.SafeL2)
		l.lastSubmittedBlock = syncStatus.SafeL2.ID()
	}
	// If it's lagging behind, catch it up.
	if l.lastSubmittedBlock.Number < syncStatus.SafeL2.Number {
		l.log.Warn("last submitted block lagged behind L2 safe head: batch submission will continue from the safe head now", "last", l.lastSubmittedBlock, "safe", syncStatus.SafeL2)
		l.lastSubmittedBlock = syncStatus.SafeL2.ID()
	}

	prevID := l.lastSubmittedBlock

	// Add all blocks to "state"
	for i := l.lastSubmittedBlock.Number + 1; i <= syncStatus.UnsafeL2.Number; i++ {
		ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
		block, err := l.cfg.L2Client.BlockByNumber(ctx, new(big.Int).SetUint64(i))
		cancel()
		if err != nil {
			l.log.Error("issue fetching L2 block", "err", err)
			return
		}
		if block.ParentHash() != prevID.Hash {
			l.log.Error("detected a reorg in L2 chain vs previous submitted information, resetting to safe head now", "safe_head", syncStatus.SafeL2)
			l.lastSubmittedBlock = syncStatus.SafeL2.ID()
			return
		}
		if err := l.s.AddL2Block(block); err != nil {
			return
		}
		prevID = eth.ToBlockID(block)
		l.lastSubmittedBlock = prevID
		// l.log.Info("added L2 block to channel", "block", prevID, "channel_id", l.ch.ID(), "tx_count", len(block.Transactions()), "time", block.Time())
	}
}

func (l *BatchSubmitter) submitTransaction(data []byte) error {
	// Query for the submitter's current nonce.
	walletAddr := crypto.PubkeyToAddress(l.cfg.PrivKey.PublicKey)
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*10)
	nonce, err := l.cfg.L1Client.NonceAt(ctx, walletAddr, nil)
	cancel()
	if err != nil {
		l.log.Error("unable to get current nonce", "err", err)
		return err
	}

	// Create the transaction
	ctx, cancel = context.WithTimeout(l.ctx, time.Second*10)
	tx, err := l.CraftTx(ctx, data, nonce)
	cancel()
	if err != nil {
		l.log.Error("unable to craft tx", "err", err)
		return err
	}

	// Construct the a closure that will update the txn with the current gas prices.
	updateGasPrice := func(ctx context.Context) (*types.Transaction, error) {
		l.log.Debug("updating batch tx gas price")
		return l.UpdateGasPrice(ctx, tx)
	}

	// Wait until one of our submitted transactions confirms. If no
	// receipt is received it's likely our gas price was too low.
	// TODO: does the tx manager nicely replace the tx?
	//  (submit a new one, that's within the channel timeout, but higher fee than previously submitted tx? Or use a cheap cancel tx?)
	ctx, cancel = context.WithTimeout(l.ctx, time.Second*time.Duration(l.cfg.ChannelTimeout))
	receipt, err := l.txMgr.Send(ctx, updateGasPrice, l.cfg.L1Client.SendTransaction)
	cancel()
	if err != nil {
		l.log.Warn("unable to publish tx", "err", err)
		return err
	}

	// The transaction was successfully submitted.
	l.log.Info("tx successfully published", "tx_hash", receipt.TxHash)
	return nil
}

func (l *BatchSubmitter) loop() {
	defer l.wg.Done()

	ticker := time.NewTicker(l.cfg.PollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			l.loadBlocksIntoState(context.TODO())
			data, _, err := l.s.TxData(eth.L1BlockRef{})
			if err == io.EOF {
				break
			} else if err != nil {
				l.log.Error("unable to get tx data", "err", err)
				break
			}
			_ = l.submitTransaction(data)

		case <-l.done:
			return
		}
	}
}
