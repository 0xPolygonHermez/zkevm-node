package migrations_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0006 struct{}

func (m migrationTest0006) InsertData(db *sql.DB) error {
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
}

func TestMigration0006(t *testing.T) {
	runMigrationTest(t, 6, migrationTest0006{})
}
