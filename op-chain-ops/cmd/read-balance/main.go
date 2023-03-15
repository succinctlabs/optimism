package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum-optimism/optimism/op-chain-ops/db"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/mattn/go-isatty"
	"os"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/log"

	"github.com/urfave/cli"
)

func main() {
	log.Root().SetHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(isatty.IsTerminal(os.Stderr.Fd()))))

	app := &cli.App{
		Name:  "read-balance",
		Usage: "read a raw balance",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "db-path",
				Usage:    "Path to database",
				Required: true,
			},
			cli.StringFlag{
				Name:     "address",
				Usage:    "Address to read the balance for",
				Required: true,
			},
			cli.IntFlag{
				Name:  "db-cache",
				Usage: "LevelDB cache size in mb",
				Value: 1024,
			},
			cli.IntFlag{
				Name:  "db-handles",
				Usage: "LevelDB number of handles",
				Value: 60,
			},
		},
		Action: func(ctx *cli.Context) error {
			dbCache := ctx.Int("db-cache")
			dbHandles := ctx.Int("db-handles")
			ldb, err := db.Open(ctx.String("db-path"), dbCache, dbHandles)
			if err != nil {
				return fmt.Errorf("cannot open database: %w", err)
			}

			hash := rawdb.ReadHeadHeaderHash(ldb)
			log.Info("Reading chain tip from database", "hash", hash)

			// Grab the header number.
			num := rawdb.ReadHeaderNumber(ldb, hash)
			if num == nil {
				return fmt.Errorf("cannot find header number for %s", hash)
			}

			// Grab the full header.
			header := rawdb.ReadHeader(ldb, hash, *num)
			log.Info("Read header from database", "number", *num)

			underlyingDB := state.NewDatabaseWithConfig(ldb, &trie.Config{
				Preimages: true,
				Cache:     1024,
			})

			// Open up the state database.
			sdb, err := state.New(header.Root, underlyingDB, nil)
			if err != nil {
				return fmt.Errorf("cannot open StateDB: %w", err)
			}

			balance := sdb.GetBalance(common.HexToAddress(ctx.String("address")))
			fmt.Println(balance.String())
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Crit("error in migration", "err", err)
	}
}

func writeJSON(outfile string, input interface{}) error {
	f, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(input)
}
