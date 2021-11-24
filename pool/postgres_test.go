package pool

import (
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/test/dbutils"
	"github.com/stretchr/testify/assert"
)

var cfg = db.Config{
	Database: "testing",
	User:     "hermez",
	Password: "password",
	Host:     "localhost",
	Port:     "5432",
}

func TestAddTx(t *testing.T) {
	// Start DB Server
	dbutils.StartPostgreSQL(cfg.Database, cfg.User, cfg.Password, "") //nolint:gosec,errcheck
	defer dbutils.StopPostgreSQL()                                    //nolint:gosec,errcheck

	err := db.RunMigrations(cfg)
	if err != nil {
		t.Error(err)
	}

	p, err := newPostgresPool(cfg)
	if err != nil {
		t.Error(err)
	}

	txRLPHash := "0xf86e8212658082520894fd8b27a263e19f0e9592180e61f0f8c9dfeb1ff6880de0b6b3a764000080850133333355a01eac4c2defc7ed767ae36bbd02613c581b8fb87d0e4f579c9ee3a7cfdb16faa7a043ce30f43d952b9d034cf8f04fecb631192a5dbc7ee2a47f1f49c0d022a8849d"
	b, err := hex.DecodeHex(txRLPHash)
	if err != nil {
		t.Error(err)
	}
	tx := new(types.Transaction)
	tx.UnmarshalBinary(b) //nolint:gosec,errcheck

	err = p.AddTx(*tx)
	if err != nil {
		t.Error(err)
	}

	sqlDB, err := db.NewSQLDB(cfg)
	if err != nil {
		t.Error(err)
	}
	defer sqlDB.Close() //nolint:gosec,errcheck

	rows, err := sqlDB.Query("SELECT hash, encoded, decoded, state FROM pool.txs")
	if err != nil {
		t.Error(err)
	}

	c := 0
	for rows.Next() {
		var hash, encoded, decoded, state string
		err := rows.Scan(&hash, &encoded, &decoded, &state)
		if err != nil {
			t.Error(err)
		}
		b, _ := tx.MarshalJSON()

		assert.Equal(t, "0xa3cff5abdf47d4feb8204a45c0a8c58fc9b9bb9b29c6588c1d206b746815e9cc", hash, "invalid hash")
		assert.Equal(t, txRLPHash, encoded, "invalid encoded")
		assert.JSONEq(t, string(b), decoded, "invalid decoded")
		assert.Equal(t, string(TxStatePending), state, "invalid tx state")
		c++
	}

	assert.Equal(t, 1, c, "invalid number of txs in the pool")
}
