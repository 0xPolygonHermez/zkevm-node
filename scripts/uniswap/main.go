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
	"github.com/hermeznetwork/hermez-core/test/testutils"
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

	deploySC(ctx, client, auth, "uniswap/v2/core/UniswapV2ERC20.bin", 1200000)
	deploySC(ctx, client, auth, "uniswap/v2/core/UniswapV2Factory.bin", 1200000)
	deploySC(ctx, client, auth, "uniswap/v2/core/UniswapV2Pair.bin", 1200000)
	deploySC(ctx, client, auth, "uniswap/v2/periphery/UniswapV2Migrator.bin", 1200000)
	deploySC(ctx, client, auth, "uniswap/v2/periphery/UniswapV2Router01.bin", 1200000)
	deploySC(ctx, client, auth, "uniswap/v2/periphery/UniswapV2Router02.bin", 1200000)
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

func deploySC(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, contractPath string, gasLimit uint64) common.Address {
	bytes, err := testutils.ReadBytecode(contractPath)
	chkErr(err)

	addr := sendTxToDeploySC(ctx, client, auth, bytes, 1200000)
	log.Debugf("%v: %v", contractPath, addr.Hex())

	return addr
}

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
