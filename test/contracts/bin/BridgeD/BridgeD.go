// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package BridgeD

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

// BridgeDMetaData contains all meta data concerning the BridgeD contract.
var BridgeDMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"exec\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080604052348015600f57600080fd5b50607d8061001e6000396000f3fe60806040526004361060205760003560e01c8063c1c0e9c414602b57600080fd5b36602657005b600080fd5b6045600080546001600160a01b0319163317905534600155565b00fea2646970667358221220aaf44f2ad2fa35f506fc332df1f7e991eab5ff654446615796a1e5d3419bee9364736f6c634300080c0033",
}

// BridgeDABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeDMetaData.ABI instead.
var BridgeDABI = BridgeDMetaData.ABI

// BridgeDBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BridgeDMetaData.Bin instead.
var BridgeDBin = BridgeDMetaData.Bin

// DeployBridgeD deploys a new Ethereum contract, binding an instance of BridgeD to it.
func DeployBridgeD(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BridgeD, error) {
	parsed, err := BridgeDMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BridgeDBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BridgeD{BridgeDCaller: BridgeDCaller{contract: contract}, BridgeDTransactor: BridgeDTransactor{contract: contract}, BridgeDFilterer: BridgeDFilterer{contract: contract}}, nil
}

// BridgeD is an auto generated Go binding around an Ethereum contract.
type BridgeD struct {
	BridgeDCaller     // Read-only binding to the contract
	BridgeDTransactor // Write-only binding to the contract
	BridgeDFilterer   // Log filterer for contract events
}

// BridgeDCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeDCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeDTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeDTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeDFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeDFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeDSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeDSession struct {
	Contract     *BridgeD          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeDCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeDCallerSession struct {
	Contract *BridgeDCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// BridgeDTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeDTransactorSession struct {
	Contract     *BridgeDTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// BridgeDRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeDRaw struct {
	Contract *BridgeD // Generic contract binding to access the raw methods on
}

// BridgeDCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeDCallerRaw struct {
	Contract *BridgeDCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeDTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeDTransactorRaw struct {
	Contract *BridgeDTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeD creates a new instance of BridgeD, bound to a specific deployed contract.
func NewBridgeD(address common.Address, backend bind.ContractBackend) (*BridgeD, error) {
	contract, err := bindBridgeD(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeD{BridgeDCaller: BridgeDCaller{contract: contract}, BridgeDTransactor: BridgeDTransactor{contract: contract}, BridgeDFilterer: BridgeDFilterer{contract: contract}}, nil
}

// NewBridgeDCaller creates a new read-only instance of BridgeD, bound to a specific deployed contract.
func NewBridgeDCaller(address common.Address, caller bind.ContractCaller) (*BridgeDCaller, error) {
	contract, err := bindBridgeD(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeDCaller{contract: contract}, nil
}

// NewBridgeDTransactor creates a new write-only instance of BridgeD, bound to a specific deployed contract.
func NewBridgeDTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeDTransactor, error) {
	contract, err := bindBridgeD(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeDTransactor{contract: contract}, nil
}

// NewBridgeDFilterer creates a new log filterer instance of BridgeD, bound to a specific deployed contract.
func NewBridgeDFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeDFilterer, error) {
	contract, err := bindBridgeD(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeDFilterer{contract: contract}, nil
}

// bindBridgeD binds a generic wrapper to an already deployed contract.
func bindBridgeD(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BridgeDMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeD *BridgeDRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeD.Contract.BridgeDCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeD *BridgeDRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeD.Contract.BridgeDTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeD *BridgeDRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeD.Contract.BridgeDTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeD *BridgeDCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeD.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeD *BridgeDTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeD.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeD *BridgeDTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeD.Contract.contract.Transact(opts, method, params...)
}

// Exec is a paid mutator transaction binding the contract method 0xc1c0e9c4.
//
// Solidity: function exec() payable returns()
func (_BridgeD *BridgeDTransactor) Exec(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeD.contract.Transact(opts, "exec")
}

// Exec is a paid mutator transaction binding the contract method 0xc1c0e9c4.
//
// Solidity: function exec() payable returns()
func (_BridgeD *BridgeDSession) Exec() (*types.Transaction, error) {
	return _BridgeD.Contract.Exec(&_BridgeD.TransactOpts)
}

// Exec is a paid mutator transaction binding the contract method 0xc1c0e9c4.
//
// Solidity: function exec() payable returns()
func (_BridgeD *BridgeDTransactorSession) Exec() (*types.Transaction, error) {
	return _BridgeD.Contract.Exec(&_BridgeD.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeD *BridgeDTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeD.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeD *BridgeDSession) Receive() (*types.Transaction, error) {
	return _BridgeD.Contract.Receive(&_BridgeD.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeD *BridgeDTransactorSession) Receive() (*types.Transaction, error) {
	return _BridgeD.Contract.Receive(&_BridgeD.TransactOpts)
}
