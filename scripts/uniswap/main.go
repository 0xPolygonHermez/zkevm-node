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
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/erc20"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/uniswap/v2/core/UniswapV2Factory"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/uniswap/v2/interface/UniswapInterfaceMulticall"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/uniswap/v2/periphery/UniswapV2Router02"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/weth"
)

const (
	networkURL          = "http://localhost:8123"
	pk                  = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	txMinedTimeoutLimit = 60 * time.Second
)

func main() {
	ctx := context.Background()

	log.Infof("connecting to %v", networkURL)
	client, err := ethclient.Dial(networkURL)
	chkErr(err)
	log.Infof("connected")

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
	wEthAddr, tx, _, err := weth.DeployWeth(auth, client)
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

	// Create uniswap pairs to allow tokens to be swapped
	tx, wEthACoinPairAddr := createPair(auth, client, factory, wEthAddr, aCoinAddr)
	log.Debugf("Uniswap Pair wEth <-> B tx: %v", tx.Hash().Hex())
	log.Debugf("Uniswap Pair wEth <-> B addr: %v", wEthACoinPairAddr.Hash().Hex())
	fmt.Println()
	tx, aCoinBCoinPairAddr := createPair(auth, client, factory, aCoinAddr, bCoinAddr)
	log.Debugf("Uniswap Pair A <-> B tx: %v", tx.Hash().Hex())
	log.Debugf("Uniswap Pair A <-> B addr: %v", aCoinBCoinPairAddr.Hash().Hex())
	fmt.Println()
	tx, bCoinCCoinPairAddr := createPair(auth, client, factory, bCoinAddr, cCoinAddr)
	log.Debugf("Uniswap Pair B <-> C tx: %v", tx.Hash().Hex())
	log.Debugf("Uniswap Pair B <-> C addr: %v", bCoinCCoinPairAddr.Hash().Hex())
	fmt.Println()
	tx, cCoinWEthPairAddr := createPair(auth, client, factory, cCoinAddr, wEthAddr)
	log.Debugf("Uniswap Pair C <-> wEth tx: %v", tx.Hash().Hex())
	log.Debugf("Uniswap Pair C <-> wEth addr: %v", cCoinWEthPairAddr.Hash().Hex())
	fmt.Println()

	// Mint balance to tokens
	tx = mintERC20(auth, client, aCoin, "500000000000000000000")
	log.Debugf("Mint A Coin tx: %v", tx.Hash().Hex())
	fmt.Println()
	tx = mintERC20(auth, client, bCoin, "600000000000000000000")
	log.Debugf("Mint B Coin tx: %v", tx.Hash().Hex())
	fmt.Println()
	tx = mintERC20(auth, client, cCoin, "700000000000000000000")
	log.Debugf("Mint C Coin tx: %v", tx.Hash().Hex())
	fmt.Println()

	// Add liquidity to the pool
	tx = addLiquidity(auth, client, router, aCoinAddr, bCoinAddr, aCoinBCoinPairAddr)
	log.Debugf("Add Liquidity to Pair A <-> B tx: %v", tx.Hash().Hex())
	fmt.Println()
	tx = addLiquidity(auth, client, router, bCoinAddr, cCoinAddr, bCoinCCoinPairAddr)
	log.Debugf("Add Liquidity to Pair B <-> C tx: %v", tx.Hash().Hex())
	fmt.Println()

	// Execute swaps
}

func getAuth(ctx context.Context, client *ethclient.Client, pkHex string) *bind.TransactOpts {
	chainID, err := client.ChainID(ctx)
	chkErr(err)
	log.Infof("chainID: %v", chainID)

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
			//log.Infof("Receipt not found yet, retrying...")
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

func deployERC20(auth *bind.TransactOpts, client *ethclient.Client, name, symbol string) (common.Address, *erc20.Erc20) {
	log.Debugf("Deploying ERC20 Token: [%v]%v", symbol, name)
	addr, tx, instance, err := erc20.DeployErc20(auth, client, name, symbol)
	chkErr(err)
	_, err = waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)
	log.Debugf("%v SC tx: %v", name, tx.Hash().Hex())
	log.Debugf("%v SC addr: %v", name, addr.Hex())
	return addr, instance
}

func mintERC20(auth *bind.TransactOpts, client *ethclient.Client, erc20sc *erc20.Erc20, amount string) *types.Transaction {
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

func createPair(auth *bind.TransactOpts, client *ethclient.Client, factory *UniswapV2Factory.UniswapV2Factory, tokenA, tokenB common.Address) (*types.Transaction, common.Address) {
	tx, err := factory.CreatePair(auth, tokenA, tokenB)
	chkErr(err)
	r, err := waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)

	pair := common.Address{}
	if r.Status == types.ReceiptStatusSuccessful {
		pair, err = factory.GetPair(nil, tokenA, tokenB)
		chkErr(err)
	}

	return tx, pair
}

func addLiquidity(auth *bind.TransactOpts, client *ethclient.Client, router *UniswapV2Router02.UniswapV2Router02, tokenA, tokenB common.Address, pairAddr common.Address) *types.Transaction {
	amount, _ := big.NewInt(0).SetString("100000000000000000000", encoding.Base10)
	tx, err := router.AddLiquidity(auth, tokenA, tokenB, amount, amount, amount, amount, pairAddr, big.NewInt(0))
	chkErr(err)
	_, err = waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)
	return tx
}

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
