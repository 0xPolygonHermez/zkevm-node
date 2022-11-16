package aggregator2

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/aggregator2/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

type proverInterface interface {
	ID() string
	IsIdle() bool
	ProveBatch(input *pb.InputProver) (string, error)
	FinalProof(inputProof string) (string, error)
	WaitRecursiveProof(ctx context.Context, proofID string) (string, error)
	WaitFinalProof(ctx context.Context, proofID string) (*pb.FinalProof, error)
}

// ethTxManager contains the methods required to send txs to
// ethereum.
type ethTxManager interface {
}

// etherman contains the methods required to interact with ethereum
type etherman interface {
	VerifyBatches2(ctx context.Context, lastVerifiedBatch, newVerifiedBatch uint64, resGetProof *pb.GetProofResponse_FinalProof, gasLimit uint64, gasPrice, nonce *big.Int) (*types.Transaction, error)
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
	CheckProofContainsCompleteSequences(ctx context.Context, proof *state.RecursiveProof, dbTx pgx.Tx) (bool, error)
	GetLastVerifiedBatch(ctx context.Context, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetVirtualBatchToRecursiveProve(ctx context.Context, lastVerfiedBatchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetRecursiveProofsToAggregate(ctx context.Context, dbTx pgx.Tx) (*state.RecursiveProof, *state.RecursiveProof, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	AddGeneratedRecursiveProof(ctx context.Context, proof *state.RecursiveProof, dbTx pgx.Tx) error
	UpdateGeneratedRecursiveProof(ctx context.Context, proof *state.RecursiveProof, dbTx pgx.Tx) error
	DeleteGeneratedRecursiveProof(ctx context.Context, batchNumber uint64, batchNumberFinal uint64, dbTx pgx.Tx) error
	DeleteUngeneratedRecursiveProofs(ctx context.Context, dbTx pgx.Tx) error
}
