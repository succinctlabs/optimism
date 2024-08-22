package utils

import (
	"context"
	"fmt"
	"math/big"
	"os"

	// "os"
	"time"

	"github.com/ethereum-optimism/optimism/op-node/cmd/batch_decoder/fetch"
	"github.com/ethereum-optimism/optimism/op-node/cmd/batch_decoder/reassemble"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/dial"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BatchDecoderConfig is a struct that holds the configuration for the batch decoder.
type BatchDecoderConfig struct {
	L2GenesisTime     uint64
	L2GenesisBlock    uint64
	L2BlockTime       uint64
	BatchInboxAddress common.Address
	StartBlock        uint64
	EndBlock          uint64
	L2ChainID         *big.Int
	L2Node            string
	L1RPC             string
	L1Beacon          string
	BatchSender       string
	DataDir           string
}

// GetAllSpanBatchesInBlockRange fetches span batches within a range of L2 blocks.
func GetAllSpanBatchesInBlockRange(config BatchDecoderConfig) ([]reassemble.SpanBatchRange, error) {
	rollupCfg, err := setupBatchDecoderConfig(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rollupClient, err := dial.DialRollupClientWithTimeout(ctx, dial.DefaultDialTimeout, nil, config.L2Node)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rollup client: %w", err)
	}

	l1Origin, finalizedL1, err := getL1Origins(rollupClient, config.StartBlock, config.EndBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to get L1 origin and finalized: %w", err)
	}

	// Clear the out directory so that loading the transaction frames is fast. Otherwise, when loading thousands of transactions,
	// this process can become quite slow.
	err = os.RemoveAll(config.DataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to clear out directory: %w", err)
	}

	// Step on the L1 and store all batches posted in config.DataDir.
	err = fetchBatches(config, rollupCfg, l1Origin, finalizedL1)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch batches: %w", err)
	}

	// Reassemble the batches into span batches from the stored transaction frames in config.DataDir.
	reassembleConfig := reassemble.Config{
		BatchInbox:    config.BatchInboxAddress,
		InDirectory:   config.DataDir,
		OutDirectory:  "",
		L2ChainID:     config.L2ChainID,
		L2GenesisTime: config.L2GenesisTime,
		L2BlockTime:   config.L2BlockTime,
	}

	// Get all span batch ranges in the given L2 block range.
	ranges, err := reassemble.GetSpanBatchRanges(reassembleConfig, rollupCfg, config.StartBlock, config.EndBlock, 1000000)
	if err != nil {
		return nil, fmt.Errorf("failed to get span batch ranges: %w", err)
	}

	return ranges, nil
}

func setupBatchDecoderConfig(config *BatchDecoderConfig) (*rollup.Config, error) {
	rollupCfg, err := rollup.LoadOPStackRollupConfig(config.L2ChainID.Uint64())
	if err != nil {
		return nil, err
	}

	if config.L2GenesisTime != rollupCfg.Genesis.L2Time {
		config.L2GenesisTime = rollupCfg.Genesis.L2Time
		fmt.Printf("L2GenesisTime overridden: %v\n", config.L2GenesisTime)
	}
	if config.L2GenesisBlock != rollupCfg.Genesis.L2.Number {
		config.L2GenesisBlock = rollupCfg.Genesis.L2.Number
		fmt.Printf("L2GenesisBlock overridden: %v\n", config.L2GenesisBlock)
	}
	if config.L2BlockTime != rollupCfg.BlockTime {
		config.L2BlockTime = rollupCfg.BlockTime
		fmt.Printf("L2BlockTime overridden: %v\n", config.L2BlockTime)
	}
	if config.BatchInboxAddress != rollupCfg.BatchInboxAddress {
		config.BatchInboxAddress = rollupCfg.BatchInboxAddress
		fmt.Printf("BatchInboxAddress overridden: %v\n", config.BatchInboxAddress)
	}

	return rollupCfg, nil
}

// Get the L1 origin corresponding to the given L2 block and the latest finalized L1 block.
func getL1Origins(rollupClient *sources.RollupClient, startBlock, endBlock uint64) (uint64, uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	output, err := rollupClient.OutputAtBlock(ctx, startBlock)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get output at start block: %w", err)
	}
	startL1Origin := output.BlockRef.L1Origin.Number

	output, err = rollupClient.OutputAtBlock(ctx, endBlock)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get output at end block: %w", err)
	}

	// TODO: Change 12 to the L1 block time
	l1BlockTime := 12
	endL1Origin := output.BlockRef.L1Origin.Number + (uint64(60/l1BlockTime) * 10)

	return startL1Origin, endL1Origin, nil
}

func fetchBatches(config BatchDecoderConfig, rollupCfg *rollup.Config, l1Origin, finalizedL1 uint64) error {
	fetchConfig := fetch.Config{
		Start:   l1Origin,
		End:     finalizedL1,
		ChainID: rollupCfg.L1ChainID,
		BatchSenders: map[common.Address]struct{}{
			common.HexToAddress(config.BatchSender): {},
		},
		BatchInbox:         config.BatchInboxAddress,
		OutDirectory:       config.DataDir,
		ConcurrentRequests: 10,
	}

	l1Client, err := ethclient.Dial(config.L1RPC)
	if err != nil {
		return fmt.Errorf("failed to dial L1 client: %w", err)
	}

	beacon, err := setupBeacon(config)
	if err != nil {
		return err
	}

	totalValid, totalInvalid := fetch.Batches(l1Client, beacon, fetchConfig)
	fmt.Printf("Fetched batches in range [%v,%v). Found %v valid & %v invalid batches\n", fetchConfig.Start, fetchConfig.End, totalValid, totalInvalid)

	return nil
}

func setupBeacon(config BatchDecoderConfig) (*sources.L1BeaconClient, error) {
	if config.L1Beacon == "" {
		fmt.Println("L1 Beacon endpoint not set. Unable to fetch post-ecotone channel frames")
		return nil, nil
	}

	beaconClient := sources.NewBeaconHTTPClient(client.NewBasicHTTPClient(config.L1Beacon, nil))
	beaconCfg := sources.L1BeaconClientConfig{FetchAllSidecars: false}
	beacon := sources.NewL1BeaconClient(beaconClient, beaconCfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := beacon.GetVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check L1 Beacon API version: %w", err)
	}

	return beacon, nil
}
