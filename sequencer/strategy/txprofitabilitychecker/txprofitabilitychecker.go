package txprofitabilitychecker

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/state"
)

// TxProfitabilityChecker interface for different profitability checkers
type TxProfitabilityChecker interface {
	IsProfitable(context.Context, []*types.Transaction) (bool, error)
}

// Base struct
type Base struct {
	EthMan etherman.EtherMan
	State  state.State

	IntervalAfterWhichBatchSentAnyway time.Duration
	MinReward                         *big.Int
}

// NewTxProfitabilityCheckerBase inits base tx profitability checker with min reward from config and ethMan
func NewTxProfitabilityCheckerBase(ethMan etherman.EtherMan, state state.State, minReward *big.Int, intervalAfterWhichBatchSentAnyway time.Duration) TxProfitabilityChecker {
	return &Base{
		EthMan: ethMan,
		State:  state,

		MinReward:                         minReward,
		IntervalAfterWhichBatchSentAnyway: intervalAfterWhichBatchSentAnyway,
	}
}

// IsProfitable checks for txs cost against the main reward
func (pc *Base) IsProfitable(ctx context.Context, txs []*types.Transaction) (bool, error) {
	if pc.IntervalAfterWhichBatchSentAnyway != 0 {
		ok, err := isNewBatchNotAppeared(ctx, pc.State, pc.IntervalAfterWhichBatchSentAnyway)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	txsGasCost := big.NewInt(0)
	for _, tx := range txs {
		txsGasCost.Add(txsGasCost, tx.Cost())
	}
	// sending amount of txs as matic reward there, bcs to calculate gas cost for tx this value is not important
	// checker can't get it before, bcs final matic amount is dependent on value from gas estimation
	//TODO improve the matic amount
	const maticAmount = 1000000000000000001
	gasCostForSendingBatch, err := pc.EthMan.EstimateSendBatchCost(ctx, txs, big.NewInt(int64(maticAmount)))
	if err != nil {
		return false, fmt.Errorf("failed to estimate gas for sending batch, err: %v", err)
	}

	reward := new(big.Int).Sub(txsGasCost, gasCostForSendingBatch)
	if reward.Cmp(big.NewInt(0)) < 0 {
		return false, nil
	}
	if reward.Cmp(pc.MinReward) < 0 {
		return false, nil
	}

	// TODO: take from config, how much matic tokens in percentage from whole reward sequencer is willing to pay to aggregator
	// to get his batch consolidated, then calculate, how much is reward minus minReward in matic tokens and multiply it
	// by percentage in config. This value we can also return from function

	return false, nil
}

// AcceptAll always returns true
type AcceptAll struct {
	State                             state.State
	IntervalAfterWhichBatchSentAnyway time.Duration
}

// NewTxProfitabilityCheckerAcceptAll inits tx profitability checker which accept all
func NewTxProfitabilityCheckerAcceptAll(state state.State, intervalAfterWhichBatchSentAnyway time.Duration) TxProfitabilityChecker {
	return &AcceptAll{
		State:                             state,
		IntervalAfterWhichBatchSentAnyway: intervalAfterWhichBatchSentAnyway,
	}
}

// IsProfitable always returns true
func (pc *AcceptAll) IsProfitable(ctx context.Context, txs []*types.Transaction) (bool, error) {
	if pc.IntervalAfterWhichBatchSentAnyway != 0 {
		ok, err := isNewBatchNotAppeared(ctx, pc.State, pc.IntervalAfterWhichBatchSentAnyway)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	return true, nil
}

func isNewBatchNotAppeared(ctx context.Context, state state.State, intervalAfterWhichBatchSentAnyway time.Duration) (bool, error) {
	batch, err := state.GetLastBatch(ctx, true)
	if err != nil {
		return false, fmt.Errorf("failed to get last batch, err: %v", err)
	}
	interval := intervalAfterWhichBatchSentAnyway * time.Minute

	if batch.ReceivedAt.Before(time.Now().Add(-interval)) {
		return true, nil
	}

	return false, nil
}
