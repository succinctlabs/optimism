package proposer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-proposer/bindings"
	"github.com/ethereum-optimism/optimism/op-proposer/metrics"
	"github.com/ethereum-optimism/optimism/op-service/dial"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
)

var (
	supportedL2OutputVersion = eth.Bytes32{}
	ErrProposerNotRunning    = errors.New("proposer is not running")
)

type L1Client interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	// CodeAt returns the code of the given account. This is needed to differentiate
	// between contract internal errors and the local chain being out of sync.
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)

	// CallContract executes an Ethereum contract call with the specified data as the
	// input.
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

type L2OOContract interface {
	Version(*bind.CallOpts) (string, error)
	NextBlockNumber(*bind.CallOpts) (*big.Int, error)
}

type RollupClient interface {
	SyncStatus(ctx context.Context) (*eth.SyncStatus, error)
	OutputAtBlock(ctx context.Context, blockNum uint64) (*eth.OutputResponse, error)
}


type DriverSetup struct {
	Log      log.Logger
	Metr     metrics.Metricer
	Cfg      ProposerConfig
	Txmgr    txmgr.TxManager
	L1Client *ethclient.Client

	// RollupProvider's RollupClient() is used to retrieve output roots from
	RollupProvider dial.RollupProvider
}

// L2OutputSubmitter is responsible for proposing outputs
type L2OutputSubmitter struct {
	DriverSetup

	wg   sync.WaitGroup
	done chan struct{}

	ctx    context.Context
	cancel context.CancelFunc

	mutex   sync.Mutex
	running bool

	l2ooContract L2OOContract
	l2ooABI      *abi.ABI

	dgfContract *bindings.DisputeGameFactoryCaller
	dgfABI      *abi.ABI

	db ProofDB
}

// NewL2OutputSubmitter creates a new L2 Output Submitter
func NewL2OutputSubmitter(setup DriverSetup) (_ *L2OutputSubmitter, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	// The above context is long-lived, and passed to the `L2OutputSubmitter` instance. This context is closed by
	// `StopL2OutputSubmitting`, but if this function returns an error or panics, we want to ensure that the context
	// doesn't leak.
	defer func() {
		if err != nil || recover() != nil {
			cancel()
		}
	}()

	if setup.Cfg.L2OutputOracleAddr != nil {
		return newL2OOSubmitter(ctx, cancel, setup)
	} else if setup.Cfg.DisputeGameFactoryAddr != nil {
		return newDGFSubmitter(ctx, cancel, setup)
	} else {
		return nil, errors.New("neither the `L2OutputOracle` nor `DisputeGameFactory` addresses were provided")
	}
}

func newL2OOSubmitter(ctx context.Context, cancel context.CancelFunc, setup DriverSetup) (*L2OutputSubmitter, error) {
	l2ooContract, err := bindings.NewL2OutputOracleCaller(*setup.Cfg.L2OutputOracleAddr, setup.L1Client)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create L2OO at address %s: %w", setup.Cfg.L2OutputOracleAddr, err)
	}

	cCtx, cCancel := context.WithTimeout(ctx, setup.Cfg.NetworkTimeout)
	defer cCancel()
	version, err := l2ooContract.Version(&bind.CallOpts{Context: cCtx})
	if err != nil {
		cancel()
		return nil, err
	}
	log.Info("Connected to L2OutputOracle", "address", setup.Cfg.L2OutputOracleAddr, "version", version)

	parsed, err := bindings.L2OutputOracleMetaData.GetAbi()
	if err != nil {
		cancel()
		return nil, err
	}

	db, err := InitDB("./proofs.db")
	if err != nil {
		cancel()
		return nil, err
	}

	return &L2OutputSubmitter{
		DriverSetup: setup,
		done:        make(chan struct{}),
		ctx:         ctx,
		cancel:      cancel,

		l2ooContract: l2ooContract,
		l2ooABI:      parsed,
		db:           db,
	}, nil
}

func newDGFSubmitter(ctx context.Context, cancel context.CancelFunc, setup DriverSetup) (*L2OutputSubmitter, error) {
	dgfCaller, err := bindings.NewDisputeGameFactoryCaller(*setup.Cfg.DisputeGameFactoryAddr, setup.L1Client)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create DGF at address %s: %w", setup.Cfg.DisputeGameFactoryAddr, err)
	}

	cCtx, cCancel := context.WithTimeout(ctx, setup.Cfg.NetworkTimeout)
	defer cCancel()
	version, err := dgfCaller.Version(&bind.CallOpts{Context: cCtx})
	if err != nil {
		cancel()
		return nil, err
	}
	log.Info("Connected to DisputeGameFactory", "address", setup.Cfg.DisputeGameFactoryAddr, "version", version)

	parsed, err := bindings.DisputeGameFactoryMetaData.GetAbi()
	if err != nil {
		cancel()
		return nil, err
	}

	return &L2OutputSubmitter{
		DriverSetup: setup,
		done:        make(chan struct{}),
		ctx:         ctx,
		cancel:      cancel,

		dgfContract: dgfCaller,
		dgfABI:      parsed,
	}, nil
}

func (l *L2OutputSubmitter) StartL2OutputSubmitting() error {
	l.Log.Info("Starting Proposer")

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.running {
		return errors.New("proposer is already running")
	}
	l.running = true

	l.wg.Add(1)
	go l.loop()

	l.Log.Info("Proposer started")
	return nil
}

func (l *L2OutputSubmitter) StopL2OutputSubmittingIfRunning() error {
	err := l.StopL2OutputSubmitting()
	if errors.Is(err, ErrProposerNotRunning) {
		return nil
	}
	return err
}

func (l *L2OutputSubmitter) StopL2OutputSubmitting() error {
	l.Log.Info("Stopping Proposer")

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if !l.running {
		return ErrProposerNotRunning
	}
	l.running = false

	l.cancel()
	close(l.done)
	l.wg.Wait()

	if l.db != nil {
		if err := l.db.CloseDB(); err != nil {
			return fmt.Errorf("error closing database: %w", err)
		}
	}

	l.Log.Info("Proposer stopped")
	return nil
}
func (l *L2OutputSubmitter) addAvailableSpanBatchesToDB(ctx context.Context) error {
	// nextBlock is equal to the highest value in the `EndBlock` column of the db, plus 1
	// ZTODO: think through off by ones
	lastEndBlock, err := l.db.GetLatestEndRequested()
	if err != nil {
		l.Log.Error("failed to get latest end requested", "err", err)
		return err
	}
	nextBlock := lastEndBlock + 1

	// use batch decoder to pull all batches from next block's L1 Origin through Finalized L1 from chain to disk
	err := l.FetchBatchesFromChain(ctx, nextBlock)
	if err != nil {
		l.Log.Error("failed to fetch batches from chain", "err", err)
		return err
	}

	for {
		// use batch decoder to reassemble the batches from disk to determine the start and end of relevant span batch
		start, end, err := l.GenerateSpanBatchRange(ctx, nextBlock)
		if err == NoSpanBatchFoundError {
			l.Log.Info("no span batch found", "nextBlock", nextBlock)
			break
		} else if err != nil {
			l.Log.Error("failed to generate span batch range", "err", err)
			return err
		}

		// the nextBlock should always be the start of a new span batch, warn if not
		if start != nextBlock {
			l.Log.Warn("start block does not match next block", "start", start, "nextBlock", nextBlock)
		}

		// insert the new span into the db to be requested in the future
		err = l.db.NewEntry("SPAN", nextBlock, end)
		if err != nil {
			l.Log.Error("failed to insert proof request", "err", err)
			return err
		}

		// ZTODO: think through off by ones
		nextBlock = end + 1
	}

	return nil
}

func parseSpanBatchResponse(data []byte) (uint64, uint64, error) {
	parts := strings.Split(string(data), " ")
	if len(parts) != 2 {
		l.Log.Error("too many parts in span batch response", "span", spanBatchData)
		return errors.New("failed to parse span range")
	}
	start, err := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing start value: %w", err)
	}
	end, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing end value: %w", err)
	}

	return start, end, nil
}

func (l *L2OutputSubmitter) updateRequestedProofs() error {
	reqs, err := l.db.GetAllRequestedProofs()
	if err != nil {
		return err
	}
	for _, req := range reqs {
		// check prover network for req id status
		// TODO: HAVE IT PING SP1 NETWORK TO ASK FOR STATUS
		switch proverNetworkResp := "SUCCESS"; proverNetworkResp {
		case "SUCCESS":
			// get the completed proof from the network
			// ZTODO
			proof := []byte("proof")

			// update the proof to the DB and update status to "COMPLETE"
			err = l.db.AddProof(req.ProverRequestID, proof)
			if err != nil {
				l.Log.Error("failed to update completed proof status", "err", err)
				return err
			}

		case "FAILED":
			// update status in db to "FAILED"
			err = l.db.UpdateProofStatus(req.ProverRequestID, "FAILED")
			if err != nil {
				l.Log.Error("failed to update failed proof status", "err", err)
				return err
			}

			if req.Type == ProofTypeAGG {
				l.Log.Error("failed to get agg proof", "req", req)
				return errors.New("failed to get agg proof")
				// ZTODO: Should we default to trying again or will it be same result?
			}

			// add two new entries for the request split in half
			tmpStart := req.StartBlock
			tmpEnd := start + ((req.EndBlock - start) / 2)
			for i := 0; i < 2; i++ {
				err = l.db.NewEntry("SPAN", tmpStart, tmpEnd)
				if err != nil {
					l.Log.Error("failed to add new proof request", "err", err)
					return err
				}

				tmpStart = tmpEnd + 1
				tmpEnd = req.EndBlock
			}
		}
	}
}

func (l *L2OutputSubmitter) queuePendingAggProofs(ctx context.Context) error {
	// Get the latest L2OO output
	from, err := l.l2ooContract.LatestOutputIndex(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get latest L2OO output: %w", err)
	}

	minTo, err := l.l2ooContract.NextOutputIndex(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get next L2OO output: %w", err)
	}

	_, err := l.db.TryCreateAggProofFromSpanProofs(from, minTo)
	if err != nil {
		return fmt.Errorf("failed to create agg proof from span proofs: %w", err)
	}

	return nil
}

func (l *L2OutputSubmitter) requestUnrequestedProofs() error {
	unrequestedProofs, err := l.db.GetAllUnrequestedProofs()
    if err != nil {
        return fmt.Errorf("failed to get unrequested proofs: %w", err)
    }

    for _, proof := range unrequestedProofs {
		if proof.Type == ProofTypeAGG {
			blockNumber, blockHash, err := l.checkpointBlockHash(ctx)
			if err != nil {
				l.Log.Error("failed to checkpoint block hash", "err", err)
				return err
			}
			l.db.AddL1BlockInfo(proof.StartBlock, proof.EndBlock, blockNumber, blockHash)
		}
        go func(p ProofRequest) {
			// TODO:
			// - implement requestProofFromKonaSP1 function
			// - figure out how to get proverReqId back
			// - determine order of operations (can't wait too long on status but can't preempt)
            proverRequestID, err := l.requestProofFromKonaSP1(ctx, p)
            if err != nil {
                l.Log.Error("failed to request proof from Kona SP1", "err", err, "proof", p)
                return
            }

            err = l.db.UpdateProofStatus(proverRequestID, "REQ")
            if err != nil {
                l.Log.Error("failed to update proof status", "err", err, "proverRequestID", proverRequestID)
            }
        }(proof)
    }

    return nil
}

func (l *L2OutputSubmitter) submitAggProofs(ctx context.Context) error {
    // Get the latest output index from the L2OutputOracle contract
    latestOutputIndex, err := l.l2ooContract.LatestOutputIndex(&bind.CallOpts{Context: ctx})
    if err != nil {
        return fmt.Errorf("failed to get latest output index: %w", err)
    }

    // Check for a completed AGG proof starting at the next index
	// TODO: Check for off by one?
    completedAggProofs, err := l.db.GetAllCompletedAggProofs(latestOutputIndex)
    if err != nil {
        return fmt.Errorf("failed to query for completed AGG proof: %w", err)
    }

    if len(completedAggProofs) == 0 {
        return nil
    }

	for _, aggProof := range completedAggProofs {
		// TODO: Off by one?
		output, err := l.FetchOutput(ctx, aggProof.EndBlock)
		if err != nil {
			return fmt.Errorf("failed to fetch output at block %d: %w", aggProof.EndBlock, err)
		}

		l.proposeOutput(ctx, output, aggProof.proof)
		l.Log.Info("AGG proof submitted on-chain", "start", aggProof.StartBlock, "end", aggProof.EndBlock)
	}

    return nil
}

// FetchL2OOOutput gets the next output proposal for the L2OO.
// It queries the L2OO for the earliest next block number that should be proposed.
// It returns the output to propose, and whether the proposal should be submitted at all.
// The passed context is expected to be a lifecycle context. A network timeout
// context will be derived from it.
func (l *L2OutputSubmitter) FetchL2OOOutput(ctx context.Context) (*eth.OutputResponse, bool, error) {
	if l.l2ooContract == nil {
		return nil, false, fmt.Errorf("L2OutputOracle contract not set, cannot fetch next output info")
	}

	cCtx, cancel := context.WithTimeout(ctx, l.Cfg.NetworkTimeout)
	defer cancel()
	callOpts := &bind.CallOpts{
		From:    l.Txmgr.From(),
		Context: cCtx,
	}
	nextCheckpointBlockBig, err := l.l2ooContract.NextBlockNumber(callOpts)
	if err != nil {
		return nil, false, fmt.Errorf("querying next block number: %w", err)
	}
	nextCheckpointBlock := nextCheckpointBlockBig.Uint64()
	// Fetch the current L2 heads
	currentBlockNumber, err := l.FetchCurrentBlockNumber(ctx)
	if err != nil {
		return nil, false, err
	}

	// Ensure that we do not submit a block in the future
	if currentBlockNumber < nextCheckpointBlock {
		l.Log.Debug("Proposer submission interval has not elapsed", "currentBlockNumber", currentBlockNumber, "nextBlockNumber", nextCheckpointBlock)
		return nil, false, nil
	}

	output, err := l.FetchOutput(ctx, nextCheckpointBlock)
	if err != nil {
		return nil, false, fmt.Errorf("fetching output: %w", err)
	}

	// Always propose if it's part of the Finalized L2 chain. Or if allowed, if it's part of the safe L2 chain.
	if output.BlockRef.Number > output.Status.FinalizedL2.Number && (!l.Cfg.AllowNonFinalized || output.BlockRef.Number > output.Status.SafeL2.Number) {
		l.Log.Debug("Not proposing yet, L2 block is not ready for proposal",
			"l2_proposal", output.BlockRef,
			"l2_safe", output.Status.SafeL2,
			"l2_finalized", output.Status.FinalizedL2,
			"allow_non_finalized", l.Cfg.AllowNonFinalized)
		return output, false, nil
	}
	return output, true, nil
}

// FetchDGFOutput gets the next output proposal for the DGF.
// The passed context is expected to be a lifecycle context. A network timeout
// context will be derived from it.
func (l *L2OutputSubmitter) FetchDGFOutput(ctx context.Context) (*eth.OutputResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, l.Cfg.NetworkTimeout)
	defer cancel()

	blockNum, err := l.FetchCurrentBlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	return l.FetchOutput(ctx, blockNum)
}

// FetchCurrentBlockNumber gets the current block number from the [L2OutputSubmitter]'s [RollupClient]. If the `AllowNonFinalized` configuration
// option is set, it will return the safe head block number, and if not, it will return the finalized head block number.
func (l *L2OutputSubmitter) FetchCurrentBlockNumber(ctx context.Context) (uint64, error) {
	rollupClient, err := l.RollupProvider.RollupClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("getting rollup client: %w", err)
	}

	status, err := rollupClient.SyncStatus(ctx)
	if err != nil {
		return 0, fmt.Errorf("getting sync status: %w", err)
	}

	// Use either the finalized or safe head depending on the config. Finalized head is default & safer.
	if l.Cfg.AllowNonFinalized {
		return status.SafeL2.Number, nil
	}
	return status.FinalizedL2.Number, nil
}

func (l *L2OutputSubmitter) FetchOutput(ctx context.Context, block uint64) (*eth.OutputResponse, error) {
	rollupClient, err := l.RollupProvider.RollupClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting rollup client: %w", err)
	}

	output, err := rollupClient.OutputAtBlock(ctx, block)
	if err != nil {
		return nil, fmt.Errorf("fetching output at block %d: %w", block, err)
	}
	if output.Version != supportedL2OutputVersion {
		return nil, fmt.Errorf("unsupported l2 output version: %v, supported: %v", output.Version, supportedL2OutputVersion)
	}
	if onum := output.BlockRef.Number; onum != block { // sanity check, e.g. in case of bad RPC caching
		return nil, fmt.Errorf("output block number %d mismatches requested %d", output.BlockRef.Number, block)
	}
	return output, nil
}

// ProposeL2OutputTxData creates the transaction data for the ProposeL2Output function
func (l *L2OutputSubmitter) ProposeL2OutputTxData(output *eth.OutputResponse, proof []byte) ([]byte, error) {
	return proposeL2OutputTxData(l.l2ooABI, output, proof)
}

// proposeL2OutputTxData creates the transaction data for the ProposeL2Output function
func proposeL2OutputTxData(abi *abi.ABI, output *eth.OutputResponse, proof []byte) ([]byte, error) {
	return abi.Pack(
		"proposeL2Output",
		output.OutputRoot,
		new(big.Int).SetUint64(output.BlockRef.Number),
		output.Status.CurrentL1.Hash,
		new(big.Int).SetUint64(output.Status.CurrentL1.Number),
		proof)
}

func (l *L2OutputSubmitter) ProposeL2OutputDGFTxData(output *eth.OutputResponse) ([]byte, *big.Int, error) {
	bond, err := l.dgfContract.InitBonds(&bind.CallOpts{}, l.Cfg.DisputeGameType)
	if err != nil {
		return nil, nil, err
	}
	data, err := proposeL2OutputDGFTxData(l.dgfABI, l.Cfg.DisputeGameType, output)
	if err != nil {
		return nil, nil, err
	}
	return data, bond, err
}

// proposeL2OutputDGFTxData creates the transaction data for the DisputeGameFactory's `create` function
func proposeL2OutputDGFTxData(abi *abi.ABI, gameType uint32, output *eth.OutputResponse) ([]byte, error) {
	return abi.Pack("create", gameType, output.OutputRoot, math.U256Bytes(new(big.Int).SetUint64(output.BlockRef.Number)))
}

// We wait until l1head advances beyond blocknum. This is used to make sure proposal tx won't
// immediately fail when checking the l1 blockhash. Note that EstimateGas uses "latest" state to
// execute the transaction by default, meaning inside the call, the head block is considered
// "pending" instead of committed. In the case l1blocknum == l1head then, blockhash(l1blocknum)
// will produce a value of 0 within EstimateGas, and the call will fail when the contract checks
// that l1blockhash matches blockhash(l1blocknum).
func (l *L2OutputSubmitter) waitForL1Head(ctx context.Context, blockNum uint64) error {
	ticker := time.NewTicker(l.Cfg.PollInterval)
	defer ticker.Stop()
	l1head, err := l.Txmgr.BlockNumber(ctx)
	if err != nil {
		return err
	}
	for l1head <= blockNum {
		l.Log.Debug("Waiting for l1 head > l1blocknum1+1", "l1head", l1head, "l1blocknum", blockNum)
		select {
		case <-ticker.C:
			l1head, err = l.Txmgr.BlockNumber(ctx)
			if err != nil {
				return err
			}
		case <-l.done:
			return fmt.Errorf("L2OutputSubmitter is done()")
		}
	}
	return nil
}

// sendTransaction creates & sends transactions through the underlying transaction manager.
func (l *L2OutputSubmitter) sendTransaction(ctx context.Context, output *eth.OutputResponse, proof []byte) error {
	err := l.waitForL1Head(ctx, output.Status.HeadL1.Number+1)
	if err != nil {
		return err
	}

	l.Log.Info("Proposing output root", "output", output.OutputRoot, "block", output.BlockRef)
	var receipt *types.Receipt
	if l.Cfg.DisputeGameFactoryAddr != nil {
		data, bond, err := l.ProposeL2OutputDGFTxData(output)
		if err != nil {
			return err
		}
		receipt, err = l.Txmgr.Send(ctx, txmgr.TxCandidate{
			TxData:   data,
			To:       l.Cfg.DisputeGameFactoryAddr,
			GasLimit: 0,
			Value:    bond,
		})
		if err != nil {
			return err
		}
	} else {
		data, err := l.ProposeL2OutputTxData(output, proof)
		if err != nil {
			return err
		}
		receipt, err = l.Txmgr.Send(ctx, txmgr.TxCandidate{
			TxData:   data,
			To:       l.Cfg.L2OutputOracleAddr,
			GasLimit: 0,
		})
		if err != nil {
			return err
		}
	}

	if receipt.Status == types.ReceiptStatusFailed {
		l.Log.Error("Proposer tx successfully published but reverted", "tx_hash", receipt.TxHash)
	} else {
		l.Log.Info("Proposer tx successfully published",
			"tx_hash", receipt.TxHash,
			"l1blocknum", output.Status.CurrentL1.Number,
			"l1blockhash", output.Status.CurrentL1.Hash)
	}
	return nil
}


// sendCheckpointTransaction creates & sends transaction to checkpoint blockhash on L2OO contract.
func (l *L2OutputSubmitter) sendCheckpointTransaction(ctx context.Context, blockNumber uint64, blockHash common.Hash) (uint64, common.Hash, error) {
	var receipt *types.Receipt
	data, err := l.CheckpointBlockHashTxData(blockNumber, blockHash)
	if err != nil {
		return 0, common.Hash{}, err
	}
	receipt, err = l.Txmgr.Send(ctx, txmgr.TxCandidate{
		TxData:   data,
		To:       l.Cfg.L2OutputOracleAddr,
		GasLimit: 0,
	})
	if err != nil {
		return 0, common.Hash{}, err
	}

	if receipt.Status == types.ReceiptStatusFailed {
		l.Log.Error("checkpoint blockhash tx successfully published but reverted", "tx_hash", receipt.TxHash)
	} else {
		l.Log.Info("checkpoint blockhash tx successfully published",
			"tx_hash", receipt.TxHash)
	}
	return blockNumber, blockHash, nil
}

// loop is responsible for creating & submitting the next outputs
func (l *L2OutputSubmitter) loop() {
	defer l.wg.Done()
	ctx := l.ctx

	if l.Cfg.WaitNodeSync {
		err := l.waitNodeSync()
		if err != nil {
			l.Log.Error("Error waiting for node sync", "err", err)
			return
		}
	}

	if l.dgfContract == nil {
		l.loopL2OO(ctx)
	} else {
		l.loopDGF(ctx)
	}
}

func (l *L2OutputSubmitter) waitNodeSync() error {
	cCtx, cancel := context.WithTimeout(l.ctx, l.Cfg.NetworkTimeout)
	defer cancel()

	l1head, err := l.Txmgr.BlockNumber(cCtx)
	if err != nil {
		return fmt.Errorf("failed to retrieve current L1 block number: %w", err)
	}

	rollupClient, err := l.RollupProvider.RollupClient(l.ctx)
	if err != nil {
		return fmt.Errorf("failed to get rollup client: %w", err)
	}

	return dial.WaitRollupSync(l.ctx, l.Log, rollupClient, l1head, time.Second*12)
}

// The loopL2OO regularly polls the L2OO for the next block to propose,
// and if the current finalized (or safe) block is past that next block, it
// proposes it.
func (l *L2OutputSubmitter) loopL2OO(ctx context.Context) {
	ticker := time.NewTicker(l.Cfg.PollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 1) Queue up any span batches that are ready to prove.
			// This is done by checking the chain for completed channels and pulling span batches out.
			// We add them to the DB as "UNREQ" so they are queued up to request later.
			err := l.addAvailableSpanBatchesToDB(ctx)
			if err != nil {
				l.Log.Error("failed to add next span batches to db", "err", err)
				break
			}

			// 2) Check the statuses of all requested proofs.
			// If it's successfully returned, we validate that we have it on disk and set status = "COMPLETE".
			// If it fails, we set status = "FAILED" (and, if it's a span proof, split the request in half to try again).
			err = l.updateRequestedProofs()
			if err != nil {
				l.Log.Error("failed to update requested proofs", "err", err)
				break
			}

			// 3) Determine if any agg proofs are ready to be submitted and queue them up.
			// This is done by checking if we have contiguous span proofs from the last on chain
			// output root through at least the submission interval.
			err = l.queuePendingAggProofs()
			if err != nil {
				l.Log.Error("failed to generate pending agg proofs", "err", err)
				break
			}

			// 4) Request all unrequested proofs from the prover network.
			// Any DB entry with status = "UNREQ" means it's queued up and ready.
			// We request all of these (both span and agg) from the prover network.
			// For agg proofs, we also checkpoint the blockhash in advance.
			err = l.requestUnrequestedProofs()
			if err != nil {
				l.Log.Error("failed to request unrequested proofs", "err", err)
				break
			}

			// 5) Submit agg proofs on chain.
			err = l.submitAggProofs(ctx)
		case <-l.done:
			return
		}
	}
}

// The loopDGF proposes a new output every proposal interval. It does _not_ query
// the DGF for when to next propose, as the DGF doesn't have the concept of a
// proposal interval, like in the L2OO case. For this reason, it has to keep track
// of the interval itself, for which it uses an internal ticker.
func (l *L2OutputSubmitter) loopDGF(ctx context.Context) {
	defer l.Log.Info("loopDGF returning")
	ticker := time.NewTicker(l.Cfg.ProposalInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			var (
				output *eth.OutputResponse
				err    error
			)
			// A note on retrying: because the proposal interval is usually much
			// larger than the interval at which to retry proposing on a failed attempt,
			// we want to keep retrying getting the output proposal until we succeed.
			for output == nil || err != nil {
				select {
				case <-l.done:
					return
				default:
				}

				output, err = l.FetchDGFOutput(ctx)
				if err != nil {
					l.Log.Warn("Error getting DGF output, retrying...", "err", err)
					time.Sleep(l.Cfg.OutputRetryInterval)
				}
			}

			l.proposeOutput(ctx, output, nil)
		case <-l.done:
			return
		}
	}
}

func (l *L2OutputSubmitter) proposeOutput(ctx context.Context, output *eth.OutputResponse, proof []byte) {
	cCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	if err := l.sendTransaction(cCtx, output, proof); err != nil {
		l.Log.Error("Failed to send proposal transaction",
			"err", err,
			"l1blocknum", output.Status.CurrentL1.Number,
			"l1blockhash", output.Status.CurrentL1.Hash,
			"l1head", output.Status.HeadL1.Number
			"proof", proof)
		return
	}
	l.Metr.RecordL2BlocksProposed(output.BlockRef)
}

func (l *L2OutputSubmitter) checkpointBlockHash(ctx context.Context) (uint64, common.Hash, error) {
	cCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	blockNumber, err := l.Txmgr.BlockNumber(cCtx)
	if err != nil {
		return 0, common.Hash{}, err
	}
	header, err := l.Txmgr.BlockHeader(cCtx)
	if err != nil {
		return 0, common.Hash{}, err
	}
	blockHash := header.Hash()

	return l.sendCheckpointTransaction(cCtx, blockNumber, blockHash)
}
