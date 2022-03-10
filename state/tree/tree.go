package tree

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// DefaultMerkleTreeArity specifies Merkle Tree arity used by default
const DefaultMerkleTreeArity = 4

// StateTree provides methods to access and modify state in merkletree
type StateTree struct {
	mt          *MerkleTree
	scCodeStore Store
}

// NewStateTree creates new StateTree
func NewStateTree(mt *MerkleTree, scCodeStore Store) *StateTree {
	return &StateTree{
		mt:          mt,
		scCodeStore: scCodeStore,
	}
}

// GetBalance returns balance
func (tree *StateTree) GetBalance(ctx context.Context, address common.Address, root []byte) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := GetKey(LeafTypeBalance, address, nil, tree.mt.arity, nil)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(ctx, r, k)
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return proof.Value, nil
}

// GetNonce returns nonce
func (tree *StateTree) GetNonce(ctx context.Context, address common.Address, root []byte) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := GetKey(LeafTypeNonce, address, nil, tree.mt.arity, nil)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(ctx, r, k)
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return proof.Value, nil
}

// GetCodeHash returns code hash
func (tree *StateTree) GetCodeHash(ctx context.Context, address common.Address, root []byte) ([]byte, error) {
	r := new(big.Int).SetBytes(root)

	key, err := GetKey(LeafTypeCode, address, nil, tree.mt.arity, nil)
	if err != nil {
		return nil, err
	}

	// this code gets only the hash of the smart contract code from the merkle tree
	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(ctx, r, k)
	if err != nil {
		return nil, err
	}
	if proof.Value == nil {
		return []byte{}, nil
	}

	var buf [32]byte
	return proof.Value.FillBytes(buf[:]), nil
}

// GetCode returns code
func (tree *StateTree) GetCode(ctx context.Context, address common.Address, root []byte) ([]byte, error) {
	scCodeHash, err := tree.GetCodeHash(ctx, address, root)
	if err != nil {
		return nil, err
	}

	// this code gets actual smart contract code from sc code storage
	scCode, err := tree.scCodeStore.Get(ctx, scCodeHash)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return []byte{}, nil
		}
		return nil, err
	}

	return scCode, nil
}

// GetStorageAt returns Storage Value at specified position
func (tree *StateTree) GetStorageAt(ctx context.Context, address common.Address, position common.Hash, root []byte) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := GetKey(LeafTypeStorage, address, position[:], tree.mt.arity, nil)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(ctx, r, k)
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return proof.Value, nil
}

// ReverseHash reverse a hash of an exisiting Merkletree node.
func (tree *StateTree) ReverseHash(root, hash []byte) ([]byte, error) {
	hashBI := new(big.Int).SetBytes(hash[:])
	rootBI := new(big.Int).SetBytes(root[:])

	proof, err := tree.mt.Get(context.Background(), rootBI, hashBI)
	if err != nil {
		return nil, err
	}

	return proof.Value.Bytes(), nil
}

// SetBalance sets balance
func (tree *StateTree) SetBalance(ctx context.Context, address common.Address, balance *big.Int, root []byte) (newRoot []byte, proof *UpdateProof, err error) {
	if balance.Cmp(big.NewInt(0)) == -1 {
		return nil, nil, fmt.Errorf("invalid balance")
	}

	r := new(big.Int).SetBytes(root)
	key, err := GetKey(LeafTypeBalance, address, nil, tree.mt.arity, nil)
	if err != nil {
		return nil, nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	updateProof, err := tree.mt.Set(ctx, r, k, balance)
	if err != nil {
		return nil, nil, err
	}

	return updateProof.NewRoot.Bytes(), updateProof, nil
}

// SetNonce sets nonce
func (tree *StateTree) SetNonce(ctx context.Context, address common.Address, nonce *big.Int, root []byte) (newRoot []byte, proof *UpdateProof, err error) {
	if nonce.Cmp(big.NewInt(0)) == -1 {
		return nil, nil, fmt.Errorf("invalid nonce")
	}

	r := new(big.Int).SetBytes(root)
	key, err := GetKey(LeafTypeNonce, address, nil, tree.mt.arity, nil)
	if err != nil {
		return nil, nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	updateProof, err := tree.mt.Set(ctx, r, k, nonce)
	if err != nil {
		return nil, nil, err
	}

	return updateProof.NewRoot.Bytes(), updateProof, nil
}

// SetCode sets smart contract code
func (tree *StateTree) SetCode(ctx context.Context, address common.Address, code []byte, root []byte) (newRoot []byte, proof *UpdateProof, err error) {
	// calculating smart contract code hash
	scCodeHashBI := big.NewInt(0)
	if len(code) > 0 {
		scCodeHashBI, err = tree.mt.scHashFunction(code)
		if err != nil {
			return nil, nil, err
		}
	}
	// we need to have exactly maxBigIntLen bytes for a key in db for
	// interoperability with prover/executor code, but big.Int Bytes()
	// can return less, so we make sure it has the right size with FillBytes.
	var scCodeHash [maxBigIntLen]byte
	scCodeHashBI.FillBytes(scCodeHash[:])

	// store smart contract code by its hash
	err = tree.scCodeStore.Set(ctx, scCodeHash[:], code)
	if err != nil {
		return nil, nil, err
	}

	// set smart contract code hash as a leaf value in merkle tree
	r := new(big.Int).SetBytes(root)
	key, err := GetKey(LeafTypeCode, address, nil, tree.mt.arity, nil)
	if err != nil {
		return nil, nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	updateProof, err := tree.mt.Set(ctx, r, k, new(big.Int).SetBytes(scCodeHash[:]))
	if err != nil {
		return nil, nil, err
	}

	return updateProof.NewRoot.Bytes(), updateProof, nil
}

// SetStorageAt sets storage value at specified position
func (tree *StateTree) SetStorageAt(ctx context.Context, address common.Address, position *big.Int, value *big.Int, root []byte) (newRoot []byte, proof *UpdateProof, err error) {
	r := new(big.Int).SetBytes(root)
	key, err := GetKey(LeafTypeStorage, address, position.Bytes(), tree.mt.arity, nil)
	if err != nil {
		return nil, nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	updateProof, err := tree.mt.Set(ctx, r, k, value)
	if err != nil {
		return nil, nil, err
	}

	return updateProof.NewRoot.Bytes(), updateProof, nil
}

// SetHashValue sets value for an specific key.
func (tree *StateTree) SetHashValue(ctx context.Context, key common.Hash, value *big.Int, root []byte) (newRoot []byte, proof *UpdateProof, err error) {
	r := new(big.Int).SetBytes(root)
	k := new(big.Int).SetBytes(key[:])
	updateProof, err := tree.mt.Set(ctx, r, k, value)
	if err != nil {
		return nil, nil, err
	}

	return updateProof.NewRoot.Bytes(), updateProof, nil
}

// SetNodeData sets data for a specific node.
func (tree *StateTree) SetNodeData(ctx context.Context, key *big.Int, value *big.Int) (err error) {
	var k [maxBigIntLen]byte
	key.FillBytes(k[:])

	return tree.mt.store.Set(ctx, k[:], value.Bytes())
}

// GetNodeData sets data for a specific node.
func (tree *StateTree) GetNodeData(ctx context.Context, key *big.Int) (*big.Int, error) {
	var k [maxBigIntLen]byte
	key.FillBytes(k[:])
	data, err := tree.mt.store.Get(ctx, k[:])
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetBytes(data), nil
}
