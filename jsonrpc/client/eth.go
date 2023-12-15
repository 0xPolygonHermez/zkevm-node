package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

// BlockNumber returns the latest block number
func (c *Client) BlockNumber(ctx context.Context) (uint64, error) {
	response, err := JSONRPCCall(c.url, "eth_blockNumber")
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

	bigBlockNumber := hex.DecodeBig(result)
	blockNumber := bigBlockNumber.Uint64()

	return blockNumber, nil
}

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
func (c *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	bn := types.LatestBlockNumber
	if number != nil {
		bn = types.BlockNumber(number.Int64())
	}

	response, err := JSONRPCCall(c.url, "eth_getBlockByNumber", bn.StringOrHex(), true)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("%v %v", response.Error.Code, response.Error.Message)
	}

	var result *types.Block
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
