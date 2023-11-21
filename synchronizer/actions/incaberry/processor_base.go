package incaberry

import (
	"reflect"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
)

// ProcessorBase is the base struct for all the processors, if reduces the boilerplate
// implementing the Name, SupportedEvents and SupportedForkIds functions
type ProcessorBase[T any] struct {
	supportedEvent    etherman.EventOrder
	supportedForkdIds *[]actions.ForkIdType
}

// Name returns the name of the struct T
func (g *ProcessorBase[T]) Name() string {
	var value T
	a := reflect.TypeOf(value)
	b := a.Name()
	return b
}

// SupportedEvents returns the supported events in the struct
func (p *ProcessorBase[T]) SupportedEvents() []etherman.EventOrder {
	return []etherman.EventOrder{p.supportedEvent}
}

// SupportedForkIds returns the supported forkIds in the struct or the dafault till incaberry forkId
func (p *ProcessorBase[T]) SupportedForkIds() []actions.ForkIdType {
	if p.supportedForkdIds != nil {
		return *p.supportedForkdIds
	}
	// returns default forkIds till incaberry forkId
	return []actions.ForkIdType{1, 2, 3, 4, 5, 6}
}
