package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/test/operations"
)

const (
	l1NetworkURL = "http://localhost:8545"
	l2NetworkURL = "http://localhost:8123"

	poeAddress            = "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"
	bridgeAddress         = "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
	maticTokenAddress     = "0x5FbDB2315678afecb367f032d93F642f64180aa3" //nolint:gosec
	globalExitRootAddress = "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0"

	l1AccHexAddress    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	l1AccHexPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	sequencerAddress    = "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"
	sequencerPrivateKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"
)

// TestStateTransition tests state transitions using the vector
func main() {
	ctx := context.Background()

	// Eth client
	fmt.Println("Connecting to l1")
	client, err := ethclient.Dial(l1NetworkURL)
	checkErr(err)

	// Get network chain id
	fmt.Println("Getting chainID")
	chainID, err := client.NetworkID(ctx)
	checkErr(err)

	// Preparing l1 acc info
	fmt.Println("Preparing authorization")
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(l1AccHexPrivateKey, "0x"))
	checkErr(err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	checkErr(err)

	// Getting l1 info
	fmt.Println("Getting L1 info")
	gasPrice, err := client.SuggestGasPrice(ctx)
	checkErr(err)

	// Send some Ether from l1Acc to sequencer acc
	fmt.Println("Transferring ETH to the sequencer")
	fromAddress := common.HexToAddress(l1AccHexAddress)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	checkErr(err)
	const gasLimit = 21000
	toAddress := common.HexToAddress(sequencerAddress)
	ethAmount, _ := big.NewInt(0).SetString("100000000000000000000", encoding.Base10)
	tx := types.NewTransaction(nonce, toAddress, ethAmount, gasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(auth.From, tx)
	checkErr(err)
	err = client.SendTransaction(ctx, signedTx)
	checkErr(err)

	// Wait eth transfer to be mined
	fmt.Println("Waiting tx to be mined")
	const txETHTransferTimeout = 5 * time.Second
	err = waitTxToBeMined(ctx, client, signedTx.Hash(), txETHTransferTimeout)
	checkErr(err)

	// Create matic maticTokenSC sc instance
	fmt.Println("Loading Matic token SC instance")
	maticTokenSC, err := operations.NewToken(common.HexToAddress(maticTokenAddress), client)
	checkErr(err)

	// Send matic to sequencer
	fmt.Println("Transferring MATIC tokens to sequencer")
	maticAmount, _ := big.NewInt(0).SetString("100000000000000000000000", encoding.Base10)
	tx, err = maticTokenSC.Transfer(auth, toAddress, maticAmount)
	checkErr(err)

	// wait matic transfer to be mined
	fmt.Println("Waiting tx to be mined")
	const txMaticTransferTimeout = 5 * time.Second
	err = waitTxToBeMined(ctx, client, tx.Hash(), txMaticTransferTimeout)
	checkErr(err)

	// Create sequencer auth
	fmt.Println("Creating sequencer authorization")
	privateKey, err = crypto.HexToECDSA(strings.TrimPrefix(sequencerPrivateKey, "0x"))
	checkErr(err)
	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	checkErr(err)

	// approve tokens to be used by PoE SC on behalf of the sequencer
	fmt.Println("Approving tokens to be used by PoE on behalf of the sequencer")
	tx, err = maticTokenSC.Approve(auth, common.HexToAddress(poeAddress), maticAmount)
	checkErr(err)
	const txApprovalTimeout = 5 * time.Second
	err = waitTxToBeMined(ctx, client, tx.Hash(), txApprovalTimeout)
	checkErr(err)

	// Register the sequencer
	fmt.Println("Registering the sequencer")
	ethermanConfig := etherman.Config{
		URL: l1NetworkURL,
	}
	etherman, err := etherman.NewEtherman(ethermanConfig, auth, common.HexToAddress(poeAddress), common.HexToAddress(bridgeAddress), common.HexToAddress(maticTokenAddress), common.HexToAddress(globalExitRootAddress))
	checkErr(err)
	tx, err = etherman.RegisterSequencer(l2NetworkURL)
	checkErr(err)

	// Wait sequencer to be registered
	fmt.Println("waiting tx to be mined")
	const txRegistrationTimeout = 5 * time.Second
	err = waitTxToBeMined(ctx, client, tx.Hash(), txRegistrationTimeout)
	checkErr(err)

	fmt.Println("Network initialized properly")
}

func waitTxToBeMined(ctx context.Context, client *ethclient.Client, hash common.Hash, timeout time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return errors.New("timeout exceed")
		}

		time.Sleep(1 * time.Second)

		_, isPending, err := client.TransactionByHash(ctx, hash)
		if err == ethereum.NotFound {
			continue
		}

		if err != nil {
			return err
		}

		if !isPending {
			r, err := client.TransactionReceipt(ctx, hash)
			if err != nil {
				return err
			}

			if r.Status == types.ReceiptStatusFailed {
				return fmt.Errorf("transaction has failed: %s", string(r.PostState))
			}

			return nil
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
