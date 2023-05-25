// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Double

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

// DoubleMetaData contains all meta data concerning the Double contract.
var DoubleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"a\",\"type\":\"int256\"}],\"name\":\"double\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610152806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80636ffa1caa14610030575b600080fd5b61004361003e366004610068565b610055565b60405190815260200160405180910390f35b6000610062826002610097565b92915050565b60006020828403121561007a57600080fd5b5035919050565b634e487b7160e01b600052601160045260246000fd5b60006001600160ff1b03818413828413808216868404861116156100bd576100bd610081565b600160ff1b60008712828116878305891216156100dc576100dc610081565b600087129250878205871284841616156100f8576100f8610081565b8785058712818416161561010e5761010e610081565b50505092909302939250505056fea26469706673582212205c4d503eca301a04694bc37b49abe6e9f39bf64a198399f967ba70f9ca5f097764736f6c634300080c0033",
}

// DoubleABI is the input ABI used to generate the binding from.
// Deprecated: Use DoubleMetaData.ABI instead.
var DoubleABI = DoubleMetaData.ABI

// DoubleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DoubleMetaData.Bin instead.
var DoubleBin = DoubleMetaData.Bin

// DeployDouble deploys a new Ethereum contract, binding an instance of Double to it.
func DeployDouble(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Double, error) {
	parsed, err := DoubleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DoubleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Double{DoubleCaller: DoubleCaller{contract: contract}, DoubleTransactor: DoubleTransactor{contract: contract}, DoubleFilterer: DoubleFilterer{contract: contract}}, nil
}

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
	parsed, err := DoubleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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
