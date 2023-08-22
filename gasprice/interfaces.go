package gasprice

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

// poolInterface contains methods to interact with the tx poolInterface.
type poolInterface interface {
	SetGasPrices(ctx context.Context, l2GasPrice uint64, l1GasPrice uint64) error
	GetGasPrices(ctx context.Context) (pool.GasPrices, error)
	DeleteGasPricesHistoryOlderThan(ctx context.Context, date time.Time) error
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetTxsByBlockNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) ([]*types.Transaction, error)
}

// ethermanInterface contains the methods required to interact with ethereum.
type ethermanInterface interface {
	GetL1GasPrice(ctx context.Context) *big.Int
}
