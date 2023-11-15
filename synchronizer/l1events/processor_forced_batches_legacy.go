package l1events

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

type StateProcessForcedBatchesLegacyInterface interface {
	AddForcedBatch(ctx context.Context, forcedBatch *state.ForcedBatch, dbTx pgx.Tx) error
}

// GlobalExitRootLegacy implements L1EventProcessor
type ProcessForcedBatchesLegacy struct {
	state StateProcessForcedBatchesLegacyInterface
}

func NewProcessForcedBatchesLegacy(state StateProcessForcedBatchesLegacyInterface) *ProcessForcedBatchesLegacy {
	return &ProcessForcedBatchesLegacy{state: state}
}

func (g *ProcessForcedBatchesLegacy) String() string {
	return "ProcessForcedBatchesLegacy"
}

func (p *ProcessForcedBatchesLegacy) SupportedForkIds() []ForkIdType {
	return []ForkIdType{1, 2, 3, 4, 5, 6}
}

func (p *ProcessForcedBatchesLegacy) Process(ctx context.Context, event etherman.EventOrder, l1Block *etherman.Block, postion int, dbTx pgx.Tx) error {
	err := p.processForcedBatch(ctx, l1Block.ForcedBatches[postion], dbTx)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProcessForcedBatchesLegacy) processForcedBatch(ctx context.Context, forcedBatch etherman.ForcedBatch, dbTx pgx.Tx) error {
	// Store forced batch into the db
	forcedB := state.ForcedBatch{
		BlockNumber:       forcedBatch.BlockNumber,
		ForcedBatchNumber: forcedBatch.ForcedBatchNumber,
		Sequencer:         forcedBatch.Sequencer,
		GlobalExitRoot:    forcedBatch.GlobalExitRoot,
		RawTxsData:        forcedBatch.RawTxsData,
		ForcedAt:          forcedBatch.ForcedAt,
	}
	log.Infof("processForcedBatch: Storing forcedBatch. BatchNumber: %d  BlockNumber: %d", forcedBatch.ForcedBatchNumber, forcedBatch.BlockNumber)
	err := p.state.AddForcedBatch(ctx, &forcedB, dbTx)
	if err != nil {
		log.Errorf("error storing the forcedBatch in processForcedBatch. BlockNumber: %d", forcedBatch.BlockNumber)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", forcedBatch.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error storing the forcedBatch in processForcedBatch. BlockNumber: %d, error: %v", forcedBatch.BlockNumber, err)
		return err
	}
	return nil
}
