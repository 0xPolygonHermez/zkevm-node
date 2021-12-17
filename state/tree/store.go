package tree

import (
	"context"
	"errors"
)

var (
	// ErrNotFound is used when the object in db is not found
	ErrNotFound = errors.New("not found")
)

type Store interface {
	Get(ctx context.Context, key []byte) ([]byte, error)
	Set(ctx context.Context, key []byte, value []byte) error
}
