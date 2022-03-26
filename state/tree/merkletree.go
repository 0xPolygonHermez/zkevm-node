package tree

import (
	"context"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/iden3/go-iden3-crypto/ff"
	poseidon "github.com/iden3/go-iden3-crypto/goldenposeidon"
)

const (
	keyItems = 4

	modeInsertFound    = "insertFound"
	modeInsertNotFound = "insertNotFound"
	modeUpdate         = "update"
	modeDeleteFound    = "deleteFound"
	modeDeleteNotFound = "deleteNotFound"
	modeDeleteLast     = "deleteLast"
	modeZeroToZero     = "zeroToZero"
)

// MerkleTree implements merkle tree
type MerkleTree struct {
	store          Store
	hashFunction   HashFunction
	scHashFunction scHashFunction
	arity          uint8
}

// UpdateProof is a proof generated on Set operation
type UpdateProof struct {
	OldRoot  []uint64
	NewRoot  []uint64
	Key      []uint64
	Siblings [][]uint64
	InsKey   []uint64
	InsValue []uint64
	IsOld0   bool
	OldValue []uint64
	NewValue []uint64
}

// Proof is a proof generated on Get operation
type Proof struct {
	Root     []uint64
	Key      []uint64
	Value    []uint64
	Siblings [][]uint64
	IsOld0   bool
	InsKey   []uint64
	InsValue []uint64
}

// HashFunction is a function interface type to specify hash function that MT should use
type HashFunction func(inp [poseidon.NROUNDSF]uint64, cap [poseidon.CAPLEN]uint64) ([poseidon.CAPLEN]uint64, error)

type scHashFunction func(code []byte) ([]uint64, error)

// NewMerkleTree creates new MerkleTree instance
func NewMerkleTree(store Store, arity uint8, hashFunction HashFunction) *MerkleTree {
	if hashFunction == nil {
		hashFunction = poseidon.Hash
	}

	scHashFunction := hashContractBytecode

	return &MerkleTree{
		store:          store,
		arity:          arity,
		hashFunction:   hashFunction,
		scHashFunction: scHashFunction,
	}
}

// SupportsDBTransactions indicates whether the store implementation supports DB transactions
func (mt *MerkleTree) SupportsDBTransactions() bool {
	return mt.store.SupportsDBTransactions()
}

// BeginDBTransaction starts a transaction block
func (mt *MerkleTree) BeginDBTransaction(ctx context.Context) error {
	return mt.store.BeginDBTransaction(ctx)
}

// Commit commits a db transaction
func (mt *MerkleTree) Commit(ctx context.Context) error {
	return mt.store.Commit(ctx)
}

// Rollback rollbacks a db transaction
func (mt *MerkleTree) Rollback(ctx context.Context) error {
	return mt.store.Rollback(ctx)
}

// Set method sets value of a leaf at a given location (key) in the tree
func (mt *MerkleTree) Set(ctx context.Context, oldRoot []uint64, key []uint64, value []uint64) (*UpdateProof, error) {
	// exit early if context is cancelled
	err := ctx.Err()
	if err != nil {
		return nil, err
	}
	newRoot := oldRoot[:]
	r := oldRoot[:]

	keys := mt.splitKey(key[:])
	level := 0

	siblings := make([][]uint64, len(keys))

	var (
		insKey, insValue        []uint64
		foundKey                []uint64
		foundVal                []uint64
		mode                    string
		foundRKey, foundOldValH []uint64
		accKey                  []uint64
	)

	oldValue := make([]uint64, 8)
	isOld0 := true

	for (!nodeIsZero(r)) && (foundKey == nil) {
		node, err := mt.getNodeData(ctx, r)
		if err != nil {
			return nil, err
		}
		siblings[level] = node[:]
		if isOneSiblings(siblings[level]) {
			foundOldValH = siblings[level][4:8]
			foundValA, err := mt.getNodeData(ctx, foundOldValH[:8])
			if err != nil {
				return nil, err
			}
			foundRKey = siblings[level][:4]
			foundVal = foundValA[:]
			foundKey = mt.joinKey(accKey[:], foundRKey)
		} else {
			r = siblings[level][keys[level]*4 : keys[level]*4+4]
			accKey = append(accKey, uint64(keys[level]))
			level++
		}
	}

	level--
	if len(accKey) > 0 {
		accKey = accKey[:len(accKey)-1]
	}

	if !nodeIsZero(value) {
		if foundKey != nil {
			if nodeIsEq(key[:], foundKey) { // Update
				mode = modeUpdate

				newValH, err := mt.hashSave(ctx, value[:], []uint64{0, 0, 0, 0})
				if err != nil {
					return nil, err
				}
				newLeafHash, err := mt.hashSave(ctx,
					[]uint64{foundRKey[0], foundRKey[1], foundRKey[2], foundRKey[3], newValH[0], newValH[1], newValH[2], newValH[3]},
					[]uint64{1, 0, 0, 0})
				if err != nil {
					return nil, err
				}

				if level >= 0 {
					for j := 0; j < 4; j++ {
						siblings[level][keys[level]*4+uint(j)] = newLeafHash[j]
					}
				} else {
					newRoot = newLeafHash[:]
				}
			} else { // insert with foundKey
				mode = modeInsertFound

				node := make([]uint64, 8)
				level2 := level + 1
				foundKeys := mt.splitKey(foundKey)
				for keys[level2] == foundKeys[level2] {
					level2++
				}

				oldKey := removeKeyBits(foundKey, level2+1)
				oldLeafHash, err := mt.hashSave(ctx,
					[]uint64{oldKey[0], oldKey[1], oldKey[2], oldKey[3], foundOldValH[0], foundOldValH[1], foundOldValH[2], foundOldValH[3]},
					[]uint64{1, 0, 0, 0},
				)
				if err != nil {
					return nil, err
				}

				insKey = foundKey
				insValue = foundVal
				isOld0 = false

				newKey := removeKeyBits(key[:], level2+1)
				newValH, err := mt.hashSave(ctx, value[:], []uint64{0, 0, 0, 0})
				if err != nil {
					return nil, err
				}
				newLeafHash, err := mt.hashSave(ctx,
					[]uint64{newKey[0], newKey[1], newKey[2], newKey[3], newValH[0], newValH[1], newValH[2], newValH[3]},
					[]uint64{1, 0, 0, 0},
				)
				if err != nil {
					return nil, err
				}
				for i := 0; i < 8; i++ {
					node[i] = 0
				}
				for j := 0; j < 4; j++ {
					node[keys[level2]*4+uint(j)] = newLeafHash[j]
					node[foundKeys[level2]*4+uint(j)] = oldLeafHash[j]
				}

				r2, err := mt.hashSave(ctx, node, []uint64{0, 0, 0, 0})
				if err != nil {
					return nil, err
				}
				level2--

				for level2 != level {
					for i := 0; i < 8; i++ {
						node[i] = 0
					}
					for j := 0; j < 4; j++ {
						node[keys[level2]*4+uint(j)] = r2[j]
					}
					r2, err = mt.hashSave(ctx, node, []uint64{0, 0, 0, 0})
					if err != nil {
						return nil, err
					}
					level2--
				}

				if level >= 0 {
					for j := 0; j < 4; j++ {
						siblings[level][keys[level]*4+uint(j)] = r2[j]
					}
				} else {
					newRoot = r2
				}
			}
		} else { // insert without foundKey
			mode = modeInsertNotFound
			newKey := removeKeyBits(key[:], level+1)
			newValH, err := mt.hashSave(ctx, value[:], []uint64{0, 0, 0, 0})
			if err != nil {
				return nil, err
			}
			newLeafHash, err := mt.hashSave(ctx,
				[]uint64{newKey[0], newKey[1], newKey[2], newKey[3], newValH[0], newValH[1], newValH[2], newValH[3]},
				[]uint64{1, 0, 0, 0},
			)
			if err != nil {
				return nil, err
			}
			if level >= 0 {
				for j := 0; j < 4; j++ {
					siblings[level][keys[level]*4+uint(j)] = newLeafHash[j]
				}
			} else {
				newRoot = newLeafHash
			}
		}
	} else {
		if (foundKey != nil) && nodeIsEq(key[:], foundKey) { // Delete
			if level >= 0 {
				for j := 0; j < 4; j++ {
					siblings[level][keys[level]*4+uint(j)] = 0
				}

				uKey := mt.getUniqueSibling(siblings[level])

				if uKey >= 0 {
					mode = modeDeleteFound
					node, err := mt.getNodeData(ctx, siblings[level][uKey*4:uKey*4+4])
					if err != nil {
						return nil, err
					}
					siblings[level+1] = node

					if isOneSiblings(siblings[level+1]) {
						valH := siblings[level+1][4:]
						valANode, err := mt.getNodeData(ctx, valH)
						if err != nil {
							return nil, err
						}
						valA := valANode[0:8]
						val := fea2scalar(valA)
						rKey := siblings[level+1][:4]

						insKey = mt.joinKey([]uint64{accKey[0], accKey[1], accKey[2], accKey[3], uint64(uKey)}, rKey)
						insValue = scalar2fea(val)
						isOld0 = false

						for (uKey >= 0) && (level >= 0) {
							level--
							if level >= 0 {
								uKey = mt.getUniqueSibling(siblings[level])
							}
						}

						oldKey := removeKeyBits(insKey, level+1)
						oldLeafHash, err := mt.hashSave(ctx,
							[]uint64{oldKey[0], oldKey[1], oldKey[2], oldKey[3], valH[0], valH[1], valH[2], valH[3]},
							[]uint64{1, 0, 0, 0},
						)
						if err != nil {
							return nil, err
						}

						if level >= 0 {
							for j := 0; j < 4; j++ {
								siblings[level][keys[level]*4+uint(j)] = oldLeafHash[j]
							}
						} else {
							newRoot = oldLeafHash
						}
					} else {
						mode = modeDeleteNotFound
					}
				} else {
					mode = modeDeleteNotFound
				}
			} else {
				mode = modeDeleteLast
				newRoot = make([]uint64, 4)
			}
		} else {
			mode = modeZeroToZero
		}
	}

	siblings = siblings[:level+1]

	for level >= 0 {
		newRoot, err = mt.hashSave(ctx, siblings[level][:8], siblings[level][8:12])
		if err != nil {
			return nil, err
		}
		level--
		if level >= 0 {
			for j := 0; j < 4; j++ {
				siblings[level][keys[level]*4+uint(j)] = newRoot[j]
			}
		}
	}

	proof := UpdateProof{
		OldRoot:  oldRoot[:],
		NewRoot:  newRoot,
		Key:      key[:],
		Siblings: siblings,
		InsKey:   insKey,
		InsValue: insValue,
		IsOld0:   isOld0,
		OldValue: oldValue,
		NewValue: value,
	}

	log.Debugw("Set", "key", h4ToString(key[:]), "value", value, "mode", mode)

	return &proof, nil
}

// Get method gets value at a given location in the tree
func (mt *MerkleTree) Get(ctx context.Context, root, key []uint64) (*Proof, error) {
	// exit early if context is cancelled
	err := ctx.Err()
	if err != nil {
		return nil, err
	}

	r := root[:]

	keys := mt.splitKey(key)
	level := 0

	siblings := make([][]uint64, len(keys))

	var (
		insKey, insValue   []uint64
		foundKey, foundVal []uint64
		accKey             []uint64
	)

	value := make([]uint64, 8)
	isOld0 := true

	for (!nodeIsZero(r)) && (foundKey == nil) {
		node, err := mt.getNodeData(ctx, r)
		if err != nil {
			return nil, err
		}
		siblings[level] = node[:]

		if isOneSiblings(siblings[level]) {
			nodeValA, err := mt.getNodeData(ctx, siblings[level][4:8])
			if err != nil {
				return nil, err
			}
			foundValA := nodeValA[0:8]
			foundRKey := siblings[level][0:4]
			foundVal = foundValA[:]
			foundKey = mt.joinKey(accKey[:], foundRKey)
		} else {
			r = siblings[level][keys[level]*4 : keys[level]*4+4]
			accKey = append(accKey, uint64(keys[level]))
			level++
		}
	}

	level--

	if foundKey != nil {
		if nodeIsEq(key, foundKey) {
			value = foundVal
		} else {
			insKey = foundKey
			insValue = foundVal
			isOld0 = false
		}
	}

	siblings = siblings[:level+1]

	return &Proof{
		Root:     root,
		Key:      key,
		Value:    value,
		Siblings: siblings,
		IsOld0:   isOld0,
		InsKey:   insKey,
		InsValue: insValue,
	}, nil
}

func nodeIsZero(a []uint64) bool {
	result := true
	for i := 0; i < len(a); i++ {
		result = result && a[i] == 0
	}
	return result
}

func (mt *MerkleTree) getUniqueSibling(a []uint64) int64 {
	nFound := 0
	var fnd int64
	for i := 0; i < len(a); i += keyItems {
		if !nodeIsZero(a[i : i+keyItems]) {
			nFound++
			fnd = int64(i / keyItems)
		}
	}
	if nFound == 1 {
		return fnd
	}
	return -1
}

func (mt *MerkleTree) hashSave(ctx context.Context, nodeData []uint64, cap []uint64) ([]uint64, error) {
	if len(nodeData) != poseidon.NROUNDSF {
		return nil, fmt.Errorf("Invalid data length of %v", nodeData)
	}
	if len(cap) != poseidon.CAPLEN {
		return nil, fmt.Errorf("Invalid cap length of %v", cap)
	}
	capIn := [4]uint64{cap[0], cap[1], cap[2], cap[3]}
	nd := [8]uint64{nodeData[0], nodeData[1], nodeData[2], nodeData[3], nodeData[4], nodeData[5], nodeData[6], nodeData[7]}
	hash, err := mt.hashFunction(nd, capIn)
	if err != nil {
		return nil, err
	}

	err = mt.setNodeData(ctx, hash[:], append(nodeData, cap...))
	if err != nil {
		return nil, err
	}

	return hash[:], nil
}

func (mt *MerkleTree) setNodeData(ctx context.Context, key []uint64, data []uint64) error {
	var dataByte []byte
	for i := 0; i < len(data); i++ {
		e := fmt.Sprintf("%016s", strconv.FormatUint(data[i], 16))
		eByte, err := hex.DecodeHex(e)
		if err != nil {
			return err
		}
		dataByte = append(dataByte, eByte...)
	}
	return mt.store.Set(ctx, h4ToScalar(key).Bytes(), dataByte)
}

func (mt *MerkleTree) getNodeData(ctx context.Context, key []uint64) ([]uint64, error) {
	dataByte, err := mt.store.Get(ctx, h4ToScalar(key).Bytes())
	if err != nil {
		return nil, err
	}

	res := make([]uint64, 12)
	for i := 0; i < 12; i++ {
		res[i] = binary.BigEndian.Uint64(dataByte[i*8 : (i+1)*8])
	}
	return res[:], nil
}

// splitKey gets the path for a given key.
func (mt *MerkleTree) splitKey(key []uint64) []uint {
	var res []uint
	auxk := make([]uint64, 4)
	copy(auxk, key)
	for i := 0; i < 64; i++ {
		for j := 0; j < 4; j++ {
			res = append(res, uint(auxk[j]&1))
			auxk[j] = auxk[j] >> 1
		}
	}
	return res
}

// joinKey joins full key from remaining key and path already used.
func (mt *MerkleTree) joinKey(bits []uint64, k []uint64) []uint64 {
	n := make([]uint64, 4)
	accs := ff.NewElement()
	for i := 0; i < len(bits); i++ {
		if bits[i] == 1 {
			accs[i%4] = accs[i%4] | (1 << n[i%4])
		}
		n[i%4]++
	}
	auxk := make([]uint64, 4)
	for i := 0; i < 4; i++ {
		auxk[i] = k[i]<<n[i] | accs[i]
	}
	return auxk
}

// isOneSiblings checks if a node is a final node (final node: [n0, n1, ... n7, 1, 0, 0, 0])
func isOneSiblings(n []uint64) bool {
	return n[8] == 1
}

// nodeIsEq returns true if the node items are equal.
func nodeIsEq(x, y []uint64) bool {
	return x[0] == y[0] &&
		x[1] == y[1] &&
		x[2] == y[2] &&
		x[3] == y[3]
}

// removeKeyBits removes bits from the key depending on the smt level.
func removeKeyBits(key []uint64, nBits int) []uint64 {
	fullLevels := nBits / keyItems
	auxk := [4]uint64{key[0], key[1], key[2], key[3]}
	for i := 0; i < 4; i++ {
		n := fullLevels
		if fullLevels*4+i < nBits {
			n++
		}
		auxk[i] = auxk[i] >> n
	}
	return auxk[:]
}
