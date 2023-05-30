// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ChainCallLevel2

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

// ChainCallLevel2MetaData contains all meta data concerning the ChainCallLevel2 contract.
var ChainCallLevel2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level3Addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"callRevert\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level3Addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"delegateCallRevert\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level3Addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"exec\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level3Addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"get\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"t\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061079a806100206000396000f3fe60806040526004361061003f5760003560e01c80634352429814610044578063d81e842314610059578063ee2d01151461008f578063ffb8d165146100a2575b600080fd5b6100576100523660046105ae565b6100b5565b005b34801561006557600080fd5b506100796100743660046105ae565b6101b7565b6040516100869190610611565b60405180910390f35b61005761009d3660046105ae565b6103d6565b6100576100b03660046105ae565b6104db565b6040516001600160a01b0382811660248301526000919084169060440160408051601f198184030181529181526020820180516001600160e01b03166324df79d560e21b179052516101079190610644565b600060405180830381855af49150503d8060008114610142576040519150601f19603f3d011682016040523d82523d6000602084013e610147565b606091505b505080915050806101b25760405162461bcd60e51b815260206004820152602a60248201527f6661696c656420746f20706572666f726d2064656c65676174652063616c6c20604482015269746f206c6576656c203360b01b60648201526084015b60405180910390fd5b505050565b6040516001600160a01b03828116602483015260609160009183919086169060440160408051601f198184030181529181526020820180516001600160e01b03166330af0bbf60e21b1790525161020e9190610644565b600060405180830381855afa9150503d8060008114610249576040519150601f19603f3d011682016040523d82523d6000602084013e61024e565b606091505b509092509050816102b25760405162461bcd60e51b815260206004820152602860248201527f6661696c656420746f20706572666f726d207374617469632063616c6c20746f604482015267206c6576656c203360c01b60648201526084016101a9565b808060200190518101906102c69190610676565b60408051600481526024810182526020810180516001600160e01b0316631b53398f60e21b17905290519194506001600160a01b038616916103089190610644565b600060405180830381855afa9150503d8060008114610343576040519150601f19603f3d011682016040523d82523d6000602084013e610348565b606091505b509092509050816103b95760405162461bcd60e51b815260206004820152603560248201527f6661696c656420746f20706572666f726d207374617469632063616c6c20746f604482015274103632bb32b6101a10333937b6903632bb32b6101960591b60648201526084016101a9565b808060200190518101906103cd9190610676565b95945050505050565b6040516001600160a01b0382811660248301526000919084169060440160408051601f198184030181529181526020820180516001600160e01b03166335db093760e11b179052516104289190610644565b6000604051808303816000865af19150503d8060008114610465576040519150601f19603f3d011682016040523d82523d6000602084013e61046a565b606091505b5050809150508061048d5760405162461bcd60e51b81526004016101a990610723565b6040516001600160a01b03838116602483015284169060440160408051601f198184030181529181526020820180516001600160e01b03166335db093760e11b179052516101079190610644565b6040516001600160a01b0382811660248301526000919084169060440160408051601f198184030181529181526020820180516001600160e01b0316636f852eab60e11b1790525161052d9190610644565b6000604051808303816000865af19150503d806000811461056a576040519150601f19603f3d011682016040523d82523d6000602084013e61056f565b606091505b505080915050806101b25760405162461bcd60e51b81526004016101a990610723565b80356001600160a01b03811681146105a957600080fd5b919050565b600080604083850312156105c157600080fd5b6105ca83610592565b91506105d860208401610592565b90509250929050565b60005b838110156105fc5781810151838201526020016105e4565b8381111561060b576000848401525b50505050565b60208152600082518060208401526106308160408501602087016105e1565b601f01601f19169190910160400192915050565b600082516106568184602087016105e1565b9190910192915050565b634e487b7160e01b600052604160045260246000fd5b60006020828403121561068857600080fd5b815167ffffffffffffffff808211156106a057600080fd5b818401915084601f8301126106b457600080fd5b8151818111156106c6576106c6610660565b604051601f8201601f19908116603f011681019083821181831017156106ee576106ee610660565b8160405282815287602084870101111561070757600080fd5b6107188360208301602088016105e1565b979650505050505050565b60208082526021908201527f6661696c656420746f20706572666f726d2063616c6c20746f206c6576656c206040820152603360f81b60608201526080019056fea264697066735822122011f1fa7d58307c11cefe4347c017092fdb121f066b4d2186bf3ee7db30a85e8664736f6c634300080c0033",
}

// ChainCallLevel2ABI is the input ABI used to generate the binding from.
// Deprecated: Use ChainCallLevel2MetaData.ABI instead.
var ChainCallLevel2ABI = ChainCallLevel2MetaData.ABI

// ChainCallLevel2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ChainCallLevel2MetaData.Bin instead.
var ChainCallLevel2Bin = ChainCallLevel2MetaData.Bin

// DeployChainCallLevel2 deploys a new Ethereum contract, binding an instance of ChainCallLevel2 to it.
func DeployChainCallLevel2(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ChainCallLevel2, error) {
	parsed, err := ChainCallLevel2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChainCallLevel2Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChainCallLevel2{ChainCallLevel2Caller: ChainCallLevel2Caller{contract: contract}, ChainCallLevel2Transactor: ChainCallLevel2Transactor{contract: contract}, ChainCallLevel2Filterer: ChainCallLevel2Filterer{contract: contract}}, nil
}

// ChainCallLevel2 is an auto generated Go binding around an Ethereum contract.
type ChainCallLevel2 struct {
	ChainCallLevel2Caller     // Read-only binding to the contract
	ChainCallLevel2Transactor // Write-only binding to the contract
	ChainCallLevel2Filterer   // Log filterer for contract events
}

// ChainCallLevel2Caller is an auto generated read-only Go binding around an Ethereum contract.
type ChainCallLevel2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainCallLevel2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ChainCallLevel2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainCallLevel2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChainCallLevel2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainCallLevel2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChainCallLevel2Session struct {
	Contract     *ChainCallLevel2  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChainCallLevel2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChainCallLevel2CallerSession struct {
	Contract *ChainCallLevel2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ChainCallLevel2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChainCallLevel2TransactorSession struct {
	Contract     *ChainCallLevel2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ChainCallLevel2Raw is an auto generated low-level Go binding around an Ethereum contract.
type ChainCallLevel2Raw struct {
	Contract *ChainCallLevel2 // Generic contract binding to access the raw methods on
}

// ChainCallLevel2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChainCallLevel2CallerRaw struct {
	Contract *ChainCallLevel2Caller // Generic read-only contract binding to access the raw methods on
}

// ChainCallLevel2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChainCallLevel2TransactorRaw struct {
	Contract *ChainCallLevel2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewChainCallLevel2 creates a new instance of ChainCallLevel2, bound to a specific deployed contract.
func NewChainCallLevel2(address common.Address, backend bind.ContractBackend) (*ChainCallLevel2, error) {
	contract, err := bindChainCallLevel2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel2{ChainCallLevel2Caller: ChainCallLevel2Caller{contract: contract}, ChainCallLevel2Transactor: ChainCallLevel2Transactor{contract: contract}, ChainCallLevel2Filterer: ChainCallLevel2Filterer{contract: contract}}, nil
}

// NewChainCallLevel2Caller creates a new read-only instance of ChainCallLevel2, bound to a specific deployed contract.
func NewChainCallLevel2Caller(address common.Address, caller bind.ContractCaller) (*ChainCallLevel2Caller, error) {
	contract, err := bindChainCallLevel2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel2Caller{contract: contract}, nil
}

// NewChainCallLevel2Transactor creates a new write-only instance of ChainCallLevel2, bound to a specific deployed contract.
func NewChainCallLevel2Transactor(address common.Address, transactor bind.ContractTransactor) (*ChainCallLevel2Transactor, error) {
	contract, err := bindChainCallLevel2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel2Transactor{contract: contract}, nil
}

// NewChainCallLevel2Filterer creates a new log filterer instance of ChainCallLevel2, bound to a specific deployed contract.
func NewChainCallLevel2Filterer(address common.Address, filterer bind.ContractFilterer) (*ChainCallLevel2Filterer, error) {
	contract, err := bindChainCallLevel2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel2Filterer{contract: contract}, nil
}

// bindChainCallLevel2 binds a generic wrapper to an already deployed contract.
func bindChainCallLevel2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChainCallLevel2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainCallLevel2 *ChainCallLevel2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainCallLevel2.Contract.ChainCallLevel2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainCallLevel2 *ChainCallLevel2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainCallLevel2.Contract.ChainCallLevel2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainCallLevel2 *ChainCallLevel2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainCallLevel2.Contract.ChainCallLevel2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainCallLevel2 *ChainCallLevel2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainCallLevel2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainCallLevel2 *ChainCallLevel2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainCallLevel2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainCallLevel2 *ChainCallLevel2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainCallLevel2.Contract.contract.Transact(opts, method, params...)
}

// Get is a free data retrieval call binding the contract method 0xd81e8423.
//
// Solidity: function get(address level3Addr, address level4Addr) view returns(string t)
func (_ChainCallLevel2 *ChainCallLevel2Caller) Get(opts *bind.CallOpts, level3Addr common.Address, level4Addr common.Address) (string, error) {
	var out []interface{}
	err := _ChainCallLevel2.contract.Call(opts, &out, "get", level3Addr, level4Addr)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Get is a free data retrieval call binding the contract method 0xd81e8423.
//
// Solidity: function get(address level3Addr, address level4Addr) view returns(string t)
func (_ChainCallLevel2 *ChainCallLevel2Session) Get(level3Addr common.Address, level4Addr common.Address) (string, error) {
	return _ChainCallLevel2.Contract.Get(&_ChainCallLevel2.CallOpts, level3Addr, level4Addr)
}

// Get is a free data retrieval call binding the contract method 0xd81e8423.
//
// Solidity: function get(address level3Addr, address level4Addr) view returns(string t)
func (_ChainCallLevel2 *ChainCallLevel2CallerSession) Get(level3Addr common.Address, level4Addr common.Address) (string, error) {
	return _ChainCallLevel2.Contract.Get(&_ChainCallLevel2.CallOpts, level3Addr, level4Addr)
}

// CallRevert is a paid mutator transaction binding the contract method 0xffb8d165.
//
// Solidity: function callRevert(address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel2 *ChainCallLevel2Transactor) CallRevert(opts *bind.TransactOpts, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel2.contract.Transact(opts, "callRevert", level3Addr, level4Addr)
}

// CallRevert is a paid mutator transaction binding the contract method 0xffb8d165.
//
// Solidity: function callRevert(address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel2 *ChainCallLevel2Session) CallRevert(level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel2.Contract.CallRevert(&_ChainCallLevel2.TransactOpts, level3Addr, level4Addr)
}

// CallRevert is a paid mutator transaction binding the contract method 0xffb8d165.
//
// Solidity: function callRevert(address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel2 *ChainCallLevel2TransactorSession) CallRevert(level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel2.Contract.CallRevert(&_ChainCallLevel2.TransactOpts, level3Addr, level4Addr)
}

// DelegateCallRevert is a paid mutator transaction binding the contract method 0x43524298.
//
// Solidity: function delegateCallRevert(address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel2 *ChainCallLevel2Transactor) DelegateCallRevert(opts *bind.TransactOpts, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel2.contract.Transact(opts, "delegateCallRevert", level3Addr, level4Addr)
}

// DelegateCallRevert is a paid mutator transaction binding the contract method 0x43524298.
//
// Solidity: function delegateCallRevert(address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel2 *ChainCallLevel2Session) DelegateCallRevert(level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel2.Contract.DelegateCallRevert(&_ChainCallLevel2.TransactOpts, level3Addr, level4Addr)
}

// DelegateCallRevert is a paid mutator transaction binding the contract method 0x43524298.
//
// Solidity: function delegateCallRevert(address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel2 *ChainCallLevel2TransactorSession) DelegateCallRevert(level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel2.Contract.DelegateCallRevert(&_ChainCallLevel2.TransactOpts, level3Addr, level4Addr)
}

// Exec is a paid mutator transaction binding the contract method 0xee2d0115.
//
// Solidity: function exec(address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel2 *ChainCallLevel2Transactor) Exec(opts *bind.TransactOpts, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel2.contract.Transact(opts, "exec", level3Addr, level4Addr)
}

// Exec is a paid mutator transaction binding the contract method 0xee2d0115.
//
// Solidity: function exec(address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel2 *ChainCallLevel2Session) Exec(level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel2.Contract.Exec(&_ChainCallLevel2.TransactOpts, level3Addr, level4Addr)
}

// Exec is a paid mutator transaction binding the contract method 0xee2d0115.
//
// Solidity: function exec(address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel2 *ChainCallLevel2TransactorSession) Exec(level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel2.Contract.Exec(&_ChainCallLevel2.TransactOpts, level3Addr, level4Addr)
}
