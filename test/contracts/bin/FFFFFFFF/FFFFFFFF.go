// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package FFFFFFFF

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

// FFFFFFFFMetaData contains all meta data concerning the FFFFFFFF contract.
var FFFFFFFFMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506000196000f3fe",
}

// FFFFFFFFABI is the input ABI used to generate the binding from.
// Deprecated: Use FFFFFFFFMetaData.ABI instead.
var FFFFFFFFABI = FFFFFFFFMetaData.ABI

// FFFFFFFFBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FFFFFFFFMetaData.Bin instead.
var FFFFFFFFBin = FFFFFFFFMetaData.Bin

// DeployFFFFFFFF deploys a new Ethereum contract, binding an instance of FFFFFFFF to it.
func DeployFFFFFFFF(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FFFFFFFF, error) {
	parsed, err := FFFFFFFFMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FFFFFFFFBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FFFFFFFF{FFFFFFFFCaller: FFFFFFFFCaller{contract: contract}, FFFFFFFFTransactor: FFFFFFFFTransactor{contract: contract}, FFFFFFFFFilterer: FFFFFFFFFilterer{contract: contract}}, nil
}

// FFFFFFFF is an auto generated Go binding around an Ethereum contract.
type FFFFFFFF struct {
	FFFFFFFFCaller     // Read-only binding to the contract
	FFFFFFFFTransactor // Write-only binding to the contract
	FFFFFFFFFilterer   // Log filterer for contract events
}

// FFFFFFFFCaller is an auto generated read-only Go binding around an Ethereum contract.
type FFFFFFFFCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FFFFFFFFTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FFFFFFFFTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FFFFFFFFFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FFFFFFFFFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FFFFFFFFSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FFFFFFFFSession struct {
	Contract     *FFFFFFFF         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FFFFFFFFCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FFFFFFFFCallerSession struct {
	Contract *FFFFFFFFCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// FFFFFFFFTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FFFFFFFFTransactorSession struct {
	Contract     *FFFFFFFFTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// FFFFFFFFRaw is an auto generated low-level Go binding around an Ethereum contract.
type FFFFFFFFRaw struct {
	Contract *FFFFFFFF // Generic contract binding to access the raw methods on
}

// FFFFFFFFCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FFFFFFFFCallerRaw struct {
	Contract *FFFFFFFFCaller // Generic read-only contract binding to access the raw methods on
}

// FFFFFFFFTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FFFFFFFFTransactorRaw struct {
	Contract *FFFFFFFFTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFFFFFFFF creates a new instance of FFFFFFFF, bound to a specific deployed contract.
func NewFFFFFFFF(address common.Address, backend bind.ContractBackend) (*FFFFFFFF, error) {
	contract, err := bindFFFFFFFF(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FFFFFFFF{FFFFFFFFCaller: FFFFFFFFCaller{contract: contract}, FFFFFFFFTransactor: FFFFFFFFTransactor{contract: contract}, FFFFFFFFFilterer: FFFFFFFFFilterer{contract: contract}}, nil
}

// NewFFFFFFFFCaller creates a new read-only instance of FFFFFFFF, bound to a specific deployed contract.
func NewFFFFFFFFCaller(address common.Address, caller bind.ContractCaller) (*FFFFFFFFCaller, error) {
	contract, err := bindFFFFFFFF(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FFFFFFFFCaller{contract: contract}, nil
}

// NewFFFFFFFFTransactor creates a new write-only instance of FFFFFFFF, bound to a specific deployed contract.
func NewFFFFFFFFTransactor(address common.Address, transactor bind.ContractTransactor) (*FFFFFFFFTransactor, error) {
	contract, err := bindFFFFFFFF(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FFFFFFFFTransactor{contract: contract}, nil
}

// NewFFFFFFFFFilterer creates a new log filterer instance of FFFFFFFF, bound to a specific deployed contract.
func NewFFFFFFFFFilterer(address common.Address, filterer bind.ContractFilterer) (*FFFFFFFFFilterer, error) {
	contract, err := bindFFFFFFFF(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FFFFFFFFFilterer{contract: contract}, nil
}

// bindFFFFFFFF binds a generic wrapper to an already deployed contract.
func bindFFFFFFFF(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FFFFFFFFMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FFFFFFFF *FFFFFFFFRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FFFFFFFF.Contract.FFFFFFFFCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FFFFFFFF *FFFFFFFFRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FFFFFFFF.Contract.FFFFFFFFTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FFFFFFFF *FFFFFFFFRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FFFFFFFF.Contract.FFFFFFFFTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FFFFFFFF *FFFFFFFFCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FFFFFFFF.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FFFFFFFF *FFFFFFFFTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FFFFFFFF.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FFFFFFFF *FFFFFFFFTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FFFFFFFF.Contract.contract.Transact(opts, method, params...)
}
