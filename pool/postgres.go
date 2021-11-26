package pool

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/hex"
)

type postgresPool struct {
	cfg db.Config
}

func newPostgresPool(cfg db.Config) (*postgresPool, error) {
	return &postgresPool{
		cfg: cfg,
	}, nil
}

func (p *postgresPool) AddTx(tx types.Transaction) error {
	// hash
	hash := tx.Hash().Hex()

	// encoded
	b, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	encoded := hex.EncodeToHex(b)

	// decoded
	b, err = tx.MarshalJSON()
	if err != nil {
		return err
	}
	decoded := string(b)

	// transaction state
	state := TxStatePending

	// get connection
	sqlDB, err := db.NewSQLDB(p.cfg)
	if err != nil {
		return err
	}
	defer sqlDB.Close() //nolint:errcheck

	// save
	sql := "INSERT INTO pool.txs (hash, encoded, decoded, state) VALUES($1, $2, $3, $4)"
	if _, err := sqlDB.Exec(sql, hash, encoded, decoded, state); err != nil {
		return err
	}
	return nil
}

func (p *postgresPool) GetPendingTxs() ([]Transaction, error) {
	// get connection
	sqlDB, err := db.NewSQLDB(p.cfg)
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close() //nolint:errcheck

	sql := "SELECT encoded, state FROM pool.txs WHERE state = $1"
	rows, err := sqlDB.Query(sql, TxStatePending)
	if err != nil {
		return nil, err
	}

	txs := []Transaction{}
	for rows.Next() {
		var encoded, state string

		if err := rows.Scan(&encoded, &state); err != nil {
			return nil, err
		}

		tx := new(Transaction)

		b, err := hex.DecodeHex(encoded)
		if err != nil {
			return nil, err
		}

		if err := tx.UnmarshalBinary(b); err != nil {
			return nil, err
		}

		tx.State = TxState(state)

		txs = append(txs, *tx)
	}

	return txs, nil
}

func (p *postgresPool) UpdateTxState(hash common.Hash, newState TxState) error {
	// get connection
	sqlDB, err := db.NewSQLDB(p.cfg)
	if err != nil {
		return err
	}
	defer sqlDB.Close() //nolint:errcheck

	// save
	sql := "UPDATE pool.txs SET state = $1 WHERE hash = $2"
	if _, err := sqlDB.Exec(sql, newState, hash.Hex()); err != nil {
		return err
	}
	return nil
}

func (p *postgresPool) CleanUpInvalidAndNonSelectedTxs() error {
	panic("not implemented yet")
}

func (p *postgresPool) SetGasPrice(gasPrice uint64) error {
	// get connection
	sqlDB, err := db.NewSQLDB(p.cfg)
	if err != nil {
		return err
	}
	defer sqlDB.Close() //nolint:errcheck

	// save
	sql := "INSERT INTO pool.gas_price (price, timestamp) VALUES ($1, $2)"
	if _, err := sqlDB.Exec(sql, gasPrice, time.Now().UTC()); err != nil {
		return err
	}
	return nil
}

func (p *postgresPool) GetGasPrice() (uint64, error) {
	// get connection
	sqlDB, err := db.NewSQLDB(p.cfg)
	if err != nil {
		return 0, err
	}
	defer sqlDB.Close() //nolint:errcheck

	// save
	sql := "SELECT price FROM pool.gas_price ORDER BY item_id DESC LIMIT 1"
	rows, err := sqlDB.Query(sql)
	if err != nil {
		return 0, err
	}

	gasPrice := uint64(0)

	for rows.Next() {
		err := rows.Scan(&gasPrice)
		if err != nil {
			return 0, err
		}
	}

	return gasPrice, nil
}
