package tree

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	getNodeByKeySQL = "SELECT data FROM %s WHERE hash = $1"
	setNodeByKeySQL = "INSERT INTO %s (hash, data) VALUES ($1, $2)"
)

const (
	merkleTreeTable = "state.merkletree"
	scCodeTreeTable = "state.sc_code"
)

// PostgresStore stores key-value pairs in memory
type PostgresStore struct {
	db        *pgxpool.Pool
	tableName string
}

// NewPostgresStore creates an instance of PostgresStore
func NewPostgresStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{db: db, tableName: merkleTreeTable}
}

// NewPostgresSCCodeStore creates an instance of PostgresStore
func NewPostgresSCCodeStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{db: db, tableName: scCodeTreeTable}
}

// Get gets value of key from the db
func (p *PostgresStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	var data []byte
	err := p.db.QueryRow(ctx, fmt.Sprintf(getNodeByKeySQL, p.tableName), key).Scan(&data)
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
	_, err := p.db.Exec(ctx, fmt.Sprintf(setNodeByKeySQL, p.tableName), key, value)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil
		}
		return err
	}
	return nil
}
