package jsonrpc

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
)

// Consumer interfaces required by the package.

// jsonrpcTxPool contains the methods required to interact with the tx pool.
type jsonRPCTxPool interface {
	AddTx(ctx context.Context, tx types.Transaction) error
	GetPendingTxs(ctx context.Context) ([]pool.Transaction, error)
	GetGasPrice(ctx context.Context) (uint64, error)
}
