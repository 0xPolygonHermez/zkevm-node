// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Creates

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

// CreatesMetaData contains all meta data concerning the Creates contract.
var CreatesMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"}],\"name\":\"add\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"bytecode\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"opCreate\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"bytecode\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"opCreate2\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"bytecode\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"opCreate2Complex\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"bytecode\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"opCreate2Value\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"bytecode\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"opCreateValue\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sendValue\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610369806100206000396000f3fe6080604052600436106100705760003560e01c8063771602f71161004e578063771602f7146100d0578063b88c4aa9146100fe578063c935aee414610111578063e3306a251461015057600080fd5b806327c845dc146100755780633c77eba3146100805780635b8e9959146100b0575b600080fd5b61007e34600155565b005b61009361008e366004610236565b610170565b6040516001600160a01b0390911681526020015b60405180910390f35b3480156100bc57600080fd5b506100936100cb366004610236565b610187565b3480156100dc57600080fd5b506100f06100eb3660046102eb565b61019d565b6040519081526020016100a7565b61009361010c366004610236565b6101b0565b34801561011d57600080fd5b5061013161012c366004610236565b6101cb565b604080516001600160a01b0390931683526020830191909152016100a7565b34801561015c57600080fd5b5061009361016b366004610236565b610208565b6000808260a06101f4f06000819055949350505050565b6000808260a06000f06000819055949350505050565b60006101a9828461030d565b9392505050565b600080620555558360a061012cf56000819055949350505050565b60008060006101dc6001600261019d565b90506000600285602088016000f59050806000556101fc6002600461019d565b90969095509350505050565b60008060028360a06000f56000819055949350505050565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561024957600080fd5b823567ffffffffffffffff8082111561026157600080fd5b818501915085601f83011261027557600080fd5b81358181111561028757610287610220565b604051601f8201601f19908116603f011681019083821181831017156102af576102af610220565b816040528281528860208487010111156102c857600080fd5b826020860160208301376000602093820184015298969091013596505050505050565b600080604083850312156102fe57600080fd5b50508035926020909101359150565b6000821982111561032e57634e487b7160e01b600052601160045260246000fd5b50019056fea2646970667358221220cdc0e0bdc2487139b3aa0666f32d3f0ed1e40a81659b28e6dea427224cc6104f64736f6c634300080c0033",
}

// CreatesABI is the input ABI used to generate the binding from.
// Deprecated: Use CreatesMetaData.ABI instead.
var CreatesABI = CreatesMetaData.ABI

// CreatesBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CreatesMetaData.Bin instead.
var CreatesBin = CreatesMetaData.Bin

// DeployCreates deploys a new Ethereum contract, binding an instance of Creates to it.
func DeployCreates(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Creates, error) {
	parsed, err := CreatesMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CreatesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Creates{CreatesCaller: CreatesCaller{contract: contract}, CreatesTransactor: CreatesTransactor{contract: contract}, CreatesFilterer: CreatesFilterer{contract: contract}}, nil
}

// Creates is an auto generated Go binding around an Ethereum contract.
type Creates struct {
	CreatesCaller     // Read-only binding to the contract
	CreatesTransactor // Write-only binding to the contract
	CreatesFilterer   // Log filterer for contract events
}

// CreatesCaller is an auto generated read-only Go binding around an Ethereum contract.
type CreatesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CreatesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CreatesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CreatesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CreatesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CreatesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CreatesSession struct {
	Contract     *Creates          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CreatesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CreatesCallerSession struct {
	Contract *CreatesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// CreatesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CreatesTransactorSession struct {
	Contract     *CreatesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// CreatesRaw is an auto generated low-level Go binding around an Ethereum contract.
type CreatesRaw struct {
	Contract *Creates // Generic contract binding to access the raw methods on
}

// CreatesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CreatesCallerRaw struct {
	Contract *CreatesCaller // Generic read-only contract binding to access the raw methods on
}

// CreatesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CreatesTransactorRaw struct {
	Contract *CreatesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCreates creates a new instance of Creates, bound to a specific deployed contract.
func NewCreates(address common.Address, backend bind.ContractBackend) (*Creates, error) {
	contract, err := bindCreates(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Creates{CreatesCaller: CreatesCaller{contract: contract}, CreatesTransactor: CreatesTransactor{contract: contract}, CreatesFilterer: CreatesFilterer{contract: contract}}, nil
}

// NewCreatesCaller creates a new read-only instance of Creates, bound to a specific deployed contract.
func NewCreatesCaller(address common.Address, caller bind.ContractCaller) (*CreatesCaller, error) {
	contract, err := bindCreates(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CreatesCaller{contract: contract}, nil
}

// NewCreatesTransactor creates a new write-only instance of Creates, bound to a specific deployed contract.
func NewCreatesTransactor(address common.Address, transactor bind.ContractTransactor) (*CreatesTransactor, error) {
	contract, err := bindCreates(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CreatesTransactor{contract: contract}, nil
}

// NewCreatesFilterer creates a new log filterer instance of Creates, bound to a specific deployed contract.
func NewCreatesFilterer(address common.Address, filterer bind.ContractFilterer) (*CreatesFilterer, error) {
	contract, err := bindCreates(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CreatesFilterer{contract: contract}, nil
}

// bindCreates binds a generic wrapper to an already deployed contract.
func bindCreates(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CreatesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Creates *CreatesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Creates.Contract.CreatesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Creates *CreatesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Creates.Contract.CreatesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Creates *CreatesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Creates.Contract.CreatesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Creates *CreatesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Creates.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Creates *CreatesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Creates.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Creates *CreatesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Creates.Contract.contract.Transact(opts, method, params...)
}

// Add is a free data retrieval call binding the contract method 0x771602f7.
//
// Solidity: function add(uint256 a, uint256 b) pure returns(uint256)
func (_Creates *CreatesCaller) Add(opts *bind.CallOpts, a *big.Int, b *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Creates.contract.Call(opts, &out, "add", a, b)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Add is a free data retrieval call binding the contract method 0x771602f7.
//
// Solidity: function add(uint256 a, uint256 b) pure returns(uint256)
func (_Creates *CreatesSession) Add(a *big.Int, b *big.Int) (*big.Int, error) {
	return _Creates.Contract.Add(&_Creates.CallOpts, a, b)
}

// Add is a free data retrieval call binding the contract method 0x771602f7.
//
// Solidity: function add(uint256 a, uint256 b) pure returns(uint256)
func (_Creates *CreatesCallerSession) Add(a *big.Int, b *big.Int) (*big.Int, error) {
	return _Creates.Contract.Add(&_Creates.CallOpts, a, b)
}

// OpCreate is a paid mutator transaction binding the contract method 0x5b8e9959.
//
// Solidity: function opCreate(bytes bytecode, uint256 length) returns(address)
func (_Creates *CreatesTransactor) OpCreate(opts *bind.TransactOpts, bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.contract.Transact(opts, "opCreate", bytecode, length)
}

// OpCreate is a paid mutator transaction binding the contract method 0x5b8e9959.
//
// Solidity: function opCreate(bytes bytecode, uint256 length) returns(address)
func (_Creates *CreatesSession) OpCreate(bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.Contract.OpCreate(&_Creates.TransactOpts, bytecode, length)
}

// OpCreate is a paid mutator transaction binding the contract method 0x5b8e9959.
//
// Solidity: function opCreate(bytes bytecode, uint256 length) returns(address)
func (_Creates *CreatesTransactorSession) OpCreate(bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.Contract.OpCreate(&_Creates.TransactOpts, bytecode, length)
}

// OpCreate2 is a paid mutator transaction binding the contract method 0xe3306a25.
//
// Solidity: function opCreate2(bytes bytecode, uint256 length) returns(address)
func (_Creates *CreatesTransactor) OpCreate2(opts *bind.TransactOpts, bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.contract.Transact(opts, "opCreate2", bytecode, length)
}

// OpCreate2 is a paid mutator transaction binding the contract method 0xe3306a25.
//
// Solidity: function opCreate2(bytes bytecode, uint256 length) returns(address)
func (_Creates *CreatesSession) OpCreate2(bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.Contract.OpCreate2(&_Creates.TransactOpts, bytecode, length)
}

// OpCreate2 is a paid mutator transaction binding the contract method 0xe3306a25.
//
// Solidity: function opCreate2(bytes bytecode, uint256 length) returns(address)
func (_Creates *CreatesTransactorSession) OpCreate2(bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.Contract.OpCreate2(&_Creates.TransactOpts, bytecode, length)
}

// OpCreate2Complex is a paid mutator transaction binding the contract method 0xc935aee4.
//
// Solidity: function opCreate2Complex(bytes bytecode, uint256 length) returns(address, uint256)
func (_Creates *CreatesTransactor) OpCreate2Complex(opts *bind.TransactOpts, bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.contract.Transact(opts, "opCreate2Complex", bytecode, length)
}

// OpCreate2Complex is a paid mutator transaction binding the contract method 0xc935aee4.
//
// Solidity: function opCreate2Complex(bytes bytecode, uint256 length) returns(address, uint256)
func (_Creates *CreatesSession) OpCreate2Complex(bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.Contract.OpCreate2Complex(&_Creates.TransactOpts, bytecode, length)
}

// OpCreate2Complex is a paid mutator transaction binding the contract method 0xc935aee4.
//
// Solidity: function opCreate2Complex(bytes bytecode, uint256 length) returns(address, uint256)
func (_Creates *CreatesTransactorSession) OpCreate2Complex(bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.Contract.OpCreate2Complex(&_Creates.TransactOpts, bytecode, length)
}

// OpCreate2Value is a paid mutator transaction binding the contract method 0xb88c4aa9.
//
// Solidity: function opCreate2Value(bytes bytecode, uint256 length) payable returns(address)
func (_Creates *CreatesTransactor) OpCreate2Value(opts *bind.TransactOpts, bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.contract.Transact(opts, "opCreate2Value", bytecode, length)
}

// OpCreate2Value is a paid mutator transaction binding the contract method 0xb88c4aa9.
//
// Solidity: function opCreate2Value(bytes bytecode, uint256 length) payable returns(address)
func (_Creates *CreatesSession) OpCreate2Value(bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.Contract.OpCreate2Value(&_Creates.TransactOpts, bytecode, length)
}

// OpCreate2Value is a paid mutator transaction binding the contract method 0xb88c4aa9.
//
// Solidity: function opCreate2Value(bytes bytecode, uint256 length) payable returns(address)
func (_Creates *CreatesTransactorSession) OpCreate2Value(bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.Contract.OpCreate2Value(&_Creates.TransactOpts, bytecode, length)
}

// OpCreateValue is a paid mutator transaction binding the contract method 0x3c77eba3.
//
// Solidity: function opCreateValue(bytes bytecode, uint256 length) payable returns(address)
func (_Creates *CreatesTransactor) OpCreateValue(opts *bind.TransactOpts, bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.contract.Transact(opts, "opCreateValue", bytecode, length)
}

// OpCreateValue is a paid mutator transaction binding the contract method 0x3c77eba3.
//
// Solidity: function opCreateValue(bytes bytecode, uint256 length) payable returns(address)
func (_Creates *CreatesSession) OpCreateValue(bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.Contract.OpCreateValue(&_Creates.TransactOpts, bytecode, length)
}

// OpCreateValue is a paid mutator transaction binding the contract method 0x3c77eba3.
//
// Solidity: function opCreateValue(bytes bytecode, uint256 length) payable returns(address)
func (_Creates *CreatesTransactorSession) OpCreateValue(bytecode []byte, length *big.Int) (*types.Transaction, error) {
	return _Creates.Contract.OpCreateValue(&_Creates.TransactOpts, bytecode, length)
}

// SendValue is a paid mutator transaction binding the contract method 0x27c845dc.
//
// Solidity: function sendValue() payable returns()
func (_Creates *CreatesTransactor) SendValue(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Creates.contract.Transact(opts, "sendValue")
}

// SendValue is a paid mutator transaction binding the contract method 0x27c845dc.
//
// Solidity: function sendValue() payable returns()
func (_Creates *CreatesSession) SendValue() (*types.Transaction, error) {
	return _Creates.Contract.SendValue(&_Creates.TransactOpts)
}

// SendValue is a paid mutator transaction binding the contract method 0x27c845dc.
//
// Solidity: function sendValue() payable returns()
func (_Creates *CreatesTransactorSession) SendValue() (*types.Transaction, error) {
	return _Creates.Contract.SendValue(&_Creates.TransactOpts)
}
