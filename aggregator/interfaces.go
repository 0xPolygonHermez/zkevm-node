package aggregator

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

// ethTxManager contains the methods required to send txs to
// ethereum.
type ethTxManager interface {
	VerifyBatch(batchNum uint64, proof *pb.GetProofResponse) error
}

// etherman contains the methods required to interact with ethereum
type etherman interface {
	GetLatestVerifiedBatchNum() (uint64, error)
}

// aggregatorTxProfitabilityChecker interface for different profitability
// checking algorithms.
type aggregatorTxProfitabilityChecker interface {
	IsProfitable(context.Context, *big.Int) (bool, error)
}

// stateInterface gathers the methods to interract with the state.
type stateInterface interface {
	GetLastVerifiedBatch(ctx context.Context, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetLastVerifiedBatchNumberSeenOnEthereum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetVirtualBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetStateRootByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (common.Hash, error)
	GetLocalExitRootByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (common.Hash, error)
	GetBlockNumVirtualBatchByBatchNum(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (uint64, error)
}
