package etrog

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

// stateProcessorL1GlobalExitRootInterface interface required from state
type stateProcessorL1GlobalExitRootInterface interface {
	AddGlobalExitRoot(ctx context.Context, exitRoot *state.GlobalExitRoot, dbTx pgx.Tx) error
}

// ProcessorL1GlobalExitRoot implements L1EventProcessor for GlobalExitRootsOrder
type ProcessorL1GlobalExitRoot struct {
	actions.ProcessorBase[ProcessorL1GlobalExitRoot]
	state stateProcessorL1GlobalExitRootInterface
}

// NewProcessorL1GlobalExitRoot new processor for GlobalExitRootsOrder
func NewProcessorL1GlobalExitRoot(state stateProcessorL1GlobalExitRootInterface) *ProcessorL1GlobalExitRoot {
	return &ProcessorL1GlobalExitRoot{
		ProcessorBase: actions.ProcessorBase[ProcessorL1GlobalExitRoot]{
			SupportedEvent:    etherman.L1InfoTreeOrder,
			SupportedForkdIds: &ForksIdOnlyEtrog},
		state: state}
}

// Process process event
func (p *ProcessorL1GlobalExitRoot) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	l1InfoTree := l1Block.L1InfoTree[order.Pos]
	ger := state.GlobalExitRoot{
		BlockNumber:     l1InfoTree.BlockNumber,
		MainnetExitRoot: l1InfoTree.MainnetExitRoot,
		RollupExitRoot:  l1InfoTree.RollupExitRoot,
		GlobalExitRoot:  l1InfoTree.GlobalExitRoot.GlobalExitRoot,
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
