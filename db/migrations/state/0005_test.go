package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0005 struct{}

func (m migrationTest0005) InsertData(db *sql.DB) error {
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
	// Insert verified batch
	const insertVerifiedBatch = `INSERT INTO state.verified_batch (
		batch_num, tx_hash, aggregator, state_root, block_num
	) VALUES (
		1, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x514910771af9ca656af840dff83e8264ecf986ca', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', 1
	);`
	_, err = db.Exec(insertVerifiedBatch)
	return err
}

func (m migrationTest0005) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	// Insert batch
	_, err := db.Exec(`INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, state_root, acc_input_hash, timestamp, coinbase, raw_txs_data)
	VALUES (2, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			'0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			$1, '0x2536C2745Ac4A584656A830f7bdCd329c94e8F30', $2)`, time.Now(), common.HexToHash("0x29e885edaf8e0000000000000000a23cf2d7d9f1"))
	assert.NoError(t, err)
	// Insert virtual batch
	const insertVirtualBatch = `INSERT INTO state.virtual_batch (
		batch_num, tx_hash, coinbase, block_num
		) VALUES (
			2, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x514910771af9ca656af840dff83e8264ecf986ca', 1);`
	_, err = db.Exec(insertVirtualBatch)
	assert.NoError(t, err)

	// Insert verified batch
	const insertVerifiedBatch = `INSERT INTO state.verified_batch (
		batch_num, tx_hash, aggregator, state_root, block_num, is_trusted
	) VALUES (
		2, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x514910771af9ca656af840dff83e8264ecf986ca', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', 1, true
	);`
	_, err = db.Exec(insertVerifiedBatch)
	assert.NoError(t, err)

	// Insert monitored_txs
	const insertMonitoredTxs = `INSERT INTO state.monitored_txs (
		owner, id, from_addr, to_addr, nonce, value, data, gas, gas_price, status, block_num, created_at, updated_at
	) VALUES (
		'0x514910771af9ca656af840dff83e8264ecf986ca', '1', '0x514910771af9ca656af840dff83e8264ecf986ca', '0x514910771af9ca656af840dff83e8264ecf986ca', 1, 0, '0x', 100, 12, 'created', 1, $1, $2
	);`
	_, err = db.Exec(insertMonitoredTxs, time.Now(), time.Now())
	assert.NoError(t, err)
}

func (m migrationTest0005) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	// Insert batch
	_, err := db.Exec(`INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, state_root, acc_input_hash, timestamp, coinbase, raw_txs_data)
	VALUES (3, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			'0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1',
			$1, '0x2536C2745Ac4A584656A830f7bdCd329c94e8F30', $2)`, time.Now(), common.HexToHash("0x29e885edaf8e0000000000000000a23cf2d7d9f1"))
	assert.NoError(t, err)
	// Insert virtual batch
	const insertVirtualBatch = `INSERT INTO state.virtual_batch (
		batch_num, tx_hash, coinbase, block_num
		) VALUES (
			3, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x514910771af9ca656af840dff83e8264ecf986ca', 1);`
	_, err = db.Exec(insertVirtualBatch)
	assert.NoError(t, err)

	// Insert verified batch
	insertVerifiedBatch := `INSERT INTO state.verified_batch (
		batch_num, tx_hash, aggregator, state_root, block_num, is_trusted
	) VALUES (
		3, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x514910771af9ca656af840dff83e8264ecf986ca', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', 1, true
	);`
	_, err = db.Exec(insertVerifiedBatch)
	assert.Error(t, err)
	insertVerifiedBatch = `INSERT INTO state.verified_batch (
		batch_num, tx_hash, aggregator, state_root, block_num
	) VALUES (
		3, '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', '0x514910771af9ca656af840dff83e8264ecf986ca', '0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1', 1
	);`
	_, err = db.Exec(insertVerifiedBatch)
	assert.NoError(t, err)

	// Insert monitored_txs
	const insertMonitoredTxs = `INSERT INTO state.monitored_txs (
		owner, id, from_addr, to_addr, nonce, value, data, gas, gas_price, status, block_num, created_at, updated_at
	) VALUES (
		'0x514910771af9ca656af840dff83e8264ecf986ca', '1', '0x514910771af9ca656af840dff83e8264ecf986ca', '0x514910771af9ca656af840dff83e8264ecf986ca', 1, 0, '0x', 100, 12, 'created', 1, $1, $2
	);`
	_, err = db.Exec(insertMonitoredTxs, time.Now(), time.Now())
	assert.Error(t, err)
}

func TestMigration0005(t *testing.T) {
	runMigrationTest(t, 5, migrationTest0005{})
}
