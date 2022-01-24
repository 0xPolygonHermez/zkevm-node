package helper

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/umbracle/fastrlp"
)

var addressPool fastrlp.ArenaPool

// CreateAddress creates an Ethereum address.
func CreateAddress(addr common.Address, nonce uint64) common.Address {
	a := addressPool.Get()
	defer addressPool.Put(a)

	v := a.NewArray()
	v.Set(a.NewBytes(addr.Bytes()))
	v.Set(a.NewUint(nonce))

	dst := v.MarshalTo(nil)
	dst = Keccak256(dst)[12:]

	return common.BytesToAddress(dst)
}

var create2Prefix = []byte{0xff}

// CreateAddress2 creates an Ethereum address following the CREATE2 Opcode.
func CreateAddress2(addr common.Address, salt [32]byte, inithash []byte) common.Address {
	return common.BytesToAddress(Keccak256(create2Prefix, addr.Bytes(), salt[:], Keccak256(inithash))[12:])
}
