// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Caller

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

// CallerMetaData contains all meta data concerning the Caller contract.
var CallerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_num\",\"type\":\"uint256\"}],\"name\":\"call\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_num\",\"type\":\"uint256\"}],\"name\":\"delegateCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"name\":\"invalidStaticCallLessParameters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"name\":\"invalidStaticCallMoreParameters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"name\":\"invalidStaticCallWithInnerCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_num\",\"type\":\"uint256\"}],\"name\":\"multiCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"preEcrecover_0\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"name\":\"staticCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506107e2806100206000396000f3fe60806040526004361061007b5760003560e01c806387b1d6ad1161004e57806387b1d6ad146100c8578063b3e554d0146100db578063c6c211e9146100fb578063fff0972f1461010e57600080fd5b80630f6be00d14610080578063351f14c5146100955780633bd9ef28146100b557806369c2b58f14610095575b600080fd5b61009361008e3660046106e8565b610123565b005b3480156100a157600080fd5b506100936100b0366004610714565b610214565b6100936100c33660046106e8565b61030b565b6100936100d6366004610714565b6103ed565b3480156100e757600080fd5b506100936100f6366004610714565b6104e8565b6100936101093660046106e8565b6105f0565b34801561011a57600080fd5b5061009361060d565b6000826001600160a01b03168260405160240161014291815260200190565b60408051601f198184030181529181526020820180516001600160e01b0316636466414b60e01b179052516101779190610738565b600060405180830381855af49150503d80600081146101b2576040519150601f19603f3d011682016040523d82523d6000602084013e6101b7565b606091505b5050809150508061020f5760405162461bcd60e51b815260206004820152601f60248201527f6661696c656420746f20706572666f726d2064656c65676174652063616c6c0060448201526064015b60405180910390fd5b505050565b60408051600481526024810182526020810180516001600160e01b03166305d77c8f60e21b17905290516000916001600160a01b038416916102569190610738565b600060405180830381855afa9150503d8060008114610291576040519150601f19603f3d011682016040523d82523d6000602084013e610296565b606091505b509091505080156103075760405162461bcd60e51b815260206004820152603560248201527f7374617469632063616c6c2077617320737570706f73656420746f206661696c6044820152742077697468206c65737320706172616d657465727360581b6064820152608401610206565b5050565b6000826001600160a01b03168260405160240161032a91815260200190565b60408051601f198184030181529181526020820180516001600160e01b0316636466414b60e01b1790525161035f9190610738565b6000604051808303816000865af19150503d806000811461039c576040519150601f19603f3d011682016040523d82523d6000602084013e6103a1565b606091505b5050809150508061020f5760405162461bcd60e51b815260206004820152601660248201527519985a5b1959081d1bc81c195c999bdc9b4818d85b1b60521b6044820152606401610206565b60408051600481526024810182526020810180516001600160e01b031663813d8a3760e01b17905290516000916060916001600160a01b0385169161043191610738565b600060405180830381855afa9150503d806000811461046c576040519150601f19603f3d011682016040523d82523d6000602084013e610471565b606091505b509092509050816104c45760405162461bcd60e51b815260206004820152601d60248201527f6661696c656420746f20706572666f726d207374617469632063616c6c0000006044820152606401610206565b6000806000838060200190518101906104dd9190610773565b505050505050505050565b60405160016024820152600260448201526000906001600160a01b0383169060640160408051601f198184030181529181526020820180516001600160e01b03166305d77c8f60e21b1790525161053f9190610738565b600060405180830381855afa9150503d806000811461057a576040519150601f19603f3d011682016040523d82523d6000602084013e61057f565b606091505b509091505080156103075760405162461bcd60e51b815260206004820152603560248201527f7374617469632063616c6c2077617320737570706f73656420746f206661696c6044820152742077697468206d6f726520706172616d657465727360581b6064820152608401610206565b6105fa828261030b565b6106048282610123565b610307826103ed565b6040805160008152602081018083527f456e9aea5e197a1f1af7a3e85a3212fa4049a3ba34c2289b4c860fc0b0c64ef390819052601c9282018390527f9242685bf161793cc25603c231bc2f568eb630ea16aa137d2664ac8038825608606083018190527f4f8ae3bd7535248d0bd448298cc2e2071e56992d0774dc340c368ae950852ada6080840181905291939290919060019060a0016020604051602081039080840390855afa1580156106c7573d6000803e3d6000fd5b50505050505050565b6001600160a01b03811681146106e557600080fd5b50565b600080604083850312156106fb57600080fd5b8235610706816106d0565b946020939093013593505050565b60006020828403121561072657600080fd5b8135610731816106d0565b9392505050565b6000825160005b81811015610759576020818601810151858301520161073f565b81811115610768576000828501525b509190910192915050565b60008060006060848603121561078857600080fd5b83519250602084015161079a816106d0565b8092505060408401519050925092509256fea2646970667358221220a82281d9696e9c0b1485742fdc267f1a20789436e6606892338b4fa4f6f1e43e64736f6c634300080c0033",
}

// CallerABI is the input ABI used to generate the binding from.
// Deprecated: Use CallerMetaData.ABI instead.
var CallerABI = CallerMetaData.ABI

// CallerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CallerMetaData.Bin instead.
var CallerBin = CallerMetaData.Bin

// DeployCaller deploys a new Ethereum contract, binding an instance of Caller to it.
func DeployCaller(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Caller, error) {
	parsed, err := CallerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CallerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Caller{CallerCaller: CallerCaller{contract: contract}, CallerTransactor: CallerTransactor{contract: contract}, CallerFilterer: CallerFilterer{contract: contract}}, nil
}

// Caller is an auto generated Go binding around an Ethereum contract.
type Caller struct {
	CallerCaller     // Read-only binding to the contract
	CallerTransactor // Write-only binding to the contract
	CallerFilterer   // Log filterer for contract events
}

// CallerCaller is an auto generated read-only Go binding around an Ethereum contract.
type CallerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CallerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CallerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CallerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CallerSession struct {
	Contract     *Caller           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CallerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CallerCallerSession struct {
	Contract *CallerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// CallerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CallerTransactorSession struct {
	Contract     *CallerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CallerRaw is an auto generated low-level Go binding around an Ethereum contract.
type CallerRaw struct {
	Contract *Caller // Generic contract binding to access the raw methods on
}

// CallerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CallerCallerRaw struct {
	Contract *CallerCaller // Generic read-only contract binding to access the raw methods on
}

// CallerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CallerTransactorRaw struct {
	Contract *CallerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCaller creates a new instance of Caller, bound to a specific deployed contract.
func NewCaller(address common.Address, backend bind.ContractBackend) (*Caller, error) {
	contract, err := bindCaller(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Caller{CallerCaller: CallerCaller{contract: contract}, CallerTransactor: CallerTransactor{contract: contract}, CallerFilterer: CallerFilterer{contract: contract}}, nil
}

// NewCallerCaller creates a new read-only instance of Caller, bound to a specific deployed contract.
func NewCallerCaller(address common.Address, caller bind.ContractCaller) (*CallerCaller, error) {
	contract, err := bindCaller(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CallerCaller{contract: contract}, nil
}

// NewCallerTransactor creates a new write-only instance of Caller, bound to a specific deployed contract.
func NewCallerTransactor(address common.Address, transactor bind.ContractTransactor) (*CallerTransactor, error) {
	contract, err := bindCaller(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CallerTransactor{contract: contract}, nil
}

// NewCallerFilterer creates a new log filterer instance of Caller, bound to a specific deployed contract.
func NewCallerFilterer(address common.Address, filterer bind.ContractFilterer) (*CallerFilterer, error) {
	contract, err := bindCaller(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CallerFilterer{contract: contract}, nil
}

// bindCaller binds a generic wrapper to an already deployed contract.
func bindCaller(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CallerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Caller *CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Caller.Contract.CallerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Caller *CallerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Caller.Contract.CallerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Caller *CallerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Caller.Contract.CallerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Caller *CallerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Caller.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Caller *CallerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Caller.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Caller *CallerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Caller.Contract.contract.Transact(opts, method, params...)
}

// Call is a paid mutator transaction binding the contract method 0x3bd9ef28.
//
// Solidity: function call(address _contract, uint256 _num) payable returns()
func (_Caller *CallerTransactor) Call(opts *bind.TransactOpts, _contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.contract.Transact(opts, "call", _contract, _num)
}

// Call is a paid mutator transaction binding the contract method 0x3bd9ef28.
//
// Solidity: function call(address _contract, uint256 _num) payable returns()
func (_Caller *CallerSession) Call(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.Contract.Call(&_Caller.TransactOpts, _contract, _num)
}

// Call is a paid mutator transaction binding the contract method 0x3bd9ef28.
//
// Solidity: function call(address _contract, uint256 _num) payable returns()
func (_Caller *CallerTransactorSession) Call(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.Contract.Call(&_Caller.TransactOpts, _contract, _num)
}

// DelegateCall is a paid mutator transaction binding the contract method 0x0f6be00d.
//
// Solidity: function delegateCall(address _contract, uint256 _num) payable returns()
func (_Caller *CallerTransactor) DelegateCall(opts *bind.TransactOpts, _contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.contract.Transact(opts, "delegateCall", _contract, _num)
}

// DelegateCall is a paid mutator transaction binding the contract method 0x0f6be00d.
//
// Solidity: function delegateCall(address _contract, uint256 _num) payable returns()
func (_Caller *CallerSession) DelegateCall(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.Contract.DelegateCall(&_Caller.TransactOpts, _contract, _num)
}

// DelegateCall is a paid mutator transaction binding the contract method 0x0f6be00d.
//
// Solidity: function delegateCall(address _contract, uint256 _num) payable returns()
func (_Caller *CallerTransactorSession) DelegateCall(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.Contract.DelegateCall(&_Caller.TransactOpts, _contract, _num)
}

// InvalidStaticCallLessParameters is a paid mutator transaction binding the contract method 0x69c2b58f.
//
// Solidity: function invalidStaticCallLessParameters(address _contract) returns()
func (_Caller *CallerTransactor) InvalidStaticCallLessParameters(opts *bind.TransactOpts, _contract common.Address) (*types.Transaction, error) {
	return _Caller.contract.Transact(opts, "invalidStaticCallLessParameters", _contract)
}

// InvalidStaticCallLessParameters is a paid mutator transaction binding the contract method 0x69c2b58f.
//
// Solidity: function invalidStaticCallLessParameters(address _contract) returns()
func (_Caller *CallerSession) InvalidStaticCallLessParameters(_contract common.Address) (*types.Transaction, error) {
	return _Caller.Contract.InvalidStaticCallLessParameters(&_Caller.TransactOpts, _contract)
}

// InvalidStaticCallLessParameters is a paid mutator transaction binding the contract method 0x69c2b58f.
//
// Solidity: function invalidStaticCallLessParameters(address _contract) returns()
func (_Caller *CallerTransactorSession) InvalidStaticCallLessParameters(_contract common.Address) (*types.Transaction, error) {
	return _Caller.Contract.InvalidStaticCallLessParameters(&_Caller.TransactOpts, _contract)
}

// InvalidStaticCallMoreParameters is a paid mutator transaction binding the contract method 0xb3e554d0.
//
// Solidity: function invalidStaticCallMoreParameters(address _contract) returns()
func (_Caller *CallerTransactor) InvalidStaticCallMoreParameters(opts *bind.TransactOpts, _contract common.Address) (*types.Transaction, error) {
	return _Caller.contract.Transact(opts, "invalidStaticCallMoreParameters", _contract)
}

// InvalidStaticCallMoreParameters is a paid mutator transaction binding the contract method 0xb3e554d0.
//
// Solidity: function invalidStaticCallMoreParameters(address _contract) returns()
func (_Caller *CallerSession) InvalidStaticCallMoreParameters(_contract common.Address) (*types.Transaction, error) {
	return _Caller.Contract.InvalidStaticCallMoreParameters(&_Caller.TransactOpts, _contract)
}

// InvalidStaticCallMoreParameters is a paid mutator transaction binding the contract method 0xb3e554d0.
//
// Solidity: function invalidStaticCallMoreParameters(address _contract) returns()
func (_Caller *CallerTransactorSession) InvalidStaticCallMoreParameters(_contract common.Address) (*types.Transaction, error) {
	return _Caller.Contract.InvalidStaticCallMoreParameters(&_Caller.TransactOpts, _contract)
}

// InvalidStaticCallWithInnerCall is a paid mutator transaction binding the contract method 0x351f14c5.
//
// Solidity: function invalidStaticCallWithInnerCall(address _contract) returns()
func (_Caller *CallerTransactor) InvalidStaticCallWithInnerCall(opts *bind.TransactOpts, _contract common.Address) (*types.Transaction, error) {
	return _Caller.contract.Transact(opts, "invalidStaticCallWithInnerCall", _contract)
}

// InvalidStaticCallWithInnerCall is a paid mutator transaction binding the contract method 0x351f14c5.
//
// Solidity: function invalidStaticCallWithInnerCall(address _contract) returns()
func (_Caller *CallerSession) InvalidStaticCallWithInnerCall(_contract common.Address) (*types.Transaction, error) {
	return _Caller.Contract.InvalidStaticCallWithInnerCall(&_Caller.TransactOpts, _contract)
}

// InvalidStaticCallWithInnerCall is a paid mutator transaction binding the contract method 0x351f14c5.
//
// Solidity: function invalidStaticCallWithInnerCall(address _contract) returns()
func (_Caller *CallerTransactorSession) InvalidStaticCallWithInnerCall(_contract common.Address) (*types.Transaction, error) {
	return _Caller.Contract.InvalidStaticCallWithInnerCall(&_Caller.TransactOpts, _contract)
}

// MultiCall is a paid mutator transaction binding the contract method 0xc6c211e9.
//
// Solidity: function multiCall(address _contract, uint256 _num) payable returns()
func (_Caller *CallerTransactor) MultiCall(opts *bind.TransactOpts, _contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.contract.Transact(opts, "multiCall", _contract, _num)
}

// MultiCall is a paid mutator transaction binding the contract method 0xc6c211e9.
//
// Solidity: function multiCall(address _contract, uint256 _num) payable returns()
func (_Caller *CallerSession) MultiCall(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.Contract.MultiCall(&_Caller.TransactOpts, _contract, _num)
}

// MultiCall is a paid mutator transaction binding the contract method 0xc6c211e9.
//
// Solidity: function multiCall(address _contract, uint256 _num) payable returns()
func (_Caller *CallerTransactorSession) MultiCall(_contract common.Address, _num *big.Int) (*types.Transaction, error) {
	return _Caller.Contract.MultiCall(&_Caller.TransactOpts, _contract, _num)
}

// PreEcrecover0 is a paid mutator transaction binding the contract method 0xfff0972f.
//
// Solidity: function preEcrecover_0() returns()
func (_Caller *CallerTransactor) PreEcrecover0(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Caller.contract.Transact(opts, "preEcrecover_0")
}

// PreEcrecover0 is a paid mutator transaction binding the contract method 0xfff0972f.
//
// Solidity: function preEcrecover_0() returns()
func (_Caller *CallerSession) PreEcrecover0() (*types.Transaction, error) {
	return _Caller.Contract.PreEcrecover0(&_Caller.TransactOpts)
}

// PreEcrecover0 is a paid mutator transaction binding the contract method 0xfff0972f.
//
// Solidity: function preEcrecover_0() returns()
func (_Caller *CallerTransactorSession) PreEcrecover0() (*types.Transaction, error) {
	return _Caller.Contract.PreEcrecover0(&_Caller.TransactOpts)
}

// StaticCall is a paid mutator transaction binding the contract method 0x87b1d6ad.
//
// Solidity: function staticCall(address _contract) payable returns()
func (_Caller *CallerTransactor) StaticCall(opts *bind.TransactOpts, _contract common.Address) (*types.Transaction, error) {
	return _Caller.contract.Transact(opts, "staticCall", _contract)
}

// StaticCall is a paid mutator transaction binding the contract method 0x87b1d6ad.
//
// Solidity: function staticCall(address _contract) payable returns()
func (_Caller *CallerSession) StaticCall(_contract common.Address) (*types.Transaction, error) {
	return _Caller.Contract.StaticCall(&_Caller.TransactOpts, _contract)
}

// StaticCall is a paid mutator transaction binding the contract method 0x87b1d6ad.
//
// Solidity: function staticCall(address _contract) payable returns()
func (_Caller *CallerTransactorSession) StaticCall(_contract common.Address) (*types.Transaction, error) {
	return _Caller.Contract.StaticCall(&_Caller.TransactOpts, _contract)
}
