package sequencer

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/pool"
)

// Consumer interfaces required by the package.

// sequencerTxPool contains the methods required to interact with the tx pool.
type sequencerTxPool interface {
	GetPendingTxs(ctx context.Context) ([]pool.Transaction, error)
	UpdateTxState(ctx context.Context, hash common.Hash, newState pool.TxState) error
	UpdateTxsState(ctx context.Context, hashes []string, newState pool.TxState) error
	SetGasPrice(ctx context.Context, gasPrice uint64) error
}
