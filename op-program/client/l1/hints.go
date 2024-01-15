package l1

import (
	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
)

const (
	HintL1BlockHeader        = "l1-block-header"
	HintL1BeaconTransactions = "l1-beacon-transactions"
	HintL1Receipts           = "l1-receipts"
)

type BlockHeaderHint common.Hash

var _ preimage.Hint = BlockHeaderHint{}

func (l BlockHeaderHint) Hint() string {
	return HintL1BlockHeader + " " + (common.Hash)(l).String()
}

type BeaconTransactionsHint eth.Bytes32

var _ preimage.Hint = BeaconTransactionsHint{}

func (l BeaconTransactionsHint) Hint() string {
	return HintL1BeaconTransactions + " " + (eth.Bytes32)(l).String()
}

type ReceiptsHint common.Hash

var _ preimage.Hint = ReceiptsHint{}

func (l ReceiptsHint) Hint() string {
	return HintL1Receipts + " " + (common.Hash)(l).String()
}
