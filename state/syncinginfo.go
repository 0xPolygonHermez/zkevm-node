package state

// SyncingInfo stores information regarding the syncing status of the node
type SyncingInfo struct {
	InitialSyncingBatch         uint64
	LastBatchNumberSeen         uint64
	LastBatchNumberConsolidated uint64
	CurrentBatchNumber          uint64
}
