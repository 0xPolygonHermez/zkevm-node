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
	"golang.org/x/crypto/sha3"
)

const (
	networkURL = "http://localhost:8123"

	accPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	txMinedTimeoutLimit = 60 * time.Second

	// compiled using http://remix.ethereum.org
	// COMPILER: 0.8.7+commit.e28d00a7
	// OPTIMIZATION: disabled
	// ../../test/contracts/emitLog.sol
	counterHexBytes = "608060405234801561001057600080fd5b50610173806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806306661abd1461003b578063d09de08a14610059575b600080fd5b610043610063565b6040516100509190610093565b60405180910390f35b610061610069565b005b60005481565b600160008082825461007b91906100ae565b92505081905550565b61008d81610104565b82525050565b60006020820190506100a86000830184610084565b92915050565b60006100b982610104565b91506100c483610104565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156100f9576100f861010e565b5b828201905092915050565b6000819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fdfea2646970667358221220fb8f8c48bca53685975855e0f6abe633c0649ccea8e22fd38f2636291e4d956364736f6c63430008070033"
	// ../../test/contracts/emitLog.sol
	emitLogHexBytes = "608060405234801561001057600080fd5b506102e8806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80637966b4f614610030575b600080fd5b61003861003a565b005b7f5e7df75d54e493185612379c616118a4c9ac802de621b010c96f74d22df4b30a60405160405180910390a160017f977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac6260405160405180910390a2600260017fbb6e4da744abea70325874159d52c1ad3e57babfae7c329a948e7dcb274deb0960405160405180910390a36003600260017f966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b760405160405180910390a46003600260017fe5562b12d9276c5c987df08afff7b1946f2d869236866ea2285c7e2e95685a64600460405161012c9190610269565b60405180910390a46002600360047fe5562b12d9276c5c987df08afff7b1946f2d869236866ea2285c7e2e95685a64600160405161016a919061024e565b60405180910390a46001600260037f966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b760405160405180910390a4600160027fbb6e4da744abea70325874159d52c1ad3e57babfae7c329a948e7dcb274deb0960405160405180910390a360017f977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac6260405160405180910390a27f5e7df75d54e493185612379c616118a4c9ac802de621b010c96f74d22df4b30a60405160405180910390a1565b6102398161028e565b82525050565b610248816102a0565b82525050565b60006020820190506102636000830184610230565b92915050565b600060208201905061027e600083018461023f565b92915050565b6000819050919050565b600061029982610284565b9050919050565b60006102ab82610284565b905091905056fea26469706673582212202a7046afe48f762569461cfa53c2a188613e57a3bb6f580e745dfa4feee9562b64736f6c63430008070033"
	// ../../test/contracts/erc20.sol
	erc20HexBytes = "60806040526040518060400160405280601381526020017f536f6c6964697479206279204578616d706c65000000000000000000000000008152506003908051906020019062000051929190620000d0565b506040518060400160405280600781526020017f534f4c4259455800000000000000000000000000000000000000000000000000815250600490805190602001906200009f929190620000d0565b506012600560006101000a81548160ff021916908360ff160217905550348015620000c957600080fd5b50620001e5565b828054620000de9062000180565b90600052602060002090601f0160209004810192826200010257600085556200014e565b82601f106200011d57805160ff19168380011785556200014e565b828001600101855582156200014e579182015b828111156200014d57825182559160200191906001019062000130565b5b5090506200015d919062000161565b5090565b5b808211156200017c57600081600090555060010162000162565b5090565b600060028204905060018216806200019957607f821691505b60208210811415620001b057620001af620001b6565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b610d9680620001f56000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c806342966c681161007157806342966c681461016857806370a082311461018457806395d89b41146101b4578063a0712d68146101d2578063a9059cbb146101ee578063dd62ed3e1461021e576100a9565b806306fdde03146100ae578063095ea7b3146100cc57806318160ddd146100fc57806323b872dd1461011a578063313ce5671461014a575b600080fd5b6100b661024e565b6040516100c39190610b06565b60405180910390f35b6100e660048036038101906100e19190610a18565b6102dc565b6040516100f39190610aeb565b60405180910390f35b6101046103ce565b6040516101119190610b28565b60405180910390f35b610134600480360381019061012f91906109c5565b6103d4565b6040516101419190610aeb565b60405180910390f35b610152610585565b60405161015f9190610b43565b60405180910390f35b610182600480360381019061017d9190610a58565b610598565b005b61019e60048036038101906101999190610958565b61066f565b6040516101ab9190610b28565b60405180910390f35b6101bc610687565b6040516101c99190610b06565b60405180910390f35b6101ec60048036038101906101e79190610a58565b610715565b005b61020860048036038101906102039190610a18565b6107ec565b6040516102159190610aeb565b60405180910390f35b61023860048036038101906102339190610985565b610909565b6040516102459190610b28565b60405180910390f35b6003805461025b90610c8c565b80601f016020809104026020016040519081016040528092919081815260200182805461028790610c8c565b80156102d45780601f106102a9576101008083540402835291602001916102d4565b820191906000526020600020905b8154815290600101906020018083116102b757829003601f168201915b505050505081565b600081600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516103bc9190610b28565b60405180910390a36001905092915050565b60005481565b600081600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546104629190610bd0565b9250508190555081600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546104b89190610bd0565b9250508190555081600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461050e9190610b7a565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516105729190610b28565b60405180910390a3600190509392505050565b600560009054906101000a900460ff1681565b80600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546105e79190610bd0565b92505081905550806000808282546105ff9190610bd0565b92505081905550600073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516106649190610b28565b60405180910390a350565b60016020528060005260406000206000915090505481565b6004805461069490610c8c565b80601f01602080910402602001604051908101604052809291908181526020018280546106c090610c8c565b801561070d5780601f106106e25761010080835404028352916020019161070d565b820191906000526020600020905b8154815290600101906020018083116106f057829003601f168201915b505050505081565b80600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546107649190610b7a565b925050819055508060008082825461077c9190610b7a565b925050819055503373ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516107e19190610b28565b60405180910390a350565b600081600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461083d9190610bd0565b9250508190555081600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546108939190610b7a565b925050819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516108f79190610b28565b60405180910390a36001905092915050565b6002602052816000526040600020602052806000526040600020600091509150505481565b60008135905061093d81610d32565b92915050565b60008135905061095281610d49565b92915050565b60006020828403121561096e5761096d610d1c565b5b600061097c8482850161092e565b91505092915050565b6000806040838503121561099c5761099b610d1c565b5b60006109aa8582860161092e565b92505060206109bb8582860161092e565b9150509250929050565b6000806000606084860312156109de576109dd610d1c565b5b60006109ec8682870161092e565b93505060206109fd8682870161092e565b9250506040610a0e86828701610943565b9150509250925092565b60008060408385031215610a2f57610a2e610d1c565b5b6000610a3d8582860161092e565b9250506020610a4e85828601610943565b9150509250929050565b600060208284031215610a6e57610a6d610d1c565b5b6000610a7c84828501610943565b91505092915050565b610a8e81610c16565b82525050565b6000610a9f82610b5e565b610aa98185610b69565b9350610ab9818560208601610c59565b610ac281610d21565b840191505092915050565b610ad681610c42565b82525050565b610ae581610c4c565b82525050565b6000602082019050610b006000830184610a85565b92915050565b60006020820190508181036000830152610b208184610a94565b905092915050565b6000602082019050610b3d6000830184610acd565b92915050565b6000602082019050610b586000830184610adc565b92915050565b600081519050919050565b600082825260208201905092915050565b6000610b8582610c42565b9150610b9083610c42565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115610bc557610bc4610cbe565b5b828201905092915050565b6000610bdb82610c42565b9150610be683610c42565b925082821015610bf957610bf8610cbe565b5b828203905092915050565b6000610c0f82610c22565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600060ff82169050919050565b60005b83811015610c77578082015181840152602081019050610c5c565b83811115610c86576000848401525b50505050565b60006002820490506001821680610ca457607f821691505b60208210811415610cb857610cb7610ced565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600080fd5b6000601f19601f8301169050919050565b610d3b81610c04565b8114610d4657600080fd5b50565b610d5281610c42565b8114610d5d57600080fd5b5056fea264697066735822122014c474e8e3584cd97ec8c18b47b4f80f3e1fe4161014922d35adfb31108ba9cc64736f6c63430008070033"
	// ../../test/contracts/storage.sol
	storageHexBytes = "608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100d9565b60405180910390f35b610073600480360381019061006e919061009d565b61007e565b005b60008054905090565b8060008190555050565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea264697066735822122066e69d9a7c6de5190ac263274e26fc172ce35f8692e67f911dc62025fc4ee53064736f6c63430008070033"
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

	var scAddr common.Address
	scAddr = deploySC(ctx, client, auth, counterHexBytes, 400000)
	sendTxsToCounterSC(ctx, client, auth, scAddr)

	scAddr = deploySC(ctx, client, auth, emitLogHexBytes, 400000)
	sendTxsToEmitLogSC(ctx, client, auth, scAddr)

	scAddr = deploySC(ctx, client, auth, erc20HexBytes, 1200000)
	sendTxsToERC20SC(ctx, client, auth, scAddr)

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

func sendTxsToERC20SC(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, scAddr common.Address) {
	// mint
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write([]byte("mint(uint256)"))
	chkErr(err)
	methodID := hash.Sum(nil)[:4]
	a, _ := big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	amount := common.LeftPadBytes(a.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, amount...)

	log.Infof("sending mint")
	scCall(ctx, client, auth, scAddr, data)
	log.Infof("mint processed successfully")

	// transfer
	hash = sha3.NewLegacyKeccak256()
	_, err = hash.Write([]byte("transfer(address,uint256)"))
	chkErr(err)
	methodID = hash.Sum(nil)[:4]
	receiver := common.LeftPadBytes(common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D").Bytes(), 32)
	a, _ = big.NewInt(0).SetString("1000000000000000000000", encoding.Base10)
	amount = common.LeftPadBytes(a.Bytes(), 32)

	data = []byte{}
	data = append(data, methodID...)
	data = append(data, receiver...)
	data = append(data, amount...)

	log.Infof("sending transfer")
	scCall(ctx, client, auth, scAddr, data)
	log.Infof("transfer processed successfully")

	// invalid transfer - no enough balance
	// TODO: uncomment this when hezcore is able to handle reverted transactions
	// hash = sha3.NewLegacyKeccak256()
	// _, err := hash.Write([]byte("transfer(address,uint256)"))
	// chkErr(err)
	// methodID = hash.Sum(nil)[:4]
	// receiver = common.LeftPadBytes(common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D").Bytes(), 32)
	// a, _ = big.NewInt(0).SetString("2000000000000000000000", encoding.Base10)
	// amount = common.LeftPadBytes(a.Bytes(), 32)

	// data = []byte{}
	// data = append(data, methodID...)
	// data = append(data, receiver...)
	// data = append(data, amount...)

	// log.Infof("sending transfer")
	// scCall(ctx, client, auth, scAddr, data)
	// log.Infof("transfer processed successfully")
}

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
