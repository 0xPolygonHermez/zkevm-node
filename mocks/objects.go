package mocks

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hermeznetwork/hermez-core/state"
)

const (
	blockNumber  = 54321
	batchNumber  = 12345
	balance      = 112233
	estimatedGas = 111222
	gasPrice     = 222333
	txNonce      = 1
	txAmount     = 987654321
	address      = "0x03e75d7DD38CCE2e20FfEE35EC914C57780A8e29"
)

var (
	block       = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(blockNumber)})
	batch       = &state.Batch{Number: batchNumber}
	txToAddress = common.HexToAddress(address)
	tx          = types.NewTransaction(txNonce, txToAddress, big.NewInt(txAmount), estimatedGas, big.NewInt(gasPrice), []byte{})
	txReceipt   = types.NewReceipt([]byte{}, false, 1234)
)
