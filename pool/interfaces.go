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
	CountTransactionsByState(ctx context.Context, state TxState) (uint64, error)
	DeleteTxsByHashes(ctx context.Context, hashes []common.Hash) error
	GetGasPrice(ctx context.Context) (uint64, error)
	GetNonce(ctx context.Context, address common.Address) (uint64, error)
	GetPendingTxHashesSince(ctx context.Context, since time.Time) ([]common.Hash, error)
	GetTopPendingTxByProfitabilityAndZkCounters(ctx context.Context, maxZkCounters ZkCounters) (*Transaction, error)
	GetTxsByFromAndNonce(ctx context.Context, from common.Address, nonce uint64) ([]Transaction, error)
	GetTxsByState(ctx context.Context, state TxState, isClaims bool, limit uint64) ([]Transaction, error)
	IsTxPending(ctx context.Context, hash common.Hash) (bool, error)
	MarkReorgedTxsAsPending(ctx context.Context) error
	SetGasPrice(ctx context.Context, gasPrice uint64) error
	UpdateTxsState(ctx context.Context, hashes []common.Hash, newState TxState) error
	UpdateTxState(ctx context.Context, hash common.Hash, newState TxState) error
}

type stateInterface interface {
	GetBalance(ctx context.Context, address common.Address, batchNumber uint64, dbTx pgx.Tx) (*big.Int, error)
	GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetNonce(ctx context.Context, address common.Address, batchNumber uint64, dbTx pgx.Tx) (uint64, error)
}
