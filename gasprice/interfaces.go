package gasprice

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

// pool contains methods to interact with the tx pool.
type pool interface {
	SetGasPrice(ctx context.Context, gasPrice uint64) error
	GetGasPrice(ctx context.Context) (uint64, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetTxsByBlockNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) ([]*types.Transaction, error)
}
