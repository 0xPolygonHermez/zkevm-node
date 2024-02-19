package elderberry

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

// PreviousProcessor is the interface that the previous processor (Etrog)
type PreviousProcessor interface {
	ProcessSequenceBatches(ctx context.Context, sequencedBatches []etherman.SequencedBatch, blockNumber uint64, l1BlockTimestamp time.Time, dbTx pgx.Tx) error
}

// ProcessorL1SequenceBatchesElderberry is the processor for SequenceBatches for Elderberry
type ProcessorL1SequenceBatchesElderberry struct {
	actions.ProcessorBase[ProcessorL1SequenceBatchesElderberry]
	previousProcessor PreviousProcessor
}

// NewProcessorL1SequenceBatchesElderberry returns instance of a processor for SequenceBatchesOrder
func NewProcessorL1SequenceBatchesElderberry(previousProcessor PreviousProcessor) *ProcessorL1SequenceBatchesElderberry {
	return &ProcessorL1SequenceBatchesElderberry{
		ProcessorBase: actions.ProcessorBase[ProcessorL1SequenceBatchesElderberry]{
			SupportedEvent:    []etherman.EventOrder{etherman.SequenceBatchesOrder},
			SupportedForkdIds: &ForksIdOnlyElderberry},
		previousProcessor: previousProcessor,
	}
}

// Process process event
func (g *ProcessorL1SequenceBatchesElderberry) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	if l1Block == nil || len(l1Block.SequencedBatches) <= order.Pos {
		return actions.ErrInvalidParams
	}
	if len(l1Block.SequencedBatches[order.Pos]) == 0 {
		log.Warnf("No sequenced batches for position")
		return nil
	}
	sbatch := l1Block.SequencedBatches[order.Pos][0]
	if sbatch.SequencedBatchElderberryData == nil {
		log.Errorf("No elderberry sequenced batch data for batch %d", sbatch.BatchNumber)
		return fmt.Errorf("no elderberry sequenced batch data for batch %d", sbatch.BatchNumber)
	}
	// We known that the MaxSequenceTimestamp is the same for all the batches
	timeLimit := time.Unix(int64(sbatch.SequencedBatchElderberryData.MaxSequenceTimestamp), 0)
	err := g.previousProcessor.ProcessSequenceBatches(ctx, l1Block.SequencedBatches[order.Pos], l1Block.BlockNumber, timeLimit, dbTx)
	return err
}
