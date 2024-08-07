package db

import (
	"context"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-proposer/proposer/db/ent"
	"github.com/ethereum-optimism/optimism/op-proposer/proposer/db/ent/proofrequest"

	_ "github.com/mattn/go-sqlite3"
)

type ProofDB struct {
	client *ent.Client
}

func InitDB(dbPath string) (*ProofDB, error) {
	client, err := ent.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %v", err)
	}

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %v", err)
	}

	return &ProofDB{client: client}, nil
}

func (db *ProofDB) CloseDB() error {
	if db.client != nil {
		if err := db.client.Close(); err != nil {
			return fmt.Errorf("error closing database: %w", err)
		}
	}
	return nil
}

func (db *ProofDB) NewEntry(proofType string, start, end uint64) error {
	// Convert string to proofrequest.Type
	var pType proofrequest.Type
	switch proofType {
	case "SPAN":
		pType = proofrequest.TypeSPAN
	case "AGG":
		pType = proofrequest.TypeAGG
	default:
		return fmt.Errorf("invalid proof type: %s", proofType)
	}

	_, err := db.client.ProofRequest.
		Create().
		SetType(pType).
		SetStartBlock(start).
		SetEndBlock(end).
		SetStatus(proofrequest.StatusUNREQ).
		Save(context.Background())

	if err != nil {
		return fmt.Errorf("failed to create new entry: %w", err)
	}

	return nil
}

func (db *ProofDB) UpdateProofStatus(proverRequestID string, newStatus string) error {
	// Convert string to proofrequest.Type
	var pStatus proofrequest.Status
	switch newStatus {
	case "UNREQ":
		pStatus = proofrequest.StatusUNREQ
	case "REQ":
		pStatus = proofrequest.StatusREQ
	case "COMPLETE":
		pStatus = proofrequest.StatusCOMPLETE
	case "FAILED":
		pStatus = proofrequest.StatusFAILED
	default:
		return fmt.Errorf("invalid proof status: %s", newStatus)
	}

	_, err := db.client.ProofRequest.Update().
		Where(proofrequest.ProverRequestID(proverRequestID)).
		SetStatus(pStatus).
		Save(context.Background())
	return err
}

func (db *ProofDB) AddL1BlockInfoToAggRequest(startBlock, endBlock, l1BlockNumber uint64, l1BlockHash string) error {
	_, err := db.client.ProofRequest.Update().
		Where(
			proofrequest.TypeEQ(proofrequest.TypeAGG),
			proofrequest.StatusEQ(proofrequest.StatusUNREQ),
			proofrequest.StartBlockEQ(startBlock),
			proofrequest.EndBlockEQ(endBlock),
		).
		SetL1BlockNumber(l1BlockNumber).
		SetL1BlockHash(l1BlockHash).
		Save(context.Background())

	if err != nil {
		return fmt.Errorf("failed to update L1 block info: %w", err)
	}

	return nil
}

func (db *ProofDB) GetLatestEndRequested() (uint64, error) {
	maxEnd, err := db.client.ProofRequest.Query().
		Order(ent.Desc(proofrequest.FieldEndBlock)).
		FirstID(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get latest end requested: %w", err)
	}
	return uint64(maxEnd), nil
}

func (db *ProofDB) GetAllRequestedProofs() ([]*ent.ProofRequest, error) {
	return db.client.ProofRequest.Query().
		Where(proofrequest.StatusEQ(proofrequest.StatusREQ)).
		All(context.Background())
}

func (db *ProofDB) GetAllUnrequestedProofs() ([]*ent.ProofRequest, error) {
	return db.client.ProofRequest.Query().
		Where(proofrequest.StatusEQ(proofrequest.StatusUNREQ)).
		All(context.Background())
}

func (db *ProofDB) GetCompletedAggProofs(startBlock uint64) (*ent.ProofRequest, error) {
	proof, err := db.client.ProofRequest.Query().
		Where(
			proofrequest.TypeEQ(proofrequest.TypeAGG),
			proofrequest.StartBlockEQ(startBlock),
			proofrequest.StatusEQ(proofrequest.StatusCOMPLETE),
		).
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query completed AGG proof: %w", err)
	}

	return proof, nil
}

func (db *ProofDB) TryCreateAggProofFromSpanProofs(latestOutputIndex, nextOutputIndex uint64) (bool, error) {
	// Start a transaction
	tx, err := db.client.Tx(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if there's already an AGG proof in progress
	// ZTODO: Think about off by ones
	count, err := tx.ProofRequest.Query().
		Where(
			proofrequest.TypeEQ(proofrequest.TypeAGG),
			proofrequest.StartBlockEQ(latestOutputIndex),
			proofrequest.StatusNEQ(proofrequest.StatusFAILED),
		).
		Count(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to check existing AGG proofs: %w", err)
	}
	if count > 0 {
		return false, nil // There's already an AGG proof in progress
	}

	// Find consecutive SPAN proofs
	start := latestOutputIndex
	var end uint64
	for {
		spanProof, err := tx.ProofRequest.Query().
			Where(
				proofrequest.TypeEQ(proofrequest.TypeSPAN),
				proofrequest.StatusEQ(proofrequest.StatusCOMPLETE),
				proofrequest.StartBlockEQ(start),
			).
			First(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				break // No more consecutive SPAN proofs
			}
			return false, fmt.Errorf("failed to query SPAN proof: %w", err)
		}
		end = spanProof.EndBlock
		start = end + 1
	}

	if end < nextOutputIndex {
		return false, nil // Not enough SPAN proofs to create an AGG proof
	}

	// Create a new AGG proof request
	_, err = tx.ProofRequest.Create().
		SetType(proofrequest.TypeAGG).
		SetStartBlock(latestOutputIndex).
		SetEndBlock(end).
		SetStatus(proofrequest.StatusUNREQ).
		Save(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to insert AGG proof request: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return true, nil
}
