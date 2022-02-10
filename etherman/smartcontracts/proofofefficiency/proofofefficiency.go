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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIGlobalExitRootManager\",\"name\":\"_globalExitRootManager\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_matic\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"_rollupVerifier\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"genesisRoot\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sequencerAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"chainID\",\"type\":\"uint32\"}],\"name\":\"RegisterSequencer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"numBatch\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"batchChainID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"lastGlobalExitRoot\",\"type\":\"bytes32\"}],\"name\":\"SendBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"numBatch\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatch\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_CHAIN_ID\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateSequencerCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentLocalExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"globalExitRootManager\",\"outputs\":[{\"internalType\":\"contractIGlobalExitRootManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBatchSent\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastVerifiedBatch\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"matic\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numSequencers\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"}],\"name\":\"registerSequencer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupVerifier\",\"outputs\":[{\"internalType\":\"contractIVerifierRollup\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maticAmount\",\"type\":\"uint256\"}],\"name\":\"sendBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"name\":\"sentBatches\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"batchHashData\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maticCollateral\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"sequencers\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"internalType\":\"uint32\",\"name\":\"chainID\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"numBatch\",\"type\":\"uint32\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofA\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"proofB\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofC\",\"type\":\"uint256[2]\"}],\"name\":\"verifyBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200295f3803806200295f83398181016040528101906200003791906200032c565b620000576200004b6200011d60201b60201c565b6200012560201b60201c565b836004806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508273ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff168152505081600760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600581905550505050506200039e565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200021b82620001ee565b9050919050565b60006200022f826200020e565b9050919050565b620002418162000222565b81146200024d57600080fd5b50565b600081519050620002618162000236565b92915050565b600062000274826200020e565b9050919050565b620002868162000267565b81146200029257600080fd5b50565b600081519050620002a6816200027b565b92915050565b6000620002b9826200020e565b9050919050565b620002cb81620002ac565b8114620002d757600080fd5b50565b600081519050620002eb81620002c0565b92915050565b6000819050919050565b6200030681620002f1565b81146200031257600080fd5b50565b6000815190506200032681620002fb565b92915050565b60008060008060808587031215620003495762000348620001e9565b5b6000620003598782880162000250565b94505060206200036c8782880162000295565b93505060406200037f87828801620002da565b9250506060620003928782880162000315565b91505092959194509250565b608051612597620003c860003960008181610398015281816109c90152610de701526125976000f3fe608060405234801561001057600080fd5b50600436106101165760003560e01c80638da5cb5b116100a2578063ca98a30811610071578063ca98a3081461029a578063d02103ca146102b8578063e8bf92ed146102d6578063f2fde38b146102f4578063f51a97c01461031057610116565b80638da5cb5b14610222578063959c2f4714610240578063ac2eba981461025e578063b6b0b0971461027c57610116565b80633b880f77116100e95780633b880f77146101a457806343ea1996146101c0578063715018a6146101de5780637fcb3653146101e85780638a4abab81461020657610116565b806306d6490f1461011b578063188ea07d146101375780631c7a07ee146101555780631dc125b214610186575b600080fd5b61013560048036038101906101309190611600565b610341565b005b61013f610691565b60405161014c919061167b565b60405180910390f35b61016f600480360381019061016a91906116f4565b6106a7565b60405161017d9291906117a9565b60405180910390f35b61018e610763565b60405161019b91906117e8565b60405180910390f35b6101be60048036038101906101b991906118ae565b6107c9565b005b6101c8610a61565b6040516101d5919061167b565b60405180910390f35b6101e6610a67565b005b6101f0610aef565b6040516101fd919061167b565b60405180910390f35b610220600480360381019061021b91906119de565b610b05565b005b61022a610db0565b6040516102379190611a36565b60405180910390f35b610248610dd9565b6040516102559190611a60565b60405180910390f35b610266610ddf565b6040516102739190611a60565b60405180910390f35b610284610de5565b6040516102919190611ada565b60405180910390f35b6102a2610e09565b6040516102af919061167b565b60405180910390f35b6102c0610e1f565b6040516102cd9190611b16565b60405180910390f35b6102de610e43565b6040516102eb9190611b52565b60405180910390f35b61030e600480360381019061030991906116f4565b610e69565b005b61032a60048036038101906103259190611b6d565b610f61565b604051610338929190611b9a565b60405180910390f35b600061034b610763565b905081811115610390576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161038790611c35565b60405180910390fd5b6103dd3330837f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16610f85909392919063ffffffff16565b600060048054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633ed691ef6040518163ffffffff1660e01b815260040160206040518083038186803b15801561044557600080fd5b505afa158015610459573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061047d9190611c6a565b9050600080600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900463ffffffff1663ffffffff161461053a57600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900463ffffffff169050610540565b6103e890505b6002600481819054906101000a900463ffffffff168092919061056290611cc6565b91906101000a81548163ffffffff021916908363ffffffff16021790555050848242338460405160200161059a959493929190611dfa565b6040516020818303038152906040528051906020012060036000600260049054906101000a900463ffffffff1663ffffffff1663ffffffff168152602001908152602001600020600001819055508260036000600260049054906101000a900463ffffffff1663ffffffff1663ffffffff168152602001908152602001600020600101819055503373ffffffffffffffffffffffffffffffffffffffff16600260049054906101000a900463ffffffff1663ffffffff167f54195e3066d1cdf2d5295f7e34f49afff2a0bff1960f6a902536c3d88b163ab38385604051610682929190611e55565b60405180910390a35050505050565b600260049054906101000a900463ffffffff1681565b60016020528060005260406000206000915090508060000180546106ca90611ead565b80601f01602080910402602001604051908101604052809291908181526020018280546106f690611ead565b80156107435780601f1061071857610100808354040283529160200191610743565b820191906000526020600020905b81548152906001019060200180831161072657829003601f168201915b5050505050908060010160009054906101000a900463ffffffff16905082565b6000600460009054906101000a900463ffffffff16600260049054906101000a900463ffffffff1660016107979190611edf565b6107a19190611f19565b63ffffffff16670de0b6b3a76400006107ba9190611f61565b67ffffffffffffffff16905090565b6001600460009054906101000a900463ffffffff166107e89190611edf565b63ffffffff168463ffffffff1614610835576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161082c90612015565b60405180910390fd5b6000600360008663ffffffff1663ffffffff16815260200190815260200160002060405180604001604052908160008201548152602001600182015481525050905060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001600554600654898b86600001518b6040516020016108bd96959493929190612035565b6040516020818303038152906040528051906020012060001c6108e091906120d4565b90506004600081819054906101000a900463ffffffff168092919061090490611cc6565b91906101000a81548163ffffffff021916908363ffffffff16021790555050866005819055508760068190555060048054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166333d6247d6006546040518263ffffffff1660e01b815260040161098c9190611a60565b600060405180830381600087803b1580156109a657600080fd5b505af11580156109ba573d6000803e3d6000fd5b50505050610a0d3383602001517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1661100e9092919063ffffffff16565b3373ffffffffffffffffffffffffffffffffffffffff168663ffffffff167fb0a69d17322b203355f6c64abe7663174dd1d3a3a10ea0fa12bb4788d3f865d560405160405180910390a35050505050505050565b6103e881565b610a6f611094565b73ffffffffffffffffffffffffffffffffffffffff16610a8d610db0565b73ffffffffffffffffffffffffffffffffffffffff1614610ae3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ada90612151565b60405180910390fd5b610aed600061109c565b565b600460009054906101000a900463ffffffff1681565b600081511415610b4a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b41906121e3565b60405180910390fd5b6000600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900463ffffffff1663ffffffff161415610cc8576002600081819054906101000a900463ffffffff1680929190610bcd90611cc6565b91906101000a81548163ffffffff021916908363ffffffff1602179055505080600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000019080519060200190610c429291906113cd565b50600260009054906101000a900463ffffffff166103e8610c639190611edf565b600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160006101000a81548163ffffffff021916908363ffffffff160217905550610d20565b80600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000019080519060200190610d1e9291906113cd565b505b7fac2ab692920559c6528fa0189844a84d5889ac26fd6eb9f5bea5c7cd699a5ff13382600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900463ffffffff16604051610da593929190612203565b60405180910390a150565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60065481565b60055481565b7f000000000000000000000000000000000000000000000000000000000000000081565b600260009054906101000a900463ffffffff1681565b60048054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b610e71611094565b73ffffffffffffffffffffffffffffffffffffffff16610e8f610db0565b73ffffffffffffffffffffffffffffffffffffffff1614610ee5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610edc90612151565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610f55576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f4c906122b3565b60405180910390fd5b610f5e8161109c565b50565b60036020528060005260406000206000915090508060000154908060010154905082565b611008846323b872dd60e01b858585604051602401610fa6939291906122d3565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050611160565b50505050565b61108f8363a9059cbb60e01b848460405160240161102d92919061230a565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050611160565b505050565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b60006111c2826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166112279092919063ffffffff16565b905060008151111561122257808060200190518101906111e2919061236b565b611221576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112189061240a565b60405180910390fd5b5b505050565b6060611236848460008561123f565b90509392505050565b606082471015611284576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161127b9061249c565b60405180910390fd5b61128d85611353565b6112cc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112c390612508565b60405180910390fd5b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516112f59190612528565b60006040518083038185875af1925050503d8060008114611332576040519150601f19603f3d011682016040523d82523d6000602084013e611337565b606091505b5091509150611347828286611366565b92505050949350505050565b600080823b905060008111915050919050565b60608315611376578290506113c6565b6000835111156113895782518084602001fd5b816040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016113bd919061253f565b60405180910390fd5b9392505050565b8280546113d990611ead565b90600052602060002090601f0160209004810192826113fb5760008555611442565b82601f1061141457805160ff1916838001178555611442565b82800160010185558215611442579182015b82811115611441578251825591602001919060010190611426565b5b50905061144f9190611453565b5090565b5b8082111561146c576000816000905550600101611454565b5090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6114d78261148e565b810181811067ffffffffffffffff821117156114f6576114f561149f565b5b80604052505050565b6000611509611470565b905061151582826114ce565b919050565b600067ffffffffffffffff8211156115355761153461149f565b5b61153e8261148e565b9050602081019050919050565b82818337600083830152505050565b600061156d6115688461151a565b6114ff565b90508281526020810184848401111561158957611588611489565b5b61159484828561154b565b509392505050565b600082601f8301126115b1576115b0611484565b5b81356115c184826020860161155a565b91505092915050565b6000819050919050565b6115dd816115ca565b81146115e857600080fd5b50565b6000813590506115fa816115d4565b92915050565b600080604083850312156116175761161661147a565b5b600083013567ffffffffffffffff8111156116355761163461147f565b5b6116418582860161159c565b9250506020611652858286016115eb565b9150509250929050565b600063ffffffff82169050919050565b6116758161165c565b82525050565b6000602082019050611690600083018461166c565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006116c182611696565b9050919050565b6116d1816116b6565b81146116dc57600080fd5b50565b6000813590506116ee816116c8565b92915050565b60006020828403121561170a5761170961147a565b5b6000611718848285016116df565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b8381101561175b578082015181840152602081019050611740565b8381111561176a576000848401525b50505050565b600061177b82611721565b611785818561172c565b935061179581856020860161173d565b61179e8161148e565b840191505092915050565b600060408201905081810360008301526117c38185611770565b90506117d2602083018461166c565b9392505050565b6117e2816115ca565b82525050565b60006020820190506117fd60008301846117d9565b92915050565b6000819050919050565b61181681611803565b811461182157600080fd5b50565b6000813590506118338161180d565b92915050565b6118428161165c565b811461184d57600080fd5b50565b60008135905061185f81611839565b92915050565b600080fd5b60008190508260206002028201111561188657611885611865565b5b92915050565b6000819050826040600202820111156118a8576118a7611865565b5b92915050565b60008060008060008061016087890312156118cc576118cb61147a565b5b60006118da89828a01611824565b96505060206118eb89828a01611824565b95505060406118fc89828a01611850565b945050606061190d89828a0161186a565b93505060a061191e89828a0161188c565b92505061012061193089828a0161186a565b9150509295509295509295565b600067ffffffffffffffff8211156119585761195761149f565b5b6119618261148e565b9050602081019050919050565b600061198161197c8461193d565b6114ff565b90508281526020810184848401111561199d5761199c611489565b5b6119a884828561154b565b509392505050565b600082601f8301126119c5576119c4611484565b5b81356119d584826020860161196e565b91505092915050565b6000602082840312156119f4576119f361147a565b5b600082013567ffffffffffffffff811115611a1257611a1161147f565b5b611a1e848285016119b0565b91505092915050565b611a30816116b6565b82525050565b6000602082019050611a4b6000830184611a27565b92915050565b611a5a81611803565b82525050565b6000602082019050611a756000830184611a51565b92915050565b6000819050919050565b6000611aa0611a9b611a9684611696565b611a7b565b611696565b9050919050565b6000611ab282611a85565b9050919050565b6000611ac482611aa7565b9050919050565b611ad481611ab9565b82525050565b6000602082019050611aef6000830184611acb565b92915050565b6000611b0082611aa7565b9050919050565b611b1081611af5565b82525050565b6000602082019050611b2b6000830184611b07565b92915050565b6000611b3c82611aa7565b9050919050565b611b4c81611b31565b82525050565b6000602082019050611b676000830184611b43565b92915050565b600060208284031215611b8357611b8261147a565b5b6000611b9184828501611850565b91505092915050565b6000604082019050611baf6000830185611a51565b611bbc60208301846117d9565b9392505050565b7f50726f6f664f66456666696369656e63793a3a73656e6442617463683a204e4f60008201527f545f454e4f5547485f4d41544943000000000000000000000000000000000000602082015250565b6000611c1f602e8361172c565b9150611c2a82611bc3565b604082019050919050565b60006020820190508181036000830152611c4e81611c12565b9050919050565b600081519050611c648161180d565b92915050565b600060208284031215611c8057611c7f61147a565b5b6000611c8e84828501611c55565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611cd18261165c565b915063ffffffff821415611ce857611ce7611c97565b5b600182019050919050565b600081519050919050565b600081905092915050565b6000611d1482611cf3565b611d1e8185611cfe565b9350611d2e81856020860161173d565b80840191505092915050565b6000819050919050565b611d55611d5082611803565b611d3a565b82525050565b6000819050919050565b611d76611d71826115ca565b611d5b565b82525050565b60008160601b9050919050565b6000611d9482611d7c565b9050919050565b6000611da682611d89565b9050919050565b611dbe611db9826116b6565b611d9b565b82525050565b60008160e01b9050919050565b6000611ddc82611dc4565b9050919050565b611df4611def8261165c565b611dd1565b82525050565b6000611e068288611d09565b9150611e128287611d44565b602082019150611e228286611d65565b602082019150611e328285611dad565b601482019150611e428284611de3565b6004820191508190509695505050505050565b6000604082019050611e6a600083018561166c565b611e776020830184611a51565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680611ec557607f821691505b60208210811415611ed957611ed8611e7e565b5b50919050565b6000611eea8261165c565b9150611ef58361165c565b92508263ffffffff03821115611f0e57611f0d611c97565b5b828201905092915050565b6000611f248261165c565b9150611f2f8361165c565b925082821015611f4257611f41611c97565b5b828203905092915050565b600067ffffffffffffffff82169050919050565b6000611f6c82611f4d565b9150611f7783611f4d565b92508167ffffffffffffffff0483118215151615611f9857611f97611c97565b5b828202905092915050565b7f50726f6f664f66456666696369656e63793a3a76657269667942617463683a2060008201527f42415443485f444f45535f4e4f545f4d41544348000000000000000000000000602082015250565b6000611fff60348361172c565b915061200a82611fa3565b604082019050919050565b6000602082019050818103600083015261202e81611ff2565b9050919050565b60006120418289611d44565b6020820191506120518288611d44565b6020820191506120618287611d44565b6020820191506120718286611d44565b6020820191506120818285611d44565b6020820191506120918284611de3565b600482019150819050979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006120df826115ca565b91506120ea836115ca565b9250826120fa576120f96120a5565b5b828206905092915050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b600061213b60208361172c565b915061214682612105565b602082019050919050565b6000602082019050818103600083015261216a8161212e565b9050919050565b7f50726f6f664f66456666696369656e63793a3a7265676973746572536571756560008201527f6e6365723a204e4f545f56414c49445f55524c00000000000000000000000000602082015250565b60006121cd60338361172c565b91506121d882612171565b604082019050919050565b600060208201905081810360008301526121fc816121c0565b9050919050565b60006060820190506122186000830186611a27565b818103602083015261222a8185611770565b9050612239604083018461166c565b949350505050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b600061229d60268361172c565b91506122a882612241565b604082019050919050565b600060208201905081810360008301526122cc81612290565b9050919050565b60006060820190506122e86000830186611a27565b6122f56020830185611a27565b61230260408301846117d9565b949350505050565b600060408201905061231f6000830185611a27565b61232c60208301846117d9565b9392505050565b60008115159050919050565b61234881612333565b811461235357600080fd5b50565b6000815190506123658161233f565b92915050565b6000602082840312156123815761238061147a565b5b600061238f84828501612356565b91505092915050565b7f5361666545524332303a204552433230206f7065726174696f6e20646964206e60008201527f6f74207375636365656400000000000000000000000000000000000000000000602082015250565b60006123f4602a8361172c565b91506123ff82612398565b604082019050919050565b60006020820190508181036000830152612423816123e7565b9050919050565b7f416464726573733a20696e73756666696369656e742062616c616e636520666f60008201527f722063616c6c0000000000000000000000000000000000000000000000000000602082015250565b600061248660268361172c565b91506124918261242a565b604082019050919050565b600060208201905081810360008301526124b581612479565b9050919050565b7f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000600082015250565b60006124f2601d8361172c565b91506124fd826124bc565b602082019050919050565b60006020820190508181036000830152612521816124e5565b9050919050565b60006125348284611d09565b915081905092915050565b600060208201905081810360008301526125598184611770565b90509291505056fea26469706673582212202938b3254aa843f022a33dde08d25bbd6d42e29154cdc3d5316c18210d20c2e064736f6c63430008090033",
}

// ProofofefficiencyABI is the input ABI used to generate the binding from.
// Deprecated: Use ProofofefficiencyMetaData.ABI instead.
var ProofofefficiencyABI = ProofofefficiencyMetaData.ABI

// ProofofefficiencyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ProofofefficiencyMetaData.Bin instead.
var ProofofefficiencyBin = ProofofefficiencyMetaData.Bin

// DeployProofofefficiency deploys a new Ethereum contract, binding an instance of Proofofefficiency to it.
func DeployProofofefficiency(auth *bind.TransactOpts, backend bind.ContractBackend, _globalExitRootManager common.Address, _matic common.Address, _rollupVerifier common.Address, genesisRoot [32]byte) (common.Address, *types.Transaction, *Proofofefficiency, error) {
	parsed, err := ProofofefficiencyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ProofofefficiencyBin), backend, _globalExitRootManager, _matic, _rollupVerifier, genesisRoot)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Proofofefficiency{ProofofefficiencyCaller: ProofofefficiencyCaller{contract: contract}, ProofofefficiencyTransactor: ProofofefficiencyTransactor{contract: contract}, ProofofefficiencyFilterer: ProofofefficiencyFilterer{contract: contract}}, nil
}

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

// DEFAULTCHAINID is a free data retrieval call binding the contract method 0x43ea1996.
//
// Solidity: function DEFAULT_CHAIN_ID() view returns(uint32)
func (_Proofofefficiency *ProofofefficiencyCaller) DEFAULTCHAINID(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "DEFAULT_CHAIN_ID")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// DEFAULTCHAINID is a free data retrieval call binding the contract method 0x43ea1996.
//
// Solidity: function DEFAULT_CHAIN_ID() view returns(uint32)
func (_Proofofefficiency *ProofofefficiencySession) DEFAULTCHAINID() (uint32, error) {
	return _Proofofefficiency.Contract.DEFAULTCHAINID(&_Proofofefficiency.CallOpts)
}

// DEFAULTCHAINID is a free data retrieval call binding the contract method 0x43ea1996.
//
// Solidity: function DEFAULT_CHAIN_ID() view returns(uint32)
func (_Proofofefficiency *ProofofefficiencyCallerSession) DEFAULTCHAINID() (uint32, error) {
	return _Proofofefficiency.Contract.DEFAULTCHAINID(&_Proofofefficiency.CallOpts)
}

// CalculateSequencerCollateral is a free data retrieval call binding the contract method 0x1dc125b2.
//
// Solidity: function calculateSequencerCollateral() view returns(uint256)
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
// Solidity: function calculateSequencerCollateral() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencySession) CalculateSequencerCollateral() (*big.Int, error) {
	return _Proofofefficiency.Contract.CalculateSequencerCollateral(&_Proofofefficiency.CallOpts)
}

// CalculateSequencerCollateral is a free data retrieval call binding the contract method 0x1dc125b2.
//
// Solidity: function calculateSequencerCollateral() view returns(uint256)
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

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCaller) GlobalExitRootManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "globalExitRootManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Proofofefficiency *ProofofefficiencySession) GlobalExitRootManager() (common.Address, error) {
	return _Proofofefficiency.Contract.GlobalExitRootManager(&_Proofofefficiency.CallOpts)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCallerSession) GlobalExitRootManager() (common.Address, error) {
	return _Proofofefficiency.Contract.GlobalExitRootManager(&_Proofofefficiency.CallOpts)
}

// LastBatchSent is a free data retrieval call binding the contract method 0x188ea07d.
//
// Solidity: function lastBatchSent() view returns(uint32)
func (_Proofofefficiency *ProofofefficiencyCaller) LastBatchSent(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "lastBatchSent")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// LastBatchSent is a free data retrieval call binding the contract method 0x188ea07d.
//
// Solidity: function lastBatchSent() view returns(uint32)
func (_Proofofefficiency *ProofofefficiencySession) LastBatchSent() (uint32, error) {
	return _Proofofefficiency.Contract.LastBatchSent(&_Proofofefficiency.CallOpts)
}

// LastBatchSent is a free data retrieval call binding the contract method 0x188ea07d.
//
// Solidity: function lastBatchSent() view returns(uint32)
func (_Proofofefficiency *ProofofefficiencyCallerSession) LastBatchSent() (uint32, error) {
	return _Proofofefficiency.Contract.LastBatchSent(&_Proofofefficiency.CallOpts)
}

// LastVerifiedBatch is a free data retrieval call binding the contract method 0x7fcb3653.
//
// Solidity: function lastVerifiedBatch() view returns(uint32)
func (_Proofofefficiency *ProofofefficiencyCaller) LastVerifiedBatch(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "lastVerifiedBatch")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// LastVerifiedBatch is a free data retrieval call binding the contract method 0x7fcb3653.
//
// Solidity: function lastVerifiedBatch() view returns(uint32)
func (_Proofofefficiency *ProofofefficiencySession) LastVerifiedBatch() (uint32, error) {
	return _Proofofefficiency.Contract.LastVerifiedBatch(&_Proofofefficiency.CallOpts)
}

// LastVerifiedBatch is a free data retrieval call binding the contract method 0x7fcb3653.
//
// Solidity: function lastVerifiedBatch() view returns(uint32)
func (_Proofofefficiency *ProofofefficiencyCallerSession) LastVerifiedBatch() (uint32, error) {
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
// Solidity: function numSequencers() view returns(uint32)
func (_Proofofefficiency *ProofofefficiencyCaller) NumSequencers(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "numSequencers")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// NumSequencers is a free data retrieval call binding the contract method 0xca98a308.
//
// Solidity: function numSequencers() view returns(uint32)
func (_Proofofefficiency *ProofofefficiencySession) NumSequencers() (uint32, error) {
	return _Proofofefficiency.Contract.NumSequencers(&_Proofofefficiency.CallOpts)
}

// NumSequencers is a free data retrieval call binding the contract method 0xca98a308.
//
// Solidity: function numSequencers() view returns(uint32)
func (_Proofofefficiency *ProofofefficiencyCallerSession) NumSequencers() (uint32, error) {
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

// SentBatches is a free data retrieval call binding the contract method 0xf51a97c0.
//
// Solidity: function sentBatches(uint32 ) view returns(bytes32 batchHashData, uint256 maticCollateral)
func (_Proofofefficiency *ProofofefficiencyCaller) SentBatches(opts *bind.CallOpts, arg0 uint32) (struct {
	BatchHashData   [32]byte
	MaticCollateral *big.Int
}, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "sentBatches", arg0)

	outstruct := new(struct {
		BatchHashData   [32]byte
		MaticCollateral *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BatchHashData = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.MaticCollateral = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// SentBatches is a free data retrieval call binding the contract method 0xf51a97c0.
//
// Solidity: function sentBatches(uint32 ) view returns(bytes32 batchHashData, uint256 maticCollateral)
func (_Proofofefficiency *ProofofefficiencySession) SentBatches(arg0 uint32) (struct {
	BatchHashData   [32]byte
	MaticCollateral *big.Int
}, error) {
	return _Proofofefficiency.Contract.SentBatches(&_Proofofefficiency.CallOpts, arg0)
}

// SentBatches is a free data retrieval call binding the contract method 0xf51a97c0.
//
// Solidity: function sentBatches(uint32 ) view returns(bytes32 batchHashData, uint256 maticCollateral)
func (_Proofofefficiency *ProofofefficiencyCallerSession) SentBatches(arg0 uint32) (struct {
	BatchHashData   [32]byte
	MaticCollateral *big.Int
}, error) {
	return _Proofofefficiency.Contract.SentBatches(&_Proofofefficiency.CallOpts, arg0)
}

// Sequencers is a free data retrieval call binding the contract method 0x1c7a07ee.
//
// Solidity: function sequencers(address ) view returns(string sequencerURL, uint32 chainID)
func (_Proofofefficiency *ProofofefficiencyCaller) Sequencers(opts *bind.CallOpts, arg0 common.Address) (struct {
	SequencerURL string
	ChainID      uint32
}, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "sequencers", arg0)

	outstruct := new(struct {
		SequencerURL string
		ChainID      uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SequencerURL = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.ChainID = *abi.ConvertType(out[1], new(uint32)).(*uint32)

	return *outstruct, err

}

// Sequencers is a free data retrieval call binding the contract method 0x1c7a07ee.
//
// Solidity: function sequencers(address ) view returns(string sequencerURL, uint32 chainID)
func (_Proofofefficiency *ProofofefficiencySession) Sequencers(arg0 common.Address) (struct {
	SequencerURL string
	ChainID      uint32
}, error) {
	return _Proofofefficiency.Contract.Sequencers(&_Proofofefficiency.CallOpts, arg0)
}

// Sequencers is a free data retrieval call binding the contract method 0x1c7a07ee.
//
// Solidity: function sequencers(address ) view returns(string sequencerURL, uint32 chainID)
func (_Proofofefficiency *ProofofefficiencyCallerSession) Sequencers(arg0 common.Address) (struct {
	SequencerURL string
	ChainID      uint32
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

// VerifyBatch is a paid mutator transaction binding the contract method 0x3b880f77.
//
// Solidity: function verifyBatch(bytes32 newLocalExitRoot, bytes32 newStateRoot, uint32 numBatch, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) VerifyBatch(opts *bind.TransactOpts, newLocalExitRoot [32]byte, newStateRoot [32]byte, numBatch uint32, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "verifyBatch", newLocalExitRoot, newStateRoot, numBatch, proofA, proofB, proofC)
}

// VerifyBatch is a paid mutator transaction binding the contract method 0x3b880f77.
//
// Solidity: function verifyBatch(bytes32 newLocalExitRoot, bytes32 newStateRoot, uint32 numBatch, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencySession) VerifyBatch(newLocalExitRoot [32]byte, newStateRoot [32]byte, numBatch uint32, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.VerifyBatch(&_Proofofefficiency.TransactOpts, newLocalExitRoot, newStateRoot, numBatch, proofA, proofB, proofC)
}

// VerifyBatch is a paid mutator transaction binding the contract method 0x3b880f77.
//
// Solidity: function verifyBatch(bytes32 newLocalExitRoot, bytes32 newStateRoot, uint32 numBatch, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) VerifyBatch(newLocalExitRoot [32]byte, newStateRoot [32]byte, numBatch uint32, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.VerifyBatch(&_Proofofefficiency.TransactOpts, newLocalExitRoot, newStateRoot, numBatch, proofA, proofB, proofC)
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

// ProofofefficiencyRegisterSequencerIterator is returned from FilterRegisterSequencer and is used to iterate over the raw logs and unpacked data for RegisterSequencer events raised by the Proofofefficiency contract.
type ProofofefficiencyRegisterSequencerIterator struct {
	Event *ProofofefficiencyRegisterSequencer // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencyRegisterSequencerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencyRegisterSequencer)
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
		it.Event = new(ProofofefficiencyRegisterSequencer)
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
func (it *ProofofefficiencyRegisterSequencerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencyRegisterSequencerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencyRegisterSequencer represents a RegisterSequencer event raised by the Proofofefficiency contract.
type ProofofefficiencyRegisterSequencer struct {
	SequencerAddress common.Address
	SequencerURL     string
	ChainID          uint32
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRegisterSequencer is a free log retrieval operation binding the contract event 0xac2ab692920559c6528fa0189844a84d5889ac26fd6eb9f5bea5c7cd699a5ff1.
//
// Solidity: event RegisterSequencer(address sequencerAddress, string sequencerURL, uint32 chainID)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterRegisterSequencer(opts *bind.FilterOpts) (*ProofofefficiencyRegisterSequencerIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "RegisterSequencer")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyRegisterSequencerIterator{contract: _Proofofefficiency.contract, event: "RegisterSequencer", logs: logs, sub: sub}, nil
}

// WatchRegisterSequencer is a free log subscription operation binding the contract event 0xac2ab692920559c6528fa0189844a84d5889ac26fd6eb9f5bea5c7cd699a5ff1.
//
// Solidity: event RegisterSequencer(address sequencerAddress, string sequencerURL, uint32 chainID)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchRegisterSequencer(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyRegisterSequencer) (event.Subscription, error) {

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "RegisterSequencer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencyRegisterSequencer)
				if err := _Proofofefficiency.contract.UnpackLog(event, "RegisterSequencer", log); err != nil {
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

// ParseRegisterSequencer is a log parse operation binding the contract event 0xac2ab692920559c6528fa0189844a84d5889ac26fd6eb9f5bea5c7cd699a5ff1.
//
// Solidity: event RegisterSequencer(address sequencerAddress, string sequencerURL, uint32 chainID)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseRegisterSequencer(log types.Log) (*ProofofefficiencyRegisterSequencer, error) {
	event := new(ProofofefficiencyRegisterSequencer)
	if err := _Proofofefficiency.contract.UnpackLog(event, "RegisterSequencer", log); err != nil {
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
	NumBatch           uint32
	Sequencer          common.Address
	BatchChainID       uint32
	LastGlobalExitRoot [32]byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterSendBatch is a free log retrieval operation binding the contract event 0x54195e3066d1cdf2d5295f7e34f49afff2a0bff1960f6a902536c3d88b163ab3.
//
// Solidity: event SendBatch(uint32 indexed numBatch, address indexed sequencer, uint32 batchChainID, bytes32 lastGlobalExitRoot)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSendBatch(opts *bind.FilterOpts, numBatch []uint32, sequencer []common.Address) (*ProofofefficiencySendBatchIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}
	var sequencerRule []interface{}
	for _, sequencerItem := range sequencer {
		sequencerRule = append(sequencerRule, sequencerItem)
	}

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "SendBatch", numBatchRule, sequencerRule)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencySendBatchIterator{contract: _Proofofefficiency.contract, event: "SendBatch", logs: logs, sub: sub}, nil
}

// WatchSendBatch is a free log subscription operation binding the contract event 0x54195e3066d1cdf2d5295f7e34f49afff2a0bff1960f6a902536c3d88b163ab3.
//
// Solidity: event SendBatch(uint32 indexed numBatch, address indexed sequencer, uint32 batchChainID, bytes32 lastGlobalExitRoot)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSendBatch(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySendBatch, numBatch []uint32, sequencer []common.Address) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}
	var sequencerRule []interface{}
	for _, sequencerItem := range sequencer {
		sequencerRule = append(sequencerRule, sequencerItem)
	}

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "SendBatch", numBatchRule, sequencerRule)
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

// ParseSendBatch is a log parse operation binding the contract event 0x54195e3066d1cdf2d5295f7e34f49afff2a0bff1960f6a902536c3d88b163ab3.
//
// Solidity: event SendBatch(uint32 indexed numBatch, address indexed sequencer, uint32 batchChainID, bytes32 lastGlobalExitRoot)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseSendBatch(log types.Log) (*ProofofefficiencySendBatch, error) {
	event := new(ProofofefficiencySendBatch)
	if err := _Proofofefficiency.contract.UnpackLog(event, "SendBatch", log); err != nil {
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
	NumBatch   uint32
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVerifyBatch is a free log retrieval operation binding the contract event 0xb0a69d17322b203355f6c64abe7663174dd1d3a3a10ea0fa12bb4788d3f865d5.
//
// Solidity: event VerifyBatch(uint32 indexed numBatch, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterVerifyBatch(opts *bind.FilterOpts, numBatch []uint32, aggregator []common.Address) (*ProofofefficiencyVerifyBatchIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}
	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "VerifyBatch", numBatchRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyVerifyBatchIterator{contract: _Proofofefficiency.contract, event: "VerifyBatch", logs: logs, sub: sub}, nil
}

// WatchVerifyBatch is a free log subscription operation binding the contract event 0xb0a69d17322b203355f6c64abe7663174dd1d3a3a10ea0fa12bb4788d3f865d5.
//
// Solidity: event VerifyBatch(uint32 indexed numBatch, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchVerifyBatch(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyVerifyBatch, numBatch []uint32, aggregator []common.Address) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}
	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "VerifyBatch", numBatchRule, aggregatorRule)
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

// ParseVerifyBatch is a log parse operation binding the contract event 0xb0a69d17322b203355f6c64abe7663174dd1d3a3a10ea0fa12bb4788d3f865d5.
//
// Solidity: event VerifyBatch(uint32 indexed numBatch, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseVerifyBatch(log types.Log) (*ProofofefficiencyVerifyBatch, error) {
	event := new(ProofofefficiencyVerifyBatch)
	if err := _Proofofefficiency.contract.UnpackLog(event, "VerifyBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
