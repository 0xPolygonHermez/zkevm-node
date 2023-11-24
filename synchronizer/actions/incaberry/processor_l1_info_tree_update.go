package incaberry

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

// stateProcessorL1InfoTreeInterface interface required from state
type stateProcessorL1InfoTreeInterface interface {
	AddGlobalExitRoot(ctx context.Context, exitRoot *state.GlobalExitRoot, dbTx pgx.Tx) error
}

// ProcessorL1InfoTreeUpdate implements L1EventProcessor for GlobalExitRootsOrder
type ProcessorL1InfoTreeUpdate struct {
	actions.ProcessorBase[ProcessorL1InfoTreeUpdate]
	state stateProcessorL1InfoTreeInterface
}

// NewProcessorL1InfoTreeUpdate new processor for GlobalExitRootsOrder
func NewProcessorL1InfoTreeUpdate(state stateProcessorL1InfoTreeInterface) *ProcessorL1InfoTreeUpdate {
	return &ProcessorL1InfoTreeUpdate{
		ProcessorBase: actions.ProcessorBase[ProcessorL1InfoTreeUpdate]{
			SupportedEvent:    []etherman.EventOrder{etherman.L1InfoTreeOrder},
			SupportedForkdIds: &actions.ForksIdToIncaberry},
		state: state}
}

// Process process event
func (p *ProcessorL1InfoTreeUpdate) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	l1InfoTree := l1Block.L1InfoTree[order.Pos]
	ger := state.GlobalExitRoot{
		BlockNumber:     l1InfoTree.BlockNumber,
		MainnetExitRoot: l1InfoTree.MainnetExitRoot,
		RollupExitRoot:  l1InfoTree.RollupExitRoot,
		GlobalExitRoot:  l1InfoTree.GlobalExitRoot.GlobalExitRoot,
		Timestamp:       l1InfoTree.MinTimestamp,
	}

	err := p.state.AddGlobalExitRoot(ctx, &ger, dbTx)
	if err != nil {
		log.Errorf("error storing the L1InfoTreeOrder(incaberry). BlockNumber: %d, error: %v", l1Block.BlockNumber, err)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", l1Block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error storing the L1InfoTreeOrder(incaberry). BlockNumber: %d, error: %v", l1Block.BlockNumber, err)
		return err
	}
	return nil
}
