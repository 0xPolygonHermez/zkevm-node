// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package BridgeA

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

// BridgeAMetaData contains all meta data concerning the BridgeA contract.
var BridgeAMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridgeB\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"bridgeC\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"bridgeD\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"}],\"name\":\"exec\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610238806100206000396000f3fe6080604052600436106100225760003560e01c8063f1dd7bae1461002e57600080fd5b3661002957005b600080fd5b61004161003c366004610173565b610043565b005b6040516001600160a01b038481166024830152838116604483015282811660648301526000919086169060840160408051601f198184030181529181526020820180516001600160e01b031663c023dcf360e01b179052516100a591906101c7565b600060405180830381855af49150503d80600081146100e0576040519150601f19603f3d011682016040523d82523d6000602084013e6100e5565b606091505b505080915050806101505760405162461bcd60e51b815260206004820152602b60248201527f6661696c656420746f20706572666f726d2064656c65676174652063616c6c2060448201526a3a3790313934b233b2902160a91b606482015260840160405180910390fd5b5050505050565b80356001600160a01b038116811461016e57600080fd5b919050565b6000806000806080858703121561018957600080fd5b61019285610157565b93506101a060208601610157565b92506101ae60408601610157565b91506101bc60608601610157565b905092959194509250565b6000825160005b818110156101e857602081860181015185830152016101ce565b818111156101f7576000828501525b50919091019291505056fea2646970667358221220a947a50743f17cd53a26643a22a9335f3f44dba656defb400ed3af1a7f61d38164736f6c634300080c0033",
}

// BridgeAABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeAMetaData.ABI instead.
var BridgeAABI = BridgeAMetaData.ABI

// BridgeABin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BridgeAMetaData.Bin instead.
var BridgeABin = BridgeAMetaData.Bin

// DeployBridgeA deploys a new Ethereum contract, binding an instance of BridgeA to it.
func DeployBridgeA(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BridgeA, error) {
	parsed, err := BridgeAMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BridgeABin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BridgeA{BridgeACaller: BridgeACaller{contract: contract}, BridgeATransactor: BridgeATransactor{contract: contract}, BridgeAFilterer: BridgeAFilterer{contract: contract}}, nil
}

// BridgeA is an auto generated Go binding around an Ethereum contract.
type BridgeA struct {
	BridgeACaller     // Read-only binding to the contract
	BridgeATransactor // Write-only binding to the contract
	BridgeAFilterer   // Log filterer for contract events
}

// BridgeACaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeACaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeATransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeATransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeAFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeAFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeASession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeASession struct {
	Contract     *BridgeA          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeACallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeACallerSession struct {
	Contract *BridgeACaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// BridgeATransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeATransactorSession struct {
	Contract     *BridgeATransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// BridgeARaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeARaw struct {
	Contract *BridgeA // Generic contract binding to access the raw methods on
}

// BridgeACallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeACallerRaw struct {
	Contract *BridgeACaller // Generic read-only contract binding to access the raw methods on
}

// BridgeATransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeATransactorRaw struct {
	Contract *BridgeATransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeA creates a new instance of BridgeA, bound to a specific deployed contract.
func NewBridgeA(address common.Address, backend bind.ContractBackend) (*BridgeA, error) {
	contract, err := bindBridgeA(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeA{BridgeACaller: BridgeACaller{contract: contract}, BridgeATransactor: BridgeATransactor{contract: contract}, BridgeAFilterer: BridgeAFilterer{contract: contract}}, nil
}

// NewBridgeACaller creates a new read-only instance of BridgeA, bound to a specific deployed contract.
func NewBridgeACaller(address common.Address, caller bind.ContractCaller) (*BridgeACaller, error) {
	contract, err := bindBridgeA(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeACaller{contract: contract}, nil
}

// NewBridgeATransactor creates a new write-only instance of BridgeA, bound to a specific deployed contract.
func NewBridgeATransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeATransactor, error) {
	contract, err := bindBridgeA(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeATransactor{contract: contract}, nil
}

// NewBridgeAFilterer creates a new log filterer instance of BridgeA, bound to a specific deployed contract.
func NewBridgeAFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeAFilterer, error) {
	contract, err := bindBridgeA(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeAFilterer{contract: contract}, nil
}

// bindBridgeA binds a generic wrapper to an already deployed contract.
func bindBridgeA(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BridgeAMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeA *BridgeARaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeA.Contract.BridgeACaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeA *BridgeARaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeA.Contract.BridgeATransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeA *BridgeARaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeA.Contract.BridgeATransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeA *BridgeACallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeA.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeA *BridgeATransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeA.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeA *BridgeATransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeA.Contract.contract.Transact(opts, method, params...)
}

// Exec is a paid mutator transaction binding the contract method 0xf1dd7bae.
//
// Solidity: function exec(address bridgeB, address bridgeC, address bridgeD, address acc) payable returns()
func (_BridgeA *BridgeATransactor) Exec(opts *bind.TransactOpts, bridgeB common.Address, bridgeC common.Address, bridgeD common.Address, acc common.Address) (*types.Transaction, error) {
	return _BridgeA.contract.Transact(opts, "exec", bridgeB, bridgeC, bridgeD, acc)
}

// Exec is a paid mutator transaction binding the contract method 0xf1dd7bae.
//
// Solidity: function exec(address bridgeB, address bridgeC, address bridgeD, address acc) payable returns()
func (_BridgeA *BridgeASession) Exec(bridgeB common.Address, bridgeC common.Address, bridgeD common.Address, acc common.Address) (*types.Transaction, error) {
	return _BridgeA.Contract.Exec(&_BridgeA.TransactOpts, bridgeB, bridgeC, bridgeD, acc)
}

// Exec is a paid mutator transaction binding the contract method 0xf1dd7bae.
//
// Solidity: function exec(address bridgeB, address bridgeC, address bridgeD, address acc) payable returns()
func (_BridgeA *BridgeATransactorSession) Exec(bridgeB common.Address, bridgeC common.Address, bridgeD common.Address, acc common.Address) (*types.Transaction, error) {
	return _BridgeA.Contract.Exec(&_BridgeA.TransactOpts, bridgeB, bridgeC, bridgeD, acc)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeA *BridgeATransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeA.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeA *BridgeASession) Receive() (*types.Transaction, error) {
	return _BridgeA.Contract.Receive(&_BridgeA.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeA *BridgeATransactorSession) Receive() (*types.Transaction, error) {
	return _BridgeA.Contract.Receive(&_BridgeA.TransactOpts)
}
