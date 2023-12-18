package eth_transfers

import (
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/params"
	"github.com/0xPolygonHermez/zkevm-node/test/benchmarks/sequencer/common/transactions"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	uniswap "github.com/0xPolygonHermez/zkevm-node/test/scripts/uniswap/pkg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	gasLimit  = 21000
	ethAmount = big.NewInt(0)
	sleepTime = 1 * time.Second
	countTxs  = 0
)

// TxSender sends eth transfer to the sequencer
func TxSender(l2Client *ethclient.Client, gasPrice *big.Int, auth *bind.TransactOpts, erc20SC *ERC20.ERC20, uniswapDeployments *uniswap.Deployments) ([]*types.Transaction, error) {
	fmt.Printf("sending tx num: %d\n", countTxs+1)
	senderNonce, err := l2Client.PendingNonceAt(params.Ctx, auth.From)
	if err != nil {
		panic(err)
	}
	tx := types.NewTx(&types.LegacyTx{
		GasPrice: gasPrice,
		Gas:      uint64(gasLimit),
		To:       &params.To,
		Value:    ethAmount,
		Data:     nil,
		Nonce:    senderNonce,
	})

	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return nil, err
	}

	err = l2Client.SendTransaction(params.Ctx, signedTx)
	for transactions.ShouldRetryError(err) {
		time.Sleep(sleepTime)
		err = l2Client.SendTransaction(params.Ctx, signedTx)
	}

	if err == nil {
		countTxs += 1
	}

	return []*types.Transaction{signedTx}, err
}
