package txselector

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

// TxSelector interface for different types of selection
type TxSelector interface {
	// SelectTxs selecting txs and returning selected txs, hashes of the selected txs (to not build array multiple times) and hashes of invalid txs
	SelectTxs(batchProcessor state.BatchProcessor, pendingTxs []pool.Transaction) ([]*types.Transaction, []string, []string, error)
}

// AcceptAllType that accept all transactions
type AcceptAll struct{}

// NewTxSelectorAcceptAll init function
func NewTxSelectorAcceptAll() TxSelector {
	return &AcceptAll{}
}

// SelectTxs selects all transactions and don't check anything
func (s *AcceptAll) SelectTxs(batchProcessor state.BatchProcessor, pendingTxs []pool.Transaction) ([]*types.Transaction, []string, []string, error) {
	selectedTxs := make([]*types.Transaction, 0, len(pendingTxs))
	selectedTxsHashes := make([]string, 0, len(pendingTxs))
	for _, tx := range pendingTxs {
		t := tx.Transaction
		selectedTxs = append(selectedTxs, &t)
		selectedTxsHashes = append(selectedTxsHashes, tx.Hash().Hex())
	}
	return selectedTxs, selectedTxsHashes, nil, nil
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
func (t *Base) SelectTxs(batchProcessor state.BatchProcessor, pendingTxs []pool.Transaction) ([]*types.Transaction, []string, []string, error) {
	sortedTxs := t.TxSorter.SortTxs(pendingTxs)
	var (
		selectedTxs                         []*types.Transaction
		selectedTxsHashes, invalidTxsHashes []string
	)
	for _, tx := range sortedTxs {
		t := tx.Transaction
		_, _, _, err := batchProcessor.CheckTransaction(&t)
		if err != nil {
			invalidTxsHashes = append(invalidTxsHashes, tx.Hash().Hex())
		} else {
			selectedTxs = append(selectedTxs, &t)
			selectedTxsHashes = append(selectedTxsHashes, t.Hash().Hex())
		}
	}

	return selectedTxs, selectedTxsHashes, invalidTxsHashes, nil
}
