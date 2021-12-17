package strategy

import (
	"sort"

	"github.com/hermeznetwork/hermez-core/pool"
)

// TxSorterType for different txs sorters types
type TxSorterType string

const (
	// ByCostAndTime sorting txs by cost and time
	ByCostAndTime TxSorterType = "bycostandtime"
	// ByCostAndNonce sorting txs by cost and nonce
	ByCostAndNonce = "bycostandnonce"
)

// TxSorter interface for for different txs sorters
type TxSorter interface {
	SortTxs(txs []pool.Transaction) []pool.Transaction
}

// TxSorterByCostAndNonce sorts by tx cost and nonce
type TxSorterByCostAndNonce struct{}

// SortTxs sorts by tx cost and nonce
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

// TxSorterByCostAndTime sorts by tx cost and time
type TxSorterByCostAndTime struct{}

// SortTxs sorts by tx cost and time
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
