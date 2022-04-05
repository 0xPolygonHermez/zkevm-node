package tree

import (
	"context"
	"encoding/binary"
	"errors"
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
	cache          *nodeCache
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

// HashFunction is a function interface type to specify hash function that MT should use.
type HashFunction func(inp [poseidon.NROUNDSF]uint64, cap [poseidon.CAPLEN]uint64) ([poseidon.CAPLEN]uint64, error)

type scHashFunction func(code []byte) ([]uint64, error)

// NewMerkleTree creates new MerkleTree instance.
func NewMerkleTree(store Store, arity uint8) *MerkleTree {
	return &MerkleTree{
		store:          store,
		arity:          arity,
		hashFunction:   poseidon.Hash,
		scHashFunction: hashContractBytecode,
		cache:          newNodeCache(),
	}
}

// SupportsDBTransactions indicates whether the store implementation supports DB transactions.
func (mt *MerkleTree) SupportsDBTransactions() bool {
	return mt.store.SupportsDBTransactions()
}

// BeginDBTransaction starts a transaction block
func (mt *MerkleTree) BeginDBTransaction(ctx context.Context) error {
	mt.cache.init()

	return mt.store.BeginDBTransaction(ctx)
}

// Commit commits a db transaction
func (mt *MerkleTree) Commit(ctx context.Context) error {
	defer mt.cache.teardown()

	err := mt.writeCacheContents(ctx)
	if err != nil {
		return err
	}

	return mt.store.Commit(ctx)
}

// Rollback rollbacks a db transaction
func (mt *MerkleTree) Rollback(ctx context.Context) error {
	mt.cache.teardown()

	return mt.store.Rollback(ctx)
}

// Set method sets value of a leaf at a given location (key) in the tree.
//
// The algorithm works as follows:
// * First we iterate through the tree nodes from the root following the path
//   associated to the given key trying to determine if it already exists in the
//   tree.
// * Then, depending on the result of the previous operation, the leaf node is
//   updated/created/reset, with different actions for each of the following
//   cases:
//   * Update: key exists
//   * Insert with found key: part of the branch nodes in the path are already
//     created
//   * Insert without found key: none of the branch nodes in the path are already
//     created
//   * Delete: node value is empty and path exists.
// * Finally, the tree is traversed backwards from the leaf to the root updating
//   or creating all the affected branch nodes (including the root itself)
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

	// in this loop we iterate through the tree nodes from the root following the
	// path associated to the given key trying to determine if exists in the tree.
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

	// now we insert the leaf taking into account different cases: update, insert
	// with found key, insert without found key and delete.
	if !nodeIsZero(value) {
		if foundKey != nil {
			if nodeIsEq(key[:], foundKey) { // Update, key path exists
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
			} else { // insert with foundKey, part of the key path already exists
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
		} else { // insert without foundKey, key path is not present
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
		if (foundKey != nil) && nodeIsEq(key[:], foundKey) { // Delete, node value is empty and key path exists
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

						auxKey := make([]uint64, 4)
						for j := 0; j < 4; j++ {
							if len(accKey) > j {
								auxKey[j] = accKey[j]
							}
						}
						insKey = mt.joinKey([]uint64{auxKey[0], auxKey[1], auxKey[2], auxKey[3], uint64(uKey)}, rKey)
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
		} else { // nothing to do, node value is empty and key path doesn't exists
			mode = modeZeroToZero
		}
	}

	siblings = siblings[:level+1]

	// now we traverse the tree backwards from the leaf to the root updating
	// or creating all the affected branch nodes (including the root itself).
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

// Get method gets value at a given location in the tree.
//
// The algorithm has a single step, we traverse the tree from the root following
// the path assciated with the given key.
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

	// this loops iterates the tree nodes from the root following the path
	// associated with the given key.
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
		// the complete path or part of it is present in the tree.
		if nodeIsEq(key, foundKey) {
			// the complete path associated with the key is present, value found.
			value = foundVal
		} else {
			// only part of the path is present, update partial results related to
			// the existing path.
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
	if !mt.cache.isActive() {
		return mt.setStoreNodeData(ctx, key, data)
	}
	err := mt.cache.set(key, data)
	if err != nil {
		log.Errorf("Error setting data in MT node cache: %v", err)
	}
	return nil
}

func (mt *MerkleTree) setStoreNodeData(ctx context.Context, key []uint64, data []uint64) error {
	dataByte, err := uint64ToByte(data)
	if err != nil {
		return err
	}

	return mt.store.Set(ctx, h4ToScalar(key).Bytes(), dataByte)
}

func (mt *MerkleTree) getNodeData(ctx context.Context, key []uint64) ([]uint64, error) {
	var (
		dataByte []byte
		err      error
	)
	if mt.cache.isActive() {
		cachedData, err := mt.cache.get(key)
		if err != nil && !errors.Is(err, errMTNodeCacheItemNotFound) {
			log.Errorf("Error getting data from MT node cache: %v", err)
		}
		if cachedData != nil {
			dataByte, err = uint64ToByte(cachedData)
			if err != nil {
				log.Errorf("Error decoding MT node cache data: %v", err)
			}
		}
	}
	if len(dataByte) == 0 {
		log.Debugf("about to call store.Get...")
		dataByte, err = mt.store.Get(ctx, h4ToScalar(key).Bytes())
		if err != nil {
			return nil, err
		}
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

// writeCacheContents writes the contents of the cache to the data store.
func (mt *MerkleTree) writeCacheContents(ctx context.Context) error {
	for k, v := range mt.cache.data {
		key, err := stringToh4(k)
		if err != nil {
			return err
		}
		err = mt.setStoreNodeData(ctx, key, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func uint64ToByte(data []uint64) ([]byte, error) {
	var dataByte []byte
	for i := 0; i < len(data); i++ {
		e := fmt.Sprintf("%016s", strconv.FormatUint(data[i], 16))
		eByte, err := hex.DecodeHex(e)
		if err != nil {
			return nil, err
		}
		dataByte = append(dataByte, eByte...)
	}
	return dataByte, nil
}
