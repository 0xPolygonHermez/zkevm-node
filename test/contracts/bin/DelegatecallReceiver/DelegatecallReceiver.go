// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package DelegatecallReceiver

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

// DelegatecallReceiverMetaData contains all meta data concerning the DelegatecallReceiver contract.
var DelegatecallReceiverMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"entrypoint\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"expectedSender\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610415806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80631754dba51461003b578063a65d69d41461006a575b600080fd5b60005461004e906001600160a01b031681565b6040516001600160a01b03909116815260200160405180910390f35b610072610074565b005b6000546001600160a01b03163381149061008f9060146100e4565b61009a3360146100e4565b6040516020016100ab9291906102b7565b604051602081830303815290604052906100e15760405162461bcd60e51b81526004016100d8919061031c565b60405180910390fd5b50565b606060006100f3836002610365565b6100fe906002610384565b67ffffffffffffffff8111156101165761011661039c565b6040519080825280601f01601f191660200182016040528015610140576020820181803683370190505b509050600360fc1b8160008151811061015b5761015b6103b2565b60200101906001600160f81b031916908160001a905350600f60fb1b8160018151811061018a5761018a6103b2565b60200101906001600160f81b031916908160001a90535060006101ae846002610365565b6101b9906001610384565b90505b6001811115610231576f181899199a1a9b1b9c1cb0b131b232b360811b85600f16601081106101ed576101ed6103b2565b1a60f81b828281518110610203576102036103b2565b60200101906001600160f81b031916908160001a90535060049490941c9361022a816103c8565b90506101bc565b5083156102805760405162461bcd60e51b815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e7460448201526064016100d8565b9392505050565b60005b838110156102a257818101518382015260200161028a565b838111156102b1576000848401525b50505050565b6e032bc3832b1ba32b229b2b73232b91608d1b8152600083516102e181600f850160208801610287565b6e01030b1ba3ab0b61039b2b73232b91608d1b600f91840191820152835161031081601e840160208801610287565b01601e01949350505050565b602081526000825180602084015261033b816040850160208701610287565b601f01601f19169190910160400192915050565b634e487b7160e01b600052601160045260246000fd5b600081600019048311821515161561037f5761037f61034f565b500290565b600082198211156103975761039761034f565b500190565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b6000816103d7576103d761034f565b50600019019056fea264697066735822122021f68a56920052d98a209122db94525247e8b3b43872363d5e85c3529943feed64736f6c634300080c0033",
}

// DelegatecallReceiverABI is the input ABI used to generate the binding from.
// Deprecated: Use DelegatecallReceiverMetaData.ABI instead.
var DelegatecallReceiverABI = DelegatecallReceiverMetaData.ABI

// DelegatecallReceiverBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DelegatecallReceiverMetaData.Bin instead.
var DelegatecallReceiverBin = DelegatecallReceiverMetaData.Bin

// DeployDelegatecallReceiver deploys a new Ethereum contract, binding an instance of DelegatecallReceiver to it.
func DeployDelegatecallReceiver(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DelegatecallReceiver, error) {
	parsed, err := DelegatecallReceiverMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DelegatecallReceiverBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DelegatecallReceiver{DelegatecallReceiverCaller: DelegatecallReceiverCaller{contract: contract}, DelegatecallReceiverTransactor: DelegatecallReceiverTransactor{contract: contract}, DelegatecallReceiverFilterer: DelegatecallReceiverFilterer{contract: contract}}, nil
}

// DelegatecallReceiver is an auto generated Go binding around an Ethereum contract.
type DelegatecallReceiver struct {
	DelegatecallReceiverCaller     // Read-only binding to the contract
	DelegatecallReceiverTransactor // Write-only binding to the contract
	DelegatecallReceiverFilterer   // Log filterer for contract events
}

// DelegatecallReceiverCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegatecallReceiverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatecallReceiverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegatecallReceiverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatecallReceiverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegatecallReceiverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegatecallReceiverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegatecallReceiverSession struct {
	Contract     *DelegatecallReceiver // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// DelegatecallReceiverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegatecallReceiverCallerSession struct {
	Contract *DelegatecallReceiverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// DelegatecallReceiverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegatecallReceiverTransactorSession struct {
	Contract     *DelegatecallReceiverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// DelegatecallReceiverRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegatecallReceiverRaw struct {
	Contract *DelegatecallReceiver // Generic contract binding to access the raw methods on
}

// DelegatecallReceiverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegatecallReceiverCallerRaw struct {
	Contract *DelegatecallReceiverCaller // Generic read-only contract binding to access the raw methods on
}

// DelegatecallReceiverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegatecallReceiverTransactorRaw struct {
	Contract *DelegatecallReceiverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegatecallReceiver creates a new instance of DelegatecallReceiver, bound to a specific deployed contract.
func NewDelegatecallReceiver(address common.Address, backend bind.ContractBackend) (*DelegatecallReceiver, error) {
	contract, err := bindDelegatecallReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DelegatecallReceiver{DelegatecallReceiverCaller: DelegatecallReceiverCaller{contract: contract}, DelegatecallReceiverTransactor: DelegatecallReceiverTransactor{contract: contract}, DelegatecallReceiverFilterer: DelegatecallReceiverFilterer{contract: contract}}, nil
}

// NewDelegatecallReceiverCaller creates a new read-only instance of DelegatecallReceiver, bound to a specific deployed contract.
func NewDelegatecallReceiverCaller(address common.Address, caller bind.ContractCaller) (*DelegatecallReceiverCaller, error) {
	contract, err := bindDelegatecallReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegatecallReceiverCaller{contract: contract}, nil
}

// NewDelegatecallReceiverTransactor creates a new write-only instance of DelegatecallReceiver, bound to a specific deployed contract.
func NewDelegatecallReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegatecallReceiverTransactor, error) {
	contract, err := bindDelegatecallReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegatecallReceiverTransactor{contract: contract}, nil
}

// NewDelegatecallReceiverFilterer creates a new log filterer instance of DelegatecallReceiver, bound to a specific deployed contract.
func NewDelegatecallReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegatecallReceiverFilterer, error) {
	contract, err := bindDelegatecallReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegatecallReceiverFilterer{contract: contract}, nil
}

// bindDelegatecallReceiver binds a generic wrapper to an already deployed contract.
func bindDelegatecallReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DelegatecallReceiverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelegatecallReceiver *DelegatecallReceiverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DelegatecallReceiver.Contract.DelegatecallReceiverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelegatecallReceiver *DelegatecallReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegatecallReceiver.Contract.DelegatecallReceiverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelegatecallReceiver *DelegatecallReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelegatecallReceiver.Contract.DelegatecallReceiverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DelegatecallReceiver *DelegatecallReceiverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DelegatecallReceiver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DelegatecallReceiver *DelegatecallReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DelegatecallReceiver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DelegatecallReceiver *DelegatecallReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DelegatecallReceiver.Contract.contract.Transact(opts, method, params...)
}

// Entrypoint is a free data retrieval call binding the contract method 0xa65d69d4.
//
// Solidity: function entrypoint() view returns()
func (_DelegatecallReceiver *DelegatecallReceiverCaller) Entrypoint(opts *bind.CallOpts) error {
	var out []interface{}
	err := _DelegatecallReceiver.contract.Call(opts, &out, "entrypoint")

	if err != nil {
		return err
	}

	return err

}

// Entrypoint is a free data retrieval call binding the contract method 0xa65d69d4.
//
// Solidity: function entrypoint() view returns()
func (_DelegatecallReceiver *DelegatecallReceiverSession) Entrypoint() error {
	return _DelegatecallReceiver.Contract.Entrypoint(&_DelegatecallReceiver.CallOpts)
}

// Entrypoint is a free data retrieval call binding the contract method 0xa65d69d4.
//
// Solidity: function entrypoint() view returns()
func (_DelegatecallReceiver *DelegatecallReceiverCallerSession) Entrypoint() error {
	return _DelegatecallReceiver.Contract.Entrypoint(&_DelegatecallReceiver.CallOpts)
}

// ExpectedSender is a free data retrieval call binding the contract method 0x1754dba5.
//
// Solidity: function expectedSender() view returns(address)
func (_DelegatecallReceiver *DelegatecallReceiverCaller) ExpectedSender(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DelegatecallReceiver.contract.Call(opts, &out, "expectedSender")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ExpectedSender is a free data retrieval call binding the contract method 0x1754dba5.
//
// Solidity: function expectedSender() view returns(address)
func (_DelegatecallReceiver *DelegatecallReceiverSession) ExpectedSender() (common.Address, error) {
	return _DelegatecallReceiver.Contract.ExpectedSender(&_DelegatecallReceiver.CallOpts)
}

// ExpectedSender is a free data retrieval call binding the contract method 0x1754dba5.
//
// Solidity: function expectedSender() view returns(address)
func (_DelegatecallReceiver *DelegatecallReceiverCallerSession) ExpectedSender() (common.Address, error) {
	return _DelegatecallReceiver.Contract.ExpectedSender(&_DelegatecallReceiver.CallOpts)
}
