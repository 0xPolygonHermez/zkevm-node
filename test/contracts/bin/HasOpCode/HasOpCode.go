// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package HasOpCode

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

// HasOpCodeMetaData contains all meta data concerning the HasOpCode contract.
var HasOpCodeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"opBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"opGasPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600080556000600155348015601857600080fd5b506080806100276000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80633ab08cf914603757806374c73639146042575b600080fd5b60403331600155565b005b60403a60005556fea264697066735822122086d3f33465f92e2f6ddc32c9acfb8512d8c86ff16e540197cd39d4f3aaf38ffc64736f6c634300080c0033",
}

// HasOpCodeABI is the input ABI used to generate the binding from.
// Deprecated: Use HasOpCodeMetaData.ABI instead.
var HasOpCodeABI = HasOpCodeMetaData.ABI

// HasOpCodeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use HasOpCodeMetaData.Bin instead.
var HasOpCodeBin = HasOpCodeMetaData.Bin

// DeployHasOpCode deploys a new Ethereum contract, binding an instance of HasOpCode to it.
func DeployHasOpCode(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *HasOpCode, error) {
	parsed, err := HasOpCodeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(HasOpCodeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &HasOpCode{HasOpCodeCaller: HasOpCodeCaller{contract: contract}, HasOpCodeTransactor: HasOpCodeTransactor{contract: contract}, HasOpCodeFilterer: HasOpCodeFilterer{contract: contract}}, nil
}

// HasOpCode is an auto generated Go binding around an Ethereum contract.
type HasOpCode struct {
	HasOpCodeCaller     // Read-only binding to the contract
	HasOpCodeTransactor // Write-only binding to the contract
	HasOpCodeFilterer   // Log filterer for contract events
}

// HasOpCodeCaller is an auto generated read-only Go binding around an Ethereum contract.
type HasOpCodeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HasOpCodeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HasOpCodeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HasOpCodeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HasOpCodeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HasOpCodeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HasOpCodeSession struct {
	Contract     *HasOpCode        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HasOpCodeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HasOpCodeCallerSession struct {
	Contract *HasOpCodeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// HasOpCodeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HasOpCodeTransactorSession struct {
	Contract     *HasOpCodeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// HasOpCodeRaw is an auto generated low-level Go binding around an Ethereum contract.
type HasOpCodeRaw struct {
	Contract *HasOpCode // Generic contract binding to access the raw methods on
}

// HasOpCodeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HasOpCodeCallerRaw struct {
	Contract *HasOpCodeCaller // Generic read-only contract binding to access the raw methods on
}

// HasOpCodeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HasOpCodeTransactorRaw struct {
	Contract *HasOpCodeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHasOpCode creates a new instance of HasOpCode, bound to a specific deployed contract.
func NewHasOpCode(address common.Address, backend bind.ContractBackend) (*HasOpCode, error) {
	contract, err := bindHasOpCode(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HasOpCode{HasOpCodeCaller: HasOpCodeCaller{contract: contract}, HasOpCodeTransactor: HasOpCodeTransactor{contract: contract}, HasOpCodeFilterer: HasOpCodeFilterer{contract: contract}}, nil
}

// NewHasOpCodeCaller creates a new read-only instance of HasOpCode, bound to a specific deployed contract.
func NewHasOpCodeCaller(address common.Address, caller bind.ContractCaller) (*HasOpCodeCaller, error) {
	contract, err := bindHasOpCode(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HasOpCodeCaller{contract: contract}, nil
}

// NewHasOpCodeTransactor creates a new write-only instance of HasOpCode, bound to a specific deployed contract.
func NewHasOpCodeTransactor(address common.Address, transactor bind.ContractTransactor) (*HasOpCodeTransactor, error) {
	contract, err := bindHasOpCode(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HasOpCodeTransactor{contract: contract}, nil
}

// NewHasOpCodeFilterer creates a new log filterer instance of HasOpCode, bound to a specific deployed contract.
func NewHasOpCodeFilterer(address common.Address, filterer bind.ContractFilterer) (*HasOpCodeFilterer, error) {
	contract, err := bindHasOpCode(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HasOpCodeFilterer{contract: contract}, nil
}

// bindHasOpCode binds a generic wrapper to an already deployed contract.
func bindHasOpCode(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := HasOpCodeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HasOpCode *HasOpCodeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HasOpCode.Contract.HasOpCodeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HasOpCode *HasOpCodeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HasOpCode.Contract.HasOpCodeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HasOpCode *HasOpCodeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HasOpCode.Contract.HasOpCodeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HasOpCode *HasOpCodeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HasOpCode.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HasOpCode *HasOpCodeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HasOpCode.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HasOpCode *HasOpCodeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HasOpCode.Contract.contract.Transact(opts, method, params...)
}

// OpBalance is a paid mutator transaction binding the contract method 0x3ab08cf9.
//
// Solidity: function opBalance() returns()
func (_HasOpCode *HasOpCodeTransactor) OpBalance(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HasOpCode.contract.Transact(opts, "opBalance")
}

// OpBalance is a paid mutator transaction binding the contract method 0x3ab08cf9.
//
// Solidity: function opBalance() returns()
func (_HasOpCode *HasOpCodeSession) OpBalance() (*types.Transaction, error) {
	return _HasOpCode.Contract.OpBalance(&_HasOpCode.TransactOpts)
}

// OpBalance is a paid mutator transaction binding the contract method 0x3ab08cf9.
//
// Solidity: function opBalance() returns()
func (_HasOpCode *HasOpCodeTransactorSession) OpBalance() (*types.Transaction, error) {
	return _HasOpCode.Contract.OpBalance(&_HasOpCode.TransactOpts)
}

// OpGasPrice is a paid mutator transaction binding the contract method 0x74c73639.
//
// Solidity: function opGasPrice() returns()
func (_HasOpCode *HasOpCodeTransactor) OpGasPrice(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HasOpCode.contract.Transact(opts, "opGasPrice")
}

// OpGasPrice is a paid mutator transaction binding the contract method 0x74c73639.
//
// Solidity: function opGasPrice() returns()
func (_HasOpCode *HasOpCodeSession) OpGasPrice() (*types.Transaction, error) {
	return _HasOpCode.Contract.OpGasPrice(&_HasOpCode.TransactOpts)
}

// OpGasPrice is a paid mutator transaction binding the contract method 0x74c73639.
//
// Solidity: function opGasPrice() returns()
func (_HasOpCode *HasOpCodeTransactorSession) OpGasPrice() (*types.Transaction, error) {
	return _HasOpCode.Contract.OpGasPrice(&_HasOpCode.TransactOpts)
}
