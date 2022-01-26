package aggregator

import (
	"context"
	"encoding/binary"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/proverclient"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/iden3/go-iden3-crypto/keccak256"
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
	// this is a batches, that were sent to ethereum to consolidate
	batchesSent := make(map[uint64]bool)

	for {
		select {
		case <-time.After(a.cfg.IntervalToConsolidateState.Duration):
			// init connection to the prover
			var opts []grpc.CallOption
			getProofClient, err := a.ZkProverClient.GenProof(a.ctx, opts...)
			if err != nil {
				log.Errorf("failed to connect to the prover, err: %v", err)
				return
			}

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
				log.Infof("waiting for the state to be synced, lastConsolidatedBatchNum: %d, lastEthConsolidatedBatchNum: %d", lastConsolidatedBatch.BatchNumber, lastConsolidatedEthBatchNum)
				continue
			}

			// 2. find next batch to consolidate
			delete(batchesSent, lastConsolidatedBatch.BatchNumber)

			batchToConsolidate, err := a.State.GetBatchByNumber(a.ctx, lastConsolidatedBatch.BatchNumber+1)

			if err != nil {
				if err == state.ErrNotFound {
					log.Infof("there are no batches to consolidate")
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

			rawTxs := batchToConsolidate.RawTxsData
			// Split transactions
			var pos int64
			var txs []string
			const (
				headerByteLength = 2
				sLength          = 32
				rLength          = 32
				vLength          = 1
			)
			for pos < int64(len(rawTxs)) {
				length := rawTxs[pos : pos+1]
				str := hex.EncodeToString(length)
				num, err := strconv.ParseInt(str, hex.Base, encoding.BitSize64)
				if err != nil {
					log.Debug("error parsing header length: ", err)
					txs = []string{}
					break
				}
				// First byte is the length and must be ignored
				const c0 = 192 // 192 is c0. This value is defined by the rlp protocol
				num = num - c0 - 1

				fullDataTx := rawTxs[pos : pos+num+rLength+sLength+vLength+headerByteLength]

				pos = pos + num + rLength + sLength + vLength + headerByteLength

				txs = append(txs, "0x"+hex.EncodeToString(fullDataTx))
			}
			log.Debug("Txs: ", txs)

			// TODO: consider putting chain id to the batch, so we will get rid of additional request to db
			seq, err := a.State.GetSequencer(a.ctx, batchToConsolidate.Sequencer)
			if err != nil {
				log.Warnf("failed to get sequencer from the state, addr: %s, err: %v", seq.Address, err)
				continue
			}
			chainID := uint32(seq.ChainID.Uint64())

			// TODO: change this, once we have dynamic exit root
			globalExitRoot := common.HexToHash("0xa116e19a7984f21055d07b606c55628a5ffbf8ae1261c1e9f4e3a61620cf810a")
			oldLocalExitRoot := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
			newLocalExitRoot := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
			fakeKeys := map[string]string{
				"0540ae2a259cb9179561cffe6a0a3852a2c1806ad894ed396a2ef16e1f10e9c7": "0000000000000000000000000000000000000000000000056bc75e2d63100000",
				"061927dd2a72763869c1d5d9336a42d12a9a2f22809c9cf1feeb2a6d1643d950": "0000000000000000000000000000000000000000000000000000000000000000",
				"03ae74d1bbdff41d14f155ec79bb389db716160c1766a49ee9c9707407f80a11": "00000000000000000000000000000000000000000000000ad78ebc5ac6200000",
				"18d749d7bcc2bc831229c19256f9e933c08b6acdaff4915be158e34cbbc8a8e1": "0000000000000000000000000000000000000000000000000000000000000000",
			}
			batchHashData := common.BytesToHash(keccak256.Hash(batchToConsolidate.RawTxsData, globalExitRoot[:]))
			oldStateRoot := common.BytesToHash(stateRootConsolidated)
			newStateRoot := common.BytesToHash(stateRootToConsolidate)
			inputProver := &proverclient.InputProver{
				Message: "calculate",
				PublicInputs: &proverclient.PublicInputs{
					OldStateRoot:     oldStateRoot.String(),
					OldLocalExitRoot: oldLocalExitRoot.String(),
					NewStateRoot:     newStateRoot.String(),
					NewLocalExitRoot: newLocalExitRoot.String(),
					SequencerAddr:    batchToConsolidate.Sequencer.String(),
					BatchHashData:    batchHashData.String(),
					ChainId:          chainID,
					BatchNum:         uint32(batchToConsolidate.BatchNumber),
				},
				GlobalExitRoot: globalExitRoot.String(),
				Txs:            txs,
				Keys:           fakeKeys,
			}
			log.Debugf("Data sent to the prover: %+v", inputProver)
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

			// Calc inpuntHash
			batchChainIDByte := make([]byte, 4)
			batchNumberByte := make([]byte, 4)
			binary.BigEndian.PutUint32(batchChainIDByte, inputProver.PublicInputs.ChainId)
			binary.BigEndian.PutUint32(batchNumberByte, inputProver.PublicInputs.BatchNum)

			internalInputHash := keccak256.Hash(
				oldStateRoot[:],
				oldLocalExitRoot[:],
				newStateRoot[:],
				newLocalExitRoot[:],
				batchToConsolidate.Sequencer[:],
				batchHashData[:],
				batchChainIDByte[:],
				batchNumberByte[:],
			)

			// InputHash must match
			internalInputHashS := "0x" + hex.EncodeToString(internalInputHash)
			if proofState.Proof.PublicInputsExtended.InputHash != internalInputHashS {
				log.Error("inputHash received from the prover (", proofState.Proof.PublicInputsExtended.InputHash,
					") doesn't match with the internal value: ", internalInputHashS)
				log.Debug("internalBatchHashData: ", batchHashData, " externalBatchHashData: ", proofState.Proof.PublicInputsExtended.PublicInputs.BatchHashData)
				log.Debug("inputProver.PublicInputs.OldStateRoot: ", inputProver.PublicInputs.OldStateRoot)
				log.Debug("inputProver.PublicInputs.OldLocalExitRoot:", inputProver.PublicInputs.OldLocalExitRoot)
				log.Debug("inputProver.PublicInputs.NewStateRoot: ", inputProver.PublicInputs.NewStateRoot)
				log.Debug("inputProver.PublicInputs.NewLocalExitRoot: ", inputProver.PublicInputs.NewLocalExitRoot)
				log.Debug("inputProver.PublicInputs.SequencerAddr: ", inputProver.PublicInputs.SequencerAddr)
				log.Debug("inputProver.PublicInputs.BatchHashData: ", inputProver.PublicInputs.BatchHashData)
				log.Debug("inputProver.PublicInputs.ChainId: ", inputProver.PublicInputs.ChainId)
				log.Debug("inputProver.PublicInputs.BatchNum: ", inputProver.PublicInputs.BatchNum)
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
