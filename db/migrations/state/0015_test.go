package migrations_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0015 struct{}

func (m migrationTest0015) InsertData(db *sql.DB) error {
	return nil
}

func (m migrationTest0015) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	indexes := []string{
		"idx_receipt_tx_index",
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

func (m migrationTest0015) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	indexes := []string{
		"idx_receipt_tx_index",
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

func TestMigration0015(t *testing.T) {
	runMigrationTest(t, 15, migrationTest0015{})
}
