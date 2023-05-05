package sequencesender

import (
	"context"
	"math/big"
	"time"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	BuildSequenceBatchesTxData(sender common.Address, sequences []ethmanTypes.Sequence) (to *common.Address, data []byte, err error)
	EstimateGasSequenceBatches(sender common.Address, sequences []ethmanTypes.Sequence) (*types.Transaction, error)
	GetLastBatchTimestamp() (uint64, error)
	GetLatestBlockTimestamp(ctx context.Context) (uint64, error)
	GetLatestBatchNumber() (uint64, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	IsBatchClosed(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (bool, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetForcedBatch(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (*state.ForcedBatch, error)
	GetTimeForLatestBatchVirtualization(ctx context.Context, dbTx pgx.Tx) (time.Time, error)
	GetLastBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
}

type ethTxManager interface {
	Add(ctx context.Context, owner, id string, from common.Address, to *common.Address, value *big.Int, data []byte, dbTx pgx.Tx) error
	ProcessPendingMonitoredTxs(ctx context.Context, owner string, failedResultHandler ethtxmanager.ResultHandler, dbTx pgx.Tx)
}
