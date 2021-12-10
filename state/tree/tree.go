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
	GetCurrentRoot() ([]byte, error)
}

// Writer interface
type Writer interface {
	SetBalance(address common.Address, balance *big.Int) (newRoot []byte, proof interface{}, err error)
	SetNonce(address common.Address, nonce *big.Int) (newRoot []byte, proof interface{}, err error)
	SetCode(address common.Address, code []byte) (newRoot []byte, proof interface{}, err error)
	SetStorageAt(address common.Address, key common.Hash, value *big.Int) (newRoot []byte, proof interface{}, err error)
	SetCurrentRoot([]byte)
}

// ReadWriter interface
type ReadWriter interface {
	Reader
	Writer
}

// StateTree provides methods to access and modify state in merkletree
type StateTree struct {
	db          *pgxpool.Pool
	mt          *MerkleTree
	currentRoot *big.Int
}

// NewStateTree creates new StateTree
func NewStateTree(db *pgxpool.Pool, root []byte) ReadWriter {
	mt := NewMerkleTree(db, mtArity, nil)
	return &StateTree{
		db:          db,
		mt:          mt,
		currentRoot: new(big.Int).SetBytes(root),
	}
}

// GetBalance returns balance
func (tree *StateTree) GetBalance(address common.Address, root []byte) (*big.Int, error) {
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
func (tree *StateTree) GetNonce(address common.Address, root []byte) (*big.Int, error) {
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
func (tree *StateTree) GetCode(address common.Address, root []byte) ([]byte, error) {
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
func (tree *StateTree) GetStorageAt(address common.Address, position common.Hash, root []byte) (*big.Int, error) {
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

// GetCurrentRoot returns current MerkleTree root hash
func (tree *StateTree) GetCurrentRoot() ([]byte, error) {
	return tree.currentRoot.Bytes(), nil
}

// SetBalance sets balance
func (tree *StateTree) SetBalance(address common.Address, balance *big.Int) (newRoot []byte, proof interface{}, err error) {
	if balance.Cmp(big.NewInt(0)) == -1 {
		return nil, nil, fmt.Errorf("invalid balance")
	}

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
func (tree *StateTree) SetNonce(address common.Address, nonce *big.Int) (newRoot []byte, proof interface{}, err error) {
	if nonce.Cmp(big.NewInt(0)) == -1 {
		return nil, nil, fmt.Errorf("invalid nonce")
	}

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
func (tree *StateTree) SetCode(address common.Address, code []byte) (newRoot []byte, proof interface{}, err error) {
	return nil, nil, fmt.Errorf("not implemented")
}

// SetStorageAt sets storage value at specified position
func (tree *StateTree) SetStorageAt(address common.Address, position common.Hash, value *big.Int) (newRoot []byte, proof interface{}, err error) {
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

// SetCurrentRoot sets current root of the state tree
func (tree *StateTree) SetCurrentRoot(root []byte) {
	tree.currentRoot = new(big.Int).SetBytes(root)
}
