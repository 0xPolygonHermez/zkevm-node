package state

import (
	"github.com/ethereum/go-ethereum/common"
)

// GlobalExitRoot struct
type GlobalExitRoot struct {
	BlockNumber     uint64
	MainnetExitRoot common.Hash
	RollupExitRoot  common.Hash
	GlobalExitRoot  common.Hash
}
