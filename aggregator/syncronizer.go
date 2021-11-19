package aggregator

import "math/big"

type SyncEvent struct {
	BatchNum  uint64
	StateRoot *big.Int
}

type VirtualBatchEvent struct {
	BatchNum uint64
}

type ConsolidatedBatchEvent struct {
	BatchNum uint64
}

type SynchronizerClient struct {
	VirtualBatchEventChan      chan VirtualBatchEvent
	ConsolidatedBatchEventChan chan ConsolidatedBatchEvent
}

func NewSynchronizerClient() SynchronizerClient {
	virtualBatchEventChan := make(chan VirtualBatchEvent)
	consolidatedBatchEventChan := make(chan ConsolidatedBatchEvent)

	return SynchronizerClient{
		VirtualBatchEventChan:      virtualBatchEventChan,
		ConsolidatedBatchEventChan: consolidatedBatchEventChan,
	}
}
