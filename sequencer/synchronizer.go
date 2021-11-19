package sequencer

import (
	"github.com/ethereum/go-ethereum/common"
)

type SyncEvent struct {
	StartingHash common.Hash
	BatchNum     uint64
}

type SynchronizerClient struct {
	SyncEventChan chan SyncEvent
}

func NewSynchronizerClient() SynchronizerClient {
	// connect to synch to get signal, that it's synced
	syncEventChan := make(chan SyncEvent)
	return SynchronizerClient{SyncEventChan: syncEventChan}
}
