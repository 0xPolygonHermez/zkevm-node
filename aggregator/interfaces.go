package aggregator

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

// ethTxManager contains the methods required to send txs to
// ethereum.
type ethTxManager interface {
	VerifyBatch(batchNum uint64, proof *pb.GetProofResponse) error
}

// etherman contains the methods required to interact with ethereum
type etherman interface {
	GetLatestVerifiedBatchNum() (uint64, error)
	GetPublicAddress() common.Address
}

// aggregatorTxProfitabilityChecker interface for different profitability
// checking algorithms.
type aggregatorTxProfitabilityChecker interface {
	IsProfitable(context.Context, *big.Int) (bool, error)
}

// proverClient is a wrapper to the prover service
type proverClientInterface interface {
	GetURI() string
	IsIdle(ctx context.Context) bool
	GetGenProofID(ctx context.Context, inputProver *pb.InputProver) (string, error)
	GetResGetProof(ctx context.Context, genProofID string, batchNumber uint64) (*pb.GetProofResponse, error)
}

// stateInterface gathers the methods to interact with the state.
type stateInterface interface {
	GetLastVerifiedBatch(ctx context.Context, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetVirtualBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	AddGeneratedProof(ctx context.Context, proof *state.Proof, dbTx pgx.Tx) error
	UpdateGeneratedProof(ctx context.Context, proof *state.Proof, dbTx pgx.Tx) error
	GetGeneratedProofByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Proof, error)
	DeleteGeneratedProof(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	DeleteUngeneratedProofs(ctx context.Context, dbTx pgx.Tx) error
	GetWIPProofByProver(ctx context.Context, prover string, dbTx pgx.Tx) (*state.Proof, error)
}
