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
	"google.golang.org/grpc"
)

// Aggregator represents an aggregator
type Aggregator struct {
	cfg Config

	State                stateInterface
	EthTxManager         ethTxManager
	Ethman               etherman
	ProverClients        []proverClientInterface
	ProfitabilityChecker aggregatorTxProfitabilityChecker
}

// NewAggregator creates a new aggregator
func NewAggregator(
	cfg Config,
	stateInterface stateInterface,
	ethTxManager ethTxManager,
	etherman etherman,
	grpcClientConns []*grpc.ClientConn,
) (Aggregator, error) {
	var profitabilityChecker aggregatorTxProfitabilityChecker
	switch cfg.TxProfitabilityCheckerType {
	case ProfitabilityBase:
		profitabilityChecker = NewTxProfitabilityCheckerBase(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration, cfg.TxProfitabilityMinReward.Int)
	case ProfitabilityAcceptAll:
		profitabilityChecker = NewTxProfitabilityCheckerAcceptAll(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration)
	}

	proverClients := make([]proverClientInterface, 0, len(cfg.ProverURIs))
	ctx := context.Background()

	a := Aggregator{
		cfg: cfg,

		State:                stateInterface,
		EthTxManager:         ethTxManager,
		Ethman:               etherman,
		ProverClients:        proverClients,
		ProfitabilityChecker: profitabilityChecker,
	}

	for _, proverURI := range cfg.ProverURIs {
		proverClient := prover.NewClient(proverURI, cfg.IntervalFrequencyToGetProofGenerationState)
		proverClients = append(proverClients, proverClient)
		grpcClientConns = append(grpcClientConns, proverClient.Prover.Conn)
		log.Infof("Connected to prover %v", proverURI)

		// Check if prover is already working in a proof generation
		proof, err := stateInterface.GetWIPProofByProver(ctx, proverURI, nil)
		if err != nil && err != state.ErrNotFound {
			log.Errorf("Error while getting WIP proof for prover %v", proverURI)
			continue
		}

		if proof != nil {
			log.Infof("Resuming WIP proof generation for batchNumber %v in prover %v", proof.BatchNumber, *proof.Prover)
			go func() {
				a.resumeWIPProofGeneration(ctx, proof, proverClient)
			}()
		}
	}

	a.ProverClients = proverClients

	return a, nil
}

func (a *Aggregator) resumeWIPProofGeneration(ctx context.Context, proof *state.Proof, prover proverClientInterface) {
	err := a.getAndStoreProof(ctx, proof, prover)
	if err != nil {
		log.Warnf("Could not resume WIP Proof Generation for prover %v and batchNumber %v", *proof.Prover, proof.BatchNumber)
	}
}

// Start starts the aggregator
func (a *Aggregator) Start(ctx context.Context) {
	// define those vars here, bcs it can be used in case <-a.ctx.Done()
	tickerVerifyBatch := time.NewTicker(a.cfg.IntervalToConsolidateState.Duration)
	tickerSendVerifiedBatch := time.NewTicker(a.cfg.IntervalToConsolidateState.Duration)
	defer tickerVerifyBatch.Stop()
	defer tickerSendVerifiedBatch.Stop()

	for i := 0; i < len(a.ProverClients); i++ {
		go func() {
			for {
				a.tryVerifyBatch(ctx, tickerVerifyBatch)
			}
		}()
		time.Sleep(10 * time.Second) //nolint:gomnd
	}

	go func() {
		for {
			a.tryToSendVerifiedBatch(ctx, tickerSendVerifiedBatch)
		}
	}()
	// Wait until context is done
	<-ctx.Done()
}

func (a *Aggregator) tryToSendVerifiedBatch(ctx context.Context, ticker *time.Ticker) {
	log.Debug("checking if network is synced")
	for !a.isSynced(ctx) {
		log.Infof("waiting for synchronizer to sync...")
		waitTick(ctx, ticker)
		continue
	}
	log.Debug("checking if there is any consolidated batch to be verified")
	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Warnf("failed to get last consolidated batch, err: %v", err)
		waitTick(ctx, ticker)
		return
	} else if err == state.ErrNotFound {
		log.Debug("no consolidated batch found")
		waitTick(ctx, ticker)
		return
	}

	batchNumberToVerify := lastVerifiedBatch.BatchNumber + 1

	proof, err := a.State.GetGeneratedProofByBatchNumber(ctx, batchNumberToVerify, nil)
	if err != nil && err != state.ErrNotFound {
		log.Warnf("failed to get last proof for batch %v, err: %v", batchNumberToVerify, err)
		waitTick(ctx, ticker)
		return
	}

	if proof != nil && proof.Proof != nil {
		log.Infof("sending verified proof to the ethereum smart contract, batchNumber %d", batchNumberToVerify)
		err := a.EthTxManager.VerifyBatch(batchNumberToVerify, proof.Proof)
		if err != nil {
			log.Errorf("proof for the batch was NOT sent, batchNumber: %v, error: %w", batchNumberToVerify, err)
		} else {
			log.Infof("proof for the batch was sent, batchNumber: %v", batchNumberToVerify)
		}

		/*
			err := a.State.DeleteGeneratedProof(ctx, batchNumberToVerify, nil)
			if err != nil {
				log.Warnf("failed to delete generated proof for batchNumber %v, err: %v", batchNumberToVerify, err)
				return
			}
		*/
	} else {
		log.Debugf("no generated proof for batchNumber %v has been found", batchNumberToVerify)
		waitTick(ctx, ticker)
		return
	}
}

func (a *Aggregator) tryVerifyBatch(ctx context.Context, ticker *time.Ticker) {
	log.Info("checking if network is synced")
	for !a.isSynced(ctx) {
		log.Infof("waiting for synchronizer to sync...")
		waitTick(ctx, ticker)
		continue
	}
	log.Info("network is synced, getting batch to verify")
	batchToVerify, err := a.getBatchToVerify(ctx)
	if err != nil || batchToVerify == nil {
		waitTick(ctx, ticker)
		return
	}
	log.Infof("checking profitability to aggregate batch, batchNumber: %d", batchToVerify.BatchNumber)
	// pass matic collateral as zero here, bcs in smart contract fee for aggregator is not defined yet
	isProfitable, err := a.ProfitabilityChecker.IsProfitable(ctx, big.NewInt(0))
	if err != nil {
		log.Warnf("failed to check aggregator profitability, err: %v", err)
		waitTick(ctx, ticker)
		return
	}

	if !isProfitable {
		log.Infof("Batch %d is not profitable, matic collateral %d", batchToVerify.BatchNumber, big.NewInt(0))
		waitTick(ctx, ticker)
		return
	}

	log.Infof("sending zki + batch to the prover, batchNumber: %d", batchToVerify.BatchNumber)
	inputProver, err := a.buildInputProver(ctx, batchToVerify)
	if err != nil {
		log.Warnf("failed to build input prover, err: %v", err)
		waitTick(ctx, ticker)
		return
	}

	log.Infof("sending a batch to the prover, OLDSTATEROOT: %s, NEWSTATEROOT: %s, BATCHNUM: %d",
		inputProver.PublicInputs.OldStateRoot, inputProver.PublicInputs.NewStateRoot, inputProver.PublicInputs.BatchNum)

	var prover proverClientInterface
	var idleProverFound bool
	// Look for a free prover
	for _, prover = range a.ProverClients {
		if prover.IsIdle(ctx) {
			idleProverFound = true
			break
		}
	}

	if !idleProverFound {
		log.Warn("all provers are busy")
		waitTick(ctx, ticker)
		return
	}

	proverURI := prover.GetURI()
	proof := &state.Proof{BatchNumber: batchToVerify.BatchNumber, Prover: &proverURI, InputProver: inputProver}

	// Avoid other thread to process the same batch
	err = a.State.AddGeneratedProof(ctx, proof, nil)
	if err != nil {
		log.Warnf("failed to create proof generation mark, err: %v", err)
		waitTick(ctx, ticker)
		return
	}

	log.Infof("Prover %s is going to be used for batchNumber: %d", prover.GetURI(), batchToVerify.BatchNumber)

	genProofID, err := prover.GetGenProofID(ctx, inputProver)
	if err != nil {
		log.Warnf("failed to get gen proof id, err: %v", err)
		err2 := a.State.DeleteGeneratedProof(ctx, proof.BatchNumber, nil)
		if err2 != nil {
			log.Errorf("failed to delete proof generation mark after error, err: %v", err2)
		}
		waitTick(ctx, ticker)
		return
	}

	proof.ProofID = &genProofID

	// Avoid other thread to process the same batch
	err = a.State.UpdateGeneratedProof(ctx, proof, nil)
	if err != nil {
		log.Warnf("failed to update proof generation mark, err: %v", err)
		waitTick(ctx, ticker)
		return
	}

	log.Infof("Proof ID for batchNumber %d: %v", proof.BatchNumber, *proof.ProofID)

	err = a.getAndStoreProof(ctx, proof, prover)
	if err != nil {
		waitTick(ctx, ticker)
		return
	}
}

func (a *Aggregator) getAndStoreProof(ctx context.Context, proof *state.Proof, prover proverClientInterface) error {
	resGetProof, err := prover.GetResGetProof(ctx, *proof.ProofID, proof.BatchNumber)
	if err != nil {
		log.Warnf("failed to get proof from prover, err: %v", err)
		err2 := a.State.DeleteGeneratedProof(ctx, proof.BatchNumber, nil)
		if err2 != nil {
			log.Errorf("failed to delete proof generation mark after error, err: %v", err2)
			return err2
		}
		return err
	}

	proof.Proof = resGetProof

	a.compareInputHashes(proof.InputProver, proof.Proof)

	// Handle local exit root in the case of the mock prover
	if proof.Proof.Public.PublicInputs.NewLocalExitRoot == "0x17c04c3760510b48c6012742c540a81aba4bca2f78b9d14bfd2f123e2e53ea3e" {
		// This local exit root comes from the mock, use the one captured by the executor instead
		log.Warnf(
			"NewLocalExitRoot looks like a mock value, using value from executor instead: %v",
			proof.InputProver.PublicInputs.NewLocalExitRoot,
		)
		resGetProof.Public.PublicInputs.NewLocalExitRoot = proof.InputProver.PublicInputs.NewLocalExitRoot
	}

	// Store proof
	err = a.State.UpdateGeneratedProof(ctx, proof, nil)
	if err != nil {
		log.Warnf("failed to store generated proof, err: %v", err)
		err2 := a.State.DeleteGeneratedProof(ctx, proof.BatchNumber, nil)
		if err2 != nil {
			log.Errorf("failed to delete proof generation mark after error, err: %v", err2)
			return err2
		}
		return err
	}

	return nil
}

func (a *Aggregator) isSynced(ctx context.Context) bool {
	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Warnf("failed to get last consolidated batch, err: %v", err)
		return false
	}
	if lastVerifiedBatch == nil {
		return false
	}
	lastVerifiedEthBatchNum, err := a.Ethman.GetLatestVerifiedBatchNum()
	if err != nil {
		log.Warnf("failed to get last eth batch, err: %v", err)
		return false
	}
	if lastVerifiedBatch.BatchNumber < lastVerifiedEthBatchNum {
		log.Infof("waiting for the state to be synced, lastVerifiedBatchNum: %d, lastVerifiedEthBatchNum: %d",
			lastVerifiedBatch.BatchNumber, lastVerifiedEthBatchNum)
		return false
	}
	return true
}

func (a *Aggregator) getBatchToVerify(ctx context.Context) (*state.Batch, error) {
	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil {
		return nil, err
	}

	batchNumberToVerify := lastVerifiedBatch.BatchNumber + 1

	// Check if a prover is already working on this batch
	_, err = a.State.GetGeneratedProofByBatchNumber(ctx, batchNumberToVerify, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		return nil, err
	}

	for !errors.Is(err, state.ErrNotFound) {
		batchNumberToVerify++
		_, err = a.State.GetGeneratedProofByBatchNumber(ctx, batchNumberToVerify, nil)
		if err != nil && !errors.Is(err, state.ErrNotFound) {
			return nil, err
		}
	}

	batchToVerify, err := a.State.GetVirtualBatchByNumber(ctx, batchNumberToVerify, nil)
	if err != nil {
		if errors.Is(err, state.ErrNotFound) {
			log.Infof("there are no batches to consolidate")
			return nil, err
		}
		log.Warnf("failed to get batch to consolidate, err: %v", err)
		return nil, err
	}
	return batchToVerify, nil
}

func (a *Aggregator) buildInputProver(ctx context.Context, batchToVerify *state.Batch) (*pb.InputProver, error) {
	previousBatch, err := a.State.GetBatchByNumber(ctx, batchToVerify.BatchNumber-1, nil)
	if err != nil && err != state.ErrStateNotSynchronized {
		return nil, fmt.Errorf("failed to get previous batch, err: %v", err)
	}

	blockTimestampByte := make([]byte, 8) //nolint:gomnd
	binary.BigEndian.PutUint64(blockTimestampByte, uint64(batchToVerify.Timestamp.Unix()))
	batchHashData := common.BytesToHash(keccak256.Hash(
		batchToVerify.BatchL2Data,
		batchToVerify.GlobalExitRoot[:],
		blockTimestampByte,
		batchToVerify.Coinbase[:],
	))
	inputProver := &pb.InputProver{
		PublicInputs: &pb.PublicInputs{
			OldStateRoot:     previousBatch.StateRoot.String(),
			OldLocalExitRoot: previousBatch.LocalExitRoot.String(),
			NewStateRoot:     batchToVerify.StateRoot.String(),
			NewLocalExitRoot: batchToVerify.LocalExitRoot.String(),
			SequencerAddr:    batchToVerify.Coinbase.String(),
			BatchHashData:    batchHashData.String(),
			BatchNum:         uint32(batchToVerify.BatchNumber),
			EthTimestamp:     uint64(batchToVerify.Timestamp.Unix()),
			AggregatorAddr:   a.Ethman.GetPublicAddress().String(),
			ChainId:          a.cfg.ChainID,
		},
		GlobalExitRoot:    batchToVerify.GlobalExitRoot.String(),
		BatchL2Data:       hex.EncodeToString(batchToVerify.BatchL2Data),
		Db:                map[string]string{},
		ContractsBytecode: map[string]string{},
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
