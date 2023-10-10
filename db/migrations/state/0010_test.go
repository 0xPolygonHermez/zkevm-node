package migrations_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0010 struct{}

func (m migrationTest0010) InsertData(db *sql.DB) error {
	return nil
}

func (m migrationTest0010) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	indexes := []string{"l2block_block_hash_idx"}
	// Check indexes adding
	for _, idx := range indexes {
		// getIndex
		const getIndex = `SELECT count(*) FROM pg_indexes WHERE indexname = $1;`
		row := db.QueryRow(getIndex, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 1, result)
	}

	constraints := []string{"sequences_pkey", "trusted_reorg_pkey", "sync_info_pkey"}
	// Check constraint adding
	for _, idx := range constraints {
		// getConstraint
		const getConstraint = `	SELECT count(*) FROM pg_constraint c WHERE c.conname = $1;`
		row := db.QueryRow(getConstraint, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 1, result)
	}
}

func (m migrationTest0010) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	indexes := []string{"l2block_block_hash_idx"}
	// Check indexes removing
	for _, idx := range indexes {
		// getIndex
		const getIndex = `SELECT count(*) FROM pg_indexes WHERE indexname = $1;`
		row := db.QueryRow(getIndex, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 0, result)
	}

	constraints := []string{"sequences_pkey", "trusted_reorg_pkey", "sync_info_pkey"}
	// Check constraint adding
	for _, idx := range constraints {
		// getConstraint
		const getConstraint = `	SELECT count(*) FROM pg_constraint c WHERE c.conname = $1;`
		row := db.QueryRow(getConstraint, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 0, result)
	}
}

func TestMigration0010(t *testing.T) {
	runMigrationTest(t, 10, migrationTest0010{})
}
