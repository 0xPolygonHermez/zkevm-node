package aggregator

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-crypto/keccak256"
	"google.golang.org/grpc"
)

// Prime field. It is the prime number used as the order in our elliptic curve
const fr = "21888242871839275222246405745257275088548364400416034343698204186575808495617"

// Aggregator represents an aggregator
type Aggregator struct {
	cfg Config

	State          stateInterface
	EthTxManager   ethTxManager
	Ethman         etherman
	ZkProverClient pb.ZKProverServiceClient

	ProfitabilityChecker aggregatorTxProfitabilityChecker

	ctx    context.Context
	cancel context.CancelFunc
}

// NewAggregator creates a new aggregator
func NewAggregator(
	cfg Config,
	state stateInterface,
	ethTxManager ethTxManager,
	etherman etherman,
	zkProverClient pb.ZKProverServiceClient,
) (Aggregator, error) {
	ctx, cancel := context.WithCancel(context.Background())

	var profitabilityChecker aggregatorTxProfitabilityChecker
	switch cfg.TxProfitabilityCheckerType {
	case ProfitabilityBase:
		profitabilityChecker = NewTxProfitabilityCheckerBase(state, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration, cfg.TxProfitabilityMinReward.Int)
	case ProfitabilityAcceptAll:
		profitabilityChecker = NewTxProfitabilityCheckerAcceptAll(state, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration)
	}
	a := Aggregator{
		cfg: cfg,

		State:                state,
		EthTxManager:         ethTxManager,
		Ethman:               etherman,
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

	// define those vars here, bcs it can be used in case <-a.ctx.Done()
	var getProofCtx context.Context
	var getProofCtxCancel context.CancelFunc
	for {
		select {
		case <-time.After(a.cfg.IntervalToConsolidateState.Duration):
			// 1. check, if state is synced
			lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(a.ctx, nil)
			var lastVerifiedBatchNum uint64
			if err != nil && err != state.ErrNotFound {
				log.Warnf("failed to get last consolidated batch, err: %v", err)
				continue
			}
			if lastVerifiedBatch != nil {
				lastVerifiedBatchNum = lastVerifiedBatch.BatchNumber
			}
			lastConsolidatedEthBatchNum, err := a.Ethman.GetLatestVerifiedBatchNum()
			if err != nil {
				log.Warnf("failed to get last eth batch, err: %v", err)
				continue
			}
			if lastVerifiedBatchNum < lastConsolidatedEthBatchNum {
				log.Infof("waiting for the state to be synced, lastConsolidatedBatchNum: %d, lastEthConsolidatedBatchNum: %d",
					lastVerifiedBatchNum, lastConsolidatedEthBatchNum)
				continue
			}

			// 2. find next batch to consolidate
			delete(batchesSent, lastVerifiedBatchNum)

			batchToVerify, err := a.State.GetBatchByNumber(a.ctx, lastVerifiedBatchNum+1, nil)

			if err != nil {
				if err == state.ErrNotFound {
					log.Infof("there are no batches to consolidate")
					continue
				}
				log.Warnf("failed to get batch to consolidate, err: %v", err)
				continue
			}

			if batchesSent[batchToVerify.BatchNumber] {
				log.Infof("batch with number %d was already sent, but not yet consolidated by synchronizer",
					batchToVerify.BatchNumber)
				continue
			}

			// 3. check if it's profitable or not
			// check is it profitable to aggregate txs or not
			// pass matic collateral as zero here, bcs in smart contract fee for aggregator is not defined
			isProfitable, err := a.ProfitabilityChecker.IsProfitable(a.ctx, big.NewInt(0))
			if err != nil {
				log.Warnf("failed to check aggregator profitability, err: %v", err)
				continue
			}

			if !isProfitable {
				log.Infof("Batch %d is not profitable, matic collateral %d", batchToVerify.BatchNumber, big.NewInt(0))
				continue
			}

			// 4. send zki + txs to the prover
			stateRootConsolidated, err := a.State.GetStateRootByBatchNumber(a.ctx, lastVerifiedBatchNum, nil)
			if err != nil && err != state.ErrNotFound {
				log.Warnf("failed to get current state root, err: %v", err)
				continue
			}

			stateRootToConsolidate, err := a.State.GetStateRootByBatchNumber(a.ctx, batchToVerify.BatchNumber, nil)
			if err != nil {
				log.Warnf("failed to get state root to consolidate, err: %v", err)
				continue
			}

			rawTxs, err := state.EncodeTransactions(batchToVerify.Transactions)
			if err != nil {
				log.Warnf("failed to encode transactions, err: %v", err)
				continue
			}
			globalExitRoot := batchToVerify.GlobalExitRoot

			oldLocalExitRoot, err := a.State.GetLocalExitRootByBatchNumber(a.ctx, lastVerifiedBatchNum, nil)
			if err != nil {
				log.Warnf("failed to get local exit root for batch %d, err: %v", lastVerifiedBatchNum, err)
				continue
			}
			newLocalExitRoot := batchToVerify.LocalExitRoot
			// TODO: change this, once it will be clear, what db means
			db := map[string]string{
				"0540ae2a259cb9179561cffe6a0a3852a2c1806ad894ed396a2ef16e1f10e9c7": "0000000000000000000000000000000000000000000000056bc75e2d63100000",
				"061927dd2a72763869c1d5d9336a42d12a9a2f22809c9cf1feeb2a6d1643d950": "0000000000000000000000000000000000000000000000000000000000000000",
				"03ae74d1bbdff41d14f155ec79bb389db716160c1766a49ee9c9707407f80a11": "00000000000000000000000000000000000000000000000ad78ebc5ac6200000",
				"18d749d7bcc2bc831229c19256f9e933c08b6acdaff4915be158e34cbbc8a8e1": "0000000000000000000000000000000000000000000000000000000000000000",
			}

			batchChainIDByte := make([]byte, 4)   //nolint:gomnd
			blockTimestampByte := make([]byte, 8) //nolint:gomnd
			binary.BigEndian.PutUint64(blockTimestampByte, uint64(batchToVerify.Timestamp.Unix()))
			batchHashData := common.BytesToHash(keccak256.Hash(
				rawTxs,
				globalExitRoot[:],
				blockTimestampByte,
				batchToVerify.Coinbase[:],
				batchChainIDByte,
			))
			oldStateRoot := stateRootConsolidated
			newStateRoot := stateRootToConsolidate
			inputProver := &pb.InputProver{
				PublicInputs: &pb.PublicInputs{
					OldStateRoot:     oldStateRoot.String(),
					OldLocalExitRoot: oldLocalExitRoot.String(),
					NewStateRoot:     newStateRoot.String(),
					NewLocalExitRoot: newLocalExitRoot.String(),
					SequencerAddr:    batchToVerify.Coinbase.String(),
					BatchHashData:    batchHashData.String(),
					BatchNum:         uint32(batchToVerify.BatchNumber),
					EthTimestamp:     uint64(batchToVerify.Timestamp.Unix()),
				},
				GlobalExitRoot:    globalExitRoot.String(),
				BatchL2Data:       hex.EncodeToString(batchToVerify.BatchL2Data),
				Db:                db,
				ContractsBytecode: db,
			}

			genProofRequest := pb.GenProofRequest{Input: inputProver}

			// init connection to the prover
			var opts []grpc.CallOption
			resGenProof, err := a.ZkProverClient.GenProof(a.ctx, &genProofRequest, opts...)
			if err != nil {
				log.Errorf("failed to connect to the prover to gen proof, err: %v", err)
				continue
			}

			log.Debugf("Data sent to the prover: %+v", inputProver)
			genProofRes := resGenProof.GetResult()
			if genProofRes != pb.GenProofResponse_RESULT_GEN_PROOF_OK {
				log.Warnf("failed to get result from the prover, batchNumber: %d, err: %v", batchToVerify.BatchNumber)
				continue
			}
			genProofID := resGenProof.GetId()

			resGetProof := &pb.GetProofResponse{
				Result: -1,
			}
			getProofCtx, getProofCtxCancel = context.WithCancel(a.ctx)
			getProofClient, err := a.ZkProverClient.GetProof(getProofCtx)
			if err != nil {
				log.Warnf("failed to init getProofClient, batchNumber: %d, err: %v", batchToVerify.BatchNumber, err)
				continue
			}
			for resGetProof.Result != pb.GetProofResponse_RESULT_GET_PROOF_COMPLETED_OK {
				err = getProofClient.Send(&pb.GetProofRequest{
					Id: genProofID,
				})
				if err != nil {
					log.Warnf("failed to send get proof request to the prover, batchNumber: %d, err: %v", batchToVerify.BatchNumber, err)
					break
				}

				resGetProof, err = getProofClient.Recv()
				if err != nil {
					log.Warnf("failed to get proof from the prover, batchNumber: %d, err: %v", batchToVerify.BatchNumber, err)
					break
				}

				resGetProofState := resGetProof.GetResult()
				if resGetProofState == pb.GetProofResponse_RESULT_GET_PROOF_ERROR ||
					resGetProofState == pb.GetProofResponse_RESULT_GET_PROOF_COMPLETED_ERROR {
					log.Fatalf("failed to get a proof for batch, batch number %d", batchToVerify.BatchNumber)
				}
				if resGetProofState == pb.GetProofResponse_RESULT_GET_PROOF_INTERNAL_ERROR {
					log.Warnf("failed to generate proof for batch, batchNumber: %v, ResGetProofState: %v", batchToVerify.BatchNumber, resGetProofState)
					break
				}

				if resGetProofState == pb.GetProofResponse_RESULT_GET_PROOF_CANCEL {
					log.Warnf("proof generation was cancelled, batchNumber: %v", batchToVerify.BatchNumber)
					break
				}

				if resGetProofState == pb.GetProofResponse_RESULT_GET_PROOF_PENDING {
					// in this case aggregator will wait, to send another request
					time.Sleep(a.cfg.IntervalFrequencyToGetProofGenerationStateInSeconds.Duration)
				}
			}

			// getProofCtxCancel call closes the connection stream with the prover. This is the only way to close it by client
			getProofCtxCancel()

			if resGetProof.GetResult() != pb.GetProofResponse_RESULT_GET_PROOF_COMPLETED_OK {
				continue
			}

			// Calc inputHash
			batchNumberByte := make([]byte, 4) //nolint:gomnd
			binary.BigEndian.PutUint32(batchNumberByte, inputProver.PublicInputs.BatchNum)
			hash := keccak256.Hash(
				oldStateRoot[:],
				oldLocalExitRoot[:],
				newStateRoot[:],
				newLocalExitRoot[:],
				batchToVerify.Coinbase[:],
				batchHashData[:],
				batchNumberByte[:],
				blockTimestampByte[:],
			)
			frB, _ := new(big.Int).SetString(fr, encoding.Base10)
			inputHashMod := new(big.Int).Mod(new(big.Int).SetBytes(hash), frB)
			internalInputHash := inputHashMod.Bytes()

			// InputHash must match
			internalInputHashS := fmt.Sprintf("0x%064s", hex.EncodeToString(internalInputHash))
			publicInputsExtended := resGetProof.GetPublic()
			if resGetProof.GetPublic().InputHash != internalInputHashS {
				log.Error("inputHash received from the prover (", publicInputsExtended.InputHash,
					") doesn't match with the internal value: ", internalInputHashS)
				log.Debug("internalBatchHashData: ", batchHashData, " externalBatchHashData: ", publicInputsExtended.PublicInputs.BatchHashData)
				log.Debug("inputProver.PublicInputs.OldStateRoot: ", inputProver.PublicInputs.OldStateRoot)
				log.Debug("inputProver.PublicInputs.OldLocalExitRoot:", inputProver.PublicInputs.OldLocalExitRoot)
				log.Debug("inputProver.PublicInputs.NewStateRoot: ", inputProver.PublicInputs.NewStateRoot)
				log.Debug("inputProver.PublicInputs.NewLocalExitRoot: ", inputProver.PublicInputs.NewLocalExitRoot)
				log.Debug("inputProver.PublicInputs.SequencerAddr: ", inputProver.PublicInputs.SequencerAddr)
				log.Debug("inputProver.PublicInputs.BatchHashData: ", inputProver.PublicInputs.BatchHashData)
				log.Debug("inputProver.PublicInputs.BatchNum: ", inputProver.PublicInputs.BatchNum)
				log.Debug("inputProver.PublicInputs.EthTimestamp: ", inputProver.PublicInputs.EthTimestamp)
			}

			// 4. send proof + txs to the SC
			err = a.EthTxManager.VerifyBatch(batchToVerify.BatchNumber, resGetProof)
			if err != nil {
				log.Warnf("failed to send request to consolidate batch to ethereum, batch number: %d, err: %v",
					batchToVerify.BatchNumber, err)
				continue
			}
			batchesSent[batchToVerify.BatchNumber] = true
		case <-a.ctx.Done():
			getProofCtxCancel()
			return
		}
	}
}

// Stop stops the aggregator
func (a *Aggregator) Stop() {
	a.cancel()
}
