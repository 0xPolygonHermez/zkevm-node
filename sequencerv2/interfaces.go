//nolint
package sequencerv2

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
	GetPendingTxs(ctx context.Context, isClaims bool, limit uint64) ([]pool.Transaction, error)
	UpdateTxState(ctx context.Context, hash common.Hash, newState pool.TxState) error
	UpdateTxsState(ctx context.Context, hashes []common.Hash, newState pool.TxState) error
	SetGasPrice(ctx context.Context, gasPrice uint64) error
	IsTxPending(ctx context.Context, hash common.Hash) (bool, error)
}

// etherman contains the methods required to interact with ethereum.
type etherman interface {
	GetAddress() common.Address
	GetSequencerAddress() common.Address

	EstimateSendBatchCost(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (*big.Int, error)
	EstimateSendBatchGas(ctx context.Context, txs []*types.Transaction, maticAmount *big.Int) (uint64, error)

	GetSequencerCollateral() (*big.Int, error)

	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)

	SequenceBatches(sequences []*sequence) error
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastBatch(ctx context.Context, isVirtual bool, txBundleID string) (*state.Batch, error)
	GetLastBatchNumber(ctx context.Context, txBundleID string) (uint64, error)
	GetLastBatchNumberSeenOnEthereum(ctx context.Context, txBundleID string) (uint64, error)
	GetLastBatchByStateRoot(ctx context.Context, stateRoot []byte, txBundleID string) (*state.Batch, error)
	NewBatchProcessor(ctx context.Context, sequencerAddress common.Address, stateRoot []byte, txBundleID string) (*state.BatchProcessor, error)
}
