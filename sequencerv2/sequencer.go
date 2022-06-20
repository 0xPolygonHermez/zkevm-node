package sequencerv2

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/hermeznetwork/hermez-core/ethermanv2/types"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencerv2/profitabilitychecker"
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
	manager txManager) (Sequencer, error) {

	checker := profitabilitychecker.New(etherman, priceGetter)
	return Sequencer{
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
		s.sequenceInProgress = s.newSequence()
	}

	// 3. Check if current sequence should be sent
	shouldSent, shouldCut := s.shouldSendSequences(ctx)
	if shouldSent {
		if shouldCut {
			cutSequence := s.closedSequences[len(s.closedSequences)-1]
			_ = s.txManager.SequenceBatches(s.closedSequences)
			s.closedSequences = []types.Sequence{cutSequence}
		} else {
			_ = s.txManager.SequenceBatches(s.closedSequences)
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
	res := s.state.ProcessBatchAndStoreLastTx(ctx, s.sequenceInProgress.Txs)
	if res.Err != nil {
		s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:len(s.sequenceInProgress.Txs)-1]
		log.Debugf("failed to process tx, hash: %s, err: %v", tx.Hash(), res.Err)
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
	lastSyncedBatchNum, err := s.state.GetLastBatchNumber(ctx, "")
	if err != nil {
		log.Errorf("failed to get last synced batch, err: %v", err)
		return false
	}
	lastEthBatchNum, err := s.state.GetLastBatchNumberSeenOnEthereum(ctx, "")
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

func (s *Sequencer) shouldSendSequences(ctx context.Context) (bool, bool) {
	estimatedGas, err := s.etherman.EstimateGasSequenceBatches(s.closedSequences)
	if err != nil && isDataForEthTxTooBig(err) {
		return true, true
	}

	if err != nil {
		log.Errorf("failed to estimate gas for sequence batches", err)
		return false, false
	}

	// todo: checkAgainstForcedBatchQueueTimeout

	// checkAgainstForcedBatchQueueTimeout
	lastL1TimeInteraction, err := s.state.GetLastL1InteractionTime(ctx)
	if err != nil {
		log.Errorf("failed to get last l1 interaction time, err: %v", err)
		return false, false
	}

	if lastL1TimeInteraction.Before(time.Now().Add(-s.cfg.LastL1InteractionTimeMaxWaitPeriod.Duration)) {
		return true, false
	}

	// check profitability
	if s.checker.IsSendSequencesProfitable(estimatedGas, s.closedSequences) {
		return true, false
	}

	return false, false
}

func (s *Sequencer) shouldCloseSequenceInProgress(ctx context.Context) bool {
	lastTimeGERUpdated, err := s.state.GetLastTimeGERUpdated(ctx)
	if err != nil {
		log.Errorf("failed to get last time GER updated, err: %v", err)
		return false
	}
	if lastTimeGERUpdated.Before(time.Now().Add(-s.cfg.LastTimeGERUpdatedMaxWaitPeriod.Duration)) {
		return true
	}

	lastTimeDeposit, err := s.state.GetLastTimeDeposit(ctx)
	if err != nil {
		log.Errorf("failed to get last time deposit, err: %v", err)
		return false
	}
	if lastTimeDeposit.Before(time.Now().Add(-s.cfg.LastTimeDepositMaxWaitPeriod.Duration)) {
		return true
	}

	lastBatchTime, err := s.state.GetLastBatchTime(ctx)
	if err != nil {
		log.Errorf("failed to get last batch time, err: %v", err)
		return false
	}
	if lastBatchTime.Before(time.Now().Add(-s.cfg.LastTimeBatchMaxWaitPeriod.Duration)) {
		return true
	}

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

func (s *Sequencer) newSequence() types.Sequence {
	return types.Sequence{}
}

func isDataForEthTxTooBig(err error) bool {
	if strings.Contains(err.Error(), errGasRequiredExceedsAllowance) ||
		errors.As(err, &core.ErrOversizedData) ||
		strings.Contains(err.Error(), errContentLengthTooLarge) {
		return true
	}
	return false
}
