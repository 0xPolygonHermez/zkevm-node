package sequencer

import (
	"context"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg Config

	Pool    pool.Pool
	State   state.State
	EthMan  etherman.EtherMan
	Address common.Address

	ctx    context.Context
	cancel context.CancelFunc
}

// NewSequencer creates a new sequencer
func NewSequencer(cfg Config, pool pool.Pool, state state.State, ethMan etherman.EtherMan) (Sequencer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	s := Sequencer{
		cfg:     cfg,
		Pool:    pool,
		State:   state,
		EthMan:  ethMan,
		Address: ethMan.GetAddress(),

		ctx:    ctx,
		cancel: cancel,
	}

	return s, nil
}

// Start starts the sequencer
func (s *Sequencer) Start() {
	// Infinite for loop:
	for {
		select {
		case <-time.After(s.cfg.IntervalToProposeBatch.Duration):
			s.tryProposeBatch()
		case <-s.ctx.Done():
			return
		}
	}
}

// Stop stops the sequencer
func (s *Sequencer) Stop() {
	s.cancel()
}

func (s *Sequencer) tryProposeBatch() {
	// 1. Wait for synchronizer to sync last batch
	lastSyncedBatchNum, err := s.State.GetLastBatchNumber(s.ctx)
	if err != nil {
		log.Errorf("failed to get last synced batch, err: %v", err)
		return
	}
	lastEthBatchNum, err := s.State.GetLastBatchNumberSeenOnEthereum(s.ctx)
	if err != nil {
		log.Errorf("failed to get last eth batch, err: %v", err)
		return
	}
	if lastSyncedBatchNum+s.cfg.SyncedBlockDif < lastEthBatchNum {
		log.Infow("waiting for the state to be synced, lastSyncedBatchNum: %d, lastEthBatchNum: %d", lastSyncedBatchNum, lastEthBatchNum)
		return
	}

	// 2. Estimate available time to run selection
	// get pending txs from the pool
	txs, err := s.Pool.GetPendingTxs(s.ctx)
	if err != nil {
		log.Errorf("failed to get pending txs, err: %v", err)
		return
	}

	if len(txs) == 0 {
		log.Infof("transactions pool is empty, waiting for the new txs...")
		return
	}

	// estimate time for selecting txs
	estimatedTime, err := s.estimateTime(txs)
	if err != nil {
		log.Errorf("failed to estimate time for selecting txs, err: %v", err)
		return
	}

	log.Infof("Estimated time for selecting txs is %dms", estimatedTime.Milliseconds())

	// 3. Run selection
	// init batch processor
	lastBatch, err := s.State.GetLastBatch(s.ctx, false)
	if err != nil {
		log.Errorf("failed to get last batch from the state, err: %v", err)
		return
	}
	bp, err := s.State.NewBatchProcessor(s.Address, lastBatch.BatchNumber)
	if err != nil {
		log.Errorf("failed to create new batch processor, err: %v", err)
		return
	}

	// select txs
	selectedTxs, err := s.selectTxs(bp, txs, estimatedTime)
	if err != nil && !strings.Contains(err.Error(), "selection took too much time") {
		log.Errorf("failed to get last batch from the state, err: %v", err)
		return
	}

	// 4. Is selection profitable?
	// check is it profitable to send selection
	isProfitable := s.isSelectionProfitable(selectedTxs)
	log.Infof("Transaction selection is profitable: %v", isProfitable)
	if isProfitable && len(selectedTxs) > 0 {
		// assume, that fee for 1 tx is 1 matic
		maticAmount := big.NewInt(int64(len(selectedTxs)))

		// YES: send selection to Ethereum
		sendBatchTx, err := s.EthMan.SendBatch(s.ctx, selectedTxs, maticAmount)
		if err != nil {
			log.Errorf("failed to send batch proposal to ethereum, err: %v", err)
			return
		}
		log.Infof("Batch proposal sent successfully: %s", sendBatchTx.Hash().Hex())

		// update txs in the pool as selected
		for _, tx := range selectedTxs {
			err := s.Pool.UpdateTxState(s.ctx, tx.Hash(), pool.TxStateSelected)
			if err != nil {
				log.Warnf("failed to update tx(%s) state to selected, err: %v", tx.Hash().Hex(), err)
			}
		}
		log.Infof("Finished updating selected transactions state in the pool")
	}
	// NO: discard selection and wait for the new batch
}

// selectTxs process txs and split valid txs into batches of txs. This process should be completed in less than selectionTime
func (s *Sequencer) selectTxs(batchProcessor state.BatchProcessor, pendingTxs []pool.Transaction, selectionTime time.Duration) ([]*types.Transaction, error) {
	start := time.Now()
	sortedTxs := s.sortTxs(pendingTxs)
	var selectedTxs []*types.Transaction
	for _, tx := range sortedTxs {
		// check if tx is valid
		_, _, _, err := batchProcessor.CheckTransaction(&tx.Transaction)
		if err != nil {
			if err = s.Pool.UpdateTxState(s.ctx, tx.Hash(), pool.TxStateInvalid); err != nil {
				return nil, err
			}
		} else {
			t := tx.Transaction
			selectedTxs = append(selectedTxs, &t)
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
