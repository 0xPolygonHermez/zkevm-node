package broadcast

import (
	"context"
	"fmt"
	"net"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Server provides the functionality of the Broadcast service.
type Server struct {
	cfg *ServerConfig

	srv *grpc.Server
	pb.UnimplementedBroadcastServiceServer
	state stateInterface
}

// NewServer is the Broadcast server constructor.
func NewServer(cfg *ServerConfig, state stateInterface) *Server {
	return &Server{
		cfg:   cfg,
		state: state,
	}
}

// SetState is the state setter.
func (s *Server) SetState(st stateInterface) {
	s.state = st
}

// Start sets up the server to process requests.
func (s *Server) Start() {
	address := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.srv = grpc.NewServer()
	pb.RegisterBroadcastServiceServer(s.srv, s)

	healthService := newHealthChecker()
	grpc_health_v1.RegisterHealthServer(s.srv, healthService)

	log.Infof("Server listening in %q", address)
	if err := s.srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Stop stops the server.
func (s *Server) Stop() {
	s.srv.Stop()
}

// Implementation of pb.BroadcastServiceServer interface methods.

// GetBatch returns a batch by batch number.
func (s *Server) GetBatch(ctx context.Context, in *pb.GetBatchRequest) (*pb.GetBatchResponse, error) {
	batch, err := s.state.GetBatchByNumber(ctx, in.BatchNumber, nil)
	if err != nil {
		return nil, err
	}
	return s.genericGetBatch(ctx, batch)
}

// GetLastBatch returns the last batch.
func (s *Server) GetLastBatch(ctx context.Context, empty *emptypb.Empty) (*pb.GetBatchResponse, error) {
	batch, err := s.state.GetLastBatch(ctx, nil)
	if err != nil {
		return nil, err
	}
	return s.genericGetBatch(ctx, batch)
}

func (s *Server) genericGetBatch(ctx context.Context, batch *state.Batch) (*pb.GetBatchResponse, error) {
	txs, err := s.state.GetEncodedTransactionsByBatchNumber(ctx, batch.BatchNumber, nil)
	if err != nil {
		return nil, err
	}
	transactions := make([]*pb.Transaction, len(txs))
	for i, tx := range txs {
		transactions[i] = &pb.Transaction{
			Encoded: tx,
		}
	}

	var forcedBatchNum uint64
	forcedBatch, err := s.state.GetForcedBatchByBatchNumber(ctx, batch.BatchNumber, nil)
	if err == nil {
		forcedBatchNum = forcedBatch.ForcedBatchNumber
	} else if err != state.ErrNotFound {
		return nil, err
	}

	var mainnetExitRoot, rollupExitRoot string
	ger, err := s.state.GetExitRootByGlobalExitRoot(ctx, batch.GlobalExitRoot, nil)
	if err == nil {
		mainnetExitRoot = ger.MainnetExitRoot.String()
		rollupExitRoot = ger.RollupExitRoot.String()
	} else if err != state.ErrNotFound {
		return nil, err
	}

	return &pb.GetBatchResponse{
		BatchNumber:       batch.BatchNumber,
		GlobalExitRoot:    batch.GlobalExitRoot.String(),
		Sequencer:         batch.Coinbase.String(),
		LocalExitRoot:     batch.LocalExitRoot.String(),
		StateRoot:         batch.StateRoot.String(),
		MainnetExitRoot:   mainnetExitRoot,
		RollupExitRoot:    rollupExitRoot,
		Timestamp:         uint64(batch.Timestamp.Unix()),
		Transactions:      transactions,
		ForcedBatchNumber: forcedBatchNum,
	}, nil
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
