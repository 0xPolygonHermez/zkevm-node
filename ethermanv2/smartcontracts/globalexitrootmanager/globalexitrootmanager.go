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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rollupAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"globalExitRootNum\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"mainnetExitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"rollupExitRoot\",\"type\":\"bytes32\"}],\"name\":\"UpdateGlobalExitRoot\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastGlobalExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"globalExitRootMap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastGlobalExitRootNum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMainnetExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newRoot\",\"type\":\"bytes32\"}],\"name\":\"updateExitRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516108a13803806108a18339818101604052810190610032919061011e565b81600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600460006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505061015e565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100eb826100c0565b9050919050565b6100fb816100e0565b811461010657600080fd5b50565b600081519050610118816100f2565b92915050565b60008060408385031215610135576101346100bb565b5b600061014385828601610109565b925050602061015485828601610109565b9150509250929050565b6107348061016d6000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c806333d6247d1161005b57806333d6247d146101175780633ed691ef146101335780635ec6a8df14610151578063a3c573eb1461016f57610088565b806301fd90441461008d578063029f2793146100ab578063257b3632146100c9578063319cf735146100f9575b600080fd5b61009561018d565b6040516100a2919061048e565b60405180910390f35b6100b3610193565b6040516100c091906104c2565b60405180910390f35b6100e360048036038101906100de919061050e565b610199565b6040516100f091906104c2565b60405180910390f35b6101016101b1565b60405161010e919061048e565b60405180910390f35b610131600480360381019061012c919061050e565b6101b7565b005b61013b6103f5565b604051610148919061048e565b60405180910390f35b610159610429565b604051610166919061057c565b60405180910390f35b61017761044f565b604051610184919061057c565b60405180910390f35b60005481565b60035481565b60026020528060005260406000206000915090505481565b60015481565b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806102605750600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b61029f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102969061061a565b60405180910390fd5b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16036102fc57806000819055505b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff160361035957806001819055505b6003600081548092919061036c90610669565b9190505550600060015460005460405160200161038a9291906106d2565b60405160208183030381529060405280519060200120905060035460026000838152602001908152602001600020819055506000546001546003547fb7c409af8cb511116b88f38824d48a0196194596241fdb2d177210d3d3b89fbf60405160405180910390a45050565b600060015460005460405160200161040e9291906106d2565b60405160208183030381529060405280519060200120905090565b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000819050919050565b61048881610475565b82525050565b60006020820190506104a3600083018461047f565b92915050565b6000819050919050565b6104bc816104a9565b82525050565b60006020820190506104d760008301846104b3565b92915050565b600080fd5b6104eb81610475565b81146104f657600080fd5b50565b600081359050610508816104e2565b92915050565b600060208284031215610524576105236104dd565b5b6000610532848285016104f9565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006105668261053b565b9050919050565b6105768161055b565b82525050565b6000602082019050610591600083018461056d565b92915050565b600082825260208201905092915050565b7f476c6f62616c45786974526f6f744d616e616765723a3a75706461746545786960008201527f74526f6f743a204f4e4c595f414c4c4f5745445f434f4e545241435453000000602082015250565b6000610604603d83610597565b915061060f826105a8565b604082019050919050565b60006020820190508181036000830152610633816105f7565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610674826104a9565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036106a6576106a561063a565b5b600182019050919050565b6000819050919050565b6106cc6106c782610475565b6106b1565b82525050565b60006106de82856106bb565b6020820191506106ee82846106bb565b602082019150819050939250505056fea2646970667358221220aadf8631dcbfd25b8114e2fd4206c0751f52b9b549730dc166ec14edb2e90aec64736f6c634300080f0033",
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
