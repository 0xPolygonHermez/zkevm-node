package sequencer

import (
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	base "github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/profitabilitychecker"
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

	// MaxTimeForBatchToBeOpen is time after which new batch should be closed
	MaxTimeForBatchToBeOpen types.Duration `mapstructure:"MaxTimeForBatchToBeOpen"`

	// BlocksAmountForTxsToBeDeleted is blocks amount after which txs will be deleted from the pool
	BlocksAmountForTxsToBeDeleted uint64 `mapstructure:"BlocksAmountForTxsToBeDeleted"`

	// FrequencyToCheckTxsForDelete is frequency with which txs will be checked for deleting
	FrequencyToCheckTxsForDelete types.Duration `mapstructure:"FrequencyToCheckTxsForDelete"`

	// MaxCumulativeGasUsed is max gas amount used by batch
	MaxCumulativeGasUsed uint64 `mapstructure:"MaxCumulativeGasUsed"`

	// MaxKeccakHashes is max keccak hashes used by batch
	MaxKeccakHashes int32 `mapstructure:"MaxKeccakHashes"`

	// MaxPoseidonHashes is max poseidon hashes batch can handle
	MaxPoseidonHashes int32 `mapstructure:"MaxPoseidonHashes"`

	// MaxPoseidonPaddings is max poseidon paddings batch can handle
	MaxPoseidonPaddings int32 `mapstructure:"MaxPoseidonPaddings"`

	// MaxMemAligns is max mem aligns batch can handle
	MaxMemAligns int32 `mapstructure:"MaxMemAligns"`

	// MaxArithmetics is max arithmetics batch can handle
	MaxArithmetics int32 `mapstructure:"MaxArithmetics"`

	// MaxBinaries is max binaries batch can handle
	MaxBinaries int32 `mapstructure:"MaxBinaries"`

	// MaxSteps is max steps batch can handle
	MaxSteps int32 `mapstructure:"MaxSteps"`

	// ProfitabilityChecker configuration
	ProfitabilityChecker profitabilitychecker.Config `mapstructure:"ProfitabilityChecker"`

	// Maximum size, in gas size, a sequence can reach
	MaxSequenceSize MaxSequenceSize `mapstructure:"MaxSequenceSize"`

	// Maximum allowed failed counter for the tx before it becomes invalid
	MaxAllowedFailedCounter uint64 `mapstructure:"MaxAllowedFailedCounter"`
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
