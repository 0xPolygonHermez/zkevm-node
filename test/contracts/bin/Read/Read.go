// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Read

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

// Readtoken is an auto generated low-level Go binding around an user-defined struct.
type Readtoken struct {
	Name     string
	Quantity *big.Int
	Address  common.Address
}

// ReadMetaData contains all meta data concerning the Read contract.
var ReadMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"Owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OwnerName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"Tokens\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"Quantity\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"Value\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"Quantity\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"internalType\":\"structRead.token\",\"name\":\"t\",\"type\":\"tuple\"}],\"name\":\"externalAddToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"externalGetOwnerName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"a\",\"type\":\"address\"}],\"name\":\"externalGetToken\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"Quantity\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"internalType\":\"structRead.token\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"externalRead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"externalReadWParams\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"Quantity\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"internalType\":\"structRead.token\",\"name\":\"t\",\"type\":\"tuple\"}],\"name\":\"publicAddToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"publicGetOwnerName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"a\",\"type\":\"address\"}],\"name\":\"publicGetToken\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"Name\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"Quantity\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"Address\",\"type\":\"address\"}],\"internalType\":\"structRead.token\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"publicRead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"p\",\"type\":\"uint256\"}],\"name\":\"publicReadWParams\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405260016002553480156200001657600080fd5b5060405162000b1238038062000b12833981016040819052620000399162000124565b80516200004e90600190602084019062000068565b5050600080546001600160a01b031916331790556200023d565b828054620000769062000200565b90600052602060002090601f0160209004810192826200009a5760008555620000e5565b82601f10620000b557805160ff1916838001178555620000e5565b82800160010185558215620000e5579182015b82811115620000e5578251825591602001919060010190620000c8565b50620000f3929150620000f7565b5090565b5b80821115620000f35760008155600101620000f8565b634e487b7160e01b600052604160045260246000fd5b600060208083850312156200013857600080fd5b82516001600160401b03808211156200015057600080fd5b818501915085601f8301126200016557600080fd5b8151818111156200017a576200017a6200010e565b604051601f8201601f19908116603f01168101908382118183101715620001a557620001a56200010e565b816040528281528886848701011115620001be57600080fd5b600093505b82841015620001e25784840186015181850187015292850192620001c3565b82841115620001f45760008684830101525b98975050505050505050565b600181811c908216806200021557607f821691505b602082108114156200023757634e487b7160e01b600052602260045260246000fd5b50919050565b6108c5806200024d6000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c8063b1ceff4e1161008c578063c191ace511610066578063c191ace5146100ef578063d05c715a146101b4578063e117176a14610169578063f1d876b4146101bc57600080fd5b8063b1ceff4e14610169578063b4a99a4e14610189578063bfa044ed1461015657600080fd5b806333eb6be8116100c857806333eb6be8146101415780636e5298b2146101565780639553400f14610141578063af3347571461012f57600080fd5b8063081540b2146100ef5780630deef63a1461010d57806310f3f91a1461012f575b600080fd5b6100f76101c5565b6040516101049190610602565b60405180910390f35b61012061011b366004610638565b610257565b60405161010493929190610653565b6002545b604051908152602001610104565b61015461014f3660046106f6565b61030b565b005b6101336101643660046107ce565b610373565b61017c610177366004610638565b610389565b60405161010491906107e7565b60005461019c906001600160a01b031681565b6040516001600160a01b039091168152602001610104565b6100f761048e565b61013360025481565b6060600180546101d49061082e565b80601f01602080910402602001604051908101604052809291908181526020018280546102009061082e565b801561024d5780601f106102225761010080835404028352916020019161024d565b820191906000526020600020905b81548152906001019060200180831161023057829003601f168201915b5050505050905090565b6003602052600090815260409020805481906102729061082e565b80601f016020809104026020016040519081016040528092919081815260200182805461029e9061082e565b80156102eb5780601f106102c0576101008083540402835291602001916102eb565b820191906000526020600020905b8154815290600101906020018083116102ce57829003601f168201915b5050505060018301546002909301549192916001600160a01b0316905083565b6040808201516001600160a01b031660009081526003602090815291902082518051849361033d92849291019061051c565b5060208201516001820155604090910151600290910180546001600160a01b0319166001600160a01b0390921691909117905550565b6000816002546103839190610869565b92915050565b6103b66040518060600160405280606081526020016000815260200160006001600160a01b031681525090565b6001600160a01b038216600090815260036020526040908190208151606081019092528054829082906103e89061082e565b80601f01602080910402602001604051908101604052809291908181526020018280546104149061082e565b80156104615780601f1061043657610100808354040283529160200191610461565b820191906000526020600020905b81548152906001019060200180831161044457829003601f168201915b5050509183525050600182015460208201526002909101546001600160a01b031660409091015292915050565b6001805461049b9061082e565b80601f01602080910402602001604051908101604052809291908181526020018280546104c79061082e565b80156105145780601f106104e957610100808354040283529160200191610514565b820191906000526020600020905b8154815290600101906020018083116104f757829003601f168201915b505050505081565b8280546105289061082e565b90600052602060002090601f01602090048101928261054a5760008555610590565b82601f1061056357805160ff1916838001178555610590565b82800160010185558215610590579182015b82811115610590578251825591602001919060010190610575565b5061059c9291506105a0565b5090565b5b8082111561059c57600081556001016105a1565b6000815180845260005b818110156105db576020818501810151868301820152016105bf565b818111156105ed576000602083870101525b50601f01601f19169290920160200192915050565b60208152600061061560208301846105b5565b9392505050565b80356001600160a01b038116811461063357600080fd5b919050565b60006020828403121561064a57600080fd5b6106158261061c565b60608152600061066660608301866105b5565b6020830194909452506001600160a01b0391909116604090910152919050565b634e487b7160e01b600052604160045260246000fd5b6040516060810167ffffffffffffffff811182821017156106bf576106bf610686565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156106ee576106ee610686565b604052919050565b6000602080838503121561070957600080fd5b823567ffffffffffffffff8082111561072157600080fd5b908401906060828703121561073557600080fd5b61073d61069c565b82358281111561074c57600080fd5b8301601f8101881361075d57600080fd5b80358381111561076f5761076f610686565b610781601f8201601f191687016106c5565b9350808452888682840101111561079757600080fd5b808683018786013760009084018601525081815282840135818501526107bf6040840161061c565b60408201529695505050505050565b6000602082840312156107e057600080fd5b5035919050565b60208152600082516060602084015261080360808401826105b5565b6020850151604085810191909152909401516001600160a01b03166060909301929092525090919050565b600181811c9082168061084257607f821691505b6020821081141561086357634e487b7160e01b600052602260045260246000fd5b50919050565b6000821982111561088a57634e487b7160e01b600052601160045260246000fd5b50019056fea2646970667358221220a0f0d200f69f0e789d27f571230dfef313a2709f90ab94a804881830d196e05164736f6c634300080c0033",
}

// ReadABI is the input ABI used to generate the binding from.
// Deprecated: Use ReadMetaData.ABI instead.
var ReadABI = ReadMetaData.ABI

// ReadBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ReadMetaData.Bin instead.
var ReadBin = ReadMetaData.Bin

// DeployRead deploys a new Ethereum contract, binding an instance of Read to it.
func DeployRead(auth *bind.TransactOpts, backend bind.ContractBackend, name string) (common.Address, *types.Transaction, *Read, error) {
	parsed, err := ReadMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ReadBin), backend, name)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Read{ReadCaller: ReadCaller{contract: contract}, ReadTransactor: ReadTransactor{contract: contract}, ReadFilterer: ReadFilterer{contract: contract}}, nil
}

// Read is an auto generated Go binding around an Ethereum contract.
type Read struct {
	ReadCaller     // Read-only binding to the contract
	ReadTransactor // Write-only binding to the contract
	ReadFilterer   // Log filterer for contract events
}

// ReadCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReadCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReadTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReadTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReadFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReadFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReadSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReadSession struct {
	Contract     *Read             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReadCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReadCallerSession struct {
	Contract *ReadCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ReadTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReadTransactorSession struct {
	Contract     *ReadTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReadRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReadRaw struct {
	Contract *Read // Generic contract binding to access the raw methods on
}

// ReadCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReadCallerRaw struct {
	Contract *ReadCaller // Generic read-only contract binding to access the raw methods on
}

// ReadTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReadTransactorRaw struct {
	Contract *ReadTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRead creates a new instance of Read, bound to a specific deployed contract.
func NewRead(address common.Address, backend bind.ContractBackend) (*Read, error) {
	contract, err := bindRead(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Read{ReadCaller: ReadCaller{contract: contract}, ReadTransactor: ReadTransactor{contract: contract}, ReadFilterer: ReadFilterer{contract: contract}}, nil
}

// NewReadCaller creates a new read-only instance of Read, bound to a specific deployed contract.
func NewReadCaller(address common.Address, caller bind.ContractCaller) (*ReadCaller, error) {
	contract, err := bindRead(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReadCaller{contract: contract}, nil
}

// NewReadTransactor creates a new write-only instance of Read, bound to a specific deployed contract.
func NewReadTransactor(address common.Address, transactor bind.ContractTransactor) (*ReadTransactor, error) {
	contract, err := bindRead(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReadTransactor{contract: contract}, nil
}

// NewReadFilterer creates a new log filterer instance of Read, bound to a specific deployed contract.
func NewReadFilterer(address common.Address, filterer bind.ContractFilterer) (*ReadFilterer, error) {
	contract, err := bindRead(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReadFilterer{contract: contract}, nil
}

// bindRead binds a generic wrapper to an already deployed contract.
func bindRead(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReadMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Read *ReadRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Read.Contract.ReadCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Read *ReadRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Read.Contract.ReadTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Read *ReadRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Read.Contract.ReadTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Read *ReadCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Read.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Read *ReadTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Read.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Read *ReadTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Read.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0xb4a99a4e.
//
// Solidity: function Owner() view returns(address)
func (_Read *ReadCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "Owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0xb4a99a4e.
//
// Solidity: function Owner() view returns(address)
func (_Read *ReadSession) Owner() (common.Address, error) {
	return _Read.Contract.Owner(&_Read.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0xb4a99a4e.
//
// Solidity: function Owner() view returns(address)
func (_Read *ReadCallerSession) Owner() (common.Address, error) {
	return _Read.Contract.Owner(&_Read.CallOpts)
}

// OwnerName is a free data retrieval call binding the contract method 0xd05c715a.
//
// Solidity: function OwnerName() view returns(string)
func (_Read *ReadCaller) OwnerName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "OwnerName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// OwnerName is a free data retrieval call binding the contract method 0xd05c715a.
//
// Solidity: function OwnerName() view returns(string)
func (_Read *ReadSession) OwnerName() (string, error) {
	return _Read.Contract.OwnerName(&_Read.CallOpts)
}

// OwnerName is a free data retrieval call binding the contract method 0xd05c715a.
//
// Solidity: function OwnerName() view returns(string)
func (_Read *ReadCallerSession) OwnerName() (string, error) {
	return _Read.Contract.OwnerName(&_Read.CallOpts)
}

// Tokens is a free data retrieval call binding the contract method 0x0deef63a.
//
// Solidity: function Tokens(address ) view returns(string Name, uint256 Quantity, address Address)
func (_Read *ReadCaller) Tokens(opts *bind.CallOpts, arg0 common.Address) (struct {
	Name     string
	Quantity *big.Int
	Address  common.Address
}, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "Tokens", arg0)

	outstruct := new(struct {
		Name     string
		Quantity *big.Int
		Address  common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Name = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Quantity = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Address = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Tokens is a free data retrieval call binding the contract method 0x0deef63a.
//
// Solidity: function Tokens(address ) view returns(string Name, uint256 Quantity, address Address)
func (_Read *ReadSession) Tokens(arg0 common.Address) (struct {
	Name     string
	Quantity *big.Int
	Address  common.Address
}, error) {
	return _Read.Contract.Tokens(&_Read.CallOpts, arg0)
}

// Tokens is a free data retrieval call binding the contract method 0x0deef63a.
//
// Solidity: function Tokens(address ) view returns(string Name, uint256 Quantity, address Address)
func (_Read *ReadCallerSession) Tokens(arg0 common.Address) (struct {
	Name     string
	Quantity *big.Int
	Address  common.Address
}, error) {
	return _Read.Contract.Tokens(&_Read.CallOpts, arg0)
}

// Value is a free data retrieval call binding the contract method 0xf1d876b4.
//
// Solidity: function Value() view returns(uint256)
func (_Read *ReadCaller) Value(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "Value")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Value is a free data retrieval call binding the contract method 0xf1d876b4.
//
// Solidity: function Value() view returns(uint256)
func (_Read *ReadSession) Value() (*big.Int, error) {
	return _Read.Contract.Value(&_Read.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0xf1d876b4.
//
// Solidity: function Value() view returns(uint256)
func (_Read *ReadCallerSession) Value() (*big.Int, error) {
	return _Read.Contract.Value(&_Read.CallOpts)
}

// ExternalGetOwnerName is a free data retrieval call binding the contract method 0xc191ace5.
//
// Solidity: function externalGetOwnerName() view returns(string)
func (_Read *ReadCaller) ExternalGetOwnerName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "externalGetOwnerName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ExternalGetOwnerName is a free data retrieval call binding the contract method 0xc191ace5.
//
// Solidity: function externalGetOwnerName() view returns(string)
func (_Read *ReadSession) ExternalGetOwnerName() (string, error) {
	return _Read.Contract.ExternalGetOwnerName(&_Read.CallOpts)
}

// ExternalGetOwnerName is a free data retrieval call binding the contract method 0xc191ace5.
//
// Solidity: function externalGetOwnerName() view returns(string)
func (_Read *ReadCallerSession) ExternalGetOwnerName() (string, error) {
	return _Read.Contract.ExternalGetOwnerName(&_Read.CallOpts)
}

// ExternalGetToken is a free data retrieval call binding the contract method 0xe117176a.
//
// Solidity: function externalGetToken(address a) view returns((string,uint256,address))
func (_Read *ReadCaller) ExternalGetToken(opts *bind.CallOpts, a common.Address) (Readtoken, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "externalGetToken", a)

	if err != nil {
		return *new(Readtoken), err
	}

	out0 := *abi.ConvertType(out[0], new(Readtoken)).(*Readtoken)

	return out0, err

}

// ExternalGetToken is a free data retrieval call binding the contract method 0xe117176a.
//
// Solidity: function externalGetToken(address a) view returns((string,uint256,address))
func (_Read *ReadSession) ExternalGetToken(a common.Address) (Readtoken, error) {
	return _Read.Contract.ExternalGetToken(&_Read.CallOpts, a)
}

// ExternalGetToken is a free data retrieval call binding the contract method 0xe117176a.
//
// Solidity: function externalGetToken(address a) view returns((string,uint256,address))
func (_Read *ReadCallerSession) ExternalGetToken(a common.Address) (Readtoken, error) {
	return _Read.Contract.ExternalGetToken(&_Read.CallOpts, a)
}

// ExternalRead is a free data retrieval call binding the contract method 0x10f3f91a.
//
// Solidity: function externalRead() view returns(uint256)
func (_Read *ReadCaller) ExternalRead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "externalRead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExternalRead is a free data retrieval call binding the contract method 0x10f3f91a.
//
// Solidity: function externalRead() view returns(uint256)
func (_Read *ReadSession) ExternalRead() (*big.Int, error) {
	return _Read.Contract.ExternalRead(&_Read.CallOpts)
}

// ExternalRead is a free data retrieval call binding the contract method 0x10f3f91a.
//
// Solidity: function externalRead() view returns(uint256)
func (_Read *ReadCallerSession) ExternalRead() (*big.Int, error) {
	return _Read.Contract.ExternalRead(&_Read.CallOpts)
}

// ExternalReadWParams is a free data retrieval call binding the contract method 0x6e5298b2.
//
// Solidity: function externalReadWParams(uint256 p) view returns(uint256)
func (_Read *ReadCaller) ExternalReadWParams(opts *bind.CallOpts, p *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "externalReadWParams", p)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExternalReadWParams is a free data retrieval call binding the contract method 0x6e5298b2.
//
// Solidity: function externalReadWParams(uint256 p) view returns(uint256)
func (_Read *ReadSession) ExternalReadWParams(p *big.Int) (*big.Int, error) {
	return _Read.Contract.ExternalReadWParams(&_Read.CallOpts, p)
}

// ExternalReadWParams is a free data retrieval call binding the contract method 0x6e5298b2.
//
// Solidity: function externalReadWParams(uint256 p) view returns(uint256)
func (_Read *ReadCallerSession) ExternalReadWParams(p *big.Int) (*big.Int, error) {
	return _Read.Contract.ExternalReadWParams(&_Read.CallOpts, p)
}

// PublicGetOwnerName is a free data retrieval call binding the contract method 0x081540b2.
//
// Solidity: function publicGetOwnerName() view returns(string)
func (_Read *ReadCaller) PublicGetOwnerName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "publicGetOwnerName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// PublicGetOwnerName is a free data retrieval call binding the contract method 0x081540b2.
//
// Solidity: function publicGetOwnerName() view returns(string)
func (_Read *ReadSession) PublicGetOwnerName() (string, error) {
	return _Read.Contract.PublicGetOwnerName(&_Read.CallOpts)
}

// PublicGetOwnerName is a free data retrieval call binding the contract method 0x081540b2.
//
// Solidity: function publicGetOwnerName() view returns(string)
func (_Read *ReadCallerSession) PublicGetOwnerName() (string, error) {
	return _Read.Contract.PublicGetOwnerName(&_Read.CallOpts)
}

// PublicGetToken is a free data retrieval call binding the contract method 0xb1ceff4e.
//
// Solidity: function publicGetToken(address a) view returns((string,uint256,address))
func (_Read *ReadCaller) PublicGetToken(opts *bind.CallOpts, a common.Address) (Readtoken, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "publicGetToken", a)

	if err != nil {
		return *new(Readtoken), err
	}

	out0 := *abi.ConvertType(out[0], new(Readtoken)).(*Readtoken)

	return out0, err

}

// PublicGetToken is a free data retrieval call binding the contract method 0xb1ceff4e.
//
// Solidity: function publicGetToken(address a) view returns((string,uint256,address))
func (_Read *ReadSession) PublicGetToken(a common.Address) (Readtoken, error) {
	return _Read.Contract.PublicGetToken(&_Read.CallOpts, a)
}

// PublicGetToken is a free data retrieval call binding the contract method 0xb1ceff4e.
//
// Solidity: function publicGetToken(address a) view returns((string,uint256,address))
func (_Read *ReadCallerSession) PublicGetToken(a common.Address) (Readtoken, error) {
	return _Read.Contract.PublicGetToken(&_Read.CallOpts, a)
}

// PublicRead is a free data retrieval call binding the contract method 0xaf334757.
//
// Solidity: function publicRead() view returns(uint256)
func (_Read *ReadCaller) PublicRead(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "publicRead")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PublicRead is a free data retrieval call binding the contract method 0xaf334757.
//
// Solidity: function publicRead() view returns(uint256)
func (_Read *ReadSession) PublicRead() (*big.Int, error) {
	return _Read.Contract.PublicRead(&_Read.CallOpts)
}

// PublicRead is a free data retrieval call binding the contract method 0xaf334757.
//
// Solidity: function publicRead() view returns(uint256)
func (_Read *ReadCallerSession) PublicRead() (*big.Int, error) {
	return _Read.Contract.PublicRead(&_Read.CallOpts)
}

// PublicReadWParams is a free data retrieval call binding the contract method 0xbfa044ed.
//
// Solidity: function publicReadWParams(uint256 p) view returns(uint256)
func (_Read *ReadCaller) PublicReadWParams(opts *bind.CallOpts, p *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Read.contract.Call(opts, &out, "publicReadWParams", p)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PublicReadWParams is a free data retrieval call binding the contract method 0xbfa044ed.
//
// Solidity: function publicReadWParams(uint256 p) view returns(uint256)
func (_Read *ReadSession) PublicReadWParams(p *big.Int) (*big.Int, error) {
	return _Read.Contract.PublicReadWParams(&_Read.CallOpts, p)
}

// PublicReadWParams is a free data retrieval call binding the contract method 0xbfa044ed.
//
// Solidity: function publicReadWParams(uint256 p) view returns(uint256)
func (_Read *ReadCallerSession) PublicReadWParams(p *big.Int) (*big.Int, error) {
	return _Read.Contract.PublicReadWParams(&_Read.CallOpts, p)
}

// ExternalAddToken is a paid mutator transaction binding the contract method 0x9553400f.
//
// Solidity: function externalAddToken((string,uint256,address) t) returns()
func (_Read *ReadTransactor) ExternalAddToken(opts *bind.TransactOpts, t Readtoken) (*types.Transaction, error) {
	return _Read.contract.Transact(opts, "externalAddToken", t)
}

// ExternalAddToken is a paid mutator transaction binding the contract method 0x9553400f.
//
// Solidity: function externalAddToken((string,uint256,address) t) returns()
func (_Read *ReadSession) ExternalAddToken(t Readtoken) (*types.Transaction, error) {
	return _Read.Contract.ExternalAddToken(&_Read.TransactOpts, t)
}

// ExternalAddToken is a paid mutator transaction binding the contract method 0x9553400f.
//
// Solidity: function externalAddToken((string,uint256,address) t) returns()
func (_Read *ReadTransactorSession) ExternalAddToken(t Readtoken) (*types.Transaction, error) {
	return _Read.Contract.ExternalAddToken(&_Read.TransactOpts, t)
}

// PublicAddToken is a paid mutator transaction binding the contract method 0x33eb6be8.
//
// Solidity: function publicAddToken((string,uint256,address) t) returns()
func (_Read *ReadTransactor) PublicAddToken(opts *bind.TransactOpts, t Readtoken) (*types.Transaction, error) {
	return _Read.contract.Transact(opts, "publicAddToken", t)
}

// PublicAddToken is a paid mutator transaction binding the contract method 0x33eb6be8.
//
// Solidity: function publicAddToken((string,uint256,address) t) returns()
func (_Read *ReadSession) PublicAddToken(t Readtoken) (*types.Transaction, error) {
	return _Read.Contract.PublicAddToken(&_Read.TransactOpts, t)
}

// PublicAddToken is a paid mutator transaction binding the contract method 0x33eb6be8.
//
// Solidity: function publicAddToken((string,uint256,address) t) returns()
func (_Read *ReadTransactorSession) PublicAddToken(t Readtoken) (*types.Transaction, error) {
	return _Read.Contract.PublicAddToken(&_Read.TransactOpts, t)
}
