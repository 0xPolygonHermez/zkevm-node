package sequencerv2

import (
	"context"
	"time"

	"github.com/hermeznetwork/hermez-core/ethermanv2/types"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
)

const (
	maxSequencesLength = 5
	maxTxsInSequence   = 5
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg Config

	pool      txPool
	state     stateInterface
	txManager txManager

	closedSequences    []types.Sequence
	sequenceInProgress types.Sequence
}

// New init sequencer
func New(
	cfg Config,
	pool txPool,
	state stateInterface,
	manager txManager) (Sequencer, error) {
	return Sequencer{
		cfg:       cfg,
		pool:      pool,
		state:     state,
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
	if s.isSequenceInProgressShouldBeClosed() {
		s.closedSequences = append(s.closedSequences, s.sequenceInProgress)
		s.sequenceInProgress = s.newSequence()
	}

	// 3. Check if current sequence should be sent
	if s.isClosedSequencesShouldBeSent() {
		_ = s.txManager.SequenceBatches(s.closedSequences)
		s.closedSequences = []types.Sequence{}
	}

	// 4. get pending tx from the pool
	tx, ok := s.getPendingTx(ctx)
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
	res := s.state.ProcessBatchAndStoreLatestTx(ctx, s.sequenceInProgress.Txs)
	if res.Err != nil {
		s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:len(s.sequenceInProgress.Txs)-1]
		log.Errorf("failed to process tx, hash: %s, err: %v", tx.Hash(), res.Err)
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

func (s *Sequencer) isClosedSequencesShouldBeSent() bool {
	return len(s.closedSequences) >= maxSequencesLength
}

func (s *Sequencer) isSequenceInProgressShouldBeClosed() bool {
	return len(s.sequenceInProgress.Txs) >= maxTxsInSequence
}

func (s *Sequencer) getPendingTx(ctx context.Context) (*pool.Transaction, bool) {
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
