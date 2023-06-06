// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ChainCallLevel4

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

// ChainCallLevel4MetaData contains all meta data concerning the ChainCallLevel4 contract.
var ChainCallLevel4MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"exec\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"execRevert\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"t\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061017c806100206000396000f3fe6080604052600436106100385760003560e01c80636d4ce63c14610044578063a6ecd35614610080578063c1c0e9c41461008a57600080fd5b3661003f57005b600080fd5b34801561005057600080fd5b50604080518082018252600481526361686f7960e01b6020820152905161007791906100f1565b60405180910390f35b6100886100a5565b005b610088600080546001600160a01b0319163317905534600155565b60405162461bcd60e51b815260206004820181905260248201527f61686f792c20746869732074782077696c6c20616c7761797320726576657274604482015260640160405180910390fd5b600060208083528351808285015260005b8181101561011e57858101830151858201604001528201610102565b81811115610130576000604083870101525b50601f01601f191692909201604001939250505056fea2646970667358221220f63edbc9a42dfa09f0dfeea3f57d754eff36ac6f42381584ef027ed215f7e4f264736f6c634300080c0033",
}

// ChainCallLevel4ABI is the input ABI used to generate the binding from.
// Deprecated: Use ChainCallLevel4MetaData.ABI instead.
var ChainCallLevel4ABI = ChainCallLevel4MetaData.ABI

// ChainCallLevel4Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ChainCallLevel4MetaData.Bin instead.
var ChainCallLevel4Bin = ChainCallLevel4MetaData.Bin

// DeployChainCallLevel4 deploys a new Ethereum contract, binding an instance of ChainCallLevel4 to it.
func DeployChainCallLevel4(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ChainCallLevel4, error) {
	parsed, err := ChainCallLevel4MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChainCallLevel4Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChainCallLevel4{ChainCallLevel4Caller: ChainCallLevel4Caller{contract: contract}, ChainCallLevel4Transactor: ChainCallLevel4Transactor{contract: contract}, ChainCallLevel4Filterer: ChainCallLevel4Filterer{contract: contract}}, nil
}

// ChainCallLevel4 is an auto generated Go binding around an Ethereum contract.
type ChainCallLevel4 struct {
	ChainCallLevel4Caller     // Read-only binding to the contract
	ChainCallLevel4Transactor // Write-only binding to the contract
	ChainCallLevel4Filterer   // Log filterer for contract events
}

// ChainCallLevel4Caller is an auto generated read-only Go binding around an Ethereum contract.
type ChainCallLevel4Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainCallLevel4Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ChainCallLevel4Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainCallLevel4Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChainCallLevel4Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainCallLevel4Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChainCallLevel4Session struct {
	Contract     *ChainCallLevel4  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChainCallLevel4CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChainCallLevel4CallerSession struct {
	Contract *ChainCallLevel4Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ChainCallLevel4TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChainCallLevel4TransactorSession struct {
	Contract     *ChainCallLevel4Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ChainCallLevel4Raw is an auto generated low-level Go binding around an Ethereum contract.
type ChainCallLevel4Raw struct {
	Contract *ChainCallLevel4 // Generic contract binding to access the raw methods on
}

// ChainCallLevel4CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChainCallLevel4CallerRaw struct {
	Contract *ChainCallLevel4Caller // Generic read-only contract binding to access the raw methods on
}

// ChainCallLevel4TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChainCallLevel4TransactorRaw struct {
	Contract *ChainCallLevel4Transactor // Generic write-only contract binding to access the raw methods on
}

// NewChainCallLevel4 creates a new instance of ChainCallLevel4, bound to a specific deployed contract.
func NewChainCallLevel4(address common.Address, backend bind.ContractBackend) (*ChainCallLevel4, error) {
	contract, err := bindChainCallLevel4(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel4{ChainCallLevel4Caller: ChainCallLevel4Caller{contract: contract}, ChainCallLevel4Transactor: ChainCallLevel4Transactor{contract: contract}, ChainCallLevel4Filterer: ChainCallLevel4Filterer{contract: contract}}, nil
}

// NewChainCallLevel4Caller creates a new read-only instance of ChainCallLevel4, bound to a specific deployed contract.
func NewChainCallLevel4Caller(address common.Address, caller bind.ContractCaller) (*ChainCallLevel4Caller, error) {
	contract, err := bindChainCallLevel4(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel4Caller{contract: contract}, nil
}

// NewChainCallLevel4Transactor creates a new write-only instance of ChainCallLevel4, bound to a specific deployed contract.
func NewChainCallLevel4Transactor(address common.Address, transactor bind.ContractTransactor) (*ChainCallLevel4Transactor, error) {
	contract, err := bindChainCallLevel4(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel4Transactor{contract: contract}, nil
}

// NewChainCallLevel4Filterer creates a new log filterer instance of ChainCallLevel4, bound to a specific deployed contract.
func NewChainCallLevel4Filterer(address common.Address, filterer bind.ContractFilterer) (*ChainCallLevel4Filterer, error) {
	contract, err := bindChainCallLevel4(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel4Filterer{contract: contract}, nil
}

// bindChainCallLevel4 binds a generic wrapper to an already deployed contract.
func bindChainCallLevel4(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChainCallLevel4MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainCallLevel4 *ChainCallLevel4Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainCallLevel4.Contract.ChainCallLevel4Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainCallLevel4 *ChainCallLevel4Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainCallLevel4.Contract.ChainCallLevel4Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainCallLevel4 *ChainCallLevel4Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainCallLevel4.Contract.ChainCallLevel4Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainCallLevel4 *ChainCallLevel4CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainCallLevel4.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainCallLevel4 *ChainCallLevel4TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainCallLevel4.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainCallLevel4 *ChainCallLevel4TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainCallLevel4.Contract.contract.Transact(opts, method, params...)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() pure returns(string t)
func (_ChainCallLevel4 *ChainCallLevel4Caller) Get(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ChainCallLevel4.contract.Call(opts, &out, "get")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() pure returns(string t)
func (_ChainCallLevel4 *ChainCallLevel4Session) Get() (string, error) {
	return _ChainCallLevel4.Contract.Get(&_ChainCallLevel4.CallOpts)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() pure returns(string t)
func (_ChainCallLevel4 *ChainCallLevel4CallerSession) Get() (string, error) {
	return _ChainCallLevel4.Contract.Get(&_ChainCallLevel4.CallOpts)
}

// Exec is a paid mutator transaction binding the contract method 0xc1c0e9c4.
//
// Solidity: function exec() payable returns()
func (_ChainCallLevel4 *ChainCallLevel4Transactor) Exec(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainCallLevel4.contract.Transact(opts, "exec")
}

// Exec is a paid mutator transaction binding the contract method 0xc1c0e9c4.
//
// Solidity: function exec() payable returns()
func (_ChainCallLevel4 *ChainCallLevel4Session) Exec() (*types.Transaction, error) {
	return _ChainCallLevel4.Contract.Exec(&_ChainCallLevel4.TransactOpts)
}

// Exec is a paid mutator transaction binding the contract method 0xc1c0e9c4.
//
// Solidity: function exec() payable returns()
func (_ChainCallLevel4 *ChainCallLevel4TransactorSession) Exec() (*types.Transaction, error) {
	return _ChainCallLevel4.Contract.Exec(&_ChainCallLevel4.TransactOpts)
}

// ExecRevert is a paid mutator transaction binding the contract method 0xa6ecd356.
//
// Solidity: function execRevert() payable returns()
func (_ChainCallLevel4 *ChainCallLevel4Transactor) ExecRevert(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainCallLevel4.contract.Transact(opts, "execRevert")
}

// ExecRevert is a paid mutator transaction binding the contract method 0xa6ecd356.
//
// Solidity: function execRevert() payable returns()
func (_ChainCallLevel4 *ChainCallLevel4Session) ExecRevert() (*types.Transaction, error) {
	return _ChainCallLevel4.Contract.ExecRevert(&_ChainCallLevel4.TransactOpts)
}

// ExecRevert is a paid mutator transaction binding the contract method 0xa6ecd356.
//
// Solidity: function execRevert() payable returns()
func (_ChainCallLevel4 *ChainCallLevel4TransactorSession) ExecRevert() (*types.Transaction, error) {
	return _ChainCallLevel4.Contract.ExecRevert(&_ChainCallLevel4.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ChainCallLevel4 *ChainCallLevel4Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainCallLevel4.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ChainCallLevel4 *ChainCallLevel4Session) Receive() (*types.Transaction, error) {
	return _ChainCallLevel4.Contract.Receive(&_ChainCallLevel4.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ChainCallLevel4 *ChainCallLevel4TransactorSession) Receive() (*types.Transaction, error) {
	return _ChainCallLevel4.Contract.Receive(&_ChainCallLevel4.TransactOpts)
}
