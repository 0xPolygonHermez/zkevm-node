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
	const maxRetries = 5
	ctx, cancel := context.WithTimeout(ctx, maxWaitSeconds*time.Second)

	connections := 0

	var executorConn *grpc.ClientConn
	var err error
	for connections < maxRetries {
		log.Infof("trying to connect to executor: %v", c.URI)
		executorConn, err = grpc.DialContext(ctx, c.URI, opts...)
		if err != nil {
			log.Infof("Retrying connection to executor #%d", connections)
			connections = connections + 1
		} else {
			log.Infof("connected to executor")
			break
		}
	}

	if connections == maxRetries {
		log.Fatalf("fail to dial: %v", err)
	}
	executorClient := pb.NewExecutorServiceClient(executorConn)
	return executorClient, executorConn, cancel
}
