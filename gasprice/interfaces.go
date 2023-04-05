package gasprice

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

// pool contains methods to interact with the tx pool.
type pool interface {
	SetGasPrice(ctx *context.RequestContext, gasPrice uint64) error
	GetGasPrice(ctx *context.RequestContext) (uint64, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastL2BlockNumber(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	GetTxsByBlockNumber(ctx *context.RequestContext, blockNumber uint64, dbTx pgx.Tx) ([]*types.Transaction, error)
}

// ethermanInterface contains the methods required to interact with ethereum.
type ethermanInterface interface {
	GetL1GasPrice(ctx *context.RequestContext) *big.Int
}
