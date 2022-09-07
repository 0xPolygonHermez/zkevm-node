package pool

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type storage interface {
	AddTx(ctx context.Context, tx Transaction) error
	CountTransactionsByStatus(ctx context.Context, status TxStatus) (uint64, error)
	DeleteTxsByHashes(ctx context.Context, hashes []common.Hash) error
	GetGasPrice(ctx context.Context) (uint64, error)
	GetNonce(ctx context.Context, address common.Address) (uint64, error)
	GetPendingTxHashesSince(ctx context.Context, since time.Time) ([]common.Hash, error)
	GetTxsByFromAndNonce(ctx context.Context, from common.Address, nonce uint64) ([]Transaction, error)
	GetTxsByStatus(ctx context.Context, state TxStatus, isClaims bool, limit uint64) ([]Transaction, error)
	IsTxPending(ctx context.Context, hash common.Hash) (bool, error)
	MarkReorgedTxsAsPending(ctx context.Context) error
	SetGasPrice(ctx context.Context, gasPrice uint64) error
	UpdateTxsStatus(ctx context.Context, hashes []string, newStatus TxStatus) error
	UpdateTxStatus(ctx context.Context, hash common.Hash, newStatus TxStatus) error
	GetTopPendingTxByProfitabilityAndZkCounters(ctx context.Context, maxZkCounters ZkCounters, limit uint64) ([]*Transaction, error)
}

type stateInterface interface {
	GetBalance(ctx context.Context, address common.Address, batchNumber uint64, dbTx pgx.Tx) (*big.Int, error)
	GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetNonce(ctx context.Context, address common.Address, batchNumber uint64, dbTx pgx.Tx) (uint64, error)
}
