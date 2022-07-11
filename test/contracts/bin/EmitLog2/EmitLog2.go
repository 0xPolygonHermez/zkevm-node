// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package EmitLog2

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

// EmitLog2MetaData contains all meta data concerning the EmitLog2 contract.
var EmitLog2MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"}],\"name\":\"LogA\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"emitLogs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b5060998061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80637966b4f614602d575b600080fd5b60336035565b005b6040516001907f977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac6290600090a256fea26469706673582212200aed1b5f98749897f641535d09c3f58975b5e30a8dfed8781f982ae3e31fd16e64736f6c634300080c0033",
}

// EmitLog2ABI is the input ABI used to generate the binding from.
// Deprecated: Use EmitLog2MetaData.ABI instead.
var EmitLog2ABI = EmitLog2MetaData.ABI

// EmitLog2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EmitLog2MetaData.Bin instead.
var EmitLog2Bin = EmitLog2MetaData.Bin

// DeployEmitLog2 deploys a new Ethereum contract, binding an instance of EmitLog2 to it.
func DeployEmitLog2(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EmitLog2, error) {
	parsed, err := EmitLog2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EmitLog2Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EmitLog2{EmitLog2Caller: EmitLog2Caller{contract: contract}, EmitLog2Transactor: EmitLog2Transactor{contract: contract}, EmitLog2Filterer: EmitLog2Filterer{contract: contract}}, nil
}

// EmitLog2 is an auto generated Go binding around an Ethereum contract.
type EmitLog2 struct {
	EmitLog2Caller     // Read-only binding to the contract
	EmitLog2Transactor // Write-only binding to the contract
	EmitLog2Filterer   // Log filterer for contract events
}

// EmitLog2Caller is an auto generated read-only Go binding around an Ethereum contract.
type EmitLog2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EmitLog2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type EmitLog2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EmitLog2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EmitLog2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EmitLog2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EmitLog2Session struct {
	Contract     *EmitLog2         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EmitLog2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EmitLog2CallerSession struct {
	Contract *EmitLog2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// EmitLog2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EmitLog2TransactorSession struct {
	Contract     *EmitLog2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// EmitLog2Raw is an auto generated low-level Go binding around an Ethereum contract.
type EmitLog2Raw struct {
	Contract *EmitLog2 // Generic contract binding to access the raw methods on
}

// EmitLog2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EmitLog2CallerRaw struct {
	Contract *EmitLog2Caller // Generic read-only contract binding to access the raw methods on
}

// EmitLog2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EmitLog2TransactorRaw struct {
	Contract *EmitLog2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewEmitLog2 creates a new instance of EmitLog2, bound to a specific deployed contract.
func NewEmitLog2(address common.Address, backend bind.ContractBackend) (*EmitLog2, error) {
	contract, err := bindEmitLog2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EmitLog2{EmitLog2Caller: EmitLog2Caller{contract: contract}, EmitLog2Transactor: EmitLog2Transactor{contract: contract}, EmitLog2Filterer: EmitLog2Filterer{contract: contract}}, nil
}

// NewEmitLog2Caller creates a new read-only instance of EmitLog2, bound to a specific deployed contract.
func NewEmitLog2Caller(address common.Address, caller bind.ContractCaller) (*EmitLog2Caller, error) {
	contract, err := bindEmitLog2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EmitLog2Caller{contract: contract}, nil
}

// NewEmitLog2Transactor creates a new write-only instance of EmitLog2, bound to a specific deployed contract.
func NewEmitLog2Transactor(address common.Address, transactor bind.ContractTransactor) (*EmitLog2Transactor, error) {
	contract, err := bindEmitLog2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EmitLog2Transactor{contract: contract}, nil
}

// NewEmitLog2Filterer creates a new log filterer instance of EmitLog2, bound to a specific deployed contract.
func NewEmitLog2Filterer(address common.Address, filterer bind.ContractFilterer) (*EmitLog2Filterer, error) {
	contract, err := bindEmitLog2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EmitLog2Filterer{contract: contract}, nil
}

// bindEmitLog2 binds a generic wrapper to an already deployed contract.
func bindEmitLog2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EmitLog2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EmitLog2 *EmitLog2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EmitLog2.Contract.EmitLog2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EmitLog2 *EmitLog2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EmitLog2.Contract.EmitLog2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EmitLog2 *EmitLog2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EmitLog2.Contract.EmitLog2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EmitLog2 *EmitLog2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EmitLog2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EmitLog2 *EmitLog2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EmitLog2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EmitLog2 *EmitLog2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EmitLog2.Contract.contract.Transact(opts, method, params...)
}

// EmitLogs is a paid mutator transaction binding the contract method 0x7966b4f6.
//
// Solidity: function emitLogs() returns()
func (_EmitLog2 *EmitLog2Transactor) EmitLogs(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EmitLog2.contract.Transact(opts, "emitLogs")
}

// EmitLogs is a paid mutator transaction binding the contract method 0x7966b4f6.
//
// Solidity: function emitLogs() returns()
func (_EmitLog2 *EmitLog2Session) EmitLogs() (*types.Transaction, error) {
	return _EmitLog2.Contract.EmitLogs(&_EmitLog2.TransactOpts)
}

// EmitLogs is a paid mutator transaction binding the contract method 0x7966b4f6.
//
// Solidity: function emitLogs() returns()
func (_EmitLog2 *EmitLog2TransactorSession) EmitLogs() (*types.Transaction, error) {
	return _EmitLog2.Contract.EmitLogs(&_EmitLog2.TransactOpts)
}

// EmitLog2LogAIterator is returned from FilterLogA and is used to iterate over the raw logs and unpacked data for LogA events raised by the EmitLog2 contract.
type EmitLog2LogAIterator struct {
	Event *EmitLog2LogA // Event containing the contract specifics and raw log

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
func (it *EmitLog2LogAIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EmitLog2LogA)
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
		it.Event = new(EmitLog2LogA)
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
func (it *EmitLog2LogAIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EmitLog2LogAIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EmitLog2LogA represents a LogA event raised by the EmitLog2 contract.
type EmitLog2LogA struct {
	A   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogA is a free log retrieval operation binding the contract event 0x977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac62.
//
// Solidity: event LogA(uint256 indexed a)
func (_EmitLog2 *EmitLog2Filterer) FilterLogA(opts *bind.FilterOpts, a []*big.Int) (*EmitLog2LogAIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _EmitLog2.contract.FilterLogs(opts, "LogA", aRule)
	if err != nil {
		return nil, err
	}
	return &EmitLog2LogAIterator{contract: _EmitLog2.contract, event: "LogA", logs: logs, sub: sub}, nil
}

// WatchLogA is a free log subscription operation binding the contract event 0x977224b24e70d33f3be87246a29c5636cfc8dd6853e175b54af01ff493ffac62.
//
// Solidity: event LogA(uint256 indexed a)
func (_EmitLog2 *EmitLog2Filterer) WatchLogA(opts *bind.WatchOpts, sink chan<- *EmitLog2LogA, a []*big.Int) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}

	logs, sub, err := _EmitLog2.contract.WatchLogs(opts, "LogA", aRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EmitLog2LogA)
				if err := _EmitLog2.contract.UnpackLog(event, "LogA", log); err != nil {
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
func (_EmitLog2 *EmitLog2Filterer) ParseLogA(log types.Log) (*EmitLog2LogA, error) {
	event := new(EmitLog2LogA)
	if err := _EmitLog2.contract.UnpackLog(event, "LogA", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
