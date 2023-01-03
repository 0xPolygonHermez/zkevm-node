package ethtxmanager

import (
	"context"
	"math/big"
	"time"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

type ethermanInterface interface {
	TrustedVerifyBatches(ctx context.Context, lastVerifiedBatch, newVerifiedBatch uint64, inputs *ethmanTypes.FinalProofInputs, gasLimit uint64, gasPrice, nonce *big.Int, noSend bool) (*types.Transaction, error)
	EstimateGasForTrustedVerifyBatches(lastVerifiedBatch, newVerifiedBatch uint64, inputs *ethmanTypes.FinalProofInputs) (uint64, error)
	SequenceBatches(ctx context.Context, sequences []ethmanTypes.Sequence, gasLimit uint64, gasPrice, nonce *big.Int, noSend bool) (*types.Transaction, error)
	EstimateGasSequenceBatches(sequences []ethmanTypes.Sequence) (*types.Transaction, error)
	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	WaitTxToBeMined(ctx context.Context, tx *types.Transaction, timeout time.Duration) error
	SendTx(ctx context.Context, tx *types.Transaction) error
	CurrentNonce(ctx context.Context) (uint64, error)
	SuggestedGasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, from common.Address, to *common.Address, value *big.Int, data []byte) (uint64, error)
	CheckTxWasMined(ctx context.Context, txHash common.Hash) (bool, *types.Receipt, error)
	SignTx(ctx context.Context, tx *types.Transaction) (*types.Transaction, error)
	GetRevertMessage(ctx context.Context, tx types.Transaction) (string, error)
}

type storageInterface interface {
	Add(ctx context.Context, mTx monitoredTx, dbTx pgx.Tx) error
	Get(ctx context.Context, owner, id string, dbTx pgx.Tx) (monitoredTx, error)
	GetByStatus(ctx context.Context, owner *string, statuses []MonitoredTxStatus, dbTx pgx.Tx) ([]monitoredTx, error)
	Update(ctx context.Context, mTx monitoredTx, dbTx pgx.Tx) error
}
