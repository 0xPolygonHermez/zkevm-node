package sequencer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg Config

	pool      txPool
	state     stateInterface
	txManager txManager
	etherman  etherman

	address common.Address
}

// New init sequencer
func New(cfg Config, txPool txPool, state stateInterface, etherman etherman, manager txManager) (*Sequencer, error) {
	addr, err := etherman.TrustedSequencer()
	if err != nil {
		return nil, fmt.Errorf("failed to get trusted sequencer address, err: %v", err)
	}

	return &Sequencer{
		cfg:       cfg,
		pool:      txPool,
		state:     state,
		etherman:  etherman,
		txManager: manager,
		address:   addr,
	}, nil
}

// Start starts the sequencer
func (s *Sequencer) Start(ctx context.Context) {
	for !s.isSynced(ctx) {
		log.Infof("waiting for synchronizer to sync...")
		time.Sleep(s.cfg.WaitPeriodPoolIsEmpty.Duration)
	}
	metrics.Register()

	worker := newWorker()
	dbManager := newDBManager(s.pool, s.state, worker)
	go dbManager.Start()

	currBatch, OldAccInputHash, OldStateRoot := s.bootstrap(ctx, dbManager)
	finalizer := newFinalizer(s.cfg.Finalizer, worker, dbManager, s.state, s.address, s.isSynced, s.cfg.MaxTxsPerBatch)
	go finalizer.Start(ctx, currBatch, OldStateRoot, OldAccInputHash)

	closingSignalsManager := newClosingSignalsManager(finalizer)
	go closingSignalsManager.Start()

	go s.trackOldTxs(ctx)
	tickerProcessTxs := time.NewTicker(s.cfg.WaitPeriodPoolIsEmpty.Duration)
	tickerSendSequence := time.NewTicker(s.cfg.WaitPeriodSendSequence.Duration)
	defer tickerProcessTxs.Stop()
	defer tickerSendSequence.Stop()

	go func() {
		for {
			s.tryToSendSequence(ctx, tickerSendSequence)
		}
	}()
	// Wait until context is done
	<-ctx.Done()
}

func (s *Sequencer) bootstrap(ctx context.Context, dbManager *dbManager) (WipBatch, common.Hash, common.Hash) {
	var (
		currBatch                     WipBatch
		oldAccInputHash, oldStateRoot common.Hash
	)
	batchNum, err := dbManager.GetLastBatchNumber(ctx)
	for err != nil {
		if errors.Is(err, state.ErrStateNotSynchronized) {
			log.Warnf("state is not synchronized, trying to get last batch num once again...")
			time.Sleep(s.cfg.WaitPeriodPoolIsEmpty.Duration)
			batchNum, err = dbManager.GetLastBatchNumber(ctx)
		} else {
			log.Fatalf("failed to get last batch number, err: %v", err)
		}
	}
	if batchNum == 0 {
		///////////////////
		// GENESIS Batch //
		///////////////////
		processingCtx := dbManager.CreateFirstBatch(ctx, s.address)
		currBatch = WipBatch{
			globalExitRoot: processingCtx.GlobalExitRoot,
			batchNumber:    processingCtx.BatchNumber,
			coinbase:       processingCtx.Coinbase,
			timestamp:      uint64(processingCtx.Timestamp.Unix()),
			txs:            make([]TxTracker, 0, s.cfg.MaxTxsPerBatch),
		}
	} else {
		// Check if synchronizer is up to date
		for !s.isSynced(ctx) {
			log.Info("wait for synchronizer to sync last batch")
			time.Sleep(time.Second)
		}
		// Revert reorged txs to pending
		if err = dbManager.MarkReorgedTxsAsPending(ctx); err != nil {
			log.Fatal(fmt.Errorf("failed to mark reorged txs as pending, err: %w", err))
		}
		// Get current wip batch
		currBatch, err = dbManager.GetWIPBatch(ctx)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to load batch from the state, err: %w", err))
		}
		// Get data for prevBatch
		lastBatch, err := dbManager.GetLastBatch(ctx)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to get last batch. err: %w", err))
		}
		isClosed, err := dbManager.IsBatchClosed(ctx, lastBatch.BatchNumber)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to check is batch closed or not, err: %w", err))
		}
		if isClosed {
			//ger, _, err := s.getLatestGer(ctx, dbTx)
			// TODO: Open New batch and create WipBatch (currBatch)
		} else {
			if lastBatch.BatchNumber == 1 {
				oldAccInputHash = lastBatch.AccInputHash
				oldStateRoot = lastBatch.StateRoot
			} else {
				n := uint(2)
				batches, err := dbManager.GetLastNBatches(ctx, n)
				if err != nil {
					log.Fatal(fmt.Errorf("failed to get last %d batches, err: %w", n, err))
				}
				oldAccInputHash = batches[1].AccInputHash
				oldStateRoot = batches[1].StateRoot
			}

		}
	}
	return currBatch, oldAccInputHash, oldStateRoot
}

func (s *Sequencer) trackOldTxs(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.FrequencyToCheckTxsForDelete.Duration)
	for {
		waitTick(ctx, ticker)
		log.Infof("trying to get txs to delete from the pool...")
		txHashes, err := s.state.GetTxsOlderThanNL1Blocks(ctx, s.cfg.BlocksAmountForTxsToBeDeleted, nil)
		if err != nil {
			log.Errorf("failed to get txs hashes to delete, err: %v", err)
			continue
		}
		log.Infof("will try to delete %d redundant txs", len(txHashes))
		err = s.pool.DeleteTxsByHashes(ctx, txHashes)
		if err != nil {
			log.Errorf("failed to delete txs from the pool, err: %v", err)
			continue
		}
		log.Infof("deleted %d selected txs from the pool", len(txHashes))
	}
}

func waitTick(ctx context.Context, ticker *time.Ticker) {
	select {
	case <-ticker.C:
		// nothing
	case <-ctx.Done():
		return
	}
}

func (s *Sequencer) isSynced(ctx context.Context) bool {
	lastSyncedBatchNum, err := s.state.GetLastVirtualBatchNum(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get last isSynced batch, err: %v", err)
		return false
	}
	lastEthBatchNum, err := s.etherman.GetLatestBatchNumber()
	if err != nil {
		log.Errorf("failed to get last eth batch, err: %v", err)
		return false
	}
	if lastSyncedBatchNum < lastEthBatchNum {
		log.Infof("waiting for the state to be isSynced, lastSyncedBatchNum: %d, lastEthBatchNum: %d", lastSyncedBatchNum, lastEthBatchNum)
		return false
	}

	return true
}

/*
func (s *Sequencer) loadSequenceFromState(ctx context.Context) error {
	// Check if synchronizer is up to date
	for !s.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}
	// Revert reorged txs to pending
	if err := s.pool.MarkReorgedTxsAsPending(ctx); err != nil {
		return fmt.Errorf("failed to mark reorged txs as pending, err: %w", err)
	}
	// Get latest info from the state
	lastBatch, err := s.state.GetWIPBatch(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get last batch, err: %w", err)
	}
	isClosed, err := s.state.IsBatchClosed(ctx, lastBatch.BatchNumber, nil)
	if err != nil {
		return fmt.Errorf("failed to check is batch closed or not, err: %w", err)
	}
	if isClosed {
		dbTx, err := s.state.BeginStateTransaction(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin state tx to open a batch, err: %w", err)
		}
		ger, _, err := s.getLatestGer(ctx, dbTx)
		if err != nil {
			if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
				return fmt.Errorf(
					"failed to rollback dbTx when getting last globalExitRoot that gave err: %s. Rollback err: %s",
					rollbackErr.Error(), err.Error(),
				)
			}
			return fmt.Errorf("failed to get latest global exit root, err: %w", err)
		}
		processingCtx := state.ProcessingContext{
			BatchNumber:    lastBatch.BatchNumber + 1,
			Coinbase:       s.address,
			Timestamp:      time.Now(),
			globalExitRoot: ger.globalExitRoot,
		}
		err = s.state.OpenBatch(ctx, processingCtx, dbTx)
		if err != nil {
			rollErr := dbTx.Rollback(ctx)
			if rollErr != nil {
				err = fmt.Errorf("failed to open a batch, err: %w. Rollback err: %v", err, rollErr)
			}
			return err
		}
		if err = dbTx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit a state tx to open a batch, err: %w", err)
		}
		s.sequenceInProgress = types.Sequence{
			globalExitRoot: processingCtx.globalExitRoot,
			Timestamp:      processingCtx.Timestamp.Unix(),
		}
	} else {
		txs, err := s.state.GetTransactionsByBatchNumber(ctx, lastBatch.BatchNumber, nil)
		if err != nil {
			return fmt.Errorf("failed to get tx by batch number, err: %w", err)
		}
		s.sequenceInProgress = types.Sequence{
			globalExitRoot: lastBatch.globalExitRoot,
			Timestamp:      lastBatch.Timestamp.Unix(),
			Txs:            txs,
		}
		// TODO: execute to get state root and LER or change open/closed logic so we always store state root and LER and add an open flag
	}

	return nil

}
*/
