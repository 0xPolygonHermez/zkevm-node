package synchronizer

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

// Config represents the configuration of the synchronizer
type Config struct {
	// SyncInterval is the delay interval between reading new rollup information
	SyncInterval types.Duration `mapstructure:"SyncInterval"`
	// SyncChunkSize is the number of blocks to sync on each chunk
	SyncChunkSize uint64 `mapstructure:"SyncChunkSize"`
	// TrustedSequencerURL is the rpc url to connect and sync the trusted state
	TrustedSequencerURL string `mapstructure:"TrustedSequencerURL"`
	// NumberOfEthereumClientsToSync Number of clients used to synchronize with L1
	NumberOfEthereumClientsToSync uint64 `mapstructue:"NumberOfEthereumClientsToSync"`
	// CapacityOfBufferingRollupInfoFromL1 Size of the buffer used to store rollup information from L1, must be >= to NumberOfEthereumClientsToSync
	CapacityOfBufferingRollupInfoFromL1 uint64 `mapstructure:"CapacityOfBufferingRollupInfoFromL1"`
}
