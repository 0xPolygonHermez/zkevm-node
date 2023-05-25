// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package FailureTest

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

// FailureTestMetaData contains all meta data concerning the FailureTest contract.
var FailureTestMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"from\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"to\",\"type\":\"uint256\"}],\"name\":\"numberChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"storeAndFail\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061016c806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80636057361d14610046578063b99f3d721461005b578063f2c9ecd81461006e575b600080fd5b61005961005436600461011d565b610083565b005b61005961006936600461011d565b6100c8565b60005460405190815260200160405180910390f35b600080549082905560408051828152602081018490527f64f52b55f0d87dc11f539f2fe367a83c370795772e3312de727215ce118b8fef910160405180910390a15050565b6100d181610083565b60405162461bcd60e51b815260206004820152601860248201527f74686973206d6574686f6420616c77617973206661696c730000000000000000604482015260640160405180910390fd5b60006020828403121561012f57600080fd5b503591905056fea26469706673582212205f0653d69b95544dde31143743c51b6e438ac63084ed44018d865ccfc3ea7c0164736f6c634300080c0033",
}

// FailureTestABI is the input ABI used to generate the binding from.
// Deprecated: Use FailureTestMetaData.ABI instead.
var FailureTestABI = FailureTestMetaData.ABI

// FailureTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FailureTestMetaData.Bin instead.
var FailureTestBin = FailureTestMetaData.Bin

// DeployFailureTest deploys a new Ethereum contract, binding an instance of FailureTest to it.
func DeployFailureTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FailureTest, error) {
	parsed, err := FailureTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FailureTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FailureTest{FailureTestCaller: FailureTestCaller{contract: contract}, FailureTestTransactor: FailureTestTransactor{contract: contract}, FailureTestFilterer: FailureTestFilterer{contract: contract}}, nil
}

// FailureTest is an auto generated Go binding around an Ethereum contract.
type FailureTest struct {
	FailureTestCaller     // Read-only binding to the contract
	FailureTestTransactor // Write-only binding to the contract
	FailureTestFilterer   // Log filterer for contract events
}

// FailureTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type FailureTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FailureTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FailureTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FailureTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FailureTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FailureTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FailureTestSession struct {
	Contract     *FailureTest      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FailureTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FailureTestCallerSession struct {
	Contract *FailureTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// FailureTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FailureTestTransactorSession struct {
	Contract     *FailureTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// FailureTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type FailureTestRaw struct {
	Contract *FailureTest // Generic contract binding to access the raw methods on
}

// FailureTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FailureTestCallerRaw struct {
	Contract *FailureTestCaller // Generic read-only contract binding to access the raw methods on
}

// FailureTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FailureTestTransactorRaw struct {
	Contract *FailureTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFailureTest creates a new instance of FailureTest, bound to a specific deployed contract.
func NewFailureTest(address common.Address, backend bind.ContractBackend) (*FailureTest, error) {
	contract, err := bindFailureTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FailureTest{FailureTestCaller: FailureTestCaller{contract: contract}, FailureTestTransactor: FailureTestTransactor{contract: contract}, FailureTestFilterer: FailureTestFilterer{contract: contract}}, nil
}

// NewFailureTestCaller creates a new read-only instance of FailureTest, bound to a specific deployed contract.
func NewFailureTestCaller(address common.Address, caller bind.ContractCaller) (*FailureTestCaller, error) {
	contract, err := bindFailureTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FailureTestCaller{contract: contract}, nil
}

// NewFailureTestTransactor creates a new write-only instance of FailureTest, bound to a specific deployed contract.
func NewFailureTestTransactor(address common.Address, transactor bind.ContractTransactor) (*FailureTestTransactor, error) {
	contract, err := bindFailureTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FailureTestTransactor{contract: contract}, nil
}

// NewFailureTestFilterer creates a new log filterer instance of FailureTest, bound to a specific deployed contract.
func NewFailureTestFilterer(address common.Address, filterer bind.ContractFilterer) (*FailureTestFilterer, error) {
	contract, err := bindFailureTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FailureTestFilterer{contract: contract}, nil
}

// bindFailureTest binds a generic wrapper to an already deployed contract.
func bindFailureTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FailureTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FailureTest *FailureTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FailureTest.Contract.FailureTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FailureTest *FailureTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FailureTest.Contract.FailureTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FailureTest *FailureTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FailureTest.Contract.FailureTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FailureTest *FailureTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FailureTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FailureTest *FailureTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FailureTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FailureTest *FailureTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FailureTest.Contract.contract.Transact(opts, method, params...)
}

// GetNumber is a free data retrieval call binding the contract method 0xf2c9ecd8.
//
// Solidity: function getNumber() view returns(uint256)
func (_FailureTest *FailureTestCaller) GetNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FailureTest.contract.Call(opts, &out, "getNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNumber is a free data retrieval call binding the contract method 0xf2c9ecd8.
//
// Solidity: function getNumber() view returns(uint256)
func (_FailureTest *FailureTestSession) GetNumber() (*big.Int, error) {
	return _FailureTest.Contract.GetNumber(&_FailureTest.CallOpts)
}

// GetNumber is a free data retrieval call binding the contract method 0xf2c9ecd8.
//
// Solidity: function getNumber() view returns(uint256)
func (_FailureTest *FailureTestCallerSession) GetNumber() (*big.Int, error) {
	return _FailureTest.Contract.GetNumber(&_FailureTest.CallOpts)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 num) returns()
func (_FailureTest *FailureTestTransactor) Store(opts *bind.TransactOpts, num *big.Int) (*types.Transaction, error) {
	return _FailureTest.contract.Transact(opts, "store", num)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 num) returns()
func (_FailureTest *FailureTestSession) Store(num *big.Int) (*types.Transaction, error) {
	return _FailureTest.Contract.Store(&_FailureTest.TransactOpts, num)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 num) returns()
func (_FailureTest *FailureTestTransactorSession) Store(num *big.Int) (*types.Transaction, error) {
	return _FailureTest.Contract.Store(&_FailureTest.TransactOpts, num)
}

// StoreAndFail is a paid mutator transaction binding the contract method 0xb99f3d72.
//
// Solidity: function storeAndFail(uint256 num) returns()
func (_FailureTest *FailureTestTransactor) StoreAndFail(opts *bind.TransactOpts, num *big.Int) (*types.Transaction, error) {
	return _FailureTest.contract.Transact(opts, "storeAndFail", num)
}

// StoreAndFail is a paid mutator transaction binding the contract method 0xb99f3d72.
//
// Solidity: function storeAndFail(uint256 num) returns()
func (_FailureTest *FailureTestSession) StoreAndFail(num *big.Int) (*types.Transaction, error) {
	return _FailureTest.Contract.StoreAndFail(&_FailureTest.TransactOpts, num)
}

// StoreAndFail is a paid mutator transaction binding the contract method 0xb99f3d72.
//
// Solidity: function storeAndFail(uint256 num) returns()
func (_FailureTest *FailureTestTransactorSession) StoreAndFail(num *big.Int) (*types.Transaction, error) {
	return _FailureTest.Contract.StoreAndFail(&_FailureTest.TransactOpts, num)
}

// FailureTestNumberChangedIterator is returned from FilterNumberChanged and is used to iterate over the raw logs and unpacked data for NumberChanged events raised by the FailureTest contract.
type FailureTestNumberChangedIterator struct {
	Event *FailureTestNumberChanged // Event containing the contract specifics and raw log

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
func (it *FailureTestNumberChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FailureTestNumberChanged)
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
		it.Event = new(FailureTestNumberChanged)
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
func (it *FailureTestNumberChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FailureTestNumberChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FailureTestNumberChanged represents a NumberChanged event raised by the FailureTest contract.
type FailureTestNumberChanged struct {
	From *big.Int
	To   *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNumberChanged is a free log retrieval operation binding the contract event 0x64f52b55f0d87dc11f539f2fe367a83c370795772e3312de727215ce118b8fef.
//
// Solidity: event numberChanged(uint256 from, uint256 to)
func (_FailureTest *FailureTestFilterer) FilterNumberChanged(opts *bind.FilterOpts) (*FailureTestNumberChangedIterator, error) {

	logs, sub, err := _FailureTest.contract.FilterLogs(opts, "numberChanged")
	if err != nil {
		return nil, err
	}
	return &FailureTestNumberChangedIterator{contract: _FailureTest.contract, event: "numberChanged", logs: logs, sub: sub}, nil
}

// WatchNumberChanged is a free log subscription operation binding the contract event 0x64f52b55f0d87dc11f539f2fe367a83c370795772e3312de727215ce118b8fef.
//
// Solidity: event numberChanged(uint256 from, uint256 to)
func (_FailureTest *FailureTestFilterer) WatchNumberChanged(opts *bind.WatchOpts, sink chan<- *FailureTestNumberChanged) (event.Subscription, error) {

	logs, sub, err := _FailureTest.contract.WatchLogs(opts, "numberChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FailureTestNumberChanged)
				if err := _FailureTest.contract.UnpackLog(event, "numberChanged", log); err != nil {
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

// ParseNumberChanged is a log parse operation binding the contract event 0x64f52b55f0d87dc11f539f2fe367a83c370795772e3312de727215ce118b8fef.
//
// Solidity: event numberChanged(uint256 from, uint256 to)
func (_FailureTest *FailureTestFilterer) ParseNumberChanged(log types.Log) (*FailureTestNumberChanged, error) {
	event := new(FailureTestNumberChanged)
	if err := _FailureTest.contract.UnpackLog(event, "numberChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
