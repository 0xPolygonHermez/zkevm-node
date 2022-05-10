// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Called

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

// CalledMetaData contains all meta data concerning the Called contract.
var CalledMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"entrypoint\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"expectedSender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610102806100206000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80631754dba5146037578063a65d69d4146065575b600080fd5b6000546049906001600160a01b031681565b6040516001600160a01b03909116815260200160405180910390f35b606b606d565b005b6000546001600160a01b0316331460ca5760405162461bcd60e51b815260206004820152601c60248201527f657870656374656453656e64657220213d206d73672e73656e64657200000000604482015260640160405180910390fd5b56fea2646970667358221220e77712b2bce25e17b9647ff59680dec6a71b33d5f366083e52d87fd17d7f7bd164736f6c634300080c0033",
}

// CalledABI is the input ABI used to generate the binding from.
// Deprecated: Use CalledMetaData.ABI instead.
var CalledABI = CalledMetaData.ABI

// CalledBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CalledMetaData.Bin instead.
var CalledBin = CalledMetaData.Bin

// DeployCalled deploys a new Ethereum contract, binding an instance of Called to it.
func DeployCalled(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Called, error) {
	parsed, err := CalledMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CalledBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Called{CalledCaller: CalledCaller{contract: contract}, CalledTransactor: CalledTransactor{contract: contract}, CalledFilterer: CalledFilterer{contract: contract}}, nil
}

// Called is an auto generated Go binding around an Ethereum contract.
type Called struct {
	CalledCaller     // Read-only binding to the contract
	CalledTransactor // Write-only binding to the contract
	CalledFilterer   // Log filterer for contract events
}

// CalledCaller is an auto generated read-only Go binding around an Ethereum contract.
type CalledCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CalledTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CalledTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CalledFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CalledFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CalledSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CalledSession struct {
	Contract     *Called           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CalledCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CalledCallerSession struct {
	Contract *CalledCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// CalledTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CalledTransactorSession struct {
	Contract     *CalledTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CalledRaw is an auto generated low-level Go binding around an Ethereum contract.
type CalledRaw struct {
	Contract *Called // Generic contract binding to access the raw methods on
}

// CalledCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CalledCallerRaw struct {
	Contract *CalledCaller // Generic read-only contract binding to access the raw methods on
}

// CalledTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CalledTransactorRaw struct {
	Contract *CalledTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCalled creates a new instance of Called, bound to a specific deployed contract.
func NewCalled(address common.Address, backend bind.ContractBackend) (*Called, error) {
	contract, err := bindCalled(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Called{CalledCaller: CalledCaller{contract: contract}, CalledTransactor: CalledTransactor{contract: contract}, CalledFilterer: CalledFilterer{contract: contract}}, nil
}

// NewCalledCaller creates a new read-only instance of Called, bound to a specific deployed contract.
func NewCalledCaller(address common.Address, caller bind.ContractCaller) (*CalledCaller, error) {
	contract, err := bindCalled(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CalledCaller{contract: contract}, nil
}

// NewCalledTransactor creates a new write-only instance of Called, bound to a specific deployed contract.
func NewCalledTransactor(address common.Address, transactor bind.ContractTransactor) (*CalledTransactor, error) {
	contract, err := bindCalled(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CalledTransactor{contract: contract}, nil
}

// NewCalledFilterer creates a new log filterer instance of Called, bound to a specific deployed contract.
func NewCalledFilterer(address common.Address, filterer bind.ContractFilterer) (*CalledFilterer, error) {
	contract, err := bindCalled(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CalledFilterer{contract: contract}, nil
}

// bindCalled binds a generic wrapper to an already deployed contract.
func bindCalled(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CalledABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Called *CalledRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Called.Contract.CalledCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Called *CalledRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Called.Contract.CalledTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Called *CalledRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Called.Contract.CalledTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Called *CalledCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Called.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Called *CalledTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Called.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Called *CalledTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Called.Contract.contract.Transact(opts, method, params...)
}

// Entrypoint is a free data retrieval call binding the contract method 0xa65d69d4.
//
// Solidity: function entrypoint() view returns()
func (_Called *CalledCaller) Entrypoint(opts *bind.CallOpts) error {
	var out []interface{}
	err := _Called.contract.Call(opts, &out, "entrypoint")

	if err != nil {
		return err
	}

	return err

}

// Entrypoint is a free data retrieval call binding the contract method 0xa65d69d4.
//
// Solidity: function entrypoint() view returns()
func (_Called *CalledSession) Entrypoint() error {
	return _Called.Contract.Entrypoint(&_Called.CallOpts)
}

// Entrypoint is a free data retrieval call binding the contract method 0xa65d69d4.
//
// Solidity: function entrypoint() view returns()
func (_Called *CalledCallerSession) Entrypoint() error {
	return _Called.Contract.Entrypoint(&_Called.CallOpts)
}

// ExpectedSender is a free data retrieval call binding the contract method 0x1754dba5.
//
// Solidity: function expectedSender() view returns(address)
func (_Called *CalledCaller) ExpectedSender(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Called.contract.Call(opts, &out, "expectedSender")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ExpectedSender is a free data retrieval call binding the contract method 0x1754dba5.
//
// Solidity: function expectedSender() view returns(address)
func (_Called *CalledSession) ExpectedSender() (common.Address, error) {
	return _Called.Contract.ExpectedSender(&_Called.CallOpts)
}

// ExpectedSender is a free data retrieval call binding the contract method 0x1754dba5.
//
// Solidity: function expectedSender() view returns(address)
func (_Called *CalledCallerSession) ExpectedSender() (common.Address, error) {
	return _Called.Contract.ExpectedSender(&_Called.CallOpts)
}
