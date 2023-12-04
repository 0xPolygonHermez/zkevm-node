package aggregator

import (
	"context"
	"math/big"
	"time"
)

// TxProfitabilityCheckerType checks profitability of batch validation
type TxProfitabilityCheckerType string

const (
	// ProfitabilityBase checks pol collateral with min reward
	ProfitabilityBase = "base"
	// ProfitabilityAcceptAll validate batch anyway and don't check anything
	ProfitabilityAcceptAll = "acceptall"
)

// TxProfitabilityCheckerBase checks pol collateral with min reward
type TxProfitabilityCheckerBase struct {
	State                             stateInterface
	IntervalAfterWhichBatchSentAnyway time.Duration
	MinReward                         *big.Int
}

// NewTxProfitabilityCheckerBase init base tx profitability checker
func NewTxProfitabilityCheckerBase(state stateInterface, interval time.Duration, minReward *big.Int) *TxProfitabilityCheckerBase {
	return &TxProfitabilityCheckerBase{
		State:                             state,
		IntervalAfterWhichBatchSentAnyway: interval,
		MinReward:                         minReward,
	}
}

// IsProfitable checks pol collateral with min reward
func (pc *TxProfitabilityCheckerBase) IsProfitable(ctx context.Context, polCollateral *big.Int) (bool, error) {
	//if pc.IntervalAfterWhichBatchSentAnyway != 0 {
	//	ok, err := isConsolidatedBatchAppeared(ctx, pc.State, pc.IntervalAfterWhichBatchSentAnyway)
	//	if err != nil {
	//		return false, err
	//	}
	//	if ok {
	//		return true, nil
	//	}
	//}

	return polCollateral.Cmp(pc.MinReward) >= 0, nil
}

// TxProfitabilityCheckerAcceptAll validate batch anyway and don't check anything
type TxProfitabilityCheckerAcceptAll struct {
	State                             stateInterface
	IntervalAfterWhichBatchSentAnyway time.Duration
}

// NewTxProfitabilityCheckerAcceptAll init tx profitability checker that accept all txs
func NewTxProfitabilityCheckerAcceptAll(state stateInterface, interval time.Duration) *TxProfitabilityCheckerAcceptAll {
	return &TxProfitabilityCheckerAcceptAll{
		State:                             state,
		IntervalAfterWhichBatchSentAnyway: interval,
	}
}

// IsProfitable validate batch anyway and don't check anything
func (pc *TxProfitabilityCheckerAcceptAll) IsProfitable(ctx context.Context, polCollateral *big.Int) (bool, error) {
	//if pc.IntervalAfterWhichBatchSentAnyway != 0 {
	//	ok, err := isConsolidatedBatchAppeared(ctx, pc.State, pc.IntervalAfterWhichBatchSentAnyway)
	//	if err != nil {
	//		return false, err
	//	}
	//	if ok {
	//		return true, nil
	//	}
	//}

	return true, nil
}

// TODO: now it's impossible to check, when batch got consolidated, bcs it's not saved
//func isConsolidatedBatchAppeared(ctx context.Context, state stateInterface, intervalAfterWhichBatchConsolidatedAnyway time.Duration) (bool, error) {
//	batch, err := state.GetLastVerifiedBatch(ctx, nil)
//	if err != nil {
//		return false, fmt.Errorf("failed to get last verified batch, err: %v", err)
//	}
//	interval := intervalAfterWhichBatchConsolidatedAnyway * time.Minute
//	if batch..Before(time.Now().Add(-interval)) {
//		return true, nil
//	}
//
//	return false, err
//}
