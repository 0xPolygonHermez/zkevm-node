package gasprice

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

// EstimatorType different gas estimator types.
type EstimatorType string

const (
	// DefaultType default gas price from config is set.
	DefaultType EstimatorType = "default"
	// LastNBatchesType calculate average gas tip from last n batches.
	LastNBatchesType EstimatorType = "lastnbatches"
	// FollowerType calculate the gas price basing on the L1 gasPrice.
	FollowerType EstimatorType = "follower"
)

// Config for gas price estimator.
type Config struct {
	Type EstimatorType `mapstructure:"Type"`

	DefaultGasPriceWei uint64         `mapstructure:"DefaultGasPriceWei"`
	MaxPrice           *big.Int       `mapstructure:"MaxPrice"`
	IgnorePrice        *big.Int       `mapstructure:"IgnorePrice"`
	CheckBlocks        int            `mapstructure:"CheckBlocks"`
	Percentile         int            `mapstructure:"Percentile"`
	UpdatePeriod       types.Duration `mapstructure:"UpdatePeriod"`
	Factor             float64        `mapstructure:"Factor"`

	// BreakEvenGasPriceCalcCfg is the configuration for the break even gas price calculation
	BreakEvenGasPriceCalcCfg BreakEvenGasPriceCalculationCfg `mapstructure:"BreakEvenGasPriceCalculationCfg"`
}

// BreakEvenGasPriceCalculationCfg has parameters for the gas price calculation.
// TODO: Add config tests
type BreakEvenGasPriceCalculationCfg struct {
	// L1GasPricePercentageForL2MinPrice is the percentage of the L1 gas price that will be used as the L2 min gas price
	L1GasPricePercentageForL2MinPrice uint64 `mapstructure:"L1GasPricePercentageForL2MinPrice"`

	// ByteGasCost is the gas cost per byte
	ByteGasCost uint64 `mapstructure:"ByteGasCost"`

	// MarginFactorPercentage is the margin factor percentage to be added to the L2 min gas price
	MarginFactorPercentage uint64 `mapstructure:"MarginFactorPercentage"`

	// TxPriceGuaranteePeriod is the period of time that the gas price will be guaranteed
	TxPriceGuaranteePeriod types.Duration `mapstructure:"TxPriceGuaranteePeriod"`

	// MaxBreakEvenGasPriceDeviationPercentage is the max allowed deviation percentage BreakEvenGasPrice on re-calculation
	MaxBreakEvenGasPriceDeviationPercentage float64 `mapstructure:"MaxBreakEvenGasPriceDeviationPercentage"`
}
