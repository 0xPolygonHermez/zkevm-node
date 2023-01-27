package jsonrpc

import (
	"math/big"

	"golang.org/x/crypto/sha3"
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
	hash := sha3.NewLegacyKeccak256()
	hash.Write(b.Bytes()) //nolint:errcheck,gosec
	keccak256Hash := hash.Sum(nil)
	return argBytes(keccak256Hash), nil
}
