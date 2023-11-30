package test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

const vectorString = `[
    {
        "nonce": "",
        "gasPrice": "0x3b9aca00",
        "gasLimit": "186a0",
        "to": "0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98",
        "value": "0x100",
        "data": "0xs5b8e9959000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000da00000000000000000000000000000000000000000000000000000000000000da608060405234801561001057600080fd5b5060bb8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063f6d97e8c14602d575b600080fd5b60336047565b604051603e91906062565b60405180910390f35b6000806003s90508091505090565b605c81607b565b82525050565b6000602082019050607560008301846055565b92915050565b600081905091905056fea2646970667358221220a33fdecaf587db45fa0e1fe4bfca25de09e35bb9a45fa6dab1bf1964244a929164736f6c63430008070033000000000000",
        "from": "0x4d5Cf5032B2a844602278b01199ED191A86c93ff",
        "l2TxHash": "failed",
        "reason": "Invalid nonce value"
    },
    {
        "nonce": "0x00",
        "gasPrice": "0x3b9aca00",
        "gasLimit": "0x186a0",
        "to": "0x005Cf5032B2a844602278b01199ED191A86c93ff",
        "value": "0x00",
        "data": "0x56d5be740000000000000000000000001275fbb540c8efc58b812ba83b0d0b8b9917ae98",
        "from": "0x4d5Cf5032B2a844602278b01199ED191A86c93ff",
        "l2TxHash": "0x8f9b0375a6b0f1bd9d54ff499921766828ae8e5314fc44a494736b5c4cc3bb56"
    },
    {
        "nonce": "0x01",
        "gasPrice": "0x3b9aca00",
        "gasLimit": "0x186a0",
        "to": "0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98",
        "value": "0x100",
        "data": "0x5b8e9959000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000da00000000000000000000000000000000000000000000000000000000000000da608060405234801561001057600080fd5b5060bb8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063f6d97e8c14602d575b600080fd5b60336047565b604051603e91906062565b60405180910390f35b600080600390508091505090565b605c81607b565b82525050565b6000602082019050607560008301846055565b92915050565b600081905091905056fea2646970667358221220a33fdecaf587db45fa0e1fe4bfca25de09e35bb9a45fa6dab1bf1964244a929164736f6c63430008070033000000000000",
        "from": "0x4d5Cf5032B2a844602278b01199ED191A86c93ff",
        "l2TxHash": "0xaa8f08e5bee683718f3f14fa352aaeb8e7de49f8b0e59f03128ef37fa6ac18e3"
    },
    {
        "nonce": "0x01",
        "gasPrice": "0x3b9aca00",
        "gasLimit": "v186a0",
        "to": "0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98",
        "value": "0x100",
        "data": "0x5b8e9959000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000da00000000000000000000000000000000000000000000000000000000000000da608060405234801561001057600080fd5b5060bb8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063f6d97e8c14602d575b600080fd5b60336047565b604051603e91906062565b60405180910390f35b600080600390508091505090565b605c81607b565b82525050565b6000602082019050607560008301846055565b92915050565b600081905091905056fea2646970667358221220a33fdecaf587db45fa0e1fe4bfca25de09e35bb9a45fa6dab1bf1964244a929164736f6c63430008070033000000000000",
        "from": "0x4d5Cf5032B2a844602278b01199ED191A86c93ff",
        "l2TxHash": "failed",
        "reason": "Invalid gasLimit value"
    },
    {
        "nonce": "0x21",
        "gasPrice": "0x3b9aca00",
        "gasLimit": "186a0",
        "to": "0x1275fbb540c8efc58b812ba83b0d0b8b9917ae98",
        "value": "0x100",
        "data": "0xs5b8e9959000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000da00000000000000000000000000000000000000000000000000000000000000da608060405234801561001057600080fd5b5060bb8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063f6d97e8c14602d575b600080fd5b60336047565b604051603e91906062565b60405180910390f35b6000806003s90508091505090565b605c81607b565b82525050565b6000602082019050607560008301846055565b92915050565b600081905091905056fea2646970667358221220a33fdecaf587db45fa0e1fe4bfca25de09e35bb9a45fa6dab1bf1964244a929164736f6c63430008070033000000000000",
        "from": "0x4d5Cf5032B2a844602278b01199ED191A86c93ff",
        "l2TxHash": "failed",
        "reason": "Invalid data value"
    }
]`

type testVector struct {
	Nonce     string `json:"nonce"`
	GasPrice  string `json:"gasPrice"`
	GasLimit  string `json:"gasLimit"`
	To        string `json:"to"`
	Value     string `json:"value"`
	Data      string `json:"data"`
	From      string `json:"from"`
	L2TxHash  string `json:"l2TxHash"`
	Reason    string `json:"reason"`
	Signature string `json:"signature"`
}

func TestL2TxHash(t *testing.T) {
	// Unmarshall the test vector
	var testVectors []testVector
	err := json.Unmarshal([]byte(vectorString), &testVectors)
	if err != nil {
		require.NoError(t, err)
	}

	// Create types.Transaction from test vector
	for x, testVector := range testVectors {
		nonce := new(big.Int).SetBytes(common.FromHex(testVector.Nonce)).Uint64()
		gasPrice := new(big.Int).SetBytes(common.FromHex(testVector.GasPrice))
		gasLimit := new(big.Int).SetBytes(common.FromHex(testVector.GasLimit)).Uint64()
		to := common.HexToAddress(testVector.To)
		value := new(big.Int).SetBytes(common.FromHex(testVector.Value))
		data := common.FromHex(testVector.Data)
		from := common.HexToAddress(testVector.From)

		if testVector.L2TxHash != "failed" {
			log.Debug("Test vector: ", x)
			log.Debugf("nonce: %x", nonce)
			log.Debugf("gasPrice: %x", gasPrice)
			log.Debugf("gasLimit: %x", gasLimit)
			log.Debugf("to: %s", to.String())
			log.Debugf("value: %x", value)
			log.Debugf("data: %s", common.Bytes2Hex(data))
			log.Debugf("from: %s", from.String())

			tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)
			require.NoError(t, err)

			hash, err := state.TestGetL2Hash(*tx, from)
			require.NoError(t, err)

			require.Equal(t, testVector.L2TxHash, hash.String())
		}
	}
}
