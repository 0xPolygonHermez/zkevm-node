package dataavailability

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type stateInterface interface {
	GetBatchL2DataByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]byte, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
}

// DABackender is the interface for a DA backend
type DABackender interface {
	Init() error
	GetData(batchNum uint64, hash common.Hash) ([]byte, error)
	PostSequence(ctx context.Context, sequences []types.Sequence) ([]byte, error)
}
