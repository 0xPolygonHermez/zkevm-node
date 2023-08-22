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

	// DefaultGasPriceWei is used to set the gas price to be used by the default gas pricer or as minimim gas price by the follower gas pricer.
	DefaultGasPriceWei uint64 `mapstructure:"DefaultGasPriceWei"`
	// MaxGasPriceWei is used to limit the gas price returned by the follower gas pricer to a maximum value. It is ignored if 0.
	MaxGasPriceWei            uint64         `mapstructure:"MaxGasPriceWei"`
	MaxPrice                  *big.Int       `mapstructure:"MaxPrice"`
	IgnorePrice               *big.Int       `mapstructure:"IgnorePrice"`
	CheckBlocks               int            `mapstructure:"CheckBlocks"`
	Percentile                int            `mapstructure:"Percentile"`
	UpdatePeriod              types.Duration `mapstructure:"UpdatePeriod"`
	CleanHistoryPeriod        types.Duration `mapstructure:"CleanHistoryPeriod"`
	CleanHistoryTimeRetention types.Duration `mapstructure:"CleanHistoryTimeRetention"`

	Factor float64 `mapstructure:"Factor"`
}
