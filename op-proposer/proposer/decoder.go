package proposer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum-optimism/optimism/op-node/cmd/batch_decoder/fetch"
	"github.com/ethereum-optimism/optimism/op-node/cmd/batch_decoder/reassemble"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum/go-ethereum/common"
)

func (l *L2OutputSubmitter) FetchBatchesFromChain(ctx context.Context, nextBlock uint64) error {
	proposerConfig := l.DriverSetup.Cfg
	l1Client := l.DriverSetup.L1Client

	l1Origin, finalizedL1, err := l.getL1OriginAndFinalized(ctx, nextBlock)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	chainID, err := l1Client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	beaconAddr := proposerConfig.beaconRpc
	var beacon *sources.L1BeaconClient
	if beaconAddr != "" {
		beaconClient := sources.NewBeaconHTTPClient(client.NewBasicHTTPClient(beaconAddr, nil))
		beaconCfg := sources.L1BeaconClientConfig{FetchAllSidecars: false}
		beacon = sources.NewL1BeaconClient(beaconClient, beaconCfg)
		_, err := beacon.GetVersion(ctx)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to check L1 Beacon API version: %w", err))
			return err
		}
	} else {
		fmt.Println("L1 Beacon endpoint not set. Unable to fetch post-ecotone channel frames")
		return err
	}
	rollupCfg, err := rollup.LoadOPStackRollupConfig(chainID)

	fetchConfig := fetch.Config{
		Start:   l1Origin,
		End:     finalizedL1,
		ChainID: chainID,
		BatchSenders: map[common.Address]struct{}{
			rollupCfg.Genesis.SystemConfig.BatcherAddr: {},
		},
		BatchInbox:         rollupCfg.BatchInboxAddress,
		OutDirectory:       proposerConfig.txCacheOutDir,
		ConcurrentRequests: proposerConfig.batchDecoderConcurrentReqs,
	}

	fetch.Batches(l1Client, beacon, fetchConfig)
	return nil
}

func (l *L2OutputSubmitter) GenerateSpanBatchRange(nextBlock uint64) (uint64, uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	chainID, err := l.DriverSetup.L1Client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	rollupCfg, err := rollup.LoadOPStackRollupConfig(chainID)

	reassembleConfig := reassemble.Config{
		BatchInbox:    rollupCfg.BatchInboxAddress,
		InDirectory:   l.DriverSetup.Cfg.txCacheOutDir,
		OutDirectory:  "",
		L2ChainID:     chainID,
		L2GenesisTime: rollupCfg.Genesis.L2Time,
		L2BlockTime:   rollupCfg.BlockTime,
	}

	return reassemble.GetSpanBatchRange(reassembleConfig, rollupCfg, nextBlock)
}

func (l *L2OutputSubmitter) getL1OriginAndFinalized(ctx context.Context, nextBlock uint64) (uint64, uint64, error) {
	cCtx, cancel := context.WithTimeout(ctx, l.Cfg.NetworkTimeout)
	defer cancel()

	rollupClient, err := l.RollupProvider.RollupClient(ctx)
	if err != nil {
		l.Log.Error("proposer unable to get rollup client", "err", err)
		return 0, 0, err
	}

	output, err := rollupClient.OutputAtBlock(cCtx, nextBlock)
	if err != nil {
		l.Log.Error("proposer unable to get sync status", "err", err)
		return 0, 0, err
	}
	l1Origin := output.BlockRef.L1Origin.Number

	// get the latest finalized L1
	status, err := rollupClient.SyncStatus(cCtx)
	if err != nil {
		l.Log.Error("proposer unable to get sync status", "err", err)
		return 0, 0, err
	}
	finalizedL1 := status.FinalizedL1.Number

	return l1Origin, finalizedL1, nil
}
