package gaspriceestimator

import (
	"math/big"
)

// Type different gas estimator types
type Type string

const (
	// DefaultType default gas price from config is set
	DefaultType Type = "default"
	// AllBatchesType calculate average gas used from all batches
	AllBatchesType = "allbatches"
	// LastNBatchesType calculate average gas tip from last n batches
	LastNBatchesType = "lastnbatches"
)

// Config for gas price estimator
type Config struct {
	Type Type `mapstructure:"Type"`

	DefaultPriceWei uint64   `mapstructure:"DefaultPriceWei"`
	MaxPrice        *big.Int `mapstructure:"MaxPrice"`
	IgnorePrice     *big.Int `mapstructure:"IgnorePrice"`
	CheckBlocks     int      `mapstructure:"CheckBlocks"`
	Percentile      int      `mapstructure:"Percentile"`
}
