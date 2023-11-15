package l1events

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

type StateGlobalExitRootInterface interface {
	AddGlobalExitRoot(ctx context.Context, exitRoot *state.GlobalExitRoot, dbTx pgx.Tx) error
}

// GlobalExitRootLegacy implements L1EventProcessor
type GlobalExitRootLegacy struct {
	state StateGlobalExitRootInterface
}

func NewProcessorGlobalExitRootLegacy(state StateGlobalExitRootInterface) *GlobalExitRootLegacy {
	return &GlobalExitRootLegacy{state: state}
}

func (g *GlobalExitRootLegacy) String() string {
	return "GlobalExitRootLegacy"
}

func (g *GlobalExitRootLegacy) SupportedForkIds() []ForkIdType {
	return []ForkIdType{1, 2, 3, 4, 5, 6}
}

func (g *GlobalExitRootLegacy) Process(ctx context.Context, event etherman.EventOrder, l1Block *etherman.Block, postion int, dbTx pgx.Tx) error {
	globalExitRoot := l1Block.GlobalExitRoots[postion]
	ger := state.GlobalExitRoot{
		BlockNumber:     globalExitRoot.BlockNumber,
		MainnetExitRoot: globalExitRoot.MainnetExitRoot,
		RollupExitRoot:  globalExitRoot.RollupExitRoot,
		GlobalExitRoot:  globalExitRoot.GlobalExitRoot,
	}
	err := g.state.AddGlobalExitRoot(ctx, &ger, dbTx)
	if err != nil {
		log.Errorf("error storing the GlobalExitRoot in processGlobalExitRoot. BlockNumber: %d, error: %v", l1Block.BlockNumber, err)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", l1Block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error storing the GlobalExitRoot in processGlobalExitRoot. BlockNumber: %d, error: %v", l1Block.BlockNumber, err)
		return err
	}
	return nil
}
