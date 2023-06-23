package state_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	forkID5 = 5
	forkID4 = 4
)

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stderr"},
	})
}

func TestDecodeRandomBatchL2Data(t *testing.T) {
	randomData := []byte("Random data")
	txs, _, _, err := state.DecodeTxs(randomData, forkID5)
	require.Error(t, err)
	assert.Equal(t, []types.Transaction{}, txs)
	t.Log("Txs decoded 1: ", txs)

	randomData = []byte("Esto es autentica basura")
	txs, _, _, err = state.DecodeTxs(randomData, forkID5)
	require.Error(t, err)
	assert.Equal(t, []types.Transaction{}, txs)
	t.Log("Txs decoded 2: ", txs)

	randomData = []byte("beef")
	txs, _, _, err = state.DecodeTxs(randomData, forkID5)
	require.Error(t, err)
	assert.Equal(t, []types.Transaction{}, txs)
	t.Log("Txs decoded 3: ", txs)
}

func TestDecodePre155BatchL2DataPreForkID5(t *testing.T) {
	pre155, err := hex.DecodeString("e480843b9aca00826163941275fbb540c8efc58b812ba83b0d0b8b9917ae98808464fbb77cb7d2a666860f3c6b8f5ef96f86c7ec5562e97fd04c2e10f3755ff3a0456f9feb246df95217bf9082f84f9e40adb0049c6664a5bb4c9cbe34ab1a73e77bab26ed1b")
	require.NoError(t, err)
	txs, _, _, err := state.DecodeTxs(pre155, forkID4)
	require.NoError(t, err)
	t.Log("Txs decoded: ", txs, len(txs))
	assert.Equal(t, 1, len(txs))
	v, r, s := txs[0].RawSignatureValues()
	assert.Equal(t, "0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98", txs[0].To().String())
	assert.Equal(t, "1b", fmt.Sprintf("%x", v))
	assert.Equal(t, "b7d2a666860f3c6b8f5ef96f86c7ec5562e97fd04c2e10f3755ff3a0456f9feb", fmt.Sprintf("%x", r))
	assert.Equal(t, "246df95217bf9082f84f9e40adb0049c6664a5bb4c9cbe34ab1a73e77bab26ed", fmt.Sprintf("%x", s))
	assert.Equal(t, uint64(24931), txs[0].Gas())
	assert.Equal(t, "64fbb77c", hex.EncodeToString(txs[0].Data()))
	assert.Equal(t, uint64(0), txs[0].Nonce())
	assert.Equal(t, new(big.Int).SetUint64(1000000000), txs[0].GasPrice())

	pre155, err = hex.DecodeString("e580843b9aca00830186a0941275fbb540c8efc58b812ba83b0d0b8b9917ae988084159278193d7bcd98c00060650f12c381cc2d4f4cc8abf54059aecd2c7aabcfcdd191ba6827b1e72f0eb0b8d5daae64962f4aafde7853e1c102de053edbedf066e6e3c2dc1b")
	require.NoError(t, err)
	txs, _, _, err = state.DecodeTxs(pre155, forkID4)
	require.NoError(t, err)
	t.Log("Txs decoded: ", txs)
	assert.Equal(t, 1, len(txs))
	assert.Equal(t, "0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98", txs[0].To().String())
	assert.Equal(t, uint64(0), txs[0].Nonce())
	assert.Equal(t, big.NewInt(0), txs[0].Value())
	assert.Equal(t, "15927819", hex.EncodeToString(txs[0].Data()))
	assert.Equal(t, uint64(100000), txs[0].Gas())
	assert.Equal(t, new(big.Int).SetUint64(1000000000), txs[0].GasPrice())
}

func TestDecodePre155BatchL2DataForkID5(t *testing.T) {
	pre155, err := hex.DecodeString("e480843b9aca00826163941275fbb540c8efc58b812ba83b0d0b8b9917ae98808464fbb77cb7d2a666860f3c6b8f5ef96f86c7ec5562e97fd04c2e10f3755ff3a0456f9feb246df95217bf9082f84f9e40adb0049c6664a5bb4c9cbe34ab1a73e77bab26ed1bff")
	require.NoError(t, err)
	txs, _, _, err := state.DecodeTxs(pre155, forkID5)
	require.NoError(t, err)
	t.Log("Txs decoded: ", txs, len(txs))
	assert.Equal(t, 1, len(txs))
	v, r, s := txs[0].RawSignatureValues()
	assert.Equal(t, "0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98", txs[0].To().String())
	assert.Equal(t, "1b", fmt.Sprintf("%x", v))
	assert.Equal(t, "b7d2a666860f3c6b8f5ef96f86c7ec5562e97fd04c2e10f3755ff3a0456f9feb", fmt.Sprintf("%x", r))
	assert.Equal(t, "246df95217bf9082f84f9e40adb0049c6664a5bb4c9cbe34ab1a73e77bab26ed", fmt.Sprintf("%x", s))
	assert.Equal(t, uint64(24931), txs[0].Gas())
	assert.Equal(t, "64fbb77c", hex.EncodeToString(txs[0].Data()))
	assert.Equal(t, uint64(0), txs[0].Nonce())
	assert.Equal(t, new(big.Int).SetUint64(1000000000), txs[0].GasPrice())

	pre155, err = hex.DecodeString("e580843b9aca00830186a0941275fbb540c8efc58b812ba83b0d0b8b9917ae988084159278193d7bcd98c00060650f12c381cc2d4f4cc8abf54059aecd2c7aabcfcdd191ba6827b1e72f0eb0b8d5daae64962f4aafde7853e1c102de053edbedf066e6e3c2dc1b")
	require.NoError(t, err)
	txs, _, _, err = state.DecodeTxs(pre155, forkID4)
	require.NoError(t, err)
	t.Log("Txs decoded: ", txs)
	assert.Equal(t, 1, len(txs))
	assert.Equal(t, "0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98", txs[0].To().String())
	assert.Equal(t, uint64(0), txs[0].Nonce())
	assert.Equal(t, big.NewInt(0), txs[0].Value())
	assert.Equal(t, "15927819", hex.EncodeToString(txs[0].Data()))
	assert.Equal(t, uint64(100000), txs[0].Gas())
	assert.Equal(t, new(big.Int).SetUint64(1000000000), txs[0].GasPrice())
}

func TestDecodePre155Tx(t *testing.T) {
	pre155 := "0xf86780843b9aca00826163941275fbb540c8efc58b812ba83b0d0b8b9917ae98808464fbb77c1ba0b7d2a666860f3c6b8f5ef96f86c7ec5562e97fd04c2e10f3755ff3a0456f9feba0246df95217bf9082f84f9e40adb0049c6664a5bb4c9cbe34ab1a73e77bab26ed"
	tx, err := state.DecodeTx(pre155)
	require.NoError(t, err)
	t.Log("Txs decoded: ", tx)
	v, r, s := tx.RawSignatureValues()
	assert.Equal(t, "0x1275fbb540c8efC58b812ba83B0D0B8b9917AE98", tx.To().String())
	assert.Equal(t, "1b", fmt.Sprintf("%x", v))
	assert.Equal(t, "b7d2a666860f3c6b8f5ef96f86c7ec5562e97fd04c2e10f3755ff3a0456f9feb", fmt.Sprintf("%x", r))
	assert.Equal(t, "246df95217bf9082f84f9e40adb0049c6664a5bb4c9cbe34ab1a73e77bab26ed", fmt.Sprintf("%x", s))
	assert.Equal(t, uint64(24931), tx.Gas())
	assert.Equal(t, "64fbb77c", hex.EncodeToString(tx.Data()))
	assert.Equal(t, uint64(0), tx.Nonce())
	assert.Equal(t, new(big.Int).SetUint64(1000000000), tx.GasPrice())
}

func TestEncodePre155BatchL2DataPreForkID5(t *testing.T) {
	pre155, err := hex.DecodeString("e480843b9aca00826163941275fbb540c8efc58b812ba83b0d0b8b9917ae98808464fbb77cb7d2a666860f3c6b8f5ef96f86c7ec5562e97fd04c2e10f3755ff3a0456f9feb246df95217bf9082f84f9e40adb0049c6664a5bb4c9cbe34ab1a73e77bab26ed1b")
	require.NoError(t, err)
	txs, _, effectivePercentages, err := state.DecodeTxs(pre155, forkID4)
	require.NoError(t, err)
	rawtxs, err := state.EncodeTransactions(txs, effectivePercentages, forkID4)
	require.NoError(t, err)
	t.Log("Txs decoded: ", txs, len(txs))
	assert.Equal(t, pre155, rawtxs)
}

func TestEncodePre155BatchL2DataForkID5(t *testing.T) {
	pre155, err := hex.DecodeString("e480843b9aca00826163941275fbb540c8efc58b812ba83b0d0b8b9917ae98808464fbb77cb7d2a666860f3c6b8f5ef96f86c7ec5562e97fd04c2e10f3755ff3a0456f9feb246df95217bf9082f84f9e40adb0049c6664a5bb4c9cbe34ab1a73e77bab26ed1bff")
	require.NoError(t, err)
	txs, _, effectivePercentages, err := state.DecodeTxs(pre155, forkID5)
	require.NoError(t, err)
	rawtxs, err := state.EncodeTransactions(txs, effectivePercentages, forkID5)
	require.NoError(t, err)
	t.Log("Txs decoded: ", txs, len(txs))
	assert.Equal(t, pre155, rawtxs)
}

func TestMaliciousTransaction(t *testing.T) {
	b := []byte{
		0xee, 0x80, 0x84, 0x3b, 0x9a, 0xca, 0x00, 0x83, 0x01, 0x86, 0xa0, 0x94,
		0x4d, 0x5c, 0xf5, 0x03, 0x2b, 0x2a, 0x84, 0x46, 0x02, 0x27, 0x8b, 0x01,
		0x19, 0x9e, 0xd1, 0x91, 0xa8, 0x6c, 0x93, 0xff, 0x88, 0x01, 0x63, 0x45,
		0x78, 0x5d, 0x8a, 0x00, 0x00, 0x80, 0x82, 0x01, 0x91, 0x80, 0x80, 0x6e,
		0x20, 0x9c, 0x61, 0xca, 0x92, 0xc2, 0xb9, 0x80, 0xd6, 0x19, 0x7e, 0x7a,
		0xc9, 0xcc, 0xc3, 0xf5, 0x47, 0xbf, 0x13, 0xbe, 0x64, 0x55, 0xdf, 0xe6,
		0x82, 0xaa, 0x5d, 0xda, 0x96, 0x55, 0xef, 0x16, 0x81, 0x9a, 0x7e, 0xdc,
		0xc3, 0xfe, 0xfe, 0xc8, 0x1c, 0xa9, 0x7c, 0x7a, 0x6f, 0x3d, 0x10, 0xec,
		0x77, 0x44, 0x40, 0xe4, 0x09, 0xad, 0xbb, 0xa6, 0x93, 0xce, 0x8b, 0x69,
		0x8d, 0x41, 0xf1, 0x1c, 0xef, 0x80, 0x84, 0x3b, 0x9a, 0xca, 0x00, 0x83,
		0x01, 0x86, 0xa0, 0x94, 0x4d, 0x5c, 0xf5, 0x03, 0x2b, 0x2a, 0x84, 0x46,
		0x02, 0x27, 0x8b, 0x01, 0x19, 0x9e, 0xd1, 0x91, 0xa8, 0x6c, 0x93, 0xff,
		0x89, 0x05, 0x6b, 0xc7, 0x5e, 0x2d, 0x63, 0x10, 0x00, 0x00, 0x80, 0x82,
		0x03, 0xe9, 0x80, 0x80, 0xfe, 0x1e, 0x96, 0xb3, 0x5c, 0x83, 0x6f, 0xbe,
		0xba, 0xc8, 0x87, 0x68, 0x11, 0x50, 0xc5, 0xfc, 0x9f, 0xda, 0xe8, 0x62,
		0xd7, 0x47, 0xaa, 0xaf, 0x8c, 0x30, 0x37, 0x3c, 0x0b, 0xec, 0xf7, 0x69,
		0x1f, 0xf0, 0xc9, 0x00, 0xaa, 0xaa, 0xc6, 0xd1, 0x56, 0x5a, 0x60, 0x3f,
		0x69, 0xb5, 0xa4, 0x5f, 0x22, 0x2e, 0xd2, 0x05, 0xf0, 0xa3, 0x6f, 0xdc,
		0x6e, 0x4e, 0x4c, 0x5a, 0x7b, 0x88, 0xd4, 0x5b, 0x1b, 0xee, 0x80, 0x84,
		0x3b, 0x9a, 0xca, 0x00, 0x83, 0x01, 0x86, 0xa0, 0x94, 0x4d, 0x5c, 0xf5,
		0x03, 0x2b, 0x2a, 0x84, 0x46, 0x02, 0x27, 0x8b, 0x01, 0x19, 0x9e, 0xd1,
		0x91, 0xa8, 0x6c, 0x93, 0xff, 0x88, 0x01, 0x63, 0x45, 0x78, 0x5d, 0x8a,
		0x00, 0x00, 0x80, 0x82, 0x01, 0x91, 0x80, 0x80, 0x6e, 0x20, 0x9c, 0x61,
		0xca, 0x92, 0xc2, 0xb9, 0x80, 0xd6, 0x19, 0x7e, 0x7a, 0xc9, 0xcc, 0xc3,
		0xf5, 0x47, 0xbf, 0x13, 0xbe, 0x64, 0x55, 0xdf, 0xe6, 0x82, 0xaa, 0x5d,
		0xda, 0x96, 0x55, 0xef, 0x16, 0x81, 0x9a, 0x7e, 0xdc, 0xc3, 0xfe, 0xfe,
		0xc8, 0x1c, 0xa9, 0x7c, 0x7a, 0x6f, 0x3d, 0x10, 0xec, 0x77, 0x44, 0x40,
		0xe4, 0x09, 0xad, 0xbb, 0xa6, 0x93, 0xce, 0x8b, 0x69, 0x8d, 0x41, 0xf1,
		0x1c, 0xef, 0x80, 0x84, 0x3b, 0x9a, 0xca, 0x00, 0x83, 0x01, 0x86, 0xa0,
		0x94, 0x4d, 0x5c, 0xf5, 0x03, 0x2b, 0x2a, 0x84, 0x46, 0x02, 0x27, 0x8b,
		0x01, 0x19, 0x9e, 0xd1, 0x91, 0xa8, 0x6c, 0x93, 0xff, 0x89, 0x05, 0x6b,
		0xc7, 0x5e, 0x2d, 0x63, 0x10, 0x00, 0x00, 0x80, 0x82, 0x03, 0xe9, 0x80,
		0x80, 0xfe, 0x1e, 0x96, 0xb3, 0x5c, 0x83, 0x6f, 0xbe, 0xba, 0xc8, 0x87,
		0x68, 0x11, 0x50, 0xc5, 0xfc, 0x9f, 0xda, 0xe8, 0x62, 0xd7, 0x47, 0xaa,
		0xaf, 0x8c, 0x30, 0x37, 0x3c, 0x0b, 0xec, 0xf7, 0x69, 0x1f, 0xf0, 0xc9,
		0x00, 0xaa, 0xaa, 0xc6, 0xd1, 0x56, 0x5a, 0x60, 0x3f, 0x69, 0xb5, 0xa4,
		0x5f, 0x22, 0x2e, 0xd2, 0x05, 0xf0, 0xa3, 0x6f, 0xdc, 0x6e, 0x4e, 0x4c,
		0x5a, 0x7b, 0x88, 0xd4, 0x5b, 0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0x00, 0x86, 0xa0, 0x94, 0x4d, 0x5c, 0xf5, 0x03, 0x2b, 0x2a,
		0x84, 0x46, 0x02, 0x27, 0x8b, 0x01, 0x19, 0x9e, 0xd1, 0x91, 0xa8, 0x6c,
		0x93, 0xff, 0x88, 0x01, 0x63, 0x45, 0x78, 0x5d, 0x8a, 0x00, 0x00, 0x80,
		0x82, 0x01, 0x91, 0x80, 0x80, 0x6e, 0x20, 0x9c, 0x61, 0xca, 0x92, 0xc2,
		0xb9, 0x80, 0xd6, 0x19, 0x7e, 0x7a, 0xc9, 0xcc, 0xc3, 0xf5, 0x47, 0xbf,
		0x13, 0xbe, 0x64, 0x55, 0xdf, 0xe6, 0x82, 0xaa, 0x5d, 0xda, 0x96, 0x55,
		0xef, 0x16, 0x81, 0x9a, 0x7e, 0xdc, 0xc3, 0xfe, 0xfe, 0xc8, 0x1c, 0xa9,
		0x7c, 0x7a, 0x6f, 0x3d, 0x10, 0xec, 0x77, 0x44, 0x40, 0xe4, 0x09, 0xad,
		0xbb, 0xa6, 0x93, 0xce, 0x8b, 0x69, 0x8d, 0x41, 0xf1, 0x1c, 0xef, 0x80,
		0x84, 0x3b, 0x9a, 0xca, 0x00, 0x83, 0x01, 0x86, 0xa0, 0x94, 0x4d, 0x5c,
		0xf5, 0x03, 0x2b, 0x2a, 0x84, 0x46, 0x02, 0x27, 0x8b, 0x01, 0x19, 0x9e,
		0xd1, 0x91, 0xa8, 0x6c, 0x93, 0xff, 0x89, 0x05, 0x6b, 0xc7, 0x5e, 0x2d,
		0x63, 0x10, 0x00, 0x00, 0x80, 0x82, 0x03, 0xe9, 0x80, 0x80, 0xfe, 0x1e,
		0x96, 0xb3, 0x5c, 0x83, 0x6f, 0xbe, 0xba, 0xc8, 0x87, 0x68, 0x11, 0x50,
		0xc5, 0xfc, 0x9f, 0xda, 0xe8, 0x62, 0xd7, 0x47, 0xaa, 0xaf, 0x8c, 0x30,
		0x37, 0x3c, 0x0b, 0xec, 0xf7, 0x69, 0x1f, 0xf0, 0xc9, 0x00, 0xaa, 0xaa,
		0xc6, 0xd1, 0x56, 0x5a, 0x60, 0x3f, 0x69, 0xb5, 0xa4, 0x5f, 0x22, 0x2e,
		0xd2, 0x05, 0xf0, 0xa3, 0x6f, 0xdc, 0x6e, 0x4e, 0x4c, 0x5a, 0x7b, 0x88,
		0xd4, 0x5b, 0x1b}

	_, _, _, err = state.DecodeTxs(b, forkID4)
	require.Error(t, err)
	require.Equal(t, err, state.ErrInvalidData)
}
