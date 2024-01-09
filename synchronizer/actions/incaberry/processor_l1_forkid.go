package incaberry

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

type stateProcessorForkIdInterface interface {
	GetLastBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetForkIDs(ctx context.Context, dbTx pgx.Tx) ([]state.ForkIDInterval, error)
	AddForkIDInterval(ctx context.Context, newForkID state.ForkIDInterval, dbTx pgx.Tx) error
	ResetForkID(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
}

type syncProcessorForkIdInterface interface {
	IsTrustedSequencer() bool
}

// ProcessorForkId implements L1EventProcessor
type ProcessorForkId struct {
	actions.ProcessorBase[ProcessorForkId]
	state stateProcessorForkIdInterface
	sync  syncProcessorForkIdInterface
}

// NewProcessorForkId returns instance of a processor for ForkIDsOrder
func NewProcessorForkId(state stateProcessorForkIdInterface, sync syncProcessorForkIdInterface) *ProcessorForkId {
	return &ProcessorForkId{
		ProcessorBase: actions.ProcessorBase[ProcessorForkId]{
			SupportedEvent:    []etherman.EventOrder{etherman.ForkIDsOrder},
			SupportedForkdIds: &actions.ForksIdAll,
		},
		state: state,
		sync:  sync}
}

// Process process event
func (p *ProcessorForkId) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	return p.processForkID(ctx, l1Block.ForkIDs[order.Pos], l1Block.BlockNumber, dbTx)
}

func (s *ProcessorForkId) processForkID(ctx context.Context, forkID etherman.ForkID, blockNumber uint64, dbTx pgx.Tx) error {
	fID := state.ForkIDInterval{
		FromBatchNumber: forkID.BatchNumber + 1,
		ToBatchNumber:   math.MaxUint64,
		ForkId:          forkID.ForkID,
		Version:         forkID.Version,
		BlockNumber:     blockNumber,
	}

	// If forkID affects to a batch from the past. State must be reseted.
	log.Debugf("ForkID: %d, synchronization must use the new forkID since batch: %d", forkID.ForkID, forkID.BatchNumber+1)
	fIds, err := s.state.GetForkIDs(ctx, dbTx)
	if err != nil {
		log.Error("error getting ForkIDTrustedReorg. Error: ", err)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state get forkID trusted state. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}
	if len(fIds) != 0 && fIds[len(fIds)-1].ForkId == fID.ForkId { // If the forkID reset was already done
		return nil
	}
	//If the forkID.batchnumber is a future batch
	latestBatchNumber, err := s.state.GetLastBatchNumber(ctx, dbTx)
	if err != nil && !errors.Is(err, state.ErrStateNotSynchronized) {
		log.Error("error getting last batch number. Error: ", err)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}
	// Add new forkID to the state
	err = s.state.AddForkIDInterval(ctx, fID, dbTx)
	if err != nil {
		log.Error("error adding new forkID interval to the state. Error: ", err)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}
	if latestBatchNumber <= forkID.BatchNumber || s.sync.IsTrustedSequencer() { //If the forkID will start in a future batch or isTrustedSequencer
		log.Infof("Just adding forkID. Skipping reset forkID. ForkID: %+v.", fID)
		return nil
	}

	log.Info("ForkID received in the permissionless node that affects to a batch from the past")
	//Reset DB only if permissionless node
	log.Debugf("ForkID: %d, Reverting synchronization to batch: %d", forkID.ForkID, forkID.BatchNumber+1)
	err = s.state.ResetForkID(ctx, forkID.BatchNumber+1, dbTx)
	if err != nil {
		log.Error("error resetting the state. Error: ", err)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}

	// Commit because it returns an error to force the resync
	err = dbTx.Commit(ctx)
	if err != nil {
		log.Error("error committing the resetted state. Error: ", err)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}

	return fmt.Errorf("new ForkID detected, reseting synchronizarion")
}
