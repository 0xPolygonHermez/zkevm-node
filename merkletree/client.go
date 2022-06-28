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
	}
	executorConn, err := grpc.Dial(c.URI, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	executorClient := pb.NewStateDBServiceClient(executorConn)
	return executorClient, executorConn
}
