package txprofitabilitychecker

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
)

// Base struct
type Base struct {
	EthMan etherman
	State  state.State

	IntervalAfterWhichBatchSentAnyway time.Duration
	MinReward                         *big.Int
	RewardPercentageToAggregator      int64
}

// NewTxProfitabilityCheckerBase inits base tx profitability checker with min reward from config and ethMan
func NewTxProfitabilityCheckerBase(
	ethMan etherman,
	state state.State, minReward *big.Int,
	intervalAfterWhichBatchSentAnyway time.Duration,
	rewardPercentageToAggregator int64,
) *Base {
	return &Base{
		EthMan: ethMan,
		State:  state,

		MinReward:                         minReward,
		IntervalAfterWhichBatchSentAnyway: intervalAfterWhichBatchSentAnyway,
		RewardPercentageToAggregator:      rewardPercentageToAggregator,
	}
}

// IsProfitable checks for txs cost against the main reward
func (pc *Base) IsProfitable(ctx context.Context, txs []*types.Transaction) (bool, *big.Int, error) {
	// sending amount of txs as matic reward there, bcs to calculate gas cost for tx this value is not important
	// checker can't get it before, bcs final matic amount is dependent on value from gas estimation
	// TODO improve the matic amount
	const maticAmount = 1000000000000000001
	gasCostForSendingBatch, err := pc.EthMan.EstimateSendBatchCost(ctx, txs, big.NewInt(int64(maticAmount)))
	if err != nil {
		return false, big.NewInt(0), fmt.Errorf("failed to estimate gas for sending batch, err: %v", err)
	}

	// calculate gas cost of all txs in a batch
	txsGasCost := big.NewInt(0)
	for _, tx := range txs {
		txsGasCost.Add(txsGasCost, tx.Cost())
	}

	reward := new(big.Int).Sub(txsGasCost, gasCostForSendingBatch)

	// if there is no batches in time, then sequencer have to forge anyway
	if pc.IntervalAfterWhichBatchSentAnyway != 0 {
		ok, err := isNewBatchNotAppeared(ctx, pc.State, pc.IntervalAfterWhichBatchSentAnyway)
		if err != nil {
			return false, big.NewInt(0), err
		}
		// TODO: how much should sequencer send for collateral, if there is no profit for sequencer?
		if ok {
			return true, big.NewInt(0), nil
		}
	}

	// if gasCostForSendingBatch is more than txsGasCost, then selection is not profitable
	if reward.Cmp(big.NewInt(0)) < 0 {
		return false, big.NewInt(0), nil
	}

	// calculate aggregator reward in ether wei
	aggregatorReward := big.NewInt(0).Mul(reward, big.NewInt(pc.RewardPercentageToAggregator))
	// bcs on previous step reward was multiplied by not adapted percentage amount (e.g. 80), it should be divided by 100
	const percentageDivider = 100
	aggregatorReward.Div(aggregatorReward, big.NewInt(percentageDivider))
	// get sequencer reward
	sequencerReward := big.NewInt(0).Sub(reward, aggregatorReward)
	// if sequencer reward is less than minimal reward from config, then it makes no sense to send a batch
	if sequencerReward.Cmp(pc.MinReward) < 0 {
		return false, big.NewInt(0), nil
	}
	// bcs price updater is not supported yet, ethToMatic ratio is hardcoded there
	const ethToMatic = 2000
	// calculate aggregator reward in matic
	aggregatorReward.Mul(aggregatorReward, big.NewInt(ethToMatic))
	// if aggregator reward is less than the collateral retrieved from the smc, then it makes no sense to propose a new batch
	collateral, err := pc.EthMan.GetCurrentSequencerCollateral()
	if err != nil {
		return false, big.NewInt(0), fmt.Errorf("failed to get current collateral amount from smc, err: %v", err)
	}
	if aggregatorReward.Cmp(collateral) < 0 {
		return false, big.NewInt(0), nil
	}
	return true, aggregatorReward, nil
}

// AcceptAll always returns true
type AcceptAll struct {
	EthMan                            etherman
	State                             state.State
	IntervalAfterWhichBatchSentAnyway time.Duration
}

// NewTxProfitabilityCheckerAcceptAll inits tx profitability checker which accept all
func NewTxProfitabilityCheckerAcceptAll(ethman etherman, state state.State, intervalAfterWhichBatchSentAnyway time.Duration) *AcceptAll {
	return &AcceptAll{
		EthMan:                            ethman,
		State:                             state,
		IntervalAfterWhichBatchSentAnyway: intervalAfterWhichBatchSentAnyway,
	}
}

// IsProfitable always returns true, until it's failed to get sequencer collateral
func (pc *AcceptAll) IsProfitable(ctx context.Context, txs []*types.Transaction) (bool, *big.Int, error) {
	collateral, err := pc.EthMan.GetCurrentSequencerCollateral()
	if err != nil {
		return false, big.NewInt(0), fmt.Errorf("failed to get current collateral amount from smc, err: %v", err)
	}
	return true, collateral, nil
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
