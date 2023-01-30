package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// This migration adds the column `eth_tx_hash` on `batch`

type migrationTest0003 struct{}

func (m migrationTest0003) InsertData(db *sql.DB) error {
	return nil
}

func (m migrationTest0003) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	const insertDebug = `INSERT INTO state.debug (error_type, timestamp, payload) VALUES ('error type', $1, 'payload stored')`
	_, err := db.Exec(insertDebug, time.Now())
	assert.NoError(t, err)
}

func (m migrationTest0003) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	const insertDebug = `INSERT INTO state.debug (error_type, timestamp, payload) VALUES ('error type', $1, 'payload stored')`
	_, err := db.Exec(insertDebug, time.Now())
	assert.Error(t, err)
}

func TestMigration0003(t *testing.T) {
	runMigrationTest(t, 3, migrationTest0003{})
}
