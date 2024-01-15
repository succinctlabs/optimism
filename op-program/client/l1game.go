package client

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum-optimism/optimism/op-program/client/l2"
	"github.com/ethereum-optimism/optimism/op-program/client/l2/engineapi"
	"github.com/ethereum-optimism/optimism/op-service/eth"

	preimage "github.com/ethereum-optimism/optimism/op-preimage"
)

// L1Game represents a state-transition proof of an "L1-output".
// We run the L1 state-transition because L1 is careless to merkleize all outputs of interest to L2s.
// But the inputs are accessible, so we process those, and establish all the details.
// The L1-output can be efficiently created off-chain.
// This prevents the "large preimage problem" by generating well-merkleized L1 receipts within the VM,
// rather than having to load 4MB pre-images into the VM.
type L1Game struct {
	L1Head common.Hash

	// There is no "prestate", since the whole chain is already implied by the L1Head.
	// The "post" that we dispute is an out-of-band interpretation of processing that L1 data our way.

	L1Claim       eth.Bytes32
	L1ClaimNumber uint64

	L1ChainConfig *params.ChainConfig
}

func (l1Game *L1Game) Run(logger log.Logger, pClient *preimage.OracleClient, hClient *preimage.HintWriter) error {
	out, err := l1Game.produceOutput(logger, pClient, hClient)
	if err != nil {
		// exit with panic, this computation should always succeed, if the iputs are valid.
		panic(fmt.Errorf("L1 computation failed: %w", err))
	}
	if out != l1Game.L1Claim {
		return fmt.Errorf("claim %s is invalid, computed %s for L1 block %d", l1Game.L1Claim, out, l1Game.L1ClaimNumber)
	}
	return nil
}

func (l1Game *L1Game) produceOutput(logger log.Logger, pClient *preimage.OracleClient, hClient *preimage.HintWriter) (eth.Bytes32, error) {
	// TODO create L1 exec engine, state etc.

	l1Oracle := NewL1PreimageOracle(pClient, hClient)

	targetBlockHash := common.Hash{} // TODO

	// walk back L1 chain to get L1 block by number
	hdr := l1Oracle.header(l1Game.L1Head)
	var beaconBlockRoot eth.Bytes32
	for hdr.Number.Uint64() > l1Game.L1ClaimNumber {
		if hdr.Number.Uint64() == l1Game.L1ClaimNumber+1 {
			if hdr.ParentBeaconRoot == nil {
				return eth.Bytes32{}, fmt.Errorf("required parent beacon root, block %d", hdr.Number)
			}
			beaconBlockRoot = eth.Bytes32(*hdr.ParentBeaconRoot)
		}
		hdr = l1Oracle.header(hdr.ParentHash)
	}

	// TODO caching oracles
	oracle := l2.NewPreimageOracle(pClient, hClient)

	preHeader := l1Oracle.header(hdr.ParentHash)
	backend, err := l2.NewOracleBackedChain(logger, oracle, l1Game.L1ChainConfig, preHeader)
	if err != nil {
		return eth.Bytes32{}, fmt.Errorf("failed to build oracle-backed chain: %w", err)
	}
	processor, err := engineapi.NewBlockProcessorFromHeader(backend, preHeader)
	if err != nil {
		return eth.Bytes32{}, fmt.Errorf("failed to make block processor: %w", err)
	}

	// Load the txs through the beacon-root:
	// the txs are merkleized with SSZ hash-tree-root into a nice binary merkle-tree,
	// allowing us to load them small parts at a time,
	// to make the preimage-loading step in the L1 EVM possible within limited gas.
	txs := l1Oracle.TransactionsByBeaconBlockRoot(beaconBlockRoot)
	for i, tx := range txs {
		if err := processor.AddTx(tx); err != nil {
			return eth.Bytes32{}, fmt.Errorf("failed to process tx %d %s: %w", i, tx.Hash(), err)
		}
	}

	block, receipts, err := processor.Assemble()
	if err != nil {
		return eth.Bytes32{}, fmt.Errorf("failed to assembly block: %w", err)
	}

	// Sanity check our computation. Exit with panic if we failed; game is invalid.
	// E.g. the chain config or EVM code may be outdated.
	if got := block.Hash(); got != targetBlockHash {
		return eth.Bytes32{}, fmt.Errorf("expected to correctly re-compute the block, got %s but expeced %s", got, targetBlockHash)
	}

	// just merkleize the tx-hashes nicely now, so they won't require MPT traversal in the L2 client.
	txHashes := make([]common.Hash, len(block.Transactions()))
	txRoots := make([]eth.Bytes32, len(block.Transactions()))
	for i, tx := range block.Transactions() {
		txHashes[i] = tx.Hash()

		data, err := tx.MarshalBinary()
		if err != nil {
			return eth.Bytes32{}, fmt.Errorf("failed to encode tx: %w", err)
		}
		txRoots[i] = sszByteListHTR(data, 30)
	}

	// merkleize the receipts
	receiptRoots := make([]eth.Bytes32, len(receipts))
	for i, r := range receipts {
		data, err := r.MarshalBinary()
		if err != nil {
			return eth.Bytes32{}, fmt.Errorf("failed to encode receipt: %w", err)
		}
		receiptRoots[i] = sszByteListHTR(data, 30)
	}

	// TODO: merkleize the above data into an "L1 output",
	// such that the preimage-client of the L2 game can load the L1 data.
	root := eth.Bytes32{}

	return root, nil
}

func sszByteListHTR(data []byte, depth uint8) eth.Bytes32 {
	// TODO merkleize the data, with length mixin etc. like SSZ
	return eth.Bytes32{}
}
