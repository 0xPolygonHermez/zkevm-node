// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ChainCallLevel3

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ChainCallLevel3MetaData contains all meta data concerning the ChainCallLevel3 contract.
var ChainCallLevel3MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"callRevert\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"delegateCallRevert\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"exec\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"get\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"t\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610631806100206000396000f3fe6080604052600436106100435760003560e01c80636bb6126e1461004f578063937de75414610064578063c2bc2efc14610077578063df0a5d56146100ad57600080fd5b3661004a57005b600080fd5b61006261005d366004610448565b6100c0565b005b610062610072366004610448565b610258565b34801561008357600080fd5b50610097610092366004610448565b61029a565b6040516100a491906104a8565b60405180910390f35b6100626100bb366004610448565b6103a1565b60408051600481526024810182526020810180516001600160e01b03166330703a7160e21b17905290516000916001600160a01b0384169161010291906104db565b6000604051808303816000865af19150503d806000811461013f576040519150601f19603f3d011682016040523d82523d6000602084013e610144565b606091505b505080915050806101705760405162461bcd60e51b8152600401610167906104f7565b60405180910390fd5b60408051600481526024810182526020810180516001600160e01b03166330703a7160e21b17905290516001600160a01b038416916101ae916104db565b600060405180830381855af49150503d80600081146101e9576040519150601f19603f3d011682016040523d82523d6000602084013e6101ee565b606091505b505080915050806102545760405162461bcd60e51b815260206004820152602a60248201527f6661696c656420746f20706572666f726d2064656c65676174652063616c6c206044820152691d1bc81b195d995b080d60b21b6064820152608401610167565b5050565b60408051600481526024810182526020810180516001600160e01b031663537669ab60e11b17905290516000916001600160a01b038416916101ae91906104db565b60408051600481526024810182526020810180516001600160e01b0316631b53398f60e21b179052905160609160009183916001600160a01b038616916102e191906104db565b600060405180830381855afa9150503d806000811461031c576040519150601f19603f3d011682016040523d82523d6000602084013e610321565b606091505b509092509050816103855760405162461bcd60e51b815260206004820152602860248201527f6661696c656420746f20706572666f726d207374617469632063616c6c20746f604482015267081b195d995b080d60c21b6064820152608401610167565b80806020019051810190610399919061054e565b949350505050565b60408051600481526024810182526020810180516001600160e01b031663537669ab60e11b17905290516000916001600160a01b038416916103e391906104db565b6000604051808303816000865af19150503d8060008114610420576040519150601f19603f3d011682016040523d82523d6000602084013e610425565b606091505b505080915050806102545760405162461bcd60e51b8152600401610167906104f7565b60006020828403121561045a57600080fd5b81356001600160a01b038116811461047157600080fd5b9392505050565b60005b8381101561049357818101518382015260200161047b565b838111156104a2576000848401525b50505050565b60208152600082518060208401526104c7816040850160208701610478565b601f01601f19169190910160400192915050565b600082516104ed818460208701610478565b9190910192915050565b60208082526021908201527f6661696c656420746f20706572666f726d2063616c6c20746f206c6576656c206040820152600d60fa1b606082015260800190565b634e487b7160e01b600052604160045260246000fd5b60006020828403121561056057600080fd5b815167ffffffffffffffff8082111561057857600080fd5b818401915084601f83011261058c57600080fd5b81518181111561059e5761059e610538565b604051601f8201601f19908116603f011681019083821181831017156105c6576105c6610538565b816040528281528760208487010111156105df57600080fd5b6105f0836020830160208801610478565b97965050505050505056fea2646970667358221220f3aa6c56dac196edcf14769b375fa15c320e946e4466ab35ec78a66e0dd21fcd64736f6c634300080c0033",
}

// ChainCallLevel3ABI is the input ABI used to generate the binding from.
// Deprecated: Use ChainCallLevel3MetaData.ABI instead.
var ChainCallLevel3ABI = ChainCallLevel3MetaData.ABI

// ChainCallLevel3Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ChainCallLevel3MetaData.Bin instead.
var ChainCallLevel3Bin = ChainCallLevel3MetaData.Bin

// DeployChainCallLevel3 deploys a new Ethereum contract, binding an instance of ChainCallLevel3 to it.
func DeployChainCallLevel3(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ChainCallLevel3, error) {
	parsed, err := ChainCallLevel3MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChainCallLevel3Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChainCallLevel3{ChainCallLevel3Caller: ChainCallLevel3Caller{contract: contract}, ChainCallLevel3Transactor: ChainCallLevel3Transactor{contract: contract}, ChainCallLevel3Filterer: ChainCallLevel3Filterer{contract: contract}}, nil
}

// ChainCallLevel3 is an auto generated Go binding around an Ethereum contract.
type ChainCallLevel3 struct {
	ChainCallLevel3Caller     // Read-only binding to the contract
	ChainCallLevel3Transactor // Write-only binding to the contract
	ChainCallLevel3Filterer   // Log filterer for contract events
}

// ChainCallLevel3Caller is an auto generated read-only Go binding around an Ethereum contract.
type ChainCallLevel3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainCallLevel3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ChainCallLevel3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainCallLevel3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChainCallLevel3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainCallLevel3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChainCallLevel3Session struct {
	Contract     *ChainCallLevel3  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChainCallLevel3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChainCallLevel3CallerSession struct {
	Contract *ChainCallLevel3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ChainCallLevel3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChainCallLevel3TransactorSession struct {
	Contract     *ChainCallLevel3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ChainCallLevel3Raw is an auto generated low-level Go binding around an Ethereum contract.
type ChainCallLevel3Raw struct {
	Contract *ChainCallLevel3 // Generic contract binding to access the raw methods on
}

// ChainCallLevel3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChainCallLevel3CallerRaw struct {
	Contract *ChainCallLevel3Caller // Generic read-only contract binding to access the raw methods on
}

// ChainCallLevel3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChainCallLevel3TransactorRaw struct {
	Contract *ChainCallLevel3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewChainCallLevel3 creates a new instance of ChainCallLevel3, bound to a specific deployed contract.
func NewChainCallLevel3(address common.Address, backend bind.ContractBackend) (*ChainCallLevel3, error) {
	contract, err := bindChainCallLevel3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel3{ChainCallLevel3Caller: ChainCallLevel3Caller{contract: contract}, ChainCallLevel3Transactor: ChainCallLevel3Transactor{contract: contract}, ChainCallLevel3Filterer: ChainCallLevel3Filterer{contract: contract}}, nil
}

// NewChainCallLevel3Caller creates a new read-only instance of ChainCallLevel3, bound to a specific deployed contract.
func NewChainCallLevel3Caller(address common.Address, caller bind.ContractCaller) (*ChainCallLevel3Caller, error) {
	contract, err := bindChainCallLevel3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel3Caller{contract: contract}, nil
}

// NewChainCallLevel3Transactor creates a new write-only instance of ChainCallLevel3, bound to a specific deployed contract.
func NewChainCallLevel3Transactor(address common.Address, transactor bind.ContractTransactor) (*ChainCallLevel3Transactor, error) {
	contract, err := bindChainCallLevel3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel3Transactor{contract: contract}, nil
}

// NewChainCallLevel3Filterer creates a new log filterer instance of ChainCallLevel3, bound to a specific deployed contract.
func NewChainCallLevel3Filterer(address common.Address, filterer bind.ContractFilterer) (*ChainCallLevel3Filterer, error) {
	contract, err := bindChainCallLevel3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel3Filterer{contract: contract}, nil
}

// bindChainCallLevel3 binds a generic wrapper to an already deployed contract.
func bindChainCallLevel3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChainCallLevel3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainCallLevel3 *ChainCallLevel3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainCallLevel3.Contract.ChainCallLevel3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainCallLevel3 *ChainCallLevel3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainCallLevel3.Contract.ChainCallLevel3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainCallLevel3 *ChainCallLevel3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainCallLevel3.Contract.ChainCallLevel3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainCallLevel3 *ChainCallLevel3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainCallLevel3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainCallLevel3 *ChainCallLevel3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainCallLevel3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainCallLevel3 *ChainCallLevel3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainCallLevel3.Contract.contract.Transact(opts, method, params...)
}

// Get is a free data retrieval call binding the contract method 0xc2bc2efc.
//
// Solidity: function get(address level4Addr) view returns(string t)
func (_ChainCallLevel3 *ChainCallLevel3Caller) Get(opts *bind.CallOpts, level4Addr common.Address) (string, error) {
	var out []interface{}
	err := _ChainCallLevel3.contract.Call(opts, &out, "get", level4Addr)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Get is a free data retrieval call binding the contract method 0xc2bc2efc.
//
// Solidity: function get(address level4Addr) view returns(string t)
func (_ChainCallLevel3 *ChainCallLevel3Session) Get(level4Addr common.Address) (string, error) {
	return _ChainCallLevel3.Contract.Get(&_ChainCallLevel3.CallOpts, level4Addr)
}

// Get is a free data retrieval call binding the contract method 0xc2bc2efc.
//
// Solidity: function get(address level4Addr) view returns(string t)
func (_ChainCallLevel3 *ChainCallLevel3CallerSession) Get(level4Addr common.Address) (string, error) {
	return _ChainCallLevel3.Contract.Get(&_ChainCallLevel3.CallOpts, level4Addr)
}

// CallRevert is a paid mutator transaction binding the contract method 0xdf0a5d56.
//
// Solidity: function callRevert(address level4Addr) payable returns()
func (_ChainCallLevel3 *ChainCallLevel3Transactor) CallRevert(opts *bind.TransactOpts, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel3.contract.Transact(opts, "callRevert", level4Addr)
}

// CallRevert is a paid mutator transaction binding the contract method 0xdf0a5d56.
//
// Solidity: function callRevert(address level4Addr) payable returns()
func (_ChainCallLevel3 *ChainCallLevel3Session) CallRevert(level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel3.Contract.CallRevert(&_ChainCallLevel3.TransactOpts, level4Addr)
}

// CallRevert is a paid mutator transaction binding the contract method 0xdf0a5d56.
//
// Solidity: function callRevert(address level4Addr) payable returns()
func (_ChainCallLevel3 *ChainCallLevel3TransactorSession) CallRevert(level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel3.Contract.CallRevert(&_ChainCallLevel3.TransactOpts, level4Addr)
}

// DelegateCallRevert is a paid mutator transaction binding the contract method 0x937de754.
//
// Solidity: function delegateCallRevert(address level4Addr) payable returns()
func (_ChainCallLevel3 *ChainCallLevel3Transactor) DelegateCallRevert(opts *bind.TransactOpts, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel3.contract.Transact(opts, "delegateCallRevert", level4Addr)
}

// DelegateCallRevert is a paid mutator transaction binding the contract method 0x937de754.
//
// Solidity: function delegateCallRevert(address level4Addr) payable returns()
func (_ChainCallLevel3 *ChainCallLevel3Session) DelegateCallRevert(level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel3.Contract.DelegateCallRevert(&_ChainCallLevel3.TransactOpts, level4Addr)
}

// DelegateCallRevert is a paid mutator transaction binding the contract method 0x937de754.
//
// Solidity: function delegateCallRevert(address level4Addr) payable returns()
func (_ChainCallLevel3 *ChainCallLevel3TransactorSession) DelegateCallRevert(level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel3.Contract.DelegateCallRevert(&_ChainCallLevel3.TransactOpts, level4Addr)
}

// Exec is a paid mutator transaction binding the contract method 0x6bb6126e.
//
// Solidity: function exec(address level4Addr) payable returns()
func (_ChainCallLevel3 *ChainCallLevel3Transactor) Exec(opts *bind.TransactOpts, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel3.contract.Transact(opts, "exec", level4Addr)
}

// Exec is a paid mutator transaction binding the contract method 0x6bb6126e.
//
// Solidity: function exec(address level4Addr) payable returns()
func (_ChainCallLevel3 *ChainCallLevel3Session) Exec(level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel3.Contract.Exec(&_ChainCallLevel3.TransactOpts, level4Addr)
}

// Exec is a paid mutator transaction binding the contract method 0x6bb6126e.
//
// Solidity: function exec(address level4Addr) payable returns()
func (_ChainCallLevel3 *ChainCallLevel3TransactorSession) Exec(level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel3.Contract.Exec(&_ChainCallLevel3.TransactOpts, level4Addr)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ChainCallLevel3 *ChainCallLevel3Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainCallLevel3.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ChainCallLevel3 *ChainCallLevel3Session) Receive() (*types.Transaction, error) {
	return _ChainCallLevel3.Contract.Receive(&_ChainCallLevel3.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ChainCallLevel3 *ChainCallLevel3TransactorSession) Receive() (*types.Transaction, error) {
	return _ChainCallLevel3.Contract.Receive(&_ChainCallLevel3.TransactOpts)
}
