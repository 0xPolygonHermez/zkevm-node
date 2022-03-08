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
