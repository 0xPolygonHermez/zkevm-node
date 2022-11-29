package prover

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	proverclientpb "github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// Prover struct
type Prover struct {
	URI     string
	Client  proverclientpb.ZKProverServiceClient
	Conn    *grpc.ClientConn
	Working bool
}

// NewProver creates a new Prover
func NewProver(proverURI string) *Prover {
	const checkWaitInSeconds = 20
	ctx := context.Background()
	tickerCheckConnection := time.NewTicker(checkWaitInSeconds * time.Second)
	opts := []grpc.DialOption{
		// TODO: once we have user and password for prover server, change this
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	proverConn, err := grpc.Dial(proverURI, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	proverClient := proverclientpb.NewZKProverServiceClient(proverConn)
	prover := &Prover{URI: proverURI, Client: proverClient, Conn: proverConn}
	prover.Working = false

	go func() {
		waitTick(ctx, tickerCheckConnection)
		for {
			prover.checkConnection(ctx, tickerCheckConnection)
		}
	}()

	return prover
}

func (p *Prover) checkConnection(ctx context.Context, ticker *time.Ticker) {
	state := p.Conn.GetState()
	log.Debugf("Checking connection to prover %v. State: %v", p.URI, state)

	if state != connectivity.Ready {
		p.Working = false
		log.Infof("Connection to prover %v seems broken. Trying to reconnect...", p.URI)
		if err := p.Conn.Close(); err != nil {
			log.Errorf("Could not properly close gRPC connection: %v", err)
		}

		opts := []grpc.DialOption{
			// TODO: once we have user and password for prover server, change this
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
		proverConn, err := grpc.Dial(p.URI, opts...)
		if err != nil {
			log.Errorf("Could not reconnect to: %v: %v", p.URI, err)
			waitTick(ctx, ticker)
		}

		p.Client = proverclientpb.NewZKProverServiceClient(proverConn)
		p.Conn = proverConn
	} else {
		p.Working = true
	}

	waitTick(ctx, ticker)
}

func waitTick(ctx context.Context, ticker *time.Ticker) {
	select {
	case <-ticker.C:
		// nothing
	case <-ctx.Done():
		return
	}
}
