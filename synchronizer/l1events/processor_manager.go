package synchronizer_l1_events

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/jackc/pgx/v4"
)

const (
	defaultForkId forkIdType = 0
)

var (
	// ErrNotFound is used when the object is not found
	ErrNotFound = errors.New("not found")
)

type forkIdType uint64

type L1EventProcessors struct {
	// forkId -> event -> processor
	processors map[forkIdType]map[etherman.EventOrder]L1EventProcessor
}

func NewL1EventProcessors() *L1EventProcessors {
	return &L1EventProcessors{
		processors: make(map[forkIdType]map[etherman.EventOrder]L1EventProcessor),
	}
}

func (p *L1EventProcessors) Set(forkId forkIdType, event etherman.EventOrder, processor L1EventProcessor) {
	if _, ok := p.processors[forkId]; !ok {
		p.processors[forkId] = make(map[etherman.EventOrder]L1EventProcessor)
	}
	p.processors[forkId][event] = processor
}

func (p *L1EventProcessors) Get(forkId forkIdType, event etherman.EventOrder) L1EventProcessor {
	if _, ok := p.processors[forkId]; !ok {
		return p.Get(defaultForkId, event)
	}
	return p.processors[forkId][event]
}

func (p *L1EventProcessors) Process(ctx context.Context, forkId forkIdType, event etherman.EventOrder, block *etherman.Block, position int, dbTx pgx.Tx) error {
	processor := p.Get(forkId, event)
	if processor == nil {
		return ErrNotFound
	}
	return processor.Process(ctx, event, block, position, dbTx)
}
