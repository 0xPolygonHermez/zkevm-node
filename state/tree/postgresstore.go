package tree

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	getNodeByKeySQL = "SELECT COALESCE(data, null) FROM %s WHERE hash = $1"
	setNodeByKeySQL = "INSERT INTO %s (hash, data) VALUES ($1, $2) ON CONFLICT ON CONSTRAINT %s DO NOTHING;"
)

const (
	merkleTreeTable  = "state.merkletree"
	scCodeTreeTable  = "state.sc_code"
	mtConstraint     = "merkletree_pkey"
	scCodeConstraint = "sc_code_pkey"
)

var (
	// ErrNilDBTransaction indicates the db transaction has not been properly initialized
	ErrNilDBTransaction = errors.New("database transaction not properly initialized")
	// ErrAlreadyInitializedDBTransaction indicates the db transaction was already initialized
	ErrAlreadyInitializedDBTransaction = errors.New("database transaction already initialized")
)

// PostgresStore stores key-value pairs in memory
type PostgresStore struct {
	db             *pgxpool.Pool
	tableName      string
	constraintName string

	mu  *sync.Mutex
	txs map[string]pgx.Tx
}

// NewPostgresStore creates an instance of PostgresStore
func NewPostgresStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{
		db: db, tableName: merkleTreeTable,
		constraintName: mtConstraint,

		mu:  new(sync.Mutex),
		txs: make(map[string]pgx.Tx),
	}
}

// NewPostgresSCCodeStore creates an instance of PostgresStore
func NewPostgresSCCodeStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{
		db:             db,
		tableName:      scCodeTreeTable,
		constraintName: scCodeConstraint,

		mu:  new(sync.Mutex),
		txs: make(map[string]pgx.Tx),
	}
}

// SupportsDBTransactions indicates whether the store implementation supports DB transactions
func (p *PostgresStore) SupportsDBTransactions() bool {
	return true
}

// BeginDBTransaction starts a transaction block
func (p *PostgresStore) BeginDBTransaction(ctx context.Context, txBundleID string) error {
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
func (p *PostgresStore) Commit(ctx context.Context, txBundleID string) error {
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
func (p *PostgresStore) Rollback(ctx context.Context, txBundleID string) error {
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

func (p *PostgresStore) queryRow(ctx context.Context, txBundleID string, sql string, args ...interface{}) pgx.Row {
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

func (p *PostgresStore) exec(ctx context.Context, txBundleID string, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error) {
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

// Get gets value of key from the db
func (p *PostgresStore) Get(ctx context.Context, key []byte, txBundleID string) ([]byte, error) {
	var data []byte
	err := p.queryRow(ctx, txBundleID, fmt.Sprintf(getNodeByKeySQL, p.tableName), key).Scan(&data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return data, nil
}

// Set inserts a key-value pair into the db.
// If record with such a key already exists its assumed that the value is correct,
// because it's a reverse hash table, and the key is a hash of the value
func (p *PostgresStore) Set(ctx context.Context, key []byte, value []byte, txBundleID string) error {
	_, err := p.exec(ctx, txBundleID, fmt.Sprintf(setNodeByKeySQL, p.tableName, p.constraintName), key, value)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil
		}
		return err
	}
	return nil
}

// Reset clears the db.
func (p *PostgresStore) Reset() error {
	_, err := p.exec(context.Background(), "", fmt.Sprintf("TRUNCATE TABLE %s;", p.tableName))
	return err
}
