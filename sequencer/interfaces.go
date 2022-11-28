package sequencer

import (
	"context"
	"math/big"
	"time"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

// txPool contains the methods required to interact with the tx pool.
type txPool interface {
	GetPendingTxs(ctx context.Context, isClaims bool, limit uint64) ([]pool.Transaction, error)
	UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus pool.TxStatus) error
	UpdateTxsStatus(ctx context.Context, hashes []string, newStatus pool.TxStatus) error
	IsTxPending(ctx context.Context, hash common.Hash) (bool, error)
	DeleteTxsByHashes(ctx context.Context, hashes []common.Hash) error
	MarkReorgedTxsAsPending(ctx context.Context) error
	GetTxs(ctx context.Context, filterStatus pool.TxStatus, isClaims bool, minGasPrice, limit uint64) ([]*pool.Transaction, error)
	GetTxFromAddressFromByHash(ctx context.Context, hash common.Hash) (common.Address, uint64, error)
	IncrementFailedCounter(ctx context.Context, hashes []string) error
}

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	EstimateGasSequenceBatches(sequences []ethmanTypes.Sequence) (*types.Transaction, error)
	GetSendSequenceFee() (*big.Int, error)
	TrustedSequencer() (common.Address, error)
	GetLatestBatchNumber() (uint64, error)
	GetLatestBlockNumber(ctx context.Context) (uint64, error)
	GetLastBatchTimestamp() (uint64, error)
	GetLatestBlockTimestamp(ctx context.Context) (uint64, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLatestGlobalExitRoot(ctx context.Context, dbTx pgx.Tx) (*state.GlobalExitRoot, error)
	GetTimeForLatestBatchVirtualization(ctx context.Context, dbTx pgx.Tx) (time.Time, error)
	GetTxsOlderThanNL1Blocks(ctx context.Context, nL1Blocks uint64, dbTx pgx.Tx) ([]common.Hash, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (txs []types.Transaction, err error)
	IsBatchClosed(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (bool, error)
	IsBatchVirtualized(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (bool, error)

	GetLastBatch(ctx context.Context, dbTx pgx.Tx) (*state.Batch, error)
	GetLastBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)

	StoreTransactions(ctx context.Context, batchNum uint64, processedTxs []*state.ProcessTransactionResponse, dbTx pgx.Tx) error
	CloseBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	ProcessSequencerBatch(ctx context.Context, batchNumber uint64, txs []types.Transaction, dbTx pgx.Tx) (*state.ProcessBatchResponse, error)

	UpdateGERInOpenBatch(ctx context.Context, ger common.Hash, dbTx pgx.Tx) error
	GetBlockNumAndMainnetExitRootByGER(ctx context.Context, ger common.Hash, dbTx pgx.Tx) (uint64, common.Hash, error)
	GetStateRootByBatchNumber(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (common.Hash, error)

	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)

	GetNonce(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) (uint64, error)
	GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
}

type txManager interface {
	SequenceBatches(sequences []ethmanTypes.Sequence) error
}

// priceGetter is for getting eth/matic price, used for the tx profitability checker
type priceGetter interface {
	Start(ctx context.Context)
	GetEthToMaticPrice(ctx context.Context) (*big.Float, error)
}

// gasPriceEstimator contains the methods required to interact with gas price estimator
type gasPriceEstimator interface {
	GetAvgGasPrice(ctx context.Context) (*big.Int, error)
}
