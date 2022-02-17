package sequencer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
)

// Consumer interfaces required by the package.

// txPool contains the methods required to interact with the tx pool.
type txPool interface {
	GetPendingTxs(ctx context.Context) ([]pool.Transaction, error)
	UpdateTxState(ctx context.Context, hash common.Hash, newState pool.TxState) error
	UpdateTxsState(ctx context.Context, hashes []string, newState pool.TxState) error
	SetGasPrice(ctx context.Context, gasPrice uint64) error
}

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	SendBatch(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*types.Transaction, error)
	GetAddress() common.Address
	EstimateSendBatchCost(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*big.Int, error)
	GetCustomChainID() (*big.Int, error)
}

// txProfitabilityChecker interface for different profitability checkers.
type txProfitabilityChecker interface {
	IsProfitable(context.Context, []*types.Transaction) (bool, *big.Int, error)
}
