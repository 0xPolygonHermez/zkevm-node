package sequencer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/ethereum/go-ethereum/common"
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg      Config
	batchCfg state.BatchConfig
	poolCfg  pool.Config

	pool     txPool
	state    stateInterface
	eventLog *event.EventLog
	etherman etherman

	address common.Address
}

// L2ReorgEvent is the event that is triggered when a reorg happens in the L2
type L2ReorgEvent struct {
	TxHashes []common.Hash
}

// ClosingSignalCh is a struct that contains all the channels that are used to receive batch closing signals
type ClosingSignalCh struct {
	ForcedBatchCh chan state.ForcedBatch
	GERCh         chan common.Hash
	L2ReorgCh     chan L2ReorgEvent
}

// New init sequencer
func New(cfg Config, batchCfg state.BatchConfig, poolCfg pool.Config, txPool txPool, state stateInterface, etherman etherman, eventLog *event.EventLog) (*Sequencer, error) {
	addr, err := etherman.TrustedSequencer()
	if err != nil {
		return nil, fmt.Errorf("failed to get trusted sequencer address, err: %v", err)
	}

	sequencer := &Sequencer{
		cfg:      cfg,
		batchCfg: batchCfg,
		poolCfg:  poolCfg,
		pool:     txPool,
		state:    state,
		etherman: etherman,
		address:  addr,
		eventLog: eventLog,
	}

	return sequencer, nil
}

// Start starts the sequencer
func (s *Sequencer) Start(ctx context.Context) {
	for !s.isSynced(ctx) {
		log.Infof("waiting for synchronizer to sync...")
		time.Sleep(s.cfg.WaitPeriodPoolIsEmpty.Duration)
	}
	metrics.Register()

	closingSignalCh := ClosingSignalCh{
		ForcedBatchCh: make(chan state.ForcedBatch),
		GERCh:         make(chan common.Hash),
		L2ReorgCh:     make(chan L2ReorgEvent),
	}

	err := s.pool.MarkWIPTxsAsPending(ctx)
	if err != nil {
		log.Fatalf("failed to mark WIP txs as pending, err: %v", err)
	}

	worker := NewWorker(s.state, s.batchCfg.Constraints)
	dbManager := newDBManager(ctx, s.cfg.DBManager, s.pool, s.state, worker, closingSignalCh, s.batchCfg.Constraints)

	// Start stream server if enabled
	if s.cfg.StreamServer.Enabled {
		streamServer, err := datastreamer.NewServer(s.cfg.StreamServer.Port, state.StreamTypeSequencer, s.cfg.StreamServer.Filename, &s.cfg.StreamServer.Log)
		if err != nil {
			log.Fatalf("failed to create stream server, err: %v", err)
		}

		dbManager.streamServer = streamServer
		err = dbManager.streamServer.Start()
		if err != nil {
			log.Fatalf("failed to start stream server, err: %v", err)
		}

		s.updateDataStreamerFile(ctx, streamServer)
	}

	go dbManager.Start()

	var streamServer *datastreamer.StreamServer = nil
	if s.cfg.StreamServer.Enabled {
		streamServer = dbManager.streamServer
	}

	finalizer := newFinalizer(s.cfg.Finalizer, s.poolCfg, worker, dbManager, s.state, s.address, s.isSynced, closingSignalCh, s.batchCfg.Constraints, s.eventLog, streamServer)
	currBatch, processingReq := s.bootstrap(ctx, dbManager, finalizer)
	go finalizer.Start(ctx, currBatch, processingReq)

	closingSignalsManager := newClosingSignalsManager(ctx, finalizer.dbManager, closingSignalCh, finalizer.cfg, s.etherman)
	go closingSignalsManager.Start()

	go s.purgeOldPoolTxs(ctx)
	tickerProcessTxs := time.NewTicker(s.cfg.WaitPeriodPoolIsEmpty.Duration)
	defer tickerProcessTxs.Stop()

	// Expire too old txs in the worker
	go func() {
		for {
			time.Sleep(s.cfg.TxLifetimeCheckTimeout.Duration)
			txTrackers := worker.ExpireTransactions(s.cfg.MaxTxLifetime.Duration)
			failedReason := ErrExpiredTransaction.Error()
			for _, txTracker := range txTrackers {
				err := s.pool.UpdateTxStatus(ctx, txTracker.Hash, pool.TxStatusFailed, false, &failedReason)
				metrics.TxProcessed(metrics.TxProcessedLabelFailed, 1)
				if err != nil {
					log.Errorf("failed to update tx status, err: %v", err)
				}
			}
		}
	}()

	// Wait until context is done
	<-ctx.Done()
}

func (s *Sequencer) updateDataStreamerFile(ctx context.Context, streamServer *datastreamer.StreamServer) {
	err := state.GenerateDataStreamerFile(ctx, streamServer, s.state)
	if err != nil {
		log.Fatalf("failed to generate data streamer file, err: %v", err)
	}
	log.Info("Data streamer file updated")
}

func (s *Sequencer) bootstrap(ctx context.Context, dbManager *dbManager, finalizer *finalizer) (*WipBatch, *state.ProcessRequest) {
	var (
		currBatch      *WipBatch
		processRequest *state.ProcessRequest
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
		timestamp := processingCtx.Timestamp
		_, oldStateRoot, err := finalizer.getLastBatchNumAndOldStateRoot(ctx)
		if err != nil {
			log.Fatalf("failed to get old state root, err: %v", err)
		}
		processRequest = &state.ProcessRequest{
			BatchNumber:    processingCtx.BatchNumber,
			OldStateRoot:   oldStateRoot,
			GlobalExitRoot: processingCtx.GlobalExitRoot,
			Coinbase:       processingCtx.Coinbase,
			Timestamp:      timestamp,
			Caller:         stateMetrics.SequencerCallerLabel,
		}
		currBatch = &WipBatch{
			globalExitRoot:     processingCtx.GlobalExitRoot,
			initialStateRoot:   oldStateRoot,
			stateRoot:          oldStateRoot,
			batchNumber:        processingCtx.BatchNumber,
			coinbase:           processingCtx.Coinbase,
			timestamp:          timestamp,
			remainingResources: getMaxRemainingResources(finalizer.batchConstraints),
		}
	} else {
		err := finalizer.syncWithState(ctx, &batchNum)
		if err != nil {
			log.Fatalf("failed to sync with state, err: %v", err)
		}
		currBatch = finalizer.batch
		processRequest = &finalizer.processRequest
	}

	return currBatch, processRequest
}

func (s *Sequencer) purgeOldPoolTxs(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.FrequencyToCheckTxsForDelete.Duration)
	for {
		waitTick(ctx, ticker)
		log.Infof("trying to get txs to delete from the pool...")
		txHashes, err := s.state.GetTxsOlderThanNL1Blocks(ctx, s.cfg.BlocksAmountForTxsToBeDeleted, nil)
		if err != nil {
			log.Errorf("failed to get txs hashes to delete, err: %v", err)
			continue
		}
		log.Infof("trying to delete %d selected txs", len(txHashes))
		err = s.pool.DeleteTransactionsByHashes(ctx, txHashes)
		if err != nil {
			log.Errorf("failed to delete selected txs from the pool, err: %v", err)
			continue
		}
		log.Infof("deleted %d selected txs from the pool", len(txHashes))

		log.Infof("trying to delete failed txs from the pool")
		// Delete failed txs older than a certain date (14 seconds per L1 block)
		err = s.pool.DeleteFailedTransactionsOlderThan(ctx, time.Now().Add(-time.Duration(s.cfg.BlocksAmountForTxsToBeDeleted*14)*time.Second)) //nolint:gomnd
		if err != nil {
			log.Errorf("failed to delete failed txs from the pool, err: %v", err)
			continue
		}
		log.Infof("failed txs deleted from the pool")
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
	lastBatchNum, err := s.state.GetLastBatchNumber(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get last batch num, err: %v", err)
		return false
	}
	if lastBatchNum > lastSyncedBatchNum {
		return true
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

func getMaxRemainingResources(constraints state.BatchConstraintsCfg) state.BatchResources {
	return state.BatchResources{
		ZKCounters: state.ZKCounters{
			CumulativeGasUsed:    constraints.MaxCumulativeGasUsed,
			UsedKeccakHashes:     constraints.MaxKeccakHashes,
			UsedPoseidonHashes:   constraints.MaxPoseidonHashes,
			UsedPoseidonPaddings: constraints.MaxPoseidonPaddings,
			UsedMemAligns:        constraints.MaxMemAligns,
			UsedArithmetics:      constraints.MaxArithmetics,
			UsedBinaries:         constraints.MaxBinaries,
			UsedSteps:            constraints.MaxSteps,
		},
		Bytes: constraints.MaxBatchBytesSize,
	}
}
