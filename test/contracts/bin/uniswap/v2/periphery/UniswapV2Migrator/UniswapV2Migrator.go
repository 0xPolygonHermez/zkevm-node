// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package UniswapV2Migrator

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

// UniswapV2MigratorMetaData contains all meta data concerning the UniswapV2Migrator contract.
var UniswapV2MigratorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factoryV1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// UniswapV2MigratorABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapV2MigratorMetaData.ABI instead.
var UniswapV2MigratorABI = UniswapV2MigratorMetaData.ABI

// UniswapV2Migrator is an auto generated Go binding around an Ethereum contract.
type UniswapV2Migrator struct {
	UniswapV2MigratorCaller     // Read-only binding to the contract
	UniswapV2MigratorTransactor // Write-only binding to the contract
	UniswapV2MigratorFilterer   // Log filterer for contract events
}

// UniswapV2MigratorCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniswapV2MigratorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2MigratorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniswapV2MigratorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2MigratorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniswapV2MigratorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2MigratorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniswapV2MigratorSession struct {
	Contract     *UniswapV2Migrator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// UniswapV2MigratorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniswapV2MigratorCallerSession struct {
	Contract *UniswapV2MigratorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// UniswapV2MigratorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniswapV2MigratorTransactorSession struct {
	Contract     *UniswapV2MigratorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// UniswapV2MigratorRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniswapV2MigratorRaw struct {
	Contract *UniswapV2Migrator // Generic contract binding to access the raw methods on
}

// UniswapV2MigratorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniswapV2MigratorCallerRaw struct {
	Contract *UniswapV2MigratorCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapV2MigratorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniswapV2MigratorTransactorRaw struct {
	Contract *UniswapV2MigratorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapV2Migrator creates a new instance of UniswapV2Migrator, bound to a specific deployed contract.
func NewUniswapV2Migrator(address common.Address, backend bind.ContractBackend) (*UniswapV2Migrator, error) {
	contract, err := bindUniswapV2Migrator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapV2Migrator{UniswapV2MigratorCaller: UniswapV2MigratorCaller{contract: contract}, UniswapV2MigratorTransactor: UniswapV2MigratorTransactor{contract: contract}, UniswapV2MigratorFilterer: UniswapV2MigratorFilterer{contract: contract}}, nil
}

// NewUniswapV2MigratorCaller creates a new read-only instance of UniswapV2Migrator, bound to a specific deployed contract.
func NewUniswapV2MigratorCaller(address common.Address, caller bind.ContractCaller) (*UniswapV2MigratorCaller, error) {
	contract, err := bindUniswapV2Migrator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2MigratorCaller{contract: contract}, nil
}

// NewUniswapV2MigratorTransactor creates a new write-only instance of UniswapV2Migrator, bound to a specific deployed contract.
func NewUniswapV2MigratorTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapV2MigratorTransactor, error) {
	contract, err := bindUniswapV2Migrator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2MigratorTransactor{contract: contract}, nil
}

// NewUniswapV2MigratorFilterer creates a new log filterer instance of UniswapV2Migrator, bound to a specific deployed contract.
func NewUniswapV2MigratorFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapV2MigratorFilterer, error) {
	contract, err := bindUniswapV2Migrator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapV2MigratorFilterer{contract: contract}, nil
}

// bindUniswapV2Migrator binds a generic wrapper to an already deployed contract.
func bindUniswapV2Migrator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UniswapV2MigratorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2Migrator *UniswapV2MigratorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2Migrator.Contract.UniswapV2MigratorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2Migrator *UniswapV2MigratorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.UniswapV2MigratorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2Migrator *UniswapV2MigratorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.UniswapV2MigratorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2Migrator *UniswapV2MigratorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2Migrator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2Migrator *UniswapV2MigratorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2Migrator *UniswapV2MigratorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.contract.Transact(opts, method, params...)
}

// Migrate is a paid mutator transaction binding the contract method 0xb7df1d25.
//
// Solidity: function migrate(address token, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns()
func (_UniswapV2Migrator *UniswapV2MigratorTransactor) Migrate(opts *bind.TransactOpts, token common.Address, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV2Migrator.contract.Transact(opts, "migrate", token, amountTokenMin, amountETHMin, to, deadline)
}

// Migrate is a paid mutator transaction binding the contract method 0xb7df1d25.
//
// Solidity: function migrate(address token, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns()
func (_UniswapV2Migrator *UniswapV2MigratorSession) Migrate(token common.Address, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.Migrate(&_UniswapV2Migrator.TransactOpts, token, amountTokenMin, amountETHMin, to, deadline)
}

// Migrate is a paid mutator transaction binding the contract method 0xb7df1d25.
//
// Solidity: function migrate(address token, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns()
func (_UniswapV2Migrator *UniswapV2MigratorTransactorSession) Migrate(token common.Address, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.Migrate(&_UniswapV2Migrator.TransactOpts, token, amountTokenMin, amountETHMin, to, deadline)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_UniswapV2Migrator *UniswapV2MigratorTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Migrator.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_UniswapV2Migrator *UniswapV2MigratorSession) Receive() (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.Receive(&_UniswapV2Migrator.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_UniswapV2Migrator *UniswapV2MigratorTransactorSession) Receive() (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.Receive(&_UniswapV2Migrator.TransactOpts)
}
