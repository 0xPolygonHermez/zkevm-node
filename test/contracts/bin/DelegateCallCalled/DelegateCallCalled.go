// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package DelegateCallCalled

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

// DelegateCallCalledMetaData contains all meta data concerning the DelegateCallCalled contract.
var DelegateCallCalledMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"num\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_num\",\"type\":\"uint256\"}],\"name\":\"setVars\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"value\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610125806100206000396000f3fe608060405260043610603a5760003560e01c80633fa4f24514603f5780634e70b1dc1460665780636466414b14607a57806367e404ce1460a2575b600080fd5b348015604a57600080fd5b50605360025481565b6040519081526020015b60405180910390f35b348015607157600080fd5b50605360005481565b60a0608536600460d7565b600055600180546001600160a01b0319163317905534600255565b005b34801560ad57600080fd5b5060015460c0906001600160a01b031681565b6040516001600160a01b039091168152602001605d565b60006020828403121560e857600080fd5b503591905056fea264697066735822122097fd95fb426fdceccd021d918999a67d34c23e6f5d5d1dc7aacc3f9482aa3a6b64736f6c634300080c0033",
}

// DelegateCallCalledABI is the input ABI used to generate the binding from.
// Deprecated: Use DelegateCallCalledMetaData.ABI instead.
var DelegateCallCalledABI = DelegateCallCalledMetaData.ABI

// DelegateCallCalledBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DelegateCallCalledMetaData.Bin instead.
var DelegateCallCalledBin = DelegateCallCalledMetaData.Bin

// DeployDelegateCallCalled deploys a new Ethereum contract, binding an instance of DelegateCallCalled to it.
func DeployDelegateCallCalled(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DelegateCallCalled, error) {
	parsed, err := DelegateCallCalledMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DelegateCallCalledBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DelegateCallCalled{DelegateCallCalledCaller: DelegateCallCalledCaller{contract: contract}, DelegateCallCalledTransactor: DelegateCallCalledTransactor{contract: contract}, DelegateCallCalledFilterer: DelegateCallCalledFilterer{contract: contract}}, nil
}

// DelegateCallCalled is an auto generated Go binding around an Ethereum contract.
type DelegateCallCalled struct {
	DelegateCallCalledCaller     // Read-only binding to the contract
	DelegateCallCalledTransactor // Write-only binding to the contract
	DelegateCallCalledFilterer   // Log filterer for contract events
}

// DelegateCallCalledCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegateCallCalledCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegateCallCalledTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegateCallCalledTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegateCallCalledFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegateCallCalledFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegateCallCalledSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegateCallCalledSession struct {
	Contract     *DelegateCallCalled // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DelegateCallCalledCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegateCallCalledCallerSession struct {
	Contract *DelegateCallCalledCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// DelegateCallCalledTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegateCallCalledTransactorSession struct {
	Contract     *DelegateCallCalledTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// DelegateCallCalledRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegateCallCalledRaw struct {
	Contract *DelegateCallCalled // Generic contract binding to access the raw methods on
}

// DelegateCallCalledCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegateCallCalledCallerRaw struct {
	Contract *DelegateCallCalledCaller // Generic read-only contract binding to access the raw methods on
}

// DelegateCallCalledTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegateCallCalledTransactorRaw struct {
	Contract *DelegateCallCalledTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegateCallCalled creates a new instance of DelegateCallCalled, bound to a specific deployed contract.
func NewDelegateCallCalled(address common.Address, backend bind.ContractBackend) (*DelegateCallCalled, error) {
	contract, err := bindDelegateCallCalled(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DelegateCallCalled{DelegateCallCalledCaller: DelegateCallCalledCaller{contract: contract}, DelegateCallCalledTransactor: DelegateCallCalledTransactor{contract: contract}, DelegateCallCalledFilterer: DelegateCallCalledFilterer{contract: contract}}, nil
}

// NewDelegateCallCalledCaller creates a new read-only instance of DelegateCallCalled, bound to a specific deployed contract.
func NewDelegateCallCalledCaller(address common.Address, caller bind.ContractCaller) (*DelegateCallCalledCaller, error) {
	contract, err := bindDelegateCallCalled(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegateCallCalledCaller{contract: contract}, nil
}

// NewDelegateCallCalledTransactor creates a new write-only instance of DelegateCallCalled, bound to a specific deployed contract.
func NewDelegateCallCalledTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegateCallCalledTransactor, error) {
	contract, err := bindDelegateCallCalled(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegateCallCalledTransactor{contract: contract}, nil
}

// NewDelegateCallCalledFilterer creates a new log filterer instance of DelegateCallCalled, bound to a specific deployed contract.
func NewDelegateCallCalledFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegateCallCalledFilterer, error) {
	contract, err := bindDelegateCallCalled(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegateCallCalledFilterer{contract: contract}, nil
}

// bindDelegateCallCalled binds a generic wrapper to an already deployed contract.
func bindDelegateCallCalled(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DelegateCallCalledMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelegateCallCalled *DelegateCallCalledRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DelegateCallCalled.Contract.DelegateCallCalledCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelegateCallCalled *DelegateCallCalledRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegateCallCalled.Contract.DelegateCallCalledTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelegateCallCalled *DelegateCallCalledRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelegateCallCalled.Contract.DelegateCallCalledTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelegateCallCalled *DelegateCallCalledCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DelegateCallCalled.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelegateCallCalled *DelegateCallCalledTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegateCallCalled.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelegateCallCalled *DelegateCallCalledTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelegateCallCalled.Contract.contract.Transact(opts, method, params...)
}

// Num is a free data retrieval call binding the contract method 0x4e70b1dc.
//
// Solidity: function num() view returns(uint256)
func (_DelegateCallCalled *DelegateCallCalledCaller) Num(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DelegateCallCalled.contract.Call(opts, &out, "num")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Num is a free data retrieval call binding the contract method 0x4e70b1dc.
//
// Solidity: function num() view returns(uint256)
func (_DelegateCallCalled *DelegateCallCalledSession) Num() (*big.Int, error) {
	return _DelegateCallCalled.Contract.Num(&_DelegateCallCalled.CallOpts)
}

// Num is a free data retrieval call binding the contract method 0x4e70b1dc.
//
// Solidity: function num() view returns(uint256)
func (_DelegateCallCalled *DelegateCallCalledCallerSession) Num() (*big.Int, error) {
	return _DelegateCallCalled.Contract.Num(&_DelegateCallCalled.CallOpts)
}

// Sender is a free data retrieval call binding the contract method 0x67e404ce.
//
// Solidity: function sender() view returns(address)
func (_DelegateCallCalled *DelegateCallCalledCaller) Sender(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DelegateCallCalled.contract.Call(opts, &out, "sender")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Sender is a free data retrieval call binding the contract method 0x67e404ce.
//
// Solidity: function sender() view returns(address)
func (_DelegateCallCalled *DelegateCallCalledSession) Sender() (common.Address, error) {
	return _DelegateCallCalled.Contract.Sender(&_DelegateCallCalled.CallOpts)
}

// Sender is a free data retrieval call binding the contract method 0x67e404ce.
//
// Solidity: function sender() view returns(address)
func (_DelegateCallCalled *DelegateCallCalledCallerSession) Sender() (common.Address, error) {
	return _DelegateCallCalled.Contract.Sender(&_DelegateCallCalled.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(uint256)
func (_DelegateCallCalled *DelegateCallCalledCaller) Value(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DelegateCallCalled.contract.Call(opts, &out, "value")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(uint256)
func (_DelegateCallCalled *DelegateCallCalledSession) Value() (*big.Int, error) {
	return _DelegateCallCalled.Contract.Value(&_DelegateCallCalled.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(uint256)
func (_DelegateCallCalled *DelegateCallCalledCallerSession) Value() (*big.Int, error) {
	return _DelegateCallCalled.Contract.Value(&_DelegateCallCalled.CallOpts)
}

// SetVars is a paid mutator transaction binding the contract method 0x6466414b.
//
// Solidity: function setVars(uint256 _num) payable returns()
func (_DelegateCallCalled *DelegateCallCalledTransactor) SetVars(opts *bind.TransactOpts, _num *big.Int) (*types.Transaction, error) {
	return _DelegateCallCalled.contract.Transact(opts, "setVars", _num)
}

// SetVars is a paid mutator transaction binding the contract method 0x6466414b.
//
// Solidity: function setVars(uint256 _num) payable returns()
func (_DelegateCallCalled *DelegateCallCalledSession) SetVars(_num *big.Int) (*types.Transaction, error) {
	return _DelegateCallCalled.Contract.SetVars(&_DelegateCallCalled.TransactOpts, _num)
}

// SetVars is a paid mutator transaction binding the contract method 0x6466414b.
//
// Solidity: function setVars(uint256 _num) payable returns()
func (_DelegateCallCalled *DelegateCallCalledTransactorSession) SetVars(_num *big.Int) (*types.Transaction, error) {
	return _DelegateCallCalled.Contract.SetVars(&_DelegateCallCalled.TransactOpts, _num)
}
