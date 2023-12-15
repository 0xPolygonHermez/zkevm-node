package sequencer

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

const (
	datastreamChannelMultiplier = 2
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg      Config
	batchCfg state.BatchConfig
	poolCfg  pool.Config

	pool     txPool
	stateI   stateInterface
	eventLog *event.EventLog
	etherman etherman
	worker   *Worker

	streamServer *datastreamer.StreamServer
	dataToStream chan state.DSL2FullBlock

	closingSignalCh ClosingSignalCh

	numberOfStateInconsistencies uint64
	address                      common.Address
}

// L2ReorgEvent is the event that is triggered when a reorg happens in the L2
type L2ReorgEvent struct {
	TxHashes []common.Hash
}

// ClosingSignalCh is a struct that contains all the channels that are used to receive batch closing signals
type ClosingSignalCh struct {
	ForcedBatchCh        chan state.ForcedBatch
	GERCh                chan common.Hash
	L1InfoTreeExitRootCh chan state.L1InfoTreeExitRootStorageEntry
	L2ReorgCh            chan L2ReorgEvent
}

// New init sequencer
func New(cfg Config, batchCfg state.BatchConfig, poolCfg pool.Config, txPool txPool, stateI stateInterface, etherman etherman, eventLog *event.EventLog) (*Sequencer, error) {
	addr, err := etherman.TrustedSequencer()
	if err != nil {
		return nil, fmt.Errorf("failed to get trusted sequencer address, err: %v", err)
	}

	sequencer := &Sequencer{
		cfg:      cfg,
		batchCfg: batchCfg,
		poolCfg:  poolCfg,
		pool:     txPool,
		stateI:   stateI,
		etherman: etherman,
		address:  addr,
		eventLog: eventLog,
	}

	sequencer.dataToStream = make(chan state.DSL2FullBlock, batchCfg.Constraints.MaxTxsPerBatch*datastreamChannelMultiplier)

	return sequencer, nil
}

// Start starts the sequencer
func (s *Sequencer) Start(ctx context.Context) {
	for !s.isSynced(ctx) {
		log.Infof("waiting for synchronizer to sync...")
		time.Sleep(s.cfg.WaitPeriodPoolIsEmpty.Duration)
	}
	metrics.Register()

	s.closingSignalCh = ClosingSignalCh{
		ForcedBatchCh:        make(chan state.ForcedBatch),
		GERCh:                make(chan common.Hash),
		L1InfoTreeExitRootCh: make(chan state.L1InfoTreeExitRootStorageEntry),
		L2ReorgCh:            make(chan L2ReorgEvent),
	}

	err := s.pool.MarkWIPTxsAsPending(ctx)
	if err != nil {
		log.Fatalf("failed to mark WIP txs as pending, err: %v", err)
	}

	s.worker = NewWorker(s.stateI, s.batchCfg.Constraints)
	//dbManager := newDBManager(ctx, s.cfg.DBManager, s.pool, s.state, worker, closingSignalCh, s.batchCfg.Constraints)

	// Start stream server if enabled
	if s.cfg.StreamServer.Enabled {
		s.streamServer, err = datastreamer.NewServer(s.cfg.StreamServer.Port, state.StreamTypeSequencer, s.cfg.StreamServer.Filename, &s.cfg.StreamServer.Log)
		if err != nil {
			log.Fatalf("failed to create stream server, err: %v", err)
		}

		err = s.streamServer.Start()
		if err != nil {
			log.Fatalf("failed to start stream server, err: %v", err)
		}

		s.updateDataStreamerFile(ctx)
	}

	go s.loadFromPool(ctx)

	go func() {
		for {
			time.Sleep(s.cfg.DBManager.L2ReorgRetrievalInterval.Duration)
			s.checkStateInconsistency(ctx)
		}
	}()

	if s.streamServer != nil {
		go s.sendDataToStreamer()
	}

	finalizer := newFinalizer(s.cfg.Finalizer, s.poolCfg, s.worker, s.pool, s.stateI, s.etherman, s.address, s.isSynced, s.closingSignalCh, s.batchCfg.Constraints, s.eventLog, s.streamServer, s.dataToStream)
	go finalizer.Start(ctx)

	closingSignalsManager := newClosingSignalsManager(ctx, s.stateI, s.closingSignalCh, finalizer.cfg, s.etherman)
	go closingSignalsManager.Start()

	go s.purgeOldPoolTxs(ctx) //TODO: Review if this function is needed as we have other go func to expire old txs in the worker

	// Expire too old txs in the worker
	go func() {
		for {
			time.Sleep(s.cfg.TxLifetimeCheckTimeout.Duration)
			txTrackers := s.worker.ExpireTransactions(s.cfg.MaxTxLifetime.Duration)
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

func (s *Sequencer) updateDataStreamerFile(ctx context.Context) {
	err := state.GenerateDataStreamerFile(ctx, s.streamServer, s.stateI, true, nil)
	if err != nil {
		log.Fatalf("failed to generate data streamer file, err: %v", err)
	}
	log.Info("Data streamer file updated")
}

func (s *Sequencer) purgeOldPoolTxs(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.FrequencyToCheckTxsForDelete.Duration)
	for {
		waitTick(ctx, ticker)
		log.Infof("trying to get txs to delete from the pool...")
		txHashes, err := s.stateI.GetTxsOlderThanNL1Blocks(ctx, s.cfg.BlocksAmountForTxsToBeDeleted, nil)
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

// checkStateInconsistency checks if state inconsistency happened
func (s *Sequencer) checkStateInconsistency(ctx context.Context) {
	stateInconsistenciesDetected, err := s.stateI.CountReorgs(ctx, nil)
	if err != nil {
		log.Error("failed to get number of reorgs: %v", err)
		return
	}

	if stateInconsistenciesDetected != s.numberOfStateInconsistencies {
		log.Warnf("New State Inconsistency detected")
		s.closingSignalCh.L2ReorgCh <- L2ReorgEvent{}
	}
}

// loadFromPool keeps loading transactions from the pool
func (s *Sequencer) loadFromPool(ctx context.Context) {
	for {
		time.Sleep(s.cfg.DBManager.PoolRetrievalInterval.Duration)

		poolTransactions, err := s.pool.GetNonWIPPendingTxs(ctx)
		if err != nil && err != pool.ErrNotFound {
			log.Errorf("load tx from pool: %v", err)
		}

		for _, tx := range poolTransactions {
			err := s.addTxToWorker(ctx, tx)
			if err != nil {
				log.Errorf("error adding transaction to worker: %v", err)
			}
		}
	}
}

func (s *Sequencer) addTxToWorker(ctx context.Context, tx pool.Transaction) error {
	txTracker, err := s.worker.NewTxTracker(tx.Transaction, tx.ZKCounters, tx.IP)
	if err != nil {
		return err
	}
	replacedTx, dropReason := s.worker.AddTxTracker(ctx, txTracker)
	if dropReason != nil {
		failedReason := dropReason.Error()
		return s.pool.UpdateTxStatus(ctx, txTracker.Hash, pool.TxStatusFailed, false, &failedReason)
	} else {
		if replacedTx != nil {
			failedReason := ErrReplacedTransaction.Error()
			error := s.pool.UpdateTxStatus(ctx, replacedTx.Hash, pool.TxStatusFailed, false, &failedReason)
			if error != nil {
				log.Warnf("error when setting as failed replacedTx(%s)", replacedTx.HashStr)
			}
		}
		return s.pool.UpdateTxWIPStatus(ctx, tx.Hash(), true)
	}
}

// sendDataToStreamer sends data to the data stream server
func (s *Sequencer) sendDataToStreamer() {
	var err error
	for {
		// Read error from previous iteration
		if err != nil {
			err = s.streamServer.RollbackAtomicOp()
			if err != nil {
				log.Errorf("failed to rollback atomic op: %v", err)
			}
			s.streamServer = nil
		}

		// Read data from channel
		fullL2Block := <-s.dataToStream

		l2Block := fullL2Block
		l2Transactions := fullL2Block.Txs

		if s.streamServer != nil {
			err = s.streamServer.StartAtomicOp()
			if err != nil {
				log.Errorf("failed to start atomic op for l2block %v: %v ", l2Block.L2BlockNumber, err)
				continue
			}

			bookMark := state.DSBookMark{
				Type:          state.BookMarkTypeL2Block,
				L2BlockNumber: l2Block.L2BlockNumber,
			}

			_, err = s.streamServer.AddStreamBookmark(bookMark.Encode())
			if err != nil {
				log.Errorf("failed to add stream bookmark for l2block %v: %v", l2Block.L2BlockNumber, err)
				continue
			}

			blockStart := state.DSL2BlockStart{
				BatchNumber:    l2Block.BatchNumber,
				L2BlockNumber:  l2Block.L2BlockNumber,
				Timestamp:      l2Block.Timestamp,
				GlobalExitRoot: l2Block.GlobalExitRoot,
				Coinbase:       l2Block.Coinbase,
				ForkID:         l2Block.ForkID,
			}

			_, err = s.streamServer.AddStreamEntry(state.EntryTypeL2BlockStart, blockStart.Encode())
			if err != nil {
				log.Errorf("failed to add stream entry for l2block %v: %v", l2Block.L2BlockNumber, err)
				continue
			}

			for _, l2Transaction := range l2Transactions {
				_, err = s.streamServer.AddStreamEntry(state.EntryTypeL2Tx, l2Transaction.Encode())
				if err != nil {
					log.Errorf("failed to add l2tx stream entry for l2block %v: %v", l2Block.L2BlockNumber, err)
					continue
				}
			}

			blockEnd := state.DSL2BlockEnd{
				L2BlockNumber: l2Block.L2BlockNumber,
				BlockHash:     l2Block.BlockHash,
				StateRoot:     l2Block.StateRoot,
			}

			_, err = s.streamServer.AddStreamEntry(state.EntryTypeL2BlockEnd, blockEnd.Encode())
			if err != nil {
				log.Errorf("failed to add stream entry for l2block %v: %v", l2Block.L2BlockNumber, err)
				continue
			}

			err = s.streamServer.CommitAtomicOp()
			if err != nil {
				log.Errorf("failed to commit atomic op for l2block %v: %v ", l2Block.L2BlockNumber, err)
				continue
			}
		}
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
	lastSyncedBatchNum, err := s.stateI.GetLastVirtualBatchNum(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get last isSynced batch, err: %v", err)
		return false
	}
	lastBatchNum, err := s.stateI.GetLastBatchNumber(ctx, nil)
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
