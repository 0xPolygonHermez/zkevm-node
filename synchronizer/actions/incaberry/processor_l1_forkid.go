package incaberry

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"

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
	UpdateForkIDBlockNumber(ctx context.Context, forkdID uint64, newBlockNumber uint64, updateMemCache bool, dbTx pgx.Tx) error
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
	if l1Block == nil {
		return errors.New("nil l1Block")
	}
	if len(l1Block.ForkIDs) <= order.Pos {
		return fmt.Errorf("ForkIDsOrder index out of range. BlockNumber: %d, ForkIDsOrder index: %d", l1Block.BlockNumber, order.Pos)
	}
	return p.processForkID(ctx, l1Block.ForkIDs[order.Pos], l1Block.BlockNumber, dbTx)
}

func getForkdFromSlice(fIds []state.ForkIDInterval, forkId uint64) (bool, state.ForkIDInterval) {
	if len(fIds) == 0 {
		return false, state.ForkIDInterval{}
	}
	for _, f := range fIds {
		if f.ForkId == forkId {
			return true, f
		}
	}
	return false, state.ForkIDInterval{}
}

func isForksSameFromBatchNumber(f1, f2 state.ForkIDInterval) bool {
	return f1.ForkId == f2.ForkId && f1.FromBatchNumber == f2.FromBatchNumber
}

func isIncommingForkIdGreatestThanLastOne(incommingForkID state.ForkIDInterval, fIds []state.ForkIDInterval) bool {
	if len(fIds) == 0 {
		return true
	}
	last := lastForkID(fIds)
	// Must be greater than the last one
	return incommingForkID.ForkId > last
}

func lastForkID(fIds []state.ForkIDInterval) uint64 {
	if len(fIds) == 0 {
		return 0
	}
	sort.Slice(fIds, func(i, j int) bool {
		return fIds[i].ForkId > fIds[j].ForkId
	})
	return fIds[0].ForkId
}

// return true if have been update or false if it's a new one
func (s *ProcessorForkId) updateForkIDIfNeeded(ctx context.Context, forkIDincomming state.ForkIDInterval, forkIDsInState []state.ForkIDInterval, dbTx pgx.Tx) (bool, error) {
	found, dbForkID := getForkdFromSlice(forkIDsInState, forkIDincomming.ForkId)
	if !found {
		// Is a new forkid
		return false, nil
	}
	if isForksSameFromBatchNumber(forkIDincomming, dbForkID) {
		if forkIDincomming.BlockNumber != dbForkID.BlockNumber {
			isLastForkId := lastForkID(forkIDsInState) == forkIDincomming.ForkId
			log.Infof("ForkID: %d, received again: same fork_id but different blockNumber old: %d, new: %d", forkIDincomming.ForkId, dbForkID.BlockNumber, forkIDincomming.BlockNumber)
			if isLastForkId {
				log.Warnf("ForkID: %d is the last one in the state. Updating BlockNumber from %d to %d", forkIDincomming.ForkId, dbForkID.BlockNumber, forkIDincomming.BlockNumber)
				err := s.state.UpdateForkIDBlockNumber(ctx, forkIDincomming.ForkId, forkIDincomming.BlockNumber, true, dbTx)
				if err != nil {
					log.Errorf("error updating forkID: %d blocknumber. Error: %v", forkIDincomming.ForkId, err)
					return true, err
				}
				return true, nil
			}
			err := fmt.Errorf("ForkID: %d, already in the state but with different blockNumber and is not last ForkID, so can't update BlockNumber. DB ForkID: %+v. New ForkID: %+v", forkIDincomming.ForkId, dbForkID, forkIDincomming)
			log.Error(err.Error())
			return true, err
		}
		log.Infof("ForkID: %d, already in the state. Skipping . ForkID: %+v.", forkIDincomming.ForkId, forkIDincomming)
		return true, nil
	}
	err := fmt.Errorf("ForkID: %d, already in the state but with different starting BatchNumber. DB ForkID: %+v. New ForkID: %+v", forkIDincomming.ForkId, dbForkID, forkIDincomming)
	log.Error(err.Error())
	return true, err
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
		log.Errorf("error getting forkIDs. Error: %v", err)
		//return s.rollback(ctx, blockNumber, "getting GetForkIDs", err, dbTx)
		return err
	}
	isUpdate, err := s.updateForkIDIfNeeded(ctx, fID, fIds, dbTx)
	if err != nil {
		log.Errorf("ForkID: %d, error updating forkID . Error: %v", forkID.ForkID, err)
		return err
		//return s.rollback(ctx, blockNumber, "updateForkIDIfNeeded", err, dbTx)
	}
	if isUpdate {
		return nil // The calling function is doing the commit
	}

	if !isIncommingForkIdGreatestThanLastOne(fID, fIds) {
		err = fmt.Errorf("ForkID: %d, received don't fit sequence, last forkid:%d ", forkID.ForkID, lastForkID(fIds))
		log.Error(err.Error())
		return err
	}
	//If the forkID.batchnumber is a future batch
	latestBatchNumber, err := s.state.GetLastBatchNumber(ctx, dbTx)
	if err != nil && !errors.Is(err, state.ErrStateNotSynchronized) {
		log.Errorf("error getting last batch number. Error: %v", err)
		//return s.rollback(ctx, blockNumber, "getting last batch number", err, dbTx)
	}
	// Add new forkID to the state
	err = s.state.AddForkIDInterval(ctx, fID, dbTx)
	if err != nil {
		return s.rollback(ctx, blockNumber, "adding new forkID interval to the state", err, dbTx)
	}
	if latestBatchNumber <= forkID.BatchNumber || s.sync.IsTrustedSequencer() { //If the forkID will start in a future batch or isTrustedSequencer
		log.Infof("Just adding forkID. Skipping reset forkID. ForkID: %+v.", fID)
		return nil // The calling function is doing the commit
	}

	log.Info("ForkID received in the permissionless node that affects to a batch from the past")
	//Reset DB only if permissionless node
	log.Debugf("ForkID: %d, Reverting synchronization to batch: %d", forkID.ForkID, forkID.BatchNumber+1)
	err = s.state.ResetForkID(ctx, forkID.BatchNumber+1, dbTx)
	if err != nil {
		// I'm wondering if it's really needed to call rollback, in fact, the caller is going to call the rollback when returns an error
		return s.rollback(ctx, blockNumber, fmt.Sprintf("resetting the state to %d", forkID.BatchNumber+1), err, dbTx)
	}

	// Commit because it returns an error to force the resync
	err = s.commit(ctx, blockNumber, dbTx)
	if err != nil {
		return err
	}

	return fmt.Errorf("new ForkID detected, reseting synchronizarion")
}
func (s *ProcessorForkId) rollback(ctx context.Context, blockNumber uint64, msg string, err error, dbTx pgx.Tx) error {
	log.Error("error %s. Error: ", msg, err)
	rollbackErr := dbTx.Rollback(ctx)
	if rollbackErr != nil {
		log.Errorf("error rolling back %s. BlockNumber: %d, rollbackErr: %s, error : %v", msg, blockNumber, rollbackErr.Error(), err)
		return rollbackErr
	}
	return err
}

func (s *ProcessorForkId) commit(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) error {
	err := dbTx.Commit(ctx)
	if err != nil {
		log.Error("error committing forkId. Error: ", err)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}
	return nil
}
