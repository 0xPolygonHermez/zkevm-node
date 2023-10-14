package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0008 struct{}

var indexes_0008 = []string{"receipt_block_num_idx"}

func (m migrationTest0008) InsertData(db *sql.DB) error {
	// Insert block to respect the FKey
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := db.Exec(addBlock, 1, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"); err != nil {
		return err
	}
	if _, err := db.Exec(addBlock, 2, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f2"); err != nil {
		return err
	}
	const insertBatch = `
		INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num) 
		VALUES (0,'0x000', '0x000', '0x000', '0x000', now(), '0x000', null, null)`

	// insert batch
	_, err := db.Exec(insertBatch)
	if err != nil {
		return err
	}

	const insertL2Block = `
		INSERT INTO state.l2block (block_num, block_hash, header, uncles, parent_hash, state_root, received_at, batch_num, created_at)
		VALUES (0, '0x001', '{}', '{}', '0x002', '0x003', now(), 0, now())`

	// insert l2 block
	_, err = db.Exec(insertL2Block)
	if err != nil {
		return err
	}

	const insertTx = `
		INSERT INTO state.transaction (hash, encoded, decoded, l2_block_num, effective_percentage)
		VALUES ('0x001', 'ABCDEF', '{}', 0, 255)`

	// insert tx
	_, err = db.Exec(insertTx)
	if err != nil {
		return err
	}
	const insertReceipt = `
		INSERT INTO state.receipt (tx_hash, type, post_state, status, cumulative_gas_used, gas_used, effective_gas_price, block_num, tx_index, contract_address)
		VALUES ('0x001', 1, null, 1, 1234, 1234, 1, 0, 0, '0x002')`

	// insert receipt
	_, err = db.Exec(insertReceipt)
	if err != nil {
		return err
	}

	return nil
}

func (m migrationTest0008) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	// Insert forkID
	const insertForkID = `INSERT INTO state.fork_id (
		from_batch_num, to_batch_num, fork_id, version, block_num) VALUES (
		1, 10, 2, 'First version', 1
	);`
	_, err := db.Exec(insertForkID)
	assert.NoError(t, err)

	_, err = db.Exec(insertForkID)
	assert.Error(t, err)
	for _, idx := range indexes_0008 {
		// getIndex
		const getIndex = `SELECT count(*) FROM pg_indexes WHERE indexname = $1;`
		row := db.QueryRow(getIndex, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 1, result)
	}
}

func (m migrationTest0008) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	// Insert forkID
	const insertForkID = `INSERT INTO state.fork_id (
		from_batch_num, to_batch_num, fork_id, version, block_num) VALUES (
		1, 10, 2, 'First version', 2
	);`
	_, err := db.Exec(insertForkID)
	assert.Error(t, err)
	for _, idx := range indexes_0008 {
		// getIndex
		const getIndex = `SELECT count(*) FROM pg_indexes WHERE indexname = $1;`
		row := db.QueryRow(getIndex, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 0, result)
	}
}

func TestMigration0008(t *testing.T) {
	runMigrationTest(t, 8, migrationTest0008{})
}
