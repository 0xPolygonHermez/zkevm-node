package bridgetree

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-crypto/keccak256"
)

// VerifyMerkleProof verifies Merkle Proof for specified leaf with its index against provided root
func VerifyMerkleProof(leaf [32]byte, proof [][32]byte, index int, root [32]byte) bool {
	value := leaf
	for i := 0; i < len(proof); i++ {
		if int(math.Floor(float64(index)/math.Pow(2, float64(i))))%2 != 0 {
			value = hash(proof[i], value)
		} else {
			value = hash(value, proof[i])
		}
	}
	return bytes.Compare(value[:], root[:]) == 0
}

// CalculateLeafValue calculates leaf value
func CalculateLeafValue(originalNetwork uint32, tokenAddress common.Address, amount *big.Int, destinationNetwork uint32, destinationAddress common.Address) [32]byte {
	var res [32]byte
	origNet := make([]byte, 4)
	binary.LittleEndian.PutUint32(origNet, originalNetwork)
	destNet := make([]byte, 4)
	binary.LittleEndian.PutUint32(destNet, destinationNetwork)
	var buf [32]byte
	copy(res[:], keccak256.Hash(origNet, tokenAddress[:], amount.FillBytes(buf[:]), destNet, destinationAddress[:]))
	return res
}
