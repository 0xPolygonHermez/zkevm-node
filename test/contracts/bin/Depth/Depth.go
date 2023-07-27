// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Depth

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

// DepthMetaData contains all meta data concerning the Depth contract.
var DepthMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"secondCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasForwarded\",\"type\":\"uint256\"}],\"name\":\"start\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526000805534801561001457600080fd5b506101c1806100246000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80633c2d70e91461003b57806349a0639314610050575b600080fd5b61004e610049366004610126565b610075565b005b61006361005e366004610150565b6100e9565b60405190815260200160405180910390f35b6040516349a0639360e01b81526001600160a01b038316600482015230906349a0639390839060240160206040518083038160008887f11580156100bd573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906100e29190610172565b6000555050565b600080636aecbc3360e01b6080526020608060046080865afa905050919050565b80356001600160a01b038116811461012157600080fd5b919050565b6000806040838503121561013957600080fd5b6101428361010a565b946020939093013593505050565b60006020828403121561016257600080fd5b61016b8261010a565b9392505050565b60006020828403121561018457600080fd5b505191905056fea2646970667358221220572542f847a41e898ee3794ab8f893d4816858903ee722c477cf6c0583b61acf64736f6c634300080c0033",
}

// DepthABI is the input ABI used to generate the binding from.
// Deprecated: Use DepthMetaData.ABI instead.
var DepthABI = DepthMetaData.ABI

// DepthBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DepthMetaData.Bin instead.
var DepthBin = DepthMetaData.Bin

// DeployDepth deploys a new Ethereum contract, binding an instance of Depth to it.
func DeployDepth(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Depth, error) {
	parsed, err := DepthMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DepthBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Depth{DepthCaller: DepthCaller{contract: contract}, DepthTransactor: DepthTransactor{contract: contract}, DepthFilterer: DepthFilterer{contract: contract}}, nil
}

// Depth is an auto generated Go binding around an Ethereum contract.
type Depth struct {
	DepthCaller     // Read-only binding to the contract
	DepthTransactor // Write-only binding to the contract
	DepthFilterer   // Log filterer for contract events
}

// DepthCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepthCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepthTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepthTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepthFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepthFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepthSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepthSession struct {
	Contract     *Depth            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepthCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepthCallerSession struct {
	Contract *DepthCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DepthTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepthTransactorSession struct {
	Contract     *DepthTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepthRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepthRaw struct {
	Contract *Depth // Generic contract binding to access the raw methods on
}

// DepthCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepthCallerRaw struct {
	Contract *DepthCaller // Generic read-only contract binding to access the raw methods on
}

// DepthTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepthTransactorRaw struct {
	Contract *DepthTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDepth creates a new instance of Depth, bound to a specific deployed contract.
func NewDepth(address common.Address, backend bind.ContractBackend) (*Depth, error) {
	contract, err := bindDepth(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Depth{DepthCaller: DepthCaller{contract: contract}, DepthTransactor: DepthTransactor{contract: contract}, DepthFilterer: DepthFilterer{contract: contract}}, nil
}

// NewDepthCaller creates a new read-only instance of Depth, bound to a specific deployed contract.
func NewDepthCaller(address common.Address, caller bind.ContractCaller) (*DepthCaller, error) {
	contract, err := bindDepth(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepthCaller{contract: contract}, nil
}

// NewDepthTransactor creates a new write-only instance of Depth, bound to a specific deployed contract.
func NewDepthTransactor(address common.Address, transactor bind.ContractTransactor) (*DepthTransactor, error) {
	contract, err := bindDepth(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepthTransactor{contract: contract}, nil
}

// NewDepthFilterer creates a new log filterer instance of Depth, bound to a specific deployed contract.
func NewDepthFilterer(address common.Address, filterer bind.ContractFilterer) (*DepthFilterer, error) {
	contract, err := bindDepth(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepthFilterer{contract: contract}, nil
}

// bindDepth binds a generic wrapper to an already deployed contract.
func bindDepth(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DepthMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Depth *DepthRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Depth.Contract.DepthCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Depth *DepthRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depth.Contract.DepthTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Depth *DepthRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Depth.Contract.DepthTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Depth *DepthCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Depth.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Depth *DepthTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depth.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Depth *DepthTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Depth.Contract.contract.Transact(opts, method, params...)
}

// SecondCall is a paid mutator transaction binding the contract method 0x49a06393.
//
// Solidity: function secondCall(address addr) returns(uint256)
func (_Depth *DepthTransactor) SecondCall(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _Depth.contract.Transact(opts, "secondCall", addr)
}

// SecondCall is a paid mutator transaction binding the contract method 0x49a06393.
//
// Solidity: function secondCall(address addr) returns(uint256)
func (_Depth *DepthSession) SecondCall(addr common.Address) (*types.Transaction, error) {
	return _Depth.Contract.SecondCall(&_Depth.TransactOpts, addr)
}

// SecondCall is a paid mutator transaction binding the contract method 0x49a06393.
//
// Solidity: function secondCall(address addr) returns(uint256)
func (_Depth *DepthTransactorSession) SecondCall(addr common.Address) (*types.Transaction, error) {
	return _Depth.Contract.SecondCall(&_Depth.TransactOpts, addr)
}

// Start is a paid mutator transaction binding the contract method 0x3c2d70e9.
//
// Solidity: function start(address addr, uint256 gasForwarded) returns()
func (_Depth *DepthTransactor) Start(opts *bind.TransactOpts, addr common.Address, gasForwarded *big.Int) (*types.Transaction, error) {
	return _Depth.contract.Transact(opts, "start", addr, gasForwarded)
}

// Start is a paid mutator transaction binding the contract method 0x3c2d70e9.
//
// Solidity: function start(address addr, uint256 gasForwarded) returns()
func (_Depth *DepthSession) Start(addr common.Address, gasForwarded *big.Int) (*types.Transaction, error) {
	return _Depth.Contract.Start(&_Depth.TransactOpts, addr, gasForwarded)
}

// Start is a paid mutator transaction binding the contract method 0x3c2d70e9.
//
// Solidity: function start(address addr, uint256 gasForwarded) returns()
func (_Depth *DepthTransactorSession) Start(addr common.Address, gasForwarded *big.Int) (*types.Transaction, error) {
	return _Depth.Contract.Start(&_Depth.TransactOpts, addr, gasForwarded)
}
