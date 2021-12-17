package aggregator

import (
	"math/big"
)

type TxProfitabilityCheckerType string

const (
	BaseProfitability = "base"
)

type TxProfitabilityChecker interface {
	IsProfitable(maticCollateral *big.Int) bool
}

type BaseTxProfitabilityChecker struct {
	MinReward *big.Int
}

func (pc *BaseTxProfitabilityChecker) IsProfitable(maticCollateral *big.Int) bool {
	return maticCollateral.Cmp(pc.MinReward) >= 0
}
