package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	blockHashValue         = "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"
	mainExitRootValue      = "0x83fc198de31e1b2b1a8212d2430fbb7766c13d9ad305637dea3759065606475d"
	rollupExitRootValue    = "0xadb91a6a1fce56eaea561002bc9a993f4e65a7710bd72f4eee3067cbd73a743c"
	globalExitRootValue    = "0x5bf4af1a651a2a74b36e6eb208481f94c69fc959f756223dfa49608061937585"
	previousBlockHashValue = "0xe865e912b504572a4d80ad018e29797e3c11f00bf9ae2549548a25779c9d7e57"
	l1InfoRootValue        = "0x2b9484b83c6398033241865b015fb9430eb3e159182a6075d00c924845cc393e"
)

// this migration changes length of the token name
type migrationTest0013 struct{}

func (m migrationTest0013) insertBlock(blockNumber uint64, db *sql.DB) error {
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := db.Exec(addBlock, blockNumber, time.Now(), blockHashValue); err != nil {
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
	if err = m.insertRowInOldTable(db, []interface{}{123, time.Now(), mainExitRootValue, rollupExitRootValue, globalExitRootValue}); err != nil {
		return err
	}
	if err = m.insertRowInOldTable(db, []interface{}{124, time.Now(), mainExitRootValue, rollupExitRootValue, globalExitRootValue}); err != nil {
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
	prevBlockHash := previousBlockHashValue
	l1InfoRoot := l1InfoRootValue
	err = m.insertRowInMigratedTable(db, []interface{}{125, time.Now(), mainExitRootValue, rollupExitRootValue, globalExitRootValue, prevBlockHash, l1InfoRoot, 1})
	assert.NoError(t, err)
	// insert duplicated l1_info_root
	err = m.insertRowInMigratedTable(db, []interface{}{126, time.Now(), mainExitRootValue, rollupExitRootValue, globalExitRootValue, prevBlockHash, l1InfoRoot, 1})
	assert.Error(t, err)

	// insert in the old way must work
	err = m.insertRowInOldTable(db, []interface{}{127, time.Now(), mainExitRootValue, rollupExitRootValue, globalExitRootValue})
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
