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
func NewClient(ctx context.Context, serverAddress string) (pb.BroadcastServiceClient, *grpc.ClientConn, context.CancelFunc) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}
	const maxWaitSeconds = 120
	ctx, cancel := context.WithTimeout(ctx, maxWaitSeconds*time.Second)
	log.Infof("connecting to broadcast service: %v", serverAddress)
	conn, err := grpc.DialContext(ctx, serverAddress, opts...)
	if err != nil {
		log.Fatalf("failed to connect to broadcast service: %v", err)
	}
	client := pb.NewBroadcastServiceClient(conn)
	log.Info("connected to broadcast service")

	return client, conn, cancel
}
