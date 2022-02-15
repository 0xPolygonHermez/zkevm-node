package gasprice

import (
	"math/big"
)

// EstimatorType different gas estimator types.
type EstimatorType string

const (
	// DefaultType default gas price from config is set.
	DefaultType EstimatorType = "default"
	// AllBatchesType calculate average gas used from all batches.
	AllBatchesType EstimatorType = "allbatches"
	// LastNBatchesType calculate average gas tip from last n batches.
	LastNBatchesType EstimatorType = "lastnbatches"
)

// Config for gas price estimator.
type Config struct {
	Type EstimatorType `mapstructure:"Type"`

	DefaultGasPriceWei uint64   `mapstructure:"DefaultGasPriceWei"`
	MaxPrice           *big.Int `mapstructure:"MaxPrice"`
	IgnorePrice        *big.Int `mapstructure:"IgnorePrice"`
	CheckBlocks        int      `mapstructure:"CheckBlocks"`
	Percentile         int      `mapstructure:"Percentile"`
}
