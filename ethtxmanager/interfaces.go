package ethtxmanager

import (
	"context"
	"math/big"
	"time"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type etherman interface {
	TrustedVerifyBatches(ctx context.Context, lastVerifiedBatch, newVerifiedBatch uint64, inputs *ethmanTypes.FinalProofInputs, gasLimit uint64, gasPrice, nonce *big.Int, noSend bool) (*types.Transaction, error)
	EstimateGasForTrustedVerifyBatches(lastVerifiedBatch, newVerifiedBatch uint64, inputs *ethmanTypes.FinalProofInputs) (uint64, error)
	SequenceBatches(ctx context.Context, sequences []ethmanTypes.Sequence, gasLimit uint64, gasPrice, nonce *big.Int, noSend bool) (*types.Transaction, error)
	EstimateGasSequenceBatches(sequences []ethmanTypes.Sequence) (*types.Transaction, error)
	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	WaitTxToBeMined(ctx context.Context, tx *types.Transaction, timeout time.Duration) error
	SendTx(ctx context.Context, tx *types.Transaction) error
	CurrentNonce(ctx context.Context) (uint64, error)
	CheckTxWasMined(ctx context.Context, txHash common.Hash) (bool, *types.Receipt, error)
}

type state interface {
	WaitSequencingTxToBeSynced(parentCtx context.Context, tx *types.Transaction, timeout time.Duration) error
	WaitVerifiedBatchToBeSynced(parentCtx context.Context, batchNumber uint64, timeout time.Duration) error
}
