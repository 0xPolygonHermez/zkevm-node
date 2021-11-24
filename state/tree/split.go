package tree

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

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
	if len(val) > 32 {
		return nil, fmt.Errorf("value size of more than 256 bits is not supported")
	}
	val256 := make([]byte, 32)
	copy(val256, val)
	v0 := val256[0:8]
	v1 := val256[8:16]
	v2 := val256[16:24]
	v3 := val256[24:32]
	return [][]byte{v0, v1, v2, v3}, nil
}
