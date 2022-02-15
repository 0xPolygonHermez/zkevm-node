package gaspriceestimator

import (
	"math/big"
)

// Type different gas estimator types
type Type string

// Config for gas price estimator
type Config struct {
	Type Type `mapstructure:"Type"`

	DefaultPriceWei uint64   `mapstructure:"DefaultPriceWei"`
	MaxPrice        *big.Int `mapstructure:"MaxPrice"`
	IgnorePrice     *big.Int `mapstructure:"IgnorePrice"`
	CheckBlocks     int      `mapstructure:"CheckBlocks"`
	Percentile      int      `mapstructure:"Percentile"`
}
