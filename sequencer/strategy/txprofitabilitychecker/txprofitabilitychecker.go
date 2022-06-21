package txprofitabilitychecker

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/hermeznetwork/hermez-core/pricegetter"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy/txselector"
)

// Base struct
type Base struct {
	EthMan      etherman
	State       stateInterface
	PriceGetter pricegetter.Client

	IntervalAfterWhichBatchSentAnyway time.Duration
	MinReward                         *big.Int
	RewardPercentageToAggregator      int64
}

// NewTxProfitabilityCheckerBase inits base tx profitability checker with min reward from config and ethMan
func NewTxProfitabilityCheckerBase(
	ethMan etherman,
	state stateInterface,
	priceGetter priceGetter,
	minReward *big.Int,
	intervalAfterWhichBatchSentAnyway time.Duration,
	rewardPercentageToAggregator int64,
) *Base {
	return &Base{
		EthMan:      ethMan,
		State:       state,
		PriceGetter: priceGetter,

		MinReward:                         minReward,
		IntervalAfterWhichBatchSentAnyway: intervalAfterWhichBatchSentAnyway,
		RewardPercentageToAggregator:      rewardPercentageToAggregator,
	}
}

// IsProfitable checks for txs cost against the main reward
func (pc *Base) IsProfitable(ctx context.Context, selectionRes txselector.SelectTxsOutput) (bool, *big.Int, error) {
	collateral, err := pc.EthMan.GetCurrentSequencerCollateral()
	if err != nil {
		return false, big.NewInt(0), fmt.Errorf("failed to get current collateral amount from smc, err: %v", err)
	}

	// if there is no batches in time, then sequencer have to forge anyway
	if pc.IntervalAfterWhichBatchSentAnyway != 0 {
		ok, err := isNewBatchNotAppeared(ctx, pc.State, pc.IntervalAfterWhichBatchSentAnyway)
		if err != nil {
			return false, big.NewInt(0), err
		}
		if ok {
			return true, collateral, nil
		}
	}

	var sentAnyway bool
	if len(selectionRes.SelectedClaimsTxs) > 0 {
		sentAnyway = true
	}

	txs := append(selectionRes.SelectedTxs, selectionRes.SelectedClaimsTxs...)

	// sending amount of txs as matic reward there, bcs to calculate gas cost for tx this value is not important
	// checker can't get it before, bcs final matic amount is dependent on value from gas estimation
	// TODO improve the matic amount
	const maticAmount uint64 = 10000000000000000001
	gasCostForSendingBatch, err := pc.EthMan.EstimateSendBatchCost(ctx, txs, new(big.Int).SetUint64(maticAmount))
	if err != nil {
		return false, big.NewInt(0), fmt.Errorf("failed to estimate gas for sending batch, err: %v", err)
	}

	// calculate gas cost of all txs in a batch
	txsGasCost := big.NewInt(0)
	for _, tx := range selectionRes.SelectedTxs {
		txsGasCost.Add(txsGasCost, tx.Cost())
	}

	reward := new(big.Int).Sub(txsGasCost, gasCostForSendingBatch)

	// if gasCostForSendingBatch is more than txsGasCost, then selection is not profitable
	if reward.Cmp(big.NewInt(0)) < 0 && !sentAnyway {
		return false, big.NewInt(0), nil
	} else if sentAnyway {
		return true, collateral, nil
	}

	// calculate aggregator reward in ether wei
	aggregatorReward := big.NewInt(0).Mul(reward, big.NewInt(pc.RewardPercentageToAggregator))
	// bcs on previous step reward was multiplied by not adapted percentage amount (e.g. 80), it should be divided by 100
	const percentageDivider = 100
	aggregatorReward.Div(aggregatorReward, big.NewInt(percentageDivider))
	// get sequencer reward
	sequencerReward := big.NewInt(0).Sub(reward, aggregatorReward)
	// if sequencer reward is less than minimal reward from config, then it makes no sense to send a batch
	if sequencerReward.Cmp(pc.MinReward) < 0 && !sentAnyway {
		return false, big.NewInt(0), nil
	} else if sentAnyway {
		return true, collateral, nil
	}
	// get price from the price updater
	price, err := pc.PriceGetter.GetEthToMaticPrice(ctx)
	if err != nil {
		return false, big.NewInt(0), err
	}
	priceInt := new(big.Int)
	price.Int(priceInt)
	// calculate aggregator reward in matic
	aggregatorReward.Mul(aggregatorReward, priceInt)
	// if aggregator reward is less than the collateral retrieved from the smc, then it makes no sense to propose a new batch
	if aggregatorReward.Cmp(collateral) < 0 && !sentAnyway {
		return false, big.NewInt(0), nil
	} else if sentAnyway {
		return true, collateral, nil
	}

	return true, aggregatorReward, nil
}

// AcceptAll always returns true
type AcceptAll struct {
	EthMan                            etherman
	State                             stateInterface
	IntervalAfterWhichBatchSentAnyway time.Duration
}

// NewTxProfitabilityCheckerAcceptAll inits tx profitability checker which accept all
func NewTxProfitabilityCheckerAcceptAll(ethman etherman, state stateInterface, intervalAfterWhichBatchSentAnyway time.Duration) *AcceptAll {
	return &AcceptAll{
		EthMan:                            ethman,
		State:                             state,
		IntervalAfterWhichBatchSentAnyway: intervalAfterWhichBatchSentAnyway,
	}
}

// IsProfitable always returns true, until it's failed to get sequencer collateral
func (pc *AcceptAll) IsProfitable(ctx context.Context, selectionRes txselector.SelectTxsOutput) (bool, *big.Int, error) {
	collateral, err := pc.EthMan.GetCurrentSequencerCollateral()
	if err != nil {
		return false, big.NewInt(0), fmt.Errorf("failed to get current collateral amount from smc, err: %v", err)
	}
	return true, collateral, nil
}

func isNewBatchNotAppeared(ctx context.Context, state stateInterface, intervalAfterWhichBatchSentAnyway time.Duration) (bool, error) {
	batch, err := state.GetLastBatch(ctx, true, "")
	if err != nil {
		return false, fmt.Errorf("failed to get last batch, err: %v", err)
	}
	interval := intervalAfterWhichBatchSentAnyway * time.Minute
	if batch.ReceivedAt.Before(time.Now().Add(-interval)) {
		return true, nil
	}

	return false, nil
}
