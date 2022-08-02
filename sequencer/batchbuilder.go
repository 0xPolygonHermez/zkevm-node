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
)

func (s *Sequencer) tryToProcessTx(ctx context.Context, ticker *time.Ticker) {
	// Check if synchronizer is up to date
	if !s.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		waitTick(ctx, ticker)
		return
	}
	log.Info("synchronizer has synced last batch, checking if current sequence should be closed")

	// Check if should close sequence
	log.Infof("checking if current sequence should be closed")
	if s.shouldCloseSequenceInProgress(ctx) {
		log.Infof("current sequence should be closed")
		err := s.closeSequence(ctx)
		if err != nil {
			log.Errorf("error closing sequence: %v", err)
			return
		}
	}

	// Get next tx from the pool
	log.Info("getting pending tx from the pool")
	tx, err := s.pool.GetTopPendingTxByProfitabilityAndZkCounters(ctx, s.calculateZkCounters())
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

	processedTxs, unprocessedTxs := state.DetermineProcessedTransactions(processBatchResp.Responses)
	// only save in DB processed transactions.
	dbTx, err = s.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for StoreTransactions, err: %v", err)
		return
	}
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
		CumulativeGasUsed:    int64(s.cfg.MaxCumulativeGasUsed) - s.sequenceInProgress.CumulativeGasUsed,
		UsedKeccakHashes:     s.cfg.MaxKeccakHashes - s.sequenceInProgress.UsedKeccakHashes,
		UsedPoseidonHashes:   s.cfg.MaxPoseidonHashes - s.sequenceInProgress.UsedKeccakHashes,
		UsedPoseidonPaddings: s.cfg.MaxPoseidonPaddings - s.sequenceInProgress.UsedPoseidonPaddings,
		UsedMemAligns:        s.cfg.MaxMemAligns - s.sequenceInProgress.UsedMemAligns,
		UsedArithmetics:      s.cfg.MaxArithmetics - s.sequenceInProgress.UsedArithmetics,
		UsedBinaries:         s.cfg.MaxBinaries - s.sequenceInProgress.UsedBinaries,
		UsedSteps:            s.cfg.MaxSteps - s.sequenceInProgress.UsedSteps,
	}
}

func (s *Sequencer) closeSequence(ctx context.Context) error {
	newSequence, err := s.newSequence(ctx)
	if err != nil {
		return fmt.Errorf("failed to create new sequence, err: %v", err)
	}
	s.sequenceInProgress = newSequence
	return nil
}

// shouldCloseSequenceInProgress checks if sequence should be closed or not
// in case it's enough blocks since last GER update, long time since last batch and sequence is profitable
func (s *Sequencer) shouldCloseSequenceInProgress(ctx context.Context) bool {
	// Check if GER needs to be updated
	numberOfBlocks, err := s.state.GetNumberOfBlocksSinceLastGERUpdate(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get last time GER updated, err: %v", err)
		return false
	}
	if numberOfBlocks >= s.cfg.WaitBlocksToUpdateGER {
		if len(s.sequenceInProgress.Txs) == 0 {
			log.Warn("TODO: update GER without closing batch as no txs have been added yet")
			return false
		}
		isProfitable := s.isSequenceProfitable(ctx)
		if isProfitable {
			log.Infof("current sequence should be closed because %d blocks have been mined since last GER and tx is profitable", numberOfBlocks)
			return true
		}
	}
	// Check if it has been to long since a batch is virtualized
	lastBatchTime, err := s.state.GetLastBatchTime(ctx, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		log.Errorf("failed to get last batch time, err: %v", err)
		return false
	}
	if lastBatchTime.Before(time.Now().Add(-s.cfg.LastTimeBatchMaxWaitPeriod.Duration)) && len(s.sequenceInProgress.Txs) > 0 {
		isProfitable := s.isSequenceProfitable(ctx)
		if isProfitable {
			log.Info(
				"current sequence should be closed because LastTimeBatchMaxWaitPeriod has been exceeded, " +
					"there are pending sequences to be sent and they are profitable")
			return true
		}
	}
	// Check ZK counters
	zkCounters := s.calculateZkCounters()
	if zkCounters.IsZkCountersBelowZero() && len(s.sequenceInProgress.Txs) != 0 {
		log.Info("closing sequence because at least some ZK counter is bellow 0")
		return true
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
