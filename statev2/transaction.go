package statev2

import "github.com/ethereum/go-ethereum/core/types"

type Transaction struct {
	BatchNumber uint64
	Header      *types.Header
	Uncles      []*types.Header
	types.Transaction
}
