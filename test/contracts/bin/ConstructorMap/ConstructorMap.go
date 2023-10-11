// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ConstructorMap

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

// ConstructorMapMetaData contains all meta data concerning the ConstructorMap contract.
var ConstructorMapMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"numbers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060005b60648110156100405760008181526020819052604090208190558061003881610046565b915050610014565b5061006f565b600060001982141561006857634e487b7160e01b600052601160045260246000fd5b5060010190565b60aa8061007d6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063d39fa23314602d575b600080fd5b604a6038366004605c565b60006020819052908152604090205481565b60405190815260200160405180910390f35b600060208284031215606d57600080fd5b503591905056fea26469706673582212207164b7e8cab7019534d840c5be1f93a98671cdbddc7ea08c6a73b67022062ee864736f6c634300080c0033",
}

// ConstructorMapABI is the input ABI used to generate the binding from.
// Deprecated: Use ConstructorMapMetaData.ABI instead.
var ConstructorMapABI = ConstructorMapMetaData.ABI

// ConstructorMapBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConstructorMapMetaData.Bin instead.
var ConstructorMapBin = ConstructorMapMetaData.Bin

// DeployConstructorMap deploys a new Ethereum contract, binding an instance of ConstructorMap to it.
func DeployConstructorMap(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ConstructorMap, error) {
	parsed, err := ConstructorMapMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConstructorMapBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConstructorMap{ConstructorMapCaller: ConstructorMapCaller{contract: contract}, ConstructorMapTransactor: ConstructorMapTransactor{contract: contract}, ConstructorMapFilterer: ConstructorMapFilterer{contract: contract}}, nil
}

// ConstructorMap is an auto generated Go binding around an Ethereum contract.
type ConstructorMap struct {
	ConstructorMapCaller     // Read-only binding to the contract
	ConstructorMapTransactor // Write-only binding to the contract
	ConstructorMapFilterer   // Log filterer for contract events
}

// ConstructorMapCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConstructorMapCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstructorMapTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConstructorMapTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstructorMapFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConstructorMapFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConstructorMapSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConstructorMapSession struct {
	Contract     *ConstructorMap   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConstructorMapCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConstructorMapCallerSession struct {
	Contract *ConstructorMapCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ConstructorMapTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConstructorMapTransactorSession struct {
	Contract     *ConstructorMapTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ConstructorMapRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConstructorMapRaw struct {
	Contract *ConstructorMap // Generic contract binding to access the raw methods on
}

// ConstructorMapCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConstructorMapCallerRaw struct {
	Contract *ConstructorMapCaller // Generic read-only contract binding to access the raw methods on
}

// ConstructorMapTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConstructorMapTransactorRaw struct {
	Contract *ConstructorMapTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConstructorMap creates a new instance of ConstructorMap, bound to a specific deployed contract.
func NewConstructorMap(address common.Address, backend bind.ContractBackend) (*ConstructorMap, error) {
	contract, err := bindConstructorMap(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConstructorMap{ConstructorMapCaller: ConstructorMapCaller{contract: contract}, ConstructorMapTransactor: ConstructorMapTransactor{contract: contract}, ConstructorMapFilterer: ConstructorMapFilterer{contract: contract}}, nil
}

// NewConstructorMapCaller creates a new read-only instance of ConstructorMap, bound to a specific deployed contract.
func NewConstructorMapCaller(address common.Address, caller bind.ContractCaller) (*ConstructorMapCaller, error) {
	contract, err := bindConstructorMap(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConstructorMapCaller{contract: contract}, nil
}

// NewConstructorMapTransactor creates a new write-only instance of ConstructorMap, bound to a specific deployed contract.
func NewConstructorMapTransactor(address common.Address, transactor bind.ContractTransactor) (*ConstructorMapTransactor, error) {
	contract, err := bindConstructorMap(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConstructorMapTransactor{contract: contract}, nil
}

// NewConstructorMapFilterer creates a new log filterer instance of ConstructorMap, bound to a specific deployed contract.
func NewConstructorMapFilterer(address common.Address, filterer bind.ContractFilterer) (*ConstructorMapFilterer, error) {
	contract, err := bindConstructorMap(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConstructorMapFilterer{contract: contract}, nil
}

// bindConstructorMap binds a generic wrapper to an already deployed contract.
func bindConstructorMap(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConstructorMapMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstructorMap *ConstructorMapRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConstructorMap.Contract.ConstructorMapCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstructorMap *ConstructorMapRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstructorMap.Contract.ConstructorMapTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstructorMap *ConstructorMapRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstructorMap.Contract.ConstructorMapTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConstructorMap *ConstructorMapCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConstructorMap.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConstructorMap *ConstructorMapTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConstructorMap.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConstructorMap *ConstructorMapTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConstructorMap.Contract.contract.Transact(opts, method, params...)
}

// Numbers is a free data retrieval call binding the contract method 0xd39fa233.
//
// Solidity: function numbers(uint256 ) view returns(uint256)
func (_ConstructorMap *ConstructorMapCaller) Numbers(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ConstructorMap.contract.Call(opts, &out, "numbers", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Numbers is a free data retrieval call binding the contract method 0xd39fa233.
//
// Solidity: function numbers(uint256 ) view returns(uint256)
func (_ConstructorMap *ConstructorMapSession) Numbers(arg0 *big.Int) (*big.Int, error) {
	return _ConstructorMap.Contract.Numbers(&_ConstructorMap.CallOpts, arg0)
}

// Numbers is a free data retrieval call binding the contract method 0xd39fa233.
//
// Solidity: function numbers(uint256 ) view returns(uint256)
func (_ConstructorMap *ConstructorMapCallerSession) Numbers(arg0 *big.Int) (*big.Int, error) {
	return _ConstructorMap.Contract.Numbers(&_ConstructorMap.CallOpts, arg0)
}
