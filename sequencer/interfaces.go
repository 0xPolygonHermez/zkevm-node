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
	OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetTimeForLatestBatchVirtualization(ctx context.Context, dbTx pgx.Tx) (time.Time, error)
	GetTxsOlderThanNL1Blocks(ctx context.Context, nL1Blocks uint64, dbTx pgx.Tx) ([]common.Hash, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (txs []types.Transaction, err error)
	IsBatchClosed(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (bool, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	ProcessSingleTx(request state.ProcessSingleTxRequest) state.ProcessBatchResponse
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
}

type txManager interface {
	SequenceBatches(ctx context.Context, sequences []ethmanTypes.Sequence) error
}

type workerInterface interface {
	GetBestFittingTx(resources BatchResources) *TxTracker
	UpdateAfterSingleSuccessfulTxExecution(from common.Address, touchedAddresses map[common.Address]*state.TouchedAddress)
	UpdateTx(txHash common.Hash, from common.Address, ZKCounters state.ZKCounters)
}

// The dbManager will need to handle the errors inside the functions which don't return error as they will be used async in the other abstractions.
// Also if dbTx is missing this needs also to be handled in the dbManager
type dbManagerInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	CreateFirstBatch(ctx context.Context, sequencerAddress common.Address) state.ProcessingContext
	GetLastBatchNumber(ctx context.Context) (uint64, error)
	StoreProcessedTransaction(ctx context.Context, dbTx pgx.Tx, batchNumber uint64, processedTx *state.ProcessTransactionResponse) error
	DeleteTxFromPool(ctx context.Context, dbTx pgx.Tx, txHash common.Hash) error
	StoreProcessedTxAndDeleteFromPool(ctx context.Context, batchNumber uint64, response *state.ProcessTransactionResponse)
	CloseBatch(ctx context.Context, params ClosingBatchParameters)
	GetWIPBatch(ctx context.Context) (wipBatch, error)
	GetLastBatch(ctx context.Context) (state.Batch, error)
	GetLastNBatches(ctx context.Context, numBatches uint) ([]*state.Batch, error)
	GetLastClosedBatch(ctx context.Context) (state.Batch, error)
	IsBatchClosed(ctx context.Context, batchNum uint64) (bool, error)
	MarkReorgedTxsAsPending(ctx context.Context) error
}
