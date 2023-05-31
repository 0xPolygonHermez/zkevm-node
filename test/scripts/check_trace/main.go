package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	url               = ""
	blockRange uint64 = 1000
)

func main() {
	client, err := ethclient.Dial(url)
	chkErr(err)

	ctx := context.Background()

	fmt.Println("getting latest block number")

	latestBlockNumber, err := client.BlockNumber(ctx)
	chkErr(err)

	fmt.Printf("latest block number: %v\n", latestBlockNumber)

	from, to := int64(latestBlockNumber)-int64(blockRange), int64(latestBlockNumber)
	if from < 0 {
		from = 0
	}

	fmt.Printf("reading blocks from %v to %v\n", from, to)
	for blockNumber := from; blockNumber <= to; blockNumber++ {
		fmt.Println()
		fmt.Printf("getting block %v\n", blockNumber)
		block, err := client.BlockByNumber(ctx, big.NewInt(blockNumber))
		chkErr(err)

		for _, tx := range block.Transactions() {
			fmt.Printf("getting default trace for tx %v", tx.Hash().String())
			getTrace(url, tx.Hash(), map[string]interface{}{
				"disableStorage":   false,
				"disableStack":     false,
				"enableMemory":     true,
				"enableReturnData": true,
			})
			chkErr(err)
			fmt.Print(": success\n")

			fmt.Printf("getting callTracer trace for tx %v", tx.Hash().String())
			getTrace(url, tx.Hash(), map[string]interface{}{
				"tracer": "callTracer",
				"tracerConfig": map[string]interface{}{
					"onlyTopCall": false,
					"withLog":     true,
				},
			})
			chkErr(err)
			fmt.Print(": success\n")
		}
	}
}

func getTrace(url string, hash common.Hash, tracerConfig map[string]interface{}) map[string]interface{} {
	response, err := client.JSONRPCCall(url, "debug_traceTransaction", hash.String(), tracerConfig)
	chkErr(err)

	if response.Error != nil {
		err := fmt.Errorf("%v %v %v", response.Error.Code, response.Error.Message, response.Error.Data)
		panic(err)
	}

	if response.Result == nil {
		panic("trace not found")
	}

	var result map[string]interface{}
	err = json.Unmarshal(response.Result, &result)
	chkErr(err)

	return result
}

func chkErr(err error) {
	if err != nil {
		panic(err)
	}
}
