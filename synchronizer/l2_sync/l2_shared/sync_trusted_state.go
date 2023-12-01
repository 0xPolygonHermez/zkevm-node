/*
Package actionshared contains shared code for actions.
Shared objects between implementations.

If some action need to change this could stop using the share object
*/
package l2_shared

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/metrics"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type ZkEVMClientInterface interface {
	BatchNumber(ctx context.Context) (uint64, error)
	BatchByNumber(ctx context.Context, number *big.Int) (*types.Batch, error)
}

type StateInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
}

type BatchExecutor interface {
	// ProcessTrustedBatch processes a trusted batch
	//ProcessTrustedBatch(ctx context.Context, trustedBatch *types.Batch, dbTx pgx.Tx) ([]*state.Batch, *common.Hash, error)
	ProcessTrustedBatch(ctx context.Context, trustedBatch *types.Batch, status TrustedState, dbTx pgx.Tx) (*TrustedState, error)
}

type SyncInterface interface {
	PendingFlushID(flushID uint64, proverID string)
	CheckFlushID(dbTx pgx.Tx) error
}
type StateRootEntry struct {
	// Last batch processed
	batchNumber uint64
	// State root for lastBatchNumber.
	// - If not closed is the intermediate state root
	StateRoot common.Hash
}
type TrustedState struct {
	// LastTrustedBatches [0] -> Current  batch, [1] -> previous batch
	LastTrustedBatches []*state.Batch
	// LastStateRoot is the last state root for LastStateRootBatchNumber
	LastStateRoot *StateRootEntry
}

type SyncTrustedStateTemplate struct {
	steps        BatchExecutor
	zkEVMClient  ZkEVMClientInterface
	state        StateInterface
	sync         SyncInterface
	TrustedState TrustedState
}

func NewSyncTrustedStateTemplate(steps BatchExecutor, zkEVMClient ZkEVMClientInterface, state StateInterface, sync SyncInterface) *SyncTrustedStateTemplate {
	return &SyncTrustedStateTemplate{
		steps:        steps,
		zkEVMClient:  zkEVMClient,
		state:        state,
		sync:         sync,
		TrustedState: TrustedState{},
	}
}

func (s *SyncTrustedStateTemplate) SyncTrustedState(ctx context.Context, latestSyncedBatch uint64) error {
	log.Info("syncTrustedState: Getting trusted state info")
	if latestSyncedBatch == 0 {
		log.Info("syncTrustedState: latestSyncedBatch is 0, assuming first batch as 1")
		latestSyncedBatch = 1
	}
	lastTrustedStateBatchNumber, err := s.zkEVMClient.BatchNumber(ctx)

	if err != nil {
		log.Warn("syncTrustedState: error syncing trusted state. Error: ", err)
		return err
	}
	log.Infof("syncTrustedState: latestSyncedBatch:%d syncTrustedState:%d", latestSyncedBatch, lastTrustedStateBatchNumber)

	if isSyncrhonizedTrustedState(lastTrustedStateBatchNumber, latestSyncedBatch) {
		log.Info("syncTrustedState: Trusted state is synchronized")
		return nil
	}
	return s.syncTrustedBatchesToFrom(ctx, latestSyncedBatch, lastTrustedStateBatchNumber)
}

func isSyncrhonizedTrustedState(lastTrustedStateBatchNumber uint64, latestSyncedBatch uint64) bool {
	return lastTrustedStateBatchNumber < latestSyncedBatch
}

func (s *SyncTrustedStateTemplate) syncTrustedBatchesToFrom(ctx context.Context, latestSyncedBatch uint64, lastTrustedStateBatchNumber uint64) error {
	batchNumberToSync := latestSyncedBatch
	for batchNumberToSync <= lastTrustedStateBatchNumber {
		start := time.Now()
		batchToSync, err := s.zkEVMClient.BatchByNumber(ctx, big.NewInt(0).SetUint64(batchNumberToSync))
		metrics.GetTrustedBatchInfoTime(time.Since(start))
		if err != nil {
			log.Warnf("syncTrustedState: failed to get batch %d from trusted state. Error: %v", batchNumberToSync, err)
			return err
		}

		dbTx, err := s.state.BeginStateTransaction(ctx)
		if err != nil {
			log.Errorf("syncTrustedState: error creating db transaction to sync trusted batch %d: %v", batchNumberToSync, err)
			return err
		}
		start = time.Now()
		cbatches, err := s.getCurrentBatches(ctx, s.TrustedState.LastTrustedBatches, batchToSync, dbTx)
		if err != nil {
			log.Errorf("syncTrustedState: error getting current batches to sync trusted batch %d: %v", batchNumberToSync, err)
			return rollback(ctx, dbTx, err)
		}
		previousStatus := TrustedState{
			LastTrustedBatches: cbatches,
			LastStateRoot:      s.TrustedState.LastStateRoot,
		}
		//cbatches, lastStateRoot, err := s.steps.ProcessTrustedBatch(ctx, batchToSync, dbTx)
		newTrustedState, err := s.steps.ProcessTrustedBatch(ctx, batchToSync, previousStatus, dbTx)
		metrics.ProcessTrustedBatchTime(time.Since(start))
		if err != nil {
			log.Errorf("syncTrustedState: error processing trusted batch %d: %v", batchNumberToSync, err)
			return rollback(ctx, dbTx, err)
		}
		log.Debug("syncTrustedState: Checking FlushID to commit trustedState data to db")
		err = s.sync.CheckFlushID(dbTx)
		if err != nil {
			log.Errorf("syncTrustedState: error checking flushID. Error: %v", err)
			return rollback(ctx, dbTx, err)
		}

		if err := dbTx.Commit(ctx); err != nil {
			log.Errorf("syncTrustedState: error committing db transaction to sync trusted batch %v: %v", batchNumberToSync, err)
			return err
		}
		//s.TrustedState.LastTrustedBatches = cbatches
		//s.TrustedState.LastStateRoot = lastStateRoot
		s.TrustedState = *newTrustedState
		batchNumberToSync++
	}

	log.Info("syncTrustedState: Trusted state fully synchronized")
	return nil
}

func rollback(ctx context.Context, dbTx pgx.Tx, err error) error {
	rollbackErr := dbTx.Rollback(ctx)
	if rollbackErr != nil {
		log.Errorf("syncTrustedState: error rolling back state. RollbackErr: %s, Error : %v", rollbackErr.Error(), err)
		return rollbackErr
	}
	return err
}

func (s *SyncTrustedStateTemplate) getCurrentBatches(ctx context.Context, batches []*state.Batch, trustedBatch *types.Batch, dbTx pgx.Tx) ([]*state.Batch, error) {
	if len(batches) == 0 || batches[0] == nil || (batches[0] != nil && uint64(trustedBatch.Number) != batches[0].BatchNumber) {
		log.Debug("Updating batch[0] value!")
		batch, err := s.state.GetBatchByNumber(ctx, uint64(trustedBatch.Number), dbTx)
		if err != nil && err != state.ErrNotFound {
			log.Warnf("failed to get batch %v from local trusted state. Error: %v", trustedBatch.Number, err)
			return nil, err
		}
		var prevBatch *state.Batch
		if len(batches) == 0 || batches[0] == nil || (batches[0] != nil && uint64(trustedBatch.Number-1) != batches[0].BatchNumber) {
			log.Debug("Updating batch[1] value!")
			prevBatch, err = s.state.GetBatchByNumber(ctx, uint64(trustedBatch.Number-1), dbTx)
			if err != nil && err != state.ErrNotFound {
				log.Warnf("failed to get prevBatch %v from local trusted state. Error: %v", trustedBatch.Number-1, err)
				return nil, err
			}
		} else {
			prevBatch = batches[0]
		}
		log.Debug("batch: ", batch)
		log.Debug("prevBatch: ", prevBatch)
		batches = []*state.Batch{batch, prevBatch}
	}
	return batches, nil
}
