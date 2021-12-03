package tree

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4/pgxpool"
)

const mtArity = 4

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
func NewReader(db *pgxpool.Pool) (Reader, error) {
	// TODO: switch to state tree backed by db
	return NewBasicTree(db), nil
}

// NewReadWriter returns object of ReadWriter interface
func NewReadWriter(db *pgxpool.Pool) (ReadWriter, error) {
	// TODO: switch to state tree backed by db
	return NewBasicTree(db), nil
}

// BasicTree is a basic in-memory implementation of StateTree
type BasicTree struct {
	mt          *MerkleTree
	currentRoot *big.Int
}

// NewBasicTree creates new BasicTree
func NewBasicTree(db *pgxpool.Pool) *BasicTree {
	mt := NewMerkleTree(db, mtArity, nil)
	return &BasicTree{mt: mt}
}

// GetBalance returns balance
func (tree *BasicTree) GetBalance(address common.Address, root []byte) (*big.Int, error) {
	r := tree.currentRoot
	if root != nil {
		r = new(big.Int).SetBytes(root)
	}
	key, err := GetKey(LeafTypeBalance, address, nil)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(context.TODO(), r, k)
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return proof.Value, nil
}

// GetNonce returns nonce
func (tree *BasicTree) GetNonce(address common.Address, root []byte) (*big.Int, error) {
	r := tree.currentRoot
	if root != nil {
		r = new(big.Int).SetBytes(root)
	}
	key, err := GetKey(LeafTypeNonce, address, nil)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(context.TODO(), r, k)
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return proof.Value, nil
}

// GetCode returns code
func (tree *BasicTree) GetCode(address common.Address, root []byte) ([]byte, error) {
	r := tree.currentRoot
	if root != nil {
		r = new(big.Int).SetBytes(root)
	}
	key, err := GetKey(LeafTypeCode, address, nil)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(context.TODO(), r, k)
	if err != nil {
		return nil, err
	}
	if proof.Value == nil {
		return []byte{}, nil
	}
	// TODO: fix how code is returned from mt
	return proof.Value.Bytes(), nil
}

// GetStorageAt returns Storage Value at specified position
func (tree *BasicTree) GetStorageAt(address common.Address, position common.Hash, root []byte) (*big.Int, error) {
	r := tree.currentRoot
	if root != nil {
		r = new(big.Int).SetBytes(root)
	}
	key, err := GetKey(LeafTypeStorage, address, position[:])
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(context.TODO(), r, k)
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return proof.Value, nil
}

// GetRoot returns current MerkleTree root hash
func (tree *BasicTree) GetRoot() ([]byte, error) {
	return tree.currentRoot.Bytes(), nil
	//return nil, fmt.Errorf("not implemented")
}

// GetRootForBatchNumber returns MerkleTree root for specified batchNumber
func (tree *BasicTree) GetRootForBatchNumber(batchNumber uint64) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

// SetBalance sets balance
func (tree *BasicTree) SetBalance(address common.Address, balance *big.Int) (newRoot []byte, proof interface{}, err error) {
	r := tree.currentRoot
	key, err := GetKey(LeafTypeBalance, address, nil)
	if err != nil {
		return nil, nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	updateProof, err := tree.mt.Set(context.TODO(), r, k, balance)
	if err != nil {
		return nil, nil, err
	}

	tree.currentRoot = updateProof.NewRoot

	return updateProof.NewRoot.Bytes(), updateProof, nil
}

// SetNonce sets nonce
func (tree *BasicTree) SetNonce(address common.Address, nonce *big.Int) (newRoot []byte, proof interface{}, err error) {
	r := tree.currentRoot
	key, err := GetKey(LeafTypeNonce, address, nil)
	if err != nil {
		return nil, nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	updateProof, err := tree.mt.Set(context.TODO(), r, k, nonce)
	if err != nil {
		return nil, nil, err
	}

	tree.currentRoot = updateProof.NewRoot

	return updateProof.NewRoot.Bytes(), updateProof, nil
}

// SetCode sets code
func (tree *BasicTree) SetCode(address common.Address, code []byte) (newRoot []byte, proof interface{}, err error) {
	return nil, nil, fmt.Errorf("not implemented")
}

// SetStorageAt sets storage value at specified position
func (tree *BasicTree) SetStorageAt(address common.Address, position common.Hash, value *big.Int) (newRoot []byte, proof interface{}, err error) {
	r := tree.currentRoot
	key, err := GetKey(LeafTypeStorage, address, position[:])
	if err != nil {
		return nil, nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	updateProof, err := tree.mt.Set(context.TODO(), r, k, value)
	if err != nil {
		return nil, nil, err
	}

	tree.currentRoot = updateProof.NewRoot

	return updateProof.NewRoot.Bytes(), updateProof, nil
}

// SetRootForBatchNumber sets root for specified batchNumber
func (tree *BasicTree) SetRootForBatchNumber(batchNumber uint64, root []byte) error {
	return fmt.Errorf("not implemented")
}
