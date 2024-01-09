package op_e2e

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	bss "github.com/ethereum-optimism/optimism/op-batcher/batcher"
	"github.com/ethereum-optimism/optimism/op-batcher/compressor"
	con "github.com/ethereum-optimism/optimism/op-conductor/conductor"
	conrpc "github.com/ethereum-optimism/optimism/op-conductor/rpc"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/wait"
	rollupNode "github.com/ethereum-optimism/optimism/op-node/node"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum-optimism/optimism/op-node/rollup/driver"
	"github.com/ethereum-optimism/optimism/op-node/rollup/sync"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	oprpc "github.com/ethereum-optimism/optimism/op-service/rpc"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
)

const (
	sequencer1Name = "sequencer1"
	sequencer2Name = "sequencer2"
	sequencer3Name = "sequencer3"
	verifierName   = "verifier"

	sequencer1Port = 9001
	sequencer2Port = 9002
	sequencer3Port = 9003

	conductor1ConsPort = 50001
	conductor2ConsPort = 50002
	conductor3ConsPort = 50003

	conductor1RpcPort = 50051
	conductor2RpcPort = 50052
	conductor3RpcPort = 50053

	localhost = "127.0.0.1"
)

type conductor struct {
	service *con.OpConductor
	client  conrpc.API
}

func setupSequencerFailoverTest(t *testing.T) (*System, map[string]*conductor) {
	InitParallel(t)
	ctx := context.Background()

	// 3 sequencers, 1 verifier, 1 active sequencer.
	cfg := sequencerFailoverSystemConfig(t)
	sys, err := cfg.Start(t)
	require.NoError(t, err)

	// 1 batcher that listens to all 3 sequencers, in started mode.
	setupBatcher(t, sys)

	// 3 conductors that connects to 1 sequencer each.
	conductors := make(map[string]*conductor)

	// initialize all conductors in paused mode
	conductorCfgs := []struct {
		consPort      int
		conductorPort int
		name          string
		bootstrap     bool
	}{
		{conductor1ConsPort, conductor1RpcPort, sequencer1Name, true}, // one in bootstrap mode so that we can form a cluster.
		{conductor2ConsPort, conductor2RpcPort, sequencer2Name, false},
		{conductor3ConsPort, conductor3RpcPort, sequencer3Name, false},
	}
	for _, cfg := range conductorCfgs {
		cfg := cfg
		nodePRC := sys.RollupNodes[cfg.name].HTTPEndpoint()
		engineRPC := sys.EthInstances[cfg.name].HTTPEndpoint()
		conductors[cfg.name] = setupConductor(t, cfg.consPort, cfg.conductorPort, cfg.name, t.TempDir(), nodePRC, engineRPC, cfg.bootstrap, *sys.RollupConfig)
	}

	// form a cluster
	c1 := conductors[sequencer1Name]
	c2 := conductors[sequencer2Name]
	c3 := conductors[sequencer3Name]

	require.NoError(t, waitForLeadershipChange(t, c1, true))
	require.NoError(t, c1.client.AddServerAsVoter(ctx, sequencer2Name, fmt.Sprintf("%s:%d", localhost, conductor2ConsPort)))
	require.NoError(t, c1.client.AddServerAsVoter(ctx, sequencer3Name, fmt.Sprintf("%s:%d", localhost, conductor3ConsPort)))
	require.True(t, leader(t, ctx, c1))
	require.False(t, leader(t, ctx, c2))
	require.False(t, leader(t, ctx, c3))

	// weirdly, batcher does not submit a batch until unsafe block 9.
	// It became normal after that and submits a batch every L1 block (2s) per configuration.
	// Since our health monitor checks on safe head progression, wait for batcher to become normal before proceeding.
	require.NoError(t, wait.ForNextSafeBlock(ctx, sys.Clients[sequencer1Name]))
	require.NoError(t, wait.ForNextSafeBlock(ctx, sys.Clients[sequencer1Name]))
	require.NoError(t, wait.ForNextSafeBlock(ctx, sys.Clients[sequencer1Name]))

	// make sure conductor reports all sequencers as healthy, this means they're syncing correctly.
	require.True(t, healthy(t, ctx, c1))
	require.True(t, healthy(t, ctx, c2))
	require.True(t, healthy(t, ctx, c3))

	// unpause all conductors
	require.NoError(t, c1.client.Resume(ctx))
	require.NoError(t, c2.client.Resume(ctx))
	require.NoError(t, c3.client.Resume(ctx))

	// final check, make sure everything is in the right place
	require.True(t, conductorActive(t, ctx, c1))
	require.True(t, conductorActive(t, ctx, c2))
	require.True(t, conductorActive(t, ctx, c3))

	require.True(t, sequencerActive(t, ctx, sys.RollupClient(sequencer1Name)))
	require.False(t, sequencerActive(t, ctx, sys.RollupClient(sequencer2Name)))
	require.False(t, sequencerActive(t, ctx, sys.RollupClient(sequencer3Name)))

	require.True(t, healthy(t, ctx, c1))
	require.True(t, healthy(t, ctx, c2))
	require.True(t, healthy(t, ctx, c3))

	return sys, conductors
}

func setupConductor(
	t *testing.T,
	consPort, conductorPort int,
	serverID, dir, nodePRC, engineRPC string,
	bootstrap bool,
	rollupCfg rollup.Config,
) *conductor {
	cfg := con.Config{
		ConsensusAddr:  localhost,
		ConsensusPort:  consPort,
		RaftServerID:   serverID,
		RaftStorageDir: dir,
		RaftBootstrap:  bootstrap,
		NodeRPC:        nodePRC,
		ExecutionRPC:   engineRPC,
		Paused:         true,
		HealthCheck: con.HealthCheckConfig{
			Interval:     1, // per test setup, l2 block time is 1s.
			SafeInterval: 4, // per test setup (l1 block time = 2s, max channel duration = 1, 2s buffer)
			MinPeerCount: 2, // per test setup, each sequencer has 2 peers
		},
		RollupCfg: rollupCfg,
		LogConfig: oplog.CLIConfig{
			Level: log.LvlInfo,
			Color: false,
		},
		RPC: oprpc.CLIConfig{
			ListenAddr: localhost,
			ListenPort: conductorPort,
		},
	}

	ctx := context.Background()
	service, err := con.New(ctx, &cfg, testlog.Logger(t, log.LvlInfo), "0.0.1")
	require.NoError(t, err)
	err = service.Start(ctx)
	require.NoError(t, err)

	rawClient, err := rpc.DialContext(ctx, service.HTTPEndpoint())
	require.NoError(t, err)
	client := conrpc.NewAPIClient(rawClient)

	return &conductor{
		service: service,
		client:  client,
	}
}

func setupBatcher(t *testing.T, sys *System) {
	var batchType uint = derive.SingularBatchType
	if sys.Cfg.DeployConfig.L2GenesisDeltaTimeOffset != nil && *sys.Cfg.DeployConfig.L2GenesisDeltaTimeOffset == hexutil.Uint64(0) {
		batchType = derive.SpanBatchType
	}
	batcherMaxL1TxSizeBytes := sys.Cfg.BatcherMaxL1TxSizeBytes
	if batcherMaxL1TxSizeBytes == 0 {
		batcherMaxL1TxSizeBytes = 240_000
	}

	// enable active sequencer follow mode.
	l2EthRpc := strings.Join([]string{
		sys.EthInstances[sequencer1Name].WSEndpoint(),
		sys.EthInstances[sequencer2Name].WSEndpoint(),
		sys.EthInstances[sequencer3Name].WSEndpoint(),
	}, ",")
	rollupRpc := strings.Join([]string{
		sys.RollupNodes[sequencer1Name].HTTPEndpoint(),
		sys.RollupNodes[sequencer2Name].HTTPEndpoint(),
		sys.RollupNodes[sequencer3Name].HTTPEndpoint(),
	}, ",")
	batcherCLIConfig := &bss.CLIConfig{
		L1EthRpc:               sys.EthInstances["l1"].WSEndpoint(),
		L2EthRpc:               l2EthRpc,
		RollupRpc:              rollupRpc,
		MaxPendingTransactions: 0,
		MaxChannelDuration:     1,
		MaxL1TxSize:            batcherMaxL1TxSizeBytes,
		CompressorConfig: compressor.CLIConfig{
			TargetL1TxSizeBytes: sys.Cfg.BatcherTargetL1TxSizeBytes,
			TargetNumFrames:     1,
			ApproxComprRatio:    0.4,
		},
		SubSafetyMargin: 0,
		PollInterval:    50 * time.Millisecond,
		TxMgrConfig:     newTxMgrConfig(sys.EthInstances["l1"].WSEndpoint(), sys.Cfg.Secrets.Batcher),
		LogConfig: oplog.CLIConfig{
			Level:  log.LvlInfo,
			Format: oplog.FormatText,
		},
		Stopped:   false,
		BatchType: batchType,
	}

	batcher, err := bss.BatcherServiceFromCLIConfig(context.Background(), "0.0.1", batcherCLIConfig, sys.Cfg.Loggers["batcher"])
	require.NoError(t, err)
	err = batcher.Start(context.Background())
	require.NoError(t, err)
	sys.BatchSubmitter = batcher
}

func sequencerFailoverSystemConfig(t *testing.T) SystemConfig {
	cfg := DefaultSystemConfig(t)
	delete(cfg.Nodes, "sequencer")
	cfg.Nodes[sequencer1Name] = sequencerCfg(sequencer1Port, true)
	cfg.Nodes[sequencer2Name] = sequencerCfg(sequencer2Port, false)
	cfg.Nodes[sequencer3Name] = sequencerCfg(sequencer3Port, false)

	delete(cfg.Loggers, "sequencer")
	cfg.Loggers[sequencer1Name] = testlog.Logger(t, log.LvlInfo).New("role", sequencer1Name)
	cfg.Loggers[sequencer2Name] = testlog.Logger(t, log.LvlInfo).New("role", sequencer2Name)
	cfg.Loggers[sequencer3Name] = testlog.Logger(t, log.LvlInfo).New("role", sequencer3Name)

	cfg.P2PTopology = map[string][]string{
		sequencer1Name: {sequencer2Name, sequencer3Name},
		sequencer2Name: {sequencer3Name, verifierName},
		sequencer3Name: {verifierName, sequencer1Name},
		verifierName:   {sequencer1Name, sequencer2Name},
	}

	return cfg
}

func sequencerCfg(port int, sequencerEnabled bool) *rollupNode.Config {
	return &rollupNode.Config{
		Driver: driver.Config{
			VerifierConfDepth:  0,
			SequencerConfDepth: 0,
			SequencerEnabled:   sequencerEnabled,
		},
		// Submitter PrivKey is set in system start for rollup nodes where sequencer = true
		RPC: rollupNode.RPCConfig{
			ListenAddr:  localhost,
			ListenPort:  port,
			EnableAdmin: true,
		},
		L1EpochPollInterval:         time.Second * 2,
		RuntimeConfigReloadInterval: time.Minute * 10,
		ConfigPersistence:           &rollupNode.DisabledConfigPersistence{},
		Sync:                        sync.Config{SyncMode: sync.CLSync},
	}
}

func waitForLeadershipChange(t *testing.T, c *conductor, leader bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			isLeader, err := c.client.Leader(ctx)
			if err != nil {
				return err
			}
			if isLeader == leader {
				return nil
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func leader(t *testing.T, ctx context.Context, con *conductor) bool {
	leader, err := con.client.Leader(ctx)
	require.NoError(t, err)
	return leader
}

func healthy(t *testing.T, ctx context.Context, con *conductor) bool {
	healthy, err := con.client.SequencerHealthy(ctx)
	require.NoError(t, err)
	return healthy
}

func conductorActive(t *testing.T, ctx context.Context, con *conductor) bool {
	active, err := con.client.Active(ctx)
	require.NoError(t, err)
	return active
}

func sequencerActive(t *testing.T, ctx context.Context, rollupClient *sources.RollupClient) bool {
	active, err := rollupClient.SequencerActive(ctx)
	require.NoError(t, err)
	return active
}
