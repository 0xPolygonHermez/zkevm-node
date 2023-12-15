/*
object TrustedBatchesRetrieve:
- It get all pending batches from trusted node to be synchronized

You must implements BatchProcessor with the code to process the batches
*/
package l2_shared

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/metrics"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

const (
	firstTrustedBatchNumber = uint64(2)
)

// ZkEVMClientInterface contains the methods required to interact with zkEVM-RPC
type ZkEVMClientInterface interface {
	BatchNumber(ctx context.Context) (uint64, error)
	BatchByNumber(ctx context.Context, number *big.Int) (*types.Batch, error)
}

// StateInterface contains the methods required to interact with the state.
type StateInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
}

// BatchProcessor is a interface with the ProcessTrustedBatch methor
//
//	this method is responsible to process a trusted batch
type BatchProcessor interface {
	// ProcessTrustedBatch processes a trusted batch
	ProcessTrustedBatch(ctx context.Context, trustedBatch *types.Batch, status TrustedState, dbTx pgx.Tx) (*TrustedState, error)
}

// SyncInterface contains the methods required to interact with the synchronizer main class.
type SyncInterface interface {
	PendingFlushID(flushID uint64, proverID string)
	CheckFlushID(dbTx pgx.Tx) error
}

// StateRootEntry is the state root entry, basically contains the batch number and the state root. The stateRoot could be a intermediate state root
type StateRootEntry struct {
	// Last batch processed
	batchNumber uint64
	// State root for lastBatchNumber.
	// - If not closed is the intermediate state root
	StateRoot common.Hash
}

// TrustedState is the trusted state, basically contains the batch cache and the last state root (could be a intermediate state root)
type TrustedState struct {
	// LastTrustedBatches [0] -> Current  batch, [1] -> previous batch
	LastTrustedBatches []*state.Batch
	// LastStateRoot is the last state root, it have the batchNumber to be sure that is the expected one
	LastStateRoot *StateRootEntry
}

// TrustedBatchesRetrieve it gets pending batches from Trusted node. It calls for each batch to BatchExecutor
//
//	and for each new batch calls the ProcessTrustedBatch method of the BatchExecutor interface
type TrustedBatchesRetrieve struct {
	batchExecutor          BatchProcessor
	zkEVMClient            ZkEVMClientInterface
	state                  StateInterface
	sync                   SyncInterface
	TrustedState           TrustedState
	firstBatchNumberToSync uint64
}

// NewSyncTrustedStateTemplate creates a new SyncTrustedStateTemplate
func NewSyncTrustedStateTemplate(batchExecutor BatchProcessor, zkEVMClient ZkEVMClientInterface, state StateInterface, sync SyncInterface) *TrustedBatchesRetrieve {
	return &TrustedBatchesRetrieve{
		batchExecutor:          batchExecutor,
		zkEVMClient:            zkEVMClient,
		state:                  state,
		sync:                   sync,
		TrustedState:           TrustedState{},
		firstBatchNumberToSync: firstTrustedBatchNumber,
	}
}

// CleanTrustedState Clean cache of TrustedBatches and StateRoot
func (s *TrustedBatchesRetrieve) CleanTrustedState() {
	s.TrustedState.LastTrustedBatches = nil
	s.TrustedState.LastStateRoot = nil
}

// SyncTrustedState sync trusted state from latestSyncedBatch to lastTrustedStateBatchNumber
func (s *TrustedBatchesRetrieve) SyncTrustedState(ctx context.Context, latestSyncedBatch uint64) error {
	log.Info("syncTrustedState: Getting trusted state info")
	if latestSyncedBatch == 0 {
		log.Info("syncTrustedState: latestSyncedBatch is 0, assuming first batch as 1")
		latestSyncedBatch = 1
	}
	lastTrustedStateBatchNumber, err := s.zkEVMClient.BatchNumber(ctx)

	if err != nil {
		log.Warn("syncTrustedState: error getting last batchNumber from Trusted Node. Error: ", err)
		return err
	}
	log.Infof("syncTrustedState: latestSyncedBatch:%d syncTrustedState:%d", latestSyncedBatch, lastTrustedStateBatchNumber)

	if isSyncrhonizedTrustedState(lastTrustedStateBatchNumber, latestSyncedBatch, s.firstBatchNumberToSync) {
		log.Info("syncTrustedState: Trusted state is synchronized")
		return nil
	}
	return s.syncTrustedBatchesToFrom(ctx, latestSyncedBatch, lastTrustedStateBatchNumber)
}

func isSyncrhonizedTrustedState(lastTrustedStateBatchNumber uint64, latestSyncedBatch uint64, firstBatchNumberToSync uint64) bool {
	if lastTrustedStateBatchNumber < firstBatchNumberToSync {
		return true
	}
	return lastTrustedStateBatchNumber < latestSyncedBatch
}

func (s *TrustedBatchesRetrieve) syncTrustedBatchesToFrom(ctx context.Context, latestSyncedBatch uint64, lastTrustedStateBatchNumber uint64) error {
	batchNumberToSync := max(latestSyncedBatch, s.firstBatchNumberToSync)
	for batchNumberToSync <= lastTrustedStateBatchNumber {
		debugPrefix := fmt.Sprintf("syncTrustedState: batch[%d/%d]", batchNumberToSync, lastTrustedStateBatchNumber)
		start := time.Now()
		batchToSync, err := s.zkEVMClient.BatchByNumber(ctx, big.NewInt(0).SetUint64(batchNumberToSync))
		metrics.GetTrustedBatchInfoTime(time.Since(start))
		if err != nil {
			log.Warnf("%s failed to get batch %d from trusted state. Error: %v", debugPrefix, batchNumberToSync, err)
			return err
		}

		dbTx, err := s.state.BeginStateTransaction(ctx)
		if err != nil {
			log.Errorf("%s error creating db transaction to sync trusted batch %d: %v", debugPrefix, batchNumberToSync, err)
			return err
		}
		start = time.Now()
		cbatches, err := s.getCurrentBatches(ctx, s.TrustedState.LastTrustedBatches, batchToSync, dbTx)
		if err != nil {
			log.Errorf("%s error getting current batches to sync trusted batch %d: %v", debugPrefix, batchNumberToSync, err)
			return rollback(ctx, dbTx, err)
		}
		previousStatus := TrustedState{
			LastTrustedBatches: cbatches,
			LastStateRoot:      s.TrustedState.LastStateRoot,
		}
		log.Debugf("%s processing trusted batch %d", debugPrefix, batchNumberToSync)
		newTrustedState, err := s.batchExecutor.ProcessTrustedBatch(ctx, batchToSync, previousStatus, dbTx)
		metrics.ProcessTrustedBatchTime(time.Since(start))
		if err != nil {
			log.Errorf("%s error processing trusted batch %d: %v", debugPrefix, batchNumberToSync, err)
			return rollback(ctx, dbTx, err)
		}
		log.Debug("%s Checking FlushID to commit trustedState data to db", debugPrefix)
		err = s.sync.CheckFlushID(dbTx)
		if err != nil {
			log.Errorf("%s error checking flushID. Error: %v", debugPrefix, err)
			return rollback(ctx, dbTx, err)
		}

		if err := dbTx.Commit(ctx); err != nil {
			log.Errorf("%s error committing db transaction to sync trusted batch %v: %v", debugPrefix, batchNumberToSync, err)
			return err
		}
		//s.TrustedState.LastTrustedBatches = cbatches
		//s.TrustedState.LastStateRoot = lastStateRoot
		s.TrustedState = *newTrustedState
		batchNumberToSync++
	}

	log.Infof("syncTrustedState: Trusted state fully synchronized from %d to %d", latestSyncedBatch, lastTrustedStateBatchNumber)
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

func (s *TrustedBatchesRetrieve) getCurrentBatches(ctx context.Context, batches []*state.Batch, trustedBatch *types.Batch, dbTx pgx.Tx) ([]*state.Batch, error) {
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
