package sequencer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy/txselector"
	"github.com/hermeznetwork/hermez-core/state"
)

// Consumer interfaces required by the package.

// txPool contains the methods required to interact with the tx pool.
type txPool interface {
	GetPendingTxs(ctx context.Context, isClaims bool, limit uint64) ([]pool.Transaction, error)
	UpdateTxState(ctx context.Context, hash common.Hash, newState pool.TxState) error
	UpdateTxsState(ctx context.Context, hashes []common.Hash, newState pool.TxState) error
	SetGasPrice(ctx context.Context, gasPrice uint64) error
}

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	SendBatch(ctx context.Context, gasLimit uint64, txs []*types.Transaction, maticAmount *big.Int) (*types.Transaction, error)
	GetAddress() common.Address
	EstimateSendBatchCost(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*big.Int, error)
	GetCustomChainID() (*big.Int, error)
	GetCurrentSequencerCollateral() (*big.Int, error)
	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

// txProfitabilityChecker interface for different profitability checkers.
type txProfitabilityChecker interface {
	IsProfitable(context.Context, txselector.SelectTxsOutput) (bool, *big.Int, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastBatch(ctx context.Context, isVirtual bool, txundleID string) (*state.Batch, error)
	GetSequencer(ctx context.Context, address common.Address, txundleID string) (*state.Sequencer, error)
	GetLastBatchNumber(ctx context.Context, txundleID string) (uint64, error)
	GetLastBatchNumberSeenOnEthereum(ctx context.Context, txundleID string) (uint64, error)
	GetLastBatchByStateRoot(ctx context.Context, stateRoot []byte, txundleID string) (*state.Batch, error)
	NewBatchProcessor(ctx context.Context, sequencerAddress common.Address, stateRoot []byte, txundleID string) (*state.BatchProcessor, error)
}
