package aggregator

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/proverclient"
)

// Consumer interfaces required by the package.

// etherman contains the methods required to interact with
// ethereum.
type etherman interface {
	ConsolidateBatch(batchNum *big.Int, proof *proverclient.ResGetProof) (*types.Transaction, error)
}

// aggregatorTxProfitabilityChecker interface for different profitability
// checking algorithms.
type aggregatorTxProfitabilityChecker interface {
	IsProfitable(context.Context, *big.Int) (bool, error)
}
