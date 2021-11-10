package jsonrpc

import "github.com/ethereum/go-ethereum/common"

// txnArgs is the transaction argument for the rpc endpoints
type txnArgs struct {
	From     *common.Address
	To       *common.Address
	Gas      *uint64
	GasPrice *[]byte
	Value    *[]byte
	Input    *[]byte
	Data     *[]byte
	Nonce    *uint64
}
