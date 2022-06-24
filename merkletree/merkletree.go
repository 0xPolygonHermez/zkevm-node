package merkletree

import (
	"context"

	"github.com/hermeznetwork/hermez-core/merkletree/pb"
)

// MerkleTree exposes the MT methods required by the consumers and translates them into
// gRPC calls using its client member.
type MerkleTree struct {
	grpcClient pb.StateDBServiceClient
}

// New is the constructor of MerkleTree.
func New(client pb.StateDBServiceClient) *MerkleTree {
	return &MerkleTree{
		grpcClient: client,
	}
}

func (m *MerkleTree) Get(ctx context.Context, root, key []uint64) (*Proof, error) {
	result, err := m.grpcClient.Get(ctx, &pb.GetRequest{
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

func (m *MerkleTree) GetProgram(ctx context.Context, hash string) (*ProgramProof, error) {
	result, err := m.grpcClient.GetProgram(ctx, &pb.GetProgramRequest{
		Hash: hash,
	})
	if err != nil {
		return nil, err
	}

	return &ProgramProof{
		Data: result.Data,
	}, nil
}
