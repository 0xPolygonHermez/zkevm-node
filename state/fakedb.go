package state

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

// FakeDB is the implementation of the fakeevm.FakeDB interface
type FakeDB struct {
	State     *State
	stateRoot []byte
	refund    uint64
}

// SetStateRoot is the stateRoot setter.
func (f *FakeDB) SetStateRoot(stateRoot []byte) {
	f.stateRoot = stateRoot
}

// CreateAccount not implemented
func (f *FakeDB) CreateAccount(common.Address) {
	log.Error("FakeDB: CreateAccount method not implemented")
}

// SubBalance not implemented
func (f *FakeDB) SubBalance(common.Address, *big.Int) {
	log.Error("FakeDB: SubBalance method not implemented")
}

// AddBalance not implemented
func (f *FakeDB) AddBalance(common.Address, *big.Int) {
	log.Error("FakeDB: AddBalance method not implemented")
}

// GetBalance returns the balance of the given address.
func (f *FakeDB) GetBalance(address common.Address) *big.Int {
	ctx := context.Background()
	balance, err := f.State.GetTree().GetBalance(ctx, address, f.stateRoot)

	if err != nil {
		log.Errorf("error on FakeDB GetBalance for address %v", address)
	}

	log.Debugf("FakeDB GetBalance for address %v", address)
	return balance
}

// GetNonce returns the nonce of the given address.
func (f *FakeDB) GetNonce(address common.Address) uint64 {
	ctx := context.Background()
	nonce, err := f.State.GetTree().GetNonce(ctx, address, f.stateRoot)

	if err != nil {
		log.Errorf("error on FakeDB GetNonce for address %v", address)
		return 0
	}

	log.Debugf("FakeDB GetNonce for address %v", address)
	return nonce.Uint64()
}

// SetNonce not implemented
func (f *FakeDB) SetNonce(common.Address, uint64) {
	log.Error("FakeDB: SetNonce method not implemented")
}

// GetCodeHash gets the hash for the code at a given address
func (f *FakeDB) GetCodeHash(address common.Address) common.Hash {
	ctx := context.Background()
	hash, err := f.State.GetTree().GetCodeHash(ctx, address, f.stateRoot)
	if err != nil {
		log.Errorf("error on FakeDB GetCodeHash for address %v, err: %v", address, err)
	}

	log.Debugf("FakeDB GetCodeHash for address %v => %v", address, common.BytesToHash(hash))
	return common.BytesToHash(hash)
}

// GetCode returns the SC code of the given address.
func (f *FakeDB) GetCode(address common.Address) []byte {
	ctx := context.Background()
	code, err := f.State.GetTree().GetCode(ctx, address, f.stateRoot)

	if err != nil {
		log.Errorf("error on FakeDB GetCode for address %v", address)
	}

	log.Debugf("FakeDB GetCode for address %v", address)
	return code
}

// SetCode not implemented
func (f *FakeDB) SetCode(common.Address, []byte) {
	log.Error("FakeDB: SetCode method not implemented")
}

// GetCodeSize get address code size
func (f *FakeDB) GetCodeSize(address common.Address) int {
	return len(f.GetCode(address))
}

// AddRefund adds gas to the refund counter
func (f *FakeDB) AddRefund(gas uint64) {
	f.refund += gas
}

// SubRefund subtracts gas from the refund counter
func (f *FakeDB) SubRefund(gas uint64) {
	if gas > f.refund {
		log.Errorf(fmt.Sprintf("Refund counter below zero (gas: %d > refund: %d)", gas, f.refund))
	}
	f.refund -= gas
}

// GetRefund returns the refund counter
func (f *FakeDB) GetRefund() uint64 {
	return f.refund
}

// GetCommittedState not implemented
func (f *FakeDB) GetCommittedState(common.Address, common.Hash) common.Hash {
	log.Error("FakeDB: GetCommittedState method not implemented")
	return ZeroHash
}

// GetState retrieves a value from the given account's storage trie.
func (f *FakeDB) GetState(address common.Address, hash common.Hash) common.Hash {
	ctx := context.Background()
	storage, err := f.State.GetTree().GetStorageAt(ctx, address, hash.Big(), f.stateRoot)

	if err != nil {
		log.Errorf("error on FakeDB GetState for address %v", address)
	}

	log.Debugf("FakeDB GetState for address %v", address)

	return common.BytesToHash(storage.Bytes())
}

// SetState not implemented
func (f *FakeDB) SetState(common.Address, common.Hash, common.Hash) {
	log.Error("FakeDB: SetState method not implemented")
}

// GetTransientState not implemented
func (f *FakeDB) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	log.Error("FakeDB: GetTransientState method not implemented")
	return ZeroHash
}

// SetTransientState not implemented
func (f *FakeDB) SetTransientState(addr common.Address, key, value common.Hash) {
	log.Error("FakeDB: SetTransientState method not implemented")
}

// Suicide not implemented
func (f *FakeDB) Suicide(common.Address) bool {
	log.Error("FakeDB: Suicide method not implemented")
	return false
}

// HasSuicided not implemented
func (f *FakeDB) HasSuicided(common.Address) bool {
	log.Error("FakeDB: HasSuicided method not implemented")
	return false
}

// Exist reports whether the given account exists in state.
// Notably this should also return true for suicided accounts.
func (f *FakeDB) Exist(address common.Address) bool {
	return !(f.GetNonce(address) == 0 && f.GetBalance(address).Int64() == 0 && f.GetCodeHash(address) == ZeroHash)
}

// Empty returns whether the given account is empty. Empty
// is defined according to EIP161 (balance = nonce = code = 0).
func (f *FakeDB) Empty(address common.Address) bool {
	return !(f.GetNonce(address) == 0 && f.GetBalance(address).Int64() == 0 && f.GetCodeHash(address) == ZeroHash)
}

// AddressInAccessList not implemented
func (f *FakeDB) AddressInAccessList(addr common.Address) bool {
	log.Error("FakeDB: AddressInAccessList method not implemented")
	return false
}

// SlotInAccessList not implemented
func (f *FakeDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	log.Error("FakeDB: SlotInAccessList method not implemented")
	return false, false
}

// AddAddressToAccessList adds the given address to the access list. This operation is safe to perform

// AddAddressToAccessList not implemented// even if the feature/fork is not active yet
func (f *FakeDB) AddAddressToAccessList(addr common.Address) {
	log.Error("FakeDB: AddAddressToAccessList method not implemented")
}

// AddSlotToAccessList adds the given (address,slot) to the access list. This operation is safe to perform

// AddSlotToAccessList not implemented// even if the feature/fork is not active yet
func (f *FakeDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	log.Error("FakeDB: AddSlotToAccessList method not implemented")
}

// Prepare not implemented
func (f *FakeDB) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	log.Error("FakeDB: Prepare method not implemented")
}

// RevertToSnapshot not implemented
func (f *FakeDB) RevertToSnapshot(int) {
	log.Error("FakeDB: RevertToSnapshot method not implemented")
}

// Snapshot not implemented
func (f *FakeDB) Snapshot() int {
	log.Error("FakeDB: Snapshot method not implemented")
	return 0
}

// AddLog not implemented
func (f *FakeDB) AddLog(*types.Log) {
	log.Error("FakeDB: AddLog method not implemented")
}

// AddPreimage not implemented
func (f *FakeDB) AddPreimage(common.Hash, []byte) {
	log.Error("FakeDB: AddPreimage method not implemented")
}
