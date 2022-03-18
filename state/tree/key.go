package tree

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/log"
	poseidon "github.com/iden3/go-iden3-crypto/goldenposeidon"
)

// Key stores key of the leaf
type Key [32]byte

var (
	capIn = [4]uint64{}
)

// keyEthAddr is the common code for all the keys related to ethereum addresses.
func keyEthAddr(ethAddr common.Address, leafType leafType, key1 [8]uint64) ([]byte, error) {
	ethAddrBI := new(big.Int).SetBytes(ethAddr.Bytes())
	ethAddrArr := scalar2fea(ethAddrBI)

	key0 := [8]uint64{
		ethAddrArr[0],
		ethAddrArr[1],
		ethAddrArr[2],
		ethAddrArr[3],
		ethAddrArr[4],
		0,
		uint64(leafType),
		0,
	}
	hk0, err := poseidon.Hash(key0, capIn)
	if err != nil {
		return nil, err
	}

	hk1, err := poseidon.Hash(key1, capIn)
	if err != nil {
		return nil, err
	}
	result, err := poseidon.Hash([8]uint64{
		hk0[0],
		hk0[1],
		hk0[2],
		hk0[3],
		hk1[0],
		hk1[1],
		hk1[2],
		hk1[3],
	}, capIn)
	if err != nil {
		return nil, err
	}

	return h4ToScalar(result[:]).Bytes(), nil
}

// KeyEthAddrBalance returns the key of balance leaf:
//   hk0: H([ethAddr[0:4], ethAddr[4:8], ethAddr[8:12], ethAddr[12:16], ethAddr[16:20], 0, 0, 0])
//   hk1: H([0, 0, 0, 0, 0, 0, 0, 0])
//   key = H([...hk0, ...hk1])
func KeyEthAddrBalance(ethAddr common.Address) ([]byte, error) {
	return keyEthAddr(ethAddr, leafTypeBalance, [8]uint64{})
}

// KeyEthAddrNonce returns the key of nonce leaf:
//   hk0: H([ethAddr[0:4], ethAddr[4:8], ethAddr[8:12], ethAddr[12:16], ethAddr[16:20], 0, 1, 0])
//   hk1: H([0, 0, 0, 0, 0, 0, 0, 0])
//   key = H([...hk0, ...hk1])
func KeyEthAddrNonce(ethAddr common.Address) ([]byte, error) {
	return keyEthAddr(ethAddr, leafTypeNonce, [8]uint64{})
}

// KeyContractCode returns the key of contract code leaf:
//   hk0: H([ethAddr[0:4], ethAddr[4:8], ethAddr[8:12], ethAddr[12:16], ethAddr[16:20], 0, 2, 0])
//   hk1: H([0, 0, 0, 0, 0, 0, 0, 0])
//   key = H([...hk0, ...hk1])
func KeyContractCode(ethAddr common.Address) ([]byte, error) {
	return keyEthAddr(ethAddr, leafTypeCode, [8]uint64{})
}

// KeyContractStorage returns the key of contract storage position leaf:
//   hk0: H([ethAddr[0:4], ethAddr[4:8], ethAddr[8:12], ethAddr[12:16], ethAddr[16:20], 0, 3, 0])
//   hk1: H([stoPos[0:4], stoPos[4:8], stoPos[8:12], stoPos[12:16], stoPos[16:20], stoPos[20:24], stoPos[24:28], stoPos[28:32])
//   key = H([...hk0, ...hk1])
func KeyContractStorage(ethAddr common.Address, storagePos []byte) ([]byte, error) {
	storageBI := new(big.Int).SetBytes(storagePos)

	storageArr := scalar2fea(storageBI)

	key1 := [8]uint64{
		storageArr[0],
		storageArr[1],
		storageArr[2],
		storageArr[3],
		storageArr[4],
		storageArr[5],
		storageArr[6],
		storageArr[7],
	}
	log.Debugf("key1: %v", key1)
	return keyEthAddr(ethAddr, leafTypeStorage, key1)
}
