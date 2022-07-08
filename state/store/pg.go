package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Pg contains basic postgresql store functionality.
type Pg struct {
	db *pgxpool.Pool

	mu  *sync.Mutex
	txs map[string]pgx.Tx
}

// NewPg creates a new postgres store.
func NewPg(db *pgxpool.Pool) *Pg {
	return &Pg{
		db: db,

		mu:  new(sync.Mutex),
		txs: make(map[string]pgx.Tx),
	}
}

// SupportsDBTransactions indicates whether the store implementation supports DB
// transactions.
func (p *Pg) SupportsDBTransactions() bool {
	return true
}

// BeginDBTransaction starts a transaction block
func (p *Pg) BeginDBTransaction(ctx context.Context, txBundleID string) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.txs[txBundleID]
	if ok {
		return fmt.Errorf("DB Tx bundle %q already exists", txBundleID)
	}
	p.txs[txBundleID] = tx

	return nil
}

// Commit commits a db transaction
func (p *Pg) Commit(ctx context.Context, txBundleID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	tx, ok := p.txs[txBundleID]

	if !ok {
		return fmt.Errorf("DB Tx bundle %q does not exist", txBundleID)
	}

	err := tx.Commit(ctx)
	delete(p.txs, txBundleID)

	return err
}

// Rollback rollbacks a db transaction
func (p *Pg) Rollback(ctx context.Context, txBundleID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	tx, ok := p.txs[txBundleID]
	if !ok {
		return fmt.Errorf("DB Tx bundle %q does not exist", txBundleID)
	}

	err := tx.Rollback(ctx)
	delete(p.txs, txBundleID)

	return err
}

// Exec executes a sql query for a given tx bundle.
func (p *Pg) Exec(ctx context.Context, txBundleID string, sql string, args ...interface{}) (commandTag pgconn.CommandTag, err error) {
	if txBundleID == "" {
		return p.db.Exec(ctx, sql, args...)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	tx, ok := p.txs[txBundleID]
	if !ok {
		return nil, fmt.Errorf("DB Tx bundle %q does not exist", txBundleID)
	}
	return tx.Exec(ctx, sql, args...)
}

// QueryRow executes a query that is expected to return at most one row
// (pgx.Row) for the given tx bundle.
func (p *Pg) QueryRow(ctx context.Context, txBundleID string, sql string, args ...interface{}) pgx.Row {
	if txBundleID == "" {
		return p.db.QueryRow(ctx, sql, args...)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	tx, ok := p.txs[txBundleID]
	if !ok {
		log.Errorf("DB Tx bundle %q does not exist", txBundleID)
	}
	return tx.QueryRow(ctx, sql, args...)
}

// Query executes a query that returns pgx.Rows for the given tx bundle.
func (p *Pg) Query(ctx context.Context, txBundleID string, sql string, args ...interface{}) (pgx.Rows, error) {
	if txBundleID == "" {
		return p.db.Query(ctx, sql, args...)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	tx, ok := p.txs[txBundleID]
	if !ok {
		log.Errorf("DB Tx bundle %q does not exist", txBundleID)
	}
	return tx.Query(ctx, sql, args...)
}
