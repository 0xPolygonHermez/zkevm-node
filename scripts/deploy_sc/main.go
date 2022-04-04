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
	"github.com/hermeznetwork/hermez-core/encoding"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/test/testutils"
	"golang.org/x/crypto/sha3"
)

const (
	networkURL = "http://localhost:8123"

	accPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	txMinedTimeoutLimit = 60 * time.Second
)

// this function sends a transaction to deploy a smartcontract to the local network
func main() {
	ctx := context.Background()

	log.Infof("connecting to %v", networkURL)
	client, err := ethclient.Dial(networkURL)
	chkErr(err)
	log.Infof("connected")

	auth := getAuth(ctx, client)

	sendEthTransaction(ctx, client, auth)

	counterHexBytes, err := testutils.ReadBytecode("counter/Counter.bin")
	chkErr(err)
	emitLogHexBytes, err := testutils.ReadBytecode("emitLog/EmitLog.bin")
	chkErr(err)
	// erc20HexBytes, err := testutils.ReadBytecode("erc20/ERC20.bin")
	// chkErr(err)
	storageHexBytes, err := testutils.ReadBytecode("storage/Storage.bin")
	chkErr(err)

	var scAddr common.Address
	scAddr = deploySC(ctx, client, auth, counterHexBytes, 400000)
	sendTxsToCounterSC(ctx, client, auth, scAddr)

	scAddr = deploySC(ctx, client, auth, emitLogHexBytes, 400000)
	sendTxsToEmitLogSC(ctx, client, auth, scAddr)

	// scAddr = deploySC(ctx, client, auth, erc20HexBytes, 1200000)
	// sendTxsToERC20SC(ctx, client, auth, scAddr)

	scAddr = deploySC(ctx, client, auth, storageHexBytes, 400000)
	sendTxsToStorageSC(ctx, client, auth, scAddr)
}

func sendEthTransaction(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts) {
	// ETH Transfer
	to := common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	amount, _ := big.NewInt(0).SetString("10000000000000000000", encoding.Base10)
	ethTransfer(ctx, client, auth, to, amount)

	// Invalid ETH Transfer - no enough balance
	// TODO: uncomment this when hezcore is able to handle reverted transactions
	// to = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	// amount, _ = big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	// ethTransfer(ctx, client, auth, to, amount)
}

func sendTxsToCounterSC(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, scAddr common.Address) {
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write([]byte("increment()"))
	chkErr(err)
	data := hash.Sum(nil)[:4]

	log.Infof("sending tx to increment counter")
	scCall(ctx, client, auth, scAddr, data)
	log.Infof("counter incremented")
}

func sendTxsToEmitLogSC(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, scAddr common.Address) {
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write([]byte("emitLogs()"))
	chkErr(err)
	data := hash.Sum(nil)[:4]

	log.Infof("sending tx to increment counter")
	scCall(ctx, client, auth, scAddr, data)
	log.Infof("counter incremented")
}

// func sendTxsToERC20SC(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, scAddr common.Address) {
// 	// mint
// 	hash := sha3.NewLegacyKeccak256()
// 	_, err := hash.Write([]byte("mint(uint256)"))
// 	chkErr(err)
// 	methodID := hash.Sum(nil)[:4]
// 	a, _ := big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
// 	amount := common.LeftPadBytes(a.Bytes(), 32)

// 	var data []byte
// 	data = append(data, methodID...)
// 	data = append(data, amount...)

// 	log.Infof("sending mint")
// 	scCall(ctx, client, auth, scAddr, data)
// 	log.Infof("mint processed successfully")

// 	// transfer
// 	hash = sha3.NewLegacyKeccak256()
// 	_, err = hash.Write([]byte("transfer(address,uint256)"))
// 	chkErr(err)
// 	methodID = hash.Sum(nil)[:4]
// 	receiver := common.LeftPadBytes(common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D").Bytes(), 32)
// 	a, _ = big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
// 	amount = common.LeftPadBytes(a.Bytes(), 32)

// 	data = []byte{}
// 	data = append(data, methodID...)
// 	data = append(data, receiver...)
// 	data = append(data, amount...)

// 	log.Infof("sending transfer")
// 	scCall(ctx, client, auth, scAddr, data)
// 	log.Infof("transfer processed successfully")

// 	// invalid transfer - no enough balance
// 	// TODO: uncomment this when hezcore is able to handle reverted transactions
// 	// hash = sha3.NewLegacyKeccak256()
// 	// _, err := hash.Write([]byte("transfer(address,uint256)"))
// 	// chkErr(err)
// 	// methodID = hash.Sum(nil)[:4]
// 	// receiver = common.LeftPadBytes(common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D").Bytes(), 32)
// 	// a, _ = big.NewInt(0).SetString("2000000000000000000000", encoding.Base10)
// 	// amount = common.LeftPadBytes(a.Bytes(), 32)

// 	// data = []byte{}
// 	// data = append(data, methodID...)
// 	// data = append(data, receiver...)
// 	// data = append(data, amount...)

// 	// log.Infof("sending transfer")
// 	// scCall(ctx, client, auth, scAddr, data)
// 	// log.Infof("transfer processed successfully")
// }

func sendTxsToStorageSC(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, scAddr common.Address) {
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write([]byte("store(uint256)"))
	chkErr(err)
	methodID := hash.Sum(nil)[:4]
	const numberToStore = 22
	number := common.LeftPadBytes(big.NewInt(numberToStore).Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, number...)

	log.Infof("sending tx to store number: %v", number)
	scCall(ctx, client, auth, scAddr, data)
	log.Infof("number stored")
}

func ethTransfer(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, to common.Address, amount *big.Int) {
	log.Infof("reading nonce for account: %v", auth.From.Hex())
	nonce, err := client.NonceAt(ctx, auth.From, nil)
	log.Infof("nonce: %v", nonce)
	chkErr(err)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	chkErr(err)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{To: &to})
	chkErr(err)

	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)

	signedTx, err := auth.Signer(auth.From, tx)
	chkErr(err)

	log.Infof("sending transfer tx")
	err = client.SendTransaction(ctx, signedTx)
	chkErr(err)
	log.Infof("tx sent: %v", signedTx.Hash().Hex())

	_, err = waitTxToBeMined(client, signedTx.Hash(), txMinedTimeoutLimit)
	chkErr(err)

	log.Infof("tx processed successfully!")
}

func scCall(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, scAddr common.Address, data []byte) {
	log.Infof("reading nonce for account: %v", auth.From.Hex())
	nonce, err := client.NonceAt(ctx, auth.From, nil)
	log.Infof("nonce: %v", nonce)
	chkErr(err)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	chkErr(err)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{To: &scAddr, Data: data})
	chkErr(err)

	tx := types.NewTransaction(nonce, scAddr, big.NewInt(0), gasLimit, gasPrice, data)

	signedTx, err := auth.Signer(auth.From, tx)
	chkErr(err)

	log.Infof("calling SC")
	err = client.SendTransaction(ctx, signedTx)
	chkErr(err)
	log.Infof("tx sent: %v", signedTx.Hash().Hex())

	_, err = waitTxToBeMined(client, signedTx.Hash(), txMinedTimeoutLimit)
	chkErr(err)

	log.Infof("tx processed successfully!")
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

func deploySC(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, scHexByte string, gasLimit uint64) common.Address {
	log.Infof("reading nonce for account: %v", auth.From.Hex())
	nonce, err := client.NonceAt(ctx, auth.From, nil)
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
			return nil, fmt.Errorf("transaction has failed: %s", string(r.PostState))
		}

		return r, nil
	}
}

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
