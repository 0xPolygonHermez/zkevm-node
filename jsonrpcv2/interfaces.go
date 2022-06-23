package jsonrpcv2

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state/runtime"
	state "github.com/hermeznetwork/hermez-core/statev2"
	"github.com/jackc/pgx/v4"
)

// jsonRPCTxPool contains the methods required to interact with the tx pool.
type jsonRPCTxPool interface {
	AddTx(ctx context.Context, tx types.Transaction) error
	GetPendingTxs(ctx context.Context, isClaims bool, limit uint64) ([]pool.Transaction, error)
	GetGasPrice(ctx context.Context) (uint64, error)
	GetPendingTxHashesSince(ctx context.Context, since time.Time) ([]common.Hash, error)
}

// gasPriceEstimator contains the methods required to interact with gas price estimator
type gasPriceEstimator interface {
	GetAvgGasPrice(ctx context.Context) (*big.Int, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	RollbackState(ctx context.Context, tx pgx.Tx) error
	CommitState(ctx context.Context, tx pgx.Tx) error

	// GetLastConsolidatedBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	// GetTransactionByHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Transaction, error)
	// GetTransactionReceipt(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*state.Receipt, error)
	GetLastBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLastBatch(ctx context.Context, dbTx pgx.Tx) (*state.Batch, error)
	EstimateGas(transaction *types.Transaction, senderAddress common.Address) (uint64, error)
	GetBalance(ctx context.Context, address common.Address, batchNumber uint64, dbTx pgx.Tx) (*big.Int, error)
	// GetBatchByHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (*state.Batch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetCode(ctx context.Context, address common.Address, batchNumber uint64, dbTx pgx.Tx) ([]byte, error)
	GetStorageAt(ctx context.Context, address common.Address, position *big.Int, batchNumber uint64, dbTx pgx.Tx) (*big.Int, error)
	// GetSyncingInfo(ctx context.Context, dbTx pgx.Tx) (state.SyncingInfo, error)
	// GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64, dbTx pgx.Tx) (*types.Transaction, error)
	// GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64, dbTx pgx.Tx) (*types.Transaction, error)
	GetNonce(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) (uint64, error)
	// GetBatchHeader(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*types.Header, error)
	// GetBatchTransactionCountByHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (uint64, error)
	// GetBatchTransactionCountByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (uint64, error)
	// GetLogs(ctx context.Context, fromBatch uint64, toBatch uint64, addresses []common.Address, topics [][]common.Hash, batchHash *common.Hash, since *time.Time, dbTx pgx.Tx) ([]*types.Log, error)
	// GetBatchHashesSince(ctx context.Context, since time.Time, dbTx pgx.Tx) ([]common.Hash, error)
	DebugTransaction(ctx context.Context, transactionHash common.Hash, tracer string) (*runtime.ExecutionResult, error)
	// ProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress, sequencerAddress common.Address, stateRoot []byte, dbTx pgx.Tx) *runtime.ExecutionResult
}

type storageInterface interface {
	NewLogFilter(filter LogFilter) (uint64, error)
	NewBlockFilter() (uint64, error)
	NewPendingTransactionFilter() (uint64, error)
	GetFilter(filterID uint64) (*Filter, error)
	UpdateFilterLastPoll(filterID uint64) error
	UninstallFilter(filterID uint64) (bool, error)
}
