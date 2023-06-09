// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package BridgeC

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

// BridgeCMetaData contains all meta data concerning the BridgeC contract.
var BridgeCMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridgeD\",\"type\":\"address\"}],\"name\":\"exec\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b506101d5806100206000396000f3fe6080604052600436106100225760003560e01c80636bb6126e1461002e57600080fd5b3661002957005b600080fd5b61004161003c366004610134565b610043565b005b60408051600481526024810182526020810180516001600160e01b03166330703a7160e21b17905290516000916001600160a01b038416916100859190610164565b600060405180830381855af49150503d80600081146100c0576040519150601f19603f3d011682016040523d82523d6000602084013e6100c5565b606091505b505080915050806101305760405162461bcd60e51b815260206004820152602b60248201527f6661696c656420746f20706572666f726d2064656c65676174652063616c6c2060448201526a1d1bc8189c9a5919d9481160aa1b606482015260840160405180910390fd5b5050565b60006020828403121561014657600080fd5b81356001600160a01b038116811461015d57600080fd5b9392505050565b6000825160005b81811015610185576020818601810151858301520161016b565b81811115610194576000828501525b50919091019291505056fea2646970667358221220bee297c7d22cf5dab83a018836e0cb0a7080854314a53ff53cefd47614053c3e64736f6c634300080c0033",
}

// BridgeCABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeCMetaData.ABI instead.
var BridgeCABI = BridgeCMetaData.ABI

// BridgeCBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BridgeCMetaData.Bin instead.
var BridgeCBin = BridgeCMetaData.Bin

// DeployBridgeC deploys a new Ethereum contract, binding an instance of BridgeC to it.
func DeployBridgeC(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BridgeC, error) {
	parsed, err := BridgeCMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BridgeCBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BridgeC{BridgeCCaller: BridgeCCaller{contract: contract}, BridgeCTransactor: BridgeCTransactor{contract: contract}, BridgeCFilterer: BridgeCFilterer{contract: contract}}, nil
}

// BridgeC is an auto generated Go binding around an Ethereum contract.
type BridgeC struct {
	BridgeCCaller     // Read-only binding to the contract
	BridgeCTransactor // Write-only binding to the contract
	BridgeCFilterer   // Log filterer for contract events
}

// BridgeCCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeCSession struct {
	Contract     *BridgeC          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeCCallerSession struct {
	Contract *BridgeCCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// BridgeCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeCTransactorSession struct {
	Contract     *BridgeCTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// BridgeCRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeCRaw struct {
	Contract *BridgeC // Generic contract binding to access the raw methods on
}

// BridgeCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeCCallerRaw struct {
	Contract *BridgeCCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeCTransactorRaw struct {
	Contract *BridgeCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeC creates a new instance of BridgeC, bound to a specific deployed contract.
func NewBridgeC(address common.Address, backend bind.ContractBackend) (*BridgeC, error) {
	contract, err := bindBridgeC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeC{BridgeCCaller: BridgeCCaller{contract: contract}, BridgeCTransactor: BridgeCTransactor{contract: contract}, BridgeCFilterer: BridgeCFilterer{contract: contract}}, nil
}

// NewBridgeCCaller creates a new read-only instance of BridgeC, bound to a specific deployed contract.
func NewBridgeCCaller(address common.Address, caller bind.ContractCaller) (*BridgeCCaller, error) {
	contract, err := bindBridgeC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeCCaller{contract: contract}, nil
}

// NewBridgeCTransactor creates a new write-only instance of BridgeC, bound to a specific deployed contract.
func NewBridgeCTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeCTransactor, error) {
	contract, err := bindBridgeC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeCTransactor{contract: contract}, nil
}

// NewBridgeCFilterer creates a new log filterer instance of BridgeC, bound to a specific deployed contract.
func NewBridgeCFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeCFilterer, error) {
	contract, err := bindBridgeC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeCFilterer{contract: contract}, nil
}

// bindBridgeC binds a generic wrapper to an already deployed contract.
func bindBridgeC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BridgeCMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeC *BridgeCRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeC.Contract.BridgeCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeC *BridgeCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeC.Contract.BridgeCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeC *BridgeCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeC.Contract.BridgeCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeC *BridgeCCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeC *BridgeCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeC *BridgeCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeC.Contract.contract.Transact(opts, method, params...)
}

// Exec is a paid mutator transaction binding the contract method 0x6bb6126e.
//
// Solidity: function exec(address bridgeD) payable returns()
func (_BridgeC *BridgeCTransactor) Exec(opts *bind.TransactOpts, bridgeD common.Address) (*types.Transaction, error) {
	return _BridgeC.contract.Transact(opts, "exec", bridgeD)
}

// Exec is a paid mutator transaction binding the contract method 0x6bb6126e.
//
// Solidity: function exec(address bridgeD) payable returns()
func (_BridgeC *BridgeCSession) Exec(bridgeD common.Address) (*types.Transaction, error) {
	return _BridgeC.Contract.Exec(&_BridgeC.TransactOpts, bridgeD)
}

// Exec is a paid mutator transaction binding the contract method 0x6bb6126e.
//
// Solidity: function exec(address bridgeD) payable returns()
func (_BridgeC *BridgeCTransactorSession) Exec(bridgeD common.Address) (*types.Transaction, error) {
	return _BridgeC.Contract.Exec(&_BridgeC.TransactOpts, bridgeD)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeC *BridgeCTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeC.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeC *BridgeCSession) Receive() (*types.Transaction, error) {
	return _BridgeC.Contract.Receive(&_BridgeC.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeC *BridgeCTransactorSession) Receive() (*types.Transaction, error) {
	return _BridgeC.Contract.Receive(&_BridgeC.TransactOpts)
}
