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
	VerifyBatches(ctx context.Context, lastVerifiedBatch, newVerifiedBatch uint64, resGetProof *pb.GetProofResponse, gasLimit uint64, gasPrice, nonce *big.Int) (*types.Transaction, error)
	EstimateGasForVerifyBatches(lastVerifiedBatch, newVerifiedBatch uint64, resGetProof *pb.GetProofResponse) (uint64, error)
	SequenceBatches(ctx context.Context, sequences []ethmanTypes.Sequence, gasLimit uint64, gasPrice, nonce *big.Int) (*types.Transaction, error)
	EstimateGasSequenceBatches(sequences []ethmanTypes.Sequence) (*types.Transaction, error)
	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	WaitTxToBeMined(ctx context.Context, tx *types.Transaction, timeout time.Duration) error
}
