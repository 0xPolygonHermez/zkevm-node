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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIGlobalExitRootManager\",\"name\":\"_globalExitRootManager\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_matic\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"_rollupVerifier\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"genesisRoot\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sequencerAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"RegisterSequencer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"batchChainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"lastGlobalExitRoot\",\"type\":\"bytes32\"}],\"name\":\"SendBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatch\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_CHAIN_ID\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateSequencerCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentLocalExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"globalExitRootManager\",\"outputs\":[{\"internalType\":\"contractIGlobalExitRootManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBatchSent\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastVerifiedBatch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"matic\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numSequencers\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"}],\"name\":\"registerSequencer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupVerifier\",\"outputs\":[{\"internalType\":\"contractIVerifierRollup\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maticAmount\",\"type\":\"uint256\"}],\"name\":\"sendBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"sentBatches\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"batchHashData\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maticCollateral\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"sequencers\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofA\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"proofB\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofC\",\"type\":\"uint256[2]\"}],\"name\":\"verifyBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b5060405162002a3e38038062002a3e83398181016040528101906200003791906200032d565b620000576200004b6200011e60201b60201c565b6200012660201b60201c565b83600460086101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508273ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff168152505081600760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600581905550505050506200039f565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200021c82620001ef565b9050919050565b600062000230826200020f565b9050919050565b620002428162000223565b81146200024e57600080fd5b50565b600081519050620002628162000237565b92915050565b600062000275826200020f565b9050919050565b620002878162000268565b81146200029357600080fd5b50565b600081519050620002a7816200027c565b92915050565b6000620002ba826200020f565b9050919050565b620002cc81620002ad565b8114620002d857600080fd5b50565b600081519050620002ec81620002c1565b92915050565b6000819050919050565b6200030781620002f2565b81146200031357600080fd5b50565b6000815190506200032781620002fc565b92915050565b600080600080608085870312156200034a5762000349620001ea565b5b60006200035a8782880162000251565b94505060206200036d8782880162000296565b93505060406200038087828801620002db565b9250506060620003938782880162000316565b91505092959194509250565b608051612675620003c96000396000818161039801528181610dff0152610ea901526126756000f3fe608060405234801561001057600080fd5b50600436106101165760003560e01c80638da5cb5b116100a2578063b6b0b09711610071578063b6b0b097146102ad578063ca98a308146102cb578063d02103ca146102e9578063e8bf92ed14610307578063f2fde38b1461032557610116565b80638da5cb5b1461023757806395297e2414610255578063959c2f4714610271578063ac2eba981461028f57610116565b806343ea1996116100e957806343ea1996146101a4578063519c1e7b146101c2578063715018a6146101f35780637fcb3653146101fd5780638a4abab81461021b57610116565b806306d6490f1461011b578063188ea07d146101375780631c7a07ee146101555780631dc125b214610186575b600080fd5b610135600480360381019061013091906116a0565b610341565b005b61013f6106e3565b60405161014c919061171f565b60405180910390f35b61016f600480360381019061016a9190611798565b6106fd565b60405161017d92919061184d565b60405180910390f35b61018e6107bd565b60405161019b919061188c565b60405180910390f35b6101ac610825565b6040516101b9919061171f565b60405180910390f35b6101dc60048036038101906101d791906118d3565b61082b565b6040516101ea929190611919565b60405180910390f35b6101fb61084f565b005b6102056108d7565b604051610212919061171f565b60405180910390f35b610235600480360381019061023091906119e3565b6108f1565b005b61023f610bb6565b60405161024c9190611a3b565b60405180910390f35b61026f600480360381019061026a9190611acb565b610bdf565b005b610279610e9b565b6040516102869190611b5a565b60405180910390f35b610297610ea1565b6040516102a49190611b5a565b60405180910390f35b6102b5610ea7565b6040516102c29190611bd4565b60405180910390f35b6102d3610ecb565b6040516102e09190611c0e565b60405180910390f35b6102f1610ee1565b6040516102fe9190611c4a565b60405180910390f35b61030f610f07565b60405161031c9190611c86565b60405180910390f35b61033f600480360381019061033a9190611798565b610f2d565b005b600061034b6107bd565b905081811115610390576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161038790611d13565b60405180910390fd5b6103dd3330837f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16611025909392919063ffffffff16565b6000600460089054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633ed691ef6040518163ffffffff1660e01b815260040160206040518083038186803b15801561044757600080fd5b505afa15801561045b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061047f9190611d48565b9050600080600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff1667ffffffffffffffff161461054857600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff16905061054e565b6103e890505b6002600481819054906101000a900467ffffffffffffffff168092919061057490611da4565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550508482423384600260049054906101000a900467ffffffffffffffff166040516020016105cc96959493929190611ebb565b6040516020818303038152906040528051906020012060036000600260049054906101000a900467ffffffffffffffff1667ffffffffffffffff1667ffffffffffffffff168152602001908152602001600020600001819055508260036000600260049054906101000a900467ffffffffffffffff1667ffffffffffffffff1667ffffffffffffffff168152602001908152602001600020600101819055503373ffffffffffffffffffffffffffffffffffffffff16600260049054906101000a900467ffffffffffffffff1667ffffffffffffffff167f84a31f45db3d6d8124a47a0544d32fff930986737dd276f558ff54382780487d83856040516106d4929190611f27565b60405180910390a35050505050565b600260049054906101000a900467ffffffffffffffff1681565b600160205280600052604060002060009150905080600001805461072090611f7f565b80601f016020809104026020016040519081016040528092919081815260200182805461074c90611f7f565b80156107995780601f1061076e57610100808354040283529160200191610799565b820191906000526020600020905b81548152906001019060200180831161077c57829003601f168201915b5050505050908060010160009054906101000a900467ffffffffffffffff16905082565b6000600460009054906101000a900467ffffffffffffffff16600260049054906101000a900467ffffffffffffffff1660016107f99190611fb1565b6108039190611fef565b670de0b6b3a76400006108169190612023565b67ffffffffffffffff16905090565b6103e881565b60036020528060005260406000206000915090508060000154908060010154905082565b6108576110ae565b73ffffffffffffffffffffffffffffffffffffffff16610875610bb6565b73ffffffffffffffffffffffffffffffffffffffff16146108cb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108c2906120b1565b60405180910390fd5b6108d560006110b6565b565b600460009054906101000a900467ffffffffffffffff1681565b600081511415610936576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161092d90612143565b60405180910390fd5b6000600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff1667ffffffffffffffff161415610aca576002600081819054906101000a900463ffffffff16809291906109c190612163565b91906101000a81548163ffffffff021916908363ffffffff1602179055505080600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000019080519060200190610a3692919061146d565b50600260009054906101000a900463ffffffff1663ffffffff166103e8610a5d9190611fb1565b600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160006101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550610b22565b80600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000019080519060200190610b2092919061146d565b505b7fcb9ff083daeaef4973981efa83966c7a440fc6eb5cda6b7bfad8cc5765cbbca83382600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160009054906101000a900467ffffffffffffffff16604051610bab93929190612190565b60405180910390a150565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6001600460009054906101000a900467ffffffffffffffff16610c029190611fb1565b67ffffffffffffffff168467ffffffffffffffff1614610c57576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c4e90612240565b60405180910390fd5b6000600360008667ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060405180604001604052908160008201548152602001600182015481525050905060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001600554600654898b8660000151604051602001610ce5959493929190612260565b6040516020818303038152906040528051906020012060001c610d0891906122ee565b90506004600081819054906101000a900467ffffffffffffffff1680929190610d3090611da4565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550508660058190555087600681905550600460089054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166333d6247d6006546040518263ffffffff1660e01b8152600401610dc29190611b5a565b600060405180830381600087803b158015610ddc57600080fd5b505af1158015610df0573d6000803e3d6000fd5b50505050610e433383602001517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1661117a9092919063ffffffff16565b3373ffffffffffffffffffffffffffffffffffffffff168667ffffffffffffffff167f2cdf1508085a46c7241a7d78c5a1ec3d9246d1ab95e1c2a33676d29e17d4222360405160405180910390a35050505050505050565b60065481565b60055481565b7f000000000000000000000000000000000000000000000000000000000000000081565b600260009054906101000a900463ffffffff1681565b600460089054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b610f356110ae565b73ffffffffffffffffffffffffffffffffffffffff16610f53610bb6565b73ffffffffffffffffffffffffffffffffffffffff1614610fa9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610fa0906120b1565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415611019576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161101090612391565b60405180910390fd5b611022816110b6565b50565b6110a8846323b872dd60e01b858585604051602401611046939291906123b1565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050611200565b50505050565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6111fb8363a9059cbb60e01b84846040516024016111999291906123e8565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050611200565b505050565b6000611262826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166112c79092919063ffffffff16565b90506000815111156112c257808060200190518101906112829190612449565b6112c1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112b8906124e8565b60405180910390fd5b5b505050565b60606112d684846000856112df565b90509392505050565b606082471015611324576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161131b9061257a565b60405180910390fd5b61132d856113f3565b61136c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611363906125e6565b60405180910390fd5b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516113959190612606565b60006040518083038185875af1925050503d80600081146113d2576040519150601f19603f3d011682016040523d82523d6000602084013e6113d7565b606091505b50915091506113e7828286611406565b92505050949350505050565b600080823b905060008111915050919050565b6060831561141657829050611466565b6000835111156114295782518084602001fd5b816040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161145d919061261d565b60405180910390fd5b9392505050565b82805461147990611f7f565b90600052602060002090601f01602090048101928261149b57600085556114e2565b82601f106114b457805160ff19168380011785556114e2565b828001600101855582156114e2579182015b828111156114e15782518255916020019190600101906114c6565b5b5090506114ef91906114f3565b5090565b5b8082111561150c5760008160009055506001016114f4565b5090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6115778261152e565b810181811067ffffffffffffffff821117156115965761159561153f565b5b80604052505050565b60006115a9611510565b90506115b5828261156e565b919050565b600067ffffffffffffffff8211156115d5576115d461153f565b5b6115de8261152e565b9050602081019050919050565b82818337600083830152505050565b600061160d611608846115ba565b61159f565b90508281526020810184848401111561162957611628611529565b5b6116348482856115eb565b509392505050565b600082601f83011261165157611650611524565b5b81356116618482602086016115fa565b91505092915050565b6000819050919050565b61167d8161166a565b811461168857600080fd5b50565b60008135905061169a81611674565b92915050565b600080604083850312156116b7576116b661151a565b5b600083013567ffffffffffffffff8111156116d5576116d461151f565b5b6116e18582860161163c565b92505060206116f28582860161168b565b9150509250929050565b600067ffffffffffffffff82169050919050565b611719816116fc565b82525050565b60006020820190506117346000830184611710565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006117658261173a565b9050919050565b6117758161175a565b811461178057600080fd5b50565b6000813590506117928161176c565b92915050565b6000602082840312156117ae576117ad61151a565b5b60006117bc84828501611783565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b838110156117ff5780820151818401526020810190506117e4565b8381111561180e576000848401525b50505050565b600061181f826117c5565b61182981856117d0565b93506118398185602086016117e1565b6118428161152e565b840191505092915050565b600060408201905081810360008301526118678185611814565b90506118766020830184611710565b9392505050565b6118868161166a565b82525050565b60006020820190506118a1600083018461187d565b92915050565b6118b0816116fc565b81146118bb57600080fd5b50565b6000813590506118cd816118a7565b92915050565b6000602082840312156118e9576118e861151a565b5b60006118f7848285016118be565b91505092915050565b6000819050919050565b61191381611900565b82525050565b600060408201905061192e600083018561190a565b61193b602083018461187d565b9392505050565b600067ffffffffffffffff82111561195d5761195c61153f565b5b6119668261152e565b9050602081019050919050565b600061198661198184611942565b61159f565b9050828152602081018484840111156119a2576119a1611529565b5b6119ad8482856115eb565b509392505050565b600082601f8301126119ca576119c9611524565b5b81356119da848260208601611973565b91505092915050565b6000602082840312156119f9576119f861151a565b5b600082013567ffffffffffffffff811115611a1757611a1661151f565b5b611a23848285016119b5565b91505092915050565b611a358161175a565b82525050565b6000602082019050611a506000830184611a2c565b92915050565b611a5f81611900565b8114611a6a57600080fd5b50565b600081359050611a7c81611a56565b92915050565b600080fd5b600081905082602060020282011115611aa357611aa2611a82565b5b92915050565b600081905082604060020282011115611ac557611ac4611a82565b5b92915050565b6000806000806000806101608789031215611ae957611ae861151a565b5b6000611af789828a01611a6d565b9650506020611b0889828a01611a6d565b9550506040611b1989828a016118be565b9450506060611b2a89828a01611a87565b93505060a0611b3b89828a01611aa9565b925050610120611b4d89828a01611a87565b9150509295509295509295565b6000602082019050611b6f600083018461190a565b92915050565b6000819050919050565b6000611b9a611b95611b908461173a565b611b75565b61173a565b9050919050565b6000611bac82611b7f565b9050919050565b6000611bbe82611ba1565b9050919050565b611bce81611bb3565b82525050565b6000602082019050611be96000830184611bc5565b92915050565b600063ffffffff82169050919050565b611c0881611bef565b82525050565b6000602082019050611c236000830184611bff565b92915050565b6000611c3482611ba1565b9050919050565b611c4481611c29565b82525050565b6000602082019050611c5f6000830184611c3b565b92915050565b6000611c7082611ba1565b9050919050565b611c8081611c65565b82525050565b6000602082019050611c9b6000830184611c77565b92915050565b7f50726f6f664f66456666696369656e63793a3a73656e6442617463683a204e4f60008201527f545f454e4f5547485f4d41544943000000000000000000000000000000000000602082015250565b6000611cfd602e836117d0565b9150611d0882611ca1565b604082019050919050565b60006020820190508181036000830152611d2c81611cf0565b9050919050565b600081519050611d4281611a56565b92915050565b600060208284031215611d5e57611d5d61151a565b5b6000611d6c84828501611d33565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611daf826116fc565b915067ffffffffffffffff821415611dca57611dc9611d75565b5b600182019050919050565b600081519050919050565b600081905092915050565b6000611df682611dd5565b611e008185611de0565b9350611e108185602086016117e1565b80840191505092915050565b6000819050919050565b611e37611e3282611900565b611e1c565b82525050565b60008160c01b9050919050565b6000611e5582611e3d565b9050919050565b611e6d611e68826116fc565b611e4a565b82525050565b60008160601b9050919050565b6000611e8b82611e73565b9050919050565b6000611e9d82611e80565b9050919050565b611eb5611eb08261175a565b611e92565b82525050565b6000611ec78289611deb565b9150611ed38288611e26565b602082019150611ee38287611e5c565b600882019150611ef38286611ea4565b601482019150611f038285611e5c565b600882019150611f138284611e5c565b600882019150819050979650505050505050565b6000604082019050611f3c6000830185611710565b611f49602083018461190a565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680611f9757607f821691505b60208210811415611fab57611faa611f50565b5b50919050565b6000611fbc826116fc565b9150611fc7836116fc565b92508267ffffffffffffffff03821115611fe457611fe3611d75565b5b828201905092915050565b6000611ffa826116fc565b9150612005836116fc565b92508282101561201857612017611d75565b5b828203905092915050565b600061202e826116fc565b9150612039836116fc565b92508167ffffffffffffffff048311821515161561205a57612059611d75565b5b828202905092915050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b600061209b6020836117d0565b91506120a682612065565b602082019050919050565b600060208201905081810360008301526120ca8161208e565b9050919050565b7f50726f6f664f66456666696369656e63793a3a7265676973746572536571756560008201527f6e6365723a204e4f545f56414c49445f55524c00000000000000000000000000602082015250565b600061212d6033836117d0565b9150612138826120d1565b604082019050919050565b6000602082019050818103600083015261215c81612120565b9050919050565b600061216e82611bef565b915063ffffffff82141561218557612184611d75565b5b600182019050919050565b60006060820190506121a56000830186611a2c565b81810360208301526121b78185611814565b90506121c66040830184611710565b949350505050565b7f50726f6f664f66456666696369656e63793a3a76657269667942617463683a2060008201527f42415443485f444f45535f4e4f545f4d41544348000000000000000000000000602082015250565b600061222a6034836117d0565b9150612235826121ce565b604082019050919050565b600060208201905081810360008301526122598161221d565b9050919050565b600061226c8288611e26565b60208201915061227c8287611e26565b60208201915061228c8286611e26565b60208201915061229c8285611e26565b6020820191506122ac8284611e26565b6020820191508190509695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006122f98261166a565b91506123048361166a565b925082612314576123136122bf565b5b828206905092915050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b600061237b6026836117d0565b91506123868261231f565b604082019050919050565b600060208201905081810360008301526123aa8161236e565b9050919050565b60006060820190506123c66000830186611a2c565b6123d36020830185611a2c565b6123e0604083018461187d565b949350505050565b60006040820190506123fd6000830185611a2c565b61240a602083018461187d565b9392505050565b60008115159050919050565b61242681612411565b811461243157600080fd5b50565b6000815190506124438161241d565b92915050565b60006020828403121561245f5761245e61151a565b5b600061246d84828501612434565b91505092915050565b7f5361666545524332303a204552433230206f7065726174696f6e20646964206e60008201527f6f74207375636365656400000000000000000000000000000000000000000000602082015250565b60006124d2602a836117d0565b91506124dd82612476565b604082019050919050565b60006020820190508181036000830152612501816124c5565b9050919050565b7f416464726573733a20696e73756666696369656e742062616c616e636520666f60008201527f722063616c6c0000000000000000000000000000000000000000000000000000602082015250565b60006125646026836117d0565b915061256f82612508565b604082019050919050565b6000602082019050818103600083015261259381612557565b9050919050565b7f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000600082015250565b60006125d0601d836117d0565b91506125db8261259a565b602082019050919050565b600060208201905081810360008301526125ff816125c3565b9050919050565b60006126128284611deb565b915081905092915050565b600060208201905081810360008301526126378184611814565b90509291505056fea26469706673582212205b04951b09cc4703275d13bf118e3bd8a467ab1a6bfe84782a5c1ac8858762e664736f6c63430008090033",
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
// Solidity: function DEFAULT_CHAIN_ID() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) DEFAULTCHAINID(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "DEFAULT_CHAIN_ID")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// DEFAULTCHAINID is a free data retrieval call binding the contract method 0x43ea1996.
//
// Solidity: function DEFAULT_CHAIN_ID() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) DEFAULTCHAINID() (uint64, error) {
	return _Proofofefficiency.Contract.DEFAULTCHAINID(&_Proofofefficiency.CallOpts)
}

// DEFAULTCHAINID is a free data retrieval call binding the contract method 0x43ea1996.
//
// Solidity: function DEFAULT_CHAIN_ID() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) DEFAULTCHAINID() (uint64, error) {
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
// Solidity: function lastBatchSent() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) LastBatchSent(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "lastBatchSent")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastBatchSent is a free data retrieval call binding the contract method 0x188ea07d.
//
// Solidity: function lastBatchSent() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) LastBatchSent() (uint64, error) {
	return _Proofofefficiency.Contract.LastBatchSent(&_Proofofefficiency.CallOpts)
}

// LastBatchSent is a free data retrieval call binding the contract method 0x188ea07d.
//
// Solidity: function lastBatchSent() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) LastBatchSent() (uint64, error) {
	return _Proofofefficiency.Contract.LastBatchSent(&_Proofofefficiency.CallOpts)
}

// LastVerifiedBatch is a free data retrieval call binding the contract method 0x7fcb3653.
//
// Solidity: function lastVerifiedBatch() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) LastVerifiedBatch(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "lastVerifiedBatch")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastVerifiedBatch is a free data retrieval call binding the contract method 0x7fcb3653.
//
// Solidity: function lastVerifiedBatch() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) LastVerifiedBatch() (uint64, error) {
	return _Proofofefficiency.Contract.LastVerifiedBatch(&_Proofofefficiency.CallOpts)
}

// LastVerifiedBatch is a free data retrieval call binding the contract method 0x7fcb3653.
//
// Solidity: function lastVerifiedBatch() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) LastVerifiedBatch() (uint64, error) {
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

// SentBatches is a free data retrieval call binding the contract method 0x519c1e7b.
//
// Solidity: function sentBatches(uint64 ) view returns(bytes32 batchHashData, uint256 maticCollateral)
func (_Proofofefficiency *ProofofefficiencyCaller) SentBatches(opts *bind.CallOpts, arg0 uint64) (struct {
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

// SentBatches is a free data retrieval call binding the contract method 0x519c1e7b.
//
// Solidity: function sentBatches(uint64 ) view returns(bytes32 batchHashData, uint256 maticCollateral)
func (_Proofofefficiency *ProofofefficiencySession) SentBatches(arg0 uint64) (struct {
	BatchHashData   [32]byte
	MaticCollateral *big.Int
}, error) {
	return _Proofofefficiency.Contract.SentBatches(&_Proofofefficiency.CallOpts, arg0)
}

// SentBatches is a free data retrieval call binding the contract method 0x519c1e7b.
//
// Solidity: function sentBatches(uint64 ) view returns(bytes32 batchHashData, uint256 maticCollateral)
func (_Proofofefficiency *ProofofefficiencyCallerSession) SentBatches(arg0 uint64) (struct {
	BatchHashData   [32]byte
	MaticCollateral *big.Int
}, error) {
	return _Proofofefficiency.Contract.SentBatches(&_Proofofefficiency.CallOpts, arg0)
}

// Sequencers is a free data retrieval call binding the contract method 0x1c7a07ee.
//
// Solidity: function sequencers(address ) view returns(string sequencerURL, uint64 chainID)
func (_Proofofefficiency *ProofofefficiencyCaller) Sequencers(opts *bind.CallOpts, arg0 common.Address) (struct {
	SequencerURL string
	ChainID      uint64
}, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "sequencers", arg0)

	outstruct := new(struct {
		SequencerURL string
		ChainID      uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SequencerURL = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.ChainID = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// Sequencers is a free data retrieval call binding the contract method 0x1c7a07ee.
//
// Solidity: function sequencers(address ) view returns(string sequencerURL, uint64 chainID)
func (_Proofofefficiency *ProofofefficiencySession) Sequencers(arg0 common.Address) (struct {
	SequencerURL string
	ChainID      uint64
}, error) {
	return _Proofofefficiency.Contract.Sequencers(&_Proofofefficiency.CallOpts, arg0)
}

// Sequencers is a free data retrieval call binding the contract method 0x1c7a07ee.
//
// Solidity: function sequencers(address ) view returns(string sequencerURL, uint64 chainID)
func (_Proofofefficiency *ProofofefficiencyCallerSession) Sequencers(arg0 common.Address) (struct {
	SequencerURL string
	ChainID      uint64
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

// VerifyBatch is a paid mutator transaction binding the contract method 0x95297e24.
//
// Solidity: function verifyBatch(bytes32 newLocalExitRoot, bytes32 newStateRoot, uint64 numBatch, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) VerifyBatch(opts *bind.TransactOpts, newLocalExitRoot [32]byte, newStateRoot [32]byte, numBatch uint64, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "verifyBatch", newLocalExitRoot, newStateRoot, numBatch, proofA, proofB, proofC)
}

// VerifyBatch is a paid mutator transaction binding the contract method 0x95297e24.
//
// Solidity: function verifyBatch(bytes32 newLocalExitRoot, bytes32 newStateRoot, uint64 numBatch, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencySession) VerifyBatch(newLocalExitRoot [32]byte, newStateRoot [32]byte, numBatch uint64, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.VerifyBatch(&_Proofofefficiency.TransactOpts, newLocalExitRoot, newStateRoot, numBatch, proofA, proofB, proofC)
}

// VerifyBatch is a paid mutator transaction binding the contract method 0x95297e24.
//
// Solidity: function verifyBatch(bytes32 newLocalExitRoot, bytes32 newStateRoot, uint64 numBatch, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) VerifyBatch(newLocalExitRoot [32]byte, newStateRoot [32]byte, numBatch uint64, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
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
	ChainID          uint64
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRegisterSequencer is a free log retrieval operation binding the contract event 0xcb9ff083daeaef4973981efa83966c7a440fc6eb5cda6b7bfad8cc5765cbbca8.
//
// Solidity: event RegisterSequencer(address sequencerAddress, string sequencerURL, uint64 chainID)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterRegisterSequencer(opts *bind.FilterOpts) (*ProofofefficiencyRegisterSequencerIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "RegisterSequencer")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyRegisterSequencerIterator{contract: _Proofofefficiency.contract, event: "RegisterSequencer", logs: logs, sub: sub}, nil
}

// WatchRegisterSequencer is a free log subscription operation binding the contract event 0xcb9ff083daeaef4973981efa83966c7a440fc6eb5cda6b7bfad8cc5765cbbca8.
//
// Solidity: event RegisterSequencer(address sequencerAddress, string sequencerURL, uint64 chainID)
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

// ParseRegisterSequencer is a log parse operation binding the contract event 0xcb9ff083daeaef4973981efa83966c7a440fc6eb5cda6b7bfad8cc5765cbbca8.
//
// Solidity: event RegisterSequencer(address sequencerAddress, string sequencerURL, uint64 chainID)
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
	NumBatch           uint64
	Sequencer          common.Address
	BatchChainID       uint64
	LastGlobalExitRoot [32]byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterSendBatch is a free log retrieval operation binding the contract event 0x84a31f45db3d6d8124a47a0544d32fff930986737dd276f558ff54382780487d.
//
// Solidity: event SendBatch(uint64 indexed numBatch, address indexed sequencer, uint64 batchChainID, bytes32 lastGlobalExitRoot)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSendBatch(opts *bind.FilterOpts, numBatch []uint64, sequencer []common.Address) (*ProofofefficiencySendBatchIterator, error) {

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

// WatchSendBatch is a free log subscription operation binding the contract event 0x84a31f45db3d6d8124a47a0544d32fff930986737dd276f558ff54382780487d.
//
// Solidity: event SendBatch(uint64 indexed numBatch, address indexed sequencer, uint64 batchChainID, bytes32 lastGlobalExitRoot)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSendBatch(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySendBatch, numBatch []uint64, sequencer []common.Address) (event.Subscription, error) {

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

// ParseSendBatch is a log parse operation binding the contract event 0x84a31f45db3d6d8124a47a0544d32fff930986737dd276f558ff54382780487d.
//
// Solidity: event SendBatch(uint64 indexed numBatch, address indexed sequencer, uint64 batchChainID, bytes32 lastGlobalExitRoot)
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
	NumBatch   uint64
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVerifyBatch is a free log retrieval operation binding the contract event 0x2cdf1508085a46c7241a7d78c5a1ec3d9246d1ab95e1c2a33676d29e17d42223.
//
// Solidity: event VerifyBatch(uint64 indexed numBatch, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterVerifyBatch(opts *bind.FilterOpts, numBatch []uint64, aggregator []common.Address) (*ProofofefficiencyVerifyBatchIterator, error) {

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

// WatchVerifyBatch is a free log subscription operation binding the contract event 0x2cdf1508085a46c7241a7d78c5a1ec3d9246d1ab95e1c2a33676d29e17d42223.
//
// Solidity: event VerifyBatch(uint64 indexed numBatch, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchVerifyBatch(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyVerifyBatch, numBatch []uint64, aggregator []common.Address) (event.Subscription, error) {

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

// ParseVerifyBatch is a log parse operation binding the contract event 0x2cdf1508085a46c7241a7d78c5a1ec3d9246d1ab95e1c2a33676d29e17d42223.
//
// Solidity: event VerifyBatch(uint64 indexed numBatch, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseVerifyBatch(log types.Log) (*ProofofefficiencyVerifyBatch, error) {
	event := new(ProofofefficiencyVerifyBatch)
	if err := _Proofofefficiency.contract.UnpackLog(event, "VerifyBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
