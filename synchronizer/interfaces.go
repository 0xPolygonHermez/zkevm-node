package synchronizer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/state"
)

// Consumer interfaces required by the package.

// gasPriceEstimator contains the methods required to interact with gas price estimator
type gasPriceEstimator interface {
	UpdateGasPriceAvg(newValue *big.Int)
}

// synchronizerEtherman contains the methods required to interact with ethereum.
type synchronizerEtherman interface {
	GetLatestProposedBatchNumber() (uint64, error)
	GetLatestConsolidatedBatchNumber() (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]state.Block, map[common.Hash][]etherman.Order, error)
	EthBlockByNumber(ctx context.Context, blockNum uint64) (*types.Block, error)
}
