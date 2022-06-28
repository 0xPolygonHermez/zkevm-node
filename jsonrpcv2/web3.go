package jsonrpcv2

import (
	"math/big"

	"github.com/iden3/go-iden3-crypto/keccak256"
)

// Web3 contains implementations for the "web3" RPC endpoints
type Web3 struct {
}

func (w *Web3) ClientVersion() (interface{}, error) {
	return "Polygon Hermez/v1.5.0", nil
}

func (w *Web3) Sha3(data argBig) (interface{}, error) {
	b := (*big.Int)(&data)
	return argBytes(keccak256.Hash(b.Bytes())), nil
}
