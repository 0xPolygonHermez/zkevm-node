package l1events

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/jackc/pgx/v4"
)

const (
	DefaultForkId ForkIdType = 0
)

var (
	// ErrNotFound is used when the object is not found
	ErrNotFound           = errors.New("not found")
	ErrForkIdNotSupported = errors.New("forkId not supported")
)

type ForkIdType uint64

type L1EventProcessors struct {
	// forkId -> event -> processor
	processors map[ForkIdType]map[etherman.EventOrder]L1EventProcessor
}

func NewL1EventProcessors() *L1EventProcessors {
	return &L1EventProcessors{
		processors: make(map[ForkIdType]map[etherman.EventOrder]L1EventProcessor),
	}
}

// Get returns the processor, if not found returns nil
func (p *L1EventProcessors) Get(forkId ForkIdType, event etherman.EventOrder) L1EventProcessor {
	if _, ok := p.processors[forkId]; !ok {
		return p.Get(DefaultForkId, event)
	}
	return p.processors[forkId][event]
}

func (p *L1EventProcessors) Process(ctx context.Context, forkId ForkIdType, event etherman.EventOrder, block *etherman.Block, position int, dbTx pgx.Tx) error {
	processor := p.Get(forkId, event)
	if processor == nil {
		return ErrNotFound
	}
	return processor.Process(ctx, event, block, position, dbTx)
}
