package pool_migrations_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0011 struct{}

func (m migrationTest0011) InsertData(db *sql.DB) error {
	return nil
}

var indexes = []string{
	"idx_transaction_from_nonce",
	"idx_transaction_status",
	"idx_transaction_hash",
}

func (m migrationTest0011) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
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

func (m migrationTest0011) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
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

func TestMigration0011(t *testing.T) {
	runMigrationTest(t, 11, migrationTest0011{})
}
