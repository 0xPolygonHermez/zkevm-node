package sequencer

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/jackc/pgx/v4"
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg Config

	Pool    pool.Pool
	State   state.State
	EthMan  etherman.EtherMan
	Address common.Address
	ChainID uint64
	strategy.TxSelector
	strategy.TxProfitabilityChecker

	ctx    context.Context
	cancel context.CancelFunc
}

// NewSequencer creates a new sequencer
func NewSequencer(cfg Config, pool pool.Pool, state state.State, ethMan etherman.EtherMan) (Sequencer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	var txSelector strategy.TxSelector
	switch cfg.Strategy.Type {
	case strategy.AcceptAll:
		txSelector = strategy.NewTxSelectorAcceptAll(cfg.Strategy)
	case strategy.Base:
		txSelector = strategy.NewTxSelectorBase(cfg.Strategy)
	}

	var txProfitabilityChecker strategy.TxProfitabilityChecker
	switch cfg.Strategy.TxProfitabilityCheckerType {
	case strategy.ProfitabilityAcceptAll:
		txProfitabilityChecker = &strategy.TxProfitabilityCheckerAcceptAll{}
	case strategy.ProfitabilityBase:
		txProfitabilityChecker = strategy.NewTxProfitabilityCheckerBase(ethMan, cfg.Strategy.MinReward.Int)
	}

	seqAddress := ethMan.GetAddress()
	chainID, err := getChainID(ctx, state, seqAddress)
	if err != nil {
		cancel()
		return Sequencer{}, fmt.Errorf("failed to get chain id for the sequencer, err: %v", err)
	}
	s := Sequencer{
		cfg:     cfg,
		Pool:    pool,
		State:   state,
		EthMan:  ethMan,
		Address: seqAddress,
		ChainID: chainID,

		TxSelector:             txSelector,
		TxProfitabilityChecker: txProfitabilityChecker,

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
	selectedTxs, selectedTxsHashes, invalidTxsHashes, err := s.TxSelector.SelectTxs(bp, txs, estimatedTime)
	if err != nil && !strings.Contains(err.Error(), "selection took too much time") {
		log.Errorf("failed to select txs, err: %v", err)
		return
	}

	if err = s.Pool.UpdateTxsState(s.ctx, invalidTxsHashes, pool.TxStateInvalid); err != nil {
		log.Errorf("failed to update txs state to invalid, err: %v", err)
		return
	}

	// 4. Is selection profitable?
	// check is it profitable to send selection
	isProfitable, err := s.TxProfitabilityChecker.IsProfitable(s.ctx, selectedTxs)
	if err != nil {
		log.Errorf("failed to check that txs are profitable or not, err: %v", err)
		return
	}
	if isProfitable && len(selectedTxs) > 0 {
		// assume, that fee for 1 tx is 1 matic
		maticAmount := big.NewInt(int64(len(selectedTxs)))
		maticAmount = big.NewInt(0).Mul(maticAmount, big.NewInt(encoding.TenToThePowerOf18))

		// YES: send selection to Ethereum
		sendBatchTx, err := s.EthMan.SendBatch(s.ctx, selectedTxs, maticAmount)
		if err != nil {
			log.Errorf("failed to send batch proposal to ethereum, err: %v", err)
			return
		}
		log.Infof("Batch proposal sent successfully: %s", sendBatchTx.Hash().Hex())

		// update txs in the pool as selected
		err = s.Pool.UpdateTxsState(s.ctx, selectedTxsHashes, pool.TxStateSelected)
		if err != nil {
			log.Warnf("failed to update txs state to selected, err: %v", err)
		}
		log.Infof("Finished updating selected transactions state in the pool")
	}
	// NO: discard selection and wait for the new batch
}

// estimateTime Estimate available time to run selection
func (s *Sequencer) estimateTime(txs []pool.Transaction) (time.Duration, error) {
	return time.Hour, nil
}

func getChainID(ctx context.Context, st state.State, seqAddress common.Address) (uint64, error) {
	const intervalToCheckSequencerRegistrationInSeconds = 3
	var (
		seq *state.Sequencer
		err error
	)
	for {
		seq, err = st.GetSequencer(ctx, seqAddress)
		if err != nil {
			if err == pgx.ErrNoRows {
				log.Warnf("make sure the address %s has been registered in the smart contract as a sequencer, err: %v", seqAddress.Hex(), err)
				lastSyncedBatchNum, err := st.GetLastBatchNumber(ctx)
				if err != nil {
					log.Errorf("failed to get last synced batch, err: %v", err)
					return 0, err
				}
				lastEthBatchNum, err := st.GetLastBatchNumberSeenOnEthereum(ctx)
				if err != nil {
					log.Errorf("failed to get last eth batch, err: %v", err)
					return 0, err
				}

				if lastEthBatchNum == 0 {
					log.Warnf("last eth batch num is 0, waiting to sync...")
				} else {
					const oneHundred = 100
					percentage := lastSyncedBatchNum * oneHundred / lastEthBatchNum
					log.Warnf("node is still syncing, synced %d%%", percentage)
				}
				time.Sleep(intervalToCheckSequencerRegistrationInSeconds * time.Second)
				continue
			} else {
				return 0, err
			}
		}
		return seq.ChainID.Uint64(), nil
	}
}
