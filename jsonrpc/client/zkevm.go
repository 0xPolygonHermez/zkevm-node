package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/hex"
)

// BatchNumber returns the latest batch number
func BatchNumber(ctx context.Context, url string) (uint64, error) {
	response, err := JSONRPCCall(url, "zkevm_batchNumber")
	if err != nil {
		return 0, err
	}

	if response.Error != nil {
		return 0, fmt.Errorf("%v %v", response.Error.Code, response.Error.Message)
	}

	var result string
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return 0, err
	}

	decodedBatchNumber, err := hex.DecodeHex(result)
	if err != nil {
		return 0, err
	}

	bigBatchNumber := big.NewInt(0).SetBytes(decodedBatchNumber)
	batchNumber := bigBatchNumber.Uint64()

	return batchNumber, nil
}
