package ethtxmanager

import (
	"context"
	"math/big"
	"time"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type etherman interface {
	SequenceBatches(sequences []ethmanTypes.Sequence, gasLimit uint64, gasPrice, nonce *big.Int) (*types.Transaction, error)
	VerifyBatch(batchNumber uint64, resGetProof *pb.GetProofResponse, gasLimit uint64, gasPrice, nonce *big.Int) (*types.Transaction, error)
	EstimateGasForVerifyBatch(batchNumber uint64, resGetProof *pb.GetProofResponse) (uint64, error)
	EstimateGasSequenceBatches(sequences []ethmanTypes.Sequence) (*types.Transaction, error)
	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	WaitTxToBeMined(hash common.Hash, timeout time.Duration) error
}

type state interface {
	WaitSequencingTxToBeSynced(parentCtx context.Context, tx *types.Transaction, timeout time.Duration) error
	WaitVerifiedBatchToBeSynced(parentCtx context.Context, batchNumber uint64, timeout time.Duration) error
}
