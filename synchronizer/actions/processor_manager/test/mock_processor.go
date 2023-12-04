package processor_manager_test

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

type ProcessorStub struct {
	name             string
	supportedEvents  []etherman.EventOrder
	supportedForkIds []actions.ForkIdType
	responseProcess  error
}

func (p *ProcessorStub) Name() string {
	return p.name
}

func (p *ProcessorStub) SupportedEvents() []etherman.EventOrder {
	return p.supportedEvents
}

func (p *ProcessorStub) SupportedForkIds() []actions.ForkIdType {
	return p.supportedForkIds
}

func (p *ProcessorStub) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	return p.responseProcess
}
