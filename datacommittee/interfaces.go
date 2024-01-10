package datacommittee

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

type stateInterface interface {
	GetBatchL2DataByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]byte, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
}
