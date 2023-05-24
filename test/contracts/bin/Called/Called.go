// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Called

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

// CalledMetaData contains all meta data concerning the Called contract.
var CalledMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getVars\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_num\",\"type\":\"uint256\"}],\"name\":\"getVarsAndVariable\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_num\",\"type\":\"uint256\"}],\"name\":\"setVars\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_num\",\"type\":\"uint256\"}],\"name\":\"setVarsViaCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061027f806100206000396000f3fe60806040526004361061003f5760003560e01c8063175df23c146100445780636466414b146100a6578063813d8a37146100d1578063f4bdc7341461010a575b600080fd5b34801561005057600080fd5b5061007961005f3660046101f5565b60005460015460025491936001600160a01b039091169290565b604080519485526001600160a01b0390931660208501529183015260608201526080015b60405180910390f35b6100cf6100b43660046101f5565b600055600180546001600160a01b0319163317905534600255565b005b3480156100dd57600080fd5b50600054600154600254604080519384526001600160a01b0390921660208401529082015260600161009d565b6100cf6101183660046101f5565b60405160248101829052600090309060440160408051601f198184030181529181526020820180516001600160e01b0316636466414b60e01b1790525161015f919061020e565b6000604051808303816000865af19150503d806000811461019c576040519150601f19603f3d011682016040523d82523d6000602084013e6101a1565b606091505b505080915050806101f15760405162461bcd60e51b815260206004820152601660248201527519985a5b1959081d1bc81c195c999bdc9b4818d85b1b60521b604482015260640160405180910390fd5b5050565b60006020828403121561020757600080fd5b5035919050565b6000825160005b8181101561022f5760208186018101518583015201610215565b8181111561023e576000828501525b50919091019291505056fea26469706673582212209d0d58578c55d453a185e4abef4bae0e18e22fd1b4c4e3d2072d42cbf3ccd4d864736f6c634300080c0033",
}

// CalledABI is the input ABI used to generate the binding from.
// Deprecated: Use CalledMetaData.ABI instead.
var CalledABI = CalledMetaData.ABI

// CalledBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CalledMetaData.Bin instead.
var CalledBin = CalledMetaData.Bin

// DeployCalled deploys a new Ethereum contract, binding an instance of Called to it.
func DeployCalled(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Called, error) {
	parsed, err := CalledMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CalledBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Called{CalledCaller: CalledCaller{contract: contract}, CalledTransactor: CalledTransactor{contract: contract}, CalledFilterer: CalledFilterer{contract: contract}}, nil
}

// Called is an auto generated Go binding around an Ethereum contract.
type Called struct {
	CalledCaller     // Read-only binding to the contract
	CalledTransactor // Write-only binding to the contract
	CalledFilterer   // Log filterer for contract events
}

// CalledCaller is an auto generated read-only Go binding around an Ethereum contract.
type CalledCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CalledTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CalledTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CalledFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CalledFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CalledSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CalledSession struct {
	Contract     *Called           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CalledCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CalledCallerSession struct {
	Contract *CalledCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// CalledTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CalledTransactorSession struct {
	Contract     *CalledTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CalledRaw is an auto generated low-level Go binding around an Ethereum contract.
type CalledRaw struct {
	Contract *Called // Generic contract binding to access the raw methods on
}

// CalledCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CalledCallerRaw struct {
	Contract *CalledCaller // Generic read-only contract binding to access the raw methods on
}

// CalledTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CalledTransactorRaw struct {
	Contract *CalledTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCalled creates a new instance of Called, bound to a specific deployed contract.
func NewCalled(address common.Address, backend bind.ContractBackend) (*Called, error) {
	contract, err := bindCalled(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Called{CalledCaller: CalledCaller{contract: contract}, CalledTransactor: CalledTransactor{contract: contract}, CalledFilterer: CalledFilterer{contract: contract}}, nil
}

// NewCalledCaller creates a new read-only instance of Called, bound to a specific deployed contract.
func NewCalledCaller(address common.Address, caller bind.ContractCaller) (*CalledCaller, error) {
	contract, err := bindCalled(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CalledCaller{contract: contract}, nil
}

// NewCalledTransactor creates a new write-only instance of Called, bound to a specific deployed contract.
func NewCalledTransactor(address common.Address, transactor bind.ContractTransactor) (*CalledTransactor, error) {
	contract, err := bindCalled(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CalledTransactor{contract: contract}, nil
}

// NewCalledFilterer creates a new log filterer instance of Called, bound to a specific deployed contract.
func NewCalledFilterer(address common.Address, filterer bind.ContractFilterer) (*CalledFilterer, error) {
	contract, err := bindCalled(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CalledFilterer{contract: contract}, nil
}

// bindCalled binds a generic wrapper to an already deployed contract.
func bindCalled(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CalledMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Called *CalledRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Called.Contract.CalledCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Called *CalledRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Called.Contract.CalledTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Called *CalledRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Called.Contract.CalledTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Called *CalledCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Called.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Called *CalledTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Called.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Called *CalledTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Called.Contract.contract.Transact(opts, method, params...)
}

// GetVars is a free data retrieval call binding the contract method 0x813d8a37.
//
// Solidity: function getVars() view returns(uint256, address, uint256)
func (_Called *CalledCaller) GetVars(opts *bind.CallOpts) (*big.Int, common.Address, *big.Int, error) {
	var out []interface{}
	err := _Called.contract.Call(opts, &out, "getVars")

	if err != nil {
		return *new(*big.Int), *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

// GetVars is a free data retrieval call binding the contract method 0x813d8a37.
//
// Solidity: function getVars() view returns(uint256, address, uint256)
func (_Called *CalledSession) GetVars() (*big.Int, common.Address, *big.Int, error) {
	return _Called.Contract.GetVars(&_Called.CallOpts)
}

// GetVars is a free data retrieval call binding the contract method 0x813d8a37.
//
// Solidity: function getVars() view returns(uint256, address, uint256)
func (_Called *CalledCallerSession) GetVars() (*big.Int, common.Address, *big.Int, error) {
	return _Called.Contract.GetVars(&_Called.CallOpts)
}

// GetVarsAndVariable is a free data retrieval call binding the contract method 0x175df23c.
//
// Solidity: function getVarsAndVariable(uint256 _num) view returns(uint256, address, uint256, uint256)
func (_Called *CalledCaller) GetVarsAndVariable(opts *bind.CallOpts, _num *big.Int) (*big.Int, common.Address, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _Called.contract.Call(opts, &out, "getVarsAndVariable", _num)

	if err != nil {
		return *new(*big.Int), *new(common.Address), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return out0, out1, out2, out3, err

}

// GetVarsAndVariable is a free data retrieval call binding the contract method 0x175df23c.
//
// Solidity: function getVarsAndVariable(uint256 _num) view returns(uint256, address, uint256, uint256)
func (_Called *CalledSession) GetVarsAndVariable(_num *big.Int) (*big.Int, common.Address, *big.Int, *big.Int, error) {
	return _Called.Contract.GetVarsAndVariable(&_Called.CallOpts, _num)
}

// GetVarsAndVariable is a free data retrieval call binding the contract method 0x175df23c.
//
// Solidity: function getVarsAndVariable(uint256 _num) view returns(uint256, address, uint256, uint256)
func (_Called *CalledCallerSession) GetVarsAndVariable(_num *big.Int) (*big.Int, common.Address, *big.Int, *big.Int, error) {
	return _Called.Contract.GetVarsAndVariable(&_Called.CallOpts, _num)
}

// SetVars is a paid mutator transaction binding the contract method 0x6466414b.
//
// Solidity: function setVars(uint256 _num) payable returns()
func (_Called *CalledTransactor) SetVars(opts *bind.TransactOpts, _num *big.Int) (*types.Transaction, error) {
	return _Called.contract.Transact(opts, "setVars", _num)
}

// SetVars is a paid mutator transaction binding the contract method 0x6466414b.
//
// Solidity: function setVars(uint256 _num) payable returns()
func (_Called *CalledSession) SetVars(_num *big.Int) (*types.Transaction, error) {
	return _Called.Contract.SetVars(&_Called.TransactOpts, _num)
}

// SetVars is a paid mutator transaction binding the contract method 0x6466414b.
//
// Solidity: function setVars(uint256 _num) payable returns()
func (_Called *CalledTransactorSession) SetVars(_num *big.Int) (*types.Transaction, error) {
	return _Called.Contract.SetVars(&_Called.TransactOpts, _num)
}

// SetVarsViaCall is a paid mutator transaction binding the contract method 0xf4bdc734.
//
// Solidity: function setVarsViaCall(uint256 _num) payable returns()
func (_Called *CalledTransactor) SetVarsViaCall(opts *bind.TransactOpts, _num *big.Int) (*types.Transaction, error) {
	return _Called.contract.Transact(opts, "setVarsViaCall", _num)
}

// SetVarsViaCall is a paid mutator transaction binding the contract method 0xf4bdc734.
//
// Solidity: function setVarsViaCall(uint256 _num) payable returns()
func (_Called *CalledSession) SetVarsViaCall(_num *big.Int) (*types.Transaction, error) {
	return _Called.Contract.SetVarsViaCall(&_Called.TransactOpts, _num)
}

// SetVarsViaCall is a paid mutator transaction binding the contract method 0xf4bdc734.
//
// Solidity: function setVarsViaCall(uint256 _num) payable returns()
func (_Called *CalledTransactorSession) SetVarsViaCall(_num *big.Int) (*types.Transaction, error) {
	return _Called.Contract.SetVarsViaCall(&_Called.TransactOpts, _num)
}
