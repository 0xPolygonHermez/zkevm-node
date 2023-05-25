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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"exec\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"get\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"t\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506104f6806100206000396000f3fe6080604052600436106100295760003560e01c80636bb6126e1461002e578063c2bc2efc14610043575b600080fd5b61004161003c36600461034e565b610079565b005b34801561004f57600080fd5b5061006361005e36600461034e565b610247565b60405161007091906103ae565b60405180910390f35b60408051600481526024810182526020810180516001600160e01b03166330703a7160e21b17905290516000916001600160a01b038416916100bb91906103e1565b6000604051808303816000865af19150503d80600081146100f8576040519150601f19603f3d011682016040523d82523d6000602084013e6100fd565b606091505b5050809150508061015f5760405162461bcd60e51b815260206004820152602160248201527f6661696c656420746f20706572666f726d2063616c6c20746f206c6576656c206044820152600d60fa1b60648201526084015b60405180910390fd5b60408051600481526024810182526020810180516001600160e01b03166330703a7160e21b17905290516001600160a01b0384169161019d916103e1565b600060405180830381855af49150503d80600081146101d8576040519150601f19603f3d011682016040523d82523d6000602084013e6101dd565b606091505b505080915050806102435760405162461bcd60e51b815260206004820152602a60248201527f6661696c656420746f20706572666f726d2064656c65676174652063616c6c206044820152691d1bc81b195d995b080d60b21b6064820152608401610156565b5050565b60408051600481526024810182526020810180516001600160e01b0316631b53398f60e21b179052905160609160009183916001600160a01b0386169161028e91906103e1565b600060405180830381855afa9150503d80600081146102c9576040519150601f19603f3d011682016040523d82523d6000602084013e6102ce565b606091505b509092509050816103325760405162461bcd60e51b815260206004820152602860248201527f6661696c656420746f20706572666f726d207374617469632063616c6c20746f604482015267081b195d995b080d60c21b6064820152608401610156565b808060200190518101906103469190610413565b949350505050565b60006020828403121561036057600080fd5b81356001600160a01b038116811461037757600080fd5b9392505050565b60005b83811015610399578181015183820152602001610381565b838111156103a8576000848401525b50505050565b60208152600082518060208401526103cd81604085016020870161037e565b601f01601f19169190910160400192915050565b600082516103f381846020870161037e565b9190910192915050565b634e487b7160e01b600052604160045260246000fd5b60006020828403121561042557600080fd5b815167ffffffffffffffff8082111561043d57600080fd5b818401915084601f83011261045157600080fd5b815181811115610463576104636103fd565b604051601f8201601f19908116603f0116810190838211818310171561048b5761048b6103fd565b816040528281528760208487010111156104a457600080fd5b6104b583602083016020880161037e565b97965050505050505056fea26469706673582212209bedf8d59efd5de6abd239994136138fbe5555047e0cce461bf5853bf28167db64736f6c634300080c0033",
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
