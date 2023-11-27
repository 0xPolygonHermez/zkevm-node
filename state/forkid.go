package state

import (
	"context"
	"fmt"
	"sort"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/jackc/pgx/v4"
)

const (
	// FORKID_BLUEBERRY is the fork id 4
	FORKID_BLUEBERRY = 4
	// FORKID_DRAGONFRUIT is the fork id 5
	FORKID_DRAGONFRUIT = 5
	// FORKID_INCABERRY is the fork id 6
	FORKID_INCABERRY = 6
	// FORKID_ETROG is the fork id 7
	FORKID_ETROG = 7
)

// ForkIDInterval is a fork id interval
type ForkIDInterval struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
	ForkId          uint64
	Version         string
	BlockNumber     uint64
}

// UpdateForkIDIntervalsInMemory updates the forkID intervals in memory
func (s *State) UpdateForkIDIntervalsInMemory(intervals []ForkIDInterval) {
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
		s.UpdateForkIDIntervalsInMemory([]ForkIDInterval{newForkID})
	} else {
		var forkIDs []ForkIDInterval
		forkIDs = oldForkIDs
		// Check to detect forkID inconsistencies
		if forkIDs[len(forkIDs)-1].ForkId+1 != newForkID.ForkId {
			log.Errorf("error checking forkID sequence. Last ForkID stored: %d. New ForkID received: %d", forkIDs[len(forkIDs)-1].ForkId, newForkID.ForkId)
			return fmt.Errorf("error checking forkID sequence. Last ForkID stored: %d. New ForkID received: %d", forkIDs[len(forkIDs)-1].ForkId, newForkID.ForkId)
		}
		forkIDs[len(forkIDs)-1].ToBatchNumber = newForkID.FromBatchNumber - 1
		err := s.UpdateForkID(ctx, forkIDs[len(forkIDs)-1], dbTx)
		if err != nil {
			log.Errorf("error updating forkID: %d. Error: %v", forkIDs[len(forkIDs)-1].ForkId, err)
			return err
		}
		forkIDs = append(forkIDs, newForkID)

		s.UpdateForkIDIntervalsInMemory(forkIDs)
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
	return s.storage.GetForkIDByBatchNumber(batchNumber)
}

// GetForkIDByBlockNumber returns the fork id for a given block number
func (s *State) GetForkIDByBlockNumber(blockNumber uint64) uint64 {
	for _, index := range sortIndexForForkdIDSortedByBlockNumber(s.cfg.ForkIDIntervals) {
		// reverse travesal
		interval := s.cfg.ForkIDIntervals[len(s.cfg.ForkIDIntervals)-1-index]
		if blockNumber > interval.BlockNumber {
			return interval.ForkId
		}
	}
	// If not found return the  fork id 1
	return 1
}

func sortIndexForForkdIDSortedByBlockNumber(forkIDs []ForkIDInterval) []int {
	sortedIndex := make([]int, len(forkIDs))
	for i := range sortedIndex {
		sortedIndex[i] = i
	}
	cmpFunc := func(i, j int) bool {
		return forkIDs[sortedIndex[i]].BlockNumber < forkIDs[sortedIndex[j]].BlockNumber
	}
	sort.Slice(sortedIndex, cmpFunc)
	return sortedIndex
}
