package sequencerv2

import (
	"context"
	"time"

	"github.com/hermeznetwork/hermez-core/ethermanv2"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

const (
	maxSequencesLength = 10
	maxTxsInSequence   = 10
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg Config

	pool      txPool
	state     stateInterface
	txManager txManager

	sequencesInProgress []*ethermanv2.Sequence
	sequenceInProgress  ethermanv2.Sequence
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

	parentBlock, err := s.state.GetLastL2Block(ctx)
	if err != nil {
		log.Errorf("failed to get latest block hash tx, err: %v", err)
	}
	parentTxHash := parentBlock.TxHash
	for {
		// 1. Wait for synchronizer to sync last batch
		if !s.isSynced(ctx) {
			select {
			case <-ticker.C:
				// nothing
			case <-ctx.Done():
				return
			}
			continue
		}

		// 2. Check if current sequence should be sent
		if s.isSequencesInProgressShouldBeSent() {
			_ = s.txManager.SequenceBatches(s.sequencesInProgress)
			s.sequenceInProgress = s.newSequence()
			s.sequencesInProgress = []*ethermanv2.Sequence{}
		}

		// 3. Check if current sequence should be closed
		if s.isSequenceInProgressShouldBeClosed() {
			s.sequencesInProgress = append(s.sequencesInProgress, &s.sequenceInProgress)
			s.sequenceInProgress = s.newSequence()
		}

		// 4. get pending tx from the pool
		tx, err := s.getPendingTx(ctx)
		if err != nil {
			log.Errorf("failed to get pending tx, err: %v", err)
			continue
		}

		if tx == nil {
			log.Infof("waiting for pending txs...")
			select {
			case <-ticker.C:
				// nothing
			case <-ctx.Done():
				return
			}
			continue
		}

		// 5. Process tx
		s.sequenceInProgress.Txs = append(s.sequenceInProgress.Txs, tx.Transaction)
		res := s.state.ProcessSequence(ctx, s.sequenceInProgress)
		if res.Err != nil {
			s.sequenceInProgress.Txs = s.sequenceInProgress.Txs[:len(s.sequenceInProgress.Txs)-1]
			log.Errorf("failed to process tx, hash: %s, err: %v", tx.Hash(), err)
			continue
		}

		// 6. Mark tx as selected in the pool
		// TODO: add correct handling in case update didn't go through
		_ = s.pool.UpdateTxState(ctx, tx.Hash(), pool.TxStateSelected)

		// 7. create new l2 block
		block := &state.L2Block{
			TxHash:       tx.Hash(),
			ParentTxHash: parentTxHash,
			ReceivedAt:   tx.ReceivedAt,
		}
		err = s.state.AddL2Block(ctx, block, "")
		if err != nil {
			log.Fatalf("failed to add L2 block for tx hash %q, err: %v", tx.Hash(), err)
		}
		parentTxHash = tx.Hash()

		// 8. TODO: broadcast tx
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

func (s *Sequencer) isSequencesInProgressShouldBeSent() bool {
	return len(s.sequencesInProgress) >= maxSequencesLength
}

func (s *Sequencer) isSequenceInProgressShouldBeClosed() bool {
	return len(s.sequenceInProgress.Txs) >= maxTxsInSequence
}

func (s *Sequencer) getPendingTx(ctx context.Context) (*pool.Transaction, error) {
	return nil, nil
}

func (s *Sequencer) newSequence() ethermanv2.Sequence {
	return ethermanv2.Sequence{}
}
