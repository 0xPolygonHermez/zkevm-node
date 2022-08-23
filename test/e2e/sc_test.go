package e2e

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/EmitLog2"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/FailureTest"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/Read"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

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

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

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
		assert.Equal(t, "execution reverted: this method always fails", err.Error())
	}
}

func TestRead(t *testing.T) {
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

	for _, network := range networks {
		log.Debugf(network.Name)
		client := operations.MustGetClient(network.URL)
		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)

		const ownerName = "this is the owner name"
		callOpts := &bind.CallOpts{Pending: false}

		log.Debug("deploying SC")
		_, scTx, sc, err := Read.DeployRead(auth, client, ownerName)
		require.NoError(t, err)

		logTx(scTx)
		err = operations.WaitTxToBeMined(client, scTx.Hash(), operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		log.Debug("read string public variable directly")
		ownerNameValue, err := sc.OwnerName(callOpts)
		require.NoError(t, err)
		require.Equal(t, ownerName, ownerNameValue)

		log.Debug("read address public variable directly")
		ownerValue, err := sc.Owner(callOpts)
		require.NoError(t, err)
		require.Equal(t, auth.From, ownerValue)

		tA := Read.Readtoken{
			Name:     "Token A",
			Quantity: big.NewInt(50),
			Address:  common.HexToAddress("0x1"),
		}

		tB := Read.Readtoken{
			Name:     "Token B",
			Quantity: big.NewInt(30),
			Address:  common.HexToAddress("0x2"),
		}

		log.Debug("public add token")
		tx, err := sc.PublicAddToken(auth, tA)
		require.NoError(t, err)
		logTx(tx)
		err = operations.WaitTxToBeMined(client, tx.Hash(), operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		log.Debug("external add token")
		tx, err = sc.ExternalAddToken(auth, tB)
		require.NoError(t, err)
		logTx(tx)
		err = operations.WaitTxToBeMined(client, tx.Hash(), operations.DefaultTimeoutTxToBeMined)
		require.NoError(t, err)

		log.Debug("read mapping public variable directly")
		tk, err := sc.Tokens(callOpts, tA.Address)
		require.NoError(t, err)
		require.Equal(t, tA.Name, tk.Name)
		require.Equal(t, tA.Quantity, tk.Quantity)
		require.Equal(t, tA.Address, tk.Address)

		tk, err = sc.Tokens(callOpts, tB.Address)
		require.NoError(t, err)
		require.Equal(t, tB.Name, tk.Name)
		require.Equal(t, tB.Quantity, tk.Quantity)
		require.Equal(t, tB.Address, tk.Address)

		log.Debug("public struct read")
		tk, err = sc.PublicGetToken(callOpts, tA.Address)
		require.NoError(t, err)
		require.Equal(t, tA.Name, tk.Name)
		require.Equal(t, tA.Quantity, tk.Quantity)
		require.Equal(t, tA.Address, tk.Address)

		log.Debug("external struct read")
		tk, err = sc.ExternalGetToken(callOpts, tB.Address)
		require.NoError(t, err)
		require.Equal(t, tB.Name, tk.Name)
		require.Equal(t, tB.Quantity, tk.Quantity)
		require.Equal(t, tB.Address, tk.Address)

		log.Debug("public uint256 read")
		value, err := sc.PublicRead(callOpts)
		require.NoError(t, err)
		require.Equal(t, 0, big.NewInt(1).Cmp(value))

		log.Debug("external uint256 read")
		value, err = sc.ExternalRead(callOpts)
		require.NoError(t, err)
		require.Equal(t, 0, big.NewInt(1).Cmp(value))

		log.Debug("public uint256 read with parameter")
		value, err = sc.PublicReadWParams(callOpts, big.NewInt(1))
		require.NoError(t, err)
		require.Equal(t, 0, big.NewInt(2).Cmp(value))

		log.Debug("external uint256 read with parameter")
		value, err = sc.ExternalReadWParams(callOpts, big.NewInt(1))
		require.NoError(t, err)
		require.Equal(t, 0, big.NewInt(2).Cmp(value))
	}
}
