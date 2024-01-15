package client

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	cldr "github.com/ethereum-optimism/optimism/op-program/client/driver"
	"github.com/ethereum-optimism/optimism/op-program/client/l1"
	"github.com/ethereum-optimism/optimism/op-program/client/l2"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

// L2Game represents a state-transition proof of an output-root,
// starting at a previously agreed upon L2 output-root, and an agreed-upon L1 state.
type L2Game struct {
	L1Head common.Hash

	// L1SuperRoot commits to the pre-agreed merkleized L1 chain,
	// L2 games are only played if we agree upon this.
	// A L1Game is played to establish the merkleized L1 state otherwise.
	// This is our source of truth for quick L1 access, without pre-image-size problems.
	L1SuperRoot common.Hash

	L2Claim            common.Hash
	L2ClaimBlockNumber uint64
	L2Prestate         common.Hash // f.k.a. L2OutputRoot, but this is more descriptive

	L2ChainConfig *params.ChainConfig
	RollupConfig  *rollup.Config
}

func (l2Game *L2Game) Run(logger log.Logger, pClient *preimage.OracleClient, hClient *preimage.HintWriter) error {
	// TODO enhance L1 preimage oracle
	l1PreimageOracle := l1.NewCachingOracle(l1.NewPreimageOracle(pClient, hClient))
	l2PreimageOracle := l2.NewCachingOracle(l2.NewPreimageOracle(pClient, hClient))

	logger.Info("L2 Game Bootstrapped", "gameL2", l2Game)
	return runDerivation(
		logger,
		l2Game.RollupConfig,
		l2Game.L2ChainConfig,
		l2Game.L1Head,
		l2Game.L2OutputRoot,
		l2Game.L2Claim,
		l2Game.L2ClaimBlockNumber,
		l1PreimageOracle,
		l2PreimageOracle,
	)
}

// runDerivation executes the L2 state transition, given a minimal interface to retrieve data.
func runDerivation(logger log.Logger, cfg *rollup.Config, l2Cfg *params.ChainConfig, l1Head common.Hash,
	l2OutputRoot common.Hash, l2Claim common.Hash, l2ClaimBlockNum uint64, l1Oracle l1.Oracle, l2Oracle l2.L2Oracle) error {

	l1Source := l1.NewOracleL1Client(logger, l1Oracle, l1Head)
	engineBackend, err := l2.NewOracleBackedL2Chain(logger, l2Oracle, l2Cfg, l2OutputRoot)
	if err != nil {
		return fmt.Errorf("failed to create oracle-backed L2 chain: %w", err)
	}
	l2Source := l2.NewOracleEngine(cfg, logger, engineBackend)

	logger.Info("Starting derivation")
	d := cldr.NewDriver(logger, cfg, l1Source, l2Source, l2ClaimBlockNum)
	for {
		if err = d.Step(context.Background()); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
	}
	return d.ValidateClaim(l2ClaimBlockNum, eth.Bytes32(l2Claim))
}
