package sequencer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/profitabilitychecker"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

const (
	errGasRequiredExceedsAllowance = "gas required exceeds allowance"
	errContentLengthTooLarge       = "content length too large"
	errTimestampMustBeInsideRange  = "Timestamp must be inside range"
	errInsuficientAllowance        = "insufficient allowance"
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg Config

	pool                  txPool
	state                 stateInterface
	txManager             txManager
	etherman              etherman
	checker               *profitabilitychecker.Checker
	reorgTrustedStateChan chan struct{}

	address                          common.Address
	lastBatchNum                     uint64
	lastBatchNumSentToL1             uint64
	lastStateRoot, lastLocalExitRoot common.Hash

	sequenceInProgress types.Sequence
}

// New init sequencer
func New(
	cfg Config,
	pool txPool,
	state stateInterface,
	etherman etherman,
	priceGetter priceGetter,
	reorgTrustedStateChan chan struct{},
	manager txManager) (*Sequencer, error) {
	checker := profitabilitychecker.New(cfg.ProfitabilityChecker, etherman, priceGetter)

	addr, err := etherman.TrustedSequencer()
	if err != nil {
		return nil, fmt.Errorf("failed to get trusted sequencer address, err: %v", err)
	}
	// TODO: check that private key used in etherman matches addr

	return &Sequencer{
		cfg:                   cfg,
		pool:                  pool,
		state:                 state,
		etherman:              etherman,
		checker:               checker,
		txManager:             manager,
		address:               addr,
		reorgTrustedStateChan: reorgTrustedStateChan,
	}, nil
}

// Start starts the sequencer
func (s *Sequencer) Start(ctx context.Context) {
	for !s.isSynced(ctx) {
		log.Infof("waiting for synchronizer to sync...")
		time.Sleep(s.cfg.WaitPeriodPoolIsEmpty.Duration)
	}
	// initialize sequence
	batchNum, err := s.state.GetLastBatchNumber(ctx, nil)
	for err != nil {
		if errors.Is(err, state.ErrStateNotSynchronized) {
			log.Warnf("state is not synchronized, trying to get last batch num once again...")
			time.Sleep(s.cfg.WaitPeriodPoolIsEmpty.Duration)
			batchNum, err = s.state.GetLastBatchNumber(ctx, nil)
		} else {
			log.Fatalf("failed to get last batch number, err: %v", err)
		}
	}
	// case A: genesis
	if batchNum == 0 {
		log.Infof("starting sequencer with genesis batch")
		processingCtx := state.ProcessingContext{
			BatchNumber:    1,
			Coinbase:       s.address,
			Timestamp:      time.Now(),
			GlobalExitRoot: state.ZeroHash,
		}
		dbTx, err := s.state.BeginStateTransaction(ctx)
		if err != nil {
			log.Fatalf("failed to begin state transaction for opening a batch, err: %v", err)
		}
		err = s.state.OpenBatch(ctx, processingCtx, dbTx)
		if err != nil {
			if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
				log.Fatalf(
					"failed to rollback dbTx when opening batch that gave err: %v. Rollback err: %v",
					rollbackErr, err,
				)
			}
			log.Fatalf("failed to open a batch, err: %v", err)
		}
		if err := dbTx.Commit(ctx); err != nil {
			log.Fatalf("failed to commit dbTx when opening batch, err: %v", err)
		}
		s.lastBatchNum = processingCtx.BatchNumber
		s.sequenceInProgress = types.Sequence{
			GlobalExitRoot:  processingCtx.GlobalExitRoot,
			Timestamp:       processingCtx.Timestamp.Unix(),
			ForceBatchesNum: 0,
			Txs:             nil,
		}
	} else {
		err = s.loadSequenceFromState(ctx)
		if err != nil {
			log.Fatalf("failed to load sequence from the state, err: %v", err)
		}
	}

	go s.trackReorg(ctx)
	go s.trackOldTxs(ctx)
	tickerProcessTxs := time.NewTicker(s.cfg.WaitPeriodPoolIsEmpty.Duration)
	tickerSendSequence := time.NewTicker(s.cfg.WaitPeriodSendSequence.Duration)
	defer tickerProcessTxs.Stop()
	defer tickerSendSequence.Stop()
	go func() {
		for {
			s.tryToProcessTx(ctx, tickerProcessTxs)
		}
	}()
	go func() {
		for {
			s.tryToSendSequence(ctx, tickerSendSequence)
		}
	}()
	// Wait until context is done
	<-ctx.Done()
}

func (s *Sequencer) trackReorg(ctx context.Context) {
	for {
		select {
		case <-s.reorgTrustedStateChan:
			const waitTime = 5 * time.Second

			err := s.pool.MarkReorgedTxsAsPending(ctx)
			for err != nil {
				time.Sleep(waitTime)
				log.Errorf("failed to mark reorged txs as pending")
				err = s.pool.MarkReorgedTxsAsPending(ctx)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *Sequencer) trackOldTxs(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.FrequencyToCheckTxsForDelete.Duration)
	for {
		waitTick(ctx, ticker)
		txHashes, err := s.state.GetTxsOlderThanNL1Blocks(ctx, s.cfg.BlocksAmountForTxsToBeDeleted, nil)
		if err != nil {
			log.Errorf("failed to get txs hashes to delete, err: %v", err)
			continue
		}
		err = s.pool.DeleteTxsByHashes(ctx, txHashes)
		if err != nil {
			log.Errorf("failed to delete txs from the pool, err: %v", err)
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
	lastSyncedBatchNum, err := s.state.GetLastVirtualBatchNum(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get last synced batch, err: %v", err)
		return false
	}
	lastEthBatchNum, err := s.etherman.GetLatestBatchNumber()
	if err != nil {
		log.Errorf("failed to get last eth batch, err: %v", err)
		return false
	}
	if lastSyncedBatchNum < lastEthBatchNum {
		log.Infof("waiting for the state to be synced, lastSyncedBatchNum: %d, lastEthBatchNum: %d", lastSyncedBatchNum, lastEthBatchNum)
		return false
	}
	return true
}

func (s *Sequencer) loadSequenceFromState(ctx context.Context) error {
	// WIP
	lastBatch, err := s.state.GetLastBatch(ctx, nil)
	if err != nil {
		return err
	}
	s.lastBatchNum = lastBatch.BatchNumber
	s.lastStateRoot = lastBatch.StateRoot
	s.lastLocalExitRoot = lastBatch.LocalExitRoot
	// s.lastBatchNumSentToL1 =
	return fmt.Errorf("NOT IMPLEMENTED: loadSequenceFromState")
	/*
		TODO: set s.[lastBatchNum, lastStateRoot, lastLocalExitRoot, closedSequences, sequenceInProgress]
		based on stateDB data AND potentially pending txs to be mined on Ethereum, as this function may be called either
		when starting the sequencer OR if there is a mismatch between state data and on memory
	*/
}
