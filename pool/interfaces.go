package pool

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
)

type storage interface {
	AddTx(ctx context.Context, tx types.Transaction, state TxState) error
	GetTxsByState(ctx context.Context, state TxState, limit uint64) ([]Transaction, error)
	UpdateTxState(ctx context.Context, hash common.Hash, newState TxState) error
	UpdateTxsState(ctx context.Context, hashes []common.Hash, newState TxState) error
	SetGasPrice(ctx context.Context, gasPrice uint64) error
	GetGasPrice(ctx context.Context) (uint64, error)
	CountTransactionsByState(ctx context.Context, state TxState) (uint64, error)
}

type stateInterface interface {
	GetLastBatch(ctx context.Context, isVirtual bool) (*state.Batch, error)
	GetNonce(ctx context.Context, address common.Address, batchNumber uint64) (uint64, error)
	GetBalance(ctx context.Context, address common.Address, batchNumber uint64) (*big.Int, error)
}
