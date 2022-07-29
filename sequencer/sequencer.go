package sequencer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/profitabilitychecker"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
)

const (
	errGasRequiredExceedsAllowance = "gas required exceeds allowance"
	errContentLengthTooLarge       = "content length too large"
	errTimestampMustBeInsideRange  = "Timestamp must be inside range"
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
	lastStateRoot, lastLocalExitRoot common.Hash

	closedSequences    []types.Sequence
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
	go s.txManager.TrackSequenceBatchesSending(ctx)
	ticker := time.NewTicker(s.cfg.WaitPeriodPoolIsEmpty.Duration)
	defer ticker.Stop()
	for {
		s.tryToProcessTx(ctx, ticker)
	}
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

func (s *Sequencer) tryToProcessTx(ctx context.Context, ticker *time.Ticker) {
	if !s.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		waitTick(ctx, ticker)
		return
	}

	log.Info("synchronizer has synced last batch, checking if current sequence should be closed")
	if s.shouldCloseSequenceInProgress(ctx) && !s.closeSequence(ctx) {
		return
	}

	log.Infof("checking if current sequence should be sent")
	shouldSent, shouldCut := s.shouldSendSequences(ctx)
	if shouldSent {
		log.Infof("current sequence should be sent")
		if shouldCut {
			log.Infof("current sequence should be cut")
			cutSequence := s.closedSequences[len(s.closedSequences)-1]
			s.txManager.SequenceBatches(s.closedSequences)
			s.closedSequences = []types.Sequence{cutSequence}
		} else {
			s.txManager.SequenceBatches(s.closedSequences)
			s.closedSequences = []types.Sequence{}
		}
	}

	log.Info("getting pending tx from the pool")
	zkCounters := s.calculateZkCounters()
	if zkCounters.IsZkCountersBelowZero() {
		s.closeSequence(ctx)
		return
	}
	tx, err := s.pool.GetTopPendingTxByProfitabilityAndZkCounters(ctx, zkCounters)
	if err == pgpoolstorage.ErrNotFound {
		log.Infof("there is no suitable pending tx in the pool, waiting...")
		waitTick(ctx, ticker)
		return
	} else if err != nil {
		log.Errorf("failed to get pending tx, err: %v", err)
		return
	}

	log.Infof("processing tx: %s", tx.Hash())
	dbTx, err := s.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for processing tx, err: %v", err)
		return
	}

	s.sequenceInProgress.Txs = append(s.sequenceInProgress.Txs, tx.Transaction)
	processBatchResp, err := s.state.ProcessSequencerBatch(ctx, s.lastBatchNum, s.sequenceInProgress.Txs, dbTx)
	if err != nil {
		s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:len(s.sequenceInProgress.Txs)-1]
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when processing tx that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
			return
		}
		log.Debugf("failed to process tx, hash: %s, err: %v", tx.Hash(), err)
		return
	}

	if err := dbTx.Commit(ctx); err != nil {
		log.Errorf("failed to commit dbTx when processing tx, err: %v", err)
		return
	}

	s.sequenceInProgress.ZkCounters = pool.ZkCounters{
		CumulativeGasUsed:    int64(processBatchResp.CumulativeGasUsed),
		UsedKeccakHashes:     int32(processBatchResp.CntKeccakHashes),
		UsedPoseidonHashes:   int32(processBatchResp.CntPoseidonHashes),
		UsedPoseidonPaddings: int32(processBatchResp.CntPoseidonPaddings),
		UsedMemAligns:        int32(processBatchResp.CntMemAligns),
		UsedArithmetics:      int32(processBatchResp.CntArithmetics),
		UsedBinaries:         int32(processBatchResp.CntBinaries),
		UsedSteps:            int32(processBatchResp.CntSteps),
	}
	s.lastStateRoot = processBatchResp.NewStateRoot
	s.lastLocalExitRoot = processBatchResp.NewLocalExitRoot

	dbTx, err = s.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for StoreTransactions, err: %v", err)
		return
	}

	processedTxs, unprocessedTxs := state.DetermineProcessedTransactions(processBatchResp.Responses)
	// only save in DB processed transactions.
	err = s.state.StoreTransactions(ctx, s.lastBatchNum, processedTxs, dbTx)
	if err != nil {
		s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:len(s.sequenceInProgress.Txs)-1]
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when StoreTransactions that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
			return
		}
		log.Errorf("failed to store transactions, err: %v", err)
		if err == state.ErrOutOfOrderProcessedTx || err == state.ErrExistingTxGreaterThanProcessedTx {
			err = s.loadSequenceFromState(ctx)
			log.Errorf("failed to load sequence from state, err: %v", err)
		}
		return
	}

	if err := dbTx.Commit(ctx); err != nil {
		log.Errorf("failed to commit dbTx when StoreTransactions, err: %v", err)
		return
	}

	var txState = pool.TxStateSelected
	var txUpdateMsg = fmt.Sprintf("Tx %q added into the state. Marking tx as selected in the pool", tx.Hash())
	if _, ok := unprocessedTxs[tx.Hash().String()]; ok {
		txState = pool.TxStatePending
		txUpdateMsg = fmt.Sprintf("Tx %q failed to be processed. Marking tx as pending to return the pool", tx.Hash())
	}
	log.Infof(txUpdateMsg)
	if err := s.pool.UpdateTxState(ctx, tx.Hash(), txState); err != nil {
		log.Errorf("failed to update tx status on the pool, err: %v", err)
		return
	}
}

func (s *Sequencer) closeSequence(ctx context.Context) bool {
	log.Infof("current sequence should be closed")
	s.closedSequences = append(s.closedSequences, s.sequenceInProgress)
	newSequence, err := s.newSequence(ctx)
	if err != nil {
		log.Errorf("failed to create new sequence, err: %v", err)
		s.closedSequences = s.closedSequences[:len(s.closedSequences)-1]
		return false
	}
	s.sequenceInProgress = newSequence
	return true
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

// shouldSendSequences check if sequencer should send sequencer. Returns two bool vars -
// first bool is for should sequencer send sequences or not
// second bool is for should sequencer cut last sequences from sequences slice bcs data to send is too big
func (s *Sequencer) shouldSendSequences(ctx context.Context) (bool, bool) {
	estimatedGas, err := s.etherman.EstimateGasSequenceBatches(s.closedSequences)
	if err != nil && isDataForEthTxTooBig(err) {
		log.Warnf("closedSequences eth data is too big, err: %v", err)
		return true, true
	}

	if err != nil {
		// while estimating gas a new block is not created and the POE SC may return
		// an error regarding timestamp verification, this must be handled
		if strings.Contains(err.Error(), errTimestampMustBeInsideRange) {
			// query the sc about the value of its lastTimestamp variable
			lastTimestamp, err := s.etherman.GetLastTimestamp()
			if err != nil {
				log.Errorf("failed to query last timestamp from SC, err: %v", err)
				return false, false
			}
			// check POE SC lastTimestamp against sequences' one
			for _, seq := range s.closedSequences {
				if seq.Timestamp < int64(lastTimestamp) {
					log.Fatalf("sequence timestamp %d is < POE SC lastTimestamp %d", seq.Timestamp, lastTimestamp)
				}
				lastTimestamp = uint64(seq.Timestamp)
			}

			log.Debug("block.timestamp is greater than seq.Timestamp. A new block must be mined before the gas can be estimated.")
			return false, false
		}

		log.Errorf("failed to estimate gas for sequence batches, err: %v", err)
		return false, false
	}

	// TODO: checkAgainstForcedBatchQueueTimeout

	lastBatchVirtualizationTime, err := s.state.GetTimeForLatestBatchVirtualization(ctx, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		log.Errorf("failed to get last l1 interaction time, err: %v", err)
		return false, false
	}

	if lastBatchVirtualizationTime.Before(time.Now().Add(-s.cfg.LastBatchVirtualizationTimeMaxWaitPeriod.Duration)) {
		// check profitability
		if s.checker.IsSendSequencesProfitable(new(big.Int).SetUint64(estimatedGas), s.closedSequences) {
			return true, false
		}
	}

	return false, false
}

// shouldCloseSequenceInProgress checks if sequence should be closed or not
// in case it's enough blocks since last GER update, long time since last batch and sequence is profitable
func (s *Sequencer) shouldCloseSequenceInProgress(ctx context.Context) bool {
	numberOfBlocks, err := s.state.GetNumberOfBlocksSinceLastGERUpdate(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get last time GER updated, err: %v", err)
		return false
	}
	if numberOfBlocks >= s.cfg.WaitBlocksToUpdateGER {
		return s.isSequenceProfitable(ctx)
	}

	lastBatchTime, err := s.state.GetLastBatchTime(ctx, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		log.Errorf("failed to get last batch time, err: %v", err)
		return false
	}
	if lastBatchTime.Before(time.Now().Add(-s.cfg.LastTimeBatchMaxWaitPeriod.Duration)) && len(s.sequenceInProgress.Txs) > 0 {
		return s.isSequenceProfitable(ctx)
	}

	return false
}

func (s *Sequencer) isSequenceProfitable(ctx context.Context) bool {
	isProfitable, err := s.checker.IsSequenceProfitable(ctx, s.sequenceInProgress)
	if err != nil {
		log.Errorf("failed to check is sequence profitable, err: %v", err)
		return false
	}

	return isProfitable
}

func (s *Sequencer) newSequence(ctx context.Context) (types.Sequence, error) {
	if s.lastStateRoot.String() != "" || s.lastLocalExitRoot.String() != "" {
		receipt := state.ProcessingReceipt{
			BatchNumber:   s.lastBatchNum,
			StateRoot:     s.lastStateRoot,
			LocalExitRoot: s.lastLocalExitRoot,
		}
		dbTx, err := s.state.BeginStateTransaction(ctx)
		if err != nil {
			return types.Sequence{}, fmt.Errorf("failed to begin state transaction to close batch, err: %v", err)
		}
		err = s.state.CloseBatch(ctx, receipt, dbTx)
		if err != nil {
			if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
				return types.Sequence{}, fmt.Errorf(
					"failed to rollback dbTx when closing batch that gave err: %v. Rollback err: %v",
					rollbackErr, err,
				)
			}
			return types.Sequence{}, fmt.Errorf("failed to close batch, err: %v", err)
		}
		if err := dbTx.Commit(ctx); err != nil {
			return types.Sequence{}, fmt.Errorf("failed to commit dbTx when close batch, err: %v", err)
		}
	} else {
		return types.Sequence{}, errors.New("lastStateRoot and lastLocalExitRoot are empty, impossible to close a batch")
	}
	// open next batch
	var gerHash common.Hash
	ger, err := s.state.GetLatestGlobalExitRoot(ctx, nil)
	if err != nil && err == state.ErrNotFound {
		gerHash = state.ZeroHash
	} else if err != nil {
		return types.Sequence{}, fmt.Errorf("failed to get latest global exit root, err: %v", err)
	} else {
		gerHash = ger.GlobalExitRoot
	}

	lastBatchNum, err := s.state.GetLastBatchNumber(ctx, nil)
	if err != nil {
		return types.Sequence{}, fmt.Errorf("failed to get last batch number, err: %v", err)
	}
	newBatchNum := lastBatchNum + 1
	dbTx, err := s.state.BeginStateTransaction(ctx)
	if err != nil {
		return types.Sequence{}, fmt.Errorf("failed to open new batch, err: %v", err)
	}
	processingCtx := state.ProcessingContext{
		BatchNumber:    newBatchNum,
		Coinbase:       s.address,
		Timestamp:      time.Now(),
		GlobalExitRoot: gerHash,
	}
	err = s.state.OpenBatch(ctx, processingCtx, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return types.Sequence{}, fmt.Errorf(
				"failed to rollback dbTx when opening batch that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		return types.Sequence{}, fmt.Errorf("failed to open new batch, err: %v", err)
	}
	if err := dbTx.Commit(ctx); err != nil {
		return types.Sequence{}, fmt.Errorf("failed to commit dbTx when opening batch, err: %v", err)
	}

	s.lastBatchNum = newBatchNum
	return types.Sequence{
		GlobalExitRoot:  processingCtx.GlobalExitRoot,
		Timestamp:       processingCtx.Timestamp.Unix(),
		ForceBatchesNum: 0,
		Txs:             nil,
	}, nil
}

func (s *Sequencer) calculateZkCounters() pool.ZkCounters {
	return pool.ZkCounters{
		CumulativeGasUsed:    s.cfg.MaxGasUsed - s.sequenceInProgress.CumulativeGasUsed,
		UsedKeccakHashes:     s.cfg.MaxKeccakHashes - s.sequenceInProgress.UsedKeccakHashes,
		UsedPoseidonHashes:   s.cfg.MaxPoseidonHashes - s.sequenceInProgress.UsedKeccakHashes,
		UsedPoseidonPaddings: s.cfg.MaxPoseidonPaddings - s.sequenceInProgress.UsedPoseidonPaddings,
		UsedMemAligns:        s.cfg.MaxMemAligns - s.sequenceInProgress.UsedMemAligns,
		UsedArithmetics:      s.cfg.MaxArithmetics - s.sequenceInProgress.UsedArithmetics,
		UsedBinaries:         s.cfg.MaxBinaries - s.sequenceInProgress.UsedBinaries,
		UsedSteps:            s.cfg.MaxSteps - s.sequenceInProgress.UsedSteps,
	}
}

func isDataForEthTxTooBig(err error) bool {
	return strings.Contains(err.Error(), errGasRequiredExceedsAllowance) ||
		errors.Is(err, core.ErrOversizedData) ||
		strings.Contains(err.Error(), errContentLengthTooLarge)
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
	return fmt.Errorf("NOT IMPLEMENTED: loadSequenceFromState")
	/*
		TODO: set s.[lastBatchNum, lastStateRoot, lastLocalExitRoot, closedSequences, sequenceInProgress]
		based on stateDB data AND potentially pending txs to be mined on Ethereum, as this function may be called either
		when starting the sequencer OR if there is a mismatch between state data and on memory
	*/
}
