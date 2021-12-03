package tree

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// maxBigIntLen is 256 bits (32 bytes)
const maxBigIntLen = 32

// SplitAddress splits address into 3 bytes array of 64bits each
func SplitAddress(address common.Address) ([][]byte, error) {
	addr, err := scalar2fea(new(big.Int).SetBytes(address[:]))
	if err != nil {
		return nil, err
	}
	return [][]byte{addr[0].Bytes(), addr[1].Bytes(), addr[2].Bytes()}, nil
}

// SplitValue splits value into 4 bytes array of 64bits each
func SplitValue(value *big.Int) ([][]byte, error) {
	val, err := scalar2fea(value)
	if err != nil {
		return nil, err
	}
	return [][]byte{val[0].Bytes(), val[1].Bytes(), val[2].Bytes(), val[3].Bytes()}, nil
}

func fea2scalar(v []*big.Int) *big.Int {
	res := v[0]
	res.Add(res, new(big.Int).Lsh(v[1], 64))
	res.Add(res, new(big.Int).Lsh(v[2], 128))
	res.Add(res, new(big.Int).Lsh(v[3], 192))
	return res
}

func scalar2fea(value *big.Int) ([]*big.Int, error) {
	val := make([]*big.Int, 4)
	mask, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFF", 16)
	val[0] = new(big.Int).And(value, mask)
	val[1] = new(big.Int).And(new(big.Int).Rsh(value, 64), mask)
	val[2] = new(big.Int).And(new(big.Int).Rsh(value, 128), mask)
	val[3] = new(big.Int).And(new(big.Int).Rsh(value, 192), mask)
	return val[:], nil
}
