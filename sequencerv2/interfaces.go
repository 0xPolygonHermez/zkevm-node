//nolint
package sequencerv2

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/ethermanv2"
	"github.com/hermeznetwork/hermez-core/pool"
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

	GetLastL2BlockHash(ctx context.Context) (common.Hash, error)
	AddL2Block(ctx context.Context, txHash common.Hash, parentTxHash common.Hash, txReceivedAt time.Time, txBundleID string) error

	ProcessSequence(ctx context.Context, inProgressSequence ethermanv2.Sequence) *runtime.ExecutionResult
}

type txManager interface {
	SequenceBatches(sequences []*ethermanv2.Sequence) error
}
