package synchronizer

import (
	context "context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_Exploratory(t *testing.T) {
	t.Skip("no real test, just exploratory")
	cfg := etherman.Config{
		URL: "http://localhost:8545",
	}

	l1Config := etherman.L1Config{
		L1ChainID:                 1337,
		ZkEVMAddr:                 common.HexToAddress("0x610178dA211FEF7D417bC0e6FeD39F05609AD788"),
		MaticAddr:                 common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"),
		GlobalExitRootManagerAddr: common.HexToAddress("0x2279B7A0a67DB372996a5FaB50D91eAA73d2eBe6"),
	}

	etherman, err := etherman.NewClient(cfg, l1Config)
	require.NoError(t, err)
	worker := newWorker(etherman)
	ch := make(chan genericResponse[responseRollupInfoByBlockRange])
	blockRange := blockRange{
		fromBlock: 100,
		toBlock:   20000,
	}
	err = worker.asyncRequestRollupInfoByBlockRange(context.Background(), ch, nil, blockRange)
	require.NoError(t, err)
	result := <-ch
	require.Equal(t, result.err.Error(), "not found")
}
