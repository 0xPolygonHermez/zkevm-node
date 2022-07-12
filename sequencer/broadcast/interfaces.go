package broadcast

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

type stateInterface interface {
	GetLastBatch(ctx context.Context, dbTx pgx.Tx) (*state.Batch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetEncodedTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (encoded []string, err error)
	GetForcedBatchByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.ForcedBatch, error)
}
