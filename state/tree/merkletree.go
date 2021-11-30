package tree

import (
	"fmt"
	"github.com/hermeznetwork/hermez-core/state/db"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"math/big"
)

const (
	cmpLt          = -1
	cmpEq          = 0
	cmpGt          = 1
	bigIntMaxBytes = 32
)

type MerkleTree struct {
	db           db.KeyValuer
	hashFunction interface{}
	arity        uint8
	maxLevels    uint16
	f            interface{}
	mask         *big.Int
}

type UpdateProof struct {
	OldRoot  *big.Int
	NewRoot  *big.Int
	Key      *big.Int
	Siblings [][]*big.Int
	InsKey   *big.Int
	InsValue *big.Int
	IsOld0   bool
	OldValue *big.Int
	NewValue *big.Int
}

type Proof struct {
	Root     *big.Int
	Key      *big.Int
	Value    *big.Int
	Siblings [][]*big.Int
	IsOld0   bool
	InsKey   *big.Int
	InsValue *big.Int
}

func NewMerkleTree(db db.KeyValuer, arity uint8, hash interface{}, F interface{}) *MerkleTree {
	return &MerkleTree{
		db:           db,
		arity:        arity,
		hashFunction: hash,
		f:            F,
		mask:         big.NewInt(1<<arity - 1),
		maxLevels:    uint16(160 / arity),
	}
}

func (mt *MerkleTree) Set(oldRoot, key, value *big.Int) (*UpdateProof, error) {
	var err error

	r := oldRoot
	keys := mt.splitKey(key)
	level := 0

	zero := big.NewInt(0)
	one := big.NewInt(1)

	accKey := big.NewInt(0)
	lastAccKey := big.NewInt(0)
	var foundKey *big.Int
	var siblings [][]*big.Int

	var insKey *big.Int
	var insValue *big.Int
	oldValue := big.NewInt(0)
	mode := ""
	newRoot := oldRoot
	isOld0 := true

	for (r.Cmp(zero) != cmpEq) && (foundKey != nil) {
		node, err := mt.getNode(r)
		if err != nil {
			return nil, err
		}
		siblings[level] = node

		if siblings[level][0].Cmp(one) == cmpEq {
			foundKey = new(big.Int).Add(
				accKey,
				new(big.Int).Mul(
					siblings[level][1],
					new(big.Int).Lsh(one, uint(level*int(mt.arity))),
				),
			)
		} else {
			r = siblings[level][keys[level]]
			lastAccKey = accKey
			accKey = new(big.Int).Add(accKey, new(big.Int).Lsh(big.NewInt(int64(keys[level])), uint(level*int(mt.arity))))
			level++
		}
	}

	level--
	accKey = lastAccKey

	if value.Cmp(zero) != cmpEq {
		v, err := scalar2fea(value)
		if err != nil {
			return nil, err
		}

		if foundKey != nil {
			if key.Cmp(foundKey) == cmpEq { // Update
				mode = "update"
				newLeaf := mt.newNode()

				newLeaf[0] = one
				newLeaf[1] = siblings[level+1][1]
				oldValue = fea2scalar(siblings[level+1][2:6])
				newLeaf[2] = v[0]
				newLeaf[3] = v[1]
				newLeaf[4] = v[2]
				newLeaf[5] = v[3]

				newLeafHash, err := mt.hashSave(newLeaf)
				if err != nil {
					return nil, err
				}

				if level >= 0 {
					siblings[level][keys[level]] = newLeafHash
				} else {
					newRoot = newLeafHash
				}
			} else { // insert with foundKey
				mode = "insertFound"
				node := mt.newNode()
				level2 := level + 1
				foundKeys := mt.splitKey(foundKey)
				for keys[level2] == foundKeys[level2] {
					level2++
				}

				oldLeaf := mt.newNode()
				oldLeaf[0] = one
				//oldLeaf[1] = F.e(Scalar.shr(Scalar.e(F.toObject(foundKey)), (level2+1)*mt.arity))
				oldLeaf[1] = new(big.Int).Rsh(foundKey, uint((level2+1)*int(mt.arity)))
				oldLeaf[2] = siblings[level+1][2]
				oldLeaf[3] = siblings[level+1][3]
				oldLeaf[4] = siblings[level+1][4]
				oldLeaf[5] = siblings[level+1][5]

				insKey = foundKey
				insValue = fea2scalar(siblings[level+1][2:6])
				isOld0 = false
				oldLeafHash, err := mt.hashSave(oldLeaf)
				if err != nil {
					return nil, err
				}

				newLeaf := mt.newNode()
				newLeaf[0] = one
				newLeaf[1] = new(big.Int).Rsh(key, uint((level2+1)*int(mt.arity)))
				newLeaf[2] = v[0]
				newLeaf[3] = v[1]
				newLeaf[4] = v[2]
				newLeaf[5] = v[3]

				newLeafHash, err := mt.hashSave(newLeaf)
				if err != nil {
					return nil, err
				}

				node[keys[level2]] = newLeafHash
				node[foundKeys[level2]] = oldLeafHash

				r2, err := mt.hashSave(node)
				if err != nil {
					return nil, err
				}
				level2--

				for level2 != level {
					for i := 0; i < (1 << mt.arity); i++ {
						node[i] = big.NewInt(0)
					}
					node[keys[level2]] = r2

					r2, err = mt.hashSave(node)
					if err != nil {
						return nil, err
					}
					level2--
				}

				if level >= 0 {
					siblings[level][keys[level]] = r2
				} else {
					newRoot = r2
				}

			}
		} else { // insert without foundKey
			mode = "insertNotFound"
			newLeaf := mt.newNode()
			newLeaf[0] = one
			newLeaf[1] = new(big.Int).Rsh(key, uint((level+1)*int(mt.arity)))
			newLeaf[2] = v[0]
			newLeaf[3] = v[1]
			newLeaf[4] = v[2]
			newLeaf[5] = v[3]
			newLeafHash, err := mt.hashSave(newLeaf)
			if err != nil {
				return nil, err
			}
			if level >= 0 {
				siblings[level][keys[level]] = newLeafHash
			} else {
				newRoot = newLeafHash
			}
		}
	} else {
		if (foundKey != nil) && (key.Cmp(foundKey) == cmpEq) { // Delete

			oldValue = fea2scalar(siblings[level+1][2:6])
			if level >= 0 {
				siblings[level][keys[level]] = zero

				uKey := mt.getUniqueSibling(siblings[level])

				if uKey >= 0 {
					mode = "deleteFound"
					node, err := mt.getNode(siblings[level][uKey])
					if err != nil {
						return nil, err
					}
					siblings[level+1] = node

					insKey = new(big.Int).Add(
						new(big.Int).Add(accKey, new(big.Int).Lsh(big.NewInt(int64(uKey)), uint(level*int(mt.arity)))),
						new(big.Int).Mul(
							siblings[level+1][1],
							new(big.Int).Lsh(one, uint((level+1)*int(mt.arity))),
						),
					)
					insV := siblings[level+1][2:6]
					insValue = fea2scalar(insV)
					isOld0 = false

					for (uKey >= 0) && (level >= 0) {
						level--
						if level >= 0 {
							uKey = mt.getUniqueSibling(siblings[level])
						}
					}

					oldLeaf := mt.newNode()
					oldLeaf[0] = one
					oldLeaf[1] = new(big.Int).Rsh(insKey, uint((level+1)*int(mt.arity)))
					oldLeaf[2] = insV[0]
					oldLeaf[3] = insV[1]
					oldLeaf[4] = insV[2]
					oldLeaf[5] = insV[3]
					oldLeafHash, err := mt.hashSave(oldLeaf)
					if err != nil {
						return nil, err
					}

					if level >= 0 {
						siblings[level][keys[level]] = oldLeafHash
					} else {
						newRoot = oldLeafHash
					}
				} else {
					mode = "deleteNotFound"
				}
			} else {
				mode = "deleteLast"
				newRoot = zero
			}
		} else {
			mode = "zeroToZero"
		}
	}

	siblings = siblings[0 : level+1]

	for level >= 0 {
		newRoot, err = mt.hashSave(siblings[level])
		if err != nil {
			return nil, err
		}
		level--
		if level >= 0 {
			siblings[level][keys[level]] = newRoot
		}
	}

	fmt.Println("mode: ", mode)

	return &UpdateProof{
		OldRoot:  oldRoot,
		NewRoot:  newRoot,
		Key:      key,
		Siblings: siblings,
		InsKey:   insKey,
		InsValue: insValue,
		IsOld0:   isOld0,
		OldValue: oldValue,
		NewValue: value,
	}, nil
}

func (mt *MerkleTree) Get(root, key *big.Int) (*Proof, error) {
	r := root

	keys := mt.splitKey(key)
	level := 0

	zero := big.NewInt(0)
	one := big.NewInt(1)

	accKey := big.NewInt(0)
	lastAccKey := big.NewInt(0)
	var foundKey *big.Int
	var siblings [][]*big.Int

	insKey := big.NewInt(0)
	insValue := big.NewInt(0)
	value := big.NewInt(0)

	isOld0 := true

	for (r.Cmp(zero) != cmpEq) && (foundKey != nil) {
		node, err := mt.getNode(r)
		if err != nil {
			return nil, err
		}
		siblings[level] = node

		if siblings[level][0].Cmp(one) == cmpEq {
			foundKey = new(big.Int).Add(
				accKey,
				new(big.Int).Mul(
					siblings[level][1],
					new(big.Int).Lsh(one, uint(level*int(mt.arity))),
				),
			)
		} else {
			r = siblings[level][keys[level]]
			lastAccKey = accKey
			accKey = new(big.Int).Add(accKey, new(big.Int).Lsh(big.NewInt(int64(keys[level])), uint(level*int(mt.arity))))
			level++
		}
	}

	level--
	accKey = lastAccKey

	if foundKey != nil {
		if key.Cmp(foundKey) == cmpEq {
			value = fea2scalar(siblings[level+1][2:6])
		} else {
			insKey = foundKey
			insValue = fea2scalar(siblings[level+1][2:6])
			isOld0 = false
		}
	}

	siblings = siblings[0 : level+1]

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

func (mt *MerkleTree) getUniqueSibling(a []*big.Int) int64 {
	nFound := 0
	zero := big.NewInt(0)
	var fnd int64
	for i := 0; i < len(a); i++ {
		if a[i].Cmp(zero) != cmpEq {
			nFound++
			fnd = int64(i)
		}
	}
	if nFound == 1 {
		return fnd
	}
	return -1
}

func (mt *MerkleTree) splitKey(key *big.Int) []uint {
	var res []uint
	auxk := key
	for i := 0; i < int(mt.maxLevels); i++ {
		res = append(res, uint(new(big.Int).And(auxk, mt.mask).Uint64()))
		auxk = new(big.Int).Rsh(auxk, uint(mt.arity))
	}
	return res
}

func (mt *MerkleTree) hashSave(a []*big.Int) (*big.Int, error) {
	hash, err := poseidon.Hash(a)
	if err != nil {
		return nil, err
	}

	// TODO: save to db

	return hash, nil
}

func (mt *MerkleTree) newNode() []*big.Int {
	node := make([]*big.Int, 1<<mt.arity)
	for i := 0; i < 1<<mt.arity; i++ {
		node[i] = big.NewInt(0)
	}
	return node
}

func (mt *MerkleTree) getNode(key *big.Int) ([]*big.Int, error) {
	nodeBytes, err := mt.db.Get(key.String())
	if err != nil {
		return nil, err
	}
	node := mt.newNode()
	for i := 0; i < 1<<int(mt.arity); i++ {
		node[i] = new(big.Int).SetBytes(nodeBytes[i*bigIntMaxBytes : (i+1)*bigIntMaxBytes])
	}
	return node, nil
}
