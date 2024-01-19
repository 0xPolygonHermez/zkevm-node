package syncinterfaces

import (
	"context"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

type PoolInterface interface {
	DeleteReorgedTransactions(ctx context.Context, txs []*ethTypes.Transaction) error
	StoreTx(ctx context.Context, tx ethTypes.Transaction, ip string, isWIP bool) error
}
