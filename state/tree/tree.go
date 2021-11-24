package tree

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Reader interface
type Reader interface {
	GetBalance(address common.Address, root []byte) (*big.Int, error)
	GetNonce(address common.Address, root []byte) (*big.Int, error)
	GetCode(address common.Address, root []byte) ([]byte, error)
	GetStorageAt(address common.Address, position common.Hash, root []byte) (*big.Int, error)
	GetRoot() ([]byte, error)
	GetRootForBatchNumber(batchNumber uint64) ([]byte, error)
}

// Writer interface
type Writer interface {
	SetBalance(address common.Address, balance *big.Int) (newRoot []byte, proof interface{}, err error)
	SetNonce(address common.Address, nonce *big.Int) (newRoot []byte, proof interface{}, err error)
	SetCode(address common.Address, code []byte) (newRoot []byte, proof interface{}, err error)
	SetStorageAt(address common.Address, key common.Hash, value *big.Int) (newRoot []byte, proof interface{}, err error)
	SetRootForBatchNumber(batchNumber uint64, root []byte) error
}

// ReadWriter interface
type ReadWriter interface {
	Reader
	Writer
}

// NewReader returns object of Reader interface
func NewReader(db interface{}) (Reader, error) {
	// TODO: switch to state tree backed by db
	return NewMemTree(), nil
}

// NewReadWriter returns object of ReadWriter interface
func NewReadWriter(db interface{}) (ReadWriter, error) {
	// TODO: switch to state tree backed by db
	return NewMemTree(), nil
}
