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

// DefaultChainId returns the default chain id that is allowed to be used by all the sequencers
func (h *Hez) DefaultChainId() (interface{}, error) { //nolint:golint
	return hex.EncodeUint64(h.defaultChainID), nil
}

// ConsolidatedBlockNumber returns current block number for consolidated batches
func (h *Hez) ConsolidatedBlockNumber() (interface{}, error) {
	ctx := context.Background()

	lastBatchNumber, err := h.state.GetLastBatchNumber(ctx, false)
	if err != nil {
		return nil, err
	}

	return hex.EncodeUint64(lastBatchNumber), nil
}
