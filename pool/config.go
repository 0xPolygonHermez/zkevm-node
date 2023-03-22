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

	// MinGasPrice is the minimum gas price allowed for a tx
	MinGasPrice uint64 `mapstructure:"MinGasPrice"`

	// MinSuggestedGasPriceInterval is the interval to look back of the suggested min gas price for a tx
	MinSuggestedGasPriceInterval types.Duration `mapstructure:"MinSuggestedGasPriceInterval"`
}
