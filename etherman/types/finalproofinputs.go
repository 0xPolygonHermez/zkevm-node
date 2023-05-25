package types

import "github.com/0xPolygonHermez/zkevm-node/aggregator/pb"

// FinalProofInputs struct
type FinalProofInputs struct {
	FinalProof       *pb.FinalProof
	NewLocalExitRoot []byte
	NewStateRoot     []byte
}
