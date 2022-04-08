package main

import (
	"context"
	"errors"
	"flag"
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
	"github.com/hermeznetwork/hermez-core/config"
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/etherman"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/bridge"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/matic"
	"github.com/hermeznetwork/hermez-core/etherman/smartcontracts/proofofefficiency"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/test/operations"
	"github.com/urfave/cli/v2"
)

const (
	l1NetworkURL = "http://localhost:8545"
	l2NetworkURL = "http://localhost:8123"

	poeAddress        = "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"
	maticTokenAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3" //nolint:gosec
	l1BridgeAddr      = "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
	l2BridgeAddr      = "0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988"

	l1AccHexAddress    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	l1AccHexPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	sequencerAddress    = "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"
	sequencerPrivateKey = "0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e"

	txTimeout = 5 * time.Second
)

// TestStateTransition tests state transitions using the vector
func main() {
	ctx := context.Background()

	app := cli.NewApp()
	var n string
	flag.StringVar(&n, "network", "local", "")
	context := cli.NewContext(app, flag.CommandLine, nil)

	config, err := config.Load(context)
	checkErr(err)

	// Eth client
	log.Infof("Connecting to l1")
	client, err := ethclient.Dial(l1NetworkURL)
	checkErr(err)

	// Get network chain id
	log.Infof("Getting chainID")
	chainID, err := client.NetworkID(ctx)
	checkErr(err)

	// Preparing l1 acc info
	log.Infof("Preparing authorization")
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(l1AccHexPrivateKey, "0x"))
	checkErr(err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	checkErr(err)

	// Getting l1 info
	log.Infof("Getting L1 info")
	gasPrice, err := client.SuggestGasPrice(ctx)
	checkErr(err)

	// Send some Ether from l1Acc to sequencer acc
	log.Infof("Transferring ETH to the sequencer")
	fromAddress := common.HexToAddress(l1AccHexAddress)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	checkErr(err)
	const gasLimit = 21000
	toAddress := common.HexToAddress(sequencerAddress)
	ethAmount, _ := big.NewInt(0).SetString("200000000000000000000", encoding.Base10)
	tx := types.NewTransaction(nonce, toAddress, ethAmount, gasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(auth.From, tx)
	checkErr(err)
	err = client.SendTransaction(ctx, signedTx)
	checkErr(err)

	// Wait eth transfer to be mined
	log.Infof("Waiting tx to be mined")
	err = waitTxToBeMined(ctx, client, signedTx.Hash(), txTimeout)
	checkErr(err)

	// Create matic maticTokenSC sc instance
	log.Infof("Loading Matic token SC instance")
	maticTokenSC, err := operations.NewToken(config.NetworkConfig.MaticAddr, client)
	checkErr(err)

	// Send matic to sequencer
	log.Infof("Transferring MATIC tokens to sequencer")
	maticAmount, _ := big.NewInt(0).SetString("200000000000000000000000", encoding.Base10)
	tx, err = maticTokenSC.Transfer(auth, toAddress, maticAmount)
	checkErr(err)

	// wait matic transfer to be mined
	log.Infof("Waiting tx to be mined")
	const txMaticTransferTimeout = 5 * time.Second
	err = waitTxToBeMined(ctx, client, tx.Hash(), txMaticTransferTimeout)
	checkErr(err)

	// Create sequencer auth
	log.Infof("Creating sequencer authorization")
	privateKey, err = crypto.HexToECDSA(strings.TrimPrefix(sequencerPrivateKey, "0x"))
	checkErr(err)
	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	checkErr(err)

	// approve tokens to be used by PoE SC on behalf of the sequencer
	log.Infof("Approving tokens to be used by PoE on behalf of the sequencer")
	tx, err = maticTokenSC.Approve(auth, config.NetworkConfig.PoEAddr, maticAmount)
	checkErr(err)
	const txApprovalTimeout = 5 * time.Second
	err = waitTxToBeMined(ctx, client, tx.Hash(), txApprovalTimeout)
	checkErr(err)

	// Register the sequencer
	log.Infof("Registering the sequencer")
	ethermanConfig := etherman.Config{
		URL: l1NetworkURL,
	}
	etherman, err := etherman.NewClient(ethermanConfig, auth, config.NetworkConfig.PoEAddr, config.NetworkConfig.MaticAddr)
	checkErr(err)
	tx, err = etherman.RegisterSequencer(l2NetworkURL)
	checkErr(err)

	// Wait sequencer to be registered
	log.Infof("waiting sequencer to be registered")
	const txRegistrationTimeout = 5 * time.Second
	err = waitTxToBeMined(ctx, client, tx.Hash(), txRegistrationTimeout)
	checkErr(err)
	log.Infof("sequencer registered")

	// Deposit funds to L2 via bridge
	log.Infof("Depositing funds to L2 via bridge")
	const destNetwork = uint32(1)
	depositAmount, _ := big.NewInt(0).SetString("1000000000000000000", encoding.Base10)
	ethAddr := common.Address{}
	sendL1Deposit(ctx, auth, client, ethAddr, depositAmount, destNetwork, &auth.From)

	// Proposing empty batch to trigger the l2 synchronization process
	forceBatchProposal(ctx, auth, client)

	// Claiming the funds deposited via bridge on L2
	// TODO: Get the values for deposit, smtProof and globalExitRoot
	deposit := &Deposit{}
	smtProof := [][32]byte{}
	globalExitRoot := &globalExitRoot{}
	sendL2Claim(ctx, auth, client, deposit, smtProof, globalExitRoot)

	log.Infof("Network initialized properly")
}

// sendL1Deposit sends a deposit from l1 to l2
func sendL1Deposit(ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, tokenAddr common.Address, amount *big.Int,
	destNetwork uint32, destAddr *common.Address,
) {
	emptyAddr := common.Address{}
	if tokenAddr == emptyAddr {
		auth.Value = amount
	}
	if destAddr == nil {
		destAddr = &auth.From
	}
	br, err := bridge.NewBridge(common.HexToAddress(l1BridgeAddr), client)
	checkErr(err)
	tx, err := br.Bridge(auth, tokenAddr, amount, destNetwork, *destAddr)
	checkErr(err)

	log.Infof("Waiting L1Deposit to be mined")
	err = waitTxToBeMined(ctx, client, tx.Hash(), txTimeout)
	checkErr(err)
}

func forceBatchProposal(ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client) {
	poeAddr := common.HexToAddress(poeAddress)
	poe, err := proofofefficiency.NewProofofefficiency(poeAddr, client)
	checkErr(err)

	maticAmount, err := poe.CalculateSequencerCollateral(&bind.CallOpts{Pending: false})
	checkErr(err)

	matic, err := matic.NewMatic(common.HexToAddress(maticTokenAddress), client)
	checkErr(err)
	txApprove, err := matic.Approve(auth, poeAddr, maticAmount)
	checkErr(err)

	err = waitTxToBeMined(ctx, client, txApprove.Hash(), txTimeout)
	checkErr(err)

	tx, err := poe.SendBatch(auth, []byte{}, maticAmount)
	checkErr(err)

	log.Infof("Waiting force batch proposal to be mined")
	err = waitTxToBeMined(ctx, client, tx.Hash(), txTimeout)
	checkErr(err)
}

func sendL2Claim(ctx context.Context, auth *bind.TransactOpts, client *ethclient.Client, deposit *Deposit, smtProof [][32]byte, globalExitRoot *globalExitRoot) {
	auth.GasPrice = big.NewInt(0)
	br, err := bridge.NewBridge(common.HexToAddress(l2BridgeAddr), client)
	checkErr(err)

	amount, _ := new(big.Int).SetString(deposit.Amount, encoding.Base10)
	tx, err := br.Claim(auth, common.HexToAddress(deposit.TokenAddr), amount, deposit.OrigNet, deposit.DestNet,
		common.HexToAddress(deposit.DestAddr), smtProof, uint32(deposit.DepositCnt), globalExitRoot.GlobalExitRootNum,
		globalExitRoot.ExitRoots[0], globalExitRoot.ExitRoots[1])
	checkErr(err)

	log.Infof("waiting L2 Claim tx to be mined")
	const txTimeout = 15 * time.Second
	err = waitTxToBeMined(ctx, client, tx.Hash(), txTimeout)
	checkErr(err)

	log.Infof("wait for the consolidation")
	const timeToWaitForTheConsolidation = 30 * time.Second
	time.Sleep(timeToWaitForTheConsolidation)
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
		log.Fatal(err)
	}
}
