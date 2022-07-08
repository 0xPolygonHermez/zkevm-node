package vectors

import (
	"math/big"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
)

type argBigInt struct {
	big.Int
}

func (a argBigInt) MarshalJSON() ([]byte, error) {
	return []byte(a.Text(hex.Base)), nil
}

func (a *argBigInt) UnmarshalJSON(input []byte) error {
	str := strings.Trim(string(input), "\"")
	if strings.ToLower(strings.TrimSpace(str)) == "null" {
		return nil
	}

	bi, err := encoding.DecodeUint256orHex(&str)
	if err != nil {
		return err
	}

	a.Int = *bi

	return nil
}
