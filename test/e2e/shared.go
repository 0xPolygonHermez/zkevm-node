//nolint:deadcode,unused,varcheck
package e2e

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	invalidParamsErrorCode = -32602
	toAddressHex           = "0x4d5Cf5032B2a844602278b01199ED191A86c93ff"
)

var (
	toAddress = common.HexToAddress(toAddressHex)
)

var networks = []struct {
	Name         string
	URL          string
	WebSocketURL string
	ChainID      uint64
	PrivateKey   string
}{
	{
		Name:         "Local L1",
		URL:          operations.DefaultL1NetworkURL,
		WebSocketURL: operations.DefaultL1NetworkWebSocketURL,
		ChainID:      operations.DefaultL1ChainID,
		PrivateKey:   operations.DefaultSequencerPrivateKey,
	},
	{
		Name:         "Local L2",
		URL:          operations.DefaultL2NetworkURL,
		WebSocketURL: operations.DefaultL2NetworkWebSocketURL,
		ChainID:      operations.DefaultL2ChainID,
		PrivateKey:   operations.DefaultSequencerPrivateKey,
	},
}

func setup() {
	var err error
	ctx := context.Background()
	err = operations.Teardown()
	if err != nil {
		panic(err)
	}

	opsCfg := operations.GetDefaultOperationsConfig()
	opsMan, err := operations.NewManager(ctx, opsCfg)
	if err != nil {
		panic(err)
	}
	err = opsMan.Setup()
	if err != nil {
		panic(err)
	}
}

func teardown() {
	err := operations.Teardown()
	if err != nil {
		panic(err)
	}
}

func createTX(client *ethclient.Client, auth *bind.TransactOpts, to common.Address, amount *big.Int) (*ethTypes.Transaction, error) {
	nonce, err := client.NonceAt(context.Background(), auth.From, nil)
	if err != nil {
		return nil, err
	}
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{From: auth.From, To: &to, Value: amount})
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	log.Infof("\nTX details:\n\tNonce:    %d\n\tGasLimit: %d\n\tGasPrice: %d", nonce, gasLimit, gasPrice)
	if gasLimit != uint64(21000) { //nolint:gomnd
		return nil, fmt.Errorf("gasLimit %d != 21000", gasLimit)
	}
	tx := ethTypes.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return nil, err
	}
	log.Infof("Sending Tx %v Nonce %v", signedTx.Hash(), signedTx.Nonce())
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func logTx(tx *ethTypes.Transaction) {
	sender, _ := state.GetSender(*tx)
	log.Debugf("********************")
	log.Debugf("Hash: %v", tx.Hash())
	log.Debugf("From: %v", sender)
	log.Debugf("Nonce: %v", tx.Nonce())
	log.Debugf("ChainId: %v", tx.ChainId())
	log.Debugf("To: %v", tx.To())
	log.Debugf("Gas: %v", tx.Gas())
	log.Debugf("GasPrice: %v", tx.GasPrice())
	log.Debugf("Cost: %v", tx.Cost())

	// b, _ := tx.MarshalBinary()
	//log.Debugf("RLP: ", hex.EncodeToHex(b))
	log.Debugf("********************")
}
