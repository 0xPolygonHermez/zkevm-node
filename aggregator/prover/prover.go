package prover

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
	"unicode"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/metrics"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

var (
	ErrBadProverResponse    = errors.New("prover returned wrong type for response")  //nolint:revive
	ErrProverInternalError  = errors.New("prover returned INTERNAL_ERROR response")  //nolint:revive
	ErrProverCompletedError = errors.New("prover returned COMPLETED_ERROR response") //nolint:revive
	ErrBadRequest           = errors.New("prover returned ERROR for a bad request")  //nolint:revive
	ErrUnspecified          = errors.New("prover returned an UNSPECIFIED response")  //nolint:revive
	ErrUnknown              = errors.New("prover returned an unknown response")      //nolint:revive
	ErrProofCanceled        = errors.New("proof has been canceled")                  //nolint:revive
	ErrUnsupportedForkID    = errors.New("prover does not support required fork ID") //nolint:revive
)

type ProverJob interface {
	Job()
}

type JobResult struct {
	ProverName string
	ProverID   string
	Tracking   string
	Job        ProverJob
	Proof      *state.Proof
	Err        error
}

type FinalJobResult struct {
	ProverName string
	ProverID   string
	Job        *FinalJob
	Proof      *pb.FinalProof
	Err        error
}

type NilJob struct {
	Tracking string
}

// Job implements the proverJob interface.
func (*NilJob) Job() {}

type FinalJob struct {
	Tracking      string
	SenderAddress string
	Proof         *state.Proof
}

// Job implements the proverJob interface.
func (*FinalJob) Job() {}

type AggregationJob struct {
	Tracking string
	Proof1   *state.Proof
	Proof2   *state.Proof
	ProofCh  chan *JobResult
}

// Job implements the proverJob interface.
func (*AggregationJob) Job() {}

type GenerationJob struct {
	Tracking    string
	Batch       *state.Batch
	InputProver *pb.InputProver
	Proof       *state.Proof
	ProofCh     chan *JobResult
}

// Job implements the proverJob interface.
func (*GenerationJob) Job() {}

// Prover abstraction of the grpc prover client.
type Prover struct {
	name                      string
	id                        string
	address                   net.Addr
	proofStatePollingInterval types.Duration
	stream                    pb.AggregatorService_ChannelServer
}

// New returns a new Prover instance.
func New(
	stream pb.AggregatorService_ChannelServer,
	addr net.Addr,
	proofStatePollingInterval types.Duration,
	forkID uint64,
) (*Prover, error) {
	p := &Prover{
		stream:                    stream,
		address:                   addr,
		proofStatePollingInterval: proofStatePollingInterval,
	}
	status, err := p.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve prover id, %w", err)
	}
	if status.ForkId != forkID {
		log.Debugf("Prover %s supports fork ID %d", p.ID(), status.ForkId)
		return nil, ErrUnsupportedForkID
	}
	p.name = status.ProverName
	p.id = status.ProverId
	return p, nil
}

// Name returns the Prover name.
func (p *Prover) Name() string { return p.name }

// ID returns the Prover ID.
func (p *Prover) ID() string { return p.id }

// Addr returns the prover IP address.
func (p *Prover) Addr() string {
	if p.address == nil {
		return ""
	}
	return p.address.String()
}

// Status gets the prover status.
func (p *Prover) Status() (*pb.GetStatusResponse, error) {
	req := &pb.AggregatorMessage{
		Request: &pb.AggregatorMessage_GetStatusRequest{
			GetStatusRequest: &pb.GetStatusRequest{},
		},
	}
	res, err := p.call(req)
	if err != nil {
		return nil, err
	}
	if msg, ok := res.Response.(*pb.ProverMessage_GetStatusResponse); ok {
		return msg.GetStatusResponse, nil
	}
	return nil, fmt.Errorf("%w, wanted %T, got %T", ErrBadProverResponse, &pb.ProverMessage_GetStatusResponse{}, res.Response)
}

// IsIdle returns true if the prover is idling.
func (p *Prover) IsIdle() (bool, error) {
	status, err := p.Status()
	if err != nil {
		return false, err
	}
	return status.Status == pb.GetStatusResponse_STATUS_IDLE, nil
}

// HandleAggregationJob takes care of producing a recursive proof aggregating 2
// proofs received in the job.  It returns the result of the job execution
// containing the proof in case of success. In case of error the JobResult Err
// field will be non-nil.
func (p *Prover) HandleAggregationJob(ctx context.Context, job *AggregationJob) *JobResult {
	proverName := p.Name()
	proverID := p.ID()

	log := log.WithFields(
		"prover", proverName,
		"proverId", proverID,
		"proverAddr", p.Addr(),
		"tracking", job.Tracking,
	)

	log.Infof("Aggregating proofs [%d-%d] and [%d-%d]",
		job.Proof1.BatchNumber, job.Proof1.BatchNumberFinal, job.Proof2.BatchNumber, job.Proof2.BatchNumberFinal)

	log = log.WithFields("batches", fmt.Sprintf("%d-%d", job.Proof1.BatchNumber, job.Proof2.BatchNumberFinal))

	jr := JobResult{
		ProverName: proverName,
		ProverID:   proverID,
		Tracking:   job.Tracking,
		Job:        job,
	}

	now := time.Now().UTC().Round(time.Microsecond)
	proof := &state.Proof{
		BatchNumber:      job.Proof1.BatchNumber,
		BatchNumberFinal: job.Proof2.BatchNumberFinal,
		Prover:           &proverID,
		InputProver:      job.Proof1.InputProver,
		GeneratingSince:  &now,
	}

	proofID, err := p.AggregatedProof(job.Proof1.Proof, job.Proof2.Proof)
	if err != nil {
		err = fmt.Errorf("failed to instruct prover to generate aggregated proof, %w", err)
		log.Error(FirstToUpper(err.Error()))
		jr.Err = err
		return &jr
	}
	proof.ProofID = proofID

	log.Infof("Proof ID for aggregated proof [%d-%d]: %v", proof.BatchNumber, proof.BatchNumberFinal, *proof.ProofID)
	log = log.WithFields("proofId", *proofID)

	aggrProof, err := p.WaitRecursiveProof(ctx, *proofID)
	if err != nil {
		err = fmt.Errorf("failed to retrieve aggregated proof from prover, %w", err)
		log.Error(FirstToUpper(err.Error()))
		jr.Err = err
		return &jr
	}
	proof.Proof = aggrProof
	jr.Proof = proof

	return &jr
}

// HandleGenerationJob takes care of producing a proof generating it from a
// batch received in the job.  It returns the result of the job execution
// containing the proof in case of success. In case of error the JobResult Err
// field will be non-nil.
func (p *Prover) HandleGenerationJob(ctx context.Context, job *GenerationJob) *JobResult {
	proverName := p.Name()
	proverID := p.ID()

	log := log.WithFields(
		"proverName", proverName,
		"proverId", proverID,
		"proverAddr", p.Addr(),
		"batch", job.Batch.BatchNumber,
		"tracking", job.Tracking,
	)

	log.Info("Generating proof")

	jr := JobResult{
		ProverName: proverName,
		ProverID:   proverID,
		Tracking:   job.Tracking,
		Job:        job,
	}

	proof := job.Proof
	proof.Prover = &proverID

	log.Info("Sending zki + batch to the prover")

	b, err := json.Marshal(job.InputProver)
	if err != nil {
		err = fmt.Errorf("failed serialize input prover, %w", err)
	}
	proof.InputProver = string(b)

	log.Infof("Sending a batch to the prover, OLDSTATEROOT: %#x, OLDBATCHNUM: %d",
		job.InputProver.PublicInputs.OldStateRoot, job.InputProver.PublicInputs.OldBatchNum)

	genProofID, err := p.BatchProof(job.InputProver)
	if err != nil {
		err = fmt.Errorf("failed instruct prover to prove a batch, %w", err)
		log.Error(FirstToUpper(err.Error()))
		jr.Err = err
		return &jr
	}
	proof.ProofID = genProofID

	log.Infof("Proof ID [%s]", *job.Proof.ProofID)
	log = log.WithFields("proofId", *genProofID)

	genProof, err := p.WaitRecursiveProof(ctx, *job.Proof.ProofID)
	if err != nil {
		err = fmt.Errorf("failed to get proof from prover, %w", err)
		log.Error(FirstToUpper(err.Error()))
		jr.Err = err
		return &jr
	}
	proof.Proof = genProof
	jr.Proof = proof

	log.Info("Proof generated")

	return &jr
}

func (p *Prover) HandleFinalJob(ctx context.Context, job *FinalJob) *FinalJobResult {
	proverName := p.Name()
	proverID := p.ID()
	log := log.WithFields(
		"prover", proverName,
		"proverId", proverID,
		"proverAddr", p.Addr(),
		"recursiveProofId", *job.Proof.ProofID,
		"batches", fmt.Sprintf("%d-%d", job.Proof.BatchNumber, job.Proof.BatchNumberFinal),
		"tracking", job.Tracking,
	)

	log.Info("Generating final proof")

	finalJobRes := FinalJobResult{
		ProverID: proverID,
		Job:      job,
	}

	finalProofID, err := p.FinalProof(job.Proof.Proof, job.SenderAddress)
	if err != nil {
		err = fmt.Errorf("failed to instruct prover to prepare final proof, %w", err)
		log.Error(FirstToUpper(err.Error()))
		finalJobRes.Err = err
		return &finalJobRes
	}
	log.Infof("Final proof ID [%s]", *finalProofID)
	log = log.WithFields("finalProofId", *finalProofID)

	proof, err := p.WaitFinalProof(ctx, *finalProofID)
	if err != nil {
		err = fmt.Errorf("failed to get final proof, %w", err)
		log.Error(FirstToUpper(err.Error()))
		finalJobRes.Err = err
		return &finalJobRes
	}
	finalJobRes.Proof = proof

	log.Infof("Final proof generated")

	return &finalJobRes
}

// BatchProof instructs the prover to generate a batch proof for the provided
// input. It returns the ID of the proof being computed.
func (p *Prover) BatchProof(input *pb.InputProver) (*string, error) {
	metrics.WorkingProver()

	req := &pb.AggregatorMessage{
		Request: &pb.AggregatorMessage_GenBatchProofRequest{
			GenBatchProofRequest: &pb.GenBatchProofRequest{Input: input},
		},
	}
	res, err := p.call(req)
	if err != nil {
		return nil, err
	}

	if msg, ok := res.Response.(*pb.ProverMessage_GenBatchProofResponse); ok {
		switch msg.GenBatchProofResponse.Result {
		case pb.Result_RESULT_UNSPECIFIED:
			return nil, fmt.Errorf("failed to generate proof %s, %w, input %v", msg.GenBatchProofResponse.String(), ErrUnspecified, input)
		case pb.Result_RESULT_OK:
			return &msg.GenBatchProofResponse.Id, nil
		case pb.Result_RESULT_ERROR:
			return nil, fmt.Errorf("failed to generate proof %s, %w, input %v", msg.GenBatchProofResponse.String(), ErrBadRequest, input)
		case pb.Result_RESULT_INTERNAL_ERROR:
			return nil, fmt.Errorf("failed to generate proof %s, %w, input %v", msg.GenBatchProofResponse.String(), ErrProverInternalError, input)
		default:
			return nil, fmt.Errorf("failed to generate proof %s, %w,input %v", msg.GenBatchProofResponse.String(), ErrUnknown, input)
		}
	}

	return nil, fmt.Errorf("%w, wanted %T, got %T", ErrBadProverResponse, &pb.ProverMessage_GenBatchProofResponse{}, res.Response)
}

// AggregatedProof instructs the prover to generate an aggregated proof from
// the two inputs provided. It returns the ID of the proof being computed.
func (p *Prover) AggregatedProof(inputProof1, inputProof2 string) (*string, error) {
	metrics.WorkingProver()

	req := &pb.AggregatorMessage{
		Request: &pb.AggregatorMessage_GenAggregatedProofRequest{
			GenAggregatedProofRequest: &pb.GenAggregatedProofRequest{
				RecursiveProof_1: inputProof1,
				RecursiveProof_2: inputProof2,
			},
		},
	}
	res, err := p.call(req)
	if err != nil {
		return nil, err
	}

	if msg, ok := res.Response.(*pb.ProverMessage_GenAggregatedProofResponse); ok {
		switch msg.GenAggregatedProofResponse.Result {
		case pb.Result_RESULT_UNSPECIFIED:
			return nil, fmt.Errorf("failed to aggregate proofs %s, %w, input 1 %s, input 2 %s",
				msg.GenAggregatedProofResponse.String(), ErrUnspecified, inputProof1, inputProof2)
		case pb.Result_RESULT_OK:
			return &msg.GenAggregatedProofResponse.Id, nil
		case pb.Result_RESULT_ERROR:
			return nil, fmt.Errorf("failed to aggregate proofs %s, %w, input 1 %s, input 2 %s",
				msg.GenAggregatedProofResponse.String(), ErrBadRequest, inputProof1, inputProof2)
		case pb.Result_RESULT_INTERNAL_ERROR:
			return nil, fmt.Errorf("failed to aggregate proofs %s, %w, input 1 %s, input 2 %s",
				msg.GenAggregatedProofResponse.String(), ErrProverInternalError, inputProof1, inputProof2)
		default:
			return nil, fmt.Errorf("failed to aggregate proofs %s, %w, input 1 %s, input 2 %s",
				msg.GenAggregatedProofResponse.String(), ErrUnknown, inputProof1, inputProof2)
		}
	}

	return nil, fmt.Errorf("%w, wanted %T, got %T", ErrBadProverResponse, &pb.ProverMessage_GenAggregatedProofResponse{}, res.Response)
}

// FinalProof instructs the prover to generate a final proof for the given
// input. It returns the ID of the proof being computed.
func (p *Prover) FinalProof(inputProof string, aggregatorAddr string) (*string, error) {
	metrics.WorkingProver()

	req := &pb.AggregatorMessage{
		Request: &pb.AggregatorMessage_GenFinalProofRequest{
			GenFinalProofRequest: &pb.GenFinalProofRequest{
				RecursiveProof: inputProof,
				AggregatorAddr: aggregatorAddr,
			},
		},
	}
	res, err := p.call(req)
	if err != nil {
		return nil, err
	}

	if msg, ok := res.Response.(*pb.ProverMessage_GenFinalProofResponse); ok {
		switch msg.GenFinalProofResponse.Result {
		case pb.Result_RESULT_UNSPECIFIED:
			return nil, fmt.Errorf("failed to generate final proof %s, %w, input %s",
				msg.GenFinalProofResponse.String(), ErrUnspecified, inputProof)
		case pb.Result_RESULT_OK:
			return &msg.GenFinalProofResponse.Id, nil
		case pb.Result_RESULT_ERROR:
			return nil, fmt.Errorf("failed to generate final proof %s, %w, input %s",
				msg.GenFinalProofResponse.String(), ErrBadRequest, inputProof)
		case pb.Result_RESULT_INTERNAL_ERROR:
			return nil, fmt.Errorf("failed to generate final proof %s, %w, input %s",
				msg.GenFinalProofResponse.String(), ErrProverInternalError, inputProof)
		default:
			return nil, fmt.Errorf("failed to generate final proof %s, %w, input %s",
				msg.GenFinalProofResponse.String(), ErrUnknown, inputProof)
		}
	}
	return nil, fmt.Errorf("%w, wanted %T, got %T", ErrBadProverResponse, &pb.ProverMessage_GenFinalProofResponse{}, res.Response)
}

// CancelProofRequest asks the prover to stop the generation of the proof
// matching the provided proofID.
func (p *Prover) CancelProofRequest(proofID string) error {
	req := &pb.AggregatorMessage{
		Request: &pb.AggregatorMessage_CancelRequest{
			CancelRequest: &pb.CancelRequest{Id: proofID},
		},
	}
	res, err := p.call(req)
	if err != nil {
		return err
	}
	if msg, ok := res.Response.(*pb.ProverMessage_CancelResponse); ok {
		switch msg.CancelResponse.Result {
		case pb.Result_RESULT_UNSPECIFIED:
			return fmt.Errorf("failed to cancel proof id [%s], %w, %s",
				proofID, ErrUnspecified, msg.CancelResponse.String())
		case pb.Result_RESULT_OK:
			return nil
		case pb.Result_RESULT_ERROR:
			return fmt.Errorf("failed to cancel proof id [%s], %w, %s",
				proofID, ErrBadRequest, msg.CancelResponse.String())
		case pb.Result_RESULT_INTERNAL_ERROR:
			return fmt.Errorf("failed to cancel proof id [%s], %w, %s",
				proofID, ErrProverInternalError, msg.CancelResponse.String())
		default:
			return fmt.Errorf("failed to cancel proof id [%s], %w, %s",
				proofID, ErrUnknown, msg.CancelResponse.String())
		}
	}
	return fmt.Errorf("%w, wanted %T, got %T", ErrBadProverResponse, &pb.ProverMessage_CancelResponse{}, res.Response)
}

// WaitRecursiveProof waits for a recursive proof to be generated by the prover
// and returns it.
func (p *Prover) WaitRecursiveProof(ctx context.Context, proofID string) (string, error) {
	res, err := p.waitProof(ctx, proofID)
	if err != nil {
		return "", err
	}
	resProof := res.Proof.(*pb.GetProofResponse_RecursiveProof)
	return resProof.RecursiveProof, nil
}

// WaitFinalProof waits for the final proof to be generated by the prover and
// returns it.
func (p *Prover) WaitFinalProof(ctx context.Context, proofID string) (*pb.FinalProof, error) {
	res, err := p.waitProof(ctx, proofID)
	if err != nil {
		return nil, err
	}
	resProof := res.Proof.(*pb.GetProofResponse_FinalProof)
	return resProof.FinalProof, nil
}

// waitProof waits for a proof to be generated by the prover and returns the
// prover response.
func (p *Prover) waitProof(ctx context.Context, proofID string) (*pb.GetProofResponse, error) {
	defer metrics.IdlingProver()

	req := &pb.AggregatorMessage{
		Request: &pb.AggregatorMessage_GetProofRequest{
			GetProofRequest: &pb.GetProofRequest{
				// TODO(pg): set Timeout field?
				Id: proofID,
			},
		},
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			res, err := p.call(req)
			if err != nil {
				return nil, err
			}
			if msg, ok := res.Response.(*pb.ProverMessage_GetProofResponse); ok {
				switch msg.GetProofResponse.Result {
				case pb.GetProofResponse_RESULT_PENDING:
					time.Sleep(p.proofStatePollingInterval.Duration)
					continue
				case pb.GetProofResponse_RESULT_UNSPECIFIED:
					return nil, fmt.Errorf("failed to get proof ID: %s, %w, prover response: %s",
						proofID, ErrUnspecified, msg.GetProofResponse.String())
				case pb.GetProofResponse_RESULT_COMPLETED_OK:
					return msg.GetProofResponse, nil
				case pb.GetProofResponse_RESULT_ERROR:
					return nil, fmt.Errorf("failed to get proof with ID %s, %w, prover response: %s",
						proofID, ErrBadRequest, msg.GetProofResponse.String())
				case pb.GetProofResponse_RESULT_COMPLETED_ERROR:
					return nil, fmt.Errorf("failed to get proof with ID %s, %w, prover response: %s",
						proofID, ErrProverCompletedError, msg.GetProofResponse.String())
				case pb.GetProofResponse_RESULT_INTERNAL_ERROR:
					return nil, fmt.Errorf("failed to get proof ID: %s, %w, prover response: %s",
						proofID, ErrProverInternalError, msg.GetProofResponse.String())
				case pb.GetProofResponse_RESULT_CANCEL:
					return nil, fmt.Errorf("proof generation was cancelled for proof ID %s, %w, prover response: %s",
						proofID, ErrProofCanceled, msg.GetProofResponse.String())
				default:
					return nil, fmt.Errorf("failed to get proof ID: %s, %w, prover response: %s",
						proofID, ErrUnknown, msg.GetProofResponse.String())
				}
			}
			return nil, fmt.Errorf("%w, wanted %T, got %T", ErrBadProverResponse, &pb.ProverMessage_GetProofResponse{}, res.Response)
		}
	}
}

// call sends a message to the prover and waits to receive the response over
// the connection stream.
func (p *Prover) call(req *pb.AggregatorMessage) (*pb.ProverMessage, error) {
	if err := p.stream.Send(req); err != nil {
		return nil, err
	}
	res, err := p.stream.Recv()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// FirstToUpper returns the string passed as argument with the first letter in
// uppercase.
func FirstToUpper(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
