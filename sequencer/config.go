package sequencer

import (
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	base "github.com/0xPolygonHermez/zkevm-node/encoding"
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

	// WaitBlocksToUpdateGER is number of blocks for sequencer to wait
	WaitBlocksToUpdateGER uint64 `mapstructure:"WaitBlocksToUpdateGER"`

	// WaitBlocksToConsiderGerFinal is number of blocks for sequencer to consider globalExitRoot final
	WaitBlocksToConsiderGerFinal uint64 `mapstructure:"WaitBlocksToConsiderGerFinal"`

	// ElapsedTimeToCloseBatchWithoutTxsDueToNewGER it's time to close a batch bcs new globalExitRoot appeared
	ElapsedTimeToCloseBatchWithoutTxsDueToNewGER types.Duration `mapstructure:"ElapsedTimeToCloseBatchWithoutTxsDueToNewGER"`

	// MaxTimeForBatchToBeOpen is time after which new batch should be closed
	MaxTimeForBatchToBeOpen types.Duration `mapstructure:"MaxTimeForBatchToBeOpen"`

	// BlocksAmountForTxsToBeDeleted is blocks amount after which txs will be deleted from the pool
	BlocksAmountForTxsToBeDeleted uint64 `mapstructure:"BlocksAmountForTxsToBeDeleted"`

	// FrequencyToCheckTxsForDelete is frequency with which txs will be checked for deleting
	FrequencyToCheckTxsForDelete types.Duration `mapstructure:"FrequencyToCheckTxsForDelete"`

	// MaxTxsPerBatch is the maximum amount of transactions in the batch
	MaxTxsPerBatch uint64 `mapstructure:"MaxTxsPerBatch"`

	// MaxBatchBytesSize is the maximum batch size in bytes
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

	// Maximum size, in gas size, a sequence can reach
	MaxSequenceSize MaxSequenceSize `mapstructure:"MaxSequenceSize"`

	// Maximum allowed failed counter for the tx before it becomes invalid
	MaxAllowedFailedCounter uint64 `mapstructure:"MaxAllowedFailedCounter"`

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

	// Finalizer's specific config properties
	Finalizer FinalizerCfg `mapstructure:"Finalizer"`
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

	// ClosingSignalsManagerWaitForL1OperationsInSec is used by the closing signals manager to wait for its operation
	ClosingSignalsManagerWaitForL1OperationsInSec types.Duration `mapstructure:"ClosingSignalsManagerWaitForL1OperationsInSec"`
}

// MaxSequenceSize is a wrapper type that parses token amount to big int
type MaxSequenceSize struct {
	*big.Int `validate:"required"`
}

// UnmarshalText unmarshal token amount from float string to big int
func (m *MaxSequenceSize) UnmarshalText(data []byte) error {
	amount, ok := new(big.Int).SetString(string(data), base.Base10)
	if !ok {
		return fmt.Errorf("failed to unmarshal string to float")
	}
	m.Int = amount

	return nil
}
