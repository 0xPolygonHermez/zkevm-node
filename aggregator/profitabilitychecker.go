package aggregator

import (
	"math/big"
)

type TxProfitabilityCheckerType string

const (
	ProfitabilityBase      = "base"
	ProfitabilityAcceptAll = "acceptall"
)

type TxProfitabilityChecker interface {
	IsProfitable(maticCollateral *big.Int) bool
}

type TxProfitabilityCheckerBase struct {
	MinReward *big.Int
}

func (pc *TxProfitabilityCheckerBase) IsProfitable(maticCollateral *big.Int) bool {
	return maticCollateral.Cmp(pc.MinReward) >= 0
}

type TxProfitablityCheckerAcceptAll struct{}

func (pc *TxProfitablityCheckerAcceptAll) IsProfitable(maticCollateral *big.Int) bool {
	return true
}
