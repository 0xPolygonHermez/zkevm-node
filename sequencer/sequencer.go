package sequencer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/pricegetter"
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

	lastSentBatchNumber uint64

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
		priceGetter, err := pricegetter.NewClient(cfg.PriceGetter)
		if err != nil {
			cancel()
			return Sequencer{}, err
		}
		priceGetter.Start(ctx)
		txProfitabilityChecker = txprofitabilitychecker.NewTxProfitabilityCheckerBase(
			ethMan,
			state,
			priceGetter,
			minReward,
			cfg.IntervalAfterWhichBatchSentAnyway.Duration,
			cfg.Strategy.TxProfitabilityChecker.RewardPercentageToAggregator)
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
	var root []byte
	// Infinite for loop:
	for {
		root = s.tryProposeBatch(root)
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

func (s *Sequencer) tryProposeBatch(root []byte) []byte {
	// 1. Wait for synchronizer to sync last batch
	if !s.isSynced() {
		return nil
	}
	// 2. get pending txs from the pool
	txs, claimTxs, ok := s.getPendingTxs()
	if !ok {
		return nil
	}

	// 3. Run selection
	selectedTxsRes, ok := s.selectTxs(txs, claimTxs, root)
	if !ok {
		return nil
	}
	// 4. Send batch to ethereum
	ok = s.sendBatchToEthereum(selectedTxsRes)
	if !ok {
		return nil
	}
	return selectedTxsRes.NewRoot
}

func (s *Sequencer) isSynced() bool {
	lastSyncedBatchNum, err := s.State.GetLastBatchNumber(s.ctx, "")
	if err != nil {
		log.Errorf("failed to get last synced batch, err: %v", err)
		return false
	}
	lastEthBatchNum, err := s.State.GetLastBatchNumberSeenOnEthereum(s.ctx, "")
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

func (s *Sequencer) getPendingTxs() ([]pool.Transaction, []pool.Transaction, bool) {
	// get txs with claims
	claimsTxs, err := s.Pool.GetPendingTxs(s.ctx, true, amountOfPendingTxsRequested)
	if err != nil {
		log.Errorf("failed to get pending txs with claims, err: %v", err)
		return nil, nil, false
	}

	// get txs without claims
	txs, err := s.Pool.GetPendingTxs(s.ctx, false, amountOfPendingTxsRequested-uint64(len(claimsTxs)))
	if err != nil {
		log.Errorf("failed to get pending txs, err: %v", err)
		return nil, nil, false
	}

	if len(txs) == 0 && len(claimsTxs) == 0 {
		log.Infof("transactions pool is empty, waiting for the new txs...")
		return nil, nil, false
	}

	return txs, claimsTxs, true
}

func (s *Sequencer) selectTxs(txs, claimsTxs []pool.Transaction, root []byte) (txselector.SelectTxsOutput, bool) {
	root, batchNumber, err := s.chooseRoot(root)
	if err != nil {
		return txselector.SelectTxsOutput{}, false
	}
	bp, err := s.State.NewBatchProcessor(s.ctx, s.Address, root, "")
	if err != nil {
		log.Errorf("failed to create new batch processor, err: %v", err)
		return txselector.SelectTxsOutput{}, false
	}

	// select txs
	selectTxsRes, err := s.TxSelector.SelectTxs(s.ctx, txselector.SelectTxsInput{BatchProcessor: bp, PendingTxs: txs, PendingClaimsTxs: claimsTxs, SequencerAddress: s.Address})
	if err != nil {
		log.Errorf("failed to select txs, err: %v", err)
		return txselector.SelectTxsOutput{}, false
	}

	if err = s.Pool.UpdateTxsState(s.ctx, selectTxsRes.InvalidTxsHashes, pool.TxStateInvalid); err != nil {
		log.Errorf("failed to update txs state to invalid, err: %v", err)
		return txselector.SelectTxsOutput{}, false
	}
	selectTxsRes.BatchNumber = batchNumber
	return selectTxsRes, true
}

// chooseRoot the sequencer is deciding how to instantiate the batch processor
func (s *Sequencer) chooseRoot(prevRoot []byte) ([]byte, uint64, error) {
	lastVirtualBatch, err := s.State.GetLastBatch(s.ctx, true, "")
	if err != nil {
		log.Errorf("failed to get last batch from the state, err: %v", err)
		return nil, 0, err
	}
	lastVirtualBatchNumber := lastVirtualBatch.Header.Number.Uint64()
	lastVirtualBatchRoot := lastVirtualBatch.Header.Root[:]
	var isFromLastVirtualBatch bool

	// check if previous root is present, if not, take root from previous synced batch
	// if lastVirtualBatchNumber == s.lastSentBatchNumber, it means, that sequencer is synced
	// and can take root from synced batch
	if prevRoot == nil || lastVirtualBatchNumber == s.lastSentBatchNumber {
		isFromLastVirtualBatch = true
	} else if lastVirtualBatchNumber > s.lastSentBatchNumber {
		// in this case sequencer is trying to get batch by root
		// if root exist, it means sequencer can use root from the synced batch
		// if not exist, than batch processor initialization should be decided by param from the config
		_, err := s.State.GetLastBatchByStateRoot(s.ctx, prevRoot, "")
		if err != nil {
			if errors.Is(err, state.ErrNotFound) {
				if s.cfg.InitBatchProcessorIfDiffType == InitBatchProcessorIfDiffTypeSynced {
					isFromLastVirtualBatch = true
				}
			} else {
				log.Errorf("failed to get batch from the state by root, err: %v", err)
				return nil, 0, err
			}
		} else {
			isFromLastVirtualBatch = true
		}
	}

	var (
		root        []byte
		batchNumber uint64
	)

	if isFromLastVirtualBatch {
		root = lastVirtualBatchRoot
		batchNumber = lastVirtualBatchNumber + 1
	} else {
		root = prevRoot
		batchNumber = s.lastSentBatchNumber + 1
	}

	return root, batchNumber, nil
}

func isDataForEthTxTooBig(err error) bool {
	if strings.Contains(err.Error(), errGasRequiredExceedsAllowance) ||
		errors.As(err, &core.ErrOversizedData) ||
		strings.Contains(err.Error(), errContentLengthTooLarge) {
		return true
	}
	return false
}

func (s *Sequencer) sendBatchToEthereum(selectionRes txselector.SelectTxsOutput) bool {
	var isSent bool
	for !isSent {
		isProfitable, aggregatorReward, err := s.TxProfitabilityChecker.IsProfitable(s.ctx, selectionRes)
		if err != nil {
			log.Errorf("failed to check that txs are profitable or not, err: %v", err)
			return false
		}

		if isProfitable && (len(selectionRes.SelectedTxs) > 0 || len(selectionRes.SelectedClaimsTxs) > 0) {
			// YES: send selection to Ethereum
			sendBatchTx, err := s.EthMan.SendBatch(s.ctx, append(selectionRes.SelectedTxs, selectionRes.SelectedClaimsTxs...), aggregatorReward)
			if err != nil {
				if isDataForEthTxTooBig(err) {
					selectionRes.SelectedTxs, selectionRes.SelectedTxsHashes = cutSelectedTxs(selectionRes.SelectedTxs, selectionRes.SelectedTxsHashes)
					if len(selectionRes.SelectedTxs) == 0 {
						return false
					}
					continue
				}
				log.Errorf("failed to send batch proposal to ethereum, err: %v", err)
				return false
			}
			// update txs in the pool as selected
			selectedTxsHashes := append(selectionRes.SelectedTxsHashes, selectionRes.SelectedClaimsTxsHashes...)
			err = s.Pool.UpdateTxsState(s.ctx, selectedTxsHashes, pool.TxStateSelected)
			if err != nil {
				// it's fatal here, bcs txs are selected and sent to ethereum, but txs were not updated in the local db.
				// probably txs should be updated manually. If sequencer doesn't fail here, those txs will be sent again
				// and sequencer will lose tokens
				log.Fatalf("failed to update txs state to selected, selectedTxsHashes: %v, err: %v", selectedTxsHashes, err)
			}
			log.Infof("finished updating selected transactions state in the pool")
			s.lastSentBatchNumber = selectionRes.BatchNumber
			log.Infof("batch proposal sent successfully: %s", sendBatchTx.Hash().Hex())
			isSent = true
		} else {
			return false
		}
	}
	return true
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
		seq, err = st.GetSequencer(ctx, seqAddress, "")
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
