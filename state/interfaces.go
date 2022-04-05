package state

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/state/tree"
)

// Consumer interfaces required by the package.

// merkletree contains the methods required to interact with the Merkle tree.
type merkletree interface {
	SupportsDBTransactions() bool
	BeginDBTransaction(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	GetBalance(ctx context.Context, address common.Address, root []byte) (*big.Int, error)
	GetNonce(ctx context.Context, address common.Address, root []byte) (*big.Int, error)
	GetCode(ctx context.Context, address common.Address, root []byte) ([]byte, error)
	GetCodeHash(ctx context.Context, address common.Address, root []byte) ([]byte, error)
	GetStorageAt(ctx context.Context, address common.Address, position *big.Int, root []byte) (*big.Int, error)

	SetBalance(ctx context.Context, address common.Address, balance *big.Int, root []byte) (newRoot []byte, proof *tree.UpdateProof, err error)
	SetNonce(ctx context.Context, address common.Address, nonce *big.Int, root []byte) (newRoot []byte, proof *tree.UpdateProof, err error)
	SetCode(ctx context.Context, address common.Address, code []byte, root []byte) (newRoot []byte, proof *tree.UpdateProof, err error)
	SetStorageAt(ctx context.Context, address common.Address, key *big.Int, value *big.Int, root []byte) (newRoot []byte, proof *tree.UpdateProof, err error)
}
