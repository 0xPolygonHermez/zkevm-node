//nolint
package sequencerv2

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethmanTypes "github.com/hermeznetwork/hermez-core/ethermanv2/types"
	pool "github.com/hermeznetwork/hermez-core/poolv2"
	state "github.com/hermeznetwork/hermez-core/statev2"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

// txPool contains the methods required to interact with the tx pool.
type txPool interface {
	GetPendingTxs(ctx context.Context, isClaims bool, limit uint64) ([]pool.Transaction, error)
	UpdateTxState(ctx context.Context, hash common.Hash, newState pool.TxState) error
	UpdateTxsState(ctx context.Context, hashes []common.Hash, newState pool.TxState) error
	SetGasPrice(ctx context.Context, gasPrice uint64) error
	IsTxPending(ctx context.Context, hash common.Hash) (bool, error)
}

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	EstimateGasSequenceBatches(sequences []ethmanTypes.Sequence) (uint64, error)
	GetSendSequenceFee() (*big.Int, error)
	TrustedSequencer() (common.Address, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLastBatchNumberSeenOnEthereum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLatestGlobalExitRoot(ctx context.Context, dbTx pgx.Tx) (*state.GlobalExitRoot, error)

	GetLastBatch(ctx context.Context, dbTx pgx.Tx) (*state.Batch, error)
	GetLastBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	StoreBatchHeader(ctx context.Context, batch state.Batch, dbTx pgx.Tx) error
	StoreTransactions(ctx context.Context, batchNum uint64, processedTxs []*state.ProcessTransactionResponse, dbTx pgx.Tx) error
	CloseBatch(ctx context.Context, batchNum uint64, stateRoot, localExitRoot common.Hash, dbTx pgx.Tx) error
	ProcessBatch(ctx context.Context, batchNumber uint64, txs []types.Transaction, dbTx pgx.Tx) (*state.ProcessBatchResponse, error)
	GetTimeForLatestBatchVirtualization(ctx context.Context, dbTx pgx.Tx) (time.Time, error)
	GetNumberOfBlocksSinceLastGERUpdate(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLastBatchTime(ctx context.Context, dbTx pgx.Tx) (time.Time, error)
}

type txManager interface {
	SequenceBatches(sequences []ethmanTypes.Sequence) error
}

// priceGetter is for getting eth/matic price, used for the tx profitability checker
type priceGetter interface {
	Start(ctx context.Context)
	GetEthToMaticPrice(ctx context.Context) (*big.Float, error)
}
