// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package OpCallAux

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

// OpCallAuxMetaData contains all meta data concerning the OpCallAux contract.
var OpCallAuxMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"}],\"name\":\"addTwo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"auxFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"auxReturn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"auxStop\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"auxUpdate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"auxUpdateValues\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"opCallSelfBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"opDelegateSelfBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addrCall\",\"type\":\"address\"}],\"name\":\"opReturnCallSelfBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600160005534801561001557600080fd5b50610215806100256000396000f3fe6080604052600436106100865760003560e01c80633182d9a9116100595780633182d9a9146100e95780636aecbc331461010f5780638b00fb6a14610127578063b6f54b251461008b578063f80efde51461013a57600080fd5b806303cfd8461461008b57806303d3e8a7146100a45780630f52d66e146100b4578063266f6aec146100d4575b600080fd5b4760009081555b60405190815260200160405180910390f35b3480156100b057600080fd5b505b005b3480156100c057600080fd5b506100926100cf36600461016e565b61015b565b3480156100e057600080fd5b506100b2600080fd5b3480156100f557600080fd5b506912121212121212121212600055640123456689610092565b34801561011b57600080fd5b50640123456689610092565b610092610135366004610190565b504790565b69121212121212121212126000553360015534600255640123456689610092565b600061016782846101b9565b9392505050565b6000806040838503121561018157600080fd5b50508035926020909101359150565b6000602082840312156101a257600080fd5b81356001600160a01b038116811461016757600080fd5b600082198211156101da57634e487b7160e01b600052601160045260246000fd5b50019056fea2646970667358221220b404792d51bc77b2b654043a0d459e75c40efdd34d1c2bf01c6012c30bee9ae064736f6c634300080c0033",
}

// OpCallAuxABI is the input ABI used to generate the binding from.
// Deprecated: Use OpCallAuxMetaData.ABI instead.
var OpCallAuxABI = OpCallAuxMetaData.ABI

// OpCallAuxBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OpCallAuxMetaData.Bin instead.
var OpCallAuxBin = OpCallAuxMetaData.Bin

// DeployOpCallAux deploys a new Ethereum contract, binding an instance of OpCallAux to it.
func DeployOpCallAux(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OpCallAux, error) {
	parsed, err := OpCallAuxMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OpCallAuxBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OpCallAux{OpCallAuxCaller: OpCallAuxCaller{contract: contract}, OpCallAuxTransactor: OpCallAuxTransactor{contract: contract}, OpCallAuxFilterer: OpCallAuxFilterer{contract: contract}}, nil
}

// OpCallAux is an auto generated Go binding around an Ethereum contract.
type OpCallAux struct {
	OpCallAuxCaller     // Read-only binding to the contract
	OpCallAuxTransactor // Write-only binding to the contract
	OpCallAuxFilterer   // Log filterer for contract events
}

// OpCallAuxCaller is an auto generated read-only Go binding around an Ethereum contract.
type OpCallAuxCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpCallAuxTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OpCallAuxTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpCallAuxFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OpCallAuxFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpCallAuxSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OpCallAuxSession struct {
	Contract     *OpCallAux        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OpCallAuxCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OpCallAuxCallerSession struct {
	Contract *OpCallAuxCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// OpCallAuxTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OpCallAuxTransactorSession struct {
	Contract     *OpCallAuxTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// OpCallAuxRaw is an auto generated low-level Go binding around an Ethereum contract.
type OpCallAuxRaw struct {
	Contract *OpCallAux // Generic contract binding to access the raw methods on
}

// OpCallAuxCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OpCallAuxCallerRaw struct {
	Contract *OpCallAuxCaller // Generic read-only contract binding to access the raw methods on
}

// OpCallAuxTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OpCallAuxTransactorRaw struct {
	Contract *OpCallAuxTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOpCallAux creates a new instance of OpCallAux, bound to a specific deployed contract.
func NewOpCallAux(address common.Address, backend bind.ContractBackend) (*OpCallAux, error) {
	contract, err := bindOpCallAux(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OpCallAux{OpCallAuxCaller: OpCallAuxCaller{contract: contract}, OpCallAuxTransactor: OpCallAuxTransactor{contract: contract}, OpCallAuxFilterer: OpCallAuxFilterer{contract: contract}}, nil
}

// NewOpCallAuxCaller creates a new read-only instance of OpCallAux, bound to a specific deployed contract.
func NewOpCallAuxCaller(address common.Address, caller bind.ContractCaller) (*OpCallAuxCaller, error) {
	contract, err := bindOpCallAux(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OpCallAuxCaller{contract: contract}, nil
}

// NewOpCallAuxTransactor creates a new write-only instance of OpCallAux, bound to a specific deployed contract.
func NewOpCallAuxTransactor(address common.Address, transactor bind.ContractTransactor) (*OpCallAuxTransactor, error) {
	contract, err := bindOpCallAux(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OpCallAuxTransactor{contract: contract}, nil
}

// NewOpCallAuxFilterer creates a new log filterer instance of OpCallAux, bound to a specific deployed contract.
func NewOpCallAuxFilterer(address common.Address, filterer bind.ContractFilterer) (*OpCallAuxFilterer, error) {
	contract, err := bindOpCallAux(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OpCallAuxFilterer{contract: contract}, nil
}

// bindOpCallAux binds a generic wrapper to an already deployed contract.
func bindOpCallAux(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OpCallAuxMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OpCallAux *OpCallAuxRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OpCallAux.Contract.OpCallAuxCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OpCallAux *OpCallAuxRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpCallAux.Contract.OpCallAuxTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OpCallAux *OpCallAuxRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OpCallAux.Contract.OpCallAuxTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OpCallAux *OpCallAuxCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OpCallAux.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OpCallAux *OpCallAuxTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpCallAux.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OpCallAux *OpCallAuxTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OpCallAux.Contract.contract.Transact(opts, method, params...)
}

// AddTwo is a paid mutator transaction binding the contract method 0x0f52d66e.
//
// Solidity: function addTwo(uint256 a, uint256 b) returns(uint256)
func (_OpCallAux *OpCallAuxTransactor) AddTwo(opts *bind.TransactOpts, a *big.Int, b *big.Int) (*types.Transaction, error) {
	return _OpCallAux.contract.Transact(opts, "addTwo", a, b)
}

// AddTwo is a paid mutator transaction binding the contract method 0x0f52d66e.
//
// Solidity: function addTwo(uint256 a, uint256 b) returns(uint256)
func (_OpCallAux *OpCallAuxSession) AddTwo(a *big.Int, b *big.Int) (*types.Transaction, error) {
	return _OpCallAux.Contract.AddTwo(&_OpCallAux.TransactOpts, a, b)
}

// AddTwo is a paid mutator transaction binding the contract method 0x0f52d66e.
//
// Solidity: function addTwo(uint256 a, uint256 b) returns(uint256)
func (_OpCallAux *OpCallAuxTransactorSession) AddTwo(a *big.Int, b *big.Int) (*types.Transaction, error) {
	return _OpCallAux.Contract.AddTwo(&_OpCallAux.TransactOpts, a, b)
}

// AuxFail is a paid mutator transaction binding the contract method 0x266f6aec.
//
// Solidity: function auxFail() returns()
func (_OpCallAux *OpCallAuxTransactor) AuxFail(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpCallAux.contract.Transact(opts, "auxFail")
}

// AuxFail is a paid mutator transaction binding the contract method 0x266f6aec.
//
// Solidity: function auxFail() returns()
func (_OpCallAux *OpCallAuxSession) AuxFail() (*types.Transaction, error) {
	return _OpCallAux.Contract.AuxFail(&_OpCallAux.TransactOpts)
}

// AuxFail is a paid mutator transaction binding the contract method 0x266f6aec.
//
// Solidity: function auxFail() returns()
func (_OpCallAux *OpCallAuxTransactorSession) AuxFail() (*types.Transaction, error) {
	return _OpCallAux.Contract.AuxFail(&_OpCallAux.TransactOpts)
}

// AuxReturn is a paid mutator transaction binding the contract method 0x6aecbc33.
//
// Solidity: function auxReturn() returns(uint256)
func (_OpCallAux *OpCallAuxTransactor) AuxReturn(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpCallAux.contract.Transact(opts, "auxReturn")
}

// AuxReturn is a paid mutator transaction binding the contract method 0x6aecbc33.
//
// Solidity: function auxReturn() returns(uint256)
func (_OpCallAux *OpCallAuxSession) AuxReturn() (*types.Transaction, error) {
	return _OpCallAux.Contract.AuxReturn(&_OpCallAux.TransactOpts)
}

// AuxReturn is a paid mutator transaction binding the contract method 0x6aecbc33.
//
// Solidity: function auxReturn() returns(uint256)
func (_OpCallAux *OpCallAuxTransactorSession) AuxReturn() (*types.Transaction, error) {
	return _OpCallAux.Contract.AuxReturn(&_OpCallAux.TransactOpts)
}

// AuxStop is a paid mutator transaction binding the contract method 0x03d3e8a7.
//
// Solidity: function auxStop() returns()
func (_OpCallAux *OpCallAuxTransactor) AuxStop(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpCallAux.contract.Transact(opts, "auxStop")
}

// AuxStop is a paid mutator transaction binding the contract method 0x03d3e8a7.
//
// Solidity: function auxStop() returns()
func (_OpCallAux *OpCallAuxSession) AuxStop() (*types.Transaction, error) {
	return _OpCallAux.Contract.AuxStop(&_OpCallAux.TransactOpts)
}

// AuxStop is a paid mutator transaction binding the contract method 0x03d3e8a7.
//
// Solidity: function auxStop() returns()
func (_OpCallAux *OpCallAuxTransactorSession) AuxStop() (*types.Transaction, error) {
	return _OpCallAux.Contract.AuxStop(&_OpCallAux.TransactOpts)
}

// AuxUpdate is a paid mutator transaction binding the contract method 0x3182d9a9.
//
// Solidity: function auxUpdate() returns(uint256)
func (_OpCallAux *OpCallAuxTransactor) AuxUpdate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpCallAux.contract.Transact(opts, "auxUpdate")
}

// AuxUpdate is a paid mutator transaction binding the contract method 0x3182d9a9.
//
// Solidity: function auxUpdate() returns(uint256)
func (_OpCallAux *OpCallAuxSession) AuxUpdate() (*types.Transaction, error) {
	return _OpCallAux.Contract.AuxUpdate(&_OpCallAux.TransactOpts)
}

// AuxUpdate is a paid mutator transaction binding the contract method 0x3182d9a9.
//
// Solidity: function auxUpdate() returns(uint256)
func (_OpCallAux *OpCallAuxTransactorSession) AuxUpdate() (*types.Transaction, error) {
	return _OpCallAux.Contract.AuxUpdate(&_OpCallAux.TransactOpts)
}

// AuxUpdateValues is a paid mutator transaction binding the contract method 0xf80efde5.
//
// Solidity: function auxUpdateValues() payable returns(uint256)
func (_OpCallAux *OpCallAuxTransactor) AuxUpdateValues(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpCallAux.contract.Transact(opts, "auxUpdateValues")
}

// AuxUpdateValues is a paid mutator transaction binding the contract method 0xf80efde5.
//
// Solidity: function auxUpdateValues() payable returns(uint256)
func (_OpCallAux *OpCallAuxSession) AuxUpdateValues() (*types.Transaction, error) {
	return _OpCallAux.Contract.AuxUpdateValues(&_OpCallAux.TransactOpts)
}

// AuxUpdateValues is a paid mutator transaction binding the contract method 0xf80efde5.
//
// Solidity: function auxUpdateValues() payable returns(uint256)
func (_OpCallAux *OpCallAuxTransactorSession) AuxUpdateValues() (*types.Transaction, error) {
	return _OpCallAux.Contract.AuxUpdateValues(&_OpCallAux.TransactOpts)
}

// OpCallSelfBalance is a paid mutator transaction binding the contract method 0xb6f54b25.
//
// Solidity: function opCallSelfBalance() payable returns(uint256)
func (_OpCallAux *OpCallAuxTransactor) OpCallSelfBalance(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpCallAux.contract.Transact(opts, "opCallSelfBalance")
}

// OpCallSelfBalance is a paid mutator transaction binding the contract method 0xb6f54b25.
//
// Solidity: function opCallSelfBalance() payable returns(uint256)
func (_OpCallAux *OpCallAuxSession) OpCallSelfBalance() (*types.Transaction, error) {
	return _OpCallAux.Contract.OpCallSelfBalance(&_OpCallAux.TransactOpts)
}

// OpCallSelfBalance is a paid mutator transaction binding the contract method 0xb6f54b25.
//
// Solidity: function opCallSelfBalance() payable returns(uint256)
func (_OpCallAux *OpCallAuxTransactorSession) OpCallSelfBalance() (*types.Transaction, error) {
	return _OpCallAux.Contract.OpCallSelfBalance(&_OpCallAux.TransactOpts)
}

// OpDelegateSelfBalance is a paid mutator transaction binding the contract method 0x03cfd846.
//
// Solidity: function opDelegateSelfBalance() payable returns(uint256)
func (_OpCallAux *OpCallAuxTransactor) OpDelegateSelfBalance(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpCallAux.contract.Transact(opts, "opDelegateSelfBalance")
}

// OpDelegateSelfBalance is a paid mutator transaction binding the contract method 0x03cfd846.
//
// Solidity: function opDelegateSelfBalance() payable returns(uint256)
func (_OpCallAux *OpCallAuxSession) OpDelegateSelfBalance() (*types.Transaction, error) {
	return _OpCallAux.Contract.OpDelegateSelfBalance(&_OpCallAux.TransactOpts)
}

// OpDelegateSelfBalance is a paid mutator transaction binding the contract method 0x03cfd846.
//
// Solidity: function opDelegateSelfBalance() payable returns(uint256)
func (_OpCallAux *OpCallAuxTransactorSession) OpDelegateSelfBalance() (*types.Transaction, error) {
	return _OpCallAux.Contract.OpDelegateSelfBalance(&_OpCallAux.TransactOpts)
}

// OpReturnCallSelfBalance is a paid mutator transaction binding the contract method 0x8b00fb6a.
//
// Solidity: function opReturnCallSelfBalance(address addrCall) payable returns(uint256)
func (_OpCallAux *OpCallAuxTransactor) OpReturnCallSelfBalance(opts *bind.TransactOpts, addrCall common.Address) (*types.Transaction, error) {
	return _OpCallAux.contract.Transact(opts, "opReturnCallSelfBalance", addrCall)
}

// OpReturnCallSelfBalance is a paid mutator transaction binding the contract method 0x8b00fb6a.
//
// Solidity: function opReturnCallSelfBalance(address addrCall) payable returns(uint256)
func (_OpCallAux *OpCallAuxSession) OpReturnCallSelfBalance(addrCall common.Address) (*types.Transaction, error) {
	return _OpCallAux.Contract.OpReturnCallSelfBalance(&_OpCallAux.TransactOpts, addrCall)
}

// OpReturnCallSelfBalance is a paid mutator transaction binding the contract method 0x8b00fb6a.
//
// Solidity: function opReturnCallSelfBalance(address addrCall) payable returns(uint256)
func (_OpCallAux *OpCallAuxTransactorSession) OpReturnCallSelfBalance(addrCall common.Address) (*types.Transaction, error) {
	return _OpCallAux.Contract.OpReturnCallSelfBalance(&_OpCallAux.TransactOpts, addrCall)
}
