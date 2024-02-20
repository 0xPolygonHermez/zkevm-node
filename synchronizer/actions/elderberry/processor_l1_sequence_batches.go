package elderberry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

var (
	ErrInvalidInitialBatchNumber = errors.New("invalid initial batch number")
)

// PreviousProcessor is the interface that the previous processor (Etrog)
type PreviousProcessor interface {
	ProcessSequenceBatches(ctx context.Context, sequencedBatches []etherman.SequencedBatch, blockNumber uint64, l1BlockTimestamp time.Time, dbTx pgx.Tx) error
}

type stateL1SequenceBatchesElderberry interface {
	GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetL2BlocksByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]state.L2Block, error)
}

// ProcessorL1SequenceBatchesElderberry is the processor for SequenceBatches for Elderberry
type ProcessorL1SequenceBatchesElderberry struct {
	actions.ProcessorBase[ProcessorL1SequenceBatchesElderberry]
	previousProcessor PreviousProcessor
	state             stateL1SequenceBatchesElderberry
}

// NewProcessorL1SequenceBatchesElderberry returns instance of a processor for SequenceBatchesOrder
func NewProcessorL1SequenceBatchesElderberry(previousProcessor PreviousProcessor, state stateL1SequenceBatchesElderberry) *ProcessorL1SequenceBatchesElderberry {
	return &ProcessorL1SequenceBatchesElderberry{
		ProcessorBase: actions.ProcessorBase[ProcessorL1SequenceBatchesElderberry]{
			SupportedEvent:    []etherman.EventOrder{etherman.SequenceBatchesOrder},
			SupportedForkdIds: &ForksIdOnlyElderberry},
		previousProcessor: previousProcessor,
		state:             state,
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
	// We need to check that the sequence match
	err := g.sanityCheckExpectedSequence(sbatch.SequencedBatchElderberryData.InitSequencedBatchNumber, dbTx)
	if err != nil {
		return err
	}
	// We known that the MaxSequenceTimestamp is the same for all the batches so we can use the first one
	timeLimit := time.Unix(int64(sbatch.SequencedBatchElderberryData.MaxSequenceTimestamp), 0)
	err = g.previousProcessor.ProcessSequenceBatches(ctx, l1Block.SequencedBatches[order.Pos], l1Block.BlockNumber, timeLimit, dbTx)
	// The last L2block timestamp must match MaxSequenceTimestamp
	if err != nil {
		return err
	}
	err = g.sanityCheckTstampLastL2Block(timeLimit, dbTx)
	if err != nil {
		return err
	}
	return err
}

func (g *ProcessorL1SequenceBatchesElderberry) sanityCheckExpectedSequence(initialBatchNumber uint64, dbTx pgx.Tx) error {
	// We need to check that the sequence match
	lastVirtualBatchNum, err := g.state.GetLastVirtualBatchNum(context.Background(), dbTx)
	if err != nil {
		log.Errorf("Error getting last virtual batch number: %s", err)
		return err
	}
	if lastVirtualBatchNum+1 != initialBatchNumber {
		log.Errorf("The last virtual batch number is not the expected one. Expected: %d (last on DB), got: %d (L1 event)", lastVirtualBatchNum+1, initialBatchNumber)
		return fmt.Errorf("the last virtual batch number is not the expected one. Expected: %d (last on DB), got: %d (L1 event) err:%w", lastVirtualBatchNum+1, initialBatchNumber, ErrInvalidInitialBatchNumber)
	}
	return nil
}

func (g *ProcessorL1SequenceBatchesElderberry) sanityCheckTstampLastL2Block(timeLimit time.Time, dbTx pgx.Tx) error {
	lastVirtualBatchNum, err := g.state.GetLastVirtualBatchNum(context.Background(), dbTx)
	if err != nil {
		log.Errorf("Error getting last virtual batch number: %s", err)
		return err
	}
	l2blocks, err := g.state.GetL2BlocksByBatchNumber(context.Background(), lastVirtualBatchNum, dbTx)
	if err != nil {
		log.Errorf("Error getting last virtual batch number: %s", err)
		return err
	}
	if len(l2blocks) == 0 {
		return nil
	}
	lastL2Block := l2blocks[len(l2blocks)-1]
	if lastL2Block.ReceivedAt != timeLimit {
		log.Errorf("The last L2 block timestamp is not the expected one. Expected: %s (L1 event), got: %s (last L2Block)", timeLimit, lastL2Block.ReceivedAt)
		return fmt.Errorf("dont match last L2 block timestamp with L1 event timestamp")
	}
	return nil
}
