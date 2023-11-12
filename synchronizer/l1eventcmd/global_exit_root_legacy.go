package synchronizer_l1_events

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

type stateGlobalExitRootWriter interface {
	AddGlobalExitRoot(ctx context.Context, exitRoot *state.GlobalExitRoot, dbTx pgx.Tx) error
}

// GlobalExitRootLegacyCmd implements L1EventExecutor
type GlobalExitRootLegacyCmd struct {
	state stateGlobalExitRootWriter
}

func (g *GlobalExitRootLegacyCmd) Execute(ctx context.Context, l1Block *etherman.Block, postion int, dbTx pgx.Tx) error {
	globalExitRoot := l1Block.GlobalExitRoots[postion]
	ger := state.GlobalExitRoot{
		BlockNumber:     globalExitRoot.BlockNumber,
		MainnetExitRoot: globalExitRoot.MainnetExitRoot,
		RollupExitRoot:  globalExitRoot.RollupExitRoot,
		GlobalExitRoot:  globalExitRoot.GlobalExitRoot,
	}
	err := g.state.AddGlobalExitRoot(ctx, &ger, dbTx)
	if err != nil {
		log.Errorf("error storing the globalExitRoot in processGlobalExitRoot. BlockNumber: %d", globalExitRoot.BlockNumber)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", globalExitRoot.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error storing the GlobalExitRoot in processGlobalExitRoot. BlockNumber: %d, error: %v", globalExitRoot.BlockNumber, err)
		return err
	}
	return nil
}
