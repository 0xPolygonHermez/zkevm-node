package synchronizer

import (
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	l1events "github.com/0xPolygonHermez/zkevm-node/synchronizer/l1events"
)

func defaultsL1EventProcessors(sync *ClientSynchronizer) *l1events.L1EventProcessors {
	p := l1events.L1EventProcessorsBuilder{}
	p.Add(etherman.GlobalExitRootsOrder,
		l1events.NewProcessorGlobalExitRootLegacy(sync.state))
	p.Add(etherman.GlobalExitRootsOrder,
		l1events.NewProcessorSequenceBatchesLegacy(sync.state, sync.etherMan, sync.pool, sync.eventLog, sync))
	return p.Build()
}
