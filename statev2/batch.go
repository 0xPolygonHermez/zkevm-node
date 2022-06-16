package statev2

import (
	"github.com/ethereum/go-ethereum/common"
)

// VerifyBatch represents a VerifyBatch
type VerifyBatch struct {
	BlockNumber       uint64
	BatchNumber       uint64
	Aggregator        common.Address
	TxHash            common.Hash
}
