package types

import (
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/context"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// PoolInterface contains the methods required to interact with the tx pool.
type PoolInterface interface {
	AddTx(ctx *context.RequestContext, tx types.Transaction, ip string) error
	GetGasPrice(ctx *context.RequestContext) (uint64, error)
	GetNonce(ctx *context.RequestContext, address common.Address) (uint64, error)
	GetPendingTxHashesSince(ctx *context.RequestContext, since time.Time) ([]common.Hash, error)
	GetPendingTxs(ctx *context.RequestContext, isClaims bool, limit uint64) ([]pool.Transaction, error)
	CountPendingTransactions(ctx *context.RequestContext) (uint64, error)
	GetTxByHash(ctx *context.RequestContext, hash common.Hash) (*pool.Transaction, error)
}

// StateInterface gathers the methods required to interact with the state.
type StateInterface interface {
	PrepareToHandleNewL2BlockEvents()
	BeginStateTransaction(ctx *context.RequestContext) (pgx.Tx, error)
	DebugTransaction(ctx *context.RequestContext, transactionHash common.Hash, traceConfig state.TraceConfig, dbTx pgx.Tx) (*runtime.ExecutionResult, error)
	EstimateGas(ctx *context.RequestContext, transaction *types.Transaction, senderAddress common.Address, l2BlockNumber *uint64, dbTx pgx.Tx) (uint64, error)
	GetBalance(ctx *context.RequestContext, address common.Address, root common.Hash) (*big.Int, error)
	GetCode(ctx *context.RequestContext, address common.Address, root common.Hash) ([]byte, error)
	GetL2BlockByHash(ctx *context.RequestContext, hash common.Hash, dbTx pgx.Tx) (*types.Block, error)
	GetL2BlockByNumber(ctx *context.RequestContext, blockNumber uint64, dbTx pgx.Tx) (*types.Block, error)
	BatchNumberByL2BlockNumber(ctx *context.RequestContext, blockNumber uint64, dbTx pgx.Tx) (uint64, error)
	GetL2BlockHashesSince(ctx *context.RequestContext, since time.Time, dbTx pgx.Tx) ([]common.Hash, error)
	GetL2BlockHeaderByNumber(ctx *context.RequestContext, blockNumber uint64, dbTx pgx.Tx) (*types.Header, error)
	GetL2BlockTransactionCountByHash(ctx *context.RequestContext, hash common.Hash, dbTx pgx.Tx) (uint64, error)
	GetL2BlockTransactionCountByNumber(ctx *context.RequestContext, blockNumber uint64, dbTx pgx.Tx) (uint64, error)
	GetLastConsolidatedL2BlockNumber(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	GetLastL2Block(ctx *context.RequestContext, dbTx pgx.Tx) (*types.Block, error)
	GetLastL2BlockNumber(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	GetLogs(ctx *context.RequestContext, fromBlock uint64, toBlock uint64, addresses []common.Address, topics [][]common.Hash, blockHash *common.Hash, since *time.Time, dbTx pgx.Tx) ([]*types.Log, error)
	GetNonce(ctx *context.RequestContext, address common.Address, root common.Hash) (uint64, error)
	GetStorageAt(ctx *context.RequestContext, address common.Address, position *big.Int, root common.Hash) (*big.Int, error)
	GetSyncingInfo(ctx *context.RequestContext, dbTx pgx.Tx) (state.SyncingInfo, error)
	GetTransactionByHash(ctx *context.RequestContext, transactionHash common.Hash, dbTx pgx.Tx) (*types.Transaction, error)
	GetTransactionByL2BlockHashAndIndex(ctx *context.RequestContext, blockHash common.Hash, index uint64, dbTx pgx.Tx) (*types.Transaction, error)
	GetTransactionByL2BlockNumberAndIndex(ctx *context.RequestContext, blockNumber uint64, index uint64, dbTx pgx.Tx) (*types.Transaction, error)
	GetTransactionReceipt(ctx *context.RequestContext, transactionHash common.Hash, dbTx pgx.Tx) (*types.Receipt, error)
	IsL2BlockConsolidated(ctx *context.RequestContext, blockNumber uint64, dbTx pgx.Tx) (bool, error)
	IsL2BlockVirtualized(ctx *context.RequestContext, blockNumber uint64, dbTx pgx.Tx) (bool, error)
	ProcessUnsignedTransaction(ctx *context.RequestContext, tx *types.Transaction, senderAddress common.Address, l2BlockNumber *uint64, noZKEVMCounters bool, dbTx pgx.Tx) (*runtime.ExecutionResult, error)
	RegisterNewL2BlockEventHandler(ctx *context.RequestContext, h state.NewL2BlockEventHandler)
	GetLastVirtualBatchNum(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	GetLastVerifiedBatch(ctx *context.RequestContext, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetLastBatchNumber(ctx *context.RequestContext, dbTx pgx.Tx) (uint64, error)
	GetBatchByNumber(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetTransactionsByBatchNumber(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) (txs []types.Transaction, err error)
	GetVirtualBatch(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) (*state.VirtualBatch, error)
	GetVerifiedBatch(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetExitRootByGlobalExitRoot(ctx *context.RequestContext, ger common.Hash, dbTx pgx.Tx) (*state.GlobalExitRoot, error)
}
