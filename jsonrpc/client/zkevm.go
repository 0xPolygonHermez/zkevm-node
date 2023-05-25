package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

// BatchNumber returns the latest batch number
func (c *Client) BatchNumber(ctx context.Context) (uint64, error) {
	response, err := JSONRPCCall(c.url, "zkevm_batchNumber")
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

	bigBatchNumber := hex.DecodeBig(result)
	batchNumber := bigBatchNumber.Uint64()

	return batchNumber, nil
}

// BatchByNumber returns a batch from the current canonical chain. If number is nil, the
// latest known batch is returned.
func (c *Client) BatchByNumber(ctx context.Context, number *big.Int) (*types.Batch, error) {
	response, err := JSONRPCCall(c.url, "zkevm_getBatchByNumber", types.ToBatchNumArg(number), true)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("%v %v", response.Error.Code, response.Error.Message)
	}

	var result *types.Batch
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
