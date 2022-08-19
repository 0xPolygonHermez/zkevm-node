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
		if err != nil {
			log.Errorf("error closing sequence: %v", err)
			log.Info("resetting sequence in progress")
			if err = s.loadSequenceFromState(ctx); err != nil {
				log.Error("error loading sequence from state: %v", err)
			}
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
	previousStateRoot, err := s.state.GetStateRootByBatchNumber(ctx, s.lastBatchNum-1, nil)
	if err != nil {
		log.Errorf("failed to get state root for batchNum %d, err: %v", s.lastBatchNum, err)
		return
	}

	processBatchResp, err := s.state.ProcessSequencerBatch(ctx, previousStateRoot, s.lastBatchNum, s.sequenceInProgress.Txs, dbTx)
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
	for _, tx := range processedTxs {
		log.Debugf("tx: %v\n", tx.Tx.Hash())
		log.Debugf("tx nonce: %v\n", tx.Tx.Nonce())
	}
	for _, tx := range unprocessedTxs {
		log.Debugf("tx: %v\n", tx.Tx.Hash())
		log.Debugf("tx nonce: %v\n", tx.Tx.Nonce())
		log.Debugf("tx error: %v\n", tx.Error)
	}
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
			if err != nil {
				log.Errorf("failed to load sequence from state, err: %v", err)
			}
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
		s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:len(s.sequenceInProgress.Txs)-1]
		// in this case tx is invalid
		if len(s.sequenceInProgress.Txs) == 0 {
			txState = pool.TxStateInvalid
			txUpdateMsg = fmt.Sprintf("Tx %q failed to be processed. Marking tx as invalid", tx.Hash())
		} else {
			// otherwise close batch and put tx as pending
			txState = pool.TxStatePending
			txUpdateMsg = fmt.Sprintf("Tx %q failed to be processed. Marking tx as pending to return the pool", tx.Hash())

			log.Infof("current sequence should be closed, so tx with hash %q can be processed", tx.Hash())
			err := s.closeSequence(ctx)
			if err != nil {
				log.Errorf("error closing sequence: %v", err)
				log.Info("resetting sequence in progress")
				if err = s.loadSequenceFromState(ctx); err != nil {
					log.Error("error loading sequence from state: %v", err)
				}
			}
		}
	}
	log.Infof(txUpdateMsg)
	if err := s.pool.UpdateTxState(ctx, tx.Hash(), txState); err != nil {
		log.Errorf("failed to update tx status on the pool, err: %v", err)
		return
	}
}

func (s *Sequencer) newSequence(ctx context.Context) (types.Sequence, error) {
	var dbTx pgx.Tx
	var err error
	if s.lastStateRoot.String() != "" || s.lastLocalExitRoot.String() != "" {
		receipt := state.ProcessingReceipt{
			BatchNumber:   s.lastBatchNum,
			StateRoot:     s.lastStateRoot,
			LocalExitRoot: s.lastLocalExitRoot,
		}
		dbTx, err = s.state.BeginStateTransaction(ctx)
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
	} else {
		return types.Sequence{}, errors.New("lastStateRoot and lastLocalExitRoot are empty, impossible to close a batch")
	}
	// open next batch
	var gerHash common.Hash
	ger, err := s.state.GetLatestGlobalExitRoot(ctx, nil)
	if err != nil && err == state.ErrNotFound {
		gerHash = state.ZeroHash
	} else if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return types.Sequence{}, fmt.Errorf(
				"failed to rollback dbTx when getting last GER that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		return types.Sequence{}, fmt.Errorf("failed to get latest global exit root, err: %v", err)
	} else {
		gerHash = ger.GlobalExitRoot
	}

	lastBatchNum, err := s.state.GetLastBatchNumber(ctx, nil)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return types.Sequence{}, fmt.Errorf(
				"failed to rollback dbTx when getting last batch num that gave err: %v. Rollback err: %v",
				rollbackErr, err,
			)
		}
		return types.Sequence{}, fmt.Errorf("failed to get last batch number, err: %v", err)
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
	blockNum, mainnetExitRoot, err := s.state.GetBlockNumAndMainnetExitRootByGER(ctx, s.sequenceInProgress.GlobalExitRoot, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get mainnetExitRoot and blockNum by ger, err: %v", err)
		return false
	}

	lastGer, err := s.state.GetLatestGlobalExitRoot(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get latest global exit root, err: %v", err)
		return false
	}

	if lastGer != nil && lastGer.MainnetExitRoot != mainnetExitRoot {
		latestBlockNumber, err := s.etherman.GetLatestBlockNumber(ctx)
		if err != nil {
			log.Errorf("failed to get latest batch number from ethereum, err: %v", err)
			return false
		}
		if latestBlockNumber-blockNum > s.cfg.WaitBlocksToUpdateGER {
			if len(s.sequenceInProgress.Txs) == 0 {
				log.Info("update GER without closing batch as no txs have been added yet")

				dbTx, err := s.state.BeginStateTransaction(ctx)
				if err != nil {
					log.Errorf("failed to begin state transaction for UpdateGERInOpenBatch tx, err: %v", err)
					return false
				}

				err = s.state.UpdateGERInOpenBatch(ctx, lastGer.GlobalExitRoot, dbTx)
				if err != nil {
					if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
						log.Errorf(
							"failed to rollback dbTx when UpdateGERInOpenBatch that gave err: %v. Rollback err: %v",
							rollbackErr, err,
						)
						return false
					}
					log.Errorf("failed to update ger in open batch, err: %v", err)
					return false
				}

				if err := dbTx.Commit(ctx); err != nil {
					log.Errorf("failed to commit dbTx when processing UpdateGERInOpenBatch, err: %v", err)
					return false
				}
			} else {
				isProfitable := s.isSequenceProfitable(ctx)
				if isProfitable {
					log.Infof("current sequence should be closed because blocks have been mined since last GER and tx is profitable")
					return true
				}
			}
		}
	}
	// Check if it has been too long since a batch is virtualized
	isPreviousBatchVirtualized, err := s.state.IsBatchVirtualized(ctx, s.lastBatchNum-1, nil)
	if err != nil {
		log.Errorf("failed to get last virtual batch num, err: %v", err)
		return false
	}
	if time.Unix(s.sequenceInProgress.Timestamp, 0).Add(s.cfg.MaxTimeForBatchToBeOpen.Duration).Before(time.Now()) &&
		isPreviousBatchVirtualized && len(s.sequenceInProgress.Txs) > 0 {
		log.Info("closing batch because there are enough time to close a batch, previous batch is virtualized and batch has txs")
		return true
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
