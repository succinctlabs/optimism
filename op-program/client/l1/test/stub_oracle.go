package test

import (
	"testing"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type StubOracle struct {
	t *testing.T

	// Blocks maps block hash to eth.BlockInfo
	Blocks map[common.Hash]eth.BlockInfo

	// Txs maps beacon block root to transactions
	Txs map[eth.Bytes32]types.Transactions

	// Rcpts maps Block hash to receipts
	Rcpts map[common.Hash]types.Receipts
}

func NewStubOracle(t *testing.T) *StubOracle {
	return &StubOracle{
		t:      t,
		Blocks: make(map[common.Hash]eth.BlockInfo),
		Txs:    make(map[eth.Bytes32]types.Transactions),
		Rcpts:  make(map[common.Hash]types.Receipts),
	}
}
func (o StubOracle) HeaderByBlockHash(blockHash common.Hash) eth.BlockInfo {
	info, ok := o.Blocks[blockHash]
	if !ok {
		o.t.Fatalf("unknown block %s", blockHash)
	}
	return info
}

func (o StubOracle) TransactionsByBeaconBlockRoot(blockRoot eth.Bytes32) types.Transactions {
	txs, ok := o.Txs[blockRoot]
	if !ok {
		o.t.Fatalf("unknown txs %s", blockRoot)
	}
	return txs
}

func (o StubOracle) ReceiptsByBlockHash(blockHash common.Hash) (eth.BlockInfo, types.Receipts) {
	rcpts, ok := o.Rcpts[blockHash]
	if !ok {
		o.t.Fatalf("unknown rcpts %s", blockHash)
	}
	return o.HeaderByBlockHash(blockHash), rcpts
}
