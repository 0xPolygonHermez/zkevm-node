package server

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/merkletree/pb"
	"github.com/0xPolygonHermez/zkevm-node/tools/zkevmprovermock/testvector"
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

// Set is the mock of the method for setting values in the tree.
func (server *StateDBMock) Set(ctx context.Context, request *pb.SetRequest) (*pb.SetResponse, error) {
	keyBIStr, err := getKeyBIStr(request.Key)
	if err != nil {
		return nil, err
	}

	oldRootStr := merkletree.H4ToString([]uint64{request.OldRoot.Fe0, request.OldRoot.Fe1, request.OldRoot.Fe2, request.OldRoot.Fe3})
	log.Debugf("Set called with key %v, value %v, root %v", keyBIStr, request.Value, oldRootStr)
	_, newRoot, err := server.tvContainer.FindSMTValue(keyBIStr, oldRootStr)
	if err != nil {
		return nil, err
	}
	h4NewRoot, err := merkletree.StringToh4(newRoot)
	if err != nil {
		return nil, err
	}
	return &pb.SetResponse{
		NewRoot: &pb.Fea{Fe0: h4NewRoot[0], Fe1: h4NewRoot[1], Fe2: h4NewRoot[2], Fe3: h4NewRoot[3]},
	}, nil
}

// Get is the mock of the method for getting values from the tree.
func (server *StateDBMock) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	keyBIStr, err := getKeyBIStr(request.Key)
	if err != nil {
		return nil, err
	}

	rootStr := merkletree.H4ToString([]uint64{request.Root.Fe0, request.Root.Fe1, request.Root.Fe2, request.Root.Fe3})

	value, _, err := server.tvContainer.FindSMTValue(keyBIStr, rootStr)
	if err != nil {
		return nil, err
	}
	valueBI, ok := new(big.Int).SetString(value, encoding.Base10)
	if !ok {
		return nil, fmt.Errorf("Could not convert base 10 %q to big.Int", value)
	}
	valueHex := hex.EncodeBig(valueBI)[2:]

	log.Debugf("Get called with key %v, root %v, returning value %v", keyBIStr, rootStr, valueHex)
	return &pb.GetResponse{
		Value: valueHex,
	}, nil
}

// SetProgram is the mock of the method for setting SC contents in the tree.
func (server *StateDBMock) SetProgram(ctx context.Context, request *pb.SetProgramRequest) (*pb.SetProgramResponse, error) {
	keyBIStr, err := getKeyBIStr(request.Key)
	if err != nil {
		return nil, err
	}

	_, err = server.tvContainer.FindBytecode(keyBIStr)
	if err != nil {
		return nil, err
	}
	return &pb.SetProgramResponse{}, nil
}

// GetProgram is the mock of the method for getting SC contents from the tree.
func (server *StateDBMock) GetProgram(ctx context.Context, request *pb.GetProgramRequest) (*pb.GetProgramResponse, error) {
	keyBIStr, err := getKeyBIStr(request.Key)
	if err != nil {
		return nil, err
	}

	bytecode, err := server.tvContainer.FindBytecode(keyBIStr)
	if err != nil {
		return nil, err
	}
	data, err := hex.DecodeHex(bytecode)
	if err != nil {
		return nil, err
	}
	return &pb.GetProgramResponse{
		Data: data,
	}, nil
}

func getKeyBIStr(key *pb.Fea) (string, error) {
	keyStr := merkletree.H4ToString([]uint64{key.Fe0, key.Fe1, key.Fe2, key.Fe3})

	if strings.HasPrefix(keyStr, "0x") { // nolint
		keyStr = keyStr[2:]
	}

	keyBI, ok := new(big.Int).SetString(keyStr, hex.Base)
	if !ok {
		return "", fmt.Errorf("Could not convert the hex string %q into big.Int", keyStr)
	}
	keyBytes := merkletree.ScalarToFilledByteSlice(keyBI)

	return new(big.Int).SetBytes(keyBytes).String(), nil
}
