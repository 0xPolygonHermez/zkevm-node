package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/state/tree"
)

// Consumer interfaces required by the package.

// merkletree contains the methods required to interact with the Merkle tree.
type merkletree interface {
	GetBalance(address common.Address, root []byte) (*big.Int, error)
	GetNonce(address common.Address, root []byte) (*big.Int, error)
	GetCode(address common.Address, root []byte) ([]byte, error)
	GetCodeHash(address common.Address, root []byte) ([]byte, error)
	GetStorageAt(address common.Address, position common.Hash, root []byte) (*big.Int, error)
	GetCurrentRoot() ([]byte, error)

	SetBalance(address common.Address, balance *big.Int) (newRoot []byte, proof *tree.UpdateProof, err error)
	SetNonce(address common.Address, nonce *big.Int) (newRoot []byte, proof *tree.UpdateProof, err error)
	SetCode(address common.Address, code []byte) (newRoot []byte, proof *tree.UpdateProof, err error)
	SetStorageAt(address common.Address, key common.Hash, value *big.Int) (newRoot []byte, proof *tree.UpdateProof, err error)
	SetCurrentRoot([]byte)
}
