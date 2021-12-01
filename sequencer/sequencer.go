package sequencer

import (
	"context"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg Config

	Pool           pool.Pool
	State          state.State
	BatchProcessor state.BatchProcessor
	EthMan         etherman.EtherMan

	ctx    context.Context
	cancel context.CancelFunc
}

// NewSequencer creates a new sequencer
func NewSequencer(cfg Config, pool pool.Pool, state state.State, ethMan etherman.EtherMan) (Sequencer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	s := Sequencer{
		cfg:    cfg,
		Pool:   pool,
		State:  state,
		EthMan: ethMan,

		ctx:    ctx,
		cancel: cancel,
	}

	return s, nil
}

// Start starts the sequencer
func (s *Sequencer) Start() {
	// Infinite for loop:
	for {
		time.Sleep(s.cfg.IntervalToProposeBatch)

		ctx := context.Background()

		// 1. Wait for synchronizer to sync last batch
		// TODO: state will provide methods to check if it is synchronized

		// 2. Estimate available time to run selection
		// get pending txs from the pool
		txs, err := s.Pool.GetPendingTxs(ctx)
		if err != nil {
			return
		}

		// estimate time for selecting txs
		estimatedTime, err := s.estimateTime(txs)
		if err != nil {
			return
		}

		log.Infof("Estimated time for selecting txs is %dms", estimatedTime.Milliseconds())

		// 3. Run selection
		// select txs
		selectedTxs, err := s.selectTxs(txs, estimatedTime)
		if err != nil && !strings.Contains(err.Error(), "selection took too much time") {
			return
		}

		// 4. Is selection profitable?
		// check is it profitable to send selection
		isProfitable := s.isSelectionProfitable(selectedTxs)
		batch := state.Batch{Transactions: selectedTxs}
		var maticAmount *big.Int //TODO calculate the amount depending on the profitability
		if isProfitable {
			// YES: send selection to Ethereum
			_, err = s.EthMan.SendBatch(ctx, batch, maticAmount)
			if err != nil {
				log.Error(err)
				return
			}
		}
		// NO: discard selection and wait for the new batch
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
		_, _, _, err := s.BatchProcessor.CheckTransaction(&tx.Transaction)
		if err != nil {
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
func (s *Sequencer) estimateTime(txs []pool.Transaction) (time.Duration, error) {
	return time.Hour, nil
}

func (s *Sequencer) isSelectionProfitable(txs []*types.Transaction) bool {
	return true
}
