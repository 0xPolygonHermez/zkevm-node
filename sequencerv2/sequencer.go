package sequencerv2

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/hermeznetwork/hermez-core/ethermanv2/types"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencerv2/profitabilitychecker"
	"github.com/hermeznetwork/hermez-core/statev2"
)

const (
	errGasRequiredExceedsAllowance = "gas required exceeds allowance"
	errContentLengthTooLarge       = "content length too large"
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg Config

	pool      txPool
	state     stateInterface
	txManager txManager
	etherman  etherman
	checker   *profitabilitychecker.Checker

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
	manager txManager) (*Sequencer, error) {
	checker := profitabilitychecker.New(cfg.ProfitabilityChecker, etherman, priceGetter)
	return &Sequencer{
		cfg:       cfg,
		pool:      pool,
		state:     state,
		etherman:  etherman,
		checker:   checker,
		txManager: manager,
	}, nil
}

// Start starts the sequencer
func (s *Sequencer) Start(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.WaitPeriodPoolIsEmpty.Duration)
	defer ticker.Stop()
	for {
		s.tryToProcessTx(ctx, ticker)
	}
}

func (s *Sequencer) tryToProcessTx(ctx context.Context, ticker *time.Ticker) {
	// 1. Wait for synchronizer to sync last batch
	if !s.isSynced(ctx) {
		waitTick(ctx, ticker)
		return
	}

	// 2. Check if current sequence should be closed
	if s.shouldCloseSequenceInProgress(ctx) {
		s.closedSequences = append(s.closedSequences, s.sequenceInProgress)
		newSequence, err := s.newSequence(ctx)
		if err != nil {
			log.Errorf("failed to create new sequence, err: %v", err)
			s.closedSequences = s.closedSequences[:len(s.closedSequences)-1]
			return
		}
		s.sequenceInProgress = newSequence
	}

	// 3. Check if current sequence should be sent
	shouldSent, shouldCut := s.shouldSendSequences(ctx)
	if shouldSent {
		if shouldCut {
			cutSequence := s.closedSequences[len(s.closedSequences)-1]
			err := s.txManager.SequenceBatches(s.closedSequences)
			if err != nil {
				log.Errorf("failed to SequenceBatches, err: %v", err)
				return
			}
			s.closedSequences = []types.Sequence{cutSequence}
		} else {
			err := s.txManager.SequenceBatches(s.closedSequences)
			if err != nil {
				log.Errorf("failed to SequenceBatches, err: %v", err)
				return
			}
			s.closedSequences = []types.Sequence{}
		}
	}

	// 4. get pending tx from the pool
	tx, ok := s.getMostProfitablePendingTx(ctx)
	if !ok {
		return
	}

	if tx == nil {
		log.Infof("waiting for pending txs...")
		waitTick(ctx, ticker)
		return
	}

	// 5. Process tx
	s.sequenceInProgress.Txs = append(s.sequenceInProgress.Txs, tx.Transaction)
	processBatchResp, err := s.state.ProcessBatch(ctx, s.sequenceInProgress.Txs)
	if err != nil {
		s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:len(s.sequenceInProgress.Txs)-1]
		log.Debugf("failed to process tx, hash: %s, err: %v", tx.Hash(), err)
		return
	}

	s.lastStateRoot = processBatchResp.NewStateRoot
	s.lastLocalExitRoot = processBatchResp.NewLocalExitRoot

	// TODO: add logic based on this response to decide which txs we include on the DB
	err = s.state.StoreTransactions(ctx, s.lastBatchNum, processBatchResp.Responses)
	if err != nil {
		log.Errorf("failed to store transactions, err: %v", err)
		return
	}

	// 6. Mark tx as selected in the pool
	// TODO: add correct handling in case update didn't go through
	_ = s.pool.UpdateTxState(ctx, tx.Hash(), pool.TxStateSelected)

	// 7. broadcast tx in a new l2 block
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
	lastSyncedBatchNum, err := s.state.GetLastVirtualBatchNum(ctx)
	if err != nil {
		log.Errorf("failed to get last synced batch, err: %v", err)
		return false
	}
	lastEthBatchNum, err := s.state.GetLastBatchNumberSeenOnEthereum(ctx)
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
		return true, true
	}

	if err != nil {
		log.Errorf("failed to estimate gas for sequence batches", err)
		return false, false
	}

	// TODO: checkAgainstForcedBatchQueueTimeout

	lastL1TimeInteraction, err := s.state.GetLastSendSequenceTime(ctx)
	if err != nil {
		log.Errorf("failed to get last l1 interaction time, err: %v", err)
		return false, false
	}

	if lastL1TimeInteraction.Before(time.Now().Add(-s.cfg.LastL1InteractionTimeMaxWaitPeriod.Duration)) {
		// check profitability
		if s.checker.IsSendSequencesProfitable(estimatedGas, s.closedSequences) {
			return true, false
		}
	}

	return false, false
}

// shouldCloseSequenceInProgress checks if sequence should be closed or not
// in case it's enough blocks since last GER update, long time since last batch and sequence is profitable
func (s *Sequencer) shouldCloseSequenceInProgress(ctx context.Context) bool {
	numberOfBlocks, err := s.state.GetNumberOfBlocksSinceLastGERUpdate(ctx)
	if err != nil {
		log.Errorf("failed to get last time GER updated, err: %v", err)
		return false
	}
	if numberOfBlocks >= s.cfg.WaitBlocksToUpdateGER {
		return s.isSequenceProfitable(ctx)
	}

	lastBatchTime, err := s.state.GetLastBatchTime(ctx)
	if err != nil {
		log.Errorf("failed to get last batch time, err: %v", err)
		return false
	}
	if lastBatchTime.Before(time.Now().Add(-s.cfg.LastTimeBatchMaxWaitPeriod.Duration)) {
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

func (s *Sequencer) getMostProfitablePendingTx(ctx context.Context) (*pool.Transaction, bool) {
	tx, err := s.pool.GetPendingTxs(ctx, false, 1)
	if err != nil {
		log.Errorf("failed to get pending tx, err: %v", err)
		return nil, false
	}
	if len(tx) == 0 {
		log.Infof("waiting for pending tx to appear...")
		return nil, false
	}
	return &tx[0], true
}

func (s *Sequencer) newSequence(ctx context.Context) (types.Sequence, error) {
	// close current batch
	if s.lastStateRoot.String() != "" || s.lastLocalExitRoot.String() != "" {
		err := s.state.CloseBatch(ctx, s.lastBatchNum, s.lastStateRoot, s.lastLocalExitRoot)
		if err != nil {
			return types.Sequence{}, fmt.Errorf("failed to close batch, err: %v", err)
		}
	}

	root, err := s.state.GetLatestGlobalExitRoot(ctx, nil)
	if err != nil {
		return types.Sequence{}, fmt.Errorf("failed to get latest global exit root, err: %v", err)
	}

	if s.lastBatchNum == 0 {
		s.lastBatchNum, err = s.state.GetLastBatchNumber(ctx)
		if err != nil {
			return types.Sequence{}, fmt.Errorf("failed to get last batch number, err: %v", err)
		}
	} else {
		s.lastBatchNum = s.lastBatchNum + 1
	}

	batchHeader := statev2.Batch{
		BatchNumber: s.lastBatchNum,
		Timestamp:   time.Now(),
	}
	err = s.state.StoreBatchHeader(ctx, batchHeader, nil)
	if err != nil {
		return types.Sequence{}, fmt.Errorf("failed to store batch header, err: %v", err)
	}

	return types.Sequence{
		GlobalExitRoot:  root.GlobalExitRoot,
		Timestamp:       time.Now().Unix(),
		ForceBatchesNum: 0,
		Txs:             nil,
	}, nil
}

func isDataForEthTxTooBig(err error) bool {
	return strings.Contains(err.Error(), errGasRequiredExceedsAllowance) ||
		errors.As(err, &core.ErrOversizedData) ||
		strings.Contains(err.Error(), errContentLengthTooLarge)
}
