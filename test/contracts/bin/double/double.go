// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package double

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

// DoubleMetaData contains all meta data concerning the Double contract.
var DoubleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"a\",\"type\":\"int256\"}],\"name\":\"double\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// DoubleABI is the input ABI used to generate the binding from.
// Deprecated: Use DoubleMetaData.ABI instead.
var DoubleABI = DoubleMetaData.ABI

// Double is an auto generated Go binding around an Ethereum contract.
type Double struct {
	DoubleCaller     // Read-only binding to the contract
	DoubleTransactor // Write-only binding to the contract
	DoubleFilterer   // Log filterer for contract events
}

// DoubleCaller is an auto generated read-only Go binding around an Ethereum contract.
type DoubleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DoubleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DoubleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DoubleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DoubleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DoubleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DoubleSession struct {
	Contract     *Double           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DoubleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DoubleCallerSession struct {
	Contract *DoubleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DoubleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DoubleTransactorSession struct {
	Contract     *DoubleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DoubleRaw is an auto generated low-level Go binding around an Ethereum contract.
type DoubleRaw struct {
	Contract *Double // Generic contract binding to access the raw methods on
}

// DoubleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DoubleCallerRaw struct {
	Contract *DoubleCaller // Generic read-only contract binding to access the raw methods on
}

// DoubleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DoubleTransactorRaw struct {
	Contract *DoubleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDouble creates a new instance of Double, bound to a specific deployed contract.
func NewDouble(address common.Address, backend bind.ContractBackend) (*Double, error) {
	contract, err := bindDouble(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Double{DoubleCaller: DoubleCaller{contract: contract}, DoubleTransactor: DoubleTransactor{contract: contract}, DoubleFilterer: DoubleFilterer{contract: contract}}, nil
}

// NewDoubleCaller creates a new read-only instance of Double, bound to a specific deployed contract.
func NewDoubleCaller(address common.Address, caller bind.ContractCaller) (*DoubleCaller, error) {
	contract, err := bindDouble(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DoubleCaller{contract: contract}, nil
}

// NewDoubleTransactor creates a new write-only instance of Double, bound to a specific deployed contract.
func NewDoubleTransactor(address common.Address, transactor bind.ContractTransactor) (*DoubleTransactor, error) {
	contract, err := bindDouble(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DoubleTransactor{contract: contract}, nil
}

// NewDoubleFilterer creates a new log filterer instance of Double, bound to a specific deployed contract.
func NewDoubleFilterer(address common.Address, filterer bind.ContractFilterer) (*DoubleFilterer, error) {
	contract, err := bindDouble(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DoubleFilterer{contract: contract}, nil
}

// bindDouble binds a generic wrapper to an already deployed contract.
func bindDouble(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DoubleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Double *DoubleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Double.Contract.DoubleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Double *DoubleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Double.Contract.DoubleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Double *DoubleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Double.Contract.DoubleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Double *DoubleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Double.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Double *DoubleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Double.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Double *DoubleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Double.Contract.contract.Transact(opts, method, params...)
}

// Double is a free data retrieval call binding the contract method 0x6ffa1caa.
//
// Solidity: function double(int256 a) pure returns(int256)
func (_Double *DoubleCaller) Double(opts *bind.CallOpts, a *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Double.contract.Call(opts, &out, "double", a)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Double is a free data retrieval call binding the contract method 0x6ffa1caa.
//
// Solidity: function double(int256 a) pure returns(int256)
func (_Double *DoubleSession) Double(a *big.Int) (*big.Int, error) {
	return _Double.Contract.Double(&_Double.CallOpts, a)
}

// Double is a free data retrieval call binding the contract method 0x6ffa1caa.
//
// Solidity: function double(int256 a) pure returns(int256)
func (_Double *DoubleCallerSession) Double(a *big.Int) (*big.Int, error) {
	return _Double.Contract.Double(&_Double.CallOpts, a)
}
