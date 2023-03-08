package broadcast

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/broadcast/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient creates a grpc client to communicates with the Broadcast server
func NewClient(ctx context.Context, serverAddress string) (pb.BroadcastServiceClient, *grpc.ClientConn, context.CancelFunc, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}
	const maxWaitSeconds = 120
	ctx2, cancel := context.WithTimeout(ctx, maxWaitSeconds*time.Second)
	log.Infof("connecting to broadcast service: %v", serverAddress)
	conn, err := grpc.DialContext(ctx2, serverAddress, opts...)
	if err != nil {
		log.Errorf("failed to connect to broadcast service: %v", err)
		return nil, nil, cancel, err
	}
	client := pb.NewBroadcastServiceClient(conn)
	log.Info("connected to broadcast service")

	return client, conn, cancel, nil
}
