package actions

import (
	"reflect"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
)

// ProcessorBase is the base struct for all the processors, if reduces the boilerplate
// implementing the Name, SupportedEvents and SupportedForkIds functions
type ProcessorBase[T any] struct {
	SupportedEvent    []etherman.EventOrder
	SupportedForkdIds *[]ForkIdType
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
	return p.SupportedEvent
}

// SupportedForkIds returns the supported forkIds in the struct or the dafault till incaberry forkId
func (p *ProcessorBase[T]) SupportedForkIds() []ForkIdType {
	if p.SupportedForkdIds != nil {
		return *p.SupportedForkdIds
	}
	// returns none
	return []ForkIdType{}
}
