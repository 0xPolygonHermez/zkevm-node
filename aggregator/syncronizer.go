package aggregator

import "math/big"

type SyncEvent struct {
	LastStateRoot *big.Int
	// ??? which type proof will be?
	Proof *big.Int
}

type SynchronizerClient struct {
	SyncEventChan chan SyncEvent
}

func NewSynchronizerClient() SynchronizerClient {
	syncEventChan := make(chan SyncEvent)
	return SynchronizerClient{SyncEventChan: syncEventChan}
}
