package aggregator

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

type proverInterface interface {
	ID() string
	IsIdle() bool
	BatchProof(input *pb.InputProver) (string, error)
	FinalProof(inputProof string) (string, error)
	WaitRecursiveProof(ctx context.Context, proofID string) (string, error)
	WaitFinalProof(ctx context.Context, proofID string) (*pb.FinalProof, error)
}

// ethTxManager contains the methods required to send txs to
// ethereum.
type ethTxManager interface {
	VerifyBatches(ctx context.Context, lastVerifiedBatch uint64, batchNum uint64, resGetProof *pb.FinalProof) (*types.Transaction, error)
}

// etherman contains the methods required to interact with ethereum
type etherman interface {
	GetLatestVerifiedBatchNum() (uint64, error)
	GetPublicAddress() (common.Address, error)
}

// aggregatorTxProfitabilityChecker interface for different profitability
// checking algorithms.
type aggregatorTxProfitabilityChecker interface {
	IsProfitable(context.Context, *big.Int) (bool, error)
}

// stateInterface gathers the methods to interact with the state.
type stateInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	CheckProofContainsCompleteSequences(ctx context.Context, proof *state.Proof, dbTx pgx.Tx) (bool, error)
	GetLastVerifiedBatch(ctx context.Context, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetVirtualBatchToProve(ctx context.Context, lastVerfiedBatchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetProofsToAggregate(ctx context.Context, dbTx pgx.Tx) (*state.Proof, *state.Proof, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	AddGeneratedProof(ctx context.Context, proof *state.Proof, dbTx pgx.Tx) error
	UpdateGeneratedProof(ctx context.Context, proof *state.Proof, dbTx pgx.Tx) error
	DeleteGeneratedProofs(ctx context.Context, batchNumber uint64, batchNumberFinal uint64, dbTx pgx.Tx) error
	DeleteUngeneratedProofs(ctx context.Context, dbTx pgx.Tx) error
}
