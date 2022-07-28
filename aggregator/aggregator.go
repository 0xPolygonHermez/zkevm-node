package aggregator

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/prover"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-crypto/keccak256"
)

// Aggregator represents an aggregator
type Aggregator struct {
	cfg Config

	State                stateInterface
	EthTxManager         ethTxManager
	Ethman               etherman
	ProverClient         proverClient
	ProfitabilityChecker aggregatorTxProfitabilityChecker

	lastVerifiedBatchNum uint64
	batchesSent          map[uint64]bool
}

// NewAggregator creates a new aggregator
func NewAggregator(
	cfg Config,
	state stateInterface,
	ethTxManager ethTxManager,
	etherman etherman,
	zkProverClient pb.ZKProverServiceClient,
) (Aggregator, error) {
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
		ProverClient:         prover.NewClient(zkProverClient),
		ProfitabilityChecker: profitabilityChecker,

		batchesSent: make(map[uint64]bool),
	}

	return a, nil
}

// Start starts the aggregator
func (a *Aggregator) Start(ctx context.Context) {
	// define those vars here, bcs it can be used in case <-a.ctx.Done()
	ticker := time.NewTicker(a.cfg.IntervalToConsolidateState.Duration)
	defer ticker.Stop()
	for {
		a.tryVerifyBatch(ctx, ticker)
	}
}

func (a *Aggregator) tryVerifyBatch(ctx context.Context, ticker *time.Ticker) {
	// 1. check, if state is synced
	for !a.isSynced(ctx) {
		log.Infof("waiting for synchronizer to sync...")
		waitTick(ctx, ticker)
		continue
	}
	// 2. find next batch to verify
	batchToVerify, err := a.getBatchToVerify(ctx)
	if err != nil {
		waitTick(ctx, ticker)
		return
	}
	// 3. check if it's profitable or not
	// check is it profitable to aggregate txs or not
	// pass matic collateral as zero here, bcs in smart contract fee for aggregator is not defined yet
	isProfitable, err := a.ProfitabilityChecker.IsProfitable(ctx, big.NewInt(0))
	if err != nil {
		log.Warnf("failed to check aggregator profitability, err: %v", err)
		return
	}

	if !isProfitable {
		log.Infof("Batch %d is not profitable, matic collateral %d", batchToVerify.BatchNumber, big.NewInt(0))
		return
	}

	// 4. send zki + txs to the prover
	inputProver, err := a.buildInputProver(ctx, batchToVerify)
	if err != nil {
		log.Warnf("failed to build input prover, err: %v", err)
		return
	}

	genProofID, err := a.ProverClient.GetGenProofID(ctx, inputProver)
	if err != nil {
		log.Warnf("failed to get gen proof id, err: %v", err)
		return
	}

	resGetProof, err := a.ProverClient.GetResGetProof(ctx, genProofID, batchToVerify.BatchNumber)
	if err != nil {
		log.Warnf("failed to get proof from prover, err: %v", err)
		return
	}
	a.compareInputHashes(inputProver, resGetProof)

	// 4. send proof + txs to the SC
	err = a.EthTxManager.VerifyBatch(batchToVerify.BatchNumber, resGetProof)
	if err != nil {
		log.Warnf("failed to send request to consolidate batch to ethereum, batch number: %d, err: %v",
			batchToVerify.BatchNumber, err)
		return
	}
	a.batchesSent[batchToVerify.BatchNumber] = true
}

func (a *Aggregator) isSynced(ctx context.Context) bool {
	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Warnf("failed to get last consolidated batch, err: %v", err)
		return false
	}
	if lastVerifiedBatch != nil {
		a.lastVerifiedBatchNum = lastVerifiedBatch.BatchNumber
	}
	lastConsolidatedEthBatchNum, err := a.Ethman.GetLatestVerifiedBatchNum()
	if err != nil {
		log.Warnf("failed to get last eth batch, err: %v", err)
		return false
	}
	if a.lastVerifiedBatchNum < lastConsolidatedEthBatchNum {
		log.Infof("waiting for the state to be synced, lastConsolidatedBatchNum: %d, lastEthConsolidatedBatchNum: %d",
			a.lastVerifiedBatchNum, lastConsolidatedEthBatchNum)
		return false
	}
	return true
}

func (a *Aggregator) getBatchToVerify(ctx context.Context) (*state.Batch, error) {
	delete(a.batchesSent, a.lastVerifiedBatchNum)

	batchToVerify, err := a.State.GetVirtualBatchByNumber(ctx, a.lastVerifiedBatchNum+1, nil)

	if err != nil {
		if errors.Is(err, state.ErrNotFound) {
			log.Infof("there are no batches to consolidate")
			return nil, err
		}
		log.Warnf("failed to get batch to consolidate, err: %v", err)
		return nil, err
	}

	if a.batchesSent[batchToVerify.BatchNumber] {
		log.Infof("batch with number %d was already sent, but not yet consolidated by synchronizer",
			batchToVerify.BatchNumber)
		return nil, nil
	}

	return batchToVerify, nil
}

func (a *Aggregator) buildInputProver(ctx context.Context, batchToVerify *state.Batch) (*pb.InputProver, error) {
	oldStateRoot, err := a.State.GetStateRootByBatchNumber(ctx, a.lastVerifiedBatchNum, nil)
	if err != nil && err != state.ErrNotFound {
		return nil, fmt.Errorf("failed to get current state root, err: %v", err)
	}

	newStateRoot, err := a.State.GetStateRootByBatchNumber(ctx, batchToVerify.BatchNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get state root to consolidate, err: %v", err)
	}

	rawTxs, err := state.EncodeTransactions(batchToVerify.Transactions)
	if err != nil {
		return nil, fmt.Errorf("failed to encode transactions, err: %v", err)
	}
	globalExitRoot := batchToVerify.GlobalExitRoot

	oldLocalExitRoot, err := a.State.GetLocalExitRootByBatchNumber(ctx, a.lastVerifiedBatchNum, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get local exit root for batch %d, err: %v", a.lastVerifiedBatchNum, err)
	}
	newLocalExitRoot := batchToVerify.LocalExitRoot
	db := map[string]string{}

	blockTimestampByte := make([]byte, 8) //nolint:gomnd
	binary.BigEndian.PutUint64(blockTimestampByte, uint64(batchToVerify.Timestamp.Unix()))
	batchHashData := common.BytesToHash(keccak256.Hash(
		rawTxs,
		globalExitRoot[:],
		blockTimestampByte,
		batchToVerify.Coinbase[:],
	))
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

	return inputProver, nil
}

func (a *Aggregator) compareInputHashes(ip *pb.InputProver, resGetProof *pb.GetProofResponse) {
	// Calc inputHash
	batchNumberByte := make([]byte, 4) //nolint:gomnd
	binary.BigEndian.PutUint32(batchNumberByte, ip.PublicInputs.BatchNum)
	blockTimestampByte := make([]byte, 8) //nolint:gomnd
	binary.BigEndian.PutUint64(blockTimestampByte, ip.PublicInputs.EthTimestamp)
	hash := keccak256.Hash(
		[]byte(ip.PublicInputs.OldStateRoot)[:],
		[]byte(ip.PublicInputs.OldLocalExitRoot)[:],
		[]byte(ip.PublicInputs.NewStateRoot)[:],
		[]byte(ip.PublicInputs.NewLocalExitRoot)[:],
		[]byte(ip.PublicInputs.SequencerAddr)[:],
		[]byte(ip.PublicInputs.BatchHashData)[:],
		batchNumberByte[:],
		blockTimestampByte[:],
	)
	// Prime field. It is the prime number used as the order in our elliptic curve
	const fr = "21888242871839275222246405745257275088548364400416034343698204186575808495617"
	frB, _ := new(big.Int).SetString(fr, encoding.Base10)
	inputHashMod := new(big.Int).Mod(new(big.Int).SetBytes(hash), frB)
	internalInputHash := inputHashMod.Bytes()

	// InputHash must match
	internalInputHashS := fmt.Sprintf("0x%064s", hex.EncodeToString(internalInputHash))
	publicInputsExtended := resGetProof.GetPublic()
	if resGetProof.GetPublic().InputHash != internalInputHashS {
		log.Error("inputHash received from the prover (", publicInputsExtended.InputHash,
			") doesn't match with the internal value: ", internalInputHashS)
		log.Debug("internalBatchHashData: ", ip.PublicInputs.BatchHashData, " externalBatchHashData: ", publicInputsExtended.PublicInputs.BatchHashData)
		log.Debug("inputProver.PublicInputs.OldStateRoot: ", ip.PublicInputs.OldStateRoot)
		log.Debug("inputProver.PublicInputs.OldLocalExitRoot:", ip.PublicInputs.OldLocalExitRoot)
		log.Debug("inputProver.PublicInputs.NewStateRoot: ", ip.PublicInputs.NewStateRoot)
		log.Debug("inputProver.PublicInputs.NewLocalExitRoot: ", ip.PublicInputs.NewLocalExitRoot)
		log.Debug("inputProver.PublicInputs.SequencerAddr: ", ip.PublicInputs.SequencerAddr)
		log.Debug("inputProver.PublicInputs.BatchHashData: ", ip.PublicInputs.BatchHashData)
		log.Debug("inputProver.PublicInputs.BatchNum: ", ip.PublicInputs.BatchNum)
		log.Debug("inputProver.PublicInputs.EthTimestamp: ", ip.PublicInputs.EthTimestamp)
	}
}

func waitTick(ctx context.Context, ticker *time.Ticker) {
	select {
	case <-ticker.C:
		// nothing
	case <-ctx.Done():
		return
	}
}
