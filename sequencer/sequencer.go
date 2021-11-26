package sequencer

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/synchronizer"
)

// Sequencer represents a sequencer
type Sequencer struct {
	Pool           pool.Pool
	State          state.State
	BatchProcessor state.BatchProcessor
	EthMan         etherman.EtherMan
	Synchronizer   synchronizer.Synchronizer

	ctx    context.Context
	cancel context.CancelFunc
}

// NewSequencer creates a new sequencer
func NewSequencer(cfg Config, pool pool.Pool, state state.State, ethMan etherman.EtherMan, sy synchronizer.Synchronizer) (Sequencer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	s := Sequencer{
		Pool:         pool,
		State:        state,
		EthMan:       ethMan,
		Synchronizer: sy,

		ctx:    ctx,
		cancel: cancel,
	}

	sy.RegisterNewConsolidatedStateHandler(s.onNewBatchPropostal)

	return s, nil
}

// Start starts the sequencer
func (s *Sequencer) Start() {
	// Infinite for loop:
	// 1. Wait for synchronizer to sync last batch
	// 2. Estimate available time to run selection
	// 3. Run selection
	// 4. Is selection profitable?
	// YES: send selection to Ethereum
	// NO: discard selection and wait for the new batch
	if err := s.Synchronizer.Sync(); err != nil {
		log.Fatal(err)
	}
}

func (s *Sequencer) onNewBatchPropostal(batchNumber uint64, root common.Hash) {
	ctx := context.Background()

	s.BatchProcessor = s.State.NewBatchProcessor(root, false)
	// get pending txs from the pool
	txs, err := s.Pool.GetPendingTxs(ctx)
	if err != nil {
		return
	}
	// estimate time for selecting txs
	estimatedTime, err := s.estimateTime()
	if err != nil {
		return
	}
	// select txs
	selectedTxs, err := s.selectTxs(txs, estimatedTime)
	if err != nil && !strings.Contains(err.Error(), "selection took too much time") {
		return
	}
	// check is it profitable to send selection
	isProfitable := s.isSelectionProfitable(selectedTxs)
	batch := state.Batch{Transactions: selectedTxs}
	if isProfitable {
		_, err = s.EthMan.SendBatch(batch)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

// Stop stops the sequencer
func (s *Sequencer) Stop() {
	s.cancel()
}

// selectTxs process txs and split valid txs into batches of txs. This process should be completed in less than selectionTime
func (s *Sequencer) selectTxs(pendingTxs []pool.Transaction, selectionTime time.Duration) ([]*types.Transaction, error) {
	start := time.Now()
	sortedTxs := s.sortTxs(pendingTxs)
	var selectedTxs []*types.Transaction
	ctx := context.Background()
	for _, tx := range sortedTxs {
		// check if tx is valid
		if err := s.BatchProcessor.CheckTransaction(tx.Transaction); err != nil {
			if err = s.Pool.UpdateTxState(ctx, tx.Hash(), pool.TxStateInvalid); err != nil {
				return nil, err
			}
		} else {
			selectedTxs = append(selectedTxs, &tx.Transaction)
		}

		elapsed := time.Since(start)
		if elapsed > selectionTime {
			return selectedTxs, nil
		}
	}
	return selectedTxs, nil
}

func (s *Sequencer) sortTxs(txs []pool.Transaction) []pool.Transaction {
	sort.Slice(txs, func(i, j int) bool {
		costI := txs[i].Cost()
		costJ := txs[j].Cost()
		if costI != costJ {
			return costI.Cmp(costJ) >= 1
		}
		return txs[i].Nonce() < txs[j].Nonce()
	})
	return txs
}

// estimateTime Estimate available time to run selection
func (s *Sequencer) estimateTime() (time.Duration, error) {
	return time.Hour, nil
}

func (s *Sequencer) isSelectionProfitable(txs []*types.Transaction) bool {
	return true
}
