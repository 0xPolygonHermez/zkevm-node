package state

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/jackc/pgx/v4"
)

// ForkIDInterval is a fork id interval
type ForkIDInterval struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
	ForkId          uint64
	Version         string
	BlockNumber     uint64
}

// UpdateForkIDIntervals updates the forkID intervals
func (s *State) UpdateForkIDIntervals(intervals []ForkIDInterval) {
	log.Infof("Updating forkIDs. Setting %d forkIDs", len(intervals))
	log.Infof("intervals: %#v", intervals)
	s.cfg.ForkIDIntervals = intervals
}

// AddForkIDInterval updates the forkID intervals
func (s *State) AddForkIDInterval(ctx context.Context, newForkID ForkIDInterval, dbTx pgx.Tx) error {
	// Add forkId to db and memori variable
	oldForkIDs, err := s.GetForkIDs(ctx, dbTx)
	if err != nil {
		log.Error("error getting oldForkIDs. Error: ", err)
		return err
	}
	if len(oldForkIDs) == 0 {
		s.UpdateForkIDIntervals([]ForkIDInterval{newForkID})
	} else {
		var forkIDs []ForkIDInterval
		forkIDs = oldForkIDs
		forkIDs[len(forkIDs)-1].ToBatchNumber = newForkID.FromBatchNumber - 1
		err := s.UpdateForkID(ctx, forkIDs[len(forkIDs)-1], dbTx)
		if err != nil {
			log.Errorf("error updating forkID: %d. Error: %v", forkIDs[len(forkIDs)-1].ForkId, err)
			return err
		}
		forkIDs = append(forkIDs, newForkID)

		s.UpdateForkIDIntervals(forkIDs)
	}
	err = s.AddForkID(ctx, newForkID, dbTx)
	if err != nil {
		log.Errorf("error adding forkID %d. Error: %v", newForkID.ForkId, err)
		return err
	}
	return nil
}

// GetForkIDByBatchNumber returns the fork id for a given batch number
func (s *State) GetForkIDByBatchNumber(batchNumber uint64) uint64 {
	// If NumBatchForkIdUpgrade is defined (!=0) we are performing forkid upgrade process
	// In this case, if the batchNumber is the next to the NumBatchForkIdUpgrade, we need to return the
	// new "future" forkId (ForkUpgradeNewForkId)
	if (s.cfg.ForkUpgradeBatchNumber) != 0 && (batchNumber > s.cfg.ForkUpgradeBatchNumber) {
		return s.cfg.ForkUpgradeNewForkId
	}

	for _, interval := range s.cfg.ForkIDIntervals {
		if batchNumber >= interval.FromBatchNumber && batchNumber <= interval.ToBatchNumber {
			return interval.ForkId
		}
	}

	// If not found return the last fork id
	return s.cfg.ForkIDIntervals[len(s.cfg.ForkIDIntervals)-1].ForkId
}
