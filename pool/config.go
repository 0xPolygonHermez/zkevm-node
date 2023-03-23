package pool

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
)

// Config is the pool configuration
type Config struct {
	// FreeClaimGasLimit is the max gas allowed use to do a free claim
	FreeClaimGasLimit uint64 `mapstructure:"FreeClaimGasLimit"`

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
}
