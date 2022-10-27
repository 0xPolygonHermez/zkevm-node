package state

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ProofStatusPending represents a proof that has not been mined yet on L1
	ProofStatusPending ProofStatus = "pending"
	// ProofStatusConfirmed represents a proof that has been mined and the state is now verified
	ProofStatusConfirmed ProofStatus = "confirmed"
)

// ProofStatus represents the state of a tx
type ProofStatus string

// Proof struct
type Proof struct {
	BatchNumber uint64
	Proof       *pb.GetProofResponse
	InputProver *pb.InputProver
	ProofID     *string
	Prover      *string

	TxHash    *common.Hash
	TxNonce   *uint64
	CreatedAt time.Time
	UpdatedAt *time.Time
	Status    ProofStatus
}
