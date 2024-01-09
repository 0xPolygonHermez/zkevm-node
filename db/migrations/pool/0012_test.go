package pool_migrations_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0012 struct{}

func (m migrationTest0012) InsertData(db *sql.DB) error {
	const insertTx = `
		INSERT INTO pool.transaction (hash, ip, received_at, from_address)
		VALUES ('0x0001', '127.0.0.1', '2023-12-07', '0x0011')`

	_, err := db.Exec(insertTx)
	if err != nil {
		return err
	}

	return nil
}

var indexesMigration12 = []string{
	"idx_transaction_l2_hash",
}

func (m migrationTest0012) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	// Check indexes adding
	for _, idx := range indexesMigration12 {
		// getIndex
		const getIndex = `SELECT count(*) FROM pg_indexes WHERE indexname = $1;`
		row := db.QueryRow(getIndex, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 1, result)
	}

	const insertTx = `
		INSERT INTO pool.transaction (hash, ip, received_at, from_address, used_sha256_hashes)
		VALUES ('0x0002', '127.0.0.1', '2023-12-07', '0x0022', 222)`

	_, err := db.Exec(insertTx)
	assert.NoError(t, err)
}

func (m migrationTest0012) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	// Check indexes removing
	for _, idx := range indexesMigration12 {
		// getIndex
		const getIndex = `SELECT count(*) FROM pg_indexes WHERE indexname = $1;`
		row := db.QueryRow(getIndex, idx)
		var result int
		assert.NoError(t, row.Scan(&result))
		assert.Equal(t, 0, result)
	}
}

func TestMigration0012(t *testing.T) {
	runMigrationTest(t, 12, migrationTest0012{})
}
