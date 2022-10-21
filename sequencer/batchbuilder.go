package sequencer

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
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

const (
	maxTxsPerBatch    uint64 = 150
	maxBatchBytesSize int    = 30000
)

type processTxResponse struct {
	processedTxs         []*state.ProcessTransactionResponse
	processedTxsHashes   []string
	unprocessedTxs       map[string]*state.ProcessTransactionResponse
	unprocessedTxsHashes []string
	isBatchProcessed     bool
}

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
			if strings.Contains(err.Error(), state.ErrClosingBatchWithoutTxs.Error()) {
				log.Warn("Current batch has not been closed since it had no txs. Trying to add more txs to avoid death lock")
			} else {
				log.Errorf("error closing sequence: %w", err)
				log.Info("resetting sequence in progress")
				if err = s.loadSequenceFromState(ctx); err != nil {
					log.Errorf("error loading sequence from state: %w", err)
				}
				return
			}
		}
	}

	// backup current sequence
	sequenceBeforeTryingToProcessNewTxs := types.Sequence{
		GlobalExitRoot: s.sequenceInProgress.GlobalExitRoot,
		StateRoot:      s.sequenceInProgress.StateRoot,
		LocalExitRoot:  s.sequenceInProgress.LocalExitRoot,
		Timestamp:      s.sequenceInProgress.Timestamp,
	}
	copy(sequenceBeforeTryingToProcessNewTxs.Txs, s.sequenceInProgress.Txs)

	getTxsLimit := maxTxsPerBatch - uint64(len(s.sequenceInProgress.Txs))

	minGasPrice, err := s.gpe.GetAvgGasPrice(ctx)
	if err != nil {
		log.Errorf("failed to get avg gas price, err: %w", err)
		return
	}

	// get txs from the pool
	appendedClaimsTxsAmount := s.appendPendingTxs(ctx, true, 0, getTxsLimit, ticker)
	appendedTxsAmount := s.appendPendingTxs(ctx, false, minGasPrice.Uint64(), getTxsLimit-appendedClaimsTxsAmount, ticker) + appendedClaimsTxsAmount

	if appendedTxsAmount == 0 {
		return
	}
	// clear txs if it bigger than expected
	encodedTxsBytesSize := math.MaxInt
	numberOfTxsInProcess := len(s.sequenceInProgress.Txs)
	for encodedTxsBytesSize > maxBatchBytesSize && numberOfTxsInProcess > 0 {
		encodedTxs, err := state.EncodeTransactions(s.sequenceInProgress.Txs)
		if err != nil {
			log.Errorf("failed to encode txs, err: %w", err)
			return
		}
		encodedTxsBytesSize = len(encodedTxs)

		if encodedTxsBytesSize > maxBatchBytesSize && numberOfTxsInProcess > 0 {
			// if only one tx overflows, than it means, tx is invalid
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
				encodedTxsBytesSize, maxBatchBytesSize)
			s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:numberOfTxsInProcess-1]
			s.isSequenceTooBig = true
			numberOfTxsInProcess = len(s.sequenceInProgress.Txs)
		}
	}

	// process batch
	log.Infof("processing batch with %d txs. %d txs are new from this iteration", len(s.sequenceInProgress.Txs), appendedTxsAmount)
	processResponse, err := s.processTxs(ctx)
	if err != nil {
		s.sequenceInProgress = sequenceBeforeTryingToProcessNewTxs
		log.Errorf("failed to process txs, err: %w", err)
		return
	}

	// reprocess the batch until:
	// - all the txs in it are processed, so the batch doesn't include invalid txs
	// - the batch is processed (certain situations may cause the entire batch to not have effect on the state)
	unprocessedTxs := processResponse.unprocessedTxs
	for !processResponse.isBatchProcessed || len(processResponse.unprocessedTxs) > 0 {
		// include only processed txs in the sequence
		s.sequenceInProgress.Txs = make([]ethTypes.Transaction, 0, len(processResponse.processedTxs))
		for i := 0; i < len(processResponse.processedTxs); i++ {
			s.sequenceInProgress.Txs = append(s.sequenceInProgress.Txs, processResponse.processedTxs[i].Tx)
		}

		if len(s.sequenceInProgress.Txs) == 0 {
			log.Infof("sequence in progress doesn't have txs, no need to send a batch")
			break
		}
		log.Infof("failed to process batch or invalid txs. Retrying with %d txs", len(s.sequenceInProgress.Txs))
		// reprocess
		processResponse, err = s.processTxs(ctx)
		if err != nil {
			s.sequenceInProgress = sequenceBeforeTryingToProcessNewTxs
			log.Errorf("failed to reprocess txs, err: %w", err)
			return
		}
		if len(processResponse.processedTxsHashes) != 0 {
			for _, hash := range processResponse.processedTxsHashes {
				delete(unprocessedTxs, hash)
			}
		}
		for _, txHash := range processResponse.unprocessedTxsHashes {
			if _, ok := unprocessedTxs[txHash]; !ok {
				unprocessedTxs[txHash] = processResponse.unprocessedTxs[txHash]
			}
		}
	}
	log.Infof("%d txs processed successfully", len(processResponse.processedTxsHashes))

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

	// only save in DB processed transactions.
	err = s.storeProcessedTransactions(ctx, processResponse.processedTxs)
	if err != nil {
		s.sequenceInProgress = sequenceBeforeTryingToProcessNewTxs
		log.Errorf("failed to store processed txs, err: %w", err)
		return
	}
	log.Infof("%d txs stored and added into the trusted state", len(processResponse.processedTxs))

	invalidTxsHashes, failedTxsHashes := s.splitInvalidAndFailedTxs(ctx, unprocessedTxs, ticker)

	// update processed txs
	s.updateTxsStatus(ctx, ticker, processResponse.processedTxsHashes, pool.TxStatusSelected)
	// update invalid txs
	s.updateTxsStatus(ctx, ticker, invalidTxsHashes, pool.TxStatusInvalid)
	// update failed txs
	s.updateTxsStatus(ctx, ticker, failedTxsHashes, pool.TxStatusFailed)
	// increment counter for failed txs
	s.incrementFailedCounter(ctx, ticker, failedTxsHashes)
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
	err := s.pool.UpdateTxsStatus(ctx, hashes, status)
	for err != nil {
		log.Errorf("failed to update txs status to %s, err: %w", status, err)
		waitTick(ctx, ticker)
		err = s.pool.UpdateTxsStatus(ctx, hashes, status)
	}
}

func (s *Sequencer) incrementFailedCounter(ctx context.Context, ticker *time.Ticker, hashes []string) {
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
				"failed to rollback dbTx when closing batch that gave err: %s. Rollback err: %s",
				rollbackErr.Error(), err.Error(),
			)
		}
		return types.Sequence{}, err
	}
	// open next batch
	gerHash, err := s.getLatestGer(ctx, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return types.Sequence{}, fmt.Errorf(
				"failed to rollback dbTx when getting last GER that gave err: %s. Rollback err: %s",
				rollbackErr.Error(), err.Error(),
			)
		}
		return types.Sequence{}, err
	}

	processingCtx, err := s.openBatch(ctx, gerHash, dbTx)
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

func (s *Sequencer) isSequenceProfitable(ctx context.Context) bool {
	isProfitable, err := s.checker.IsSequenceProfitable(ctx, s.sequenceInProgress)
	if err != nil {
		log.Errorf("failed to check is sequence profitable, err: %w", err)
		return false
	}

	return isProfitable
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

	processBatchResp, err := s.state.ProcessSequencerBatch(ctx, lastBatchNumber, s.sequenceInProgress.Txs, dbTx)
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

	s.sequenceInProgress.StateRoot = processBatchResp.NewStateRoot
	s.sequenceInProgress.LocalExitRoot = processBatchResp.NewLocalExitRoot

	processedTxs, processedTxsHashes, unprocessedTxs, unprocessedTxsHashes := state.DetermineProcessedTransactions(processBatchResp.Responses)

	response := processTxResponse{
		processedTxs:         processedTxs,
		processedTxsHashes:   processedTxsHashes,
		unprocessedTxs:       unprocessedTxs,
		unprocessedTxsHashes: unprocessedTxsHashes,
		isBatchProcessed:     processBatchResp.IsBatchProcessed,
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

func (s *Sequencer) updateGerInBatch(ctx context.Context, lastGer *state.GlobalExitRoot) error {
	log.Info("update GER without closing batch as no txs have been added yet")

	dbTx, err := s.state.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for UpdateGERInOpenBatch tx, err: %w", err)
		return err
	}

	err = s.state.UpdateGERInOpenBatch(ctx, lastGer.GlobalExitRoot, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			log.Errorf(
				"failed to rollback dbTx when UpdateGERInOpenBatch that gave err: %w. Rollback err: %w",
				rollbackErr, err,
			)
			return err
		}
		log.Errorf("failed to update ger in open batch, err: %w", err)
		return err
	}

	if err := dbTx.Commit(ctx); err != nil {
		log.Errorf("failed to commit dbTx when processing UpdateGERInOpenBatch, err: %w", err.Error())
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

func (s *Sequencer) getLatestGer(ctx context.Context, dbTx pgx.Tx) (common.Hash, error) {
	ger, err := s.state.GetLatestGlobalExitRoot(ctx, dbTx)
	if err != nil && errors.Is(err, state.ErrNotFound) {
		return state.ZeroHash, nil
	} else if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get latest global exit root, err: %w", err)
	} else {
		return ger.GlobalExitRoot, nil
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
			log.Info("there is no suitable pending or failed txs in the pool, waiting...")
			waitTick(ctx, ticker)
			return 0
		}
	} else if err != nil {
		log.Errorf("failed to get pending tx, err: %w", err)
		return 0
	}
	for i := 0; i < len(pendTxs); i++ {
		if pendTxs[i].FailedCounter > s.cfg.MaxAllowedFailedCounter {
			hash := pendTxs[i].Transaction.Hash().String()
			log.Warnf("mark tx with hash %s as invalid, failed counter %d exceeded max %d from config",
				hash, pendTxs[i].FailedCounter, s.cfg.MaxAllowedFailedCounter)
			s.updateTxsStatus(ctx, ticker, []string{hash}, pool.TxStatusInvalid)
			continue
		}
		s.sequenceInProgress.Txs = append(s.sequenceInProgress.Txs, pendTxs[i].Transaction)
	}

	return uint64(len(pendTxs))
}
