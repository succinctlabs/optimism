package genesis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum-optimism/optimism/op-chain-ops/genesis"
	"github.com/ethereum-optimism/optimism/op-service/jsonutil"
)

var (
	l1RPCFlag = &cli.StringFlag{
		Name:  "l1-rpc",
		Usage: "RPC URL for an Ethereum L1 node. Cannot be used with --l1-starting-block",
	}
	l1StartingBlockFlag = &cli.PathFlag{
		Name:  "l1-starting-block",
		Usage: "Path to a JSON file containing the L1 starting block. Overrides the need for using an L1 RPC to fetch the block. Cannot be used with --l1-rpc",
	}
	deployConfigFlag = &cli.PathFlag{
		Name:     "deploy-config",
		Usage:    "Path to deploy config file",
		Required: true,
	}
	l1DeploymentsFlag = &cli.PathFlag{
		Name:  "l1-deployments",
		Usage: "Path to L1 deployments JSON file as in superchain-registry",
	}
	outfileL2Flag = &cli.PathFlag{
		Name:  "outfile.l2",
		Usage: "Path to L2 genesis output file",
	}
	outfileRollupFlag = &cli.PathFlag{
		Name:  "outfile.rollup",
		Usage: "Path to rollup output file",
	}

	l1AllocsFlag = &cli.StringFlag{
		Name:  "l1-allocs",
		Usage: "Path to L1 genesis state dump",
	}
	outfileL1Flag = &cli.StringFlag{
		Name:  "outfile.l1",
		Usage: "Path to L1 genesis output file",
	}
	l2AllocsFlag = &cli.StringFlag{
		Name:  "l2-allocs",
		Usage: "Path to L2 genesis state dump",
	}

	l1Flags = []cli.Flag{
		deployConfigFlag,
		l1AllocsFlag,
		l1DeploymentsFlag,
		outfileL1Flag,
	}

	l2Flags = []cli.Flag{
		l1RPCFlag,
		l1StartingBlockFlag,
		deployConfigFlag,
		l2AllocsFlag,
		outfileL2Flag,
		outfileRollupFlag,
	}
)

var Subcommands = cli.Commands{
	{
		Name:  "l1",
		Usage: "Generates a L1 genesis state file",
		Flags: l1Flags,
		Action: func(ctx *cli.Context) error {
			deployConfig := ctx.String("deploy-config")
			config, err := genesis.NewDeployConfig(deployConfig)
			if err != nil {
				return err
			}

			var deployments *genesis.L1Deployments
			if l1Deployments := ctx.String("l1-deployments"); l1Deployments != "" {
				deployments, err = genesis.NewL1Deployments(l1Deployments)
				if err != nil {
					return err
				}
			}

			if deployments != nil {
				config.SetDeployments(deployments)
			}

			if err := config.Check(); err != nil {
				return fmt.Errorf("deploy config at %s invalid: %w", deployConfig, err)
			}

			// Check the addresses after setting the deployments
			if err := config.CheckAddresses(); err != nil {
				return fmt.Errorf("deploy config at %s invalid: %w", deployConfig, err)
			}

			var dump *state.Dump
			if l1Allocs := ctx.String("l1-allocs"); l1Allocs != "" {
				dump, err = genesis.NewStateDump(l1Allocs)
				if err != nil {
					return err
				}
			}

			l1Genesis, err := genesis.BuildL1DeveloperGenesis(config, dump, deployments)
			if err != nil {
				return err
			}

			return jsonutil.WriteJSON(ctx.String("outfile.l1"), l1Genesis, 0o666)
		},
	},
	{
		Name:  "l2",
		Usage: "Generates an L2 genesis file and rollup config suitable for a deployed network",
		Description: "Generating the L2 genesis depends on knowledge of L1 contract addresses for the bridge to be secure. " +
			"A deploy config and L2 genesis allocs file are used to create the L2 genesis. " +
			"An L1 starting block is necessary, it can either be fetched dynamically using config in the deploy config " +
			"or it can be provided as a JSON file.",
		Flags: l2Flags,
		Action: func(ctx *cli.Context) error {
			deployConfig := ctx.Path("deploy-config")
			log.Info("Deploy config", "path", deployConfig)
			config, err := genesis.NewDeployConfig(deployConfig)
			if err != nil {
				return err
			}

			l1StartBlockPath := ctx.Path("l1-starting-block")
			l1RPC := ctx.String("l1-rpc")

			if l1StartBlockPath == "" && l1RPC == "" {
				return errors.New("must specify either --l1-starting-block or --l1-rpc")
			}
			if l1StartBlockPath != "" && l1RPC != "" {
				return errors.New("cannot specify both --l1-starting-block and --l1-rpc")
			}

			var l1StartBlock *types.Block
			if l1StartBlockPath != "" {
				if l1StartBlock, err = readBlockJSON(l1StartBlockPath); err != nil {
					return fmt.Errorf("cannot read L1 starting block at %s: %w", l1StartBlockPath, err)
				}
			}

			var l2Allocs *genesis.ForgeAllocs
			if l2AllocsPath := ctx.String("l2-allocs"); l2AllocsPath != "" {
				l2Allocs, err = genesis.LoadForgeAllocs(l2AllocsPath)
				if err != nil {
					return err
				}
			} else {
				return errors.New("missing l2-allocs")
			}

			if l1RPC != "" {
				client, err := ethclient.Dial(l1RPC)
				if err != nil {
					return fmt.Errorf("cannot dial %s: %w", l1RPC, err)
				}

				if config.L1StartingBlockTag == nil {
					l1StartBlock, err = client.BlockByNumber(context.Background(), nil)
					if err != nil {
						return fmt.Errorf("cannot fetch latest block: %w", err)
					}
					tag := rpc.BlockNumberOrHashWithHash(l1StartBlock.Hash(), true)
					config.L1StartingBlockTag = (*genesis.MarshalableRPCBlockNumberOrHash)(&tag)
				} else if config.L1StartingBlockTag.BlockHash != nil {
					l1StartBlock, err = client.BlockByHash(context.Background(), *config.L1StartingBlockTag.BlockHash)
					if err != nil {
						return fmt.Errorf("cannot fetch block by hash: %w", err)
					}
				} else if config.L1StartingBlockTag.BlockNumber != nil {
					l1StartBlock, err = client.BlockByNumber(context.Background(), big.NewInt(config.L1StartingBlockTag.BlockNumber.Int64()))
					if err != nil {
						return fmt.Errorf("cannot fetch block by number: %w", err)
					}
				}
			}

			// Ensure that there is a starting L1 block
			if l1StartBlock == nil {
				return errors.New("no starting L1 block")
			}

			// Sanity check the config. Do this after filling in the L1StartingBlockTag
			// if it is not defined.
			if err := config.Check(); err != nil {
				return err
			}

			log.Info("Using L1 Start Block", "number", l1StartBlock.Number(), "hash", l1StartBlock.Hash().Hex())

			// Build the L2 genesis block
			l2Genesis, err := genesis.BuildL2Genesis(config, l2Allocs, l1StartBlock)
			if err != nil {
				return fmt.Errorf("error creating l2 genesis: %w", err)
			}

			l2GenesisBlock := l2Genesis.ToBlock()
			rollupConfig, err := config.RollupConfig(l1StartBlock, l2GenesisBlock.Hash(), l2GenesisBlock.Number().Uint64())
			if err != nil {
				return err
			}
			if err := rollupConfig.Check(); err != nil {
				return fmt.Errorf("generated rollup config does not pass validation: %w", err)
			}

			if err := jsonutil.WriteJSON(ctx.String("outfile.l2"), l2Genesis, 0o666); err != nil {
				return err
			}
			return jsonutil.WriteJSON(ctx.String("outfile.rollup"), rollupConfig, 0o666)
		},
	},
}

// rpcBlock represents the JSON serialization of a block from an Ethereum RPC.
type rpcBlock struct {
	Hash         common.Hash         `json:"hash"`
	Transactions []rpcTransaction    `json:"transactions"`
	UncleHashes  []common.Hash       `json:"uncles"`
	Withdrawals  []*types.Withdrawal `json:"withdrawals,omitempty"`
}

// rpcTransaction represents the JSON serialization of a transaction from an Ethereum RPC.
type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

// txExtraInfo includes extra information about a transaction that is returned from
// and Ethereum RPC endpoint.
type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

// readBlockJSON will read a JSON file from disk containing a serialized block.
// This logic can break if the block format changes but there is no modular way
// to turn a block into JSON in go-ethereum.
func readBlockJSON(path string) (*types.Block, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("block file at %s not found: %w", path, err)
	}

	var header types.Header
	if err := json.Unmarshal(raw, &header); err != nil {
		return nil, fmt.Errorf("cannot unmarshal block: %w", err)
	}

	var body rpcBlock
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}

	if len(body.UncleHashes) > 0 {
		return nil, fmt.Errorf("cannot unmarshal block with uncles")
	}

	txs := make([]*types.Transaction, len(body.Transactions))
	for i, tx := range body.Transactions {
		txs[i] = tx.tx
	}
	return types.NewBlockWithHeader(&header).WithBody(txs, nil).WithWithdrawals(body.Withdrawals), nil
}
