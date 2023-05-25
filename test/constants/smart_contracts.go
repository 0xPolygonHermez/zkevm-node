package constants

import "github.com/ethereum/go-ethereum/crypto"

var (
	ForcedBatchSignatureHash = crypto.Keccak256Hash([]byte("ForceBatch(uint64,bytes32,address,bytes)"))
)
