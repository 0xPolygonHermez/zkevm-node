package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// GlobalExitRoot struct
type GlobalExitRoot struct {
	GlobalExitRootNum *big.Int
	MainnetExitRoot   common.Hash
	RollupExitRoot    common.Hash
}
