package l1

import (
	"github.com/hashicorp/golang-lru/v2/simplelru"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum-optimism/optimism/op-service/eth"
)

// Cache size is quite high as retrieving data from the pre-image oracle can be quite expensive
const cacheSize = 2000

// CachingOracle is an implementation of Oracle that delegates to another implementation, adding caching of all results
type CachingOracle struct {
	oracle Oracle
	blocks *simplelru.LRU[common.Hash, eth.BlockInfo]
	txs    *simplelru.LRU[eth.Bytes32, types.Transactions]
	rcpts  *simplelru.LRU[uint64, types.Receipts]
}

func NewCachingOracle(oracle Oracle) *CachingOracle {
	blockLRU, _ := simplelru.NewLRU[common.Hash, eth.BlockInfo](cacheSize, nil)
	txsLRU, _ := simplelru.NewLRU[eth.Bytes32, types.Transactions](cacheSize, nil)
	rcptsLRU, _ := simplelru.NewLRU[uint64, types.Receipts](cacheSize, nil)
	return &CachingOracle{
		oracle: oracle,
		blocks: blockLRU,
		txs:    txsLRU,
		rcpts:  rcptsLRU,
	}
}

func (o *CachingOracle) HeaderByBlockHash(blockHash common.Hash) eth.BlockInfo {
	block, ok := o.blocks.Get(blockHash)
	if ok {
		return block
	}
	block = o.oracle.HeaderByBlockHash(blockHash)
	o.blocks.Add(blockHash, block)
	return block
}

func (o *CachingOracle) TransactionsByBeaconBlockRoot(blockRoot eth.Bytes32) types.Transactions {
	txs, ok := o.txs.Get(blockRoot)
	if ok {
		return txs
	}
	txs = o.oracle.TransactionsByBeaconBlockRoot(blockRoot)
	o.txs.Add(blockRoot, txs)
	return txs
}

func (o *CachingOracle) ReceiptsByBlockNum(blockNumber uint64) types.Receipts {
	rcpts, ok := o.rcpts.Get(blockNumber)
	if ok {
		return rcpts
	}
	rcpts := o.oracle.ReceiptsByBlockNum(blockNumber)
	o.rcpts.Add(blockNumber, rcpts)
	return rcpts
}
