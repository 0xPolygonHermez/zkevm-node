package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0013 struct{}

func (m migrationTest0013) insertBlock(blockNumber uint64, db *sql.DB) error {
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := db.Exec(addBlock, blockNumber, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"); err != nil {
		return err
	}
	return nil
}

func (m migrationTest0013) insertRowInOldTable(db *sql.DB, args []interface{}) error {
	insert := `
        INSERT INTO state.exit_root (block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root) 
                             VALUES ($1,        $2,        $3,                $4,               $5  );`

	_, err := db.Exec(insert, args...)
	return err
}

func (m migrationTest0013) InsertData(db *sql.DB) error {
	var err error
	if err = m.insertBlock(uint64(123), db); err != nil {
		return err
	}
	if err = m.insertBlock(uint64(124), db); err != nil {
		return err
	}
	if err = m.insertRowInOldTable(db, []interface{}{123, time.Now(), "mer", "rer", "ger"}); err != nil {
		return err
	}
	if err = m.insertRowInOldTable(db, []interface{}{124, time.Now(), "mer", "rer", "ger"}); err != nil {
		return err
	}

	return nil
}
func (m migrationTest0013) insertRowInMigratedTable(db *sql.DB, args []interface{}) error {
	insert := `
        INSERT INTO state.exit_root (block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index) 
                             VALUES ($1,        $2,        $3,                $4,               $5,               $6,              $7,           $8);`

	_, err := db.Exec(insert, args...)
	return err
}

func (m migrationTest0013) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	err := m.insertBlock(uint64(125), db)
	assert.NoError(t, err)
	err = m.insertBlock(uint64(126), db)
	assert.NoError(t, err)
	err = m.insertBlock(uint64(127), db)
	assert.NoError(t, err)
	prevBlockHash := "prev_block_hash"
	l1InfoRoot := "l1inforoot"
	err = m.insertRowInMigratedTable(db, []interface{}{125, time.Now(), "mer", "rer", "ger", prevBlockHash, l1InfoRoot, 1})
	assert.NoError(t, err)
	// insert duplicated l1_info_root
	err = m.insertRowInMigratedTable(db, []interface{}{126, time.Now(), "mer", "rer", "ger", prevBlockHash, l1InfoRoot, 1})
	assert.Error(t, err)

	// insert in the old way must work
	err = m.insertRowInOldTable(db, []interface{}{127, time.Now(), "mer", "rer", "ger"})
	assert.NoError(t, err)

	sqlSelect := `SELECT prev_block_hash, l1_info_root FROM state.exit_root WHERE l1_info_tree_index = $1`
	currentPrevBlockHash := ""
	currentL1InfoRoot := ""
	err = db.QueryRow(sqlSelect, 1).Scan(&currentPrevBlockHash, &currentL1InfoRoot)
	assert.NoError(t, err)
	assert.Equal(t, prevBlockHash, currentPrevBlockHash)
	assert.Equal(t, l1InfoRoot, currentL1InfoRoot)
}

func (m migrationTest0013) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	sqlSelect := `SELECT count(id) FROM state.exit_root`
	count := 0
	err := db.QueryRow(sqlSelect).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 4, count)
}

func TestMigration0013(t *testing.T) {
	runMigrationTest(t, 13, migrationTest0013{})
}
