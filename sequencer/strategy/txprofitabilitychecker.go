package strategy

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// TxProfitabilityCheckerType for different profitability checkers types
type TxProfitabilityCheckerType string

const (
	// ProfitabilityBase type that checks sum of costs of txs against min reward
	ProfitabilityBase = "base"
	// ProfitabilityAcceptAll validate batch anyway and don't check anything
	ProfitabilityAcceptAll = "acceptall"
)

// TxProfitabilityChecker interface for different profitability checkers
type TxProfitabilityChecker interface {
	IsProfitable([]*types.Transaction) bool
}

// TxProfitabilityCheckerBase struct
type TxProfitabilityCheckerBase struct {
	MinReward *big.Int
}

// IsProfitable checks for txs cost against the main reward
func (pc *TxProfitabilityCheckerBase) IsProfitable(txs []*types.Transaction) bool {
	sum := big.NewInt(0)
	for _, tx := range txs {
		sum.Add(sum, tx.Cost())
		if sum.Cmp(pc.MinReward) > 0 {
			return true
		}
	}

	return false
}

// TxProfitabilityCheckerAcceptAll always returns true
type TxProfitabilityCheckerAcceptAll struct{}

// IsProfitable always returns true
func (pc *TxProfitabilityCheckerAcceptAll) IsProfitable(txs []*types.Transaction) bool {
	return true
}
