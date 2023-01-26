// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package polygonzkevmglobalexitroot

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

// PolygonzkevmglobalexitrootMetaData contains all meta data concerning the Polygonzkevmglobalexitroot contract.
var PolygonzkevmglobalexitrootMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"mainnetExitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"rollupExitRoot\",\"type\":\"bytes32\"}],\"name\":\"UpdateGlobalExitRoot\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastGlobalExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"globalExitRootMap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rollupAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMainnetExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newRoot\",\"type\":\"bytes32\"}],\"name\":\"updateExitRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506104f3806100206000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80633ed691ef1161005b5780633ed691ef146100e7578063485cc955146101205780635ec6a8df14610133578063a3c573eb1461015e57600080fd5b806301fd90441461008d578063257b3632146100a9578063319cf735146100c957806333d6247d146100d2575b600080fd5b61009660015481565b6040519081526020015b60405180910390f35b6100966100b7366004610455565b60036020526000908152604090205481565b61009660025481565b6100e56100e0366004610455565b610171565b005b61009660025460015460408051602081019390935282015260009060600160405160208183030381529060405280519060200120905090565b6100e561012e36600461048a565b6102f3565b600554610146906001600160a01b031681565b6040516001600160a01b0390911681526020016100a0565b600454610146906001600160a01b031681565b6005546001600160a01b031633148061019457506004546001600160a01b031633145b6102315760405162461bcd60e51b815260206004820152604260248201527f506f6c79676f6e5a6b45564d476c6f62616c45786974526f6f743a3a7570646160448201527f746545786974526f6f743a204f6e6c7920616c6c6f77656420636f6e7472616360648201527f7473000000000000000000000000000000000000000000000000000000000000608482015260a4015b60405180910390fd5b6005546001600160a01b031633036102495760018190555b6004546001600160a01b031633036102615760028190555b60025460015460408051602081019390935282015260009060600160405160208183030381529060405280519060200120905060036000828152602001908152602001600020546000036102ef57600081815260036020526040808220429055600154600254915190927f61014378f82a0d809aefaf87a8ac9505b89c321808287a6e7810f29304c1fce391a35b5050565b600054610100900460ff16158080156103135750600054600160ff909116105b8061032d5750303b15801561032d575060005460ff166001145b61039f5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610228565b6000805460ff1916600117905580156103c2576000805461ff0019166101001790555b600580546001600160a01b038086167fffffffffffffffffffffffff00000000000000000000000000000000000000009283161790925560048054928516929091169190911790558015610450576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050565b60006020828403121561046757600080fd5b5035919050565b80356001600160a01b038116811461048557600080fd5b919050565b6000806040838503121561049d57600080fd5b6104a68361046e565b91506104b46020840161046e565b9050925092905056fea26469706673582212204755a4cbf50d99cfafb99c5429e784137da120dba2e519fdf39654b3c4a8a5d864736f6c634300080f0033",
}

// PolygonzkevmglobalexitrootABI is the input ABI used to generate the binding from.
// Deprecated: Use PolygonzkevmglobalexitrootMetaData.ABI instead.
var PolygonzkevmglobalexitrootABI = PolygonzkevmglobalexitrootMetaData.ABI

// PolygonzkevmglobalexitrootBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PolygonzkevmglobalexitrootMetaData.Bin instead.
var PolygonzkevmglobalexitrootBin = PolygonzkevmglobalexitrootMetaData.Bin

// DeployPolygonzkevmglobalexitroot deploys a new Ethereum contract, binding an instance of Polygonzkevmglobalexitroot to it.
func DeployPolygonzkevmglobalexitroot(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Polygonzkevmglobalexitroot, error) {
	parsed, err := PolygonzkevmglobalexitrootMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PolygonzkevmglobalexitrootBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Polygonzkevmglobalexitroot{PolygonzkevmglobalexitrootCaller: PolygonzkevmglobalexitrootCaller{contract: contract}, PolygonzkevmglobalexitrootTransactor: PolygonzkevmglobalexitrootTransactor{contract: contract}, PolygonzkevmglobalexitrootFilterer: PolygonzkevmglobalexitrootFilterer{contract: contract}}, nil
}

// Polygonzkevmglobalexitroot is an auto generated Go binding around an Ethereum contract.
type Polygonzkevmglobalexitroot struct {
	PolygonzkevmglobalexitrootCaller     // Read-only binding to the contract
	PolygonzkevmglobalexitrootTransactor // Write-only binding to the contract
	PolygonzkevmglobalexitrootFilterer   // Log filterer for contract events
}

// PolygonzkevmglobalexitrootCaller is an auto generated read-only Go binding around an Ethereum contract.
type PolygonzkevmglobalexitrootCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolygonzkevmglobalexitrootTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PolygonzkevmglobalexitrootTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolygonzkevmglobalexitrootFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PolygonzkevmglobalexitrootFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolygonzkevmglobalexitrootSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PolygonzkevmglobalexitrootSession struct {
	Contract     *Polygonzkevmglobalexitroot // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// PolygonzkevmglobalexitrootCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PolygonzkevmglobalexitrootCallerSession struct {
	Contract *PolygonzkevmglobalexitrootCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// PolygonzkevmglobalexitrootTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PolygonzkevmglobalexitrootTransactorSession struct {
	Contract     *PolygonzkevmglobalexitrootTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// PolygonzkevmglobalexitrootRaw is an auto generated low-level Go binding around an Ethereum contract.
type PolygonzkevmglobalexitrootRaw struct {
	Contract *Polygonzkevmglobalexitroot // Generic contract binding to access the raw methods on
}

// PolygonzkevmglobalexitrootCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PolygonzkevmglobalexitrootCallerRaw struct {
	Contract *PolygonzkevmglobalexitrootCaller // Generic read-only contract binding to access the raw methods on
}

// PolygonzkevmglobalexitrootTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PolygonzkevmglobalexitrootTransactorRaw struct {
	Contract *PolygonzkevmglobalexitrootTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPolygonzkevmglobalexitroot creates a new instance of Polygonzkevmglobalexitroot, bound to a specific deployed contract.
func NewPolygonzkevmglobalexitroot(address common.Address, backend bind.ContractBackend) (*Polygonzkevmglobalexitroot, error) {
	contract, err := bindPolygonzkevmglobalexitroot(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Polygonzkevmglobalexitroot{PolygonzkevmglobalexitrootCaller: PolygonzkevmglobalexitrootCaller{contract: contract}, PolygonzkevmglobalexitrootTransactor: PolygonzkevmglobalexitrootTransactor{contract: contract}, PolygonzkevmglobalexitrootFilterer: PolygonzkevmglobalexitrootFilterer{contract: contract}}, nil
}

// NewPolygonzkevmglobalexitrootCaller creates a new read-only instance of Polygonzkevmglobalexitroot, bound to a specific deployed contract.
func NewPolygonzkevmglobalexitrootCaller(address common.Address, caller bind.ContractCaller) (*PolygonzkevmglobalexitrootCaller, error) {
	contract, err := bindPolygonzkevmglobalexitroot(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmglobalexitrootCaller{contract: contract}, nil
}

// NewPolygonzkevmglobalexitrootTransactor creates a new write-only instance of Polygonzkevmglobalexitroot, bound to a specific deployed contract.
func NewPolygonzkevmglobalexitrootTransactor(address common.Address, transactor bind.ContractTransactor) (*PolygonzkevmglobalexitrootTransactor, error) {
	contract, err := bindPolygonzkevmglobalexitroot(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmglobalexitrootTransactor{contract: contract}, nil
}

// NewPolygonzkevmglobalexitrootFilterer creates a new log filterer instance of Polygonzkevmglobalexitroot, bound to a specific deployed contract.
func NewPolygonzkevmglobalexitrootFilterer(address common.Address, filterer bind.ContractFilterer) (*PolygonzkevmglobalexitrootFilterer, error) {
	contract, err := bindPolygonzkevmglobalexitroot(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmglobalexitrootFilterer{contract: contract}, nil
}

// bindPolygonzkevmglobalexitroot binds a generic wrapper to an already deployed contract.
func bindPolygonzkevmglobalexitroot(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PolygonzkevmglobalexitrootABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Polygonzkevmglobalexitroot.Contract.PolygonzkevmglobalexitrootCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitroot.Contract.PolygonzkevmglobalexitrootTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitroot.Contract.PolygonzkevmglobalexitrootTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Polygonzkevmglobalexitroot.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitroot.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitroot.Contract.contract.Transact(opts, method, params...)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootCaller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonzkevmglobalexitroot.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootSession) BridgeAddress() (common.Address, error) {
	return _Polygonzkevmglobalexitroot.Contract.BridgeAddress(&_Polygonzkevmglobalexitroot.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootCallerSession) BridgeAddress() (common.Address, error) {
	return _Polygonzkevmglobalexitroot.Contract.BridgeAddress(&_Polygonzkevmglobalexitroot.CallOpts)
}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootCaller) GetLastGlobalExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Polygonzkevmglobalexitroot.contract.Call(opts, &out, "getLastGlobalExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootSession) GetLastGlobalExitRoot() ([32]byte, error) {
	return _Polygonzkevmglobalexitroot.Contract.GetLastGlobalExitRoot(&_Polygonzkevmglobalexitroot.CallOpts)
}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootCallerSession) GetLastGlobalExitRoot() ([32]byte, error) {
	return _Polygonzkevmglobalexitroot.Contract.GetLastGlobalExitRoot(&_Polygonzkevmglobalexitroot.CallOpts)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootCaller) GlobalExitRootMap(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Polygonzkevmglobalexitroot.contract.Call(opts, &out, "globalExitRootMap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootSession) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Polygonzkevmglobalexitroot.Contract.GlobalExitRootMap(&_Polygonzkevmglobalexitroot.CallOpts, arg0)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootCallerSession) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Polygonzkevmglobalexitroot.Contract.GlobalExitRootMap(&_Polygonzkevmglobalexitroot.CallOpts, arg0)
}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootCaller) LastMainnetExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Polygonzkevmglobalexitroot.contract.Call(opts, &out, "lastMainnetExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootSession) LastMainnetExitRoot() ([32]byte, error) {
	return _Polygonzkevmglobalexitroot.Contract.LastMainnetExitRoot(&_Polygonzkevmglobalexitroot.CallOpts)
}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootCallerSession) LastMainnetExitRoot() ([32]byte, error) {
	return _Polygonzkevmglobalexitroot.Contract.LastMainnetExitRoot(&_Polygonzkevmglobalexitroot.CallOpts)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootCaller) LastRollupExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Polygonzkevmglobalexitroot.contract.Call(opts, &out, "lastRollupExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootSession) LastRollupExitRoot() ([32]byte, error) {
	return _Polygonzkevmglobalexitroot.Contract.LastRollupExitRoot(&_Polygonzkevmglobalexitroot.CallOpts)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootCallerSession) LastRollupExitRoot() ([32]byte, error) {
	return _Polygonzkevmglobalexitroot.Contract.LastRollupExitRoot(&_Polygonzkevmglobalexitroot.CallOpts)
}

// RollupAddress is a free data retrieval call binding the contract method 0x5ec6a8df.
//
// Solidity: function rollupAddress() view returns(address)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootCaller) RollupAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonzkevmglobalexitroot.contract.Call(opts, &out, "rollupAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RollupAddress is a free data retrieval call binding the contract method 0x5ec6a8df.
//
// Solidity: function rollupAddress() view returns(address)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootSession) RollupAddress() (common.Address, error) {
	return _Polygonzkevmglobalexitroot.Contract.RollupAddress(&_Polygonzkevmglobalexitroot.CallOpts)
}

// RollupAddress is a free data retrieval call binding the contract method 0x5ec6a8df.
//
// Solidity: function rollupAddress() view returns(address)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootCallerSession) RollupAddress() (common.Address, error) {
	return _Polygonzkevmglobalexitroot.Contract.RollupAddress(&_Polygonzkevmglobalexitroot.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _rollupAddress, address _bridgeAddress) returns()
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootTransactor) Initialize(opts *bind.TransactOpts, _rollupAddress common.Address, _bridgeAddress common.Address) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitroot.contract.Transact(opts, "initialize", _rollupAddress, _bridgeAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _rollupAddress, address _bridgeAddress) returns()
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootSession) Initialize(_rollupAddress common.Address, _bridgeAddress common.Address) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitroot.Contract.Initialize(&_Polygonzkevmglobalexitroot.TransactOpts, _rollupAddress, _bridgeAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _rollupAddress, address _bridgeAddress) returns()
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootTransactorSession) Initialize(_rollupAddress common.Address, _bridgeAddress common.Address) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitroot.Contract.Initialize(&_Polygonzkevmglobalexitroot.TransactOpts, _rollupAddress, _bridgeAddress)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootTransactor) UpdateExitRoot(opts *bind.TransactOpts, newRoot [32]byte) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitroot.contract.Transact(opts, "updateExitRoot", newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootSession) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitroot.Contract.UpdateExitRoot(&_Polygonzkevmglobalexitroot.TransactOpts, newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootTransactorSession) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Polygonzkevmglobalexitroot.Contract.UpdateExitRoot(&_Polygonzkevmglobalexitroot.TransactOpts, newRoot)
}

// PolygonzkevmglobalexitrootInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Polygonzkevmglobalexitroot contract.
type PolygonzkevmglobalexitrootInitializedIterator struct {
	Event *PolygonzkevmglobalexitrootInitialized // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmglobalexitrootInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmglobalexitrootInitialized)
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
		it.Event = new(PolygonzkevmglobalexitrootInitialized)
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
func (it *PolygonzkevmglobalexitrootInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmglobalexitrootInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmglobalexitrootInitialized represents a Initialized event raised by the Polygonzkevmglobalexitroot contract.
type PolygonzkevmglobalexitrootInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootFilterer) FilterInitialized(opts *bind.FilterOpts) (*PolygonzkevmglobalexitrootInitializedIterator, error) {

	logs, sub, err := _Polygonzkevmglobalexitroot.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmglobalexitrootInitializedIterator{contract: _Polygonzkevmglobalexitroot.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PolygonzkevmglobalexitrootInitialized) (event.Subscription, error) {

	logs, sub, err := _Polygonzkevmglobalexitroot.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmglobalexitrootInitialized)
				if err := _Polygonzkevmglobalexitroot.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootFilterer) ParseInitialized(log types.Log) (*PolygonzkevmglobalexitrootInitialized, error) {
	event := new(PolygonzkevmglobalexitrootInitialized)
	if err := _Polygonzkevmglobalexitroot.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonzkevmglobalexitrootUpdateGlobalExitRootIterator is returned from FilterUpdateGlobalExitRoot and is used to iterate over the raw logs and unpacked data for UpdateGlobalExitRoot events raised by the Polygonzkevmglobalexitroot contract.
type PolygonzkevmglobalexitrootUpdateGlobalExitRootIterator struct {
	Event *PolygonzkevmglobalexitrootUpdateGlobalExitRoot // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmglobalexitrootUpdateGlobalExitRootIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmglobalexitrootUpdateGlobalExitRoot)
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
		it.Event = new(PolygonzkevmglobalexitrootUpdateGlobalExitRoot)
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
func (it *PolygonzkevmglobalexitrootUpdateGlobalExitRootIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmglobalexitrootUpdateGlobalExitRootIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmglobalexitrootUpdateGlobalExitRoot represents a UpdateGlobalExitRoot event raised by the Polygonzkevmglobalexitroot contract.
type PolygonzkevmglobalexitrootUpdateGlobalExitRoot struct {
	MainnetExitRoot [32]byte
	RollupExitRoot  [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpdateGlobalExitRoot is a free log retrieval operation binding the contract event 0x61014378f82a0d809aefaf87a8ac9505b89c321808287a6e7810f29304c1fce3.
//
// Solidity: event UpdateGlobalExitRoot(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootFilterer) FilterUpdateGlobalExitRoot(opts *bind.FilterOpts, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (*PolygonzkevmglobalexitrootUpdateGlobalExitRootIterator, error) {

	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Polygonzkevmglobalexitroot.contract.FilterLogs(opts, "UpdateGlobalExitRoot", mainnetExitRootRule, rollupExitRootRule)
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmglobalexitrootUpdateGlobalExitRootIterator{contract: _Polygonzkevmglobalexitroot.contract, event: "UpdateGlobalExitRoot", logs: logs, sub: sub}, nil
}

// WatchUpdateGlobalExitRoot is a free log subscription operation binding the contract event 0x61014378f82a0d809aefaf87a8ac9505b89c321808287a6e7810f29304c1fce3.
//
// Solidity: event UpdateGlobalExitRoot(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootFilterer) WatchUpdateGlobalExitRoot(opts *bind.WatchOpts, sink chan<- *PolygonzkevmglobalexitrootUpdateGlobalExitRoot, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (event.Subscription, error) {

	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Polygonzkevmglobalexitroot.contract.WatchLogs(opts, "UpdateGlobalExitRoot", mainnetExitRootRule, rollupExitRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmglobalexitrootUpdateGlobalExitRoot)
				if err := _Polygonzkevmglobalexitroot.contract.UnpackLog(event, "UpdateGlobalExitRoot", log); err != nil {
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
func (_Polygonzkevmglobalexitroot *PolygonzkevmglobalexitrootFilterer) ParseUpdateGlobalExitRoot(log types.Log) (*PolygonzkevmglobalexitrootUpdateGlobalExitRoot, error) {
	event := new(PolygonzkevmglobalexitrootUpdateGlobalExitRoot)
	if err := _Polygonzkevmglobalexitroot.contract.UnpackLog(event, "UpdateGlobalExitRoot", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
