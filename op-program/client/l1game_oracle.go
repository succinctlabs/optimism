package client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	"github.com/ethereum-optimism/optimism/op-program/client/ssz"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

type L1PreimageOracle struct {
	oracle *preimage.OracleClient
	hint   *preimage.HintWriter
	fn     ssz.PairPreimageFn
}

func NewL1PreimageOracle(pClient *preimage.OracleClient, hClient *preimage.HintWriter) *L1PreimageOracle {
	return &L1PreimageOracle{oracle: pClient, hint: hClient, fn: ssz.OracleSSZ(pClient)}
}

func (o *L1PreimageOracle) header(blockHash common.Hash) *types.Header {
	o.hint.Hint(BlockHeaderHint(blockHash))
	headerRlp := o.oracle.Get(preimage.Keccak256Key(blockHash))
	var header types.Header
	if err := rlp.DecodeBytes(headerRlp, &header); err != nil {
		panic(fmt.Errorf("invalid block header %s: %w", blockHash, err))
	}
	return &header
}

func (p *L1PreimageOracle) TransactionsByBeaconBlockRoot(blockRoot eth.Bytes32) types.Transactions {
	p.hint.Hint(BeaconTransactionsHint(blockRoot))

	// Note: every EL block is embedded in a CL beacon block.
	// Skipped beacon-chain slots do not have a block, and do not have beacon-block roots embedded into the EL.

	// TODO traverse sha2 tree with sha2 preimage calls
	_ = p.oracle.Get(preimage.Sha256Key(blockRoot))

	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#beaconblockheader
	// BeaconBlockHeader: (5 fields)
	//    0: slot: Slot
	//    1: proposer_index: ValidatorIndex
	//    2: parent_root: Root
	//    3: state_root: Root
	//    4: body_root: Root
	// repeat to get binary tree of depth 3, with above 5 leaf values

	// TODO we may have to adjust remaining work depending on future L1 upgrades, and then need to load the slot

	bodyRoot := p.traverseBranchSSZ(blockRoot, 3, 4)

	// Now traverse the beacon-block body, to get to the root of the execution payload
	// Note: the BeaconBlockBody is growing, we may have to extend the depth after a certain L1 future upgrade
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#beaconblockbody
	// BeaconBlockBody: (11 fields as of Deneb)
	//    0: randao_reveal: BLSSignature
	//    1: eth1_data: Eth1Data  # Eth1 data vote
	//    2: graffiti: Bytes32  # Arbitrary data
	//    # Operations
	//    3: proposer_slashings: List[ProposerSlashing, MAX_PROPOSER_SLASHINGS]
	//    4: attester_slashings: List[AttesterSlashing, MAX_ATTESTER_SLASHINGS]
	//    5: attestations: List[Attestation, MAX_ATTESTATIONS]
	//    6: deposits: List[Deposit, MAX_DEPOSITS]
	//    7: voluntary_exits: List[SignedVoluntaryExit, MAX_VOLUNTARY_EXITS]
	//    8: sync_aggregate: SyncAggregate
	//    # Execution
	//    9: execution_payload: ExecutionPayload
	//    # Capella operations
	//    A: bls_to_execution_changes: List[...]

	executionPayloadRoot := p.fn.TraverseBranch(bodyRoot, 4, 0xA)

	// Now traverse the execution-payload (header is just a summary), to get to the root of the transactions list
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#executionpayloadheader
	// ExecutionPayloadHeader: (15 fields as of Deneb)
	//    # Execution block header fields
	//    0: parent_hash: Hash32
	//    1: fee_recipient: ExecutionAddress
	//    2: state_root: Bytes32
	//    3: receipts_root: Bytes32
	//    4: logs_bloom: ByteVector[BYTES_PER_LOGS_BLOOM]
	//    5: prev_randao: Bytes32
	//    6: block_number: uint64
	//    7: gas_limit: uint64
	//    8: gas_used: uint64
	//    9: timestamp: uint64
	//    A: extra_data: ByteList[MAX_EXTRA_DATA_BYTES]
	//    B: base_fee_per_gas: uint256
	//    # Extra payload fields
	//    C: block_hash: Hash32  # Hash of execution block
	//    D: transactions_root: Root
	//    E: withdrawals_root: Root  # [New in Capella]

	transactionsRoot := p.fn.TraverseBranch(executionPayloadRoot, 4, 0xD)

	// expand the list into a list of SSZ-roots of opaque txs
	// transactions: List[Transaction, MAX_TRANSACTIONS_PER_PAYLOAD]
	// MAX_TRANSACTIONS_PER_PAYLOAD == 2**20
	txRoots := ssz.ListLeaves[eth.Bytes32](p.fn, transactionsRoot, 20)

	// now load each of the txs
	opaqueTxs := make([]hexutil.Bytes, len(txRoots))
	for i, txRoot := range txRoots {
		// Transaction: ByteList[MAX_BYTES_PER_TRANSACTION]
		// MAX_BYTES_PER_TRANSACTION == 2**30
		opaqueTxs[i] = ssz.BytesList(p.fn, txRoot, 30)
	}

	txs, err := eth.DecodeTransactions(opaqueTxs)
	if err != nil {
		panic(fmt.Errorf("failed to decode list of txs: %w", err))
	}

	return txs
}
