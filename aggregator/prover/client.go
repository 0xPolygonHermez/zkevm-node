package prover

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"google.golang.org/grpc"
)

// Client wrapper for the prover client
type Client struct {
	Prover                                     *Prover
	IntervalFrequencyToGetProofGenerationState types.Duration
}

// NewClient inits prover wrapper client
func NewClient(proverURI string, intervalFrequencyToGetProofGenerationState types.Duration) *Client {
	return &Client{
		Prover: NewProver(proverURI),
		IntervalFrequencyToGetProofGenerationState: intervalFrequencyToGetProofGenerationState,
	}
}

// GetURI return the URI of the prover
func (c *Client) GetURI() string {
	return c.Prover.URI
}

// IsIdle indicates the prover is ready to process requests
func (c *Client) IsIdle(ctx context.Context) bool {
	if !c.Prover.Working {
		return false
	}
	var opts []grpc.CallOption
	status, err := c.Prover.Client.GetStatus(ctx, &pb.GetStatusRequest{}, opts...)
	if err != nil || status.State != pb.GetStatusResponse_STATUS_PROVER_IDLE {
		return false
	}
	return true
}

// GetGenProofID get id of generation proof request
func (c *Client) GetGenProofID(ctx context.Context, inputProver *pb.InputProver) (string, error) {
	genProofRequest := pb.GenProofRequest{Input: inputProver}
	// init connection to the prover
	var opts []grpc.CallOption
	resGenProof, err := c.Prover.Client.GenProof(ctx, &genProofRequest, opts...)
	if err != nil {
		return "", fmt.Errorf("failed to connect to the prover to gen proof, err: %v", err)
	}

	log.Debugf("Data sent to the prover: %+v", inputProver)
	genProofRes := resGenProof.GetResult()
	if genProofRes != pb.GenProofResponse_RESULT_GEN_PROOF_OK {
		return "", fmt.Errorf("failed to get result from the prover, batchNumber: %d, err: %v", inputProver.PublicInputs.BatchNum, err)
	}
	genProofID := resGenProof.GetId()

	return genProofID, err
}

// GetResGetProof get result from proof generation
func (c *Client) GetResGetProof(ctx context.Context, genProofID string, batchNumber uint64) (*pb.GetProofResponse, error) {
	resGetProof := &pb.GetProofResponse{Result: -1}
	getProofCtx, getProofCtxCancel := context.WithCancel(ctx)
	defer getProofCtxCancel()
	getProofClient, err := c.Prover.Client.GetProof(getProofCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to init getProofClient, batchNumber: %d, err: %v", batchNumber, err)
	}
	for resGetProof.Result != pb.GetProofResponse_RESULT_GET_PROOF_COMPLETED_OK {
		err = getProofClient.Send(&pb.GetProofRequest{
			Id: genProofID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to send get proof request to the prover, batchNumber: %d, err: %v", batchNumber, err)
		}

		resGetProof, err = getProofClient.Recv()
		if err != nil {
			return nil, fmt.Errorf("failed to get proof from the prover, batchNumber: %d, err: %v", batchNumber, err)
		}

		resGetProofState := resGetProof.GetResult()
		if resGetProofState == pb.GetProofResponse_RESULT_GET_PROOF_ERROR ||
			resGetProofState == pb.GetProofResponse_RESULT_GET_PROOF_COMPLETED_ERROR {
			log.Fatalf("failed to get a proof for batch, batch number %d", batchNumber)
		}
		if resGetProofState == pb.GetProofResponse_RESULT_GET_PROOF_INTERNAL_ERROR {
			return nil, fmt.Errorf("failed to generate proof for batch, batchNumber: %v, ResGetProofState: %v", batchNumber, resGetProofState)
		}

		if resGetProofState == pb.GetProofResponse_RESULT_GET_PROOF_CANCEL {
			log.Warnf("proof generation was cancelled, batchNumber: %v", batchNumber)
			break
		}

		if resGetProofState == pb.GetProofResponse_RESULT_GET_PROOF_PENDING {
			// in this case aggregator will wait, to send another request
			time.Sleep(c.IntervalFrequencyToGetProofGenerationState.Duration)
		}
	}

	// getProofCtxCancel call closes the connection stream with the prover. This is the only way to close it by client
	getProofCtxCancel()

	return resGetProof, nil
}
