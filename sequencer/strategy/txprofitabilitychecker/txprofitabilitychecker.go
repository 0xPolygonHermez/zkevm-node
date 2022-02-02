package txprofitabilitychecker

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/state"
)

// TxProfitabilityChecker interface for different profitability checkers
type TxProfitabilityChecker interface {
	IsProfitable(context.Context, []*types.Transaction) (bool, *big.Int, error)
}

// Base struct
type Base struct {
	EthMan etherman.EtherMan
	State  state.State

	IntervalAfterWhichBatchSentAnyway time.Duration
	MinReward                         *big.Int
	RewardPercentageToAggregator      int64
}

// NewTxProfitabilityCheckerBase inits base tx profitability checker with min reward from config and ethMan
func NewTxProfitabilityCheckerBase(
	ethMan etherman.EtherMan,
	state state.State, minReward *big.Int,
	intervalAfterWhichBatchSentAnyway time.Duration,
	rewardPercentageToAggregator int64,
) TxProfitabilityChecker {
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

	return true, aggregatorReward, nil
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
func (pc *AcceptAll) IsProfitable(ctx context.Context, txs []*types.Transaction) (bool, *big.Int, error) {
	// TODO until gas calculation and price updater is not implemented, this value will be hardcoded
	maticReward := big.NewInt(int64(len(txs)))
	maticReward.Mul(maticReward, big.NewInt(encoding.TenToThePowerOf18))
	if pc.IntervalAfterWhichBatchSentAnyway != 0 {
		ok, err := isNewBatchNotAppeared(ctx, pc.State, pc.IntervalAfterWhichBatchSentAnyway)
		if err != nil {
			return false, maticReward, err
		}
		if ok {
			return true, maticReward, nil
		}
	}

	return true, maticReward, nil
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
