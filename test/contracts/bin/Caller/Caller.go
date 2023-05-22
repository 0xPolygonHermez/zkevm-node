// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Caller

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

// CallerMetaData contains all meta data concerning the Caller contract.
var CallerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_num\",\"type\":\"uint256\"}],\"name\":\"call\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_num\",\"type\":\"uint256\"}],\"name\":\"delegateCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_num\",\"type\":\"uint256\"}],\"name\":\"multiCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"name\":\"staticCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610493806100206000396000f3fe60806040526004361061003f5760003560e01c80630f6be00d146100445780633bd9ef281461005957806387b1d6ad1461006c578063c6c211e91461007f575b600080fd5b610057610052366004610399565b610092565b005b610057610067366004610399565b610183565b61005761007a3660046103c5565b610265565b61005761008d366004610399565b610360565b6000826001600160a01b0316826040516024016100b191815260200190565b60408051601f198184030181529181526020820180516001600160e01b0316636466414b60e01b179052516100e691906103e9565b600060405180830381855af49150503d8060008114610121576040519150601f19603f3d011682016040523d82523d6000602084013e610126565b606091505b5050809150508061017e5760405162461bcd60e51b815260206004820152601f60248201527f6661696c656420746f20706572666f726d2064656c65676174652063616c6c0060448201526064015b60405180910390fd5b505050565b6000826001600160a01b0316826040516024016101a291815260200190565b60408051601f198184030181529181526020820180516001600160e01b0316636466414b60e01b179052516101d791906103e9565b6000604051808303816000865af19150503d8060008114610214576040519150601f19603f3d011682016040523d82523d6000602084013e610219565b606091505b5050809150508061017e5760405162461bcd60e51b815260206004820152601660248201527519985a5b1959081d1bc81c195c999bdc9b4818d85b1b60521b6044820152606401610175565b60408051600481526024810182526020810180516001600160e01b031663813d8a3760e01b17905290516000916060916001600160a01b038516916102a9916103e9565b600060405180830381855afa9150503d80600081146102e4576040519150601f19603f3d011682016040523d82523d6000602084013e6102e9565b606091505b5090925090508161033c5760405162461bcd60e51b815260206004820152601d60248201527f6661696c656420746f20706572666f726d207374617469632063616c6c0000006044820152606401610175565b6000806000838060200190518101906103559190610424565b505050505050505050565b61036a8282610183565b6103748282610092565b61037d82610265565b5050565b6001600160a01b038116811461039657600080fd5b50565b600080604083850312156103ac57600080fd5b82356103b781610381565b946020939093013593505050565b6000602082840312156103d757600080fd5b81356103e281610381565b9392505050565b6000825160005b8181101561040a57602081860181015185830152016103f0565b81811115610419576000828501525b509190910192915050565b60008060006060848603121561043957600080fd5b83519250602084015161044b81610381565b8092505060408401519050925092509256fea26469706673582212207c864726dd059eb16f2efd3e6e7c97283b977c6501d217b22a8dddc40488364c64736f6c634300080c0033",
}

// CallerABI is the input ABI used to generate the binding from.
// Deprecated: Use CallerMetaData.ABI instead.
var CallerABI = CallerMetaData.ABI

// CallerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CallerMetaData.Bin instead.
var CallerBin = CallerMetaData.Bin

// DeployCaller deploys a new Ethereum contract, binding an instance of Caller to it.
func DeployCaller(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Caller, error) {
	parsed, err := CallerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CallerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Caller{CallerCaller: CallerCaller{contract: contract}, CallerTransactor: CallerTransactor{contract: contract}, CallerFilterer: CallerFilterer{contract: contract}}, nil
}

// Caller is an auto generated Go binding around an Ethereum contract.
type Caller struct {
	CallerCaller     // Read-only binding to the contract
	CallerTransactor // Write-only binding to the contract
	CallerFilterer   // Log filterer for contract events
}

// CallerCaller is an auto generated read-only Go binding around an Ethereum contract.
type CallerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CallerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CallerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CallerSession struct {
	Contract     *Caller           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CallerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CallerCallerSession struct {
	Contract *CallerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// CallerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CallerTransactorSession struct {
	Contract     *CallerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CallerRaw is an auto generated low-level Go binding around an Ethereum contract.
type CallerRaw struct {
	Contract *Caller // Generic contract binding to access the raw methods on
}

// CallerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CallerCallerRaw struct {
	Contract *CallerCaller // Generic read-only contract binding to access the raw methods on
}

// CallerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CallerTransactorRaw struct {
	Contract *CallerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCaller creates a new instance of Caller, bound to a specific deployed contract.
func NewCaller(address common.Address, backend bind.ContractBackend) (*Caller, error) {
	contract, err := bindCaller(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Caller{CallerCaller: CallerCaller{contract: contract}, CallerTransactor: CallerTransactor{contract: contract}, CallerFilterer: CallerFilterer{contract: contract}}, nil
}

// NewCallerCaller creates a new read-only instance of Caller, bound to a specific deployed contract.
func NewCallerCaller(address common.Address, caller bind.ContractCaller) (*CallerCaller, error) {
	contract, err := bindCaller(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CallerCaller{contract: contract}, nil
}

// NewCallerTransactor creates a new write-only instance of Caller, bound to a specific deployed contract.
func NewCallerTransactor(address common.Address, transactor bind.ContractTransactor) (*CallerTransactor, error) {
	contract, err := bindCaller(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CallerTransactor{contract: contract}, nil
}

// NewCallerFilterer creates a new log filterer instance of Caller, bound to a specific deployed contract.
func NewCallerFilterer(address common.Address, filterer bind.ContractFilterer) (*CallerFilterer, error) {
	contract, err := bindCaller(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CallerFilterer{contract: contract}, nil
}

// bindCaller binds a generic wrapper to an already deployed contract.
func bindCaller(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CallerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Caller *CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Caller.Contract.CallerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Caller *CallerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Caller.Contract.CallerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Caller *CallerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Caller.Contract.CallerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Caller *CallerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Caller.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Caller *CallerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Caller.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Caller *CallerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Caller.Contract.contract.Transact(opts, method, params...)
}

// Call is a paid mutator transaction binding the contract method 0x3bd9ef28.
//
// Solidity: function call(address _contract, uint256 _num) payable returns()
func (_Caller *CallerTransactor) Call(opts *bind.TransactOpts, _contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.contract.Transact(opts, "call", _contract, _num)
}

// Call is a paid mutator transaction binding the contract method 0x3bd9ef28.
//
// Solidity: function call(address _contract, uint256 _num) payable returns()
func (_Caller *CallerSession) Call(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.Contract.Call(&_Caller.TransactOpts, _contract, _num)
}

// Call is a paid mutator transaction binding the contract method 0x3bd9ef28.
//
// Solidity: function call(address _contract, uint256 _num) payable returns()
func (_Caller *CallerTransactorSession) Call(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.Contract.Call(&_Caller.TransactOpts, _contract, _num)
}

// DelegateCall is a paid mutator transaction binding the contract method 0x0f6be00d.
//
// Solidity: function delegateCall(address _contract, uint256 _num) payable returns()
func (_Caller *CallerTransactor) DelegateCall(opts *bind.TransactOpts, _contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.contract.Transact(opts, "delegateCall", _contract, _num)
}

// DelegateCall is a paid mutator transaction binding the contract method 0x0f6be00d.
//
// Solidity: function delegateCall(address _contract, uint256 _num) payable returns()
func (_Caller *CallerSession) DelegateCall(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.Contract.DelegateCall(&_Caller.TransactOpts, _contract, _num)
}

// DelegateCall is a paid mutator transaction binding the contract method 0x0f6be00d.
//
// Solidity: function delegateCall(address _contract, uint256 _num) payable returns()
func (_Caller *CallerTransactorSession) DelegateCall(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.Contract.DelegateCall(&_Caller.TransactOpts, _contract, _num)
}

// MultiCall is a paid mutator transaction binding the contract method 0xc6c211e9.
//
// Solidity: function multiCall(address _contract, uint256 _num) payable returns()
func (_Caller *CallerTransactor) MultiCall(opts *bind.TransactOpts, _contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.contract.Transact(opts, "multiCall", _contract, _num)
}

// MultiCall is a paid mutator transaction binding the contract method 0xc6c211e9.
//
// Solidity: function multiCall(address _contract, uint256 _num) payable returns()
func (_Caller *CallerSession) MultiCall(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.Contract.MultiCall(&_Caller.TransactOpts, _contract, _num)
}

// MultiCall is a paid mutator transaction binding the contract method 0xc6c211e9.
//
// Solidity: function multiCall(address _contract, uint256 _num) payable returns()
func (_Caller *CallerTransactorSession) MultiCall(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.Contract.MultiCall(&_Caller.TransactOpts, _contract, _num)
}

// StaticCall is a paid mutator transaction binding the contract method 0x87b1d6ad.
//
// Solidity: function staticCall(address _contract) payable returns()
func (_Caller *CallerTransactor) StaticCall(opts *bind.TransactOpts, _contract common.Address) (*types.Transaction, error) {
	return _Caller.contract.Transact(opts, "staticCall", _contract)
}

// StaticCall is a paid mutator transaction binding the contract method 0x87b1d6ad.
//
// Solidity: function staticCall(address _contract) payable returns()
func (_Caller *CallerSession) StaticCall(_contract common.Address) (*types.Transaction, error) {
	return _Caller.Contract.StaticCall(&_Caller.TransactOpts, _contract)
}

// StaticCall is a paid mutator transaction binding the contract method 0x87b1d6ad.
//
// Solidity: function staticCall(address _contract) payable returns()
func (_Caller *CallerTransactorSession) StaticCall(_contract common.Address) (*types.Transaction, error) {
	return _Caller.Contract.StaticCall(&_Caller.TransactOpts, _contract)
}
