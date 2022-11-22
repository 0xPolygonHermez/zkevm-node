package state

import "github.com/0xPolygonHermez/zkevm-node/aggregator/pb"

// Proof struct
type Proof struct {
	BatchNumber      uint64
	BatchNumberFinal uint64
	Proof            *pb.GetProofResponse_RecursiveProof
	InputProver      string
	ProofID          *string
	Prover           *string
	Generating       bool
}
