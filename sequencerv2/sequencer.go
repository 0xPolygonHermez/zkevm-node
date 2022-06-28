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
	if !s.isSynced(ctx) {
		log.Infof("wait for synchronizer to sync last batch")
		waitTick(ctx, ticker)
		return
	}

	log.Infof("synchronizer has synced last batch, checking if current sequence should be closed")
	if s.shouldCloseSequenceInProgress(ctx) {
		log.Infof("current sequence should be closed")
		s.closedSequences = append(s.closedSequences, s.sequenceInProgress)
		newSequence, err := s.newSequence(ctx)
		if err != nil {
			log.Errorf("failed to create new sequence, err: %v", err)
			s.closedSequences = s.closedSequences[:len(s.closedSequences)-1]
			return
		}
		s.sequenceInProgress = newSequence
	}

	log.Infof("checking if current sequence should be sent")
	shouldSent, shouldCut := s.shouldSendSequences(ctx)
	if shouldSent {
		log.Infof("current sequence should be sent")
		if shouldCut {
			log.Infof("current sequence should be cut")
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

	log.Infof("getting pending tx from the pool")
	tx, ok := s.getMostProfitablePendingTx(ctx)
	if !ok {
		return
	}

	if tx == nil {
		log.Infof("waiting for pending txs...")
		waitTick(ctx, ticker)
		return
	}

	log.Infof("processing tx")
	s.sequenceInProgress.Txs = append(s.sequenceInProgress.Txs, tx.Transaction)
	_, err := s.state.ProcessBatch(ctx, s.sequenceInProgress.Txs)
	if err != nil {
		s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:len(s.sequenceInProgress.Txs)-1]
		log.Debugf("failed to process tx, hash: %s, err: %v", tx.Hash(), err)
		return
	}

	log.Infof("marking tx as selected in the pool")
	// TODO: add correct handling in case update didn't go through
	_ = s.pool.UpdateTxState(ctx, tx.Hash(), pool.TxStateSelected)

	log.Infof("TODO: broadcast tx in a new l2 block")
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
	root, err := s.state.GetLatestGlobalExitRoot(ctx, nil)
	if err != nil {
		return types.Sequence{}, err
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
