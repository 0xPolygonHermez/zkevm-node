package vectors

import (
	"github.com/ethereum/go-ethereum/common"
)

// L1InfoTree holds the test vector for the merkle tree
type L1InfoTree struct {
	PreviousLeafValues []common.Hash `json:"previousLeafValues"`
	CurrentRoot         common.Hash   `json:"currentRoot"`
	NewLeafValue        common.Hash   `json:"newLeafValue"`
	NewRoot             common.Hash   `json:"newRoot"`
}
