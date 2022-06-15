package state

import (
	"context"
)

// Consumer interfaces required by the package.

// statetree contains the methods required to interact with the Merkle tree.
type statetree interface {
	BeginDBTransaction(ctx context.Context, txBundleID string) error
	Commit(ctx context.Context, txBundleID string) error
	Rollback(ctx context.Context, txBundleID string) error
}
