package sequencer

import (
	"time"

	eth "github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

type Sequencer struct {
	Pool               pool.Pool
	State              state.State
	BatchProcessor     state.BatchProcessor
	EthClient          eth.Client
	SynchronizerClient SynchronizerClient
}

func NewSequencer(pool pool.Pool, state state.State, bp state.BatchProcessor, ethClient eth.Client, syncClient SynchronizerClient) (Sequencer, error) {
	return Sequencer{
		Pool:               pool,
		State:              state,
		EthClient:          ethClient,
		BatchProcessor:     bp,
		SynchronizerClient: syncClient,
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
}

type batch struct {
	txs []pool.Transaction
}

// selectTxs process txs and split valid txs into batches of txs. This process should be completed in less than selectionTime
func (s *Sequencer) selectTxs(pendingTxs []pool.Transaction, selectionTime time.Duration) ([]batch, error) {
	panic("not implemented")
}

// estimateTime Estimate available time to run selection
func (s *Sequencer) estimateTime() (time.Duration, error) {
	panic("not implemented")
}

func (s *Sequencer) isSelectionProfitable(b batch) bool {
	panic("not implemented")
}
