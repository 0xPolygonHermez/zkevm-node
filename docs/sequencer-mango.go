package docs

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/ethereum/go-ethereum/common"
)

type WorkerPool map[common.Address]AddrQueue // Replace map for sorted map. TODO: find good library

type AddrQueue struct {
	CurrentNonce   uint64
	CurrentBalance *big.Int
	ReadyTxs       []TxTracker
	NotReadyTxs    []TxTracker
}

type TxTracker struct {
	Nonce      uint64
	Benefit    *big.Int        // GasLimit * GasPrice
	ZKCounters pool.ZkCounters // To check if it fits into a batch
	Size       uint64          // To check if it fits into a batch
	Gas        uint64          // To check if it fits into a batch
	Efficiency float64         // To sort. TODO: calculate Benefit / Cost. Cost = some formula taking into account ZKC and Byte Size
	RawTx      []byte
}
