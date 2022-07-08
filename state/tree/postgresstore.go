package tree

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/state/store"
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
	*store.Pg

	tableName      string
	constraintName string
}

// NewPostgresStore creates an instance of PostgresStore
func NewPostgresStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{
		Pg: store.NewPg(db),

		tableName:      merkleTreeTable,
		constraintName: mtConstraint,
	}
}

// NewPostgresSCCodeStore creates an instance of PostgresStore
func NewPostgresSCCodeStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{
		Pg: store.NewPg(db),

		tableName:      scCodeTreeTable,
		constraintName: scCodeConstraint,
	}
}

// Get gets value of key from the db
func (p *PostgresStore) Get(ctx context.Context, key []byte, txBundleID string) ([]byte, error) {
	var data []byte
	err := p.QueryRow(ctx, txBundleID, fmt.Sprintf(getNodeByKeySQL, p.tableName), key).Scan(&data)
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
	_, err := p.Exec(ctx, txBundleID, fmt.Sprintf(setNodeByKeySQL, p.tableName, p.constraintName), key, value)
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
	_, err := p.Exec(context.Background(), "", fmt.Sprintf("TRUNCATE TABLE %s;", p.tableName))
	return err
}
