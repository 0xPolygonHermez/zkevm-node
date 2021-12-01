package encoding

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/hermeznetwork/hermez-core/hex"
)

const (
	// BitSize64 64 bits
	BitSize64 = 64
)

// DecodeUint64orHex decodes a string uint64 or hex string into a uint64
func DecodeUint64orHex(val *string) (uint64, error) {
	if val == nil {
		return 0, nil
	}

	str := *val
	base := 10
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
		base = 16
	}
	return strconv.ParseUint(str, base, BitSize64)
}

// DecodeUint256orHex decodes a string uint256 or hex string into a bit.Int
func DecodeUint256orHex(val *string) (*big.Int, error) {
	if val == nil {
		return nil, nil
	}

	str := *val
	base := 10
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
		base = 16
	}
	b, ok := new(big.Int).SetString(str, base)
	if !ok {
		return nil, fmt.Errorf("could not parse")
	}
	return b, nil
}

// DecodeInt64orHex decodes a string int64 or hex string into a int64
func DecodeInt64orHex(val *string) (int64, error) {
	i, err := DecodeUint64orHex(val)
	return int64(i), err
}

// DecodeBytes decodes a hex string into a []byte
func DecodeBytes(val *string) ([]byte, error) {
	if val == nil {
		return []byte{}, nil
	}

	str := strings.TrimPrefix(*val, "0x")

	return hex.DecodeString(str)
}

// EncodeUint64 encodes a uint64 into a hex string
func EncodeUint64(b uint64) *string {
	res := fmt.Sprintf("0x%x", b)
	return &res
}

// EncodeBytes encodes a []bytes into a hex string
func EncodeBytes(b []byte) *string {
	res := "0x" + hex.EncodeToString(b)
	return &res
}

// EncodeBigInt encodes a big.Int into a hex string
func EncodeBigInt(b *big.Int) *string {
	res := "0x" + b.Text(hex.Base)
	return &res
}
