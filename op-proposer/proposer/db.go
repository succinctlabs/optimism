package proposer

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type ProofType string

const (
	ProofTypeSPAN ProofType = "SPAN"
	ProofTypeAGG  ProofType = "AGG"
)

type ProofStatus string

const (
	ProofStatusUNREQ    ProofStatus = "UNREQ"
	ProofStatusREQ      ProofStatus = "REQ"
	ProofStatusFAILED   ProofStatus = "FAILED"
	ProofStatusCOMPLETE ProofStatus = "COMPLETE"
)

type ProofRequest struct {
	Type             ProofType
	StartBlock       uint64
	EndBlock         uint64
	Status           ProofStatus
	ProverRequestID  string
	ProofRequestTime int64
	L1BlockNumber    uint64
	L1BlockHash      string
}

type ProofDB struct {
	db *sql.DB
}

func InitDB(dbPath string) (ProofDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return ProofDB{}, err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS proof_requests (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            type TEXT NOT NULL CHECK(type IN ('SPAN', 'AGG')),
            start_block_number INTEGER NOT NULL,
            end_block_number INTEGER NOT NULL,
            status TEXT NOT NULL CHECK(status IN ('UNREQ', 'REQ', 'FAILED', 'COMPLETE')),
            prover_request_id TEXT,
            proof_request_time INTEGER,
            l1_block_number INTEGER,
            l1_block_hash TEXT
        )
    `)
	if err != nil {
		return ProofDB{}, err
	}

	return ProofDB{db}, nil
}

// TODO: Flip this so it's an arg on the DB?
func (db ProofDB) newEntry(pr ProofRequest) error {
	_, err := db.db.Exec(`
        INSERT INTO proof_requests
        (type, start_block_number, end_block_number, status, prover_request_id, proof_request_time, l1_block_number, l1_block_hash)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, pr.Type, pr.StartBlock, pr.EndBlock, pr.Status, pr.ProverRequestID, pr.ProofRequestTime, pr.L1BlockNumber, pr.L1BlockHash)
	return err
}

func (db ProofDB) updateProofStatus(proverRequestID string, newStatus ProofStatus) error {
	query := `
        UPDATE proof_requests
        SET status = ?
        WHERE prover_request_id = ?
    `

	result, err := db.db.Exec(query, string(newStatus), proverRequestID)
	if err != nil {
		return fmt.Errorf("failed to update proof status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no proof request found with ProverRequestID: %s", proverRequestID)
	}

	return nil
}

func (db ProofDB) addL1BlockInfo(proofRequest ProofRequest, l1BlockNumber uint64, l1BlockHash string) error {
	query := `
		UPDATE proof_requests
		SET l1_block_number = ?, l1_block_hash = ?
		WHERE type = 'AGG' AND status = 'UNREQ' AND start_block_number = ? AND end_block_number = ?
	`

	result, err := db.db.Exec(query, l1BlockNumber, l1BlockHash, proofRequest.StartBlock, proofRequest.EndBlock)
	if err != nil {
		return fmt.Errorf("failed to update proof status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no proof request found with ProverRequestID: %s", proofRequest.ProverRequestID)
	}

	return nil
}

func (db ProofDB) getLatestEndRequested() (uint64, error) {
	var endBlock uint64
	// claude says err := l.db.QueryRow("SELECT COALESCE(MAX(end_block_number), 0) FROM proof_requests").Scan(&nextBlock)
	err := db.db.QueryRow("SELECT MAX(end_block_number) FROM proof_requests").Scan(&endBlock)
	return endBlock, err
}

func (db ProofDB) getRequestedProofs() ([]ProofRequest, error) {
	rows, err := db.db.Query(`
		SELECT type, start_block_number, end_block_number, status, prover_request_id, proof_request_time, l1_block_number, l1_block_hash
		FROM proof_requests
		WHERE status = 'REQ'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToProofRequests(rows)
}

func (db ProofDB) getUnrequestedProofs() ([]ProofRequest, error) {
	rows, err := db.db.Query(`
		SELECT type, start_block_number, end_block_number, status, prover_request_id, proof_request_time, l1_block_number, l1_block_hash
		FROM proof_requests
		WHERE status = 'UNREQ'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToProofRequests(rows)
}

func (db *ProofDB) getCompletedAggProofRequest(startBlock uint64) (*ProofRequest, error) {
	query := `
        SELECT type, start_block_number, end_block_number, status, prover_request_id, proof_request_time, l1_block_number, l1_block_hash
        FROM proof_requests
        WHERE type = 'AGG' AND start_block_number = ? AND status = 'COMPLETE'
        LIMIT 1
    `

	var proof ProofRequest
	err := db.db.QueryRow(query, startBlock).Scan(
		&proof.Type, &proof.StartBlock, &proof.EndBlock, &proof.Status, &proof.ProverRequestID,
		&proof.ProofRequestTime, &proof.L1BlockNumber, &proof.L1BlockHash,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query completed AGG proof: %w", err)
	}

	return &proof, nil
}

func (db *ProofDB) tryCreateAggProofFromSpanProofs(latestOutputIndex, nextOutputIndex uint64) (bool, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if there's already an AGG proof in progress
	var count int
	err = tx.QueryRow(`
        SELECT COUNT(*) FROM proof_requests
        WHERE type = 'AGG' AND start_block_number = ? AND status != 'FAILED'
    `, latestOutputIndex).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check existing AGG proofs: %w", err)
	}
	if count > 0 {
		return false, nil // There's already an AGG proof in progress
	}

	// Find consecutive SPAN proofs
	// TODO: Handle off by ones. I think output index is where we start (we haven't proceessed txs, but we sill add 1 below?)
	start := latestOutputIndex
	var end uint64
	for {
		var spanEnd uint64
		err = tx.QueryRow(`
            SELECT end_block_number FROM proof_requests
            WHERE type = 'SPAN' AND status = 'COMPLETE' AND start_block_number = ?
        `, start).Scan(&spanEnd)
		if err == sql.ErrNoRows {
			break // No more consecutive SPAN proofs
		}
		if err != nil {
			return false, fmt.Errorf("failed to query SPAN proof: %w", err)
		}
		end = spanEnd
		start = spanEnd + 1
	}

	if end < nextOutputIndex {
		return false, nil // Not enough SPAN proofs to create an AGG proof
	}

	// Create a new AGG proof request
	_, err = tx.Exec(`
        INSERT INTO proof_requests (type, start_block_number, end_block_number, status)
        VALUES ('AGG', ?, ?, 'UNREQ')
    `, latestOutputIndex, end)
	if err != nil {
		return false, fmt.Errorf("failed to insert AGG proof request: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return true, nil
}

func rowsToProofRequests(rows *sql.Rows) ([]ProofRequest, error) {
	var prs []ProofRequest
	for rows.Next() {
		var pr ProofRequest
		if err := rows.Scan(&pr.Type, &pr.StartBlock, &pr.EndBlock, &pr.Status, &pr.ProverRequestID, &pr.ProofRequestTime, &pr.L1BlockNumber, &pr.L1BlockHash); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, nil
}
