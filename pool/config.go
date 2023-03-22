package pool

import "github.com/0xPolygonHermez/zkevm-node/db"

// Config is the pool configuration
type Config struct {
	// FreeClaimGasLimit is the max gas allowed use to do a free claim
	FreeClaimGasLimit  uint64    `mapstructure:"FreeClaimGasLimit"`
	MaxTxBytesSize     uint64    `mapstructure:"MaxTxBytesSize"`
	MaxTxDataBytesSize int       `mapstructure:"MaxTxDataBytesSize"`
	DB                 db.Config `mapstructure:"DB"`
}
