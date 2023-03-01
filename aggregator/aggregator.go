package aggregator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/metrics"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/prover"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/peer"
)

const (
	mockedStateRoot     = "0x090bcaf734c4f06c93954a827b45a6e8c67b8e0fd1e0a35a1c5982d6961828f9"
	mockedLocalExitRoot = "0x17c04c3760510b48c6012742c540a81aba4bca2f78b9d14bfd2f123e2e53ea3e"

	ethTxManagerOwner = "aggregator"
	monitoredIDFormat = "proof-from-%v-to-%v"
)

// ErrNotValidForFinal is returned for proof not valid to be used as final.
var ErrNotValidForFinal error = errors.New("proof not valid to be sent as final")

type proverJob interface {
	Proof()
}

type nilJob struct {
	tracking string
}

// Proof implements the proverJob interface.
func (nilJob) Proof() {}

type aggregationJob struct {
	tracking string
	proof1   *state.Proof
	proof2   *state.Proof
	proofCh  chan jobResult
}

// Proof implements the proverJob interface.
func (aggregationJob) Proof() {}

type generationJob struct {
	tracking string
	batch    *state.Batch
	proof    *state.Proof
	proofCh  chan jobResult
}

// Proof implements the proverJob interface.
func (generationJob) Proof() {}

type finalJob struct {
	tracking string
	proof    *state.Proof
}

// Proof implements the proverJob interface.
func (finalJob) Proof() {}

type jobResult struct {
	proverName string
	proverID   string
	tracking   string
	job        proverJob
	proof      *state.Proof
	err        error
}

type finalJobResult struct {
	proverName string
	proverID   string
	job        *finalJob
	proof      *pb.FinalProof
	err        error
}

type proverClient struct {
	name     string
	id       string
	addr     string
	tracking string
	ctx      context.Context
	jobChan  chan proverJob
}

// Aggregator represents an aggregator.
type Aggregator struct {
	pb.UnimplementedAggregatorServiceServer

	cfg Config

	State                stateInterface
	EthTxManager         ethTxManager
	Ethman               etherman
	ProfitabilityChecker aggregatorTxProfitabilityChecker

	proversCh          chan proverClient
	finalJobCh         chan *finalJob
	finalProofCh       chan finalJobResult
	verifyProofTimeOut chan struct{}
	verifyProofTimer   *time.Timer
	srv                *grpc.Server
	ctx                context.Context
	exit               context.CancelFunc
}

// New creates a new aggregator.
func New(
	cfg Config,
	stateInterface stateInterface,
	ethTxManager ethTxManager,
	etherman etherman,
) (*Aggregator, error) {
	var profitabilityChecker aggregatorTxProfitabilityChecker
	switch cfg.TxProfitabilityCheckerType {
	case ProfitabilityBase:
		profitabilityChecker = NewTxProfitabilityCheckerBase(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration, cfg.TxProfitabilityMinReward.Int)
	case ProfitabilityAcceptAll:
		profitabilityChecker = NewTxProfitabilityCheckerAcceptAll(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration)
	}

	a := &Aggregator{
		State:                stateInterface,
		EthTxManager:         ethTxManager,
		Ethman:               etherman,
		ProfitabilityChecker: profitabilityChecker,
		cfg:                  cfg,
		proversCh:            make(chan proverClient),
		finalJobCh:           make(chan *finalJob),
		finalProofCh:         make(chan finalJobResult),
	}

	return a, nil
}

// Start starts the aggregator.
func (a *Aggregator) Start(ctx context.Context) error {
	defer a.Stop()
	var cancel context.CancelFunc
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel = context.WithCancel(ctx)
	a.ctx = ctx
	a.exit = cancel

	metrics.Register()

	// process monitored batch verifications before starting
	a.EthTxManager.ProcessPendingMonitoredTxs(ctx, ethTxManagerOwner, func(result ethtxmanager.MonitoredTxResult, dbTx pgx.Tx) {
		a.handleMonitoredTxResult(result)
	}, nil)

	// delete ungenerated recursive proofs
	err := a.State.DeleteUngeneratedProofs(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to initialize proofs cache, %w", err)
	}

	address := fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen, %w", err)
	}

	a.srv = grpc.NewServer()
	pb.RegisterAggregatorServiceServer(a.srv, a)

	healthService := newHealthChecker()
	grpchealth.RegisterHealthServer(a.srv, healthService)

	go func() {
		log.Infof("Server listening on port %d", a.cfg.Port)
		if err := a.srv.Serve(lis); err != nil {
			a.exit()
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	a.verifyProofTimeOut = make(chan struct{})
	a.resetTimer()

	go a.handleFinalProof()
	go a.aggregate()

	// Wait until context is done
	<-ctx.Done()
	return ctx.Err()
}

// Stop stops the Aggregator server.
func (a *Aggregator) Stop() {
	close(a.finalProofCh)
	a.verifyProofTimer.Stop()
	a.srv.Stop()
	a.exit()
}

// Channel implements the bi-directional communication channel between the
// Prover client and the Aggregator server.
func (a *Aggregator) Channel(stream pb.AggregatorService_ChannelServer) error {
	metrics.ConnectedProver()
	defer metrics.DisconnectedProver()

	ctx := stream.Context()
	var proverAddr net.Addr
	p, ok := peer.FromContext(ctx)
	if ok {
		proverAddr = p.Addr
	}
	prover, err := prover.New(stream, proverAddr, a.cfg.ProofStatePollingInterval)
	if err != nil {
		return err
	}

	proverName := prover.Name()
	proverID := prover.ID()

	log := log.WithFields(
		"prover", proverName,
		"proverId", proverID,
		"proverAddr", proverAddr.String(),
	)
	log.Debug("Establishing stream connection with prover")

	jobChan := make(chan proverJob)

	// poll the prover to check when it's idle
	for {
		isIdle, err := prover.IsIdle()
		if err != nil {
			return fmt.Errorf("failed to check prover status, %w", err)
		}
		if !isIdle {
			time.Sleep(a.cfg.ProofStatePollingInterval.Duration)
			continue
		}

		tracking := uuid.NewString()[:8]

		select {
		case <-a.ctx.Done():
			return a.ctx.Err()
		case <-ctx.Done():
			return ctx.Err()
		default:
			//send the readiness message to the aggregator
			log.Debugf("prover ready to receive jobs, tracking [%s]", tracking)
			proverMsg := proverClient{
				name:     proverName,
				id:       proverID,
				addr:     prover.Addr(),
				tracking: tracking,
				ctx:      ctx,
				jobChan:  jobChan,
			}
			a.proversCh <- proverMsg

			// wait for the response in the job channel
			log.Debugf("waiting for job, tracking [%s]", tracking)
		jobsLoop:
			for proverJob := range jobChan {
				var proof *state.Proof
				var proofCh chan jobResult
				var err error

				switch job := proverJob.(type) {
				case *nilJob:
					log := log.WithFields("tracking", job.tracking)
					log.Debug("nothing to prove")

					// nothing to do, wait a bit and retry
					time.Sleep(a.cfg.ProofStatePollingInterval.Duration)
					break jobsLoop

				case *finalJob:
					proof, err := a.handleFinalJob(ctx, prover, job)
					finalJobRes := finalJobResult{
						proverID: proverID,
						job:      job,
						proof:    proof,
						err:      err,
					}

					select {
					case <-a.ctx.Done():
						return a.ctx.Err()
					case a.finalProofCh <- finalJobRes:
						break jobsLoop
					}

				case *aggregationJob:
					proofCh = job.proofCh
					proof, err = a.handleAggregationJob(ctx, prover, job)

				case *generationJob:
					proofCh = job.proofCh
					proof, err = a.handleGenerationJob(ctx, prover, job)
				}

				jr := jobResult{
					proverName: proverName,
					proverID:   proverID,
					tracking:   tracking,
					job:        proverJob,
					proof:      proof,
					err:        err,
				}

				select {
				case <-a.ctx.Done():
					return a.ctx.Err()
				case proofCh <- jr:
				}
			} // jobsLoop
		} // select
	} // main for
}

func (a *Aggregator) handleFinalProof() {
	ctx := a.ctx

	for result := range a.finalProofCh {
		inputProof := result.job.proof
		finalProof := result.proof

		log := log.WithFields(
			"prover", result.proverName,
			"proverId", result.proverID,
			"proofId", *inputProof.ProofID,
			"batches", fmt.Sprintf("%d-%d", inputProof.BatchNumber, inputProof.BatchNumberFinal),
			"tracking", result.job.tracking,
		)

		// mock prover sanity check
		if string(finalProof.Public.NewStateRoot) == mockedStateRoot && string(finalProof.Public.NewLocalExitRoot) == mockedLocalExitRoot {
			// This local exit root and state root come from the mock
			// prover, use the one captured by the executor instead
			finalBatch, err := a.State.GetBatchByNumber(a.ctx, inputProof.BatchNumberFinal, nil)
			if err != nil {
				err := fmt.Errorf("failed to retrieve batch with number [%d], %w", inputProof.BatchNumberFinal, err)
				log.Error(err)
				a.enableFinal()
				continue
			}
			log.Debugf("NewLocalExitRoot and NewStateRoot look like a mock values, using values from executor instead: LER: %v, SR: %v",
				finalBatch.LocalExitRoot.TerminalString(), finalBatch.StateRoot.TerminalString())
			finalProof.Public.NewStateRoot = finalBatch.StateRoot.Bytes()
			finalProof.Public.NewLocalExitRoot = finalBatch.LocalExitRoot.Bytes()
		}

		finalBatch, err := a.State.GetBatchByNumber(ctx, inputProof.BatchNumberFinal, nil)
		if err != nil {
			log.Errorf("Failed to retrieve batch with number [%d]: %v", inputProof.BatchNumberFinal, err)
			continue
		}

		inputs := ethmanTypes.FinalProofInputs{
			FinalProof:       finalProof,
			NewLocalExitRoot: finalBatch.LocalExitRoot.Bytes(),
			NewStateRoot:     finalBatch.StateRoot.Bytes(),
		}

		log.Infof("Final proof inputs: NewLocalExitRoot [%#x], NewStateRoot [%#x]", inputs.NewLocalExitRoot, inputs.NewStateRoot)

		// add batch verification to be monitored
		sender := common.HexToAddress(a.cfg.SenderAddress)
		to, data, err := a.Ethman.BuildTrustedVerifyBatchesTxData(inputProof.BatchNumber-1, inputProof.BatchNumberFinal, &inputs)
		if err != nil {
			log.Errorf("Error estimating batch verification to add to eth tx manager: %v", err)
			a.handleFailureToAddVerifyBatchToBeMonitored(ctx, inputProof)
			continue
		}

		monitoredTxID := fmt.Sprintf(monitoredIDFormat, inputProof.BatchNumber, inputProof.BatchNumberFinal)

		err = a.EthTxManager.Add(ctx, ethTxManagerOwner, monitoredTxID, sender, to, nil, data, nil)
		if err != nil {
			log := log.WithFields("tx", monitoredTxID)
			log.Errorf("Error to add batch verification tx to eth tx manager: %v", err)
			a.handleFailureToAddVerifyBatchToBeMonitored(ctx, inputProof)
			continue
		}

		// process monitored batch verifications before starting a next cycle
		a.EthTxManager.ProcessPendingMonitoredTxs(ctx, ethTxManagerOwner, func(result ethtxmanager.MonitoredTxResult, dbTx pgx.Tx) {
			a.handleMonitoredTxResult(result)
		}, nil)

		a.resetTimer()
	}
}

func (a *Aggregator) handleFailureToAddVerifyBatchToBeMonitored(ctx context.Context, proof *state.Proof) {
	log := log.WithFields("proofId", proof.ProofID, "batches", fmt.Sprintf("%d-%d", proof.BatchNumber, proof.BatchNumberFinal))
	proof.GeneratingSince = nil
	err := a.State.UpdateGeneratedProof(ctx, proof, nil)
	if err != nil {
		log.Errorf("failed updating proof state (false), err: %v", err)
	}
	// a.endProofVerification()
}

func (a *Aggregator) handleMonitoredTxResult(result ethtxmanager.MonitoredTxResult) {
	resLog := log.WithFields("owner", ethTxManagerOwner, "txId", result.ID)
	if result.Status == ethtxmanager.MonitoredTxStatusFailed {
		resLog.Fatal("failed to send batch verification, TODO: review this fatal and define what to do in this case")
	}

	idSlice := strings.Split(result.ID, "-")
	proofBatchNumberStr := idSlice[2]
	proofBatchNumber, err := strconv.ParseUint(proofBatchNumberStr, encoding.Base10, 0)
	if err != nil {
		resLog.Errorf("failed to read final proof batch number from monitored tx: %v", err)
	}

	proofBatchNumberFinalStr := idSlice[4]
	proofBatchNumberFinal, err := strconv.ParseUint(proofBatchNumberFinalStr, encoding.Base10, 0)
	if err != nil {
		resLog.Errorf("failed to read final proof batch number final from monitored tx: %v", err)
	}

	log := log.WithFields("txId", result.ID, "batches", fmt.Sprintf("%d-%d", proofBatchNumber, proofBatchNumberFinal))
	log.Info("Final proof verified")

	// wait for the synchronizer to catch up the verified batches
	log.Debug("A final proof has been sent, waiting for the network to be synced")
	for !a.isSynced(a.ctx, &proofBatchNumberFinal) {
		log.Info("Waiting for synchronizer to sync...")
		time.Sleep(a.cfg.RetryTime.Duration)
	}

	// network is synced with the final proof, we can safely delete all recursive
	// proofs up to the last synced batch
	err = a.State.CleanupGeneratedProofs(a.ctx, proofBatchNumberFinal, nil)
	if err != nil {
		log.Errorf("failed to store proof aggregation result: %v", err)
	}
}

// aggregate is the Aggregator main loop. Here it receives messages from idling
// Provers and it feeds them with jobs. Once a Prover gets a job, the Aggregator
// waits for the result and processes the proof.
func (a *Aggregator) aggregate() {
	for {
		select {
		case <-a.ctx.Done():
			return

		// here we receive idle provers
		case prover := <-a.proversCh:
			log := log.WithFields(
				"proverId", prover.id,
				"proverAddr", prover.addr,
				"tracking", prover.tracking,
			)

			// dedicated channel to receive the proof
			proofCh := make(chan jobResult)

			err := a.feedProver(prover, proofCh)
			if err != nil {
				log.Error(err)
			}

			// wait for the proof
			go func() {
				for {
					select {
					case <-a.ctx.Done():
						return
					case <-prover.ctx.Done():
						return
					case result := <-proofCh:
						log := log.WithFields("batches", fmt.Sprintf("%d-%d", result.proof.BatchNumber, result.proof.BatchNumberFinal))

						if err := a.handleProof(a.ctx, result); err != nil {
							log.Error(err)
						}
						return
					}
				}
			}()
		} // select
	} // for
}

// feedProver prepares the next job to be scheduled to a Prover. If it's time
// to send the final proof, it checks if the eligible proof is in memory or if
// not it retrieves it from the state.
func (a *Aggregator) feedProver(prover proverClient, proofCh chan jobResult) error {
	log := log.WithFields(
		"prover", prover.name,
		"proverId", prover.id,
		"proverAddr", prover.addr,
	)
	ctx := prover.ctx

	sendJob := func(pJob proverJob) error {
		select {
		case <-a.ctx.Done():
			return a.ctx.Err()
		case <-ctx.Done():
			return ctx.Err()
		case prover.jobChan <- pJob:
		}
		return nil
	}

	select {
	case <-a.verifyProofTimeOut:
		log.Debug("Time to send the final proof")

		select {
		// before looking for a proof into the state, we listen if the
		// eligible proof has just been produced by a prover
		case fj := <-a.finalJobCh:
			log.Debugf("received proof valid for final, tracking [%s] ", fj.tracking)
			a.reserveFinal()
			return sendJob(fj)

		default:
			log.Debug("Check if there is a previous batch eligible to be final")
			proof, err := a.eligibleFinalProof(ctx, nil)
			if errors.Is(err, state.ErrNotFound) {
				// nothing to verify, swallow the error and try to feed the
				// prover to make a regular (non-final) proof
				log.Debug("No proofs ready to verify")
			} else if err != nil {
				return err
			} else {
				a.reserveFinal()
				fj := &finalJob{
					tracking: prover.tracking,
					proof:    proof,
				}
				return sendJob(fj)
			}
		}
	default:
	}

	log = log.WithFields("tracking", prover.tracking)

	proof1, proof2, err := a.getAndLockProofsToAggregate(ctx)
	if errors.Is(err, state.ErrNotFound) {
		log.Debug("No proofs to aggregate, trying to generate from batch")

		batch, proof, err := a.getAndLockBatchToProve(ctx, prover.id)
		if errors.Is(err, state.ErrNotFound) {
			log.Debug("No batches to generate proof from")
			// nothing to generate, swallow the error and send a nil job
			return sendJob(&nilJob{tracking: prover.tracking})
		}
		if err != nil {
			return fmt.Errorf("failed to get batch to prove, %w", err)
		}

		log.Debugf("Sending job for proof generation from batch [%d]", batch.BatchNumber)
		pJob := &generationJob{
			tracking: prover.tracking,
			batch:    batch,
			proof:    proof,
			proofCh:  proofCh,
		}
		return sendJob(pJob)
	}
	if err != nil {
		return fmt.Errorf("failed to get proofs to aggregate, %w", err)
	}

	log.Debugf("Sending job for aggregating proofs of batches [%d-%d]", proof1.BatchNumber, proof2.BatchNumberFinal)
	pJob := &aggregationJob{
		tracking: prover.tracking,
		proof1:   proof1,
		proof2:   proof2,
		proofCh:  proofCh,
	}
	return sendJob(pJob)
}

// handleProof takes care of storing the generated proof into the state. If
// it's time to send the final proof and the proof in hand is the eligible one,
// then it sends it over a channel to be verified.
func (a *Aggregator) handleProof(ctx context.Context, result jobResult) error {
	log := log.WithFields(
		"prover", result.proverName,
		"proverId", result.proverID,
		"proofId", *result.proof.ProofID,
		"tracking", result.tracking,
	)

	if result.err != nil {
		return a.rollbackFailedJob(ctx, result)
	}

	dbTx, err := a.State.BeginStateTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction to store proof aggregation result, %w", err)
	}

	var validForFinal *time.Time

	select {
	case <-a.verifyProofTimeOut:
		log.Debug("Time to send the final proof, checking if the current proof can be sent as final")

		_, err := a.eligibleFinalProof(a.ctx, result.proof)
		if errors.Is(err, ErrNotValidForFinal) {
			// proof is not valid for final, carry on storing it
			log.Debug(err.Error())
		} else if err != nil {
			return fmt.Errorf("failed to validate job for final proof, %w", err)
		}

		// if the proof is eligible to be final, it needs to be reserved
		// by setting the GeneratingSince timestamp
		now := time.Now().UTC().Round(time.Microsecond)
		validForFinal = &now
	default:
	}

	switch job := result.job.(type) {
	case *aggregationJob:
		err = a.State.DeleteGeneratedProofs(ctx, job.proof1.BatchNumber, job.proof2.BatchNumberFinal, dbTx)
		if err != nil {
			if err := dbTx.Rollback(ctx); err != nil {
				return fmt.Errorf("failed to rollback failing to delete aggregation input proofs, %w", err)
			}
			return fmt.Errorf("failed to delete aggregation input proofs, %w", err)
		}

		result.proof.GeneratingSince = validForFinal
		err := a.State.AddGeneratedProof(ctx, result.proof, dbTx)
		if err != nil {
			if err := dbTx.Rollback(ctx); err != nil {
				return fmt.Errorf("failed to rollback failing to store proof aggregation result, %w", err)
			}
			return fmt.Errorf("failed to store proof aggregation result, %w", err)
		}

	case *generationJob:
		// Store the proof, if it's a proof valid for final, keep it reserved
		result.proof.GeneratingSince = validForFinal
		err := a.State.UpdateGeneratedProof(ctx, result.proof, dbTx)
		if err != nil {
			if err := dbTx.Rollback(ctx); err != nil {
				return fmt.Errorf("failed to rollback failing updating proof in progress, %w", err)
			}
			return fmt.Errorf("failed to to store batch proof result, %w", err)
		}
	}

	if err := dbTx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit proof job, %w", err)
	}

	// send the eligible final proof
	if validForFinal != nil {
		fj := &finalJob{
			tracking: result.tracking,
			proof:    result.proof,
		}

		select {
		case <-a.ctx.Done():
			return a.ctx.Err()
		case a.finalJobCh <- fj:
		}
	}
	return nil
}

func (a *Aggregator) rollbackFailedJob(ctx context.Context, result jobResult) error {
	dbTx, err := a.State.BeginStateTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction to rollback failed job, %w", err)
	}

	switch job := result.job.(type) {
	case *aggregationJob:
		err := a.unlockProofsToAggregate(ctx, job.proof1, job.proof2, dbTx)
		if err != nil {
			if err := dbTx.Rollback(ctx); err != nil {
				return fmt.Errorf("failed to rollback failing unlockProofsToAggregate, %w", err)
			}
			return fmt.Errorf("failed to unlock aggregated proofs, %w", err)
		}
		if err := dbTx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit unlocking aggregated proofs, %w", err)
		}

		return fmt.Errorf("failed to aggregate proofs, %w", result.err)

	case *generationJob:
		log.Errorf("Failed to generate proof: %v", result.err)

		err := a.State.DeleteGeneratedProofs(ctx, job.proof.BatchNumber, job.proof.BatchNumberFinal, dbTx)
		if err != nil {
			if err := dbTx.Rollback(ctx); err != nil {
				return fmt.Errorf("failed to rollback failing deleting proof in progress, %w", err)
			}
			return fmt.Errorf("failed to delete proof in progress, %w", err)
		}
		if err := dbTx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit rollback on proof generation job, %w", err)
		}

		return fmt.Errorf("failed to generate proof, %w", result.err)
	}
	return nil
}

// eligibleFinalProof returns a proof which is suitable to be used to generate
// the final proof.  If the `proof` argument is not nil, then the provided
// proof is checked for eligiblity, otherwise if `proof` is nil, a valid proof
// is retrieved from the state.
func (a *Aggregator) eligibleFinalProof(ctx context.Context, proof *state.Proof) (*state.Proof, error) {
	log.Debug("Checking if network is synced")
	for !a.isSynced(ctx, nil) { // TODO(pg): check if nil batch is ok here
		log.Debug("Waiting for synchronizer to sync...")
		time.Sleep(a.cfg.RetryTime.Duration)
	}
	log.Debug("Network synced")

	var lastBatchNumber uint64
	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		return nil, fmt.Errorf("failed to get last verified batch, %w", err)
	}
	if lastVerifiedBatch != nil {
		lastBatchNumber = lastVerifiedBatch.BatchNumber
	}

	if proof == nil {
		proof, err = a.getAndLockProofReadyToVerify(ctx, lastBatchNumber)
		if err != nil {
			return nil, err
		}
	} else {
		batchNumberToVerify := lastBatchNumber + 1

		if proof.BatchNumber != batchNumberToVerify {
			batchNumberStr := fmt.Sprintf("%d", proof.BatchNumber)
			if proof.BatchNumber != proof.BatchNumberFinal {
				batchNumberStr = fmt.Sprintf("%s-%d", batchNumberStr, proof.BatchNumberFinal)
			}
			return nil, fmt.Errorf("%w: batch number [%s] is not the following to last verfied batch number [%d]",
				ErrNotValidForFinal, batchNumberStr, lastBatchNumber)
		}

		completeSeq, err := a.State.CheckProofContainsCompleteSequences(ctx, proof, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to check if proof [%d-%d] contains complete sequences, %w", proof.BatchNumber, proof.BatchNumberFinal, err)
		}
		if !completeSeq {
			return nil, fmt.Errorf("%w: proof [%d-%d] does not contain complete sequences", ErrNotValidForFinal, proof.BatchNumber, proof.BatchNumberFinal)
		}
	}

	log.Debugf("Proof ID [%s] for batches [%d-%d] is valid for final", *proof.ProofID, proof.BatchNumber, proof.BatchNumberFinal)

	return proof, nil
}

func (a *Aggregator) handleAggregationJob(ctx context.Context, prover proverInterface, job *aggregationJob) (*state.Proof, error) {
	proverName := prover.Name()
	proverID := prover.ID()

	log := log.WithFields(
		"prover", proverName,
		"proverId", proverID,
		"proverAddr", prover.Addr(),
		"tracking", job.tracking,
	)

	log.Infof("Aggregating proofs [%d-%d] and [%d-%d]",
		job.proof1.BatchNumber, job.proof1.BatchNumberFinal, job.proof2.BatchNumber, job.proof2.BatchNumberFinal)

	log = log.WithFields("batches", fmt.Sprintf("%d-%d", job.proof1.BatchNumber, job.proof2.BatchNumberFinal))

	now := time.Now().UTC().Round(time.Microsecond)
	proof := &state.Proof{
		BatchNumber:      job.proof1.BatchNumber,
		BatchNumberFinal: job.proof2.BatchNumberFinal,
		Prover:           &proverID,
		InputProver:      job.proof1.InputProver,
		GeneratingSince:  &now,
	}

	proofID, err := prover.AggregatedProof(job.proof1.Proof, job.proof2.Proof)
	if err != nil {
		return nil, fmt.Errorf("failed to instruct prover to generate aggregated proof, %w", err)
	}
	proof.ProofID = proofID

	log.Infof("Proof ID for aggregated proof [%d-%d]: %v", proof.BatchNumber, proof.BatchNumberFinal, *proof.ProofID)
	log = log.WithFields("proofId", *proofID)

	aggrProof, err := prover.WaitRecursiveProof(ctx, *proofID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve aggregated proof from prover, %w", err)
	}
	proof.Proof = aggrProof

	return proof, nil
}

func (a *Aggregator) handleGenerationJob(ctx context.Context, prover proverInterface, job *generationJob) (*state.Proof, error) {
	proverName := prover.Name()
	proverID := prover.ID()

	log := log.WithFields(
		"proverName", proverName,
		"proverId", proverID,
		"proverAddr", prover.Addr(),
		"batch", job.batch.BatchNumber,
		"tracking", job.tracking,
	)

	log.Info("Generating proof")

	proof := job.proof
	proof.Prover = &proverID

	log.Info("Sending zki + batch to the prover")

	inputProver, err := a.buildInputProver(ctx, job.batch)
	if err != nil {
		return nil, fmt.Errorf("failed to build input prover, %w", err)
	}

	b, err := json.Marshal(inputProver)
	if err != nil {
		return nil, fmt.Errorf("failed serialize input prover, %w", err)
	}
	proof.InputProver = string(b)

	log.Infof("Sending a batch to the prover, OLDSTATEROOT: %#x, OLDBATCHNUM: %d",
		inputProver.PublicInputs.OldStateRoot, inputProver.PublicInputs.OldBatchNum)

	genProofID, err := prover.BatchProof(inputProver)
	if err != nil {
		return nil, fmt.Errorf("failed instruct prover to prove a batch, %w", err)
	}
	proof.ProofID = genProofID

	log.Infof("Proof ID [%s]", *job.proof.ProofID)
	log = log.WithFields("proofId", *genProofID)

	genProof, err := prover.WaitRecursiveProof(ctx, *job.proof.ProofID)
	if err != nil {
		return nil, fmt.Errorf("failed to get proof from prover, %w", err)
	}
	proof.Proof = genProof

	log.Info("Proof generated")

	return proof, nil
}

func (a *Aggregator) handleFinalJob(ctx context.Context, prover proverInterface, job *finalJob) (*pb.FinalProof, error) {
	proverName := prover.Name()
	proverID := prover.ID()
	log := log.WithFields(
		"prover", proverName,
		"proverId", proverID,
		"proverAddr", prover.Addr(),
		"recursiveProofId", *job.proof.ProofID,
		"batches", fmt.Sprintf("%d-%d", job.proof.BatchNumber, job.proof.BatchNumberFinal),
		"tracking", job.tracking,
	)

	log.Info("Generating final proof")

	finalProofID, err := prover.FinalProof(job.proof.Proof, a.cfg.SenderAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to instruct prover to prepare final proof, %w", err)
	}
	log.Infof("Final proof ID [%s]", *finalProofID)
	log = log.WithFields("finalProofId", *finalProofID)

	proof, err := prover.WaitFinalProof(ctx, *finalProofID)
	if err != nil {
		return nil, fmt.Errorf("failed to get final proof, %w", err)
	}
	log.Infof("Final proof generated")

	return proof, nil
}

func (a *Aggregator) getAndLockProofReadyToVerify(ctx context.Context, lastVerifiedBatchNum uint64) (*state.Proof, error) {
	// Get proof ready to be verified
	proofToVerify, err := a.State.GetProofReadyToVerify(ctx, lastVerifiedBatchNum, nil)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC().Round(time.Microsecond)
	proofToVerify.GeneratingSince = &now

	err = a.State.UpdateGeneratedProof(ctx, proofToVerify, nil)
	if err != nil {
		return nil, err
	}

	return proofToVerify, nil
}

func (a *Aggregator) unlockProofsToAggregate(ctx context.Context, proof1, proof2 *state.Proof, dbTx pgx.Tx) error {
	proof1.GeneratingSince = nil
	err := a.State.UpdateGeneratedProof(ctx, proof1, dbTx)
	if err != nil {
		return err
	}

	proof2.GeneratingSince = nil
	err = a.State.UpdateGeneratedProof(ctx, proof2, dbTx)
	if err != nil {
		return err
	}

	return nil
}

func (a *Aggregator) getAndLockProofsToAggregate(ctx context.Context) (*state.Proof, *state.Proof, error) {
	proof1, proof2, err := a.State.GetProofsToAggregate(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	// Set proofs in aggregating state in a single transaction
	// TODO(pg) create a state.UpdateGeneratedProofs method
	dbTx, err := a.State.BeginStateTransaction(ctx)
	if err != nil {
		return nil, nil, err
	}

	now := time.Now().UTC().Round(time.Microsecond)
	proof1.GeneratingSince = &now
	err = a.State.UpdateGeneratedProof(ctx, proof1, dbTx)
	if err != nil {
		log.Errorf("Failed to set proof aggregation state, err: %v", err)
		dbTx.Rollback(ctx) //nolint:errcheck
		return nil, nil, err
	}

	proof2.GeneratingSince = &now
	err = a.State.UpdateGeneratedProof(ctx, proof2, dbTx)
	if err != nil {
		log.Errorf("Failed to set proof aggregation state, err: %v", err)
		dbTx.Rollback(ctx) //nolint:errcheck
		return nil, nil, err
	}

	if err := dbTx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to lock proofs to aggregate, %w", err)
	}

	return proof1, proof2, nil
}

func (a *Aggregator) getAndLockBatchToProve(ctx context.Context, proverID string) (*state.Batch, *state.Proof, error) {
	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	// Get virtual batch pending to generate proof
	batchToVerify, err := a.State.GetVirtualBatchToProve(ctx, lastVerifiedBatch.BatchNumber, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Infof("Found virtual batch [%d] pending to generate proof", batchToVerify.BatchNumber)

	log.Infof("Checking profitability to aggregate batch, batchNumber: %d", batchToVerify.BatchNumber)

	// pass matic collateral as zero here, bcs in smart contract fee for aggregator is not defined yet
	isProfitable, err := a.ProfitabilityChecker.IsProfitable(ctx, big.NewInt(0))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check aggregator profitability, %w", err)
	}

	if !isProfitable {
		return nil, nil, fmt.Errorf("batch %d is not profitable, matic collateral %d, %w", batchToVerify.BatchNumber, big.NewInt(0), err)
	}

	now := time.Now().UTC().Round(time.Microsecond)
	proof := &state.Proof{
		BatchNumber:      batchToVerify.BatchNumber,
		BatchNumberFinal: batchToVerify.BatchNumber,
		Prover:           &proverID,
		GeneratingSince:  &now,
	}

	// Avoid other provers to process the same batch
	err = a.State.AddGeneratedProof(ctx, proof, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to add batch proof, %w", err)
	}

	return batchToVerify, proof, nil
}

func (a *Aggregator) resetTimer() {
	a.verifyProofTimer = time.AfterFunc(a.cfg.VerifyProofInterval.Duration, func() {
		a.enableFinal()
	})
}

func (a *Aggregator) enableFinal() { close(a.verifyProofTimeOut) }

func (a *Aggregator) reserveFinal() { a.verifyProofTimeOut = make(chan struct{}) }

// isSynced checks if the state is synchronized with L1. If a batch number is
// provided, it makes sure that the state is synced with that batch.
func (a *Aggregator) isSynced(ctx context.Context, batchNum *uint64) bool {
	// get latest verified batch as seen by the synchronizer
	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err == state.ErrNotFound {
		return false
	}
	if err != nil {
		log.Warnf("Failed to get last consolidated batch: %v", err)
		return false
	}

	if lastVerifiedBatch == nil {
		return false
	}

	if batchNum != nil && lastVerifiedBatch.BatchNumber < *batchNum {
		log.Infof("Waiting for the state to be synced, lastVerifiedBatchNum: %d, waiting for batch: %d", lastVerifiedBatch.BatchNumber, batchNum)
		return false
	}

	// latest verified batch in L1
	lastVerifiedEthBatchNum, err := a.Ethman.GetLatestVerifiedBatchNum()
	if err != nil {
		log.Warnf("Failed to get last eth batch, err: %v", err)
		return false
	}

	// check if L2 is synced with L1
	if lastVerifiedBatch.BatchNumber < lastVerifiedEthBatchNum {
		log.Infof("Waiting for the state to be synced, lastVerifiedBatchNum: %d, lastVerifiedEthBatchNum: %d, waiting for batch",
			lastVerifiedBatch.BatchNumber, lastVerifiedEthBatchNum)
		return false
	}
	return true
}

func (a *Aggregator) buildInputProver(ctx context.Context, batchToVerify *state.Batch) (*pb.InputProver, error) {
	previousBatch, err := a.State.GetBatchByNumber(ctx, batchToVerify.BatchNumber-1, nil)
	if err != nil && err != state.ErrStateNotSynchronized {
		return nil, fmt.Errorf("failed to get previous batch, %w", err)
	}

	inputProver := &pb.InputProver{
		PublicInputs: &pb.PublicInputs{
			OldStateRoot:    previousBatch.StateRoot.Bytes(),
			OldAccInputHash: previousBatch.AccInputHash.Bytes(),
			OldBatchNum:     previousBatch.BatchNumber,
			ChainId:         a.cfg.ChainID,
			BatchL2Data:     batchToVerify.BatchL2Data,
			GlobalExitRoot:  batchToVerify.GlobalExitRoot.Bytes(),
			EthTimestamp:    uint64(batchToVerify.Timestamp.Unix()),
			SequencerAddr:   batchToVerify.Coinbase.String(),
			AggregatorAddr:  a.cfg.SenderAddress,
		},
		Db:                map[string]string{},
		ContractsBytecode: map[string]string{},
	}

	return inputProver, nil
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
func (hc *healthChecker) Check(ctx context.Context, req *grpchealth.HealthCheckRequest) (*grpchealth.HealthCheckResponse, error) {
	log.Info("Serving the Check request for health check")
	return &grpchealth.HealthCheckResponse{
		Status: grpchealth.HealthCheckResponse_SERVING,
	}, nil
}

// Watch returns the current status of the server for stream gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (hc *healthChecker) Watch(req *grpchealth.HealthCheckRequest, server grpchealth.Health_WatchServer) error {
	log.Info("Serving the Watch request for health check")
	return server.Send(&grpchealth.HealthCheckResponse{
		Status: grpchealth.HealthCheckResponse_SERVING,
	})
}
