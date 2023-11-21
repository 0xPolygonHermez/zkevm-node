package actions

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/jackc/pgx/v4"
)

// ForkIdType is the type of the forkId
type ForkIdType uint64

const (
	// WildcardForkId It match for all forkIds
	WildcardForkId ForkIdType = 0
)

// L1EventProcessor is the interface for a processor of L1 events
// The main function is Process that must execute the event
type L1EventProcessor interface {
	// Name of the processor
	Name() string
	// SupportedForkIds list of forkId that support (you could use WildcardForkId)
	SupportedForkIds() []ForkIdType
	// SupportedEvents list of events that support (typically one)
	SupportedEvents() []etherman.EventOrder
	// Process a incomming event
	Process(ctx context.Context, event etherman.EventOrder, l1Block *etherman.Block, position int, dbTx pgx.Tx) error
}
