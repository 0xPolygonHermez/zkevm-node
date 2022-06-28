package ethtxmanager

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethmanTypes "github.com/hermeznetwork/hermez-core/ethermanv2/types"
)

type etherman interface {
	SequenceBatches(sequences []ethmanTypes.Sequence, gasLimit uint64) (*types.Transaction, error)
	EstimateGasSequenceBatches(sequences []ethmanTypes.Sequence) (uint64, error)
	GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}
