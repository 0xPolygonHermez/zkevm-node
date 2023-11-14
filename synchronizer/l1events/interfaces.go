package l1events

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/jackc/pgx/v4"
)

type L1EventProcessor interface {
	String() string
	SupportedForkIds() []forkIdType
	Process(ctx context.Context, event etherman.EventOrder, l1Block *etherman.Block, postion int, dbTx pgx.Tx) error
}
