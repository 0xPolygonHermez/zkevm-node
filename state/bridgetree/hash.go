package bridgetree

import "github.com/iden3/go-iden3-crypto/keccak256"

func hash(left [32]byte, right [32]byte) [32]byte {
	var res [32]byte
	copy(res[:], keccak256.Hash(left[:], right[:]))
	return res
}
