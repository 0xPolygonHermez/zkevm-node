// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package BridgeB

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

// BridgeBMetaData contains all meta data concerning the BridgeB contract.
var BridgeBMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridgeC\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"bridgeD\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"acc\",\"type\":\"address\"}],\"name\":\"exec\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b506102b1806100206000396000f3fe6080604052600436106100225760003560e01c8063c023dcf31461002e57600080fd5b3661002957005b600080fd5b61004161003c3660046101fd565b610043565b005b6040516001600160a01b0383811660248301526000919085169060440160408051601f198184030181529181526020820180516001600160e01b03166335db093760e11b179052516100959190610240565b6000604051808303816000865af19150503d80600081146100d2576040519150601f19603f3d011682016040523d82523d6000602084013e6100d7565b606091505b5050809150508061013a5760405162461bcd60e51b815260206004820152602260248201527f6661696c656420746f20706572666f726d2063616c6c20746f20627269646765604482015261204360f01b60648201526084015b60405180910390fd5b6040516001600160a01b038316903490600081818185875af1925050503d8060008114610183576040519150601f19603f3d011682016040523d82523d6000602084013e610188565b606091505b505080915050806101db5760405162461bcd60e51b815260206004820152601d60248201527f6661696c656420746f20706572666f726d2063616c6c20746f206163630000006044820152606401610131565b50505050565b80356001600160a01b03811681146101f857600080fd5b919050565b60008060006060848603121561021257600080fd5b61021b846101e1565b9250610229602085016101e1565b9150610237604085016101e1565b90509250925092565b6000825160005b818110156102615760208186018101518583015201610247565b81811115610270576000828501525b50919091019291505056fea2646970667358221220a6d0c41f2ec16a0138e910ef7f4d6468924a72bc507ef5ab01aa99a682f7172664736f6c634300080c0033",
}

// BridgeBABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeBMetaData.ABI instead.
var BridgeBABI = BridgeBMetaData.ABI

// BridgeBBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BridgeBMetaData.Bin instead.
var BridgeBBin = BridgeBMetaData.Bin

// DeployBridgeB deploys a new Ethereum contract, binding an instance of BridgeB to it.
func DeployBridgeB(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BridgeB, error) {
	parsed, err := BridgeBMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BridgeBBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BridgeB{BridgeBCaller: BridgeBCaller{contract: contract}, BridgeBTransactor: BridgeBTransactor{contract: contract}, BridgeBFilterer: BridgeBFilterer{contract: contract}}, nil
}

// BridgeB is an auto generated Go binding around an Ethereum contract.
type BridgeB struct {
	BridgeBCaller     // Read-only binding to the contract
	BridgeBTransactor // Write-only binding to the contract
	BridgeBFilterer   // Log filterer for contract events
}

// BridgeBCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeBCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeBTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeBTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeBFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeBFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeBSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeBSession struct {
	Contract     *BridgeB          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeBCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeBCallerSession struct {
	Contract *BridgeBCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// BridgeBTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeBTransactorSession struct {
	Contract     *BridgeBTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// BridgeBRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeBRaw struct {
	Contract *BridgeB // Generic contract binding to access the raw methods on
}

// BridgeBCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeBCallerRaw struct {
	Contract *BridgeBCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeBTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeBTransactorRaw struct {
	Contract *BridgeBTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeB creates a new instance of BridgeB, bound to a specific deployed contract.
func NewBridgeB(address common.Address, backend bind.ContractBackend) (*BridgeB, error) {
	contract, err := bindBridgeB(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeB{BridgeBCaller: BridgeBCaller{contract: contract}, BridgeBTransactor: BridgeBTransactor{contract: contract}, BridgeBFilterer: BridgeBFilterer{contract: contract}}, nil
}

// NewBridgeBCaller creates a new read-only instance of BridgeB, bound to a specific deployed contract.
func NewBridgeBCaller(address common.Address, caller bind.ContractCaller) (*BridgeBCaller, error) {
	contract, err := bindBridgeB(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeBCaller{contract: contract}, nil
}

// NewBridgeBTransactor creates a new write-only instance of BridgeB, bound to a specific deployed contract.
func NewBridgeBTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeBTransactor, error) {
	contract, err := bindBridgeB(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeBTransactor{contract: contract}, nil
}

// NewBridgeBFilterer creates a new log filterer instance of BridgeB, bound to a specific deployed contract.
func NewBridgeBFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeBFilterer, error) {
	contract, err := bindBridgeB(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeBFilterer{contract: contract}, nil
}

// bindBridgeB binds a generic wrapper to an already deployed contract.
func bindBridgeB(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BridgeBMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeB *BridgeBRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeB.Contract.BridgeBCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeB *BridgeBRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeB.Contract.BridgeBTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeB *BridgeBRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeB.Contract.BridgeBTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeB *BridgeBCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeB.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeB *BridgeBTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeB.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeB *BridgeBTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeB.Contract.contract.Transact(opts, method, params...)
}

// Exec is a paid mutator transaction binding the contract method 0xc023dcf3.
//
// Solidity: function exec(address bridgeC, address bridgeD, address acc) payable returns()
func (_BridgeB *BridgeBTransactor) Exec(opts *bind.TransactOpts, bridgeC common.Address, bridgeD common.Address, acc common.Address) (*types.Transaction, error) {
	return _BridgeB.contract.Transact(opts, "exec", bridgeC, bridgeD, acc)
}

// Exec is a paid mutator transaction binding the contract method 0xc023dcf3.
//
// Solidity: function exec(address bridgeC, address bridgeD, address acc) payable returns()
func (_BridgeB *BridgeBSession) Exec(bridgeC common.Address, bridgeD common.Address, acc common.Address) (*types.Transaction, error) {
	return _BridgeB.Contract.Exec(&_BridgeB.TransactOpts, bridgeC, bridgeD, acc)
}

// Exec is a paid mutator transaction binding the contract method 0xc023dcf3.
//
// Solidity: function exec(address bridgeC, address bridgeD, address acc) payable returns()
func (_BridgeB *BridgeBTransactorSession) Exec(bridgeC common.Address, bridgeD common.Address, acc common.Address) (*types.Transaction, error) {
	return _BridgeB.Contract.Exec(&_BridgeB.TransactOpts, bridgeC, bridgeD, acc)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeB *BridgeBTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeB.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeB *BridgeBSession) Receive() (*types.Transaction, error) {
	return _BridgeB.Contract.Receive(&_BridgeB.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeB *BridgeBTransactorSession) Receive() (*types.Transaction, error) {
	return _BridgeB.Contract.Receive(&_BridgeB.TransactOpts)
}
