package tree

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dgraph-io/ristretto"
	"github.com/hermeznetwork/hermez-core/state/store"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PgRistrettoStore uses PostgreSQL with a ristretto cache in front.
type PgRistrettoStore struct {
	*store.Pg

	tableName      string
	constraintName string
	cache          *ristretto.Cache
}

// NewPgRistrettoStore creates an instance of PgRistrettoStore.
func NewPgRistrettoStore(db *pgxpool.Pool, cache *ristretto.Cache) *PgRistrettoStore {
	return &PgRistrettoStore{
		Pg: store.NewPg(db),

		tableName:      merkleTreeTable,
		constraintName: mtConstraint,
		cache:          cache,
	}
}

// NewPgRistrettoSCCodeStore creates an instance of PgRistrettoStore.
func NewPgRistrettoSCCodeStore(db *pgxpool.Pool, cache *ristretto.Cache) *PgRistrettoStore {
	return &PgRistrettoStore{
		Pg: store.NewPg(db),

		tableName:      scCodeTreeTable,
		constraintName: scCodeConstraint,
		cache:          cache,
	}
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
	err := p.QueryRow(ctx, txBundleID, fmt.Sprintf(getNodeByKeySQL, p.tableName), key).Scan(&data)
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
	_, err := p.Exec(ctx, txBundleID, fmt.Sprintf(setNodeByKeySQL, p.tableName, p.constraintName), key, value)
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

	_, err := p.Exec(context.Background(), "", fmt.Sprintf("TRUNCATE TABLE %s;", p.tableName))
	return err
}
