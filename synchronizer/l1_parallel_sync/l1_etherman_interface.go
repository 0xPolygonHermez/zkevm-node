package l1_parallel_sync

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

// L1ParallelEthermanInterface is an interface for the etherman package
type L1ParallelEthermanInterface interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*ethTypes.Header, error)
	GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]etherman.Block, map[common.Hash][]etherman.Order, error)
	EthBlockByNumber(ctx context.Context, blockNumber uint64) (*ethTypes.Block, error)
	GetLatestBatchNumber() (uint64, error)
	GetTrustedSequencerURL() (string, error)
	VerifyGenBlockNumber(ctx context.Context, genBlockNumber uint64) (bool, error)
	GetLatestVerifiedBatchNum() (uint64, error)
}
