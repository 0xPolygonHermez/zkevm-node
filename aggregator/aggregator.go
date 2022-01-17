package aggregator

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

// Aggregator represents an aggregator
type Aggregator struct {
	cfg Config

	State          state.State
	BatchProcessor state.BatchProcessor
	EtherMan       etherman.EtherMan
	ZkProverClient proverclient.ZKProverClient

	ProfitabilityChecker TxProfitabilityChecker

	ctx    context.Context
	cancel context.CancelFunc
}

// NewAggregator creates a new aggregator
func NewAggregator(
	cfg Config,
	state state.State,
	ethMan etherman.EtherMan,
	zkProverClient proverclient.ZKProverClient,
) (Aggregator, error) {
	ctx, cancel := context.WithCancel(context.Background())

	var profitabilityChecker TxProfitabilityChecker
	switch cfg.TxProfitabilityCheckerType {
	case ProfitabilityBase:
		profitabilityChecker = NewTxProfitabilityCheckerBase(state, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration, cfg.TxProfitabilityMinReward.Int)
	case ProfitabilityAcceptAll:
		profitabilityChecker = NewTxProfitabilityCheckerAcceptAll(state, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration)
	}
	a := Aggregator{
		cfg: cfg,

		State:                state,
		EtherMan:             ethMan,
		ZkProverClient:       zkProverClient,
		ProfitabilityChecker: profitabilityChecker,

		ctx:    ctx,
		cancel: cancel,
	}

	return a, nil
}

// Start starts the aggregator
func (a *Aggregator) Start() {
	// init connection to the prover
	var opts []grpc.CallOption
	getProofClient, err := a.ZkProverClient.GenProof(a.ctx, opts...)
	if err != nil {
		log.Errorf("failed to connect to the prover, err: %v", err)
		return
	}

	// this is a batches, that were sent to ethereum to consolidate
	batchesSent := make(map[uint64]bool)

	for {
		select {
		case <-time.After(a.cfg.IntervalToConsolidateState.Duration):
			// 1. check, if state is synced
			lastConsolidatedBatch, err := a.State.GetLastBatch(a.ctx, false)
			if err != nil {
				log.Warnf("failed to get last consolidated batch, err: %v", err)
				continue
			}
			lastConsolidatedEthBatchNum, err := a.State.GetLastBatchNumberConsolidatedOnEthereum(a.ctx)
			if err != nil {
				log.Warnf("failed to get last eth batch, err: %v", err)
				continue
			}
			if lastConsolidatedBatch.BatchNumber < lastConsolidatedEthBatchNum {
				log.Infow("waiting for the state to be synced, lastConsolidatedBatchNum: %d, lastEthConsolidatedBatchNum: %d", lastConsolidatedBatch.BatchNumber, lastConsolidatedEthBatchNum)
				continue
			}

			// 2. find next batch to consolidate
			delete(batchesSent, lastConsolidatedBatch.BatchNumber)

			batchToConsolidate, err := a.State.GetBatchByNumber(a.ctx, lastConsolidatedBatch.BatchNumber+1)

			if err != nil {
				if err == pgx.ErrNoRows {
					log.Infof("there is no batches to consolidate")
					continue
				}
				log.Warnf("failed to get batch to consolidate, err: %v", err)
				continue
			}

			if batchesSent[batchToConsolidate.BatchNumber] {
				log.Infof("batch with number %d was already sent, but not yet consolidated by synchronizer",
					batchToConsolidate.BatchNumber)
				continue
			}

			// 3. check if it's profitable or not
			// check is it profitable to aggregate txs or not
			isProfitable, err := a.ProfitabilityChecker.IsProfitable(a.ctx, batchToConsolidate.MaticCollateral)
			if err != nil {
				log.Warnf("failed to check aggregator profitability, err: %v", err)
				continue
			}

			if !isProfitable {
				log.Info("Batch %d is not profitable, matic collateral %v", batchToConsolidate.BatchNumber, batchToConsolidate.MaticCollateral)
				continue
			}

			// 4. send zki + txs to the prover
			stateRootConsolidated, err := a.State.GetStateRootByBatchNumber(lastConsolidatedBatch.BatchNumber)
			if err != nil {
				log.Warnf("failed to get current state root, err: %v", err)
				continue
			}

			stateRootToConsolidate, err := a.State.GetStateRootByBatchNumber(batchToConsolidate.BatchNumber)
			if err != nil {
				log.Warnf("failed to get state root to consolidate, err: %v", err)
				continue
			}

			txsHashes := make([]string, 0, len(batchToConsolidate.Transactions))
			for _, tx := range batchToConsolidate.Transactions {
				txsHashes = append(txsHashes, tx.Hash().String())
			}

			// TODO: consider putting chain id to the batch, so we will get rid of additional request to db
			seq, err := a.State.GetSequencer(a.ctx, batchToConsolidate.Sequencer)
			if err != nil {
				log.Warnf("failed to get sequencer from the state, addr: %s, err: %v", seq.Address, err)
				continue
			}
			chainID := uint32(seq.ChainID.Uint64())

			// TODO: change this, once we have exit root
			fakeLastGlobalExitRoot, _ := new(big.Int).SetString("1234123412341234123412341234123412341234123412341234123412341234", 16)
			fakeKeys := map[string]string{
				"0540ae2a259cb9179561cffe6a0a3852a2c1806ad894ed396a2ef16e1f10e9c7": "0000000000000000000000000000000000000000000000056bc75e2d63100000",
				"061927dd2a72763869c1d5d9336a42d12a9a2f22809c9cf1feeb2a6d1643d950": "0000000000000000000000000000000000000000000000000000000000000000",
				"03ae74d1bbdff41d14f155ec79bb389db716160c1766a49ee9c9707407f80a11": "00000000000000000000000000000000000000000000000ad78ebc5ac6200000",
				"18d749d7bcc2bc831229c19256f9e933c08b6acdaff4915be158e34cbbc8a8e1": "0000000000000000000000000000000000000000000000000000000000000000",
			}
			inputProver := &proverclient.InputProver{
				Message: "calculate",
				PublicInputs: &proverclient.PublicInputs{
					OldStateRoot:     common.BytesToHash(stateRootConsolidated).String(),
					OldLocalExitRoot: fakeLastGlobalExitRoot.String(),
					NewStateRoot:     common.BytesToHash(stateRootToConsolidate).String(),
					NewLocalExitRoot: fakeLastGlobalExitRoot.String(),
					SequencerAddr:    batchToConsolidate.Sequencer.String(),
					BatchHashData:    batchToConsolidate.BatchHash.String(),
					ChainId:          chainID,
					BatchNum:         uint32(batchToConsolidate.BatchNumber),
				},
				GlobalExitRoot: fakeLastGlobalExitRoot.String(),
				Txs:            txsHashes,
				Keys:           fakeKeys,
			}
			err = getProofClient.Send(inputProver)
			if err != nil {
				log.Warnf("failed to send batch to the prover, batchNumber: %v, err: %v", batchToConsolidate.BatchNumber, err)
				continue
			}
			proofState, err := getProofClient.Recv()
			if err != nil {
				log.Warnf("failed to get proof from the prover, batchNumber: %v, err: %v", batchToConsolidate.BatchNumber, err)
				continue
			}
			// 4. send proof + txs to the SC
			batchNum := new(big.Int).SetUint64(batchToConsolidate.BatchNumber)
			h, err := a.EtherMan.ConsolidateBatch(batchNum, proofState.Proof)
			if err != nil {
				log.Warnf("failed to send request to consolidate batch to ethereum, batch number: %d, err: %v",
					batchToConsolidate.BatchNumber, err)
				continue
			}
			batchesSent[batchToConsolidate.BatchNumber] = true

			log.Infof("Batch %d consolidated: %s", batchToConsolidate.BatchNumber, h.Hash())
		case <-a.ctx.Done():
			return
		}
	}
}

// Stop stops the aggregator
func (a *Aggregator) Stop() {
	a.cancel()
}
