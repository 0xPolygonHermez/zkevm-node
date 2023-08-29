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

	// IntervalToRefreshGasPrices is the time to wait to refresh the gas prices
	IntervalToRefreshGasPrices types.Duration `mapstructure:"IntervalToRefreshGasPrices"`

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

	// AccountQueue represents the maximum number of non-executable transaction slots permitted per account
	AccountQueue uint64 `mapstructure:"AccountQueue"`

	// GlobalQueue represents the maximum number of non-executable transaction slots for all accounts
	GlobalQueue uint64 `mapstructure:"GlobalQueue"`
}
