package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum-optimism/optimism/op-node/cmd/batch_decoder/fetch"
	"github.com/ethereum-optimism/optimism/op-node/cmd/batch_decoder/reassemble"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/dial"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "batch-decoder"
	app.Usage = "Optimism Batch Decoding Utility"
	app.Commands = []*cli.Command{
		{
			Name:  "fetch",
			Usage: "Fetches batches in the specified range",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:     "start",
					Required: true,
					Usage:    "First block (inclusive) to fetch",
				},
				&cli.IntFlag{
					Name:     "end",
					Required: true,
					Usage:    "Last block (exclusive) to fetch",
				},
				&cli.StringFlag{
					Name:     "inbox",
					Required: true,
					Usage:    "Batch Inbox Address",
				},
				&cli.StringFlag{
					Name:     "sender",
					Required: true,
					Usage:    "Batch Sender Address",
				},
				&cli.StringFlag{
					Name:  "out",
					Value: "/tmp/batch_decoder/transactions_cache",
					Usage: "Cache directory for the found transactions",
				},
				&cli.StringFlag{
					Name:     "l1",
					Required: true,
					Usage:    "L1 RPC URL",
					EnvVars:  []string{"L1_RPC"},
				},
				&cli.StringFlag{
					Name:     "l1.beacon",
					Required: false,
					Usage:    "Address of L1 Beacon-node HTTP endpoint to use",
					EnvVars:  []string{"L1_BEACON"},
				},
				&cli.IntFlag{
					Name:  "concurrent-requests",
					Value: 10,
					Usage: "Concurrency level when fetching L1",
				},
			},
			Action: func(cliCtx *cli.Context) error {
				l1Client, err := ethclient.Dial(cliCtx.String("l1"))
				if err != nil {
					log.Fatal(err)
				}
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				chainID, err := l1Client.ChainID(ctx)
				if err != nil {
					log.Fatal(err)
				}
				beaconAddr := cliCtx.String("l1.beacon")
				var beacon *sources.L1BeaconClient
				if beaconAddr != "" {
					beaconClient := sources.NewBeaconHTTPClient(client.NewBasicHTTPClient(beaconAddr, nil))
					beaconCfg := sources.L1BeaconClientConfig{FetchAllSidecars: false}
					beacon = sources.NewL1BeaconClient(beaconClient, beaconCfg)
					_, err := beacon.GetVersion(ctx)
					if err != nil {
						log.Fatal(fmt.Errorf("failed to check L1 Beacon API version: %w", err))
					}
				} else {
					fmt.Println("L1 Beacon endpoint not set. Unable to fetch post-ecotone channel frames")
				}
				config := fetch.Config{
					Start:   uint64(cliCtx.Int("start")),
					End:     uint64(cliCtx.Int("end")),
					ChainID: chainID,
					BatchSenders: map[common.Address]struct{}{
						common.HexToAddress(cliCtx.String("sender")): {},
					},
					BatchInbox:         common.HexToAddress(cliCtx.String("inbox")),
					OutDirectory:       cliCtx.String("out"),
					ConcurrentRequests: uint64(cliCtx.Int("concurrent-requests")),
				}
				totalValid, totalInvalid := fetch.Batches(l1Client, beacon, config)
				fmt.Printf("Fetched batches in range [%v,%v). Found %v valid & %v invalid batches\n", config.Start, config.End, totalValid, totalInvalid)
				fmt.Printf("Fetch Config: Chain ID: %v. Inbox Address: %v. Valid Senders: %v.\n", config.ChainID, config.BatchInbox, config.BatchSenders)
				fmt.Printf("Wrote transactions with batches to %v\n", config.OutDirectory)
				return nil
			},
		},
		{
			Name:  "reassemble",
			Usage: "Reassembles channels from fetched batch transactions and decode batches",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "in",
					Value: "/tmp/batch_decoder/transactions_cache",
					Usage: "Cache directory for the found transactions",
				},
				&cli.StringFlag{
					Name:  "out",
					Value: "/tmp/batch_decoder/channel_cache",
					Usage: "Cache directory for the found channels",
				},
				&cli.Uint64Flag{
					Name:  "l2-chain-id",
					Value: 10,
					Usage: "L2 chain id for span batch derivation. Default value from op-mainnet.",
				},
				&cli.Uint64Flag{
					Name:  "l2-genesis-timestamp",
					Value: 1686068903,
					Usage: "L2 genesis time for span batch derivation. Default value from op-mainnet. " +
						"Superchain-registry prioritized when given value is inconsistent.",
				},
				&cli.Uint64Flag{
					Name:  "l2-block-time",
					Value: 2,
					Usage: "L2 block time for span batch derivation. Default value from op-mainnet. " +
						"Superchain-registry prioritized when given value is inconsistent.",
				},
				&cli.StringFlag{
					Name:  "inbox",
					Value: "0xFF00000000000000000000000000000000000010",
					Usage: "Batch Inbox Address. Default value from op-mainnet. " +
						"Superchain-registry prioritized when given value is inconsistent.",
				},
			},
			Action: func(cliCtx *cli.Context) error {
				var (
					L2GenesisTime     uint64         = cliCtx.Uint64("l2-genesis-timestamp")
					L2BlockTime       uint64         = cliCtx.Uint64("l2-block-time")
					BatchInboxAddress common.Address = common.HexToAddress(cliCtx.String("inbox"))
				)
				L2ChainID := new(big.Int).SetUint64(cliCtx.Uint64("l2-chain-id"))
				rollupCfg, err := rollup.LoadOPStackRollupConfig(L2ChainID.Uint64())
				if err == nil {
					// prioritize superchain config
					if L2GenesisTime != rollupCfg.Genesis.L2Time {
						L2GenesisTime = rollupCfg.Genesis.L2Time
						fmt.Printf("L2GenesisTime overridden: %v\n", L2GenesisTime)
					}
					if L2BlockTime != rollupCfg.BlockTime {
						L2BlockTime = rollupCfg.BlockTime
						fmt.Printf("L2BlockTime overridden: %v\n", L2BlockTime)
					}
					if BatchInboxAddress != rollupCfg.BatchInboxAddress {
						BatchInboxAddress = rollupCfg.BatchInboxAddress
						fmt.Printf("BatchInboxAddress overridden: %v\n", BatchInboxAddress)
					}
				}
				config := reassemble.Config{
					BatchInbox:    BatchInboxAddress,
					InDirectory:   cliCtx.String("in"),
					OutDirectory:  cliCtx.String("out"),
					L2ChainID:     L2ChainID,
					L2GenesisTime: L2GenesisTime,
					L2BlockTime:   L2BlockTime,
				}
				reassemble.Channels(config, rollupCfg)
				return nil
			},
		},
		{
			Name:  "get-range",
			Usage: "For a given L2 block number, gets the full range of the span batch that it's a part of",
			Flags: []cli.Flag{
				&cli.Uint64Flag{
					Name:  "start",
					Usage: "The first block in the range ",
				},
				&cli.Uint64Flag{
					Name:  "end",
					Usage: "L2 block number to get the span batch range for.",
				},
				&cli.StringFlag{
					Name:  "in",
					Value: "/tmp/batch_decoder/transactions_cache_new",
					Usage: "Cache directory for the found transactions",
				},
				&cli.StringFlag{
					Name:     "l2",
					Required: true,
					Usage:    "L2 RPC URL",
					EnvVars:  []string{"L2_RPC"},
				},
				&cli.StringFlag{
					Name:     "l2.node",
					Required: true,
					Usage:    "L2 node URL",
					EnvVars:  []string{"L2_NODE"},
				},
				&cli.StringFlag{
					Name:     "l1",
					Required: true,
					Usage:    "L1 RPC URL",
					EnvVars:  []string{"L1_RPC"},
				},
				&cli.StringFlag{
					Name:     "l1.beacon",
					Required: false,
					Usage:    "Address of L1 Beacon-node HTTP endpoint to use",
					EnvVars:  []string{"L1_BEACON"},
				},
				&cli.StringFlag{
					Name:     "sender",
					Required: true,
					Usage:    "Batch Sender Address",
				},
				&cli.Uint64Flag{
					Name:  "l2-chain-id",
					Value: 10,
					Usage: "L2 chain id for span batch derivation. Default value from op-mainnet.",
				},
				&cli.Uint64Flag{
					Name:  "l2-genesis-timestamp",
					Value: 1686068903,
					Usage: "L2 genesis time for span batch derivation. Default value from op-mainnet. " +
						"Superchain-registry prioritized when given value is inconsistent.",
				},
				&cli.Uint64Flag{
					Name:  "l2-genesis-block",
					Value: 105235063,
					Usage: "L2 genesis block for span batch derivation. Default value from op-mainnet. " +
						"Superchain-registry prioritized when given value is inconsistent.",
				},
				&cli.Uint64Flag{
					Name:  "l2-block-time",
					Value: 2,
					Usage: "L2 block time for span batch derivation. Default value from op-mainnet. " +
						"Superchain-registry prioritized when given value is inconsistent.",
				},
				&cli.StringFlag{
					Name:  "inbox",
					Value: "0xFF00000000000000000000000000000000000010",
					Usage: "Batch Inbox Address. Default value from op-mainnet. " +
						"Superchain-registry prioritized when given value is inconsistent.",
				},
			},
			Action: func(cliCtx *cli.Context) error {
				var (
					L2GenesisTime     uint64         = cliCtx.Uint64("l2-genesis-timestamp")
					L2GenesisBlock    uint64         = cliCtx.Uint64("l2-genesis-block")
					L2BlockTime       uint64         = cliCtx.Uint64("l2-block-time")
					BatchInboxAddress common.Address = common.HexToAddress(cliCtx.String("inbox"))
					StartBlock        uint64         = cliCtx.Uint64("start")
					EndBlock          uint64         = cliCtx.Uint64("end")
				)

				// Get the L1 origin corresponding to the start block
				// nextBlock is equal to the highest value in the `EndBlock` column of the db, plus 1
				L2ChainID := new(big.Int).SetUint64(cliCtx.Uint64("l2-chain-id"))
				rollupCfg, err := rollup.LoadOPStackRollupConfig(L2ChainID.Uint64())
				if err == nil {
					// prioritize superchain config
					if L2GenesisTime != rollupCfg.Genesis.L2Time {
						L2GenesisTime = rollupCfg.Genesis.L2Time
						fmt.Printf("L2GenesisTime overridden: %v\n", L2GenesisTime)
					}
					if L2GenesisBlock != rollupCfg.Genesis.L2.Number {
						L2GenesisBlock = rollupCfg.Genesis.L2.Number
						fmt.Printf("L2GenesisBlock overridden: %v\n", L2GenesisBlock)
					}
					if L2BlockTime != rollupCfg.BlockTime {
						L2BlockTime = rollupCfg.BlockTime
						fmt.Printf("L2BlockTime overridden: %v\n", L2BlockTime)
					}
					if BatchInboxAddress != rollupCfg.BatchInboxAddress {
						BatchInboxAddress = rollupCfg.BatchInboxAddress
						fmt.Printf("BatchInboxAddress overridden: %v\n", BatchInboxAddress)
					}
				}
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				rollupClient, err := dial.DialRollupClientWithTimeout(ctx, dial.DefaultDialTimeout, nil, cliCtx.String("l2.node"))
				if err != nil {
					log.Fatal(err)
				}

				l1Origin, finalizedL1, err := getL1OriginAndFinalized(rollupClient, StartBlock)
				if err != nil {
					log.Fatal(err)
				}

				fetchConfig := fetch.Config{
					Start:   l1Origin,
					End:     finalizedL1,
					ChainID: rollupCfg.L1ChainID,
					BatchSenders: map[common.Address]struct{}{
						common.HexToAddress(cliCtx.String("sender")): {},
					},
					BatchInbox:         BatchInboxAddress,
					OutDirectory:       cliCtx.String("in"),
					ConcurrentRequests: 10,
				}

				l1Client, err := ethclient.Dial(cliCtx.String("l1"))
				if err != nil {
					log.Fatal(err)
				}
				beaconAddr := cliCtx.String("l1.beacon")
				var beacon *sources.L1BeaconClient
				if beaconAddr != "" {
					beaconClient := sources.NewBeaconHTTPClient(client.NewBasicHTTPClient(beaconAddr, nil))
					beaconCfg := sources.L1BeaconClientConfig{FetchAllSidecars: false}
					beacon = sources.NewL1BeaconClient(beaconClient, beaconCfg)
					_, err := beacon.GetVersion(ctx)
					if err != nil {
						log.Fatal(fmt.Errorf("failed to check L1 Beacon API version: %w", err))
					}
				} else {
					fmt.Println("L1 Beacon endpoint not set. Unable to fetch post-ecotone channel frames")
				}
				totalValid, totalInvalid := fetch.Batches(l1Client, beacon, fetchConfig)
				fmt.Printf("Fetched batches in range [%v,%v). Found %v valid & %v invalid batches\n", fetchConfig.Start, fetchConfig.End, totalValid, totalInvalid)

				config := reassemble.Config{
					BatchInbox:    BatchInboxAddress,
					InDirectory:   cliCtx.String("in"),
					OutDirectory:  "",
					L2ChainID:     L2ChainID,
					L2GenesisTime: L2GenesisTime,
					L2BlockTime:   L2BlockTime,
				}

				ranges, err := GetSpanBatchRanges(config, rollupCfg, StartBlock, EndBlock, 1000000)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Span batch ranges: %v\n", ranges)

				return nil
			},
		},
		{
			Name:  "force-close",
			Usage: "Create the tx data which will force close a channel",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "id",
					Required: true,
					Usage:    "ID of the channel to close",
				},
				&cli.StringFlag{
					Name:  "inbox",
					Value: "0x0000000000000000000000000000000000000000",
					Usage: "(Optional) Batch Inbox Address",
				},
				&cli.StringFlag{
					Name:  "in",
					Value: "/tmp/batch_decoder/transactions_cache",
					Usage: "Cache directory for the found transactions",
				},
			},
			Action: func(cliCtx *cli.Context) error {
				var id derive.ChannelID
				if err := (&id).UnmarshalText([]byte(cliCtx.String("id"))); err != nil {
					log.Fatal(err)
				}
				frames := reassemble.LoadFrames(cliCtx.String("in"), common.HexToAddress(cliCtx.String("inbox")))
				var filteredFrames []derive.Frame
				for _, frame := range frames {
					if frame.Frame.ID == id {
						filteredFrames = append(filteredFrames, frame.Frame)
					}
				}
				data, err := derive.ForceCloseTxData(filteredFrames)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%x\n", data)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// Get a list of span batch ranges for a given L2 block range
func GetSpanBatchRanges(config reassemble.Config, rollupCfg *rollup.Config, startBlock uint64, endBlock uint64, maxSpanBatchDeviation uint64) ([]uint64, error) {
	var ranges []uint64
	ranges = append(ranges, startBlock)

	start, end, err := reassemble.GetSpanBatchRange(config, rollupCfg, startBlock, 1000000)
	currentStart := startBlock + 100
	if err != nil {
		fmt.Printf("Error getting span batch range: %v\n", err)
	} else {
		fmt.Printf("Span batch range: [%v,%v)\n", start, end)
		currentStart = end
	}

	for currentStart < endBlock {
		ranges = append(ranges, currentStart)
		spanStart, spanEnd, err := reassemble.GetSpanBatchRange(config, rollupCfg, currentStart, 1000000)
		if err != nil {
			// If we hit an error, log it as a warning
			fmt.Printf("Error getting span batch range: %v\n", err)
			currentStart += 100
		} else {
			fmt.Printf("Span batch range: [%v,%v)\n", spanStart, spanEnd)
			currentStart = spanEnd + 1
		}
	}
	ranges = append(ranges, endBlock)
	return ranges, nil
}

// Get the L1 origin corresponding to the given L2 block and the latest finalized L1 block.
func getL1OriginAndFinalized(rollupClient *sources.RollupClient, l2Block uint64) (uint64, uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	output, err := rollupClient.OutputAtBlock(ctx, l2Block)
	if err != nil {
		return 0, 0, err
	}
	l1Origin := output.BlockRef.L1Origin.Number

	// get the latest finalized L1
	status, err := rollupClient.SyncStatus(ctx)
	if err != nil {
		return 0, 0, err
	}
	finalizedL1 := status.FinalizedL1.Number

	return l1Origin, finalizedL1, nil
}
