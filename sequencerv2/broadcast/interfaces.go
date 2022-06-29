package broadcast

import (
	"context"

	"github.com/hermeznetwork/hermez-core/statev2"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

type stateInterface interface {
	GetLastBatch(ctx context.Context, dbTx pgx.Tx) (*statev2.Batch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*statev2.Batch, error)
	GetEncodedTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (encoded []string, err error)
}
