package l1

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

type Oracle interface {
	// HeaderByBlockHash retrieves the block header with the given hash.
	HeaderByBlockHash(blockHash common.Hash) eth.BlockInfo

	// TransactionsByBeaconBlockRoot retrieves the transactions from the beacon block with the given block-root.
	TransactionsByBeaconBlockRoot(blockRoot eth.Bytes32) types.Transactions

	// ReceiptsByBlockNum retrieves the receipts from the canonical block with the given block number
	ReceiptsByBlockNum(blockNumber uint64) types.Receipts
}

// PreimageOracle implements Oracle using by interfacing with the pure preimage.Oracle
// to fetch pre-images to decode into the requested data.
type PreimageOracle struct {
	oracle preimage.Oracle
	hint   preimage.Hinter
}

var _ Oracle = (*PreimageOracle)(nil)

func NewPreimageOracle(raw preimage.Oracle, hint preimage.Hinter) *PreimageOracle {
	return &PreimageOracle{
		oracle: raw,
		hint:   hint,
	}
}

func (p *PreimageOracle) headerByBlockHash(blockHash common.Hash) *types.Header {
	p.hint.Hint(BlockHeaderHint(blockHash))
	headerRlp := p.oracle.Get(preimage.Keccak256Key(blockHash))
	var header types.Header
	if err := rlp.DecodeBytes(headerRlp, &header); err != nil {
		panic(fmt.Errorf("invalid block header %s: %w", blockHash, err))
	}
	return &header
}

func (p *PreimageOracle) HeaderByBlockHash(blockHash common.Hash) eth.BlockInfo {
	return eth.HeaderBlockInfo(p.headerByBlockHash(blockHash))
}

func (p *PreimageOracle) ReceiptsByBlockNum(blockNumber uint64) types.Receipts {
	// TODO traverse L1 source-truth root to find receipts data of given block number
	// TODO need to hint requirement of L1 output by block number first
	l1OutputRoot := eth.Bytes32{}

	_ = l1OutputRoot // TODO get below properties, through the l1 Output root

	l1BlockHash := common.Hash{}

	txHashesRoot := eth.Bytes32{}
	receiptsRoot := eth.Bytes32{}

	p.hint.Hint(ReceiptsHint(l1BlockHash))

	// TODO assumes we put the tx hashes in a nice SSZ list
	txHashes := listLeaves[common.Hash](p, txHashesRoot, 20)

	// TODO assumes we format it like an SSZ list, just like the txs
	// number of receipts matches number of txs, so we use the same capacity as used for the txs list
	receiptRoots := listLeaves[eth.Bytes32](p, receiptsRoot, 20)

	// now load each of the receipts
	opaqueReceipts := make([]hexutil.Bytes, len(receiptRoots))
	for i, recRoot := range receiptRoots {
		// TODO: is 2**30 enough capacity for a raw receipt?
		opaqueReceipts[i] = p.bytesList(recRoot, 30)
	}

	id := eth.BlockID{Hash: l1BlockHash, Number: blockNumber}
	receipts, err := eth.DecodeRawReceipts(id, opaqueReceipts, txHashes)
	if err != nil {
		panic(fmt.Errorf("bad receipts data for block %s: %w", l1BlockHash, err))
	}

	return receipts
}
