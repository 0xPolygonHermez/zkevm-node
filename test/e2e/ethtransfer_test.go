package e2e

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func Test1000EthTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	// TODO: use opsman to spin up/down all the containers
	// defer func() {
	// 	require.NoError(t, operations.Teardown())
	// }()
	// operations.NewManager()

	// Load account with balance on local genesis
	auth, err := operations.GetAuth("0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d", big.NewInt(1001))
	require.NoError(t, err)
	// Load eth client
	client, err := ethclient.Dial("http://localhost:8124")
	require.NoError(t, err)
	// Send txs
	nTxs := 1001
	amount := big.NewInt(10000)
	toAddress := common.HexToAddress("0x0000000000000000000000000000000000000001")
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(1000000)
	log.Infof("Sending %d transactions...", nTxs)
	var lastTxHash common.Hash
	for i := 0; i < nTxs; i++ {
		nonce := uint64(i + 1)
		tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		err = client.SendTransaction(context.Background(), signedTx)
		require.NoError(t, err)
		if i == nTxs-1 {
			lastTxHash = signedTx.Hash()
		}
	}
	log.Infof("%d transactions sent without error. Waiting for all the transactions to be mined", nTxs)
	timeout := 3 * time.Minute
	err = operations.WaitTxToBeMined(client, lastTxHash, timeout)
	require.NoError(t, err)
}
