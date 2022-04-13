package pool

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/hex"
)

func Test_IsClaimTx(t *testing.T) {

	l2GlobalExitRootManagerAddr := common.HexToAddress("0x00000000000000000000000000000001")
	differentAddr := common.HexToAddress("0x00000000000000000000000000000002")
	claimData := hex.DecodeHexToBig(bridgeClaimMethodSignature).Bytes()

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
			Name: "To address as Any address other than l2GlobalExitRootManagerAddr address",
			Tx: Transaction{
				Transaction: *types.NewTx(&types.LegacyTx{Nonce: 1, To: &differentAddr, Value: big.NewInt(0), Gas: 0, GasPrice: big.NewInt(0), Data: claimData}),
			},
			expectedResult: false,
		},
		{
			Name: "To address as l2GlobalExitRootManagerAddr address",
			Tx: Transaction{
				Transaction: *types.NewTx(&types.LegacyTx{Nonce: 1, To: &l2GlobalExitRootManagerAddr, Value: big.NewInt(0), Gas: 0, GasPrice: big.NewInt(0), Data: claimData}),
			},
			expectedResult: false,
		},
		{
			Name: "To address as l2GlobalExitRootManagerAddr address",
			Tx: Transaction{
				Transaction: *types.NewTx(&types.LegacyTx{Nonce: 1, To: &l2GlobalExitRootManagerAddr, Value: big.NewInt(0), Gas: 0, GasPrice: big.NewInt(0), Data: claimData}),
			},
			expectedResult: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			result := testCase.Tx.IsClaimTx(l2GlobalExitRootManagerAddr)
			if result != testCase.expectedResult {
				t.Errorf("Invalid result, expected: %v, found: %v", testCase.expectedResult, result)
			}
		})
	}
}
