package txselector

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state/runtime"
)

// Consumer interfaces required by the package.

// batchProcessor includes the methods required to process batches.
type batchProcessor interface {
	ProcessTransaction(ctx context.Context, tx *types.Transaction, sequencerAddress common.Address) *runtime.ExecutionResult
}
