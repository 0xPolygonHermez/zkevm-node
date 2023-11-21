package incaberry

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

type stateL1TrustedVerifyBatchInterface interface {
	GetLastVerifiedBatch(ctx context.Context, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	AddVerifiedBatch(ctx context.Context, verifiedBatch *state.VerifiedBatch, dbTx pgx.Tx) error
}

// ProcessorL1TrustedVerifyBatch implements L1EventProcessor
type ProcessorL1TrustedVerifyBatch struct {
	ProcessorBase[ProcessorL1TrustedVerifyBatch]
	state stateL1TrustedVerifyBatchInterface
}

// NewProcessorL1TrustedVerifyBatch returns instance of a processor for TrustedVerifyBatchOrder
func NewProcessorL1TrustedVerifyBatch(state stateL1TrustedVerifyBatchInterface) *ProcessorL1TrustedVerifyBatch {
	return &ProcessorL1TrustedVerifyBatch{
		ProcessorBase: ProcessorBase[ProcessorL1TrustedVerifyBatch]{supportedEvent: etherman.TrustedVerifyBatchOrder},
		state:         state}
}

// Process process event
func (p *ProcessorL1TrustedVerifyBatch) Process(ctx context.Context, event etherman.EventOrder, l1Block *etherman.Block, position int, dbTx pgx.Tx) error {
	return p.processTrustedVerifyBatches(ctx, l1Block.VerifiedBatches[position], dbTx)
}

func (p *ProcessorL1TrustedVerifyBatch) processTrustedVerifyBatches(ctx context.Context, lastVerifiedBatch etherman.VerifiedBatch, dbTx pgx.Tx) error {
	lastVBatch, err := p.state.GetLastVerifiedBatch(ctx, dbTx)
	if err != nil {
		log.Errorf("error getting lastVerifiedBatch stored in db in processTrustedVerifyBatches. Processing synced blockNumber: %d", lastVerifiedBatch.BlockNumber)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. Processing synced blockNumber: %d, rollbackErr: %s, error : %v", lastVerifiedBatch.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting lastVerifiedBatch stored in db in processTrustedVerifyBatches. Processing synced blockNumber: %d, error: %v", lastVerifiedBatch.BlockNumber, err)
		return err
	}
	nbatches := lastVerifiedBatch.BatchNumber - lastVBatch.BatchNumber
	batch, err := p.state.GetBatchByNumber(ctx, lastVerifiedBatch.BatchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting GetBatchByNumber stored in db in processTrustedVerifyBatches. Processing batchNumber: %d", lastVerifiedBatch.BatchNumber)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. Processing batchNumber: %d, rollbackErr: %s, error : %v", lastVerifiedBatch.BatchNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting GetBatchByNumber stored in db in processTrustedVerifyBatches. Processing batchNumber: %d, error: %v", lastVerifiedBatch.BatchNumber, err)
		return err
	}

	// Checks that calculated state root matches with the verified state root in the smc
	if batch.StateRoot != lastVerifiedBatch.StateRoot {
		log.Warn("nbatches: ", nbatches)
		log.Warnf("Batch from db: %+v", batch)
		log.Warnf("Verified Batch: %+v", lastVerifiedBatch)
		log.Errorf("error: stateRoot calculated and state root verified don't match in processTrustedVerifyBatches. Processing batchNumber: %d", lastVerifiedBatch.BatchNumber)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. Processing batchNumber: %d, rollbackErr: %v", lastVerifiedBatch.BatchNumber, rollbackErr)
			return rollbackErr
		}
		log.Errorf("error: stateRoot calculated and state root verified don't match in processTrustedVerifyBatches. Processing batchNumber: %d", lastVerifiedBatch.BatchNumber)
		return fmt.Errorf("error: stateRoot calculated and state root verified don't match in processTrustedVerifyBatches. Processing batchNumber: %d", lastVerifiedBatch.BatchNumber)
	}
	var i uint64
	for i = 1; i <= nbatches; i++ {
		verifiedB := state.VerifiedBatch{
			BlockNumber: lastVerifiedBatch.BlockNumber,
			BatchNumber: lastVBatch.BatchNumber + i,
			Aggregator:  lastVerifiedBatch.Aggregator,
			StateRoot:   lastVerifiedBatch.StateRoot,
			TxHash:      lastVerifiedBatch.TxHash,
			IsTrusted:   true,
		}
		log.Infof("processTrustedVerifyBatches: Storing verifiedB. BlockNumber: %d, BatchNumber: %d", verifiedB.BlockNumber, verifiedB.BatchNumber)
		err = p.state.AddVerifiedBatch(ctx, &verifiedB, dbTx)
		if err != nil {
			log.Errorf("error storing the verifiedB in processTrustedVerifyBatches. verifiedBatch: %+v, lastVerifiedBatch: %+v", verifiedB, lastVerifiedBatch)
			rollbackErr := dbTx.Rollback(ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", lastVerifiedBatch.BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			log.Errorf("error storing the verifiedB in processTrustedVerifyBatches. BlockNumber: %d, error: %v", lastVerifiedBatch.BlockNumber, err)
			return err
		}
	}
	return nil
}
