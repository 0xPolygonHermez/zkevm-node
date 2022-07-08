package pool

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

type storage interface {
	AddTx(ctx context.Context, tx Transaction) error
	GetTxsByState(ctx context.Context, state TxState, isClaims bool, limit uint64) ([]Transaction, error)
	UpdateTxState(ctx context.Context, hash common.Hash, newState TxState) error
	UpdateTxsState(ctx context.Context, hashes []common.Hash, newState TxState) error
	SetGasPrice(ctx context.Context, gasPrice uint64) error
	GetGasPrice(ctx context.Context) (uint64, error)
	CountTransactionsByState(ctx context.Context, state TxState) (uint64, error)
	GetPendingTxHashesSince(ctx context.Context, since time.Time) ([]common.Hash, error)
	IsTxPending(ctx context.Context, hash common.Hash) (bool, error)
	DeleteTxsByHashes(ctx context.Context, hashes []common.Hash) error
	MarkReorgedTxsAsPending(ctx context.Context) error
	GetTopPendingTxByProfitabilityAndZkCounters(ctx context.Context, maxZkCounters ZkCounters) (*Transaction, error)
}

type stateInterface interface {
	GetLastBatch(ctx context.Context, isVirtual bool, txBundleID string) (*state.Batch, error)
	GetNonce(ctx context.Context, address common.Address, batchNumber uint64, txBundleID string) (uint64, error)
	GetBalance(ctx context.Context, address common.Address, batchNumber uint64, txBundleID string) (*big.Int, error)
}
