package broadcast

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

type stateInterface interface {
	GetLastBatch(ctx context.Context, tx pgx.Tx) (*Batch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, tx pgx.Tx) (*Batch, error)
	GetEncodedTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, tx pgx.Tx) (encoded []string, err error)
}

// This should be moved into the state package

// Batch represents a Batch
type Batch struct {
	BatchNumber    uint64
	GlobalExitRoot common.Hash
	RawTxsData     []byte
	Timestamp      time.Time
}
