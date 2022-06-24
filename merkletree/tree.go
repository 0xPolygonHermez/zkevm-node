package merkletree

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// StateTree provides methods to access and modify state in merkletree
type StateTree struct {
	mt *MerkleTree
}

// NewStateTree creates new StateTree.
func NewStateTree(mt *MerkleTree) *StateTree {
	return &StateTree{
		mt: mt,
	}
}

// GetBalance returns balance
func (tree *StateTree) GetBalance(ctx context.Context, address common.Address, root []byte, txBundleID string) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyEthAddrBalance(address)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(ctx, scalarToh4(r), scalarToh4(k))
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return fea2scalar(proof.Value), nil
}

// GetNonce returns nonce
func (tree *StateTree) GetNonce(ctx context.Context, address common.Address, root []byte, txBundleID string) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyEthAddrNonce(address)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(ctx, scalarToh4(r), scalarToh4(k))
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return fea2scalar(proof.Value), nil
}

// GetCodeHash returns code hash
func (tree *StateTree) GetCodeHash(ctx context.Context, address common.Address, root []byte, txBundleID string) ([]byte, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyContractCode(address)
	if err != nil {
		return nil, err
	}
	// this code gets only the hash of the smart contract code from the merkle tree
	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(ctx, scalarToh4(r), scalarToh4(k))
	if err != nil {
		return nil, err
	}
	if proof.Value == nil {
		return nil, nil
	}

	valueBi := fea2scalar(proof.Value)
	return ScalarToFilledByteSlice(valueBi), nil
}

// GetCode returns code
func (tree *StateTree) GetCode(ctx context.Context, address common.Address, root []byte, txBundleID string) ([]byte, error) {
	scCodeHash, err := tree.GetCodeHash(ctx, address, root, txBundleID)
	if err != nil {
		return nil, err
	}

	// this code gets actual smart contract code from sc code storage
	scCode, err := tree.mt.GetProgram(ctx, common.Bytes2Hex(scCodeHash))
	if err != nil {
		return nil, err
	}

	return scCode.Data, nil
}

// GetStorageAt returns Storage Value at specified position
func (tree *StateTree) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, root []byte, txBundleID string) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyContractStorage(address, position.Bytes())
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(ctx, scalarToh4(r), scalarToh4(k))
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return fea2scalar(proof.Value), nil
}
