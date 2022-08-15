// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Read

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

// ReadMetaData contains all meta data concerning the Read contract.
var ReadMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"externalRead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"externalReadWParams\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"publicRead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"publicReadWParams\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600160005534801561001557600080fd5b5060f3806100246000396000f3fe6080604052348015600f57600080fd5b506004361060465760003560e01c806310f3f91a14604b5780636e5298b2146061578063af33475714604b578063bfa044ed146061575b600080fd5b6000545b60405190815260200160405180910390f35b604f606c3660046080565b600081600054607a91906098565b92915050565b600060208284031215609157600080fd5b5035919050565b6000821982111560b857634e487b7160e01b600052601160045260246000fd5b50019056fea26469706673582212201c06ac59eab70c5210664e3c901557a235d2c7823a61abb1fab9fa73602dedc764736f6c634300080c0033",
}

// ReadABI is the input ABI used to generate the binding from.
// Deprecated: Use ReadMetaData.ABI instead.
var ReadABI = ReadMetaData.ABI

// ReadBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ReadMetaData.Bin instead.
var ReadBin = ReadMetaData.Bin

// DeployRead deploys a new Ethereum contract, binding an instance of Read to it.
func DeployRead(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Read, error) {
	parsed, err := ReadMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ReadBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Read{ReadCaller: ReadCaller{contract: contract}, ReadTransactor: ReadTransactor{contract: contract}, ReadFilterer: ReadFilterer{contract: contract}}, nil
}

// Read is an auto generated Go binding around an Ethereum contract.
type Read struct {
	ReadCaller     // Read-only binding to the contract
	ReadTransactor // Write-only binding to the contract
	ReadFilterer   // Log filterer for contract events
}

// ReadCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReadCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReadTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReadTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReadFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReadFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReadSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReadSession struct {
	Contract     *Read             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReadCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReadCallerSession struct {
	Contract *ReadCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ReadTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReadTransactorSession struct {
	Contract     *ReadTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReadRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReadRaw struct {
	Contract *Read // Generic contract binding to access the raw methods on
}

// ReadCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReadCallerRaw struct {
	Contract *ReadCaller // Generic read-only contract binding to access the raw methods on
}

// ReadTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReadTransactorRaw struct {
	Contract *ReadTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRead creates a new instance of Read, bound to a specific deployed contract.
func NewRead(address common.Address, backend bind.ContractBackend) (*Read, error) {
	contract, err := bindRead(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Read{ReadCaller: ReadCaller{contract: contract}, ReadTransactor: ReadTransactor{contract: contract}, ReadFilterer: ReadFilterer{contract: contract}}, nil
}

// NewReadCaller creates a new read-only instance of Read, bound to a specific deployed contract.
func NewReadCaller(address common.Address, caller bind.ContractCaller) (*ReadCaller, error) {
	contract, err := bindRead(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReadCaller{contract: contract}, nil
}

// NewReadTransactor creates a new write-only instance of Read, bound to a specific deployed contract.
func NewReadTransactor(address common.Address, transactor bind.ContractTransactor) (*ReadTransactor, error) {
	contract, err := bindRead(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReadTransactor{contract: contract}, nil
}

// NewReadFilterer creates a new log filterer instance of Read, bound to a specific deployed contract.
func NewReadFilterer(address common.Address, filterer bind.ContractFilterer) (*ReadFilterer, error) {
	contract, err := bindRead(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReadFilterer{contract: contract}, nil
}

// bindRead binds a generic wrapper to an already deployed contract.
func bindRead(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ReadABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Read *ReadRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Read.Contract.ReadCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Read *ReadRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Read.Contract.ReadTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Read *ReadRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Read.Contract.ReadTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Read *ReadCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Read.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Read *ReadTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Read.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Read *ReadTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Read.Contract.contract.Transact(opts, method, params...)
}

// ExternalRead is a free data retrieval call binding the contract method 0x10f3f91a.
//
// Solidity: function externalRead() view returns(uint256)
func (_Read *ReadCaller) ExternalRead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "externalRead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExternalRead is a free data retrieval call binding the contract method 0x10f3f91a.
//
// Solidity: function externalRead() view returns(uint256)
func (_Read *ReadSession) ExternalRead() (*big.Int, error) {
	return _Read.Contract.ExternalRead(&_Read.CallOpts)
}

// ExternalRead is a free data retrieval call binding the contract method 0x10f3f91a.
//
// Solidity: function externalRead() view returns(uint256)
func (_Read *ReadCallerSession) ExternalRead() (*big.Int, error) {
	return _Read.Contract.ExternalRead(&_Read.CallOpts)
}

// ExternalReadWParams is a free data retrieval call binding the contract method 0x6e5298b2.
//
// Solidity: function externalReadWParams(uint256 p) view returns(uint256)
func (_Read *ReadCaller) ExternalReadWParams(opts *bind.CallOpts, p *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "externalReadWParams", p)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExternalReadWParams is a free data retrieval call binding the contract method 0x6e5298b2.
//
// Solidity: function externalReadWParams(uint256 p) view returns(uint256)
func (_Read *ReadSession) ExternalReadWParams(p *big.Int) (*big.Int, error) {
	return _Read.Contract.ExternalReadWParams(&_Read.CallOpts, p)
}

// ExternalReadWParams is a free data retrieval call binding the contract method 0x6e5298b2.
//
// Solidity: function externalReadWParams(uint256 p) view returns(uint256)
func (_Read *ReadCallerSession) ExternalReadWParams(p *big.Int) (*big.Int, error) {
	return _Read.Contract.ExternalReadWParams(&_Read.CallOpts, p)
}

// PublicRead is a free data retrieval call binding the contract method 0xaf334757.
//
// Solidity: function publicRead() view returns(uint256)
func (_Read *ReadCaller) PublicRead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "publicRead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PublicRead is a free data retrieval call binding the contract method 0xaf334757.
//
// Solidity: function publicRead() view returns(uint256)
func (_Read *ReadSession) PublicRead() (*big.Int, error) {
	return _Read.Contract.PublicRead(&_Read.CallOpts)
}

// PublicRead is a free data retrieval call binding the contract method 0xaf334757.
//
// Solidity: function publicRead() view returns(uint256)
func (_Read *ReadCallerSession) PublicRead() (*big.Int, error) {
	return _Read.Contract.PublicRead(&_Read.CallOpts)
}

// PublicReadWParams is a free data retrieval call binding the contract method 0xbfa044ed.
//
// Solidity: function publicReadWParams(uint256 p) view returns(uint256)
func (_Read *ReadCaller) PublicReadWParams(opts *bind.CallOpts, p *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "publicReadWParams", p)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PublicReadWParams is a free data retrieval call binding the contract method 0xbfa044ed.
//
// Solidity: function publicReadWParams(uint256 p) view returns(uint256)
func (_Read *ReadSession) PublicReadWParams(p *big.Int) (*big.Int, error) {
	return _Read.Contract.PublicReadWParams(&_Read.CallOpts, p)
}

// PublicReadWParams is a free data retrieval call binding the contract method 0xbfa044ed.
//
// Solidity: function publicReadWParams(uint256 p) view returns(uint256)
func (_Read *ReadCallerSession) PublicReadWParams(p *big.Int) (*big.Int, error) {
	return _Read.Contract.PublicReadWParams(&_Read.CallOpts, p)
}
