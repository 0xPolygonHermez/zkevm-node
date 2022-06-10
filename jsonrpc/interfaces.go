package jsonrpc

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/runtime"
)

// Consumer interfaces required by the package.

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
	GetLastConsolidatedBatchNumber(ctx context.Context, txBundleID string) (uint64, error)
	GetTransactionByHash(ctx context.Context, transactionHash common.Hash, txBundleID string) (*types.Transaction, error)
	GetTransactionReceipt(ctx context.Context, transactionHash common.Hash, txBundleID string) (*state.Receipt, error)
	GetLastBatchNumber(ctx context.Context, txBundleID string) (uint64, error)
	GetLastBatch(ctx context.Context, isVirtual bool, txBundleID string) (*state.Batch, error)
	NewBatchProcessor(ctx context.Context, sequencerAddress common.Address, stateRoot []byte, txBundleID string) (BatchProcessorInterface, error)
	EstimateGas(transaction *types.Transaction, senderAddress common.Address) (uint64, error)
	GetBalance(ctx context.Context, address common.Address, batchNumber uint64, txBundleID string) (*big.Int, error)
	GetBatchByHash(ctx context.Context, hash common.Hash, txBundleID string) (*state.Batch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, txBundleID string) (*state.Batch, error)
	GetCode(ctx context.Context, address common.Address, batchNumber uint64, txBundleID string) ([]byte, error)
	GetStorageAt(ctx context.Context, address common.Address, position *big.Int, batchNumber uint64, txBundleID string) (*big.Int, error)
	GetSyncingInfo(ctx context.Context, txBundleID string) (state.SyncingInfo, error)
	GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64, txBundleID string) (*types.Transaction, error)
	GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64, txBundleID string) (*types.Transaction, error)
	GetNonce(ctx context.Context, address common.Address, batchNumber uint64, txBundleID string) (uint64, error)
	GetBatchHeader(ctx context.Context, batchNumber uint64, txBundleID string) (*types.Header, error)
	ReplayTransaction(transactionHash common.Hash, traceMode []string) *runtime.ExecutionResult
	ReplayBatchTransactions(batchNumber uint64, traceMode []string) ([]*runtime.ExecutionResult, error)
	GetBatchTransactionCountByHash(ctx context.Context, hash common.Hash, txBundleID string) (uint64, error)
	GetBatchTransactionCountByNumber(ctx context.Context, batchNumber uint64, txBundleID string) (uint64, error)
	GetLogs(ctx context.Context, fromBatch uint64, toBatch uint64, addresses []common.Address, topics [][]common.Hash, batchHash *common.Hash, since *time.Time, txBundleID string) ([]*types.Log, error)
	GetBatchHashesSince(ctx context.Context, since time.Time, txBundleID string) ([]common.Hash, error)
}

type storageInterface interface {
	NewLogFilter(filter LogFilter) (uint64, error)
	NewBlockFilter() (uint64, error)
	NewPendingTransactionFilter() (uint64, error)
	GetFilter(filterID uint64) (*Filter, error)
	UpdateFilterLastPoll(filterID uint64) error
	UninstallFilter(filterID uint64) (bool, error)
}

type BatchProcessorInterface interface {
	ProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress, sequencerAddress common.Address) *runtime.ExecutionResult
}
