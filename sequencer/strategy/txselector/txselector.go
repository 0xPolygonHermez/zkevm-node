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
	SelectTxs(ctx context.Context, input SelectTxsInput) (SelectTxsOutput, error)
}

// AcceptAll that accept all transactions
type AcceptAll struct{}

// NewTxSelectorAcceptAll init function
func NewTxSelectorAcceptAll() TxSelector {
	return &AcceptAll{}
}

// SelectTxs selects all transactions and don't check anything
func (s *AcceptAll) SelectTxs(cxt context.Context, input SelectTxsInput) (SelectTxsOutput, error) {
	pendingClaimsTxs := input.PendingClaimsTxs
	selectedClaimsTxs := make([]*types.Transaction, 0, len(pendingClaimsTxs))
	selectedClaimsTxsHashes := make([]common.Hash, 0, len(pendingClaimsTxs))
	for _, tx := range input.PendingClaimsTxs {
		t := tx.Transaction
		selectedClaimsTxs = append(selectedClaimsTxs, &t)
		selectedClaimsTxsHashes = append(selectedClaimsTxsHashes, tx.Hash())
	}

	pendingTxs := input.PendingTxs
	selectedTxs := make([]*types.Transaction, 0, len(pendingTxs))
	selectedTxsHashes := make([]common.Hash, 0, len(pendingTxs))

	for _, tx := range pendingTxs {
		t := tx.Transaction
		selectedTxs = append(selectedTxs, &t)
		selectedTxsHashes = append(selectedTxsHashes, tx.Hash())
	}
	return SelectTxsOutput{
		SelectedTxs:             selectedTxs,
		SelectedClaimsTxs:       selectedClaimsTxs,
		SelectedTxsHashes:       selectedTxsHashes,
		SelectedClaimsTxsHashes: selectedClaimsTxsHashes,
	}, nil
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

type SelectTxsInput struct {
	BatchProcessor               batchProcessor
	PendingTxs, PendingClaimsTxs []pool.Transaction
	SequencerAddress             common.Address
}

type SelectTxsOutput struct {
	SelectedTxs, SelectedClaimsTxs                               []*types.Transaction
	SelectedTxsHashes, SelectedClaimsTxsHashes, InvalidTxsHashes []common.Hash
	NewRoot                                                      []byte
	BatchNumber                                                  uint64
}

type selectTxsInternalOutput struct {
	selectedTxs                         []*types.Transaction
	selectedTxsHashes, invalidTxsHashes []common.Hash
	newRoot                             []byte
	isGasExceeded                       bool
}

func (b *Base) selectTxs(ctx context.Context, batchProcessor batchProcessor, pendingTxs []pool.Transaction, sequencerAddress common.Address) (selectTxsInternalOutput, error) {
	sortedTxs := b.TxSorter.SortTxs(pendingTxs)
	var (
		selectedTxs                         []*types.Transaction
		selectedTxsHashes, invalidTxsHashes []common.Hash
		result                              *runtime.ExecutionResult
		root                                []byte
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
				return selectTxsInternalOutput{
					selectedTxs:       selectedTxs,
					selectedTxsHashes: selectedTxsHashes,
					invalidTxsHashes:  invalidTxsHashes,
					newRoot:           result.StateRoot,
					isGasExceeded:     true,
				}, nil
			} else {
				return selectTxsInternalOutput{}, err
			}
		} else {
			selectedTxs = append(selectedTxs, &t)
			selectedTxsHashes = append(selectedTxsHashes, t.Hash())
		}
	}

	if result != nil {
		root = result.StateRoot
	}

	return selectTxsInternalOutput{
		selectedTxs:       selectedTxs,
		selectedTxsHashes: selectedTxsHashes,
		invalidTxsHashes:  invalidTxsHashes,
		newRoot:           root,
	}, nil
}

// SelectTxs process txs and split valid txs into batches of txs. This process should be completed in less than selectionTime
func (b *Base) SelectTxs(ctx context.Context, input SelectTxsInput) (SelectTxsOutput, error) {
	selectTxsClaimsInternalOutput, err := b.selectTxs(ctx, input.BatchProcessor, input.PendingClaimsTxs, input.SequencerAddress)
	if err != nil {
		return SelectTxsOutput{}, err
	}

	if selectTxsClaimsInternalOutput.isGasExceeded {
		return SelectTxsOutput{
			SelectedClaimsTxs:       selectTxsClaimsInternalOutput.selectedTxs,
			SelectedClaimsTxsHashes: selectTxsClaimsInternalOutput.selectedTxsHashes,
			InvalidTxsHashes:        selectTxsClaimsInternalOutput.invalidTxsHashes,
			NewRoot:                 selectTxsClaimsInternalOutput.newRoot,
		}, nil
	}

	selectTxsInternalOutput, err := b.selectTxs(ctx, input.BatchProcessor, input.PendingTxs, input.SequencerAddress)
	if err != nil {
		return SelectTxsOutput{}, err
	}

	return SelectTxsOutput{
		SelectedTxs:             selectTxsInternalOutput.selectedTxs,
		SelectedClaimsTxs:       selectTxsClaimsInternalOutput.selectedTxs,
		SelectedTxsHashes:       selectTxsInternalOutput.selectedTxsHashes,
		SelectedClaimsTxsHashes: selectTxsClaimsInternalOutput.selectedTxsHashes,
		InvalidTxsHashes:        append(selectTxsInternalOutput.invalidTxsHashes, selectTxsClaimsInternalOutput.invalidTxsHashes...),
		NewRoot:                 selectTxsInternalOutput.newRoot,
	}, nil
}
