package aggregator

import (
	"math/big"
)

type TxProfitabilityCheckerType string

const (
	// ProfitabilityBase checks matic collateral with min reward
	ProfitabilityBase = "base"
	// ProfitabilityAcceptAll validate batch anyway and don't check anything
	ProfitabilityAcceptAll = "acceptall"
)

// TxProfitabilityChecker interface for different profitability checking algorithms
type TxProfitabilityChecker interface {
	IsProfitable(maticCollateral *big.Int) bool
}

// TxProfitabilityCheckerBase checks matic collateral with min reward
type TxProfitabilityCheckerBase struct {
	MinReward *big.Int
}

// IsProfitable checks matic collateral with min reward
func (pc *TxProfitabilityCheckerBase) IsProfitable(maticCollateral *big.Int) bool {
	return maticCollateral.Cmp(pc.MinReward) >= 0
}

// TxProfitabilityCheckerAcceptAll validate batch anyway and don't check anything
type TxProfitabilityCheckerAcceptAll struct{}

// TxProfitabilityCheckerAcceptAll validate batch anyway and don't check anything
func (pc *TxProfitabilityCheckerAcceptAll) IsProfitable(maticCollateral *big.Int) bool {
	return true
}
