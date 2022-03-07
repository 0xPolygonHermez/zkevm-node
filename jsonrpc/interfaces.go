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
	GetPendingTxs(ctx context.Context) ([]pool.Transaction, error)
	GetGasPrice(ctx context.Context) (uint64, error)
}

// gasPriceEstimator contains the methods required to interact with gas price estimator
type gasPriceEstimator interface {
	GetAvgGasPrice(ctx context.Context) (*big.Int, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastConsolidatedBatchNumber(ctx context.Context) (uint64, error)
	GetTransactionByHash(ctx context.Context, transactionHash common.Hash) (*types.Transaction, error)
	GetTransactionReceipt(ctx context.Context, transactionHash common.Hash) (*state.Receipt, error)
	GetLastBatchNumber(ctx context.Context) (uint64, error)
	GetLastBatch(ctx context.Context, isVirtual bool) (*state.Batch, error)
	NewBatchProcessor(ctx context.Context, sequencerAddress common.Address, lastBatchNumber uint64) (*state.BasicBatchProcessor, error)
	EstimateGas(transaction *types.Transaction) uint64
	GetBalance(ctx context.Context, address common.Address, batchNumber uint64) (*big.Int, error)
	GetBatchByHash(ctx context.Context, hash common.Hash) (*state.Batch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64) (*state.Batch, error)
	GetCode(ctx context.Context, address common.Address, batchNumber uint64) ([]byte, error)
	GetStorageAt(ctx context.Context, address common.Address, position common.Hash, batchNumber uint64) (*big.Int, error)
	GetTransactionByBatchHashAndIndex(ctx context.Context, batchHash common.Hash, index uint64) (*types.Transaction, error)
	GetTransactionByBatchNumberAndIndex(ctx context.Context, batchNumber uint64, index uint64) (*types.Transaction, error)
	GetNonce(ctx context.Context, address common.Address, batchNumber uint64) (uint64, error)
	GetBatchHeader(ctx context.Context, batchNumber uint64) (*types.Header, error)
}
