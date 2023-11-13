package synchronizer

import (
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	l1events "github.com/0xPolygonHermez/zkevm-node/synchronizer/l1events"
)

func defaultsL1EventProcessors(state stateInterface) *l1events.L1EventProcessors {
	p := l1events.NewL1EventProcessors()
	p.Set(l1events.DefaultForkId, etherman.GlobalExitRootsOrder, &l1events.GlobalExitRootLegacy{})
	return p
}
