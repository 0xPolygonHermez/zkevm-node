package pool

import (
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func Test_IsClaimTx(t *testing.T) {
	l2BridgeAddr := common.HexToAddress("0x00000000000000000000000000000001")
	differentAddr := common.HexToAddress("0x00000000000000000000000000000002")
	claimData, err := hex.DecodeHex(BridgeClaimMethodSignature)
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		Name           string
		Tx             Transaction
		expectedResult bool
	}{
		{
			Name: "To address as nil",
			Tx: Transaction{
				Transaction: *types.NewTx(&types.LegacyTx{Nonce: 1, To: nil, Value: big.NewInt(0), Gas: 0, GasPrice: big.NewInt(0), Data: claimData}),
			},
			expectedResult: false,
		},
		{
			Name: "To address as Zeroaddress",
			Tx: Transaction{
				Transaction: *types.NewTx(&types.LegacyTx{Nonce: 1, To: &common.Address{}, Value: big.NewInt(0), Gas: 0, GasPrice: big.NewInt(0), Data: claimData}),
			},
			expectedResult: false,
		},
		{
			Name: "To address as Any address other than l2BridgeAddr address",
			Tx: Transaction{
				Transaction: *types.NewTx(&types.LegacyTx{Nonce: 1, To: &differentAddr, Value: big.NewInt(0), Gas: 0, GasPrice: big.NewInt(0), Data: claimData}),
			},
			expectedResult: false,
		},
		{
			Name: "To address as l2BridgeAddr address",
			Tx: Transaction{
				Transaction: *types.NewTx(&types.LegacyTx{Nonce: 1, To: &l2BridgeAddr, Value: big.NewInt(0), Gas: 50000, GasPrice: big.NewInt(0), Data: claimData}),
			},
			expectedResult: true,
		},
		{
			Name: "More Gas than 150K",
			Tx: Transaction{
				Transaction: *types.NewTx(&types.LegacyTx{Nonce: 1, To: &l2BridgeAddr, Value: big.NewInt(0), Gas: 160000, GasPrice: big.NewInt(0), Data: claimData}),
			},
			expectedResult: false,
		},
		{
			Name: "Tx with Gas 150K",
			Tx: Transaction{
				Transaction: *types.NewTx(&types.LegacyTx{Nonce: 1, To: &l2BridgeAddr, Value: big.NewInt(0), Gas: 150000, GasPrice: big.NewInt(0), Data: claimData}),
			},
			expectedResult: true,
		},
	}
	const freeClaimGasLimit uint64 = 150000
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			result := testCase.Tx.IsClaimTx(l2BridgeAddr, freeClaimGasLimit)
			if result != testCase.expectedResult {
				t.Errorf("Invalid result, expected: %v, found: %v", testCase.expectedResult, result)
			}
		})
	}
}
