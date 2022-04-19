package tree

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/state/tree/pb"
)

// Adapter exposes the MT methods required by the state and translates them into
// gRPC calls using its client member.
type Adapter struct {
	grpcClient pb.MTServiceClient
}

// NewAdapter is the constructor of Adapter.
func NewAdapter(client pb.MTServiceClient) *Adapter {
	return &Adapter{
		grpcClient: client,
	}
}

// SupportsDBTransactions indicates whether the store implementation supports DB transactions
func (m *Adapter) SupportsDBTransactions() bool {
	return false
}

// BeginDBTransaction starts a transaction block
func (m *Adapter) BeginDBTransaction(ctx context.Context, txBundleID string) error {
	return ErrDBTxsNotSupported
}

// Commit commits a db transaction
func (m *Adapter) Commit(ctx context.Context, txBundleID string) error {
	return ErrDBTxsNotSupported
}

// Rollback rollbacks a db transaction
func (m *Adapter) Rollback(ctx context.Context, txBundleID string) error {
	return ErrDBTxsNotSupported
}

// GetBalance returns balance.
func (m *Adapter) GetBalance(ctx context.Context, address common.Address, root []byte, txBundleID string) (*big.Int, error) {
	result, err := m.grpcClient.GetBalance(ctx, &pb.CommonGetRequest{
		EthAddress: address.String(),
		Root:       hex.EncodeToString(root),
	})
	if err != nil {
		return nil, err
	}

	balance, ok := new(big.Int).SetString(result.Balance, 10)
	if !ok {
		return nil, fmt.Errorf("Could not initialize balance from %q", result.Balance)
	}

	return balance, nil
}

// GetNonce returns nonce.
func (m *Adapter) GetNonce(ctx context.Context, address common.Address, root []byte, txBundleID string) (*big.Int, error) {
	result, err := m.grpcClient.GetNonce(ctx, &pb.CommonGetRequest{
		EthAddress: address.String(),
		Root:       hex.EncodeToString(root),
	})
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetUint64(result.Nonce), nil
}

// GetCode returns code.
func (m *Adapter) GetCode(ctx context.Context, address common.Address, root []byte, txBundleID string) ([]byte, error) {
	result, err := m.grpcClient.GetCode(ctx, &pb.CommonGetRequest{
		EthAddress: address.String(),
		Root:       hex.EncodeToString(root),
	})
	if err != nil {
		return nil, err
	}

	code, err := hex.DecodeString(result.Code)
	if err != nil {
		return nil, err
	}

	return code, nil
}

// GetCodeHash returns code hash.
func (m *Adapter) GetCodeHash(ctx context.Context, address common.Address, root []byte, txBundleID string) ([]byte, error) {
	result, err := m.grpcClient.GetCodeHash(ctx, &pb.CommonGetRequest{
		EthAddress: address.String(),
		Root:       hex.EncodeToString(root),
	})
	if err != nil {
		return nil, err
	}

	hash, err := hex.DecodeString(result.Hash)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

// GetStorageAt returns Storage Value at specified position.
func (m *Adapter) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, root []byte, txBundleID string) (*big.Int, error) {
	result, err := m.grpcClient.GetStorageAt(ctx, &pb.GetStorageAtRequest{
		EthAddress: address.String(),
		Root:       hex.EncodeToString(root),
		Position:   position.Uint64(),
	})
	if err != nil {
		return nil, err
	}

	value, ok := new(big.Int).SetString(result.Value, 10)
	if !ok {
		return nil, fmt.Errorf("Could not initialize storage value from %q", result.Value)
	}

	return value, nil
}

// SetBalance sets balance.
func (m *Adapter) SetBalance(ctx context.Context, address common.Address, balance *big.Int, root []byte, txBundleID string) (newRoot []byte, proof *UpdateProof, err error) {
	result, err := m.grpcClient.SetBalance(ctx, &pb.SetBalanceRequest{
		EthAddress: address.String(),
		Balance:    balance.String(),
		Root:       hex.EncodeToString(root),
	})

	if err != nil {
		return nil, nil, err
	}

	if !result.Success {
		return nil, nil, fmt.Errorf("Could not set balance %d for address %q", balance, address.String())
	}

	newRoot, err = hex.DecodeString(result.NewRoot)
	if err != nil {
		return nil, nil, err
	}

	return newRoot, proof, nil
}

// SetNonce sets nonce.
func (m *Adapter) SetNonce(ctx context.Context, address common.Address, nonce *big.Int, root []byte, txBundleID string) (newRoot []byte, proof *UpdateProof, err error) {
	result, err := m.grpcClient.SetNonce(ctx, &pb.SetNonceRequest{
		EthAddress: address.String(),
		Nonce:      nonce.Uint64(),
		Root:       hex.EncodeToString(root),
	})

	if err != nil {
		return nil, nil, err
	}

	if !result.Success {
		return nil, nil, fmt.Errorf("Could not set nonce %d for address %q", nonce, address.String())
	}

	newRoot, err = hex.DecodeString(result.NewRoot)
	if err != nil {
		return nil, nil, err
	}

	return newRoot, proof, nil
}

// SetCode sets smart contract code.
func (m *Adapter) SetCode(ctx context.Context, address common.Address, code []byte, root []byte, txBundleID string) (newRoot []byte, proof *UpdateProof, err error) {
	result, err := m.grpcClient.SetCode(ctx, &pb.SetCodeRequest{
		EthAddress: address.String(),
		Code:       hex.EncodeToString(code),
		Root:       hex.EncodeToString(root),
	})

	if err != nil {
		return nil, nil, err
	}

	if !result.Success {
		return nil, nil, fmt.Errorf("Could not set code %q for address %q", code, address.String())
	}

	newRoot, err = hex.DecodeString(result.NewRoot)
	if err != nil {
		return nil, nil, err
	}

	return newRoot, proof, nil
}

// SetStorageAt sets storage value at specified position.
func (m *Adapter) SetStorageAt(ctx context.Context, address common.Address, position *big.Int, value *big.Int, root []byte, txBundleID string) (newRoot []byte, proof *UpdateProof, err error) {
	result, err := m.grpcClient.SetStorageAt(ctx, &pb.SetStorageAtRequest{
		EthAddress: address.String(),
		Position:   position.String(),
		Value:      value.String(),
		Root:       hex.EncodeToString(root),
	})

	if err != nil {
		return nil, nil, err
	}

	if !result.Success {
		return nil, nil, fmt.Errorf("Could not set storage at position %q for address %q and value %d", position.String(), address.String(), value)
	}

	newRoot, err = hex.DecodeString(result.NewRoot)
	if err != nil {
		return nil, nil, err
	}

	return newRoot, proof, nil
}
