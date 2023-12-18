package pkg

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/ERC20"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/WETH"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/uniswap/v2/core/UniswapV2Factory"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/uniswap/v2/interface/UniswapInterfaceMulticall"
	"github.com/0xPolygonHermez/zkevm-node/test/contracts/bin/uniswap/v2/periphery/UniswapV2Router02"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	txTimeout = 60 * time.Second
)

var (
	executedTransctionsCount uint64 = 0
)

func DeployContractsAndAddLiquidity(client *ethclient.Client, auth *bind.TransactOpts) Deployments {
	ctx := context.Background()
	fmt.Println()
	balance, err := client.BalanceAt(ctx, auth.From, nil)
	ChkErr(err)
	log.Debugf("ETH Balance for %v: %v", auth.From, balance)
	// Deploy ERC20 Tokens to be swapped
	aCoinAddr, aCoin := deployERC20(auth, client, "A COIN", "ACO")
	fmt.Println()
	bCoinAddr, bCoin := deployERC20(auth, client, "B COIN", "BCO")
	fmt.Println()
	cCoinAddr, cCoin := deployERC20(auth, client, "C COIN", "CCO")
	fmt.Println()
	// Deploy wETH Token, it's required by uniswap to swap ETH by tokens
	log.Debugf("Deploying wEth SC")
	wEthAddr, tx, wethSC, err := WETH.DeployWETH(auth, client)
	err = WaitForTransactionAndIncrementNonce(client, auth, err, ctx, tx)
	log.Debugf("wEth SC tx: %v", tx.Hash().Hex())
	log.Debugf("wEth SC addr: %v", wEthAddr.Hex())
	fmt.Println()
	// Deploy Uniswap Factory
	log.Debugf("Deploying Uniswap Factory")
	factoryAddr, tx, factory, err := UniswapV2Factory.DeployUniswapV2Factory(auth, client, auth.From)
	err = WaitForTransactionAndIncrementNonce(client, auth, err, ctx, tx)
	log.Debugf("Uniswap Factory SC tx: %v", tx.Hash().Hex())
	log.Debugf("Uniswap Factory SC addr: %v", factoryAddr.Hex())
	fmt.Println()
	// Deploy Uniswap Router
	log.Debugf("Deploying Uniswap Router")
	routerAddr, tx, router, err := UniswapV2Router02.DeployUniswapV2Router02(auth, client, factoryAddr, wEthAddr)
	err = WaitForTransactionAndIncrementNonce(client, auth, err, ctx, tx)
	log.Debugf("Uniswap Router SC tx: %v", tx.Hash().Hex())
	log.Debugf("Uniswap Router SC addr: %v", routerAddr.Hex())
	fmt.Println()
	// Deploy Uniswap Interface Multicall
	log.Debugf("Deploying Uniswap Multicall")
	multicallAddr, tx, _, err := UniswapInterfaceMulticall.DeployUniswapInterfaceMulticall(auth, client)
	err = WaitForTransactionAndIncrementNonce(client, auth, err, ctx, tx)
	log.Debugf("Uniswap Interface Multicall SC tx: %v", tx.Hash().Hex())
	log.Debugf("Uniswap Interface Multicall SC addr: %v", multicallAddr.Hex())
	fmt.Println()
	// Mint balance to tokens
	log.Debugf("Minting ERC20 Tokens")
	aMintAmount := "1000000000000000000000"
	tx = mintERC20(auth, client, aCoin, aMintAmount)
	log.Debugf("Mint A Coin tx: %v", tx.Hash().Hex())
	fmt.Println()
	bMintAmount := "1000000000000000000000"
	tx = mintERC20(auth, client, bCoin, bMintAmount)
	log.Debugf("Mint B Coin tx: %v", tx.Hash().Hex())
	fmt.Println()
	cMintAmount := "1000000000000000000000"
	tx = mintERC20(auth, client, cCoin, cMintAmount)
	log.Debugf("Mint C Coin tx: %v", tx.Hash().Hex())
	fmt.Println()
	// wrapping eth
	wethDepositoAmount := "0000000000000000"
	log.Debugf("Depositing %v ETH for account %v on token wEth", wethDepositoAmount, auth.From)
	auth.Value, _ = big.NewInt(0).SetString(wethDepositoAmount, encoding.Base10)
	tx, err = wethSC.Deposit(auth)
	err = WaitForTransactionAndIncrementNonce(client, auth, err, ctx, tx)
	value, err := aCoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("before allowance aCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	value, err = bCoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("before allowance bCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	value, err = cCoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("before allowance cCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	// Add allowance
	approveERC20(auth, client, aCoin, routerAddr, aMintAmount)
	fmt.Println()
	approveERC20(auth, client, bCoin, routerAddr, bMintAmount)
	fmt.Println()
	approveERC20(auth, client, cCoin, routerAddr, cMintAmount)
	fmt.Println()
	approveERC20(auth, client, wethSC, routerAddr, wethDepositoAmount)
	fmt.Println()
	const liquidityAmount = "10000000000000000000"
	value, err = aCoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("before adding liquidity A, B aCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	value, err = bCoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("before adding liquidity A, B bCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	value, err = cCoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("before adding liquidity A, B cCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	// Add liquidity to the pool
	tx = addLiquidity(auth, client, router, aCoinAddr, bCoinAddr, liquidityAmount)
	log.Debugf("Add Liquidity to Pair A <-> B tx: %v", tx.Hash().Hex())
	fmt.Println()
	value, err = aCoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("before adding liquidity B, C aCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	value, err = bCoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("before adding liquidity B, C bCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	value, err = cCoin.BalanceOf(&bind.CallOpts{}, auth.From)
	ChkErr(err)
	log.Debugf("before adding liquidity B, C cCoin.balanceOf[%s]: %d", auth.From.Hex(), value)
	tx = addLiquidity(auth, client, router, bCoinAddr, cCoinAddr, liquidityAmount)
	log.Debugf("Add Liquidity to Pair B <-> C tx: %v", tx.Hash().Hex())
	fmt.Println()

	return Deployments{
		ACoin:     aCoin,
		ACoinAddr: aCoinAddr,
		BCoin:     bCoin,
		BCoinAddr: bCoinAddr,
		CCoin:     cCoin,
		CCoinAddr: cCoinAddr,
		Router:    router,
		Factory:   factory,
	}
}

func WaitForTransactionAndIncrementNonce(l2Client *ethclient.Client, auth *bind.TransactOpts, err error, ctx context.Context, tx *types.Transaction) error {
	ChkErr(err)
	err = operations.WaitTxToBeMined(ctx, l2Client, tx, txTimeout)
	ChkErr(err)
	executedTransctionsCount++
	auth.Nonce = nil
	auth.Value = nil

	return err
}

func GetAuth(ctx context.Context, client *ethclient.Client, pkHex string) *bind.TransactOpts {
	chainID, err := client.ChainID(ctx)
	ChkErr(err)
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(pkHex, "0x"))
	ChkErr(err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	ChkErr(err)
	senderNonce, err := client.PendingNonceAt(ctx, auth.From)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(senderNonce))
	return auth
}

func ChkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetExecutedTransactionsCount() uint64 {
	return executedTransctionsCount
}

func deployERC20(auth *bind.TransactOpts, client *ethclient.Client, name, symbol string) (common.Address, *ERC20.ERC20) {
	ctx := context.Background()
	log.Debugf("Deploying ERC20 Token: [%v]%v", symbol, name)
	addr, tx, instance, err := ERC20.DeployERC20(auth, client, name, symbol)
	err = WaitForTransactionAndIncrementNonce(client, auth, err, ctx, tx)
	log.Debugf("%v SC tx: %v", name, tx.Hash().Hex())
	log.Debugf("%v SC addr: %v", name, addr.Hex())
	return addr, instance
}

func mintERC20(auth *bind.TransactOpts, client *ethclient.Client, erc20sc *ERC20.ERC20, amount string) *types.Transaction {
	ctx := context.Background()
	name, err := erc20sc.Name(nil)
	ChkErr(err)
	log.Debugf("Minting %v tokens for account %v on token %v", amount, auth.From, name)
	mintAmount, _ := big.NewInt(0).SetString(amount, encoding.Base10)
	tx, err := erc20sc.Mint(auth, mintAmount)
	err = WaitForTransactionAndIncrementNonce(client, auth, err, ctx, tx)
	return tx
}

func approveERC20(auth *bind.TransactOpts, client *ethclient.Client,
	sc interface {
		Name(opts *bind.CallOpts) (string, error)
		Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)
	},
	routerAddr common.Address,
	amount string) {
	ctx := context.Background()
	name, err := sc.Name(nil)
	ChkErr(err)
	a, _ := big.NewInt(0).SetString(amount, encoding.Base10)
	log.Debugf("Approving %v tokens to be used by the router for %v on behalf of account %v", a.Text(encoding.Base10), name, auth.From)
	tx, err := sc.Approve(auth, routerAddr, a)
	err = WaitForTransactionAndIncrementNonce(client, auth, err, ctx, tx)
	log.Debugf("Approval %v tx: %v", name, tx.Hash().Hex())
}

func addLiquidity(auth *bind.TransactOpts, client *ethclient.Client, router *UniswapV2Router02.UniswapV2Router02, tokenA, tokenB common.Address, amount string) *types.Transaction {
	ctx := context.Background()
	a, _ := big.NewInt(0).SetString(amount, encoding.Base10)
	log.Debugf("Adding liquidity(%v) for tokens A: %v, B:%v, Recipient: %v", amount, tokenA.Hex(), tokenB.Hex(), auth.From.Hex())
	tx, err := router.AddLiquidity(auth, tokenA, tokenB, a, a, a, a, auth.From, getDeadline())
	err = WaitForTransactionAndIncrementNonce(client, auth, err, ctx, tx)
	return tx
}

func getDeadline() *big.Int {
	const deadLinelimit = 5 * time.Minute
	return big.NewInt(time.Now().UTC().Add(deadLinelimit).Unix())
}
