package pool

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

type storage interface {
	AddTx(ctx context.Context, tx Transaction) error
	CountTransactionsByStatus(ctx context.Context, status TxStatus) (uint64, error)
	DeleteTransactionsByHashes(ctx context.Context, hashes []common.Hash) error
	GetGasPrice(ctx context.Context) (uint64, error)
	GetNonce(ctx context.Context, address common.Address) (uint64, error)
	GetPendingTxHashesSince(ctx context.Context, since time.Time) ([]common.Hash, error)
	GetTxsByFromAndNonce(ctx context.Context, from common.Address, nonce uint64) ([]Transaction, error)
	GetTxsByStatus(ctx context.Context, state TxStatus, isClaims bool, limit uint64) ([]Transaction, error)
	GetNonWIPTxsByStatus(ctx context.Context, status TxStatus, isClaims bool, limit uint64) ([]Transaction, error)
	IsTxPending(ctx context.Context, hash common.Hash) (bool, error)
	SetGasPrice(ctx context.Context, gasPrice uint64) error
	UpdateTxsStatus(ctx context.Context, hashes []string, newStatus TxStatus) error
	UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus TxStatus, isWIP bool) error
	UpdateTxWIPStatus(ctx context.Context, hash common.Hash, isWIP bool) error
	GetTxs(ctx context.Context, filterStatus TxStatus, isClaims bool, minGasPrice, limit uint64) ([]*Transaction, error)
	GetTxFromAddressFromByHash(ctx context.Context, hash common.Hash) (common.Address, uint64, error)
	GetTxByHash(ctx context.Context, hash common.Hash) (*Transaction, error)
	GetTxZkCountersByHash(ctx context.Context, hash common.Hash) (*state.ZKCounters, error)
	DeleteTransactionByHash(ctx context.Context, hash common.Hash) error
	MarkWIPTxsAsPending(ctx context.Context) error
	MinGasPriceSince(ctx context.Context, timestamp time.Time) (uint64, error)
	DepositCountExists(ctx context.Context, depositCount uint64) (bool, error)
}

type stateInterface interface {
	GetBalance(ctx context.Context, address common.Address, batchNumber uint64, dbTx pgx.Tx) (*big.Int, error)
	GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetNonce(ctx context.Context, address common.Address, batchNumber uint64, dbTx pgx.Tx) (uint64, error)
	GetTransactionByHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Transaction, error)
	PreProcessTransaction(ctx context.Context, tx *types.Transaction, dbTx pgx.Tx) (*state.ProcessBatchResponse, error)
	AddEvent(ctx context.Context, event *state.Event, dbTx pgx.Tx) error
}
