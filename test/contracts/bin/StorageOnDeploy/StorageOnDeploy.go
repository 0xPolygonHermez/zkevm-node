// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package StorageOnDeploy

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

// StorageOnDeployMetaData contains all meta data concerning the StorageOnDeploy contract.
var StorageOnDeployMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506104d260005560788060236000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80632e64cec114602d575b600080fd5b60005460405190815260200160405180910390f3fea26469706673582212200eeb2dc70c5ef23ca9bb84a56e85db3f9f780572e30c74d437af9b138be6e51564736f6c634300080c0033",
}

// StorageOnDeployABI is the input ABI used to generate the binding from.
// Deprecated: Use StorageOnDeployMetaData.ABI instead.
var StorageOnDeployABI = StorageOnDeployMetaData.ABI

// StorageOnDeployBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StorageOnDeployMetaData.Bin instead.
var StorageOnDeployBin = StorageOnDeployMetaData.Bin

// DeployStorageOnDeploy deploys a new Ethereum contract, binding an instance of StorageOnDeploy to it.
func DeployStorageOnDeploy(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StorageOnDeploy, error) {
	parsed, err := StorageOnDeployMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StorageOnDeployBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StorageOnDeploy{StorageOnDeployCaller: StorageOnDeployCaller{contract: contract}, StorageOnDeployTransactor: StorageOnDeployTransactor{contract: contract}, StorageOnDeployFilterer: StorageOnDeployFilterer{contract: contract}}, nil
}

// StorageOnDeploy is an auto generated Go binding around an Ethereum contract.
type StorageOnDeploy struct {
	StorageOnDeployCaller     // Read-only binding to the contract
	StorageOnDeployTransactor // Write-only binding to the contract
	StorageOnDeployFilterer   // Log filterer for contract events
}

// StorageOnDeployCaller is an auto generated read-only Go binding around an Ethereum contract.
type StorageOnDeployCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageOnDeployTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StorageOnDeployTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageOnDeployFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StorageOnDeployFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageOnDeploySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StorageOnDeploySession struct {
	Contract     *StorageOnDeploy  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StorageOnDeployCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StorageOnDeployCallerSession struct {
	Contract *StorageOnDeployCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// StorageOnDeployTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StorageOnDeployTransactorSession struct {
	Contract     *StorageOnDeployTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// StorageOnDeployRaw is an auto generated low-level Go binding around an Ethereum contract.
type StorageOnDeployRaw struct {
	Contract *StorageOnDeploy // Generic contract binding to access the raw methods on
}

// StorageOnDeployCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StorageOnDeployCallerRaw struct {
	Contract *StorageOnDeployCaller // Generic read-only contract binding to access the raw methods on
}

// StorageOnDeployTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StorageOnDeployTransactorRaw struct {
	Contract *StorageOnDeployTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStorageOnDeploy creates a new instance of StorageOnDeploy, bound to a specific deployed contract.
func NewStorageOnDeploy(address common.Address, backend bind.ContractBackend) (*StorageOnDeploy, error) {
	contract, err := bindStorageOnDeploy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StorageOnDeploy{StorageOnDeployCaller: StorageOnDeployCaller{contract: contract}, StorageOnDeployTransactor: StorageOnDeployTransactor{contract: contract}, StorageOnDeployFilterer: StorageOnDeployFilterer{contract: contract}}, nil
}

// NewStorageOnDeployCaller creates a new read-only instance of StorageOnDeploy, bound to a specific deployed contract.
func NewStorageOnDeployCaller(address common.Address, caller bind.ContractCaller) (*StorageOnDeployCaller, error) {
	contract, err := bindStorageOnDeploy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StorageOnDeployCaller{contract: contract}, nil
}

// NewStorageOnDeployTransactor creates a new write-only instance of StorageOnDeploy, bound to a specific deployed contract.
func NewStorageOnDeployTransactor(address common.Address, transactor bind.ContractTransactor) (*StorageOnDeployTransactor, error) {
	contract, err := bindStorageOnDeploy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StorageOnDeployTransactor{contract: contract}, nil
}

// NewStorageOnDeployFilterer creates a new log filterer instance of StorageOnDeploy, bound to a specific deployed contract.
func NewStorageOnDeployFilterer(address common.Address, filterer bind.ContractFilterer) (*StorageOnDeployFilterer, error) {
	contract, err := bindStorageOnDeploy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StorageOnDeployFilterer{contract: contract}, nil
}

// bindStorageOnDeploy binds a generic wrapper to an already deployed contract.
func bindStorageOnDeploy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StorageOnDeployMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageOnDeploy *StorageOnDeployRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageOnDeploy.Contract.StorageOnDeployCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageOnDeploy *StorageOnDeployRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageOnDeploy.Contract.StorageOnDeployTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageOnDeploy *StorageOnDeployRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageOnDeploy.Contract.StorageOnDeployTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageOnDeploy *StorageOnDeployCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageOnDeploy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageOnDeploy *StorageOnDeployTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageOnDeploy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageOnDeploy *StorageOnDeployTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageOnDeploy.Contract.contract.Transact(opts, method, params...)
}

// Retrieve is a free data retrieval call binding the contract method 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (_StorageOnDeploy *StorageOnDeployCaller) Retrieve(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StorageOnDeploy.contract.Call(opts, &out, "retrieve")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Retrieve is a free data retrieval call binding the contract method 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (_StorageOnDeploy *StorageOnDeploySession) Retrieve() (*big.Int, error) {
	return _StorageOnDeploy.Contract.Retrieve(&_StorageOnDeploy.CallOpts)
}

// Retrieve is a free data retrieval call binding the contract method 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (_StorageOnDeploy *StorageOnDeployCallerSession) Retrieve() (*big.Int, error) {
	return _StorageOnDeploy.Contract.Retrieve(&_StorageOnDeploy.CallOpts)
}
