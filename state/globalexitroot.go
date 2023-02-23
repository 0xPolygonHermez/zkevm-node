package state

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// GlobalExitRoot struct
type GlobalExitRoot struct {
	BlockNumber     uint64
	Timestamp       time.Time
	MainnetExitRoot common.Hash
	RollupExitRoot  common.Hash
	GlobalExitRoot  common.Hash
}
