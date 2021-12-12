package tree

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// maxBigIntLen is 256 bits (32 bytes)
const maxBigIntLen = 32

// splitAddress splits address into 3 byte array of 64bits each
func splitAddress(address common.Address) ([][]byte, error) {
	addr, err := scalar2fea(new(big.Int).SetBytes(address[:]))
	if err != nil {
		return nil, err
	}
	return [][]byte{addr[0].Bytes(), addr[1].Bytes(), addr[2].Bytes()}, nil
}

// splitValue splits 256bit value into 4 byte arrays of 64bits each
func splitValue(value *big.Int) ([][]byte, error) {
	val, err := scalar2fea(value)
	if err != nil {
		return nil, err
	}
	return [][]byte{val[0].Bytes(), val[1].Bytes(), val[2].Bytes(), val[3].Bytes()}, nil
}

// fea2scalar converts array of 64bit big.Int values into one 256bit big.Int
func fea2scalar(v []*big.Int) *big.Int {
	res := new(big.Int).Set(v[0])
	res.Add(res, new(big.Int).Lsh(v[1], 64))
	res.Add(res, new(big.Int).Lsh(v[2], 128))
	res.Add(res, new(big.Int).Lsh(v[3], 192))
	return res
}

// fea2scalar splits 256bit big.Int into array of 64bit big.Int values
func scalar2fea(value *big.Int) ([]*big.Int, error) {
	val := make([]*big.Int, 4)
	mask, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFF", 16)
	val[0] = new(big.Int).And(value, mask)
	val[1] = new(big.Int).And(new(big.Int).Rsh(value, 64), mask)
	val[2] = new(big.Int).And(new(big.Int).Rsh(value, 128), mask)
	val[3] = new(big.Int).And(new(big.Int).Rsh(value, 192), mask)
	return val[:], nil
}
