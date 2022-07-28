package server

import (
	"net"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree/pb"
	"github.com/0xPolygonHermez/zkevm-node/zkevmprovermock/testvector"
	"google.golang.org/grpc"
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
