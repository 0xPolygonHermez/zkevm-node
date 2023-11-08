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

	// L1SynchronizationMode define how to synchronize with L1:
	// - parallel: Request data to L1 in parallel, and process sequentially. The advantage is that executor is not blocked waiting for L1 data
	// - sequential: Request data to L1 and execute
	L1SynchronizationMode string `jsonschema:"enum=sequential,enum=parallel"`
	// L1ParallelSynchronization Configuration for parallel mode (if L1SynchronizationMode equal to 'parallel')
	L1ParallelSynchronization L1ParallelSynchronizationConfig
}

// L1ParallelSynchronizationConfig Configuration for parallel mode (if UL1SynchronizationMode equal to 'parallel')
type L1ParallelSynchronizationConfig struct {
	// MaxClients Number of clients used to synchronize with L1
	MaxClients uint64
	// MaxPendingNoProcessedBlocks Size of the buffer used to store rollup information from L1, must be >= to NumberOfEthereumClientsToSync
	// sugested twice of NumberOfParallelOfEthereumClients
	MaxPendingNoProcessedBlocks uint64

	// RequestLastBlockPeriod is the time to wait to request the
	// last block to L1 to known if we need to retrieve more data.
	// This value only apply when the system is synchronized
	RequestLastBlockPeriod types.Duration

	// Consumer Configuration for the consumer of rollup information from L1
	PerformanceWarning L1PerformanceCheckConfig

	// RequestLastBlockTimeout Timeout for request LastBlock On L1
	RequestLastBlockTimeout types.Duration
	// RequestLastBlockMaxRetries Max number of retries to request LastBlock On L1
	RequestLastBlockMaxRetries int
	// StatisticsPeriod how ofter show a log with statistics (0 is disabled)
	StatisticsPeriod types.Duration
	// TimeOutMainLoop is the timeout for the main loop of the L1 synchronizer when is not updated
	TimeOutMainLoop types.Duration
	// RollupInfoRetriesSpacing is the minimum time between retries to request rollup info (it will sleep for fulfill this time) to avoid spamming L1
	RollupInfoRetriesSpacing types.Duration
	// FallbackToSequentialModeOnSynchronized if true switch to sequential mode if the system is synchronized
	FallbackToSequentialModeOnSynchronized bool
}

// L1PerformanceCheckConfig Configuration for the consumer of rollup information from L1
type L1PerformanceCheckConfig struct {
	// AceptableInacctivityTime is the expected maximum time that the consumer
	// could wait until new data is produced. If the time is greater it emmit a log to warn about
	// that. The idea is keep working the consumer as much as possible, so if the producer is not
	// fast enought then you could increse the number of parallel clients to sync with L1
	AceptableInacctivityTime types.Duration
	// ApplyAfterNumRollupReceived is the number of iterations to
	// start checking the time waiting for new rollup info data
	ApplyAfterNumRollupReceived int
}
