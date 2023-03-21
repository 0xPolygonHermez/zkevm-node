package sequencer

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

// Config represents the configuration of a sequencer
type Config struct {
	// WaitPeriodSendSequence is the time the sequencer waits until
	// trying to send a sequence to L1
	WaitPeriodSendSequence types.Duration `mapstructure:"WaitPeriodSendSequence"`

	// WaitPeriodPoolIsEmpty is the time the sequencer waits until
	// trying to add new txs to the state
	WaitPeriodPoolIsEmpty types.Duration `mapstructure:"WaitPeriodPoolIsEmpty"`

	// LastBatchVirtualizationTimeMaxWaitPeriod is time since sequences should be sent
	LastBatchVirtualizationTimeMaxWaitPeriod types.Duration `mapstructure:"LastBatchVirtualizationTimeMaxWaitPeriod"`

	// BlocksAmountForTxsToBeDeleted is blocks amount after which txs will be deleted from the pool
	BlocksAmountForTxsToBeDeleted uint64 `mapstructure:"BlocksAmountForTxsToBeDeleted"`

	// FrequencyToCheckTxsForDelete is frequency with which txs will be checked for deleting
	FrequencyToCheckTxsForDelete types.Duration `mapstructure:"FrequencyToCheckTxsForDelete"`

	// MaxTxsPerBatch is the maximum amount of transactions in the batch
	MaxTxsPerBatch uint64 `mapstructure:"MaxTxsPerBatch"`

	// MaxBatchBytesSize is the maximum batch size in bytes
	// (subtracted bits of all types.Sequence fields excluding BatchL2Data from MaxTxSizeForL1)
	MaxBatchBytesSize uint64 `mapstructure:"MaxBatchBytesSize"`

	// MaxCumulativeGasUsed is max gas amount used by batch
	MaxCumulativeGasUsed uint64 `mapstructure:"MaxCumulativeGasUsed"`

	// MaxKeccakHashes is max keccak hashes used by batch
	MaxKeccakHashes uint32 `mapstructure:"MaxKeccakHashes"`

	// MaxPoseidonHashes is max poseidon hashes batch can handle
	MaxPoseidonHashes uint32 `mapstructure:"MaxPoseidonHashes"`

	// MaxPoseidonPaddings is max poseidon paddings batch can handle
	MaxPoseidonPaddings uint32 `mapstructure:"MaxPoseidonPaddings"`

	// MaxMemAligns is max mem aligns batch can handle
	MaxMemAligns uint32 `mapstructure:"MaxMemAligns"`

	// MaxArithmetics is max arithmetics batch can handle
	MaxArithmetics uint32 `mapstructure:"MaxArithmetics"`

	// MaxBinaries is max binaries batch can handle
	MaxBinaries uint32 `mapstructure:"MaxBinaries"`

	// MaxSteps is max steps batch can handle
	MaxSteps uint32 `mapstructure:"MaxSteps"`

	// WeightBatchBytesSize is the cost weight for the BatchBytesSize batch resource
	WeightBatchBytesSize int `mapstructure:"WeightBatchBytesSize"`

	// WeightCumulativeGasUsed is the cost weight for the CumulativeGasUsed batch resource
	WeightCumulativeGasUsed int `mapstructure:"WeightCumulativeGasUsed"`

	// WeightKeccakHashes is the cost weight for the KeccakHashes batch resource
	WeightKeccakHashes int `mapstructure:"WeightKeccakHashes"`

	// WeightPoseidonHashes is the cost weight for the PoseidonHashes batch resource
	WeightPoseidonHashes int `mapstructure:"WeightPoseidonHashes"`

	// WeightPoseidonPaddings is the cost weight for the PoseidonPaddings batch resource
	WeightPoseidonPaddings int `mapstructure:"WeightPoseidonPaddings"`

	// WeightMemAligns is the cost weight for the MemAligns batch resource
	WeightMemAligns int `mapstructure:"WeightMemAligns"`

	// WeightArithmetics is the cost weight for the Arithmetics batch resource
	WeightArithmetics int `mapstructure:"WeightArithmetics"`

	// WeightBinaries is the cost weight for the Binaries batch resource
	WeightBinaries int `mapstructure:"WeightBinaries"`

	// WeightSteps is the cost weight for the Steps batch resource
	WeightSteps int `mapstructure:"WeightSteps"`

	// TxLifetimeCheckTimeout is the time the sequencer waits to check txs lifetime
	TxLifetimeCheckTimeout types.Duration `mapstructure:"TxLifetimeCheckTimeout"`

	// MaxTxLifetime is the time a tx can be in the sequencer memory
	MaxTxLifetime types.Duration `mapstructure:"MaxTxLifetime"`

	// MaxTxSizeForL1 is the maximum size a single transaction can have. This field has
	// non-trivial consequences: larger transactions than 128KB are significantly harder and
	// more expensive to propagate; larger transactions also take more resources
	// to validate whether they fit into the pool or not.
	MaxTxSizeForL1 uint64 `mapstructure:"MaxTxSizeForL1"`

	// Finalizer's specific config properties
	Finalizer FinalizerCfg `mapstructure:"Finalizer"`

	// DBManager's specific config properties
	DBManager DBManagerCfg `mapstructure:"DBManager"`
}

// FinalizerCfg contains the finalizer's configuration properties
type FinalizerCfg struct {
	// GERDeadlineTimeoutInSec is the time the finalizer waits after receiving closing signal to update Global Exit Root
	GERDeadlineTimeoutInSec types.Duration `mapstructure:"GERDeadlineTimeoutInSec"`

	// SendingToL1DeadlineTimeoutInSec is the time the finalizer waits after receiving closing signal to process Forced Batches
	ForcedBatchDeadlineTimeoutInSec types.Duration `mapstructure:"ForcedBatchDeadlineTimeoutInSec"`

	// SendingToL1DeadlineTimeoutInSec is the time the finalizer waits after receiving closing signal to sends a batch to L1
	SendingToL1DeadlineTimeoutInSec types.Duration `mapstructure:"SendingToL1DeadlineTimeoutInSec"`

	// SleepDurationInMs is the time the finalizer sleeps between each iteration, if there are no transactions to be processed
	SleepDurationInMs types.Duration `mapstructure:"SleepDurationInMs"`

	// ResourcePercentageToCloseBatch is the percentage window of the resource left out for the batch to be closed
	ResourcePercentageToCloseBatch uint32 `mapstructure:"ResourcePercentageToCloseBatch"`

	// GERFinalityNumberOfBlocks is number of blocks to consider GER final
	GERFinalityNumberOfBlocks uint64 `mapstructure:"GERFinalityNumberOfBlocks"`

	// ClosingSignalsManagerWaitForCheckingL1Timeout is used by the closing signals manager to wait for its operation
	ClosingSignalsManagerWaitForCheckingL1Timeout types.Duration `mapstructure:"ClosingSignalsManagerWaitForCheckingL1Timeout"`

	// ClosingSignalsManagerWaitForCheckingGER is used by the closing signals manager to wait for its operation
	ClosingSignalsManagerWaitForCheckingGER types.Duration `mapstructure:"ClosingSignalsManagerWaitForCheckingGER"`

	// ClosingSignalsManagerWaitForCheckingL1Timeout is used by the closing signals manager to wait for its operation
	ClosingSignalsManagerWaitForCheckingForcedBatches types.Duration `mapstructure:"ClosingSignalsManagerWaitForCheckingForcedBatches"`

	// ForcedBatchesFinalityNumberOfBlocks is number of blocks to consider GER final
	ForcedBatchesFinalityNumberOfBlocks uint64 `mapstructure:"ForcedBatchesFinalityNumberOfBlocks"`

	// SenderAddress defines which private key the eth tx manager needs to use
	// to sign the L1 txs
	SenderAddress string `mapstructure:"SenderAddress"`

	// PrivateKeys defines all the key store files that are going
	// to be read in order to provide the private keys to sign the L1 txs
	PrivateKeys []types.KeystoreFileConfig `mapstructure:"PrivateKeys"`
}

// DBManagerCfg contains the DBManager's configuration properties
type DBManagerCfg struct {
	PoolRetrievalInterval types.Duration `mapstructure:"PoolRetrievalInterval"`
}
