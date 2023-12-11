package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestL2BlockStartEncode(t *testing.T) {
	l2BlockStart := state.DSL2BlockStart{
		BatchNumber:    1,                           // 8 bytes
		L2BlockNumber:  2,                           // 8 bytes
		Timestamp:      3,                           // 8 bytes
		GlobalExitRoot: common.HexToHash("0x04"),    // 32 bytes
		Coinbase:       common.HexToAddress("0x05"), // 20 bytes
		ForkID:         5,
	}

	encoded := l2BlockStart.Encode()
	expected := []byte{1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 0}

	assert.Equal(t, expected, encoded)
}

func TestL2TransactionEncode(t *testing.T) {
	l2Transaction := state.DSL2Transaction{
		EffectiveGasPricePercentage: 128,                          // 1 byte
		IsValid:                     1,                            // 1 byte
		StateRoot:                   common.HexToHash("0x010203"), // 32 bytes
		EncodedLength:               5,                            // 4 bytes
		Encoded:                     []byte{1, 2, 3, 4, 5},        // 5 bytes
	}

	encoded := l2Transaction.Encode()
	expected := []byte{128, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 5, 0, 0, 0, 1, 2, 3, 4, 5}
	assert.Equal(t, expected, encoded)
}

func TestL2BlockEndEncode(t *testing.T) {
	l2BlockEnd := state.DSL2BlockEnd{
		L2BlockNumber: 1,                        // 8 bytes
		BlockHash:     common.HexToHash("0x02"), // 32 bytes
		StateRoot:     common.HexToHash("0x03"), // 32 bytes
	}

	encoded := l2BlockEnd.Encode()
	expected := []byte{1, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3}

	assert.Equal(t, expected, encoded)
}

func TestCalculateSCPosition(t *testing.T) {
	a := time.Now()
	blockNumber := uint64(2934867)
	expected := common.HexToHash("0xaa93c484856be45716623765b429a967296594ca362e61e91d671fb422e0f744")
	position := state.GetSystemSCPosition(blockNumber)
	assert.Equal(t, expected, common.BytesToHash(position))
	b := time.Now()

	c := b.Sub(a)
	fmt.Println(c)
}
