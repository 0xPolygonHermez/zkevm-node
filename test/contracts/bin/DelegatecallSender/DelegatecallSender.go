// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package DelegatecallSender

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

// DelegatecallSenderMetaData contains all meta data concerning the DelegatecallSender contract.
var DelegatecallSenderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"expectedSender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610359806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80631754dba51461003b578063f55332ab1461006a575b600080fd5b60005461004e906001600160a01b031681565b6040516001600160a01b03909116815260200160405180910390f35b61007d6100783660046101b1565b61007f565b005b600080546001600160a01b0319163317815560408051600481526024810182526020810180516001600160e01b03166329975a7560e21b179052905182916001600160a01b038516916100d29190610211565b600060405180830381855af49150503d806000811461010d576040519150601f19603f3d011682016040523d82523d6000602084013e610112565b606091505b5091509150816101675760448151101561012b57600080fd5b600481019050808060200190518101906101459190610243565b60405162461bcd60e51b815260040161015e91906102f0565b60405180910390fd5b816101ac5760405162461bcd60e51b815260206004820152601560248201527419195b1959d85d19590818d85b1b0819985a5b1959605a1b604482015260640161015e565b505050565b6000602082840312156101c357600080fd5b81356001600160a01b03811681146101da57600080fd5b9392505050565b60005b838110156101fc5781810151838201526020016101e4565b8381111561020b576000848401525b50505050565b600082516102238184602087016101e1565b9190910192915050565b634e487b7160e01b600052604160045260246000fd5b60006020828403121561025557600080fd5b815167ffffffffffffffff8082111561026d57600080fd5b818401915084601f83011261028157600080fd5b8151818111156102935761029361022d565b604051601f8201601f19908116603f011681019083821181831017156102bb576102bb61022d565b816040528281528760208487010111156102d457600080fd5b6102e58360208301602088016101e1565b979650505050505050565b602081526000825180602084015261030f8160408501602087016101e1565b601f01601f1916919091016040019291505056fea26469706673582212204f0f7adb219400c7c8870f1d43741cb5a4764fba74d7a50a050044ddda02166564736f6c634300080c0033",
}

// DelegatecallSenderABI is the input ABI used to generate the binding from.
// Deprecated: Use DelegatecallSenderMetaData.ABI instead.
var DelegatecallSenderABI = DelegatecallSenderMetaData.ABI

// DelegatecallSenderBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DelegatecallSenderMetaData.Bin instead.
var DelegatecallSenderBin = DelegatecallSenderMetaData.Bin

// DeployDelegatecallSender deploys a new Ethereum contract, binding an instance of DelegatecallSender to it.
func DeployDelegatecallSender(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DelegatecallSender, error) {
	parsed, err := DelegatecallSenderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DelegatecallSenderBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DelegatecallSender{DelegatecallSenderCaller: DelegatecallSenderCaller{contract: contract}, DelegatecallSenderTransactor: DelegatecallSenderTransactor{contract: contract}, DelegatecallSenderFilterer: DelegatecallSenderFilterer{contract: contract}}, nil
}

// DelegatecallSender is an auto generated Go binding around an Ethereum contract.
type DelegatecallSender struct {
	DelegatecallSenderCaller     // Read-only binding to the contract
	DelegatecallSenderTransactor // Write-only binding to the contract
	DelegatecallSenderFilterer   // Log filterer for contract events
}

// DelegatecallSenderCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegatecallSenderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatecallSenderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegatecallSenderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatecallSenderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegatecallSenderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatecallSenderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegatecallSenderSession struct {
	Contract     *DelegatecallSender // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DelegatecallSenderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegatecallSenderCallerSession struct {
	Contract *DelegatecallSenderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// DelegatecallSenderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegatecallSenderTransactorSession struct {
	Contract     *DelegatecallSenderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// DelegatecallSenderRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegatecallSenderRaw struct {
	Contract *DelegatecallSender // Generic contract binding to access the raw methods on
}

// DelegatecallSenderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegatecallSenderCallerRaw struct {
	Contract *DelegatecallSenderCaller // Generic read-only contract binding to access the raw methods on
}

// DelegatecallSenderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegatecallSenderTransactorRaw struct {
	Contract *DelegatecallSenderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegatecallSender creates a new instance of DelegatecallSender, bound to a specific deployed contract.
func NewDelegatecallSender(address common.Address, backend bind.ContractBackend) (*DelegatecallSender, error) {
	contract, err := bindDelegatecallSender(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DelegatecallSender{DelegatecallSenderCaller: DelegatecallSenderCaller{contract: contract}, DelegatecallSenderTransactor: DelegatecallSenderTransactor{contract: contract}, DelegatecallSenderFilterer: DelegatecallSenderFilterer{contract: contract}}, nil
}

// NewDelegatecallSenderCaller creates a new read-only instance of DelegatecallSender, bound to a specific deployed contract.
func NewDelegatecallSenderCaller(address common.Address, caller bind.ContractCaller) (*DelegatecallSenderCaller, error) {
	contract, err := bindDelegatecallSender(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegatecallSenderCaller{contract: contract}, nil
}

// NewDelegatecallSenderTransactor creates a new write-only instance of DelegatecallSender, bound to a specific deployed contract.
func NewDelegatecallSenderTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegatecallSenderTransactor, error) {
	contract, err := bindDelegatecallSender(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegatecallSenderTransactor{contract: contract}, nil
}

// NewDelegatecallSenderFilterer creates a new log filterer instance of DelegatecallSender, bound to a specific deployed contract.
func NewDelegatecallSenderFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegatecallSenderFilterer, error) {
	contract, err := bindDelegatecallSender(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegatecallSenderFilterer{contract: contract}, nil
}

// bindDelegatecallSender binds a generic wrapper to an already deployed contract.
func bindDelegatecallSender(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DelegatecallSenderABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelegatecallSender *DelegatecallSenderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DelegatecallSender.Contract.DelegatecallSenderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelegatecallSender *DelegatecallSenderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegatecallSender.Contract.DelegatecallSenderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelegatecallSender *DelegatecallSenderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelegatecallSender.Contract.DelegatecallSenderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelegatecallSender *DelegatecallSenderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DelegatecallSender.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelegatecallSender *DelegatecallSenderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegatecallSender.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelegatecallSender *DelegatecallSenderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelegatecallSender.Contract.contract.Transact(opts, method, params...)
}

// ExpectedSender is a free data retrieval call binding the contract method 0x1754dba5.
//
// Solidity: function expectedSender() view returns(address)
func (_DelegatecallSender *DelegatecallSenderCaller) ExpectedSender(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DelegatecallSender.contract.Call(opts, &out, "expectedSender")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ExpectedSender is a free data retrieval call binding the contract method 0x1754dba5.
//
// Solidity: function expectedSender() view returns(address)
func (_DelegatecallSender *DelegatecallSenderSession) ExpectedSender() (common.Address, error) {
	return _DelegatecallSender.Contract.ExpectedSender(&_DelegatecallSender.CallOpts)
}

// ExpectedSender is a free data retrieval call binding the contract method 0x1754dba5.
//
// Solidity: function expectedSender() view returns(address)
func (_DelegatecallSender *DelegatecallSenderCallerSession) ExpectedSender() (common.Address, error) {
	return _DelegatecallSender.Contract.ExpectedSender(&_DelegatecallSender.CallOpts)
}

// Call is a paid mutator transaction binding the contract method 0xf55332ab.
//
// Solidity: function call(address target) returns()
func (_DelegatecallSender *DelegatecallSenderTransactor) Call(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _DelegatecallSender.contract.Transact(opts, "call", target)
}

// Call is a paid mutator transaction binding the contract method 0xf55332ab.
//
// Solidity: function call(address target) returns()
func (_DelegatecallSender *DelegatecallSenderSession) Call(target common.Address) (*types.Transaction, error) {
	return _DelegatecallSender.Contract.Call(&_DelegatecallSender.TransactOpts, target)
}

// Call is a paid mutator transaction binding the contract method 0xf55332ab.
//
// Solidity: function call(address target) returns()
func (_DelegatecallSender *DelegatecallSenderTransactorSession) Call(target common.Address) (*types.Transaction, error) {
	return _DelegatecallSender.Contract.Call(&_DelegatecallSender.TransactOpts, target)
}
