package tree

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	getNodeByKeySQL = "SELECT data FROM state.merkletree WHERE hash = $1"
	setNodeByKeySQL = "INSERT INTO state.merkletree (hash, data) VALUES ($1, $2)"
)

// PostgresStore stores key-value pairs in memory
type PostgresStore struct {
	db *pgxpool.Pool
}

// NewPostgresStore creates an instance of PostgresStore
func NewPostgresStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{db: db}
}

// Get gets value of key from the db
func (p *PostgresStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	var data []byte
	err := p.db.QueryRow(ctx, getNodeByKeySQL, key).Scan(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Set inserts a key-value pair into the db.
// If record with such a key already exists its assumed that the value is correct,
// because it's a reverse hash table, and the key is a hash of the value
func (p *PostgresStore) Set(ctx context.Context, key []byte, value []byte) error {
	_, err := p.db.Exec(ctx, setNodeByKeySQL, key, value)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil
		}
		return err
	}
	return nil
}
