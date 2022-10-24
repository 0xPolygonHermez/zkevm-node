package pool

import (
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// TxStatusPending represents a tx that has not been processed
	TxStatusPending TxStatus = "pending"
	// TxStatusInvalid represents an invalid tx
	TxStatusInvalid TxStatus = "invalid"
	// TxStatusSelected represents a tx that has been	 selected
	TxStatusSelected TxStatus = "selected"
	// TxStatusFailed represents a tx that has been failed after processing, but can be processed in the future
	TxStatusFailed TxStatus = "failed"
)

// TxStatus represents the state of a tx
type TxStatus string

// String returns a representation of the tx state in a string format
func (s TxStatus) String() string {
	return string(s)
}

// Transaction represents a pool tx
type Transaction struct {
	types.Transaction
	Status   TxStatus
	IsClaims bool
	ZkCounters
	FailedCounter uint64
	ReceivedAt    time.Time
}

// ZkCounters counters for the tx
type ZkCounters struct {
	CumulativeGasUsed    int64
	UsedKeccakHashes     int32
	UsedPoseidonHashes   int32
	UsedPoseidonPaddings int32
	UsedMemAligns        int32
	UsedArithmetics      int32
	UsedBinaries         int32
	UsedSteps            int32
}

// IsZkCountersBelowZero checks if any of the counters are below zero
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

// SumUpZkCounters sum ups zk counters with passed tx zk counters
func (zkc *ZkCounters) SumUpZkCounters(txZkCounters ZkCounters) {
	zkc.CumulativeGasUsed += txZkCounters.CumulativeGasUsed
	zkc.UsedKeccakHashes += txZkCounters.UsedKeccakHashes
	zkc.UsedPoseidonHashes += txZkCounters.UsedPoseidonHashes
	zkc.UsedPoseidonPaddings += txZkCounters.UsedPoseidonPaddings
	zkc.UsedMemAligns += txZkCounters.UsedMemAligns
	zkc.UsedArithmetics += txZkCounters.UsedArithmetics
	zkc.UsedBinaries += txZkCounters.UsedBinaries
	zkc.UsedSteps += txZkCounters.UsedSteps
}

// IsClaimTx checks, if tx is a claim tx
func (tx *Transaction) IsClaimTx(l2BridgeAddr common.Address) bool {
	if tx.To() == nil {
		return false
	}

	if *tx.To() == l2BridgeAddr &&
		strings.HasPrefix("0x"+common.Bytes2Hex(tx.Data()), bridgeClaimMethodSignature) {
		return true
	}
	return false
}
