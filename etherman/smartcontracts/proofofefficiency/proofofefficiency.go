// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package proofofefficiency

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

// ProofofefficiencyMetaData contains all meta data concerning the Proofofefficiency contract.
var ProofofefficiencyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractBridgeInterface\",\"name\":\"_bridge\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_matic\",\"type\":\"address\"},{\"internalType\":\"contractVerifierRollupInterface\",\"name\":\"_rollupVerifier\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"batchNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"}],\"name\":\"SendBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sequencerAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"}],\"name\":\"SetSequencer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"batchNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatch\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bridge\",\"outputs\":[{\"internalType\":\"contractBridgeInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateSequencerCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentLocalExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBatchSent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastGlobalExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastVerifiedBatch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"matic\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numSequencers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"}],\"name\":\"registerSequencer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupVerifier\",\"outputs\":[{\"internalType\":\"contractVerifierRollupInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maticAmount\",\"type\":\"uint256\"}],\"name\":\"sendBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"sentBatches\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sequencerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"batchL2HashData\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maticCollateral\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"sequencers\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainID\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"batchNum\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofA\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"proofB\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofC\",\"type\":\"uint256[2]\"}],\"name\":\"verifyBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ProofofefficiencyABI is the input ABI used to generate the binding from.
// Deprecated: Use ProofofefficiencyMetaData.ABI instead.
var ProofofefficiencyABI = ProofofefficiencyMetaData.ABI

// Proofofefficiency is an auto generated Go binding around an Ethereum contract.
type Proofofefficiency struct {
	ProofofefficiencyCaller     // Read-only binding to the contract
	ProofofefficiencyTransactor // Write-only binding to the contract
	ProofofefficiencyFilterer   // Log filterer for contract events
}

// ProofofefficiencyCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProofofefficiencyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofofefficiencyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProofofefficiencyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofofefficiencyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProofofefficiencyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofofefficiencySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProofofefficiencySession struct {
	Contract     *Proofofefficiency // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ProofofefficiencyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProofofefficiencyCallerSession struct {
	Contract *ProofofefficiencyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ProofofefficiencyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProofofefficiencyTransactorSession struct {
	Contract     *ProofofefficiencyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ProofofefficiencyRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProofofefficiencyRaw struct {
	Contract *Proofofefficiency // Generic contract binding to access the raw methods on
}

// ProofofefficiencyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProofofefficiencyCallerRaw struct {
	Contract *ProofofefficiencyCaller // Generic read-only contract binding to access the raw methods on
}

// ProofofefficiencyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProofofefficiencyTransactorRaw struct {
	Contract *ProofofefficiencyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProofofefficiency creates a new instance of Proofofefficiency, bound to a specific deployed contract.
func NewProofofefficiency(address common.Address, backend bind.ContractBackend) (*Proofofefficiency, error) {
	contract, err := bindProofofefficiency(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Proofofefficiency{ProofofefficiencyCaller: ProofofefficiencyCaller{contract: contract}, ProofofefficiencyTransactor: ProofofefficiencyTransactor{contract: contract}, ProofofefficiencyFilterer: ProofofefficiencyFilterer{contract: contract}}, nil
}

// NewProofofefficiencyCaller creates a new read-only instance of Proofofefficiency, bound to a specific deployed contract.
func NewProofofefficiencyCaller(address common.Address, caller bind.ContractCaller) (*ProofofefficiencyCaller, error) {
	contract, err := bindProofofefficiency(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyCaller{contract: contract}, nil
}

// NewProofofefficiencyTransactor creates a new write-only instance of Proofofefficiency, bound to a specific deployed contract.
func NewProofofefficiencyTransactor(address common.Address, transactor bind.ContractTransactor) (*ProofofefficiencyTransactor, error) {
	contract, err := bindProofofefficiency(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyTransactor{contract: contract}, nil
}

// NewProofofefficiencyFilterer creates a new log filterer instance of Proofofefficiency, bound to a specific deployed contract.
func NewProofofefficiencyFilterer(address common.Address, filterer bind.ContractFilterer) (*ProofofefficiencyFilterer, error) {
	contract, err := bindProofofefficiency(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyFilterer{contract: contract}, nil
}

// bindProofofefficiency binds a generic wrapper to an already deployed contract.
func bindProofofefficiency(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ProofofefficiencyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Proofofefficiency *ProofofefficiencyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Proofofefficiency.Contract.ProofofefficiencyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Proofofefficiency *ProofofefficiencyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.ProofofefficiencyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Proofofefficiency *ProofofefficiencyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.ProofofefficiencyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Proofofefficiency *ProofofefficiencyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Proofofefficiency.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Proofofefficiency *ProofofefficiencyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Proofofefficiency *ProofofefficiencyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.contract.Transact(opts, method, params...)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCaller) Bridge(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "bridge")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_Proofofefficiency *ProofofefficiencySession) Bridge() (common.Address, error) {
	return _Proofofefficiency.Contract.Bridge(&_Proofofefficiency.CallOpts)
}

// Bridge is a free data retrieval call binding the contract method 0xe78cea92.
//
// Solidity: function bridge() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCallerSession) Bridge() (common.Address, error) {
	return _Proofofefficiency.Contract.Bridge(&_Proofofefficiency.CallOpts)
}

// CalculateSequencerCollateral is a free data retrieval call binding the contract method 0x1dc125b2.
//
// Solidity: function calculateSequencerCollateral() pure returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCaller) CalculateSequencerCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "calculateSequencerCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateSequencerCollateral is a free data retrieval call binding the contract method 0x1dc125b2.
//
// Solidity: function calculateSequencerCollateral() pure returns(uint256)
func (_Proofofefficiency *ProofofefficiencySession) CalculateSequencerCollateral() (*big.Int, error) {
	return _Proofofefficiency.Contract.CalculateSequencerCollateral(&_Proofofefficiency.CallOpts)
}

// CalculateSequencerCollateral is a free data retrieval call binding the contract method 0x1dc125b2.
//
// Solidity: function calculateSequencerCollateral() pure returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCallerSession) CalculateSequencerCollateral() (*big.Int, error) {
	return _Proofofefficiency.Contract.CalculateSequencerCollateral(&_Proofofefficiency.CallOpts)
}

// CurrentLocalExitRoot is a free data retrieval call binding the contract method 0x959c2f47.
//
// Solidity: function currentLocalExitRoot() view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencyCaller) CurrentLocalExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "currentLocalExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CurrentLocalExitRoot is a free data retrieval call binding the contract method 0x959c2f47.
//
// Solidity: function currentLocalExitRoot() view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencySession) CurrentLocalExitRoot() ([32]byte, error) {
	return _Proofofefficiency.Contract.CurrentLocalExitRoot(&_Proofofefficiency.CallOpts)
}

// CurrentLocalExitRoot is a free data retrieval call binding the contract method 0x959c2f47.
//
// Solidity: function currentLocalExitRoot() view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencyCallerSession) CurrentLocalExitRoot() ([32]byte, error) {
	return _Proofofefficiency.Contract.CurrentLocalExitRoot(&_Proofofefficiency.CallOpts)
}

// CurrentStateRoot is a free data retrieval call binding the contract method 0xac2eba98.
//
// Solidity: function currentStateRoot() view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencyCaller) CurrentStateRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "currentStateRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CurrentStateRoot is a free data retrieval call binding the contract method 0xac2eba98.
//
// Solidity: function currentStateRoot() view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencySession) CurrentStateRoot() ([32]byte, error) {
	return _Proofofefficiency.Contract.CurrentStateRoot(&_Proofofefficiency.CallOpts)
}

// CurrentStateRoot is a free data retrieval call binding the contract method 0xac2eba98.
//
// Solidity: function currentStateRoot() view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencyCallerSession) CurrentStateRoot() ([32]byte, error) {
	return _Proofofefficiency.Contract.CurrentStateRoot(&_Proofofefficiency.CallOpts)
}

// LastBatchSent is a free data retrieval call binding the contract method 0x188ea07d.
//
// Solidity: function lastBatchSent() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCaller) LastBatchSent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "lastBatchSent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastBatchSent is a free data retrieval call binding the contract method 0x188ea07d.
//
// Solidity: function lastBatchSent() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencySession) LastBatchSent() (*big.Int, error) {
	return _Proofofefficiency.Contract.LastBatchSent(&_Proofofefficiency.CallOpts)
}

// LastBatchSent is a free data retrieval call binding the contract method 0x188ea07d.
//
// Solidity: function lastBatchSent() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCallerSession) LastBatchSent() (*big.Int, error) {
	return _Proofofefficiency.Contract.LastBatchSent(&_Proofofefficiency.CallOpts)
}

// LastGlobalExitRoot is a free data retrieval call binding the contract method 0x0acd922c.
//
// Solidity: function lastGlobalExitRoot() view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencyCaller) LastGlobalExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "lastGlobalExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastGlobalExitRoot is a free data retrieval call binding the contract method 0x0acd922c.
//
// Solidity: function lastGlobalExitRoot() view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencySession) LastGlobalExitRoot() ([32]byte, error) {
	return _Proofofefficiency.Contract.LastGlobalExitRoot(&_Proofofefficiency.CallOpts)
}

// LastGlobalExitRoot is a free data retrieval call binding the contract method 0x0acd922c.
//
// Solidity: function lastGlobalExitRoot() view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencyCallerSession) LastGlobalExitRoot() ([32]byte, error) {
	return _Proofofefficiency.Contract.LastGlobalExitRoot(&_Proofofefficiency.CallOpts)
}

// LastVerifiedBatch is a free data retrieval call binding the contract method 0x7fcb3653.
//
// Solidity: function lastVerifiedBatch() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCaller) LastVerifiedBatch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "lastVerifiedBatch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastVerifiedBatch is a free data retrieval call binding the contract method 0x7fcb3653.
//
// Solidity: function lastVerifiedBatch() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencySession) LastVerifiedBatch() (*big.Int, error) {
	return _Proofofefficiency.Contract.LastVerifiedBatch(&_Proofofefficiency.CallOpts)
}

// LastVerifiedBatch is a free data retrieval call binding the contract method 0x7fcb3653.
//
// Solidity: function lastVerifiedBatch() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCallerSession) LastVerifiedBatch() (*big.Int, error) {
	return _Proofofefficiency.Contract.LastVerifiedBatch(&_Proofofefficiency.CallOpts)
}

// Matic is a free data retrieval call binding the contract method 0xb6b0b097.
//
// Solidity: function matic() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCaller) Matic(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "matic")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Matic is a free data retrieval call binding the contract method 0xb6b0b097.
//
// Solidity: function matic() view returns(address)
func (_Proofofefficiency *ProofofefficiencySession) Matic() (common.Address, error) {
	return _Proofofefficiency.Contract.Matic(&_Proofofefficiency.CallOpts)
}

// Matic is a free data retrieval call binding the contract method 0xb6b0b097.
//
// Solidity: function matic() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCallerSession) Matic() (common.Address, error) {
	return _Proofofefficiency.Contract.Matic(&_Proofofefficiency.CallOpts)
}

// NumSequencers is a free data retrieval call binding the contract method 0xca98a308.
//
// Solidity: function numSequencers() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCaller) NumSequencers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "numSequencers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumSequencers is a free data retrieval call binding the contract method 0xca98a308.
//
// Solidity: function numSequencers() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencySession) NumSequencers() (*big.Int, error) {
	return _Proofofefficiency.Contract.NumSequencers(&_Proofofefficiency.CallOpts)
}

// NumSequencers is a free data retrieval call binding the contract method 0xca98a308.
//
// Solidity: function numSequencers() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCallerSession) NumSequencers() (*big.Int, error) {
	return _Proofofefficiency.Contract.NumSequencers(&_Proofofefficiency.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Proofofefficiency *ProofofefficiencySession) Owner() (common.Address, error) {
	return _Proofofefficiency.Contract.Owner(&_Proofofefficiency.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCallerSession) Owner() (common.Address, error) {
	return _Proofofefficiency.Contract.Owner(&_Proofofefficiency.CallOpts)
}

// RollupVerifier is a free data retrieval call binding the contract method 0xe8bf92ed.
//
// Solidity: function rollupVerifier() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCaller) RollupVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "rollupVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RollupVerifier is a free data retrieval call binding the contract method 0xe8bf92ed.
//
// Solidity: function rollupVerifier() view returns(address)
func (_Proofofefficiency *ProofofefficiencySession) RollupVerifier() (common.Address, error) {
	return _Proofofefficiency.Contract.RollupVerifier(&_Proofofefficiency.CallOpts)
}

// RollupVerifier is a free data retrieval call binding the contract method 0xe8bf92ed.
//
// Solidity: function rollupVerifier() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCallerSession) RollupVerifier() (common.Address, error) {
	return _Proofofefficiency.Contract.RollupVerifier(&_Proofofefficiency.CallOpts)
}

// SentBatches is a free data retrieval call binding the contract method 0x84b0be8c.
//
// Solidity: function sentBatches(uint256 ) view returns(address sequencerAddress, bytes32 batchL2HashData, uint256 maticCollateral)
func (_Proofofefficiency *ProofofefficiencyCaller) SentBatches(opts *bind.CallOpts, arg0 *big.Int) (struct {
	SequencerAddress common.Address
	BatchL2HashData  [32]byte
	MaticCollateral  *big.Int
}, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "sentBatches", arg0)

	outstruct := new(struct {
		SequencerAddress common.Address
		BatchL2HashData  [32]byte
		MaticCollateral  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SequencerAddress = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.BatchL2HashData = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.MaticCollateral = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// SentBatches is a free data retrieval call binding the contract method 0x84b0be8c.
//
// Solidity: function sentBatches(uint256 ) view returns(address sequencerAddress, bytes32 batchL2HashData, uint256 maticCollateral)
func (_Proofofefficiency *ProofofefficiencySession) SentBatches(arg0 *big.Int) (struct {
	SequencerAddress common.Address
	BatchL2HashData  [32]byte
	MaticCollateral  *big.Int
}, error) {
	return _Proofofefficiency.Contract.SentBatches(&_Proofofefficiency.CallOpts, arg0)
}

// SentBatches is a free data retrieval call binding the contract method 0x84b0be8c.
//
// Solidity: function sentBatches(uint256 ) view returns(address sequencerAddress, bytes32 batchL2HashData, uint256 maticCollateral)
func (_Proofofefficiency *ProofofefficiencyCallerSession) SentBatches(arg0 *big.Int) (struct {
	SequencerAddress common.Address
	BatchL2HashData  [32]byte
	MaticCollateral  *big.Int
}, error) {
	return _Proofofefficiency.Contract.SentBatches(&_Proofofefficiency.CallOpts, arg0)
}

// Sequencers is a free data retrieval call binding the contract method 0x1c7a07ee.
//
// Solidity: function sequencers(address ) view returns(string sequencerURL, uint256 chainID)
func (_Proofofefficiency *ProofofefficiencyCaller) Sequencers(opts *bind.CallOpts, arg0 common.Address) (struct {
	SequencerURL string
	ChainID      *big.Int
}, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "sequencers", arg0)

	outstruct := new(struct {
		SequencerURL string
		ChainID      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SequencerURL = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.ChainID = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Sequencers is a free data retrieval call binding the contract method 0x1c7a07ee.
//
// Solidity: function sequencers(address ) view returns(string sequencerURL, uint256 chainID)
func (_Proofofefficiency *ProofofefficiencySession) Sequencers(arg0 common.Address) (struct {
	SequencerURL string
	ChainID      *big.Int
}, error) {
	return _Proofofefficiency.Contract.Sequencers(&_Proofofefficiency.CallOpts, arg0)
}

// Sequencers is a free data retrieval call binding the contract method 0x1c7a07ee.
//
// Solidity: function sequencers(address ) view returns(string sequencerURL, uint256 chainID)
func (_Proofofefficiency *ProofofefficiencyCallerSession) Sequencers(arg0 common.Address) (struct {
	SequencerURL string
	ChainID      *big.Int
}, error) {
	return _Proofofefficiency.Contract.Sequencers(&_Proofofefficiency.CallOpts, arg0)
}

// RegisterSequencer is a paid mutator transaction binding the contract method 0x8a4abab8.
//
// Solidity: function registerSequencer(string sequencerURL) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) RegisterSequencer(opts *bind.TransactOpts, sequencerURL string) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "registerSequencer", sequencerURL)
}

// RegisterSequencer is a paid mutator transaction binding the contract method 0x8a4abab8.
//
// Solidity: function registerSequencer(string sequencerURL) returns()
func (_Proofofefficiency *ProofofefficiencySession) RegisterSequencer(sequencerURL string) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.RegisterSequencer(&_Proofofefficiency.TransactOpts, sequencerURL)
}

// RegisterSequencer is a paid mutator transaction binding the contract method 0x8a4abab8.
//
// Solidity: function registerSequencer(string sequencerURL) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) RegisterSequencer(sequencerURL string) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.RegisterSequencer(&_Proofofefficiency.TransactOpts, sequencerURL)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Proofofefficiency *ProofofefficiencySession) RenounceOwnership() (*types.Transaction, error) {
	return _Proofofefficiency.Contract.RenounceOwnership(&_Proofofefficiency.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Proofofefficiency.Contract.RenounceOwnership(&_Proofofefficiency.TransactOpts)
}

// SendBatch is a paid mutator transaction binding the contract method 0x06d6490f.
//
// Solidity: function sendBatch(bytes transactions, uint256 maticAmount) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) SendBatch(opts *bind.TransactOpts, transactions []byte, maticAmount *big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "sendBatch", transactions, maticAmount)
}

// SendBatch is a paid mutator transaction binding the contract method 0x06d6490f.
//
// Solidity: function sendBatch(bytes transactions, uint256 maticAmount) returns()
func (_Proofofefficiency *ProofofefficiencySession) SendBatch(transactions []byte, maticAmount *big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SendBatch(&_Proofofefficiency.TransactOpts, transactions, maticAmount)
}

// SendBatch is a paid mutator transaction binding the contract method 0x06d6490f.
//
// Solidity: function sendBatch(bytes transactions, uint256 maticAmount) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) SendBatch(transactions []byte, maticAmount *big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SendBatch(&_Proofofefficiency.TransactOpts, transactions, maticAmount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Proofofefficiency *ProofofefficiencySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.TransferOwnership(&_Proofofefficiency.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.TransferOwnership(&_Proofofefficiency.TransactOpts, newOwner)
}

// VerifyBatch is a paid mutator transaction binding the contract method 0xa152b62e.
//
// Solidity: function verifyBatch(bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256 batchNum, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) VerifyBatch(opts *bind.TransactOpts, newLocalExitRoot [32]byte, newStateRoot [32]byte, batchNum *big.Int, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "verifyBatch", newLocalExitRoot, newStateRoot, batchNum, proofA, proofB, proofC)
}

// VerifyBatch is a paid mutator transaction binding the contract method 0xa152b62e.
//
// Solidity: function verifyBatch(bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256 batchNum, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencySession) VerifyBatch(newLocalExitRoot [32]byte, newStateRoot [32]byte, batchNum *big.Int, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.VerifyBatch(&_Proofofefficiency.TransactOpts, newLocalExitRoot, newStateRoot, batchNum, proofA, proofB, proofC)
}

// VerifyBatch is a paid mutator transaction binding the contract method 0xa152b62e.
//
// Solidity: function verifyBatch(bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256 batchNum, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) VerifyBatch(newLocalExitRoot [32]byte, newStateRoot [32]byte, batchNum *big.Int, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.VerifyBatch(&_Proofofefficiency.TransactOpts, newLocalExitRoot, newStateRoot, batchNum, proofA, proofB, proofC)
}

// ProofofefficiencyOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Proofofefficiency contract.
type ProofofefficiencyOwnershipTransferredIterator struct {
	Event *ProofofefficiencyOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencyOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencyOwnershipTransferred)
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
		it.Event = new(ProofofefficiencyOwnershipTransferred)
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
func (it *ProofofefficiencyOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencyOwnershipTransferred represents a OwnershipTransferred event raised by the Proofofefficiency contract.
type ProofofefficiencyOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ProofofefficiencyOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyOwnershipTransferredIterator{contract: _Proofofefficiency.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencyOwnershipTransferred)
				if err := _Proofofefficiency.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseOwnershipTransferred(log types.Log) (*ProofofefficiencyOwnershipTransferred, error) {
	event := new(ProofofefficiencyOwnershipTransferred)
	if err := _Proofofefficiency.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencySendBatchIterator is returned from FilterSendBatch and is used to iterate over the raw logs and unpacked data for SendBatch events raised by the Proofofefficiency contract.
type ProofofefficiencySendBatchIterator struct {
	Event *ProofofefficiencySendBatch // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencySendBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencySendBatch)
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
		it.Event = new(ProofofefficiencySendBatch)
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
func (it *ProofofefficiencySendBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencySendBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencySendBatch represents a SendBatch event raised by the Proofofefficiency contract.
type ProofofefficiencySendBatch struct {
	BatchNum  *big.Int
	Sequencer common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSendBatch is a free log retrieval operation binding the contract event 0xe99ca860ab96f187daa5c05e5d805b16a9dd0242a7e113c065d2a13ce30159d0.
//
// Solidity: event SendBatch(uint256 indexed batchNum, address indexed sequencer)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSendBatch(opts *bind.FilterOpts, batchNum []*big.Int, sequencer []common.Address) (*ProofofefficiencySendBatchIterator, error) {

	var batchNumRule []interface{}
	for _, batchNumItem := range batchNum {
		batchNumRule = append(batchNumRule, batchNumItem)
	}
	var sequencerRule []interface{}
	for _, sequencerItem := range sequencer {
		sequencerRule = append(sequencerRule, sequencerItem)
	}

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "SendBatch", batchNumRule, sequencerRule)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencySendBatchIterator{contract: _Proofofefficiency.contract, event: "SendBatch", logs: logs, sub: sub}, nil
}

// WatchSendBatch is a free log subscription operation binding the contract event 0xe99ca860ab96f187daa5c05e5d805b16a9dd0242a7e113c065d2a13ce30159d0.
//
// Solidity: event SendBatch(uint256 indexed batchNum, address indexed sequencer)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSendBatch(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySendBatch, batchNum []*big.Int, sequencer []common.Address) (event.Subscription, error) {

	var batchNumRule []interface{}
	for _, batchNumItem := range batchNum {
		batchNumRule = append(batchNumRule, batchNumItem)
	}
	var sequencerRule []interface{}
	for _, sequencerItem := range sequencer {
		sequencerRule = append(sequencerRule, sequencerItem)
	}

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "SendBatch", batchNumRule, sequencerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencySendBatch)
				if err := _Proofofefficiency.contract.UnpackLog(event, "SendBatch", log); err != nil {
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

// ParseSendBatch is a log parse operation binding the contract event 0xe99ca860ab96f187daa5c05e5d805b16a9dd0242a7e113c065d2a13ce30159d0.
//
// Solidity: event SendBatch(uint256 indexed batchNum, address indexed sequencer)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseSendBatch(log types.Log) (*ProofofefficiencySendBatch, error) {
	event := new(ProofofefficiencySendBatch)
	if err := _Proofofefficiency.contract.UnpackLog(event, "SendBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencySetSequencerIterator is returned from FilterSetSequencer and is used to iterate over the raw logs and unpacked data for SetSequencer events raised by the Proofofefficiency contract.
type ProofofefficiencySetSequencerIterator struct {
	Event *ProofofefficiencySetSequencer // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencySetSequencerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencySetSequencer)
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
		it.Event = new(ProofofefficiencySetSequencer)
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
func (it *ProofofefficiencySetSequencerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencySetSequencerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencySetSequencer represents a SetSequencer event raised by the Proofofefficiency contract.
type ProofofefficiencySetSequencer struct {
	SequencerAddress common.Address
	SequencerURL     string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSetSequencer is a free log retrieval operation binding the contract event 0x3c5b8bb3cdafd1a15eaa861069e1567306bab24615ef281949384e54b80bd77d.
//
// Solidity: event SetSequencer(address sequencerAddress, string sequencerURL)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSetSequencer(opts *bind.FilterOpts) (*ProofofefficiencySetSequencerIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "SetSequencer")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencySetSequencerIterator{contract: _Proofofefficiency.contract, event: "SetSequencer", logs: logs, sub: sub}, nil
}

// WatchSetSequencer is a free log subscription operation binding the contract event 0x3c5b8bb3cdafd1a15eaa861069e1567306bab24615ef281949384e54b80bd77d.
//
// Solidity: event SetSequencer(address sequencerAddress, string sequencerURL)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSetSequencer(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySetSequencer) (event.Subscription, error) {

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "SetSequencer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencySetSequencer)
				if err := _Proofofefficiency.contract.UnpackLog(event, "SetSequencer", log); err != nil {
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

// ParseSetSequencer is a log parse operation binding the contract event 0x3c5b8bb3cdafd1a15eaa861069e1567306bab24615ef281949384e54b80bd77d.
//
// Solidity: event SetSequencer(address sequencerAddress, string sequencerURL)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseSetSequencer(log types.Log) (*ProofofefficiencySetSequencer, error) {
	event := new(ProofofefficiencySetSequencer)
	if err := _Proofofefficiency.contract.UnpackLog(event, "SetSequencer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencyVerifyBatchIterator is returned from FilterVerifyBatch and is used to iterate over the raw logs and unpacked data for VerifyBatch events raised by the Proofofefficiency contract.
type ProofofefficiencyVerifyBatchIterator struct {
	Event *ProofofefficiencyVerifyBatch // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencyVerifyBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencyVerifyBatch)
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
		it.Event = new(ProofofefficiencyVerifyBatch)
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
func (it *ProofofefficiencyVerifyBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencyVerifyBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencyVerifyBatch represents a VerifyBatch event raised by the Proofofefficiency contract.
type ProofofefficiencyVerifyBatch struct {
	BatchNum   *big.Int
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVerifyBatch is a free log retrieval operation binding the contract event 0x60f152ed8e8084aa9231a20aae908de8d039a671f90e86d6b92c87889d952b22.
//
// Solidity: event VerifyBatch(uint256 indexed batchNum, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterVerifyBatch(opts *bind.FilterOpts, batchNum []*big.Int, aggregator []common.Address) (*ProofofefficiencyVerifyBatchIterator, error) {

	var batchNumRule []interface{}
	for _, batchNumItem := range batchNum {
		batchNumRule = append(batchNumRule, batchNumItem)
	}
	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "VerifyBatch", batchNumRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyVerifyBatchIterator{contract: _Proofofefficiency.contract, event: "VerifyBatch", logs: logs, sub: sub}, nil
}

// WatchVerifyBatch is a free log subscription operation binding the contract event 0x60f152ed8e8084aa9231a20aae908de8d039a671f90e86d6b92c87889d952b22.
//
// Solidity: event VerifyBatch(uint256 indexed batchNum, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchVerifyBatch(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyVerifyBatch, batchNum []*big.Int, aggregator []common.Address) (event.Subscription, error) {

	var batchNumRule []interface{}
	for _, batchNumItem := range batchNum {
		batchNumRule = append(batchNumRule, batchNumItem)
	}
	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "VerifyBatch", batchNumRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencyVerifyBatch)
				if err := _Proofofefficiency.contract.UnpackLog(event, "VerifyBatch", log); err != nil {
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

// ParseVerifyBatch is a log parse operation binding the contract event 0x60f152ed8e8084aa9231a20aae908de8d039a671f90e86d6b92c87889d952b22.
//
// Solidity: event VerifyBatch(uint256 indexed batchNum, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseVerifyBatch(log types.Log) (*ProofofefficiencyVerifyBatch, error) {
	event := new(ProofofefficiencyVerifyBatch)
	if err := _Proofofefficiency.contract.UnpackLog(event, "VerifyBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
