package tree

import (
	"context"
	"errors"
)

var (
	// ErrNotFound is used when the object in db is not found
	ErrNotFound = errors.New("not found")
)

// Store interface
type Store interface {
	BeginDBTransaction(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Get(ctx context.Context, key []byte) ([]byte, error)
	Set(ctx context.Context, key []byte, value []byte) error
}
