// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Interaction

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

// InteractionMetaData contains all meta data concerning the Interaction contract.
var InteractionMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_counter\",\"type\":\"address\"}],\"name\":\"setCounterAddr\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610183806100206000396000f3fe6080604052600436106100295760003560e01c8063a87d942c1461002e578063ec39b42914610055575b600080fd5b34801561003a57600080fd5b50610043610087565b60405190815260200160405180910390f35b610085610063366004610104565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b005b60008060009054906101000a90046001600160a01b03166001600160a01b03166306661abd6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156100db573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100ff9190610134565b905090565b60006020828403121561011657600080fd5b81356001600160a01b038116811461012d57600080fd5b9392505050565b60006020828403121561014657600080fd5b505191905056fea264697066735822122058d432243ff91dfaa2dc1ed0a7a18603bef86a1c9a5fb323673b24e70a6d472664736f6c634300080c0033",
}

// InteractionABI is the input ABI used to generate the binding from.
// Deprecated: Use InteractionMetaData.ABI instead.
var InteractionABI = InteractionMetaData.ABI

// InteractionBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use InteractionMetaData.Bin instead.
var InteractionBin = InteractionMetaData.Bin

// DeployInteraction deploys a new Ethereum contract, binding an instance of Interaction to it.
func DeployInteraction(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Interaction, error) {
	parsed, err := InteractionMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(InteractionBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Interaction{InteractionCaller: InteractionCaller{contract: contract}, InteractionTransactor: InteractionTransactor{contract: contract}, InteractionFilterer: InteractionFilterer{contract: contract}}, nil
}

// Interaction is an auto generated Go binding around an Ethereum contract.
type Interaction struct {
	InteractionCaller     // Read-only binding to the contract
	InteractionTransactor // Write-only binding to the contract
	InteractionFilterer   // Log filterer for contract events
}

// InteractionCaller is an auto generated read-only Go binding around an Ethereum contract.
type InteractionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InteractionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type InteractionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InteractionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type InteractionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InteractionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type InteractionSession struct {
	Contract     *Interaction      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// InteractionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type InteractionCallerSession struct {
	Contract *InteractionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// InteractionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type InteractionTransactorSession struct {
	Contract     *InteractionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// InteractionRaw is an auto generated low-level Go binding around an Ethereum contract.
type InteractionRaw struct {
	Contract *Interaction // Generic contract binding to access the raw methods on
}

// InteractionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type InteractionCallerRaw struct {
	Contract *InteractionCaller // Generic read-only contract binding to access the raw methods on
}

// InteractionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type InteractionTransactorRaw struct {
	Contract *InteractionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInteraction creates a new instance of Interaction, bound to a specific deployed contract.
func NewInteraction(address common.Address, backend bind.ContractBackend) (*Interaction, error) {
	contract, err := bindInteraction(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Interaction{InteractionCaller: InteractionCaller{contract: contract}, InteractionTransactor: InteractionTransactor{contract: contract}, InteractionFilterer: InteractionFilterer{contract: contract}}, nil
}

// NewInteractionCaller creates a new read-only instance of Interaction, bound to a specific deployed contract.
func NewInteractionCaller(address common.Address, caller bind.ContractCaller) (*InteractionCaller, error) {
	contract, err := bindInteraction(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InteractionCaller{contract: contract}, nil
}

// NewInteractionTransactor creates a new write-only instance of Interaction, bound to a specific deployed contract.
func NewInteractionTransactor(address common.Address, transactor bind.ContractTransactor) (*InteractionTransactor, error) {
	contract, err := bindInteraction(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InteractionTransactor{contract: contract}, nil
}

// NewInteractionFilterer creates a new log filterer instance of Interaction, bound to a specific deployed contract.
func NewInteractionFilterer(address common.Address, filterer bind.ContractFilterer) (*InteractionFilterer, error) {
	contract, err := bindInteraction(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InteractionFilterer{contract: contract}, nil
}

// bindInteraction binds a generic wrapper to an already deployed contract.
func bindInteraction(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := InteractionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Interaction *InteractionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Interaction.Contract.InteractionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Interaction *InteractionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Interaction.Contract.InteractionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Interaction *InteractionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Interaction.Contract.InteractionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Interaction *InteractionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Interaction.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Interaction *InteractionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Interaction.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Interaction *InteractionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Interaction.Contract.contract.Transact(opts, method, params...)
}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256)
func (_Interaction *InteractionCaller) GetCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Interaction.contract.Call(opts, &out, "getCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256)
func (_Interaction *InteractionSession) GetCount() (*big.Int, error) {
	return _Interaction.Contract.GetCount(&_Interaction.CallOpts)
}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256)
func (_Interaction *InteractionCallerSession) GetCount() (*big.Int, error) {
	return _Interaction.Contract.GetCount(&_Interaction.CallOpts)
}

// SetCounterAddr is a paid mutator transaction binding the contract method 0xec39b429.
//
// Solidity: function setCounterAddr(address _counter) payable returns()
func (_Interaction *InteractionTransactor) SetCounterAddr(opts *bind.TransactOpts, _counter common.Address) (*types.Transaction, error) {
	return _Interaction.contract.Transact(opts, "setCounterAddr", _counter)
}

// SetCounterAddr is a paid mutator transaction binding the contract method 0xec39b429.
//
// Solidity: function setCounterAddr(address _counter) payable returns()
func (_Interaction *InteractionSession) SetCounterAddr(_counter common.Address) (*types.Transaction, error) {
	return _Interaction.Contract.SetCounterAddr(&_Interaction.TransactOpts, _counter)
}

// SetCounterAddr is a paid mutator transaction binding the contract method 0xec39b429.
//
// Solidity: function setCounterAddr(address _counter) payable returns()
func (_Interaction *InteractionTransactorSession) SetCounterAddr(_counter common.Address) (*types.Transaction, error) {
	return _Interaction.Contract.SetCounterAddr(&_Interaction.TransactOpts, _counter)
}
