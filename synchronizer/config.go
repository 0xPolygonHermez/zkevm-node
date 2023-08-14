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

	// L1ParallelSynchronization Use new L1 synchronization that do in parallel request to L1 and process the data
	// If false use the legacy sequential mode
	UseParallelModeForL1Synchronization bool `mapstructure:"UseParallelModeForL1Synchronization"`
	// L1ParallelSynchronization Configuration for parallel mode (if UseParallelModeForL1Synchronization is true)
	L1ParallelSynchronization L1ParallelSynchronizationConfig `mapstructure:"L1ParallelSynchronization"`
}

// L1ParallelSynchronizationConfig Configuration for parallel mode (if UseParallelModeForL1Synchronization is true)
type L1ParallelSynchronizationConfig struct {
	// NumberOfParallelOfEthereumClients Number of clients used to synchronize with L1
	// (if UseParallelModeForL1Synchronization is true)
	NumberOfParallelOfEthereumClients uint64 `mapstructue:"NumberOfParallelOfEthereumClients"`
	// CapacityOfBufferingRollupInfoFromL1 Size of the buffer used to store rollup information from L1, must be >= to NumberOfEthereumClientsToSync
	// sugested twice of NumberOfParallelOfEthereumClients
	// (if UseParallelModeForL1Synchronization is true)
	CapacityOfBufferingRollupInfoFromL1 uint64 `mapstructure:"CapacityOfBufferingRollupInfoFromL1"`
}
