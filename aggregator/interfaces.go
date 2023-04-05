package aggregator

import (
	stdContext "context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	"github.com/0xPolygonHermez/zkevm-node/context"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

type proverInterface interface {
	Name() string
	ID() string
	Addr() string
	IsIdle() (bool, error)
	BatchProof(input *pb.InputProver) (*string, error)
	AggregatedProof(inputProof1, inputProof2 string) (*string, error)
	FinalProof(inputProof string, aggregatorAddr string) (*string, error)
	WaitRecursiveProof(ctx stdContext.Context, proofID string) (string, error)
	WaitFinalProof(ctx stdContext.Context, proofID string) (*pb.FinalProof, error)
}

// ethTxManager contains the methods required to send txs to
// ethereum.
type ethTxManager interface {
	Add(ctx *context.RequestContext, owner, id string, from common.Address, to *common.Address, value *big.Int, data []byte, dbTx pgx.Tx) error
	Result(ctx *context.RequestContext, owner, id string, dbTx pgx.Tx) (ethtxmanager.MonitoredTxResult, error)
	ResultsByStatus(ctx *context.RequestContext, owner string, statuses []ethtxmanager.MonitoredTxStatus, dbTx pgx.Tx) ([]ethtxmanager.MonitoredTxResult, error)
	ProcessPendingMonitoredTxs(ctx *context.RequestContext, owner string, failedResultHandler ethtxmanager.ResultHandler, dbTx pgx.Tx)
}

// etherman contains the methods required to interact with ethereum
type etherman interface {
	GetLatestVerifiedBatchNum() (uint64, error)
	BuildTrustedVerifyBatchesTxData(lastVerifiedBatch, newVerifiedBatch uint64, inputs *ethmanTypes.FinalProofInputs) (to *common.Address, data []byte, err error)
}

// aggregatorTxProfitabilityChecker interface for different profitability
// checking algorithms.
type aggregatorTxProfitabilityChecker interface {
	IsProfitable(*context.RequestContext, *big.Int) (bool, error)
}

// stateInterface gathers the methods to interact with the state.
type stateInterface interface {
	BeginStateTransaction(ctx *context.RequestContext) (pgx.Tx, error)
	CheckProofContainsCompleteSequences(ctx *context.RequestContext, proof *state.Proof, dbTx pgx.Tx) (bool, error)
	GetLastVerifiedBatch(ctx *context.RequestContext, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetProofReadyToVerify(ctx *context.RequestContext, lastVerfiedBatchNumber uint64, dbTx pgx.Tx) (*state.Proof, error)
	GetVirtualBatchToProve(ctx *context.RequestContext, lastVerfiedBatchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetProofsToAggregate(ctx *context.RequestContext, dbTx pgx.Tx) (*state.Proof, *state.Proof, error)
	GetBatchByNumber(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	AddGeneratedProof(ctx *context.RequestContext, proof *state.Proof, dbTx pgx.Tx) error
	UpdateGeneratedProof(ctx *context.RequestContext, proof *state.Proof, dbTx pgx.Tx) error
	DeleteGeneratedProofs(ctx *context.RequestContext, batchNumber uint64, batchNumberFinal uint64, dbTx pgx.Tx) error
	DeleteUngeneratedProofs(ctx *context.RequestContext, dbTx pgx.Tx) error
	CleanupGeneratedProofs(ctx *context.RequestContext, batchNumber uint64, dbTx pgx.Tx) error
	CleanupLockedProofs(ctx *context.RequestContext, duration string, dbTx pgx.Tx) (int64, error)
}
