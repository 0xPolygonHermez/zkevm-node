package ethtxmanager

import (
	"context"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type etherman interface {
	SequenceBatches(sequences []state.Sequence, gasLimit uint64, gasPrice, nonce *big.Int) (*types.Transaction, error)
	VerifyBatch(batchNumber uint64, resGetProof *pb.GetProofResponse, gasLimit uint64, gasPrice, nonce *big.Int) (*types.Transaction, error)
	EstimateGasForVerifyBatch(batchNumber uint64, resGetProof *pb.GetProofResponse) (uint64, error)
	EstimateGasSequenceBatches(sequences []state.Sequence) (*types.Transaction, error)
	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	WaitTxToBeMined(hash common.Hash, timeout time.Duration) error
}
