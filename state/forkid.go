package state

// ForkIDInterval is a fork id interval
type ForkIDInterval struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
	ForkId          uint64
}

// GetForkIDByBatchNumber returns the fork id for a given batch number
func GetForkIDByBatchNumber(intervals []ForkIDInterval, batchNumber uint64) uint64 {
	for _, interval := range intervals {
		if batchNumber >= interval.FromBatchNumber && batchNumber <= interval.ToBatchNumber {
			return interval.ForkId
		}
	}
	return 1
}
