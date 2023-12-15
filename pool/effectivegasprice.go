package pool

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/state"
)

var (
	// ErrEffectiveGasPriceEmpty happens when the effectiveGasPrice or gasPrice is nil or zero
	ErrEffectiveGasPriceEmpty = errors.New("effectiveGasPrice or gasPrice cannot be nil or zero")

	// ErrEffectiveGasPriceIsZero happens when the calculated EffectiveGasPrice is zero
	ErrEffectiveGasPriceIsZero = errors.New("effectiveGasPrice cannot be zero")
)

// EffectiveGasPrice implements the effective gas prices calculations and checks
type EffectiveGasPrice struct {
	cfg                EffectiveGasPriceCfg
	minGasPriceAllowed uint64
}

// NewEffectiveGasPrice creates and initializes an instance of EffectiveGasPrice
func NewEffectiveGasPrice(cfg EffectiveGasPriceCfg, minGasPriceAllowed uint64) *EffectiveGasPrice {
	return &EffectiveGasPrice{
		cfg:                cfg,
		minGasPriceAllowed: minGasPriceAllowed,
	}
}

// IsEnabled return if effectiveGasPrice calculation is enabled
func (e *EffectiveGasPrice) IsEnabled() bool {
	return e.cfg.Enabled
}

// GetFinalDeviation return the value for the config parameter FinalDeviationPct
func (e *EffectiveGasPrice) GetFinalDeviation() uint64 {
	return e.cfg.FinalDeviationPct
}

// GetTxAndL2GasPrice return the tx gas price and l2 suggested gas price to use in egp calculations
// If egp is disabled we will use a "simulated" tx and l2 gas price, that is calculated using the L2GasPriceSuggesterFactor config param
func (e *EffectiveGasPrice) GetTxAndL2GasPrice(txGasPrice *big.Int, l1GasPrice uint64, l2GasPrice uint64) (egpTxGasPrice *big.Int, egpL2GasPrice uint64) {
	if !e.cfg.Enabled {
		// If egp is not enabled we use the L2GasPriceSuggesterFactor to calculate the "simulated" suggested L2 gas price
		gp := new(big.Int).SetUint64(uint64(e.cfg.L2GasPriceSuggesterFactor * float64(l1GasPrice)))
		return gp, gp.Uint64()
	} else {
		return txGasPrice, l2GasPrice
	}
}

// CalculateBreakEvenGasPrice calculates the break even gas price for a transaction
func (e *EffectiveGasPrice) CalculateBreakEvenGasPrice(rawTx []byte, txGasPrice *big.Int, txGasUsed uint64, l1GasPrice uint64) (*big.Int, error) {
	const (
		// constants used in calculation of BreakEvenGasPrice
		signatureBytesLength           = 65
		effectivePercentageBytesLength = 1
		constBytesTx                   = signatureBytesLength + effectivePercentageBytesLength
	)

	if l1GasPrice == 0 {
		return nil, ErrZeroL1GasPrice
	}

	if txGasUsed == 0 {
		// Returns tx.GasPrice as the breakEvenGasPrice
		return txGasPrice, nil
	}

	// Get L2 Min Gas Price
	l2MinGasPrice := uint64(float64(l1GasPrice) * e.cfg.L1GasPriceFactor)
	if l2MinGasPrice < e.minGasPriceAllowed {
		l2MinGasPrice = e.minGasPriceAllowed
	}

	txZeroBytes := uint64(bytes.Count(rawTx, []byte{0}))
	txNonZeroBytes := uint64(len(rawTx)) - txZeroBytes

	// Calculate BreakEvenGasPrice
	totalTxPrice := (txGasUsed * l2MinGasPrice) +
		((constBytesTx+txNonZeroBytes)*e.cfg.ByteGasCost+txZeroBytes*e.cfg.ZeroByteGasCost)*l1GasPrice
	breakEvenGasPrice := new(big.Int).SetUint64(uint64(float64(totalTxPrice/txGasUsed) * e.cfg.NetProfit))

	return breakEvenGasPrice, nil
}

// CalculateEffectiveGasPrice calculates the final effective gas price for a tx
func (e *EffectiveGasPrice) CalculateEffectiveGasPrice(rawTx []byte, txGasPrice *big.Int, txGasUsed uint64, l1GasPrice uint64, l2GasPrice uint64) (*big.Int, error) {
	breakEvenGasPrice, err := e.CalculateBreakEvenGasPrice(rawTx, txGasPrice, txGasUsed, l1GasPrice)

	if err != nil {
		return nil, err
	}

	bfL2GasPrice := new(big.Float).SetUint64(l2GasPrice)
	bfTxGasPrice := new(big.Float).SetInt(txGasPrice)

	ratioPriority := new(big.Float).SetFloat64(1.0)

	if bfTxGasPrice.Cmp(bfL2GasPrice) == 1 {
		//ratioPriority = (txGasPrice / l2GasPrice)
		ratioPriority = new(big.Float).Quo(bfTxGasPrice, bfL2GasPrice)
	}

	bfEffectiveGasPrice := new(big.Float).Mul(new(big.Float).SetInt(breakEvenGasPrice), ratioPriority)

	effectiveGasPrice := new(big.Int)
	bfEffectiveGasPrice.Int(effectiveGasPrice)

	if effectiveGasPrice.Cmp(new(big.Int).SetUint64(0)) == 0 {
		return nil, ErrEffectiveGasPriceIsZero
	}

	return effectiveGasPrice, nil
}

// CalculateEffectiveGasPricePercentage calculates the gas price's effective percentage
func (e *EffectiveGasPrice) CalculateEffectiveGasPricePercentage(gasPrice *big.Int, effectiveGasPrice *big.Int) (uint8, error) {
	const bits = 256
	var bitsBigInt = big.NewInt(bits)

	if effectiveGasPrice == nil || gasPrice == nil ||
		gasPrice.Cmp(big.NewInt(0)) == 0 || effectiveGasPrice.Cmp(big.NewInt(0)) == 0 {
		return 0, ErrEffectiveGasPriceEmpty
	}

	if gasPrice.Cmp(effectiveGasPrice) <= 0 {
		return state.MaxEffectivePercentage, nil
	}

	// Simulate Ceil with integer division
	b := new(big.Int).Mul(effectiveGasPrice, bitsBigInt)
	b = b.Add(b, gasPrice)
	b = b.Sub(b, big.NewInt(1)) //nolint:gomnd
	b = b.Div(b, gasPrice)
	// At this point we have a percentage between 1-256, we need to sub 1 to have it between 0-255 (byte)
	b = b.Sub(b, big.NewInt(1)) //nolint:gomnd

	return uint8(b.Uint64()), nil
}
