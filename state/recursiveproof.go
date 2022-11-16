package state

import "github.com/0xPolygonHermez/zkevm-node/aggregator2/pb"

// Proof struct
type RecursiveProof struct {
	BatchNumber      uint64
	BatchNumberFinal uint64
	Proof            *pb.GetProofResponse_RecursiveProof
	InputProver      string
	ProofID          *string
	Prover           *string
	Generating       bool
}
