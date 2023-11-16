package l1events

import (
	"reflect"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
)

type ProcessorBase[T any] struct {
	supportedEvent etherman.EventOrder
}

func (g *ProcessorBase[T]) Name() string {
	var value T
	a := reflect.TypeOf(value)
	b := a.Name()
	return b
}

func (p *ProcessorBase[T]) SupportedEvents() []etherman.EventOrder {
	return []etherman.EventOrder{p.supportedEvent}
}

func (p *ProcessorBase[T]) SupportedForkIds() []ForkIdType {
	return []ForkIdType{1, 2, 3, 4, 5, 6}
}
