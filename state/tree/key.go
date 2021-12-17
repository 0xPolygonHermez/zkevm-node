package tree

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-crypto/poseidon"
)

// Key stores key of the leaf
type Key [32]byte

// GetKey calculates Key for the provided leaf type, address, and in case of LeafTypeStorage also storagePosition.
// For other leaf types leave storagePosition = nil
func GetKey(leafType LeafType, address common.Address, storagePosition []byte, arity uint8, hashFunction HashFunction) ([]byte, error) {
	poseidonInputsNum := 1 << arity

	addr, err := splitAddress(address)
	if err != nil {
		return nil, err
	}
	inputs := make([]*big.Int, poseidonInputsNum)

	// initialize with zeroes
	for i := 0; i < poseidonInputsNum; i++ {
		inputs[i] = big.NewInt(0)
	}

	inputs[0].SetBytes(addr[0])
	inputs[1].SetBytes(addr[1])
	inputs[2].SetBytes(addr[2])
	inputs[3].SetUint64(uint64(leafType))

	if leafType == LeafTypeStorage {
		posBigInt := big.NewInt(0).SetBytes(storagePosition)
		pos, err := splitValue(posBigInt)
		if err != nil {
			return nil, err
		}
		inputs[4].SetBytes(pos[0])
		inputs[5].SetBytes(pos[1])
		inputs[6].SetBytes(pos[2])
		inputs[7].SetBytes(pos[3])
	}

	if hashFunction == nil {
		hashFunction = poseidon.Hash
	}

	hash, err := hashFunction(inputs)
	if err != nil {
		return nil, err
	}
	hashBytes := hash.Bytes()
	return hashBytes, nil
}
