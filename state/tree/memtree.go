package tree

//
//// MemTree is a basic in-memory implementation of StateTree
//type MemTree struct {
//	//mt merkletree.Merkletree
//	mem map[Key][]byte
//}
//
//// NewMemTree creates new MemTree
//func NewMemTree() ReadWriter {
//	mem := make(map[Key][]byte)
//	return &MemTree{mem: mem}
//}
//
//// GetBalance returns balance
//func (tree *MemTree) GetBalance(address common.Address, root []byte) (*big.Int, error) {
//	if root != nil {
//		return nil, fmt.Errorf("not implemented")
//	}
//	key, err := GetKey(LeafTypeBalance, address, nil)
//	if err != nil {
//		return nil, err
//	}
//	val := tree.mem[key]
//	if val == nil {
//		return big.NewInt(0), nil
//	}
//	return big.NewInt(0).SetBytes(val), nil
//}
//
//// GetNonce returns nonce
//func (tree *MemTree) GetNonce(address common.Address, root []byte) (*big.Int, error) {
//	if root != nil {
//		return nil, fmt.Errorf("not implemented")
//	}
//	key, err := GetKey(LeafTypeNonce, address, nil)
//	if err != nil {
//		return nil, err
//	}
//	val := tree.mem[key]
//	if val == nil {
//		return big.NewInt(0), nil
//	}
//	return big.NewInt(0).SetBytes(val), nil
//}
//
//// GetCode returns code
//func (tree *MemTree) GetCode(address common.Address, root []byte) ([]byte, error) {
//	if root != nil {
//		return nil, fmt.Errorf("not implemented")
//	}
//	key, err := GetKey(LeafTypeCode, address, nil)
//	if err != nil {
//		return nil, err
//	}
//	val := tree.mem[key]
//	if len(val) == 0 {
//		return nil, nil
//	}
//	return val, nil
//}
//
//// GetStorageAt returns Storage Value at specified position
//func (tree *MemTree) GetStorageAt(address common.Address, position common.Hash, root []byte) (*big.Int, error) {
//	if root != nil {
//		return nil, fmt.Errorf("not implemented")
//	}
//	key, err := GetKey(LeafTypeStorage, address, position[:])
//	if err != nil {
//		return nil, err
//	}
//	val := tree.mem[key]
//	if val == nil {
//		return big.NewInt(0), nil
//	}
//	return big.NewInt(0).SetBytes(val), nil
//}
//
//// GetRoot returns current MerkleTree root hash
//func (tree *MemTree) GetRoot() ([]byte, error) {
//	return nil, fmt.Errorf("not implemented")
//}
//
//// GetRootForBatchNumber returns MerkleTree root for specified batchNumber
//func (tree *MemTree) GetRootForBatchNumber(batchNumber uint64) ([]byte, error) {
//	return nil, fmt.Errorf("not implemented")
//}
//
//// SetBalance sets balance
//func (tree *MemTree) SetBalance(address common.Address, balance *big.Int) (newRoot []byte, proof interface{}, err error) {
//	key, err := GetKey(LeafTypeBalance, address, nil)
//	if err != nil {
//		return nil, nil, err
//	}
//	tree.mem[key] = balance.Bytes()
//	return nil, nil, nil
//}
//
//// SetNonce sets nonce
//func (tree *MemTree) SetNonce(address common.Address, nonce *big.Int) (newRoot []byte, proof interface{}, err error) {
//	key, err := GetKey(LeafTypeNonce, address, nil)
//	if err != nil {
//		return nil, nil, err
//	}
//	tree.mem[key] = nonce.Bytes()
//	return nil, nil, nil
//}
//
//// SetCode sets code
//func (tree *MemTree) SetCode(address common.Address, code []byte) (newRoot []byte, proof interface{}, err error) {
//	key, err := GetKey(LeafTypeCode, address, nil)
//	if err != nil {
//		return nil, nil, err
//	}
//	tree.mem[key] = code
//	return nil, nil, nil
//}
//
//// SetStorageAt sets storage value at specified position
//func (tree *MemTree) SetStorageAt(address common.Address, position common.Hash, value *big.Int) (newRoot []byte, proof interface{}, err error) {
//	key, err := GetKey(LeafTypeStorage, address, position[:])
//	if err != nil {
//		return nil, nil, err
//	}
//	tree.mem[key] = value.Bytes()
//	return nil, nil, nil
//}
//
//// SetRootForBatchNumber sets root for specified batchNumber
//func (tree *MemTree) SetRootForBatchNumber(batchNumber uint64, root []byte) error {
//	return fmt.Errorf("not implemented")
//}
