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

	scHexBytes = "608060405234801561001057600080fd5b507ff24055cc366c1d53b75688e6ba6eb419d022fb5f4956c5d66a6ef54adbf3270a60405161003e9061006e565b60405180910390a16100c8565b600061005860118361008e565b91506100638261009f565b602082019050919050565b600060208201905081810360008301526100878161004b565b9050919050565b600082825260208201905092915050565b7f636f6e7472616374206372656174656421000000000000000000000000000000600082015250565b610227806100d76000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b6040516100509190610187565b60405180910390f35b610073600480360381019061006e919061014b565b61007e565b005b60008054905090565b600080549050816000819055507f4d4ebfeb041dd7cbdf868d9af07944bc3ea69fc4d89afe3f5b2d293c4340d222816040516100ba9190610187565b60405180910390a17fac3e966f295f2d5312f973dc6d42f30a6dc1c1f76ab8ee91cc8ca5dad1fa60fd826040516100f19190610187565b60405180910390a17f2db947ef788961acc438340dbcb4e242f80d026b621b7c98ee30619950390382818360405161012a9291906101a2565b60405180910390a15050565b600081359050610145816101da565b92915050565b600060208284031215610161576101606101d5565b5b600061016f84828501610136565b91505092915050565b610181816101cb565b82525050565b600060208201905061019c6000830184610178565b92915050565b60006040820190506101b76000830185610178565b6101c46020830184610178565b9392505050565b6000819050919050565b600080fd5b6101e3816101cb565b81146101ee57600080fd5b5056fea26469706673582212200992efa9c7a486665e5086be4b1e4de6edac58bf642d3392c5b8eaf3d6f7913b64736f6c63430008070033"
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
