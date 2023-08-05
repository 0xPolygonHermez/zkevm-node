package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	var networks = []struct {
		Name       string
		URL        string
		ChainID    uint64
		PrivateKey string
	}{
		// {Name: "Local L1", URL: operations.DefaultL1NetworkURL, ChainID: operations.DefaultL1ChainID, PrivateKey: operations.DefaultSequencerPrivateKey},
		// {Name: "Local L2", URL: "http://34.90.16.102:8545/", ChainID: 1234, PrivateKey: operations.DefaultSequencerPrivateKey},
		{Name: "Local L2", URL: "https://rpc.public.zkevm-test.net/", ChainID: 1442, PrivateKey: "5d120f76469fb621f9b74587096f9bb292a27ebe346a2c84071010030356c20c"},
	}

	for _, network := range networks {
		ctx := context.Background()

		log.Infof("connecting to %v: %v", network.Name, network.URL)
		client, err := ethclient.Dial(network.URL)
		chkErr(err)
		log.Infof("connected")

		auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)
		chkErr(err)

		const receiverAddr = "0xdc6Bf73BC0A688bC11D5234Cb0F2719672Babd30"

		balance, err := client.BalanceAt(ctx, auth.From, nil)
		chkErr(err)
		log.Debugf("ETH Balance for %v: %v", auth.From, balance)

		// Valid ETH Transfer
		balance, err = client.BalanceAt(ctx, auth.From, nil)
		log.Debugf("ETH Balance for %v: %v", auth.From, balance)
		chkErr(err)
		transferAmount := big.NewInt(1)
		log.Debugf("Transfer Amount: %v", transferAmount)

		nonce, err := client.NonceAt(ctx, auth.From, nil)
		chkErr(err)
		// var lastTxHash common.Hash
		for i := 0; i < 100000; i++ {
			nonce := nonce + uint64(i)
			log.Debugf("Sending TX to transfer ETH")
			to := common.HexToAddress(receiverAddr)
			tx := ethTransfer(ctx, client, auth, to, transferAmount, &nonce)
			fmt.Println("tx sent: ", tx.Hash().String())
			// lastTxHash = tx.Hash()
			time.Sleep(500 * time.Millisecond)
		}

		// err = operations.WaitTxToBeMined(client, lastTxHash, txTimeout)
		// chkErr(err)
	}
}

func ethTransfer(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, to common.Address, amount *big.Int, nonce *uint64) *types.Transaction {
	if nonce == nil {
		log.Infof("reading nonce for account: %v", auth.From.Hex())
		var err error
		n, err := client.NonceAt(ctx, auth.From, nil)
		log.Infof("nonce: %v", n)
		chkErr(err)
		nonce = &n
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	chkErr(err)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{To: &to})
	chkErr(err)

	tx := types.NewTransaction(*nonce, to, amount, gasLimit, gasPrice, nil)

	signedTx, err := auth.Signer(auth.From, tx)
	chkErr(err)

	log.Infof("sending transfer tx")
	err = client.SendTransaction(ctx, signedTx)
	chkErr(err)
	log.Infof("tx sent: %v", signedTx.Hash().Hex())

	rlp, err := signedTx.MarshalBinary()
	chkErr(err)

	log.Infof("tx rlp: %v", hex.EncodeToHex(rlp))

	return signedTx
}

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
