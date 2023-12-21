package syncinterfaces

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
)

// ZkEVMClientInterface contains the methods required to interact with zkEVM-RPC
type ZKEVMClientTrustedBatchesGetter interface {
	BatchNumber(ctx context.Context) (uint64, error)
	BatchByNumber(ctx context.Context, number *big.Int) (*types.Batch, error)
}
