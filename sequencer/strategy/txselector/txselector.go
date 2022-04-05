package txselector

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/runtime"
)

// TxSelector interface for different types of selection
type TxSelector interface {
	// SelectTxs selecting txs and returning selected txs, hashes of the selected txs (to not build array multiple times) and hashes of invalid txs
	SelectTxs(ctx context.Context, batchProcessor batchProcessor, pendingTxs []pool.Transaction, sequencerAddress common.Address) ([]*types.Transaction, []common.Hash, []common.Hash, []byte, error)
}

// AcceptAll that accept all transactions
type AcceptAll struct{}

// NewTxSelectorAcceptAll init function
func NewTxSelectorAcceptAll() TxSelector {
	return &AcceptAll{}
}

// SelectTxs selects all transactions and don't check anything
func (s *AcceptAll) SelectTxs(cxt context.Context, batchProcessor batchProcessor, pendingTxs []pool.Transaction, sequencerAddress common.Address) ([]*types.Transaction, []common.Hash, []common.Hash, []byte, error) {
	selectedTxs := make([]*types.Transaction, 0, len(pendingTxs))
	selectedTxsHashes := make([]common.Hash, 0, len(pendingTxs))
	for _, tx := range pendingTxs {
		t := tx.Transaction
		selectedTxs = append(selectedTxs, &t)
		selectedTxsHashes = append(selectedTxsHashes, tx.Hash())
	}
	return selectedTxs, selectedTxsHashes, nil, nil, nil
}

// Base tx selector with basic selection algorithm. Accepts different tx sorting and tx profitability checking structs
type Base struct {
	TxSorter TxSorter
}

// NewTxSelectorBase init function
func NewTxSelectorBase(cfg Config) TxSelector {
	var sorter TxSorter

	switch cfg.TxSorterType {
	case ByCostAndTime:
		sorter = &TxSorterByCostAndTime{}
	case ByCostAndNonce:
		sorter = &TxSorterByCostAndNonce{}
	}

	return &Base{
		TxSorter: sorter,
	}
}

// SelectTxs process txs and split valid txs into batches of txs. This process should be completed in less than selectionTime
func (b *Base) SelectTxs(ctx context.Context, batchProcessor batchProcessor, pendingTxs []pool.Transaction, sequencerAddress common.Address) ([]*types.Transaction, []common.Hash, []common.Hash, []byte, error) {
	sortedTxs := b.TxSorter.SortTxs(pendingTxs)
	var (
		selectedTxs                         []*types.Transaction
		selectedTxsHashes, invalidTxsHashes []common.Hash
		result                              *runtime.ExecutionResult
	)
	for _, tx := range sortedTxs {
		t := tx.Transaction
		result = batchProcessor.ProcessTransaction(ctx, &t, sequencerAddress)
		if result.Failed() {
			err := result.Err
			if state.InvalidTxErrors[err.Error()] {
				invalidTxsHashes = append(invalidTxsHashes, tx.Hash())
			} else if errors.Is(err, state.ErrNonceIsBiggerThanAccountNonce) {
				// this means, that this tx could be valid in the future, but can't be selected at this moment
				continue
			} else if errors.Is(err, state.ErrInvalidCumulativeGas) {
				// this means, that cumulative gas from txs is exceeded max amount
				return selectedTxs, selectedTxsHashes, invalidTxsHashes, result.StateRoot, nil
			} else {
				return nil, nil, nil, nil, err
			}
		} else {
			selectedTxs = append(selectedTxs, &t)
			selectedTxsHashes = append(selectedTxsHashes, t.Hash())
		}
	}

	return selectedTxs, selectedTxsHashes, invalidTxsHashes, result.StateRoot, nil
}
