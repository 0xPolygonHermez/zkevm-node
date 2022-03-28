package tree

import (
	"context"
	"crypto/sha256"
)

const (
	byte32len = 32
)

type kvPair struct {
	Key   []byte
	Value []byte
}

// MemStore stores key-value pairs in memory
type MemStore struct {
	kv map[[byte32len]byte]kvPair
}

// NewMemStore creates an instance of MemStore
func NewMemStore() *MemStore {
	kv := make(map[[byte32len]byte]kvPair)
	return &MemStore{kv}
}

// SupportsDBTransactions indicates whether the store implementation supports DB transactions
func (m *MemStore) SupportsDBTransactions() bool {
	return false
}

// BeginDBTransaction starts a transaction block
func (m *MemStore) BeginDBTransaction(ctx context.Context) error {
	return ErrDBTxsNotSupported
}

// Commit commits a db transaction
func (m *MemStore) Commit(ctx context.Context) error {
	return ErrDBTxsNotSupported
}

// Rollback rollbacks a db transaction
func (m *MemStore) Rollback(ctx context.Context) error {
	return ErrDBTxsNotSupported
}

// Get gets value of key from the memory
func (m *MemStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	k := sha256.Sum256(key)
	kv, ok := m.kv[k]
	if !ok {
		return nil, ErrNotFound
	}
	return kv.Value, nil
}

// Set sets value of key in the memory
func (m *MemStore) Set(ctx context.Context, key []byte, value []byte) error {
	k := sha256.Sum256(key)
	kv := kvPair{key, value}
	m.kv[k] = kv
	return nil
}

// Reset clears the stored data.
func (m *MemStore) Reset() error {
	m.kv = make(map[[byte32len]byte]kvPair)
	return nil
}
