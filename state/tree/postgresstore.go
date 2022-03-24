package tree

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	dbTx           pgx.Tx
	tableName      string
	constraintName string
}

// NewPostgresStore creates an instance of PostgresStore
func NewPostgresStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{db: db, tableName: merkleTreeTable, constraintName: mtConstraint}
}

// NewPostgresSCCodeStore creates an instance of PostgresStore
func NewPostgresSCCodeStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{db: db, tableName: scCodeTreeTable, constraintName: scCodeConstraint}
}

// BeginDBTransaction starts a transaction block
func (p *PostgresStore) BeginDBTransaction(ctx context.Context) error {
	if p.dbTx != nil {
		return ErrAlreadyInitializedDBTransaction
	}

	dbTx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	p.dbTx = dbTx
	return nil
}

// Commit commits a db transaction
func (p *PostgresStore) Commit(ctx context.Context) error {
	if p.dbTx != nil {
		err := p.dbTx.Commit(ctx)
		p.dbTx = nil
		return err
	}

	return ErrNilDBTransaction
}

// Rollback rollbacks a db transaction
func (p *PostgresStore) Rollback(ctx context.Context) error {
	if p.dbTx != nil {
		err := p.dbTx.Rollback(ctx)
		p.dbTx = nil
		return err
	}

	return ErrNilDBTransaction
}

func (p *PostgresStore) exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error) {
	if p.dbTx != nil {
		return p.dbTx.Exec(ctx, sql, arguments...)
	}
	return p.db.Exec(ctx, sql, arguments...)
}

func (p *PostgresStore) queryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if p.dbTx != nil {
		return p.dbTx.QueryRow(ctx, sql, args...)
	}
	return p.db.QueryRow(ctx, sql, args...)
}

// Get gets value of key from the db
func (p *PostgresStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	data := []byte{}
	err := p.queryRow(ctx, fmt.Sprintf(getNodeByKeySQL, p.tableName), key).Scan(&data)
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
func (p *PostgresStore) Set(ctx context.Context, key []byte, value []byte) error {
	_, err := p.exec(ctx, fmt.Sprintf(setNodeByKeySQL, p.tableName, p.constraintName), key, value)
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
	_, err := p.exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE %s;", p.tableName))
	return err
}
