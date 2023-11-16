package l1events

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

type StateProcessorGlobalExitRootInterface interface {
	AddGlobalExitRoot(ctx context.Context, exitRoot *state.GlobalExitRoot, dbTx pgx.Tx) error
}

// GlobalExitRootLegacy implements L1EventProcessor
type ProcessorGlobalExitRoot struct {
	ProcessorBase[ProcessorGlobalExitRoot]
	state StateProcessorGlobalExitRootInterface
}

func NewProcessorGlobalExitRoot(state StateProcessorGlobalExitRootInterface) *ProcessorGlobalExitRoot {
	return &ProcessorGlobalExitRoot{
		ProcessorBase: ProcessorBase[ProcessorGlobalExitRoot]{supportedEvent: etherman.GlobalExitRootsOrder},
		state:         state}
}

func (p *ProcessorGlobalExitRoot) Process(ctx context.Context, event etherman.EventOrder, l1Block *etherman.Block, postion int, dbTx pgx.Tx) error {
	globalExitRoot := l1Block.GlobalExitRoots[postion]
	ger := state.GlobalExitRoot{
		BlockNumber:     globalExitRoot.BlockNumber,
		MainnetExitRoot: globalExitRoot.MainnetExitRoot,
		RollupExitRoot:  globalExitRoot.RollupExitRoot,
		GlobalExitRoot:  globalExitRoot.GlobalExitRoot,
	}
	err := p.state.AddGlobalExitRoot(ctx, &ger, dbTx)
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
