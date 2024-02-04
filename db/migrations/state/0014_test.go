package migrations_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0014 struct{}

func (m migrationTest0014) InsertData(db *sql.DB) error {
	return nil
}

func (m migrationTest0014) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	indexes := []string{
		"idx_batch_global_exit_root",
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

func (m migrationTest0014) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	indexes := []string{
		"idx_batch_global_exit_root",
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

func TestMigration0014(t *testing.T) {
	runMigrationTest(t, 14, migrationTest0014{})
}
