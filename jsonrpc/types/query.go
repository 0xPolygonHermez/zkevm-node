package types

import (
	"github.com/ethereum/go-ethereum/common"
)

// LogFilterRequest represents a log filter request.
type LogFilterRequest struct {
	BlockHash *common.Hash  `json:"blockHash,omitempty"`
	FromBlock *string       `json:"fromBlock,omitempty"`
	ToBlock   *string       `json:"toBlock,omitempty"`
	Address   interface{}   `json:"address,omitempty"`
	Topics    []interface{} `json:"topics,omitempty"`
}
