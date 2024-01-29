package state

// SyncingInfo stores information regarding the syncing status of the node
type SyncingInfo struct {
	InitialSyncingBlock uint64 // L2Block corresponding to InitialSyncingBatch
	CurrentBlockNumber  uint64 // last L2Block in state

	InitialSyncingBatch         uint64
	LastBatchNumberSeen         uint64
	LastBatchNumberConsolidated uint64
	CurrentBatchNumber          uint64
	// IsSynchronizing indicates if the node is syncing (true -> syncing, false -> fully synced)
	IsSynchronizing bool
}
