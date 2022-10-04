// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ConditionalLoop

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

// ConditionalLoopMetaData contains all meta data concerning the ConditionalLoop contract.
var ConditionalLoopMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"times\",\"type\":\"uint256\"}],\"name\":\"ExecuteLoop\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610155806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80636fa13d2414610030575b600080fd5b61004361003e3660046100dd565b610055565b60405190815260200160405180910390f35b60008082116100aa5760405162461bcd60e51b815260206004820152601e60248201527f74696d6573206e65656420746f20626520626967676572207468616e20300000604482015260640160405180910390fd5b600060015b8381116100d657816100c0816100f6565b92505080806100ce906100f6565b9150506100af565b5092915050565b6000602082840312156100ef57600080fd5b5035919050565b600060001982141561011857634e487b7160e01b600052601160045260246000fd5b506001019056fea26469706673582212207ad76f1fb30df3c2a250e2140d470fadf98052ee74d6f3847f0a305aee22743c64736f6c634300080c0033",
}

// ConditionalLoopABI is the input ABI used to generate the binding from.
// Deprecated: Use ConditionalLoopMetaData.ABI instead.
var ConditionalLoopABI = ConditionalLoopMetaData.ABI

// ConditionalLoopBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConditionalLoopMetaData.Bin instead.
var ConditionalLoopBin = ConditionalLoopMetaData.Bin

// DeployConditionalLoop deploys a new Ethereum contract, binding an instance of ConditionalLoop to it.
func DeployConditionalLoop(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ConditionalLoop, error) {
	parsed, err := ConditionalLoopMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConditionalLoopBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConditionalLoop{ConditionalLoopCaller: ConditionalLoopCaller{contract: contract}, ConditionalLoopTransactor: ConditionalLoopTransactor{contract: contract}, ConditionalLoopFilterer: ConditionalLoopFilterer{contract: contract}}, nil
}

// ConditionalLoop is an auto generated Go binding around an Ethereum contract.
type ConditionalLoop struct {
	ConditionalLoopCaller     // Read-only binding to the contract
	ConditionalLoopTransactor // Write-only binding to the contract
	ConditionalLoopFilterer   // Log filterer for contract events
}

// ConditionalLoopCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConditionalLoopCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConditionalLoopTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConditionalLoopTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConditionalLoopFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConditionalLoopFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConditionalLoopSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConditionalLoopSession struct {
	Contract     *ConditionalLoop  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConditionalLoopCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConditionalLoopCallerSession struct {
	Contract *ConditionalLoopCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ConditionalLoopTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConditionalLoopTransactorSession struct {
	Contract     *ConditionalLoopTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ConditionalLoopRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConditionalLoopRaw struct {
	Contract *ConditionalLoop // Generic contract binding to access the raw methods on
}

// ConditionalLoopCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConditionalLoopCallerRaw struct {
	Contract *ConditionalLoopCaller // Generic read-only contract binding to access the raw methods on
}

// ConditionalLoopTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConditionalLoopTransactorRaw struct {
	Contract *ConditionalLoopTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConditionalLoop creates a new instance of ConditionalLoop, bound to a specific deployed contract.
func NewConditionalLoop(address common.Address, backend bind.ContractBackend) (*ConditionalLoop, error) {
	contract, err := bindConditionalLoop(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConditionalLoop{ConditionalLoopCaller: ConditionalLoopCaller{contract: contract}, ConditionalLoopTransactor: ConditionalLoopTransactor{contract: contract}, ConditionalLoopFilterer: ConditionalLoopFilterer{contract: contract}}, nil
}

// NewConditionalLoopCaller creates a new read-only instance of ConditionalLoop, bound to a specific deployed contract.
func NewConditionalLoopCaller(address common.Address, caller bind.ContractCaller) (*ConditionalLoopCaller, error) {
	contract, err := bindConditionalLoop(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConditionalLoopCaller{contract: contract}, nil
}

// NewConditionalLoopTransactor creates a new write-only instance of ConditionalLoop, bound to a specific deployed contract.
func NewConditionalLoopTransactor(address common.Address, transactor bind.ContractTransactor) (*ConditionalLoopTransactor, error) {
	contract, err := bindConditionalLoop(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConditionalLoopTransactor{contract: contract}, nil
}

// NewConditionalLoopFilterer creates a new log filterer instance of ConditionalLoop, bound to a specific deployed contract.
func NewConditionalLoopFilterer(address common.Address, filterer bind.ContractFilterer) (*ConditionalLoopFilterer, error) {
	contract, err := bindConditionalLoop(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConditionalLoopFilterer{contract: contract}, nil
}

// bindConditionalLoop binds a generic wrapper to an already deployed contract.
func bindConditionalLoop(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConditionalLoopABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConditionalLoop *ConditionalLoopRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConditionalLoop.Contract.ConditionalLoopCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConditionalLoop *ConditionalLoopRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConditionalLoop.Contract.ConditionalLoopTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConditionalLoop *ConditionalLoopRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConditionalLoop.Contract.ConditionalLoopTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConditionalLoop *ConditionalLoopCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConditionalLoop.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConditionalLoop *ConditionalLoopTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConditionalLoop.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConditionalLoop *ConditionalLoopTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConditionalLoop.Contract.contract.Transact(opts, method, params...)
}

// ExecuteLoop is a free data retrieval call binding the contract method 0x6fa13d24.
//
// Solidity: function ExecuteLoop(uint256 times) pure returns(uint256)
func (_ConditionalLoop *ConditionalLoopCaller) ExecuteLoop(opts *bind.CallOpts, times *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ConditionalLoop.contract.Call(opts, &out, "ExecuteLoop", times)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExecuteLoop is a free data retrieval call binding the contract method 0x6fa13d24.
//
// Solidity: function ExecuteLoop(uint256 times) pure returns(uint256)
func (_ConditionalLoop *ConditionalLoopSession) ExecuteLoop(times *big.Int) (*big.Int, error) {
	return _ConditionalLoop.Contract.ExecuteLoop(&_ConditionalLoop.CallOpts, times)
}

// ExecuteLoop is a free data retrieval call binding the contract method 0x6fa13d24.
//
// Solidity: function ExecuteLoop(uint256 times) pure returns(uint256)
func (_ConditionalLoop *ConditionalLoopCallerSession) ExecuteLoop(times *big.Int) (*big.Int, error) {
	return _ConditionalLoop.Contract.ExecuteLoop(&_ConditionalLoop.CallOpts, times)
}
