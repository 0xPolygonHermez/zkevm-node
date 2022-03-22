package tree

import (
	"context"
	"fmt"
	"sync"

	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/ristretto"
)

// BadgerRistrettoStore uses BadgerDB with a ristretto cache in front
// (badger's default).
type BadgerRistrettoStore struct {
	db    *badger.DB
	cache *ristretto.Cache
	lock  sync.RWMutex
}

// NewBadgerDB returns a badger db configured to use the given dataDir.
func NewBadgerDB(dataDir string) (*badger.DB, error) {
	const defaultCompactors = 20

	opts := badger.DefaultOptions(dataDir)
	opts.ValueDir = dataDir
	opts.Logger = nil
	opts.WithSyncWrites(false)
	opts.WithNumCompactors(defaultCompactors)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// BeginDBTransaction starts a transaction block
func (b *BadgerRistrettoStore) BeginDBTransaction(ctx context.Context) error {
	return ErrDBTxsNotSupported
}

// Commit commits a db transaction
func (b *BadgerRistrettoStore) Commit(ctx context.Context) error {
	return ErrDBTxsNotSupported
}

// Rollback rollbacks a db transaction
func (b *BadgerRistrettoStore) Rollback(ctx context.Context) error {
	return ErrDBTxsNotSupported
}

// NewBadgerRistrettoStore creates an instance of BadgerRistrettoStore.
func NewBadgerRistrettoStore(db *badger.DB, cache *ristretto.Cache) *BadgerRistrettoStore {
	return &BadgerRistrettoStore{
		db:    db,
		cache: cache,
	}
}

// Get gets value of key, first trying the cache, then the db.
func (b *BadgerRistrettoStore) Get(ctx context.Context, key []byte) (data []byte, err error) {
	value, found := b.cache.Get(key)
	if found {
		data, ok := value.([]byte)
		if !ok {
			return nil, fmt.Errorf("Could not cast data to []byte for key %q", key)
		}
		return data, nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()
	err = b.db.View(func(txn *badger.Txn) error {
		item, e := txn.Get(key)
		if e != nil {
			return e
		}
		data, e = item.ValueCopy(nil)
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	b.cache.Set(key, data, cacheDefaultCost)
	return data, nil
}

// Set inserts a key-value pair into the db.
// If record with such a key already exists its assumed that the value is correct,
// because it's a reverse hash table, and the key is a hash of the value.
func (b *BadgerRistrettoStore) Set(ctx context.Context, key []byte, value []byte) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
}

// Reset clears the db and the cache.
func (b *BadgerRistrettoStore) Reset() error {
	b.cache.Clear()
	return b.db.DropAll()
}
