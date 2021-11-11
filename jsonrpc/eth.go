package jsonrpc

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var (
	block = map[string]interface{}{
		"difficulty":       "0x2",
		"extraData":        "0x506172697479205465636820417574686f726974790000000000000000000000bf254a15d223a85e9c9f708b09426b9c5d37f3de96f58cedc9ebd1c3df4f259b07b1bec96e22d5cce45c8e2f38b3027cd189bd1ef2e6e94865d437720e3b47df01",
		"gasLimit":         "0x7a1200",
		"gasUsed":          "0x0",
		"hash":             "0xb0550d9b3033c5bc04407af24af13f943d5636e3674e1e93855219a6fe2ed885",
		"logsBloom":        "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"miner":            "0x0000000000000000000000000000000000000000",
		"mixHash":          "0x0000000000000000000000000000000000000000000000000000000000000000",
		"nonce":            "0x0000000000000000",
		"number":           "0x1b4",
		"parentHash":       "0xa49ec4e60be8e6334c0123578286d228e4ca16b29db4a8b15081194f91c93413",
		"receiptsRoot":     "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"sha3Uncles":       "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"size":             "0x260",
		"stateRoot":        "0x5d6cded585e73c4e322c30c2f782a336316f17dd85a4863b9d838d2d4b8b3008",
		"timestamp":        "0x5c53297a",
		"totalDifficulty":  "0x369",
		"transactions":     []interface{}{},
		"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"uncles":           []interface{}{},
	}
	tx = map[string]interface{}{
		"blockHash":        "0xb1c40f3aea853a830ecfa07cf68031a47163381d5641ab41ea723e42e6839f71",
		"blockNumber":      "0x58fb10",
		"from":             "0xf2955cd74ab4b24faf4bc732061c332e7296e756",
		"gas":              "0x186a0",
		"gasPrice":         "0x59682f08",
		"hash":             "0x08a9451bd9e5128c3f118eb1745b26aa50e5705a2911f3395e890615dae394ed",
		"input":            "0xa5977fbb000000000000000000000000f2955cd74ab4b24faf4bc732061c332e7296e756000000000000000000000000f4b2cbc3ba04c478f0dc824f4806ac39982dce7300000000000000000000000000000000000000000000000000000000009896800000000000000000000000000000000000000000000000000000000000000061000000000000000000000000000000000000000000000000000000000000126400000000000000000000000000000000000000000000000000000000000dbba0",
		"nonce":            "0x1264",
		"to":               "0xb0bbd74d211948e1d2917c17fdb47898e9b75b26",
		"transactionIndex": "0x3a",
		"value":            "0x0",
		"type":             "0x0",
		"v":                "0x2e",
		"r":                "0xb7798059a9058ca7c6857d1c551573eb7c73952ffb9300789631ae2c9eea5387",
		"s":                "0x2f6bb3ccacd48523dbb308a58e720a75a7dd7998218c7c07ed275d9af026fbc2",
	}
)

// Eth is the eth jsonrpc endpoint
type Eth struct{}

// BlockNumber returns current block number
func (e *Eth) BlockNumber() (interface{}, error) {
	return "0x58faca", nil
}

// ChainId returns the chain id of the client
func (e *Eth) ChainId() (interface{}, error) {
	return "0x99999999", nil // 2576980377
}

func (e *Eth) EstimateGas(arg *txnArgs, rawNum *BlockNumber) (interface{}, error) {
	b, _ := json.Marshal(arg)
	fmt.Println("arg", string(b))
	fmt.Println("rawNum", rawNum)
	return "0x0", nil
}

// GasPrice returns the average gas price based on the last x blocks
func (e *Eth) GasPrice() (interface{}, error) {
	return "0x0", nil
}

// GetBalance returns the account's balance at the referenced block
func (e *Eth) GetBalance(address common.Address, number *BlockNumber) (interface{}, error) {
	fmt.Println("address", address)
	fmt.Println("number", number)
	return "0xa157a16b4775f0ae2", nil
}

// GetBlockByHash returns information about a block by hash
func (e *Eth) GetBlockByHash(hash common.Hash, fullTx bool) (interface{}, error) {
	fmt.Println("hash", hash)
	fmt.Println("fullTx", fullTx)
	return block, nil
}

// GetBlockByNumber returns information about a block by block number
func (e *Eth) GetBlockByNumber(number BlockNumber, fullTx bool) (interface{}, error) {
	fmt.Println("number", number)
	fmt.Println("fullTx", fullTx)
	return block, nil
}

// GetCode returns account code at given block number
func (e *Eth) GetCode(address common.Address, number *BlockNumber) (interface{}, error) {
	fmt.Println("address", address)
	fmt.Println("number", number)
	return "0x", nil
}

func (e *Eth) GetTransactionByBlockHashAndIndex(hash common.Hash, index Index) (interface{}, error) {
	fmt.Println("hash", hash)
	fmt.Println("index", index)
	return tx, nil
}

func (e *Eth) GetTransactionByBlockNumberAndIndex(number *BlockNumber, index Index) (interface{}, error) {
	fmt.Println("number", number)
	fmt.Println("index", index)
	return tx, nil
}

// GetTransactionByHash returns a transaction by his hash
func (e *Eth) GetTransactionByHash(hash common.Hash) (interface{}, error) {
	fmt.Println("hash", hash)
	return tx, nil
}

// GetTransactionCount returns account nonce
func (e *Eth) GetTransactionCount(address common.Address, number *BlockNumber) (interface{}, error) {
	fmt.Println("address", address)
	fmt.Println("number", number)
	return "0x1265", nil
}

// GetTransactionReceipt returns a transaction receipt by his hash
func (e *Eth) GetTransactionReceipt(hash common.Hash) (interface{}, error) {
	fmt.Println("hash", hash)
	return nil, nil
}

// SendRawTransaction sends a raw transaction
func (e *Eth) SendRawTransaction(input string) (interface{}, error) {
	fmt.Println("input", input)
	return nil, nil
}
