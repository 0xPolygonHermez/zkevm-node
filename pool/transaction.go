package pool

import (
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// TxStatePending represents a tx that has not been processed
	TxStatePending TxState = "pending"
	// TxStateInvalid represents an invalid tx
	TxStateInvalid TxState = "invalid"
	// TxStateSelected represents a tx that has been selected
	TxStateSelected TxState = "selected"
)

// TxState represents the state of a tx
type TxState string

// String returns a representation of the tx state in a string format
func (s TxState) String() string {
	return string(s)
}

// Transaction represents a pool tx
type Transaction struct {
	types.Transaction
	State    TxState
	IsClaims bool
	ZkCounters
	ReceivedAt time.Time
}

type ZkCounters struct {
	CumulativeGasUsed    uint64
	UsedKeccakHashes     uint32
	UsedPoseidonHashes   uint32
	UsedPoseidonPaddings uint32
	UsedMemAligns        uint32
	UsedArithmetics      uint32
	UsedBinaries         uint32
	UsedSteps            uint32
}

func (zkc *ZkCounters) IsZkCountersBelowZero() bool {
	return zkc.CumulativeGasUsed < 0 ||
		zkc.UsedArithmetics < 0 ||
		zkc.UsedSteps < 0 ||
		zkc.UsedBinaries < 0 ||
		zkc.UsedMemAligns < 0 ||
		zkc.UsedPoseidonPaddings < 0 ||
		zkc.UsedPoseidonHashes < 0 ||
		zkc.UsedKeccakHashes < 0
}

// IsClaimTx checks, if tx is a claim tx
func (tx *Transaction) IsClaimTx(l2GlobalExitRootManagerAddr common.Address) bool {
	if tx.To() == nil {
		return false
	}

	if *tx.To() == l2GlobalExitRootManagerAddr &&
		strings.HasPrefix("0x"+common.Bytes2Hex(tx.Data()), bridgeClaimMethodSignature) {
		return true
	}
	return false
}
