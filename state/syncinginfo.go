package state

// SyncingInfo stores information regarding the syncing status of the node
type SyncingInfo struct {
	LastBatchNumberSeen         uint64 `json:"lastBatchNumberSeen"`
	LastBatchNumberConsolidated uint64 `json:"lastBatchNumberConsolidated"`
	InitialSyncingBatch         uint64 `json:"initialSyncingBatch"`
}
