package migrations_test

import (
	"database/sql"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0017 struct{}

func (m migrationTest0017) InsertData(db *sql.DB) error {
	const insertBatch0 = `
		INSERT INTO state.receipt (tx_hash, type, post_state, status, cumulative_gas_used, gas_used, block_num, tx_index, contract_address, effective_gas_price, im_state_root) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	// insert batch
	_, err := db.Exec(insertBatch0, "0x0000", 0, common.Hex2Bytes("0x123456"), 0, 0, 0, 0, 0, "0x0000", 0, common.Hex2Bytes("0x123456"))
	if err != nil {
		return err
	}

	return nil
}

func (m migrationTest0017) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	indexes := []string{
		"l2block_block_num_idx",
	}
	// Check indexes adding
	for _, idx := range indexes {
		// getIndex
		const getIndex = `SELECT count(*) FROM pg_indexes WHERE indexname = $1;`
		row := db.QueryRow(getIndex, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 1, result)
	}
}

func (m migrationTest0017) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	indexes := []string{
		"l2block_block_num_idx",
	}
	// Check indexes removing
	for _, idx := range indexes {
		// getIndex
		const getIndex = `SELECT count(*) FROM pg_indexes WHERE indexname = $1;`
		row := db.QueryRow(getIndex, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 0, result)
	}
}

func TestMigration0017(t *testing.T) {
	runMigrationTest(t, 17, migrationTest0017{})
}
