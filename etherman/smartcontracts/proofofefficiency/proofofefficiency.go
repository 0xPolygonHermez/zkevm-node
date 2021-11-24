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
	Bin: "0x60a06040523480156200001157600080fd5b5060405162001cc238038062001cc28339818101604052810190620000379190620002ea565b620000576200004b6200011660201b60201c565b6200011e60201b60201c565b82600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff168152505080600a60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050505062000346565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200021482620001e7565b9050919050565b6000620002288262000207565b9050919050565b6200023a816200021b565b81146200024657600080fd5b50565b6000815190506200025a816200022f565b92915050565b60006200026d8262000207565b9050919050565b6200027f8162000260565b81146200028b57600080fd5b50565b6000815190506200029f8162000274565b92915050565b6000620002b28262000207565b9050919050565b620002c481620002a5565b8114620002d057600080fd5b50565b600081519050620002e481620002b9565b92915050565b600080600060608486031215620003065762000305620001e2565b5b6000620003168682870162000249565b935050602062000329868287016200028e565b92505060406200033c86828701620002d3565b9150509250925092565b6080516119606200036260003960006109b701526119606000f3fe608060405234801561001057600080fd5b50600436106101165760003560e01c80638da5cb5b116100a2578063b6b0b09711610071578063b6b0b097146102ae578063ca98a308146102cc578063e78cea92146102ea578063e8bf92ed14610308578063f2fde38b1461032657610116565b80638da5cb5b14610238578063959c2f4714610256578063a152b62e14610274578063ac2eba981461029057610116565b80631dc125b2116100e95780631dc125b2146101a4578063715018a6146101c25780637fcb3653146101cc57806384b0be8c146101ea5780638a4abab81461021c57610116565b806306d6490f1461011b5780630acd922c14610137578063188ea07d146101555780631c7a07ee14610173575b600080fd5b61013560048036038101906101309190610e22565b610342565b005b61013f6103a4565b60405161014c9190610e97565b60405180910390f35b61015d6103aa565b60405161016a9190610ec1565b60405180910390f35b61018d60048036038101906101889190610f3a565b6103b0565b60405161019b929190610fef565b60405180910390f35b6101ac61045c565b6040516101b99190610ec1565b60405180910390f35b6101ca61046c565b005b6101d46104f4565b6040516101e19190610ec1565b60405180910390f35b61020460048036038101906101ff919061101f565b6104fa565b6040516102139392919061105b565b60405180910390f35b61023660048036038101906102319190611133565b610544565b005b610240610725565b60405161024d919061117c565b60405180910390f35b61025e61074e565b60405161026b9190610e97565b60405180910390f35b61028e6004803603810190610289919061120c565b610754565b005b6102986109af565b6040516102a59190610e97565b60405180910390f35b6102b66109b5565b6040516102c391906112fa565b60405180910390f35b6102d46109d9565b6040516102e19190610ec1565b60405180910390f35b6102f26109df565b6040516102ff9190611336565b60405180910390f35b610310610a05565b60405161031d9190611372565b60405180910390f35b610340600480360381019061033b9190610f3a565b610a2b565b005b60036000815480929190610355906113bc565b91905055503373ffffffffffffffffffffffffffffffffffffffff166003547fe99ca860ab96f187daa5c05e5d805b16a9dd0242a7e113c065d2a13ce30159d060405160405180910390a35050565b60095481565b60035481565b60016020528060005260406000206000915090508060000180546103d390611434565b80601f01602080910402602001604051908101604052809291908181526020018280546103ff90611434565b801561044c5780601f106104215761010080835404028352916020019161044c565b820191906000526020600020905b81548152906001019060200180831161042f57829003601f168201915b5050505050908060010154905082565b6000670de0b6b3a7640000905090565b610474610b23565b73ffffffffffffffffffffffffffffffffffffffff16610492610725565b73ffffffffffffffffffffffffffffffffffffffff16146104e8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104df906114b2565b60405180910390fd5b6104f26000610b2b565b565b60055481565b60046020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060010154908060020154905083565b600081511415610589576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161058090611544565b60405180910390fd5b6000600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010154141561069157600260008154809291906105e7906113bc565b919050555080600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000019080519060200190610642929190610bef565b50600254600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101819055506106e9565b80600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000190805190602001906106e7929190610bef565b505b7f3c5b8bb3cdafd1a15eaa861069e1567306bab24615ef281949384e54b80bd77d338260405161071a929190611564565b60405180910390a150565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60085481565b60016005546107639190611594565b84146107a4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161079b9061165c565b60405180910390fd5b6000600460008681526020019081526020016000206040518060600160405290816000820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820154815260200160028201548152505090506000816000015190506000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010154905060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000160026007546008546009548d8f898b602001518a6040516020016108ce989796959493929190611706565b6040516020818303038152906040526040516108ea91906117df565b602060405180830381855afa158015610907573d6000803e3d6000fd5b5050506040513d601f19601f8201168201806040525081019061092a919061180b565b60001c6109379190611867565b90506005600081548092919061094c906113bc565b919050555088600781905550896008819055503373ffffffffffffffffffffffffffffffffffffffff16887f60f152ed8e8084aa9231a20aae908de8d039a671f90e86d6b92c87889d952b2260405160405180910390a350505050505050505050565b60075481565b7f000000000000000000000000000000000000000000000000000000000000000081565b60025481565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600a60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b610a33610b23565b73ffffffffffffffffffffffffffffffffffffffff16610a51610725565b73ffffffffffffffffffffffffffffffffffffffff1614610aa7576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a9e906114b2565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610b17576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b0e9061190a565b60405180910390fd5b610b2081610b2b565b50565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b828054610bfb90611434565b90600052602060002090601f016020900481019282610c1d5760008555610c64565b82601f10610c3657805160ff1916838001178555610c64565b82800160010185558215610c64579182015b82811115610c63578251825591602001919060010190610c48565b5b509050610c719190610c75565b5090565b5b80821115610c8e576000816000905550600101610c76565b5090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610cf982610cb0565b810181811067ffffffffffffffff82111715610d1857610d17610cc1565b5b80604052505050565b6000610d2b610c92565b9050610d378282610cf0565b919050565b600067ffffffffffffffff821115610d5757610d56610cc1565b5b610d6082610cb0565b9050602081019050919050565b82818337600083830152505050565b6000610d8f610d8a84610d3c565b610d21565b905082815260208101848484011115610dab57610daa610cab565b5b610db6848285610d6d565b509392505050565b600082601f830112610dd357610dd2610ca6565b5b8135610de3848260208601610d7c565b91505092915050565b6000819050919050565b610dff81610dec565b8114610e0a57600080fd5b50565b600081359050610e1c81610df6565b92915050565b60008060408385031215610e3957610e38610c9c565b5b600083013567ffffffffffffffff811115610e5757610e56610ca1565b5b610e6385828601610dbe565b9250506020610e7485828601610e0d565b9150509250929050565b6000819050919050565b610e9181610e7e565b82525050565b6000602082019050610eac6000830184610e88565b92915050565b610ebb81610dec565b82525050565b6000602082019050610ed66000830184610eb2565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610f0782610edc565b9050919050565b610f1781610efc565b8114610f2257600080fd5b50565b600081359050610f3481610f0e565b92915050565b600060208284031215610f5057610f4f610c9c565b5b6000610f5e84828501610f25565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610fa1578082015181840152602081019050610f86565b83811115610fb0576000848401525b50505050565b6000610fc182610f67565b610fcb8185610f72565b9350610fdb818560208601610f83565b610fe481610cb0565b840191505092915050565b600060408201905081810360008301526110098185610fb6565b90506110186020830184610eb2565b9392505050565b60006020828403121561103557611034610c9c565b5b600061104384828501610e0d565b91505092915050565b61105581610efc565b82525050565b6000606082019050611070600083018661104c565b61107d6020830185610e88565b61108a6040830184610eb2565b949350505050565b600067ffffffffffffffff8211156110ad576110ac610cc1565b5b6110b682610cb0565b9050602081019050919050565b60006110d66110d184611092565b610d21565b9050828152602081018484840111156110f2576110f1610cab565b5b6110fd848285610d6d565b509392505050565b600082601f83011261111a57611119610ca6565b5b813561112a8482602086016110c3565b91505092915050565b60006020828403121561114957611148610c9c565b5b600082013567ffffffffffffffff81111561116757611166610ca1565b5b61117384828501611105565b91505092915050565b6000602082019050611191600083018461104c565b92915050565b6111a081610e7e565b81146111ab57600080fd5b50565b6000813590506111bd81611197565b92915050565b600080fd5b6000819050826020600202820111156111e4576111e36111c3565b5b92915050565b600081905082604060020282011115611206576112056111c3565b5b92915050565b600080600080600080610160878903121561122a57611229610c9c565b5b600061123889828a016111ae565b965050602061124989828a016111ae565b955050604061125a89828a01610e0d565b945050606061126b89828a016111c8565b93505060a061127c89828a016111ea565b92505061012061128e89828a016111c8565b9150509295509295509295565b6000819050919050565b60006112c06112bb6112b684610edc565b61129b565b610edc565b9050919050565b60006112d2826112a5565b9050919050565b60006112e4826112c7565b9050919050565b6112f4816112d9565b82525050565b600060208201905061130f60008301846112eb565b92915050565b6000611320826112c7565b9050919050565b61133081611315565b82525050565b600060208201905061134b6000830184611327565b92915050565b600061135c826112c7565b9050919050565b61136c81611351565b82525050565b60006020820190506113876000830184611363565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006113c782610dec565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156113fa576113f961138d565b5b600182019050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061144c57607f821691505b602082108114156114605761145f611405565b5b50919050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b600061149c602083610f72565b91506114a782611466565b602082019050919050565b600060208201905081810360008301526114cb8161148f565b9050919050565b7f50726f6f664f66456666696369656e63793a3a7265676973746572536571756560008201527f6e6365723a204e4f545f56414c49445f55524c00000000000000000000000000602082015250565b600061152e603383610f72565b9150611539826114d2565b604082019050919050565b6000602082019050818103600083015261155d81611521565b9050919050565b6000604082019050611579600083018561104c565b818103602083015261158b8184610fb6565b90509392505050565b600061159f82610dec565b91506115aa83610dec565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156115df576115de61138d565b5b828201905092915050565b7f50726f6f664f66456666696369656e63793a3a76657269667942617463683a2060008201527f42415443485f444f45535f4e4f545f4d41544348000000000000000000000000602082015250565b6000611646603483610f72565b9150611651826115ea565b604082019050919050565b6000602082019050818103600083015261167581611639565b9050919050565b6000819050919050565b61169761169282610e7e565b61167c565b82525050565b60008160601b9050919050565b60006116b58261169d565b9050919050565b60006116c7826116aa565b9050919050565b6116df6116da82610efc565b6116bc565b82525050565b6000819050919050565b6117006116fb82610dec565b6116e5565b82525050565b6000611712828b611686565b602082019150611722828a611686565b6020820191506117328289611686565b6020820191506117428288611686565b6020820191506117528287611686565b60208201915061176282866116ce565b6014820191506117728285611686565b60208201915061178282846116ef565b6020820191508190509998505050505050505050565b600081519050919050565b600081905092915050565b60006117b982611798565b6117c381856117a3565b93506117d3818560208601610f83565b80840191505092915050565b60006117eb82846117ae565b915081905092915050565b60008151905061180581611197565b92915050565b60006020828403121561182157611820610c9c565b5b600061182f848285016117f6565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061187282610dec565b915061187d83610dec565b92508261188d5761188c611838565b5b828206905092915050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b60006118f4602683610f72565b91506118ff82611898565b604082019050919050565b60006020820190508181036000830152611923816118e7565b905091905056fea26469706673582212207c0cea9002998559280978578f0f03570354a8722e0102eff71d9d3438ae1f0964736f6c63430008090033",
}

// ProofofefficiencyABI is the input ABI used to generate the binding from.
// Deprecated: Use ProofofefficiencyMetaData.ABI instead.
var ProofofefficiencyABI = ProofofefficiencyMetaData.ABI

// ProofofefficiencyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ProofofefficiencyMetaData.Bin instead.
var ProofofefficiencyBin = ProofofefficiencyMetaData.Bin

// DeployProofofefficiency deploys a new Ethereum contract, binding an instance of Proofofefficiency to it.
func DeployProofofefficiency(auth *bind.TransactOpts, backend bind.ContractBackend, _bridge common.Address, _matic common.Address, _rollupVerifier common.Address) (common.Address, *types.Transaction, *Proofofefficiency, error) {
	parsed, err := ProofofefficiencyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ProofofefficiencyBin), backend, _bridge, _matic, _rollupVerifier)
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
