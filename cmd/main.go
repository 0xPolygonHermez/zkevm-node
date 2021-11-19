package main

import (
	"github.com/hermeznetwork/hermez-core/jsonrpc"
	"github.com/hermeznetwork/hermez-core/mocks"
)

func main() {
	runJSONRpcServer()
}

func runJSONRpcServer() {
	c := jsonrpc.Config{
		Host: "",
		Port: 8123,

		ChainID: 2576980377, // 0x99999999
	}
	p := mocks.NewPool()
	s := mocks.NewState()

	jsonrpc.NewServer(c, p, s).Start()
}
