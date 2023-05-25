package pool

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
)

// Config is the pool configuration
type Config struct {
	// IntervalToRefreshBlockedAddresses is the time it takes to sync the
	// blocked address list from db to memory
	IntervalToRefreshBlockedAddresses types.Duration `mapstructure:"IntervalToRefreshBlockedAddresses"`

	// MaxTxBytesSize is the max size of a transaction in bytes
	MaxTxBytesSize uint64 `mapstructure:"MaxTxBytesSize"`

	// MaxTxDataBytesSize is the max size of the data field of a transaction in bytes
	MaxTxDataBytesSize int `mapstructure:"MaxTxDataBytesSize"`

	// DB is the database configuration
	DB db.Config `mapstructure:"DB"`

	// DefaultMinGasPriceAllowed is the default min gas price to suggest
	DefaultMinGasPriceAllowed uint64 `mapstructure:"DefaultMinGasPriceAllowed"`

	// MinAllowedGasPriceInterval is the interval to look back of the suggested min gas price for a tx
	MinAllowedGasPriceInterval types.Duration `mapstructure:"MinAllowedGasPriceInterval"`

	// PollMinAllowedGasPriceInterval is the interval to poll the suggested min gas price for a tx
	PollMinAllowedGasPriceInterval types.Duration `mapstructure:"PollMinAllowedGasPriceInterval"`

	// EffectiveGasPrice is the configuration for the break even and effective gas price calculation
	EffectiveGasPrice EffectiveGasPrice `mapstructure:"EffectiveGasPrice"`
}

// EffectiveGasPrice has parameters for the effective gas price calculation.
type EffectiveGasPrice struct {
	// L1GasPricePercentageForL2MinPrice is the percentage of the L1 gas price that will be used as the L2 min gas price
	L1GasPricePercentageForL2MinPrice uint64 `mapstructure:"L1GasPricePercentageForL2MinPrice"`

	// ByteGasCost is the gas cost per byte
	ByteGasCost uint64 `mapstructure:"ByteGasCost"`

	// MarginFactorPercentage is the margin factor percentage to be added to the L2 min gas price
	MarginFactorPercentage uint64 `mapstructure:"MarginFactorPercentage"`
}
