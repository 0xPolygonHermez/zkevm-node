package jsonrpc

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

type StateWrapper struct {
	*state.State
}

// NewBatchProcessor returns an interface of BatchProcessor instead of the real implementation
func (s *StateWrapper) NewBatchProcessor(ctx context.Context, sequencerAddress common.Address, stateRoot []byte, txBundleID string) (BatchProcessorInterface, error) {
	return s.State.NewBatchProcessor(ctx, sequencerAddress, stateRoot, txBundleID)
}
