package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// This migration creates the fiat table

type migrationTest0004 struct{}

func (m migrationTest0004) InsertData(db *sql.DB) error {
	// Insert block to respect the FKey
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := db.Exec(addBlock, 1, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"); err != nil {
		return err
	}
	const addForcedBatch = "INSERT INTO state.forced_batch (forced_batch_num, global_exit_root, raw_txs_data, coinbase, timestamp, block_num, batch_num) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if _, err := db.Exec(addForcedBatch, 1, "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1", "", "0x2536C2745Ac4A584656A830f7bdCd329c94e8F30", time.Now(), 1, 1); err != nil {
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
	return nil
}

func (m migrationTest0004) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	// Insert batch
	_, err := db.Exec(`INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, state_root, acc_input_hash, timestamp, coinbase, raw_txs_data, forced_batch_num)
	VALUES (2, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			'0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			$1, '0x2536C2745Ac4A584656A830f7bdCd329c94e8F30', $2, 1)`, time.Now(), common.HexToHash("0x29e885edaf8e0000000000000000a23cf2d7d9f1"))
	assert.NoError(t, err)
	addForcedBatch := "INSERT INTO state.forced_batch (forced_batch_num, global_exit_root, raw_txs_data, coinbase, timestamp, block_num, batch_num) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err = db.Exec(addForcedBatch, 1, "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1", "", "0x2536C2745Ac4A584656A830f7bdCd329c94e8F30", time.Now(), 1, 2)
	assert.Error(t, err)
	addForcedBatch = "INSERT INTO state.forced_batch (forced_batch_num, global_exit_root, raw_txs_data, coinbase, timestamp, block_num) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = db.Exec(addForcedBatch, 2, "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1", "", "0x2536C2745Ac4A584656A830f7bdCd329c94e8F30", time.Now(), 1)
	assert.NoError(t, err)
}

func (m migrationTest0004) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	const addForcedBatch = "INSERT INTO state.forced_batch (forced_batch_num, global_exit_root, raw_txs_data, coinbase, timestamp, block_num, batch_num) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := db.Exec(addForcedBatch, 3, "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1", "", "0x2536C2745Ac4A584656A830f7bdCd329c94e8F30", time.Now(), 1, 1)
	assert.NoError(t, err)
	// Insert batch
	_, err = db.Exec(`INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, state_root, acc_input_hash, timestamp, coinbase, raw_txs_data, forced_batch_num)
	VALUES (3, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			'0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			$1, '0x2536C2745Ac4A584656A830f7bdCd329c94e8F30', $2, 1)`, time.Now(), common.HexToHash("0x29e885edaf8e0000000000000000a23cf2d7d9f1"))
	assert.Error(t, err)
	_, err = db.Exec(`INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, state_root, acc_input_hash, timestamp, coinbase, raw_txs_data)
	VALUES (3, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			'0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			$1, '0x2536C2745Ac4A584656A830f7bdCd329c94e8F30', $2)`, time.Now(), common.HexToHash("0x29e885edaf8e0000000000000000a23cf2d7d9f1"))
	assert.NoError(t, err)
}

func TestMigration0004(t *testing.T) {
	runMigrationTest(t, 4, migrationTest0004{})
}
