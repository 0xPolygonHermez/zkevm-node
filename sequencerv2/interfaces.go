//nolint
package sequencerv2

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/ethermanv2"
	"github.com/hermeznetwork/hermez-core/pool"
	"github.com/hermeznetwork/hermez-core/state"
	"github.com/hermeznetwork/hermez-core/state/runtime"
)

// Consumer interfaces required by the package.

// txPool contains the methods required to interact with the tx pool.
type txPool interface {
	UpdateTxState(ctx context.Context, hash common.Hash, newState pool.TxState) error
	IsTxPending(ctx context.Context, hash common.Hash) (bool, error)
}

// stateInterface gathers the methods required to interact with the state.
type stateInterface interface {
	GetLastBatchNumber(ctx context.Context, txBundleID string) (uint64, error)
	GetLastBatchNumberSeenOnEthereum(ctx context.Context, txBundleID string) (uint64, error)

	GetLastL2Block(ctx context.Context) (*state.L2Block, error)
	AddL2Block(ctx context.Context, block *state.L2Block, txBundleID string) error

	ProcessSequence(ctx context.Context, inProgressSequence ethermanv2.Sequence) *runtime.ExecutionResult
}

type txManager interface {
	SequenceBatches(sequences []*ethermanv2.Sequence) error
}
