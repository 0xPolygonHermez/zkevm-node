package etrog

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type BatchStepsExecutorEtrogStateInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error
	CloseBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error
	ProcessBatch(ctx context.Context, request state.ProcessRequest, updateMerkleTree bool) (*state.ProcessBatchResponse, error)
	StoreTransaction(ctx context.Context, batchNumber uint64, processedTx *state.ProcessTransactionResponse, coinbase common.Address, timestamp uint64, egpLog *state.EffectiveGasPriceLog, dbTx pgx.Tx) (*ethTypes.Header, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetForkIDByBatchNumber(batchNumber uint64) uint64
	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	UpdateBatchL2Data(ctx context.Context, batchNumber uint64, batchL2Data []byte, dbTx pgx.Tx) error
}

type BatchStepsExecutorEtrog struct {
	state BatchStepsExecutorEtrogStateInterface
}

// NewSyncTrustedStateEtrogExecutor creates a new prcessor for sync with L2 batches
func NewSyncTrustedStateEtrogExecutor(zkEVMClient l2_shared.ZkEVMClientInterface,
	state l2_shared.StateInterface, stateBatchExecutor BatchStepsExecutorEtrogStateInterface,
	sync l2_shared.SyncInterface) *l2_shared.SyncTrustedStateTemplate {
	executorSteps := &BatchStepsExecutorEtrog{state: stateBatchExecutor}
	executor := &l2_shared.SyncTrustedStateBatchExecutorTemplate{
		Steps: executorSteps,
	}
	a := l2_shared.NewSyncTrustedStateTemplate(executor, zkEVMClient, state, sync)
	return a
}

// FullProcess process a batch that is not on database, so is the first time we process it
func (b *BatchStepsExecutorEtrog) FullProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	return nil, ErrNotImplemented
}

// IncrementalProcess process a batch that we have processed before, and we have the intermediate state root, so is going to be process only new Tx
func (b *BatchStepsExecutorEtrog) IncrementalProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	return nil, ErrNotImplemented
}

// ReProcess process a batch that we have processed before, but we don't have the intermediate state root, so we need to reprocess it
func (b *BatchStepsExecutorEtrog) ReProcess(ctx context.Context, data *l2_shared.ProcessData, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	return nil, ErrNotImplemented
}

// CloseBatch close a batch
func (b *BatchStepsExecutorEtrog) CloseBatch(ctx context.Context, trustedBatch *types.Batch, dbTx pgx.Tx) error {
	return ErrNotImplemented
}
