// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package destruct

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
)

// DestructMetaData contains all meta data concerning the Destruct contract.
var DestructMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"close\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DestructABI is the input ABI used to generate the binding from.
// Deprecated: Use DestructMetaData.ABI instead.
var DestructABI = DestructMetaData.ABI

// Destruct is an auto generated Go binding around an Ethereum contract.
type Destruct struct {
	DestructCaller     // Read-only binding to the contract
	DestructTransactor // Write-only binding to the contract
	DestructFilterer   // Log filterer for contract events
}

// DestructCaller is an auto generated read-only Go binding around an Ethereum contract.
type DestructCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DestructTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DestructTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DestructFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DestructFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DestructSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DestructSession struct {
	Contract     *Destruct         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DestructCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DestructCallerSession struct {
	Contract *DestructCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// DestructTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DestructTransactorSession struct {
	Contract     *DestructTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DestructRaw is an auto generated low-level Go binding around an Ethereum contract.
type DestructRaw struct {
	Contract *Destruct // Generic contract binding to access the raw methods on
}

// DestructCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DestructCallerRaw struct {
	Contract *DestructCaller // Generic read-only contract binding to access the raw methods on
}

// DestructTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DestructTransactorRaw struct {
	Contract *DestructTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDestruct creates a new instance of Destruct, bound to a specific deployed contract.
func NewDestruct(address common.Address, backend bind.ContractBackend) (*Destruct, error) {
	contract, err := bindDestruct(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Destruct{DestructCaller: DestructCaller{contract: contract}, DestructTransactor: DestructTransactor{contract: contract}, DestructFilterer: DestructFilterer{contract: contract}}, nil
}

// NewDestructCaller creates a new read-only instance of Destruct, bound to a specific deployed contract.
func NewDestructCaller(address common.Address, caller bind.ContractCaller) (*DestructCaller, error) {
	contract, err := bindDestruct(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DestructCaller{contract: contract}, nil
}

// NewDestructTransactor creates a new write-only instance of Destruct, bound to a specific deployed contract.
func NewDestructTransactor(address common.Address, transactor bind.ContractTransactor) (*DestructTransactor, error) {
	contract, err := bindDestruct(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DestructTransactor{contract: contract}, nil
}

// NewDestructFilterer creates a new log filterer instance of Destruct, bound to a specific deployed contract.
func NewDestructFilterer(address common.Address, filterer bind.ContractFilterer) (*DestructFilterer, error) {
	contract, err := bindDestruct(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DestructFilterer{contract: contract}, nil
}

// bindDestruct binds a generic wrapper to an already deployed contract.
func bindDestruct(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DestructABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Destruct *DestructRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Destruct.Contract.DestructCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Destruct *DestructRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Destruct.Contract.DestructTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Destruct *DestructRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Destruct.Contract.DestructTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Destruct *DestructCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Destruct.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Destruct *DestructTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Destruct.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Destruct *DestructTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Destruct.Contract.contract.Transact(opts, method, params...)
}

// Retrieve is a free data retrieval call binding the contract method 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (_Destruct *DestructCaller) Retrieve(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Destruct.contract.Call(opts, &out, "retrieve")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Retrieve is a free data retrieval call binding the contract method 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (_Destruct *DestructSession) Retrieve() (*big.Int, error) {
	return _Destruct.Contract.Retrieve(&_Destruct.CallOpts)
}

// Retrieve is a free data retrieval call binding the contract method 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (_Destruct *DestructCallerSession) Retrieve() (*big.Int, error) {
	return _Destruct.Contract.Retrieve(&_Destruct.CallOpts)
}

// Close is a paid mutator transaction binding the contract method 0x43d726d6.
//
// Solidity: function close() returns()
func (_Destruct *DestructTransactor) Close(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Destruct.contract.Transact(opts, "close")
}

// Close is a paid mutator transaction binding the contract method 0x43d726d6.
//
// Solidity: function close() returns()
func (_Destruct *DestructSession) Close() (*types.Transaction, error) {
	return _Destruct.Contract.Close(&_Destruct.TransactOpts)
}

// Close is a paid mutator transaction binding the contract method 0x43d726d6.
//
// Solidity: function close() returns()
func (_Destruct *DestructTransactorSession) Close() (*types.Transaction, error) {
	return _Destruct.Contract.Close(&_Destruct.TransactOpts)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 num) returns()
func (_Destruct *DestructTransactor) Store(opts *bind.TransactOpts, num *big.Int) (*types.Transaction, error) {
	return _Destruct.contract.Transact(opts, "store", num)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 num) returns()
func (_Destruct *DestructSession) Store(num *big.Int) (*types.Transaction, error) {
	return _Destruct.Contract.Store(&_Destruct.TransactOpts, num)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 num) returns()
func (_Destruct *DestructTransactorSession) Store(num *big.Int) (*types.Transaction, error) {
	return _Destruct.Contract.Store(&_Destruct.TransactOpts, num)
}
