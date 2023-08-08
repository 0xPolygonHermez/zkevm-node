package merkletree

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree/hashdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewMTDBServiceClient creates a new MTDB client.
func NewMTDBServiceClient(ctx context.Context, c Config) (hashdb.HashDBServiceClient, *grpc.ClientConn, context.CancelFunc) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}
	const maxWaitSeconds = 120
	ctx, cancel := context.WithTimeout(ctx, maxWaitSeconds*time.Second)

	log.Infof("trying to connect to merkletree: %v", c.URI)
	mtDBConn, err := grpc.DialContext(ctx, c.URI, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	log.Infof("connected to merkletree")

	mtDBClient := hashdb.NewHashDBServiceClient(mtDBConn)
	return mtDBClient, mtDBConn, cancel
}
