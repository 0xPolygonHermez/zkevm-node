package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/log"
	ERC20 "github.com/hermeznetwork/hermez-core/test/contracts/bin/ERC20"
	WETH "github.com/hermeznetwork/hermez-core/test/contracts/bin/WETH"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/uniswap/v2/core/UniswapV2Factory"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/uniswap/v2/core/UniswapV2Pair"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/uniswap/v2/interface/UniswapInterfaceMulticall"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/uniswap/v2/periphery/UniswapV2Router02"
)

const (
	networkURL          = "http://localhost:8123"
	pk                  = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	txMinedTimeoutLimit = 60 * time.Second
)

func main() {
	ctx := context.Background()

	// if you want to test using goerli network
	// pk := "" // replace by your account in goerli with ETH balance to send the transactions
	// networkURL := "" // replace by your goerli infura project url

	log.Infof("connecting to %v", networkURL)
	client, err := ethclient.Dial(networkURL)
	chkErr(err)
	log.Infof("connected")

	chainID, err := client.ChainID(ctx)
	chkErr(err)
	log.Infof("chainID: %v", chainID)

	auth := getAuth(ctx, client, pk)
	fmt.Println()

	// Deploy ERC20 Tokens to be swapped
	aCoinAddr, aCoin := deployERC20(auth, client, "A COIN", "ACO")
	fmt.Println()
	bCoinAddr, bCoin := deployERC20(auth, client, "B COIN", "BCO")
	fmt.Println()
	cCoinAddr, cCoin := deployERC20(auth, client, "C COIN", "CCO")
	fmt.Println()

	// Deploy wETH Token, it's required by uniswap to swap ETH by tokens
	wEthAddr, tx, wethSC, err := WETH.DeployWETH(auth, client)
	chkErr(err)
	_, err = waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)
	log.Debugf("wEth SC tx: %v", tx.Hash().Hex())
	log.Debugf("wEth SC addr: %v", wEthAddr.Hex())
	fmt.Println()

	// Deploy Uniswap Factory
	factoryAddr, tx, factory, err := UniswapV2Factory.DeployUniswapV2Factory(auth, client, auth.From)
	chkErr(err)
	_, err = waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)
	log.Debugf("Uniswap Factory SC tx: %v", tx.Hash().Hex())
	log.Debugf("Uniswap Factory SC addr: %v", factoryAddr.Hex())
	fmt.Println()

	// Deploy Uniswap Router
	routerAddr, tx, router, err := UniswapV2Router02.DeployUniswapV2Router02(auth, client, factoryAddr, wEthAddr)
	chkErr(err)
	_, err = waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)
	log.Debugf("Uniswap Router SC tx: %v", tx.Hash().Hex())
	log.Debugf("Uniswap Router SC addr: %v", routerAddr.Hex())
	fmt.Println()

	// Deploy Uniswap Interface Multicall
	multicallAddr, tx, _, err := UniswapInterfaceMulticall.DeployUniswapInterfaceMulticall(auth, client)
	chkErr(err)
	_, err = waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)
	log.Debugf("Uniswap Interface Multicall SC tx: %v", tx.Hash().Hex())
	log.Debugf("Uniswap Interface Multicall SC addr: %v", multicallAddr.Hex())
	fmt.Println()

	// Mint balance to tokens
	aMintAmount := "500000000000000000000"
	tx = mintERC20(auth, client, aCoin, aMintAmount)
	log.Debugf("Mint A Coin tx: %v", tx.Hash().Hex())
	fmt.Println()
	bMintAmount := "600000000000000000000"
	tx = mintERC20(auth, client, bCoin, bMintAmount)
	log.Debugf("Mint B Coin tx: %v", tx.Hash().Hex())
	fmt.Println()
	cMintAmount := "700000000000000000000"
	tx = mintERC20(auth, client, cCoin, cMintAmount)
	log.Debugf("Mint C Coin tx: %v", tx.Hash().Hex())
	fmt.Println()

	// wrapping eth
	wethDepositAmount, _ := big.NewInt(0).SetString("200000000000000000000", encoding.Base10)
	log.Debugf("Depositing %v ETH for account %v on token wEth", wethDepositAmount.Text(encoding.Base10), auth.From)
	wAuth := getAuth(ctx, client, pk)
	wAuth.Value = wethDepositAmount
	tx, err = wethSC.Deposit(auth)
	chkErr(err)
	_, err = waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)

	// Add allowance
	approveERC20(auth, client, aCoin, routerAddr, aMintAmount)
	fmt.Println()
	approveERC20(auth, client, bCoin, routerAddr, bMintAmount)
	fmt.Println()
	approveERC20(auth, client, cCoin, routerAddr, cMintAmount)
	fmt.Println()
	approveERC20(auth, client, wethSC, routerAddr, wethDepositAmount.Text(encoding.Base10))
	fmt.Println()

	liquidityAmount := "100000000000000000000"

	// Add liquidity to the pool
	tx = addLiquidity(auth, client, router, aCoinAddr, bCoinAddr, liquidityAmount)
	log.Debugf("Add Liquidity to Pair A <-> B tx: %v", tx.Hash().Hex())
	fmt.Println()

	tx = addLiquidity(auth, client, router, bCoinAddr, cCoinAddr, liquidityAmount)
	log.Debugf("Add Liquidity to Pair B <-> C tx: %v", tx.Hash().Hex())
	fmt.Println()

	// Execute swaps
	log.Debugf("Swapping A <-> B")
	pairAddr, err := factory.GetPair(nil, aCoinAddr, bCoinAddr)
	chkErr(err)
	log.Debugf("Swapping A <-> B pair: %v", pairAddr.Hex())
	pairSC, err := UniswapV2Pair.NewUniswapV2Pair(pairAddr, client)
	chkErr(err)

	pairReserves, err := pairSC.GetReserves(nil)
	chkErr(err)
	log.Debugf("Swapping A <-> B reserves: 0: %v 1: %v Block Timestamp: %v", pairReserves.Reserve0, pairReserves.Reserve1, pairReserves.BlockTimestampLast)

	const sourceAmount = 1000
	exactAmountIn := big.NewInt(sourceAmount)
	amountOut, err := router.GetAmountOut(nil, exactAmountIn, pairReserves.Reserve0, pairReserves.Reserve1)
	chkErr(err)

	tx, err = router.SwapExactTokensForTokens(auth, exactAmountIn, amountOut, []common.Address{aCoinAddr, bCoinAddr}, auth.From, getDeadline())
	chkErr(err)
	log.Debugf("Swapping A <-> B tx: %v", tx.Hash().Hex())
	_, err = waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)
}

func getAuth(ctx context.Context, client *ethclient.Client, pkHex string) *bind.TransactOpts {
	chainID, err := client.ChainID(ctx)
	chkErr(err)
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(pkHex, "0x"))
	chkErr(err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	chkErr(err)

	return auth
}

func waitTxToBeMined(client *ethclient.Client, hash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	log.Infof("waiting tx to be mined: %v", hash.Hex())
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return nil, errors.New("timeout exceed")
		}

		r, err := client.TransactionReceipt(context.Background(), hash)
		if errors.Is(err, ethereum.NotFound) {
			time.Sleep(1 * time.Second)
			continue
		}
		if err != nil {
			log.Errorf("Failed to get tx receipt: %v", err)
			return nil, err
		}

		if r.Status == types.ReceiptStatusFailed {
			log.Errorf("tx mined[FAILED]: %v", hash.Hex())
		} else {
			log.Infof("tx mined[SUCCESS]: %v", hash.Hex())
		}

		return r, nil
	}
}

func deployERC20(auth *bind.TransactOpts, client *ethclient.Client, name, symbol string) (common.Address, *ERC20.ERC20) {
	log.Debugf("Deploying ERC20 Token: [%v]%v", symbol, name)
	addr, tx, instance, err := ERC20.DeployERC20(auth, client, name, symbol)
	chkErr(err)
	_, err = waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)
	log.Debugf("%v SC tx: %v", name, tx.Hash().Hex())
	log.Debugf("%v SC addr: %v", name, addr.Hex())
	return addr, instance
}

func mintERC20(auth *bind.TransactOpts, client *ethclient.Client, erc20sc *ERC20.ERC20, amount string) *types.Transaction {
	name, err := erc20sc.Name(nil)
	chkErr(err)
	log.Debugf("Minting %v tokens for account %v on token %v", amount, auth.From, name)
	mintAmount, _ := big.NewInt(0).SetString(amount, encoding.Base10)
	tx, err := erc20sc.Mint(auth, mintAmount)
	chkErr(err)
	_, err = waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)
	return tx
}

func approveERC20(auth *bind.TransactOpts, client *ethclient.Client,
	sc interface {
		Name(opts *bind.CallOpts) (string, error)
		Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)
	},
	routerAddr common.Address,
	amount string) {
	name, err := sc.Name(nil)
	chkErr(err)

	a, _ := big.NewInt(0).SetString(amount, encoding.Base10)

	log.Debugf("Approving %v tokens to be used by the router for %v on behalf of account %v", a.Text(encoding.Base10), name, auth.From)
	tx, err := sc.Approve(auth, routerAddr, a)
	chkErr(err)
	_, err = waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)
	log.Debugf("Approval %v tx: %v", name, tx.Hash().Hex())
}

func addLiquidity(auth *bind.TransactOpts, client *ethclient.Client, router *UniswapV2Router02.UniswapV2Router02, tokenA, tokenB common.Address, amount string) *types.Transaction {
	a, _ := big.NewInt(0).SetString(amount, encoding.Base10)
	log.Debugf("Adding liquidity(%v) for tokens A: %v, B:%v, Recipient: %v", amount, tokenA.Hex(), tokenB.Hex(), auth.From.Hex())
	tx, err := router.AddLiquidity(auth, tokenA, tokenB, a, a, a, a, auth.From, getDeadline())
	chkErr(err)
	_, err = waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)
	return tx
}

func getDeadline() *big.Int {
	const deadLinelimit = 5 * time.Minute
	return big.NewInt(time.Now().UTC().Add(deadLinelimit).Unix())
}

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
