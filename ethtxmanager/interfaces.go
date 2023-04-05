package ethtxmanager

import (
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/context"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

type ethermanInterface interface {
	GetTx(ctx *context.RequestContext, txHash common.Hash) (*types.Transaction, bool, error)
	GetTxReceipt(ctx *context.RequestContext, txHash common.Hash) (*types.Receipt, error)
	WaitTxToBeMined(ctx *context.RequestContext, tx *types.Transaction, timeout time.Duration) (bool, error)
	SendTx(ctx *context.RequestContext, tx *types.Transaction) error
	CurrentNonce(ctx *context.RequestContext, account common.Address) (uint64, error)
	SuggestedGasPrice(ctx *context.RequestContext) (*big.Int, error)
	EstimateGas(ctx *context.RequestContext, from common.Address, to *common.Address, value *big.Int, data []byte) (uint64, error)
	CheckTxWasMined(ctx *context.RequestContext, txHash common.Hash) (bool, *types.Receipt, error)
	SignTx(ctx *context.RequestContext, sender common.Address, tx *types.Transaction) (*types.Transaction, error)
	GetRevertMessage(ctx *context.RequestContext, tx *types.Transaction) (string, error)
}

type storageInterface interface {
	Add(ctx *context.RequestContext, mTx monitoredTx, dbTx pgx.Tx) error
	Get(ctx *context.RequestContext, owner, id string, dbTx pgx.Tx) (monitoredTx, error)
	GetByStatus(ctx *context.RequestContext, owner *string, statuses []MonitoredTxStatus, dbTx pgx.Tx) ([]monitoredTx, error)
	GetByBlock(ctx *context.RequestContext, fromBlock, toBlock *uint64, dbTx pgx.Tx) ([]monitoredTx, error)
	Update(ctx *context.RequestContext, mTx monitoredTx, dbTx pgx.Tx) error
}

type stateInterface interface {
	GetLastBlock(ctx *context.RequestContext, dbTx pgx.Tx) (*state.Block, error)
}
