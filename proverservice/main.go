package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/0xPolygonHermez/zkevm-node/proverservice/pb"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		os.Exit(1)
	}

	s := grpc.NewServer()

	zkProverServiceServer := NewZkProverServiceServer()
	pb.RegisterZKProverServiceServer(s, zkProverServiceServer)

	go func() {
		fmt.Println("starting health service...")
		http.HandleFunc("/health", health)
		if err = http.ListenAndServe(":50052", nil); err != nil {
			fmt.Printf("failed to serve: %v\n", err)
		}
	}()

	fmt.Println("start a service...")
	if err = s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
