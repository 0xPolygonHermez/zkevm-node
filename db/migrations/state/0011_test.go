package migrations_test

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

func (m migrationTest0011) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	indexes := []string{
		"l2block_created_at_idx",
		"log_log_index_idx",
		"log_topic0_idx",
		"log_topic1_idx",
		"log_topic2_idx",
		"log_topic3_idx",
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

	// Check column egp_log exists in state.transactions table
	const getFinalDeviationColumn = `SELECT count(*) FROM information_schema.columns WHERE table_name='transaction' and column_name='egp_log'`
	row := db.QueryRow(getFinalDeviationColumn)
	var result int
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 1, result)
}

func (m migrationTest0011) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	indexes := []string{
		"l2block_created_at_idx",
		"log_log_index_idx",
		"log_topic0_idx",
		"log_topic1_idx",
		"log_topic2_idx",
		"log_topic3_idx",
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

	// Check column egp_log doesn't exists in state.transactions table
	const getFinalDeviationColumn = `SELECT count(*) FROM information_schema.columns WHERE table_name='transaction' and column_name='egp_log'`
	row := db.QueryRow(getFinalDeviationColumn)
	var result int
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 0, result)
}

func TestMigration0011(t *testing.T) {
	runMigrationTest(t, 11, migrationTest0011{})
}
