package proposer

import (
	"context"
	"errors"
	"fmt"

	"github.com/ava-labs/coreth/accounts/abi/bind"
	"github.com/ethereum-optimism/optimism/op-proposer/proposer/db/ent"
	"github.com/ethereum-optimism/optimism/op-proposer/proposer/db/ent/proofrequest"
)

func (l *L2OutputSubmitter) ProcessPendingProofs() error {
	reqs, err := l.db.GetAllPendingProofs()
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

		// ZTODO: insert timeout logic using l.DriverSetup.Cfg.MaxProofTime.
		case "FAILED", "TIMEOUT":
			// update status in db to "FAILED"
			err = l.db.UpdateProofStatus(req.ProverRequestID, "FAILED")
			if err != nil {
				l.Log.Error("failed to update failed proof status", "err", err)
				return err
			}

			if req.Type == proofrequest.TypeAGG {
				l.Log.Error("failed to get agg proof", "req", req)
				return errors.New("failed to get agg proof")
				// ZTODO: Should we default to trying again or will it be same result?
			}

			// add two new entries for the request split in half
			tmpStart := req.StartBlock
			tmpEnd := tmpStart + ((req.EndBlock - tmpStart) / 2)
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

	return nil
}

func (l *L2OutputSubmitter) RequestQueuedProofs(ctx context.Context) error {
	unrequestedProofs, err := l.db.GetAllUnrequestedProofs()
	if err != nil {
		return fmt.Errorf("failed to get unrequested proofs: %w", err)
	}

	for _, proof := range unrequestedProofs {
		if proof.Type == proofrequest.TypeAGG {
			blockNumber, blockHash, err := l.checkpointBlockHash(ctx)
			if err != nil {
				l.Log.Error("failed to checkpoint block hash", "err", err)
				return err
			}
			l.db.AddL1BlockInfo(proof.StartBlock, proof.EndBlock, blockNumber, blockHash)
		}
		go func(p ent.ProofRequest) {
			err = l.db.UpdateProofStatus(proof.ID, "REQ")
			if err != nil {
				l.Log.Error("failed to update proof status", "err", err, "proverRequestID", proverRequestID)
				return
			}

			proverRequestID, err := l.RequestKonaProof(p)
			if err != nil {
				err = l.db.UpdateProofStatus(proof.ID, "FAILED")
				if err != nil {
					l.Log.Error("failed to revert proof status", "err", err, "proverRequestID", proverRequestID)
				}
				l.Log.Error("failed to request proof from Kona SP1", "err", err, "proof", p)
				return
			}

			err = l.db.SetProverRequestID(proof.ID, proverRequestID)
			if err != nil {
				l.Log.Error("failed to set prover request ID", "err", err, "proverRequestID", proverRequestID
			}
		}(proof)
	}

	return nil
}

func (l *L2OutputSubmitter) DeriveAggProofs(ctx context.Context) error {
	// Get the latest L2OO output
	from, err := l.l2ooContract.LatestOutputIndex(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get latest L2OO output: %w", err)
	}

	minTo, err := l.l2ooContract.NextOutputIndex(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get next L2OO output: %w", err)
	}

	// ZTODO: Need to turn these big ints into uint64
	_, err = l.db.TryCreateAggProofFromSpanProofs(from, minTo)
	if err != nil {
		return fmt.Errorf("failed to create agg proof from span proofs: %w", err)
	}

	return nil
}

func (l *L2OutputSubmitter) RequestKonaProof(p ent.ProofRequest) (string, error) {
	// TODO:
	// - implement requestProofFromKonaSP1 function
	// - figure out how to get proverReqId back
	// - determine order of operations (can't wait too long on status but can't preempt)
}
