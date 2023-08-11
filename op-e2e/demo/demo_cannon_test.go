package demo

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/ethereum-optimism/optimism/op-chain-ops/deployer"
	"github.com/ethereum-optimism/optimism/op-chain-ops/genesis"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/disputegame"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestCreateGame(t *testing.T) {
	ctx := context.Background()
	var deployments genesis.L1Deployments
	data, err := os.ReadFile("../../.devnet/addresses.json")
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &deployments))

	client, err := ethclient.DialContext(ctx, "http://localhost:8545")
	require.NoError(t, err)
	gameFactory := disputegame.NewFactoryHelper(t, ctx, &deployments, client)
	game := gameFactory.StartCannonGame(ctx, common.Hash{0xaa})
	game.LogGameData(ctx)

	t.Logf("Deployer addr: %v", deployer.TestAddress)
	claimCount := int64(2)
	for claimCount <= 7 {
		game.WaitForClaimCount(ctx, claimCount)
		claimCount++
		game.LogGameData(ctx)
	}

	for claimCount <= 30 {
		game.WaitForClaimCount(ctx, claimCount)
		game.Attack(ctx, claimCount-1, common.Hash{0xbb})
		claimCount += 2
		game.LogGameData(ctx)
	}

	game.WaitForClaimAtMaxDepth(ctx, true)
	game.LogGameData(ctx)
}
