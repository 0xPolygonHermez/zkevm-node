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
	SequenceBatches(sequences []state.Sequence, gas uint64, nonce *big.Int) (*types.Transaction, error)
	VerifyBatch(batchNumber uint64, resGetProof *pb.GetProofResponse, gasLimit uint64, nonce *big.Int) (*types.Transaction, error)
	EstimateGasForVerifyBatch(batchNumber uint64, resGetProof *pb.GetProofResponse) (uint64, error)
	EstimateGasSequenceBatches(sequences []state.Sequence) (*types.Transaction, error)
	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	WaitTxToBeMined(ctx context.Context, tx *types.Transaction, timeout time.Duration) error
}

type stateInterface interface {
	GetSequencesWithoutGroup(ctx context.Context, dbTx pgx.Tx) ([]state.Sequence, error)
	GetPendingSequenceGroups(ctx context.Context, dbTx pgx.Tx) ([]state.SequenceGroup, error)
	GetLastSequenceGroup(ctx context.Context, dbTx pgx.Tx) (*state.SequenceGroup, error)
	GetSequencesByBatchNums(ctx context.Context, batchNumbers []uint64, dbTx pgx.Tx) ([]state.Sequence, error)
	AddSequenceGroup(ctx context.Context, sequenceGroup state.SequenceGroup, dbTx pgx.Tx) error
	SetSequenceGroupAsConfirmed(ctx context.Context, txHash common.Hash, dbTx pgx.Tx) error
	UpdateSequenceGroupTx(ctx context.Context, oldTxHash, newTxHash common.Hash, dbTx pgx.Tx) error

	GetPendingProofs(ctx context.Context, dbTx pgx.Tx) ([]state.Proof, error)
	UpdateProofTx(ctx context.Context, batchNumber uint64, newTxHash common.Hash, dbTx pgx.Tx) error
	SetProofAsConfirmed(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
}
