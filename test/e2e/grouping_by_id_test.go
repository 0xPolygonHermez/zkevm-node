package e2e

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestSequenceSize(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	defer func() { require.NoError(t, operations.Teardown()) }()

	err := operations.Teardown()
	require.NoError(t, err)
	opsCfg := operations.GetDefaultOperationsConfig()
	opsCfg.State.MaxCumulativeGasUsed = 80000000000
	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsman.Setup()
	require.NoError(t, err)

	batch, err := opsman.State().GetBatchByNumber(ctx, 0, nil)
	require.NoError(t, err)

	auth, err := operations.GetAuth(operations.DefaultSequencerPrivateKey, operations.DefaultL2ChainID)
	require.NoError(t, err)

	amount := big.NewInt(10000)
	toAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")

	payload := make([]byte, 1000) // 10mb

	_, err = rand.Read(payload)
	require.NoError(t, err)

	client, err := ethclient.Dial("http://localhost:8545")
	require.NoError(t, err)

	gasPrice, err := client.SuggestGasPrice(ctx)
	require.NoError(t, err)

	tx := ethtypes.NewTransaction(1, toAddress, amount, 21000, gasPrice, payload)

	fmt.Printf("%+v", tx)

	sequences := []types.Sequence{
		{
			GlobalExitRoot: batch.GlobalExitRoot,
			Timestamp:      batch.Timestamp.Unix(),
			Txs:            []ethtypes.Transaction{*tx},
		},
	}

	ethman, err := etherman.NewClient(etherman.Config{
		URL: "http://localhost:8545",
	},
		auth,
		common.HexToAddress("0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"),
		common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"),
		common.HexToAddress("0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"))
	require.NoError(t, err)

	// Check if can be send
	tx, err = ethman.EstimateGasSequenceBatches(sequences)

	//require.ErrorContains(t, err, "gas required exceeds allowance")

	fmt.Println(tx.GasPrice())

}
