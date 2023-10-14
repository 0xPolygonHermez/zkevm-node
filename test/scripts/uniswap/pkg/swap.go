package pkg

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/uniswap/v2/core/UniswapV2Factory"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/uniswap/v2/core/UniswapV2Pair"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/uniswap/v2/periphery/UniswapV2Router02"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SwapTokens(client *ethclient.Client, auth *bind.TransactOpts, deployments Deployments) []*types.Transaction {
	transactions := make([]*types.Transaction, 0, 2)
	// Execute swaps
	const swapExactAmountInNumber = 10
	swapExactAmountIn := big.NewInt(swapExactAmountInNumber)
	swapExactAmountIn2 := big.NewInt(swapExactAmountInNumber - 1)
	value, err := deployments.ACoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("before first swap aCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	value, err = deployments.BCoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("before first swap bCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	value, err = deployments.CCoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("before first swap cCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	log.Debugf("Swaping tokens from A <-> B")
	res := SwapExactTokensForTokens(auth, client, deployments.Factory, deployments.Router, deployments.ACoinAddr, deployments.BCoinAddr, swapExactAmountIn)
	transactions = append(transactions, res...)
	fmt.Println()
	value, err = deployments.ACoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("after first swap aCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	value, err = deployments.BCoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("after first swap bCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	value, err = deployments.CCoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("after first swap cCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	log.Debugf("Swaping tokens from B <-> C")
	res = SwapExactTokensForTokens(auth, client, deployments.Factory, deployments.Router, deployments.BCoinAddr, deployments.CCoinAddr, swapExactAmountIn2)
	transactions = append(transactions, res...)
	fmt.Println()

	return transactions
}

func SwapExactTokensForTokens(auth *bind.TransactOpts, client *ethclient.Client,
	factory *UniswapV2Factory.UniswapV2Factory, router *UniswapV2Router02.UniswapV2Router02,
	tokenA, tokenB common.Address, exactAmountIn *big.Int) []*types.Transaction {
	ctx := context.Background()
	logPrefix := fmt.Sprintf("SwapExactTokensForTokens %v <-> %v", tokenA.Hex(), tokenB.Hex())
	pairAddr, err := factory.GetPair(nil, tokenA, tokenB)
	ChkErr(err)
	log.Debug(logPrefix, " pair: ", pairAddr.Hex())
	pairSC, err := UniswapV2Pair.NewUniswapV2Pair(pairAddr, client)
	ChkErr(err)
	pairReserves, err := pairSC.GetReserves(nil)
	ChkErr(err)
	log.Debug(logPrefix, " reserves 0: ", pairReserves.Reserve0, " 1: ", pairReserves.Reserve1, " Block Timestamp: ", pairReserves.BlockTimestampLast)
	amountOut, err := router.GetAmountOut(nil, exactAmountIn, pairReserves.Reserve0, pairReserves.Reserve1)
	ChkErr(err)
	log.Debug(logPrefix, " exactAmountIn: ", exactAmountIn, " amountOut: ", amountOut)
	tx, err := router.SwapExactTokensForTokens(auth, exactAmountIn, amountOut, []common.Address{tokenA, tokenB}, auth.From, getDeadline())
	err = WaitForTransactionAndIncrementNonce(client, auth, err, ctx, tx)
	return []*types.Transaction{tx}
}
