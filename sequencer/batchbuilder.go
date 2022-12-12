package sequencer

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

const (
	getTxsLimit = uint64(1)
)

type processTxResponse struct {
	processedTxs         []*state.ProcessTransactionResponse
	processedTxsHashes   []string
	unprocessedTxs       map[string]*state.ProcessTransactionResponse
	unprocessedTxsHashes []string
	isBatchProcessed     bool
}

func (s *Sequencer) tryToProcessTx(ctx context.Context, ticker *time.Ticker) {
	start := time.Now()
	defer s.observeProcessingTime(start)

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
			log.Errorf("error closing sequence: %w", err)
			log.Info("resetting sequence in progress")
			if err = s.loadSequenceFromState(ctx); err != nil {
				log.Errorf("error loading sequence from state: %w", err)
			}
			return
		}
	}

	// backup current sequence
	sequenceBeforeTryingToProcessNewTxs := s.backupSequence()

	minGasPrice, err := s.gpe.GetAvgGasPrice(ctx)
	metrics.AverageGasPrice(float64(minGasPrice.Uint64()))
	if err != nil {
		log.Errorf("failed to get avg gas price, err: %w", err)
		return
	}

	// get txs from the pool
	txAmount := s.appendPendingTxs(ctx, true, 0, getTxsLimit, ticker)
	if txAmount == 0 {
		txAmount = s.appendPendingTxs(ctx, false, minGasPrice.Uint64(), getTxsLimit, ticker)
	}

	if txAmount == 0 {
		return
	}
	// clear txs if it bigger than expected
	err = s.cleanTxsIfTxsDataIsBiggerThanExpected(ctx, ticker)
	if err != nil {
		return
	}
	// process batch
	log.Infof("processing tx number [%d] for current batch...", len(s.sequenceInProgress.Txs))
	processResponse, err := s.processTxs(ctx)
	if err != nil {
		s.sequenceInProgress = sequenceBeforeTryingToProcessNewTxs
		log.Errorf("failed to process tx number: [%d], err: %w", err)
		return
	}
	if len(processResponse.unprocessedTxs) > 0 {
		// removing last tx as it failed processing
		s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:len(s.sequenceInProgress.Txs)-1]
	}
	log.Infof("tx number [%d] processed successfully", len(processResponse.processedTxsHashes))

	// If after processing new txs the sequence is equal or smaller, revert changes and close sequence
	if len(s.sequenceInProgress.Txs) <= len(sequenceBeforeTryingToProcessNewTxs.Txs) && len(s.sequenceInProgress.Txs) > 0 {
		log.Infof(
			"current sequence should be closed because after trying to add txs to it, it went from having %d valid txs to %d",
			len(sequenceBeforeTryingToProcessNewTxs.Txs), len(s.sequenceInProgress.Txs),
		)
		s.sequenceInProgress = sequenceBeforeTryingToProcessNewTxs
		if err := s.closeSequence(ctx); err != nil {
			log.Errorf("error closing sequence: %w", err)
		}
		return
	}

	if len(processResponse.processedTxs) > 0 {
		// only save in DB processed transactions.
		err = s.storeProcessedTransactions(ctx, processResponse.processedTxs)
		if err != nil {
			s.sequenceInProgress = sequenceBeforeTryingToProcessNewTxs
			log.Errorf("failed to store processed txs, err: %w", err)
			return
		}
		log.Infof("tx number [%d] stored and added into the trusted state", len(s.sequenceInProgress.Txs))
	}

	s.updateTxsInPool(ctx, ticker, processResponse, processResponse.unprocessedTxs)
}

func (s *Sequencer) observeProcessingTime(start time.Time) {
	elapsed := time.Since(start)
	metrics.ProcessingTime(elapsed)
}

func (s *Sequencer) updateTxsInPool(
	ctx context.Context,
	ticker *time.Ticker,
	processResponse processTxResponse,
	unprocessedTxs map[string]*state.ProcessTransactionResponse,
) {
	invalidTxsHashes, failedTxsHashes := s.splitInvalidAndFailedTxs(ctx, unprocessedTxs, ticker)

	metrics.TxProcessed(metrics.TxProcessedLabelSuccessful, float64(len(processResponse.processedTxsHashes)))
	metrics.TxProcessed(metrics.TxProcessedLabelInvalid, float64(len(invalidTxsHashes)))
	metrics.TxProcessed(metrics.TxProcessedLabelFailed, float64(len(failedTxsHashes)))

	// update processed txs
	s.updateTxsStatus(ctx, ticker, processResponse.processedTxsHashes, pool.TxStatusSelected)
	// update invalid txs
	s.updateTxsStatus(ctx, ticker, invalidTxsHashes, pool.TxStatusInvalid)
	// update failed txs
	s.updateTxsStatus(ctx, ticker, failedTxsHashes, pool.TxStatusFailed)
	// increment counter for failed txs
	s.incrementFailedCounter(ctx, ticker, failedTxsHashes)
}

func (s *Sequencer) cleanTxsIfTxsDataIsBiggerThanExpected(ctx context.Context, ticker *time.Ticker) error {
	encodedTxsBytesSize := math.MaxInt
	numberOfTxsInProcess := len(s.sequenceInProgress.Txs)
	for encodedTxsBytesSize > s.cfg.MaxBatchBytesSize && numberOfTxsInProcess > 0 {
		encodedTxs, err := state.EncodeTransactions(s.sequenceInProgress.Txs)
		if err != nil {
			log.Errorf("failed to encode txs, err: %w", err)
			return err
		}
		encodedTxsBytesSize = len(encodedTxs)
		if encodedTxsBytesSize > s.cfg.MaxBatchBytesSize && numberOfTxsInProcess > 0 {
			// if only one tx overflows, that it means, tx is invalid
			if numberOfTxsInProcess == 1 {
				err = s.pool.UpdateTxStatus(ctx, s.sequenceInProgress.Txs[0].Hash(), pool.TxStatusInvalid)
				for err != nil {
					log.Errorf("failed to update tx with hash: %s to status: %s",
						s.sequenceInProgress.Txs[0].Hash().String(), pool.TxStatusInvalid)
					err = s.pool.UpdateTxStatus(ctx, s.sequenceInProgress.Txs[0].Hash(), pool.TxStatusInvalid)
					waitTick(ctx, ticker)
				}
			}
			log.Infof("decreasing amount of sent txs, bcs encodedTxsBytesSize > maxBatchBytesSize, encodedTxsBytesSize: %d, maxBatchBytesSize: %d",
				encodedTxsBytesSize, s.cfg.MaxBatchBytesSize)
			s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:numberOfTxsInProcess-1]
			updatedNumberTxsInProgress := len(s.sequenceInProgress.Txs)
			if updatedNumberTxsInProgress != 0 {
				s.sequenceInProgress.IsSequenceTooBig = true
			}
			numberOfTxsInProcess = updatedNumberTxsInProgress
		}
	}
	return nil
}

func (s *Sequencer) splitInvalidAndFailedTxs(ctx context.Context, unprocessedTxs map[string]*state.ProcessTransactionResponse, ticker *time.Ticker) ([]string, []string) {
	invalidTxsHashes := []string{}
	failedTxsHashes := []string{}
	for _, tx := range unprocessedTxs {
		isTxNonceLessThanAccountNonce, err := s.isTxNonceLessThanAccountNonce(ctx, tx)
		for err != nil {
			log.Errorf("failed to compare account nonce and tx nonce, err: %w", err)
			isTxNonceLessThanAccountNonce, err = s.isTxNonceLessThanAccountNonce(ctx, tx)
			waitTick(ctx, ticker)
		}
		if isTxNonceLessThanAccountNonce {
			log.Infof("tx with hash %s is invalid, account nonce > tx nonce")
			invalidTxsHashes = append(invalidTxsHashes, tx.Tx.Hash().String())
		} else {
			failedTxsHashes = append(failedTxsHashes, tx.Tx.Hash().String())
		}
	}

	return invalidTxsHashes, failedTxsHashes
}

func (s *Sequencer) updateTxsStatus(ctx context.Context, ticker *time.Ticker, hashes []string, status pool.TxStatus) {
	if len(hashes) == 0 {
		return
	}
	err := s.pool.UpdateTxsStatus(ctx, hashes, status)
	for err != nil {
		log.Errorf("failed to update txs status to %s, err: %w", status, err)
		waitTick(ctx, ticker)
		err = s.pool.UpdateTxsStatus(ctx, hashes, status)
	}
}

func (s *Sequencer) incrementFailedCounter(ctx context.Context, ticker *time.Ticker, hashes []string) {
	if len(hashes) == 0 {
		return
	}
	err := s.pool.IncrementFailedCounter(ctx, hashes)
	for err != nil {
		log.Errorf("failed to increment failed tx counter, err: %w", err)
		waitTick(ctx, ticker)
		err = s.pool.IncrementFailedCounter(ctx, hashes)
	}
}

func (s *Sequencer) isTxNonceLessThanAccountNonce(ctx context.Context, tx *state.ProcessTransactionResponse) (bool, error) {
	fromAddr, txNonce, err := s.pool.GetTxFromAddressFromByHash(ctx, tx.Tx.Hash())
	if err != nil {
		return false, fmt.Errorf("failed to get from addr, err: %w", err)
	}

	lastL2BlockNumber, err := s.state.GetLastL2BlockNumber(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to get last l2 block number, err: %w", err)
	}

	accNonce, err := s.state.GetNonce(ctx, fromAddr, lastL2BlockNumber, nil)
	if err != nil {
		return false, fmt.Errorf("failed to get nonce for the account, err: %w", err)
	}

	return txNonce < accNonce, nil
}

func (s *Sequencer) newSequence(ctx context.Context) (types.Sequence, error) {
	var (
		dbTx pgx.Tx
		err  error
	)

	// It is necessary to pass the batch without txs to the executor in order to update the State
	if len(s.sequenceInProgress.Txs) == 0 {
		// backup current sequence
		sequenceBeforeTryingToProcessNewTxs := s.backupSequence()
		_, err = s.processTxs(ctx)
		for err != nil {
			s.sequenceInProgress = sequenceBeforeTryingToProcessNewTxs
			log.Errorf("failed to process txs, err: %w", err)
			_, err = s.processTxs(ctx)
		}
	}

	if s.sequenceInProgress.StateRoot.String() == "" || s.sequenceInProgress.LocalExitRoot.String() == "" {
		return types.Sequence{}, errors.New("state root and local exit root must have value to close batch")
	}
	dbTx, err = s.state.BeginStateTransaction(ctx)
	if err != nil {
		return types.Sequence{}, fmt.Errorf("failed to begin state transaction to close batch, err: %w", err)
	}

	lastBatchNumber, err := s.state.GetLastBatchNumber(ctx, dbTx)
	if err != nil {
		return types.Sequence{}, fmt.Errorf("failed to get last batch number, err: %w", err)
	}
	err = s.closeBatch(ctx, lastBatchNumber, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return types.Sequence{}, fmt.Errorf(
				"failed to rollback dbTx when closing batch that gave err: %s. Rollback err: %w",
				rollbackErr.Error(), err,
			)
		}
		return types.Sequence{}, err
	}
	// open next batch
	gerHash, _, err := s.getLatestGer(ctx, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return types.Sequence{}, fmt.Errorf(
				"failed to rollback dbTx when getting last GER that gave err: %s. Rollback err: %s",
				rollbackErr.Error(), err.Error(),
			)
		}
		return types.Sequence{}, err
	}

	processingCtx, err := s.openBatch(ctx, gerHash.GlobalExitRoot, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return types.Sequence{}, fmt.Errorf(
				"failed to rollback dbTx when getting last batch num that gave err: %s. Rollback err: %s",
				rollbackErr.Error(), err.Error(),
			)
		}
		return types.Sequence{}, err
	}
	if err := dbTx.Commit(ctx); err != nil {
		return types.Sequence{}, err
	}
	return types.Sequence{
		GlobalExitRoot: processingCtx.GlobalExitRoot,
		Timestamp:      processingCtx.Timestamp.Unix(),
		Txs:            []ethTypes.Transaction{},
	}, nil
}

func (s *Sequencer) closeSequence(ctx context.Context) error {
	newSequence, err := s.newSequence(ctx)
	if err != nil {
		return fmt.Errorf("failed to create new sequence, err: %w", err)
	}
	s.sequenceInProgress = newSequence
	return nil
}

func (s *Sequencer) processTxs(ctx context.Context) (processTxResponse, error) {
	dbTx, err := s.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for processing tx, err: %w", err)
		return processTxResponse{}, err
	}

	lastBatchNumber, err := s.state.GetLastBatchNumber(ctx, dbTx)
	if err != nil {
		log.Errorf("failed to get last batch number, err: %w", err)
		return processTxResponse{}, err
	}
	lenOfTxs := len(s.sequenceInProgress.Txs)
	isFirstTx := lenOfTxs == 1
	lastTx := s.sequenceInProgress.Txs[len(s.sequenceInProgress.Txs)-1]
	s.lastProcessBatchResp, err = s.state.ProcessSingleTransaction(ctx, lastBatchNumber, lastTx, isFirstTx, s.lastProcessBatchResp, dbTx, state.SequencerCallerLabel)
	if err != nil {
		if err == state.ErrBatchAlreadyClosed || err == state.ErrInvalidBatchNumber {
			log.Warnf("unexpected state local vs DB: %w", err)
			log.Info("reloading local sequence")
			errLoadSeq := s.loadSequenceFromState(ctx)
			if errLoadSeq != nil {
				log.Errorf("error loading sequence from state: %w", errLoadSeq)
			}
		}
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when processing tx that gave err: %w. Rollback err: %v",
				rollbackErr, err,
			)
			return processTxResponse{}, err
		}
		log.Errorf("failed processing batch, err: %w", err)
		return processTxResponse{}, err
	}

	if err := dbTx.Commit(ctx); err != nil {
		log.Errorf("failed to commit dbTx when processing tx, err: %w", err)
		return processTxResponse{}, err
	}

	s.sequenceInProgress.StateRoot = s.lastProcessBatchResp.NewStateRoot
	s.sequenceInProgress.LocalExitRoot = s.lastProcessBatchResp.NewLocalExitRoot

	processedTxs, processedTxsHashes, unprocessedTxs, unprocessedTxsHashes := state.DetermineProcessedTransactions(s.lastProcessBatchResp.Responses)

	response := processTxResponse{
		processedTxs:         processedTxs,
		processedTxsHashes:   processedTxsHashes,
		unprocessedTxs:       unprocessedTxs,
		unprocessedTxsHashes: unprocessedTxsHashes,
		isBatchProcessed:     s.lastProcessBatchResp.IsBatchProcessed,
	}

	return response, nil
}

func (s *Sequencer) storeProcessedTransactions(ctx context.Context, processedTxs []*state.ProcessTransactionResponse) error {
	dbTx, err := s.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for StoreTransactions, err: %w", err)
		return err
	}

	lastBatchNumber, err := s.state.GetLastBatchNumber(ctx, dbTx)
	if err != nil {
		log.Errorf("failed to get last batch number, err: %w", err)
		return err
	}

	err = s.state.StoreTransactions(ctx, lastBatchNumber, processedTxs, dbTx)
	if err != nil {
		s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:len(s.sequenceInProgress.Txs)-len(processedTxs)]
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when StoreTransactions that gave err: %w. Rollback err: %w",
				rollbackErr, err,
			)
			return err
		}
		log.Errorf("failed to store transactions, err: %w", err)
		if err == state.ErrOutOfOrderProcessedTx || err == state.ErrExistingTxGreaterThanProcessedTx {
			err = s.loadSequenceFromState(ctx)
			if err != nil {
				log.Errorf("failed to load sequence from state, err: %w", err)
			}
		}
		return err
	}

	if err := dbTx.Commit(ctx); err != nil {
		log.Errorf("failed to commit dbTx when StoreTransactions, err: %w", err)
		return err
	}

	return nil
}

func (s *Sequencer) closeBatch(ctx context.Context, lastBatchNumber uint64, dbTx pgx.Tx) error {
	receipt := state.ProcessingReceipt{
		BatchNumber:   lastBatchNumber,
		StateRoot:     s.sequenceInProgress.StateRoot,
		LocalExitRoot: s.sequenceInProgress.LocalExitRoot,
		Txs:           s.sequenceInProgress.Txs,
	}
	err := s.state.CloseBatch(ctx, receipt, dbTx)
	if err != nil {
		return fmt.Errorf("failed to close batch, err: %w", err)
	}

	return nil
}

func (s *Sequencer) getLatestGer(ctx context.Context, dbTx pgx.Tx) (state.GlobalExitRoot, time.Time, error) {
	lastBlockNumber, err := s.etherman.GetLatestBlockNumber(ctx)
	if err != nil {
		return state.GlobalExitRoot{}, time.Time{}, fmt.Errorf("failed to get latest eth block number, err: %w", err)
	}
	maxBlockNumber := lastBlockNumber - s.cfg.WaitBlocksToConsiderGerFinal
	ger, receivedAt, err := s.state.GetLatestGlobalExitRoot(ctx, maxBlockNumber, dbTx)
	if err != nil && errors.Is(err, state.ErrNotFound) {
		return state.GlobalExitRoot{}, time.Time{}, nil
	} else if err != nil {
		return state.GlobalExitRoot{}, time.Time{}, fmt.Errorf("failed to get latest global exit root, err: %w", err)
	} else {
		return ger, receivedAt, nil
	}
}

func (s *Sequencer) openBatch(ctx context.Context, gerHash common.Hash, dbTx pgx.Tx) (state.ProcessingContext, error) {
	lastBatchNum, err := s.state.GetLastBatchNumber(ctx, nil)
	if err != nil {
		return state.ProcessingContext{}, fmt.Errorf("failed to get last batch number, err: %w", err)
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
		return state.ProcessingContext{}, fmt.Errorf("failed to open new batch, err: %w", err)
	}

	return processingCtx, nil
}

func (s *Sequencer) appendPendingTxs(ctx context.Context, isClaims bool, minGasPrice, getTxsLimit uint64, ticker *time.Ticker) uint64 {
	pendTxs, err := s.pool.GetTxs(ctx, pool.TxStatusPending, isClaims, minGasPrice, getTxsLimit)
	if err == pgpoolstorage.ErrNotFound || len(pendTxs) == 0 {
		pendTxs, err = s.pool.GetTxs(ctx, pool.TxStatusFailed, isClaims, minGasPrice, getTxsLimit)
		if err == pgpoolstorage.ErrNotFound || len(pendTxs) == 0 {
			log.Infof("there is no suitable pending or failed txs in the pool, isClaims: %t, minGasPrice: %d, waiting...", isClaims, minGasPrice)
			if !isClaims {
				waitTick(ctx, ticker)
			}
			return 0
		}
	} else if err != nil {
		log.Errorf("failed to get pending tx, err: %w", err)
		return 0
	}
	var invalidTxsCounter int
	for i := 0; i < len(pendTxs); i++ {
		if pendTxs[i].FailedCounter > s.cfg.MaxAllowedFailedCounter {
			hash := pendTxs[i].Transaction.Hash().String()
			log.Warnf("mark tx with hash %s as invalid, failed counter %d exceeded max %d from config",
				hash, pendTxs[i].FailedCounter, s.cfg.MaxAllowedFailedCounter)
			s.updateTxsStatus(ctx, ticker, []string{hash}, pool.TxStatusInvalid)
			invalidTxsCounter++
			continue
		}
		s.sequenceInProgress.Txs = append(s.sequenceInProgress.Txs, pendTxs[i].Transaction)
	}

	return uint64(len(pendTxs) - invalidTxsCounter)
}

func (s *Sequencer) backupSequence() types.Sequence {
	backupSequence := types.Sequence{
		GlobalExitRoot: s.sequenceInProgress.GlobalExitRoot,
		StateRoot:      s.sequenceInProgress.StateRoot,
		LocalExitRoot:  s.sequenceInProgress.LocalExitRoot,
		Timestamp:      s.sequenceInProgress.Timestamp,
	}

	copy(backupSequence.Txs, s.sequenceInProgress.Txs)

	return backupSequence
}
