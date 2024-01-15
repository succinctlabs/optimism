package l1

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

var (
	ErrNotFound     = ethereum.NotFound
	ErrUnknownLabel = errors.New("unknown label")
)

type OracleL1Client struct {
	logger               log.Logger
	oracle               Oracle
	head                 eth.L1BlockRef
	hashByNum            map[uint64]common.Hash
	earliestIndexedBlock eth.L1BlockRef
}

func NewOracleL1Client(logger log.Logger, oracle Oracle, l1Head common.Hash) *OracleL1Client {
	head := eth.InfoToL1BlockRef(oracle.HeaderByBlockHash(l1Head))
	logger.Info("L1 head loaded", "hash", head.Hash, "number", head.Number)
	return &OracleL1Client{
		logger:               logger,
		oracle:               oracle,
		head:                 head,
		hashByNum:            map[uint64]common.Hash{head.Number: head.Hash},
		earliestIndexedBlock: head,
	}
}

func (o *OracleL1Client) L1BlockRefByLabel(ctx context.Context, label eth.BlockLabel) (eth.L1BlockRef, error) {
	if label != eth.Unsafe && label != eth.Safe && label != eth.Finalized {
		return eth.L1BlockRef{}, fmt.Errorf("%w: %s", ErrUnknownLabel, label)
	}
	// The L1 head is pre-agreed and unchanging so it can be used for all of unsafe, safe and finalized
	return o.head, nil
}

func (o *OracleL1Client) L1BlockRefByNumber(ctx context.Context, number uint64) (eth.L1BlockRef, error) {
	if number > o.head.Number {
		return eth.L1BlockRef{}, fmt.Errorf("%w: block number %d", ErrNotFound, number)
	}
	hash, ok := o.hashByNum[number]
	if ok {
		return o.L1BlockRefByHash(ctx, hash)
	}
	block := o.earliestIndexedBlock
	o.logger.Info("Extending block by number lookup", "from", block.Number, "to", number)
	for block.Number > number {
		block = eth.InfoToL1BlockRef(o.oracle.HeaderByBlockHash(block.ParentHash))
		o.hashByNum[block.Number] = block.Hash
		o.earliestIndexedBlock = block
	}
	return block, nil
}

func (o *OracleL1Client) L1BlockRefByHash(ctx context.Context, hash common.Hash) (eth.L1BlockRef, error) {
	return eth.InfoToL1BlockRef(o.oracle.HeaderByBlockHash(hash)), nil
}

func (o *OracleL1Client) InfoByHash(ctx context.Context, hash common.Hash) (eth.BlockInfo, error) {
	return o.oracle.HeaderByBlockHash(hash), nil
}

// FetchReceipts fetches receipts of the given L1 block. The block must be canonical.
func (o *OracleL1Client) FetchReceipts(ctx context.Context, blockHash common.Hash) (eth.BlockInfo, types.Receipts, error) {
	info, err := o.InfoByHash(ctx, blockHash)
	if err != nil {
		return nil, nil, err
	}
	receipts := o.oracle.ReceiptsByBlockNum(info.NumberU64())
	expectedRoot := info.ReceiptHash()
	computedRoot := types.DeriveSha(receipts, trie.NewStackTrie(nil))
	if expectedRoot != computedRoot {
		panic("loaded receipts don't match expected receipts")
	}
	return info, receipts, nil
}

type DencunL1Info interface {
	eth.BlockInfo
	ParentBeaconRoot() eth.Bytes32
}

type ExtendedInfo interface {
	TransactionsHash() common.Hash // TODO we should add the MPT txs root attribute to the block-info so we can sanity-check tx lists
}

func (o *OracleL1Client) InfoAndTxsByHash(ctx context.Context, hash common.Hash) (eth.BlockInfo, types.Transactions, error) {
	info := o.oracle.HeaderByBlockHash(hash)
	// We need the next block, to access the parent-beacon-block-root
	next, err := o.L1BlockRefByNumber(ctx, info.NumberU64()+1)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get ")
	}
	if hash != next.ParentHash {
		panic("cannot retrieve non-canonical txs")
	}
	nextInfo, err := o.InfoByHash(ctx, next.Hash)
	beaconRoot := nextInfo.(DencunL1Info).ParentBeaconRoot() // TODO: need Dencun branch merged for this to actually eb there
	txs := o.oracle.TransactionsByBeaconBlockRoot(beaconRoot)
	// sanity check the txs are what we expected
	expectedRoot := info.(ExtendedInfo).TransactionsHash()
	computedRoot := types.DeriveSha(txs, trie.NewStackTrie(nil))
	if expectedRoot != computedRoot {
		panic("loaded txs don't match expected txs")
	}
	return info, txs, nil
}
