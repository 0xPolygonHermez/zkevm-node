package state

import "time"

// Proof struct
type Proof struct {
	BatchNumber      uint64
	BatchNumberFinal uint64
	Proof            string
	InputProver      string
	ProofID          *string
	Prover           *string
	Generating       bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
