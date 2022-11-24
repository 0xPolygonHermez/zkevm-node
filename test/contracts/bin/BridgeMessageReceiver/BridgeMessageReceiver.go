// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package BridgeMessageReceiver

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

// BridgeMessageReceiverMetaData contains all meta data concerning the BridgeMessageReceiver contract.
var BridgeMessageReceiverMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"originAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"originNetwork\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onMessageReceived\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506101ac806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80631806b5f214610030575b600080fd5b61004361003e36600461008b565b610057565b604051901515815260200160405180910390f35b6000815160001461006a5750600161006e565b5060005b9392505050565b634e487b7160e01b600052604160045260246000fd5b6000806000606084860312156100a057600080fd5b83356001600160a01b03811681146100b757600080fd5b9250602084013563ffffffff811681146100d057600080fd5b9150604084013567ffffffffffffffff808211156100ed57600080fd5b818601915086601f83011261010157600080fd5b81358181111561011357610113610075565b604051601f8201601f19908116603f0116810190838211818310171561013b5761013b610075565b8160405282815289602084870101111561015457600080fd5b826020860160208301376000602084830101528095505050505050925092509256fea2646970667358221220d0b5aec353a3b80b79514a07d05caf090c80c1c44bbb21e623bb4dec0eed0dad64736f6c634300080c0033",
}

// BridgeMessageReceiverABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeMessageReceiverMetaData.ABI instead.
var BridgeMessageReceiverABI = BridgeMessageReceiverMetaData.ABI

// BridgeMessageReceiverBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BridgeMessageReceiverMetaData.Bin instead.
var BridgeMessageReceiverBin = BridgeMessageReceiverMetaData.Bin

// DeployBridgeMessageReceiver deploys a new Ethereum contract, binding an instance of BridgeMessageReceiver to it.
func DeployBridgeMessageReceiver(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BridgeMessageReceiver, error) {
	parsed, err := BridgeMessageReceiverMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BridgeMessageReceiverBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BridgeMessageReceiver{BridgeMessageReceiverCaller: BridgeMessageReceiverCaller{contract: contract}, BridgeMessageReceiverTransactor: BridgeMessageReceiverTransactor{contract: contract}, BridgeMessageReceiverFilterer: BridgeMessageReceiverFilterer{contract: contract}}, nil
}

// BridgeMessageReceiver is an auto generated Go binding around an Ethereum contract.
type BridgeMessageReceiver struct {
	BridgeMessageReceiverCaller     // Read-only binding to the contract
	BridgeMessageReceiverTransactor // Write-only binding to the contract
	BridgeMessageReceiverFilterer   // Log filterer for contract events
}

// BridgeMessageReceiverCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeMessageReceiverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeMessageReceiverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeMessageReceiverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeMessageReceiverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeMessageReceiverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeMessageReceiverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeMessageReceiverSession struct {
	Contract     *BridgeMessageReceiver // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BridgeMessageReceiverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeMessageReceiverCallerSession struct {
	Contract *BridgeMessageReceiverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// BridgeMessageReceiverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeMessageReceiverTransactorSession struct {
	Contract     *BridgeMessageReceiverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// BridgeMessageReceiverRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeMessageReceiverRaw struct {
	Contract *BridgeMessageReceiver // Generic contract binding to access the raw methods on
}

// BridgeMessageReceiverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeMessageReceiverCallerRaw struct {
	Contract *BridgeMessageReceiverCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeMessageReceiverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeMessageReceiverTransactorRaw struct {
	Contract *BridgeMessageReceiverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeMessageReceiver creates a new instance of BridgeMessageReceiver, bound to a specific deployed contract.
func NewBridgeMessageReceiver(address common.Address, backend bind.ContractBackend) (*BridgeMessageReceiver, error) {
	contract, err := bindBridgeMessageReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeMessageReceiver{BridgeMessageReceiverCaller: BridgeMessageReceiverCaller{contract: contract}, BridgeMessageReceiverTransactor: BridgeMessageReceiverTransactor{contract: contract}, BridgeMessageReceiverFilterer: BridgeMessageReceiverFilterer{contract: contract}}, nil
}

// NewBridgeMessageReceiverCaller creates a new read-only instance of BridgeMessageReceiver, bound to a specific deployed contract.
func NewBridgeMessageReceiverCaller(address common.Address, caller bind.ContractCaller) (*BridgeMessageReceiverCaller, error) {
	contract, err := bindBridgeMessageReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeMessageReceiverCaller{contract: contract}, nil
}

// NewBridgeMessageReceiverTransactor creates a new write-only instance of BridgeMessageReceiver, bound to a specific deployed contract.
func NewBridgeMessageReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeMessageReceiverTransactor, error) {
	contract, err := bindBridgeMessageReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeMessageReceiverTransactor{contract: contract}, nil
}

// NewBridgeMessageReceiverFilterer creates a new log filterer instance of BridgeMessageReceiver, bound to a specific deployed contract.
func NewBridgeMessageReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeMessageReceiverFilterer, error) {
	contract, err := bindBridgeMessageReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeMessageReceiverFilterer{contract: contract}, nil
}

// bindBridgeMessageReceiver binds a generic wrapper to an already deployed contract.
func bindBridgeMessageReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BridgeMessageReceiverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeMessageReceiver *BridgeMessageReceiverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeMessageReceiver.Contract.BridgeMessageReceiverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeMessageReceiver *BridgeMessageReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeMessageReceiver.Contract.BridgeMessageReceiverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeMessageReceiver *BridgeMessageReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeMessageReceiver.Contract.BridgeMessageReceiverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeMessageReceiver *BridgeMessageReceiverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeMessageReceiver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeMessageReceiver *BridgeMessageReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeMessageReceiver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeMessageReceiver *BridgeMessageReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeMessageReceiver.Contract.contract.Transact(opts, method, params...)
}

// OnMessageReceived is a free data retrieval call binding the contract method 0x1806b5f2.
//
// Solidity: function onMessageReceived(address originAddress, uint32 originNetwork, bytes data) view returns(bool)
func (_BridgeMessageReceiver *BridgeMessageReceiverCaller) OnMessageReceived(opts *bind.CallOpts, originAddress common.Address, originNetwork uint32, data []byte) (bool, error) {
	var out []interface{}
	err := _BridgeMessageReceiver.contract.Call(opts, &out, "onMessageReceived", originAddress, originNetwork, data)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// OnMessageReceived is a free data retrieval call binding the contract method 0x1806b5f2.
//
// Solidity: function onMessageReceived(address originAddress, uint32 originNetwork, bytes data) view returns(bool)
func (_BridgeMessageReceiver *BridgeMessageReceiverSession) OnMessageReceived(originAddress common.Address, originNetwork uint32, data []byte) (bool, error) {
	return _BridgeMessageReceiver.Contract.OnMessageReceived(&_BridgeMessageReceiver.CallOpts, originAddress, originNetwork, data)
}

// OnMessageReceived is a free data retrieval call binding the contract method 0x1806b5f2.
//
// Solidity: function onMessageReceived(address originAddress, uint32 originNetwork, bytes data) view returns(bool)
func (_BridgeMessageReceiver *BridgeMessageReceiverCallerSession) OnMessageReceived(originAddress common.Address, originNetwork uint32, data []byte) (bool, error) {
	return _BridgeMessageReceiver.Contract.OnMessageReceived(&_BridgeMessageReceiver.CallOpts, originAddress, originNetwork, data)
}
