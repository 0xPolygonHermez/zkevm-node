package tree

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/dgraph-io/ristretto"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PgRistrettoStore uses PostgreSQL with a ristretto cache in front.
type PgRistrettoStore struct {
	db             *pgxpool.Pool
	tableName      string
	constraintName string
	cache          *ristretto.Cache

	mu  *sync.Mutex
	txs map[string]pgx.Tx
}

// NewPgRistrettoStore creates an instance of PgRistrettoStore.
func NewPgRistrettoStore(db *pgxpool.Pool, cache *ristretto.Cache) *PgRistrettoStore {
	return &PgRistrettoStore{
		db:             db,
		tableName:      merkleTreeTable,
		constraintName: mtConstraint,
		cache:          cache,

		mu:  new(sync.Mutex),
		txs: make(map[string]pgx.Tx),
	}
}

// NewPgRistrettoSCCodeStore creates an instance of PgRistrettoStore.
func NewPgRistrettoSCCodeStore(db *pgxpool.Pool, cache *ristretto.Cache) *PgRistrettoStore {
	return &PgRistrettoStore{
		db:             db,
		tableName:      scCodeTreeTable,
		constraintName: scCodeConstraint,
		cache:          cache,

		mu:  new(sync.Mutex),
		txs: make(map[string]pgx.Tx),
	}
}

// SupportsDBTransactions indicates whether the store implementation supports DB transactions
func (p *PgRistrettoStore) SupportsDBTransactions() bool {
	return true
}

// BeginDBTransaction starts a transaction block
func (p *PgRistrettoStore) BeginDBTransaction(ctx context.Context, txBundleID string) error {
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
func (p *PgRistrettoStore) Commit(ctx context.Context, txBundleID string) error {
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
func (p *PgRistrettoStore) Rollback(ctx context.Context, txBundleID string) error {
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

func (p *PgRistrettoStore) exec(ctx context.Context, txBundleID string, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error) {
	if txBundleID == "" {
		return p.db.Exec(ctx, sql, arguments...)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	tx, ok := p.txs[txBundleID]
	if !ok {
		return nil, fmt.Errorf("DB Tx bundle %q does not exist", txBundleID)
	}
	return tx.Exec(ctx, sql, arguments...)
}

func (p *PgRistrettoStore) queryRow(ctx context.Context, txBundleID string, sql string, args ...interface{}) pgx.Row {
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

// Get gets value of key, first trying the cache, then the db.
func (p *PgRistrettoStore) Get(ctx context.Context, key []byte, txBundleID string) ([]byte, error) {
	value, found := p.cache.Get(key)
	if found {
		data, ok := value.([]byte)
		if !ok {
			return nil, fmt.Errorf("Could not cast data to []byte for key %q", key)
		}
		return data, nil
	}
	var data []byte
	err := p.queryRow(ctx, txBundleID, fmt.Sprintf(getNodeByKeySQL, p.tableName), key).Scan(&data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	p.cache.Set(key, data, cacheDefaultCost)
	return data, nil
}

// Set inserts a key-value pair into the db.
// If record with such a key already exists its assumed that the value is correct,
// because it's a reverse hash table, and the key is a hash of the value.
func (p *PgRistrettoStore) Set(ctx context.Context, key []byte, value []byte, txBundleID string) error {
	_, err := p.exec(ctx, txBundleID, fmt.Sprintf(setNodeByKeySQL, p.tableName, p.constraintName), key, value)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil
		}
		return err
	}
	return nil
}

// Reset clears the db and the cache.
func (p *PgRistrettoStore) Reset() error {
	p.cache.Clear()

	_, err := p.exec(context.Background(), "", fmt.Sprintf("TRUNCATE TABLE %s;", p.tableName))
	return err
}
