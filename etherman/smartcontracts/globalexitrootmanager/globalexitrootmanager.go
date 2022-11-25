// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package globalexitrootmanager

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

// GlobalexitrootmanagerMetaData contains all meta data concerning the Globalexitrootmanager contract.
var GlobalexitrootmanagerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"mainnetExitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"rollupExitRoot\",\"type\":\"bytes32\"}],\"name\":\"UpdateGlobalExitRoot\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastGlobalExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"globalExitRootMap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rollupAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMainnetExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newRoot\",\"type\":\"bytes32\"}],\"name\":\"updateExitRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506104cd806100206000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80633ed691ef1161005b5780633ed691ef146100e7578063485cc955146101205780635ec6a8df14610133578063a3c573eb1461015e57600080fd5b806301fd90441461008d578063257b3632146100a9578063319cf735146100c957806333d6247d146100d2575b600080fd5b61009660015481565b6040519081526020015b60405180910390f35b6100966100b736600461042f565b60036020526000908152604090205481565b61009660025481565b6100e56100e036600461042f565b610171565b005b61009660025460015460408051602081019390935282015260009060600160405160208183030381529060405280519060200120905090565b6100e561012e366004610464565b6102cd565b600554610146906001600160a01b031681565b6040516001600160a01b0390911681526020016100a0565b600454610146906001600160a01b031681565b6005546001600160a01b031633148061019457506004546001600160a01b031633145b61020b5760405162461bcd60e51b815260206004820152603d60248201527f476c6f62616c45786974526f6f744d616e616765723a3a75706461746545786960448201527f74526f6f743a204f4e4c595f414c4c4f5745445f434f4e54524143545300000060648201526084015b60405180910390fd5b6005546001600160a01b031633036102235760018190555b6004546001600160a01b0316330361023b5760028190555b60025460015460408051602081019390935282015260009060600160405160208183030381529060405280519060200120905060036000828152602001908152602001600020546000036102c957600081815260036020526040808220429055600154600254915190927f61014378f82a0d809aefaf87a8ac9505b89c321808287a6e7810f29304c1fce391a35b5050565b600054610100900460ff16158080156102ed5750600054600160ff909116105b806103075750303b158015610307575060005460ff166001145b6103795760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610202565b6000805460ff19166001179055801561039c576000805461ff0019166101001790555b600580546001600160a01b038086167fffffffffffffffffffffffff0000000000000000000000000000000000000000928316179092556004805492851692909116919091179055801561042a576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b60006020828403121561044157600080fd5b5035919050565b80356001600160a01b038116811461045f57600080fd5b919050565b6000806040838503121561047757600080fd5b61048083610448565b915061048e60208401610448565b9050925092905056fea264697066735822122017f0e386b65d4b488e1ae82e0a2d700352d5ea3125ecc66d095224034edccd2164736f6c634300080f0033",
}

// GlobalexitrootmanagerABI is the input ABI used to generate the binding from.
// Deprecated: Use GlobalexitrootmanagerMetaData.ABI instead.
var GlobalexitrootmanagerABI = GlobalexitrootmanagerMetaData.ABI

// GlobalexitrootmanagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GlobalexitrootmanagerMetaData.Bin instead.
var GlobalexitrootmanagerBin = GlobalexitrootmanagerMetaData.Bin

// DeployGlobalexitrootmanager deploys a new Ethereum contract, binding an instance of Globalexitrootmanager to it.
func DeployGlobalexitrootmanager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Globalexitrootmanager, error) {
	parsed, err := GlobalexitrootmanagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GlobalexitrootmanagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Globalexitrootmanager{GlobalexitrootmanagerCaller: GlobalexitrootmanagerCaller{contract: contract}, GlobalexitrootmanagerTransactor: GlobalexitrootmanagerTransactor{contract: contract}, GlobalexitrootmanagerFilterer: GlobalexitrootmanagerFilterer{contract: contract}}, nil
}

// Globalexitrootmanager is an auto generated Go binding around an Ethereum contract.
type Globalexitrootmanager struct {
	GlobalexitrootmanagerCaller     // Read-only binding to the contract
	GlobalexitrootmanagerTransactor // Write-only binding to the contract
	GlobalexitrootmanagerFilterer   // Log filterer for contract events
}

// GlobalexitrootmanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type GlobalexitrootmanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GlobalexitrootmanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GlobalexitrootmanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GlobalexitrootmanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GlobalexitrootmanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GlobalexitrootmanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GlobalexitrootmanagerSession struct {
	Contract     *Globalexitrootmanager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// GlobalexitrootmanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GlobalexitrootmanagerCallerSession struct {
	Contract *GlobalexitrootmanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// GlobalexitrootmanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GlobalexitrootmanagerTransactorSession struct {
	Contract     *GlobalexitrootmanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// GlobalexitrootmanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type GlobalexitrootmanagerRaw struct {
	Contract *Globalexitrootmanager // Generic contract binding to access the raw methods on
}

// GlobalexitrootmanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GlobalexitrootmanagerCallerRaw struct {
	Contract *GlobalexitrootmanagerCaller // Generic read-only contract binding to access the raw methods on
}

// GlobalexitrootmanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GlobalexitrootmanagerTransactorRaw struct {
	Contract *GlobalexitrootmanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGlobalexitrootmanager creates a new instance of Globalexitrootmanager, bound to a specific deployed contract.
func NewGlobalexitrootmanager(address common.Address, backend bind.ContractBackend) (*Globalexitrootmanager, error) {
	contract, err := bindGlobalexitrootmanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Globalexitrootmanager{GlobalexitrootmanagerCaller: GlobalexitrootmanagerCaller{contract: contract}, GlobalexitrootmanagerTransactor: GlobalexitrootmanagerTransactor{contract: contract}, GlobalexitrootmanagerFilterer: GlobalexitrootmanagerFilterer{contract: contract}}, nil
}

// NewGlobalexitrootmanagerCaller creates a new read-only instance of Globalexitrootmanager, bound to a specific deployed contract.
func NewGlobalexitrootmanagerCaller(address common.Address, caller bind.ContractCaller) (*GlobalexitrootmanagerCaller, error) {
	contract, err := bindGlobalexitrootmanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GlobalexitrootmanagerCaller{contract: contract}, nil
}

// NewGlobalexitrootmanagerTransactor creates a new write-only instance of Globalexitrootmanager, bound to a specific deployed contract.
func NewGlobalexitrootmanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*GlobalexitrootmanagerTransactor, error) {
	contract, err := bindGlobalexitrootmanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GlobalexitrootmanagerTransactor{contract: contract}, nil
}

// NewGlobalexitrootmanagerFilterer creates a new log filterer instance of Globalexitrootmanager, bound to a specific deployed contract.
func NewGlobalexitrootmanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*GlobalexitrootmanagerFilterer, error) {
	contract, err := bindGlobalexitrootmanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GlobalexitrootmanagerFilterer{contract: contract}, nil
}

// bindGlobalexitrootmanager binds a generic wrapper to an already deployed contract.
func bindGlobalexitrootmanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GlobalexitrootmanagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Globalexitrootmanager *GlobalexitrootmanagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Globalexitrootmanager.Contract.GlobalexitrootmanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Globalexitrootmanager *GlobalexitrootmanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.GlobalexitrootmanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Globalexitrootmanager *GlobalexitrootmanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.GlobalexitrootmanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Globalexitrootmanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.contract.Transact(opts, method, params...)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) BridgeAddress() (common.Address, error) {
	return _Globalexitrootmanager.Contract.BridgeAddress(&_Globalexitrootmanager.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) BridgeAddress() (common.Address, error) {
	return _Globalexitrootmanager.Contract.BridgeAddress(&_Globalexitrootmanager.CallOpts)
}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) GetLastGlobalExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "getLastGlobalExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) GetLastGlobalExitRoot() ([32]byte, error) {
	return _Globalexitrootmanager.Contract.GetLastGlobalExitRoot(&_Globalexitrootmanager.CallOpts)
}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) GetLastGlobalExitRoot() ([32]byte, error) {
	return _Globalexitrootmanager.Contract.GetLastGlobalExitRoot(&_Globalexitrootmanager.CallOpts)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) GlobalExitRootMap(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "globalExitRootMap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Globalexitrootmanager.Contract.GlobalExitRootMap(&_Globalexitrootmanager.CallOpts, arg0)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Globalexitrootmanager.Contract.GlobalExitRootMap(&_Globalexitrootmanager.CallOpts, arg0)
}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) LastMainnetExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "lastMainnetExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) LastMainnetExitRoot() ([32]byte, error) {
	return _Globalexitrootmanager.Contract.LastMainnetExitRoot(&_Globalexitrootmanager.CallOpts)
}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) LastMainnetExitRoot() ([32]byte, error) {
	return _Globalexitrootmanager.Contract.LastMainnetExitRoot(&_Globalexitrootmanager.CallOpts)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) LastRollupExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "lastRollupExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) LastRollupExitRoot() ([32]byte, error) {
	return _Globalexitrootmanager.Contract.LastRollupExitRoot(&_Globalexitrootmanager.CallOpts)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) LastRollupExitRoot() ([32]byte, error) {
	return _Globalexitrootmanager.Contract.LastRollupExitRoot(&_Globalexitrootmanager.CallOpts)
}

// RollupAddress is a free data retrieval call binding the contract method 0x5ec6a8df.
//
// Solidity: function rollupAddress() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) RollupAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "rollupAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RollupAddress is a free data retrieval call binding the contract method 0x5ec6a8df.
//
// Solidity: function rollupAddress() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) RollupAddress() (common.Address, error) {
	return _Globalexitrootmanager.Contract.RollupAddress(&_Globalexitrootmanager.CallOpts)
}

// RollupAddress is a free data retrieval call binding the contract method 0x5ec6a8df.
//
// Solidity: function rollupAddress() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) RollupAddress() (common.Address, error) {
	return _Globalexitrootmanager.Contract.RollupAddress(&_Globalexitrootmanager.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _rollupAddress, address _bridgeAddress) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactor) Initialize(opts *bind.TransactOpts, _rollupAddress common.Address, _bridgeAddress common.Address) (*types.Transaction, error) {
	return _Globalexitrootmanager.contract.Transact(opts, "initialize", _rollupAddress, _bridgeAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _rollupAddress, address _bridgeAddress) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) Initialize(_rollupAddress common.Address, _bridgeAddress common.Address) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.Initialize(&_Globalexitrootmanager.TransactOpts, _rollupAddress, _bridgeAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _rollupAddress, address _bridgeAddress) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactorSession) Initialize(_rollupAddress common.Address, _bridgeAddress common.Address) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.Initialize(&_Globalexitrootmanager.TransactOpts, _rollupAddress, _bridgeAddress)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactor) UpdateExitRoot(opts *bind.TransactOpts, newRoot [32]byte) (*types.Transaction, error) {
	return _Globalexitrootmanager.contract.Transact(opts, "updateExitRoot", newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.UpdateExitRoot(&_Globalexitrootmanager.TransactOpts, newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactorSession) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.UpdateExitRoot(&_Globalexitrootmanager.TransactOpts, newRoot)
}

// GlobalexitrootmanagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Globalexitrootmanager contract.
type GlobalexitrootmanagerInitializedIterator struct {
	Event *GlobalexitrootmanagerInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GlobalexitrootmanagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GlobalexitrootmanagerInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GlobalexitrootmanagerInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GlobalexitrootmanagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GlobalexitrootmanagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GlobalexitrootmanagerInitialized represents a Initialized event raised by the Globalexitrootmanager contract.
type GlobalexitrootmanagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*GlobalexitrootmanagerInitializedIterator, error) {

	logs, sub, err := _Globalexitrootmanager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GlobalexitrootmanagerInitializedIterator{contract: _Globalexitrootmanager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GlobalexitrootmanagerInitialized) (event.Subscription, error) {

	logs, sub, err := _Globalexitrootmanager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GlobalexitrootmanagerInitialized)
				if err := _Globalexitrootmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) ParseInitialized(log types.Log) (*GlobalexitrootmanagerInitialized, error) {
	event := new(GlobalexitrootmanagerInitialized)
	if err := _Globalexitrootmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GlobalexitrootmanagerUpdateGlobalExitRootIterator is returned from FilterUpdateGlobalExitRoot and is used to iterate over the raw logs and unpacked data for UpdateGlobalExitRoot events raised by the Globalexitrootmanager contract.
type GlobalexitrootmanagerUpdateGlobalExitRootIterator struct {
	Event *GlobalexitrootmanagerUpdateGlobalExitRoot // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GlobalexitrootmanagerUpdateGlobalExitRootIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GlobalexitrootmanagerUpdateGlobalExitRoot)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GlobalexitrootmanagerUpdateGlobalExitRoot)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GlobalexitrootmanagerUpdateGlobalExitRootIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GlobalexitrootmanagerUpdateGlobalExitRootIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GlobalexitrootmanagerUpdateGlobalExitRoot represents a UpdateGlobalExitRoot event raised by the Globalexitrootmanager contract.
type GlobalexitrootmanagerUpdateGlobalExitRoot struct {
	MainnetExitRoot [32]byte
	RollupExitRoot  [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpdateGlobalExitRoot is a free log retrieval operation binding the contract event 0x61014378f82a0d809aefaf87a8ac9505b89c321808287a6e7810f29304c1fce3.
//
// Solidity: event UpdateGlobalExitRoot(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) FilterUpdateGlobalExitRoot(opts *bind.FilterOpts, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (*GlobalexitrootmanagerUpdateGlobalExitRootIterator, error) {

	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Globalexitrootmanager.contract.FilterLogs(opts, "UpdateGlobalExitRoot", mainnetExitRootRule, rollupExitRootRule)
	if err != nil {
		return nil, err
	}
	return &GlobalexitrootmanagerUpdateGlobalExitRootIterator{contract: _Globalexitrootmanager.contract, event: "UpdateGlobalExitRoot", logs: logs, sub: sub}, nil
}

// WatchUpdateGlobalExitRoot is a free log subscription operation binding the contract event 0x61014378f82a0d809aefaf87a8ac9505b89c321808287a6e7810f29304c1fce3.
//
// Solidity: event UpdateGlobalExitRoot(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) WatchUpdateGlobalExitRoot(opts *bind.WatchOpts, sink chan<- *GlobalexitrootmanagerUpdateGlobalExitRoot, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (event.Subscription, error) {

	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Globalexitrootmanager.contract.WatchLogs(opts, "UpdateGlobalExitRoot", mainnetExitRootRule, rollupExitRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GlobalexitrootmanagerUpdateGlobalExitRoot)
				if err := _Globalexitrootmanager.contract.UnpackLog(event, "UpdateGlobalExitRoot", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdateGlobalExitRoot is a log parse operation binding the contract event 0x61014378f82a0d809aefaf87a8ac9505b89c321808287a6e7810f29304c1fce3.
//
// Solidity: event UpdateGlobalExitRoot(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) ParseUpdateGlobalExitRoot(log types.Log) (*GlobalexitrootmanagerUpdateGlobalExitRoot, error) {
	event := new(GlobalexitrootmanagerUpdateGlobalExitRoot)
	if err := _Globalexitrootmanager.contract.UnpackLog(event, "UpdateGlobalExitRoot", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
