package main

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/hermeznetwork/hermez-core/proverservice/api/proverservice"
)

type zkProverServiceServer struct {
	proverservice.ZKProverServiceServer
	id int
}

const (
	serverProtoVersion = "1"
	serverVersion      = "1"
)

var mockProof = &proverservice.Proof{
	ProofA: []string{"0", "0"},
	ProofB: []*proverservice.ProofB{{Proofs: []string{"0", "0"}}, {Proofs: []string{"0", "0"}}},
	ProofC: []string{"0", "0"},
}

func (zkp *zkProverServiceServer) GenProof(ctx context.Context, request *proverservice.GenProofRequest) (*proverservice.GenProofResponse, error) {
	zkp.id++
	return &proverservice.GenProofResponse{
		Id:     strconv.Itoa(zkp.id),
		Result: proverservice.GenProofResponse_RESULT_GEN_PROOF_OK,
	}, nil
}

func (zkp *zkProverServiceServer) GetProof(server proverservice.ZKProverService_GetProofServer) error {
	newStateRoot, _ := new(big.Int).SetString("1212121212121212121212121212121212121212121212121212121212121212", 16)
	newLocalExitRoot, _ := new(big.Int).SetString("1234123412341234123412341234123412341234123412341234123412341234", 16)
	publicInputs := &proverservice.PublicInputs{
		NewStateRoot:     newStateRoot.String(),
		NewLocalExitRoot: newLocalExitRoot.String(),
	}

	resp := &proverservice.GetProofResponse{
		Id:     strconv.Itoa(zkp.id),
		Proof:  mockProof,
		Public: &proverservice.PublicInputsExtended{PublicInputs: publicInputs},
		Result: proverservice.GetProofResponse_RESULT_GET_PROOF_COMPLETED_OK,
	}
	err := server.Send(resp)
	if err != nil {
		fmt.Printf("Get proof err: %v\n", err)
	}
	return err
}

func (zkp *zkProverServiceServer) GetStatus(ctx context.Context, request *proverservice.GetStatusRequest) (*proverservice.GetStatusResponse, error) {
	return &proverservice.GetStatusResponse{
		State:                     proverservice.GetStatusResponse_STATUS_PROVER_IDLE,
		LastComputedRequestId:     strconv.Itoa(zkp.id),
		LastComputedEndTime:       uint64(time.Now().Unix()),
		CurrentComputingRequestId: strconv.Itoa(zkp.id + 1),
		CurrentComputingStartTime: 0,
		VersionProto:              serverProtoVersion,
		VersionServer:             serverVersion,
		PendingRequestQueueIds:    nil,
	}, nil
}

func (zkp *zkProverServiceServer) Cancel(ctx context.Context, request *proverservice.CancelRequest) (*proverservice.CancelResponse, error) {
	return &proverservice.CancelResponse{Result: proverservice.CancelResponse_RESULT_CANCEL_OK}, nil
}

func (zkp *zkProverServiceServer) Execute(server proverservice.ZKProverService_ExecuteServer) error {
	return nil
}

func (zkp *zkProverServiceServer) SynchronizeBatchProposal(ctx context.Context, request *proverservice.SynchronizeBatchProposalRequest) (*proverservice.SynchronizeBatchProposalResponse, error) {
	return nil, nil
}
