package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0002 struct{}

func (m migrationTest0002) InsertData(db *sql.DB) error {
	// Insert block to respect the FKey
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := db.Exec(addBlock, 1, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"); err != nil {
		return err
	}
	// Insert batches
	for i := 0; i < 4; i++ {
		_, err := db.Exec(`INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, state_root, acc_input_hash, timestamp, coinbase, raw_txs_data)
	VALUES ($1, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			'0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			$2, '0x2536C2745Ac4A584656A830f7bdCd329c94e8F30', $3)`, i, time.Now(), common.HexToHash("0x29e885edaf8e0000000000000000a23cf2d7d9f1"))
		if err != nil {
			return err
		}
	}
	// Insert proof
	const insertProof = `INSERT INTO state.proof (
		batch_num, batch_num_final, proof, proof_id, input_prover, prover, generating
	) VALUES (
		1, 1, '{"test": "test"}','proof_identifier','{"test": "test"}','prover 1', true
	);`
	_, err := db.Exec(insertProof)
	if err != nil {
		return err
	}
	// Insert virtual batch
	const insertVirtualBatch = `INSERT INTO state.virtual_batch (
		batch_num, tx_hash, coinbase, block_num
		) VALUES (
			1, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x514910771af9ca656af840dff83e8264ecf986ca', 1);`
	_, err = db.Exec(insertVirtualBatch)
	if err != nil {
		return err
	}
	return nil
}

var indexes = []string{"transaction_l2_block_num_idx", "l2block_batch_num_idx", "l2block_received_at_idx",
	"batch_timestamp_idx", "log_tx_hash_idx", "log_address_idx"}

func (m migrationTest0002) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	for _, idx := range indexes {
		// getIndex
		const getIndex = `SELECT count(*) FROM pg_indexes WHERE indexname = $1;`
		row := db.QueryRow(getIndex, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 1, result)
	}
	// Insert new proof
	const insertNewProof = `INSERT INTO state.proof (
		batch_num, batch_num_final, proof, proof_id, input_prover, prover, generating_since, prover_id, updated_at, created_at
	) VALUES (
		2, 2, '{"test": "test"}','proof_identifier','{"test": "test"}','prover 1', $1, 'prover identifier', $1, $1
	);`
	_, err := db.Exec(insertNewProof, time.Now())
	assert.NoError(t, err)
	const insertOldProof = `INSERT INTO state.proof (
		batch_num, batch_num_final, proof, proof_id, input_prover, prover, generating
	) VALUES (
		3, 3, '{"test": "test"}','proof_identifier','{"test": "test"}','prover 1', true
	);`
	_, err = db.Exec(insertOldProof)
	assert.Error(t, err)
	// Insert virtual batch
	const insertVirtualBatch = `INSERT INTO state.virtual_batch (
		batch_num, tx_hash, coinbase, block_num, sequencer_addr)
		VALUES (2, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x514910771af9ca656af840dff83e8264ecf986ca', 1, '0x514910771af9ca656af840dff83e8264ecf986ca');`
	_, err = db.Exec(insertVirtualBatch)
	assert.NoError(t, err)
	// Insert reorg
	const insertReorg = `INSERT INTO state.trusted_reorg (batch_num, reason)
		VALUES (2, 'reason of the trusted reorg');`
	_, err = db.Exec(insertReorg)
	assert.NoError(t, err)
}

func (m migrationTest0002) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	for _, idx := range indexes {
		// getIndex
		const getIndex = `SELECT count(*) FROM pg_indexes WHERE indexname = $1;`
		row := db.QueryRow(getIndex, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 0, result)
	}
	// Insert new proof
	const insertNewProof = `INSERT INTO state.proof (
		batch_num, batch_num_final, proof, proof_id, input_prover, prover, generating_since, prover_id, updated_at, created_at
	) VALUES (
		3, 3, '{"test": "test"}','proof_identifier','{"test": "test"}','prover 1', $1, 'prover identifier', $1, $1
	);`
	_, err := db.Exec(insertNewProof, time.Now())
	assert.Error(t, err)
	const insertOldProof = `INSERT INTO state.proof (
		batch_num, batch_num_final, proof, proof_id, input_prover, prover, generating
	) VALUES (
		3, 3, '{"test": "test"}','proof_identifier','{"test": "test"}','prover 1', true
	);`
	_, err = db.Exec(insertOldProof)
	assert.NoError(t, err)
	// Insert virtual batch
	insertVirtualBatch := `INSERT INTO state.virtual_batch (
		batch_num, tx_hash, coinbase, block_num, sequencer_addr,
		) VALUES (
			3, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x514910771af9ca656af840dff83e8264ecf986ca', 1, '0x514910771af9ca656af840dff83e8264ecf986ca');`
	_, err = db.Exec(insertVirtualBatch)
	assert.Error(t, err)
	// Insert virtual batch
	insertVirtualBatch = `INSERT INTO state.virtual_batch (
		batch_num, tx_hash, coinbase, block_num)
		VALUES (3, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x514910771af9ca656af840dff83e8264ecf986ca', 1);`
	_, err = db.Exec(insertVirtualBatch)
	assert.NoError(t, err)
	// Insert reorg
	const insertReorg = `INSERT INTO state.trusted_reorg (batch_num, reason)
		VALUES (2, 'reason of the trusted reorg');`
	_, err = db.Exec(insertReorg)
	assert.Error(t, err)
}

func TestMigration0002(t *testing.T) {
	runMigrationTest(t, 2, migrationTest0002{})
}
