package jsonrpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

// Consumer interfaces required by the package.

// jsonRPCTxPool contains the methods required to interact with the tx pool.
type jsonRPCTxPool interface {
	AddTx(ctx context.Context, tx types.Transaction) error
	GetPendingTxs(ctx context.Context, isClaims bool, limit uint64) ([]pool.Transaction, error)
	GetGasPrice(ctx context.Context) (uint64, error)
}

// gasPriceEstimator contains the methods required to interact with gas price estimator
type gasPriceEstimator interface {
	GetAvgGasPrice(ctx context.Context) (*big.Int, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastConsolidatedBatchNumber(ctx context.Context, txundleID string) (uint64, error)
	GetTransactionByHash(ctx context.Context, transactionHash common.Hash, txundleID string) (*types.Transaction, error)
	GetTransactionReceipt(ctx context.Context, transactionHash common.Hash, txundleID string) (*state.Receipt, error)
	GetLastBatchNumber(ctx context.Context, txundleID string) (uint64, error)
	GetLastBatch(ctx context.Context, isVirtual bool, txundleID string) (*state.Batch, error)
	NewBatchProcessor(ctx context.Context, sequencerAddress common.Address, stateRoot []byte, txundleID string) (*state.BatchProcessor, error)
	EstimateGas(transaction *types.Transaction, txundleID string) (uint64, error)
	GetBalance(ctx context.Context, address common.Address, batchNumber uint64, txundleID string) (*big.Int, error)
	GetBatchByHash(ctx context.Context, hash common.Hash, txundleID string) (*state.Batch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, txundleID string) (*state.Batch, error)
	GetCode(ctx context.Context, address common.Address, batchNumber uint64, txundleID string) ([]byte, error)
	GetStorageAt(ctx context.Context, address common.Address, position *big.Int, batchNumber uint64, txundleID string) (*big.Int, error)
	GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64, txundleID string) (*types.Transaction, error)
	GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64, txundleID string) (*types.Transaction, error)
	GetNonce(ctx context.Context, address common.Address, batchNumber uint64, txundleID string) (uint64, error)
	GetBatchHeader(ctx context.Context, batchNumber uint64, txundleID string) (*types.Header, error)
	GetLogs(ctx context.Context, fromBatch uint64, toBatch uint64, addresses []common.Address, topics [][]common.Hash, batchHash *common.Hash, txundleID string) ([]*types.Log, error)
}
