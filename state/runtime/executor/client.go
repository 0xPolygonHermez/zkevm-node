package executor

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const maxMsgSize = 100000000

// NewExecutorClient is the executor client constructor.
func NewExecutorClient(ctx context.Context, c Config) (pb.ExecutorServiceClient, *grpc.ClientConn, context.CancelFunc) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
		grpc.WithBlock(),
	}
	const maxWaitSeconds = 120
	ctx, cancel := context.WithTimeout(ctx, maxWaitSeconds*time.Second)

	log.Infof("trying to connect to executor: %v", c.URI)
	executorConn, err := grpc.DialContext(ctx, c.URI, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	log.Infof("connected to executor")

	executorClient := pb.NewExecutorServiceClient(executorConn)
	return executorClient, executorConn, cancel
}
