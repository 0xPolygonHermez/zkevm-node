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
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	errGasRequiredExceedsAllowance = "gas required exceeds allowance"
	errContentLengthTooLarge       = "content length too large"
	errTimestampMustBeInsideRange  = "Timestamp must be inside range"
	errInsufficientAllowance       = "insufficient allowance"
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg Config

	pool      txPool
	state     stateInterface
	txManager txManager
	etherman  etherman
	checker   *profitabilitychecker.Checker
	gpe       gasPriceEstimator

	address          common.Address
	isSequenceTooBig bool

	sequenceInProgress types.Sequence
}

// New init sequencer
func New(
	cfg Config,
	txPool txPool,
	state stateInterface,
	etherman etherman,
	priceGetter priceGetter,
	manager txManager,
	gpe gasPriceEstimator) (*Sequencer, error) {
	checker := profitabilitychecker.New(cfg.ProfitabilityChecker, etherman, priceGetter)

	addr, err := etherman.TrustedSequencer()
	if err != nil {
		return nil, fmt.Errorf("failed to get trusted sequencer address, err: %v", err)
	}
	// TODO: check that private key used in etherman matches addr

	return &Sequencer{
		cfg:       cfg,
		pool:      txPool,
		state:     state,
		etherman:  etherman,
		checker:   checker,
		txManager: manager,
		gpe:       gpe,
		address:   addr,
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
		s.createFirstBatch(ctx)
	} else {
		err = s.loadSequenceFromState(ctx)
		if err != nil {
			log.Fatalf("failed to load sequence from the state, err: %v", err)
		}
	}

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
	lastBatch, err := s.state.GetLastBatch(ctx, nil)
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
		ger, err := s.getLatestGer(ctx, dbTx)
		if err != nil {
			if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
				return fmt.Errorf(
					"failed to rollback dbTx when getting last GER that gave err: %s. Rollback err: %s",
					rollbackErr.Error(), err.Error(),
				)
			}
			return fmt.Errorf("failed to get latest global exit root, err: %w", err)
		}
		processingCtx := state.ProcessingContext{
			BatchNumber:    lastBatch.BatchNumber + 1,
			Coinbase:       s.address,
			Timestamp:      time.Now(),
			GlobalExitRoot: ger,
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
			GlobalExitRoot: processingCtx.GlobalExitRoot,
			Timestamp:      processingCtx.Timestamp.Unix(),
		}
	} else {
		txs, err := s.state.GetTransactionsByBatchNumber(ctx, lastBatch.BatchNumber, nil)
		if err != nil {
			return fmt.Errorf("failed to get tx by batch number, err: %w", err)
		}
		s.sequenceInProgress = types.Sequence{
			GlobalExitRoot: lastBatch.GlobalExitRoot,
			Timestamp:      lastBatch.Timestamp.Unix(),
			Txs:            txs,
		}
		// TODO: execute to get state root and LER or change open/closed logic so we always store state root and LER and add an open flag
	}

	return nil
	/*
		TODO: deal with ongoing L1 txs
	*/
}

func (s *Sequencer) createFirstBatch(ctx context.Context) {
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
	s.sequenceInProgress = types.Sequence{
		GlobalExitRoot: processingCtx.GlobalExitRoot,
		Timestamp:      processingCtx.Timestamp.Unix(),
		Txs:            []ethTypes.Transaction{},
	}
}
