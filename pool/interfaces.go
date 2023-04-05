package pool

import (
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/context"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

type storage interface {
	AddTx(ctx *context.RequestContext, tx Transaction) error
	CountTransactionsByStatus(ctx *context.RequestContext, status TxStatus) (uint64, error)
	DeleteTransactionsByHashes(ctx *context.RequestContext, hashes []common.Hash) error
	GetGasPrice(ctx *context.RequestContext) (uint64, error)
	GetNonce(ctx *context.RequestContext, address common.Address) (uint64, error)
	GetPendingTxHashesSince(ctx *context.RequestContext, since time.Time) ([]common.Hash, error)
	GetTxsByFromAndNonce(ctx *context.RequestContext, from common.Address, nonce uint64) ([]Transaction, error)
	GetTxsByStatus(ctx *context.RequestContext, state TxStatus, isClaims bool, limit uint64) ([]Transaction, error)
	GetNonWIPTxsByStatus(ctx *context.RequestContext, status TxStatus, isClaims bool, limit uint64) ([]Transaction, error)
	IsTxPending(ctx *context.RequestContext, hash common.Hash) (bool, error)
	SetGasPrice(ctx *context.RequestContext, gasPrice uint64) error
	UpdateTxsStatus(ctx *context.RequestContext, hashes []string, newStatus TxStatus) error
	UpdateTxStatus(ctx *context.RequestContext, hash common.Hash, newStatus TxStatus, isWIP bool) error
	UpdateTxWIPStatus(ctx *context.RequestContext, hash common.Hash, isWIP bool) error
	GetTxs(ctx *context.RequestContext, filterStatus TxStatus, isClaims bool, minGasPrice, limit uint64) ([]*Transaction, error)
	GetTxFromAddressFromByHash(ctx *context.RequestContext, hash common.Hash) (common.Address, uint64, error)
	GetTxByHash(ctx *context.RequestContext, hash common.Hash) (*Transaction, error)
	GetTxZkCountersByHash(ctx *context.RequestContext, hash common.Hash) (*state.ZKCounters, error)
	DeleteTransactionByHash(ctx *context.RequestContext, hash common.Hash) error
	MarkWIPTxsAsPending(ctx *context.RequestContext) error
	MinGasPriceSince(ctx *context.RequestContext, timestamp time.Time) (uint64, error)
	DepositCountExists(ctx *context.RequestContext, depositCount uint64) (bool, error)
}

type stateInterface interface {
	GetBalance(ctx *context.RequestContext, address common.Address, root common.Hash) (*big.Int, error)
	GetLastL2Block(ctx *context.RequestContext, dbTx pgx.Tx) (*types.Block, error)
	GetNonce(ctx *context.RequestContext, address common.Address, root common.Hash) (uint64, error)
	GetTransactionByHash(ctx *context.RequestContext, transactionHash common.Hash, dbTx pgx.Tx) (*types.Transaction, error)
	PreProcessTransaction(ctx *context.RequestContext, tx *types.Transaction, dbTx pgx.Tx) (*state.ProcessBatchResponse, error)
}
