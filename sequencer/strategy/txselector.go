package strategy

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

// TxSelector interface for different types of selection
type TxSelector interface {
	// SelectTxs selecting txs and returning selected txs, hashes of the selected txs (to not build array multiple times) and hashes of invalid txs
	SelectTxs(batchProcessor state.BatchProcessor, pendingTxs []pool.Transaction, selectionTime time.Duration) ([]*types.Transaction, []string, []string, error)
}

// TxSelectorAcceptAll that accept all transactions
type TxSelectorAcceptAll struct {
	Strategy Strategy
}

// NewTxSelectorAcceptAll init function
func NewTxSelectorAcceptAll(strategy Strategy) TxSelector {
	return &TxSelectorAcceptAll{Strategy: strategy}
}

// SelectTxs selects all transactions and don't check anything
func (s *TxSelectorAcceptAll) SelectTxs(batchProcessor state.BatchProcessor, pendingTxs []pool.Transaction, selectionTime time.Duration) ([]*types.Transaction, []string, []string, error) {
	selectedTxs := make([]*types.Transaction, 0, len(pendingTxs))
	selectedTxsHashes := make([]string, 0, len(pendingTxs))
	for _, tx := range pendingTxs {
		selectedTxs = append(selectedTxs, &tx.Transaction)
		selectedTxsHashes = append(selectedTxsHashes, tx.Hash().Hex())
	}
	return selectedTxs, selectedTxsHashes, nil, nil
}

// IsProfitable always returns true
func (s *TxSelectorAcceptAll) IsProfitable([]*types.Transaction) bool {
	return true
}

// TxSelectorBase tx selector with basic selection algorithm. Accepts different tx sorting and tx profitability checking structs
type TxSelectorBase struct {
	Strategy               Strategy
	TxSorter               TxSorter
	TxProfitabilityChecker TxProfitabilityChecker
}

// NewTxSelectorBase init function
func NewTxSelectorBase(strategy Strategy) TxSelector {
	var (
		sorter               TxSorter
		profitabilityChecker TxProfitabilityChecker
	)

	switch strategy.TxSorterType {
	case ByCostAndTime:
		sorter = &TxSorterByCostAndTime{}
	case ByCostAndNonce:
		sorter = &TxSorterByCostAndNonce{}
	}

	switch strategy.TxProfitabilityCheckerType {
	case ProfitabilityBase:
		profitabilityChecker = &TxProfitabilityCheckerBase{MinReward: strategy.MinReward.Int}
	case ProfitabilityAcceptAll:
		profitabilityChecker = &TxProfitabilityCheckerAcceptAll{}
	}
	return &TxSelectorBase{
		Strategy:               strategy,
		TxSorter:               sorter,
		TxProfitabilityChecker: profitabilityChecker,
	}
}

// SelectTxs process txs and split valid txs into batches of txs. This process should be completed in less than selectionTime
func (t *TxSelectorBase) SelectTxs(batchProcessor state.BatchProcessor, pendingTxs []pool.Transaction, selectionTime time.Duration) ([]*types.Transaction, []string, []string, error) {
	start := time.Now()
	sortedTxs := t.TxSorter.SortTxs(pendingTxs)
	var (
		selectedTxs                         []*types.Transaction
		selectedTxsHashes, invalidTxsHashes []string
	)
	for _, tx := range sortedTxs {
		_, _, _, err := batchProcessor.CheckTransaction(&tx.Transaction)
		if err != nil {
			invalidTxsHashes = append(invalidTxsHashes, tx.Hash().Hex())
		} else {
			t := tx.Transaction
			selectedTxs = append(selectedTxs, &t)
			selectedTxsHashes = append(selectedTxsHashes, t.Hash().Hex())
		}
	}

	elapsed := time.Since(start)
	if elapsed.Milliseconds()+t.Strategy.PossibleTimeToSendTx.Milliseconds() > selectionTime.Milliseconds() {
		return nil, nil, nil, fmt.Errorf("selection took too much time, expected %d, possible time to send %d, actual %d", selectionTime, t.Strategy.PossibleTimeToSendTx, elapsed)
	}

	return selectedTxs, selectedTxsHashes, invalidTxsHashes, nil
}
