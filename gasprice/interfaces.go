package gasprice

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

// Consumer interfaces required by the package.

// pool contains methods to interact with the tx pool.
type pool interface {
	SetGasPrice(ctx context.Context, gasPrice uint64) error
	GetGasPrice(ctx context.Context) (uint64, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastBatchNumber(ctx context.Context) (uint64, error)
	GetTxsByBatchNum(ctx context.Context, batchNum uint64) ([]*types.Transaction, error)
}
