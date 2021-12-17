package sequencer

import (
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

type TxSelector interface {
	SelectTxs(batchProcessor state.BatchProcessor, pendingTxs []pool.Transaction, selectionTime time.Duration) ([]*types.Transaction, []pool.Transaction, error)
	IsProfitable([]*types.Transaction) bool
}

type TxSelectorAcceptAll struct {
	Strategy Strategy
}

func NewTxSelectorAcceptAll(strategy Strategy) TxSelector {
	return &TxSelectorAcceptAll{Strategy: strategy}
}

func (s *TxSelectorAcceptAll) SelectTxs(batchProcessor state.BatchProcessor, pendingTxs []pool.Transaction, selectionTime time.Duration) ([]*types.Transaction, []pool.Transaction, error) {
	selectedTxs := make([]*types.Transaction, 0, len(pendingTxs))
	for _, tx := range pendingTxs {
		selectedTxs = append(selectedTxs, &tx.Transaction)
	}
	return selectedTxs, nil, nil
}

func (s *TxSelectorAcceptAll) IsProfitable([]*types.Transaction) bool {
	return true
}

type TxSelectorBase struct {
	Strategy               Strategy
	TxSorter               TxSorter
	TxProfitabilityChecker TxProfitabilityChecker
}

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
	case BaseProfitability:
		profitabilityChecker = &BaseTxProfitabilityChecker{MinReward: new(big.Int).SetUint64(strategy.MinReward)}
	}
	return &TxSelectorBase{
		Strategy:               strategy,
		TxSorter:               sorter,
		TxProfitabilityChecker: profitabilityChecker,
	}
}

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
			selectedTxs = append(selectedTxs, &tx.Transaction)
		}
	}

	elapsed := time.Since(start)
	if elapsed.Milliseconds()+t.Strategy.PossibleTimeToSendTx.Milliseconds() > selectionTime.Milliseconds() {
		return nil, nil, fmt.Errorf("selection took too much time, expected %d, possible time to send %d, actual %d", selectionTime, t.Strategy.PossibleTimeToSendTx, elapsed)
	}

	return selectedTxs, invalidTxs, nil
}

func (t *TxSelectorBase) IsProfitable(transactions []*types.Transaction) bool {
	return t.TxProfitabilityChecker.IsProfitable(transactions)
}

type TxSorter interface {
	SortTxs(txs []pool.Transaction) []pool.Transaction
}

type TxSorterByCostAndNonce struct{}

func (s *TxSorterByCostAndNonce) SortTxs(txs []pool.Transaction) []pool.Transaction {
	sort.Slice(txs, func(i, j int) bool {
		costI := txs[i].Cost()
		costJ := txs[j].Cost()
		if costI != costJ {
			return costI.Cmp(costJ) >= 1
		}
		return txs[i].Nonce() < txs[j].Nonce()
	})
	return txs
}

type TxSorterByCostAndTime struct{}

func (s *TxSorterByCostAndTime) SortTxs(txs []pool.Transaction) []pool.Transaction {
	sort.Slice(txs, func(i, j int) bool {
		costI := txs[i].Cost()
		costJ := txs[j].Cost()
		if costI != costJ {
			return costI.Cmp(costJ) >= 1
		}
		return txs[i].ReceivedAt.Before(txs[j].ReceivedAt)
	})
	return txs
}

type TxProfitabilityChecker interface {
	IsProfitable([]*types.Transaction) bool
}

type BaseTxProfitabilityChecker struct {
	MinReward *big.Int
}

func (pc *BaseTxProfitabilityChecker) IsProfitable(txs []*types.Transaction) bool {
	sum := big.NewInt(0)
	for _, tx := range txs {
		sum.Add(sum, tx.Cost())
		if sum.Cmp(pc.MinReward) > 0 {
			return true
		}
	}

	return false
}
