package e2e

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/EmitLog2"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmitLog2(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var err error
	operations.Teardown()
	defer func() { require.NoError(t, operations.Teardown()) }()

	ctx := context.Background()
	opsCfg := getDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsMan.Setup()
	require.NoError(t, err)

	l1Client, l2Client, err := getClients()
	require.NoError(t, err)

	l1Auth, l2Auth, err := getAuth()
	require.NoError(t, err)

	test := func(t *testing.T, auth *bind.TransactOpts, client *ethclient.Client, gasLimit *uint64) {
		if gasLimit != nil {
			auth.GasPrice = big.NewInt(1)
			auth.GasLimit = *gasLimit
		} else {
			auth.GasPrice = nil
			auth.GasLimit = 0
		}

		scAddr, scTx, sc, err := EmitLog2.DeployEmitLog2(auth, client)
		require.NoError(t, err)

		log.Debug(scTx.Hash())
		err = operations.WaitTxToBeMined(client, scTx.Hash(), defaultTimeoutTxToBeMined)
		require.NoError(t, err)

		scCallTx, err := sc.EmitLogs(auth)
		require.NoError(t, err)

		log.Debug(scCallTx.Hash())
		err = operations.WaitTxToBeMined(l1Client, scCallTx.Hash(), defaultTimeoutTxToBeMined)
		require.NoError(t, err)

		scCallTxReceipt, err := l1Client.TransactionReceipt(ctx, scCallTx.Hash())
		require.NoError(t, err)

		filterBlock := scCallTxReceipt.BlockNumber
		logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: filterBlock, ToBlock: filterBlock,
			Addresses: []common.Address{scAddr},
		})
		require.NoError(t, err)
		assert.Equal(t, 3, len(logs))
	}

	gasLimit := uint64(400000)
	test(t, l1Auth, l1Client, nil)
	test(t, l2Auth, l2Client, &gasLimit)
}
