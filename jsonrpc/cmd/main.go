package main

import (
	"github.com/hermeznetwork/hermez-core/jsonrpc"
)

func main() {
	cfg := jsonrpc.Config{
		Host: "",
		Port: 8123,
	}
	server := jsonrpc.NewServer(cfg)
	server.Start()
}
