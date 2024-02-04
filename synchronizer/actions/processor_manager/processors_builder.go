package processor_manager

import (
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	// "github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
)

// L1EventProcessorsBuilder is a builder for L1EventProcessors
// how to use:
//
//		p := L1EventProcessorsBuilder{}
//		p.Add(etherman.GlobalExitRootsOrder, l1events.NewGlobalExitRootLegacy(state))
//	 	p.Set....
//		return p.Build()
type L1EventProcessorsBuilder struct {
	result *L1EventProcessors
}

// NewL1EventProcessorsBuilder returns a new L1EventProcessorsBuilder instance
func NewL1EventProcessorsBuilder() L1EventProcessorsBuilder {
	return L1EventProcessorsBuilder{}
}

// Build return the L1EventProcessors builded
func (p *L1EventProcessorsBuilder) Build() *L1EventProcessors {
	return p.result
}

// Register register a L1EventProcessor. It ask to the processor the supported forkId and events
// if there are a previous object register it will panic
func (p *L1EventProcessorsBuilder) Register(processor actions.L1EventProcessor) {
	p.createResultIfNeeded()
	for _, forkID := range processor.SupportedForkIds() {
		for _, event := range processor.SupportedEvents() {
			p.Set(forkID, event, processor, true)
		}
	}
}

// Set add a L1EventProcessor. If param panicIfExists is true, will panic if already exists the object
//
//	the only use to panicIfExists=false is to override a processor in a unitttest
func (p *L1EventProcessorsBuilder) Set(forkID actions.ForkIdType, event etherman.EventOrder, processor actions.L1EventProcessor, panicIfExists bool) {
	p.createResultIfNeeded()
	if _, ok := p.result.processors[forkID]; !ok {
		p.result.processors[forkID] = make(map[etherman.EventOrder]actions.L1EventProcessor)
	}
	if _, ok := p.result.processors[forkID][event]; ok && panicIfExists {
		panic("processor already set")
	}
	p.result.processors[forkID][event] = processor
}

func (p *L1EventProcessorsBuilder) createResultIfNeeded() {
	if p.result == nil {
		p.result = NewL1EventProcessors()
	}
}
