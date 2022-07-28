package prover

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
)

type Client struct {
	ZkProverClient                                      pb.ZKProverServiceClient
	IntervalFrequencyToGetProofGenerationStateInSeconds types.Duration
}

func NewClient(pc pb.ZKProverServiceClient) *Client {
	return &Client{ZkProverClient: pc}
}

func (c *Client) GetGenProofID(ctx context.Context, inputProver *pb.InputProver) (string, error) {
	genProofRequest := pb.GenProofRequest{Input: inputProver}
	// init connection to the prover
	var opts []grpc.CallOption
	resGenProof, err := c.ZkProverClient.GenProof(ctx, &genProofRequest, opts...)
	if err != nil {
		return "", fmt.Errorf("failed to connect to the prover to gen proof, err: %v", err)
	}

	log.Debugf("Data sent to the prover: %+v", inputProver)
	genProofRes := resGenProof.GetResult()
	if genProofRes != pb.GenProofResponse_RESULT_GEN_PROOF_OK {
		return "", fmt.Errorf("failed to get result from the prover, batchNumber: %d, err: %v", inputProver.PublicInputs.BatchNum)
	}
	genProofID := resGenProof.GetId()

	return genProofID, err
}

func (c *Client) GetResGetProof(ctx context.Context, genProofID string, batchNumber uint64) (*pb.GetProofResponse, error) {
	resGetProof := &pb.GetProofResponse{Result: -1}
	getProofCtx, getProofCtxCancel := context.WithCancel(ctx)
	defer getProofCtxCancel()
	getProofClient, err := c.ZkProverClient.GetProof(getProofCtx)
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
			time.Sleep(c.IntervalFrequencyToGetProofGenerationStateInSeconds.Duration)
		}
	}

	// getProofCtxCancel call closes the connection stream with the prover. This is the only way to close it by client
	getProofCtxCancel()

	return resGetProof, nil
}
