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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rollupAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"globalExitRootNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"mainnetExitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"rollupExitRoot\",\"type\":\"bytes32\"}],\"name\":\"UpdateGlobalExitRoot\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastGlobalExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"globalExitRootMap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastGlobalExitRootNum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMainnetExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newRoot\",\"type\":\"bytes32\"}],\"name\":\"updateExitRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162000dd638038062000dd6833981810160405281019062000037919062000217565b620000576200004b620000e160201b60201c565b620000e960201b60201c565b81600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050506200025e565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620001df82620001b2565b9050919050565b620001f181620001d2565b8114620001fd57600080fd5b50565b6000815190506200021181620001e6565b92915050565b60008060408385031215620002315762000230620001ad565b5b6000620002418582860162000200565b9250506020620002548582860162000200565b9150509250929050565b610b68806200026e6000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80633ed691ef116100715780633ed691ef146101545780635ec6a8df14610172578063715018a6146101905780638da5cb5b1461019a578063a3c573eb146101b8578063f2fde38b146101d6576100a9565b806301fd9044146100ae578063029f2793146100cc578063257b3632146100ea578063319cf7351461011a57806333d6247d14610138575b600080fd5b6100b66101f2565b6040516100c3919061076a565b60405180910390f35b6100d46101f8565b6040516100e1919061079e565b60405180910390f35b61010460048036038101906100ff91906107ea565b6101fe565b604051610111919061079e565b60405180910390f35b610122610216565b60405161012f919061076a565b60405180910390f35b610152600480360381019061014d91906107ea565b61021c565b005b61015c61045c565b604051610169919061076a565b60405180910390f35b61017a610490565b6040516101879190610858565b60405180910390f35b6101986104b6565b005b6101a261053e565b6040516101af9190610858565b60405180910390f35b6101c0610567565b6040516101cd9190610858565b60405180910390f35b6101f060048036038101906101eb919061089f565b61058d565b005b60015481565b60045481565b60036020528060005260406000206000915090505481565b60025481565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806102c55750600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b610304576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102fb9061094f565b60405180910390fd5b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141561036257806001819055505b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614156103c057806002819055505b600460008154809291906103d39061099e565b919050555060006002546001546040516020016103f1929190610a08565b60405160208183030381529060405280519060200120905060045460036000838152602001908152602001600020819055506001546002546004547fb7c409af8cb511116b88f38824d48a0196194596241fdb2d177210d3d3b89fbf60405160405180910390a45050565b6000600254600154604051602001610475929190610a08565b60405160208183030381529060405280519060200120905090565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6104be610685565b73ffffffffffffffffffffffffffffffffffffffff166104dc61053e565b73ffffffffffffffffffffffffffffffffffffffff1614610532576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161052990610a80565b60405180910390fd5b61053c600061068d565b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b610595610685565b73ffffffffffffffffffffffffffffffffffffffff166105b361053e565b73ffffffffffffffffffffffffffffffffffffffff1614610609576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161060090610a80565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610679576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161067090610b12565b60405180910390fd5b6106828161068d565b50565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000819050919050565b61076481610751565b82525050565b600060208201905061077f600083018461075b565b92915050565b6000819050919050565b61079881610785565b82525050565b60006020820190506107b3600083018461078f565b92915050565b600080fd5b6107c781610751565b81146107d257600080fd5b50565b6000813590506107e4816107be565b92915050565b600060208284031215610800576107ff6107b9565b5b600061080e848285016107d5565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061084282610817565b9050919050565b61085281610837565b82525050565b600060208201905061086d6000830184610849565b92915050565b61087c81610837565b811461088757600080fd5b50565b60008135905061089981610873565b92915050565b6000602082840312156108b5576108b46107b9565b5b60006108c38482850161088a565b91505092915050565b600082825260208201905092915050565b7f476c6f62616c45786974526f6f744d616e616765723a3a75706461746545786960008201527f74526f6f743a204f4e4c595f414c4c4f5745445f434f4e545241435453000000602082015250565b6000610939603d836108cc565b9150610944826108dd565b604082019050919050565b600060208201905081810360008301526109688161092c565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006109a982610785565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156109dc576109db61096f565b5b600182019050919050565b6000819050919050565b610a026109fd82610751565b6109e7565b82525050565b6000610a1482856109f1565b602082019150610a2482846109f1565b6020820191508190509392505050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b6000610a6a6020836108cc565b9150610a7582610a34565b602082019050919050565b60006020820190508181036000830152610a9981610a5d565b9050919050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b6000610afc6026836108cc565b9150610b0782610aa0565b604082019050919050565b60006020820190508181036000830152610b2b81610aef565b905091905056fea2646970667358221220733bb63d3645393b2926015ec5300dfaffd454e31a2f4292d75f52060270af5f64736f6c63430008090033",
}

// GlobalexitrootmanagerABI is the input ABI used to generate the binding from.
// Deprecated: Use GlobalexitrootmanagerMetaData.ABI instead.
var GlobalexitrootmanagerABI = GlobalexitrootmanagerMetaData.ABI

// GlobalexitrootmanagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GlobalexitrootmanagerMetaData.Bin instead.
var GlobalexitrootmanagerBin = GlobalexitrootmanagerMetaData.Bin

// DeployGlobalexitrootmanager deploys a new Ethereum contract, binding an instance of Globalexitrootmanager to it.
func DeployGlobalexitrootmanager(auth *bind.TransactOpts, backend bind.ContractBackend, _rollupAddress common.Address, _bridgeAddress common.Address) (common.Address, *types.Transaction, *Globalexitrootmanager, error) {
	parsed, err := GlobalexitrootmanagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GlobalexitrootmanagerBin), backend, _rollupAddress, _bridgeAddress)
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

// LastGlobalExitRootNum is a free data retrieval call binding the contract method 0x029f2793.
//
// Solidity: function lastGlobalExitRootNum() view returns(uint256)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) LastGlobalExitRootNum(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "lastGlobalExitRootNum")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastGlobalExitRootNum is a free data retrieval call binding the contract method 0x029f2793.
//
// Solidity: function lastGlobalExitRootNum() view returns(uint256)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) LastGlobalExitRootNum() (*big.Int, error) {
	return _Globalexitrootmanager.Contract.LastGlobalExitRootNum(&_Globalexitrootmanager.CallOpts)
}

// LastGlobalExitRootNum is a free data retrieval call binding the contract method 0x029f2793.
//
// Solidity: function lastGlobalExitRootNum() view returns(uint256)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) LastGlobalExitRootNum() (*big.Int, error) {
	return _Globalexitrootmanager.Contract.LastGlobalExitRootNum(&_Globalexitrootmanager.CallOpts)
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

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Globalexitrootmanager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) Owner() (common.Address, error) {
	return _Globalexitrootmanager.Contract.Owner(&_Globalexitrootmanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Globalexitrootmanager *GlobalexitrootmanagerCallerSession) Owner() (common.Address, error) {
	return _Globalexitrootmanager.Contract.Owner(&_Globalexitrootmanager.CallOpts)
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

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Globalexitrootmanager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.RenounceOwnership(&_Globalexitrootmanager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.RenounceOwnership(&_Globalexitrootmanager.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Globalexitrootmanager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.TransferOwnership(&_Globalexitrootmanager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Globalexitrootmanager *GlobalexitrootmanagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Globalexitrootmanager.Contract.TransferOwnership(&_Globalexitrootmanager.TransactOpts, newOwner)
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

// GlobalexitrootmanagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Globalexitrootmanager contract.
type GlobalexitrootmanagerOwnershipTransferredIterator struct {
	Event *GlobalexitrootmanagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *GlobalexitrootmanagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GlobalexitrootmanagerOwnershipTransferred)
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
		it.Event = new(GlobalexitrootmanagerOwnershipTransferred)
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
func (it *GlobalexitrootmanagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GlobalexitrootmanagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GlobalexitrootmanagerOwnershipTransferred represents a OwnershipTransferred event raised by the Globalexitrootmanager contract.
type GlobalexitrootmanagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GlobalexitrootmanagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Globalexitrootmanager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GlobalexitrootmanagerOwnershipTransferredIterator{contract: _Globalexitrootmanager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GlobalexitrootmanagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Globalexitrootmanager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GlobalexitrootmanagerOwnershipTransferred)
				if err := _Globalexitrootmanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) ParseOwnershipTransferred(log types.Log) (*GlobalexitrootmanagerOwnershipTransferred, error) {
	event := new(GlobalexitrootmanagerOwnershipTransferred)
	if err := _Globalexitrootmanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
	GlobalExitRootNum *big.Int
	MainnetExitRoot   [32]byte
	RollupExitRoot    [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUpdateGlobalExitRoot is a free log retrieval operation binding the contract event 0xb7c409af8cb511116b88f38824d48a0196194596241fdb2d177210d3d3b89fbf.
//
// Solidity: event UpdateGlobalExitRoot(uint256 indexed globalExitRootNum, bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) FilterUpdateGlobalExitRoot(opts *bind.FilterOpts, globalExitRootNum []*big.Int, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (*GlobalexitrootmanagerUpdateGlobalExitRootIterator, error) {

	var globalExitRootNumRule []interface{}
	for _, globalExitRootNumItem := range globalExitRootNum {
		globalExitRootNumRule = append(globalExitRootNumRule, globalExitRootNumItem)
	}
	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Globalexitrootmanager.contract.FilterLogs(opts, "UpdateGlobalExitRoot", globalExitRootNumRule, mainnetExitRootRule, rollupExitRootRule)
	if err != nil {
		return nil, err
	}
	return &GlobalexitrootmanagerUpdateGlobalExitRootIterator{contract: _Globalexitrootmanager.contract, event: "UpdateGlobalExitRoot", logs: logs, sub: sub}, nil
}

// WatchUpdateGlobalExitRoot is a free log subscription operation binding the contract event 0xb7c409af8cb511116b88f38824d48a0196194596241fdb2d177210d3d3b89fbf.
//
// Solidity: event UpdateGlobalExitRoot(uint256 indexed globalExitRootNum, bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) WatchUpdateGlobalExitRoot(opts *bind.WatchOpts, sink chan<- *GlobalexitrootmanagerUpdateGlobalExitRoot, globalExitRootNum []*big.Int, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (event.Subscription, error) {

	var globalExitRootNumRule []interface{}
	for _, globalExitRootNumItem := range globalExitRootNum {
		globalExitRootNumRule = append(globalExitRootNumRule, globalExitRootNumItem)
	}
	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Globalexitrootmanager.contract.WatchLogs(opts, "UpdateGlobalExitRoot", globalExitRootNumRule, mainnetExitRootRule, rollupExitRootRule)
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

// ParseUpdateGlobalExitRoot is a log parse operation binding the contract event 0xb7c409af8cb511116b88f38824d48a0196194596241fdb2d177210d3d3b89fbf.
//
// Solidity: event UpdateGlobalExitRoot(uint256 indexed globalExitRootNum, bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Globalexitrootmanager *GlobalexitrootmanagerFilterer) ParseUpdateGlobalExitRoot(log types.Log) (*GlobalexitrootmanagerUpdateGlobalExitRoot, error) {
	event := new(GlobalexitrootmanagerUpdateGlobalExitRoot)
	if err := _Globalexitrootmanager.contract.UnpackLog(event, "UpdateGlobalExitRoot", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
