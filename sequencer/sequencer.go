package sequencer

import (
	"context"
	"errors"
	"fmt"
	"sync"
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
	dbManager dbManagerStateInterface
	txManager txManager
	etherman  etherman

	address common.Address
}

// L2ReorgEvent is the event that is triggered when a reorg happens in the L2
type L2ReorgEvent struct {
	TxHashes []common.Hash
}

// ClosingSignalCh is a struct that contains all the channels that are used to receive batch closing signals
type ClosingSignalCh struct {
	ForcedBatchCh        chan state.ForcedBatch
	GERCh                chan common.Hash
	L2ReorgCh            chan L2ReorgEvent
	SendingToL1TimeoutCh chan bool
}

// TxsStore is a struct that contains the channel and the wait group for the txs to be stored in order
type TxsStore struct {
	Ch chan *txToStore
	Wg *sync.WaitGroup
}

// txToStore represents a transaction to store.
type txToStore struct {
	txResponse               *state.ProcessTransactionResponse
	batchNumber              uint64
	previousL2BlockStateRoot common.Hash
}

// New init sequencer
func New(cfg Config, txPool txPool, state stateInterface, dbManager dbManagerStateInterface, etherman etherman, manager txManager) (*Sequencer, error) {
	addr, err := etherman.TrustedSequencer()
	if err != nil {
		return nil, fmt.Errorf("failed to get trusted sequencer address, err: %v", err)
	}

	return &Sequencer{
		cfg:       cfg,
		pool:      txPool,
		state:     state,
		dbManager: dbManager,
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

	closingSignalCh := ClosingSignalCh{
		ForcedBatchCh:        make(chan state.ForcedBatch),
		GERCh:                make(chan common.Hash),
		L2ReorgCh:            make(chan L2ReorgEvent),
		SendingToL1TimeoutCh: make(chan bool),
	}

	txsStore := TxsStore{
		Ch: make(chan *txToStore),
		Wg: new(sync.WaitGroup),
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

	batchResourceWeights := batchResourceWeights{
		WeightBatchBytesSize:    s.cfg.WeightBatchBytesSize,
		WeightCumulativeGasUsed: s.cfg.WeightCumulativeGasUsed,
		WeightKeccakHashes:      s.cfg.WeightKeccakHashes,
		WeightPoseidonHashes:    s.cfg.WeightPoseidonHashes,
		WeightPoseidonPaddings:  s.cfg.WeightPoseidonPaddings,
		WeightMemAligns:         s.cfg.WeightMemAligns,
		WeightArithmetics:       s.cfg.WeightArithmetics,
		WeightBinaries:          s.cfg.WeightBinaries,
		WeightSteps:             s.cfg.WeightSteps,
	}

	worker := NewWorker(s.cfg, s.state, batchConstraints, batchResourceWeights)
	dbManager := newDBManager(ctx, s.pool, s.dbManager, worker, closingSignalCh, txsStore)
	go dbManager.Start()

	finalizer := newFinalizer(s.cfg.Finalizer, worker, dbManager, s.state, s.address, s.isSynced, closingSignalCh, txsStore, batchConstraints)
	currBatch, OldAccInputHash, OldStateRoot := s.bootstrap(ctx, dbManager, finalizer)
	go finalizer.Start(ctx, currBatch, OldStateRoot, OldAccInputHash)

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

func (s *Sequencer) bootstrap(ctx context.Context, dbManager *dbManager, finalizer *finalizer) (*WipBatch, common.Hash, common.Hash) {
	var (
		currBatch                     *WipBatch
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
		currBatch = &WipBatch{
			globalExitRoot: processingCtx.GlobalExitRoot,
			batchNumber:    processingCtx.BatchNumber,
			coinbase:       processingCtx.Coinbase,
			timestamp:      uint64(processingCtx.Timestamp.Unix()),
			txs:            make([]TxTracker, 0, s.cfg.MaxTxsPerBatch),
		}

		if err != nil {
			return nil, common.Hash{}, common.Hash{}
		}
	} else {
		// Check if synchronizer is up-to-date
		for !s.isSynced(ctx) {
			log.Info("wait for synchronizer to sync last batch")
			time.Sleep(time.Second)
		}
		finalizer.closeAndOpenNewBatch(ctx)
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
