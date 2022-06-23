package statev2

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// ForcedBatch represents a ForcedBatch
type ForcedBatch struct {
	BlockNumber       uint64
	ForcedBatchNumber uint64
	Sequencer         common.Address
	GlobalExitRoot    common.Hash
	RawTxsData        string
	ForcedAt          time.Time
}
