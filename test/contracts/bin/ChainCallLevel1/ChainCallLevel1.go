// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ChainCallLevel1

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

// ChainCallLevel1MetaData contains all meta data concerning the ChainCallLevel1 contract.
var ChainCallLevel1MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level2Addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"level3Addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"callRevert\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level2Addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"level3Addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"delegateCallRevert\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level2Addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"level3Addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"delegateTransfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"level2Addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"level3Addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"level4Addr\",\"type\":\"address\"}],\"name\":\"exec\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061073d806100206000396000f3fe6080604052600436106100435760003560e01c80630c0cc7d21461004f5780636949127414610064578063b04b2e2714610077578063c023dcf31461008a57600080fd5b3661004a57005b600080fd5b61006261005d36600461052e565b61009d565b005b61006261007236600461052e565b61016e565b61006261008536600461052e565b610223565b61006261009836600461052e565b61027d565b6040516001600160a01b03838116602483015282811660448301526000919085169060640160408051601f198184030181529181526020820180516001600160e01b0316600162472e9b60e01b0319179052516100fa919061059d565b6000604051808303816000865af19150503d8060008114610137576040519150601f19603f3d011682016040523d82523d6000602084013e61013c565b606091505b505080915050806101685760405162461bcd60e51b815260040161015f906105b9565b60405180910390fd5b50505050565b6040516001600160a01b0383811660248301526000919085169060440160408051601f198184030181529181526020820180516001600160e01b0316635effa1f760e11b179052516101c0919061059d565b600060405180830381855af49150503d80600081146101fb576040519150601f19603f3d011682016040523d82523d6000602084013e610200565b606091505b505080915050806101685760405162461bcd60e51b815260040161015f906105fa565b6040516001600160a01b03838116602483015282811660448301526000919085169060640160408051601f198184030181529181526020820180516001600160e01b031663086a485360e31b179052516101c0919061059d565b6040516001600160a01b03838116602483015282811660448301526000919085169060640160408051601f198184030181529181526020820180516001600160e01b031663ee2d011560e01b179052516102d7919061059d565b6000604051808303816000865af19150503d8060008114610314576040519150601f19603f3d011682016040523d82523d6000602084013e610319565b606091505b5050809150508061033c5760405162461bcd60e51b815260040161015f906105b9565b6040516001600160a01b038481166024830152838116604483015285169060640160408051601f198184030181529181526020820180516001600160e01b031663ee2d011560e01b17905251610392919061059d565b600060405180830381855af49150503d80600081146103cd576040519150601f19603f3d011682016040523d82523d6000602084013e6103d2565b606091505b505080915050806103f55760405162461bcd60e51b815260040161015f906105fa565b6040516001600160a01b03848116602483015283811660448301526060919086169060640160408051601f198184030181529181526020820180516001600160e01b031663d81e842360e01b1790525161044f919061059d565b600060405180830381855afa9150503d806000811461048a576040519150601f19603f3d011682016040523d82523d6000602084013e61048f565b606091505b509092509050816104f35760405162461bcd60e51b815260206004820152602860248201527f6661696c656420746f20706572666f726d207374617469632063616c6c20746f604482015267103632bb32b6101960c11b606482015260840161015f565b606081806020019051810190610509919061065a565b50505050505050565b80356001600160a01b038116811461052957600080fd5b919050565b60008060006060848603121561054357600080fd5b61054c84610512565b925061055a60208501610512565b915061056860408501610512565b90509250925092565b60005b8381101561058c578181015183820152602001610574565b838111156101685750506000910152565b600082516105af818460208701610571565b9190910192915050565b60208082526021908201527f6661696c656420746f20706572666f726d2063616c6c20746f206c6576656c206040820152601960f91b606082015260800190565b6020808252602a908201527f6661696c656420746f20706572666f726d2064656c65676174652063616c6c206040820152693a37903632bb32b6101960b11b606082015260800190565b634e487b7160e01b600052604160045260246000fd5b60006020828403121561066c57600080fd5b815167ffffffffffffffff8082111561068457600080fd5b818401915084601f83011261069857600080fd5b8151818111156106aa576106aa610644565b604051601f8201601f19908116603f011681019083821181831017156106d2576106d2610644565b816040528281528760208487010111156106eb57600080fd5b6106fc836020830160208801610571565b97965050505050505056fea2646970667358221220b59f48fa8573e9cdbbd8c29f3523b5c3f2f9c3c0cb8877d698821da35114e3c864736f6c634300080c0033",
}

// ChainCallLevel1ABI is the input ABI used to generate the binding from.
// Deprecated: Use ChainCallLevel1MetaData.ABI instead.
var ChainCallLevel1ABI = ChainCallLevel1MetaData.ABI

// ChainCallLevel1Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ChainCallLevel1MetaData.Bin instead.
var ChainCallLevel1Bin = ChainCallLevel1MetaData.Bin

// DeployChainCallLevel1 deploys a new Ethereum contract, binding an instance of ChainCallLevel1 to it.
func DeployChainCallLevel1(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ChainCallLevel1, error) {
	parsed, err := ChainCallLevel1MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ChainCallLevel1Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ChainCallLevel1{ChainCallLevel1Caller: ChainCallLevel1Caller{contract: contract}, ChainCallLevel1Transactor: ChainCallLevel1Transactor{contract: contract}, ChainCallLevel1Filterer: ChainCallLevel1Filterer{contract: contract}}, nil
}

// ChainCallLevel1 is an auto generated Go binding around an Ethereum contract.
type ChainCallLevel1 struct {
	ChainCallLevel1Caller     // Read-only binding to the contract
	ChainCallLevel1Transactor // Write-only binding to the contract
	ChainCallLevel1Filterer   // Log filterer for contract events
}

// ChainCallLevel1Caller is an auto generated read-only Go binding around an Ethereum contract.
type ChainCallLevel1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainCallLevel1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ChainCallLevel1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainCallLevel1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChainCallLevel1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainCallLevel1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChainCallLevel1Session struct {
	Contract     *ChainCallLevel1  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChainCallLevel1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChainCallLevel1CallerSession struct {
	Contract *ChainCallLevel1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ChainCallLevel1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChainCallLevel1TransactorSession struct {
	Contract     *ChainCallLevel1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ChainCallLevel1Raw is an auto generated low-level Go binding around an Ethereum contract.
type ChainCallLevel1Raw struct {
	Contract *ChainCallLevel1 // Generic contract binding to access the raw methods on
}

// ChainCallLevel1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChainCallLevel1CallerRaw struct {
	Contract *ChainCallLevel1Caller // Generic read-only contract binding to access the raw methods on
}

// ChainCallLevel1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChainCallLevel1TransactorRaw struct {
	Contract *ChainCallLevel1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewChainCallLevel1 creates a new instance of ChainCallLevel1, bound to a specific deployed contract.
func NewChainCallLevel1(address common.Address, backend bind.ContractBackend) (*ChainCallLevel1, error) {
	contract, err := bindChainCallLevel1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel1{ChainCallLevel1Caller: ChainCallLevel1Caller{contract: contract}, ChainCallLevel1Transactor: ChainCallLevel1Transactor{contract: contract}, ChainCallLevel1Filterer: ChainCallLevel1Filterer{contract: contract}}, nil
}

// NewChainCallLevel1Caller creates a new read-only instance of ChainCallLevel1, bound to a specific deployed contract.
func NewChainCallLevel1Caller(address common.Address, caller bind.ContractCaller) (*ChainCallLevel1Caller, error) {
	contract, err := bindChainCallLevel1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel1Caller{contract: contract}, nil
}

// NewChainCallLevel1Transactor creates a new write-only instance of ChainCallLevel1, bound to a specific deployed contract.
func NewChainCallLevel1Transactor(address common.Address, transactor bind.ContractTransactor) (*ChainCallLevel1Transactor, error) {
	contract, err := bindChainCallLevel1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel1Transactor{contract: contract}, nil
}

// NewChainCallLevel1Filterer creates a new log filterer instance of ChainCallLevel1, bound to a specific deployed contract.
func NewChainCallLevel1Filterer(address common.Address, filterer bind.ContractFilterer) (*ChainCallLevel1Filterer, error) {
	contract, err := bindChainCallLevel1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainCallLevel1Filterer{contract: contract}, nil
}

// bindChainCallLevel1 binds a generic wrapper to an already deployed contract.
func bindChainCallLevel1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChainCallLevel1MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainCallLevel1 *ChainCallLevel1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainCallLevel1.Contract.ChainCallLevel1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainCallLevel1 *ChainCallLevel1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.ChainCallLevel1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainCallLevel1 *ChainCallLevel1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.ChainCallLevel1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainCallLevel1 *ChainCallLevel1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainCallLevel1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainCallLevel1 *ChainCallLevel1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainCallLevel1 *ChainCallLevel1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.contract.Transact(opts, method, params...)
}

// CallRevert is a paid mutator transaction binding the contract method 0x0c0cc7d2.
//
// Solidity: function callRevert(address level2Addr, address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel1 *ChainCallLevel1Transactor) CallRevert(opts *bind.TransactOpts, level2Addr common.Address, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel1.contract.Transact(opts, "callRevert", level2Addr, level3Addr, level4Addr)
}

// CallRevert is a paid mutator transaction binding the contract method 0x0c0cc7d2.
//
// Solidity: function callRevert(address level2Addr, address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel1 *ChainCallLevel1Session) CallRevert(level2Addr common.Address, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.CallRevert(&_ChainCallLevel1.TransactOpts, level2Addr, level3Addr, level4Addr)
}

// CallRevert is a paid mutator transaction binding the contract method 0x0c0cc7d2.
//
// Solidity: function callRevert(address level2Addr, address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel1 *ChainCallLevel1TransactorSession) CallRevert(level2Addr common.Address, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.CallRevert(&_ChainCallLevel1.TransactOpts, level2Addr, level3Addr, level4Addr)
}

// DelegateCallRevert is a paid mutator transaction binding the contract method 0xb04b2e27.
//
// Solidity: function delegateCallRevert(address level2Addr, address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel1 *ChainCallLevel1Transactor) DelegateCallRevert(opts *bind.TransactOpts, level2Addr common.Address, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel1.contract.Transact(opts, "delegateCallRevert", level2Addr, level3Addr, level4Addr)
}

// DelegateCallRevert is a paid mutator transaction binding the contract method 0xb04b2e27.
//
// Solidity: function delegateCallRevert(address level2Addr, address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel1 *ChainCallLevel1Session) DelegateCallRevert(level2Addr common.Address, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.DelegateCallRevert(&_ChainCallLevel1.TransactOpts, level2Addr, level3Addr, level4Addr)
}

// DelegateCallRevert is a paid mutator transaction binding the contract method 0xb04b2e27.
//
// Solidity: function delegateCallRevert(address level2Addr, address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel1 *ChainCallLevel1TransactorSession) DelegateCallRevert(level2Addr common.Address, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.DelegateCallRevert(&_ChainCallLevel1.TransactOpts, level2Addr, level3Addr, level4Addr)
}

// DelegateTransfer is a paid mutator transaction binding the contract method 0x69491274.
//
// Solidity: function delegateTransfer(address level2Addr, address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel1 *ChainCallLevel1Transactor) DelegateTransfer(opts *bind.TransactOpts, level2Addr common.Address, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel1.contract.Transact(opts, "delegateTransfer", level2Addr, level3Addr, level4Addr)
}

// DelegateTransfer is a paid mutator transaction binding the contract method 0x69491274.
//
// Solidity: function delegateTransfer(address level2Addr, address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel1 *ChainCallLevel1Session) DelegateTransfer(level2Addr common.Address, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.DelegateTransfer(&_ChainCallLevel1.TransactOpts, level2Addr, level3Addr, level4Addr)
}

// DelegateTransfer is a paid mutator transaction binding the contract method 0x69491274.
//
// Solidity: function delegateTransfer(address level2Addr, address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel1 *ChainCallLevel1TransactorSession) DelegateTransfer(level2Addr common.Address, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.DelegateTransfer(&_ChainCallLevel1.TransactOpts, level2Addr, level3Addr, level4Addr)
}

// Exec is a paid mutator transaction binding the contract method 0xc023dcf3.
//
// Solidity: function exec(address level2Addr, address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel1 *ChainCallLevel1Transactor) Exec(opts *bind.TransactOpts, level2Addr common.Address, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel1.contract.Transact(opts, "exec", level2Addr, level3Addr, level4Addr)
}

// Exec is a paid mutator transaction binding the contract method 0xc023dcf3.
//
// Solidity: function exec(address level2Addr, address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel1 *ChainCallLevel1Session) Exec(level2Addr common.Address, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.Exec(&_ChainCallLevel1.TransactOpts, level2Addr, level3Addr, level4Addr)
}

// Exec is a paid mutator transaction binding the contract method 0xc023dcf3.
//
// Solidity: function exec(address level2Addr, address level3Addr, address level4Addr) payable returns()
func (_ChainCallLevel1 *ChainCallLevel1TransactorSession) Exec(level2Addr common.Address, level3Addr common.Address, level4Addr common.Address) (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.Exec(&_ChainCallLevel1.TransactOpts, level2Addr, level3Addr, level4Addr)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ChainCallLevel1 *ChainCallLevel1Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainCallLevel1.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ChainCallLevel1 *ChainCallLevel1Session) Receive() (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.Receive(&_ChainCallLevel1.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ChainCallLevel1 *ChainCallLevel1TransactorSession) Receive() (*types.Transaction, error) {
	return _ChainCallLevel1.Contract.Receive(&_ChainCallLevel1.TransactOpts)
}
