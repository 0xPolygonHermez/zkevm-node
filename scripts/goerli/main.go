package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/scripts"
	"github.com/hermeznetwork/hermez-core/test/contracts/bin/FailureTest"
)

const (
	networkURL = "http://localhost:8123"
	// pk         = "0xdfd01798f92667dbf91df722434e8fbe96af0211d4d1b82bbbbc8f1def7a814f"
	txTimeout = 60 * time.Second
	two       = 2
	four      = 4
)

func main() {
	ctx := context.Background()

	// if you want to test using goerli network
	pk := ""         // replace this by your goerli account private key
	networkURL := "" // replace this by your goerli infura url

	log.Infof("connecting to %v", networkURL)
	client, err := ethclient.Dial(networkURL)
	chkErr(err)
	log.Infof("connected")

	chainID, err := client.ChainID(ctx)
	chkErr(err)
	log.Infof("chainID: %v", chainID)

	auth := getAuth(ctx, client, pk)
	fmt.Println()

	balance, err := client.BalanceAt(ctx, auth.From, nil)
	chkErr(err)

	log.Debugf("ETH Balance for %v: %v", auth.From, balance)

	// deploy FailureTest SC
	failureTestSCAddr, tx, failureTestSC, err := FailureTest.DeployFailureTest(auth, client)
	chkErr(err)
	_, err = scripts.WaitTxToBeMined(client, tx.Hash(), txTimeout)
	chkErr(err)

	// check values
	checkValues(ctx, client, failureTestSCAddr, failureTestSC)

	// update number with valid tx
	tx, err = failureTestSC.Store(auth, big.NewInt(two))
	chkErr(err)
	_, err = scripts.WaitTxToBeMined(client, tx.Hash(), txTimeout)
	chkErr(err)

	// check values
	checkValues(ctx, client, failureTestSCAddr, failureTestSC)

	// update number with invalid tx
	_, err = failureTestSC.StoreAndFail(auth, big.NewInt(four))
	chkErr(err)
	// _, err = scripts.WaitTxToBeMined(client, tx.Hash(), txTimeout)
	// chkErr(err)

	// check values
	checkValues(ctx, client, failureTestSCAddr, failureTestSC)
}

func checkValues(ctx context.Context, client *ethclient.Client, failureTestSCAddr common.Address, failureTestSC *FailureTest.FailureTest) {
	number, err := failureTestSC.GetNumber(nil)
	chkErr(err)
	log.Debugf("Number: %v", number)
	logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
		Addresses: []common.Address{failureTestSCAddr},
	})
	chkErr(err)
	log.Debugf("Logs: %v", len(logs))
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

func chkErr(err error) {
	if err != nil {
		log.Error(err)
	}
}
