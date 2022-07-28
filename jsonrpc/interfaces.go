package jsonrpc

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

	GetLastConsolidatedL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetTransactionByHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Transaction, error)
	GetTransactionReceipt(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Receipt, error)
	GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLastL2Block(ctx context.Context, dbTx pgx.Tx) (*types.Block, error)
	GetLastL2BlockHeader(ctx context.Context, dbTx pgx.Tx) (*types.Header, error)
	EstimateGas(transaction *types.Transaction, senderAddress common.Address) (uint64, error)
	GetBalance(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) (*big.Int, error)
	GetL2BlockByHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (*types.Block, error)
	GetL2BlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*types.Block, error)
	GetCode(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) ([]byte, error)
	GetStorageAt(ctx context.Context, address common.Address, position *big.Int, blockNumber uint64, dbTx pgx.Tx) (*big.Int, error)
	GetSyncingInfo(ctx context.Context, dbTx pgx.Tx) (state.SyncingInfo, error)
	GetTransactionByL2BlockHashAndIndex(ctx context.Context, blockHash common.Hash, index uint64, dbTx pgx.Tx) (*types.Transaction, error)
	GetTransactionByL2BlockNumberAndIndex(ctx context.Context, blockNumber uint64, index uint64, dbTx pgx.Tx) (*types.Transaction, error)
	GetNonce(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) (uint64, error)
	GetL2BlockHeaderByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*types.Header, error)
	GetL2BlockTransactionCountByHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (uint64, error)
	GetL2BlockTransactionCountByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (uint64, error)
	GetLogs(ctx context.Context, fromBlock uint64, toBlock uint64, addresses []common.Address, topics [][]common.Hash, blockHash *common.Hash, since *time.Time, dbTx pgx.Tx) ([]*types.Log, error)
	GetL2BlockHashesSince(ctx context.Context, since time.Time, dbTx pgx.Tx) ([]common.Hash, error)
	DebugTransaction(ctx context.Context, transactionHash common.Hash, tracer string) (*runtime.ExecutionResult, error)
	ProcessUnsignedTransaction(ctx context.Context, tx *types.Transaction, senderAddress common.Address, blockNumber uint64, dbTx pgx.Tx) *runtime.ExecutionResult
	IsBatchConsolidated(ctx context.Context, batchId int, dbTx pgx.Tx) (bool, error)
	IsBatchVirtualized(ctx context.Context, batchId int, dbTx pgx.Tx) (bool, error)
}

type storageInterface interface {
	NewLogFilter(filter LogFilter) (uint64, error)
	NewBlockFilter() (uint64, error)
	NewPendingTransactionFilter() (uint64, error)
	GetFilter(filterID uint64) (*Filter, error)
	UpdateFilterLastPoll(filterID uint64) error
	UninstallFilter(filterID uint64) (bool, error)
}
