package sequencer

import (
	"context"
	"time"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

// txPool contains the methods required to interact with the tx pool.
type txPool interface {
	DeleteTxsByHashes(ctx context.Context, hashes []common.Hash) error
}

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	EstimateGasSequenceBatches(sequences []ethmanTypes.Sequence) (*types.Transaction, error)
	TrustedSequencer() (common.Address, error)
	GetLatestBatchNumber() (uint64, error)
	GetLastBatchTimestamp() (uint64, error)
	GetLatestBlockTimestamp(ctx context.Context) (uint64, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetTimeForLatestBatchVirtualization(ctx context.Context, dbTx pgx.Tx) (time.Time, error)
	GetTxsOlderThanNL1Blocks(ctx context.Context, nL1Blocks uint64, dbTx pgx.Tx) ([]common.Hash, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (txs []types.Transaction, err error)
	IsBatchClosed(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (bool, error)
}

type txManager interface {
	SequenceBatches(ctx context.Context, sequences []ethmanTypes.Sequence) error
}

type workerInterface interface {
	GetBestFittingTx(resources RemainingResources) *TxTracker
}

type dbManagerInterface interface {
}

type executorInterface interface {
	ExecuteTransaction(processBatchRequest state.ProcessBatchRequest) state.ProcessBatchResponse
}
