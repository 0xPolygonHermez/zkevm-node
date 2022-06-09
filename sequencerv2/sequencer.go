//nolint
package sequencerv2

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Sequence represents an operation sent to the PoE smart contract to be
// processed.
type Sequence struct {
	globalExitRoot  common.Hash
	timestamp       uint64
	forceBatchesNum uint64
	txs             []types.Transaction
}
