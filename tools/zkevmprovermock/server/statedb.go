package server

import (
	"context"
	"net"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/merkletree/pb"
	"github.com/0xPolygonHermez/zkevm-node/tools/zkevmprovermock/testvector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StateDBMock represents a StateDB mock server.
type StateDBMock struct {
	// address is the address on which the gRPC server will listen, eg. 0.0.0.0:50061
	address string

	tvContainer *testvector.Container

	// srv is an insance of the gRPC server.
	srv *grpc.Server
	// embedding an instance of pb.UnimplementedStateDBServiceServer will allow us
	// to implement all the required method interfaces.
	pb.UnimplementedStateDBServiceServer
}

// NewStateDBMock is the StateDBMock constructor.
func NewStateDBMock(address string, tvContainer *testvector.Container) *StateDBMock {
	return &StateDBMock{
		address:     address,
		tvContainer: tvContainer,
	}
}

// Start sets up the stateDB server to process requests.
func (server *StateDBMock) Start() {
	lis, err := net.Listen("tcp", server.address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server.srv = grpc.NewServer()
	pb.RegisterStateDBServiceServer(server.srv, server)

	log.Infof("StateDB mock server: listening at %s", server.address)
	if err := server.srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Stop stops the server.
func (server *StateDBMock) Stop() {
	log.Info("StateDB mock server: stopping...")
	server.srv.Stop()
}

// Set is the mock of the method for setting values in the tree.
func (server *StateDBMock) Set(ctx context.Context, request *pb.SetRequest) (*pb.SetResponse, error) {
	log.Infof("Set called")

	_, newRoot, err := server.tvContainer.FindE2EGenesisRaw(request.Value, request.OldRoot.String())
	if err != nil {
		return nil, err
	}
	feaNewRoot, err := merkletree.String2fea(newRoot)
	if err != nil {
		return nil, err
	}
	return &pb.SetResponse{
		NewRoot: &pb.Fea{Fe0: feaNewRoot[0], Fe1: feaNewRoot[1], Fe2: feaNewRoot[2], Fe3: feaNewRoot[3]},
	}, nil
}

// Get is the mock of the method for getting values from the tree.
func (server *StateDBMock) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}

// SetProgram is the mock of the method for setting SC contents in the tree.
func (server *StateDBMock) SetProgram(ctx context.Context, request *pb.SetProgramRequest) (*pb.SetProgramResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetProgram not implemented")
}

// GetProgram is the mock of the method for getting SC contents from the tree.
func (server *StateDBMock) GetProgram(ctx context.Context, request *pb.GetProgramRequest) (*pb.GetProgramResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProgram not implemented")
}
