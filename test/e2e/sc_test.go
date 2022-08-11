package e2e

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/EmitLog2"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/FailureTest"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmitLog2(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var err error
	err = operations.Teardown()
	require.NoError(t, err)

	defer func() { require.NoError(t, operations.Teardown()) }()

	ctx := context.Background()
	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsMan.Setup()
	require.NoError(t, err)

	l1Client, l2Client, err := operations.GetL1AndL2Clients()
	require.NoError(t, err)

	l1Auth, l2Auth, err := operations.GetL1AndL2Authorizations()
	require.NoError(t, err)

	test := func(t *testing.T, auth *bind.TransactOpts, client *ethclient.Client) {
		scAddr, scTx, sc, err := EmitLog2.DeployEmitLog2(auth, client)
		require.NoError(t, err)

		logTx(scTx)
		err = operations.WaitTxToBeMined(client, scTx.Hash(), operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		scCallTx, err := sc.EmitLogs(auth)
		require.NoError(t, err)

		logTx(scCallTx)
		err = operations.WaitTxToBeMined(client, scCallTx.Hash(), operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		scCallTxReceipt, err := client.TransactionReceipt(ctx, scCallTx.Hash())
		require.NoError(t, err)

		filterBlock := scCallTxReceipt.BlockNumber
		logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: filterBlock, ToBlock: filterBlock,
			Addresses: []common.Address{scAddr},
		})
		require.NoError(t, err)
		assert.Equal(t, 3, len(logs))

		j, err := json.Marshal(logs)
		require.NoError(t, err)
		log.Debug(string(j))

		_, err = sc.ParseLog(logs[0])
		require.NoError(t, err)
		logA, err := sc.ParseLogA(logs[1])
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(1), logA.A)
		logABCD, err := sc.ParseLogABCD(logs[2])
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(1), logABCD.A)
		assert.Equal(t, big.NewInt(2), logABCD.B)
		assert.Equal(t, big.NewInt(3), logABCD.C)
		assert.Equal(t, big.NewInt(4), logABCD.D)
	}

	log.Debug("testing l1")
	test(t, l1Auth, l1Client)

	log.Debug("testing l2")
	test(t, l2Auth, l2Client)
}

func TestFailureTest(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var err error
	err = operations.Teardown()
	require.NoError(t, err)

	defer func() { require.NoError(t, operations.Teardown()) }()

	ctx := context.Background()
	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsMan.Setup()
	require.NoError(t, err)

	l1Client, l2Client, err := operations.GetL1AndL2Clients()
	require.NoError(t, err)

	l1Auth, l2Auth, err := operations.GetL1AndL2Authorizations()
	require.NoError(t, err)

	test := func(t *testing.T, auth *bind.TransactOpts, client *ethclient.Client) {
		log.Debug("deploying SC")
		_, scTx, sc, err := FailureTest.DeployFailureTest(auth, client)
		require.NoError(t, err)

		logTx(scTx)
		err = operations.WaitTxToBeMined(client, scTx.Hash(), operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		log.Debug("storing value")
		scCallTx, err := sc.Store(auth, big.NewInt(1))
		require.NoError(t, err)

		logTx(scCallTx)
		err = operations.WaitTxToBeMined(client, scCallTx.Hash(), operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		log.Debug("storing value with revert")
		_, err = sc.StoreAndFail(auth, big.NewInt(2))
		assert.Equal(t, err.Error(), "execution reverted: this method always fails")
	}

	log.Debug("testing l1")
	test(t, l1Auth, l1Client)

	log.Debug("testing l2")
	test(t, l2Auth, l2Client)
}

func logTx(tx *types.Transaction) {
	b, _ := tx.MarshalBinary()
	log.Debug(tx.Hash(), " ", hex.EncodeToHex(b))
}
