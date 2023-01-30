package migrations_test

import (
	"testing"
	"time"
	"database/sql"

	"github.com/stretchr/testify/assert"
)

// This migration creates a different proof table dropping all the information.

type migrationTest0002 struct{}

func (m migrationTest0002) InsertData(db *sql.DB) error {
	// Insert block to respect the FKey
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := db.Exec(addBlock, 1, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"); err != nil {
		return err
	}
	// Insert batch to respect the FKey
	_, err := db.Exec("INSERT INTO state.batch (batch_num) VALUES (1)")
	if err != nil {
		return err
	}
	// Insert old proof
	const insertProof = `INSERT INTO state.proof (
		batch_num, proof, proof_id, input_prover, prover
	) VALUES (
		1,'{}','proof_identifier','{}','prover 1'
	);`
	_, err = db.Exec(insertProof)
	return err
}


func (m migrationTest0002) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	// Insert new proof
	const insertProof = `INSERT INTO state.proof (
		batch_num, batch_num_final, proof, proof_id, input_prover, prover, generating
	) VALUES (
		1, 1, '{}','proof_identifier','{}','prover 1', true
	);`
	_, err := db.Exec(insertProof)
	assert.NoError(t, err)
}

func (m migrationTest0002) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	// Insert new proof
	const insertNewProof = `INSERT INTO state.proof (
		batch_num, batch_num_final, proof, proof_id, input_prover, prover, generating
	) VALUES (
		1, 1, '{}','proof_identifier','{}','prover 1', true
	);`
	_, err := db.Exec(insertNewProof)
	assert.Error(t, err)

	// Insert old proof
	const insertProof = `INSERT INTO state.proof (
		batch_num, proof, proof_id, input_prover, prover
	) VALUES (
		1,'{}','proof_identifier','{}','prover 1'
	);`
	_, err = db.Exec(insertProof)
	assert.NoError(t, err)
}

func TestMigration0002(t *testing.T) {
	runMigrationTest(t, 2, migrationTest0002{})
}
