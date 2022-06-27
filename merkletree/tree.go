package merkletree

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/merkletree/pb"
)

// StateTree provides methods to access and modify state in merkletree
type StateTree struct {
	grpcClient pb.StateDBServiceClient
}

// NewStateTree creates new StateTree.
func NewStateTree(client pb.StateDBServiceClient) *StateTree {
	return &StateTree{
		grpcClient: client,
	}
}

// GetBalance returns balance
func (tree *StateTree) GetBalance(ctx context.Context, address common.Address, root []byte) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyEthAddrBalance(address)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.get(ctx, scalarToh4(r), scalarToh4(k))
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return fea2scalar(proof.Value), nil
}

// GetNonce returns nonce
func (tree *StateTree) GetNonce(ctx context.Context, address common.Address, root []byte) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyEthAddrNonce(address)
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.get(ctx, scalarToh4(r), scalarToh4(k))
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return fea2scalar(proof.Value), nil
}

// GetCodeHash returns code hash
func (tree *StateTree) GetCodeHash(ctx context.Context, address common.Address, root []byte) ([]byte, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyContractCode(address)
	if err != nil {
		return nil, err
	}
	// this code gets only the hash of the smart contract code from the merkle tree
	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.get(ctx, scalarToh4(r), scalarToh4(k))
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
func (tree *StateTree) GetCode(ctx context.Context, address common.Address, root []byte) ([]byte, error) {
	scCodeHash, err := tree.GetCodeHash(ctx, address, root)
	if err != nil {
		return nil, err
	}

	// this code gets actual smart contract code from sc code storage
	scCode, err := tree.getProgram(ctx, common.Bytes2Hex(scCodeHash))
	if err != nil {
		return nil, err
	}

	return scCode.Data, nil
}

// GetStorageAt returns Storage Value at specified position
func (tree *StateTree) GetStorageAt(ctx context.Context, address common.Address, position *big.Int, root []byte) (*big.Int, error) {
	r := new(big.Int).SetBytes(root)

	key, err := KeyContractStorage(address, position.Bytes())
	if err != nil {
		return nil, err
	}

	k := new(big.Int).SetBytes(key[:])
	proof, err := tree.get(ctx, scalarToh4(r), scalarToh4(k))
	if err != nil {
		return nil, err
	}
	if proof == nil || proof.Value == nil {
		return big.NewInt(0), nil
	}
	return fea2scalar(proof.Value), nil
}

func (tree *StateTree) get(ctx context.Context, root, key []uint64) (*Proof, error) {
	result, err := tree.grpcClient.Get(ctx, &pb.GetRequest{
		Root:    &pb.Fea{Fe0: root[0], Fe1: root[1], Fe2: root[2], Fe3: root[3]},
		Key:     &pb.Fea{Fe0: key[0], Fe1: key[1], Fe2: key[2], Fe3: key[3]},
		Details: false,
	})
	if err != nil {
		return nil, err
	}

	value, err := stringToh4(result.Value)
	if err != nil {
		return nil, err
	}
	return &Proof{
		Root:  []uint64{result.Root.Fe0, result.Root.Fe1, result.Root.Fe2, result.Root.Fe3},
		Key:   key,
		Value: value,
	}, nil
}

func (tree *StateTree) getProgram(ctx context.Context, hash string) (*ProgramProof, error) {
	result, err := tree.grpcClient.GetProgram(ctx, &pb.GetProgramRequest{
		Hash: hash,
	})
	if err != nil {
		return nil, err
	}

	return &ProgramProof{
		Data: result.Data,
	}, nil
}
