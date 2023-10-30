package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0012 struct{}

func (m migrationTest0012) InsertData(db *sql.DB) error {
	addMonitoredTx := `
        INSERT INTO state.monitored_txs (owner, id, from_addr, to_addr, nonce, value, data, gas, gas_price, status, history, block_num, created_at, updated_at) 
                                VALUES (   $1, $2,        $3,      $4,    $5,    $6,   $7,  $8,        $9,    $10,     $11,       $12,        $13,        $14);`

	args := []interface{}{
		"owner", "id1", common.HexToAddress("0x111").String(), common.HexToAddress("0x222").String(), 333, 444,
		[]byte{5, 5, 5}, 666, 777, "status", []string{common.HexToHash("0x888").String()}, 999, time.Now(), time.Now(),
	}
	if _, err := db.Exec(addMonitoredTx, args...); err != nil {
		return err
	}

	return nil
}

func (m migrationTest0012) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	addMonitoredTx := `
	INSERT INTO state.monitored_txs (owner, id, from_addr, to_addr, nonce, value, data, gas, gas_price, status, history, block_num, created_at, updated_at, gas_offset) 
                            VALUES (   $1, $2,        $3,      $4,    $5,    $6,   $7,  $8,        $9,    $10,     $11,       $12,        $13,        $14,        $15);`

	args := []interface{}{
		"owner", "id2", common.HexToAddress("0x111").String(), common.HexToAddress("0x222").String(), 333, 444,
		[]byte{5, 5, 5}, 666, 777, "status", []string{common.HexToHash("0x888").String()}, 999, time.Now(), time.Now(),
		101010,
	}
	_, err := db.Exec(addMonitoredTx, args...)
	assert.NoError(t, err)

	gasOffset := 999

	getGasOffsetQuery := `SELECT gas_offset FROM state.monitored_txs WHERE id = $1`
	err = db.QueryRow(getGasOffsetQuery, "id1").Scan(&gasOffset)
	assert.NoError(t, err)
	assert.Equal(t, 0, gasOffset)

	err = db.QueryRow(getGasOffsetQuery, "id2").Scan(&gasOffset)
	assert.NoError(t, err)
	assert.Equal(t, 101010, gasOffset)
}

func (m migrationTest0012) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {

}

func TestMigration0012(t *testing.T) {
	runMigrationTest(t, 12, migrationTest0012{})
}
