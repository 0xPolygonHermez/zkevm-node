// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package UniswapV2Migrator

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

// UniswapV2MigratorMetaData contains all meta data concerning the UniswapV2Migrator contract.
var UniswapV2MigratorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factoryV1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60c060405234801561001057600080fd5b506040516109443803806109448339818101604052604081101561003357600080fd5b5080516020909101516001600160601b0319606092831b8116608052911b1660a05260805160601c60a05160601c6108bd6100876000398061030752806103775280610409525080608352506108bd6000f3fe6080604052600436106100225760003560e01c8063b7df1d251461002e57610029565b3661002957005b600080fd5b34801561003a57600080fd5b5061007d600480360360a081101561005157600080fd5b506001600160a01b0381358116916020810135916040820135916060810135909116906080013561007f565b005b60007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166306f2bf62876040518263ffffffff1660e01b815260040180826001600160a01b03166001600160a01b0316815260200191505060206040518083038186803b1580156100f757600080fd5b505afa15801561010b573d6000803e3d6000fd5b505050506040513d602081101561012157600080fd5b5051604080516370a0823160e01b815233600482015290519192506000916001600160a01b038416916370a08231916024808301926020929190829003018186803b15801561016f57600080fd5b505afa158015610183573d6000803e3d6000fd5b505050506040513d602081101561019957600080fd5b5051604080516323b872dd60e01b81523360048201523060248201526044810183905290519192506001600160a01b038416916323b872dd916064808201926020929091908290030181600087803b1580156101f457600080fd5b505af1158015610208573d6000803e3d6000fd5b505050506040513d602081101561021e57600080fd5b5051610268576040805162461bcd60e51b81526020600482015260146024820152731514905394d1915497d19493d357d1905253115160621b604482015290519081900360640190fd5b60408051637c45f8ad60e11b81526004810183905260016024820181905260448201526000196064820152815160009283926001600160a01b0387169263f88bf15a9260848084019391929182900301818787803b1580156102c957600080fd5b505af11580156102dd573d6000803e3d6000fd5b505050506040513d60408110156102f357600080fd5b508051602090910151909250905061032c897f000000000000000000000000000000000000000000000000000000000000000083610462565b6040805163f305d71960e01b81526001600160a01b038b8116600483015260248201849052604482018b9052606482018a9052888116608483015260a48201889052915160009283927f00000000000000000000000000000000000000000000000000000000000000009091169163f305d71991879160c480830192606092919082900301818588803b1580156103c257600080fd5b505af11580156103d6573d6000803e3d6000fd5b50505050506040513d60608110156103ed57600080fd5b5080516020909101519092509050818311156104415761042f8b7f00000000000000000000000000000000000000000000000000000000000000006000610462565b61043c8b338486036105b6565b610455565b808411156104555761045533828603610703565b5050505050505050505050565b604080516001600160a01b038481166024830152604480830185905283518084039091018152606490920183526020820180516001600160e01b031663095ea7b360e01b178152925182516000946060949389169392918291908083835b602083106104df5780518252601f1990920191602091820191016104c0565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610541576040519150601f19603f3d011682016040523d82523d6000602084013e610546565b606091505b5091509150818015610574575080511580610574575080806020019051602081101561057157600080fd5b50515b6105af5760405162461bcd60e51b815260040180806020018281038252602b815260200180610830602b913960400191505060405180910390fd5b5050505050565b604080516001600160a01b038481166024830152604480830185905283518084039091018152606490920183526020820180516001600160e01b031663a9059cbb60e01b178152925182516000946060949389169392918291908083835b602083106106335780518252601f199092019160209182019101610614565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610695576040519150601f19603f3d011682016040523d82523d6000602084013e61069a565b606091505b50915091508180156106c85750805115806106c857508080602001905160208110156106c557600080fd5b50515b6105af5760405162461bcd60e51b815260040180806020018281038252602d81526020018061085b602d913960400191505060405180910390fd5b604080516000808252602082019092526001600160a01b0384169083906040518082805190602001908083835b6020831061074f5780518252601f199092019160209182019101610730565b6001836020036101000a03801982511681845116808217855250505050505090500191505060006040518083038185875af1925050503d80600081146107b1576040519150601f19603f3d011682016040523d82523d6000602084013e6107b6565b606091505b50509050806107f65760405162461bcd60e51b81526004018080602001828103825260348152602001806107fc6034913960400191505060405180910390fd5b50505056fe5472616e7366657248656c7065723a3a736166655472616e736665724554483a20455448207472616e73666572206661696c65645472616e7366657248656c7065723a3a73616665417070726f76653a20617070726f7665206661696c65645472616e7366657248656c7065723a3a736166655472616e736665723a207472616e73666572206661696c6564a26469706673582212201899c4ec7f4361c60e41fea7ed99f592748923c045a19f1a2483842d87af714f64736f6c63430006060033",
}

// UniswapV2MigratorABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapV2MigratorMetaData.ABI instead.
var UniswapV2MigratorABI = UniswapV2MigratorMetaData.ABI

// UniswapV2MigratorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UniswapV2MigratorMetaData.Bin instead.
var UniswapV2MigratorBin = UniswapV2MigratorMetaData.Bin

// DeployUniswapV2Migrator deploys a new Ethereum contract, binding an instance of UniswapV2Migrator to it.
func DeployUniswapV2Migrator(auth *bind.TransactOpts, backend bind.ContractBackend, _factoryV1 common.Address, _router common.Address) (common.Address, *types.Transaction, *UniswapV2Migrator, error) {
	parsed, err := UniswapV2MigratorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UniswapV2MigratorBin), backend, _factoryV1, _router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UniswapV2Migrator{UniswapV2MigratorCaller: UniswapV2MigratorCaller{contract: contract}, UniswapV2MigratorTransactor: UniswapV2MigratorTransactor{contract: contract}, UniswapV2MigratorFilterer: UniswapV2MigratorFilterer{contract: contract}}, nil
}

// UniswapV2Migrator is an auto generated Go binding around an Ethereum contract.
type UniswapV2Migrator struct {
	UniswapV2MigratorCaller     // Read-only binding to the contract
	UniswapV2MigratorTransactor // Write-only binding to the contract
	UniswapV2MigratorFilterer   // Log filterer for contract events
}

// UniswapV2MigratorCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniswapV2MigratorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2MigratorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniswapV2MigratorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2MigratorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniswapV2MigratorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2MigratorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniswapV2MigratorSession struct {
	Contract     *UniswapV2Migrator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// UniswapV2MigratorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniswapV2MigratorCallerSession struct {
	Contract *UniswapV2MigratorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// UniswapV2MigratorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniswapV2MigratorTransactorSession struct {
	Contract     *UniswapV2MigratorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// UniswapV2MigratorRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniswapV2MigratorRaw struct {
	Contract *UniswapV2Migrator // Generic contract binding to access the raw methods on
}

// UniswapV2MigratorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniswapV2MigratorCallerRaw struct {
	Contract *UniswapV2MigratorCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapV2MigratorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniswapV2MigratorTransactorRaw struct {
	Contract *UniswapV2MigratorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapV2Migrator creates a new instance of UniswapV2Migrator, bound to a specific deployed contract.
func NewUniswapV2Migrator(address common.Address, backend bind.ContractBackend) (*UniswapV2Migrator, error) {
	contract, err := bindUniswapV2Migrator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapV2Migrator{UniswapV2MigratorCaller: UniswapV2MigratorCaller{contract: contract}, UniswapV2MigratorTransactor: UniswapV2MigratorTransactor{contract: contract}, UniswapV2MigratorFilterer: UniswapV2MigratorFilterer{contract: contract}}, nil
}

// NewUniswapV2MigratorCaller creates a new read-only instance of UniswapV2Migrator, bound to a specific deployed contract.
func NewUniswapV2MigratorCaller(address common.Address, caller bind.ContractCaller) (*UniswapV2MigratorCaller, error) {
	contract, err := bindUniswapV2Migrator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2MigratorCaller{contract: contract}, nil
}

// NewUniswapV2MigratorTransactor creates a new write-only instance of UniswapV2Migrator, bound to a specific deployed contract.
func NewUniswapV2MigratorTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapV2MigratorTransactor, error) {
	contract, err := bindUniswapV2Migrator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2MigratorTransactor{contract: contract}, nil
}

// NewUniswapV2MigratorFilterer creates a new log filterer instance of UniswapV2Migrator, bound to a specific deployed contract.
func NewUniswapV2MigratorFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapV2MigratorFilterer, error) {
	contract, err := bindUniswapV2Migrator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapV2MigratorFilterer{contract: contract}, nil
}

// bindUniswapV2Migrator binds a generic wrapper to an already deployed contract.
func bindUniswapV2Migrator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UniswapV2MigratorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2Migrator *UniswapV2MigratorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2Migrator.Contract.UniswapV2MigratorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2Migrator *UniswapV2MigratorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.UniswapV2MigratorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2Migrator *UniswapV2MigratorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.UniswapV2MigratorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2Migrator *UniswapV2MigratorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2Migrator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2Migrator *UniswapV2MigratorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2Migrator *UniswapV2MigratorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.contract.Transact(opts, method, params...)
}

// Migrate is a paid mutator transaction binding the contract method 0xb7df1d25.
//
// Solidity: function migrate(address token, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns()
func (_UniswapV2Migrator *UniswapV2MigratorTransactor) Migrate(opts *bind.TransactOpts, token common.Address, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV2Migrator.contract.Transact(opts, "migrate", token, amountTokenMin, amountETHMin, to, deadline)
}

// Migrate is a paid mutator transaction binding the contract method 0xb7df1d25.
//
// Solidity: function migrate(address token, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns()
func (_UniswapV2Migrator *UniswapV2MigratorSession) Migrate(token common.Address, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.Migrate(&_UniswapV2Migrator.TransactOpts, token, amountTokenMin, amountETHMin, to, deadline)
}

// Migrate is a paid mutator transaction binding the contract method 0xb7df1d25.
//
// Solidity: function migrate(address token, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns()
func (_UniswapV2Migrator *UniswapV2MigratorTransactorSession) Migrate(token common.Address, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.Migrate(&_UniswapV2Migrator.TransactOpts, token, amountTokenMin, amountETHMin, to, deadline)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_UniswapV2Migrator *UniswapV2MigratorTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Migrator.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_UniswapV2Migrator *UniswapV2MigratorSession) Receive() (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.Receive(&_UniswapV2Migrator.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_UniswapV2Migrator *UniswapV2MigratorTransactorSession) Receive() (*types.Transaction, error) {
	return _UniswapV2Migrator.Contract.Receive(&_UniswapV2Migrator.TransactOpts)
}
