package proposer

import (
	"context"
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
		// ZTODO: communicate with server to get this response
		// it'll come back as json
		// parse into status & proof & proof req time
		switch proverNetworkResp := "SUCCESS"; proverNetworkResp {
		case "SUCCESS":
			// get the completed proof from the network
			proof := []byte("proof")

			// update the proof to the DB and update status to "COMPLETE"
			err = l.db.AddProof(req.ID, proof)
			if err != nil {
				l.Log.Error("failed to update completed proof status", "err", err)
				return err
			}

		// ZTODO: insert timeout logic using l.DriverSetup.Cfg.MaxProofTime.
		// this needs to be adapted so we have all requested proofs included those without proverRequestID
		// then we can accurately see if they have timed out in native mode
		case "FAILED", "TIMEOUT":
			// update status in db to "FAILED"
			err = l.db.UpdateProofStatus(req.ID, "FAILED")
			if err != nil {
				l.Log.Error("failed to update failed proof status", "err", err)
				return err
			}

			// If an AGG proof failed, we're in trouble.
			// Try again.
			if req.Type == proofrequest.TypeAGG {
				l.Log.Error("failed to get agg proof, adding to db to retry", "req", req)

				err = l.db.NewEntry("AGG", req.StartBlock, req.EndBlock)
				if err != nil {
					l.Log.Error("failed to add new proof request", "err")
					return err
				}
			}

			// If a SPAN proof failed, assume it was too big.
			// Therefore, create two new entries for the original proof split in half.
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
				l.Log.Error("failed to update proof status", "err", err)
				return
			}

			err = l.RequestKonaProof(p)
			if err != nil {
				err = l.db.UpdateProofStatus(proof.ID, "FAILED")
				if err != nil {
					l.Log.Error("failed to revert proof status", "err", err, "proverRequestID", proof.ID)
				}
				l.Log.Error("failed to request proof from Kona SP1", "err", err, "proof", p)
			}
		}(proof)
	}

	return nil
}

// Use the L2OO contract to look up the range of blocks that the next proof must cover.
// Check the DB to see if we have sufficient span proofs to request an agg proof that covers this range.
// If so, queue up the agg proof in the DB to be requested later.
func (l *L2OutputSubmitter) DeriveAggProofs(ctx context.Context) error {
	latest, err := l.l2ooContract.LatestOutputIndex(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get latest L2OO output: %w", err)
	}
	from := latest.Uint64() + 1

	minTo, err := l.l2ooContract.NextOutputIndex(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get next L2OO output: %w", err)
	}

	_, err = l.db.TryCreateAggProofFromSpanProofs(from, minTo.Uint64())
	if err != nil {
		return fmt.Errorf("failed to create agg proof from span proofs: %w", err)
	}

	return nil
}

func (l *L2OutputSubmitter) RequestKonaProof(p ent.ProofRequest) error {
	prevConfirmedBlock := p.StartBlock - 1
	var proofId string
	var err error

	if p.Type == proofrequest.TypeAGG {
		subproofs, err := l.db.GetSubproofs(p.StartBlock, p.EndBlock)
		if err != nil {
			return fmt.Errorf("failed to get subproofs: %w", err)
		}

		proofId, err = l.RequestAggProof(prevConfirmedBlock, p.EndBlock, subproofs)
		if err != nil {
			return fmt.Errorf("failed to request AGG proof: %w", err)
		}
	} else if p.Type == proofrequest.TypeSPAN {
		proofId, err = l.RequestSpanProof(prevConfirmedBlock, p.EndBlock)
		if err != nil {
			return fmt.Errorf("failed to request SPAN proof: %w", err)
		}
	} else {
		return fmt.Errorf("unknown proof type: %d", p.Type)
	}

	err = l.db.SetProverRequestID(p.ID, proofId)
	if err != nil {
		return fmt.Errorf("failed to set prover request ID: %w", err)
	}

	return nil
}

func (l *L2OutputSubmitter) RequestSpanProof(start, end uint64) (string, error) {
	// use l.DriverSetup.Cfg.KonaServerURL
	return "", nil
}

func (l *L2OutputSubmitter) RequestAggProof(start, end uint64, subproofs []*ent.ProofRequest) (string, error) {
	// use l.DriverSetup.Cfg.KonaServerURL
	return "", nil
}
