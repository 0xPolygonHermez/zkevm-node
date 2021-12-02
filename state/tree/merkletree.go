package tree

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hermeznetwork/hermez-core/hex"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/jackc/pgx/v4/pgxpool"
	"math/big"
)

const (
	cmpLt = -1
	cmpEq = 0
	cmpGt = 1
)

const (
	getNodeByKeySQL    = "SELECT key, data FROM merkletree WHERE key = $1"
	setNodeByKeySQL    = "INSERT INTO merkletree (key, data) VALUES ($1, $2)"
	deleteNodeByKeySQL = "DELETE FROM merkletree WHERE key = $1"
	checkNodeExistsSQL = "SELECT COUNT(*) as exists FROM merkletree WHERE key = $1"
)

type MerkleTree struct {
	db           *pgxpool.Pool
	hashFunction interface{}
	arity        uint8
	maxLevels    uint16
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

// NodeItem is a model of node stored in the database
type NodeItem struct {
	Key  []byte
	Data []byte
}

func NewMerkleTree(db *pgxpool.Pool, arity uint8, hash interface{}) *MerkleTree {
	return &MerkleTree{
		db:           db,
		arity:        arity,
		hashFunction: hash,
		mask:         big.NewInt(1<<arity - 1),
		maxLevels:    uint16(160 / arity),
	}
}

func (mt *MerkleTree) Set(ctx context.Context, oldRoot *big.Int, key *big.Int, value *big.Int) (*UpdateProof, error) {
	var err error

	// exit early if context is cancelled
	err = ctx.Err()
	if err != nil {
		return nil, err
	}

	zero := big.NewInt(0)
	one := big.NewInt(1)

	r := big.NewInt(0)

	if oldRoot != nil {
		r = oldRoot
	}
	keys := mt.splitKey(key)
	level := 0

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
		node, err := mt.getNodeData(ctx, r)
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
				newLeaf := mt.newNodeData()

				newLeaf[0] = one
				newLeaf[1] = siblings[level+1][1]
				oldValue = fea2scalar(siblings[level+1][2:6])
				newLeaf[2] = v[0]
				newLeaf[3] = v[1]
				newLeaf[4] = v[2]
				newLeaf[5] = v[3]

				newLeafHash, err := mt.hashSave(ctx, newLeaf)
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
				node := mt.newNodeData()
				level2 := level + 1
				foundKeys := mt.splitKey(foundKey)
				for keys[level2] == foundKeys[level2] {
					level2++
				}

				oldLeaf := mt.newNodeData()
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
				oldLeafHash, err := mt.hashSave(ctx, oldLeaf)
				if err != nil {
					return nil, err
				}

				newLeaf := mt.newNodeData()
				newLeaf[0] = one
				newLeaf[1] = new(big.Int).Rsh(key, uint((level2+1)*int(mt.arity)))
				newLeaf[2] = v[0]
				newLeaf[3] = v[1]
				newLeaf[4] = v[2]
				newLeaf[5] = v[3]

				newLeafHash, err := mt.hashSave(ctx, newLeaf)
				if err != nil {
					return nil, err
				}

				node[keys[level2]] = newLeafHash
				node[foundKeys[level2]] = oldLeafHash

				r2, err := mt.hashSave(ctx, node)
				if err != nil {
					return nil, err
				}
				level2--

				for level2 != level {
					for i := 0; i < (1 << mt.arity); i++ {
						node[i] = big.NewInt(0)
					}
					node[keys[level2]] = r2

					r2, err = mt.hashSave(ctx, node)
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
			newLeaf := mt.newNodeData()
			newLeaf[0] = one
			newLeaf[1] = new(big.Int).Rsh(key, uint((level+1)*int(mt.arity)))
			newLeaf[2] = v[0]
			newLeaf[3] = v[1]
			newLeaf[4] = v[2]
			newLeaf[5] = v[3]
			newLeafHash, err := mt.hashSave(ctx, newLeaf)
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
					node, err := mt.getNodeData(ctx, siblings[level][uKey])
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

					oldLeaf := mt.newNodeData()
					oldLeaf[0] = one
					oldLeaf[1] = new(big.Int).Rsh(insKey, uint((level+1)*int(mt.arity)))
					oldLeaf[2] = insV[0]
					oldLeaf[3] = insV[1]
					oldLeaf[4] = insV[2]
					oldLeaf[5] = insV[3]
					oldLeafHash, err := mt.hashSave(ctx, oldLeaf)
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
		newRoot, err = mt.hashSave(ctx, siblings[level])
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

func (mt *MerkleTree) Get(ctx context.Context, root, key *big.Int) (*Proof, error) {
	// exit early if context is cancelled
	err := ctx.Err()
	if err != nil {
		return nil, err
	}

	zero := big.NewInt(0)
	one := big.NewInt(1)

	r := big.NewInt(0)

	if root != nil {
		r = zero
	}

	keys := mt.splitKey(key)
	level := 0

	accKey := big.NewInt(0)
	lastAccKey := big.NewInt(0)
	var foundKey *big.Int
	var siblings [][]*big.Int

	insKey := big.NewInt(0)
	insValue := big.NewInt(0)
	value := big.NewInt(0)

	isOld0 := true

	for (r.Cmp(zero) != cmpEq) && (foundKey != nil) {
		node, err := mt.getNodeData(ctx, r)
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

func (mt *MerkleTree) hashSave(ctx context.Context, nodeData []*big.Int) (*big.Int, error) {
	var data []*big.Int
	//for i := range nodeData {
	//	data = append(data, new(big.Int).SetBytes(utils.SwapEndianness(nodeData[i].Bytes())))
	//}

	data = nodeData

	hash, err := poseidon.Hash(data)
	if err != nil {
		return nil, err
	}

	//hash.SetBytes(utils.SwapEndianness(hash.Bytes()))

	fmt.Printf("set node Key: %+v Data: %+v\n", hex.EncodeToString(hash.Bytes()), data)

	//err = mt.setNodeData(ctx, hash, nodeData)
	err = mt.setNodeData(ctx, hash, data)
	if err != nil {
		return nil, err
	}

	//hash.SetBytes(utils.SwapEndianness(hash.Bytes()))

	return hash, nil
}

func (mt *MerkleTree) newNodeData() []*big.Int {
	node := make([]*big.Int, 1<<mt.arity)
	for i := 0; i < 1<<mt.arity; i++ {
		node[i] = big.NewInt(0)
	}
	return node
}

func (mt *MerkleTree) getNodeData(ctx context.Context, key *big.Int) ([]*big.Int, error) {
	fmt.Printf("Get node Key: %+v\n", hex.EncodeToHex(key.Bytes()))
	var node NodeItem
	err := mt.db.QueryRow(ctx, getNodeByKeySQL, key.Bytes()).Scan(&node.Key, &node.Data)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	// parse bytes into []*big.Int
	nodeData := mt.newNodeData()
	for i := 0; i < len(node.Data)/maxBigIntLen; i++ {
		nodeData[i] = new(big.Int).SetBytes(node.Data[i*maxBigIntLen : (i+1)*maxBigIntLen])
	}
	fmt.Printf("Got nodeData: %+v\n", nodeData)
	return nodeData, nil
}

func (mt *MerkleTree) setNodeData(ctx context.Context, key *big.Int, data []*big.Int) error {
	var exists int
	err := mt.db.QueryRow(ctx, checkNodeExistsSQL, key.Bytes()).Scan(&exists)
	if exists != 0 {
		fmt.Println("Item already exists, key: ", hex.EncodeToHex(key.Bytes()))
		// item already exists, no need to do anything
		//return nil
		_, err = mt.db.Exec(ctx, deleteNodeByKeySQL, key.Bytes())
		if err != nil {
			return err
		}
		fmt.Println("Item deleted, key: ", hex.EncodeToHex(key.Bytes()))
	}

	var buf bytes.Buffer
	for i := 0; i < len(data); i++ {
		var b [maxBigIntLen]byte
		d := data[i].FillBytes(b[:])
		buf.Write(d)
		//fmt.Println(hex.EncodeToHex(b[:]))
	}
	fmt.Printf("Set node Key: %+v Data: %+v\n", hex.EncodeToHex(key.Bytes()), hex.EncodeToHex(buf.Bytes()))
	// insert node into the database
	_, err = mt.db.Exec(ctx, setNodeByKeySQL, key.Bytes(), buf.Bytes())
	if err != nil {
		return err
	}

	_, err = mt.getNodeData(ctx, key)
	if err != nil {
		return err
	}
	//fmt.Printf("check got data form db: ")

	return nil
}
