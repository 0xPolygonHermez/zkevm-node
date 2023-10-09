package sequencer

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
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
	cfg Config

	pool         txPool
	state        stateInterface
	eventLog     *event.EventLog
	ethTxManager ethTxManager
	etherman     etherman

	address common.Address
}

// batchConstraints represents the constraints for a batch
type batchConstraints struct {
	MaxTxsPerBatch       uint64
	MaxBatchBytesSize    uint64
	MaxCumulativeGasUsed uint64
	MaxKeccakHashes      uint32
	MaxPoseidonHashes    uint32
	MaxPoseidonPaddings  uint32
	MaxMemAligns         uint32
	MaxArithmetics       uint32
	MaxBinaries          uint32
	MaxSteps             uint32
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
func New(cfg Config, txPool txPool, state stateInterface, etherman etherman, manager ethTxManager, eventLog *event.EventLog) (*Sequencer, error) {
	addr, err := etherman.TrustedSequencer()
	if err != nil {
		return nil, fmt.Errorf("failed to get trusted sequencer address, err: %v", err)
	}

	sequencer := &Sequencer{
		cfg:          cfg,
		pool:         txPool,
		state:        state,
		etherman:     etherman,
		ethTxManager: manager,
		address:      addr,
		eventLog:     eventLog,
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

	batchConstraints := batchConstraints{
		MaxTxsPerBatch:       s.cfg.MaxTxsPerBatch,
		MaxBatchBytesSize:    s.cfg.MaxBatchBytesSize,
		MaxCumulativeGasUsed: s.cfg.MaxCumulativeGasUsed,
		MaxKeccakHashes:      s.cfg.MaxKeccakHashes,
		MaxPoseidonHashes:    s.cfg.MaxPoseidonHashes,
		MaxPoseidonPaddings:  s.cfg.MaxPoseidonPaddings,
		MaxMemAligns:         s.cfg.MaxMemAligns,
		MaxArithmetics:       s.cfg.MaxArithmetics,
		MaxBinaries:          s.cfg.MaxBinaries,
		MaxSteps:             s.cfg.MaxSteps,
	}

	err := s.pool.MarkWIPTxsAsPending(ctx)
	if err != nil {
		log.Fatalf("failed to mark WIP txs as pending, err: %v", err)
	}

	worker := NewWorker(s.state)

	dbManager := newDBManager(ctx, s.cfg.DBManager, s.pool, s.state, worker, closingSignalCh, batchConstraints)

	// Start stream server if enabled
	if s.cfg.StreamServer.Enabled {
		streamServer, err := datastreamer.New(s.cfg.StreamServer.Port, state.StreamTypeSequencer, s.cfg.StreamServer.Filename, &s.cfg.StreamServer.Log)
		if err != nil {
			log.Fatalf("failed to create stream server, err: %v", err)
		}

		// Set entities definition
		entriesDefinition := map[datastreamer.EntryType]datastreamer.EntityDefinition{
			state.EntryTypeL2BlockStart: {
				Name:       "L2BlockStart",
				StreamType: state.StreamTypeSequencer,
				Definition: reflect.TypeOf(state.DSL2BlockStart{}),
			},
			state.EntryTypeL2Tx: {
				Name:       "L2Transaction",
				StreamType: state.StreamTypeSequencer,
				Definition: reflect.TypeOf(state.DSL2Transaction{}),
			},
			state.EntryTypeL2BlockEnd: {
				Name:       "L2BlockEnd",
				StreamType: state.StreamTypeSequencer,
				Definition: reflect.TypeOf(state.DSL2BlockEnd{}),
			},
		}

		streamServer.SetEntriesDef(entriesDefinition)

		s.updateDataStreamerFile(ctx, &streamServer)

		dbManager.streamServer = &streamServer
		err = dbManager.streamServer.Start()
		if err != nil {
			log.Fatalf("failed to start stream server, err: %v", err)
		}
	}

	go dbManager.Start()

	finalizer := newFinalizer(s.cfg.Finalizer, s.cfg.EffectiveGasPrice, worker, dbManager, s.state, s.address, s.isSynced, closingSignalCh, batchConstraints, s.eventLog)
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
	var currentL2Block uint64
	var currentTxIndex uint64
	var err error

	header := streamServer.GetHeader()

	if header.TotalEntries == 0 {
		// Get Genesis block
		genesisL2Block, err := s.state.GetDSGenesisBlock(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		err = streamServer.StartAtomicOp()
		if err != nil {
			log.Fatal(err)
		}

		genesisBlock := state.DSL2BlockStart{
			BatchNumber:    genesisL2Block.BatchNumber,
			L2BlockNumber:  genesisL2Block.L2BlockNumber,
			Timestamp:      genesisL2Block.Timestamp,
			GlobalExitRoot: genesisL2Block.GlobalExitRoot,
			Coinbase:       genesisL2Block.Coinbase,
			ForkID:         genesisL2Block.ForkID,
		}

		log.Infof("Genesis block: %+v", genesisBlock)

		_, err = streamServer.AddStreamEntry(1, genesisBlock.Encode())
		if err != nil {
			log.Fatal(err)
		}

		genesisBlockEnd := state.DSL2BlockEnd{
			L2BlockNumber: genesisL2Block.L2BlockNumber,
			BlockHash:     genesisL2Block.BlockHash,
			StateRoot:     genesisL2Block.StateRoot,
		}

		_, err = streamServer.AddStreamEntry(state.EntryTypeL2BlockEnd, genesisBlockEnd.Encode())
		if err != nil {
			log.Fatal(err)
		}

		err = streamServer.CommitAtomicOp()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		latestEntry, err := streamServer.GetEntry(header.TotalEntries - 1)
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("Latest entry: %+v", latestEntry)

		switch latestEntry.EntryType {
		case state.EntryTypeL2BlockStart:
			log.Info("Latest entry type is L2BlockStart")
			currentL2Block = binary.LittleEndian.Uint64(latestEntry.Data[8:16])
		case state.EntryTypeL2Tx:
			log.Info("Latest entry type is L2Tx")
			for latestEntry.EntryType == state.EntryTypeL2Tx {
				currentTxIndex++
				latestEntry, err = streamServer.GetEntry(header.TotalEntries - currentTxIndex)
				if err != nil {
					log.Fatal(err)
				}
			}
			if latestEntry.EntryType != state.EntryTypeL2BlockStart {
				log.Fatal("Latest entry is not a L2BlockStart")
			}
			currentL2Block = binary.LittleEndian.Uint64(latestEntry.Data[8:16])
		case state.EntryTypeL2BlockEnd:
			log.Info("Latest entry type is L2BlockEnd")
			currentL2Block = binary.LittleEndian.Uint64(latestEntry.Data[0:8])
		}
	}

	log.Infof("Current transaction index: %d", currentTxIndex)
	log.Infof("Current L2 block number: %d", currentL2Block)

	var limit uint64 = 1000
	var offset uint64 = currentL2Block
	var entry uint64 = header.TotalEntries
	var l2blocks []*state.DSL2Block

	if entry > 0 {
		entry--
	}

	for err == nil {
		log.Infof("Current entry number: %d", entry)

		l2blocks, err = s.state.GetDSL2Blocks(ctx, limit, offset, nil)
		offset += limit
		if len(l2blocks) == 0 {
			break
		}
		// Get transactions for all the retrieved l2 blocks
		l2Transactions, err := s.state.GetDSL2Transactions(ctx, l2blocks[0].L2BlockNumber, l2blocks[len(l2blocks)-1].L2BlockNumber, nil)
		if err != nil {
			log.Fatal(err)
		}

		err = streamServer.StartAtomicOp()
		if err != nil {
			log.Fatal(err)
		}

		for x, l2block := range l2blocks {
			if currentTxIndex > 0 {
				x += int(currentTxIndex)
				currentTxIndex = 0
			}

			blockStart := state.DSL2BlockStart{
				BatchNumber:    l2block.BatchNumber,
				L2BlockNumber:  l2block.L2BlockNumber,
				Timestamp:      l2block.Timestamp,
				GlobalExitRoot: l2block.GlobalExitRoot,
				Coinbase:       l2block.Coinbase,
				ForkID:         l2block.ForkID,
			}

			_, err = streamServer.AddStreamEntry(state.EntryTypeL2BlockStart, blockStart.Encode())
			if err != nil {
				log.Fatal(err)
			}

			entry, err = streamServer.AddStreamEntry(state.EntryTypeL2Tx, l2Transactions[x].Encode())
			if err != nil {
				log.Fatal(err)
			}

			blockEnd := state.DSL2BlockEnd{
				L2BlockNumber: l2block.L2BlockNumber,
				BlockHash:     l2block.BlockHash,
				StateRoot:     l2block.StateRoot,
			}

			_, err = streamServer.AddStreamEntry(state.EntryTypeL2BlockEnd, blockEnd.Encode())
			if err != nil {
				log.Fatal(err)
			}
		}
		err = streamServer.CommitAtomicOp()
		if err != nil {
			log.Fatal(err)
		}
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
		log.Infof("will try to delete %d redundant txs", len(txHashes))
		err = s.pool.DeleteTransactionsByHashes(ctx, txHashes)
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

func getMaxRemainingResources(constraints batchConstraints) state.BatchResources {
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
