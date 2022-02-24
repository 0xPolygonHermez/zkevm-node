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
	ctx context.Context

	grpcClient pb.MTServiceClient
}

// NewAdapter is the constructor of Adapter.
func NewAdapter(ctx context.Context, client pb.MTServiceClient) *Adapter {
	return &Adapter{
		ctx:        ctx,
		grpcClient: client,
	}
}

// GetBalance returns balance.
func (m *Adapter) GetBalance(address common.Address, root []byte) (*big.Int, error) {
	result, err := m.grpcClient.GetBalance(m.ctx, &pb.GetBalanceRequest{
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
func (m *Adapter) GetNonce(address common.Address, root []byte) (*big.Int, error) {
	result, err := m.grpcClient.GetNonce(m.ctx, &pb.GetNonceRequest{
		EthAddress: address.String(),
		Root:       hex.EncodeToString(root),
	})
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetUint64(result.Nonce), nil
}

// GetCode returns code.
func (m *Adapter) GetCode(address common.Address, root []byte) ([]byte, error) {
	result, err := m.grpcClient.GetCode(m.ctx, &pb.GetCodeRequest{
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
func (m *Adapter) GetCodeHash(address common.Address, root []byte) ([]byte, error) {
	result, err := m.grpcClient.GetCodeHash(m.ctx, &pb.GetCodeHashRequest{
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
func (m *Adapter) GetStorageAt(address common.Address, position common.Hash, root []byte) (*big.Int, error) {
	result, err := m.grpcClient.GetStorageAt(m.ctx, &pb.GetStorageAtRequest{
		EthAddress: address.String(),
		Root:       hex.EncodeToString(root),
		Position:   position.Big().Uint64(),
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
func (m *Adapter) SetBalance(address common.Address, balance *big.Int, root []byte) (newRoot []byte, proof *UpdateProof, err error) {
	result, err := m.grpcClient.SetBalance(m.ctx, &pb.SetBalanceRequest{
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

	if result.Data == nil {
		return nil, nil, fmt.Errorf("No data returned in gRPC call SetBalance")
	}

	newRoot, err = hex.DecodeString(result.Data.NewRoot)
	if err != nil {
		return nil, nil, err
	}

	return newRoot, proof, nil
}

// SetNonce sets nonce.
func (m *Adapter) SetNonce(address common.Address, nonce *big.Int, root []byte) (newRoot []byte, proof *UpdateProof, err error) {
	result, err := m.grpcClient.SetNonce(m.ctx, &pb.SetNonceRequest{
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

	if result.Data == nil {
		return nil, nil, fmt.Errorf("No data returned in gRPC call SetNonce")
	}

	newRoot, err = hex.DecodeString(result.Data.NewRoot)
	if err != nil {
		return nil, nil, err
	}

	return newRoot, proof, nil
}

// SetCode sets smart contract code.
func (m *Adapter) SetCode(address common.Address, code []byte, root []byte) (newRoot []byte, proof *UpdateProof, err error) {
	result, err := m.grpcClient.SetCode(m.ctx, &pb.SetCodeRequest{
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

	if result.Data == nil {
		return nil, nil, fmt.Errorf("No data returned in gRPC call SetCode")
	}

	newRoot, err = hex.DecodeString(result.Data.NewRoot)
	if err != nil {
		return nil, nil, err
	}

	return newRoot, proof, nil
}

// SetStorageAt sets storage value at specified position.
func (m *Adapter) SetStorageAt(address common.Address, position *big.Int, value *big.Int, root []byte) (newRoot []byte, proof *UpdateProof, err error) {
	result, err := m.grpcClient.SetStorageAt(m.ctx, &pb.SetStorageAtRequest{
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

	if result.Data == nil {
		return nil, nil, fmt.Errorf("No data returned in gRPC call SetStorageAt")
	}

	newRoot, err = hex.DecodeString(result.Data.NewRoot)
	if err != nil {
		return nil, nil, err
	}

	return newRoot, proof, nil
}
