package l1events

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/jackc/pgx/v4"
)

type L1EventProcessor interface {
	// Name of the processor
	Name() string
	// List of forkId that support
	SupportedForkIds() []ForkIdType
	// List of events that support (typically one)
	SupportedEvents() []etherman.EventOrder
	// Process a incomming event
	Process(ctx context.Context, event etherman.EventOrder, l1Block *etherman.Block, postion int, dbTx pgx.Tx) error
}
