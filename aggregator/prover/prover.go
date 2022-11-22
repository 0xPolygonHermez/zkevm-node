package prover

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

// Prover abstraction of the grpc prover client.
type Prover struct {
	id                                         string
	IntervalFrequencyToGetProofGenerationState types.Duration
	stream                                     pb.AggregatorService_ChannelServer
}

// New returns a new Prover instance.
func New(stream pb.AggregatorService_ChannelServer, intervalFrequencyToGetProofGenerationState types.Duration) (*Prover, error) {
	p := &Prover{
		stream: stream,
		IntervalFrequencyToGetProofGenerationState: intervalFrequencyToGetProofGenerationState,
	}
	status, err := p.Status()
	if err != nil {
		return nil, err
	}
	p.id = status.ProverId
	return p, nil
}

// ID returns the Prover ID.
func (p *Prover) ID() string { return p.id }

// Status gets the prover status.
func (p *Prover) Status() (*pb.GetStatusResponse, error) {
	req := &pb.AggregatorMessage{
		Request: &pb.AggregatorMessage_GetStatusRequest{
			GetStatusRequest: &pb.GetStatusRequest{},
		},
	}
	res, err := p.Call(req)
	if err != nil {
		return nil, err
	}
	if msg, ok := res.Response.(*pb.ProverMessage_GetStatusResponse); ok {
		return msg.GetStatusResponse, nil
	}
	return nil, errors.New("Bad response") // FIXME(pg)
}

// IsIdle returns true if the prover is idling.
func (p *Prover) IsIdle() bool {
	status, err := p.Status()
	if err != nil {
		log.Warnf("Error asking status for prover ID %s: %w", p.ID, err)
		return false
	}
	return status.Status == pb.GetStatusResponse_IDLE
}

// BatchProof instructs the prover to generate a batch proof for the provided
// input. It returns the ID of the proof being computed.
func (p *Prover) BatchProof(input *pb.InputProver) (string, error) {
	req := &pb.AggregatorMessage{
		Request: &pb.AggregatorMessage_GenBatchProofRequest{
			GenBatchProofRequest: &pb.GenBatchProofRequest{Input: input},
		},
	}
	res, err := p.Call(req)
	if err != nil {
		return "", err
	}
	if msg, ok := res.Response.(*pb.ProverMessage_GenBatchProofResponse); ok {
		switch msg.GenBatchProofResponse.Result {
		case pb.Result_OK:
			return msg.GenBatchProofResponse.Id, nil
		case pb.Result_ERROR:
			return "", errors.New("GenBatchProofResponse.Result: ERROR")
		case pb.Result_INTERNAL_ERROR:
			return "", errors.New("GenBatchProofResponse.Result: INTERNAL_ERROR")
		case pb.Result_UNSPECIFIED:
			return "", errors.New("GenBatchProofResponse.Result: UNSPECIFIED")
		default:
			return "", errors.New("GenBatchProofResponse.Result: UNKNOWN")
		}
	}

	return "", errors.New("GenBatchProofResponse.Result: UNKNOWN")
}

// AggregatedProof instructs the prover to generate an aggregated proof from
// the two inputs provided. It returns the ID of the proof being computed.
func (p *Prover) AggregatedProof(inputProof1, inputProof2 string) (string, error) {
	req := &pb.AggregatorMessage{
		Request: &pb.AggregatorMessage_GenAggregatedProofRequest{
			GenAggregatedProofRequest: &pb.GenAggregatedProofRequest{
				RecursiveProof_1: inputProof1,
				RecursiveProof_2: inputProof2,
			},
		},
	}
	res, err := p.Call(req)
	if err != nil {
		return "", err
	}
	if msg, ok := res.Response.(*pb.ProverMessage_GenAggregatedProofResponse); ok {
		switch msg.GenAggregatedProofResponse.Result {
		case pb.Result_UNSPECIFIED:
			// TODO(pg): handle this case
		case pb.Result_OK:
			return msg.GenAggregatedProofResponse.Id, nil
		case pb.Result_ERROR:
			return "", errors.New("Prover error")
		case pb.Result_INTERNAL_ERROR:
			return "", errors.New("Prover internal error")
		}
	}
	return "", errors.New("Bad response") // FIXME(pg)
}

// FinalProof instructs the prover to generate a final proof for the given
// input. It returns the ID of the proof being computed.
func (p *Prover) FinalProof(inputProof string) (string, error) {
	req := &pb.AggregatorMessage{
		Request: &pb.AggregatorMessage_GenFinalProofRequest{
			GenFinalProofRequest: &pb.GenFinalProofRequest{RecursiveProof: inputProof},
		},
	}
	res, err := p.Call(req)
	if err != nil {
		return "", err
	}
	if msg, ok := res.Response.(*pb.ProverMessage_GenFinalProofResponse); ok {
		switch msg.GenFinalProofResponse.Result {
		case pb.Result_UNSPECIFIED:
			// TODO(pg): handle this case
		case pb.Result_OK:
			return msg.GenFinalProofResponse.Id, nil
		case pb.Result_ERROR:
			return "", errors.New("Prover error")
		case pb.Result_INTERNAL_ERROR:
			return "", errors.New("Prover internal error")
		}
	}
	return "", errors.New("Bad response") // FIXME(pg)
}

// CancelProofRequest asks the prover to stop the generation of the proof
// matching the provided proofID.
func (p *Prover) CancelProofRequest(proofID string) error {
	req := &pb.AggregatorMessage{
		Request: &pb.AggregatorMessage_CancelRequest{
			CancelRequest: &pb.CancelRequest{Id: proofID},
		},
	}
	res, err := p.Call(req)
	if err != nil {
		return err
	}
	if msg, ok := res.Response.(*pb.ProverMessage_CancelResponse); ok {
		// TODO(pg): handle all cases
		switch msg.CancelResponse.Result {
		case pb.Result_UNSPECIFIED:
		case pb.Result_OK:
			return nil
		case pb.Result_ERROR:
			return errors.New("Prover error")
		case pb.Result_INTERNAL_ERROR:
			return errors.New("Prover internal error")
		}
	}
	return errors.New("Bad response") // FIXME(pg)
}

// WaitRecursiveProof waits for a recursive proof to be generated by the prover
// and returns it.
func (p *Prover) WaitRecursiveProof(ctx context.Context, proofID string) (*pb.GetProofResponse_RecursiveProof, error) {
	res, err := p.WaitProof(ctx, proofID)
	if err != nil {
		return nil, err
	}
	resProof := res.Proof.(*pb.GetProofResponse_RecursiveProof)
	return resProof, nil
}

// WaitFinalProof waits for the final proof to be generated by the prover and
// returns it.
func (p *Prover) WaitFinalProof(ctx context.Context, proofID string) (*pb.GetProofResponse_FinalProof, error) {
	res, err := p.WaitProof(ctx, proofID)
	if err != nil {
		return nil, err
	}
	resProof := res.Proof.(*pb.GetProofResponse_FinalProof)
	return resProof, nil
}

// waitProof waits for a proof to be generated by the prover and returns the
// prover response.
func (p *Prover) WaitProof(ctx context.Context, proofID string) (*pb.GetProofResponse, error) {
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
			res, err := p.Call(req)
			if err != nil {
				return nil, err
			}
			if msg, ok := res.Response.(*pb.ProverMessage_GetProofResponse); ok {
				switch msg.GetProofResponse.Result {
				case pb.GetProofResponse_PENDING:
					time.Sleep(p.IntervalFrequencyToGetProofGenerationState.Duration)
					continue
				case pb.GetProofResponse_UNSPECIFIED:
					return nil, fmt.Errorf("Failed to generate proof ID: %s, ResGetProofState: %v", proofID, msg.GetProofResponse)
				case pb.GetProofResponse_COMPLETED_OK:
					return msg.GetProofResponse, nil
				case pb.GetProofResponse_ERROR, pb.GetProofResponse_COMPLETED_ERROR:
					log.Fatalf("Failed to get proof with ID %s", proofID)
				case pb.GetProofResponse_INTERNAL_ERROR:
					return nil, fmt.Errorf("Failed to generate proof ID: %s, ResGetProofState: %v", proofID, msg.GetProofResponse)
				case pb.GetProofResponse_CANCEL:
					log.Warnf("Proof generation was cancelled for proof ID %s", proofID)
					return msg.GetProofResponse, nil
				}
			}
		}
	}
}

// Call sends a message to the prover and waits to receive the response over
// the connection stream.
func (p *Prover) Call(req *pb.AggregatorMessage) (*pb.ProverMessage, error) {
	if err := p.stream.Send(req); err != nil {
		return nil, err
	}
	res, err := p.stream.Recv()
	if err != nil {
		return nil, err
	}
	return res, nil
}
