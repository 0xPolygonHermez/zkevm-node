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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractBridgeInterface\",\"name\":\"_bridge\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_matic\",\"type\":\"address\"},{\"internalType\":\"contractVerifierRollupInterface\",\"name\":\"_rollupVerifier\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"genesisRoot\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sequencerAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"chainID\",\"type\":\"uint32\"}],\"name\":\"RegisterSequencer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"batchNum\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"}],\"name\":\"SendBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"batchNum\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatch\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_CHAIN_ID\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridge\",\"outputs\":[{\"internalType\":\"contractBridgeInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateSequencerCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentLocalExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBatchSent\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastVerifiedBatch\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"matic\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numSequencers\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"}],\"name\":\"registerSequencer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupVerifier\",\"outputs\":[{\"internalType\":\"contractVerifierRollupInterface\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maticAmount\",\"type\":\"uint256\"}],\"name\":\"sendBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"name\":\"sentBatches\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"sequencerAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"chainID\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"batchHashData\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maticCollateral\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"sequencers\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"internalType\":\"uint32\",\"name\":\"chainID\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"batchNum\",\"type\":\"uint32\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofA\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"proofB\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofC\",\"type\":\"uint256[2]\"}],\"name\":\"verifyBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001fae38038062001fae83398181016040528101906200003791906200032c565b620000576200004b6200011d60201b60201c565b6200012560201b60201c565b836004806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508273ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff168152505081600760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600581905550505050506200039e565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200021b82620001ee565b9050919050565b60006200022f826200020e565b9050919050565b620002418162000222565b81146200024d57600080fd5b50565b600081519050620002618162000236565b92915050565b600062000274826200020e565b9050919050565b620002868162000267565b81146200029257600080fd5b50565b600081519050620002a6816200027b565b92915050565b6000620002b9826200020e565b9050919050565b620002cb81620002ac565b8114620002d757600080fd5b50565b600081519050620002eb81620002c0565b92915050565b6000819050919050565b6200030681620002f1565b81146200031257600080fd5b50565b6000815190506200032681620002fb565b92915050565b60008060008060808587031215620003495762000348620001e9565b5b6000620003598782880162000250565b94505060206200036c8782880162000295565b93505060406200037f87828801620002da565b9250506060620003928782880162000315565b91505092959194509250565b608051611bf4620003ba6000396000610c7e0152611bf46000f3fe608060405234801561001057600080fd5b50600436106101165760003560e01c80638da5cb5b116100a2578063ca98a30811610071578063ca98a3081461029a578063e78cea92146102b8578063e8bf92ed146102d6578063f2fde38b146102f4578063f51a97c01461031057610116565b80638da5cb5b14610222578063959c2f4714610240578063ac2eba981461025e578063b6b0b0971461027c57610116565b80633b880f77116100e95780633b880f77146101a457806343ea1996146101c0578063715018a6146101de5780637fcb3653146101e85780638a4abab81461020657610116565b806306d6490f1461011b578063188ea07d146101375780631c7a07ee146101555780631dc125b214610186575b600080fd5b61013560048036038101906101309190611157565b610343565b005b61013f6105fd565b60405161014c91906111d2565b60405180910390f35b61016f600480360381019061016a919061124b565b610613565b60405161017d929190611300565b60405180910390f35b61018e6106cf565b60405161019b919061133f565b60405180910390f35b6101be60048036038101906101b99190611405565b6106df565b005b6101c86108f8565b6040516101d591906111d2565b60405180910390f35b6101e66108fe565b005b6101f0610986565b6040516101fd91906111d2565b60405180910390f35b610220600480360381019061021b9190611535565b61099c565b005b61022a610c47565b604051610237919061158d565b60405180910390f35b610248610c70565b60405161025591906115b7565b60405180910390f35b610266610c76565b60405161027391906115b7565b60405180910390f35b610284610c7c565b6040516102919190611631565b60405180910390f35b6102a2610ca0565b6040516102af91906111d2565b60405180910390f35b6102c0610cb6565b6040516102cd919061166d565b60405180910390f35b6102de610cda565b6040516102eb91906116a9565b60405180910390f35b61030e6004803603810190610309919061124b565b610d00565b005b61032a600480360381019061032591906116c4565b610df8565b60405161033a94939291906116f1565b60405180910390f35b600061034d6106cf565b90506002600481819054906101000a900463ffffffff168092919061037190611765565b91906101000a81548163ffffffff021916908363ffffffff160217905550508060036000600260049054906101000a900463ffffffff1663ffffffff1663ffffffff168152602001908152602001600020600201819055503360036000600260049054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900463ffffffff1663ffffffff161461054657600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900463ffffffff1660036000600260049054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060000160146101000a81548163ffffffff021916908363ffffffff16021790555061059c565b6103e860036000600260049054906101000a900463ffffffff1663ffffffff1663ffffffff16815260200190815260200160002060000160146101000a81548163ffffffff021916908363ffffffff1602179055505b3373ffffffffffffffffffffffffffffffffffffffff16600260049054906101000a900463ffffffff1663ffffffff167facbc4bfeebf66c85e73298070a2b645714fa69a7ef8890ed3bd6e9c581f7cdaf60405160405180910390a3505050565b600260049054906101000a900463ffffffff1681565b6001602052806000526040600020600091509050806000018054610636906117c1565b80601f0160208091040260200160405190810160405280929190818152602001828054610662906117c1565b80156106af5780601f10610684576101008083540402835291602001916106af565b820191906000526020600020905b81548152906001019060200180831161069257829003601f168201915b5050505050908060010160009054906101000a900463ffffffff16905082565b6000670de0b6b3a7640000905090565b6001600460009054906101000a900463ffffffff166106fe91906117f3565b63ffffffff168463ffffffff161461074b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107429061189f565b60405180910390fd5b6000600360008663ffffffff1663ffffffff1681526020019081526020016000206040518060800160405290816000820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016000820160149054906101000a900463ffffffff1663ffffffff1663ffffffff1681526020016001820154815260200160028201548152505090506000600554600654888a8560000151866040015187602001518c60405160200161083a98979695949392919061195e565b6040516020818303038152906040528051906020012060001c90506004600081819054906101000a900463ffffffff168092919061087790611765565b91906101000a81548163ffffffff021916908363ffffffff1602179055505086600581905550876006819055503373ffffffffffffffffffffffffffffffffffffffff168663ffffffff167fb0a69d17322b203355f6c64abe7663174dd1d3a3a10ea0fa12bb4788d3f865d560405160405180910390a35050505050505050565b6103e881565b610906610e58565b73ffffffffffffffffffffffffffffffffffffffff16610924610c47565b73ffffffffffffffffffffffffffffffffffffffff161461097a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161097190611a3c565b60405180910390fd5b6109846000610e60565b565b600460009054906101000a900463ffffffff1681565b6000815114156109e1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109d890611ace565b60405180910390fd5b6000600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900463ffffffff1663ffffffff161415610b5f576002600081819054906101000a900463ffffffff1680929190610a6490611765565b91906101000a81548163ffffffff021916908363ffffffff1602179055505080600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000019080519060200190610ad9929190610f24565b50600260009054906101000a900463ffffffff166103e8610afa91906117f3565b600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160006101000a81548163ffffffff021916908363ffffffff160217905550610bb7565b80600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000019080519060200190610bb5929190610f24565b505b7fac2ab692920559c6528fa0189844a84d5889ac26fd6eb9f5bea5c7cd699a5ff13382600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900463ffffffff16604051610c3c93929190611aee565b60405180910390a150565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60065481565b60055481565b7f000000000000000000000000000000000000000000000000000000000000000081565b600260009054906101000a900463ffffffff1681565b60048054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b610d08610e58565b73ffffffffffffffffffffffffffffffffffffffff16610d26610c47565b73ffffffffffffffffffffffffffffffffffffffff1614610d7c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d7390611a3c565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610dec576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610de390611b9e565b60405180910390fd5b610df581610e60565b50565b60036020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060000160149054906101000a900463ffffffff16908060010154908060020154905084565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b828054610f30906117c1565b90600052602060002090601f016020900481019282610f525760008555610f99565b82601f10610f6b57805160ff1916838001178555610f99565b82800160010185558215610f99579182015b82811115610f98578251825591602001919060010190610f7d565b5b509050610fa69190610faa565b5090565b5b80821115610fc3576000816000905550600101610fab565b5090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61102e82610fe5565b810181811067ffffffffffffffff8211171561104d5761104c610ff6565b5b80604052505050565b6000611060610fc7565b905061106c8282611025565b919050565b600067ffffffffffffffff82111561108c5761108b610ff6565b5b61109582610fe5565b9050602081019050919050565b82818337600083830152505050565b60006110c46110bf84611071565b611056565b9050828152602081018484840111156110e0576110df610fe0565b5b6110eb8482856110a2565b509392505050565b600082601f83011261110857611107610fdb565b5b81356111188482602086016110b1565b91505092915050565b6000819050919050565b61113481611121565b811461113f57600080fd5b50565b6000813590506111518161112b565b92915050565b6000806040838503121561116e5761116d610fd1565b5b600083013567ffffffffffffffff81111561118c5761118b610fd6565b5b611198858286016110f3565b92505060206111a985828601611142565b9150509250929050565b600063ffffffff82169050919050565b6111cc816111b3565b82525050565b60006020820190506111e760008301846111c3565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000611218826111ed565b9050919050565b6112288161120d565b811461123357600080fd5b50565b6000813590506112458161121f565b92915050565b60006020828403121561126157611260610fd1565b5b600061126f84828501611236565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b838110156112b2578082015181840152602081019050611297565b838111156112c1576000848401525b50505050565b60006112d282611278565b6112dc8185611283565b93506112ec818560208601611294565b6112f581610fe5565b840191505092915050565b6000604082019050818103600083015261131a81856112c7565b905061132960208301846111c3565b9392505050565b61133981611121565b82525050565b60006020820190506113546000830184611330565b92915050565b6000819050919050565b61136d8161135a565b811461137857600080fd5b50565b60008135905061138a81611364565b92915050565b611399816111b3565b81146113a457600080fd5b50565b6000813590506113b681611390565b92915050565b600080fd5b6000819050826020600202820111156113dd576113dc6113bc565b5b92915050565b6000819050826040600202820111156113ff576113fe6113bc565b5b92915050565b600080600080600080610160878903121561142357611422610fd1565b5b600061143189828a0161137b565b965050602061144289828a0161137b565b955050604061145389828a016113a7565b945050606061146489828a016113c1565b93505060a061147589828a016113e3565b92505061012061148789828a016113c1565b9150509295509295509295565b600067ffffffffffffffff8211156114af576114ae610ff6565b5b6114b882610fe5565b9050602081019050919050565b60006114d86114d384611494565b611056565b9050828152602081018484840111156114f4576114f3610fe0565b5b6114ff8482856110a2565b509392505050565b600082601f83011261151c5761151b610fdb565b5b813561152c8482602086016114c5565b91505092915050565b60006020828403121561154b5761154a610fd1565b5b600082013567ffffffffffffffff81111561156957611568610fd6565b5b61157584828501611507565b91505092915050565b6115878161120d565b82525050565b60006020820190506115a2600083018461157e565b92915050565b6115b18161135a565b82525050565b60006020820190506115cc60008301846115a8565b92915050565b6000819050919050565b60006115f76115f26115ed846111ed565b6115d2565b6111ed565b9050919050565b6000611609826115dc565b9050919050565b600061161b826115fe565b9050919050565b61162b81611610565b82525050565b60006020820190506116466000830184611622565b92915050565b6000611657826115fe565b9050919050565b6116678161164c565b82525050565b6000602082019050611682600083018461165e565b92915050565b6000611693826115fe565b9050919050565b6116a381611688565b82525050565b60006020820190506116be600083018461169a565b92915050565b6000602082840312156116da576116d9610fd1565b5b60006116e8848285016113a7565b91505092915050565b6000608082019050611706600083018761157e565b61171360208301866111c3565b61172060408301856115a8565b61172d6060830184611330565b95945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611770826111b3565b915063ffffffff82141561178757611786611736565b5b600182019050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806117d957607f821691505b602082108114156117ed576117ec611792565b5b50919050565b60006117fe826111b3565b9150611809836111b3565b92508263ffffffff0382111561182257611821611736565b5b828201905092915050565b7f50726f6f664f66456666696369656e63793a3a76657269667942617463683a2060008201527f42415443485f444f45535f4e4f545f4d41544348000000000000000000000000602082015250565b6000611889603483611283565b91506118948261182d565b604082019050919050565b600060208201905081810360008301526118b88161187c565b9050919050565b6000819050919050565b6118da6118d58261135a565b6118bf565b82525050565b60008160601b9050919050565b60006118f8826118e0565b9050919050565b600061190a826118ed565b9050919050565b61192261191d8261120d565b6118ff565b82525050565b60008160e01b9050919050565b600061194082611928565b9050919050565b611958611953826111b3565b611935565b82525050565b600061196a828b6118c9565b60208201915061197a828a6118c9565b60208201915061198a82896118c9565b60208201915061199a82886118c9565b6020820191506119aa8287611911565b6014820191506119ba82866118c9565b6020820191506119ca8285611947565b6004820191506119da8284611947565b6004820191508190509998505050505050505050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b6000611a26602083611283565b9150611a31826119f0565b602082019050919050565b60006020820190508181036000830152611a5581611a19565b9050919050565b7f50726f6f664f66456666696369656e63793a3a7265676973746572536571756560008201527f6e6365723a204e4f545f56414c49445f55524c00000000000000000000000000602082015250565b6000611ab8603383611283565b9150611ac382611a5c565b604082019050919050565b60006020820190508181036000830152611ae781611aab565b9050919050565b6000606082019050611b03600083018661157e565b8181036020830152611b1581856112c7565b9050611b2460408301846111c3565b949350505050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b6000611b88602683611283565b9150611b9382611b2c565b604082019050919050565b60006020820190508181036000830152611bb781611b7b565b905091905056fea2646970667358221220f1940d4a1eb52343482348362f1f364b6a4de22eaf3fa78c530eaaeb8045194664736f6c63430008090033",
}

// ProofofefficiencyABI is the input ABI used to generate the binding from.
// Deprecated: Use ProofofefficiencyMetaData.ABI instead.
var ProofofefficiencyABI = ProofofefficiencyMetaData.ABI

// ProofofefficiencyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ProofofefficiencyMetaData.Bin instead.
var ProofofefficiencyBin = ProofofefficiencyMetaData.Bin

// DeployProofofefficiency deploys a new Ethereum contract, binding an instance of Proofofefficiency to it.
func DeployProofofefficiency(auth *bind.TransactOpts, backend bind.ContractBackend, _bridge common.Address, _matic common.Address, _rollupVerifier common.Address, genesisRoot [32]byte) (common.Address, *types.Transaction, *Proofofefficiency, error) {
	parsed, err := ProofofefficiencyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ProofofefficiencyBin), backend, _bridge, _matic, _rollupVerifier, genesisRoot)
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
// Solidity: function sentBatches(uint32 ) view returns(address sequencerAddress, uint32 chainID, bytes32 batchHashData, uint256 maticCollateral)
func (_Proofofefficiency *ProofofefficiencyCaller) SentBatches(opts *bind.CallOpts, arg0 uint32) (struct {
	SequencerAddress common.Address
	ChainID          uint32
	BatchHashData    [32]byte
	MaticCollateral  *big.Int
}, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "sentBatches", arg0)

	outstruct := new(struct {
		SequencerAddress common.Address
		ChainID          uint32
		BatchHashData    [32]byte
		MaticCollateral  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SequencerAddress = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ChainID = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.BatchHashData = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.MaticCollateral = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// SentBatches is a free data retrieval call binding the contract method 0xf51a97c0.
//
// Solidity: function sentBatches(uint32 ) view returns(address sequencerAddress, uint32 chainID, bytes32 batchHashData, uint256 maticCollateral)
func (_Proofofefficiency *ProofofefficiencySession) SentBatches(arg0 uint32) (struct {
	SequencerAddress common.Address
	ChainID          uint32
	BatchHashData    [32]byte
	MaticCollateral  *big.Int
}, error) {
	return _Proofofefficiency.Contract.SentBatches(&_Proofofefficiency.CallOpts, arg0)
}

// SentBatches is a free data retrieval call binding the contract method 0xf51a97c0.
//
// Solidity: function sentBatches(uint32 ) view returns(address sequencerAddress, uint32 chainID, bytes32 batchHashData, uint256 maticCollateral)
func (_Proofofefficiency *ProofofefficiencyCallerSession) SentBatches(arg0 uint32) (struct {
	SequencerAddress common.Address
	ChainID          uint32
	BatchHashData    [32]byte
	MaticCollateral  *big.Int
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
// Solidity: function verifyBatch(bytes32 newLocalExitRoot, bytes32 newStateRoot, uint32 batchNum, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) VerifyBatch(opts *bind.TransactOpts, newLocalExitRoot [32]byte, newStateRoot [32]byte, batchNum uint32, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "verifyBatch", newLocalExitRoot, newStateRoot, batchNum, proofA, proofB, proofC)
}

// VerifyBatch is a paid mutator transaction binding the contract method 0x3b880f77.
//
// Solidity: function verifyBatch(bytes32 newLocalExitRoot, bytes32 newStateRoot, uint32 batchNum, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencySession) VerifyBatch(newLocalExitRoot [32]byte, newStateRoot [32]byte, batchNum uint32, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.VerifyBatch(&_Proofofefficiency.TransactOpts, newLocalExitRoot, newStateRoot, batchNum, proofA, proofB, proofC)
}

// VerifyBatch is a paid mutator transaction binding the contract method 0x3b880f77.
//
// Solidity: function verifyBatch(bytes32 newLocalExitRoot, bytes32 newStateRoot, uint32 batchNum, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) VerifyBatch(newLocalExitRoot [32]byte, newStateRoot [32]byte, batchNum uint32, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
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
	BatchNum  uint32
	Sequencer common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSendBatch is a free log retrieval operation binding the contract event 0xacbc4bfeebf66c85e73298070a2b645714fa69a7ef8890ed3bd6e9c581f7cdaf.
//
// Solidity: event SendBatch(uint32 indexed batchNum, address indexed sequencer)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSendBatch(opts *bind.FilterOpts, batchNum []uint32, sequencer []common.Address) (*ProofofefficiencySendBatchIterator, error) {

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

// WatchSendBatch is a free log subscription operation binding the contract event 0xacbc4bfeebf66c85e73298070a2b645714fa69a7ef8890ed3bd6e9c581f7cdaf.
//
// Solidity: event SendBatch(uint32 indexed batchNum, address indexed sequencer)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSendBatch(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySendBatch, batchNum []uint32, sequencer []common.Address) (event.Subscription, error) {

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

// ParseSendBatch is a log parse operation binding the contract event 0xacbc4bfeebf66c85e73298070a2b645714fa69a7ef8890ed3bd6e9c581f7cdaf.
//
// Solidity: event SendBatch(uint32 indexed batchNum, address indexed sequencer)
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
	BatchNum   uint32
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVerifyBatch is a free log retrieval operation binding the contract event 0xb0a69d17322b203355f6c64abe7663174dd1d3a3a10ea0fa12bb4788d3f865d5.
//
// Solidity: event VerifyBatch(uint32 indexed batchNum, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterVerifyBatch(opts *bind.FilterOpts, batchNum []uint32, aggregator []common.Address) (*ProofofefficiencyVerifyBatchIterator, error) {

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

// WatchVerifyBatch is a free log subscription operation binding the contract event 0xb0a69d17322b203355f6c64abe7663174dd1d3a3a10ea0fa12bb4788d3f865d5.
//
// Solidity: event VerifyBatch(uint32 indexed batchNum, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchVerifyBatch(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyVerifyBatch, batchNum []uint32, aggregator []common.Address) (event.Subscription, error) {

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

// ParseVerifyBatch is a log parse operation binding the contract event 0xb0a69d17322b203355f6c64abe7663174dd1d3a3a10ea0fa12bb4788d3f865d5.
//
// Solidity: event VerifyBatch(uint32 indexed batchNum, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseVerifyBatch(log types.Log) (*ProofofefficiencyVerifyBatch, error) {
	event := new(ProofofefficiencyVerifyBatch)
	if err := _Proofofefficiency.contract.UnpackLog(event, "VerifyBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
