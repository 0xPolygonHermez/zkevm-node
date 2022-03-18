package main

import (
	"fmt"
	"net"
	"os"

	"github.com/hermeznetwork/hermez-core/proverservice/pb"
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
	fmt.Println("start a service...")
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
