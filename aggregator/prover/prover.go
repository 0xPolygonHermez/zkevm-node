package prover

import (
	"github.com/0xPolygonHermez/zkevm-node/log"
	proverclientpb "github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Prover struct
type Prover struct {
	Client proverclientpb.ZKProverServiceClient
	Conn   *grpc.ClientConn
}

// NewProver creates a new Prover
func NewProver(proverURI string) Prover {
	opts := []grpc.DialOption{
		// TODO: once we have user and password for prover server, change this
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	proverConn, err := grpc.Dial(proverURI, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	proverClient := proverclientpb.NewZKProverServiceClient(proverConn)
	return Prover{Client: proverClient, Conn: proverConn}
}
