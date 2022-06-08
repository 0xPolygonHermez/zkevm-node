package synchronizerv2

import (
	"context"

	"github.com/hermeznetwork/hermez-core/state"
)

// localEtherman contains the methods required to interact with ethereum.
type localEtherman interface {
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastBlock(ctx context.Context, txBundleID string) (*state.Block, error)
	SetGenesis(ctx context.Context, genesis state.Genesis, txBundleID string) error
}
