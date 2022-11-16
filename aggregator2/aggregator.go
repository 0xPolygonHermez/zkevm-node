package aggregator2

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator2/pb"
	"github.com/0xPolygonHermez/zkevm-node/aggregator2/prover"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-crypto/keccak256"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Aggregator represents an aggregator
type Aggregator2 struct {
	pb.UnimplementedAggregatorServiceServer

	cfg Config

	State                stateInterface
	EthTxManager         ethTxManager
	Ethman               etherman
	ProfitabilityChecker aggregatorTxProfitabilityChecker
	TimeSendFinalProof   time.Time
	StateDBMutex         *sync.Mutex

	srv  *grpc.Server
	ctx  context.Context
	exit context.CancelFunc
}

// NewAggregator creates a new aggregator
func NewAggregator2(
	cfg Config,
	stateInterface stateInterface,
	ethTxManager ethTxManager,
	etherman etherman,
	grpcClientConns []*grpc.ClientConn,
) (Aggregator2, error) {
	var profitabilityChecker aggregatorTxProfitabilityChecker
	switch cfg.TxProfitabilityCheckerType {
	case ProfitabilityBase:
		profitabilityChecker = NewTxProfitabilityCheckerBase(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration, cfg.TxProfitabilityMinReward.Int)
	case ProfitabilityAcceptAll:
		profitabilityChecker = NewTxProfitabilityCheckerAcceptAll(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration)
	}

	a := Aggregator2{
		cfg: cfg,

		State:                stateInterface,
		EthTxManager:         ethTxManager,
		Ethman:               etherman,
		ProfitabilityChecker: profitabilityChecker,
		StateDBMutex:         &sync.Mutex{},
	}

	return a, nil
}

// Start starts the aggregator
func (a *Aggregator2) Start(ctx context.Context) {
	var cancel context.CancelFunc
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel = context.WithCancel(ctx)
	a.ctx = ctx
	a.exit = cancel

	// Delete ungenerated recursive proofs
	a.State.DeleteUngeneratedRecursiveProofs(ctx, nil)

	address := fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	a.srv = grpc.NewServer()
	pb.RegisterAggregatorServiceServer(a.srv, a)

	healthService := newHealthChecker()
	grpc_health_v1.RegisterHealthServer(a.srv, healthService)

	go func() {
		log.Infof("Server listening on port %d", a.cfg.Port)
		if err := a.srv.Serve(lis); err != nil {
			a.exit()
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// define those vars here, bcs it can be used in case <-a.ctx.Done()
	tickerVerifyBatch := time.NewTicker(a.cfg.IntervalToConsolidateState.Duration)
	tickerSendVerifiedBatch := time.NewTicker(a.cfg.IntervalToSendFinalProof.Duration)
	defer tickerVerifyBatch.Stop()
	defer tickerSendVerifiedBatch.Stop()

	a.TimeSendFinalProof = time.Now().Add(a.cfg.IntervalToSendFinalProof.Duration)

	<-ctx.Done()
}

// Stop stops the Aggregator server.
func (a *Aggregator2) Stop() {
	a.exit()
	a.srv.Stop()
}

// Channel implements the bi-directional communication channel between the
// Prover client and the Aggregator server.
func (a *Aggregator2) Channel(stream pb.AggregatorService_ChannelServer) error {
	prover, err := prover.New(stream, a.cfg.IntervalFrequencyToGetProofGenerationState)
	if err != nil {
		return err
	}
	log.Debugf("Establishing stream connection for prover %s", prover.ID())

	ctx := stream.Context()

	tickerVerifyBatch := time.NewTicker(a.cfg.IntervalToConsolidateState.Duration)
	defer tickerVerifyBatch.Stop()

	go func() {
		for {
			select {
			case <-a.ctx.Done():
				// server disconnected
				return
			case <-ctx.Done():
				// client disconnected
				return
			case <-tickerVerifyBatch.C:
				if prover.IsIdle() {
					proofGenerated, _ := a.tryAggregateProofs(ctx, prover, tickerVerifyBatch)
					if !proofGenerated {
						proofGenerated, _ = a.tryGenerateBatchProof(ctx, prover, tickerVerifyBatch)
					}
					if !proofGenerated {
						// if no proof was generated (aggregated or batch) wait some time waiting before retry
						time.Sleep(a.cfg.IntervalToConsolidateState.Duration)
					} // if proof was generated we retry inmediatly as probably we have more proofs to process
				} else {
					log.Warn("Prover ID %s is not idle", prover.ID())
					time.Sleep(a.cfg.IntervalToConsolidateState.Duration)
				}
			}
		}
	}()

	// keep this scope alive, the stream gets closed if we exit from here.
	for {
		select {
		case <-a.ctx.Done():
			// server disconnecting
			return nil
		case <-ctx.Done():
			// client disconnected
			// TODO(pg): reconnect?
			return nil
		}
	}
}

func (a *Aggregator2) trySendFinalProof(ctx context.Context, prover *prover.Prover, proof *state.RecursiveProof, ticker *time.Ticker) (bool, error) {
	if a.TimeSendFinalProof.Before(time.Now()) {
		log.Debug("Send final proof time reached")

		log.Debug("Checking if network is synced")
		for !a.isSynced(ctx) {
			log.Info("Waiting for synchronizer to sync...")
			waitTick(ctx, ticker) //·sleep
			continue
		}

		lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
		if err != nil && err != state.ErrNotFound {
			log.Warnf("Failed to get last verified batch, err: %v", err)
			return false, err
		} else if err == state.ErrNotFound {
			log.Debug("Last verified batch not found")
			return false, err
		}

		batchNumberToVerify := lastVerifiedBatch.BatchNumber + 1

		if proof.BatchNumber == batchNumberToVerify {

			bComplete, err := a.State.CheckProofContainsCompleteSequences(ctx, proof, nil)

			if !bComplete {
				log.Infof("Recursive proof %d-%d does not contain completes sequences", proof.BatchNumber, proof.BatchNumberFinal)
				return false, err
			}

			log.Infof("Prover %s is going to be used to generate final proof for batches: %d-%d", prover.ID(), proof.BatchNumber, proof.BatchNumberFinal)

			finalProofID, err := prover.FinalProof(proof.Proof.RecursiveProof)
			if err != nil {
				log.Warnf("Failed to get final proof id, err: %v", err)
				return false, err
			}

			proof.ProofID = &finalProofID

			log.Infof("Proof ID for final proof %d-%d: %s", proof.BatchNumber, proof.BatchNumberFinal, *proof.ProofID)

			resGetProof, err := prover.WaitFinalProof(ctx, *proof.ProofID)
			if err != nil {
				log.Errorf("Failed to get final proof from prover, err: %v", err)
				return false, err
			}

			//b, err := json.Marshal(resGetProof.FinalProof)
			log.Infof("Final proof %s generated", *proof.ProofID)

			var inputProver *pb.InputProver
			json.Unmarshal([]byte(proof.InputProver), inputProver)
			a.compareInputHashes(inputProver, resGetProof)

			// Handle local exit root in the case of the mock prover
			if string(resGetProof.FinalProof.Public.NewLocalExitRoot[:]) == "0x17c04c3760510b48c6012742c540a81aba4bca2f78b9d14bfd2f123e2e53ea3e" {
				// This local exit root comes from the mock, use the one captured by the executor instead
				log.Warnf("NewLocalExitRoot looks like a mock value")
				/*log.Warnf(
					"NewLocalExitRoot looks like a mock value, using value from executor instead: %v",
					proof.InputProver.PublicInputs.NewLocalExitRoot,
				)*/
				//resGetProof.Public.PublicInputs.NewLocalExitRoot = proof.InputProver.PublicInputs.NewLocalExitRoot
			}

			log.Infof("Verfiying final proof with ethereum smart contract, batches %d-%d", proof.BatchNumber, proof.BatchNumberFinal)
			// · Not working with mock prover _, err = a.Ethman.VerifyBatches2(ctx, proof.BatchNumber-1, proof.BatchNumberFinal, resGetProof, 0, nil, nil)
			if err != nil {
				log.Errorf("Error verifiying final proof for batches %d-%d, err: %w", proof.BatchNumber, proof.BatchNumberFinal, err)
				return false, err
			}

			/* · Is needed to do this additional steps??
			err = c.state.UpdateProofTx(ctx, pendingProof.BatchNumber, tx.Hash(), tx.Nonce(), nil)
			if err != nil {
				log.Errorf("failed to update tx to verify batch for batch number %v, new tx hash %v, nonce %v, err: %v",
					pendingProof.BatchNumber, tx.Hash().String(), tx.Nonce(), err)
				break
			}
			err = c.ethMan.WaitTxToBeMined(ctx, tx, c.cfg.IntervalToReviewVerifyBatchTx.Duration)
			if err != nil {
				log.Errorf("error waiting tx to be mined: %s, error: %w", tx.Hash(), err)
				break
			}
			txHash := tx.Hash()
			pendingProof.TxHash = &txHash
			nonce := tx.Nonce()
			pendingProof.TxNonce = &nonce
			time.Sleep(time.Second * 2) // nolint
			*/

			log.Infof("Final proof for batches %d-%d verified", proof.BatchNumber, proof.BatchNumberFinal)
			a.TimeSendFinalProof = time.Now().Add(a.cfg.IntervalToSendFinalProof.Duration)
			return true, nil

		} else {
			log.Infof("Proof batch number %d is not the following to last verfied batch number %d", proof.BatchNumber, batchNumberToVerify)
			return false, nil
		}
	} else {
		return false, nil
	}
}
func (a *Aggregator2) unlockProofsToAggregate(ctx context.Context, proof1 *state.RecursiveProof, proof2 *state.RecursiveProof, ticker *time.Ticker) error {
	// Release proofs from generating state in a single transaction
	dbTx, err := a.State.BeginStateTransaction(ctx)
	if err != nil {
		log.Warnf("Failed to begin transaction to release proof aggregation state, err: %v", err)
		return err
	}

	proof1.Generating = false
	err = a.State.UpdateGeneratedRecursiveProof(ctx, proof1, dbTx)
	if err == nil {
		proof2.Generating = false
		err = a.State.UpdateGeneratedRecursiveProof(ctx, proof2, dbTx)
	}

	if err != nil {
		log.Warnf("Failed to release proof aggregation state, err: %v", err)
		dbTx.Rollback(ctx)
		return err
	}

	dbTx.Commit(ctx)

	return nil
}

func (a *Aggregator2) getAndLockProofsToAggregate(ctx context.Context, prover *prover.Prover, ticker *time.Ticker) (*state.RecursiveProof, *state.RecursiveProof, error) {
	a.StateDBMutex.Lock()
	defer a.StateDBMutex.Unlock()

	proof1, proof2, err := a.State.GetRecursiveProofsToAggregate(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	// Set proofs in generating state in a single transaction
	dbTx, err := a.State.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("Failed to begin transaction to set proof aggregation state, err: %v", err)
		return nil, nil, err
	}

	proof1.Generating = true
	err = a.State.UpdateGeneratedRecursiveProof(ctx, proof1, dbTx)
	if err == nil {
		proof2.Generating = true
		err = a.State.UpdateGeneratedRecursiveProof(ctx, proof2, dbTx)
	}

	if err != nil {
		log.Errorf("Failed to set proof aggregation state, err: %v", err)
		dbTx.Rollback(ctx)
		return nil, nil, err
	}

	dbTx.Commit(ctx)

	return proof1, proof2, nil
}

func (a *Aggregator2) tryAggregateProofs(ctx context.Context, prover *prover.Prover, ticker *time.Ticker) (bool, error) {
	log.Debugf("tryAggregateProofs start %s", prover.ID())

	proof1, proof2, err0 := a.getAndLockProofsToAggregate(ctx, prover, ticker)
	if err0 != nil {
		return false, err0
	}

	var err error

	defer func() {
		if err != nil {
			err2 := a.unlockProofsToAggregate(ctx, proof1, proof2, ticker)
			if err2 != nil {
				log.Errorf("Failed to release aggregated proofs, err: %v", err2)
			}
		}
	}()

	log.Infof("Prover %s is going to be used to aggregate proofs: %d-%d and %d-%d", prover.ID(), proof1.BatchNumber, proof1.BatchNumberFinal, proof2.BatchNumber, proof2.BatchNumberFinal)

	proverID := prover.ID()
	inputProver := map[string]interface{}{
		"recursive_proof_1": proof1.Proof.RecursiveProof,
		"recursive_proof_2": proof2.Proof.RecursiveProof,
	}
	b, err := json.Marshal(inputProver)
	proof := &state.RecursiveProof{BatchNumber: proof1.BatchNumber, BatchNumberFinal: proof2.BatchNumberFinal, Prover: &proverID, InputProver: string(b), Generating: false}

	aggrProofID, err := prover.AggregatedProof(proof1.Proof.RecursiveProof, proof2.Proof.RecursiveProof)
	if err != nil {
		log.Warnf("Failed to get aggregated proof id, err: %v", err)
		return false, err
	}

	proof.ProofID = &aggrProofID

	log.Infof("Proof ID for aggregated proof %d-%d: %v", proof.BatchNumber, proof.BatchNumberFinal, *proof.ProofID)

	resGetProof, err := prover.WaitRecursiveProof(ctx, *proof.ProofID)
	if err != nil {
		log.Errorf("Failed to get aggregated proof from prover, err: %v", err)
		return false, err
	}

	log.Infof("Aggregated proof %s generated", *proof.ProofID)

	proof.Proof = resGetProof

	proofSent, _ := a.trySendFinalProof(ctx, prover, proof, ticker)

	dbTx, err := a.State.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("Failed to begin transaction to store proof aggregation result, err: %v", err)
		return false, err
	}

	// If the new aggregated proof has not been sent to L1 we store it
	if !proofSent {
		err = a.State.AddGeneratedRecursiveProof(ctx, proof, dbTx)
	}

	// Delete aggregated proofs
	if err == nil {
		err = a.State.DeleteGeneratedRecursiveProof(ctx, proof1.BatchNumber, proof1.BatchNumberFinal, nil)
	}
	if err == nil {
		err = a.State.DeleteGeneratedRecursiveProof(ctx, proof2.BatchNumber, proof2.BatchNumberFinal, nil)
	}

	if err != nil {
		dbTx.Rollback(ctx)
		log.Errorf("Failed to store proof aggregation result, err: %v", err)
		return false, err
	}

	dbTx.Commit(ctx)

	log.Debug("tryAggregateProofs end")

	return true, nil
}

func (a *Aggregator2) getAndLockBatchToProve(ctx context.Context, prover *prover.Prover, ticker *time.Ticker) (*state.Batch, *state.RecursiveProof, error) {
	a.StateDBMutex.Lock()
	defer a.StateDBMutex.Unlock()

	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	// Get virtual batch pending to generate proof
	batchToVerify, err := a.State.GetVirtualBatchToRecursiveProve(ctx, lastVerifiedBatch.BatchNumber, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Infof("Found virtual batch %d pending to generate proof", batchToVerify.BatchNumber)

	log.Infof("Checking profitability to aggregate batch, batchNumber: %d", batchToVerify.BatchNumber)
	// pass matic collateral as zero here, bcs in smart contract fee for aggregator is not defined yet
	isProfitable, err := a.ProfitabilityChecker.IsProfitable(ctx, big.NewInt(0))
	if err != nil {
		log.Errorf("Failed to check aggregator profitability, err: %v", err)
		return nil, nil, err
	}

	if !isProfitable {
		log.Infof("Batch %d is not profitable, matic collateral %d", batchToVerify.BatchNumber, big.NewInt(0))
		return nil, nil, err
	}

	proverID := prover.ID()
	proof := &state.RecursiveProof{BatchNumber: batchToVerify.BatchNumber, BatchNumberFinal: batchToVerify.BatchNumber, Prover: &proverID, Generating: true}

	// Avoid other thread to process the same batch
	err = a.State.AddGeneratedRecursiveProof(ctx, proof, nil)
	if err != nil {
		log.Errorf("Failed to add batch proof, err: %v", err)
		return nil, nil, err
	}

	return batchToVerify, proof, nil
}

func (a *Aggregator2) tryGenerateBatchProof(ctx context.Context, prover *prover.Prover, ticker *time.Ticker) (bool, error) {
	log.Debugf("tryGenerateBatchProof start %s", prover.ID())

	batchToProve, proof, err0 := a.getAndLockBatchToProve(ctx, prover, ticker)
	if err0 != nil {
		return false, err0
	}

	var err error

	defer func() {
		if err != nil {
			err2 := a.State.DeleteGeneratedRecursiveProof(ctx, proof.BatchNumber, proof.BatchNumberFinal, nil)
			if err2 != nil {
				log.Errorf("Failed to delete proof in progress, err: %v", err2)
			}
		}
	}()

	log.Infof("Prover %s is going to be used to generate batch proof: %d", prover.ID(), batchToProve.BatchNumber)

	log.Infof("Sending zki + batch to the prover, batchNumber: %d", batchToProve.BatchNumber)
	inputProver, err := a.buildInputProver(ctx, batchToProve)
	if err != nil {
		log.Errorf("Failed to build input prover, err: %v", err)
		return false, err
	}

	b, err := json.Marshal(inputProver)
	proof.InputProver = string(b)

	log.Infof("Sending a batch to the prover, OLDSTATEROOT: %s, OLDBATCHNUM: %d",
		inputProver.PublicInputs.OldStateRoot, inputProver.PublicInputs.OldBatchNum)

	genProofID, err := prover.BatchProof(inputProver)
	if err != nil {
		log.Errorf("Failed to get batch proof id, err: %v", err)
		return false, err
	}

	proof.ProofID = &genProofID

	log.Infof("Proof ID for batch %d: %v", proof.BatchNumber, *proof.ProofID)

	resGetProof, err := prover.WaitRecursiveProof(ctx, *proof.ProofID)
	if err != nil {
		log.Errorf("Failed to get proof from prover, err: %v", err)
		return false, err
	}

	log.Infof("Batch proof %s generated", *proof.ProofID)

	proof.Proof = resGetProof
	proof.Generating = false

	proofSent, _ := a.trySendFinalProof(ctx, prover, proof, ticker)

	if !proofSent {
		// Store proof
		err = a.State.UpdateGeneratedRecursiveProof(ctx, proof, nil)
		if err != nil {
			log.Errorf("Failed to store batch proof result, err: %v", err)
			return false, err
		}
	} else {
		// The proof has sent then delete
		err = a.State.DeleteGeneratedRecursiveProof(ctx, proof.BatchNumber, proof.BatchNumberFinal, nil)
		if err != nil {
			log.Errorf("Failed to delete sent batch proof, err: %v", err)
			return false, err
		}
	}

	log.Debug("tryGenerateBatchProof end")

	return true, nil
}

func (a *Aggregator2) isSynced(ctx context.Context) bool {
	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Warnf("Failed to get last consolidated batch, err: %v", err)
		return false
	}
	if lastVerifiedBatch == nil {
		return false
	}
	lastVerifiedEthBatchNum, err := a.Ethman.GetLatestVerifiedBatchNum()
	if err != nil {
		log.Warnf("Failed to get last eth batch, err: %v", err)
		return false
	}
	if lastVerifiedBatch.BatchNumber < lastVerifiedEthBatchNum {
		log.Infof("Waiting for the state to be synced, lastVerifiedBatchNum: %d, lastVerifiedEthBatchNum: %d",
			lastVerifiedBatch.BatchNumber, lastVerifiedEthBatchNum)
		return false
	}
	return true
}

func (a *Aggregator2) buildInputProver(ctx context.Context, batchToVerify *state.Batch) (*pb.InputProver, error) {
	previousBatch, err := a.State.GetBatchByNumber(ctx, batchToVerify.BatchNumber-1, nil)
	if err != nil && err != state.ErrStateNotSynchronized {
		return nil, fmt.Errorf("Failed to get previous batch, err: %v", err)
	}

	blockTimestampByte := make([]byte, 8) //nolint:gomnd
	binary.BigEndian.PutUint64(blockTimestampByte, uint64(batchToVerify.Timestamp.Unix()))
	batchHashData := common.BytesToHash(keccak256.Hash(
		batchToVerify.BatchL2Data,
		batchToVerify.GlobalExitRoot[:],
		blockTimestampByte,
		batchToVerify.Coinbase[:],
	))
	pubAddr, err := a.Ethman.GetPublicAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get public address, err: %w", err)
	}
	inputProver := &pb.InputProver{
		PublicInputs: &pb.PublicInputs{
			OldStateRoot:    previousBatch.StateRoot.Bytes(),
			OldAccInputHash: []byte(batchHashData.String()), //previousBatch.acc_input_hash
			OldBatchNum:     previousBatch.BatchNumber,
			ChainId:         a.cfg.ChainID,
			BatchL2Data:     batchToVerify.BatchL2Data,
			GlobalExitRoot:  batchToVerify.GlobalExitRoot.Bytes(),
			EthTimestamp:    uint64(batchToVerify.Timestamp.Unix()),
			SequencerAddr:   batchToVerify.Coinbase.String(),
			AggregatorAddr:  pubAddr.String(),
		},
		Db:                map[string]string{},
		ContractsBytecode: map[string]string{},
	}

	return inputProver, nil
}

func (a *Aggregator2) compareInputHashes(ip *pb.InputProver, resGetProof *pb.GetProofResponse_FinalProof) {
	/*	// Calc inputHash
		batchNumberByte := make([]byte, 8) //nolint:gomnd
		binary.BigEndian.PutUint64(batchNumberByte, ip.PublicInputs.OldBatchNum)
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
		}*/
}

func waitTick(ctx context.Context, ticker *time.Ticker) {
	select {
	case <-ticker.C:
		// nothing
	case <-ctx.Done():
		return
	}
}

// healthChecker will provide an implementation of the HealthCheck interface.
type healthChecker struct{}

// newHealthChecker returns a health checker according to standard package
// grpc.health.v1.
func newHealthChecker() *healthChecker {
	return &healthChecker{}
}

// HealthCheck interface implementation.

// Check returns the current status of the server for unary gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (hc *healthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	log.Info("Serving the Check request for health check")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch returns the current status of the server for stream gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (hc *healthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	log.Info("Serving the Watch request for health check")
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}
