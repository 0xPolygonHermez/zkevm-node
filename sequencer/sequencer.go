package sequencer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy/txprofitabilitychecker"
	"github.com/hermeznetwork/hermez-core/sequencer/strategy/txselector"
	"github.com/hermeznetwork/hermez-core/state"
)

const (
	amountOfPendingTxsRequested = 30000
	percentageToCutSelectedTxs  = 80
	fullPercentage              = 100

	errGasRequiredExceedsAllowance = "gas required exceeds allowance"
	errContentLengthTooLarge       = "content length too large"
	errOversizedData               = "oversized data"
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg Config

	Pool    txPool
	State   stateInterface
	EthMan  etherman
	Address common.Address
	ChainID uint64

	txselector.TxSelector
	TxProfitabilityChecker txProfitabilityChecker

	ctx    context.Context
	cancel context.CancelFunc
}

// NewSequencer creates a new sequencer
func NewSequencer(cfg Config, pool txPool, state stateInterface, ethMan etherman) (Sequencer, error) {
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
		txProfitabilityChecker = txprofitabilitychecker.NewTxProfitabilityCheckerAcceptAll(ethMan, state, cfg.IntervalAfterWhichBatchSentAnyway.Duration)
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
	ticker := time.NewTicker(s.cfg.IntervalToProposeBatch.Duration)
	defer ticker.Stop()
	// Infinite for loop:
	for {
		s.tryProposeBatch()
		select {
		case <-ticker.C:
			// nothing
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
	if !s.isSynced() {
		return
	}
	// 2. get pending txs from the pool
	txs, ok := s.getPendingTxs()
	if !ok {
		return
	}
	// 3. Run selection
	selectedTxs, selectedTxsHashes, ok := s.selectTxs(txs)
	if !ok {
		return
	}
	// 4. Is selection profitable?
	// check is it profitable to send selection
	var isSent bool
	for !isSent {
		selectedTxs, selectedTxsHashes, isSent = s.sendBatchToEthereum(selectedTxs, selectedTxsHashes)
		if selectedTxs == nil || len(selectedTxs) == 0 {
			return
		}
	}
	// NO: discard selection and wait for the new batch
}

func (s *Sequencer) isSynced() bool {
	lastSyncedBatchNum, err := s.State.GetLastBatchNumber(s.ctx)
	if err != nil {
		log.Errorf("failed to get last synced batch, err: %v", err)
		return false
	}
	lastEthBatchNum, err := s.State.GetLastBatchNumberSeenOnEthereum(s.ctx)
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

func (s *Sequencer) getPendingTxs() ([]pool.Transaction, bool) {
	txs, err := s.Pool.GetPendingTxs(s.ctx, amountOfPendingTxsRequested)
	if err != nil {
		log.Errorf("failed to get pending txs, err: %v", err)
		return nil, false
	}

	if len(txs) == 0 {
		log.Infof("transactions pool is empty, waiting for the new txs...")
		return nil, false
	}

	return txs, true
}

func (s *Sequencer) selectTxs(txs []pool.Transaction) ([]*types.Transaction, []common.Hash, bool) {
	lastVirtualBatch, err := s.State.GetLastBatch(s.ctx, true)
	if err != nil {
		log.Errorf("failed to get last batch from the state, err: %v", err)
		return nil, nil, false
	}
	// init batch processor
	bp, err := s.State.NewBatchProcessor(s.ctx, s.Address, lastVirtualBatch.Header.Root[:])
	if err != nil {
		log.Errorf("failed to create new batch processor, err: %v", err)
		return nil, nil, false
	}

	// select txs
	selectedTxs, selectedTxsHashes, invalidTxsHashes, err := s.TxSelector.SelectTxs(s.ctx, bp, txs, s.Address)
	if err != nil {
		log.Errorf("failed to select txs, err: %v", err)
		return nil, nil, false
	}

	if err = s.Pool.UpdateTxsState(s.ctx, invalidTxsHashes, pool.TxStateInvalid); err != nil {
		log.Errorf("failed to update txs state to invalid, err: %v", err)
		return nil, nil, false
	}
	return selectedTxs, selectedTxsHashes, true
}

func isDataForEthTxTooBig(err error) bool {
	if strings.Contains(err.Error(), errGasRequiredExceedsAllowance) ||
		strings.Contains(err.Error(), errOversizedData) ||
		strings.Contains(err.Error(), errContentLengthTooLarge) {
		return true
	}
	return false
}

func (s *Sequencer) sendBatchToEthereum(selectedTxs []*types.Transaction, selectedTxsHashes []common.Hash) ([]*types.Transaction, []common.Hash, bool) {
	isProfitable, aggregatorReward, err := s.TxProfitabilityChecker.IsProfitable(s.ctx, selectedTxs)
	if err != nil {
		log.Errorf("failed to check that txs are profitable or not, err: %v", err)
		return nil, nil, false
	}
	if isProfitable && len(selectedTxs) > 0 {
		// YES: send selection to Ethereum
		sendBatchTx, err := s.EthMan.SendBatch(s.ctx, selectedTxs, aggregatorReward)
		if err != nil {
			if isDataForEthTxTooBig(err) {
				selectedTxs, selectedTxsHashes = cutSelectedTxs(selectedTxs, selectedTxsHashes)
				return selectedTxs, selectedTxsHashes, false
			}
			log.Errorf("failed to send batch proposal to ethereum, err: %v", err)
			return nil, nil, false
		}
		// update txs in the pool as selected
		err = s.Pool.UpdateTxsState(s.ctx, selectedTxsHashes, pool.TxStateSelected)
		if err != nil {
			// it's fatal here, bcs txs are selected and sent to ethereum, but txs were not updated in the local db.
			// probably txs should be updated manually. If sequencer don't fail here, those txs will be sent again
			// and sequencer will lose tokens
			log.Fatalf("failed to update txs state to selected, selectedTxsHashes: %v, err: %v", selectedTxsHashes, err)
		}
		log.Infof("finished updating selected transactions state in the pool")
		log.Infof("batch proposal sent successfully: %s", sendBatchTx.Hash().Hex())
	}
	return nil, nil, true
}

func cutSelectedTxs(selectedTxs []*types.Transaction, selectedTxsHashes []common.Hash) ([]*types.Transaction, []common.Hash) {
	cutSelectedTxs := len(selectedTxs) * percentageToCutSelectedTxs / fullPercentage
	selectedTxs = selectedTxs[:cutSelectedTxs]
	selectedTxsHashes = selectedTxsHashes[:cutSelectedTxs]
	return selectedTxs, selectedTxsHashes
}

func getChainID(ctx context.Context, st stateInterface, ethMan etherman, seqAddress common.Address) (uint64, error) {
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
