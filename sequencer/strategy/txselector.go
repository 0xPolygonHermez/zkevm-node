package strategy

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

// TxSelector interface for different types of selection
type TxSelector interface {
	SelectTxs(batchProcessor state.BatchProcessor, pendingTxs []pool.Transaction, selectionTime time.Duration) ([]*types.Transaction, []pool.Transaction, error)
	IsProfitable([]*types.Transaction) bool
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
func (s *TxSelectorAcceptAll) SelectTxs(batchProcessor state.BatchProcessor, pendingTxs []pool.Transaction, selectionTime time.Duration) ([]*types.Transaction, []pool.Transaction, error) {
	selectedTxs := make([]*types.Transaction, 0, len(pendingTxs))
	for _, tx := range pendingTxs {
		selectedTxs = append(selectedTxs, &tx.Transaction)
	}
	return selectedTxs, nil, nil
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
		profitabilityChecker = &TxProfitabilityCheckerBase{MinReward: new(big.Int).SetUint64(strategy.MinReward)}
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
func (t *TxSelectorBase) SelectTxs(batchProcessor state.BatchProcessor, pendingTxs []pool.Transaction, selectionTime time.Duration) ([]*types.Transaction, []pool.Transaction, error) {
	start := time.Now()
	sortedTxs := t.TxSorter.SortTxs(pendingTxs)
	var (
		selectedTxs []*types.Transaction
		invalidTxs  []pool.Transaction
	)
	for _, tx := range sortedTxs {
		_, _, _, err := batchProcessor.CheckTransaction(&tx.Transaction)
		if err != nil {
			invalidTxs = append(invalidTxs, tx)
		} else {
			t := tx.Transaction
			selectedTxs = append(selectedTxs, &t)
		}
	}

	elapsed := time.Since(start)
	if elapsed.Milliseconds()+t.Strategy.PossibleTimeToSendTx.Milliseconds() > selectionTime.Milliseconds() {
		return nil, nil, fmt.Errorf("selection took too much time, expected %d, possible time to send %d, actual %d", selectionTime, t.Strategy.PossibleTimeToSendTx, elapsed)
	}

	return selectedTxs, invalidTxs, nil
}

// IsProfitable checks profitability for base tx selector
func (t *TxSelectorBase) IsProfitable(transactions []*types.Transaction) bool {
	return t.TxProfitabilityChecker.IsProfitable(transactions)
}
