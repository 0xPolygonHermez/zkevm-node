package synchronizerv2

import (
	"context"

	"github.com/hermeznetwork/hermez-core/ethermanv2"
	state "github.com/hermeznetwork/hermez-core/statev2"
)

// ethermanInterface contains the methods required to interact with ethereum.
type ethermanInterface interface {
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastBlock(ctx context.Context, txBundleID string) (*ethermanv2.Block, error)
	AddGlobalExitRoot(ctx context.Context, exitRoot *state.GlobalExitRoot, txBundleID string) error
	AddForcedBatch(ctx context.Context, forcedBatch *state.ForcedBatch, txBundleID string) error
}
