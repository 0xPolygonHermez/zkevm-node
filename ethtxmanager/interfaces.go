package ethtxmanager

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

type etherman interface {
	// SequenceBatches(ctx context.Context, sequences []ethmanTypes.Sequence, gasLimit uint64, gasPrice, nonce *big.Int) (*types.Transaction, error)
	VerifyBatches(ctx context.Context, lastVerifiedBatch, newVerifiedBatch uint64, resGetProof *pb.GetProofResponse, gasLimit uint64, gasPrice, nonce *big.Int) (*types.Transaction, error)
	EstimateGasForVerifyBatches(lastVerifiedBatch, newVerifiedBatch uint64, resGetProof *pb.GetProofResponse) (uint64, error)
	// EstimateGasSequenceBatches(sequences []ethmanTypes.Sequence) (*types.Transaction, error)
	SequenceBatches(ctx context.Context, sequences []state.Sequence, gasLimit uint64, gasPrice, nonce *big.Int) (*types.Transaction, error)
	// VerifyBatch(ctx context.Context, batchNumber uint64, resGetProof *pb.GetProofResponse, gasLimit uint64, gasPrice, nonce *big.Int) (*types.Transaction, error)
	// EstimateGasForVerifyBatch(batchNumber uint64, resGetProof *pb.GetProofResponse) (uint64, error)
	EstimateGasSequenceBatches(sequences []state.Sequence) (*types.Transaction, error)
	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	WaitTxToBeMined(ctx context.Context, tx *types.Transaction, timeout time.Duration) error
}

type stateInterface interface {
	GetSequencesWithoutGroup(ctx context.Context, dbTx pgx.Tx) ([]state.Sequence, error)
	GetPendingSequenceGroups(ctx context.Context, dbTx pgx.Tx) ([]state.SequenceGroup, error)
	GetLastSequenceGroup(ctx context.Context, dbTx pgx.Tx) (*state.SequenceGroup, error)
	GetSequencesByBatchNums(ctx context.Context, fromBatchNumber, toBatchNumber uint64, dbTx pgx.Tx) ([]state.Sequence, error)
	SetSequenceGroupAsConfirmed(ctx context.Context, txHash common.Hash, dbTx pgx.Tx) error
	UpdateSequenceGroupTx(ctx context.Context, oldTxHash, newTxHash common.Hash, dbTx pgx.Tx) error

	GetPendingProofs(ctx context.Context, dbTx pgx.Tx) ([]state.Proof, error)
	UpdateProofTx(ctx context.Context, batchNumber uint64, newTxHash common.Hash, nonce uint64, dbTx pgx.Tx) error
	SetProofAsConfirmed(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
}
