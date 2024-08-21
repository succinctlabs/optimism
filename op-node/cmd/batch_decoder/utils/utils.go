package utils

import (
	"context"
	"fmt"
	"math/big"
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

// Get the span batches within a range of L2 blocks. First, fetch the span batches from the L1 origin of the start block to the finalized L1 block. (This can take a while)
// Then, reassemble the span batches within the range for the given L2 blocks.
func GetAllSpanBatchesInBlockRange(config BatchDecoderConfig) ([][2]uint64, error) {
	// Get the L1 origin corresponding to the start block
	// nextBlock is equal to the highest value in the `EndBlock` column of the db, plus 1
	rollupCfg, err := rollup.LoadOPStackRollupConfig(config.L2ChainID.Uint64())
	if err == nil {
		// prioritize superchain config
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
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rollupClient, err := dial.DialRollupClientWithTimeout(ctx, dial.DefaultDialTimeout, nil, config.L2Node)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rollup client: %w", err)
	}

	l1Origin, finalizedL1, err := getL1OriginAndFinalized(rollupClient, config.StartBlock, config.EndBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to get L1 origin and finalized: %w", err)
	}

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
		return nil, fmt.Errorf("failed to dial L1 client: %w", err)
	}
	var beacon *sources.L1BeaconClient
	if config.L1Beacon != "" {
		beaconClient := sources.NewBeaconHTTPClient(client.NewBasicHTTPClient(config.L1Beacon, nil))
		beaconCfg := sources.L1BeaconClientConfig{FetchAllSidecars: false}
		beacon = sources.NewL1BeaconClient(beaconClient, beaconCfg)
		_, err := beacon.GetVersion(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to check L1 Beacon API version: %w", err)
		}
	} else {
		fmt.Println("L1 Beacon endpoint not set. Unable to fetch post-ecotone channel frames")
	}
	totalValid, totalInvalid := fetch.Batches(l1Client, beacon, fetchConfig)
	fmt.Printf("Fetched batches in range [%v,%v). Found %v valid & %v invalid batches\n", fetchConfig.Start, fetchConfig.End, totalValid, totalInvalid)

	reassembleConfig := reassemble.Config{
		BatchInbox:    config.BatchInboxAddress,
		InDirectory:   config.DataDir,
		OutDirectory:  "",
		L2ChainID:     config.L2ChainID,
		L2GenesisTime: config.L2GenesisTime,
		L2BlockTime:   config.L2BlockTime,
	}

	ranges, err := GetSpanBatchRanges(reassembleConfig, rollupCfg, config.StartBlock, config.EndBlock, 1000000)
	if err != nil {
		return nil, fmt.Errorf("failed to get span batch ranges: %w", err)
	}

	return ranges, nil
}

// Get a list of span batch ranges for a given L2 block range.
// If the end block is reached, the last span batch range will end at end block, instead of the span batch's end block.
func GetSpanBatchRanges(config reassemble.Config, rollupCfg *rollup.Config, startBlock uint64, endBlock uint64, maxSpanBatchDeviation uint64) ([][2]uint64, error) {
	var ranges [][2]uint64
	currentStart := startBlock

	for currentStart < endBlock {
		_, spanEnd, err := reassemble.GetSpanBatchRange(config, rollupCfg, currentStart, 1000000)
		batchEnd := spanEnd
		if err != nil {
			// If we hit an error, log it as a warning
			fmt.Printf("Error getting span batch range: %v\n", err)
			batchEnd = currentStart + 100
		}
		if batchEnd > endBlock {
			batchEnd = endBlock
		}
		ranges = append(ranges, [2]uint64{currentStart, batchEnd})
		currentStart = batchEnd + 1
	}
	return ranges, nil
}

// Get the L1 origin corresponding to the given L2 block and the latest finalized L1 block.
func getL1OriginAndFinalized(rollupClient *sources.RollupClient, startBlock uint64, endBlock uint64) (uint64, uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	output, err := rollupClient.OutputAtBlock(ctx, startBlock)
	if err != nil {
		return 0, 0, err
	}
	startL1Origin := output.BlockRef.L1Origin.Number

	// Get L1 origin for the end block
	output, err = rollupClient.OutputAtBlock(ctx, endBlock)
	if err != nil {
		return 0, 0, err
	}
	endL1Origin := output.BlockRef.L1Origin.Number

	return startL1Origin, endL1Origin, nil
}
