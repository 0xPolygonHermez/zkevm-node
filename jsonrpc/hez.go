package jsonrpc

import (
	"context"

	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/state"
)

// Hez contains implementations for the "hez" RPC endpoints
type Hez struct {
	defaultChainID uint64
	state          state.State
}

// DefaultChainID returns the default chain id that is allowed to be used by all the sequencers
func (h *Hez) DefaultChainID() (interface{}, error) {
	return hex.EncodeUint64(h.defaultChainID), nil
}

// ConsolidatedBlockNumber returns current block number for consolidated batches
func (h *Hez) ConsolidatedBlockNumber() (interface{}, error) {
	ctx := context.Background()

	lastConsolidatedBatch, err := h.state.GetLastBatch(ctx, false)
	if err != nil {
		return nil, err
	}

	lastConsolidatedBatchNumber := uint64(0)
	if lastConsolidatedBatch != nil {
		lastConsolidatedBatchNumber = lastConsolidatedBatch.BatchNumber
	}

	return hex.EncodeUint64(lastConsolidatedBatchNumber), nil
}
