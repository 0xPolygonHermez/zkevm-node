// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package uniswap

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

// UniswapMetaData contains all meta data concerning the Uniswap contract.
var UniswapMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"indexed\":true,\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"amount\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"indexed\":true,\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"amount0\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"amount1\",\"type\":\"uint128\"}],\"name\":\"Collect\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"amount0\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"amount1\",\"type\":\"uint128\"}],\"name\":\"CollectProtocol\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"paid0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"paid1\",\"type\":\"uint256\"}],\"name\":\"Flash\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"observationCardinalityNextOld\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"observationCardinalityNextNew\",\"type\":\"uint16\"}],\"name\":\"IncreaseObservationCardinalityNext\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96\",\"type\":\"uint160\"},{\"indexed\":false,\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"}],\"name\":\"Initialize\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"indexed\":true,\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"amount\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"feeProtocol0Old\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"feeProtocol1Old\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"feeProtocol0New\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"feeProtocol1New\",\"type\":\"uint8\"}],\"name\":\"SetFeeProtocol\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"amount0\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"amount1\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96\",\"type\":\"uint160\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"liquidity\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"}],\"name\":\"Swap\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"internalType\":\"uint128\",\"name\":\"amount\",\"type\":\"uint128\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"internalType\":\"uint128\",\"name\":\"amount0Requested\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"amount1Requested\",\"type\":\"uint128\"}],\"name\":\"collect\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"amount0\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"amount1\",\"type\":\"uint128\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"amount0Requested\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"amount1Requested\",\"type\":\"uint128\"}],\"name\":\"collectProtocol\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"amount0\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"amount1\",\"type\":\"uint128\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fee\",\"outputs\":[{\"internalType\":\"uint24\",\"name\":\"\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeGrowthGlobal0X128\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeGrowthGlobal1X128\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"flash\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"observationCardinalityNext\",\"type\":\"uint16\"}],\"name\":\"increaseObservationCardinalityNext\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96\",\"type\":\"uint160\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"liquidity\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxLiquidityPerTick\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"internalType\":\"uint128\",\"name\":\"amount\",\"type\":\"uint128\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"observations\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"blockTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"int56\",\"name\":\"tickCumulative\",\"type\":\"int56\"},{\"internalType\":\"uint160\",\"name\":\"secondsPerLiquidityCumulativeX128\",\"type\":\"uint160\"},{\"internalType\":\"bool\",\"name\":\"initialized\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"secondsAgos\",\"type\":\"uint32[]\"}],\"name\":\"observe\",\"outputs\":[{\"internalType\":\"int56[]\",\"name\":\"tickCumulatives\",\"type\":\"int56[]\"},{\"internalType\":\"uint160[]\",\"name\":\"secondsPerLiquidityCumulativeX128s\",\"type\":\"uint160[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"positions\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"liquidity\",\"type\":\"uint128\"},{\"internalType\":\"uint256\",\"name\":\"feeGrowthInside0LastX128\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeGrowthInside1LastX128\",\"type\":\"uint256\"},{\"internalType\":\"uint128\",\"name\":\"tokensOwed0\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"tokensOwed1\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"protocolFees\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"token0\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"token1\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"feeProtocol0\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"feeProtocol1\",\"type\":\"uint8\"}],\"name\":\"setFeeProtocol\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"slot0\",\"outputs\":[{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96\",\"type\":\"uint160\"},{\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"},{\"internalType\":\"uint16\",\"name\":\"observationIndex\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"observationCardinality\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"observationCardinalityNext\",\"type\":\"uint16\"},{\"internalType\":\"uint8\",\"name\":\"feeProtocol\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"unlocked\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"}],\"name\":\"snapshotCumulativesInside\",\"outputs\":[{\"internalType\":\"int56\",\"name\":\"tickCumulativeInside\",\"type\":\"int56\"},{\"internalType\":\"uint160\",\"name\":\"secondsPerLiquidityInsideX128\",\"type\":\"uint160\"},{\"internalType\":\"uint32\",\"name\":\"secondsInside\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"zeroForOne\",\"type\":\"bool\"},{\"internalType\":\"int256\",\"name\":\"amountSpecified\",\"type\":\"int256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"swap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"amount0\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"amount1\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int16\",\"name\":\"\",\"type\":\"int16\"}],\"name\":\"tickBitmap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tickSpacing\",\"outputs\":[{\"internalType\":\"int24\",\"name\":\"\",\"type\":\"int24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int24\",\"name\":\"\",\"type\":\"int24\"}],\"name\":\"ticks\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"liquidityGross\",\"type\":\"uint128\"},{\"internalType\":\"int128\",\"name\":\"liquidityNet\",\"type\":\"int128\"},{\"internalType\":\"uint256\",\"name\":\"feeGrowthOutside0X128\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeGrowthOutside1X128\",\"type\":\"uint256\"},{\"internalType\":\"int56\",\"name\":\"tickCumulativeOutside\",\"type\":\"int56\"},{\"internalType\":\"uint160\",\"name\":\"secondsPerLiquidityOutsideX128\",\"type\":\"uint160\"},{\"internalType\":\"uint32\",\"name\":\"secondsOutside\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"initialized\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token0\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token1\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// UniswapABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapMetaData.ABI instead.
var UniswapABI = UniswapMetaData.ABI

// Uniswap is an auto generated Go binding around an Ethereum contract.
type Uniswap struct {
	UniswapCaller     // Read-only binding to the contract
	UniswapTransactor // Write-only binding to the contract
	UniswapFilterer   // Log filterer for contract events
}

// UniswapCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniswapCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniswapTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniswapFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniswapSession struct {
	Contract     *Uniswap          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UniswapCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniswapCallerSession struct {
	Contract *UniswapCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// UniswapTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniswapTransactorSession struct {
	Contract     *UniswapTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// UniswapRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniswapRaw struct {
	Contract *Uniswap // Generic contract binding to access the raw methods on
}

// UniswapCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniswapCallerRaw struct {
	Contract *UniswapCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniswapTransactorRaw struct {
	Contract *UniswapTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswap creates a new instance of Uniswap, bound to a specific deployed contract.
func NewUniswap(address common.Address, backend bind.ContractBackend) (*Uniswap, error) {
	contract, err := bindUniswap(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Uniswap{UniswapCaller: UniswapCaller{contract: contract}, UniswapTransactor: UniswapTransactor{contract: contract}, UniswapFilterer: UniswapFilterer{contract: contract}}, nil
}

// NewUniswapCaller creates a new read-only instance of Uniswap, bound to a specific deployed contract.
func NewUniswapCaller(address common.Address, caller bind.ContractCaller) (*UniswapCaller, error) {
	contract, err := bindUniswap(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapCaller{contract: contract}, nil
}

// NewUniswapTransactor creates a new write-only instance of Uniswap, bound to a specific deployed contract.
func NewUniswapTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapTransactor, error) {
	contract, err := bindUniswap(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapTransactor{contract: contract}, nil
}

// NewUniswapFilterer creates a new log filterer instance of Uniswap, bound to a specific deployed contract.
func NewUniswapFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapFilterer, error) {
	contract, err := bindUniswap(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapFilterer{contract: contract}, nil
}

// bindUniswap binds a generic wrapper to an already deployed contract.
func bindUniswap(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UniswapABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Uniswap *UniswapRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Uniswap.Contract.UniswapCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Uniswap *UniswapRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Uniswap.Contract.UniswapTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Uniswap *UniswapRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Uniswap.Contract.UniswapTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Uniswap *UniswapCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Uniswap.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Uniswap *UniswapTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Uniswap.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Uniswap *UniswapTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Uniswap.Contract.contract.Transact(opts, method, params...)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_Uniswap *UniswapCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_Uniswap *UniswapSession) Factory() (common.Address, error) {
	return _Uniswap.Contract.Factory(&_Uniswap.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_Uniswap *UniswapCallerSession) Factory() (common.Address, error) {
	return _Uniswap.Contract.Factory(&_Uniswap.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint24)
func (_Uniswap *UniswapCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint24)
func (_Uniswap *UniswapSession) Fee() (*big.Int, error) {
	return _Uniswap.Contract.Fee(&_Uniswap.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint24)
func (_Uniswap *UniswapCallerSession) Fee() (*big.Int, error) {
	return _Uniswap.Contract.Fee(&_Uniswap.CallOpts)
}

// FeeGrowthGlobal0X128 is a free data retrieval call binding the contract method 0xf3058399.
//
// Solidity: function feeGrowthGlobal0X128() view returns(uint256)
func (_Uniswap *UniswapCaller) FeeGrowthGlobal0X128(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "feeGrowthGlobal0X128")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeGrowthGlobal0X128 is a free data retrieval call binding the contract method 0xf3058399.
//
// Solidity: function feeGrowthGlobal0X128() view returns(uint256)
func (_Uniswap *UniswapSession) FeeGrowthGlobal0X128() (*big.Int, error) {
	return _Uniswap.Contract.FeeGrowthGlobal0X128(&_Uniswap.CallOpts)
}

// FeeGrowthGlobal0X128 is a free data retrieval call binding the contract method 0xf3058399.
//
// Solidity: function feeGrowthGlobal0X128() view returns(uint256)
func (_Uniswap *UniswapCallerSession) FeeGrowthGlobal0X128() (*big.Int, error) {
	return _Uniswap.Contract.FeeGrowthGlobal0X128(&_Uniswap.CallOpts)
}

// FeeGrowthGlobal1X128 is a free data retrieval call binding the contract method 0x46141319.
//
// Solidity: function feeGrowthGlobal1X128() view returns(uint256)
func (_Uniswap *UniswapCaller) FeeGrowthGlobal1X128(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "feeGrowthGlobal1X128")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeGrowthGlobal1X128 is a free data retrieval call binding the contract method 0x46141319.
//
// Solidity: function feeGrowthGlobal1X128() view returns(uint256)
func (_Uniswap *UniswapSession) FeeGrowthGlobal1X128() (*big.Int, error) {
	return _Uniswap.Contract.FeeGrowthGlobal1X128(&_Uniswap.CallOpts)
}

// FeeGrowthGlobal1X128 is a free data retrieval call binding the contract method 0x46141319.
//
// Solidity: function feeGrowthGlobal1X128() view returns(uint256)
func (_Uniswap *UniswapCallerSession) FeeGrowthGlobal1X128() (*big.Int, error) {
	return _Uniswap.Contract.FeeGrowthGlobal1X128(&_Uniswap.CallOpts)
}

// Liquidity is a free data retrieval call binding the contract method 0x1a686502.
//
// Solidity: function liquidity() view returns(uint128)
func (_Uniswap *UniswapCaller) Liquidity(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "liquidity")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Liquidity is a free data retrieval call binding the contract method 0x1a686502.
//
// Solidity: function liquidity() view returns(uint128)
func (_Uniswap *UniswapSession) Liquidity() (*big.Int, error) {
	return _Uniswap.Contract.Liquidity(&_Uniswap.CallOpts)
}

// Liquidity is a free data retrieval call binding the contract method 0x1a686502.
//
// Solidity: function liquidity() view returns(uint128)
func (_Uniswap *UniswapCallerSession) Liquidity() (*big.Int, error) {
	return _Uniswap.Contract.Liquidity(&_Uniswap.CallOpts)
}

// MaxLiquidityPerTick is a free data retrieval call binding the contract method 0x70cf754a.
//
// Solidity: function maxLiquidityPerTick() view returns(uint128)
func (_Uniswap *UniswapCaller) MaxLiquidityPerTick(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "maxLiquidityPerTick")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxLiquidityPerTick is a free data retrieval call binding the contract method 0x70cf754a.
//
// Solidity: function maxLiquidityPerTick() view returns(uint128)
func (_Uniswap *UniswapSession) MaxLiquidityPerTick() (*big.Int, error) {
	return _Uniswap.Contract.MaxLiquidityPerTick(&_Uniswap.CallOpts)
}

// MaxLiquidityPerTick is a free data retrieval call binding the contract method 0x70cf754a.
//
// Solidity: function maxLiquidityPerTick() view returns(uint128)
func (_Uniswap *UniswapCallerSession) MaxLiquidityPerTick() (*big.Int, error) {
	return _Uniswap.Contract.MaxLiquidityPerTick(&_Uniswap.CallOpts)
}

// Observations is a free data retrieval call binding the contract method 0x252c09d7.
//
// Solidity: function observations(uint256 ) view returns(uint32 blockTimestamp, int56 tickCumulative, uint160 secondsPerLiquidityCumulativeX128, bool initialized)
func (_Uniswap *UniswapCaller) Observations(opts *bind.CallOpts, arg0 *big.Int) (struct {
	BlockTimestamp                    uint32
	TickCumulative                    *big.Int
	SecondsPerLiquidityCumulativeX128 *big.Int
	Initialized                       bool
}, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "observations", arg0)

	outstruct := new(struct {
		BlockTimestamp                    uint32
		TickCumulative                    *big.Int
		SecondsPerLiquidityCumulativeX128 *big.Int
		Initialized                       bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BlockTimestamp = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.TickCumulative = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.SecondsPerLiquidityCumulativeX128 = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Initialized = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// Observations is a free data retrieval call binding the contract method 0x252c09d7.
//
// Solidity: function observations(uint256 ) view returns(uint32 blockTimestamp, int56 tickCumulative, uint160 secondsPerLiquidityCumulativeX128, bool initialized)
func (_Uniswap *UniswapSession) Observations(arg0 *big.Int) (struct {
	BlockTimestamp                    uint32
	TickCumulative                    *big.Int
	SecondsPerLiquidityCumulativeX128 *big.Int
	Initialized                       bool
}, error) {
	return _Uniswap.Contract.Observations(&_Uniswap.CallOpts, arg0)
}

// Observations is a free data retrieval call binding the contract method 0x252c09d7.
//
// Solidity: function observations(uint256 ) view returns(uint32 blockTimestamp, int56 tickCumulative, uint160 secondsPerLiquidityCumulativeX128, bool initialized)
func (_Uniswap *UniswapCallerSession) Observations(arg0 *big.Int) (struct {
	BlockTimestamp                    uint32
	TickCumulative                    *big.Int
	SecondsPerLiquidityCumulativeX128 *big.Int
	Initialized                       bool
}, error) {
	return _Uniswap.Contract.Observations(&_Uniswap.CallOpts, arg0)
}

// Observe is a free data retrieval call binding the contract method 0x883bdbfd.
//
// Solidity: function observe(uint32[] secondsAgos) view returns(int56[] tickCumulatives, uint160[] secondsPerLiquidityCumulativeX128s)
func (_Uniswap *UniswapCaller) Observe(opts *bind.CallOpts, secondsAgos []uint32) (struct {
	TickCumulatives                    []*big.Int
	SecondsPerLiquidityCumulativeX128s []*big.Int
}, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "observe", secondsAgos)

	outstruct := new(struct {
		TickCumulatives                    []*big.Int
		SecondsPerLiquidityCumulativeX128s []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TickCumulatives = *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	outstruct.SecondsPerLiquidityCumulativeX128s = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Observe is a free data retrieval call binding the contract method 0x883bdbfd.
//
// Solidity: function observe(uint32[] secondsAgos) view returns(int56[] tickCumulatives, uint160[] secondsPerLiquidityCumulativeX128s)
func (_Uniswap *UniswapSession) Observe(secondsAgos []uint32) (struct {
	TickCumulatives                    []*big.Int
	SecondsPerLiquidityCumulativeX128s []*big.Int
}, error) {
	return _Uniswap.Contract.Observe(&_Uniswap.CallOpts, secondsAgos)
}

// Observe is a free data retrieval call binding the contract method 0x883bdbfd.
//
// Solidity: function observe(uint32[] secondsAgos) view returns(int56[] tickCumulatives, uint160[] secondsPerLiquidityCumulativeX128s)
func (_Uniswap *UniswapCallerSession) Observe(secondsAgos []uint32) (struct {
	TickCumulatives                    []*big.Int
	SecondsPerLiquidityCumulativeX128s []*big.Int
}, error) {
	return _Uniswap.Contract.Observe(&_Uniswap.CallOpts, secondsAgos)
}

// Positions is a free data retrieval call binding the contract method 0x514ea4bf.
//
// Solidity: function positions(bytes32 ) view returns(uint128 liquidity, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128, uint128 tokensOwed0, uint128 tokensOwed1)
func (_Uniswap *UniswapCaller) Positions(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Liquidity                *big.Int
	FeeGrowthInside0LastX128 *big.Int
	FeeGrowthInside1LastX128 *big.Int
	TokensOwed0              *big.Int
	TokensOwed1              *big.Int
}, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "positions", arg0)

	outstruct := new(struct {
		Liquidity                *big.Int
		FeeGrowthInside0LastX128 *big.Int
		FeeGrowthInside1LastX128 *big.Int
		TokensOwed0              *big.Int
		TokensOwed1              *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Liquidity = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FeeGrowthInside0LastX128 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.FeeGrowthInside1LastX128 = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TokensOwed0 = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.TokensOwed1 = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Positions is a free data retrieval call binding the contract method 0x514ea4bf.
//
// Solidity: function positions(bytes32 ) view returns(uint128 liquidity, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128, uint128 tokensOwed0, uint128 tokensOwed1)
func (_Uniswap *UniswapSession) Positions(arg0 [32]byte) (struct {
	Liquidity                *big.Int
	FeeGrowthInside0LastX128 *big.Int
	FeeGrowthInside1LastX128 *big.Int
	TokensOwed0              *big.Int
	TokensOwed1              *big.Int
}, error) {
	return _Uniswap.Contract.Positions(&_Uniswap.CallOpts, arg0)
}

// Positions is a free data retrieval call binding the contract method 0x514ea4bf.
//
// Solidity: function positions(bytes32 ) view returns(uint128 liquidity, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128, uint128 tokensOwed0, uint128 tokensOwed1)
func (_Uniswap *UniswapCallerSession) Positions(arg0 [32]byte) (struct {
	Liquidity                *big.Int
	FeeGrowthInside0LastX128 *big.Int
	FeeGrowthInside1LastX128 *big.Int
	TokensOwed0              *big.Int
	TokensOwed1              *big.Int
}, error) {
	return _Uniswap.Contract.Positions(&_Uniswap.CallOpts, arg0)
}

// ProtocolFees is a free data retrieval call binding the contract method 0x1ad8b03b.
//
// Solidity: function protocolFees() view returns(uint128 token0, uint128 token1)
func (_Uniswap *UniswapCaller) ProtocolFees(opts *bind.CallOpts) (struct {
	Token0 *big.Int
	Token1 *big.Int
}, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "protocolFees")

	outstruct := new(struct {
		Token0 *big.Int
		Token1 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Token0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Token1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ProtocolFees is a free data retrieval call binding the contract method 0x1ad8b03b.
//
// Solidity: function protocolFees() view returns(uint128 token0, uint128 token1)
func (_Uniswap *UniswapSession) ProtocolFees() (struct {
	Token0 *big.Int
	Token1 *big.Int
}, error) {
	return _Uniswap.Contract.ProtocolFees(&_Uniswap.CallOpts)
}

// ProtocolFees is a free data retrieval call binding the contract method 0x1ad8b03b.
//
// Solidity: function protocolFees() view returns(uint128 token0, uint128 token1)
func (_Uniswap *UniswapCallerSession) ProtocolFees() (struct {
	Token0 *big.Int
	Token1 *big.Int
}, error) {
	return _Uniswap.Contract.ProtocolFees(&_Uniswap.CallOpts)
}

// Slot0 is a free data retrieval call binding the contract method 0x3850c7bd.
//
// Solidity: function slot0() view returns(uint160 sqrtPriceX96, int24 tick, uint16 observationIndex, uint16 observationCardinality, uint16 observationCardinalityNext, uint8 feeProtocol, bool unlocked)
func (_Uniswap *UniswapCaller) Slot0(opts *bind.CallOpts) (struct {
	SqrtPriceX96               *big.Int
	Tick                       *big.Int
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
	Unlocked                   bool
}, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "slot0")

	outstruct := new(struct {
		SqrtPriceX96               *big.Int
		Tick                       *big.Int
		ObservationIndex           uint16
		ObservationCardinality     uint16
		ObservationCardinalityNext uint16
		FeeProtocol                uint8
		Unlocked                   bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SqrtPriceX96 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Tick = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ObservationIndex = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.ObservationCardinality = *abi.ConvertType(out[3], new(uint16)).(*uint16)
	outstruct.ObservationCardinalityNext = *abi.ConvertType(out[4], new(uint16)).(*uint16)
	outstruct.FeeProtocol = *abi.ConvertType(out[5], new(uint8)).(*uint8)
	outstruct.Unlocked = *abi.ConvertType(out[6], new(bool)).(*bool)

	return *outstruct, err

}

// Slot0 is a free data retrieval call binding the contract method 0x3850c7bd.
//
// Solidity: function slot0() view returns(uint160 sqrtPriceX96, int24 tick, uint16 observationIndex, uint16 observationCardinality, uint16 observationCardinalityNext, uint8 feeProtocol, bool unlocked)
func (_Uniswap *UniswapSession) Slot0() (struct {
	SqrtPriceX96               *big.Int
	Tick                       *big.Int
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
	Unlocked                   bool
}, error) {
	return _Uniswap.Contract.Slot0(&_Uniswap.CallOpts)
}

// Slot0 is a free data retrieval call binding the contract method 0x3850c7bd.
//
// Solidity: function slot0() view returns(uint160 sqrtPriceX96, int24 tick, uint16 observationIndex, uint16 observationCardinality, uint16 observationCardinalityNext, uint8 feeProtocol, bool unlocked)
func (_Uniswap *UniswapCallerSession) Slot0() (struct {
	SqrtPriceX96               *big.Int
	Tick                       *big.Int
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
	Unlocked                   bool
}, error) {
	return _Uniswap.Contract.Slot0(&_Uniswap.CallOpts)
}

// SnapshotCumulativesInside is a free data retrieval call binding the contract method 0xa38807f2.
//
// Solidity: function snapshotCumulativesInside(int24 tickLower, int24 tickUpper) view returns(int56 tickCumulativeInside, uint160 secondsPerLiquidityInsideX128, uint32 secondsInside)
func (_Uniswap *UniswapCaller) SnapshotCumulativesInside(opts *bind.CallOpts, tickLower *big.Int, tickUpper *big.Int) (struct {
	TickCumulativeInside          *big.Int
	SecondsPerLiquidityInsideX128 *big.Int
	SecondsInside                 uint32
}, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "snapshotCumulativesInside", tickLower, tickUpper)

	outstruct := new(struct {
		TickCumulativeInside          *big.Int
		SecondsPerLiquidityInsideX128 *big.Int
		SecondsInside                 uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TickCumulativeInside = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.SecondsPerLiquidityInsideX128 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.SecondsInside = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// SnapshotCumulativesInside is a free data retrieval call binding the contract method 0xa38807f2.
//
// Solidity: function snapshotCumulativesInside(int24 tickLower, int24 tickUpper) view returns(int56 tickCumulativeInside, uint160 secondsPerLiquidityInsideX128, uint32 secondsInside)
func (_Uniswap *UniswapSession) SnapshotCumulativesInside(tickLower *big.Int, tickUpper *big.Int) (struct {
	TickCumulativeInside          *big.Int
	SecondsPerLiquidityInsideX128 *big.Int
	SecondsInside                 uint32
}, error) {
	return _Uniswap.Contract.SnapshotCumulativesInside(&_Uniswap.CallOpts, tickLower, tickUpper)
}

// SnapshotCumulativesInside is a free data retrieval call binding the contract method 0xa38807f2.
//
// Solidity: function snapshotCumulativesInside(int24 tickLower, int24 tickUpper) view returns(int56 tickCumulativeInside, uint160 secondsPerLiquidityInsideX128, uint32 secondsInside)
func (_Uniswap *UniswapCallerSession) SnapshotCumulativesInside(tickLower *big.Int, tickUpper *big.Int) (struct {
	TickCumulativeInside          *big.Int
	SecondsPerLiquidityInsideX128 *big.Int
	SecondsInside                 uint32
}, error) {
	return _Uniswap.Contract.SnapshotCumulativesInside(&_Uniswap.CallOpts, tickLower, tickUpper)
}

// TickBitmap is a free data retrieval call binding the contract method 0x5339c296.
//
// Solidity: function tickBitmap(int16 ) view returns(uint256)
func (_Uniswap *UniswapCaller) TickBitmap(opts *bind.CallOpts, arg0 int16) (*big.Int, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "tickBitmap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TickBitmap is a free data retrieval call binding the contract method 0x5339c296.
//
// Solidity: function tickBitmap(int16 ) view returns(uint256)
func (_Uniswap *UniswapSession) TickBitmap(arg0 int16) (*big.Int, error) {
	return _Uniswap.Contract.TickBitmap(&_Uniswap.CallOpts, arg0)
}

// TickBitmap is a free data retrieval call binding the contract method 0x5339c296.
//
// Solidity: function tickBitmap(int16 ) view returns(uint256)
func (_Uniswap *UniswapCallerSession) TickBitmap(arg0 int16) (*big.Int, error) {
	return _Uniswap.Contract.TickBitmap(&_Uniswap.CallOpts, arg0)
}

// TickSpacing is a free data retrieval call binding the contract method 0xd0c93a7c.
//
// Solidity: function tickSpacing() view returns(int24)
func (_Uniswap *UniswapCaller) TickSpacing(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "tickSpacing")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TickSpacing is a free data retrieval call binding the contract method 0xd0c93a7c.
//
// Solidity: function tickSpacing() view returns(int24)
func (_Uniswap *UniswapSession) TickSpacing() (*big.Int, error) {
	return _Uniswap.Contract.TickSpacing(&_Uniswap.CallOpts)
}

// TickSpacing is a free data retrieval call binding the contract method 0xd0c93a7c.
//
// Solidity: function tickSpacing() view returns(int24)
func (_Uniswap *UniswapCallerSession) TickSpacing() (*big.Int, error) {
	return _Uniswap.Contract.TickSpacing(&_Uniswap.CallOpts)
}

// Ticks is a free data retrieval call binding the contract method 0xf30dba93.
//
// Solidity: function ticks(int24 ) view returns(uint128 liquidityGross, int128 liquidityNet, uint256 feeGrowthOutside0X128, uint256 feeGrowthOutside1X128, int56 tickCumulativeOutside, uint160 secondsPerLiquidityOutsideX128, uint32 secondsOutside, bool initialized)
func (_Uniswap *UniswapCaller) Ticks(opts *bind.CallOpts, arg0 *big.Int) (struct {
	LiquidityGross                 *big.Int
	LiquidityNet                   *big.Int
	FeeGrowthOutside0X128          *big.Int
	FeeGrowthOutside1X128          *big.Int
	TickCumulativeOutside          *big.Int
	SecondsPerLiquidityOutsideX128 *big.Int
	SecondsOutside                 uint32
	Initialized                    bool
}, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "ticks", arg0)

	outstruct := new(struct {
		LiquidityGross                 *big.Int
		LiquidityNet                   *big.Int
		FeeGrowthOutside0X128          *big.Int
		FeeGrowthOutside1X128          *big.Int
		TickCumulativeOutside          *big.Int
		SecondsPerLiquidityOutsideX128 *big.Int
		SecondsOutside                 uint32
		Initialized                    bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.LiquidityGross = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.LiquidityNet = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.FeeGrowthOutside0X128 = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.FeeGrowthOutside1X128 = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.TickCumulativeOutside = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.SecondsPerLiquidityOutsideX128 = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.SecondsOutside = *abi.ConvertType(out[6], new(uint32)).(*uint32)
	outstruct.Initialized = *abi.ConvertType(out[7], new(bool)).(*bool)

	return *outstruct, err

}

// Ticks is a free data retrieval call binding the contract method 0xf30dba93.
//
// Solidity: function ticks(int24 ) view returns(uint128 liquidityGross, int128 liquidityNet, uint256 feeGrowthOutside0X128, uint256 feeGrowthOutside1X128, int56 tickCumulativeOutside, uint160 secondsPerLiquidityOutsideX128, uint32 secondsOutside, bool initialized)
func (_Uniswap *UniswapSession) Ticks(arg0 *big.Int) (struct {
	LiquidityGross                 *big.Int
	LiquidityNet                   *big.Int
	FeeGrowthOutside0X128          *big.Int
	FeeGrowthOutside1X128          *big.Int
	TickCumulativeOutside          *big.Int
	SecondsPerLiquidityOutsideX128 *big.Int
	SecondsOutside                 uint32
	Initialized                    bool
}, error) {
	return _Uniswap.Contract.Ticks(&_Uniswap.CallOpts, arg0)
}

// Ticks is a free data retrieval call binding the contract method 0xf30dba93.
//
// Solidity: function ticks(int24 ) view returns(uint128 liquidityGross, int128 liquidityNet, uint256 feeGrowthOutside0X128, uint256 feeGrowthOutside1X128, int56 tickCumulativeOutside, uint160 secondsPerLiquidityOutsideX128, uint32 secondsOutside, bool initialized)
func (_Uniswap *UniswapCallerSession) Ticks(arg0 *big.Int) (struct {
	LiquidityGross                 *big.Int
	LiquidityNet                   *big.Int
	FeeGrowthOutside0X128          *big.Int
	FeeGrowthOutside1X128          *big.Int
	TickCumulativeOutside          *big.Int
	SecondsPerLiquidityOutsideX128 *big.Int
	SecondsOutside                 uint32
	Initialized                    bool
}, error) {
	return _Uniswap.Contract.Ticks(&_Uniswap.CallOpts, arg0)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_Uniswap *UniswapCaller) Token0(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "token0")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_Uniswap *UniswapSession) Token0() (common.Address, error) {
	return _Uniswap.Contract.Token0(&_Uniswap.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_Uniswap *UniswapCallerSession) Token0() (common.Address, error) {
	return _Uniswap.Contract.Token0(&_Uniswap.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_Uniswap *UniswapCaller) Token1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Uniswap.contract.Call(opts, &out, "token1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_Uniswap *UniswapSession) Token1() (common.Address, error) {
	return _Uniswap.Contract.Token1(&_Uniswap.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_Uniswap *UniswapCallerSession) Token1() (common.Address, error) {
	return _Uniswap.Contract.Token1(&_Uniswap.CallOpts)
}

// Burn is a paid mutator transaction binding the contract method 0xa34123a7.
//
// Solidity: function burn(int24 tickLower, int24 tickUpper, uint128 amount) returns(uint256 amount0, uint256 amount1)
func (_Uniswap *UniswapTransactor) Burn(opts *bind.TransactOpts, tickLower *big.Int, tickUpper *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Uniswap.contract.Transact(opts, "burn", tickLower, tickUpper, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xa34123a7.
//
// Solidity: function burn(int24 tickLower, int24 tickUpper, uint128 amount) returns(uint256 amount0, uint256 amount1)
func (_Uniswap *UniswapSession) Burn(tickLower *big.Int, tickUpper *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Uniswap.Contract.Burn(&_Uniswap.TransactOpts, tickLower, tickUpper, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xa34123a7.
//
// Solidity: function burn(int24 tickLower, int24 tickUpper, uint128 amount) returns(uint256 amount0, uint256 amount1)
func (_Uniswap *UniswapTransactorSession) Burn(tickLower *big.Int, tickUpper *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Uniswap.Contract.Burn(&_Uniswap.TransactOpts, tickLower, tickUpper, amount)
}

// Collect is a paid mutator transaction binding the contract method 0x4f1eb3d8.
//
// Solidity: function collect(address recipient, int24 tickLower, int24 tickUpper, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_Uniswap *UniswapTransactor) Collect(opts *bind.TransactOpts, recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _Uniswap.contract.Transact(opts, "collect", recipient, tickLower, tickUpper, amount0Requested, amount1Requested)
}

// Collect is a paid mutator transaction binding the contract method 0x4f1eb3d8.
//
// Solidity: function collect(address recipient, int24 tickLower, int24 tickUpper, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_Uniswap *UniswapSession) Collect(recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _Uniswap.Contract.Collect(&_Uniswap.TransactOpts, recipient, tickLower, tickUpper, amount0Requested, amount1Requested)
}

// Collect is a paid mutator transaction binding the contract method 0x4f1eb3d8.
//
// Solidity: function collect(address recipient, int24 tickLower, int24 tickUpper, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_Uniswap *UniswapTransactorSession) Collect(recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _Uniswap.Contract.Collect(&_Uniswap.TransactOpts, recipient, tickLower, tickUpper, amount0Requested, amount1Requested)
}

// CollectProtocol is a paid mutator transaction binding the contract method 0x85b66729.
//
// Solidity: function collectProtocol(address recipient, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_Uniswap *UniswapTransactor) CollectProtocol(opts *bind.TransactOpts, recipient common.Address, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _Uniswap.contract.Transact(opts, "collectProtocol", recipient, amount0Requested, amount1Requested)
}

// CollectProtocol is a paid mutator transaction binding the contract method 0x85b66729.
//
// Solidity: function collectProtocol(address recipient, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_Uniswap *UniswapSession) CollectProtocol(recipient common.Address, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _Uniswap.Contract.CollectProtocol(&_Uniswap.TransactOpts, recipient, amount0Requested, amount1Requested)
}

// CollectProtocol is a paid mutator transaction binding the contract method 0x85b66729.
//
// Solidity: function collectProtocol(address recipient, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_Uniswap *UniswapTransactorSession) CollectProtocol(recipient common.Address, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _Uniswap.Contract.CollectProtocol(&_Uniswap.TransactOpts, recipient, amount0Requested, amount1Requested)
}

// Flash is a paid mutator transaction binding the contract method 0x490e6cbc.
//
// Solidity: function flash(address recipient, uint256 amount0, uint256 amount1, bytes data) returns()
func (_Uniswap *UniswapTransactor) Flash(opts *bind.TransactOpts, recipient common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _Uniswap.contract.Transact(opts, "flash", recipient, amount0, amount1, data)
}

// Flash is a paid mutator transaction binding the contract method 0x490e6cbc.
//
// Solidity: function flash(address recipient, uint256 amount0, uint256 amount1, bytes data) returns()
func (_Uniswap *UniswapSession) Flash(recipient common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _Uniswap.Contract.Flash(&_Uniswap.TransactOpts, recipient, amount0, amount1, data)
}

// Flash is a paid mutator transaction binding the contract method 0x490e6cbc.
//
// Solidity: function flash(address recipient, uint256 amount0, uint256 amount1, bytes data) returns()
func (_Uniswap *UniswapTransactorSession) Flash(recipient common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _Uniswap.Contract.Flash(&_Uniswap.TransactOpts, recipient, amount0, amount1, data)
}

// IncreaseObservationCardinalityNext is a paid mutator transaction binding the contract method 0x32148f67.
//
// Solidity: function increaseObservationCardinalityNext(uint16 observationCardinalityNext) returns()
func (_Uniswap *UniswapTransactor) IncreaseObservationCardinalityNext(opts *bind.TransactOpts, observationCardinalityNext uint16) (*types.Transaction, error) {
	return _Uniswap.contract.Transact(opts, "increaseObservationCardinalityNext", observationCardinalityNext)
}

// IncreaseObservationCardinalityNext is a paid mutator transaction binding the contract method 0x32148f67.
//
// Solidity: function increaseObservationCardinalityNext(uint16 observationCardinalityNext) returns()
func (_Uniswap *UniswapSession) IncreaseObservationCardinalityNext(observationCardinalityNext uint16) (*types.Transaction, error) {
	return _Uniswap.Contract.IncreaseObservationCardinalityNext(&_Uniswap.TransactOpts, observationCardinalityNext)
}

// IncreaseObservationCardinalityNext is a paid mutator transaction binding the contract method 0x32148f67.
//
// Solidity: function increaseObservationCardinalityNext(uint16 observationCardinalityNext) returns()
func (_Uniswap *UniswapTransactorSession) IncreaseObservationCardinalityNext(observationCardinalityNext uint16) (*types.Transaction, error) {
	return _Uniswap.Contract.IncreaseObservationCardinalityNext(&_Uniswap.TransactOpts, observationCardinalityNext)
}

// Initialize is a paid mutator transaction binding the contract method 0xf637731d.
//
// Solidity: function initialize(uint160 sqrtPriceX96) returns()
func (_Uniswap *UniswapTransactor) Initialize(opts *bind.TransactOpts, sqrtPriceX96 *big.Int) (*types.Transaction, error) {
	return _Uniswap.contract.Transact(opts, "initialize", sqrtPriceX96)
}

// Initialize is a paid mutator transaction binding the contract method 0xf637731d.
//
// Solidity: function initialize(uint160 sqrtPriceX96) returns()
func (_Uniswap *UniswapSession) Initialize(sqrtPriceX96 *big.Int) (*types.Transaction, error) {
	return _Uniswap.Contract.Initialize(&_Uniswap.TransactOpts, sqrtPriceX96)
}

// Initialize is a paid mutator transaction binding the contract method 0xf637731d.
//
// Solidity: function initialize(uint160 sqrtPriceX96) returns()
func (_Uniswap *UniswapTransactorSession) Initialize(sqrtPriceX96 *big.Int) (*types.Transaction, error) {
	return _Uniswap.Contract.Initialize(&_Uniswap.TransactOpts, sqrtPriceX96)
}

// Mint is a paid mutator transaction binding the contract method 0x3c8a7d8d.
//
// Solidity: function mint(address recipient, int24 tickLower, int24 tickUpper, uint128 amount, bytes data) returns(uint256 amount0, uint256 amount1)
func (_Uniswap *UniswapTransactor) Mint(opts *bind.TransactOpts, recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Uniswap.contract.Transact(opts, "mint", recipient, tickLower, tickUpper, amount, data)
}

// Mint is a paid mutator transaction binding the contract method 0x3c8a7d8d.
//
// Solidity: function mint(address recipient, int24 tickLower, int24 tickUpper, uint128 amount, bytes data) returns(uint256 amount0, uint256 amount1)
func (_Uniswap *UniswapSession) Mint(recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Uniswap.Contract.Mint(&_Uniswap.TransactOpts, recipient, tickLower, tickUpper, amount, data)
}

// Mint is a paid mutator transaction binding the contract method 0x3c8a7d8d.
//
// Solidity: function mint(address recipient, int24 tickLower, int24 tickUpper, uint128 amount, bytes data) returns(uint256 amount0, uint256 amount1)
func (_Uniswap *UniswapTransactorSession) Mint(recipient common.Address, tickLower *big.Int, tickUpper *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Uniswap.Contract.Mint(&_Uniswap.TransactOpts, recipient, tickLower, tickUpper, amount, data)
}

// SetFeeProtocol is a paid mutator transaction binding the contract method 0x8206a4d1.
//
// Solidity: function setFeeProtocol(uint8 feeProtocol0, uint8 feeProtocol1) returns()
func (_Uniswap *UniswapTransactor) SetFeeProtocol(opts *bind.TransactOpts, feeProtocol0 uint8, feeProtocol1 uint8) (*types.Transaction, error) {
	return _Uniswap.contract.Transact(opts, "setFeeProtocol", feeProtocol0, feeProtocol1)
}

// SetFeeProtocol is a paid mutator transaction binding the contract method 0x8206a4d1.
//
// Solidity: function setFeeProtocol(uint8 feeProtocol0, uint8 feeProtocol1) returns()
func (_Uniswap *UniswapSession) SetFeeProtocol(feeProtocol0 uint8, feeProtocol1 uint8) (*types.Transaction, error) {
	return _Uniswap.Contract.SetFeeProtocol(&_Uniswap.TransactOpts, feeProtocol0, feeProtocol1)
}

// SetFeeProtocol is a paid mutator transaction binding the contract method 0x8206a4d1.
//
// Solidity: function setFeeProtocol(uint8 feeProtocol0, uint8 feeProtocol1) returns()
func (_Uniswap *UniswapTransactorSession) SetFeeProtocol(feeProtocol0 uint8, feeProtocol1 uint8) (*types.Transaction, error) {
	return _Uniswap.Contract.SetFeeProtocol(&_Uniswap.TransactOpts, feeProtocol0, feeProtocol1)
}

// Swap is a paid mutator transaction binding the contract method 0x128acb08.
//
// Solidity: function swap(address recipient, bool zeroForOne, int256 amountSpecified, uint160 sqrtPriceLimitX96, bytes data) returns(int256 amount0, int256 amount1)
func (_Uniswap *UniswapTransactor) Swap(opts *bind.TransactOpts, recipient common.Address, zeroForOne bool, amountSpecified *big.Int, sqrtPriceLimitX96 *big.Int, data []byte) (*types.Transaction, error) {
	return _Uniswap.contract.Transact(opts, "swap", recipient, zeroForOne, amountSpecified, sqrtPriceLimitX96, data)
}

// Swap is a paid mutator transaction binding the contract method 0x128acb08.
//
// Solidity: function swap(address recipient, bool zeroForOne, int256 amountSpecified, uint160 sqrtPriceLimitX96, bytes data) returns(int256 amount0, int256 amount1)
func (_Uniswap *UniswapSession) Swap(recipient common.Address, zeroForOne bool, amountSpecified *big.Int, sqrtPriceLimitX96 *big.Int, data []byte) (*types.Transaction, error) {
	return _Uniswap.Contract.Swap(&_Uniswap.TransactOpts, recipient, zeroForOne, amountSpecified, sqrtPriceLimitX96, data)
}

// Swap is a paid mutator transaction binding the contract method 0x128acb08.
//
// Solidity: function swap(address recipient, bool zeroForOne, int256 amountSpecified, uint160 sqrtPriceLimitX96, bytes data) returns(int256 amount0, int256 amount1)
func (_Uniswap *UniswapTransactorSession) Swap(recipient common.Address, zeroForOne bool, amountSpecified *big.Int, sqrtPriceLimitX96 *big.Int, data []byte) (*types.Transaction, error) {
	return _Uniswap.Contract.Swap(&_Uniswap.TransactOpts, recipient, zeroForOne, amountSpecified, sqrtPriceLimitX96, data)
}

// UniswapBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the Uniswap contract.
type UniswapBurnIterator struct {
	Event *UniswapBurn // Event containing the contract specifics and raw log

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
func (it *UniswapBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapBurn)
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
		it.Event = new(UniswapBurn)
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
func (it *UniswapBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapBurn represents a Burn event raised by the Uniswap contract.
type UniswapBurn struct {
	Owner     common.Address
	TickLower *big.Int
	TickUpper *big.Int
	Amount    *big.Int
	Amount0   *big.Int
	Amount1   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0x0c396cd989a39f4459b5fa1aed6a9a8dcdbc45908acfd67e028cd568da98982c.
//
// Solidity: event Burn(address indexed owner, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount, uint256 amount0, uint256 amount1)
func (_Uniswap *UniswapFilterer) FilterBurn(opts *bind.FilterOpts, owner []common.Address, tickLower []*big.Int, tickUpper []*big.Int) (*UniswapBurnIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var tickLowerRule []interface{}
	for _, tickLowerItem := range tickLower {
		tickLowerRule = append(tickLowerRule, tickLowerItem)
	}
	var tickUpperRule []interface{}
	for _, tickUpperItem := range tickUpper {
		tickUpperRule = append(tickUpperRule, tickUpperItem)
	}

	logs, sub, err := _Uniswap.contract.FilterLogs(opts, "Burn", ownerRule, tickLowerRule, tickUpperRule)
	if err != nil {
		return nil, err
	}
	return &UniswapBurnIterator{contract: _Uniswap.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0x0c396cd989a39f4459b5fa1aed6a9a8dcdbc45908acfd67e028cd568da98982c.
//
// Solidity: event Burn(address indexed owner, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount, uint256 amount0, uint256 amount1)
func (_Uniswap *UniswapFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *UniswapBurn, owner []common.Address, tickLower []*big.Int, tickUpper []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var tickLowerRule []interface{}
	for _, tickLowerItem := range tickLower {
		tickLowerRule = append(tickLowerRule, tickLowerItem)
	}
	var tickUpperRule []interface{}
	for _, tickUpperItem := range tickUpper {
		tickUpperRule = append(tickUpperRule, tickUpperItem)
	}

	logs, sub, err := _Uniswap.contract.WatchLogs(opts, "Burn", ownerRule, tickLowerRule, tickUpperRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapBurn)
				if err := _Uniswap.contract.UnpackLog(event, "Burn", log); err != nil {
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

// ParseBurn is a log parse operation binding the contract event 0x0c396cd989a39f4459b5fa1aed6a9a8dcdbc45908acfd67e028cd568da98982c.
//
// Solidity: event Burn(address indexed owner, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount, uint256 amount0, uint256 amount1)
func (_Uniswap *UniswapFilterer) ParseBurn(log types.Log) (*UniswapBurn, error) {
	event := new(UniswapBurn)
	if err := _Uniswap.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapCollectIterator is returned from FilterCollect and is used to iterate over the raw logs and unpacked data for Collect events raised by the Uniswap contract.
type UniswapCollectIterator struct {
	Event *UniswapCollect // Event containing the contract specifics and raw log

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
func (it *UniswapCollectIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapCollect)
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
		it.Event = new(UniswapCollect)
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
func (it *UniswapCollectIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapCollectIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapCollect represents a Collect event raised by the Uniswap contract.
type UniswapCollect struct {
	Owner     common.Address
	Recipient common.Address
	TickLower *big.Int
	TickUpper *big.Int
	Amount0   *big.Int
	Amount1   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCollect is a free log retrieval operation binding the contract event 0x70935338e69775456a85ddef226c395fb668b63fa0115f5f20610b388e6ca9c0.
//
// Solidity: event Collect(address indexed owner, address recipient, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount0, uint128 amount1)
func (_Uniswap *UniswapFilterer) FilterCollect(opts *bind.FilterOpts, owner []common.Address, tickLower []*big.Int, tickUpper []*big.Int) (*UniswapCollectIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	var tickLowerRule []interface{}
	for _, tickLowerItem := range tickLower {
		tickLowerRule = append(tickLowerRule, tickLowerItem)
	}
	var tickUpperRule []interface{}
	for _, tickUpperItem := range tickUpper {
		tickUpperRule = append(tickUpperRule, tickUpperItem)
	}

	logs, sub, err := _Uniswap.contract.FilterLogs(opts, "Collect", ownerRule, tickLowerRule, tickUpperRule)
	if err != nil {
		return nil, err
	}
	return &UniswapCollectIterator{contract: _Uniswap.contract, event: "Collect", logs: logs, sub: sub}, nil
}

// WatchCollect is a free log subscription operation binding the contract event 0x70935338e69775456a85ddef226c395fb668b63fa0115f5f20610b388e6ca9c0.
//
// Solidity: event Collect(address indexed owner, address recipient, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount0, uint128 amount1)
func (_Uniswap *UniswapFilterer) WatchCollect(opts *bind.WatchOpts, sink chan<- *UniswapCollect, owner []common.Address, tickLower []*big.Int, tickUpper []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	var tickLowerRule []interface{}
	for _, tickLowerItem := range tickLower {
		tickLowerRule = append(tickLowerRule, tickLowerItem)
	}
	var tickUpperRule []interface{}
	for _, tickUpperItem := range tickUpper {
		tickUpperRule = append(tickUpperRule, tickUpperItem)
	}

	logs, sub, err := _Uniswap.contract.WatchLogs(opts, "Collect", ownerRule, tickLowerRule, tickUpperRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapCollect)
				if err := _Uniswap.contract.UnpackLog(event, "Collect", log); err != nil {
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

// ParseCollect is a log parse operation binding the contract event 0x70935338e69775456a85ddef226c395fb668b63fa0115f5f20610b388e6ca9c0.
//
// Solidity: event Collect(address indexed owner, address recipient, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount0, uint128 amount1)
func (_Uniswap *UniswapFilterer) ParseCollect(log types.Log) (*UniswapCollect, error) {
	event := new(UniswapCollect)
	if err := _Uniswap.contract.UnpackLog(event, "Collect", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapCollectProtocolIterator is returned from FilterCollectProtocol and is used to iterate over the raw logs and unpacked data for CollectProtocol events raised by the Uniswap contract.
type UniswapCollectProtocolIterator struct {
	Event *UniswapCollectProtocol // Event containing the contract specifics and raw log

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
func (it *UniswapCollectProtocolIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapCollectProtocol)
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
		it.Event = new(UniswapCollectProtocol)
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
func (it *UniswapCollectProtocolIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapCollectProtocolIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapCollectProtocol represents a CollectProtocol event raised by the Uniswap contract.
type UniswapCollectProtocol struct {
	Sender    common.Address
	Recipient common.Address
	Amount0   *big.Int
	Amount1   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCollectProtocol is a free log retrieval operation binding the contract event 0x596b573906218d3411850b26a6b437d6c4522fdb43d2d2386263f86d50b8b151.
//
// Solidity: event CollectProtocol(address indexed sender, address indexed recipient, uint128 amount0, uint128 amount1)
func (_Uniswap *UniswapFilterer) FilterCollectProtocol(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*UniswapCollectProtocolIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Uniswap.contract.FilterLogs(opts, "CollectProtocol", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &UniswapCollectProtocolIterator{contract: _Uniswap.contract, event: "CollectProtocol", logs: logs, sub: sub}, nil
}

// WatchCollectProtocol is a free log subscription operation binding the contract event 0x596b573906218d3411850b26a6b437d6c4522fdb43d2d2386263f86d50b8b151.
//
// Solidity: event CollectProtocol(address indexed sender, address indexed recipient, uint128 amount0, uint128 amount1)
func (_Uniswap *UniswapFilterer) WatchCollectProtocol(opts *bind.WatchOpts, sink chan<- *UniswapCollectProtocol, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Uniswap.contract.WatchLogs(opts, "CollectProtocol", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapCollectProtocol)
				if err := _Uniswap.contract.UnpackLog(event, "CollectProtocol", log); err != nil {
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

// ParseCollectProtocol is a log parse operation binding the contract event 0x596b573906218d3411850b26a6b437d6c4522fdb43d2d2386263f86d50b8b151.
//
// Solidity: event CollectProtocol(address indexed sender, address indexed recipient, uint128 amount0, uint128 amount1)
func (_Uniswap *UniswapFilterer) ParseCollectProtocol(log types.Log) (*UniswapCollectProtocol, error) {
	event := new(UniswapCollectProtocol)
	if err := _Uniswap.contract.UnpackLog(event, "CollectProtocol", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapFlashIterator is returned from FilterFlash and is used to iterate over the raw logs and unpacked data for Flash events raised by the Uniswap contract.
type UniswapFlashIterator struct {
	Event *UniswapFlash // Event containing the contract specifics and raw log

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
func (it *UniswapFlashIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapFlash)
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
		it.Event = new(UniswapFlash)
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
func (it *UniswapFlashIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapFlashIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapFlash represents a Flash event raised by the Uniswap contract.
type UniswapFlash struct {
	Sender    common.Address
	Recipient common.Address
	Amount0   *big.Int
	Amount1   *big.Int
	Paid0     *big.Int
	Paid1     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFlash is a free log retrieval operation binding the contract event 0xbdbdb71d7860376ba52b25a5028beea23581364a40522f6bcfb86bb1f2dca633.
//
// Solidity: event Flash(address indexed sender, address indexed recipient, uint256 amount0, uint256 amount1, uint256 paid0, uint256 paid1)
func (_Uniswap *UniswapFilterer) FilterFlash(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*UniswapFlashIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Uniswap.contract.FilterLogs(opts, "Flash", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &UniswapFlashIterator{contract: _Uniswap.contract, event: "Flash", logs: logs, sub: sub}, nil
}

// WatchFlash is a free log subscription operation binding the contract event 0xbdbdb71d7860376ba52b25a5028beea23581364a40522f6bcfb86bb1f2dca633.
//
// Solidity: event Flash(address indexed sender, address indexed recipient, uint256 amount0, uint256 amount1, uint256 paid0, uint256 paid1)
func (_Uniswap *UniswapFilterer) WatchFlash(opts *bind.WatchOpts, sink chan<- *UniswapFlash, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Uniswap.contract.WatchLogs(opts, "Flash", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapFlash)
				if err := _Uniswap.contract.UnpackLog(event, "Flash", log); err != nil {
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

// ParseFlash is a log parse operation binding the contract event 0xbdbdb71d7860376ba52b25a5028beea23581364a40522f6bcfb86bb1f2dca633.
//
// Solidity: event Flash(address indexed sender, address indexed recipient, uint256 amount0, uint256 amount1, uint256 paid0, uint256 paid1)
func (_Uniswap *UniswapFilterer) ParseFlash(log types.Log) (*UniswapFlash, error) {
	event := new(UniswapFlash)
	if err := _Uniswap.contract.UnpackLog(event, "Flash", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapIncreaseObservationCardinalityNextIterator is returned from FilterIncreaseObservationCardinalityNext and is used to iterate over the raw logs and unpacked data for IncreaseObservationCardinalityNext events raised by the Uniswap contract.
type UniswapIncreaseObservationCardinalityNextIterator struct {
	Event *UniswapIncreaseObservationCardinalityNext // Event containing the contract specifics and raw log

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
func (it *UniswapIncreaseObservationCardinalityNextIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapIncreaseObservationCardinalityNext)
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
		it.Event = new(UniswapIncreaseObservationCardinalityNext)
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
func (it *UniswapIncreaseObservationCardinalityNextIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapIncreaseObservationCardinalityNextIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapIncreaseObservationCardinalityNext represents a IncreaseObservationCardinalityNext event raised by the Uniswap contract.
type UniswapIncreaseObservationCardinalityNext struct {
	ObservationCardinalityNextOld uint16
	ObservationCardinalityNextNew uint16
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterIncreaseObservationCardinalityNext is a free log retrieval operation binding the contract event 0xac49e518f90a358f652e4400164f05a5d8f7e35e7747279bc3a93dbf584e125a.
//
// Solidity: event IncreaseObservationCardinalityNext(uint16 observationCardinalityNextOld, uint16 observationCardinalityNextNew)
func (_Uniswap *UniswapFilterer) FilterIncreaseObservationCardinalityNext(opts *bind.FilterOpts) (*UniswapIncreaseObservationCardinalityNextIterator, error) {

	logs, sub, err := _Uniswap.contract.FilterLogs(opts, "IncreaseObservationCardinalityNext")
	if err != nil {
		return nil, err
	}
	return &UniswapIncreaseObservationCardinalityNextIterator{contract: _Uniswap.contract, event: "IncreaseObservationCardinalityNext", logs: logs, sub: sub}, nil
}

// WatchIncreaseObservationCardinalityNext is a free log subscription operation binding the contract event 0xac49e518f90a358f652e4400164f05a5d8f7e35e7747279bc3a93dbf584e125a.
//
// Solidity: event IncreaseObservationCardinalityNext(uint16 observationCardinalityNextOld, uint16 observationCardinalityNextNew)
func (_Uniswap *UniswapFilterer) WatchIncreaseObservationCardinalityNext(opts *bind.WatchOpts, sink chan<- *UniswapIncreaseObservationCardinalityNext) (event.Subscription, error) {

	logs, sub, err := _Uniswap.contract.WatchLogs(opts, "IncreaseObservationCardinalityNext")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapIncreaseObservationCardinalityNext)
				if err := _Uniswap.contract.UnpackLog(event, "IncreaseObservationCardinalityNext", log); err != nil {
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

// ParseIncreaseObservationCardinalityNext is a log parse operation binding the contract event 0xac49e518f90a358f652e4400164f05a5d8f7e35e7747279bc3a93dbf584e125a.
//
// Solidity: event IncreaseObservationCardinalityNext(uint16 observationCardinalityNextOld, uint16 observationCardinalityNextNew)
func (_Uniswap *UniswapFilterer) ParseIncreaseObservationCardinalityNext(log types.Log) (*UniswapIncreaseObservationCardinalityNext, error) {
	event := new(UniswapIncreaseObservationCardinalityNext)
	if err := _Uniswap.contract.UnpackLog(event, "IncreaseObservationCardinalityNext", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapInitializeIterator is returned from FilterInitialize and is used to iterate over the raw logs and unpacked data for Initialize events raised by the Uniswap contract.
type UniswapInitializeIterator struct {
	Event *UniswapInitialize // Event containing the contract specifics and raw log

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
func (it *UniswapInitializeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapInitialize)
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
		it.Event = new(UniswapInitialize)
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
func (it *UniswapInitializeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapInitializeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapInitialize represents a Initialize event raised by the Uniswap contract.
type UniswapInitialize struct {
	SqrtPriceX96 *big.Int
	Tick         *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterInitialize is a free log retrieval operation binding the contract event 0x98636036cb66a9c19a37435efc1e90142190214e8abeb821bdba3f2990dd4c95.
//
// Solidity: event Initialize(uint160 sqrtPriceX96, int24 tick)
func (_Uniswap *UniswapFilterer) FilterInitialize(opts *bind.FilterOpts) (*UniswapInitializeIterator, error) {

	logs, sub, err := _Uniswap.contract.FilterLogs(opts, "Initialize")
	if err != nil {
		return nil, err
	}
	return &UniswapInitializeIterator{contract: _Uniswap.contract, event: "Initialize", logs: logs, sub: sub}, nil
}

// WatchInitialize is a free log subscription operation binding the contract event 0x98636036cb66a9c19a37435efc1e90142190214e8abeb821bdba3f2990dd4c95.
//
// Solidity: event Initialize(uint160 sqrtPriceX96, int24 tick)
func (_Uniswap *UniswapFilterer) WatchInitialize(opts *bind.WatchOpts, sink chan<- *UniswapInitialize) (event.Subscription, error) {

	logs, sub, err := _Uniswap.contract.WatchLogs(opts, "Initialize")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapInitialize)
				if err := _Uniswap.contract.UnpackLog(event, "Initialize", log); err != nil {
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

// ParseInitialize is a log parse operation binding the contract event 0x98636036cb66a9c19a37435efc1e90142190214e8abeb821bdba3f2990dd4c95.
//
// Solidity: event Initialize(uint160 sqrtPriceX96, int24 tick)
func (_Uniswap *UniswapFilterer) ParseInitialize(log types.Log) (*UniswapInitialize, error) {
	event := new(UniswapInitialize)
	if err := _Uniswap.contract.UnpackLog(event, "Initialize", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the Uniswap contract.
type UniswapMintIterator struct {
	Event *UniswapMint // Event containing the contract specifics and raw log

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
func (it *UniswapMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapMint)
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
		it.Event = new(UniswapMint)
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
func (it *UniswapMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapMint represents a Mint event raised by the Uniswap contract.
type UniswapMint struct {
	Sender    common.Address
	Owner     common.Address
	TickLower *big.Int
	TickUpper *big.Int
	Amount    *big.Int
	Amount0   *big.Int
	Amount1   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x7a53080ba414158be7ec69b987b5fb7d07dee101fe85488f0853ae16239d0bde.
//
// Solidity: event Mint(address sender, address indexed owner, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount, uint256 amount0, uint256 amount1)
func (_Uniswap *UniswapFilterer) FilterMint(opts *bind.FilterOpts, owner []common.Address, tickLower []*big.Int, tickUpper []*big.Int) (*UniswapMintIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var tickLowerRule []interface{}
	for _, tickLowerItem := range tickLower {
		tickLowerRule = append(tickLowerRule, tickLowerItem)
	}
	var tickUpperRule []interface{}
	for _, tickUpperItem := range tickUpper {
		tickUpperRule = append(tickUpperRule, tickUpperItem)
	}

	logs, sub, err := _Uniswap.contract.FilterLogs(opts, "Mint", ownerRule, tickLowerRule, tickUpperRule)
	if err != nil {
		return nil, err
	}
	return &UniswapMintIterator{contract: _Uniswap.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x7a53080ba414158be7ec69b987b5fb7d07dee101fe85488f0853ae16239d0bde.
//
// Solidity: event Mint(address sender, address indexed owner, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount, uint256 amount0, uint256 amount1)
func (_Uniswap *UniswapFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *UniswapMint, owner []common.Address, tickLower []*big.Int, tickUpper []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var tickLowerRule []interface{}
	for _, tickLowerItem := range tickLower {
		tickLowerRule = append(tickLowerRule, tickLowerItem)
	}
	var tickUpperRule []interface{}
	for _, tickUpperItem := range tickUpper {
		tickUpperRule = append(tickUpperRule, tickUpperItem)
	}

	logs, sub, err := _Uniswap.contract.WatchLogs(opts, "Mint", ownerRule, tickLowerRule, tickUpperRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapMint)
				if err := _Uniswap.contract.UnpackLog(event, "Mint", log); err != nil {
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

// ParseMint is a log parse operation binding the contract event 0x7a53080ba414158be7ec69b987b5fb7d07dee101fe85488f0853ae16239d0bde.
//
// Solidity: event Mint(address sender, address indexed owner, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount, uint256 amount0, uint256 amount1)
func (_Uniswap *UniswapFilterer) ParseMint(log types.Log) (*UniswapMint, error) {
	event := new(UniswapMint)
	if err := _Uniswap.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapSetFeeProtocolIterator is returned from FilterSetFeeProtocol and is used to iterate over the raw logs and unpacked data for SetFeeProtocol events raised by the Uniswap contract.
type UniswapSetFeeProtocolIterator struct {
	Event *UniswapSetFeeProtocol // Event containing the contract specifics and raw log

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
func (it *UniswapSetFeeProtocolIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapSetFeeProtocol)
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
		it.Event = new(UniswapSetFeeProtocol)
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
func (it *UniswapSetFeeProtocolIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapSetFeeProtocolIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapSetFeeProtocol represents a SetFeeProtocol event raised by the Uniswap contract.
type UniswapSetFeeProtocol struct {
	FeeProtocol0Old uint8
	FeeProtocol1Old uint8
	FeeProtocol0New uint8
	FeeProtocol1New uint8
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterSetFeeProtocol is a free log retrieval operation binding the contract event 0x973d8d92bb299f4af6ce49b52a8adb85ae46b9f214c4c4fc06ac77401237b133.
//
// Solidity: event SetFeeProtocol(uint8 feeProtocol0Old, uint8 feeProtocol1Old, uint8 feeProtocol0New, uint8 feeProtocol1New)
func (_Uniswap *UniswapFilterer) FilterSetFeeProtocol(opts *bind.FilterOpts) (*UniswapSetFeeProtocolIterator, error) {

	logs, sub, err := _Uniswap.contract.FilterLogs(opts, "SetFeeProtocol")
	if err != nil {
		return nil, err
	}
	return &UniswapSetFeeProtocolIterator{contract: _Uniswap.contract, event: "SetFeeProtocol", logs: logs, sub: sub}, nil
}

// WatchSetFeeProtocol is a free log subscription operation binding the contract event 0x973d8d92bb299f4af6ce49b52a8adb85ae46b9f214c4c4fc06ac77401237b133.
//
// Solidity: event SetFeeProtocol(uint8 feeProtocol0Old, uint8 feeProtocol1Old, uint8 feeProtocol0New, uint8 feeProtocol1New)
func (_Uniswap *UniswapFilterer) WatchSetFeeProtocol(opts *bind.WatchOpts, sink chan<- *UniswapSetFeeProtocol) (event.Subscription, error) {

	logs, sub, err := _Uniswap.contract.WatchLogs(opts, "SetFeeProtocol")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapSetFeeProtocol)
				if err := _Uniswap.contract.UnpackLog(event, "SetFeeProtocol", log); err != nil {
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

// ParseSetFeeProtocol is a log parse operation binding the contract event 0x973d8d92bb299f4af6ce49b52a8adb85ae46b9f214c4c4fc06ac77401237b133.
//
// Solidity: event SetFeeProtocol(uint8 feeProtocol0Old, uint8 feeProtocol1Old, uint8 feeProtocol0New, uint8 feeProtocol1New)
func (_Uniswap *UniswapFilterer) ParseSetFeeProtocol(log types.Log) (*UniswapSetFeeProtocol, error) {
	event := new(UniswapSetFeeProtocol)
	if err := _Uniswap.contract.UnpackLog(event, "SetFeeProtocol", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapSwapIterator is returned from FilterSwap and is used to iterate over the raw logs and unpacked data for Swap events raised by the Uniswap contract.
type UniswapSwapIterator struct {
	Event *UniswapSwap // Event containing the contract specifics and raw log

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
func (it *UniswapSwapIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapSwap)
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
		it.Event = new(UniswapSwap)
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
func (it *UniswapSwapIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapSwapIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapSwap represents a Swap event raised by the Uniswap contract.
type UniswapSwap struct {
	Sender       common.Address
	Recipient    common.Address
	Amount0      *big.Int
	Amount1      *big.Int
	SqrtPriceX96 *big.Int
	Liquidity    *big.Int
	Tick         *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSwap is a free log retrieval operation binding the contract event 0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67.
//
// Solidity: event Swap(address indexed sender, address indexed recipient, int256 amount0, int256 amount1, uint160 sqrtPriceX96, uint128 liquidity, int24 tick)
func (_Uniswap *UniswapFilterer) FilterSwap(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*UniswapSwapIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Uniswap.contract.FilterLogs(opts, "Swap", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &UniswapSwapIterator{contract: _Uniswap.contract, event: "Swap", logs: logs, sub: sub}, nil
}

// WatchSwap is a free log subscription operation binding the contract event 0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67.
//
// Solidity: event Swap(address indexed sender, address indexed recipient, int256 amount0, int256 amount1, uint160 sqrtPriceX96, uint128 liquidity, int24 tick)
func (_Uniswap *UniswapFilterer) WatchSwap(opts *bind.WatchOpts, sink chan<- *UniswapSwap, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Uniswap.contract.WatchLogs(opts, "Swap", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapSwap)
				if err := _Uniswap.contract.UnpackLog(event, "Swap", log); err != nil {
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

// ParseSwap is a log parse operation binding the contract event 0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67.
//
// Solidity: event Swap(address indexed sender, address indexed recipient, int256 amount0, int256 amount1, uint160 sqrtPriceX96, uint128 liquidity, int24 tick)
func (_Uniswap *UniswapFilterer) ParseSwap(log types.Log) (*UniswapSwap, error) {
	event := new(UniswapSwap)
	if err := _Uniswap.contract.UnpackLog(event, "Swap", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
