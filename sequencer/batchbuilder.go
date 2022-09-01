package sequencer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

func (s *Sequencer) tryToProcessTx(ctx context.Context, ticker *time.Ticker) {
	// Check if synchronizer is up to date
	if !s.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		waitTick(ctx, ticker)
		return
	}
	log.Info("synchronizer has synced last batch, checking if current sequence should be closed")

	// Check if sequence should be close
	log.Infof("checking if current sequence should be closed")
	if s.shouldCloseSequenceInProgress(ctx) {
		log.Infof("current sequence should be closed")
		err := s.closeSequence(ctx)
		if errors.As(err, &state.ErrClosingBatchWithoutTxs) {
			log.Info("current sequence can't be closed without transactions")
			waitTick(ctx, ticker)
			return
		} else if err != nil {
			log.Errorf("error closing sequence: %v", err)
			log.Info("resetting sequence in progress")
			if err = s.loadSequenceFromState(ctx); err != nil {
				log.Error("error loading sequence from state: %v", err)
			}
			return
		}
	}

	// Get next txs from the pool
	log.Info("getting pending txs from the pool")
	for !s.isZkCountersMoreThanMax(s.sumZkCounters) {
		pendTx, err := s.pool.GetTopPendingTxByProfitabilityAndZkCounters(ctx, s.remainingZkCounters(s.sumZkCounters), s.pendingTxsHashes)
		if err == pgpoolstorage.ErrNotFound {
			if len(s.pendingTxs) > 0 {
				log.Info("there is no suitable pending tx in the pool, proceed to process pending txs...")
				break
			}
			log.Infof("there is no suitable pending tx in the pool, waiting...")
			waitTick(ctx, ticker)
			continue
		} else if err != nil {
			log.Errorf("failed to get pending tx, err: %v", err)
			return
		}
		log.Infof("adding pending txs to pending tx array, hash: %s", pendTx.Hash().String())
		s.sumZkCounters.SumUpZkCounters(pendTx.ZkCounters)
		s.pendingTxs = append(s.pendingTxs, pendTx)
		s.pendingTxsHashes = append(s.pendingTxsHashes, pendTx.Hash().String())
	}

	for _, tx := range s.pendingTxs {
		log.Infof("processing tx: %s", tx.Hash())
	}
	processedTxs, processedTxsHashes, unprocessedTxs, err := s.processTxs(ctx, s.pendingTxs)
	if err != nil {
		log.Errorf("failed to process txs, err: %w", err)
		return
	}
	// only save in DB processed transactions.
	err = s.storeProcessedTransactions(ctx, processedTxs)
	if err != nil {
		log.Errorf("failed to store processed txs, err: %w", err)
		return
	}

	// update processed txs
	err = s.pool.UpdateTxsStatus(ctx, processedTxsHashes, pool.TxStatusSelected)
	for err != nil {
		log.Errorf("failed to update txs state to selected, err: %w")
		waitTick(ctx, ticker)
		err = s.pool.UpdateTxsStatus(ctx, processedTxsHashes, pool.TxStatusSelected)
	}
	// get rid of unprocessed txs in sequenceInProgressTxs slice
	var tempSequenceInProgressTxs []ethTypes.Transaction
	for _, tx := range s.sequenceInProgress.Txs {
		if _, ok := unprocessedTxs[tx.Hash().String()]; !ok {
			tempSequenceInProgressTxs = append(tempSequenceInProgressTxs, tx)
		}
	}

	s.sequenceInProgress.Txs = tempSequenceInProgressTxs
	s.pendingTxs = []*pool.Transaction{}
	s.pendingTxsHashes = []string{}
	s.sumZkCounters = pool.ZkCounters{}
}

func (s *Sequencer) newSequence(ctx context.Context) (types.Sequence, error) {
	var (
		dbTx pgx.Tx
		err  error
	)
	if s.lastStateRoot.String() != "" || s.lastLocalExitRoot.String() != "" {
		dbTx, err = s.state.BeginStateTransaction(ctx)
		if err != nil {
			return types.Sequence{}, fmt.Errorf("failed to begin state transaction to close batch, err: %v", err)
		}
		err = s.closeBatch(ctx, dbTx)
		if err != nil {
			return types.Sequence{}, err
		}
	} else {
		return types.Sequence{}, errors.New("lastStateRoot and lastLocalExitRoot are empty, impossible to close a batch")
	}
	// open next batch
	gerHash, err := s.getLatestGer(ctx, dbTx)
	if err != nil {
		return types.Sequence{}, err
	}

	processingCtx, err := s.openBatch(ctx, gerHash, dbTx)
	if err != nil {
		return types.Sequence{}, err
	}
	return types.Sequence{
		GlobalExitRoot:  processingCtx.GlobalExitRoot,
		Timestamp:       processingCtx.Timestamp.Unix(),
		ForceBatchesNum: 0,
		Txs:             nil,
	}, nil
}

// remainingZkCounters calculates difference between max values in cfg and passed zk counters
func (s *Sequencer) remainingZkCounters(zkCounters pool.ZkCounters) pool.ZkCounters {
	return pool.ZkCounters{
		CumulativeGasUsed:    int64(s.cfg.MaxCumulativeGasUsed) - zkCounters.CumulativeGasUsed,
		UsedKeccakHashes:     s.cfg.MaxKeccakHashes - zkCounters.UsedKeccakHashes,
		UsedPoseidonHashes:   s.cfg.MaxPoseidonHashes - zkCounters.UsedKeccakHashes,
		UsedPoseidonPaddings: s.cfg.MaxPoseidonPaddings - zkCounters.UsedPoseidonPaddings,
		UsedMemAligns:        s.cfg.MaxMemAligns - zkCounters.UsedMemAligns,
		UsedArithmetics:      s.cfg.MaxArithmetics - zkCounters.UsedArithmetics,
		UsedBinaries:         s.cfg.MaxBinaries - zkCounters.UsedBinaries,
		UsedSteps:            s.cfg.MaxSteps - zkCounters.UsedSteps,
	}
}

func (s *Sequencer) closeSequence(ctx context.Context) error {
	newSequence, err := s.newSequence(ctx)
	if err != nil {
		return err
	}
	s.sequenceInProgress = newSequence
	return nil
}

func (s *Sequencer) isSequenceProfitable(ctx context.Context) bool {
	isProfitable, err := s.checker.IsSequenceProfitable(ctx, s.sequenceInProgress)
	if err != nil {
		log.Errorf("failed to check is sequence profitable, err: %v", err)
		return false
	}

	return isProfitable
}

func (s *Sequencer) processTxs(ctx context.Context, pendingTxs []*pool.Transaction) (
	[]*state.ProcessTransactionResponse, []common.Hash, map[string]*state.ProcessTransactionResponse, error) {
	dbTx, err := s.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for processing tx, err: %v", err)
		return nil, nil, nil, err
	}

	for _, tx := range pendingTxs {
		s.sequenceInProgress.Txs = append(s.sequenceInProgress.Txs, tx.Transaction)
	}

	processBatchResp, err := s.state.ProcessSequencerBatch(ctx, s.lastBatchNum, s.sequenceInProgress.Txs, dbTx)
	if err != nil {
		s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:len(s.sequenceInProgress.Txs)-len(pendingTxs)]
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when processing tx that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
			return nil, nil, nil, err
		}
		for _, tx := range pendingTxs {
			log.Debugf("failed to process tx, hash: %s, err: %v", tx.Hash(), err)
		}
		return nil, nil, nil, err
	}

	if err := dbTx.Commit(ctx); err != nil {
		log.Errorf("failed to commit dbTx when processing tx, err: %v", err)
		return nil, nil, nil, err
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

	processedTxs, processedTxsHashes, unprocessedTxs := state.DetermineProcessedTransactions(processBatchResp.Responses)
	return processedTxs, processedTxsHashes, unprocessedTxs, nil
}

func (s *Sequencer) storeProcessedTransactions(ctx context.Context, processedTxs []*state.ProcessTransactionResponse) error {
	dbTx, err := s.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for StoreTransactions, err: %v", err)
		return err
	}
	err = s.state.StoreTransactions(ctx, s.lastBatchNum, processedTxs, dbTx)
	if err != nil {
		s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:len(s.sequenceInProgress.Txs)-len(processedTxs)]
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when StoreTransactions that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
			return err
		}
		log.Errorf("failed to store transactions, err: %v", err)
		if err == state.ErrOutOfOrderProcessedTx || err == state.ErrExistingTxGreaterThanProcessedTx {
			err = s.loadSequenceFromState(ctx)
			if err != nil {
				log.Errorf("failed to load sequence from state, err: %v", err)
			}
		}
		return err
	}

	if err := dbTx.Commit(ctx); err != nil {
		log.Errorf("failed to commit dbTx when StoreTransactions, err: %v", err)
		return err
	}

	return nil
}

func (s *Sequencer) updateGerInBatch(ctx context.Context, lastGer *state.GlobalExitRoot) error {
	log.Info("update GER without closing batch as no txs have been added yet")

	dbTx, err := s.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for UpdateGERInOpenBatch tx, err: %v", err)
		return err
	}

	err = s.state.UpdateGERInOpenBatch(ctx, lastGer.GlobalExitRoot, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when UpdateGERInOpenBatch that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
			return err
		}
		log.Errorf("failed to update ger in open batch, err: %v", err)
		return err
	}

	if err := dbTx.Commit(ctx); err != nil {
		log.Errorf("failed to commit dbTx when processing UpdateGERInOpenBatch, err: %v", err)
		return err
	}

	return nil
}

func (s *Sequencer) closeBatch(ctx context.Context, dbTx pgx.Tx) error {
	receipt := state.ProcessingReceipt{
		BatchNumber:   s.lastBatchNum,
		StateRoot:     s.lastStateRoot,
		LocalExitRoot: s.lastLocalExitRoot,
	}
	err := s.state.CloseBatch(ctx, receipt, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf(
				"failed to rollback dbTx when closing batch that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		return fmt.Errorf("failed to close batch, err: %v", err)
	}

	return nil
}

func (s *Sequencer) getLatestGer(ctx context.Context, dbTx pgx.Tx) (common.Hash, error) {
	ger, err := s.state.GetLatestGlobalExitRoot(ctx, dbTx)
	if err != nil && err == state.ErrNotFound {
		return state.ZeroHash, nil
	} else if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return common.Hash{}, fmt.Errorf(
				"failed to rollback dbTx when getting last GER that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		return common.Hash{}, fmt.Errorf("failed to get latest global exit root, err: %v", err)
	} else {
		return ger.GlobalExitRoot, nil
	}
}

func (s *Sequencer) openBatch(ctx context.Context, gerHash common.Hash, dbTx pgx.Tx) (state.ProcessingContext, error) {
	lastBatchNum, err := s.state.GetLastBatchNumber(ctx, nil)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return state.ProcessingContext{}, fmt.Errorf(
				"failed to rollback dbTx when getting last batch num that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		return state.ProcessingContext{}, fmt.Errorf("failed to get last batch number, err: %v", err)
	}
	newBatchNum := lastBatchNum + 1
	processingCtx := state.ProcessingContext{
		BatchNumber:    newBatchNum,
		Coinbase:       s.address,
		Timestamp:      time.Now(),
		GlobalExitRoot: gerHash,
	}
	err = s.state.OpenBatch(ctx, processingCtx, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return state.ProcessingContext{}, fmt.Errorf(
				"failed to rollback dbTx when opening batch that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		return state.ProcessingContext{}, fmt.Errorf("failed to open new batch, err: %v", err)
	}
	if err := dbTx.Commit(ctx); err != nil {
		return state.ProcessingContext{}, fmt.Errorf("failed to commit dbTx when opening batch, err: %v", err)
	}

	s.lastBatchNum = newBatchNum

	return processingCtx, nil
}

func (s *Sequencer) isZkCountersMoreThanMax(sumCounters pool.ZkCounters) bool {
	return s.cfg.MaxCumulativeGasUsed <= uint64(sumCounters.CumulativeGasUsed) ||
		s.cfg.MaxKeccakHashes <= sumCounters.UsedKeccakHashes ||
		s.cfg.MaxPoseidonHashes <= sumCounters.UsedPoseidonHashes ||
		s.cfg.MaxPoseidonPaddings <= sumCounters.UsedPoseidonPaddings ||
		s.cfg.MaxMemAligns <= sumCounters.UsedMemAligns ||
		s.cfg.MaxArithmetics <= sumCounters.UsedArithmetics ||
		s.cfg.MaxBinaries <= sumCounters.UsedBinaries ||
		s.cfg.MaxSteps <= sumCounters.UsedSteps
}
