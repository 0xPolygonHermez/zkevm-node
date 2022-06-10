package synchronizerv2

import (
	"github.com/hermeznetwork/hermez-core/config"
)

// Config represents the configuration of the synchronizer
type Config struct {
	// SyncInterval is the delay interval between reading new rollup information
	SyncInterval config.Duration `mapstructure:"SyncInterval"`

	// SyncChunkSize is the number of blocks to sync on each chunk
	SyncChunkSize uint64 `mapstructure:"SyncChunkSize"`
}
