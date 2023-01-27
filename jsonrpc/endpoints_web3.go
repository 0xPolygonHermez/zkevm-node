package jsonrpc

import (
	"math/big"

	"github.com/iden3/go-iden3-crypto/keccak256"
)

// Web3Endpoints contains implementations for the "web3" RPC endpoints
type Web3Endpoints struct {
}

// ClientVersion returns the client version.
func (e *Web3Endpoints) ClientVersion() (interface{}, rpcError) {
	return "Polygon Hermez zkEVM/v2.0.0", nil
}

// Sha3 returns the keccak256 hash of the given data.
func (e *Web3Endpoints) Sha3(data argBig) (interface{}, rpcError) {
	b := (*big.Int)(&data)
	return argBytes(keccak256.Hash(b.Bytes())), nil
}
