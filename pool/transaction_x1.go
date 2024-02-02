package pool

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

var l2BridgeAddr common.Address

// SetL2BridgeAddr sets the L2 bridge address
func SetL2BridgeAddr(value common.Address) {
	l2BridgeAddr = value
}

// IsClaimTx checks, if tx is a claim tx
func (tx *Transaction) IsClaimTx(freeClaimGasLimit uint64) bool {
	if tx.To() == nil {
		return false
	}

	txGas := tx.Gas()
	if txGas > freeClaimGasLimit {
		return false
	}

	if *tx.To() == l2BridgeAddr &&
		strings.HasPrefix("0x"+common.Bytes2Hex(tx.Data()), BridgeClaimMethodSignature) {
		return true
	}
	return false
}
