package tree

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/hex"
)

// DefaultMerkleTreeArity specifies Merkle Tree arity used by default
const DefaultMerkleTreeArity = 1

var (
	// ErrDBTxsNotSupported indicates db transactions are not supported
	ErrDBTxsNotSupported = errors.New("transactions are not supported")
)

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

// SupportsDBTransactions indicates whether the store implementation supports DB transactions
func (tree *StateTree) SupportsDBTransactions() bool {
	return tree.mt.SupportsDBTransactions() && tree.scCodeStore.SupportsDBTransactions()
}

// BeginDBTransaction starts a transaction block
func (tree *StateTree) BeginDBTransaction(ctx context.Context, txBundleID string) error {
	if err := tree.mt.BeginDBTransaction(ctx, txBundleID); err != nil {
		return err
	}
	return tree.scCodeStore.BeginDBTransaction(ctx, txBundleID)
}

// Commit commits a db transaction
func (tree *StateTree) Commit(ctx context.Context, txBundleID string) error {
	if err := tree.mt.Commit(ctx, txBundleID); err != nil {
		return err
	}
	return tree.scCodeStore.Commit(ctx, txBundleID)
}

// Rollback rollbacks a db transaction
func (tree *StateTree) Rollback(ctx context.Context, txBundleID string) error {
	if err := tree.mt.Rollback(ctx, txBundleID); err != nil {
		return err
	}
	return tree.scCodeStore.Rollback(ctx, txBundleID)
}

// GetBalance returns balance
func (tree *StateTree) GetBalance(ctx context.Context, address common.Address, root []byte, txBundleID string) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyEthAddrBalance(address)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(ctx, scalarToh4(r), scalarToh4(k), txBundleID)
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
	proof, err := tree.mt.Get(ctx, scalarToh4(r), scalarToh4(k), txBundleID)
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
	proof, err := tree.mt.Get(ctx, scalarToh4(r), scalarToh4(k), txBundleID)
	if err != nil {
		return nil, err
	}
	if proof.Value == nil {
		return nil, nil
	}

	valueBi := fea2scalar(proof.Value)
	var valueBuff [maxBigIntLen]byte
	return valueBi.FillBytes(valueBuff[:]), nil
}

// GetCode returns code
func (tree *StateTree) GetCode(ctx context.Context, address common.Address, root []byte, txBundleID string) ([]byte, error) {
	scCodeHash, err := tree.GetCodeHash(ctx, address, root, txBundleID)
	if err != nil {
		return nil, err
	}

	// this code gets actual smart contract code from sc code storage
	scCode, err := tree.scCodeStore.Get(ctx, scCodeHash, txBundleID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return []byte{}, nil
		}
		return nil, err
	}

	return scCode, nil
}

// GetStorageAt returns Storage Value at specified position
func (tree *StateTree) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, root []byte, txBundleID string) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyContractStorage(address, position.Bytes())
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.mt.Get(ctx, scalarToh4(r), scalarToh4(k), txBundleID)
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return fea2scalar(proof.Value), nil
}

// ReverseHash reverse a hash of an exisiting Merkletree node.
func (tree *StateTree) ReverseHash(root, hash []byte) ([]byte, error) {
	hashBI := new(big.Int).SetBytes(hash[:])
	rootBI := new(big.Int).SetBytes(root[:])

	proof, err := tree.mt.Get(context.Background(), scalarToh4(rootBI), scalarToh4(hashBI), "")
	if err != nil {
		return nil, err
	}

	return fea2scalar(proof.Value).Bytes(), nil
}

// SetBalance sets balance
func (tree *StateTree) SetBalance(ctx context.Context, address common.Address, balance *big.Int, root []byte, txBundleID string) (newRoot []byte, proof *UpdateProof, err error) {
	if balance.Cmp(big.NewInt(0)) == -1 {
		return nil, nil, fmt.Errorf("invalid balance")
	}

	r := new(big.Int).SetBytes(root)
	key, err := KeyEthAddrBalance(address)
	if err != nil {
		return nil, nil, err
	}

	k := new(big.Int).SetBytes(key)
	balanceH8 := scalar2fea(balance)

	updateProof, err := tree.mt.Set(ctx, scalarToh4(r), scalarToh4(k), balanceH8, txBundleID)
	if err != nil {
		return nil, nil, err
	}

	rootBI := h4ToScalar(updateProof.NewRoot)
	return rootBI.Bytes(), updateProof, nil
}

// SetNonce sets nonce
func (tree *StateTree) SetNonce(ctx context.Context, address common.Address, nonce *big.Int, root []byte, txBundleID string) (newRoot []byte, proof *UpdateProof, err error) {
	if nonce.Cmp(big.NewInt(0)) == -1 {
		return nil, nil, fmt.Errorf("invalid nonce")
	}

	r := new(big.Int).SetBytes(root)
	key, err := KeyEthAddrNonce(address)
	if err != nil {
		return nil, nil, err
	}

	k := new(big.Int).SetBytes(key[:])

	nonceH8 := scalar2fea(nonce)

	updateProof, err := tree.mt.Set(ctx, scalarToh4(r), scalarToh4(k), nonceH8, txBundleID)
	if err != nil {
		return nil, nil, err
	}

	return h4ToScalar(updateProof.NewRoot).Bytes(), updateProof, nil
}

// SetCode sets smart contract code
func (tree *StateTree) SetCode(ctx context.Context, address common.Address, code []byte, root []byte, txBundleID string) (newRoot []byte, proof *UpdateProof, err error) {
	// calculating smart contract code hash
	scCodeHash4, err := tree.mt.scHashFunction(code)
	if err != nil {
		return nil, nil, err
	}

	scCodeHash, err := hex.DecodeHex(h4ToString(scCodeHash4))
	if err != nil {
		return nil, nil, err
	}

	// store smart contract code by its hash
	err = tree.scCodeStore.Set(ctx, scCodeHash[:], code, txBundleID)
	if err != nil {
		return nil, nil, err
	}

	// set smart contract code hash as a leaf value in merkle tree
	r := new(big.Int).SetBytes(root)
	key, err := KeyContractCode(address)
	if err != nil {
		return nil, nil, err
	}
	k := new(big.Int).SetBytes(key[:])
	scCodeHashBI := new(big.Int).SetBytes(scCodeHash[:])
	scCodeHashH8 := scalar2fea(scCodeHashBI)

	updateProof, err := tree.mt.Set(ctx, scalarToh4(r), scalarToh4(k), scCodeHashH8, txBundleID)
	if err != nil {
		return nil, nil, err
	}

	return h4ToScalar(updateProof.NewRoot).Bytes(), updateProof, nil
}

// SetStorageAt sets storage value at specified position
func (tree *StateTree) SetStorageAt(ctx context.Context, address common.Address, position *big.Int, value *big.Int, root []byte, txBundleID string) (newRoot []byte, proof *UpdateProof, err error) {
	r := new(big.Int).SetBytes(root)
	key, err := KeyContractStorage(address, position.Bytes())
	if err != nil {
		return nil, nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	valueH8 := scalar2fea(value)
	updateProof, err := tree.mt.Set(ctx, scalarToh4(r), scalarToh4(k), valueH8, txBundleID)
	if err != nil {
		return nil, nil, err
	}

	return h4ToScalar(updateProof.NewRoot).Bytes(), updateProof, nil
}

// SetHashValue sets value for an specific key.
func (tree *StateTree) SetHashValue(ctx context.Context, key common.Hash, value *big.Int, root []byte) (newRoot []byte, proof *UpdateProof, err error) {
	r := new(big.Int).SetBytes(root)
	k := new(big.Int).SetBytes(key[:])
	valueH8 := scalar2fea(value)
	updateProof, err := tree.mt.Set(ctx, scalarToh4(r), scalarToh4(k), valueH8, "")
	if err != nil {
		return nil, nil, err
	}

	return h4ToScalar(updateProof.NewRoot).Bytes(), updateProof, nil
}

// SetNodeData sets data for a specific node.
func (tree *StateTree) SetNodeData(ctx context.Context, key *big.Int, value *big.Int) (err error) {
	var k [maxBigIntLen]byte
	key.FillBytes(k[:])

	return tree.mt.store.Set(ctx, k[:], value.Bytes(), "")
}

// GetNodeData sets data for a specific node.
func (tree *StateTree) GetNodeData(ctx context.Context, key *big.Int) (*big.Int, error) {
	var k [maxBigIntLen]byte
	key.FillBytes(k[:])
	data, err := tree.mt.store.Get(ctx, k[:], "")
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetBytes(data), nil
}
