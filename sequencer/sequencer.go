package sequencer

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type Sequencer struct {
	Pool               pool.Pool
	State              state.State
	BatchProcessor     state.BatchProcessor
	EthClient          eth.Client
	SynchronizerClient SynchronizerClient

	ctx          context.Context
	cancel       context.CancelFunc
	waitDuration time.Duration
}

func NewSequencer(cfg Config, pool pool.Pool, state state.State, bp state.BatchProcessor, ethClient eth.Client, syncClient SynchronizerClient) (Sequencer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	waitDuration := time.Duration(cfg.sendBatchFrequencyInSeconds) * time.Second
	return Sequencer{
		Pool:               pool,
		State:              state,
		EthClient:          ethClient,
		BatchProcessor:     bp,
		SynchronizerClient: syncClient,

		ctx:          ctx,
		cancel:       cancel,
		waitDuration: waitDuration,
	}, nil
}

func (s *Sequencer) Start() {
	// Infinite for loop:
	// 1. Wait for synchronizer to sync last batch
	// 2. Estimate available time to run selection
	// 3. Run selection
	// 4. Is selection profitable?
	// YES: send selection to Ethereum
	// NO: discard selection and wait for the new batch
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(s.waitDuration):
				// get pending txs from the pool
				txs, err := s.Pool.GetPendingTxs()
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
				if isProfitable {
					s.EthClient.SendBatch(selectedTxs, len(selectedTxs))
				}
			}
		}
	}()
}

func (s *Sequencer) Stop() {
	s.cancel()
}

type batch struct {
	txs []pool.Transaction
}

// selectTxs process txs and split valid txs into batches of txs. This process should be completed in less than selectionTime
func (s *Sequencer) selectTxs(pendingTxs []pool.Transaction, selectionTime time.Duration) ([]types.Transaction, error) {
	start := time.Now()
	sortedTxs := s.sortTxs(pendingTxs)
	var selectedTxs []types.Transaction
	for _, tx := range sortedTxs {
		// check if tx is valid
		if err := s.BatchProcessor.CheckTransaction(tx.Transaction); err != nil {
			if err = s.Pool.UpdateTxState(tx.Hash(), pool.TxStateInvalid); err != nil {
				return nil, err
			}
		} else {
			selectedTxs = append(selectedTxs, tx.Transaction)
		}
	}
	elapsed := time.Since(start)
	if elapsed > selectionTime {
		return nil, errors.New("selection took too much time")
	}
	return selectedTxs, nil
}

func (s *Sequencer) sortTxs(txs []pool.Transaction) []pool.Transaction {
	sort.Slice(txs, func(i, j int) bool {
		if txs[i].Cost() != txs[j].Cost() {
			return txs[i].Cost().Cmp(txs[j].Cost()) >= 1
		}
		return txs[i].Nonce() < txs[j].Nonce()
	})
	return txs
}

// estimateTime Estimate available time to run selection
func (s *Sequencer) estimateTime() (time.Duration, error) {
	return time.Hour, nil
}

func (s *Sequencer) isSelectionProfitable(txs []types.Transaction) bool {
	return true
}
