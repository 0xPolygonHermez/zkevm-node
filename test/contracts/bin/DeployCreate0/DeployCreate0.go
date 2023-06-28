// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package DeployCreate0

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

// DeployCreate0MetaData contains all meta data concerning the DeployCreate0 contract.
var DeployCreate0MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506000806000f050603f8060246000396000f3fe6080604052600080fdfea26469706673582212208401af456e61e3b8cef3f26d4eaad7a4a75983147d147f7b6ec8312e0255cb9864736f6c634300080c0033",
}

// DeployCreate0ABI is the input ABI used to generate the binding from.
// Deprecated: Use DeployCreate0MetaData.ABI instead.
var DeployCreate0ABI = DeployCreate0MetaData.ABI

// DeployCreate0Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DeployCreate0MetaData.Bin instead.
var DeployCreate0Bin = DeployCreate0MetaData.Bin

// DeployDeployCreate0 deploys a new Ethereum contract, binding an instance of DeployCreate0 to it.
func DeployDeployCreate0(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DeployCreate0, error) {
	parsed, err := DeployCreate0MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DeployCreate0Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DeployCreate0{DeployCreate0Caller: DeployCreate0Caller{contract: contract}, DeployCreate0Transactor: DeployCreate0Transactor{contract: contract}, DeployCreate0Filterer: DeployCreate0Filterer{contract: contract}}, nil
}

// DeployCreate0 is an auto generated Go binding around an Ethereum contract.
type DeployCreate0 struct {
	DeployCreate0Caller     // Read-only binding to the contract
	DeployCreate0Transactor // Write-only binding to the contract
	DeployCreate0Filterer   // Log filterer for contract events
}

// DeployCreate0Caller is an auto generated read-only Go binding around an Ethereum contract.
type DeployCreate0Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployCreate0Transactor is an auto generated write-only Go binding around an Ethereum contract.
type DeployCreate0Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployCreate0Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DeployCreate0Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployCreate0Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DeployCreate0Session struct {
	Contract     *DeployCreate0    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DeployCreate0CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DeployCreate0CallerSession struct {
	Contract *DeployCreate0Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// DeployCreate0TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DeployCreate0TransactorSession struct {
	Contract     *DeployCreate0Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// DeployCreate0Raw is an auto generated low-level Go binding around an Ethereum contract.
type DeployCreate0Raw struct {
	Contract *DeployCreate0 // Generic contract binding to access the raw methods on
}

// DeployCreate0CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DeployCreate0CallerRaw struct {
	Contract *DeployCreate0Caller // Generic read-only contract binding to access the raw methods on
}

// DeployCreate0TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DeployCreate0TransactorRaw struct {
	Contract *DeployCreate0Transactor // Generic write-only contract binding to access the raw methods on
}

// NewDeployCreate0 creates a new instance of DeployCreate0, bound to a specific deployed contract.
func NewDeployCreate0(address common.Address, backend bind.ContractBackend) (*DeployCreate0, error) {
	contract, err := bindDeployCreate0(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DeployCreate0{DeployCreate0Caller: DeployCreate0Caller{contract: contract}, DeployCreate0Transactor: DeployCreate0Transactor{contract: contract}, DeployCreate0Filterer: DeployCreate0Filterer{contract: contract}}, nil
}

// NewDeployCreate0Caller creates a new read-only instance of DeployCreate0, bound to a specific deployed contract.
func NewDeployCreate0Caller(address common.Address, caller bind.ContractCaller) (*DeployCreate0Caller, error) {
	contract, err := bindDeployCreate0(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DeployCreate0Caller{contract: contract}, nil
}

// NewDeployCreate0Transactor creates a new write-only instance of DeployCreate0, bound to a specific deployed contract.
func NewDeployCreate0Transactor(address common.Address, transactor bind.ContractTransactor) (*DeployCreate0Transactor, error) {
	contract, err := bindDeployCreate0(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DeployCreate0Transactor{contract: contract}, nil
}

// NewDeployCreate0Filterer creates a new log filterer instance of DeployCreate0, bound to a specific deployed contract.
func NewDeployCreate0Filterer(address common.Address, filterer bind.ContractFilterer) (*DeployCreate0Filterer, error) {
	contract, err := bindDeployCreate0(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DeployCreate0Filterer{contract: contract}, nil
}

// bindDeployCreate0 binds a generic wrapper to an already deployed contract.
func bindDeployCreate0(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DeployCreate0MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DeployCreate0 *DeployCreate0Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DeployCreate0.Contract.DeployCreate0Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DeployCreate0 *DeployCreate0Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployCreate0.Contract.DeployCreate0Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DeployCreate0 *DeployCreate0Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DeployCreate0.Contract.DeployCreate0Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DeployCreate0 *DeployCreate0CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DeployCreate0.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DeployCreate0 *DeployCreate0TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DeployCreate0.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DeployCreate0 *DeployCreate0TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DeployCreate0.Contract.contract.Transact(opts, method, params...)
}
