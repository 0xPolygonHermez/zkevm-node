package state

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/state/tree"
)

// Consumer interfaces required by the package.

// statetree contains the methods required to interact with the Merkle tree.
type statetree interface {
	SupportsDBTransactions() bool
	BeginDBTransaction(ctx context.Context, txBundleID string) error
	Commit(ctx context.Context, txBundleID string) error
	Rollback(ctx context.Context, txBundleID string) error
	GetBalance(ctx context.Context, address common.Address, root []byte, txBundleID string) (*big.Int, error)
	GetNonce(ctx context.Context, address common.Address, root []byte, txBundleID string) (*big.Int, error)
	GetCode(ctx context.Context, address common.Address, root []byte, txBundleID string) ([]byte, error)
	GetCodeHash(ctx context.Context, address common.Address, root []byte, txBundleID string) ([]byte, error)
	GetStorageAt(ctx context.Context, address common.Address, position *big.Int, root []byte, txBundleID string) (*big.Int, error)

	SetBalance(ctx context.Context, address common.Address, balance *big.Int, root []byte, txBundleID string) (newRoot []byte, proof *tree.UpdateProof, err error)
	SetNonce(ctx context.Context, address common.Address, nonce *big.Int, root []byte, txBundleID string) (newRoot []byte, proof *tree.UpdateProof, err error)
	SetCode(ctx context.Context, address common.Address, code []byte, root []byte, txBundleID string) (newRoot []byte, proof *tree.UpdateProof, err error)
	SetStorageAt(ctx context.Context, address common.Address, position *big.Int, value *big.Int, root []byte, txBundleID string) (newRoot []byte, proof *tree.UpdateProof, err error)
}
