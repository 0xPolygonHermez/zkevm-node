package syncinterfaces

type EthermanGetLatestBatchNumber interface {
	GetLatestBatchNumber() (uint64, error)
}
