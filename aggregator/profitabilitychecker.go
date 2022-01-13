package aggregator

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/hermeznetwork/hermez-core/state"
)

// TxProfitabilityCheckerType checks profitability of batch validation
type TxProfitabilityCheckerType string

const (
	// ProfitabilityBase checks matic collateral with min reward
	ProfitabilityBase = "base"
	// ProfitabilityAcceptAll validate batch anyway and don't check anything
	ProfitabilityAcceptAll = "acceptall"
)

// TxProfitabilityChecker interface for different profitability checking algorithms
type TxProfitabilityChecker interface {
	IsProfitable(ctx context.Context, maticCollateral *big.Int) (bool, error)
}

// TxProfitabilityCheckerBase checks matic collateral with min reward
type TxProfitabilityCheckerBase struct {
	State                             state.State
	IntervalAfterWhichBatchSentAnyway time.Duration
	MinReward                         *big.Int
}

// NewTxProfitabilityCheckerBase init base tx profitability checker
func NewTxProfitabilityCheckerBase(state state.State, interval time.Duration, minReward *big.Int) TxProfitabilityChecker {
	return &TxProfitabilityCheckerBase{
		State:                             state,
		IntervalAfterWhichBatchSentAnyway: interval,
		MinReward:                         minReward,
	}
}

// IsProfitable checks matic collateral with min reward
func (pc *TxProfitabilityCheckerBase) IsProfitable(ctx context.Context, maticCollateral *big.Int) (bool, error) {
	if pc.IntervalAfterWhichBatchSentAnyway != 0 {
		ok, err := isConsolidatedBatchAppeared(ctx, pc.State, pc.IntervalAfterWhichBatchSentAnyway)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	return maticCollateral.Cmp(pc.MinReward) >= 0, nil
}

// TxProfitabilityCheckerAcceptAll validate batch anyway and don't check anything
type TxProfitabilityCheckerAcceptAll struct {
	State                             state.State
	IntervalAfterWhichBatchSentAnyway time.Duration
}

// NewTxProfitabilityCheckerAcceptAll init tx profitability checker that accept all txs
func NewTxProfitabilityCheckerAcceptAll(state state.State, interval time.Duration) TxProfitabilityChecker {
	return &TxProfitabilityCheckerAcceptAll{
		State:                             state,
		IntervalAfterWhichBatchSentAnyway: interval,
	}
}

// IsProfitable validate batch anyway and don't check anything
func (pc *TxProfitabilityCheckerAcceptAll) IsProfitable(ctx context.Context, maticCollateral *big.Int) (bool, error) {
	if pc.IntervalAfterWhichBatchSentAnyway != 0 {
		ok, err := isConsolidatedBatchAppeared(ctx, pc.State, pc.IntervalAfterWhichBatchSentAnyway)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	return true, nil
}

func isConsolidatedBatchAppeared(ctx context.Context, state state.State, intervalAfterWhichBatchConsolidatedAnyway time.Duration) (bool, error) {
	batch, err := state.GetLastBatch(ctx, false)
	if err != nil {
		return false, fmt.Errorf("failed to get last consolidated batch, err: %v", err)
	}
	interval := intervalAfterWhichBatchConsolidatedAnyway * time.Minute
	if batch.ConsolidatedAt.Before(time.Now().Add(-interval)) {
		return true, nil
	}

	return false, err
}
