package synchronizer

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

// Config represents the configuration of the synchronizer
type Config struct {
	// IsRollup indicates if the sequence sender is supposed to use a rollup consensus (if false it asumes validium)
	IsRollup bool `mapstructure:"IsRollup"`
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

	// TimeForCheckLastBlockOnL1Time is the time to wait to request the
	// last block to L1 to known if we need to retrieve more data.
	// This value only apply when the system is synchronized
	TimeForCheckLastBlockOnL1Time types.Duration `mapstructure:"TimeForCheckLastBlockOnL1Time"`

	// Consumer Configuration for the consumer of rollup information from L1
	PerformanceCheck L1PerformanceCheckConfig `mapstructure:"PerformanceCheck"`

	// TimeoutForRequestLastBlockOnL1 Timeout for request LastBlock On L1
	TimeoutForRequestLastBlockOnL1 types.Duration `mapstructure:"TimeoutForRequestLastBlockOnL1"`
	// MaxNumberOfRetriesForRequestLastBlockOnL1 Max number of retries to request LastBlock On L1
	MaxNumberOfRetriesForRequestLastBlockOnL1 int `mapstructure:"MaxNumberOfRetriesForRequestLastBlockOnL1"`
	// TimeForShowUpStatisticsLog how ofter show a log with statistics (0 is disabled)
	TimeForShowUpStatisticsLog types.Duration `mapstructure:"TimeForShowUpStatisticsLog"`
	// TimeOutMainLoop is the timeout for the main loop of the L1 synchronizer when is not updated
	TimeOutMainLoop types.Duration `mapstructure:"TimeOutMainLoop"`
	// MinTimeBetweenRetriesForRollupInfo is the minimum time between retries to request rollup info (it will sleep for fulfill this time) to avoid spamming L1
	MinTimeBetweenRetriesForRollupInfo types.Duration `mapstructure:"MinTimeBetweenRetriesForRollupInfo"`
}

// L1PerformanceCheckConfig Configuration for the consumer of rollup information from L1
type L1PerformanceCheckConfig struct {
	// AcceptableTimeWaitingForNewRollupInfo is the expected maximum time that the consumer
	// could wait until new data is produced. If the time is greater it emmit a log to warn about
	// that. The idea is keep working the consumer as much as possible, so if the producer is not
	// fast enought then you could increse the number of parallel clients to sync with L1
	AcceptableTimeWaitingForNewRollupInfo types.Duration `mapstructure:"AcceptableTimeWaitingForNewRollupInfo"`
	// NumIterationsBeforeStartCheckingTimeWaitinfForNewRollupInfo is the number of iterations to
	// start checking the time waiting for new rollup info data
	NumIterationsBeforeStartCheckingTimeWaitinfForNewRollupInfo int `mapstructure:"NumIterationsBeforeStartCheckingTimeWaitinfForNewRollupInfo"`
}
