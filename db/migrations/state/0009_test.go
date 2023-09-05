package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0009 struct{}

func (m migrationTest0009) InsertData(db *sql.DB) error {
	// Insert block to respect the FKey
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := db.Exec(addBlock, 1, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"); err != nil {
		return err
	}
	if _, err := db.Exec(addBlock, 2, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f2"); err != nil {
		return err
	}

	return nil
}

func (m migrationTest0009) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	// Insert forkID
	const insertForkID = `INSERT INTO state.fork_id (
		from_batch_num, to_batch_num, fork_id, version, block_num) VALUES (
		1, 10, 2, 'First version', 1
	);`
	_, err := db.Exec(insertForkID)
	assert.Error(t, err)

	const insertForkID2 = `INSERT INTO state.fork_id (
		from_batch_num, to_batch_num, fork_id, version) VALUES (
		1, 10, 2, 'First version'
	);`
	_, err = db.Exec(insertForkID2)
	assert.NoError(t, err)
}

func (m migrationTest0009) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	// Insert forkID
	const insertForkID = `INSERT INTO state.fork_id (
		from_batch_num, to_batch_num, fork_id, version, block_num) VALUES (
		1, 10, 1, 'First version', 2
	);`
	_, err := db.Exec(insertForkID)
	assert.NoError(t, err)
	const insertForkID2 = `INSERT INTO state.fork_id (
		from_batch_num, to_batch_num, fork_id, version, block_num) VALUES (
		1, 10, 1, 'First version'
	);`
	_, err = db.Exec(insertForkID2)
	assert.Error(t, err)
}

func TestMigration0009(t *testing.T) {
	runMigrationTest(t, 9, migrationTest0009{})
}
