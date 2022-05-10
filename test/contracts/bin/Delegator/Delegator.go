// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Delegator

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

// DelegatorMetaData contains all meta data concerning the Delegator contract.
var DelegatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_expectedSender\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"expectedSender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161028e38038061028e83398101604081905261002f91610054565b600080546001600160a01b0319166001600160a01b0392909216919091179055610084565b60006020828403121561006657600080fd5b81516001600160a01b038116811461007d57600080fd5b9392505050565b6101fb806100936000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80631754dba51461003b578063f55332ab1461006a575b600080fd5b60005461004e906001600160a01b031681565b6040516001600160a01b03909116815260200160405180910390f35b61007d61007836600461015a565b61007f565b005b60408051600481526024810182526020810180516001600160e01b03166329975a7560e21b17905290516000916001600160a01b038416916100c1919061018a565b600060405180830381855af49150503d80600081146100fc576040519150601f19603f3d011682016040523d82523d6000602084013e610101565b606091505b50509050806101565760405162461bcd60e51b815260206004820152601c60248201527f657870656374656453656e64657220213d206d73672e73656e64657200000000604482015260640160405180910390fd5b5050565b60006020828403121561016c57600080fd5b81356001600160a01b038116811461018357600080fd5b9392505050565b6000825160005b818110156101ab5760208186018101518583015201610191565b818111156101ba576000828501525b50919091019291505056fea2646970667358221220104afe8ac68c5f3d48010505c00610059b8dba04a90ef22b9befcc54301f9a4d64736f6c634300080c0033",
}

// DelegatorABI is the input ABI used to generate the binding from.
// Deprecated: Use DelegatorMetaData.ABI instead.
var DelegatorABI = DelegatorMetaData.ABI

// DelegatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DelegatorMetaData.Bin instead.
var DelegatorBin = DelegatorMetaData.Bin

// DeployDelegator deploys a new Ethereum contract, binding an instance of Delegator to it.
func DeployDelegator(auth *bind.TransactOpts, backend bind.ContractBackend, _expectedSender common.Address) (common.Address, *types.Transaction, *Delegator, error) {
	parsed, err := DelegatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DelegatorBin), backend, _expectedSender)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Delegator{DelegatorCaller: DelegatorCaller{contract: contract}, DelegatorTransactor: DelegatorTransactor{contract: contract}, DelegatorFilterer: DelegatorFilterer{contract: contract}}, nil
}

// Delegator is an auto generated Go binding around an Ethereum contract.
type Delegator struct {
	DelegatorCaller     // Read-only binding to the contract
	DelegatorTransactor // Write-only binding to the contract
	DelegatorFilterer   // Log filterer for contract events
}

// DelegatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegatorSession struct {
	Contract     *Delegator        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DelegatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegatorCallerSession struct {
	Contract *DelegatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// DelegatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegatorTransactorSession struct {
	Contract     *DelegatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// DelegatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegatorRaw struct {
	Contract *Delegator // Generic contract binding to access the raw methods on
}

// DelegatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegatorCallerRaw struct {
	Contract *DelegatorCaller // Generic read-only contract binding to access the raw methods on
}

// DelegatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegatorTransactorRaw struct {
	Contract *DelegatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegator creates a new instance of Delegator, bound to a specific deployed contract.
func NewDelegator(address common.Address, backend bind.ContractBackend) (*Delegator, error) {
	contract, err := bindDelegator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Delegator{DelegatorCaller: DelegatorCaller{contract: contract}, DelegatorTransactor: DelegatorTransactor{contract: contract}, DelegatorFilterer: DelegatorFilterer{contract: contract}}, nil
}

// NewDelegatorCaller creates a new read-only instance of Delegator, bound to a specific deployed contract.
func NewDelegatorCaller(address common.Address, caller bind.ContractCaller) (*DelegatorCaller, error) {
	contract, err := bindDelegator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegatorCaller{contract: contract}, nil
}

// NewDelegatorTransactor creates a new write-only instance of Delegator, bound to a specific deployed contract.
func NewDelegatorTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegatorTransactor, error) {
	contract, err := bindDelegator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegatorTransactor{contract: contract}, nil
}

// NewDelegatorFilterer creates a new log filterer instance of Delegator, bound to a specific deployed contract.
func NewDelegatorFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegatorFilterer, error) {
	contract, err := bindDelegator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegatorFilterer{contract: contract}, nil
}

// bindDelegator binds a generic wrapper to an already deployed contract.
func bindDelegator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DelegatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegator *DelegatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegator.Contract.DelegatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegator *DelegatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegator.Contract.DelegatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegator *DelegatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegator.Contract.DelegatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegator *DelegatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegator *DelegatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegator *DelegatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegator.Contract.contract.Transact(opts, method, params...)
}

// ExpectedSender is a free data retrieval call binding the contract method 0x1754dba5.
//
// Solidity: function expectedSender() view returns(address)
func (_Delegator *DelegatorCaller) ExpectedSender(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Delegator.contract.Call(opts, &out, "expectedSender")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ExpectedSender is a free data retrieval call binding the contract method 0x1754dba5.
//
// Solidity: function expectedSender() view returns(address)
func (_Delegator *DelegatorSession) ExpectedSender() (common.Address, error) {
	return _Delegator.Contract.ExpectedSender(&_Delegator.CallOpts)
}

// ExpectedSender is a free data retrieval call binding the contract method 0x1754dba5.
//
// Solidity: function expectedSender() view returns(address)
func (_Delegator *DelegatorCallerSession) ExpectedSender() (common.Address, error) {
	return _Delegator.Contract.ExpectedSender(&_Delegator.CallOpts)
}

// Call is a paid mutator transaction binding the contract method 0xf55332ab.
//
// Solidity: function call(address target) returns()
func (_Delegator *DelegatorTransactor) Call(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _Delegator.contract.Transact(opts, "call", target)
}

// Call is a paid mutator transaction binding the contract method 0xf55332ab.
//
// Solidity: function call(address target) returns()
func (_Delegator *DelegatorSession) Call(target common.Address) (*types.Transaction, error) {
	return _Delegator.Contract.Call(&_Delegator.TransactOpts, target)
}

// Call is a paid mutator transaction binding the contract method 0xf55332ab.
//
// Solidity: function call(address target) returns()
func (_Delegator *DelegatorTransactorSession) Call(target common.Address) (*types.Transaction, error) {
	return _Delegator.Contract.Call(&_Delegator.TransactOpts, target)
}
