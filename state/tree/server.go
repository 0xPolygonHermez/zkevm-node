package tree

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/tree/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Server provides the functionality of the MerkleTree service.
type Server struct {
	cfg   *Config
	stree *StateTree

	srv *grpc.Server
	pb.UnimplementedMTServiceServer
}

// NewServer is the MT server constructor.
func NewServer(cfg *Config, stree *StateTree) *Server {
	return &Server{
		cfg:   cfg,
		stree: stree,
	}
}

// Start sets up the server to process requests.
func (s *Server) Start() {
	address := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.srv = grpc.NewServer()
	pb.RegisterMTServiceServer(s.srv, s)

	healthService := newHealthChecker()
	grpc_health_v1.RegisterHealthServer(s.srv, healthService)

	if err := s.srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Stop stops the server.
func (s *Server) Stop() {
	s.srv.Stop()
}

// Implementation of pb.MTServiceServer interface methods.

// Getters.

// GetBalance gets the balance for a given address at a given root.
func (s *Server) GetBalance(ctx context.Context, in *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	root, err := hex.DecodeString(in.Root)
	if err != nil {
		return nil, err
	}

	balance, err := s.stree.GetBalance(common.HexToAddress(in.EthAddress), root)

	if err != nil {
		return nil, err
	}

	return &pb.GetBalanceResponse{
		Balance: balance.String(),
	}, nil
}

// GetNonce gets nonce for a given address at a given root.
func (s *Server) GetNonce(ctx context.Context, in *pb.GetNonceRequest) (*pb.GetNonceResponse, error) {
	root, err := hex.DecodeString(in.Root)
	if err != nil {
		return nil, err
	}

	nonce, err := s.stree.GetNonce(common.HexToAddress(in.EthAddress), root)
	if err != nil {
		return nil, err
	}

	return &pb.GetNonceResponse{
		Nonce: nonce.Uint64(),
	}, nil
}

// GetCode gets the code for a given address at a given root.
func (s *Server) GetCode(ctx context.Context, in *pb.GetCodeRequest) (*pb.GetCodeResponse, error) {
	root, err := hex.DecodeString(in.Root)
	if err != nil {
		return nil, err
	}

	code, err := s.stree.GetCode(common.HexToAddress(in.EthAddress), root)
	if err != nil {
		return nil, err
	}

	return &pb.GetCodeResponse{
		Code: hex.EncodeToString(code),
	}, nil
}

// GetCodeHash gets code hash for a given address at a given root.
func (s *Server) GetCodeHash(ctx context.Context, in *pb.GetCodeHashRequest) (*pb.GetCodeHashResponse, error) {
	root, err := hex.DecodeString(in.Root)
	if err != nil {
		return nil, err
	}

	hash, err := s.stree.GetCodeHash(common.HexToAddress(in.EthAddress), root)
	if err != nil {
		return nil, err
	}

	return &pb.GetCodeHashResponse{
		Hash: hex.EncodeToString(hash),
	}, nil
}

// GetStorageAt gets smart contract storage for a given address and position at a given root.
func (s *Server) GetStorageAt(ctx context.Context, in *pb.GetStorageAtRequest) (*pb.GetStorageAtResponse, error) {
	root, err := hex.DecodeString(in.Root)
	if err != nil {
		return nil, err
	}

	positionBI := new(big.Int).SetUint64(in.Position)
	value, err := s.stree.GetStorageAt(common.HexToAddress(in.EthAddress), common.BigToHash(positionBI), root)
	if err != nil {
		return nil, err
	}

	return &pb.GetStorageAtResponse{
		Value: value.String(),
	}, nil
}

// ReverseHash reverse a hash of an exisiting Merkletree node.
func (s *Server) ReverseHash(ctx context.Context, in *pb.ReverseHashRequest) (*pb.ReverseHashResponse, error) {
	return nil, nil
}

// Setters

// SetBalance sets the balance for an account at a root.
func (s *Server) SetBalance(ctx context.Context, in *pb.SetBalanceRequest) (*pb.SetBalanceResponse, error) {
	return nil, nil
}

// SetNonce sets the nonce of an account at a root.
func (s *Server) SetNonce(ctx context.Context, in *pb.SetNonceRequest) (*pb.SetNonceResponse, error) {
	return nil, nil
}

// SetCode sets the code for an account at a root.
func (s *Server) SetCode(ctx context.Context, in *pb.SetCodeRequest) (*pb.SetCodeResponse, error) {
	return nil, nil
}

// SetStorageAt sets smart contract storage for an account and position at a root.
func (s *Server) SetStorageAt(ctx context.Context, in *pb.SetStorageAtRequest) (*pb.SetStorageAtResponse, error) {
	return nil, nil
}

// SetHashValue set an entry of the reverse hash table.
func (s *Server) SetHashValue(ctx context.Context, in *pb.SetHashValueRequest) (*pb.SetHashValueResponse, error) {
	return nil, nil
}

// SetHashValueBulk sets many entries of the reverse hash table.
func (s *Server) SetHashValueBulk(ctx context.Context, in *pb.SetHashValueBulkRequest) (*pb.SetHashValueBulkResponse, error) {
	return nil, nil
}

// HealthChecker will provide an implementation of the HealthCheck interface.
type healthChecker struct{}

// NewHealthChecker returns a health checker according to standard package
// grpc.health.v1.
func newHealthChecker() *healthChecker {
	return &healthChecker{}
}

// HealthCheck interface implementation.

// Check returns the current status of the server for unary gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (s *healthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	log.Info("Serving the Check request for health check")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch returns the current status of the server for stream gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (s *healthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	log.Info("Serving the Watch request for health check")
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}
