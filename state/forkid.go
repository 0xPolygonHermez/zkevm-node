package state

import "github.com/0xPolygonHermez/zkevm-node/log"

// ForkIDInterval is a fork id interval
type ForkIDInterval struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
	ForkId          uint64
	Version         string
}

// UpdateForkIDIntervals updates the forkID intervals
func (s *State) UpdateForkIDIntervals(intervals []ForkIDInterval) {
	log.Infof("Updating forkIDs. Setting %d forkIDs", len(intervals))
	log.Infof("intervals: %#v", intervals)
	s.cfg.ForkIDIntervals = intervals
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
