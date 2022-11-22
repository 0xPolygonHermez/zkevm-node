// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package globalexitrootmanager

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

// GlobalexitrootmanagerMetaData contains all meta data concerning the Globalexitrootmanager contract.
var GlobalexitrootmanagerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"mainnetExitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"rollupExitRoot\",\"type\":\"bytes32\"}],\"name\":\"UpdateGlobalExitRoot\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastGlobalExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"globalExitRootMap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rollupAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMainnetExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newRoot\",\"type\":\"bytes32\"}],\"name\":\"updateExitRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506109fc806100206000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80633ed691ef1161005b5780633ed691ef14610115578063485cc955146101335780635ec6a8df1461014f578063a3c573eb1461016d57610088565b806301fd90441461008d578063257b3632146100ab578063319cf735146100db57806333d6247d146100f9575b600080fd5b61009561018b565b6040516100a29190610662565b60405180910390f35b6100c560048036038101906100c091906106ae565b610191565b6040516100d291906106f4565b60405180910390f35b6100e36101a9565b6040516100f09190610662565b60405180910390f35b610113600480360381019061010e91906106ae565b6101af565b005b61011d6103ec565b60405161012a9190610662565b60405180910390f35b61014d6004803603810190610148919061076d565b610420565b005b6101576105da565b60405161016491906107bc565b60405180910390f35b610175610600565b60405161018291906107bc565b60405180910390f35b60015481565b60036020528060005260406000206000915090505481565b60025481565b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806102585750600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b610297576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161028e9061085a565b60405180910390fd5b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16036102f457806001819055505b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff160361035157806002819055505b600060025460015460405160200161036a92919061089b565b60405160208183030381529060405280519060200120905060006003600083815260200190815260200160002054036103e8574260036000838152602001908152602001600020819055506001546002547f61014378f82a0d809aefaf87a8ac9505b89c321808287a6e7810f29304c1fce360405160405180910390a35b5050565b600060025460015460405160200161040592919061089b565b60405160208183030381529060405280519060200120905090565b60008060019054906101000a900460ff161590508080156104515750600160008054906101000a900460ff1660ff16105b8061047e575061046030610626565b15801561047d5750600160008054906101000a900460ff1660ff16145b5b6104bd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104b490610939565b60405180910390fd5b60016000806101000a81548160ff021916908360ff16021790555080156104fa576001600060016101000a81548160ff0219169083151502179055505b82600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555081600460006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080156105d55760008060016101000a81548160ff0219169083151502179055507f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb384740249860016040516105cc91906109ab565b60405180910390a15b505050565b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b6000819050919050565b61065c81610649565b82525050565b60006020820190506106776000830184610653565b92915050565b600080fd5b61068b81610649565b811461069657600080fd5b50565b6000813590506106a881610682565b92915050565b6000602082840312156106c4576106c361067d565b5b60006106d284828501610699565b91505092915050565b6000819050919050565b6106ee816106db565b82525050565b600060208201905061070960008301846106e5565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061073a8261070f565b9050919050565b61074a8161072f565b811461075557600080fd5b50565b60008135905061076781610741565b92915050565b600080604083850312156107845761078361067d565b5b600061079285828601610758565b92505060206107a385828601610758565b9150509250929050565b6107b68161072f565b82525050565b60006020820190506107d160008301846107ad565b92915050565b600082825260208201905092915050565b7f476c6f62616c45786974526f6f744d616e616765723a3a75706461746545786960008201527f74526f6f743a204f4e4c595f414c4c4f5745445f434f4e545241435453000000602082015250565b6000610844603d836107d7565b915061084f826107e8565b604082019050919050565b6000602082019050818103600083015261087381610837565b9050919050565b6000819050919050565b61089561089082610649565b61087a565b82525050565b60006108a78285610884565b6020820191506108b78284610884565b6020820191508190509392505050565b7f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160008201527f647920696e697469616c697a6564000000000000000000000000000000000000602082015250565b6000610923602e836107d7565b915061092e826108c7565b604082019050919050565b6000602082019050818103600083015261095281610916565b9050919050565b6000819050919050565b600060ff82169050919050565b6000819050919050565b600061099561099061098b84610959565b610970565b610963565b9050919050565b6109a58161097a565b82525050565b60006020820190506109c0600083018461099c565b9291505056fea2646970667358221220ce90ec27b6a17fb89184bc1febbcbe40c138e347ca03dfb64936c2918411dbe864736f6c634300080f0033",
}

// GlobalexitrootmanagerABI is the input ABI used to generate the binding from.
// Deprecated: Use GlobalexitrootmanagerMetaData.ABI instead.
var GlobalexitrootmanagerABI = GlobalexitrootmanagerMetaData.ABI

// GlobalexitrootmanagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GlobalexitrootmanagerMetaData.Bin instead.
var GlobalexitrootmanagerBin = GlobalexitrootmanagerMetaData.Bin

// DeployGlobalexitrootmanager deploys a new Ethereum contract, binding an instance of Globalexitrootmanager to it.
func DeployGlobalexitrootmanager(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Globalexitrootmanager, error) {
	parsed, err := GlobalexitrootmanagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GlobalexitrootmanagerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Globalexitrootmanager{GlobalexitrootmanagerCaller: GlobalexitrootmanagerCaller{contract: contract}, GlobalexitrootmanagerTransactor: GlobalexitrootmanagerTransactor{contract: contract}, GlobalexitrootmanagerFilterer: GlobalexitrootmanagerFilterer{contract: contract}}, nil
}

// Globalexitrootmanager is an auto generated Go binding around an Ethereum contract.
type Globalexitrootmanager struct {
	GlobalexitrootmanagerCaller     // Read-only binding to the contract
	GlobalexitrootmanagerTransactor // Write-only binding to the contract
	GlobalexitrootmanagerFilterer   // Log filterer for contract events
}

// GlobalexitrootmanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type GlobalexitrootmanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GlobalexitrootmanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GlobalexitrootmanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GlobalexitrootmanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GlobalexitrootmanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GlobalexitrootmanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GlobalexitrootmanagerSession struct {
	Contract     *Globalexitrootmanager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// GlobalexitrootmanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GlobalexitrootmanagerCallerSession struct {
	Contract *GlobalexitrootmanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// GlobalexitrootmanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GlobalexitrootmanagerTransactorSession struct {
	Contract     *GlobalexitrootmanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// GlobalexitrootmanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type GlobalexitrootmanagerRaw struct {
	Contract *Globalexitrootmanager // Generic contract binding to access the raw methods on
}

// GlobalexitrootmanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GlobalexitrootmanagerCallerRaw struct {
	Contract *GlobalexitrootmanagerCaller // Generic read-only contract binding to access the raw methods on
}

// GlobalexitrootmanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GlobalexitrootmanagerTransactorRaw struct {
	Contract *GlobalexitrootmanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGlobalexitrootmanager creates a new instance of Globalexitrootmanager, bound to a specific deployed contract.
func NewGlobalexitrootmanager(address common.Address, backend bind.ContractBackend) (*Globalexitrootmanager, error) {
	contract, err := bindGlobalexitrootmanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Globalexitrootmanager{GlobalexitrootmanagerCaller: GlobalexitrootmanagerCaller{contract: contract}, GlobalexitrootmanagerTransactor: GlobalexitrootmanagerTransactor{contract: contract}, GlobalexitrootmanagerFilterer: GlobalexitrootmanagerFilterer{contract: contract}}, nil
}

// NewGlobalexitrootmanagerCaller creates a new read-only instance of Globalexitrootmanager, bound to a specific deployed contract.
func NewGlobalexitrootmanagerCaller(address common.Address, caller bind.ContractCaller) (*GlobalexitrootmanagerCaller, error) {
	contract, err := bindGlobalexitrootmanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GlobalexitrootmanagerCaller{contract: contract}, nil
}

// NewGlobalexitrootmanagerTransactor creates a new write-only instance of Globalexitrootmanager, bound to a specific deployed contract.
func NewGlobalexitrootmanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*GlobalexitrootmanagerTransactor, error) {
	contract, err := bindGlobalexitrootmanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GlobalexitrootmanagerTransactor{contract: contract}, nil
}

// NewGlobalexitrootmanagerFilterer creates a new log filterer instance of Globalexitrootmanager, bound to a specific deployed contract.
func NewGlobalexitrootmanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*GlobalexitrootmanagerFilterer, error) {
	contract, err := bindGlobalexitrootmanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GlobalexitrootmanagerFilterer{contract: contract}, nil
}

// bindGlobalexitrootmanager binds a generic wrapper to an already deployed contract.
func bindGlobalexitrootmanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GlobalexitrootmanagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Globalexitrootmanager *GlobalexitrootmanagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Globalexitrootmanager.Contract.GlobalexitrootmanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Globalexitrootmanager *GlobalexitrootmanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.GlobalexitrootmanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Globalexitrootmanager *GlobalexitrootmanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.GlobalexitrootmanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Globalexitrootmanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.contract.Transact(opts, method, params...)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) BridgeAddress() (common.Address, error) {
	return _Globalexitrootmanager.Contract.BridgeAddress(&_Globalexitrootmanager.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) BridgeAddress() (common.Address, error) {
	return _Globalexitrootmanager.Contract.BridgeAddress(&_Globalexitrootmanager.CallOpts)
}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) GetLastGlobalExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "getLastGlobalExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) GetLastGlobalExitRoot() ([32]byte, error) {
	return _Globalexitrootmanager.Contract.GetLastGlobalExitRoot(&_Globalexitrootmanager.CallOpts)
}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) GetLastGlobalExitRoot() ([32]byte, error) {
	return _Globalexitrootmanager.Contract.GetLastGlobalExitRoot(&_Globalexitrootmanager.CallOpts)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) GlobalExitRootMap(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "globalExitRootMap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Globalexitrootmanager.Contract.GlobalExitRootMap(&_Globalexitrootmanager.CallOpts, arg0)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Globalexitrootmanager.Contract.GlobalExitRootMap(&_Globalexitrootmanager.CallOpts, arg0)
}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) LastMainnetExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "lastMainnetExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) LastMainnetExitRoot() ([32]byte, error) {
	return _Globalexitrootmanager.Contract.LastMainnetExitRoot(&_Globalexitrootmanager.CallOpts)
}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) LastMainnetExitRoot() ([32]byte, error) {
	return _Globalexitrootmanager.Contract.LastMainnetExitRoot(&_Globalexitrootmanager.CallOpts)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) LastRollupExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "lastRollupExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) LastRollupExitRoot() ([32]byte, error) {
	return _Globalexitrootmanager.Contract.LastRollupExitRoot(&_Globalexitrootmanager.CallOpts)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) LastRollupExitRoot() ([32]byte, error) {
	return _Globalexitrootmanager.Contract.LastRollupExitRoot(&_Globalexitrootmanager.CallOpts)
}

// RollupAddress is a free data retrieval call binding the contract method 0x5ec6a8df.
//
// Solidity: function rollupAddress() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) RollupAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "rollupAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RollupAddress is a free data retrieval call binding the contract method 0x5ec6a8df.
//
// Solidity: function rollupAddress() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) RollupAddress() (common.Address, error) {
	return _Globalexitrootmanager.Contract.RollupAddress(&_Globalexitrootmanager.CallOpts)
}

// RollupAddress is a free data retrieval call binding the contract method 0x5ec6a8df.
//
// Solidity: function rollupAddress() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) RollupAddress() (common.Address, error) {
	return _Globalexitrootmanager.Contract.RollupAddress(&_Globalexitrootmanager.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _rollupAddress, address _bridgeAddress) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactor) Initialize(opts *bind.TransactOpts, _rollupAddress common.Address, _bridgeAddress common.Address) (*types.Transaction, error) {
	return _Globalexitrootmanager.contract.Transact(opts, "initialize", _rollupAddress, _bridgeAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _rollupAddress, address _bridgeAddress) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) Initialize(_rollupAddress common.Address, _bridgeAddress common.Address) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.Initialize(&_Globalexitrootmanager.TransactOpts, _rollupAddress, _bridgeAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _rollupAddress, address _bridgeAddress) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactorSession) Initialize(_rollupAddress common.Address, _bridgeAddress common.Address) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.Initialize(&_Globalexitrootmanager.TransactOpts, _rollupAddress, _bridgeAddress)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactor) UpdateExitRoot(opts *bind.TransactOpts, newRoot [32]byte) (*types.Transaction, error) {
	return _Globalexitrootmanager.contract.Transact(opts, "updateExitRoot", newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.UpdateExitRoot(&_Globalexitrootmanager.TransactOpts, newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactorSession) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.UpdateExitRoot(&_Globalexitrootmanager.TransactOpts, newRoot)
}

// GlobalexitrootmanagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Globalexitrootmanager contract.
type GlobalexitrootmanagerInitializedIterator struct {
	Event *GlobalexitrootmanagerInitialized // Event containing the contract specifics and raw log

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
func (it *GlobalexitrootmanagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GlobalexitrootmanagerInitialized)
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
		it.Event = new(GlobalexitrootmanagerInitialized)
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
func (it *GlobalexitrootmanagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GlobalexitrootmanagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GlobalexitrootmanagerInitialized represents a Initialized event raised by the Globalexitrootmanager contract.
type GlobalexitrootmanagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*GlobalexitrootmanagerInitializedIterator, error) {

	logs, sub, err := _Globalexitrootmanager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GlobalexitrootmanagerInitializedIterator{contract: _Globalexitrootmanager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GlobalexitrootmanagerInitialized) (event.Subscription, error) {

	logs, sub, err := _Globalexitrootmanager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GlobalexitrootmanagerInitialized)
				if err := _Globalexitrootmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) ParseInitialized(log types.Log) (*GlobalexitrootmanagerInitialized, error) {
	event := new(GlobalexitrootmanagerInitialized)
	if err := _Globalexitrootmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GlobalexitrootmanagerUpdateGlobalExitRootIterator is returned from FilterUpdateGlobalExitRoot and is used to iterate over the raw logs and unpacked data for UpdateGlobalExitRoot events raised by the Globalexitrootmanager contract.
type GlobalexitrootmanagerUpdateGlobalExitRootIterator struct {
	Event *GlobalexitrootmanagerUpdateGlobalExitRoot // Event containing the contract specifics and raw log

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
func (it *GlobalexitrootmanagerUpdateGlobalExitRootIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GlobalexitrootmanagerUpdateGlobalExitRoot)
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
		it.Event = new(GlobalexitrootmanagerUpdateGlobalExitRoot)
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
func (it *GlobalexitrootmanagerUpdateGlobalExitRootIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GlobalexitrootmanagerUpdateGlobalExitRootIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GlobalexitrootmanagerUpdateGlobalExitRoot represents a UpdateGlobalExitRoot event raised by the Globalexitrootmanager contract.
type GlobalexitrootmanagerUpdateGlobalExitRoot struct {
	MainnetExitRoot [32]byte
	RollupExitRoot  [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpdateGlobalExitRoot is a free log retrieval operation binding the contract event 0x61014378f82a0d809aefaf87a8ac9505b89c321808287a6e7810f29304c1fce3.
//
// Solidity: event UpdateGlobalExitRoot(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) FilterUpdateGlobalExitRoot(opts *bind.FilterOpts, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (*GlobalexitrootmanagerUpdateGlobalExitRootIterator, error) {

	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Globalexitrootmanager.contract.FilterLogs(opts, "UpdateGlobalExitRoot", mainnetExitRootRule, rollupExitRootRule)
	if err != nil {
		return nil, err
	}
	return &GlobalexitrootmanagerUpdateGlobalExitRootIterator{contract: _Globalexitrootmanager.contract, event: "UpdateGlobalExitRoot", logs: logs, sub: sub}, nil
}

// WatchUpdateGlobalExitRoot is a free log subscription operation binding the contract event 0x61014378f82a0d809aefaf87a8ac9505b89c321808287a6e7810f29304c1fce3.
//
// Solidity: event UpdateGlobalExitRoot(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) WatchUpdateGlobalExitRoot(opts *bind.WatchOpts, sink chan<- *GlobalexitrootmanagerUpdateGlobalExitRoot, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (event.Subscription, error) {

	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Globalexitrootmanager.contract.WatchLogs(opts, "UpdateGlobalExitRoot", mainnetExitRootRule, rollupExitRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GlobalexitrootmanagerUpdateGlobalExitRoot)
				if err := _Globalexitrootmanager.contract.UnpackLog(event, "UpdateGlobalExitRoot", log); err != nil {
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

// ParseUpdateGlobalExitRoot is a log parse operation binding the contract event 0x61014378f82a0d809aefaf87a8ac9505b89c321808287a6e7810f29304c1fce3.
//
// Solidity: event UpdateGlobalExitRoot(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) ParseUpdateGlobalExitRoot(log types.Log) (*GlobalexitrootmanagerUpdateGlobalExitRoot, error) {
	event := new(GlobalexitrootmanagerUpdateGlobalExitRoot)
	if err := _Globalexitrootmanager.contract.UnpackLog(event, "UpdateGlobalExitRoot", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
