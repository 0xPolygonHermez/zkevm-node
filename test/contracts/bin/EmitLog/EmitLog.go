// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package EmitLog

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

// EmitLogMetaData contains all meta data concerning the EmitLog contract.
var EmitLogMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[],\"name\":\"Log\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"}],\"name\":\"LogA\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"}],\"name\":\"LogAB\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"}],\"name\":\"LogABC\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"b\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"c\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"d\",\"type\":\"uint256\"}],\"name\":\"LogABCD\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"emitLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061025e806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80637966b4f614610030575b600080fd5b61003861003a565b005b6040517f5e7df75d54e493185612379c616118a4c9ac802de621b010c96f74d22df4b30a90600090a16040516001907f977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac6290600090a26040516002906001907fbb6e4da744abea70325874159d52c1ad3e57babfae7c329a948e7dcb274deb0990600090a36003600260017f966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b760405160405180910390a46003600260017fe5562b12d9276c5c987df08afff7b1946f2d869236866ea2285c7e2e95685a64600460405161012891815260200190565b60405180910390a46002600360047fe5562b12d9276c5c987df08afff7b1946f2d869236866ea2285c7e2e95685a64600160405161016891815260200190565b60405180910390a46001600260037f966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b760405160405180910390a46040516001906002907fbb6e4da744abea70325874159d52c1ad3e57babfae7c329a948e7dcb274deb0990600090a36040516001907f977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac6290600090a26040517f5e7df75d54e493185612379c616118a4c9ac802de621b010c96f74d22df4b30a90600090a156fea26469706673582212204c9779aabd3e8dfef2bffe2d88075ca74555273b14d3766908a126f2cf434b8964736f6c634300080c0033",
}

// EmitLogABI is the input ABI used to generate the binding from.
// Deprecated: Use EmitLogMetaData.ABI instead.
var EmitLogABI = EmitLogMetaData.ABI

// EmitLogBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EmitLogMetaData.Bin instead.
var EmitLogBin = EmitLogMetaData.Bin

// DeployEmitLog deploys a new Ethereum contract, binding an instance of EmitLog to it.
func DeployEmitLog(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EmitLog, error) {
	parsed, err := EmitLogMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EmitLogBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EmitLog{EmitLogCaller: EmitLogCaller{contract: contract}, EmitLogTransactor: EmitLogTransactor{contract: contract}, EmitLogFilterer: EmitLogFilterer{contract: contract}}, nil
}

// EmitLog is an auto generated Go binding around an Ethereum contract.
type EmitLog struct {
	EmitLogCaller     // Read-only binding to the contract
	EmitLogTransactor // Write-only binding to the contract
	EmitLogFilterer   // Log filterer for contract events
}

// EmitLogCaller is an auto generated read-only Go binding around an Ethereum contract.
type EmitLogCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EmitLogTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EmitLogTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EmitLogFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EmitLogFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EmitLogSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EmitLogSession struct {
	Contract     *EmitLog          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EmitLogCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EmitLogCallerSession struct {
	Contract *EmitLogCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// EmitLogTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EmitLogTransactorSession struct {
	Contract     *EmitLogTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// EmitLogRaw is an auto generated low-level Go binding around an Ethereum contract.
type EmitLogRaw struct {
	Contract *EmitLog // Generic contract binding to access the raw methods on
}

// EmitLogCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EmitLogCallerRaw struct {
	Contract *EmitLogCaller // Generic read-only contract binding to access the raw methods on
}

// EmitLogTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EmitLogTransactorRaw struct {
	Contract *EmitLogTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEmitLog creates a new instance of EmitLog, bound to a specific deployed contract.
func NewEmitLog(address common.Address, backend bind.ContractBackend) (*EmitLog, error) {
	contract, err := bindEmitLog(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EmitLog{EmitLogCaller: EmitLogCaller{contract: contract}, EmitLogTransactor: EmitLogTransactor{contract: contract}, EmitLogFilterer: EmitLogFilterer{contract: contract}}, nil
}

// NewEmitLogCaller creates a new read-only instance of EmitLog, bound to a specific deployed contract.
func NewEmitLogCaller(address common.Address, caller bind.ContractCaller) (*EmitLogCaller, error) {
	contract, err := bindEmitLog(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EmitLogCaller{contract: contract}, nil
}

// NewEmitLogTransactor creates a new write-only instance of EmitLog, bound to a specific deployed contract.
func NewEmitLogTransactor(address common.Address, transactor bind.ContractTransactor) (*EmitLogTransactor, error) {
	contract, err := bindEmitLog(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EmitLogTransactor{contract: contract}, nil
}

// NewEmitLogFilterer creates a new log filterer instance of EmitLog, bound to a specific deployed contract.
func NewEmitLogFilterer(address common.Address, filterer bind.ContractFilterer) (*EmitLogFilterer, error) {
	contract, err := bindEmitLog(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EmitLogFilterer{contract: contract}, nil
}

// bindEmitLog binds a generic wrapper to an already deployed contract.
func bindEmitLog(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EmitLogMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EmitLog *EmitLogRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EmitLog.Contract.EmitLogCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EmitLog *EmitLogRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EmitLog.Contract.EmitLogTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EmitLog *EmitLogRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EmitLog.Contract.EmitLogTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EmitLog *EmitLogCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EmitLog.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EmitLog *EmitLogTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EmitLog.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EmitLog *EmitLogTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EmitLog.Contract.contract.Transact(opts, method, params...)
}

// EmitLogs is a paid mutator transaction binding the contract method 0x7966b4f6.
//
// Solidity: function emitLogs() returns()
func (_EmitLog *EmitLogTransactor) EmitLogs(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EmitLog.contract.Transact(opts, "emitLogs")
}

// EmitLogs is a paid mutator transaction binding the contract method 0x7966b4f6.
//
// Solidity: function emitLogs() returns()
func (_EmitLog *EmitLogSession) EmitLogs() (*types.Transaction, error) {
	return _EmitLog.Contract.EmitLogs(&_EmitLog.TransactOpts)
}

// EmitLogs is a paid mutator transaction binding the contract method 0x7966b4f6.
//
// Solidity: function emitLogs() returns()
func (_EmitLog *EmitLogTransactorSession) EmitLogs() (*types.Transaction, error) {
	return _EmitLog.Contract.EmitLogs(&_EmitLog.TransactOpts)
}

// EmitLogLogIterator is returned from FilterLog and is used to iterate over the raw logs and unpacked data for Log events raised by the EmitLog contract.
type EmitLogLogIterator struct {
	Event *EmitLogLog // Event containing the contract specifics and raw log

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
func (it *EmitLogLogIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EmitLogLog)
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
		it.Event = new(EmitLogLog)
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
func (it *EmitLogLogIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EmitLogLogIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EmitLogLog represents a Log event raised by the EmitLog contract.
type EmitLogLog struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLog is a free log retrieval operation binding the contract event 0x5e7df75d54e493185612379c616118a4c9ac802de621b010c96f74d22df4b30a.
//
// Solidity: event Log()
func (_EmitLog *EmitLogFilterer) FilterLog(opts *bind.FilterOpts) (*EmitLogLogIterator, error) {

	logs, sub, err := _EmitLog.contract.FilterLogs(opts, "Log")
	if err != nil {
		return nil, err
	}
	return &EmitLogLogIterator{contract: _EmitLog.contract, event: "Log", logs: logs, sub: sub}, nil
}

// WatchLog is a free log subscription operation binding the contract event 0x5e7df75d54e493185612379c616118a4c9ac802de621b010c96f74d22df4b30a.
//
// Solidity: event Log()
func (_EmitLog *EmitLogFilterer) WatchLog(opts *bind.WatchOpts, sink chan<- *EmitLogLog) (event.Subscription, error) {

	logs, sub, err := _EmitLog.contract.WatchLogs(opts, "Log")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EmitLogLog)
				if err := _EmitLog.contract.UnpackLog(event, "Log", log); err != nil {
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

// ParseLog is a log parse operation binding the contract event 0x5e7df75d54e493185612379c616118a4c9ac802de621b010c96f74d22df4b30a.
//
// Solidity: event Log()
func (_EmitLog *EmitLogFilterer) ParseLog(log types.Log) (*EmitLogLog, error) {
	event := new(EmitLogLog)
	if err := _EmitLog.contract.UnpackLog(event, "Log", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EmitLogLogAIterator is returned from FilterLogA and is used to iterate over the raw logs and unpacked data for LogA events raised by the EmitLog contract.
type EmitLogLogAIterator struct {
	Event *EmitLogLogA // Event containing the contract specifics and raw log

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
func (it *EmitLogLogAIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EmitLogLogA)
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
		it.Event = new(EmitLogLogA)
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
func (it *EmitLogLogAIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EmitLogLogAIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EmitLogLogA represents a LogA event raised by the EmitLog contract.
type EmitLogLogA struct {
	A   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogA is a free log retrieval operation binding the contract event 0x977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac62.
//
// Solidity: event LogA(uint256 indexed a)
func (_EmitLog *EmitLogFilterer) FilterLogA(opts *bind.FilterOpts, a []*big.Int) (*EmitLogLogAIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _EmitLog.contract.FilterLogs(opts, "LogA", aRule)
	if err != nil {
		return nil, err
	}
	return &EmitLogLogAIterator{contract: _EmitLog.contract, event: "LogA", logs: logs, sub: sub}, nil
}

// WatchLogA is a free log subscription operation binding the contract event 0x977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac62.
//
// Solidity: event LogA(uint256 indexed a)
func (_EmitLog *EmitLogFilterer) WatchLogA(opts *bind.WatchOpts, sink chan<- *EmitLogLogA, a []*big.Int) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _EmitLog.contract.WatchLogs(opts, "LogA", aRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EmitLogLogA)
				if err := _EmitLog.contract.UnpackLog(event, "LogA", log); err != nil {
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

// ParseLogA is a log parse operation binding the contract event 0x977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac62.
//
// Solidity: event LogA(uint256 indexed a)
func (_EmitLog *EmitLogFilterer) ParseLogA(log types.Log) (*EmitLogLogA, error) {
	event := new(EmitLogLogA)
	if err := _EmitLog.contract.UnpackLog(event, "LogA", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EmitLogLogABIterator is returned from FilterLogAB and is used to iterate over the raw logs and unpacked data for LogAB events raised by the EmitLog contract.
type EmitLogLogABIterator struct {
	Event *EmitLogLogAB // Event containing the contract specifics and raw log

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
func (it *EmitLogLogABIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EmitLogLogAB)
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
		it.Event = new(EmitLogLogAB)
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
func (it *EmitLogLogABIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EmitLogLogABIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EmitLogLogAB represents a LogAB event raised by the EmitLog contract.
type EmitLogLogAB struct {
	A   *big.Int
	B   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogAB is a free log retrieval operation binding the contract event 0xbb6e4da744abea70325874159d52c1ad3e57babfae7c329a948e7dcb274deb09.
//
// Solidity: event LogAB(uint256 indexed a, uint256 indexed b)
func (_EmitLog *EmitLogFilterer) FilterLogAB(opts *bind.FilterOpts, a []*big.Int, b []*big.Int) (*EmitLogLogABIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}

	logs, sub, err := _EmitLog.contract.FilterLogs(opts, "LogAB", aRule, bRule)
	if err != nil {
		return nil, err
	}
	return &EmitLogLogABIterator{contract: _EmitLog.contract, event: "LogAB", logs: logs, sub: sub}, nil
}

// WatchLogAB is a free log subscription operation binding the contract event 0xbb6e4da744abea70325874159d52c1ad3e57babfae7c329a948e7dcb274deb09.
//
// Solidity: event LogAB(uint256 indexed a, uint256 indexed b)
func (_EmitLog *EmitLogFilterer) WatchLogAB(opts *bind.WatchOpts, sink chan<- *EmitLogLogAB, a []*big.Int, b []*big.Int) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}

	logs, sub, err := _EmitLog.contract.WatchLogs(opts, "LogAB", aRule, bRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EmitLogLogAB)
				if err := _EmitLog.contract.UnpackLog(event, "LogAB", log); err != nil {
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

// ParseLogAB is a log parse operation binding the contract event 0xbb6e4da744abea70325874159d52c1ad3e57babfae7c329a948e7dcb274deb09.
//
// Solidity: event LogAB(uint256 indexed a, uint256 indexed b)
func (_EmitLog *EmitLogFilterer) ParseLogAB(log types.Log) (*EmitLogLogAB, error) {
	event := new(EmitLogLogAB)
	if err := _EmitLog.contract.UnpackLog(event, "LogAB", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EmitLogLogABCIterator is returned from FilterLogABC and is used to iterate over the raw logs and unpacked data for LogABC events raised by the EmitLog contract.
type EmitLogLogABCIterator struct {
	Event *EmitLogLogABC // Event containing the contract specifics and raw log

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
func (it *EmitLogLogABCIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EmitLogLogABC)
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
		it.Event = new(EmitLogLogABC)
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
func (it *EmitLogLogABCIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EmitLogLogABCIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EmitLogLogABC represents a LogABC event raised by the EmitLog contract.
type EmitLogLogABC struct {
	A   *big.Int
	B   *big.Int
	C   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogABC is a free log retrieval operation binding the contract event 0x966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b7.
//
// Solidity: event LogABC(uint256 indexed a, uint256 indexed b, uint256 indexed c)
func (_EmitLog *EmitLogFilterer) FilterLogABC(opts *bind.FilterOpts, a []*big.Int, b []*big.Int, c []*big.Int) (*EmitLogLogABCIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}
	var cRule []interface{}
	for _, cItem := range c {
		cRule = append(cRule, cItem)
	}

	logs, sub, err := _EmitLog.contract.FilterLogs(opts, "LogABC", aRule, bRule, cRule)
	if err != nil {
		return nil, err
	}
	return &EmitLogLogABCIterator{contract: _EmitLog.contract, event: "LogABC", logs: logs, sub: sub}, nil
}

// WatchLogABC is a free log subscription operation binding the contract event 0x966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b7.
//
// Solidity: event LogABC(uint256 indexed a, uint256 indexed b, uint256 indexed c)
func (_EmitLog *EmitLogFilterer) WatchLogABC(opts *bind.WatchOpts, sink chan<- *EmitLogLogABC, a []*big.Int, b []*big.Int, c []*big.Int) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}
	var cRule []interface{}
	for _, cItem := range c {
		cRule = append(cRule, cItem)
	}

	logs, sub, err := _EmitLog.contract.WatchLogs(opts, "LogABC", aRule, bRule, cRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EmitLogLogABC)
				if err := _EmitLog.contract.UnpackLog(event, "LogABC", log); err != nil {
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

// ParseLogABC is a log parse operation binding the contract event 0x966018f1afaee50c6bcf5eb4ae089eeb650bd1deb473395d69dd307ef2e689b7.
//
// Solidity: event LogABC(uint256 indexed a, uint256 indexed b, uint256 indexed c)
func (_EmitLog *EmitLogFilterer) ParseLogABC(log types.Log) (*EmitLogLogABC, error) {
	event := new(EmitLogLogABC)
	if err := _EmitLog.contract.UnpackLog(event, "LogABC", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EmitLogLogABCDIterator is returned from FilterLogABCD and is used to iterate over the raw logs and unpacked data for LogABCD events raised by the EmitLog contract.
type EmitLogLogABCDIterator struct {
	Event *EmitLogLogABCD // Event containing the contract specifics and raw log

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
func (it *EmitLogLogABCDIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EmitLogLogABCD)
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
		it.Event = new(EmitLogLogABCD)
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
func (it *EmitLogLogABCDIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EmitLogLogABCDIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EmitLogLogABCD represents a LogABCD event raised by the EmitLog contract.
type EmitLogLogABCD struct {
	A   *big.Int
	B   *big.Int
	C   *big.Int
	D   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogABCD is a free log retrieval operation binding the contract event 0xe5562b12d9276c5c987df08afff7b1946f2d869236866ea2285c7e2e95685a64.
//
// Solidity: event LogABCD(uint256 indexed a, uint256 indexed b, uint256 indexed c, uint256 d)
func (_EmitLog *EmitLogFilterer) FilterLogABCD(opts *bind.FilterOpts, a []*big.Int, b []*big.Int, c []*big.Int) (*EmitLogLogABCDIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}
	var cRule []interface{}
	for _, cItem := range c {
		cRule = append(cRule, cItem)
	}

	logs, sub, err := _EmitLog.contract.FilterLogs(opts, "LogABCD", aRule, bRule, cRule)
	if err != nil {
		return nil, err
	}
	return &EmitLogLogABCDIterator{contract: _EmitLog.contract, event: "LogABCD", logs: logs, sub: sub}, nil
}

// WatchLogABCD is a free log subscription operation binding the contract event 0xe5562b12d9276c5c987df08afff7b1946f2d869236866ea2285c7e2e95685a64.
//
// Solidity: event LogABCD(uint256 indexed a, uint256 indexed b, uint256 indexed c, uint256 d)
func (_EmitLog *EmitLogFilterer) WatchLogABCD(opts *bind.WatchOpts, sink chan<- *EmitLogLogABCD, a []*big.Int, b []*big.Int, c []*big.Int) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}
	var cRule []interface{}
	for _, cItem := range c {
		cRule = append(cRule, cItem)
	}

	logs, sub, err := _EmitLog.contract.WatchLogs(opts, "LogABCD", aRule, bRule, cRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EmitLogLogABCD)
				if err := _EmitLog.contract.UnpackLog(event, "LogABCD", log); err != nil {
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

// ParseLogABCD is a log parse operation binding the contract event 0xe5562b12d9276c5c987df08afff7b1946f2d869236866ea2285c7e2e95685a64.
//
// Solidity: event LogABCD(uint256 indexed a, uint256 indexed b, uint256 indexed c, uint256 d)
func (_EmitLog *EmitLogFilterer) ParseLogABCD(log types.Log) (*EmitLogLogABCD, error) {
	event := new(EmitLogLogABCD)
	if err := _EmitLog.contract.UnpackLog(event, "LogABCD", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
