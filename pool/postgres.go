package pool

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/db"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresPool is an implementation of the Pool interface
// that uses a postgres database to store the data
type PostgresPool struct {
	db *pgxpool.Pool
}

// NewPostgresPool creates and initializes an instance of PostgresPool
func NewPostgresPool(cfg db.Config) (*PostgresPool, error) {
	dbPool, err := db.NewSQLDB(cfg)
	if err != nil {
		return nil, err
	}

	return &PostgresPool{
		db: dbPool,
	}, nil
}

// AddTx adds a transaction to the pool table with pending state
func (p *PostgresPool) AddTx(ctx context.Context, tx types.Transaction) error {
	hash := tx.Hash().Hex()

	b, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	encoded := hex.EncodeToHex(b)

	b, err = tx.MarshalJSON()
	if err != nil {
		return err
	}
	decoded := string(b)

	receivedAt := time.Now()
	sql := "INSERT INTO pool.txs (hash, encoded, decoded, state, received_at) VALUES($1, $2, $3, $4, $5)"
	if _, err := p.db.Exec(ctx, sql, hash, encoded, decoded, TxStatePending, receivedAt); err != nil {
		return err
	}
	return nil
}

// GetPendingTxs returns an array of transactions with all
// the transactions which have the state equals pending
func (p *PostgresPool) GetPendingTxs(ctx context.Context) ([]Transaction, error) {
	sql := "SELECT encoded, state, received_at FROM pool.txs WHERE state = $1"
	rows, err := p.db.Query(ctx, sql, TxStatePending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]Transaction, 0, len(rows.RawValues()))
	for rows.Next() {
		var (
			encoded, state string
			receivedAt     time.Time
		)

		if err := rows.Scan(&encoded, &state, &receivedAt); err != nil {
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
		tx.ReceivedAt = receivedAt
		txs = append(txs, *tx)
	}

	return txs, nil
}

// UpdateTxState updates a transaction state accordingly to the
// provided state and hash
func (p *PostgresPool) UpdateTxState(ctx context.Context, hash common.Hash, newState TxState) error {
	sql := "UPDATE pool.txs SET state = $1 WHERE hash = $2"
	if _, err := p.db.Exec(ctx, sql, newState, hash.Hex()); err != nil {
		return err
	}
	return nil
}

// UpdateTxsState updates transactions state accordingly to the provided state and hashes
func (p *PostgresPool) UpdateTxsState(ctx context.Context, hashes []string, newState TxState) error {
	sql := "UPDATE pool.txs SET state = $1 WHERE hash = ANY ($2)"
	if _, err := p.db.Exec(ctx, sql, newState, hashes); err != nil {
		return err
	}
	return nil
}

// CleanUpInvalidAndNonSelectedTxs removes from the transaction pool table
// the invalid and Non selected transactions
func (p *PostgresPool) CleanUpInvalidAndNonSelectedTxs(ctx context.Context) error {
	panic("not implemented yet")
}

// SetGasPrice allows an external component to define the gas price
func (p *PostgresPool) SetGasPrice(ctx context.Context, gasPrice uint64) error {
	sql := "INSERT INTO pool.gas_price (price, timestamp) VALUES ($1, $2)"
	if _, err := p.db.Exec(ctx, sql, gasPrice, time.Now().UTC()); err != nil {
		return err
	}
	return nil
}

// GetGasPrice returns the current gas price
func (p *PostgresPool) GetGasPrice(ctx context.Context) (uint64, error) {
	sql := "SELECT price FROM pool.gas_price ORDER BY item_id DESC LIMIT 1"
	rows, err := p.db.Query(ctx, sql)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
	} else if err != nil {
		return 0, err
	}

	defer rows.Close()

	gasPrice := uint64(0)

	for rows.Next() {
		err := rows.Scan(&gasPrice)
		if err != nil {
			return 0, err
		}
	}

	return gasPrice, nil
}
