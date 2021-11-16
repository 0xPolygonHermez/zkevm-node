package main

import (
	"github.com/hermeznetwork/hermez-core/jsonrpc"
)

func main() {
	runJSONRpcServer()
}

func runJSONRpcServer() {
	jsonrpc.NewServer(jsonrpc.Config{
		Host: "",
		Port: 8123,

		ChainID: 2576980377, // 0x99999999
	}).Start()
}
