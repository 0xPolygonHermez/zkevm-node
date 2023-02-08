package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0006 struct{}

func (m migrationTest0006) InsertData(db *sql.DB) error {
	// Insert block to respect the FKey
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := db.Exec(addBlock, 1, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"); err != nil {
		return err
	}
	// Insert batch
	_, err := db.Exec(`INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, state_root, acc_input_hash, timestamp, coinbase, raw_txs_data)
		VALUES (1, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			'0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			$1, '0x2536C2745Ac4A584656A830f7bdCd329c94e8F30', $2)`, time.Now(), common.HexToHash("0x29e885edaf8e0000000000000000a23cf2d7d9f1"))
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

func (m migrationTest0006) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	for _, idx := range indexes {
		// getIndex
		const getIndex = `SELECT count(*) FROM pg_indexes WHERE indexname = $1;`
		row := db.QueryRow(getIndex, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 1, result)
	}
	// Insert batch
	_, err := db.Exec(`INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, state_root, acc_input_hash, timestamp, coinbase, raw_txs_data)
		VALUES (2, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			'0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			$1, '0x2536C2745Ac4A584656A830f7bdCd329c94e8F30', $2)`, time.Now(), common.HexToHash("0x29e885edaf8e0000000000000000a23cf2d7d9f1"))
	assert.NoError(t, err)
	// Insert virtual batch
	const insertVirtualBatch = `INSERT INTO state.virtual_batch (
		batch_num, tx_hash, coinbase, block_num, sequencer_addr)
		VALUES (2, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x514910771af9ca656af840dff83e8264ecf986ca', 1, '0x514910771af9ca656af840dff83e8264ecf986ca');`
	_, err = db.Exec(insertVirtualBatch)
	assert.NoError(t, err)
}

func (m migrationTest0006) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	for _, idx := range indexes {
		// getIndex
		const getIndex = `SELECT count(*) FROM pg_indexes WHERE indexname = $1;`
		row := db.QueryRow(getIndex, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 0, result)
	}
	// Insert batch
	_, err := db.Exec(`INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, state_root, acc_input_hash, timestamp, coinbase, raw_txs_data)
		VALUES (3, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			'0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			$1, '0x2536C2745Ac4A584656A830f7bdCd329c94e8F30', $2)`, time.Now(), common.HexToHash("0x29e885edaf8e0000000000000000a23cf2d7d9f1"))
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
}

func TestMigration0006(t *testing.T) {
	runMigrationTest(t, 6, migrationTest0006{})
}
