package server

import (
	"context"
	"net"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"github.com/0xPolygonHermez/zkevm-node/tools/zkevmprovermock/testvector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ExecutorMock represents and Executor mock server
type ExecutorMock struct {
	// address is the address on which the gRPC server will listen, eg. 0.0.0.0:50061
	address string

	tvContainer *testvector.Container

	// srv is an insance of the gRPC server.
	srv *grpc.Server

	// embedding an instance of pb.UnimplementedExecutorServiceServer will allow us
	// to implement all the required method interfaces.
	pb.UnimplementedExecutorServiceServer
}

// NewExecutorMock is the ExecutorMock constructor.
func NewExecutorMock(address string, tvContainer *testvector.Container) *ExecutorMock {
	return &ExecutorMock{
		address:     address,
		tvContainer: tvContainer,
	}
}

// Start sets up the stateDB server to process requests.
func (server *ExecutorMock) Start() {
	lis, err := net.Listen("tcp", server.address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server.srv = grpc.NewServer()
	pb.RegisterExecutorServiceServer(server.srv, server)

	log.Infof("Executor mock server: listening at %s", server.address)
	if err := server.srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Stop stops the server.
func (server *ExecutorMock) Stop() {
	log.Info("Executor mock server: stopping...")
	server.srv.Stop()
}

// ProcessBatch implements the ProcessBatch gRPC method.
func (server *ExecutorMock) ProcessBatch(ctx context.Context, request *pb.ProcessBatchRequest) (*pb.ProcessBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessBatch not implemented")
}
