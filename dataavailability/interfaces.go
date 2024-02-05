package dataavailability

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type stateInterface interface {
	GetBatchL2DataByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]byte, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
}

// BatchDataProvider is used to retrieve batch data
type BatchDataProvider interface {
	// GetBatchL2Data retrieve the data of a batch from the DA backend. The returned data must be the pre-image of the hash
	GetBatchL2Data(batchNum uint64, hash common.Hash) ([]byte, error)
}

// SequenceSender is used to send provided sequence of batches
type SequenceSender interface {
	// PostSequence sends the sequence data to the data availability backend, and returns the dataAvailabilityMessage
	// as expected by the contract
	PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error)
}

// DABackender is the interface needed to implement in order to
// integrate a DA service
type DABackender interface {
	BatchDataProvider
	SequenceSender
	// Init initializes the DABackend
	Init() error
}

// ZKEVMClientTrustedBatchesGetter contains the methods required to interact with zkEVM-RPC
type ZKEVMClientTrustedBatchesGetter interface {
	BatchByNumber(ctx context.Context, number *big.Int) (*types.Batch, error)
}
