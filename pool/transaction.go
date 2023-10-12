package pool

import (
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// TxStatusPending represents a tx that has not been processed
	TxStatusPending TxStatus = "pending"
	// TxStatusInvalid represents an invalid tx
	TxStatusInvalid TxStatus = "invalid"
	// TxStatusSelected represents a tx that has been selected
	TxStatusSelected TxStatus = "selected"
	// TxStatusFailed represents a tx that has been failed after processing
	TxStatusFailed TxStatus = "failed"
)

// TxStatus represents the state of a tx
type TxStatus string

// String returns a representation of the tx state in a string format
func (s TxStatus) String() string {
	return string(s)
}

// TxStatusUpdateInfo represents the information needed to update the status of a tx
type TxStatusUpdateInfo struct {
	Hash         common.Hash
	NewStatus    TxStatus
	IsWIP        bool
	FailedReason *string
}

// Transaction represents a pool tx
type Transaction struct {
	types.Transaction
	Status TxStatus
	state.ZKCounters
	ReceivedAt            time.Time
	PreprocessedStateRoot common.Hash
	IsWIP                 bool
	IP                    string
	FailedReason          *string
	IsClaims              bool
}

// NewTransaction creates a new transaction
func NewTransaction(tx types.Transaction, ip string, isWIP bool, p *Pool) *Transaction {
	poolTx := Transaction{
		Transaction: tx,
		Status:      TxStatusPending,
		ReceivedAt:  time.Now(),
		IsWIP:       isWIP,
		IP:          ip,
	}

	poolTx.IsClaims = poolTx.IsClaimTx(p.l2BridgeAddr, p.cfg.FreeClaimGasLimit)

	return &poolTx
}

// IsClaimTx checks, if tx is a claim tx
func (tx *Transaction) IsClaimTx(l2BridgeAddr common.Address, freeClaimGasLimit uint64) bool {
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
