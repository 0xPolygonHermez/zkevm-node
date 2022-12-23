package ethtxmanager

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresStorage hold txs to be managed
type PostgresStorage struct {
	*pgxpool.Pool
}

// NewPostgresStorage creates a new instance of storage that use
// postgres to store data
func NewPostgresStorage(dbCfg db.Config) (storageInterface, error) {
	db, err := db.NewSQLDB(dbCfg)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db,
	}, nil
}

// Add persist a monitored tx
func (s *PostgresStorage) Add(ctx context.Context, mTx monitoredTx, dbTx pgx.Tx) error {
	conn := s.dbConn(dbTx)
	cmd := `
        INSERT INTO txman.monitored_txs (id, from_addr, to_addr, nonce, value, data, gas, gas_price, status, history, created_at, updated_at)
                                 VALUES ($1,        $2,      $3,    $4,    $5,   $6,  $7,       $8,      $9,     $10,        $11,        $12)`

	_, err := conn.Exec(ctx, cmd,
		mTx.id, mTx.from.String(), mTx.toStringPtr(),
		mTx.nonce, mTx.valueU64Ptr(), mTx.dataStringPtr(),
		mTx.gas, mTx.gasPrice.Uint64(), string(mTx.status),
		mTx.historyStringSlice(), time.Now().UTC().Round(time.Microsecond),
		time.Now().UTC().Round(time.Microsecond))

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.ConstraintName == "monitored_txs_pkey" {
			return ErrAlreadyExists
		} else {
			return err
		}
	}

	return nil
}

// Get loads a persisted monitored tx
func (s *PostgresStorage) Get(ctx context.Context, id string, dbTx pgx.Tx) (monitoredTx, error) {
	conn := s.dbConn(dbTx)
	cmd := `
        SELECT id, from_addr, to_addr, nonce, value, data, gas, gas_price, status, history, created_at, updated_at
          FROM txman.monitored_txs
         WHERE id = $1`

	mTx := monitoredTx{}

	row := conn.QueryRow(ctx, cmd, id)
	err := s.scanMtx(row, &mTx)
	if errors.Is(err, pgx.ErrNoRows) {
		return mTx, ErrNotFound
	} else if err != nil {
		return mTx, err
	}

	return mTx, nil
}

// GetByStatus loads all monitored tx that match the provided status
func (s *PostgresStorage) GetByStatus(ctx context.Context, statuses []MonitoredTxStatus, dbTx pgx.Tx) ([]monitoredTx, error) {
	hasStatusToFilter := len(statuses) > 0

	conn := s.dbConn(dbTx)
	cmd := `
        SELECT id, from_addr, to_addr, nonce, value, data, gas, gas_price, status, history, created_at, updated_at
          FROM txman.monitored_txs`
	if hasStatusToFilter {
		cmd += `
         WHERE status = ANY($1)`
	}
	cmd += `
         ORDER BY created_at`

	mTxs := []monitoredTx{}

	var rows pgx.Rows
	var err error
	if hasStatusToFilter {
		rows, err = conn.Query(ctx, cmd, statuses)
	} else {
		rows, err = conn.Query(ctx, cmd)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return []monitoredTx{}, nil
	}

	for rows.Next() {
		mTx := monitoredTx{}
		err := s.scanMtx(rows, &mTx)
		if err != nil {
			return nil, err
		}
		mTxs = append(mTxs, mTx)
	}

	return mTxs, nil
}

// Update a persisted monitored tx
func (s *PostgresStorage) Update(ctx context.Context, mTx monitoredTx, dbTx pgx.Tx) error {
	conn := s.dbConn(dbTx)
	cmd := `
        UPDATE txman.monitored_txs
           SET from_addr = $2
             , to_addr = $3
             , nonce = $4
             , value = $5
             , data = $6
             , gas = $7
             , gas_price = $8
             , status = $9
             , history = $10
             , updated_at = $11
         WHERE id = $1`

	_, err := conn.Exec(ctx, cmd,
		mTx.id, mTx.from.String(), mTx.toStringPtr(),
		mTx.nonce, mTx.valueU64Ptr(), mTx.dataStringPtr(),
		mTx.gas, mTx.gasPrice.Uint64(), string(mTx.status),
		mTx.historyStringSlice(), time.Now().UTC().Round(time.Microsecond))

	if err != nil {
		return err
	}

	return nil
}

// scanMtx scans a row and fill the provided instance of monitoredTx with
// the row data
func (s *PostgresStorage) scanMtx(row pgx.Row, mTx *monitoredTx) error {
	// id, from, to, nonce, value, data, gas, gas_price, status, history, created_at, updated_at
	var from, status string
	var to, data *string
	var history []string
	var value *uint64
	var gasPrice uint64

	err := row.Scan(&mTx.id, &from, &to, &mTx.nonce, &value, &data, &mTx.gas, &gasPrice, &status, &history, &mTx.createdAt, &mTx.updatedAt)
	if err != nil {
		return err
	}

	mTx.from = common.HexToAddress(from)
	if to != nil {
		tmp := common.HexToAddress(*to)
		mTx.to = &tmp
	}
	if value != nil {
		tmp := *value
		mTx.value = big.NewInt(0).SetUint64(tmp)
	}
	if data != nil {
		tmp := *data
		bytes, err := hex.DecodeString(tmp)
		if err != nil {
			return err
		}
		mTx.data = bytes
	}
	mTx.gasPrice = big.NewInt(0).SetUint64(gasPrice)
	mTx.status = MonitoredTxStatus(status)
	h := make(map[common.Hash]bool, len(history))
	for _, txHash := range history {
		h[common.HexToHash(txHash)] = true
	}
	mTx.history = h

	return nil
}

// dbConn represents an instance of an object that can
// connect to a postgres db to execute sql commands and query data
type dbConn interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// dbConn determines which db connection to use, dbTx or the main pgxpool
func (p *PostgresStorage) dbConn(dbTx pgx.Tx) dbConn {
	if dbTx != nil {
		return dbTx
	}
	return p
}
