package state

import (
	"context"
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
}

// SetStateRoot is the stateRoot setter.
func (f *FakeDB) SetStateRoot(stateRoot []byte) {
	f.stateRoot = stateRoot
}

func (f *FakeDB) CreateAccount(common.Address) {
	panic("not implemented yet")
}

func (f *FakeDB) SubBalance(common.Address, *big.Int) {
	panic("not implemented yet")
}
func (f *FakeDB) AddBalance(common.Address, *big.Int) {
	panic("not implemented yet")
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

func (f *FakeDB) SetNonce(common.Address, uint64) {
	panic("not implemented yet")
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

func (f *FakeDB) SetCode(common.Address, []byte) {
	panic("not implemented yet")
}
func (f *FakeDB) GetCodeSize(common.Address) int {
	panic("not implemented yet")
}

func (f *FakeDB) AddRefund(uint64) {
	panic("not implemented yet")
}
func (f *FakeDB) SubRefund(uint64) {
	panic("not implemented yet")
}
func (f *FakeDB) GetRefund() uint64 {
	panic("not implemented yet")
}

func (f *FakeDB) GetCommittedState(common.Address, common.Hash) common.Hash {
	panic("not implemented yet")
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

func (f *FakeDB) SetState(common.Address, common.Hash, common.Hash) {
	panic("not implemented yet")
}

func (f *FakeDB) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	panic("not implemented yet")
}
func (f *FakeDB) SetTransientState(addr common.Address, key, value common.Hash) {
	panic("not implemented yet")
}

func (f *FakeDB) Suicide(common.Address) bool {
	panic("not implemented yet")
}
func (f *FakeDB) HasSuicided(common.Address) bool {
	panic("not implemented yet")
}

// Exist reports whether the given account exists in state.
// Notably this should also return true for suicided accounts.
func (f *FakeDB) Exist(address common.Address) bool {
	return !(f.GetNonce(address) == 0 && f.GetBalance(address).Int64() == 0 && f.GetCodeHash(address) == ZeroHash)
}

// Empty returns whether the given account is empty. Empty
// is defined according to EIP161 (balance = nonce = code = 0).
func (f *FakeDB) Empty(common.Address) bool {
	panic("not implemented yet")
}

func (f *FakeDB) AddressInAccessList(addr common.Address) bool {
	panic("not implemented yet")
}
func (f *FakeDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	panic("not implemented yet")
}

// AddAddressToAccessList adds the given address to the access list. This operation is safe to perform
// even if the feature/fork is not active yet
func (f *FakeDB) AddAddressToAccessList(addr common.Address) {
	panic("not implemented yet")
}

// AddSlotToAccessList adds the given (address,slot) to the access list. This operation is safe to perform
// even if the feature/fork is not active yet
func (f *FakeDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	panic("not implemented yet")
}

func (f *FakeDB) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	panic("not implemented yet")
}

func (f *FakeDB) RevertToSnapshot(int) {
	panic("not implemented yet")
}
func (f *FakeDB) Snapshot() int {
	panic("not implemented yet")
}

func (f *FakeDB) AddLog(*types.Log) {
	panic("not implemented yet")
}
func (f *FakeDB) AddPreimage(common.Hash, []byte) {
	panic("not implemented yet")
}
