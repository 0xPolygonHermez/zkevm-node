package state

import (
	"math/big"

	"github.com/0xPolygon/eth-state-transition/types"
	"github.com/ethereum/go-ethereum/common"
)

type SnapshotWriter interface {
	Snapshot

	Commit(objs []*Object) (SnapshotWriter, []byte)
}

type Snapshot interface {
	GetCode(hash common.Hash) ([]byte, bool)
	GetStorage(root common.Hash, key types.Hash) common.Hash
	GetAccount(addr common.Address) (*Account, error)
}

// Object is the serialization of the radix object (can be merged to StateObject?).
type Object struct {
	Address  common.Address
	CodeHash common.Hash
	Balance  *big.Int
	Root     common.Hash
	Nonce    uint64
	Deleted  bool

	// TODO: Move this to executor
	DirtyCode bool
	Code      []byte

	Storage []*StorageObject
}

// StorageObject is an entry in the storage
type StorageObject struct {
	Deleted bool
	Key     []byte
	Val     []byte
}

type Result struct {
	Logs            []*Log
	Success         bool
	GasUsed         uint64
	ContractAddress common.Address
	ReturnValue     []byte
}

type Log struct {
	Address common.Address
	Topics  []common.Hash
	Data    []byte
}
