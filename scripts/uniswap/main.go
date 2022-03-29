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
	"github.com/hermeznetwork/hermez-core/log"
	erc20 "github.com/hermeznetwork/hermez-core/test/contracts/bin/ERC20"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/uniswap/v2/core/UniswapV2Factory"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/uniswap/v2/periphery/UniswapV2Router02"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/weth"
)

const (
	networkURL = "http://localhost:8123"

	accPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	txMinedTimeoutLimit = 60 * time.Second
)

func main() {
	ctx := context.Background()

	log.Infof("connecting to %v", networkURL)
	client, err := ethclient.Dial(networkURL)
	chkErr(err)
	log.Infof("connected")

	auth := getAuth(ctx, client)

	// Deploy ERC20 Tokens to be swapped
	aCoinAddr, tx := deployERC20(auth, client, "A COIN", "ACO")
	log.Debugf("A Coin SC tx: %v", tx.Hash().Hex())
	log.Debugf("A Coin SC addr: %v", aCoinAddr.Hex())

	bCoinAddr, tx := deployERC20(auth, client, "B COIN", "BCO")
	log.Debugf("B Coin SC tx: %v", tx.Hash().Hex())
	log.Debugf("B Coin SC addr: %v", bCoinAddr.Hex())

	cCoinAddr, tx := deployERC20(auth, client, "C COIN", "CCO")
	log.Debugf("C Coin SC tx: %v", tx.Hash().Hex())
	log.Debugf("C Coin SC addr: %v", cCoinAddr.Hex())

	// Deploy wETH Token
	wethAddr, tx, _, err := weth.DeployWeth(auth, client)
	chkErr(err)
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
	_, tx, _, err = UniswapV2Router02.DeployUniswapV2Router02(auth, client, factoryAddr, wethAddr)
	chkErr(err)
	waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	log.Debugf("Uniswap Factory SC tx: %v", tx.Hash().Hex())
	log.Debugf("Uniswap Factory SC addr: %v", factoryAddr.Hex())
}

func getAuth(ctx context.Context, client *ethclient.Client) *bind.TransactOpts {
	chainID, err := client.ChainID(ctx)
	chkErr(err)
	log.Infof("chainID: %v", chainID)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(accPrivateKey, "0x"))
	chkErr(err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	chkErr(err)

	return auth
}

func sendTxToDeploySC(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, scHexByte string, gasLimit uint64) common.Address {
	log.Infof("reading nonce for account: %v", auth.From.Hex())
	nonce, err := client.NonceAt(ctx, auth.From, nil)
	log.Infof("nonce: %v", nonce)
	chkErr(err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       nil,
		Value:    new(big.Int),
		Gas:      gasLimit,
		GasPrice: new(big.Int).SetUint64(1),
		Data:     common.Hex2Bytes(scHexByte),
	})

	signedTx, err := auth.Signer(auth.From, tx)
	chkErr(err)

	log.Infof("sending tx to deploy sc")

	err = client.SendTransaction(ctx, signedTx)
	chkErr(err)
	log.Infof("tx sent: %v", signedTx.Hash().Hex())

	r, err := waitTxToBeMined(client, signedTx.Hash(), txMinedTimeoutLimit)
	chkErr(err)

	log.Infof("SC Deployed to address: %v", r.ContractAddress.Hex())

	return r.ContractAddress
}

func waitTxToBeMined(client *ethclient.Client, hash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	log.Infof("waiting tx to be mined")
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return nil, errors.New("timeout exceed")
		}

		r, err := client.TransactionReceipt(context.Background(), hash)
		if errors.Is(err, ethereum.NotFound) {
			log.Infof("Receipt not found yet, retrying...")
			time.Sleep(1 * time.Second)
			continue
		}
		if err != nil {
			log.Errorf("Failed to get tx receipt: %v", err)
			return nil, err
		}

		if r.Status == types.ReceiptStatusFailed {
			return nil, fmt.Errorf("transaction has failed: %s", hash.Hex())
		}

		return r, nil
	}
}

func deployERC20(auth *bind.TransactOpts, client *ethclient.Client, name, symbol string) (common.Address, *types.Transaction) {
	addr, tx, _, err := erc20.DeployErc20(auth, client, name, symbol)
	chkErr(err)
	waitTxToBeMined(client, tx.Hash(), txMinedTimeoutLimit)
	return addr, tx
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
