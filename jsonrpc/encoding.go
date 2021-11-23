package jsonrpc

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/hermeznetwork/hermez-core/jsonrpc/hex"
)

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
	return strconv.ParseUint(str, base, bitSize64)
}

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

func DecodeInt64orHex(val *string) (int64, error) {
	i, err := DecodeUint64orHex(val)
	return int64(i), err
}

func DecodeBytes(val *string) ([]byte, error) {
	if val == nil {
		return []byte{}, nil
	}

	str := strings.TrimPrefix(*val, "0x")

	return hex.DecodeString(str)
}

func EncodeUint64(b uint64) *string {
	res := fmt.Sprintf("0x%x", b)
	return &res
}

func EncodeBytes(b []byte) *string {
	res := "0x" + hex.EncodeToString(b)
	return &res
}

func EncodeBigInt(b *big.Int) *string {
	res := "0x" + b.Text(hexBase)
	return &res
}
