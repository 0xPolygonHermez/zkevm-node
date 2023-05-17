// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package DelegateCallCaller

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

// DelegateCallCallerMetaData contains all meta data concerning the DelegateCallCaller contract.
var DelegateCallCallerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"num\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_num\",\"type\":\"uint256\"}],\"name\":\"setVars\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"value\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610279806100206000396000f3fe60806040526004361061003f5760003560e01c80633fa4f245146100445780634e70b1dc1461006d57806367e404ce14610083578063d1e0f308146100bb575b600080fd5b34801561005057600080fd5b5061005a60025481565b6040519081526020015b60405180910390f35b34801561007957600080fd5b5061005a60005481565b34801561008f57600080fd5b506001546100a3906001600160a01b031681565b6040516001600160a01b039091168152602001610064565b6100ce6100c9366004610183565b6100dc565b6040516100649291906101eb565b60006060600080856001600160a01b03168560405160240161010091815260200190565b60408051601f198184030181529181526020820180516001600160e01b0316636466414b60e01b179052516101359190610227565b600060405180830381855af49150503d8060008114610170576040519150601f19603f3d011682016040523d82523d6000602084013e610175565b606091505b509097909650945050505050565b6000806040838503121561019657600080fd5b82356001600160a01b03811681146101ad57600080fd5b946020939093013593505050565b60005b838110156101d65781810151838201526020016101be565b838111156101e5576000848401525b50505050565b821515815260406020820152600082518060408401526102128160608501602087016101bb565b601f01601f1916919091016060019392505050565b600082516102398184602087016101bb565b919091019291505056fea264697066735822122095fca7f3f89f0c00aaa8d6a39ba474f7765cac44c0d51a4b3b1619f398926d8e64736f6c634300080c0033",
}

// DelegateCallCallerABI is the input ABI used to generate the binding from.
// Deprecated: Use DelegateCallCallerMetaData.ABI instead.
var DelegateCallCallerABI = DelegateCallCallerMetaData.ABI

// DelegateCallCallerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DelegateCallCallerMetaData.Bin instead.
var DelegateCallCallerBin = DelegateCallCallerMetaData.Bin

// DeployDelegateCallCaller deploys a new Ethereum contract, binding an instance of DelegateCallCaller to it.
func DeployDelegateCallCaller(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DelegateCallCaller, error) {
	parsed, err := DelegateCallCallerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DelegateCallCallerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DelegateCallCaller{DelegateCallCallerCaller: DelegateCallCallerCaller{contract: contract}, DelegateCallCallerTransactor: DelegateCallCallerTransactor{contract: contract}, DelegateCallCallerFilterer: DelegateCallCallerFilterer{contract: contract}}, nil
}

// DelegateCallCaller is an auto generated Go binding around an Ethereum contract.
type DelegateCallCaller struct {
	DelegateCallCallerCaller     // Read-only binding to the contract
	DelegateCallCallerTransactor // Write-only binding to the contract
	DelegateCallCallerFilterer   // Log filterer for contract events
}

// DelegateCallCallerCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegateCallCallerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegateCallCallerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegateCallCallerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegateCallCallerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegateCallCallerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegateCallCallerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegateCallCallerSession struct {
	Contract     *DelegateCallCaller // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DelegateCallCallerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegateCallCallerCallerSession struct {
	Contract *DelegateCallCallerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// DelegateCallCallerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegateCallCallerTransactorSession struct {
	Contract     *DelegateCallCallerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// DelegateCallCallerRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegateCallCallerRaw struct {
	Contract *DelegateCallCaller // Generic contract binding to access the raw methods on
}

// DelegateCallCallerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegateCallCallerCallerRaw struct {
	Contract *DelegateCallCallerCaller // Generic read-only contract binding to access the raw methods on
}

// DelegateCallCallerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegateCallCallerTransactorRaw struct {
	Contract *DelegateCallCallerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegateCallCaller creates a new instance of DelegateCallCaller, bound to a specific deployed contract.
func NewDelegateCallCaller(address common.Address, backend bind.ContractBackend) (*DelegateCallCaller, error) {
	contract, err := bindDelegateCallCaller(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DelegateCallCaller{DelegateCallCallerCaller: DelegateCallCallerCaller{contract: contract}, DelegateCallCallerTransactor: DelegateCallCallerTransactor{contract: contract}, DelegateCallCallerFilterer: DelegateCallCallerFilterer{contract: contract}}, nil
}

// NewDelegateCallCallerCaller creates a new read-only instance of DelegateCallCaller, bound to a specific deployed contract.
func NewDelegateCallCallerCaller(address common.Address, caller bind.ContractCaller) (*DelegateCallCallerCaller, error) {
	contract, err := bindDelegateCallCaller(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegateCallCallerCaller{contract: contract}, nil
}

// NewDelegateCallCallerTransactor creates a new write-only instance of DelegateCallCaller, bound to a specific deployed contract.
func NewDelegateCallCallerTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegateCallCallerTransactor, error) {
	contract, err := bindDelegateCallCaller(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegateCallCallerTransactor{contract: contract}, nil
}

// NewDelegateCallCallerFilterer creates a new log filterer instance of DelegateCallCaller, bound to a specific deployed contract.
func NewDelegateCallCallerFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegateCallCallerFilterer, error) {
	contract, err := bindDelegateCallCaller(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegateCallCallerFilterer{contract: contract}, nil
}

// bindDelegateCallCaller binds a generic wrapper to an already deployed contract.
func bindDelegateCallCaller(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DelegateCallCallerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelegateCallCaller *DelegateCallCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DelegateCallCaller.Contract.DelegateCallCallerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelegateCallCaller *DelegateCallCallerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegateCallCaller.Contract.DelegateCallCallerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelegateCallCaller *DelegateCallCallerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelegateCallCaller.Contract.DelegateCallCallerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelegateCallCaller *DelegateCallCallerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DelegateCallCaller.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelegateCallCaller *DelegateCallCallerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegateCallCaller.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelegateCallCaller *DelegateCallCallerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelegateCallCaller.Contract.contract.Transact(opts, method, params...)
}

// Num is a free data retrieval call binding the contract method 0x4e70b1dc.
//
// Solidity: function num() view returns(uint256)
func (_DelegateCallCaller *DelegateCallCallerCaller) Num(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DelegateCallCaller.contract.Call(opts, &out, "num")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Num is a free data retrieval call binding the contract method 0x4e70b1dc.
//
// Solidity: function num() view returns(uint256)
func (_DelegateCallCaller *DelegateCallCallerSession) Num() (*big.Int, error) {
	return _DelegateCallCaller.Contract.Num(&_DelegateCallCaller.CallOpts)
}

// Num is a free data retrieval call binding the contract method 0x4e70b1dc.
//
// Solidity: function num() view returns(uint256)
func (_DelegateCallCaller *DelegateCallCallerCallerSession) Num() (*big.Int, error) {
	return _DelegateCallCaller.Contract.Num(&_DelegateCallCaller.CallOpts)
}

// Sender is a free data retrieval call binding the contract method 0x67e404ce.
//
// Solidity: function sender() view returns(address)
func (_DelegateCallCaller *DelegateCallCallerCaller) Sender(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DelegateCallCaller.contract.Call(opts, &out, "sender")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Sender is a free data retrieval call binding the contract method 0x67e404ce.
//
// Solidity: function sender() view returns(address)
func (_DelegateCallCaller *DelegateCallCallerSession) Sender() (common.Address, error) {
	return _DelegateCallCaller.Contract.Sender(&_DelegateCallCaller.CallOpts)
}

// Sender is a free data retrieval call binding the contract method 0x67e404ce.
//
// Solidity: function sender() view returns(address)
func (_DelegateCallCaller *DelegateCallCallerCallerSession) Sender() (common.Address, error) {
	return _DelegateCallCaller.Contract.Sender(&_DelegateCallCaller.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(uint256)
func (_DelegateCallCaller *DelegateCallCallerCaller) Value(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DelegateCallCaller.contract.Call(opts, &out, "value")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(uint256)
func (_DelegateCallCaller *DelegateCallCallerSession) Value() (*big.Int, error) {
	return _DelegateCallCaller.Contract.Value(&_DelegateCallCaller.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(uint256)
func (_DelegateCallCaller *DelegateCallCallerCallerSession) Value() (*big.Int, error) {
	return _DelegateCallCaller.Contract.Value(&_DelegateCallCaller.CallOpts)
}

// SetVars is a paid mutator transaction binding the contract method 0xd1e0f308.
//
// Solidity: function setVars(address _contract, uint256 _num) payable returns(bool, bytes)
func (_DelegateCallCaller *DelegateCallCallerTransactor) SetVars(opts *bind.TransactOpts, _contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _DelegateCallCaller.contract.Transact(opts, "setVars", _contract, _num)
}

// SetVars is a paid mutator transaction binding the contract method 0xd1e0f308.
//
// Solidity: function setVars(address _contract, uint256 _num) payable returns(bool, bytes)
func (_DelegateCallCaller *DelegateCallCallerSession) SetVars(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _DelegateCallCaller.Contract.SetVars(&_DelegateCallCaller.TransactOpts, _contract, _num)
}

// SetVars is a paid mutator transaction binding the contract method 0xd1e0f308.
//
// Solidity: function setVars(address _contract, uint256 _num) payable returns(bool, bytes)
func (_DelegateCallCaller *DelegateCallCallerTransactorSession) SetVars(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _DelegateCallCaller.Contract.SetVars(&_DelegateCallCaller.TransactOpts, _contract, _num)
}
