package state

// ForkIDInterval is a fork id interval
type ForkIDInterval struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
	ForkId          uint64
	Version         string
}

// GetForkIDByBatchNumber returns the fork id for a given batch number
func GetForkIDByBatchNumber(intervals []ForkIDInterval, batchNumber uint64) uint64 {
	for _, interval := range intervals {
		if batchNumber >= interval.FromBatchNumber && batchNumber <= interval.ToBatchNumber {
			return interval.ForkId
		}
	}

	// If not found return the last fork id
	return intervals[len(intervals)-1].ForkId
}
