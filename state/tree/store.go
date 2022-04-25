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
	SupportsDBTransactions() bool
	BeginDBTransaction(ctx context.Context, txBundleID string) error
	Commit(ctx context.Context, txBundleID string) error
	Rollback(ctx context.Context, txBundleID string) error
	Get(ctx context.Context, key []byte, txBundleID string) ([]byte, error)
	Set(ctx context.Context, key []byte, value []byte, txBundleID string) error
}
