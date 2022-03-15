package main

import (
	"fmt"
	"net"

	"github.com/hermeznetwork/hermez-core/proverservice/api/proverservice"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()
	zkProverServiceServer := &zkProverServiceServer{id: 0}
	proverservice.RegisterZKProverServiceServer(s, zkProverServiceServer)
	fmt.Println("start a service...")
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
