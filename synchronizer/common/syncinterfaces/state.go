package syncinterfaces

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

// StateInterface contains the methods required to interact with the state.
type StateBeginTransactionInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
}

type StateGetBatchByNumberInterface interface {
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
}
