package sequencer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy/txprofitabilitychecker"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy/txselector"
	"github.com/hermeznetwork/hermez-core/state"
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg Config

	Pool    txPool
	State   state.State
	EthMan  etherman
	Address common.Address
	ChainID uint64

	txselector.TxSelector
	TxProfitabilityChecker txProfitabilityChecker

	ctx    context.Context
	cancel context.CancelFunc
}

// NewSequencer creates a new sequencer
func NewSequencer(cfg Config, pool txPool, state state.State, ethMan etherman) (Sequencer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	var txSelector txselector.TxSelector
	switch cfg.Strategy.TxSelector.Type {
	case txselector.AcceptAllType:
		txSelector = txselector.NewTxSelectorAcceptAll()
	case txselector.BaseType:
		txSelector = txselector.NewTxSelectorBase(cfg.Strategy.TxSelector)
	}

	var txProfitabilityChecker txProfitabilityChecker
	switch cfg.Strategy.TxProfitabilityChecker.Type {
	case txprofitabilitychecker.AcceptAllType:
		txProfitabilityChecker = txprofitabilitychecker.NewTxProfitabilityCheckerAcceptAll(state, cfg.IntervalAfterWhichBatchSentAnyway.Duration)
	case txprofitabilitychecker.BaseType:
		minReward := new(big.Int).Mul(cfg.Strategy.TxProfitabilityChecker.MinReward.Int, big.NewInt(encoding.TenToThePowerOf18))
		txProfitabilityChecker = txprofitabilitychecker.NewTxProfitabilityCheckerBase(ethMan, state, minReward, cfg.IntervalAfterWhichBatchSentAnyway.Duration, cfg.Strategy.TxProfitabilityChecker.RewardPercentageToAggregator)
	}

	seqAddress := ethMan.GetAddress()
	var chainID uint64
	if cfg.AllowNonRegistered {
		chainID = cfg.DefaultChainID
	} else {
		var err error
		chainID, err = getChainID(ctx, state, ethMan, seqAddress)
		if err != nil {
			cancel()
			return Sequencer{}, fmt.Errorf("failed to get chain id for the sequencer, err: %v", err)
		}
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
		log.Infof("waiting for the state to be synced, lastSyncedBatchNum: %d, lastEthBatchNum: %d", lastSyncedBatchNum, lastEthBatchNum)
		return
	}

	// 2. get pending txs from the pool
	txs, err := s.Pool.GetPendingTxs(s.ctx)
	if err != nil {
		log.Errorf("failed to get pending txs, err: %v", err)
		return
	}

	if len(txs) == 0 {
		log.Infof("transactions pool is empty, waiting for the new txs...")
		return
	}

	// 3. Run selection
	// init batch processor
	lastVirtualBatch, err := s.State.GetLastBatch(s.ctx, true)
	if err != nil {
		log.Errorf("failed to get last batch from the state, err: %v", err)
		return
	}
	bp, err := s.State.NewBatchProcessor(s.Address, lastVirtualBatch.Number().Uint64())
	if err != nil {
		log.Errorf("failed to create new batch processor, err: %v", err)
		return
	}

	// select txs
	selectedTxs, selectedTxsHashes, invalidTxsHashes, err := s.TxSelector.SelectTxs(bp, txs, s.Address)
	if err != nil {
		log.Errorf("failed to select txs, err: %v", err)
		return
	}

	if err = s.Pool.UpdateTxsState(s.ctx, invalidTxsHashes, pool.TxStateInvalid); err != nil {
		log.Errorf("failed to update txs state to invalid, err: %v", err)
		return
	}

	// 4. Is selection profitable?
	// check is it profitable to send selection
	isProfitable, aggregatorReward, err := s.TxProfitabilityChecker.IsProfitable(s.ctx, selectedTxs)
	if err != nil {
		log.Errorf("failed to check that txs are profitable or not, err: %v", err)
		return
	}
	if isProfitable && len(selectedTxs) > 0 {
		// YES: send selection to Ethereum
		sendBatchTx, err := s.EthMan.SendBatch(s.ctx, selectedTxs, aggregatorReward)
		if err != nil {
			log.Errorf("failed to send batch proposal to ethereum, err: %v", err)
			return
		}
		log.Infof("batch proposal sent successfully: %s", sendBatchTx.Hash().Hex())

		// update txs in the pool as selected
		err = s.Pool.UpdateTxsState(s.ctx, selectedTxsHashes, pool.TxStateSelected)
		if err != nil {
			log.Warnf("failed to update txs state to selected, err: %v", err)
		}
		log.Infof("finished updating selected transactions state in the pool")
	}
	// NO: discard selection and wait for the new batch
}

func getChainID(ctx context.Context, st state.State, ethMan etherman, seqAddress common.Address) (uint64, error) {
	const intervalToCheckSequencerRegistrationInSeconds = 3
	var (
		seq *state.Sequencer
		err error
	)
	for {
		seq, err = st.GetSequencer(ctx, seqAddress)
		if err == nil && seq != nil {
			return seq.ChainID.Uint64(), nil
		}
		if !errors.Is(err, state.ErrNotFound) {
			return 0, err
		}
		chainID, err := ethMan.GetCustomChainID()
		if err != nil {
			return 0, err
		}
		if chainID.Uint64() != 0 {
			return chainID.Uint64(), nil
		}

		log.Warnf("make sure the address %s has been registered in the smart contract as a sequencer, err: %v", seqAddress.Hex(), err)

		time.Sleep(intervalToCheckSequencerRegistrationInSeconds * time.Second)
	}
}
