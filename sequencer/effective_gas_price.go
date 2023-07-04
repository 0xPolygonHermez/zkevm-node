package sequencer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// CalculateTxBreakEvenGasPrice calculates the break even gas price for a transaction
func (f *finalizer) CalculateTxBreakEvenGasPrice(tx *TxTracker) (*big.Int, error) {
	const (
		// constants used in calculation of BreakEvenGasPrice
		signatureBytesLength           = 65
		effectivePercentageBytesLength = 1
		totalRlpFieldsLength           = signatureBytesLength + effectivePercentageBytesLength
	)

	if tx.L1GasPrice == 0 {
		log.Warn("CalculateTxBreakEvenGasPrice: L1 gas price 0. Skipping estimation...")
		return nil, ErrZeroL1GasPrice
	}

	// Get L2 Min Gas Price
	l2MinGasPrice := uint64(float64(tx.L1GasPrice) * f.effectiveGasPriceCfg.L1GasPriceFactor)
	if l2MinGasPrice < f.dbManager.GetDefaultMinGasPriceAllowed() {
		l2MinGasPrice = f.dbManager.GetDefaultMinGasPriceAllowed()
	}

	// Calculate BreakEvenGasPrice
	totalTxPrice := (tx.BatchResources.ZKCounters.CumulativeGasUsed * l2MinGasPrice) + (totalRlpFieldsLength * tx.BatchResources.Bytes * tx.L1GasPrice)
	breakEvenGasPrice := big.NewInt(0).SetUint64(uint64(float64(totalTxPrice/tx.BatchResources.ZKCounters.CumulativeGasUsed) * f.effectiveGasPriceCfg.MarginFactor))

	return breakEvenGasPrice, nil
}

func (f *finalizer) CompareTxBreakEvenGasPrice(ctx context.Context, tx *TxTracker, newGasUsed uint64) error {
	// Increase nunber of executions related to gas price
	tx.EffectiveGasPriceProcessCount++

	newBreakEvenGasPrice, err := f.CalculateTxBreakEvenGasPrice(tx)
	if err != nil {
		log.Errorf("failed to calculate breakEvenPrice with new gasUsed: %s", err.Error())
		return err
	}

	// if newBreakEvenGasPrice < tx.BreakEvenGasPrice
	if newBreakEvenGasPrice.Cmp(tx.BreakEvenGasPrice) == -1 {
		// Compute the difference
		diff := new(big.Int).Sub(tx.BreakEvenGasPrice, newBreakEvenGasPrice)
		// Compute deviation of breakEvenPrice
		deviation := new(big.Int).Div(new(big.Int).Mul(tx.BreakEvenGasPrice, f.maxBreakEvenGasPriceDeviationPercentage), big.NewInt(100)) //nolint:gomnd

		// tx.BreakEvenGasPrice - newBreakEventGasPrice is greater than the max deviation allowed
		if diff.Cmp(deviation) == 1 {
			if tx.EffectiveGasPriceProcessCount < 2 { //nolint:gomnd
				tx.BreakEvenGasPrice = newBreakEvenGasPrice
				return ErrEffectiveGasPriceReprocess
			} else {
				tx.BreakEvenGasPrice = tx.GasPrice
				tx.IsEffectiveGasPriceFinalExecution = true
				ev := &event.Event{
					ReceivedAt:  time.Now(),
					Source:      event.Source_Node,
					Component:   event.Component_Sequencer,
					Level:       event.Level_Critical,
					EventID:     event.EventID_FinalizerBreakEvenGasPriceBigDifference,
					Description: fmt.Sprintf("The difference: %s between the breakEvenGasPrice and the newBreakEvenGasPrice is more than %d %%", diff.String(), f.effectiveGasPriceCfg.MaxBreakEvenGasPriceDeviationPercentage),
					Json: struct {
						transactionHash               string
						preExecutionBreakEvenGasPrice string
						newBreakEvenGasPrice          string
						diff                          string
						deviation                     string
					}{
						transactionHash:               tx.Hash.String(),
						preExecutionBreakEvenGasPrice: tx.BreakEvenGasPrice.String(),
						newBreakEvenGasPrice:          newBreakEvenGasPrice.String(),
						diff:                          diff.String(),
						deviation:                     deviation.String(),
					},
				}
				err = f.eventLog.LogEvent(ctx, ev)
				if err != nil {
					log.Errorf("failed to log event: %s", err.Error())
				}
				return ErrEffectiveGasPriceReprocess
			}
		} // TODO: Review this check regarding tx.GasPrice being nil
	} else if tx.GasPrice != nil && newBreakEvenGasPrice.Cmp(tx.GasPrice) == 1 {
		tx.BreakEvenGasPrice = tx.GasPrice
		tx.IsEffectiveGasPriceFinalExecution = true
		return ErrEffectiveGasPriceReprocess
	}

	return nil
}

// CalculateEffectiveGasPricePercentage calculates the gas price's effective percentage
func CalculateEffectiveGasPricePercentage(gasPrice *big.Int, breakEven *big.Int) (uint8, error) {
	const bits = 256
	var bitsBigInt = big.NewInt(bits)

	if breakEven == nil || gasPrice == nil ||
		gasPrice.Cmp(big.NewInt(0)) == 0 || breakEven.Cmp(big.NewInt(0)) == 0 {
		return 0, ErrBreakEvenGasPriceEmpty
	}

	if gasPrice.Cmp(breakEven) <= 0 {
		return state.MaxEffectivePercentage, nil
	}

	// Simulate Ceil with integer division
	b := new(big.Int).Mul(breakEven, bitsBigInt)
	b = b.Add(b, gasPrice)
	b = b.Sub(b, big.NewInt(1)) //nolint:gomnd
	b = b.Div(b, gasPrice)

	return uint8(b.Uint64()), nil
}
