package integration_tests

import (
	"context"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum-optimism/optimism/proxyd"
	ms "github.com/ethereum-optimism/optimism/proxyd/tools/mockserver/handler"
	"github.com/stretchr/testify/require"
)

func setup_failover(t *testing.T) (map[string]nodeContext, *proxyd.BackendGroup, *ProxydHTTPClient, func()) {
	// setup mock servers
	node1 := NewMockBackend(nil)
	node2 := NewMockBackend(nil)

	dir, err := os.Getwd()
	require.NoError(t, err)

	responses := path.Join(dir, "testdata/consensus_responses.yml")

	h1 := ms.MockedHandler{
		Overrides:    []*ms.MethodTemplate{},
		Autoload:     true,
		AutoloadFile: responses,
	}
	h2 := ms.MockedHandler{
		Overrides:    []*ms.MethodTemplate{},
		Autoload:     true,
		AutoloadFile: responses,
	}

	require.NoError(t, os.Setenv("NODE1_URL", node1.URL()))
	require.NoError(t, os.Setenv("NODE2_URL", node2.URL()))

	node1.SetHandler(http.HandlerFunc(h1.Handler))
	node2.SetHandler(http.HandlerFunc(h2.Handler))

	// setup proxyd
	config := ReadConfig("fallback")
	svr, shutdown, err := proxyd.Start(config)
	require.NoError(t, err)

	// expose the proxyd client
	client := NewProxydClient("http://127.0.0.1:8545")

	// expose the backend group
	bg := svr.BackendGroups["node"]
	require.NotNil(t, bg)
	require.NotNil(t, bg.Consensus)
	require.Equal(t, 2, len(bg.Backends)) // should match config

	// convenient mapping to access the nodes by name
	nodes := map[string]nodeContext{
		"normal": {
			mockBackend: node1,
			backend:     bg.Backends[0],
			handler:     &h1,
		},
		"fallback": {
			mockBackend: node2,
			backend:     bg.Backends[1],
			handler:     &h2,
		},
	}

	return nodes, bg, client, shutdown
}

func TestFallback(t *testing.T) {
	nodes, bg, client, shutdown := setup_failover(t)
	defer nodes["normal"].mockBackend.Close()
	defer nodes["fallback"].mockBackend.Close()
	defer shutdown()

	ctx := context.Background()

	// Use Update to Advance the Candidate iteration
	update := func() {
		for _, be := range bg.Backends {
			bg.Consensus.UpdateBackend(ctx, be)
		}
		bg.Consensus.UpdateBackendGroupConsensus(ctx)
	}

	override := func(node string, method string, block string, response string) {
		if _, ok := nodes[node]; !ok {
			t.Fatalf("node %s does not exist in the nodes map", node)
		}
		nodes[node].handler.AddOverride(&ms.MethodTemplate{
			Method:   method,
			Block:    block,
			Response: response,
		})
	}

	overrideBlock := func(node string, blockRequest string, blockResponse string) {
		override(node,
			"eth_getBlockByNumber",
			blockRequest,
			buildResponse(map[string]string{
				"number": blockResponse,
				"hash":   "hash_" + blockResponse,
			}))
	}

	overrideBlockHash := func(node string, blockRequest string, number string, hash string) {
		override(node,
			"eth_getBlockByNumber",
			blockRequest,
			buildResponse(map[string]string{
				"number": number,
				"hash":   hash,
			}))
	}

	overridePeerCount := func(node string, count int) {
		override(node, "net_peerCount", "", buildResponse(hexutil.Uint64(count).String()))
	}

	overrideNotInSync := func(node string) {
		override(node, "eth_syncing", "", buildResponse(map[string]string{
			"startingblock": "0x0",
			"currentblock":  "0x0",
			"highestblock":  "0x100",
		}))
	}

	// force ban node2 and make sure node1 is the only one in consensus
	useOnlyFallback := func() {
		overridePeerCount("normal", 0)
		update()

		consensusGroup := bg.Consensus.GetConsensusGroup()
		require.Equal(t, 1, len(consensusGroup))
		require.Contains(t, consensusGroup, nodes["fallback"].backend)
		require.NotContains(t, consensusGroup, nodes["normal"].backend)
		nodes["fallback"].mockBackend.Reset()
	}

	containsFallbackNode := func(backends []*proxyd.Backend) bool {
		for _, be := range backends {
			// Note: Currently checks for name but would like to expose fallback better
			if be.Name == "fallback" {
				return true
			}
		}
		return false
	}

	recordLastUpdates := func(backends []*proxyd.Backend) []time.Time {
		lastUpdated := []time.Time{}
		for _, be := range backends {
			lastUpdated = append(lastUpdated, bg.Consensus.GetLastUpdate(be))
		}
		return lastUpdated
	}

	// convenient methods to manipulate state and mock responses
	reset := func() {
		for _, node := range nodes {
			node.handler.ResetOverrides()
			node.mockBackend.Reset()
		}
		bg.Consensus.ClearListeners()
		bg.Consensus.Reset()
		// 	// Require starting without a fallback node, and fallback is false
		// 	require.Equal(t, false, containsFallbackNode(bg.Consensus.GetConsensusGroup()))
		// 	require.Equal(t, false, bg.Consensus.GetFallbackMode())

		// 	consensusGroup := bg.Consensus.GetConsensusGroup()
		// 	require.Equal(t, "fallback", nodes["fallback"].backend.Name)
		// 	require.Equal(t, "normal", nodes["normal"].backend.Name)
		// 	require.Contains(t, consensusGroup, nodes["normal"].backend)
		// 	require.NotContains(t, consensusGroup, nodes["fallback"].backend)
		// 	// Not sure if I need these
		// 	nodes["failover"].mockBackend.Reset()
		// 	nodes["normal"].mockBackend.Reset()
	}

	// NOTE: Use debug test above the Run call to step through calls

	// TODO: Rename to with no fallback mode, no fallback backends should be returned
	t.Run("initial", func(t *testing.T) {
		reset()
		update()
		require.False(t, bg.Consensus.GetFallbackMode())
		require.Equal(t, 1, len(bg.Consensus.GetConsensusGroup()))
		require.False(t, containsFallbackNode(bg.Consensus.GetConsensusGroup()))

		// Check the backends in the Consensus Group to verify if fallback was turned on
		_, statusCode, err := client.SendRPC("eth_getBlockByNumber", []interface{}{"0x101", false})

		// TODO: Delete these later consensus at block 0x101
		require.Equal(t, 200, statusCode)
		require.Nil(t, err, "error not nil")
		require.Equal(t, "0x101", bg.Consensus.GetLatestBlockNumber().String())
		require.Equal(t, "0xe1", bg.Consensus.GetSafeBlockNumber().String())
		require.Equal(t, "0xc1", bg.Consensus.GetFinalizedBlockNumber().String())
		// TODO: Remove these, just here so compiler doesn't complain
		overridePeerCount("fallback", 0)
		overrideNotInSync("normal")
		overrideBlock("normal", "safe", "0xb1")
		overrideBlockHash("fallback", "0x102", "0x102", "wrong_hash")
		// overrideBlock("node1")
	})

	t.Run("trigger single node failing once", func(t *testing.T) {
		reset()
		useOnlyFallback()
		update()
		// require.Equal(t, false, bg.Consensus.GetFallbackMode())
		require.Equal(t, 1, len(bg.Consensus.GetConsensusGroup()))
		require.True(t, containsFallbackNode(bg.Consensus.GetConsensusGroup()))
	})

	t.Run("trigger single node failing continously", func(t *testing.T) {
		reset()
		useOnlyFallback()

		for i := 0; i < 10; i++ {
			update()
			// require.Equal(t, false, bg.Consensus.GetFallbackMode())
			require.Equal(t, 1, len(bg.Consensus.GetConsensusGroup()))
			require.True(t, containsFallbackNode(bg.Consensus.GetConsensusGroup()))
		}
	})

	t.Run("trigger single node failing continously", func(t *testing.T) {
		reset()
		useOnlyFallback()
		for i := 0; i < 10; i++ {
			update()
			// require.Equal(t, false, bg.Consensus.GetFallbackMode())
			require.Equal(t, 1, len(bg.Consensus.GetConsensusGroup()))
			require.True(t, containsFallbackNode(bg.Consensus.GetConsensusGroup()))
		}
	})

	t.Run("trigger healthy -> fallback -> healthy", func(t *testing.T) {
		reset()
		update()
		require.Equal(t, 1, len(bg.Consensus.GetConsensusGroup()))
		require.False(t, containsFallbackNode(bg.Consensus.GetConsensusGroup()))

		useOnlyFallback()
		update()
		require.Equal(t, 1, len(bg.Consensus.GetConsensusGroup()))
		require.True(t, containsFallbackNode(bg.Consensus.GetConsensusGroup()))

		reset()
		update()
		require.Equal(t, 1, len(bg.Consensus.GetConsensusGroup()))
		require.False(t, containsFallbackNode(bg.Consensus.GetConsensusGroup()))
	})
	t.Run("Ensure fallback is excluded from consensus", func(t *testing.T) {
		reset()
		// first poll
		update()
		first_update := recordLastUpdates(bg.Backends)
		useOnlyFallback()
		update()
		second_update := recordLastUpdates(bg.Backends)

		// consensus at block 0x101
		require.Equal(t, first_update[0], second_update[0])
		require.NotEqual(t, first_update[1], second_update[1])

		// as a default we use:
		// - latest at 0x101 [257]
		// - safe at 0xe1 [225]
		// - finalized at 0xc1 [193]

		// consensus at block 0x101
		require.Equal(t, "0x101", bg.Consensus.GetLatestBlockNumber().String())
		require.Equal(t, "0xe1", bg.Consensus.GetSafeBlockNumber().String())
		require.Equal(t, "0xc1", bg.Consensus.GetFinalizedBlockNumber().String())
	})
}
