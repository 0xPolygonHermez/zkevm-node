package synchronizer

import (
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	l1events "github.com/0xPolygonHermez/zkevm-node/synchronizer/l1events"
)

func defaultsL1EventProcessors(sync *ClientSynchronizer) *l1events.L1EventProcessors {
	p := l1events.L1EventProcessorsBuilder{}
	p.Add(etherman.GlobalExitRootsOrder,
		l1events.NewProcessorGlobalExitRoot(sync.state))
	p.Add(etherman.SequenceBatchesOrder,
		l1events.NewProcessorSequenceBatches(sync.state, sync.etherMan, sync.pool, sync.eventLog, sync))
	p.Add(etherman.ForcedBatchesOrder,
		l1events.NewProcessForcedBatches(sync.state))
	p.Add(etherman.SequenceForceBatchesOrder,
		l1events.NewProcessSequenceForcedBatches(sync.state, sync))
	p.Add(etherman.TrustedVerifyBatchOrder,
		l1events.NewProcessorTrustedVerifyBatch(sync.state))
	p.Add(etherman.ForkIDsOrder,
		l1events.NewProcessorForkId(sync.state, sync))
	return p.Build()
}
