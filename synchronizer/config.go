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

	// CheckForLastBlockOnL1Time is the time to wait to request the
	// last block to L1 to known if we need to retrieve more data.
	// This value only apply when the system is synchronized
	CheckForLastBlockOnL1Time types.Duration `mapstructure:"CheckForLastBlockOnL1Time"`

	// Consumer Configuration for the consumer of rollup information from L1
	PerformanceCheck L1PerformanceCheckConfig `mapstructure:"PerformanceCheck"`
}

// L1PerformanceCheckConfig Configuration for the consumer of rollup information from L1
type L1PerformanceCheckConfig struct {
	// AcceptableTimeWaitingForNewRollupInfo is the expected maximum time that the consumer
	// could wait until new data are produced. If the time is greater it emmit a log to warn about
	// that. The idea is keep working the consumer as much as possible, so if the producer is not
	// fast enought then you could increse the number of parallel clients to sync with L1
	AcceptableTimeWaitingForNewRollupInfo types.Duration `mapstructure:"AcceptableTimeWaitingForNewRollupInfo"`
	// NumIterationsBeforeStartCheckingTimeWaitinfForNewRollupInfo is the number of iterations to
	// start checking the time waiting for new rollup info data
	NumIterationsBeforeStartCheckingTimeWaitinfForNewRollupInfo int `mapstructure:"NumIterationsBeforeStartCheckingTimeWaitinfForNewRollupInfo"`
}
