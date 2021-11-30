package tree

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// maxBigIntLen is 256 bits (32 bytes)
const maxBigIntLen = 32
const splitBigIntLen = 8

// SplitAddress splits address into 3 bytes array of 64bits each
func SplitAddress(address common.Address) ([][]byte, error) {
	a0 := address[0:8]
	a1 := address[8:16]
	a2 := address[16:20]
	return [][]byte{a0, a1, a2}, nil
}

// SplitValue splits value into 4 bytes array of 64bits each
func SplitValue(value *big.Int) ([][]byte, error) {
	val := value.Bytes()
	if len(val) > maxBigIntLen {
		return nil, fmt.Errorf("value size of more than 256 bits is not supported")
	}
	val256 := make([]byte, maxBigIntLen)
	copy(val256, val)
	v0 := val256[0:8]
	v1 := val256[8:16]
	v2 := val256[16:24]
	v3 := val256[24:32]
	return [][]byte{v0, v1, v2, v3}, nil
}

func fea2scalar(v []*big.Int) *big.Int {
	var buf bytes.Buffer
	for i := 0; i < len(v); i++ {
		var b [splitBigIntLen]byte
		copy(b[:], v[i].Bytes())
		buf.Write(b[:])
	}
	return new(big.Int).SetBytes(buf.Bytes())
}

func scalar2fea(value *big.Int) ([]*big.Int, error) {
	val := make([]*big.Int, 4)
	v, err := SplitValue(value)
	if err != nil {
		return nil, err
	}
	for i := 0; i < 4; i++ {
		val[i] = new(big.Int).SetBytes(v[i])
	}
	return val, nil
}
