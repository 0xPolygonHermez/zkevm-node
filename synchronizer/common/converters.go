package common

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// RpcBatchToStateBatch converts a rpc batch to a state batch
func RpcBatchToStateBatch(rpcBatch *types.Batch) *state.Batch {
	if rpcBatch == nil {
		return nil
	}
	batch := &state.Batch{
		BatchNumber:    uint64(rpcBatch.Number),
		Coinbase:       rpcBatch.Coinbase,
		StateRoot:      rpcBatch.StateRoot,
		BatchL2Data:    rpcBatch.BatchL2Data,
		GlobalExitRoot: rpcBatch.GlobalExitRoot,
		LocalExitRoot:  rpcBatch.LocalExitRoot,
		Timestamp:      time.Unix(int64(rpcBatch.Timestamp), 0),
		WIP:            !rpcBatch.Closed,
	}
	if rpcBatch.ForcedBatchNumber != nil {
		batch.ForcedBatchNum = (*uint64)(rpcBatch.ForcedBatchNumber)
	}
	return batch
}
