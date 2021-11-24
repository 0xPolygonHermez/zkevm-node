package pool

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Pool represents a pool of transactions
type Pool interface {
	AddTx(ctx context.Context, tx types.Transaction) error
	GetPendingTxs(ctx context.Context) ([]Transaction, error)
	UpdateTxState(ctx context.Context, hash common.Hash, newState TxState) error
	CleanUpInvalidAndNonSelectedTxs(ctx context.Context) error
	GetGasPrice(ctx context.Context) (uint64, error)
}

// NewPool creates a new pool
func NewPool(cfg Config) (Pool, error) {
	return newPostgresPool(cfg)
}
