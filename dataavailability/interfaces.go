package dataavailability

import (
	"context"

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
	// Init initializes the DABackend
	Init() error
	// GetData retrieve the data of a batch from the DA backend. The returned data must be the pre-image of the hash
	GetData(batchNum uint64, hash common.Hash) ([]byte, error)
	// PostSequence sends the sequence data to the data availability backend, and returns the dataAvailabilityMessage
	// as expected by the contract
	PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error)
}
