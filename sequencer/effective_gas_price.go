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
func (f *finalizer) CalculateTxBreakEvenGasPrice(tx *TxTracker, gasUsed uint64) (*big.Int, error) {
	const (
		// constants used in calculation of BreakEvenGasPrice
		signatureBytesLength           = 65
		effectivePercentageBytesLength = 1
		totalRlpFieldsLength           = signatureBytesLength + effectivePercentageBytesLength
	)

	if tx.L1GasPrice == 0 {
		log.Warn("CalculateTxBreakEvenGasPrice: L1 gas price 0. Skipping estimation for tx %s", tx.HashStr)
		return nil, ErrZeroL1GasPrice
	}

	if gasUsed == 0 {
		// Returns tx.GasPrice as the breakEvenGasPrice
		return tx.GasPrice, nil
	}

	// Get L2 Min Gas Price
	l2MinGasPrice := uint64(float64(tx.L1GasPrice) * f.effectiveGasPriceCfg.L1GasPriceFactor)
	if l2MinGasPrice < f.defaultMinGasPriceAllowed {
		l2MinGasPrice = f.defaultMinGasPriceAllowed
	}

	// Calculate BreakEvenGasPrice
	totalTxPrice := (gasUsed * l2MinGasPrice) + ((totalRlpFieldsLength + tx.BatchResources.Bytes) * f.effectiveGasPriceCfg.ByteGasCost * tx.L1GasPrice)
	breakEvenGasPrice := big.NewInt(0).SetUint64(uint64(float64(totalTxPrice/gasUsed) * f.effectiveGasPriceCfg.MarginFactor))

	return breakEvenGasPrice, nil
}

// CompareTxBreakEvenGasPrice calculates the newBreakEvenGasPrice with the newGasUsed and compares it with
// the tx.BreakEvenGasPrice. It returns ErrEffectiveGasPriceReprocess if the tx needs to be reprocessed with
// the tx.BreakEvenGasPrice updated, otherwise it returns nil
func (f *finalizer) CompareTxBreakEvenGasPrice(ctx context.Context, tx *TxTracker, newGasUsed uint64) error {
	// Increase nunber of executions related to gas price
	tx.EffectiveGasPriceProcessCount++

	newBreakEvenGasPrice, err := f.CalculateTxBreakEvenGasPrice(tx, newGasUsed)
	if err != nil {
		log.Errorf("failed to calculate breakEvenPrice with new gasUsed for tx %s, error: %s", tx.HashStr, err.Error())
		return err
	}

	// if newBreakEvenGasPrice >= tx.GasPrice then we do a final reprocess using tx.GasPrice
	if newBreakEvenGasPrice.Cmp(tx.GasPrice) >= 0 {
		tx.BreakEvenGasPrice = tx.GasPrice
		tx.IsEffectiveGasPriceFinalExecution = true
		return ErrEffectiveGasPriceReprocess
	} else { //newBreakEvenGasPrice < tx.GasPrice
		// Compute the abosulte difference between tx.BreakEvenGasPrice - newBreakEvenGasPrice
		diff := new(big.Int).Abs(new(big.Int).Sub(tx.BreakEvenGasPrice, newBreakEvenGasPrice))
		// Compute max difference allowed of breakEvenGasPrice
		maxDiff := new(big.Int).Div(new(big.Int).Mul(tx.BreakEvenGasPrice, f.maxBreakEvenGasPriceDeviationPercentage), big.NewInt(100)) //nolint:gomnd

		// if diff is greater than the maxDiff allowed
		if diff.Cmp(maxDiff) == 1 {
			if tx.EffectiveGasPriceProcessCount < 2 { //nolint:gomnd
				// it is the first process of the tx we reprocess it with the newBreakEvenGasPrice
				tx.BreakEvenGasPrice = newBreakEvenGasPrice
				return ErrEffectiveGasPriceReprocess
			} else {
				// it is the second process attempt. It makes no sense to have a big diff at
				// this point, for this reason we do a final reprocess using tx.GasPrice.
				// Also we generate a critical event as this tx needs to be analized since
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
						deviation:                     maxDiff.String(),
					},
				}
				err = f.eventLog.LogEvent(ctx, ev)
				if err != nil {
					log.Errorf("failed to log event: %s", err.Error())
				}
				return ErrEffectiveGasPriceReprocess
			}
		} // if the diff < maxDiff it is ok, no reprocess of the tx is needed
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
	// At this point we have a percentage between 1-256, we need to sub 1 to have it between 0-255 (byte)
	b = b.Sub(b, big.NewInt(1)) //nolint:gomnd

	return uint8(b.Uint64()), nil
}
