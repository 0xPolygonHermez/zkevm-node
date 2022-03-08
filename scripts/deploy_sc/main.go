package main

import (
	"context"
	"crypto/ecdsa"
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
)

const (
	networkURL = "http://localhost:8123"

	accPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	gasLimit      = 400000

	txMinedTimeoutLimit = 60 * time.Second

	// compiled using http://remix.ethereum.org
	// COMPILER: 0.8.7+commit.e28d00a7
	// OPTIMIZATION: disabled
	// ../../test/contracts/emitLog.sol
	scHexBytes = "608060405234801561001057600080fd5b507f5e7df75d54e493185612379c616118a4c9ac802de621b010c96f74d22df4b30a60405160405180910390a160017f977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac6260405160405180910390a2600260017fbb6e4da744abea70325874159d52c1ad3e57babfae7c329a948e7dcb274deb0960405160405180910390a36003600260017f966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b760405160405180910390a46003600260017fe5562b12d9276c5c987df08afff7b1946f2d869236866ea2285c7e2e95685a6460046040516101039190610243565b60405180910390a46002600360047fe5562b12d9276c5c987df08afff7b1946f2d869236866ea2285c7e2e95685a6460016040516101419190610228565b60405180910390a46001600260037f966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b760405160405180910390a4600160027fbb6e4da744abea70325874159d52c1ad3e57babfae7c329a948e7dcb274deb0960405160405180910390a360017f977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac6260405160405180910390a27f5e7df75d54e493185612379c616118a4c9ac802de621b010c96f74d22df4b30a60405160405180910390a161028c565b61021381610268565b82525050565b6102228161027a565b82525050565b600060208201905061023d600083018461020a565b92915050565b60006020820190506102586000830184610219565b92915050565b6000819050919050565b60006102738261025e565b9050919050565b60006102858261025e565b9050919050565b603f8061029a6000396000f3fe6080604052600080fdfea2646970667358221220762c67d81efb5d60dba1d35e07b0924d0b098edb99abd3d76793806defeaabba64736f6c63430008070033"
)

// this function sends a transaction to deploy a smartcontract to the local network
func main() {
	ctx := context.Background()

	log.Infof("connecting to %v", networkURL)

	client, err := ethclient.Dial(networkURL)
	chkErr(err)

	log.Infof("connected")
	log.Infof("getting chainID")

	chainID, err := client.ChainID(ctx)
	chkErr(err)

	log.Infof("chainID: %v", chainID)
	log.Infof("reading private key...")

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(accPrivateKey, "0x"))
	chkErr(err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Infof("account: %v", fromAddress.Hex())
	log.Infof("reading nonce")
	nonce, err := client.NonceAt(ctx, fromAddress, nil)
	log.Infof("nonce: %v", nonce)
	chkErr(err)

	// we need to use this method to send `TO` field as `NULL`,
	// so the explorer can detect this is a smart contract creation
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       nil,
		Value:    new(big.Int),
		Gas:      gasLimit,
		GasPrice: new(big.Int).SetUint64(1),
		Data:     common.Hex2Bytes(scHexBytes),
	})

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	chkErr(err)

	signedTx, err := auth.Signer(auth.From, tx)
	chkErr(err)

	log.Infof("sending tx to deploy sc")

	err = client.SendTransaction(ctx, signedTx)
	chkErr(err)

	log.Infof("tx sent: %v", signedTx.Hash().Hex())
	log.Infof("waiting tx to be mined")

	err = waitTxToBeMined(client, signedTx.Hash(), txMinedTimeoutLimit)
	chkErr(err)

	log.Infof("SC Deployed")
}

func waitTxToBeMined(client *ethclient.Client, hash common.Hash, timeout time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return errors.New("timeout exceed")
		}

		r, err := client.TransactionReceipt(context.Background(), hash)
		if errors.Is(err, ethereum.NotFound) {
			log.Infof("Receipt not found yet, retrying...")
			time.Sleep(1 * time.Second)
			continue
		}
		if err != nil {
			log.Errorf("Failed to get tx receipt: %v", err)
			return err
		}

		if r.Status == types.ReceiptStatusFailed {
			return fmt.Errorf("transaction has failed: %s", string(r.PostState))
		}

		return nil
	}
}

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
