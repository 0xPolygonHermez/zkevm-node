package ethtxmanager

import (
	"context"

	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type etherman interface {
	SequenceBatches(sequences []ethmanTypes.Sequence, gasLimit uint64) (*types.Transaction, error)
	VerifyBatch(batchNumber uint64, resGetProof *pb.GetProofResponse, gasLimit uint64) (*types.Transaction, error)
	EstimateGasForVerifyBatch(batchNumber uint64, resGetProof *pb.GetProofResponse) (uint64, error)
	EstimateGasSequenceBatches(sequences []ethmanTypes.Sequence) (uint64, error)
	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}
