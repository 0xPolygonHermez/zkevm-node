package aggregator

import "math/big"

type SyncEvent struct {
	BatchNum uint64
	LastStateRoot *big.Int
	// ??? which type ZKI will be?
	ZKI *big.Int
}

type SynchronizerClient struct {
	SyncEventChan chan SyncEvent
}

func NewSynchronizerClient() SynchronizerClient {
	syncEventChan := make(chan SyncEvent)
	return SynchronizerClient{SyncEventChan: syncEventChan}
}
