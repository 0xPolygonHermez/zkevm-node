package merkletree

import (
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/merkletree/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewStateDBServiceClient(c Config) (pb.StateDBServiceClient, *grpc.ClientConn) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}
	stateDBConn, err := grpc.Dial(c.URI, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	stateDBClient := pb.NewStateDBServiceClient(stateDBConn)
	return stateDBClient, stateDBConn
}
