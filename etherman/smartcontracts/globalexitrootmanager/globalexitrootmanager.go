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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"globalExitRootNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"mainnetExitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"rollupExitRoot\",\"type\":\"bytes32\"}],\"name\":\"UpdateGlobalExitRoot\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastGlobalExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"globalExitRootMap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rollupAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastGlobalExitRootNum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMainnetExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newRoot\",\"type\":\"bytes32\"}],\"name\":\"updateExitRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061066c806100206000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c806333d6247d11610076578063485cc9551161005b578063485cc955146101445780635ec6a8df14610157578063a3c573eb1461019c57600080fd5b806333d6247d146100f65780633ed691ef1461010b57600080fd5b806301fd9044146100a8578063029f2793146100c4578063257b3632146100cd578063319cf735146100ed575b600080fd5b6100b160015481565b6040519081526020015b60405180910390f35b6100b160045481565b6100b16100db366004610562565b60036020526000908152604090205481565b6100b160025481565b610109610104366004610562565b6101bc565b005b6100b160025460015460408051602081019390935282015260009060600160405160208183030381529060405280519060200120905090565b6101096101523660046105a4565b610381565b6006546101779073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100bb565b6005546101779073ffffffffffffffffffffffffffffffffffffffff1681565b60065473ffffffffffffffffffffffffffffffffffffffff163314806101f9575060055473ffffffffffffffffffffffffffffffffffffffff1633145b61028a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603d60248201527f476c6f62616c45786974526f6f744d616e616765723a3a75706461746545786960448201527f74526f6f743a204f4e4c595f414c4c4f5745445f434f4e54524143545300000060648201526084015b60405180910390fd5b60065473ffffffffffffffffffffffffffffffffffffffff1633036102af5760018190555b60055473ffffffffffffffffffffffffffffffffffffffff1633036102d45760028190555b600480549060006102e4836105d7565b9091555050600254600154604080516020810193909352820152600090606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152828252805160209182012060045460008281526003909352928220839055600154600254919550939092917fb7c409af8cb511116b88f38824d48a0196194596241fdb2d177210d3d3b89fbf91a45050565b600054610100900460ff16158080156103a15750600054600160ff909116105b806103bb5750303b1580156103bb575060005460ff166001145b610447576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610281565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905580156104a557600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b6006805473ffffffffffffffffffffffffffffffffffffffff8086167fffffffffffffffffffffffff0000000000000000000000000000000000000000928316179092556005805492851692909116919091179055801561055d57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b60006020828403121561057457600080fd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461059f57600080fd5b919050565b600080604083850312156105b757600080fd5b6105c08361057b565b91506105ce6020840161057b565b90509250929050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361062f577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b506001019056fea2646970667358221220d9ccfdddba433da7a9a14980d28db4b102be77f8af6cb74001850e88e278061864736f6c634300080f0033",
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

// LastGlobalExitRootNum is a free data retrieval call binding the contract method 0x029f2793.
//
// Solidity: function lastGlobalExitRootNum() view returns(uint256)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) LastGlobalExitRootNum(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "lastGlobalExitRootNum")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastGlobalExitRootNum is a free data retrieval call binding the contract method 0x029f2793.
//
// Solidity: function lastGlobalExitRootNum() view returns(uint256)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) LastGlobalExitRootNum() (*big.Int, error) {
	return _Globalexitrootmanager.Contract.LastGlobalExitRootNum(&_Globalexitrootmanager.CallOpts)
}

// LastGlobalExitRootNum is a free data retrieval call binding the contract method 0x029f2793.
//
// Solidity: function lastGlobalExitRootNum() view returns(uint256)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) LastGlobalExitRootNum() (*big.Int, error) {
	return _Globalexitrootmanager.Contract.LastGlobalExitRootNum(&_Globalexitrootmanager.CallOpts)
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
	GlobalExitRootNum *big.Int
	MainnetExitRoot   [32]byte
	RollupExitRoot    [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUpdateGlobalExitRoot is a free log retrieval operation binding the contract event 0xb7c409af8cb511116b88f38824d48a0196194596241fdb2d177210d3d3b89fbf.
//
// Solidity: event UpdateGlobalExitRoot(uint256 indexed globalExitRootNum, bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) FilterUpdateGlobalExitRoot(opts *bind.FilterOpts, globalExitRootNum []*big.Int, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (*GlobalexitrootmanagerUpdateGlobalExitRootIterator, error) {

	var globalExitRootNumRule []interface{}
	for _, globalExitRootNumItem := range globalExitRootNum {
		globalExitRootNumRule = append(globalExitRootNumRule, globalExitRootNumItem)
	}
	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Globalexitrootmanager.contract.FilterLogs(opts, "UpdateGlobalExitRoot", globalExitRootNumRule, mainnetExitRootRule, rollupExitRootRule)
	if err != nil {
		return nil, err
	}
	return &GlobalexitrootmanagerUpdateGlobalExitRootIterator{contract: _Globalexitrootmanager.contract, event: "UpdateGlobalExitRoot", logs: logs, sub: sub}, nil
}

// WatchUpdateGlobalExitRoot is a free log subscription operation binding the contract event 0xb7c409af8cb511116b88f38824d48a0196194596241fdb2d177210d3d3b89fbf.
//
// Solidity: event UpdateGlobalExitRoot(uint256 indexed globalExitRootNum, bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) WatchUpdateGlobalExitRoot(opts *bind.WatchOpts, sink chan<- *GlobalexitrootmanagerUpdateGlobalExitRoot, globalExitRootNum []*big.Int, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (event.Subscription, error) {

	var globalExitRootNumRule []interface{}
	for _, globalExitRootNumItem := range globalExitRootNum {
		globalExitRootNumRule = append(globalExitRootNumRule, globalExitRootNumItem)
	}
	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Globalexitrootmanager.contract.WatchLogs(opts, "UpdateGlobalExitRoot", globalExitRootNumRule, mainnetExitRootRule, rollupExitRootRule)
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

// ParseUpdateGlobalExitRoot is a log parse operation binding the contract event 0xb7c409af8cb511116b88f38824d48a0196194596241fdb2d177210d3d3b89fbf.
//
// Solidity: event UpdateGlobalExitRoot(uint256 indexed globalExitRootNum, bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) ParseUpdateGlobalExitRoot(log types.Log) (*GlobalexitrootmanagerUpdateGlobalExitRoot, error) {
	event := new(GlobalexitrootmanagerUpdateGlobalExitRoot)
	if err := _Globalexitrootmanager.contract.UnpackLog(event, "UpdateGlobalExitRoot", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
