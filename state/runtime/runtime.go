package runtime

import "github.com/ethereum/go-ethereum/common"

// TxContext is the context of the transaction
type TxContext struct {
	GasPrice   common.Hash
	Origin     common.Address
	Coinbase   common.Address
	Number     int64
	Timestamp  int64
	GasLimit   int64
	ChainID    int64
	Difficulty common.Hash
}
