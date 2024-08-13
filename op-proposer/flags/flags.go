package flags

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"

	opservice "github.com/ethereum-optimism/optimism/op-service"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
	"github.com/ethereum-optimism/optimism/op-service/oppprof"
	oprpc "github.com/ethereum-optimism/optimism/op-service/rpc"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
)

const EnvVarPrefix = "OP_PROPOSER"

func prefixEnvVars(name string) []string {
	return opservice.PrefixEnvVar(EnvVarPrefix, name)
}

var (
	// Required Flags
	L1EthRpcFlag = &cli.StringFlag{
		Name:    "l1-eth-rpc",
		Usage:   "HTTP provider URL for L1",
		EnvVars: prefixEnvVars("L1_ETH_RPC"),
	}
	RollupRpcFlag = &cli.StringFlag{
		Name:    "rollup-rpc",
		Usage:   "HTTP provider URL for the rollup node. A comma-separated list enables the active rollup provider.",
		EnvVars: prefixEnvVars("ROLLUP_RPC"),
	}
	BeaconRpcFlag = &cli.StringFlag{
		Name:    "beacon-rpc",
		Usage:   "HTTP provider URL for the beacon node",
		EnvVars: prefixEnvVars("BEACON_RPC"),
	}
	L2ChainIDFlag = &cli.Uint64Flag{
		Name:    "l2-chain-id",
		Usage:   "Chain ID of the L2 chain",
		EnvVars: prefixEnvVars("L2_CHAIN_ID"),
	}

	// Optional flags
	L2OOAddressFlag = &cli.StringFlag{
		Name:    "l2oo-address",
		Usage:   "Address of the L2OutputOracle contract",
		EnvVars: prefixEnvVars("L2OO_ADDRESS"),
	}
	PollIntervalFlag = &cli.DurationFlag{
		Name:    "poll-interval",
		Usage:   "How frequently to poll L2 for new blocks (legacy L2OO)",
		Value:   12 * time.Second,
		EnvVars: prefixEnvVars("POLL_INTERVAL"),
	}
	AllowNonFinalizedFlag = &cli.BoolFlag{
		Name:    "allow-non-finalized",
		Usage:   "Allow the proposer to submit proposals for L2 blocks derived from non-finalized L1 blocks.",
		EnvVars: prefixEnvVars("ALLOW_NON_FINALIZED"),
	}
	DisputeGameFactoryAddressFlag = &cli.StringFlag{
		Name:    "game-factory-address",
		Usage:   "Address of the DisputeGameFactory contract",
		EnvVars: prefixEnvVars("GAME_FACTORY_ADDRESS"),
	}
	ProposalIntervalFlag = &cli.DurationFlag{
		Name:    "proposal-interval",
		Usage:   "Interval between submitting L2 output proposals when the dispute game factory address is set",
		EnvVars: prefixEnvVars("PROPOSAL_INTERVAL"),
	}
	OutputRetryIntervalFlag = &cli.DurationFlag{
		Name:    "output-retry-interval",
		Usage:   "Interval between retrying output fetching (DGF)",
		Value:   12 * time.Second,
		EnvVars: prefixEnvVars("OUTPUT_RETRY_INTERVAL"),
	}
	DisputeGameTypeFlag = &cli.UintFlag{
		Name:    "game-type",
		Usage:   "Dispute game type to create via the configured DisputeGameFactory",
		Value:   0,
		EnvVars: prefixEnvVars("GAME_TYPE"),
	}
	ActiveSequencerCheckDurationFlag = &cli.DurationFlag{
		Name:    "active-sequencer-check-duration",
		Usage:   "The duration between checks to determine the active sequencer endpoint.",
		Value:   2 * time.Minute,
		EnvVars: prefixEnvVars("ACTIVE_SEQUENCER_CHECK_DURATION"),
	}
	WaitNodeSyncFlag = &cli.BoolFlag{
		Name: "wait-node-sync",
		Usage: "Indicates if, during startup, the proposer should wait for the rollup node to sync to " +
			"the current L1 tip before proceeding with its driver loop.",
		Value:   false,
		EnvVars: prefixEnvVars("WAIT_NODE_SYNC"),
	}
	DbPathFlag = &cli.StringFlag{
		Name:  "db-path",
		Usage: "Path to the database used to track ZK proof generation",
		// ZTODO: Decide on better default path here
		Value:   "./proofs.db",
		EnvVars: prefixEnvVars("DB_PATH"),
	}
	MaxSpanBatchDeviationFlag = &cli.Uint64Flag{
		Name:    "max-span-batch-deviation",
		Usage:   "If we find a span batch this far ahead of our target, we assume an error and fill in the gap",
		Value:   600,
		EnvVars: prefixEnvVars("MAX_SPAN_BATCH_DEVIATION"),
	}
	MaxBlockRangePerSpanProofFlag = &cli.Uint64Flag{
		Name:    "max-block-range-per-span-proof",
		Usage:   "Maximum number of blocks to include in a single span proof",
		Value:   50,
		EnvVars: prefixEnvVars("MAX_BLOCK_RANGE_PER_SPAN_PROOF"),
	}
	// ZTODO: Rename to something like ProofExpiryTime
	MaxProofTimeFlag = &cli.Uint64Flag{
		Name:    "max-proof-time",
		Usage:   "Maximum time in seconds to spend generating a proof before giving up",
		Value:   14400,
		EnvVars: prefixEnvVars("MAX_PROOF_TIME"),
	}
	KonaServerUrlFlag = &cli.StringFlag{
		Name:    "kona-server-url",
		Usage:   "URL of the Kona server to request proofs from",
		Value:   "http://127.0.0.1:3000",
		EnvVars: prefixEnvVars("KONA_SERVER_URL"),
	}
	MaxConcurrentProofRequestsFlag = &cli.Uint64Flag{
		Name:    "max-concurrent-proof-requests",
		Usage:   "Maximum number of proofs to generate concurrently",
		Value:   20,
		EnvVars: prefixEnvVars("MAX_CONCURRENT_PROOF_REQUESTS"),
	}
	TxCacheOutDirFlag = &cli.StringFlag{
		Name:    "tx-cache-out-dir",
		Usage:   "Cache directory for the found transactions to determine span batch boundaries",
		Value:   "/tmp/batch_decoder/transactions_cache",
		EnvVars: prefixEnvVars("TX_CACHE_OUT_DIR"),
	}
	BatchDecoderConcurrentReqsFlag = &cli.Uint64Flag{
		Name:    "batch-decoder-concurrent-reqs",
		Usage:   "Concurrency level when fetching transactions to determine span batch boundaries",
		Value:   10,
		EnvVars: prefixEnvVars("BATCH_DECODER_CONCURRENT_REQS"),
	}
	// Legacy Flags
	L2OutputHDPathFlag = txmgr.L2OutputHDPathFlag
)

var requiredFlags = []cli.Flag{
	L1EthRpcFlag,
	RollupRpcFlag,
	BeaconRpcFlag,
	L2ChainIDFlag,
}

var optionalFlags = []cli.Flag{
	L2OOAddressFlag,
	PollIntervalFlag,
	AllowNonFinalizedFlag,
	L2OutputHDPathFlag,
	DisputeGameFactoryAddressFlag,
	ProposalIntervalFlag,
	OutputRetryIntervalFlag,
	DisputeGameTypeFlag,
	ActiveSequencerCheckDurationFlag,
	WaitNodeSyncFlag,
	DbPathFlag,
	MaxSpanBatchDeviationFlag,
	MaxBlockRangePerSpanProofFlag,
	MaxProofTimeFlag,
	TxCacheOutDirFlag,
	BatchDecoderConcurrentReqsFlag,
	KonaServerUrlFlag,
	MaxConcurrentProofRequestsFlag,
}

func init() {
	optionalFlags = append(optionalFlags, oprpc.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, oplog.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, opmetrics.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, oppprof.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, txmgr.CLIFlags(EnvVarPrefix)...)

	Flags = append(requiredFlags, optionalFlags...)
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func CheckRequired(ctx *cli.Context) error {
	for _, f := range requiredFlags {
		if !ctx.IsSet(f.Names()[0]) {
			return fmt.Errorf("flag %s is required", f.Names()[0])
		}
	}
	return nil
}
