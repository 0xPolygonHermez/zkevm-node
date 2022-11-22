package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	aggregator "github.com/0xPolygonHermez/zkevm-node/aggregator"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := aggregator.ServerConfig{
		Host:                       "0.0.0.0",
		Port:                       8888,
		IntervalToConsolidateState: types.NewDuration(time.Second),
	}
	ctx := context.Background()

	srv := aggregator.NewServer(&cfg)
	srv.Start()

	// connect
	opts := []grpc.DialOption{
		// TODO: once we have user and password for prover server, change this
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	srvConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	// client 1
	cli := pb.NewAggregatorServiceClient(srvConn)
	stream, err := cli.Channel(ctx)
	if err != nil {
		panic(err)
	}

	req, err := stream.Recv()
	if err != nil {
		panic(err)
	}
	log.Debug("[client0] received request from server")
	if _, ok := req.Request.(*pb.AggregatorMessage_GetStatusRequest); ok {
		msg := &pb.ProverMessage{
			Response: &pb.ProverMessage_GetStatusResponse{
				GetStatusResponse: &pb.GetStatusResponse{
					ProverId: "id1",
					Status:   pb.GetStatusResponse_BOOTING,
				},
			},
		}
		if err := stream.Send(msg); err != nil {
			panic(err)
		}
	} else {
		panic("bad request")
	}

	// client 2
	cli2 := pb.NewAggregatorServiceClient(srvConn)
	stream2, err := cli2.Channel(ctx)
	if err != nil {
		panic(err)
	}
	req, err = stream2.Recv()
	if err != nil {
		panic(err)
	}
	log.Debug("[client1] received request from server")
	if _, ok := req.Request.(*pb.AggregatorMessage_GetStatusRequest); ok {
		msg := &pb.ProverMessage{
			Response: &pb.ProverMessage_GetStatusResponse{
				GetStatusResponse: &pb.GetStatusResponse{
					ProverId: "id2",
					Status:   pb.GetStatusResponse_COMPUTING,
				},
			},
		}
		if err := stream2.Send(msg); err != nil {
			panic(err)
		}
	} else {
		panic("bad request")
	}

	time.Sleep(time.Second)

	// client 1 again
	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				panic(err)
			}

			log.Debug("[client0] received request from server")
			if _, ok := req.Request.(*pb.AggregatorMessage_GetStatusRequest); ok {
				msg := &pb.ProverMessage{
					Response: &pb.ProverMessage_GetStatusResponse{
						GetStatusResponse: &pb.GetStatusResponse{
							ProverId: "id1",
							Status:   pb.GetStatusResponse_Status(rand.Int31n(5)),
						},
					},
				}
				if err := stream.Send(msg); err != nil {
					panic(err)
				}
			} else {
				panic("bad request")
			}
		}
	}()

	// client 2 again
	go func() {
		for {
			req, err := stream2.Recv()
			if err != nil {
				panic(err)
			}

			log.Debug("[client1] received request from server")
			if _, ok := req.Request.(*pb.AggregatorMessage_GetStatusRequest); ok {
				msg := &pb.ProverMessage{
					Response: &pb.ProverMessage_GetStatusResponse{
						GetStatusResponse: &pb.GetStatusResponse{
							ProverId: "id2",
							Status:   pb.GetStatusResponse_Status(rand.Int31n(5)),
						},
					},
				}
				if err := stream2.Send(msg); err != nil {
					panic(err)
				}
			} else {
				panic("bad request")
			}
		}
	}()

	for {
	}
}
