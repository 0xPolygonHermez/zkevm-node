package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/log"
	erc20 "github.com/hermeznetwork/hermez-core/test/contracts/bin/ERC20"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/uniswap/v2/core/UniswapV2Factory"
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

	// Deploy ERC20 Tokens to be swapped
	aCoinAddr := deployERC20(auth, client, "A COIN", "ACO")
	bCoinAddr := deployERC20(auth, client, "B COIN", "BCO")
	cCoinAddr := deployERC20(auth, client, "C COIN", "CCO")

	// Deploy wETH Token, it's required by uniswap to swap ETH by tokens
	wethAddr, tx, _, err := weth.DeployWeth(auth, client)
	chkErr(err)
	waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	log.Debugf("wEth SC tx: %v", tx.Hash().Hex())
	log.Debugf("wEth SC addr: %v", wethAddr.Hex())

	// Deploy Uniswap Factory
	factoryAddr, tx, factory, err := UniswapV2Factory.DeployUniswapV2Factory(auth, client, auth.From)
	chkErr(err)
	waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	log.Debugf("Uniswap Factory SC tx: %v", tx.Hash().Hex())
	log.Debugf("Uniswap Factory SC addr: %v", factoryAddr.Hex())

	// Create uniswap pairs to allow tokens to be swapped
	tx = createPair(factory, auth, client, aCoinAddr, bCoinAddr)
	log.Debugf("Uniswap Pair A <-> B tx: %v", tx.Hash().Hex())
	tx = createPair(factory, auth, client, bCoinAddr, cCoinAddr)
	log.Debugf("Uniswap Pair B <-> C tx: %v", tx.Hash().Hex())

	// Deploy Uniswap Router
	routerAddr, tx, _, err := UniswapV2Router02.DeployUniswapV2Router02(auth, client, factoryAddr, wethAddr)
	chkErr(err)
	waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	log.Debugf("Uniswap Router SC tx: %v", tx.Hash().Hex())
	log.Debugf("Uniswap Router SC addr: %v", routerAddr.Hex())

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

func deployERC20(auth *bind.TransactOpts, client *ethclient.Client, name, symbol string) common.Address {
	addr, tx, _, err := erc20.DeployErc20(auth, client, name, symbol)
	chkErr(err)
	_, err = waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	chkErr(err)
	log.Debugf("%v SC tx: %v", name, tx.Hash().Hex())
	log.Debugf("%v SC addr: %v", name, addr.Hex())

	return addr
}

func createPair(factory *UniswapV2Factory.UniswapV2Factory, auth *bind.TransactOpts, client *ethclient.Client, tokenA, tokenB common.Address) *types.Transaction {
	tx, err := factory.CreatePair(auth, tokenA, tokenB)
	chkErr(err)
	waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	return tx
}

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
