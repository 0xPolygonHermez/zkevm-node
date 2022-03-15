package sequencer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

// Consumer interfaces required by the package.

// txPool contains the methods required to interact with the tx pool.
type txPool interface {
	GetPendingTxs(ctx context.Context) ([]pool.Transaction, error)
	UpdateTxState(ctx context.Context, hash common.Hash, newState pool.TxState) error
	UpdateTxsState(ctx context.Context, hashes []common.Hash, newState pool.TxState) error
	SetGasPrice(ctx context.Context, gasPrice uint64) error
}

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	SendBatch(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*types.Transaction, error)
	GetAddress() common.Address
	EstimateSendBatchCost(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*big.Int, error)
	GetCustomChainID() (*big.Int, error)
	GetCurrentSequencerCollateral() (*big.Int, error)
}

// txProfitabilityChecker interface for different profitability checkers.
type txProfitabilityChecker interface {
	IsProfitable(context.Context, []*types.Transaction) (bool, *big.Int, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastBatch(ctx context.Context, isVirtual bool) (*state.Batch, error)
	GetSequencer(ctx context.Context, address common.Address) (*state.Sequencer, error)
	GetLastBatchNumber(ctx context.Context) (uint64, error)
	GetLastBatchNumberSeenOnEthereum(ctx context.Context) (uint64, error)
	NewBatchProcessor(ctx context.Context, sequencerAddress common.Address, lastBatchNumber uint64) (*state.BasicBatchProcessor, error)
}
