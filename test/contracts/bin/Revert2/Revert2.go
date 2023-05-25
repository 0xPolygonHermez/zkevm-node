// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Revert2

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

// Revert2MetaData contains all meta data concerning the Revert2 contract.
var Revert2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"generateError\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b5060ae8061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80634abbb40a14602d575b600080fd5b60336035565b005b60405162461bcd60e51b8152602060048201526014602482015273546f646179206973206e6f74206a7565726e657360601b604482015260640160405180910390fdfea2646970667358221220b10f6c7dcbccc8fe7166fa5bad433bf84bfb174e480a2594f33ebeb4d996073e64736f6c634300080c0033",
}

// Revert2ABI is the input ABI used to generate the binding from.
// Deprecated: Use Revert2MetaData.ABI instead.
var Revert2ABI = Revert2MetaData.ABI

// Revert2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Revert2MetaData.Bin instead.
var Revert2Bin = Revert2MetaData.Bin

// DeployRevert2 deploys a new Ethereum contract, binding an instance of Revert2 to it.
func DeployRevert2(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Revert2, error) {
	parsed, err := Revert2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Revert2Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Revert2{Revert2Caller: Revert2Caller{contract: contract}, Revert2Transactor: Revert2Transactor{contract: contract}, Revert2Filterer: Revert2Filterer{contract: contract}}, nil
}

// Revert2 is an auto generated Go binding around an Ethereum contract.
type Revert2 struct {
	Revert2Caller     // Read-only binding to the contract
	Revert2Transactor // Write-only binding to the contract
	Revert2Filterer   // Log filterer for contract events
}

// Revert2Caller is an auto generated read-only Go binding around an Ethereum contract.
type Revert2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Revert2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Revert2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Revert2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Revert2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Revert2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Revert2Session struct {
	Contract     *Revert2          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Revert2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Revert2CallerSession struct {
	Contract *Revert2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// Revert2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Revert2TransactorSession struct {
	Contract     *Revert2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// Revert2Raw is an auto generated low-level Go binding around an Ethereum contract.
type Revert2Raw struct {
	Contract *Revert2 // Generic contract binding to access the raw methods on
}

// Revert2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Revert2CallerRaw struct {
	Contract *Revert2Caller // Generic read-only contract binding to access the raw methods on
}

// Revert2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Revert2TransactorRaw struct {
	Contract *Revert2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewRevert2 creates a new instance of Revert2, bound to a specific deployed contract.
func NewRevert2(address common.Address, backend bind.ContractBackend) (*Revert2, error) {
	contract, err := bindRevert2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Revert2{Revert2Caller: Revert2Caller{contract: contract}, Revert2Transactor: Revert2Transactor{contract: contract}, Revert2Filterer: Revert2Filterer{contract: contract}}, nil
}

// NewRevert2Caller creates a new read-only instance of Revert2, bound to a specific deployed contract.
func NewRevert2Caller(address common.Address, caller bind.ContractCaller) (*Revert2Caller, error) {
	contract, err := bindRevert2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Revert2Caller{contract: contract}, nil
}

// NewRevert2Transactor creates a new write-only instance of Revert2, bound to a specific deployed contract.
func NewRevert2Transactor(address common.Address, transactor bind.ContractTransactor) (*Revert2Transactor, error) {
	contract, err := bindRevert2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Revert2Transactor{contract: contract}, nil
}

// NewRevert2Filterer creates a new log filterer instance of Revert2, bound to a specific deployed contract.
func NewRevert2Filterer(address common.Address, filterer bind.ContractFilterer) (*Revert2Filterer, error) {
	contract, err := bindRevert2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Revert2Filterer{contract: contract}, nil
}

// bindRevert2 binds a generic wrapper to an already deployed contract.
func bindRevert2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Revert2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Revert2 *Revert2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Revert2.Contract.Revert2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Revert2 *Revert2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Revert2.Contract.Revert2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Revert2 *Revert2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Revert2.Contract.Revert2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Revert2 *Revert2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Revert2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Revert2 *Revert2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Revert2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Revert2 *Revert2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Revert2.Contract.contract.Transact(opts, method, params...)
}

// GenerateError is a paid mutator transaction binding the contract method 0x4abbb40a.
//
// Solidity: function generateError() returns()
func (_Revert2 *Revert2Transactor) GenerateError(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Revert2.contract.Transact(opts, "generateError")
}

// GenerateError is a paid mutator transaction binding the contract method 0x4abbb40a.
//
// Solidity: function generateError() returns()
func (_Revert2 *Revert2Session) GenerateError() (*types.Transaction, error) {
	return _Revert2.Contract.GenerateError(&_Revert2.TransactOpts)
}

// GenerateError is a paid mutator transaction binding the contract method 0x4abbb40a.
//
// Solidity: function generateError() returns()
func (_Revert2 *Revert2TransactorSession) GenerateError() (*types.Transaction, error) {
	return _Revert2.Contract.GenerateError(&_Revert2.TransactOpts)
}
