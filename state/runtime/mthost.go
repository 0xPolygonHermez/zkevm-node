package runtime

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/tree"
)

type MtHost struct {
	tree      tree.ReadWriter
	stateRoot []byte
}

// NewMerkleTree creates new MerkleTree instance
func NewMerkleTreeHost(tree tree.ReadWriter) *MtHost {
	root, err := tree.GetCurrentRoot()
	if err != nil {
		log.Error("error creating NewMerkleTreeHost", err)
		return nil
	}
	return &MtHost{tree: tree, stateRoot: root}
}

func (m *MtHost) AccountExists(address common.Address) bool {
	panic("AccountExists not implemented ")
}

func (m *MtHost) GetStorage(address common.Address, key common.Hash) common.Hash {
	storage, err := m.tree.GetStorageAt(address, key, m.stateRoot)
	if err != nil {
		log.Error(err)
		return common.BigToHash(new(big.Int))
	}
	return common.BigToHash(storage)
}

func (m *MtHost) SetStorage(address common.Address, key common.Hash, value common.Hash, config *ForksInTime) StorageStatus {
	root, _, err := m.tree.SetStorageAt(address, key, value.Big())
	if err != nil {
		log.Error("error setting storage", err)
	}
	m.stateRoot = root

	// TODO: Calculate this
	return StorageModified
}

func (m *MtHost) GetBalance(address common.Address) *big.Int {
	balance, err := m.tree.GetBalance(address, m.stateRoot)
	if err != nil {
		log.Error("error getting balance", err)
		return nil
	}
	return balance
}

func (m *MtHost) GetCodeSize(address common.Address) int {
	code, err := m.tree.GetCode(address, m.stateRoot)
	if err != nil {
		log.Error("error getting code size", err)
		return 0
	}
	return len(code)
}

func (m *MtHost) GetCodeHash(address common.Address) common.Hash {
	panic("GetCodeHash not implemented ")
}

func (m *MtHost) GetCode(address common.Address) []byte {
	code, err := m.tree.GetCode(address, m.stateRoot)
	if err != nil {
		log.Error("error getting code", err)
		return nil
	}
	return code
}

func (m *MtHost) Selfdestruct(address common.Address, beneficiary common.Address) {
	panic("Selfdestruct not implemented ")
}

func (m *MtHost) GetTxContext() TxContext {
	panic("GetTxContext not implemented ")
}

func (m *MtHost) GetBlockHash(number int64) common.Hash {
	panic("GetBlockHash not implemented ")
}

func (m *MtHost) EmitLog(address common.Address, topics []common.Hash, data []byte) {
	panic("EmitLog not implemented ")
}

func (m *MtHost) Callx(*Contract, Host) *ExecutionResult {
	panic("Callx not implemented ")
}

func (m *MtHost) Empty(address common.Address) bool {
	panic("Empty not implemented ")
}

func (m *MtHost) GetNonce(address common.Address) uint64 {
	nonce, err := m.tree.GetNonce(address, m.stateRoot)
	if err != nil {
		log.Error("error getting nonce", err)
		return 0
	}
	return nonce.Uint64()
}
