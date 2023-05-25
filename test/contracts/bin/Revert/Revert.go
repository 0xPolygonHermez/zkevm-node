// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Revert

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

// RevertMetaData contains all meta data concerning the Revert contract.
var RevertMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]",
	Bin: "0x6080604052348015600f57600080fd5b5060405162461bcd60e51b815260206004820152601460248201527f546f646179206973206e6f74206a7565726e6573000000000000000000000000604482015260640160405180910390fdfe",
}

// RevertABI is the input ABI used to generate the binding from.
// Deprecated: Use RevertMetaData.ABI instead.
var RevertABI = RevertMetaData.ABI

// RevertBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RevertMetaData.Bin instead.
var RevertBin = RevertMetaData.Bin

// DeployRevert deploys a new Ethereum contract, binding an instance of Revert to it.
func DeployRevert(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Revert, error) {
	parsed, err := RevertMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RevertBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Revert{RevertCaller: RevertCaller{contract: contract}, RevertTransactor: RevertTransactor{contract: contract}, RevertFilterer: RevertFilterer{contract: contract}}, nil
}

// Revert is an auto generated Go binding around an Ethereum contract.
type Revert struct {
	RevertCaller     // Read-only binding to the contract
	RevertTransactor // Write-only binding to the contract
	RevertFilterer   // Log filterer for contract events
}

// RevertCaller is an auto generated read-only Go binding around an Ethereum contract.
type RevertCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RevertTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RevertTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RevertFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RevertFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RevertSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RevertSession struct {
	Contract     *Revert           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RevertCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RevertCallerSession struct {
	Contract *RevertCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RevertTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RevertTransactorSession struct {
	Contract     *RevertTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RevertRaw is an auto generated low-level Go binding around an Ethereum contract.
type RevertRaw struct {
	Contract *Revert // Generic contract binding to access the raw methods on
}

// RevertCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RevertCallerRaw struct {
	Contract *RevertCaller // Generic read-only contract binding to access the raw methods on
}

// RevertTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RevertTransactorRaw struct {
	Contract *RevertTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRevert creates a new instance of Revert, bound to a specific deployed contract.
func NewRevert(address common.Address, backend bind.ContractBackend) (*Revert, error) {
	contract, err := bindRevert(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Revert{RevertCaller: RevertCaller{contract: contract}, RevertTransactor: RevertTransactor{contract: contract}, RevertFilterer: RevertFilterer{contract: contract}}, nil
}

// NewRevertCaller creates a new read-only instance of Revert, bound to a specific deployed contract.
func NewRevertCaller(address common.Address, caller bind.ContractCaller) (*RevertCaller, error) {
	contract, err := bindRevert(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RevertCaller{contract: contract}, nil
}

// NewRevertTransactor creates a new write-only instance of Revert, bound to a specific deployed contract.
func NewRevertTransactor(address common.Address, transactor bind.ContractTransactor) (*RevertTransactor, error) {
	contract, err := bindRevert(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RevertTransactor{contract: contract}, nil
}

// NewRevertFilterer creates a new log filterer instance of Revert, bound to a specific deployed contract.
func NewRevertFilterer(address common.Address, filterer bind.ContractFilterer) (*RevertFilterer, error) {
	contract, err := bindRevert(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RevertFilterer{contract: contract}, nil
}

// bindRevert binds a generic wrapper to an already deployed contract.
func bindRevert(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RevertMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Revert *RevertRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Revert.Contract.RevertCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Revert *RevertRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Revert.Contract.RevertTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Revert *RevertRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Revert.Contract.RevertTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Revert *RevertCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Revert.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Revert *RevertTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Revert.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Revert *RevertTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Revert.Contract.contract.Transact(opts, method, params...)
}
