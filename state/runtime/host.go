package runtime

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Host is the execution host
type Host interface {
	AccountExists(ctx context.Context, addr common.Address) bool
	GetStorage(ctx context.Context, addr common.Address, key common.Hash) common.Hash
	SetStorage(ctx context.Context, addr common.Address, key *big.Int, value *big.Int, config *ForksInTime) StorageStatus
	GetBalance(ctx context.Context, addr common.Address) *big.Int
	GetCodeSize(ctx context.Context, addr common.Address) int
	GetCodeHash(ctx context.Context, addr common.Address) common.Hash
	GetCode(ctx context.Context, addr common.Address) []byte
	Selfdestruct(ctx context.Context, addr common.Address, beneficiary common.Address)
	GetTxContext() TxContext
	GetBlockHash(number int64) common.Hash
	EmitLog(addr common.Address, topics []common.Hash, data []byte)
	Callx(context.Context, *Contract, Host) *ExecutionResult
	Empty(ctx context.Context, addr common.Address) bool
	GetNonce(ctx context.Context, addr common.Address) uint64
}

// StorageStatus is the status of the storage access
type StorageStatus int

const (
	// StorageUnchanged if the data has not changed
	StorageUnchanged StorageStatus = iota
	// StorageModified if the value has been modified
	StorageModified
	// StorageModifiedAgain if the value has been modified before in the txn
	StorageModifiedAgain
	// StorageAdded if this is a new entry in the storage
	StorageAdded
	// StorageDeleted if the storage was deleted
	StorageDeleted
)

func (s StorageStatus) String() string {
	switch s {
	case StorageUnchanged:
		return "StorageUnchanged"
	case StorageModified:
		return "StorageModified"
	case StorageModifiedAgain:
		return "StorageModifiedAgain"
	case StorageAdded:
		return "StorageAdded"
	case StorageDeleted:
		return "StorageDeleted"
	default:
		panic("BUG: storage status not found")
	}
}
