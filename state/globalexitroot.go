package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// GlobalExitRoot struct
type GlobalExitRoot struct {
	BlockNumber       uint64
	GlobalExitRootNum *big.Int
	MainnetExitRoot   common.Hash
	RollupExitRoot    common.Hash
	GlobalExitRoot    common.Hash
}
