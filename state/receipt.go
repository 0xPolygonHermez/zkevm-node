package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Receipt represents the results of a transaction.
type Receipt struct {
	types.Receipt
	From common.Address `json:"from"`
	To   common.Address `json:"to"`
}
