package actions

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum-optimism/optimism/op-bindings/predeploys"
	"github.com/ethereum-optimism/optimism/op-chain-ops/genesis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum/go-ethereum/log"
)

func TestEcotoneNetworkUpgradeTransactions(gt *testing.T) {
	t := NewDefaultTesting(gt)
	dp := e2eutils.MakeDeployParams(t, defaultRollupTestParams)
	offset := hexutil.Uint64(0)
	ecotoneOffset := hexutil.Uint64(2)
	dp.DeployConfig.L2GenesisCanyonTimeOffset = &offset
	dp.DeployConfig.L2GenesisDeltaTimeOffset = &offset
	dp.DeployConfig.L2GenesisEcotoneTimeOffset = &ecotoneOffset

	sd := e2eutils.Setup(t, dp, defaultAlloc)
	log := testlog.Logger(t, log.LvlDebug)
	_, _, _, sequencer, engine, verifier, _, _ := setupReorgTestActors(t, dp, sd, log)

	// start op-nodes
	sequencer.ActL2PipelineFull(t)
	verifier.ActL2PipelineFull(t)

	// Get current implementations addresses (by slot) for L1Block + GasPriceOracle
	initialGasPriceOracleAddress, err := engine.EthClient().StorageAt(context.Background(), predeploys.GasPriceOracleAddr, genesis.ImplementationSlot, nil)
	require.NoError(t, err)
	initialL1BlockAddress, err := engine.EthClient().StorageAt(context.Background(), predeploys.L1BlockAddr, genesis.ImplementationSlot, nil)
	require.NoError(t, err)

	// Build to the ecotone block
	sequencer.ActBuildL2ToEcotone(t)
	block := sequencer.L2Unsafe()

	// get latest block
	latestBlock, err := engine.EthClient().BlockByNumber(context.TODO(), nil)
	require.NoError(t, err)
	require.Equal(t, block.Number, latestBlock.Number().Uint64())

	contractAddy := common.HexToAddress("0xA83144c24063155FBcBF8E51f96dbD4378e18907")
	code, err := engine.EthClient().CodeAt(context.TODO(), contractAddy, latestBlock.Number())
	require.NoError(t, err)
	fmt.Println("code", code)

	transactions := latestBlock.Transactions()

	// L1Block: setInfo + 2 deploys + 2 upgradeTo
	require.Equal(t, 5, len(transactions))

	// All transactions are successful
	for i := 1; i < 5; i++ {
		txn := transactions[i]
		receipt, err := engine.EthClient().TransactionReceipt(context.Background(), txn.Hash())
		fmt.Println("transaction", i, receipt.GasUsed)
		fmt.Println("transaction", i, receipt.CumulativeGasUsed)
		require.NoError(t, err)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	}

	l1BlockDeployerAddress := common.HexToAddress("0x4210000000000000000000000000000000000000")
	gasPriceOracleDeployerAddress := common.HexToAddress("0x4210000000000000000000000000000000000001")
	expectedL1BlockAddress := crypto.CreateAddress(l1BlockDeployerAddress, 0)
	expectedGasPriceOracleAddress := crypto.CreateAddress(gasPriceOracleDeployerAddress, 0)

	// Gas Price Oracle Proxy is updated
	updatedGasPriceOracleAddress, err := engine.EthClient().StorageAt(context.Background(), predeploys.GasPriceOracleAddr, genesis.ImplementationSlot, latestBlock.Number())
	require.NoError(t, err)
	assert.Equal(t, expectedGasPriceOracleAddress, common.BytesToAddress(updatedGasPriceOracleAddress))
	assert.NotEqualf(t, initialGasPriceOracleAddress, updatedGasPriceOracleAddress, "Gas Price Oracle Proxy address should have changed")

	// L1Block Proxy is updated
	updatedL1BlockAddress, err := engine.EthClient().StorageAt(context.Background(), predeploys.L1BlockAddr, genesis.ImplementationSlot, latestBlock.Number())
	require.NoError(t, err)
	assert.Equal(t, expectedL1BlockAddress, common.BytesToAddress(updatedL1BlockAddress))
	assert.NotEqualf(t, initialL1BlockAddress, updatedL1BlockAddress, "L1Block Proxy address should have changed")

	// Get gas price from oracle
	gasPriceOracle, err := bindings.NewGasPriceOracleCaller(predeploys.GasPriceOracleAddr, engine.EthClient())
	require.NoError(t, err)
	cost, err := gasPriceOracle.GetL1Fee(nil, []byte{0, 1, 2, 3, 4})
	require.NoError(t, err)
	assert.Greater(t, cost.Uint64(), uint64(0))

	// Get L1Block info
	l1Block, err := bindings.NewL1BlockCaller(predeploys.L1BlockAddr, engine.EthClient())
	require.NoError(t, err)
	l1BlockInfo, err := l1Block.Timestamp(nil)
	require.NoError(t, err)
	assert.Greater(t, l1BlockInfo, uint64(0))
}
