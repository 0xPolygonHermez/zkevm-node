// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mockverifier

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

// MockverifierMetaData contains all meta data concerning the Mockverifier contract.
var MockverifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[1]\",\"name\":\"input\",\"type\":\"uint256[1]\"}],\"name\":\"verifyProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610101806100206000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806343753b4d14602d575b600080fd5b60426038366004606c565b6001949350505050565b604051901515815260200160405180910390f35b8060408101831015606657600080fd5b92915050565b600080600080610120808688031215608357600080fd5b608b87876056565b945060c0860187811115609d57600080fd5b60408701945060ab88826056565b93505086818701111560bc57600080fd5b5092959194509261010001915056fea2646970667358221220569286c06a35bb3b6e29da343de6015d1fc7ba685f6b59f888f00e18129e37df64736f6c634300080f0033",
}

// MockverifierABI is the input ABI used to generate the binding from.
// Deprecated: Use MockverifierMetaData.ABI instead.
var MockverifierABI = MockverifierMetaData.ABI

// MockverifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockverifierMetaData.Bin instead.
var MockverifierBin = MockverifierMetaData.Bin

// DeployMockverifier deploys a new Ethereum contract, binding an instance of Mockverifier to it.
func DeployMockverifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Mockverifier, error) {
	parsed, err := MockverifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockverifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Mockverifier{MockverifierCaller: MockverifierCaller{contract: contract}, MockverifierTransactor: MockverifierTransactor{contract: contract}, MockverifierFilterer: MockverifierFilterer{contract: contract}}, nil
}

// Mockverifier is an auto generated Go binding around an Ethereum contract.
type Mockverifier struct {
	MockverifierCaller     // Read-only binding to the contract
	MockverifierTransactor // Write-only binding to the contract
	MockverifierFilterer   // Log filterer for contract events
}

// MockverifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockverifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockverifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockverifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockverifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockverifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockverifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockverifierSession struct {
	Contract     *Mockverifier     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MockverifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockverifierCallerSession struct {
	Contract *MockverifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// MockverifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockverifierTransactorSession struct {
	Contract     *MockverifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// MockverifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockverifierRaw struct {
	Contract *Mockverifier // Generic contract binding to access the raw methods on
}

// MockverifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockverifierCallerRaw struct {
	Contract *MockverifierCaller // Generic read-only contract binding to access the raw methods on
}

// MockverifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockverifierTransactorRaw struct {
	Contract *MockverifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockverifier creates a new instance of Mockverifier, bound to a specific deployed contract.
func NewMockverifier(address common.Address, backend bind.ContractBackend) (*Mockverifier, error) {
	contract, err := bindMockverifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Mockverifier{MockverifierCaller: MockverifierCaller{contract: contract}, MockverifierTransactor: MockverifierTransactor{contract: contract}, MockverifierFilterer: MockverifierFilterer{contract: contract}}, nil
}

// NewMockverifierCaller creates a new read-only instance of Mockverifier, bound to a specific deployed contract.
func NewMockverifierCaller(address common.Address, caller bind.ContractCaller) (*MockverifierCaller, error) {
	contract, err := bindMockverifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockverifierCaller{contract: contract}, nil
}

// NewMockverifierTransactor creates a new write-only instance of Mockverifier, bound to a specific deployed contract.
func NewMockverifierTransactor(address common.Address, transactor bind.ContractTransactor) (*MockverifierTransactor, error) {
	contract, err := bindMockverifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockverifierTransactor{contract: contract}, nil
}

// NewMockverifierFilterer creates a new log filterer instance of Mockverifier, bound to a specific deployed contract.
func NewMockverifierFilterer(address common.Address, filterer bind.ContractFilterer) (*MockverifierFilterer, error) {
	contract, err := bindMockverifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockverifierFilterer{contract: contract}, nil
}

// bindMockverifier binds a generic wrapper to an already deployed contract.
func bindMockverifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MockverifierABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mockverifier *MockverifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mockverifier.Contract.MockverifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mockverifier *MockverifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mockverifier.Contract.MockverifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mockverifier *MockverifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mockverifier.Contract.MockverifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mockverifier *MockverifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mockverifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mockverifier *MockverifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mockverifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mockverifier *MockverifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mockverifier.Contract.contract.Transact(opts, method, params...)
}

// VerifyProof is a free data retrieval call binding the contract method 0x43753b4d.
//
// Solidity: function verifyProof(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[1] input) view returns(bool)
func (_Mockverifier *MockverifierCaller) VerifyProof(opts *bind.CallOpts, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [1]*big.Int) (bool, error) {
	var out []interface{}
	err := _Mockverifier.contract.Call(opts, &out, "verifyProof", a, b, c, input)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyProof is a free data retrieval call binding the contract method 0x43753b4d.
//
// Solidity: function verifyProof(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[1] input) view returns(bool)
func (_Mockverifier *MockverifierSession) VerifyProof(a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [1]*big.Int) (bool, error) {
	return _Mockverifier.Contract.VerifyProof(&_Mockverifier.CallOpts, a, b, c, input)
}

// VerifyProof is a free data retrieval call binding the contract method 0x43753b4d.
//
// Solidity: function verifyProof(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[1] input) view returns(bool)
func (_Mockverifier *MockverifierCallerSession) VerifyProof(a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [1]*big.Int) (bool, error) {
	return _Mockverifier.Contract.VerifyProof(&_Mockverifier.CallOpts, a, b, c, input)
}
