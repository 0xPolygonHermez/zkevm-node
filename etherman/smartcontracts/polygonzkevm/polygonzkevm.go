// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package polygonzkevm

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

// PolygonRollupBaseEtrogBatchData is an auto generated low-level Go binding around an user-defined struct.
type PolygonRollupBaseEtrogBatchData struct {
	Transactions         []byte
	ForcedGlobalExitRoot [32]byte
	ForcedTimestamp      uint64
	ForcedBlockHashL1    [32]byte
}

// PolygonzkevmMetaData contains all meta data concerning the Polygonzkevm contract.
var PolygonzkevmMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRootV2\",\"name\":\"_globalExitRootManager\",\"type\":\"address\"},{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"_pol\",\"type\":\"address\"},{\"internalType\":\"contractIPolygonZkEVMBridgeV2\",\"name\":\"_bridgeAddress\",\"type\":\"address\"},{\"internalType\":\"contractPolygonRollupManager\",\"name\":\"_rollupManager\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BatchAlreadyVerified\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BatchNotSequencedOrNotSequenceEnd\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExceedMaxVerifyBatches\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumBatchBelowLastVerifiedBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumBatchDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalPendingStateNumInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBatchNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBatchTimeoutNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBatchesAlreadyActive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBatchesOverflow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForcedDataDoesNotMatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasTokenNetworkMustBeZeroOnEther\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GlobalExitRootNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HaltTimeoutNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HugeTokenMetadataNotSupported\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitNumBatchAboveLastVerifiedBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitNumBatchDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitializeTransaction\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeBatchTimeTarget\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeForceBatchTimeout\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeMultiplierBatchFee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewAccInputHashDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewPendingStateTimeoutMustBeLower\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewStateRootNotInsidePrime\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewTrustedAggregatorTimeoutMustBeLower\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughMaticAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughPOLAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OldAccInputHashDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OldStateRootDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPendingAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyRollupManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyTrustedAggregator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyTrustedSequencer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateNotConsolidable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateTimeoutExceedHaltAggregationTimeout\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SequenceZeroBatches\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SequencedTimestampBelowForcedTimestamp\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SequencedTimestampInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StoredRootMustBeDifferentThanNewRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransactionsLengthAboveMax\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TrustedAggregatorTimeoutExceedHaltAggregationTimeout\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TrustedAggregatorTimeoutNotExpired\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AcceptAdminRole\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"ActivateForceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"forceBatchNum\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"lastGlobalExitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"}],\"name\":\"ForceBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"lastGlobalExitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"}],\"name\":\"InitialSequenceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"l1InfoRoot\",\"type\":\"bytes32\"}],\"name\":\"SequenceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"}],\"name\":\"SequenceForceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newforceBatchTimeout\",\"type\":\"uint64\"}],\"name\":\"SetForceBatchTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTrustedSequencer\",\"type\":\"address\"}],\"name\":\"SetTrustedSequencer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"newTrustedSequencerURL\",\"type\":\"string\"}],\"name\":\"SetTrustedSequencerURL\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newPendingAdmin\",\"type\":\"address\"}],\"name\":\"TransferAdminRole\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatches\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"GLOBAL_EXIT_ROOT_MANAGER_L2\",\"outputs\":[{\"internalType\":\"contractIBasePolygonZkEVMGlobalExitRoot\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_BRIDGE_LIST_LEN_LEN\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_BRIDGE_PARAMS\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_CONSTANT_BYTES\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_DATA_LEN_EMPTY_METADATA\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_EFFECTIVE_PERCENTAGE\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"\",\"type\":\"bytes1\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SIGNATURE_INITIALIZE_TX_R\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SIGNATURE_INITIALIZE_TX_S\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SIGNATURE_INITIALIZE_TX_V\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptAdminRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activateForceBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMBridgeV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculatePolPerForceBatch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"polAmount\",\"type\":\"uint256\"}],\"name\":\"forceBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"forceBatchTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"forcedBatches\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gapLastTimestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasTokenNetwork\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"networkID\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_gasTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_gasTokenNetwork\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"_gasTokenMetadata\",\"type\":\"bytes\"}],\"name\":\"generateInitializeTransaction\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"globalExitRootManager\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRootV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"networkID\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_gasTokenAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_networkName\",\"type\":\"string\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isForcedBatchAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastAccInputHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastForceBatch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastForceBatchSequenced\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"onVerifyBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pol\",\"outputs\":[{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupManager\",\"outputs\":[{\"internalType\":\"contractPolygonRollupManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"forcedGlobalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"forcedTimestamp\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"forcedBlockHashL1\",\"type\":\"bytes32\"}],\"internalType\":\"structPolygonRollupBaseEtrog.BatchData[]\",\"name\":\"batches\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"l2Coinbase\",\"type\":\"address\"}],\"name\":\"sequenceBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"forcedGlobalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"forcedTimestamp\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"forcedBlockHashL1\",\"type\":\"bytes32\"}],\"internalType\":\"structPolygonRollupBaseEtrog.BatchData[]\",\"name\":\"batches\",\"type\":\"tuple[]\"}],\"name\":\"sequenceForceBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newforceBatchTimeout\",\"type\":\"uint64\"}],\"name\":\"setForceBatchTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newTrustedSequencer\",\"type\":\"address\"}],\"name\":\"setTrustedSequencer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newTrustedSequencerURL\",\"type\":\"string\"}],\"name\":\"setTrustedSequencerURL\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newPendingAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedSequencer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedSequencerURL\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b506040516200464138038062004641833981016040819052620000359162000071565b6001600160a01b0393841660a052918316608052821660c0521660e052620000d9565b6001600160a01b03811681146200006e57600080fd5b50565b600080600080608085870312156200008857600080fd5b8451620000958162000058565b6020860151909450620000a88162000058565b6040860151909350620000bb8162000058565b6060860151909250620000ce8162000058565b939692955090935050565b60805160a05160c05160e051614494620001ad6000396000818161051a015281816109a501528181610b1201528181610e4f015281816113ea01528181611a8d01528181611b8b015281816121dc01528181612b4901528181612bca01528181612bec0152612d04015260008181610669015281816110a901528181611d5b01528181611e63015261266f01526000818161072d0152818161125e0152818161232701526126f101526000818161075f0152818161083d01528181611ad6015281816122fb0152612c9a01526144946000f3fe608060405234801561001057600080fd5b50600436106102e85760003560e01c80636ff512cc11610191578063c7fffd4b116100e3578063e46761c411610097578063ecef3f9911610071578063ecef3f99146107b4578063f35dda47146107c7578063f851a440146107cf57600080fd5b8063e46761c41461075a578063e7a7ed0214610781578063eaeb077b146107a157600080fd5b8063cfa8ed47116100c8578063cfa8ed4714610708578063d02103ca14610728578063d7bc90ff1461074f57600080fd5b8063c7fffd4b146106ed578063c89e42df146106f557600080fd5b80639f26f84011610145578063ada8f9191161011f578063ada8f9191461069e578063b0afe154146106b1578063c754c7ed146106bd57600080fd5b80639f26f84014610651578063a3c573eb14610664578063a652f26c1461068b57600080fd5b80637a5460c5116101765780637a5460c5146105f25780638c3d73011461062e5780639e0018771461063657600080fd5b80636ff512cc146105cc57806371257022146105df57600080fd5b80633cbc795b1161024a57806352bdeb6d116101fe578063676870d2116101d8578063676870d21461059b5780636b8616ce146105a35780636e05d2cd146105c357600080fd5b806352bdeb6d1461054f578063542028d51461058b5780635ec919581461059357600080fd5b8063456052671161022f57806345605267146104ed57806349b7b802146105155780634e4877061461053c57600080fd5b80633cbc795b1461045757806340b5de6c1461049557600080fd5b806326782247116102a15780632b35b3ac116102865780632b35b3ac1461040057806332c2d1531461041d5780633c351e101461043257600080fd5b8063267822471461038e5780632a9f208f146103d357600080fd5b806305835f37116102d257806305835f3714610323578063107bf28c1461036c57806311e892d41461037457600080fd5b8062d0295d146102ed5780630350896314610308575b600080fd5b6102f56107f5565b6040519081526020015b60405180910390f35b610310602081565b60405161ffff90911681526020016102ff565b61035f6040518060400160405280600881526020017f80808401c9c3809400000000000000000000000000000000000000000000000081525081565b6040516102ff91906136e1565b61035f610915565b61037c60f981565b60405160ff90911681526020016102ff565b6001546103ae9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016102ff565b6007546103e79067ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020016102ff565b60085461040d9060ff1681565b60405190151581526020016102ff565b61043061042b366004613736565b6109a3565b005b6008546103ae90610100900473ffffffffffffffffffffffffffffffffffffffff1681565b600854610480907501000000000000000000000000000000000000000000900463ffffffff1681565b60405163ffffffff90911681526020016102ff565b6104bc7fff0000000000000000000000000000000000000000000000000000000000000081565b6040517fff0000000000000000000000000000000000000000000000000000000000000090911681526020016102ff565b6007546103e790700100000000000000000000000000000000900467ffffffffffffffff1681565b6103ae7f000000000000000000000000000000000000000000000000000000000000000081565b61043061054a366004613778565b610a72565b61035f6040518060400160405280600281526020017f80b800000000000000000000000000000000000000000000000000000000000081525081565b61035f610c8c565b610430610c99565b610310601f81565b6102f56105b1366004613778565b60066020526000908152604090205481565b6102f560055481565b6104306105da366004613795565b610d83565b6104306105ed3660046138de565b610e4d565b61035f6040518060400160405280600281526020017f80b900000000000000000000000000000000000000000000000000000000000081525081565b6104306115e5565b6103ae73a40d5f56745a118d0906a34e69aec8c0db1cb8fa81565b61043061065f3660046139d7565b6116b8565b6103ae7f000000000000000000000000000000000000000000000000000000000000000081565b61035f610699366004613a19565b611c5a565b6104306106ac366004613795565b61203f565b6102f56405ca1ab1e081565b6007546103e7907801000000000000000000000000000000000000000000000000900467ffffffffffffffff1681565b61037c60e481565b610430610703366004613a8e565b612109565b6002546103ae9073ffffffffffffffffffffffffffffffffffffffff1681565b6103ae7f000000000000000000000000000000000000000000000000000000000000000081565b6102f5635ca1ab1e81565b6103ae7f000000000000000000000000000000000000000000000000000000000000000081565b6007546103e79068010000000000000000900467ffffffffffffffff1681565b6104306107af366004613ac3565b61219c565b6104306107c2366004613b3b565b6125a4565b61037c601b81565b6000546103ae9062010000900473ffffffffffffffffffffffffffffffffffffffff1681565b6040517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152600090819073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa158015610884573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108a89190613b87565b6007549091506000906108e79067ffffffffffffffff700100000000000000000000000000000000820481169168010000000000000000900416613bcf565b67ffffffffffffffff169050806000036109045760009250505090565b61090e8183613bf7565b9250505090565b6004805461092290613c32565b80601f016020809104026020016040519081016040528092919081815260200182805461094e90613c32565b801561099b5780601f106109705761010080835404028352916020019161099b565b820191906000526020600020905b81548152906001019060200180831161097e57829003601f168201915b505050505081565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163314610a12576040517fb9b3a2c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff168367ffffffffffffffff167f9c72852172521097ba7e1482e6b44b351323df0155f97f4ea18fcec28e1f596684604051610a6591815260200190565b60405180910390a3505050565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff163314610ac9576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b62093a8067ffffffffffffffff82161115610b10576040517ff5e37f2f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166315064c966040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b7b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b9f9190613c85565b610c085760075467ffffffffffffffff7801000000000000000000000000000000000000000000000000909104811690821610610c08576040517ff5e37f2f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6007805477ffffffffffffffffffffffffffffffffffffffffffffffff16780100000000000000000000000000000000000000000000000067ffffffffffffffff8416908102919091179091556040519081527fa7eb6cb8a613eb4e8bddc1ac3d61ec6cf10898760f0b187bcca794c6ca6fa40b906020015b60405180910390a150565b6003805461092290613c32565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff163314610cf0576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60085460ff1615610d2d576040517ff6ba91a100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600880547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556040517f854dd6ce5a1445c4c54388b21cffd11cf5bba1b9e763aec48ce3da75d617412f90600090a1565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff163314610dda576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527ff54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc090602001610c81565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163314610ebc576040517fb9b3a2c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600054610100900460ff1615808015610edc5750600054600160ff909116105b80610ef65750303b158015610ef6575060005460ff166001145b610f87576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084015b60405180910390fd5b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790558015610fe557600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b606073ffffffffffffffffffffffffffffffffffffffff8516156111fa5761100c85612dd9565b61101586612ee9565b61101e87612ff0565b60405160200161103093929190613ca7565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290527f318aee3d00000000000000000000000000000000000000000000000000000000825273ffffffffffffffffffffffffffffffffffffffff878116600484015290925060009182917f0000000000000000000000000000000000000000000000000000000000000000169063318aee3d9060240160408051808303816000875af11580156110f1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111159190613ce0565b915091508163ffffffff166000146111b257600880547fffffffffffffff000000000000000000000000000000000000000000000000ff1661010073ffffffffffffffffffffffffffffffffffffffff8416027fffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffff1617750100000000000000000000000000000000000000000063ffffffff8516021790556111f7565b600880547fffffffffffffffffffffff0000000000000000000000000000000000000000ff1661010073ffffffffffffffffffffffffffffffffffffffff8a16021790555b50505b600854600090611247908890610100810473ffffffffffffffffffffffffffffffffffffffff16907501000000000000000000000000000000000000000000900463ffffffff1685611c5a565b9050600081805190602001209050600042905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16633ed691ef6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156112c7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112eb9190613b87565b90506000808483858f6112ff600143613d1a565b60408051602081019790975286019490945260608086019390935260c09190911b7fffffffffffffffff000000000000000000000000000000000000000000000000166080850152901b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016608883015240609c82015260bc01604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815290829052805160209091012060058190557f9a908e73000000000000000000000000000000000000000000000000000000008252600160048301526024820181905291507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690639a908e73906044016020604051808303816000875af1158015611448573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061146c9190613d2d565b508c600060026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508b600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555088600390816114fe9190613d98565b50600461150b8982613d98565b5062069780600760186101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055507f060116213bcbf54ca19fd649dc84b59ab2bbd200ab199770e4d923e222a28e7f85838e60405161156b93929190613eb2565b60405180910390a150505050505080156115dc57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314611636576040517fd1ec4b2300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600154600080547fffffffffffffffffffff0000000000000000000000000000000000000000ffff1673ffffffffffffffffffffffffffffffffffffffff9092166201000081029290921790556040519081527f056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e9060200160405180910390a1565b60085460ff166116f4576040517f24eff8c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806000819003611730576040517fcb591a5f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103e881111561176c576040517fb59f753a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60075467ffffffffffffffff6801000000000000000082048116916117a7918491700100000000000000000000000000000000900416613ef1565b11156117df576040517fc630a00d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60075460055470010000000000000000000000000000000090910467ffffffffffffffff169060005b83811015611a8757600086868381811061182457611824613f04565b90506020028101906118369190613f33565b61183f90613f71565b90508361184b81613ffa565b825180516020918201208185015160408087015160608801519151959a509295506000946118b8948794929101938452602084019290925260c01b7fffffffffffffffff000000000000000000000000000000000000000000000000166040830152604882015260680190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152815160209283012067ffffffffffffffff8916600090815260069093529120549091508114611941576040517fce3d755e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8616600090815260066020526040812055611966600188613d1a565b84036119d55742600760189054906101000a900467ffffffffffffffff1684604001516119939190614021565b67ffffffffffffffff1611156119d5576040517fc44a082100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208381015160408086015160608088015183519586018b90529285018790528481019390935260c01b7fffffffffffffffff0000000000000000000000000000000000000000000000001660808401523390911b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166088830152609c82015260bc016040516020818303038152906040528051906020012094505050508080611a7f90614042565b915050611808565b50611afd7f000000000000000000000000000000000000000000000000000000000000000084611ab56107f5565b611abf919061407a565b73ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001691906130e3565b60058190556007805467ffffffffffffffff8416700100000000000000000000000000000000027fffffffffffffffff0000000000000000ffffffffffffffffffffffffffffffff9091161790556040517f9a908e7300000000000000000000000000000000000000000000000000000000815260009073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690639a908e7390611bd7908790869060040167ffffffffffffffff929092168252602082015260400190565b6020604051808303816000875af1158015611bf6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c1a9190613d2d565b60405190915067ffffffffffffffff8216907f648a61dd2438f072f5a1960939abd30f37aea80d2e94c9792ad142d3e0a490a490600090a2505050505050565b6060600085858573a40d5f56745a118d0906a34e69aec8c0db1cb8fa600087604051602401611c8e96959493929190614091565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167ff811bff7000000000000000000000000000000000000000000000000000000001790528351909150606090600003611ddf5760f9601f8351611d2391906140f4565b6040518060400160405280600881526020017f80808401c9c380940000000000000000000000000000000000000000000000008152507f00000000000000000000000000000000000000000000000000000000000000006040518060400160405280600281526020017f80b800000000000000000000000000000000000000000000000000000000000081525060e487604051602001611dc9979695949392919061410f565b6040516020818303038152906040529050611ee3565b815161ffff1015611e1c576040517f248b8f8200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b815160f9611e2b6020836140f4565b6040518060400160405280600881526020017f80808401c9c380940000000000000000000000000000000000000000000000008152507f00000000000000000000000000000000000000000000000000000000000000006040518060400160405280600281526020017f80b90000000000000000000000000000000000000000000000000000000000008152508588604051602001611ed097969594939291906141f2565b6040516020818303038152906040529150505b805160208083019190912060408051600080825293810180835292909252601b908201526405ca1ab1e06060820152635ca1ab1e608082015260019060a0016020604051602081039080840390855afa158015611f44573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff8116611fbc576040517fcd16196600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040516000906120029084906405ca1ab1e090635ca1ab1e90601b907fff00000000000000000000000000000000000000000000000000000000000000906020016142d5565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190529450505050505b949350505050565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff163314612096576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527fa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce690602001610c81565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff163314612160576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600361216c8282613d98565b507f6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b2081604051610c8191906136e1565b60085460ff166121d8576040517f24eff8c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663604691696040518163ffffffff1660e01b8152600401602060405180830381865afa158015612245573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122699190613b87565b9050818111156122a5576040517f2354600f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6113888311156122e1576040517fa29a6c7c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61232373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000163330846131bc565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16633ed691ef6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612390573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123b49190613b87565b600780549192506801000000000000000090910467ffffffffffffffff169060086123de83613ffa565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550508484604051612415929190614331565b604051908190039020814261242b600143613d1a565b60408051602081019590955284019290925260c01b7fffffffffffffffff000000000000000000000000000000000000000000000000166060830152406068820152608801604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152918152815160209283012060075468010000000000000000900467ffffffffffffffff166000908152600690935291205532330361253e576007546040805183815233602082015260609181018290526000918101919091526801000000000000000090910467ffffffffffffffff16907ff94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc9319060800160405180910390a261259d565b600760089054906101000a900467ffffffffffffffff1667ffffffffffffffff167ff94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc931823388886040516125949493929190614341565b60405180910390a25b5050505050565b60025473ffffffffffffffffffffffffffffffffffffffff1633146125f5576040517f11e7be1500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b816000819003612631576040517fcb591a5f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103e881111561266d576040517fb59f753a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166379e2cf976040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156126d557600080fd5b505af11580156126e9573d6000803e3d6000fd5b5050505060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16635ca1e1656040518163ffffffff1660e01b8152600401602060405180830381865afa15801561275a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061277e9190613b87565b600754600554919250429170010000000000000000000000000000000090910467ffffffffffffffff16908160005b86811015612aaf5760008a8a838181106127c9576127c9613f04565b90506020028101906127db9190613f33565b6127e490613f71565b8051805160209091012060408201519192509067ffffffffffffffff16156129c9578561281081613ffa565b9650506000818360200151846040015185606001516040516020016128739493929190938452602084019290925260c01b7fffffffffffffffff000000000000000000000000000000000000000000000000166040830152604882015260680190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152815160209283012067ffffffffffffffff8a166000908152600690935291205490915081146128fc576040517fce3d755e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208381015160408086015160608088015183519586018c90529285018790528481019390935260c01b7fffffffffffffffff000000000000000000000000000000000000000000000000166080840152908d901b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166088830152609c82015260bc01604051602081830303815290604052805190602001209550600660008867ffffffffffffffff1667ffffffffffffffff1681526020019081526020016000206000905550612a9a565b8151516201d4c01015612a08576040517fa29a6c7c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516020810187905290810182905260608082018a905260c089901b7fffffffffffffffff0000000000000000000000000000000000000000000000001660808301528b901b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001660888201526000609c82015260bc016040516020818303038152906040528051906020012094505b50508080612aa790614042565b9150506127ad565b5060075467ffffffffffffffff6801000000000000000090910481169084161115612b06576040517fc630a00d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60058290558567ffffffffffffffff84811690831614612bc4576000612b2c8386613bcf565b9050612b4267ffffffffffffffff821683613d1a565b9150612b7b7f00000000000000000000000000000000000000000000000000000000000000008267ffffffffffffffff16611ab56107f5565b50600780547fffffffffffffffff0000000000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000067ffffffffffffffff8716021790555b612cc2337f0000000000000000000000000000000000000000000000000000000000000000837f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663477fa2706040518163ffffffff1660e01b8152600401602060405180830381865afa158015612c55573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612c799190613b87565b612c83919061407a565b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169291906131bc565b6040517f9a908e7300000000000000000000000000000000000000000000000000000000815267ffffffffffffffff88166004820152602481018490526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690639a908e73906044016020604051808303816000875af1158015612d62573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d869190613d2d565b90508067ffffffffffffffff167f3e54d0825ed78523037d00a81759237eb436ce774bd546993ee67a1b67b6e76688604051612dc491815260200190565b60405180910390a25050505050505050505050565b60408051600481526024810182526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f06fdde03000000000000000000000000000000000000000000000000000000001790529051606091600091829173ffffffffffffffffffffffffffffffffffffffff861691612e5b91906143b2565b600060405180830381855afa9150503d8060008114612e96576040519150601f19603f3d011682016040523d82523d6000602084013e612e9b565b606091505b509150915081612ee0576040518060400160405280600781526020017f4e4f5f4e414d4500000000000000000000000000000000000000000000000000815250612037565b61203781613220565b60408051600481526024810182526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f95d89b41000000000000000000000000000000000000000000000000000000001790529051606091600091829173ffffffffffffffffffffffffffffffffffffffff861691612f6b91906143b2565b600060405180830381855afa9150503d8060008114612fa6576040519150601f19603f3d011682016040523d82523d6000602084013e612fab565b606091505b509150915081612ee0576040518060400160405280600981526020017f4e4f5f53594d424f4c0000000000000000000000000000000000000000000000815250612037565b60408051600481526024810182526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f313ce5670000000000000000000000000000000000000000000000000000000017905290516000918291829173ffffffffffffffffffffffffffffffffffffffff86169161307191906143b2565b600060405180830381855afa9150503d80600081146130ac576040519150601f19603f3d011682016040523d82523d6000602084013e6130b1565b606091505b50915091508180156130c4575080516020145b6130cf576012612037565b8080602001905181019061203791906143c4565b60405173ffffffffffffffffffffffffffffffffffffffff83166024820152604481018290526131b79084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526133f6565b505050565b60405173ffffffffffffffffffffffffffffffffffffffff8085166024830152831660448201526064810182905261321a9085907f23b872dd0000000000000000000000000000000000000000000000000000000090608401613135565b50505050565b60606040825110613245578180602001905181019061323f91906143e7565b92915050565b81516020036133bd5760005b602081108015613298575082818151811061326e5761326e613f04565b01602001517fff000000000000000000000000000000000000000000000000000000000000001615155b156132af57806132a781614042565b915050613251565b806000036132f257505060408051808201909152601281527f4e4f545f56414c49445f454e434f44494e4700000000000000000000000000006020820152919050565b60008167ffffffffffffffff81111561330d5761330d6137c4565b6040519080825280601f01601f191660200182016040528015613337576020820181803683370190505b50905060005b828110156133b55784818151811061335757613357613f04565b602001015160f81c60f81b82828151811061337457613374613f04565b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350806133ad81614042565b91505061333d565b509392505050565b505060408051808201909152601281527f4e4f545f56414c49445f454e434f44494e470000000000000000000000000000602082015290565b6000613458826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166135029092919063ffffffff16565b8051909150156131b757808060200190518101906134769190613c85565b6131b7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610f7e565b60606120378484600085856000808673ffffffffffffffffffffffffffffffffffffffff16858760405161353691906143b2565b60006040518083038185875af1925050503d8060008114613573576040519150601f19603f3d011682016040523d82523d6000602084013e613578565b606091505b509150915061358987838387613594565b979650505050505050565b6060831561362a5782516000036136235773ffffffffffffffffffffffffffffffffffffffff85163b613623576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610f7e565b5081612037565b612037838381511561363f5781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f7e91906136e1565b60005b8381101561368e578181015183820152602001613676565b50506000910152565b600081518084526136af816020860160208601613673565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006136f46020830184613697565b9392505050565b67ffffffffffffffff8116811461371157600080fd5b50565b73ffffffffffffffffffffffffffffffffffffffff8116811461371157600080fd5b60008060006060848603121561374b57600080fd5b8335613756816136fb565b925060208401359150604084013561376d81613714565b809150509250925092565b60006020828403121561378a57600080fd5b81356136f4816136fb565b6000602082840312156137a757600080fd5b81356136f481613714565b63ffffffff8116811461371157600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561383a5761383a6137c4565b604052919050565b600067ffffffffffffffff82111561385c5761385c6137c4565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f83011261389957600080fd5b81356138ac6138a782613842565b6137f3565b8181528460208386010111156138c157600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c087890312156138f757600080fd5b863561390281613714565b9550602087013561391281613714565b94506040870135613922816137b2565b9350606087013561393281613714565b9250608087013567ffffffffffffffff8082111561394f57600080fd5b61395b8a838b01613888565b935060a089013591508082111561397157600080fd5b5061397e89828a01613888565b9150509295509295509295565b60008083601f84011261399d57600080fd5b50813567ffffffffffffffff8111156139b557600080fd5b6020830191508360208260051b85010111156139d057600080fd5b9250929050565b600080602083850312156139ea57600080fd5b823567ffffffffffffffff811115613a0157600080fd5b613a0d8582860161398b565b90969095509350505050565b60008060008060808587031215613a2f57600080fd5b8435613a3a816137b2565b93506020850135613a4a81613714565b92506040850135613a5a816137b2565b9150606085013567ffffffffffffffff811115613a7657600080fd5b613a8287828801613888565b91505092959194509250565b600060208284031215613aa057600080fd5b813567ffffffffffffffff811115613ab757600080fd5b61203784828501613888565b600080600060408486031215613ad857600080fd5b833567ffffffffffffffff80821115613af057600080fd5b818601915086601f830112613b0457600080fd5b813581811115613b1357600080fd5b876020828501011115613b2557600080fd5b6020928301989097509590910135949350505050565b600080600060408486031215613b5057600080fd5b833567ffffffffffffffff811115613b6757600080fd5b613b738682870161398b565b909450925050602084013561376d81613714565b600060208284031215613b9957600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b67ffffffffffffffff828116828216039080821115613bf057613bf0613ba0565b5092915050565b600082613c2d577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b600181811c90821680613c4657607f821691505b602082108103613c7f577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600060208284031215613c9757600080fd5b815180151581146136f457600080fd5b606081526000613cba6060830186613697565b8281036020840152613ccc8186613697565b91505060ff83166040830152949350505050565b60008060408385031215613cf357600080fd5b8251613cfe816137b2565b6020840151909250613d0f81613714565b809150509250929050565b8181038181111561323f5761323f613ba0565b600060208284031215613d3f57600080fd5b81516136f4816136fb565b601f8211156131b757600081815260208120601f850160051c81016020861015613d715750805b601f850160051c820191505b81811015613d9057828155600101613d7d565b505050505050565b815167ffffffffffffffff811115613db257613db26137c4565b613dc681613dc08454613c32565b84613d4a565b602080601f831160018114613e195760008415613de35750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555613d90565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015613e6657888601518255948401946001909101908401613e47565b5085821015613ea257878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b606081526000613ec56060830186613697565b905083602083015273ffffffffffffffffffffffffffffffffffffffff83166040830152949350505050565b8082018082111561323f5761323f613ba0565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81833603018112613f6757600080fd5b9190910192915050565b600060808236031215613f8357600080fd5b6040516080810167ffffffffffffffff8282108183111715613fa757613fa76137c4565b816040528435915080821115613fbc57600080fd5b50613fc936828601613888565b825250602083013560208201526040830135613fe4816136fb565b6040820152606092830135928101929092525090565b600067ffffffffffffffff80831681810361401757614017613ba0565b6001019392505050565b67ffffffffffffffff818116838216019080821115613bf057613bf0613ba0565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361407357614073613ba0565b5060010190565b808202811582820484141761323f5761323f613ba0565b600063ffffffff808916835273ffffffffffffffffffffffffffffffffffffffff8089166020850152818816604085015280871660608501528086166080850152505060c060a08301526140e860c0830184613697565b98975050505050505050565b61ffff818116838216019080821115613bf057613bf0613ba0565b60007fff00000000000000000000000000000000000000000000000000000000000000808a60f81b1683527fffff0000000000000000000000000000000000000000000000000000000000008960f01b1660018401528751614178816003860160208c01613673565b80840190507fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008860601b16600382015286516141bb816017840160208b01613673565b808201915050818660f81b166017820152845191506141e1826018830160208801613673565b016018019998505050505050505050565b7fff000000000000000000000000000000000000000000000000000000000000008860f81b16815260007fffff000000000000000000000000000000000000000000000000000000000000808960f01b166001840152875161425b816003860160208c01613673565b80840190507fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008860601b166003820152865161429e816017840160208b01613673565b808201915050818660f01b166017820152845191506142c4826019830160208801613673565b016019019998505050505050505050565b600086516142e7818460208b01613673565b9190910194855250602084019290925260f81b7fff000000000000000000000000000000000000000000000000000000000000009081166040840152166041820152604201919050565b8183823760009101908152919050565b84815273ffffffffffffffffffffffffffffffffffffffff8416602082015260606040820152816060820152818360808301376000818301608090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601019392505050565b60008251613f67818460208701613673565b6000602082840312156143d657600080fd5b815160ff811681146136f457600080fd5b6000602082840312156143f957600080fd5b815167ffffffffffffffff81111561441057600080fd5b8201601f8101841361442157600080fd5b805161442f6138a782613842565b81815285602083850101111561444457600080fd5b614455826020830160208601613673565b9594505050505056fea2646970667358221220c9b142edbaf63c56761aad22240658b45bae78a2ed1e483abe43944f9fae087664736f6c63430008140033",
}

// PolygonzkevmABI is the input ABI used to generate the binding from.
// Deprecated: Use PolygonzkevmMetaData.ABI instead.
var PolygonzkevmABI = PolygonzkevmMetaData.ABI

// PolygonzkevmBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PolygonzkevmMetaData.Bin instead.
var PolygonzkevmBin = PolygonzkevmMetaData.Bin

// DeployPolygonzkevm deploys a new Ethereum contract, binding an instance of Polygonzkevm to it.
func DeployPolygonzkevm(auth *bind.TransactOpts, backend bind.ContractBackend, _globalExitRootManager common.Address, _pol common.Address, _bridgeAddress common.Address, _rollupManager common.Address) (common.Address, *types.Transaction, *Polygonzkevm, error) {
	parsed, err := PolygonzkevmMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PolygonzkevmBin), backend, _globalExitRootManager, _pol, _bridgeAddress, _rollupManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Polygonzkevm{PolygonzkevmCaller: PolygonzkevmCaller{contract: contract}, PolygonzkevmTransactor: PolygonzkevmTransactor{contract: contract}, PolygonzkevmFilterer: PolygonzkevmFilterer{contract: contract}}, nil
}

// Polygonzkevm is an auto generated Go binding around an Ethereum contract.
type Polygonzkevm struct {
	PolygonzkevmCaller     // Read-only binding to the contract
	PolygonzkevmTransactor // Write-only binding to the contract
	PolygonzkevmFilterer   // Log filterer for contract events
}

// PolygonzkevmCaller is an auto generated read-only Go binding around an Ethereum contract.
type PolygonzkevmCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolygonzkevmTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PolygonzkevmTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolygonzkevmFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PolygonzkevmFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolygonzkevmSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PolygonzkevmSession struct {
	Contract     *Polygonzkevm     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PolygonzkevmCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PolygonzkevmCallerSession struct {
	Contract *PolygonzkevmCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// PolygonzkevmTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PolygonzkevmTransactorSession struct {
	Contract     *PolygonzkevmTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// PolygonzkevmRaw is an auto generated low-level Go binding around an Ethereum contract.
type PolygonzkevmRaw struct {
	Contract *Polygonzkevm // Generic contract binding to access the raw methods on
}

// PolygonzkevmCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PolygonzkevmCallerRaw struct {
	Contract *PolygonzkevmCaller // Generic read-only contract binding to access the raw methods on
}

// PolygonzkevmTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PolygonzkevmTransactorRaw struct {
	Contract *PolygonzkevmTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPolygonzkevm creates a new instance of Polygonzkevm, bound to a specific deployed contract.
func NewPolygonzkevm(address common.Address, backend bind.ContractBackend) (*Polygonzkevm, error) {
	contract, err := bindPolygonzkevm(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Polygonzkevm{PolygonzkevmCaller: PolygonzkevmCaller{contract: contract}, PolygonzkevmTransactor: PolygonzkevmTransactor{contract: contract}, PolygonzkevmFilterer: PolygonzkevmFilterer{contract: contract}}, nil
}

// NewPolygonzkevmCaller creates a new read-only instance of Polygonzkevm, bound to a specific deployed contract.
func NewPolygonzkevmCaller(address common.Address, caller bind.ContractCaller) (*PolygonzkevmCaller, error) {
	contract, err := bindPolygonzkevm(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmCaller{contract: contract}, nil
}

// NewPolygonzkevmTransactor creates a new write-only instance of Polygonzkevm, bound to a specific deployed contract.
func NewPolygonzkevmTransactor(address common.Address, transactor bind.ContractTransactor) (*PolygonzkevmTransactor, error) {
	contract, err := bindPolygonzkevm(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmTransactor{contract: contract}, nil
}

// NewPolygonzkevmFilterer creates a new log filterer instance of Polygonzkevm, bound to a specific deployed contract.
func NewPolygonzkevmFilterer(address common.Address, filterer bind.ContractFilterer) (*PolygonzkevmFilterer, error) {
	contract, err := bindPolygonzkevm(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmFilterer{contract: contract}, nil
}

// bindPolygonzkevm binds a generic wrapper to an already deployed contract.
func bindPolygonzkevm(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PolygonzkevmMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Polygonzkevm *PolygonzkevmRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Polygonzkevm.Contract.PolygonzkevmCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Polygonzkevm *PolygonzkevmRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.PolygonzkevmTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Polygonzkevm *PolygonzkevmRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.PolygonzkevmTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Polygonzkevm *PolygonzkevmCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Polygonzkevm.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Polygonzkevm *PolygonzkevmTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Polygonzkevm *PolygonzkevmTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.contract.Transact(opts, method, params...)
}

// GLOBALEXITROOTMANAGERL2 is a free data retrieval call binding the contract method 0x9e001877.
//
// Solidity: function GLOBAL_EXIT_ROOT_MANAGER_L2() view returns(address)
func (_Polygonzkevm *PolygonzkevmCaller) GLOBALEXITROOTMANAGERL2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "GLOBAL_EXIT_ROOT_MANAGER_L2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GLOBALEXITROOTMANAGERL2 is a free data retrieval call binding the contract method 0x9e001877.
//
// Solidity: function GLOBAL_EXIT_ROOT_MANAGER_L2() view returns(address)
func (_Polygonzkevm *PolygonzkevmSession) GLOBALEXITROOTMANAGERL2() (common.Address, error) {
	return _Polygonzkevm.Contract.GLOBALEXITROOTMANAGERL2(&_Polygonzkevm.CallOpts)
}

// GLOBALEXITROOTMANAGERL2 is a free data retrieval call binding the contract method 0x9e001877.
//
// Solidity: function GLOBAL_EXIT_ROOT_MANAGER_L2() view returns(address)
func (_Polygonzkevm *PolygonzkevmCallerSession) GLOBALEXITROOTMANAGERL2() (common.Address, error) {
	return _Polygonzkevm.Contract.GLOBALEXITROOTMANAGERL2(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGELISTLENLEN is a free data retrieval call binding the contract method 0x11e892d4.
//
// Solidity: function INITIALIZE_TX_BRIDGE_LIST_LEN_LEN() view returns(uint8)
func (_Polygonzkevm *PolygonzkevmCaller) INITIALIZETXBRIDGELISTLENLEN(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_BRIDGE_LIST_LEN_LEN")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// INITIALIZETXBRIDGELISTLENLEN is a free data retrieval call binding the contract method 0x11e892d4.
//
// Solidity: function INITIALIZE_TX_BRIDGE_LIST_LEN_LEN() view returns(uint8)
func (_Polygonzkevm *PolygonzkevmSession) INITIALIZETXBRIDGELISTLENLEN() (uint8, error) {
	return _Polygonzkevm.Contract.INITIALIZETXBRIDGELISTLENLEN(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGELISTLENLEN is a free data retrieval call binding the contract method 0x11e892d4.
//
// Solidity: function INITIALIZE_TX_BRIDGE_LIST_LEN_LEN() view returns(uint8)
func (_Polygonzkevm *PolygonzkevmCallerSession) INITIALIZETXBRIDGELISTLENLEN() (uint8, error) {
	return _Polygonzkevm.Contract.INITIALIZETXBRIDGELISTLENLEN(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMS is a free data retrieval call binding the contract method 0x05835f37.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS() view returns(bytes)
func (_Polygonzkevm *PolygonzkevmCaller) INITIALIZETXBRIDGEPARAMS(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_BRIDGE_PARAMS")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// INITIALIZETXBRIDGEPARAMS is a free data retrieval call binding the contract method 0x05835f37.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS() view returns(bytes)
func (_Polygonzkevm *PolygonzkevmSession) INITIALIZETXBRIDGEPARAMS() ([]byte, error) {
	return _Polygonzkevm.Contract.INITIALIZETXBRIDGEPARAMS(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMS is a free data retrieval call binding the contract method 0x05835f37.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS() view returns(bytes)
func (_Polygonzkevm *PolygonzkevmCallerSession) INITIALIZETXBRIDGEPARAMS() ([]byte, error) {
	return _Polygonzkevm.Contract.INITIALIZETXBRIDGEPARAMS(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS is a free data retrieval call binding the contract method 0x7a5460c5.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS() view returns(bytes)
func (_Polygonzkevm *PolygonzkevmCaller) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS is a free data retrieval call binding the contract method 0x7a5460c5.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS() view returns(bytes)
func (_Polygonzkevm *PolygonzkevmSession) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS() ([]byte, error) {
	return _Polygonzkevm.Contract.INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS is a free data retrieval call binding the contract method 0x7a5460c5.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS() view returns(bytes)
func (_Polygonzkevm *PolygonzkevmCallerSession) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS() ([]byte, error) {
	return _Polygonzkevm.Contract.INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA is a free data retrieval call binding the contract method 0x52bdeb6d.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA() view returns(bytes)
func (_Polygonzkevm *PolygonzkevmCaller) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA is a free data retrieval call binding the contract method 0x52bdeb6d.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA() view returns(bytes)
func (_Polygonzkevm *PolygonzkevmSession) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA() ([]byte, error) {
	return _Polygonzkevm.Contract.INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA is a free data retrieval call binding the contract method 0x52bdeb6d.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA() view returns(bytes)
func (_Polygonzkevm *PolygonzkevmCallerSession) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA() ([]byte, error) {
	return _Polygonzkevm.Contract.INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXCONSTANTBYTES is a free data retrieval call binding the contract method 0x03508963.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES() view returns(uint16)
func (_Polygonzkevm *PolygonzkevmCaller) INITIALIZETXCONSTANTBYTES(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_CONSTANT_BYTES")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// INITIALIZETXCONSTANTBYTES is a free data retrieval call binding the contract method 0x03508963.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES() view returns(uint16)
func (_Polygonzkevm *PolygonzkevmSession) INITIALIZETXCONSTANTBYTES() (uint16, error) {
	return _Polygonzkevm.Contract.INITIALIZETXCONSTANTBYTES(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXCONSTANTBYTES is a free data retrieval call binding the contract method 0x03508963.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES() view returns(uint16)
func (_Polygonzkevm *PolygonzkevmCallerSession) INITIALIZETXCONSTANTBYTES() (uint16, error) {
	return _Polygonzkevm.Contract.INITIALIZETXCONSTANTBYTES(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXCONSTANTBYTESEMPTYMETADATA is a free data retrieval call binding the contract method 0x676870d2.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA() view returns(uint16)
func (_Polygonzkevm *PolygonzkevmCaller) INITIALIZETXCONSTANTBYTESEMPTYMETADATA(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// INITIALIZETXCONSTANTBYTESEMPTYMETADATA is a free data retrieval call binding the contract method 0x676870d2.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA() view returns(uint16)
func (_Polygonzkevm *PolygonzkevmSession) INITIALIZETXCONSTANTBYTESEMPTYMETADATA() (uint16, error) {
	return _Polygonzkevm.Contract.INITIALIZETXCONSTANTBYTESEMPTYMETADATA(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXCONSTANTBYTESEMPTYMETADATA is a free data retrieval call binding the contract method 0x676870d2.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA() view returns(uint16)
func (_Polygonzkevm *PolygonzkevmCallerSession) INITIALIZETXCONSTANTBYTESEMPTYMETADATA() (uint16, error) {
	return _Polygonzkevm.Contract.INITIALIZETXCONSTANTBYTESEMPTYMETADATA(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXDATALENEMPTYMETADATA is a free data retrieval call binding the contract method 0xc7fffd4b.
//
// Solidity: function INITIALIZE_TX_DATA_LEN_EMPTY_METADATA() view returns(uint8)
func (_Polygonzkevm *PolygonzkevmCaller) INITIALIZETXDATALENEMPTYMETADATA(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_DATA_LEN_EMPTY_METADATA")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// INITIALIZETXDATALENEMPTYMETADATA is a free data retrieval call binding the contract method 0xc7fffd4b.
//
// Solidity: function INITIALIZE_TX_DATA_LEN_EMPTY_METADATA() view returns(uint8)
func (_Polygonzkevm *PolygonzkevmSession) INITIALIZETXDATALENEMPTYMETADATA() (uint8, error) {
	return _Polygonzkevm.Contract.INITIALIZETXDATALENEMPTYMETADATA(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXDATALENEMPTYMETADATA is a free data retrieval call binding the contract method 0xc7fffd4b.
//
// Solidity: function INITIALIZE_TX_DATA_LEN_EMPTY_METADATA() view returns(uint8)
func (_Polygonzkevm *PolygonzkevmCallerSession) INITIALIZETXDATALENEMPTYMETADATA() (uint8, error) {
	return _Polygonzkevm.Contract.INITIALIZETXDATALENEMPTYMETADATA(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXEFFECTIVEPERCENTAGE is a free data retrieval call binding the contract method 0x40b5de6c.
//
// Solidity: function INITIALIZE_TX_EFFECTIVE_PERCENTAGE() view returns(bytes1)
func (_Polygonzkevm *PolygonzkevmCaller) INITIALIZETXEFFECTIVEPERCENTAGE(opts *bind.CallOpts) ([1]byte, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_EFFECTIVE_PERCENTAGE")

	if err != nil {
		return *new([1]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)

	return out0, err

}

// INITIALIZETXEFFECTIVEPERCENTAGE is a free data retrieval call binding the contract method 0x40b5de6c.
//
// Solidity: function INITIALIZE_TX_EFFECTIVE_PERCENTAGE() view returns(bytes1)
func (_Polygonzkevm *PolygonzkevmSession) INITIALIZETXEFFECTIVEPERCENTAGE() ([1]byte, error) {
	return _Polygonzkevm.Contract.INITIALIZETXEFFECTIVEPERCENTAGE(&_Polygonzkevm.CallOpts)
}

// INITIALIZETXEFFECTIVEPERCENTAGE is a free data retrieval call binding the contract method 0x40b5de6c.
//
// Solidity: function INITIALIZE_TX_EFFECTIVE_PERCENTAGE() view returns(bytes1)
func (_Polygonzkevm *PolygonzkevmCallerSession) INITIALIZETXEFFECTIVEPERCENTAGE() ([1]byte, error) {
	return _Polygonzkevm.Contract.INITIALIZETXEFFECTIVEPERCENTAGE(&_Polygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXR is a free data retrieval call binding the contract method 0xb0afe154.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_R() view returns(bytes32)
func (_Polygonzkevm *PolygonzkevmCaller) SIGNATUREINITIALIZETXR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "SIGNATURE_INITIALIZE_TX_R")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SIGNATUREINITIALIZETXR is a free data retrieval call binding the contract method 0xb0afe154.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_R() view returns(bytes32)
func (_Polygonzkevm *PolygonzkevmSession) SIGNATUREINITIALIZETXR() ([32]byte, error) {
	return _Polygonzkevm.Contract.SIGNATUREINITIALIZETXR(&_Polygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXR is a free data retrieval call binding the contract method 0xb0afe154.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_R() view returns(bytes32)
func (_Polygonzkevm *PolygonzkevmCallerSession) SIGNATUREINITIALIZETXR() ([32]byte, error) {
	return _Polygonzkevm.Contract.SIGNATUREINITIALIZETXR(&_Polygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXS is a free data retrieval call binding the contract method 0xd7bc90ff.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_S() view returns(bytes32)
func (_Polygonzkevm *PolygonzkevmCaller) SIGNATUREINITIALIZETXS(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "SIGNATURE_INITIALIZE_TX_S")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SIGNATUREINITIALIZETXS is a free data retrieval call binding the contract method 0xd7bc90ff.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_S() view returns(bytes32)
func (_Polygonzkevm *PolygonzkevmSession) SIGNATUREINITIALIZETXS() ([32]byte, error) {
	return _Polygonzkevm.Contract.SIGNATUREINITIALIZETXS(&_Polygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXS is a free data retrieval call binding the contract method 0xd7bc90ff.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_S() view returns(bytes32)
func (_Polygonzkevm *PolygonzkevmCallerSession) SIGNATUREINITIALIZETXS() ([32]byte, error) {
	return _Polygonzkevm.Contract.SIGNATUREINITIALIZETXS(&_Polygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXV is a free data retrieval call binding the contract method 0xf35dda47.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_V() view returns(uint8)
func (_Polygonzkevm *PolygonzkevmCaller) SIGNATUREINITIALIZETXV(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "SIGNATURE_INITIALIZE_TX_V")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// SIGNATUREINITIALIZETXV is a free data retrieval call binding the contract method 0xf35dda47.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_V() view returns(uint8)
func (_Polygonzkevm *PolygonzkevmSession) SIGNATUREINITIALIZETXV() (uint8, error) {
	return _Polygonzkevm.Contract.SIGNATUREINITIALIZETXV(&_Polygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXV is a free data retrieval call binding the contract method 0xf35dda47.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_V() view returns(uint8)
func (_Polygonzkevm *PolygonzkevmCallerSession) SIGNATUREINITIALIZETXV() (uint8, error) {
	return _Polygonzkevm.Contract.SIGNATUREINITIALIZETXV(&_Polygonzkevm.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Polygonzkevm *PolygonzkevmCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Polygonzkevm *PolygonzkevmSession) Admin() (common.Address, error) {
	return _Polygonzkevm.Contract.Admin(&_Polygonzkevm.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Polygonzkevm *PolygonzkevmCallerSession) Admin() (common.Address, error) {
	return _Polygonzkevm.Contract.Admin(&_Polygonzkevm.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Polygonzkevm *PolygonzkevmCaller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Polygonzkevm *PolygonzkevmSession) BridgeAddress() (common.Address, error) {
	return _Polygonzkevm.Contract.BridgeAddress(&_Polygonzkevm.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Polygonzkevm *PolygonzkevmCallerSession) BridgeAddress() (common.Address, error) {
	return _Polygonzkevm.Contract.BridgeAddress(&_Polygonzkevm.CallOpts)
}

// CalculatePolPerForceBatch is a free data retrieval call binding the contract method 0x00d0295d.
//
// Solidity: function calculatePolPerForceBatch() view returns(uint256)
func (_Polygonzkevm *PolygonzkevmCaller) CalculatePolPerForceBatch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "calculatePolPerForceBatch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculatePolPerForceBatch is a free data retrieval call binding the contract method 0x00d0295d.
//
// Solidity: function calculatePolPerForceBatch() view returns(uint256)
func (_Polygonzkevm *PolygonzkevmSession) CalculatePolPerForceBatch() (*big.Int, error) {
	return _Polygonzkevm.Contract.CalculatePolPerForceBatch(&_Polygonzkevm.CallOpts)
}

// CalculatePolPerForceBatch is a free data retrieval call binding the contract method 0x00d0295d.
//
// Solidity: function calculatePolPerForceBatch() view returns(uint256)
func (_Polygonzkevm *PolygonzkevmCallerSession) CalculatePolPerForceBatch() (*big.Int, error) {
	return _Polygonzkevm.Contract.CalculatePolPerForceBatch(&_Polygonzkevm.CallOpts)
}

// ForceBatchTimeout is a free data retrieval call binding the contract method 0xc754c7ed.
//
// Solidity: function forceBatchTimeout() view returns(uint64)
func (_Polygonzkevm *PolygonzkevmCaller) ForceBatchTimeout(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "forceBatchTimeout")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ForceBatchTimeout is a free data retrieval call binding the contract method 0xc754c7ed.
//
// Solidity: function forceBatchTimeout() view returns(uint64)
func (_Polygonzkevm *PolygonzkevmSession) ForceBatchTimeout() (uint64, error) {
	return _Polygonzkevm.Contract.ForceBatchTimeout(&_Polygonzkevm.CallOpts)
}

// ForceBatchTimeout is a free data retrieval call binding the contract method 0xc754c7ed.
//
// Solidity: function forceBatchTimeout() view returns(uint64)
func (_Polygonzkevm *PolygonzkevmCallerSession) ForceBatchTimeout() (uint64, error) {
	return _Polygonzkevm.Contract.ForceBatchTimeout(&_Polygonzkevm.CallOpts)
}

// ForcedBatches is a free data retrieval call binding the contract method 0x6b8616ce.
//
// Solidity: function forcedBatches(uint64 ) view returns(bytes32)
func (_Polygonzkevm *PolygonzkevmCaller) ForcedBatches(opts *bind.CallOpts, arg0 uint64) ([32]byte, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "forcedBatches", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ForcedBatches is a free data retrieval call binding the contract method 0x6b8616ce.
//
// Solidity: function forcedBatches(uint64 ) view returns(bytes32)
func (_Polygonzkevm *PolygonzkevmSession) ForcedBatches(arg0 uint64) ([32]byte, error) {
	return _Polygonzkevm.Contract.ForcedBatches(&_Polygonzkevm.CallOpts, arg0)
}

// ForcedBatches is a free data retrieval call binding the contract method 0x6b8616ce.
//
// Solidity: function forcedBatches(uint64 ) view returns(bytes32)
func (_Polygonzkevm *PolygonzkevmCallerSession) ForcedBatches(arg0 uint64) ([32]byte, error) {
	return _Polygonzkevm.Contract.ForcedBatches(&_Polygonzkevm.CallOpts, arg0)
}

// GapLastTimestamp is a free data retrieval call binding the contract method 0x2a9f208f.
//
// Solidity: function gapLastTimestamp() view returns(uint64)
func (_Polygonzkevm *PolygonzkevmCaller) GapLastTimestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "gapLastTimestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GapLastTimestamp is a free data retrieval call binding the contract method 0x2a9f208f.
//
// Solidity: function gapLastTimestamp() view returns(uint64)
func (_Polygonzkevm *PolygonzkevmSession) GapLastTimestamp() (uint64, error) {
	return _Polygonzkevm.Contract.GapLastTimestamp(&_Polygonzkevm.CallOpts)
}

// GapLastTimestamp is a free data retrieval call binding the contract method 0x2a9f208f.
//
// Solidity: function gapLastTimestamp() view returns(uint64)
func (_Polygonzkevm *PolygonzkevmCallerSession) GapLastTimestamp() (uint64, error) {
	return _Polygonzkevm.Contract.GapLastTimestamp(&_Polygonzkevm.CallOpts)
}

// GasTokenAddress is a free data retrieval call binding the contract method 0x3c351e10.
//
// Solidity: function gasTokenAddress() view returns(address)
func (_Polygonzkevm *PolygonzkevmCaller) GasTokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "gasTokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GasTokenAddress is a free data retrieval call binding the contract method 0x3c351e10.
//
// Solidity: function gasTokenAddress() view returns(address)
func (_Polygonzkevm *PolygonzkevmSession) GasTokenAddress() (common.Address, error) {
	return _Polygonzkevm.Contract.GasTokenAddress(&_Polygonzkevm.CallOpts)
}

// GasTokenAddress is a free data retrieval call binding the contract method 0x3c351e10.
//
// Solidity: function gasTokenAddress() view returns(address)
func (_Polygonzkevm *PolygonzkevmCallerSession) GasTokenAddress() (common.Address, error) {
	return _Polygonzkevm.Contract.GasTokenAddress(&_Polygonzkevm.CallOpts)
}

// GasTokenNetwork is a free data retrieval call binding the contract method 0x3cbc795b.
//
// Solidity: function gasTokenNetwork() view returns(uint32)
func (_Polygonzkevm *PolygonzkevmCaller) GasTokenNetwork(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "gasTokenNetwork")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GasTokenNetwork is a free data retrieval call binding the contract method 0x3cbc795b.
//
// Solidity: function gasTokenNetwork() view returns(uint32)
func (_Polygonzkevm *PolygonzkevmSession) GasTokenNetwork() (uint32, error) {
	return _Polygonzkevm.Contract.GasTokenNetwork(&_Polygonzkevm.CallOpts)
}

// GasTokenNetwork is a free data retrieval call binding the contract method 0x3cbc795b.
//
// Solidity: function gasTokenNetwork() view returns(uint32)
func (_Polygonzkevm *PolygonzkevmCallerSession) GasTokenNetwork() (uint32, error) {
	return _Polygonzkevm.Contract.GasTokenNetwork(&_Polygonzkevm.CallOpts)
}

// GenerateInitializeTransaction is a free data retrieval call binding the contract method 0xa652f26c.
//
// Solidity: function generateInitializeTransaction(uint32 networkID, address _gasTokenAddress, uint32 _gasTokenNetwork, bytes _gasTokenMetadata) view returns(bytes)
func (_Polygonzkevm *PolygonzkevmCaller) GenerateInitializeTransaction(opts *bind.CallOpts, networkID uint32, _gasTokenAddress common.Address, _gasTokenNetwork uint32, _gasTokenMetadata []byte) ([]byte, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "generateInitializeTransaction", networkID, _gasTokenAddress, _gasTokenNetwork, _gasTokenMetadata)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GenerateInitializeTransaction is a free data retrieval call binding the contract method 0xa652f26c.
//
// Solidity: function generateInitializeTransaction(uint32 networkID, address _gasTokenAddress, uint32 _gasTokenNetwork, bytes _gasTokenMetadata) view returns(bytes)
func (_Polygonzkevm *PolygonzkevmSession) GenerateInitializeTransaction(networkID uint32, _gasTokenAddress common.Address, _gasTokenNetwork uint32, _gasTokenMetadata []byte) ([]byte, error) {
	return _Polygonzkevm.Contract.GenerateInitializeTransaction(&_Polygonzkevm.CallOpts, networkID, _gasTokenAddress, _gasTokenNetwork, _gasTokenMetadata)
}

// GenerateInitializeTransaction is a free data retrieval call binding the contract method 0xa652f26c.
//
// Solidity: function generateInitializeTransaction(uint32 networkID, address _gasTokenAddress, uint32 _gasTokenNetwork, bytes _gasTokenMetadata) view returns(bytes)
func (_Polygonzkevm *PolygonzkevmCallerSession) GenerateInitializeTransaction(networkID uint32, _gasTokenAddress common.Address, _gasTokenNetwork uint32, _gasTokenMetadata []byte) ([]byte, error) {
	return _Polygonzkevm.Contract.GenerateInitializeTransaction(&_Polygonzkevm.CallOpts, networkID, _gasTokenAddress, _gasTokenNetwork, _gasTokenMetadata)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Polygonzkevm *PolygonzkevmCaller) GlobalExitRootManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "globalExitRootManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Polygonzkevm *PolygonzkevmSession) GlobalExitRootManager() (common.Address, error) {
	return _Polygonzkevm.Contract.GlobalExitRootManager(&_Polygonzkevm.CallOpts)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Polygonzkevm *PolygonzkevmCallerSession) GlobalExitRootManager() (common.Address, error) {
	return _Polygonzkevm.Contract.GlobalExitRootManager(&_Polygonzkevm.CallOpts)
}

// IsForcedBatchAllowed is a free data retrieval call binding the contract method 0x2b35b3ac.
//
// Solidity: function isForcedBatchAllowed() view returns(bool)
func (_Polygonzkevm *PolygonzkevmCaller) IsForcedBatchAllowed(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "isForcedBatchAllowed")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsForcedBatchAllowed is a free data retrieval call binding the contract method 0x2b35b3ac.
//
// Solidity: function isForcedBatchAllowed() view returns(bool)
func (_Polygonzkevm *PolygonzkevmSession) IsForcedBatchAllowed() (bool, error) {
	return _Polygonzkevm.Contract.IsForcedBatchAllowed(&_Polygonzkevm.CallOpts)
}

// IsForcedBatchAllowed is a free data retrieval call binding the contract method 0x2b35b3ac.
//
// Solidity: function isForcedBatchAllowed() view returns(bool)
func (_Polygonzkevm *PolygonzkevmCallerSession) IsForcedBatchAllowed() (bool, error) {
	return _Polygonzkevm.Contract.IsForcedBatchAllowed(&_Polygonzkevm.CallOpts)
}

// LastAccInputHash is a free data retrieval call binding the contract method 0x6e05d2cd.
//
// Solidity: function lastAccInputHash() view returns(bytes32)
func (_Polygonzkevm *PolygonzkevmCaller) LastAccInputHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "lastAccInputHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastAccInputHash is a free data retrieval call binding the contract method 0x6e05d2cd.
//
// Solidity: function lastAccInputHash() view returns(bytes32)
func (_Polygonzkevm *PolygonzkevmSession) LastAccInputHash() ([32]byte, error) {
	return _Polygonzkevm.Contract.LastAccInputHash(&_Polygonzkevm.CallOpts)
}

// LastAccInputHash is a free data retrieval call binding the contract method 0x6e05d2cd.
//
// Solidity: function lastAccInputHash() view returns(bytes32)
func (_Polygonzkevm *PolygonzkevmCallerSession) LastAccInputHash() ([32]byte, error) {
	return _Polygonzkevm.Contract.LastAccInputHash(&_Polygonzkevm.CallOpts)
}

// LastForceBatch is a free data retrieval call binding the contract method 0xe7a7ed02.
//
// Solidity: function lastForceBatch() view returns(uint64)
func (_Polygonzkevm *PolygonzkevmCaller) LastForceBatch(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "lastForceBatch")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastForceBatch is a free data retrieval call binding the contract method 0xe7a7ed02.
//
// Solidity: function lastForceBatch() view returns(uint64)
func (_Polygonzkevm *PolygonzkevmSession) LastForceBatch() (uint64, error) {
	return _Polygonzkevm.Contract.LastForceBatch(&_Polygonzkevm.CallOpts)
}

// LastForceBatch is a free data retrieval call binding the contract method 0xe7a7ed02.
//
// Solidity: function lastForceBatch() view returns(uint64)
func (_Polygonzkevm *PolygonzkevmCallerSession) LastForceBatch() (uint64, error) {
	return _Polygonzkevm.Contract.LastForceBatch(&_Polygonzkevm.CallOpts)
}

// LastForceBatchSequenced is a free data retrieval call binding the contract method 0x45605267.
//
// Solidity: function lastForceBatchSequenced() view returns(uint64)
func (_Polygonzkevm *PolygonzkevmCaller) LastForceBatchSequenced(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "lastForceBatchSequenced")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastForceBatchSequenced is a free data retrieval call binding the contract method 0x45605267.
//
// Solidity: function lastForceBatchSequenced() view returns(uint64)
func (_Polygonzkevm *PolygonzkevmSession) LastForceBatchSequenced() (uint64, error) {
	return _Polygonzkevm.Contract.LastForceBatchSequenced(&_Polygonzkevm.CallOpts)
}

// LastForceBatchSequenced is a free data retrieval call binding the contract method 0x45605267.
//
// Solidity: function lastForceBatchSequenced() view returns(uint64)
func (_Polygonzkevm *PolygonzkevmCallerSession) LastForceBatchSequenced() (uint64, error) {
	return _Polygonzkevm.Contract.LastForceBatchSequenced(&_Polygonzkevm.CallOpts)
}

// NetworkName is a free data retrieval call binding the contract method 0x107bf28c.
//
// Solidity: function networkName() view returns(string)
func (_Polygonzkevm *PolygonzkevmCaller) NetworkName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "networkName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NetworkName is a free data retrieval call binding the contract method 0x107bf28c.
//
// Solidity: function networkName() view returns(string)
func (_Polygonzkevm *PolygonzkevmSession) NetworkName() (string, error) {
	return _Polygonzkevm.Contract.NetworkName(&_Polygonzkevm.CallOpts)
}

// NetworkName is a free data retrieval call binding the contract method 0x107bf28c.
//
// Solidity: function networkName() view returns(string)
func (_Polygonzkevm *PolygonzkevmCallerSession) NetworkName() (string, error) {
	return _Polygonzkevm.Contract.NetworkName(&_Polygonzkevm.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_Polygonzkevm *PolygonzkevmCaller) PendingAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "pendingAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_Polygonzkevm *PolygonzkevmSession) PendingAdmin() (common.Address, error) {
	return _Polygonzkevm.Contract.PendingAdmin(&_Polygonzkevm.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_Polygonzkevm *PolygonzkevmCallerSession) PendingAdmin() (common.Address, error) {
	return _Polygonzkevm.Contract.PendingAdmin(&_Polygonzkevm.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Polygonzkevm *PolygonzkevmCaller) Pol(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "pol")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Polygonzkevm *PolygonzkevmSession) Pol() (common.Address, error) {
	return _Polygonzkevm.Contract.Pol(&_Polygonzkevm.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Polygonzkevm *PolygonzkevmCallerSession) Pol() (common.Address, error) {
	return _Polygonzkevm.Contract.Pol(&_Polygonzkevm.CallOpts)
}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Polygonzkevm *PolygonzkevmCaller) RollupManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "rollupManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Polygonzkevm *PolygonzkevmSession) RollupManager() (common.Address, error) {
	return _Polygonzkevm.Contract.RollupManager(&_Polygonzkevm.CallOpts)
}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Polygonzkevm *PolygonzkevmCallerSession) RollupManager() (common.Address, error) {
	return _Polygonzkevm.Contract.RollupManager(&_Polygonzkevm.CallOpts)
}

// TrustedSequencer is a free data retrieval call binding the contract method 0xcfa8ed47.
//
// Solidity: function trustedSequencer() view returns(address)
func (_Polygonzkevm *PolygonzkevmCaller) TrustedSequencer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "trustedSequencer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TrustedSequencer is a free data retrieval call binding the contract method 0xcfa8ed47.
//
// Solidity: function trustedSequencer() view returns(address)
func (_Polygonzkevm *PolygonzkevmSession) TrustedSequencer() (common.Address, error) {
	return _Polygonzkevm.Contract.TrustedSequencer(&_Polygonzkevm.CallOpts)
}

// TrustedSequencer is a free data retrieval call binding the contract method 0xcfa8ed47.
//
// Solidity: function trustedSequencer() view returns(address)
func (_Polygonzkevm *PolygonzkevmCallerSession) TrustedSequencer() (common.Address, error) {
	return _Polygonzkevm.Contract.TrustedSequencer(&_Polygonzkevm.CallOpts)
}

// TrustedSequencerURL is a free data retrieval call binding the contract method 0x542028d5.
//
// Solidity: function trustedSequencerURL() view returns(string)
func (_Polygonzkevm *PolygonzkevmCaller) TrustedSequencerURL(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Polygonzkevm.contract.Call(opts, &out, "trustedSequencerURL")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TrustedSequencerURL is a free data retrieval call binding the contract method 0x542028d5.
//
// Solidity: function trustedSequencerURL() view returns(string)
func (_Polygonzkevm *PolygonzkevmSession) TrustedSequencerURL() (string, error) {
	return _Polygonzkevm.Contract.TrustedSequencerURL(&_Polygonzkevm.CallOpts)
}

// TrustedSequencerURL is a free data retrieval call binding the contract method 0x542028d5.
//
// Solidity: function trustedSequencerURL() view returns(string)
func (_Polygonzkevm *PolygonzkevmCallerSession) TrustedSequencerURL() (string, error) {
	return _Polygonzkevm.Contract.TrustedSequencerURL(&_Polygonzkevm.CallOpts)
}

// AcceptAdminRole is a paid mutator transaction binding the contract method 0x8c3d7301.
//
// Solidity: function acceptAdminRole() returns()
func (_Polygonzkevm *PolygonzkevmTransactor) AcceptAdminRole(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygonzkevm.contract.Transact(opts, "acceptAdminRole")
}

// AcceptAdminRole is a paid mutator transaction binding the contract method 0x8c3d7301.
//
// Solidity: function acceptAdminRole() returns()
func (_Polygonzkevm *PolygonzkevmSession) AcceptAdminRole() (*types.Transaction, error) {
	return _Polygonzkevm.Contract.AcceptAdminRole(&_Polygonzkevm.TransactOpts)
}

// AcceptAdminRole is a paid mutator transaction binding the contract method 0x8c3d7301.
//
// Solidity: function acceptAdminRole() returns()
func (_Polygonzkevm *PolygonzkevmTransactorSession) AcceptAdminRole() (*types.Transaction, error) {
	return _Polygonzkevm.Contract.AcceptAdminRole(&_Polygonzkevm.TransactOpts)
}

// ActivateForceBatches is a paid mutator transaction binding the contract method 0x5ec91958.
//
// Solidity: function activateForceBatches() returns()
func (_Polygonzkevm *PolygonzkevmTransactor) ActivateForceBatches(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygonzkevm.contract.Transact(opts, "activateForceBatches")
}

// ActivateForceBatches is a paid mutator transaction binding the contract method 0x5ec91958.
//
// Solidity: function activateForceBatches() returns()
func (_Polygonzkevm *PolygonzkevmSession) ActivateForceBatches() (*types.Transaction, error) {
	return _Polygonzkevm.Contract.ActivateForceBatches(&_Polygonzkevm.TransactOpts)
}

// ActivateForceBatches is a paid mutator transaction binding the contract method 0x5ec91958.
//
// Solidity: function activateForceBatches() returns()
func (_Polygonzkevm *PolygonzkevmTransactorSession) ActivateForceBatches() (*types.Transaction, error) {
	return _Polygonzkevm.Contract.ActivateForceBatches(&_Polygonzkevm.TransactOpts)
}

// ForceBatch is a paid mutator transaction binding the contract method 0xeaeb077b.
//
// Solidity: function forceBatch(bytes transactions, uint256 polAmount) returns()
func (_Polygonzkevm *PolygonzkevmTransactor) ForceBatch(opts *bind.TransactOpts, transactions []byte, polAmount *big.Int) (*types.Transaction, error) {
	return _Polygonzkevm.contract.Transact(opts, "forceBatch", transactions, polAmount)
}

// ForceBatch is a paid mutator transaction binding the contract method 0xeaeb077b.
//
// Solidity: function forceBatch(bytes transactions, uint256 polAmount) returns()
func (_Polygonzkevm *PolygonzkevmSession) ForceBatch(transactions []byte, polAmount *big.Int) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.ForceBatch(&_Polygonzkevm.TransactOpts, transactions, polAmount)
}

// ForceBatch is a paid mutator transaction binding the contract method 0xeaeb077b.
//
// Solidity: function forceBatch(bytes transactions, uint256 polAmount) returns()
func (_Polygonzkevm *PolygonzkevmTransactorSession) ForceBatch(transactions []byte, polAmount *big.Int) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.ForceBatch(&_Polygonzkevm.TransactOpts, transactions, polAmount)
}

// Initialize is a paid mutator transaction binding the contract method 0x71257022.
//
// Solidity: function initialize(address _admin, address sequencer, uint32 networkID, address _gasTokenAddress, string sequencerURL, string _networkName) returns()
func (_Polygonzkevm *PolygonzkevmTransactor) Initialize(opts *bind.TransactOpts, _admin common.Address, sequencer common.Address, networkID uint32, _gasTokenAddress common.Address, sequencerURL string, _networkName string) (*types.Transaction, error) {
	return _Polygonzkevm.contract.Transact(opts, "initialize", _admin, sequencer, networkID, _gasTokenAddress, sequencerURL, _networkName)
}

// Initialize is a paid mutator transaction binding the contract method 0x71257022.
//
// Solidity: function initialize(address _admin, address sequencer, uint32 networkID, address _gasTokenAddress, string sequencerURL, string _networkName) returns()
func (_Polygonzkevm *PolygonzkevmSession) Initialize(_admin common.Address, sequencer common.Address, networkID uint32, _gasTokenAddress common.Address, sequencerURL string, _networkName string) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.Initialize(&_Polygonzkevm.TransactOpts, _admin, sequencer, networkID, _gasTokenAddress, sequencerURL, _networkName)
}

// Initialize is a paid mutator transaction binding the contract method 0x71257022.
//
// Solidity: function initialize(address _admin, address sequencer, uint32 networkID, address _gasTokenAddress, string sequencerURL, string _networkName) returns()
func (_Polygonzkevm *PolygonzkevmTransactorSession) Initialize(_admin common.Address, sequencer common.Address, networkID uint32, _gasTokenAddress common.Address, sequencerURL string, _networkName string) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.Initialize(&_Polygonzkevm.TransactOpts, _admin, sequencer, networkID, _gasTokenAddress, sequencerURL, _networkName)
}

// OnVerifyBatches is a paid mutator transaction binding the contract method 0x32c2d153.
//
// Solidity: function onVerifyBatches(uint64 lastVerifiedBatch, bytes32 newStateRoot, address aggregator) returns()
func (_Polygonzkevm *PolygonzkevmTransactor) OnVerifyBatches(opts *bind.TransactOpts, lastVerifiedBatch uint64, newStateRoot [32]byte, aggregator common.Address) (*types.Transaction, error) {
	return _Polygonzkevm.contract.Transact(opts, "onVerifyBatches", lastVerifiedBatch, newStateRoot, aggregator)
}

// OnVerifyBatches is a paid mutator transaction binding the contract method 0x32c2d153.
//
// Solidity: function onVerifyBatches(uint64 lastVerifiedBatch, bytes32 newStateRoot, address aggregator) returns()
func (_Polygonzkevm *PolygonzkevmSession) OnVerifyBatches(lastVerifiedBatch uint64, newStateRoot [32]byte, aggregator common.Address) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.OnVerifyBatches(&_Polygonzkevm.TransactOpts, lastVerifiedBatch, newStateRoot, aggregator)
}

// OnVerifyBatches is a paid mutator transaction binding the contract method 0x32c2d153.
//
// Solidity: function onVerifyBatches(uint64 lastVerifiedBatch, bytes32 newStateRoot, address aggregator) returns()
func (_Polygonzkevm *PolygonzkevmTransactorSession) OnVerifyBatches(lastVerifiedBatch uint64, newStateRoot [32]byte, aggregator common.Address) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.OnVerifyBatches(&_Polygonzkevm.TransactOpts, lastVerifiedBatch, newStateRoot, aggregator)
}

// SequenceBatches is a paid mutator transaction binding the contract method 0xecef3f99.
//
// Solidity: function sequenceBatches((bytes,bytes32,uint64,bytes32)[] batches, address l2Coinbase) returns()
func (_Polygonzkevm *PolygonzkevmTransactor) SequenceBatches(opts *bind.TransactOpts, batches []PolygonRollupBaseEtrogBatchData, l2Coinbase common.Address) (*types.Transaction, error) {
	return _Polygonzkevm.contract.Transact(opts, "sequenceBatches", batches, l2Coinbase)
}

// SequenceBatches is a paid mutator transaction binding the contract method 0xecef3f99.
//
// Solidity: function sequenceBatches((bytes,bytes32,uint64,bytes32)[] batches, address l2Coinbase) returns()
func (_Polygonzkevm *PolygonzkevmSession) SequenceBatches(batches []PolygonRollupBaseEtrogBatchData, l2Coinbase common.Address) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.SequenceBatches(&_Polygonzkevm.TransactOpts, batches, l2Coinbase)
}

// SequenceBatches is a paid mutator transaction binding the contract method 0xecef3f99.
//
// Solidity: function sequenceBatches((bytes,bytes32,uint64,bytes32)[] batches, address l2Coinbase) returns()
func (_Polygonzkevm *PolygonzkevmTransactorSession) SequenceBatches(batches []PolygonRollupBaseEtrogBatchData, l2Coinbase common.Address) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.SequenceBatches(&_Polygonzkevm.TransactOpts, batches, l2Coinbase)
}

// SequenceForceBatches is a paid mutator transaction binding the contract method 0x9f26f840.
//
// Solidity: function sequenceForceBatches((bytes,bytes32,uint64,bytes32)[] batches) returns()
func (_Polygonzkevm *PolygonzkevmTransactor) SequenceForceBatches(opts *bind.TransactOpts, batches []PolygonRollupBaseEtrogBatchData) (*types.Transaction, error) {
	return _Polygonzkevm.contract.Transact(opts, "sequenceForceBatches", batches)
}

// SequenceForceBatches is a paid mutator transaction binding the contract method 0x9f26f840.
//
// Solidity: function sequenceForceBatches((bytes,bytes32,uint64,bytes32)[] batches) returns()
func (_Polygonzkevm *PolygonzkevmSession) SequenceForceBatches(batches []PolygonRollupBaseEtrogBatchData) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.SequenceForceBatches(&_Polygonzkevm.TransactOpts, batches)
}

// SequenceForceBatches is a paid mutator transaction binding the contract method 0x9f26f840.
//
// Solidity: function sequenceForceBatches((bytes,bytes32,uint64,bytes32)[] batches) returns()
func (_Polygonzkevm *PolygonzkevmTransactorSession) SequenceForceBatches(batches []PolygonRollupBaseEtrogBatchData) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.SequenceForceBatches(&_Polygonzkevm.TransactOpts, batches)
}

// SetForceBatchTimeout is a paid mutator transaction binding the contract method 0x4e487706.
//
// Solidity: function setForceBatchTimeout(uint64 newforceBatchTimeout) returns()
func (_Polygonzkevm *PolygonzkevmTransactor) SetForceBatchTimeout(opts *bind.TransactOpts, newforceBatchTimeout uint64) (*types.Transaction, error) {
	return _Polygonzkevm.contract.Transact(opts, "setForceBatchTimeout", newforceBatchTimeout)
}

// SetForceBatchTimeout is a paid mutator transaction binding the contract method 0x4e487706.
//
// Solidity: function setForceBatchTimeout(uint64 newforceBatchTimeout) returns()
func (_Polygonzkevm *PolygonzkevmSession) SetForceBatchTimeout(newforceBatchTimeout uint64) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.SetForceBatchTimeout(&_Polygonzkevm.TransactOpts, newforceBatchTimeout)
}

// SetForceBatchTimeout is a paid mutator transaction binding the contract method 0x4e487706.
//
// Solidity: function setForceBatchTimeout(uint64 newforceBatchTimeout) returns()
func (_Polygonzkevm *PolygonzkevmTransactorSession) SetForceBatchTimeout(newforceBatchTimeout uint64) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.SetForceBatchTimeout(&_Polygonzkevm.TransactOpts, newforceBatchTimeout)
}

// SetTrustedSequencer is a paid mutator transaction binding the contract method 0x6ff512cc.
//
// Solidity: function setTrustedSequencer(address newTrustedSequencer) returns()
func (_Polygonzkevm *PolygonzkevmTransactor) SetTrustedSequencer(opts *bind.TransactOpts, newTrustedSequencer common.Address) (*types.Transaction, error) {
	return _Polygonzkevm.contract.Transact(opts, "setTrustedSequencer", newTrustedSequencer)
}

// SetTrustedSequencer is a paid mutator transaction binding the contract method 0x6ff512cc.
//
// Solidity: function setTrustedSequencer(address newTrustedSequencer) returns()
func (_Polygonzkevm *PolygonzkevmSession) SetTrustedSequencer(newTrustedSequencer common.Address) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.SetTrustedSequencer(&_Polygonzkevm.TransactOpts, newTrustedSequencer)
}

// SetTrustedSequencer is a paid mutator transaction binding the contract method 0x6ff512cc.
//
// Solidity: function setTrustedSequencer(address newTrustedSequencer) returns()
func (_Polygonzkevm *PolygonzkevmTransactorSession) SetTrustedSequencer(newTrustedSequencer common.Address) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.SetTrustedSequencer(&_Polygonzkevm.TransactOpts, newTrustedSequencer)
}

// SetTrustedSequencerURL is a paid mutator transaction binding the contract method 0xc89e42df.
//
// Solidity: function setTrustedSequencerURL(string newTrustedSequencerURL) returns()
func (_Polygonzkevm *PolygonzkevmTransactor) SetTrustedSequencerURL(opts *bind.TransactOpts, newTrustedSequencerURL string) (*types.Transaction, error) {
	return _Polygonzkevm.contract.Transact(opts, "setTrustedSequencerURL", newTrustedSequencerURL)
}

// SetTrustedSequencerURL is a paid mutator transaction binding the contract method 0xc89e42df.
//
// Solidity: function setTrustedSequencerURL(string newTrustedSequencerURL) returns()
func (_Polygonzkevm *PolygonzkevmSession) SetTrustedSequencerURL(newTrustedSequencerURL string) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.SetTrustedSequencerURL(&_Polygonzkevm.TransactOpts, newTrustedSequencerURL)
}

// SetTrustedSequencerURL is a paid mutator transaction binding the contract method 0xc89e42df.
//
// Solidity: function setTrustedSequencerURL(string newTrustedSequencerURL) returns()
func (_Polygonzkevm *PolygonzkevmTransactorSession) SetTrustedSequencerURL(newTrustedSequencerURL string) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.SetTrustedSequencerURL(&_Polygonzkevm.TransactOpts, newTrustedSequencerURL)
}

// TransferAdminRole is a paid mutator transaction binding the contract method 0xada8f919.
//
// Solidity: function transferAdminRole(address newPendingAdmin) returns()
func (_Polygonzkevm *PolygonzkevmTransactor) TransferAdminRole(opts *bind.TransactOpts, newPendingAdmin common.Address) (*types.Transaction, error) {
	return _Polygonzkevm.contract.Transact(opts, "transferAdminRole", newPendingAdmin)
}

// TransferAdminRole is a paid mutator transaction binding the contract method 0xada8f919.
//
// Solidity: function transferAdminRole(address newPendingAdmin) returns()
func (_Polygonzkevm *PolygonzkevmSession) TransferAdminRole(newPendingAdmin common.Address) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.TransferAdminRole(&_Polygonzkevm.TransactOpts, newPendingAdmin)
}

// TransferAdminRole is a paid mutator transaction binding the contract method 0xada8f919.
//
// Solidity: function transferAdminRole(address newPendingAdmin) returns()
func (_Polygonzkevm *PolygonzkevmTransactorSession) TransferAdminRole(newPendingAdmin common.Address) (*types.Transaction, error) {
	return _Polygonzkevm.Contract.TransferAdminRole(&_Polygonzkevm.TransactOpts, newPendingAdmin)
}

// PolygonzkevmAcceptAdminRoleIterator is returned from FilterAcceptAdminRole and is used to iterate over the raw logs and unpacked data for AcceptAdminRole events raised by the Polygonzkevm contract.
type PolygonzkevmAcceptAdminRoleIterator struct {
	Event *PolygonzkevmAcceptAdminRole // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmAcceptAdminRoleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmAcceptAdminRole)
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
		it.Event = new(PolygonzkevmAcceptAdminRole)
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
func (it *PolygonzkevmAcceptAdminRoleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmAcceptAdminRoleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmAcceptAdminRole represents a AcceptAdminRole event raised by the Polygonzkevm contract.
type PolygonzkevmAcceptAdminRole struct {
	NewAdmin common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAcceptAdminRole is a free log retrieval operation binding the contract event 0x056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e.
//
// Solidity: event AcceptAdminRole(address newAdmin)
func (_Polygonzkevm *PolygonzkevmFilterer) FilterAcceptAdminRole(opts *bind.FilterOpts) (*PolygonzkevmAcceptAdminRoleIterator, error) {

	logs, sub, err := _Polygonzkevm.contract.FilterLogs(opts, "AcceptAdminRole")
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmAcceptAdminRoleIterator{contract: _Polygonzkevm.contract, event: "AcceptAdminRole", logs: logs, sub: sub}, nil
}

// WatchAcceptAdminRole is a free log subscription operation binding the contract event 0x056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e.
//
// Solidity: event AcceptAdminRole(address newAdmin)
func (_Polygonzkevm *PolygonzkevmFilterer) WatchAcceptAdminRole(opts *bind.WatchOpts, sink chan<- *PolygonzkevmAcceptAdminRole) (event.Subscription, error) {

	logs, sub, err := _Polygonzkevm.contract.WatchLogs(opts, "AcceptAdminRole")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmAcceptAdminRole)
				if err := _Polygonzkevm.contract.UnpackLog(event, "AcceptAdminRole", log); err != nil {
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

// ParseAcceptAdminRole is a log parse operation binding the contract event 0x056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e.
//
// Solidity: event AcceptAdminRole(address newAdmin)
func (_Polygonzkevm *PolygonzkevmFilterer) ParseAcceptAdminRole(log types.Log) (*PolygonzkevmAcceptAdminRole, error) {
	event := new(PolygonzkevmAcceptAdminRole)
	if err := _Polygonzkevm.contract.UnpackLog(event, "AcceptAdminRole", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonzkevmActivateForceBatchesIterator is returned from FilterActivateForceBatches and is used to iterate over the raw logs and unpacked data for ActivateForceBatches events raised by the Polygonzkevm contract.
type PolygonzkevmActivateForceBatchesIterator struct {
	Event *PolygonzkevmActivateForceBatches // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmActivateForceBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmActivateForceBatches)
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
		it.Event = new(PolygonzkevmActivateForceBatches)
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
func (it *PolygonzkevmActivateForceBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmActivateForceBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmActivateForceBatches represents a ActivateForceBatches event raised by the Polygonzkevm contract.
type PolygonzkevmActivateForceBatches struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterActivateForceBatches is a free log retrieval operation binding the contract event 0x854dd6ce5a1445c4c54388b21cffd11cf5bba1b9e763aec48ce3da75d617412f.
//
// Solidity: event ActivateForceBatches()
func (_Polygonzkevm *PolygonzkevmFilterer) FilterActivateForceBatches(opts *bind.FilterOpts) (*PolygonzkevmActivateForceBatchesIterator, error) {

	logs, sub, err := _Polygonzkevm.contract.FilterLogs(opts, "ActivateForceBatches")
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmActivateForceBatchesIterator{contract: _Polygonzkevm.contract, event: "ActivateForceBatches", logs: logs, sub: sub}, nil
}

// WatchActivateForceBatches is a free log subscription operation binding the contract event 0x854dd6ce5a1445c4c54388b21cffd11cf5bba1b9e763aec48ce3da75d617412f.
//
// Solidity: event ActivateForceBatches()
func (_Polygonzkevm *PolygonzkevmFilterer) WatchActivateForceBatches(opts *bind.WatchOpts, sink chan<- *PolygonzkevmActivateForceBatches) (event.Subscription, error) {

	logs, sub, err := _Polygonzkevm.contract.WatchLogs(opts, "ActivateForceBatches")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmActivateForceBatches)
				if err := _Polygonzkevm.contract.UnpackLog(event, "ActivateForceBatches", log); err != nil {
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

// ParseActivateForceBatches is a log parse operation binding the contract event 0x854dd6ce5a1445c4c54388b21cffd11cf5bba1b9e763aec48ce3da75d617412f.
//
// Solidity: event ActivateForceBatches()
func (_Polygonzkevm *PolygonzkevmFilterer) ParseActivateForceBatches(log types.Log) (*PolygonzkevmActivateForceBatches, error) {
	event := new(PolygonzkevmActivateForceBatches)
	if err := _Polygonzkevm.contract.UnpackLog(event, "ActivateForceBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonzkevmForceBatchIterator is returned from FilterForceBatch and is used to iterate over the raw logs and unpacked data for ForceBatch events raised by the Polygonzkevm contract.
type PolygonzkevmForceBatchIterator struct {
	Event *PolygonzkevmForceBatch // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmForceBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmForceBatch)
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
		it.Event = new(PolygonzkevmForceBatch)
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
func (it *PolygonzkevmForceBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmForceBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmForceBatch represents a ForceBatch event raised by the Polygonzkevm contract.
type PolygonzkevmForceBatch struct {
	ForceBatchNum      uint64
	LastGlobalExitRoot [32]byte
	Sequencer          common.Address
	Transactions       []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterForceBatch is a free log retrieval operation binding the contract event 0xf94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc931.
//
// Solidity: event ForceBatch(uint64 indexed forceBatchNum, bytes32 lastGlobalExitRoot, address sequencer, bytes transactions)
func (_Polygonzkevm *PolygonzkevmFilterer) FilterForceBatch(opts *bind.FilterOpts, forceBatchNum []uint64) (*PolygonzkevmForceBatchIterator, error) {

	var forceBatchNumRule []interface{}
	for _, forceBatchNumItem := range forceBatchNum {
		forceBatchNumRule = append(forceBatchNumRule, forceBatchNumItem)
	}

	logs, sub, err := _Polygonzkevm.contract.FilterLogs(opts, "ForceBatch", forceBatchNumRule)
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmForceBatchIterator{contract: _Polygonzkevm.contract, event: "ForceBatch", logs: logs, sub: sub}, nil
}

// WatchForceBatch is a free log subscription operation binding the contract event 0xf94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc931.
//
// Solidity: event ForceBatch(uint64 indexed forceBatchNum, bytes32 lastGlobalExitRoot, address sequencer, bytes transactions)
func (_Polygonzkevm *PolygonzkevmFilterer) WatchForceBatch(opts *bind.WatchOpts, sink chan<- *PolygonzkevmForceBatch, forceBatchNum []uint64) (event.Subscription, error) {

	var forceBatchNumRule []interface{}
	for _, forceBatchNumItem := range forceBatchNum {
		forceBatchNumRule = append(forceBatchNumRule, forceBatchNumItem)
	}

	logs, sub, err := _Polygonzkevm.contract.WatchLogs(opts, "ForceBatch", forceBatchNumRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmForceBatch)
				if err := _Polygonzkevm.contract.UnpackLog(event, "ForceBatch", log); err != nil {
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

// ParseForceBatch is a log parse operation binding the contract event 0xf94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc931.
//
// Solidity: event ForceBatch(uint64 indexed forceBatchNum, bytes32 lastGlobalExitRoot, address sequencer, bytes transactions)
func (_Polygonzkevm *PolygonzkevmFilterer) ParseForceBatch(log types.Log) (*PolygonzkevmForceBatch, error) {
	event := new(PolygonzkevmForceBatch)
	if err := _Polygonzkevm.contract.UnpackLog(event, "ForceBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonzkevmInitialSequenceBatchesIterator is returned from FilterInitialSequenceBatches and is used to iterate over the raw logs and unpacked data for InitialSequenceBatches events raised by the Polygonzkevm contract.
type PolygonzkevmInitialSequenceBatchesIterator struct {
	Event *PolygonzkevmInitialSequenceBatches // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmInitialSequenceBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmInitialSequenceBatches)
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
		it.Event = new(PolygonzkevmInitialSequenceBatches)
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
func (it *PolygonzkevmInitialSequenceBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmInitialSequenceBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmInitialSequenceBatches represents a InitialSequenceBatches event raised by the Polygonzkevm contract.
type PolygonzkevmInitialSequenceBatches struct {
	Transactions       []byte
	LastGlobalExitRoot [32]byte
	Sequencer          common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterInitialSequenceBatches is a free log retrieval operation binding the contract event 0x060116213bcbf54ca19fd649dc84b59ab2bbd200ab199770e4d923e222a28e7f.
//
// Solidity: event InitialSequenceBatches(bytes transactions, bytes32 lastGlobalExitRoot, address sequencer)
func (_Polygonzkevm *PolygonzkevmFilterer) FilterInitialSequenceBatches(opts *bind.FilterOpts) (*PolygonzkevmInitialSequenceBatchesIterator, error) {

	logs, sub, err := _Polygonzkevm.contract.FilterLogs(opts, "InitialSequenceBatches")
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmInitialSequenceBatchesIterator{contract: _Polygonzkevm.contract, event: "InitialSequenceBatches", logs: logs, sub: sub}, nil
}

// WatchInitialSequenceBatches is a free log subscription operation binding the contract event 0x060116213bcbf54ca19fd649dc84b59ab2bbd200ab199770e4d923e222a28e7f.
//
// Solidity: event InitialSequenceBatches(bytes transactions, bytes32 lastGlobalExitRoot, address sequencer)
func (_Polygonzkevm *PolygonzkevmFilterer) WatchInitialSequenceBatches(opts *bind.WatchOpts, sink chan<- *PolygonzkevmInitialSequenceBatches) (event.Subscription, error) {

	logs, sub, err := _Polygonzkevm.contract.WatchLogs(opts, "InitialSequenceBatches")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmInitialSequenceBatches)
				if err := _Polygonzkevm.contract.UnpackLog(event, "InitialSequenceBatches", log); err != nil {
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

// ParseInitialSequenceBatches is a log parse operation binding the contract event 0x060116213bcbf54ca19fd649dc84b59ab2bbd200ab199770e4d923e222a28e7f.
//
// Solidity: event InitialSequenceBatches(bytes transactions, bytes32 lastGlobalExitRoot, address sequencer)
func (_Polygonzkevm *PolygonzkevmFilterer) ParseInitialSequenceBatches(log types.Log) (*PolygonzkevmInitialSequenceBatches, error) {
	event := new(PolygonzkevmInitialSequenceBatches)
	if err := _Polygonzkevm.contract.UnpackLog(event, "InitialSequenceBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonzkevmInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Polygonzkevm contract.
type PolygonzkevmInitializedIterator struct {
	Event *PolygonzkevmInitialized // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmInitialized)
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
		it.Event = new(PolygonzkevmInitialized)
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
func (it *PolygonzkevmInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmInitialized represents a Initialized event raised by the Polygonzkevm contract.
type PolygonzkevmInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Polygonzkevm *PolygonzkevmFilterer) FilterInitialized(opts *bind.FilterOpts) (*PolygonzkevmInitializedIterator, error) {

	logs, sub, err := _Polygonzkevm.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmInitializedIterator{contract: _Polygonzkevm.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Polygonzkevm *PolygonzkevmFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PolygonzkevmInitialized) (event.Subscription, error) {

	logs, sub, err := _Polygonzkevm.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmInitialized)
				if err := _Polygonzkevm.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Polygonzkevm *PolygonzkevmFilterer) ParseInitialized(log types.Log) (*PolygonzkevmInitialized, error) {
	event := new(PolygonzkevmInitialized)
	if err := _Polygonzkevm.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonzkevmSequenceBatchesIterator is returned from FilterSequenceBatches and is used to iterate over the raw logs and unpacked data for SequenceBatches events raised by the Polygonzkevm contract.
type PolygonzkevmSequenceBatchesIterator struct {
	Event *PolygonzkevmSequenceBatches // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmSequenceBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmSequenceBatches)
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
		it.Event = new(PolygonzkevmSequenceBatches)
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
func (it *PolygonzkevmSequenceBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmSequenceBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmSequenceBatches represents a SequenceBatches event raised by the Polygonzkevm contract.
type PolygonzkevmSequenceBatches struct {
	NumBatch   uint64
	L1InfoRoot [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSequenceBatches is a free log retrieval operation binding the contract event 0x3e54d0825ed78523037d00a81759237eb436ce774bd546993ee67a1b67b6e766.
//
// Solidity: event SequenceBatches(uint64 indexed numBatch, bytes32 l1InfoRoot)
func (_Polygonzkevm *PolygonzkevmFilterer) FilterSequenceBatches(opts *bind.FilterOpts, numBatch []uint64) (*PolygonzkevmSequenceBatchesIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	logs, sub, err := _Polygonzkevm.contract.FilterLogs(opts, "SequenceBatches", numBatchRule)
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmSequenceBatchesIterator{contract: _Polygonzkevm.contract, event: "SequenceBatches", logs: logs, sub: sub}, nil
}

// WatchSequenceBatches is a free log subscription operation binding the contract event 0x3e54d0825ed78523037d00a81759237eb436ce774bd546993ee67a1b67b6e766.
//
// Solidity: event SequenceBatches(uint64 indexed numBatch, bytes32 l1InfoRoot)
func (_Polygonzkevm *PolygonzkevmFilterer) WatchSequenceBatches(opts *bind.WatchOpts, sink chan<- *PolygonzkevmSequenceBatches, numBatch []uint64) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	logs, sub, err := _Polygonzkevm.contract.WatchLogs(opts, "SequenceBatches", numBatchRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmSequenceBatches)
				if err := _Polygonzkevm.contract.UnpackLog(event, "SequenceBatches", log); err != nil {
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

// ParseSequenceBatches is a log parse operation binding the contract event 0x3e54d0825ed78523037d00a81759237eb436ce774bd546993ee67a1b67b6e766.
//
// Solidity: event SequenceBatches(uint64 indexed numBatch, bytes32 l1InfoRoot)
func (_Polygonzkevm *PolygonzkevmFilterer) ParseSequenceBatches(log types.Log) (*PolygonzkevmSequenceBatches, error) {
	event := new(PolygonzkevmSequenceBatches)
	if err := _Polygonzkevm.contract.UnpackLog(event, "SequenceBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonzkevmSequenceForceBatchesIterator is returned from FilterSequenceForceBatches and is used to iterate over the raw logs and unpacked data for SequenceForceBatches events raised by the Polygonzkevm contract.
type PolygonzkevmSequenceForceBatchesIterator struct {
	Event *PolygonzkevmSequenceForceBatches // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmSequenceForceBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmSequenceForceBatches)
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
		it.Event = new(PolygonzkevmSequenceForceBatches)
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
func (it *PolygonzkevmSequenceForceBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmSequenceForceBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmSequenceForceBatches represents a SequenceForceBatches event raised by the Polygonzkevm contract.
type PolygonzkevmSequenceForceBatches struct {
	NumBatch uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSequenceForceBatches is a free log retrieval operation binding the contract event 0x648a61dd2438f072f5a1960939abd30f37aea80d2e94c9792ad142d3e0a490a4.
//
// Solidity: event SequenceForceBatches(uint64 indexed numBatch)
func (_Polygonzkevm *PolygonzkevmFilterer) FilterSequenceForceBatches(opts *bind.FilterOpts, numBatch []uint64) (*PolygonzkevmSequenceForceBatchesIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	logs, sub, err := _Polygonzkevm.contract.FilterLogs(opts, "SequenceForceBatches", numBatchRule)
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmSequenceForceBatchesIterator{contract: _Polygonzkevm.contract, event: "SequenceForceBatches", logs: logs, sub: sub}, nil
}

// WatchSequenceForceBatches is a free log subscription operation binding the contract event 0x648a61dd2438f072f5a1960939abd30f37aea80d2e94c9792ad142d3e0a490a4.
//
// Solidity: event SequenceForceBatches(uint64 indexed numBatch)
func (_Polygonzkevm *PolygonzkevmFilterer) WatchSequenceForceBatches(opts *bind.WatchOpts, sink chan<- *PolygonzkevmSequenceForceBatches, numBatch []uint64) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	logs, sub, err := _Polygonzkevm.contract.WatchLogs(opts, "SequenceForceBatches", numBatchRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmSequenceForceBatches)
				if err := _Polygonzkevm.contract.UnpackLog(event, "SequenceForceBatches", log); err != nil {
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

// ParseSequenceForceBatches is a log parse operation binding the contract event 0x648a61dd2438f072f5a1960939abd30f37aea80d2e94c9792ad142d3e0a490a4.
//
// Solidity: event SequenceForceBatches(uint64 indexed numBatch)
func (_Polygonzkevm *PolygonzkevmFilterer) ParseSequenceForceBatches(log types.Log) (*PolygonzkevmSequenceForceBatches, error) {
	event := new(PolygonzkevmSequenceForceBatches)
	if err := _Polygonzkevm.contract.UnpackLog(event, "SequenceForceBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonzkevmSetForceBatchTimeoutIterator is returned from FilterSetForceBatchTimeout and is used to iterate over the raw logs and unpacked data for SetForceBatchTimeout events raised by the Polygonzkevm contract.
type PolygonzkevmSetForceBatchTimeoutIterator struct {
	Event *PolygonzkevmSetForceBatchTimeout // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmSetForceBatchTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmSetForceBatchTimeout)
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
		it.Event = new(PolygonzkevmSetForceBatchTimeout)
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
func (it *PolygonzkevmSetForceBatchTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmSetForceBatchTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmSetForceBatchTimeout represents a SetForceBatchTimeout event raised by the Polygonzkevm contract.
type PolygonzkevmSetForceBatchTimeout struct {
	NewforceBatchTimeout uint64
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSetForceBatchTimeout is a free log retrieval operation binding the contract event 0xa7eb6cb8a613eb4e8bddc1ac3d61ec6cf10898760f0b187bcca794c6ca6fa40b.
//
// Solidity: event SetForceBatchTimeout(uint64 newforceBatchTimeout)
func (_Polygonzkevm *PolygonzkevmFilterer) FilterSetForceBatchTimeout(opts *bind.FilterOpts) (*PolygonzkevmSetForceBatchTimeoutIterator, error) {

	logs, sub, err := _Polygonzkevm.contract.FilterLogs(opts, "SetForceBatchTimeout")
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmSetForceBatchTimeoutIterator{contract: _Polygonzkevm.contract, event: "SetForceBatchTimeout", logs: logs, sub: sub}, nil
}

// WatchSetForceBatchTimeout is a free log subscription operation binding the contract event 0xa7eb6cb8a613eb4e8bddc1ac3d61ec6cf10898760f0b187bcca794c6ca6fa40b.
//
// Solidity: event SetForceBatchTimeout(uint64 newforceBatchTimeout)
func (_Polygonzkevm *PolygonzkevmFilterer) WatchSetForceBatchTimeout(opts *bind.WatchOpts, sink chan<- *PolygonzkevmSetForceBatchTimeout) (event.Subscription, error) {

	logs, sub, err := _Polygonzkevm.contract.WatchLogs(opts, "SetForceBatchTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmSetForceBatchTimeout)
				if err := _Polygonzkevm.contract.UnpackLog(event, "SetForceBatchTimeout", log); err != nil {
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

// ParseSetForceBatchTimeout is a log parse operation binding the contract event 0xa7eb6cb8a613eb4e8bddc1ac3d61ec6cf10898760f0b187bcca794c6ca6fa40b.
//
// Solidity: event SetForceBatchTimeout(uint64 newforceBatchTimeout)
func (_Polygonzkevm *PolygonzkevmFilterer) ParseSetForceBatchTimeout(log types.Log) (*PolygonzkevmSetForceBatchTimeout, error) {
	event := new(PolygonzkevmSetForceBatchTimeout)
	if err := _Polygonzkevm.contract.UnpackLog(event, "SetForceBatchTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonzkevmSetTrustedSequencerIterator is returned from FilterSetTrustedSequencer and is used to iterate over the raw logs and unpacked data for SetTrustedSequencer events raised by the Polygonzkevm contract.
type PolygonzkevmSetTrustedSequencerIterator struct {
	Event *PolygonzkevmSetTrustedSequencer // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmSetTrustedSequencerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmSetTrustedSequencer)
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
		it.Event = new(PolygonzkevmSetTrustedSequencer)
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
func (it *PolygonzkevmSetTrustedSequencerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmSetTrustedSequencerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmSetTrustedSequencer represents a SetTrustedSequencer event raised by the Polygonzkevm contract.
type PolygonzkevmSetTrustedSequencer struct {
	NewTrustedSequencer common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedSequencer is a free log retrieval operation binding the contract event 0xf54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc0.
//
// Solidity: event SetTrustedSequencer(address newTrustedSequencer)
func (_Polygonzkevm *PolygonzkevmFilterer) FilterSetTrustedSequencer(opts *bind.FilterOpts) (*PolygonzkevmSetTrustedSequencerIterator, error) {

	logs, sub, err := _Polygonzkevm.contract.FilterLogs(opts, "SetTrustedSequencer")
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmSetTrustedSequencerIterator{contract: _Polygonzkevm.contract, event: "SetTrustedSequencer", logs: logs, sub: sub}, nil
}

// WatchSetTrustedSequencer is a free log subscription operation binding the contract event 0xf54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc0.
//
// Solidity: event SetTrustedSequencer(address newTrustedSequencer)
func (_Polygonzkevm *PolygonzkevmFilterer) WatchSetTrustedSequencer(opts *bind.WatchOpts, sink chan<- *PolygonzkevmSetTrustedSequencer) (event.Subscription, error) {

	logs, sub, err := _Polygonzkevm.contract.WatchLogs(opts, "SetTrustedSequencer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmSetTrustedSequencer)
				if err := _Polygonzkevm.contract.UnpackLog(event, "SetTrustedSequencer", log); err != nil {
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

// ParseSetTrustedSequencer is a log parse operation binding the contract event 0xf54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc0.
//
// Solidity: event SetTrustedSequencer(address newTrustedSequencer)
func (_Polygonzkevm *PolygonzkevmFilterer) ParseSetTrustedSequencer(log types.Log) (*PolygonzkevmSetTrustedSequencer, error) {
	event := new(PolygonzkevmSetTrustedSequencer)
	if err := _Polygonzkevm.contract.UnpackLog(event, "SetTrustedSequencer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonzkevmSetTrustedSequencerURLIterator is returned from FilterSetTrustedSequencerURL and is used to iterate over the raw logs and unpacked data for SetTrustedSequencerURL events raised by the Polygonzkevm contract.
type PolygonzkevmSetTrustedSequencerURLIterator struct {
	Event *PolygonzkevmSetTrustedSequencerURL // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmSetTrustedSequencerURLIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmSetTrustedSequencerURL)
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
		it.Event = new(PolygonzkevmSetTrustedSequencerURL)
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
func (it *PolygonzkevmSetTrustedSequencerURLIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmSetTrustedSequencerURLIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmSetTrustedSequencerURL represents a SetTrustedSequencerURL event raised by the Polygonzkevm contract.
type PolygonzkevmSetTrustedSequencerURL struct {
	NewTrustedSequencerURL string
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedSequencerURL is a free log retrieval operation binding the contract event 0x6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b20.
//
// Solidity: event SetTrustedSequencerURL(string newTrustedSequencerURL)
func (_Polygonzkevm *PolygonzkevmFilterer) FilterSetTrustedSequencerURL(opts *bind.FilterOpts) (*PolygonzkevmSetTrustedSequencerURLIterator, error) {

	logs, sub, err := _Polygonzkevm.contract.FilterLogs(opts, "SetTrustedSequencerURL")
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmSetTrustedSequencerURLIterator{contract: _Polygonzkevm.contract, event: "SetTrustedSequencerURL", logs: logs, sub: sub}, nil
}

// WatchSetTrustedSequencerURL is a free log subscription operation binding the contract event 0x6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b20.
//
// Solidity: event SetTrustedSequencerURL(string newTrustedSequencerURL)
func (_Polygonzkevm *PolygonzkevmFilterer) WatchSetTrustedSequencerURL(opts *bind.WatchOpts, sink chan<- *PolygonzkevmSetTrustedSequencerURL) (event.Subscription, error) {

	logs, sub, err := _Polygonzkevm.contract.WatchLogs(opts, "SetTrustedSequencerURL")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmSetTrustedSequencerURL)
				if err := _Polygonzkevm.contract.UnpackLog(event, "SetTrustedSequencerURL", log); err != nil {
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

// ParseSetTrustedSequencerURL is a log parse operation binding the contract event 0x6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b20.
//
// Solidity: event SetTrustedSequencerURL(string newTrustedSequencerURL)
func (_Polygonzkevm *PolygonzkevmFilterer) ParseSetTrustedSequencerURL(log types.Log) (*PolygonzkevmSetTrustedSequencerURL, error) {
	event := new(PolygonzkevmSetTrustedSequencerURL)
	if err := _Polygonzkevm.contract.UnpackLog(event, "SetTrustedSequencerURL", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonzkevmTransferAdminRoleIterator is returned from FilterTransferAdminRole and is used to iterate over the raw logs and unpacked data for TransferAdminRole events raised by the Polygonzkevm contract.
type PolygonzkevmTransferAdminRoleIterator struct {
	Event *PolygonzkevmTransferAdminRole // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmTransferAdminRoleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmTransferAdminRole)
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
		it.Event = new(PolygonzkevmTransferAdminRole)
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
func (it *PolygonzkevmTransferAdminRoleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmTransferAdminRoleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmTransferAdminRole represents a TransferAdminRole event raised by the Polygonzkevm contract.
type PolygonzkevmTransferAdminRole struct {
	NewPendingAdmin common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterTransferAdminRole is a free log retrieval operation binding the contract event 0xa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce6.
//
// Solidity: event TransferAdminRole(address newPendingAdmin)
func (_Polygonzkevm *PolygonzkevmFilterer) FilterTransferAdminRole(opts *bind.FilterOpts) (*PolygonzkevmTransferAdminRoleIterator, error) {

	logs, sub, err := _Polygonzkevm.contract.FilterLogs(opts, "TransferAdminRole")
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmTransferAdminRoleIterator{contract: _Polygonzkevm.contract, event: "TransferAdminRole", logs: logs, sub: sub}, nil
}

// WatchTransferAdminRole is a free log subscription operation binding the contract event 0xa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce6.
//
// Solidity: event TransferAdminRole(address newPendingAdmin)
func (_Polygonzkevm *PolygonzkevmFilterer) WatchTransferAdminRole(opts *bind.WatchOpts, sink chan<- *PolygonzkevmTransferAdminRole) (event.Subscription, error) {

	logs, sub, err := _Polygonzkevm.contract.WatchLogs(opts, "TransferAdminRole")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmTransferAdminRole)
				if err := _Polygonzkevm.contract.UnpackLog(event, "TransferAdminRole", log); err != nil {
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

// ParseTransferAdminRole is a log parse operation binding the contract event 0xa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce6.
//
// Solidity: event TransferAdminRole(address newPendingAdmin)
func (_Polygonzkevm *PolygonzkevmFilterer) ParseTransferAdminRole(log types.Log) (*PolygonzkevmTransferAdminRole, error) {
	event := new(PolygonzkevmTransferAdminRole)
	if err := _Polygonzkevm.contract.UnpackLog(event, "TransferAdminRole", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonzkevmVerifyBatchesIterator is returned from FilterVerifyBatches and is used to iterate over the raw logs and unpacked data for VerifyBatches events raised by the Polygonzkevm contract.
type PolygonzkevmVerifyBatchesIterator struct {
	Event *PolygonzkevmVerifyBatches // Event containing the contract specifics and raw log

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
func (it *PolygonzkevmVerifyBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonzkevmVerifyBatches)
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
		it.Event = new(PolygonzkevmVerifyBatches)
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
func (it *PolygonzkevmVerifyBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonzkevmVerifyBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonzkevmVerifyBatches represents a VerifyBatches event raised by the Polygonzkevm contract.
type PolygonzkevmVerifyBatches struct {
	NumBatch   uint64
	StateRoot  [32]byte
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVerifyBatches is a free log retrieval operation binding the contract event 0x9c72852172521097ba7e1482e6b44b351323df0155f97f4ea18fcec28e1f5966.
//
// Solidity: event VerifyBatches(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Polygonzkevm *PolygonzkevmFilterer) FilterVerifyBatches(opts *bind.FilterOpts, numBatch []uint64, aggregator []common.Address) (*PolygonzkevmVerifyBatchesIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Polygonzkevm.contract.FilterLogs(opts, "VerifyBatches", numBatchRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &PolygonzkevmVerifyBatchesIterator{contract: _Polygonzkevm.contract, event: "VerifyBatches", logs: logs, sub: sub}, nil
}

// WatchVerifyBatches is a free log subscription operation binding the contract event 0x9c72852172521097ba7e1482e6b44b351323df0155f97f4ea18fcec28e1f5966.
//
// Solidity: event VerifyBatches(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Polygonzkevm *PolygonzkevmFilterer) WatchVerifyBatches(opts *bind.WatchOpts, sink chan<- *PolygonzkevmVerifyBatches, numBatch []uint64, aggregator []common.Address) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Polygonzkevm.contract.WatchLogs(opts, "VerifyBatches", numBatchRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonzkevmVerifyBatches)
				if err := _Polygonzkevm.contract.UnpackLog(event, "VerifyBatches", log); err != nil {
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

// ParseVerifyBatches is a log parse operation binding the contract event 0x9c72852172521097ba7e1482e6b44b351323df0155f97f4ea18fcec28e1f5966.
//
// Solidity: event VerifyBatches(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Polygonzkevm *PolygonzkevmFilterer) ParseVerifyBatches(log types.Log) (*PolygonzkevmVerifyBatches, error) {
	event := new(PolygonzkevmVerifyBatches)
	if err := _Polygonzkevm.contract.UnpackLog(event, "VerifyBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
