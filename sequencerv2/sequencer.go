// nolint
package sequencerv2

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/ethermanv2"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg Config

	pool      txPool
	state     stateInterface
	ethMan    etherman
	txManager txManager

	address  common.Address
	sequence ethermanv2.Sequence
}

// NewSequencer init sequencer
func NewSequencer(
	cfg Config,
	pool txPool,
	state stateInterface,
	ethMan etherman,
	manager txManager) (Sequencer, error) {

	seqAddress := ethMan.GetAddress()
	return Sequencer{
		cfg:       cfg,
		pool:      pool,
		state:     state,
		ethMan:    ethMan,
		address:   seqAddress,
		txManager: manager,
	}, nil
}

// Start starts the sequencer
func (s *Sequencer) Start(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.IntervalToProposeBatch.Duration)
	defer ticker.Stop()
	var root []byte
	// Infinite for loop:
	for {
		root = s.tryToSendSequenceBatches(ctx, root)
		select {
		case <-ticker.C:
			// nothing
		case <-ctx.Done():
			return
		}
	}
}

func (s *Sequencer) tryToSendSequenceBatches(ctx context.Context, root []byte) []byte {
	// 1. Wait for synchronizer to sync last batch
	if !s.isSynced(ctx) {
		return nil
	}
	// 2. get pending tx from the pool
	tx, err := s.getPendingTx(ctx)
	if err != nil {
		return nil
	}

	// 3. Process tx
	s.sequence.Txs = append(s.sequence.Txs, tx.Transaction)
	res := s.state.ProcessSequence(ctx, s.sequence)
	if res.Err != nil {
		return nil
	}
	// 4. Send sequence to ethereum
	err = s.txManager.SequenceBatches([]*ethermanv2.Sequence{&s.sequence})
	if err != nil {
		return nil
	}
	return root
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
	if lastSyncedBatchNum+s.cfg.SyncedBlockDif < lastEthBatchNum {
		log.Infof("waiting for the state to be synced, lastSyncedBatchNum: %d, lastEthBatchNum: %d", lastSyncedBatchNum, lastEthBatchNum)
		return false
	}
	return true
}

func (s *Sequencer) getPendingTx(ctx context.Context) (*pool.Transaction, error) {
	return nil, nil
}
