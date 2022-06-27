//nolint
package sequencerv2

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethmanTypes "github.com/hermeznetwork/hermez-core/ethermanv2/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/statev2"
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
	EstimateGasSequenceBatches(sequences []ethmanTypes.Sequence) (*big.Int, error)
	GetSendSequenceFee() (*big.Int, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastVirtualBatchNum(ctx context.Context) (uint64, error)
	GetLastBatchNumberSeenOnEthereum(ctx context.Context) (uint64, error)
	GetLatestGlobalExitRoot(ctx context.Context, dbTx pgx.Tx) (*statev2.GlobalExitRoot, error)

	ProcessBatch(ctx context.Context, txs []types.Transaction) (*statev2.ProcessBatchResponse, error)
	GetLastSendSequenceTime(ctx context.Context) (time.Time, error)
	GetNumberOfBlocksSinceLastGERUpdate(ctx context.Context) (uint64, error)
	GetLastBatchTime(ctx context.Context) (time.Time, error)
}

type txManager interface {
	SequenceBatches(sequences []ethmanTypes.Sequence) error
}

// priceGetter is for getting eth/matic price, used for the tx profitability checker
type priceGetter interface {
	Start(ctx context.Context)
	GetEthToMaticPrice(ctx context.Context) (*big.Float, error)
}
