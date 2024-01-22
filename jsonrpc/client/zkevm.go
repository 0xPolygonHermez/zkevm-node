package client

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/ethereum/go-ethereum/common"
)

// BatchNumber returns the latest batch number
func (c *Client) BatchNumber(ctx context.Context) (uint64, error) {
	response, err := JSONRPCCall(c.url, "zkevm_batchNumber")
	if err != nil {
		return 0, err
	}

	if response.Error != nil {
		return 0, response.Error.RPCError()
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
	bn := types.LatestBatchNumber
	if number != nil {
		bn = types.BatchNumber(number.Int64())
	}
	response, err := JSONRPCCall(c.url, "zkevm_getBatchByNumber", bn.StringOrHex(), true)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, response.Error.RPCError()
	}

	var result *types.Batch
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ExitRootsByGER returns the exit roots accordingly to the provided Global Exit Root
func (c *Client) ExitRootsByGER(ctx context.Context, globalExitRoot common.Hash) (*types.ExitRoots, error) {
	response, err := JSONRPCCall(c.url, "zkevm_getExitRootsByGER", globalExitRoot.String())
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, response.Error.RPCError()
	}

	var result *types.ExitRoots
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetLatestGlobalExitRoot returns the latest global exit root
func (c *Client) GetLatestGlobalExitRoot(ctx context.Context) (common.Hash, error) {
	response, err := JSONRPCCall(c.url, "zkevm_getLatestGlobalExitRoot")
	if err != nil {
		return common.Hash{}, err
	}

	if response.Error != nil {
		return common.Hash{}, response.Error.RPCError()
	}

	var result string
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return common.Hash{}, err
	}

	return common.HexToHash(result), nil
}
