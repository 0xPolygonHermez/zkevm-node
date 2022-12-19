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

// ProofOfEfficiencyBatchData is an auto generated low-level Go binding around an user-defined struct.
type ProofOfEfficiencyBatchData struct {
	Transactions       []byte
	GlobalExitRoot     [32]byte
	Timestamp          uint64
	MinForcedTimestamp uint64
}

// ProofOfEfficiencyForcedBatchData is an auto generated low-level Go binding around an user-defined struct.
type ProofOfEfficiencyForcedBatchData struct {
	Transactions       []byte
	GlobalExitRoot     [32]byte
	MinForcedTimestamp uint64
}

// ProofOfEfficiencyInitializePackedParameters is an auto generated low-level Go binding around an user-defined struct.
type ProofOfEfficiencyInitializePackedParameters struct {
	Admin                    common.Address
	ChainID                  uint64
	TrustedSequencer         common.Address
	PendingStateTimeout      uint64
	ForceBatchAllowed        bool
	TrustedAggregator        common.Address
	TrustedAggregatorTimeout uint64
}

// ProofofefficiencyMetaData contains all meta data concerning the Proofofefficiency contract.
var ProofofefficiencyMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"ConsolidatePendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EmergencyStateActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EmergencyStateDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"forceBatchNum\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"lastGlobalExitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"}],\"name\":\"ForceBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"OverridePendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"storedStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"provedStateRoot\",\"type\":\"bytes32\"}],\"name\":\"ProveNonDeterministicPendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"}],\"name\":\"SequenceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"}],\"name\":\"SequenceForceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"SetAdmin\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"newForceBatchAllowed\",\"type\":\"bool\"}],\"name\":\"SetForceBatchAllowed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newPendingStateTimeout\",\"type\":\"uint64\"}],\"name\":\"SetPendingStateTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTrustedAggregator\",\"type\":\"address\"}],\"name\":\"SetTrustedAggregator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newTrustedAggregatorTimeout\",\"type\":\"uint64\"}],\"name\":\"SetTrustedAggregatorTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTrustedSequencer\",\"type\":\"address\"}],\"name\":\"SetTrustedSequencer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"newTrustedSequencerURL\",\"type\":\"string\"}],\"name\":\"SetTrustedSequencerURL\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"TrustedVerifyBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatches\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"FORCE_BATCH_TIMEOUT\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"HALT_AGGREGATION_TIMEOUT\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_TRANSACTIONS_BYTE_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_VERIFY_BATCHES\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MULTIPLIER_BATCH_FEE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"VERIFY_BATCH_TIME_TARGET\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequencedBatchNum\",\"type\":\"uint64\"}],\"name\":\"activateEmergencyState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"batchFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"batchNumToStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"contractIBridge\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateRewardPerBatch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chainID\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"consolidatePendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deactivateEmergencyState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maticAmount\",\"type\":\"uint256\"}],\"name\":\"forceBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"forceBatchAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"forcedBatches\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentBatchFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"getInputSnarkBytes\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastVerifiedBatch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"globalExitRootManager\",\"outputs\":[{\"internalType\":\"contractIGlobalExitRootManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIGlobalExitRootManager\",\"name\":\"_globalExitRootManager\",\"type\":\"address\"},{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"_matic\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"_rollupVerifier\",\"type\":\"address\"},{\"internalType\":\"contractIBridge\",\"name\":\"_bridgeAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"trustedSequencer\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateTimeout\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"forceBatchAllowed\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"trustedAggregator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"trustedAggregatorTimeout\",\"type\":\"uint64\"}],\"internalType\":\"structProofOfEfficiency.InitializePackedParameters\",\"name\":\"initializePackedParameters\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"genesisRoot\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_trustedSequencerURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_networkName\",\"type\":\"string\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isEmergencyState\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"isPendingStateConsolidable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBatchSequenced\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastForceBatch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastForceBatchSequenced\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastPendingState\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastPendingStateConsolidated\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastTimestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastVerifiedBatch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"matic\",\"outputs\":[{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"initPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofA\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"proofB\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofC\",\"type\":\"uint256[2]\"}],\"name\":\"overridePendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingStateTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"pendingStateTransitions\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"initPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofA\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"proofB\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofC\",\"type\":\"uint256[2]\"}],\"name\":\"proveNonDeterministicPendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupVerifier\",\"outputs\":[{\"internalType\":\"contractIVerifierRollup\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"globalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minForcedTimestamp\",\"type\":\"uint64\"}],\"internalType\":\"structProofOfEfficiency.BatchData[]\",\"name\":\"batches\",\"type\":\"tuple[]\"}],\"name\":\"sequenceBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"globalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"minForcedTimestamp\",\"type\":\"uint64\"}],\"internalType\":\"structProofOfEfficiency.ForcedBatchData[]\",\"name\":\"batches\",\"type\":\"tuple[]\"}],\"name\":\"sequenceForceBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"sequencedBatches\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"accInputHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequencedTimestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"previousLastBatchSequenced\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"newForceBatchAllowed\",\"type\":\"bool\"}],\"name\":\"setForceBatchAllowed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newPendingStateTimeout\",\"type\":\"uint64\"}],\"name\":\"setPendingStateTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newTrustedAggregator\",\"type\":\"address\"}],\"name\":\"setTrustedAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newTrustedAggregatorTimeout\",\"type\":\"uint64\"}],\"name\":\"setTrustedAggregatorTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newTrustedSequencer\",\"type\":\"address\"}],\"name\":\"setTrustedSequencer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newTrustedSequencerURL\",\"type\":\"string\"}],\"name\":\"setTrustedSequencerURL\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedAggregator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedAggregatorTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedSequencer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedSequencerURL\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofA\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"proofB\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofC\",\"type\":\"uint256[2]\"}],\"name\":\"trustedVerifyBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofA\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"proofB\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofC\",\"type\":\"uint256[2]\"}],\"name\":\"verifyBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50615e7880620000216000396000f3fe608060405234801561001057600080fd5b506004361061038e5760003560e01c80638c4a0af7116101de578063d8d1091b1161010f578063e8bf92ed116100ad578063f1d7b21c1161007c578063f1d7b21c14610876578063f2fde38b1461087e578063f851a44014610891578063f8b823e4146108a457600080fd5b8063e8bf92ed1461082a578063eaeb077b1461083d578063edc4112114610850578063f14916d61461086357600080fd5b8063dbc16976116100e9578063dbc16976146107ec578063e11f3f18146107f4578063e217cfd614610807578063e7a7ed021461081057600080fd5b8063d8d1091b146107ab578063d8f54db0146107be578063d939b315146107d257600080fd5b8063adc879e91161017c578063c0ed84e011610156578063c0ed84e014610763578063c89e42df1461076b578063cfa8ed471461077e578063d02103ca1461079857600080fd5b8063adc879e9146106d1578063b4d63f58146106eb578063b6b0b0971461074b57600080fd5b80639c9f3dfe116101b85780639c9f3dfe146106995780639f0d039d146106ac578063a3c573eb146106b4578063ab9fc5ef146106c757600080fd5b80638c4a0af71461066d5780638da5cb5b1461068057806399f5634e1461069157600080fd5b80634a1a89a7116102c3578063704b6c02116102615780637fcb3653116102305780637fcb3653146105c8578063837a4738146105db578063841b24d7146106495780638b48931e1461066357600080fd5b8063704b6c0214610587578063715018a61461059a5780637215541a146105a257806375c508b3146105b557600080fd5b8063542028d51161029d578063542028d51461053957806360943d6a146105415780636b8616ce146105545780636ff512cc1461057457600080fd5b80634a1a89a7146104ec5780634a910e6a146105065780635392c5e01461051957600080fd5b8063383b3be811610330578063423fa8561161030a578063423fa8561461049257806345605267146104ac578063458c0477146104c65780634834a343146104d957600080fd5b8063383b3be814610457578063394218e91461046a5780633c1582671461047f57600080fd5b806319d8ac611161036c57806319d8ac61146103ef578063220d78991461040257806329878983146104155780632d0889d31461044057600080fd5b8063107bf28c14610393578063137f1edf146103b157806315064c96146103d2575b600080fd5b61039b6108ad565b6040516103a891906152b6565b60405180910390f35b6103ba61070881565b6040516001600160401b0390911681526020016103a8565b6065546103df9060ff1681565b60405190151581526020016103a8565b6068546103ba906001600160401b031681565b61039b6104103660046152e5565b61093b565b606a54610428906001600160a01b031681565b6040516001600160a01b0390911681526020016103a8565b61044961ea6081565b6040519081526020016103a8565b6103df610465366004615332565b610afe565b61047d610478366004615332565b610b45565b005b61047d61048d36600461546f565b610d6e565b6068546103ba90600160401b90046001600160401b031681565b6068546103ba90600160801b90046001600160401b031681565b6072546103ba906001600160401b031681565b61047d6104e73660046155ad565b6116d0565b6072546103ba90600160401b90046001600160401b031681565b61047d610514366004615332565b611ab0565b610449610527366004615332565b606d6020526000908152604090205481565b61039b611d73565b61047d61054f36600461564f565b611d80565b610449610562366004615332565b60666020526000908152604090205481565b61047d610582366004615726565b612116565b61047d610595366004615726565b6121f0565b61047d6122a9565b61047d6105b0366004615332565b6122bd565b61047d6105c3366004615743565b612578565b6069546103ba906001600160401b031681565b61061e6105e93660046157e1565b6071602052600090815260409020805460018201546002909201546001600160401b0380831693600160401b90930416919084565b604080516001600160401b0395861681529490931660208501529183015260608201526080016103a8565b6072546103ba90600160c01b90046001600160401b031681565b6103ba62093a8081565b61047d61067b366004615808565b612668565b6033546001600160a01b0316610428565b610449612720565b61047d6106a7366004615332565b61281a565b607454610449565b607054610428906001600160a01b031681565b6103ba6206978081565b606c546103ba90600160a81b90046001600160401b031681565b6107266106f9366004615332565b606760205260009081526040902080546001909101546001600160401b0380821691600160401b90041683565b604080519384526001600160401b0392831660208501529116908201526060016103a8565b6065546104289061010090046001600160a01b031681565b6103ba612a31565b61047d610779366004615825565b612a7e565b60695461042890600160401b90046001600160a01b031681565b606c54610428906001600160a01b031681565b61047d6107b9366004615861565b612b25565b606c546103df90600160a01b900460ff1681565b6072546103ba90600160801b90046001600160401b031681565b61047d6131a3565b61047d610802366004615743565b6132f6565b6103ba6103e881565b6068546103ba90600160c01b90046001600160401b031681565b606b54610428906001600160a01b031681565b61047d61084b366004615953565b6134cb565b61047d61085e3660046155ad565b61390b565b61047d610871366004615726565b613a94565b610449600b81565b61047d61088c366004615726565b613b4d565b607354610428906001600160a01b031681565b61044960745481565b606f80546108ba90615997565b80601f01602080910402602001604051908101604052809291908181526020018280546108e690615997565b80156109335780601f1061090857610100808354040283529160200191610933565b820191906000526020600020905b81548152906001019060200180831161091657829003601f168201915b505050505081565b6001600160401b0380861660008181526067602052604080822054938816825290205460609291158061096d57508115155b6109f25760405162461bcd60e51b815260206004820152604560248201527f50726f6f664f66456666696369656e63793a3a676574496e707574536e61726b60448201527f42797465733a206f6c64416363496e7075744861736820646f6573206e6f7420606482015264195e1a5cdd60da1b608482015260a4015b60405180910390fd5b80610a735760405162461bcd60e51b815260206004820152604560248201527f50726f6f664f66456666696369656e63793a3a676574496e707574536e61726b60448201527f42797465733a206e6577416363496e7075744861736820646f6573206e6f7420606482015264195e1a5cdd60da1b608482015260a4016109e9565b606c54604080516bffffffffffffffffffffffff193360601b166020820152603481019790975260548701939093526001600160c01b031960c0998a1b81166074880152600160a81b909104891b8116607c870152608486019490945260a485015260c4840194909452509290931b90911660e4830152805180830360cc01815260ec909201905290565b6072546001600160401b0382811660009081526071602052604081205490924292610b3492600160801b909204811691166159e7565b6001600160401b0316111592915050565b6073546001600160a01b03163314610bb05760405162461bcd60e51b815260206004820152602860248201527f50726f6f664f66456666696369656e63793a3a6f6e6c7941646d696e3a206f6e604482015267363c9030b236b4b760c11b60648201526084016109e9565b62093a806001600160401b0382161115610c455760405162461bcd60e51b815260206004820152604a60248201527f50726f6f664f66456666696369656e63793a3a73657450656e64696e6753746160448201527f746554696d656f75743a206578636565642068616c74206167677265676174696064820152691bdb881d1a5b595bdd5d60b21b608482015260a4016109e9565b60655460ff16610d00576072546001600160401b03600160c01b909104811690821610610d005760405162461bcd60e51b815260206004820152604960248201527f50726f6f664f66456666696369656e63793a3a7365745472757374656441676760448201527f72656761746f7254696d656f75743a206e65772074696d656f7574206d75737460648201527f206265206c6f7765720000000000000000000000000000000000000000000000608482015260a4016109e9565b6072805477ffffffffffffffffffffffffffffffffffffffffffffffff16600160c01b6001600160401b038416908102919091179091556040519081527f1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1906020015b60405180910390a150565b60655460ff1615610df25760405162461bcd60e51b815260206004820152604260248201527f456d657267656e63794d616e616765723a3a69664e6f74456d657267656e637960448201527f53746174653a206f6e6c79206966206e6f7420656d657267656e637920737461606482015261746560f01b608482015260a4016109e9565b606954600160401b90046001600160a01b03163314610e795760405162461bcd60e51b815260206004820152603f60248201527f50726f6f664f66456666696369656e63793a3a6f6e6c7954727573746564536560448201527f7175656e6365723a206f6e6c7920747275737465642073657175656e6365720060648201526084016109e9565b805180610ef95760405162461bcd60e51b815260206004820152604260248201527f50726f6f664f66456666696369656e63793a3a73657175656e6365426174636860448201527f65733a204174206c65617374206d7573742073657175656e63652031206261746064820152610c6d60f31b608482015260a4016109e9565b6103e88110610f965760405162461bcd60e51b815260206004820152604560248201527f50726f6f664f66456666696369656e63793a3a73657175656e6365426174636860448201527f65733a2043616e6e6f742073657175656e63652074686174206d616e7920626160648201527f7463686573000000000000000000000000000000000000000000000000000000608482015260a4016109e9565b6068546001600160401b03600160401b82048116600081815260676020526040812054838516949293600160801b90930490921691905b858110156114e5576000878281518110610fe957610fe9615a12565b60200260200101519050600081606001516001600160401b031611156111cb578361101381615a28565b94505060008160000151805190602001208260200151836060015160405160200161105e93929190928352602083019190915260c01b6001600160c01b031916604082015260480190565b60408051601f1981840301815291815281516020928301206001600160401b03881660009081526066909352912054909150811461110f5760405162461bcd60e51b815260206004820152604260248201527f50726f6f664f66456666696369656e63793a3a73657175656e6365426174636860448201527f65733a20466f7263656420626174636865732064617461206d757374206d61746064820152610c6d60f31b608482015260a4016109e9565b81606001516001600160401b031682604001516001600160401b031610156111c55760405162461bcd60e51b815260206004820152605d60248201527f50726f6f664f66456666696369656e63793a3a73657175656e6365426174636860448201527f65733a20466f7263656420626174636865732074696d657374616d70206d757360648201527f7420626520626967676572206f7220657175616c207468616e206d696e000000608482015260a4016109e9565b5061137e565b6020810151158061126d5750606c5460208201516040517f257b36320000000000000000000000000000000000000000000000000000000081526001600160a01b039092169163257b3632916112279160040190815260200190565b6020604051808303816000875af1158015611246573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061126a9190615a4e565b15155b6112df5760405162461bcd60e51b815260206004820152603f60248201527f50726f6f664f66456666696369656e63793a3a73657175656e6365426174636860448201527f65733a20476c6f62616c206578697420726f6f74206d7573742065786973740060648201526084016109e9565b80515161ea601161137e5760405162461bcd60e51b815260206004820152604a60248201527f50726f6f664f664566666963696550656e64696e67537461746563793a3a736560448201527f7175656e6365426174636865733a205472616e73616374696f6e73206279746560648201527f73206f766572666c6f7700000000000000000000000000000000000000000000608482015260a4016109e9565b856001600160401b031681604001516001600160401b0316101580156113b157504281604001516001600160401b031611155b6114495760405162461bcd60e51b815260206004820152604260248201527f50726f6f664f66456666696369656e63793a3a73657175656e6365426174636860448201527f65733a2054696d657374616d70206d75737420626520696e736964652072616e60648201527f6765000000000000000000000000000000000000000000000000000000000000608482015260a4016109e9565b805180516020918201208183015160408085015181519485018890529084019290925260608084019190915260c09190911b6001600160c01b031916608083015233901b6bffffffffffffffffffffffff19166088820152609c0160405160208183030381529060405280519060200120925084806114c790615a28565b955050806040015195505080806114dd90615a67565b915050610fcd565b506068546001600160401b03600160c01b909104811690831611156115725760405162461bcd60e51b815260206004820152603a60248201527f50726f6f664f66456666696369656e63793a3a73657175656e6365426174636860448201527f65733a20466f7263652062617463686573206f766572666c6f7700000000000060648201526084016109e9565b60685460009061159290600160801b90046001600160401b031684615a80565b6115a5906001600160401b031687615aa8565b60408051606081018252848152426001600160401b03908116602080840191825260688054600160401b9081900485168688019081528c861660008181526067909552979093209551865592516001909501805492519585166fffffffffffffffffffffffffffffffff199384161795851684029590951790945583548b841691161793029290921767ffffffffffffffff60801b1916600160801b92871692909202919091179055607454909150611681903390309084906116689190615abf565b60655461010090046001600160a01b0316929190613bda565b611689613c8b565b606854604051600160401b9091046001600160401b0316907f303446e6a8cb73c83dff421c0b1d5e5ce0719dab1bff13660fc254e58cc17fce90600090a250505050505050565b60655460ff16156117545760405162461bcd60e51b815260206004820152604260248201527f456d657267656e63794d616e616765723a3a69664e6f74456d657267656e637960448201527f53746174653a206f6e6c79206966206e6f7420656d657267656e637920737461606482015261746560f01b608482015260a4016109e9565b6072546001600160401b03878116600090815260676020526040902060010154429261178b92600160c01b909104811691166159e7565b6001600160401b0316111561182e5760405162461bcd60e51b815260206004820152604860248201527f50726f6f664f66456666696369656e63793a3a7665726966794261746368657360448201527f3a20747275737465642061676772656761746f722074696d656f7574206e6f7460648201527f2065787069726564000000000000000000000000000000000000000000000000608482015260a4016109e9565b6103e861183b8888615a80565b6001600160401b0316106118c15760405162461bcd60e51b815260206004820152604160248201527f50726f6f664f66456666696369656e63793a3a7665726966794261746368657360448201527f3a2063616e6e6f74207665726966792074686174206d616e79206261746368656064820152607360f81b608482015260a4016109e9565b6118d18888888888888888613d2f565b6118da866142b4565b607254600160801b90046001600160401b03166000036119ab576069805467ffffffffffffffff19166001600160401b038881169182179092556000908152606d60205260409020859055607254161561194857607280546fffffffffffffffffffffffffffffffff191690555b606c546040516333d6247d60e01b8152600481018790526001600160a01b03909116906333d6247d90602401600060405180830381600087803b15801561198e57600080fd5b505af11580156119a2573d6000803e3d6000fd5b50505050611a65565b6119b3613c8b565b607280546001600160401b03169060006119cc83615a28565b82546001600160401b039182166101009390930a92830292820219169190911790915560408051608081018252428316815289831660208083019182528284018b8152606084018b8152607254871660009081526071909352949091209251835492518616600160401b026fffffffffffffffffffffffffffffffff199093169516949094171781559151600183015551600290910155505b60405184815233906001600160401b038816907f9c72852172521097ba7e1482e6b44b351323df0155f97f4ea18fcec28e1f5966906020015b60405180910390a35050505050505050565b6001600160401b03811615801590611add57506072546001600160401b03600160401b9091048116908216115b8015611af857506072546001600160401b0390811690821611155b611b905760405162461bcd60e51b815260206004820152604860248201527f50726f6f664f66456666696369656e63793a3a636f6e736f6c6964617465506560448201527f6e64696e6753746174653a2070656e64696e6753746174654e756d206d75737460648201527f20696e76616c6964000000000000000000000000000000000000000000000000608482015260a4016109e9565b606a546001600160a01b03163314611c4357611bab81610afe565b611c435760405162461bcd60e51b815260206004820152605960248201527f50726f6f664f66456666696369656e63793a3a636f6e736f6c6964617465506560448201527f6e64696e6753746174653a2070656e64696e67207374617465206973206e6f7460648201527f20726561647920746f20626520636f6e736f6c69646174656400000000000000608482015260a4016109e9565b6001600160401b038181166000818152607160209081526040808320805460698054600160401b9283900490981667ffffffffffffffff19909816881790556002820154878652606d9094529382902092909255607280546fffffffffffffffff000000000000000019169390940292909217909255606c54600183015491516333d6247d60e01b815260048101929092529192916001600160a01b0316906333d6247d90602401600060405180830381600087803b158015611d0557600080fd5b505af1158015611d19573d6000803e3d6000fd5b50505050826001600160401b0316816001600160401b03167f328d3c6c0fd6f1be0515e422f2d87e59f25922cbc2233568515a0c4bc3f8510e8460020154604051611d6691815260200190565b60405180910390a3505050565b606e80546108ba90615997565b600054610100900460ff1615808015611da05750600054600160ff909116105b80611dba5750303b158015611dba575060005460ff166001145b611e2c5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084016109e9565b6000805460ff191660011790558015611e4f576000805461ff0019166101001790555b606c80546001600160a01b03199081166001600160a01b038c811691909117909255606580547fffffffffffffffffffffff0000000000000000000000000000000000000000ff166101008c851602179055606b805482168a841617905560708054909116918816919091179055611eca6020860186615726565b607380546001600160a01b0319166001600160a01b0392909216919091179055611efa6060860160408701615726565b606980546001600160a01b0392909216600160401b027fffffffff0000000000000000000000000000000000000000ffffffffffffffff909216919091179055611f4a60c0860160a08701615726565b606a80546001600160a01b0319166001600160a01b039290921691909117905560008052606d6020527fda90043ba5b4096ba14704bc227ab0d3167da15b887e62ab2e76e37daa711356849055611fa760e0860160c08701615332565b607280546001600160401b0392909216600160c01b0277ffffffffffffffffffffffffffffffffffffffffffffffff909216919091179055611fef6040860160208701615332565b606c80546001600160401b0392909216600160a81b027fffffff0000000000000000ffffffffffffffffffffffffffffffffffffffffff90921691909117905561203f6080860160608701615332565b607280546001600160401b0392909216600160801b0267ffffffffffffffff60801b1990921691909117905561207b60a0860160808701615808565b606c8054911515600160a01b0260ff60a01b19909216919091179055606e6120a38482615b24565b50606f6120b08382615b24565b50670de0b6b3a76400006074556120c56144d6565b801561210b576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050505050505050565b6073546001600160a01b031633146121815760405162461bcd60e51b815260206004820152602860248201527f50726f6f664f66456666696369656e63793a3a6f6e6c7941646d696e3a206f6e604482015267363c9030b236b4b760c11b60648201526084016109e9565b606980547fffffffff0000000000000000000000000000000000000000ffffffffffffffff16600160401b6001600160a01b038416908102919091179091556040519081527ff54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc090602001610d63565b6073546001600160a01b0316331461225b5760405162461bcd60e51b815260206004820152602860248201527f50726f6f664f66456666696369656e63793a3a6f6e6c7941646d696e3a206f6e604482015267363c9030b236b4b760c11b60648201526084016109e9565b607380546001600160a01b0319166001600160a01b0383169081179091556040519081527f5a272403b402d892977df56625f4164ccaf70ca3863991c43ecfe76a6905b0a190602001610d63565b6122b161455c565b6122bb60006145b6565b565b6033546001600160a01b0316331461256d576072546000906001600160401b03161561230e57506072546001600160401b03908116600090815260716020526040902054600160401b90041661231c565b506069546001600160401b03165b80826001600160401b0316116123c05760405162461bcd60e51b815260206004820152604160248201527f50726f6f664f66456666696369656e63793a3a6163746976617465456d65726760448201527f656e637953746174653a20426174636820616c7265616479207665726966696560648201527f6400000000000000000000000000000000000000000000000000000000000000608482015260a4016109e9565b6068546001600160401b03600160401b90910481169083161180159061240257506001600160401b038083166000908152606760205260409020600101541615155b61249a5760405162461bcd60e51b815260206004820152605560248201527f50726f6f664f66456666696369656e63793a3a6163746976617465456d65726760448201527f656e637953746174653a204261746368206e6f742073657175656e636564206f60648201527f72206e6f7420656e64206f662073657175656e63650000000000000000000000608482015260a4016109e9565b6001600160401b0380831660009081526067602052604090206001015442916124c89162093a8091166159e7565b6001600160401b0316111561256b5760405162461bcd60e51b815260206004820152605260248201527f50726f6f664f66456666696369656e63793a3a6163746976617465456d65726760448201527f656e637953746174653a204167677265676174696f6e2068616c742074696d6560648201527f6f7574206973206e6f7420657870697265640000000000000000000000000000608482015260a4016109e9565b505b612575614608565b50565b60655460ff16156125fc5760405162461bcd60e51b815260206004820152604260248201527f456d657267656e63794d616e616765723a3a69664e6f74456d657267656e637960448201527f53746174653a206f6e6c79206966206e6f7420656d657267656e637920737461606482015261746560f01b608482015260a4016109e9565b61260d898989898989898989614678565b6001600160401b0386166000908152606d60209081526040918290205482519081529081018690527f1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010910160405180910390a161210b614608565b6073546001600160a01b031633146126d35760405162461bcd60e51b815260206004820152602860248201527f50726f6f664f66456666696369656e63793a3a6f6e6c7941646d696e3a206f6e604482015267363c9030b236b4b760c11b60648201526084016109e9565b606c8054821515600160a01b0260ff60a01b199091161790556040517fbacda50a4a8575be1d91a7ebe29ee45056f3a94f12a2281eb6b43afa33bcefe690610d6390831515815260200190565b6065546040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015260009182916101009091046001600160a01b0316906370a0823190602401602060405180830381865afa15801561278a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127ae9190615a4e565b905060006127ba612a31565b6068546001600160401b03600160401b82048116916127ea91600160801b8204811691600160c01b900416615a80565b6127f491906159e7565b6127fe9190615a80565b6001600160401b031690506128138183615bf9565b9250505090565b6073546001600160a01b031633146128855760405162461bcd60e51b815260206004820152602860248201527f50726f6f664f66456666696369656e63793a3a6f6e6c7941646d696e3a206f6e604482015267363c9030b236b4b760c11b60648201526084016109e9565b62093a806001600160401b038216111561291a5760405162461bcd60e51b815260206004820152604a60248201527f50726f6f664f66456666696369656e63793a3a73657450656e64696e6753746160448201527f746554696d656f75743a206578636565642068616c74206167677265676174696064820152691bdb881d1a5b595bdd5d60b21b608482015260a4016109e9565b60655460ff166129d6576072546001600160401b03600160801b9091048116908216106129d65760405162461bcd60e51b8152602060048201526044602482018190527f50726f6f664f66456666696369656e63793a3a73657450656e64696e67537461908201527f746554696d656f75743a206e65772074696d656f7574206d757374206265206c60648201527f6f77657200000000000000000000000000000000000000000000000000000000608482015260a4016109e9565b6072805467ffffffffffffffff60801b1916600160801b6001600160401b038416908102919091179091556040519081527fc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c7590602001610d63565b6072546000906001600160401b031615612a6e57506072546001600160401b03908116600090815260716020526040902054600160401b90041690565b506069546001600160401b031690565b6073546001600160a01b03163314612ae95760405162461bcd60e51b815260206004820152602860248201527f50726f6f664f66456666696369656e63793a3a6f6e6c7941646d696e3a206f6e604482015267363c9030b236b4b760c11b60648201526084016109e9565b606e612af58282615b24565b507f6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b2081604051610d6391906152b6565b60655460ff1615612ba95760405162461bcd60e51b815260206004820152604260248201527f456d657267656e63794d616e616765723a3a69664e6f74456d657267656e637960448201527f53746174653a206f6e6c79206966206e6f7420656d657267656e637920737461606482015261746560f01b608482015260a4016109e9565b606c54600160a01b900460ff161515600114612c3e5760405162461bcd60e51b815260206004820152604860248201527f50726f6f664f66456666696369656e63793a3a6973466f72636542617463684160448201527f6c6c6f7765643a206f6e6c7920696620666f72636520626174636820697320616064820152677661696c61626c6560c01b608482015260a4016109e9565b805180612cbe5760405162461bcd60e51b815260206004820152604260248201527f50726f6f664f66456666696369656e63793a3a73657175656e6365466f72636560448201527f42617463683a204d75737420666f726365206174206c656173742031206261746064820152610c6d60f31b608482015260a4016109e9565b6103e88110612d3f5760405162461bcd60e51b815260206004820152604160248201527f50726f6f664f66456666696369656e63793a3a7665726966794261746368657360448201527f3a2063616e6e6f74207665726966792074686174206d616e79206261746368656064820152607360f81b608482015260a4016109e9565b6068546001600160401b03600160c01b8204811691612d67918491600160801b900416615c0d565b1115612ddb5760405162461bcd60e51b815260206004820152603a60248201527f50726f6f664f66456666696369656e63793a3a73657175656e6365466f72636560448201527f42617463683a20466f72636520626174636820696e76616c696400000000000060648201526084016109e9565b6068546001600160401b03600160401b820481166000818152606760205260408120549193600160801b9004909216915b8481101561309c576000868281518110612e2857612e28615a12565b602002602001015190508380612e3d90615a28565b825180516020918201208185015160408087015181519485019390935283015260c01b6001600160c01b03191660608201529095506000915060680160408051601f1981840301815291815281516020928301206001600160401b038816600090815260669093529120549091508114612f455760405162461bcd60e51b815260206004820152604760248201527f50726f6f664f66456666696369656e63793a3a73657175656e6365466f72636560448201527f426174636865733a20466f7263656420626174636865732064617461206d757360648201527f74206d6174636800000000000000000000000000000000000000000000000000608482015260a4016109e9565b612f50600188615aa8565b830361300d5742620697808360400151612f6a91906159e7565b6001600160401b0316111561300d5760405162461bcd60e51b815260206004820152604c60248201527f50726f6f664f66456666696369656e63793a3a73657175656e6365466f72636560448201527f42617463683a20466f72636564206261746368206973206e6f7420696e20746960648201527f6d656f757420706572696f640000000000000000000000000000000000000000608482015260a4016109e9565b8151805160209182012081840151604080519384018890528301919091526060808301919091524260c01b6001600160c01b031916608083015233901b6bffffffffffffffffffffffff19166088820152609c01604051602081830303815290604052805190602001209350858061308490615a28565b9650505050808061309490615a67565b915050612e0c565b506068805467ffffffffffffffff1916426001600160401b03908116918217808455604080516060810182528681526020808201958652600160401b9384900485168284019081528a861660008181526067909352848320935184559651600193909301805491519387166fffffffffffffffffffffffffffffffff199092169190911792861685029290921790915585547fffffffffffffffff00000000000000000000000000000000ffffffffffffffff1694830267ffffffffffffffff60801b191694909417600160801b88851602179485905551930416917f648a61dd2438f072f5a1960939abd30f37aea80d2e94c9792ad142d3e0a490a49190a25050505050565b60655460ff1661321b5760405162461bcd60e51b815260206004820152603b60248201527f456d657267656e63794d616e616765723a3a6966456d657267656e637953746160448201527f74653a206f6e6c7920696620656d657267656e6379207374617465000000000060648201526084016109e9565b6073546001600160a01b031633146132865760405162461bcd60e51b815260206004820152602860248201527f50726f6f664f66456666696369656e63793a3a6f6e6c7941646d696e3a206f6e604482015267363c9030b236b4b760c11b60648201526084016109e9565b607060009054906101000a90046001600160a01b03166001600160a01b031663dbc169766040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156132d657600080fd5b505af11580156132ea573d6000803e3d6000fd5b505050506122bb614e28565b606a546001600160a01b031633146133805760405162461bcd60e51b815260206004820152604160248201527f50726f6f664f66456666696369656e63793a3a6f6e6c7954727573746564416760448201527f6772656761746f723a206f6e6c7920747275737465642041676772656761746f6064820152603960f91b608482015260a4016109e9565b613391898989898989898989614678565b6069805467ffffffffffffffff19166001600160401b038881169182179092556000908152606d6020526040902085905560725416156133e557607280546fffffffffffffffffffffffffffffffff191690555b606c546040516333d6247d60e01b8152600481018790526001600160a01b03909116906333d6247d90602401600060405180830381600087803b15801561342b57600080fd5b505af115801561343f573d6000803e3d6000fd5b50506072805477ffffffffffffffffffffffffffffffffffffffffffffffff167a093a80000000000000000000000000000000000000000000000000179055505060405184815233906001600160401b038816907fcc1b5520188bf1dd3e63f98164b577c4d75c11a619ddea692112f0d1aec4cf729060200160405180910390a3505050505050505050565b60655460ff161561354f5760405162461bcd60e51b815260206004820152604260248201527f456d657267656e63794d616e616765723a3a69664e6f74456d657267656e637960448201527f53746174653a206f6e6c79206966206e6f7420656d657267656e637920737461606482015261746560f01b608482015260a4016109e9565b606c54600160a01b900460ff1615156001146135e45760405162461bcd60e51b815260206004820152604860248201527f50726f6f664f66456666696369656e63793a3a6973466f72636542617463684160448201527f6c6c6f7765643a206f6e6c7920696620666f72636520626174636820697320616064820152677661696c61626c6560c01b608482015260a4016109e9565b60006135ef60745490565b9050818111156136675760405162461bcd60e51b815260206004820152602f60248201527f50726f6f664f66456666696369656e63793a3a666f72636542617463683a206e60448201527f6f7420656e6f756768206d61746963000000000000000000000000000000000060648201526084016109e9565b61ea608351106136df5760405162461bcd60e51b815260206004820152603a60248201527f50726f6f664f66456666696369656e63793a3a666f72636542617463683a205460448201527f72616e73616374696f6e73206279746573206f766572666c6f7700000000000060648201526084016109e9565b6065546136fc9061010090046001600160a01b0316333084613bda565b606c54604080517f3ed691ef00000000000000000000000000000000000000000000000000000000815290516000926001600160a01b031691633ed691ef9160048083019260209291908290030181865afa15801561375f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906137839190615a4e565b60688054919250600160c01b9091046001600160401b03169060186137a783615a28565b91906101000a8154816001600160401b0302191690836001600160401b03160217905550508380519060200120814260405160200161380693929190928352602083019190915260c01b6001600160c01b031916604082015260480190565b60408051808303601f190181529181528151602092830120606854600160c01b90046001600160401b0316600090815260669093529120553233036138aa57606854604080518381523360208201526060918101829052600091810191909152600160c01b9091046001600160401b0316907ff94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc9319060800160405180910390a2613905565b606860189054906101000a90046001600160401b03166001600160401b03167ff94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc9318233876040516138fc93929190615c25565b60405180910390a25b50505050565b606a546001600160a01b031633146139955760405162461bcd60e51b815260206004820152604160248201527f50726f6f664f66456666696369656e63793a3a6f6e6c7954727573746564416760448201527f6772656761746f723a206f6e6c7920747275737465642041676772656761746f6064820152603960f91b608482015260a4016109e9565b6139a58888888888888888613d2f565b6069805467ffffffffffffffff19166001600160401b038881169182179092556000908152606d6020526040902085905560725416156139f957607280546fffffffffffffffffffffffffffffffff191690555b606c546040516333d6247d60e01b8152600481018790526001600160a01b03909116906333d6247d90602401600060405180830381600087803b158015613a3f57600080fd5b505af1158015613a53573d6000803e3d6000fd5b50506040518681523392506001600160401b03891691507f0c0ce073a7d7b5850c04ccc4b20ee7d3179d5f57d0ac44399565792c0f72fce790602001611a9e565b6073546001600160a01b03163314613aff5760405162461bcd60e51b815260206004820152602860248201527f50726f6f664f66456666696369656e63793a3a6f6e6c7941646d696e3a206f6e604482015267363c9030b236b4b760c11b60648201526084016109e9565b606a80546001600160a01b0319166001600160a01b0383169081179091556040519081527f61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca90602001610d63565b613b5561455c565b6001600160a01b038116613bd15760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084016109e9565b612575816145b6565b6040516001600160a01b03808516602483015283166044820152606481018290526139059085907f23b872dd00000000000000000000000000000000000000000000000000000000906084015b60408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152614ed5565b6072546001600160401b03600160401b82048116911611156122bb57607254600090613cc890600160401b90046001600160401b031660016159e7565b9050613cd381610afe565b1561257557607254600090600290613cf59084906001600160401b0316615a80565b613cff9190615c56565b613d0990836159e7565b9050613d1481610afe565b15613d2657613d2281611ab0565b5050565b613d2282611ab0565b600080613d3a612a31565b90506001600160401b038a1615613ec6576072546001600160401b03908116908b161115613df65760405162461bcd60e51b815260206004820152605d60248201527f50726f6f664f66456666696369656e63793a3a7665726966794261746368657360448201527f3a2070656e64696e6753746174654e756d206d757374206265206c657373206f60648201527f7220657175616c207468616e206c61737450656e64696e675374617465000000608482015260a4016109e9565b6001600160401b03808b1660009081526071602052604090206002810154815490945090918b8116600160401b9092041614613ec05760405162461bcd60e51b815260206004820152605160248201527f50726f6f664f66456666696369656e63793a3a7665726966794261746368657360448201527f3a20696e69744e756d4261746368206d757374206d617463682074686520706560648201527f6e64696e67207374617465206261746368000000000000000000000000000000608482015260a4016109e9565b50614033565b6001600160401b0389166000908152606d6020526040902054915081613f7a5760405162461bcd60e51b815260206004820152604860248201527f50726f6f664f66456666696369656e63793a3a7665726966794261746368657360448201527f3a20696e69744e756d426174636820737461746520726f6f7420646f6573206e60648201527f6f74206578697374000000000000000000000000000000000000000000000000608482015260a4016109e9565b806001600160401b0316896001600160401b031611156140335760405162461bcd60e51b815260206004820152606260248201527f50726f6f664f66456666696369656e63793a3a7665726966794261746368657360448201527f3a20696e69744e756d4261746368206d757374206265206c657373206f72206560648201527f7175616c207468616e2063757272656e744c61737456657269666965644261746084820152610c6d60f31b60a482015260c4016109e9565b806001600160401b0316886001600160401b0316116140e05760405162461bcd60e51b815260206004820152605c60248201527f50726f6f664f66456666696369656e63793a3a7665726966794261746368657360448201527f3a2066696e616c4e65774261746368206d75737420626520626967676572207460648201527f68616e2063757272656e744c6173745665726966696564426174636800000000608482015260a4016109e9565b60006140ef8a8a8a868b61093b565b905060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f00000016002836040516141249190615c7c565b602060405180830381855afa158015614141573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906141649190615a4e565b61416e9190615c98565b606b546040805160208101825283815290516343753b4d60e01b81529293506001600160a01b03909116916343753b4d916141b2918b918b918b9190600401615cac565b602060405180830381865afa1580156141cf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906141f39190615d26565b6142655760405162461bcd60e51b815260206004820152602f60248201527f50726f6f664f66456666696369656e63793a3a7665726966794261746368657360448201527f3a20494e56414c49445f50524f4f46000000000000000000000000000000000060648201526084016109e9565b6142a633614273858d615a80565b6001600160401b0316614284612720565b61428e9190615abf565b60655461010090046001600160a01b03169190614fbf565b505050505050505050505050565b60006142be612a31565b9050816000806142ce8484615a80565b6001600160401b031690505b836001600160401b0316836001600160401b031614614378576001600160401b0380851660009081526067602052604090206001810154909161070891614322911642615aa8565b111561435d57600181015461434790600160401b90046001600160401b031686615a80565b61435a906001600160401b031684615c0d565b92505b60010154600160401b90046001600160401b031692506142da565b60006143848383615aa8565b90508281101561442257600061439a8285615aa8565b6074549091505b60208211156143ed576d04ee2d6d415b85acef81000000006143c56020600b615e27565b6143cf9083615abf565b6143d99190615bf9565b90506143e6602083615aa8565b91506143a1565b6143f882600a615e36565b61440383600b615e36565b61440d9083615abf565b6144179190615bf9565b607455506144ce9050565b600061442e8483615aa8565b6074549091505b6020821115614481576d04ee2d6d415b85acef81000000006144596020600b615e27565b6144639083615abf565b61446d9190615bf9565b905061447a602083615aa8565b9150614435565b61448c82600a615e36565b61449783600b615e36565b6144a19083615abf565b6144ab9190615bf9565b9050806074546074546144be9190615abf565b6144c89190615bf9565b60745550505b505050505050565b600054610100900460ff166145535760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e6700000000000000000000000000000000000000000060648201526084016109e9565b6122bb336145b6565b6033546001600160a01b031633146122bb5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016109e9565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b607060009054906101000a90046001600160a01b03166001600160a01b0316632072f6c56040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561465857600080fd5b505af115801561466c573d6000803e3d6000fd5b505050506122bb615008565b60006001600160401b038a161561484e576072546001600160401b03908116908b16111561475a5760405162461bcd60e51b815260206004820152607160248201527f50726f6f664f66456666696369656e63793a3a70726f76654e6f6e446574657260448201527f6d696e697374696350656e64696e6753746174653a2070656e64696e6753746160648201527f74654e756d206d757374206265206c657373206f7220657175616c207468616e60848201527f206c61737450656e64696e67537461746500000000000000000000000000000060a482015260c4016109e9565b506001600160401b03808a1660009081526071602052604090206002810154815490928a8116600160401b90920416146148485760405162461bcd60e51b815260206004820152606560248201527f50726f6f664f66456666696369656e63793a3a70726f76654e6f6e446574657260448201527f6d696e697374696350656e64696e6753746174653a20696e69744e756d42617460648201527f6368206d757374206d61746368207468652070656e64696e672073746174652060848201527f626174636800000000000000000000000000000000000000000000000000000060a482015260c4016109e9565b506149d2565b506001600160401b0387166000908152606d6020526040902054806149015760405162461bcd60e51b815260206004820152605c60248201527f50726f6f664f66456666696369656e63793a3a70726f76654e6f6e446574657260448201527f6d696e697374696350656e64696e6753746174653a20696e69744e756d42617460648201527f636820737461746520726f6f7420646f6573206e6f7420657869737400000000608482015260a4016109e9565b6069546001600160401b0390811690891611156149d25760405162461bcd60e51b815260206004820152607660248201527f50726f6f664f66456666696369656e63793a3a70726f76654e6f6e446574657260448201527f6d696e697374696350656e64696e6753746174653a20696e69744e756d42617460648201527f6368206d757374206265206c657373206f7220657175616c207468616e20637560848201527f7272656e744c617374566572696669656442617463680000000000000000000060a482015260c4016109e9565b6072546001600160401b03908116908a1611801590614a025750896001600160401b0316896001600160401b0316115b8015614a2357506072546001600160401b03600160401b9091048116908a16115b614abb5760405162461bcd60e51b815260206004820152605460248201527f50726f6f664f66456666696369656e63793a3a70726f76654e6f6e446574657260448201527f6d696e697374696350656e64696e6753746174653a2066696e616c50656e646960648201527f6e6753746174654e756d20696e636f7272656374000000000000000000000000608482015260a4016109e9565b6001600160401b03898116600090815260716020526040902054600160401b9004811690881614614ba05760405162461bcd60e51b815260206004820152606f60248201527f50726f6f664f66456666696369656e63793a3a70726f76654e6f6e446574657260448201527f6d696e697374696350656e64696e6753746174653a2066696e616c4e6577426160648201527f746368206d75737420626520657175616c207468616e2063757272656e744c6160848201527f737456657269666965644261746368000000000000000000000000000000000060a482015260c4016109e9565b6000614baf898989858a61093b565b905060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001600283604051614be49190615c7c565b602060405180830381855afa158015614c01573d6000803e3d6000fd5b5050506040513d601f19601f82011682018060405250810190614c249190615a4e565b614c2e9190615c98565b606b546040805160208101825283815290516343753b4d60e01b81529293506001600160a01b03909116916343753b4d91614c72918a918a918a9190600401615cac565b602060405180830381865afa158015614c8f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190614cb39190615d26565b614d4b5760405162461bcd60e51b815260206004820152604360248201527f50726f6f664f66456666696369656e63793a3a70726f76654e6f6e446574657260448201527f6d696e697374696350656e64696e6753746174653a20494e56414c49445f505260648201527f4f4f460000000000000000000000000000000000000000000000000000000000608482015260a4016109e9565b6001600160401b038b166000908152607160205260409020600201548790036142a65760405162461bcd60e51b815260206004820152606760248201527f50726f6f664f66456666696369656e63793a3a70726f76654e6f6e446574657260448201527f6d696e697374696350656e64696e6753746174653a2073746f72656420726f6f60648201527f74206d75737420626520646966666572656e74207468616e206e65772073746160848201527f746520726f6f740000000000000000000000000000000000000000000000000060a482015260c4016109e9565b60655460ff16614ea05760405162461bcd60e51b815260206004820152603b60248201527f456d657267656e63794d616e616765723a3a6966456d657267656e637953746160448201527f74653a206f6e6c7920696620656d657267656e6379207374617465000000000060648201526084016109e9565b6065805460ff191690556040517f1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b390600090a1565b6000614f2a826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166150c49092919063ffffffff16565b805190915015614fba5780806020019051810190614f489190615d26565b614fba5760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f7420737563636565640000000000000000000000000000000000000000000060648201526084016109e9565b505050565b6040516001600160a01b038316602482015260448101829052614fba9084907fa9059cbb0000000000000000000000000000000000000000000000000000000090606401613c27565b60655460ff161561508c5760405162461bcd60e51b815260206004820152604260248201527f456d657267656e63794d616e616765723a3a69664e6f74456d657267656e637960448201527f53746174653a206f6e6c79206966206e6f7420656d657267656e637920737461606482015261746560f01b608482015260a4016109e9565b6065805460ff191660011790556040517f2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a549790600090a1565b60606150d384846000856150dd565b90505b9392505050565b6060824710156151555760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c000000000000000000000000000000000000000000000000000060648201526084016109e9565b6001600160a01b0385163b6151ac5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016109e9565b600080866001600160a01b031685876040516151c89190615c7c565b60006040518083038185875af1925050503d8060008114615205576040519150601f19603f3d011682016040523d82523d6000602084013e61520a565b606091505b509150915061521a828286615225565b979650505050505050565b606083156152345750816150d6565b8251156152445782518084602001fd5b8160405162461bcd60e51b81526004016109e991906152b6565b60005b83811015615279578181015183820152602001615261565b838111156139055750506000910152565b600081518084526152a281602086016020860161525e565b601f01601f19169290920160200192915050565b6020815260006150d6602083018461528a565b80356001600160401b03811681146152e057600080fd5b919050565b600080600080600060a086880312156152fd57600080fd5b615306866152c9565b9450615314602087016152c9565b94979496505050506040830135926060810135926080909101359150565b60006020828403121561534457600080fd5b6150d6826152c9565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b03811182821017156153855761538561534d565b60405290565b604051606081016001600160401b03811182821017156153855761538561534d565b604051601f8201601f191681016001600160401b03811182821017156153d5576153d561534d565b604052919050565b60006001600160401b038211156153f6576153f661534d565b5060051b60200190565b600082601f83011261541157600080fd5b81356001600160401b0381111561542a5761542a61534d565b61543d601f8201601f19166020016153ad565b81815284602083860101111561545257600080fd5b816020850160208301376000918101602001919091529392505050565b6000602080838503121561548257600080fd5b82356001600160401b038082111561549957600080fd5b818501915085601f8301126154ad57600080fd5b81356154c06154bb826153dd565b6153ad565b81815260059190911b830184019084810190888311156154df57600080fd5b8585015b83811015615578578035858111156154fb5760008081fd5b86016080818c03601f19018113156155135760008081fd5b61551b615363565b898301358881111561552d5760008081fd5b61553b8e8c83870101615400565b8252506040808401358b83015260606155558186016152c9565b828401526155648486016152c9565b9083015250855250509186019186016154e3565b5098975050505050505050565b806040810183101561559657600080fd5b92915050565b806080810183101561559657600080fd5b6000806000806000806000806101a0898b0312156155ca57600080fd5b6155d3896152c9565b97506155e160208a016152c9565b96506155ef60408a016152c9565b9550606089013594506080890135935061560c8a60a08b01615585565b925061561b8a60e08b0161559c565b915061562b8a6101608b01615585565b90509295985092959890939650565b6001600160a01b038116811461257557600080fd5b600080600080600080600080888a036101c081121561566d57600080fd5b89356156788161563a565b985060208a01356156888161563a565b975060408a01356156988161563a565b965060608a01356156a88161563a565b955060e0607f19820112156156bc57600080fd5b5060808901935061016089013592506101808901356001600160401b03808211156156e657600080fd5b6156f28c838d01615400565b93506101a08b013591508082111561570957600080fd5b506157168b828c01615400565b9150509295985092959890939650565b60006020828403121561573857600080fd5b81356150d68161563a565b60008060008060008060008060006101c08a8c03121561576257600080fd5b61576b8a6152c9565b985061577960208b016152c9565b975061578760408b016152c9565b965061579560608b016152c9565b955060808a0135945060a08a013593506157b28b60c08c01615585565b92506157c28b6101008c0161559c565b91506157d28b6101808c01615585565b90509295985092959850929598565b6000602082840312156157f357600080fd5b5035919050565b801515811461257557600080fd5b60006020828403121561581a57600080fd5b81356150d6816157fa565b60006020828403121561583757600080fd5b81356001600160401b0381111561584d57600080fd5b61585984828501615400565b949350505050565b6000602080838503121561587457600080fd5b82356001600160401b038082111561588b57600080fd5b818501915085601f83011261589f57600080fd5b81356158ad6154bb826153dd565b81815260059190911b830184019084810190888311156158cc57600080fd5b8585015b83811015615578578035858111156158e85760008081fd5b86016060818c03601f19018113156159005760008081fd5b61590861538b565b898301358881111561591a5760008081fd5b6159288e8c83870101615400565b8252506040808401358b8301526159408385016152c9565b90820152855250509186019186016158d0565b6000806040838503121561596657600080fd5b82356001600160401b0381111561597c57600080fd5b61598885828601615400565b95602094909401359450505050565b600181811c908216806159ab57607f821691505b6020821081036159cb57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b60006001600160401b03808316818516808303821115615a0957615a096159d1565b01949350505050565b634e487b7160e01b600052603260045260246000fd5b60006001600160401b03808316818103615a4457615a446159d1565b6001019392505050565b600060208284031215615a6057600080fd5b5051919050565b600060018201615a7957615a796159d1565b5060010190565b60006001600160401b0383811690831681811015615aa057615aa06159d1565b039392505050565b600082821015615aba57615aba6159d1565b500390565b6000816000190483118215151615615ad957615ad96159d1565b500290565b601f821115614fba57600081815260208120601f850160051c81016020861015615b055750805b601f850160051c820191505b818110156144ce57828155600101615b11565b81516001600160401b03811115615b3d57615b3d61534d565b615b5181615b4b8454615997565b84615ade565b602080601f831160018114615b865760008415615b6e5750858301515b600019600386901b1c1916600185901b1785556144ce565b600085815260208120601f198616915b82811015615bb557888601518255948401946001909101908401615b96565b5085821015615bd35787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b634e487b7160e01b600052601260045260246000fd5b600082615c0857615c08615be3565b500490565b60008219821115615c2057615c206159d1565b500190565b8381526001600160a01b0383166020820152606060408201526000615c4d606083018461528a565b95945050505050565b60006001600160401b0380841680615c7057615c70615be3565b92169190910492915050565b60008251615c8e81846020870161525e565b9190910192915050565b600082615ca757615ca7615be3565b500690565b61012081016040808784376000838201818152879190815b6002811015615ce457848483379084018281529284019290600101615cc4565b5050828760c0870137610100850181815286935091505b6001811015615d1a578251825260209283019290910190600101615cfb565b50505095945050505050565b600060208284031215615d3857600080fd5b81516150d6816157fa565b600181815b80851115615d7e578160001904821115615d6457615d646159d1565b80851615615d7157918102915b93841c9390800290615d48565b509250929050565b600082615d9557506001615596565b81615da257506000615596565b8160018114615db85760028114615dc257615dde565b6001915050615596565b60ff841115615dd357615dd36159d1565b50506001821b615596565b5060208310610133831016604e8410600b8410161715615e01575081810a615596565b615e0b8383615d43565b8060001904821115615e1f57615e1f6159d1565b029392505050565b60006150d660ff841683615d86565b60006150d68383615d8656fea26469706673582212200dcdd6e2039faa0c7979293eff6e1f1020ae168d989673e3ec2a0911663b442664736f6c634300080f0033",
}

// ProofofefficiencyABI is the input ABI used to generate the binding from.
// Deprecated: Use ProofofefficiencyMetaData.ABI instead.
var ProofofefficiencyABI = ProofofefficiencyMetaData.ABI

// ProofofefficiencyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ProofofefficiencyMetaData.Bin instead.
var ProofofefficiencyBin = ProofofefficiencyMetaData.Bin

// DeployProofofefficiency deploys a new Ethereum contract, binding an instance of Proofofefficiency to it.
func DeployProofofefficiency(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Proofofefficiency, error) {
	parsed, err := ProofofefficiencyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ProofofefficiencyBin), backend)
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

// FORCEBATCHTIMEOUT is a free data retrieval call binding the contract method 0xab9fc5ef.
//
// Solidity: function FORCE_BATCH_TIMEOUT() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) FORCEBATCHTIMEOUT(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "FORCE_BATCH_TIMEOUT")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// FORCEBATCHTIMEOUT is a free data retrieval call binding the contract method 0xab9fc5ef.
//
// Solidity: function FORCE_BATCH_TIMEOUT() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) FORCEBATCHTIMEOUT() (uint64, error) {
	return _Proofofefficiency.Contract.FORCEBATCHTIMEOUT(&_Proofofefficiency.CallOpts)
}

// FORCEBATCHTIMEOUT is a free data retrieval call binding the contract method 0xab9fc5ef.
//
// Solidity: function FORCE_BATCH_TIMEOUT() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) FORCEBATCHTIMEOUT() (uint64, error) {
	return _Proofofefficiency.Contract.FORCEBATCHTIMEOUT(&_Proofofefficiency.CallOpts)
}

// HALTAGGREGATIONTIMEOUT is a free data retrieval call binding the contract method 0x8b48931e.
//
// Solidity: function HALT_AGGREGATION_TIMEOUT() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) HALTAGGREGATIONTIMEOUT(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "HALT_AGGREGATION_TIMEOUT")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// HALTAGGREGATIONTIMEOUT is a free data retrieval call binding the contract method 0x8b48931e.
//
// Solidity: function HALT_AGGREGATION_TIMEOUT() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) HALTAGGREGATIONTIMEOUT() (uint64, error) {
	return _Proofofefficiency.Contract.HALTAGGREGATIONTIMEOUT(&_Proofofefficiency.CallOpts)
}

// HALTAGGREGATIONTIMEOUT is a free data retrieval call binding the contract method 0x8b48931e.
//
// Solidity: function HALT_AGGREGATION_TIMEOUT() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) HALTAGGREGATIONTIMEOUT() (uint64, error) {
	return _Proofofefficiency.Contract.HALTAGGREGATIONTIMEOUT(&_Proofofefficiency.CallOpts)
}

// MAXTRANSACTIONSBYTELENGTH is a free data retrieval call binding the contract method 0x2d0889d3.
//
// Solidity: function MAX_TRANSACTIONS_BYTE_LENGTH() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCaller) MAXTRANSACTIONSBYTELENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "MAX_TRANSACTIONS_BYTE_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXTRANSACTIONSBYTELENGTH is a free data retrieval call binding the contract method 0x2d0889d3.
//
// Solidity: function MAX_TRANSACTIONS_BYTE_LENGTH() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencySession) MAXTRANSACTIONSBYTELENGTH() (*big.Int, error) {
	return _Proofofefficiency.Contract.MAXTRANSACTIONSBYTELENGTH(&_Proofofefficiency.CallOpts)
}

// MAXTRANSACTIONSBYTELENGTH is a free data retrieval call binding the contract method 0x2d0889d3.
//
// Solidity: function MAX_TRANSACTIONS_BYTE_LENGTH() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCallerSession) MAXTRANSACTIONSBYTELENGTH() (*big.Int, error) {
	return _Proofofefficiency.Contract.MAXTRANSACTIONSBYTELENGTH(&_Proofofefficiency.CallOpts)
}

// MAXVERIFYBATCHES is a free data retrieval call binding the contract method 0xe217cfd6.
//
// Solidity: function MAX_VERIFY_BATCHES() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) MAXVERIFYBATCHES(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "MAX_VERIFY_BATCHES")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MAXVERIFYBATCHES is a free data retrieval call binding the contract method 0xe217cfd6.
//
// Solidity: function MAX_VERIFY_BATCHES() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) MAXVERIFYBATCHES() (uint64, error) {
	return _Proofofefficiency.Contract.MAXVERIFYBATCHES(&_Proofofefficiency.CallOpts)
}

// MAXVERIFYBATCHES is a free data retrieval call binding the contract method 0xe217cfd6.
//
// Solidity: function MAX_VERIFY_BATCHES() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) MAXVERIFYBATCHES() (uint64, error) {
	return _Proofofefficiency.Contract.MAXVERIFYBATCHES(&_Proofofefficiency.CallOpts)
}

// MULTIPLIERBATCHFEE is a free data retrieval call binding the contract method 0xf1d7b21c.
//
// Solidity: function MULTIPLIER_BATCH_FEE() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCaller) MULTIPLIERBATCHFEE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "MULTIPLIER_BATCH_FEE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MULTIPLIERBATCHFEE is a free data retrieval call binding the contract method 0xf1d7b21c.
//
// Solidity: function MULTIPLIER_BATCH_FEE() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencySession) MULTIPLIERBATCHFEE() (*big.Int, error) {
	return _Proofofefficiency.Contract.MULTIPLIERBATCHFEE(&_Proofofefficiency.CallOpts)
}

// MULTIPLIERBATCHFEE is a free data retrieval call binding the contract method 0xf1d7b21c.
//
// Solidity: function MULTIPLIER_BATCH_FEE() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCallerSession) MULTIPLIERBATCHFEE() (*big.Int, error) {
	return _Proofofefficiency.Contract.MULTIPLIERBATCHFEE(&_Proofofefficiency.CallOpts)
}

// VERIFYBATCHTIMETARGET is a free data retrieval call binding the contract method 0x137f1edf.
//
// Solidity: function VERIFY_BATCH_TIME_TARGET() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) VERIFYBATCHTIMETARGET(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "VERIFY_BATCH_TIME_TARGET")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// VERIFYBATCHTIMETARGET is a free data retrieval call binding the contract method 0x137f1edf.
//
// Solidity: function VERIFY_BATCH_TIME_TARGET() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) VERIFYBATCHTIMETARGET() (uint64, error) {
	return _Proofofefficiency.Contract.VERIFYBATCHTIMETARGET(&_Proofofefficiency.CallOpts)
}

// VERIFYBATCHTIMETARGET is a free data retrieval call binding the contract method 0x137f1edf.
//
// Solidity: function VERIFY_BATCH_TIME_TARGET() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) VERIFYBATCHTIMETARGET() (uint64, error) {
	return _Proofofefficiency.Contract.VERIFYBATCHTIMETARGET(&_Proofofefficiency.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Proofofefficiency *ProofofefficiencySession) Admin() (common.Address, error) {
	return _Proofofefficiency.Contract.Admin(&_Proofofefficiency.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCallerSession) Admin() (common.Address, error) {
	return _Proofofefficiency.Contract.Admin(&_Proofofefficiency.CallOpts)
}

// BatchFee is a free data retrieval call binding the contract method 0xf8b823e4.
//
// Solidity: function batchFee() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCaller) BatchFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "batchFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BatchFee is a free data retrieval call binding the contract method 0xf8b823e4.
//
// Solidity: function batchFee() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencySession) BatchFee() (*big.Int, error) {
	return _Proofofefficiency.Contract.BatchFee(&_Proofofefficiency.CallOpts)
}

// BatchFee is a free data retrieval call binding the contract method 0xf8b823e4.
//
// Solidity: function batchFee() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCallerSession) BatchFee() (*big.Int, error) {
	return _Proofofefficiency.Contract.BatchFee(&_Proofofefficiency.CallOpts)
}

// BatchNumToStateRoot is a free data retrieval call binding the contract method 0x5392c5e0.
//
// Solidity: function batchNumToStateRoot(uint64 ) view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencyCaller) BatchNumToStateRoot(opts *bind.CallOpts, arg0 uint64) ([32]byte, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "batchNumToStateRoot", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BatchNumToStateRoot is a free data retrieval call binding the contract method 0x5392c5e0.
//
// Solidity: function batchNumToStateRoot(uint64 ) view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencySession) BatchNumToStateRoot(arg0 uint64) ([32]byte, error) {
	return _Proofofefficiency.Contract.BatchNumToStateRoot(&_Proofofefficiency.CallOpts, arg0)
}

// BatchNumToStateRoot is a free data retrieval call binding the contract method 0x5392c5e0.
//
// Solidity: function batchNumToStateRoot(uint64 ) view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencyCallerSession) BatchNumToStateRoot(arg0 uint64) ([32]byte, error) {
	return _Proofofefficiency.Contract.BatchNumToStateRoot(&_Proofofefficiency.CallOpts, arg0)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCaller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Proofofefficiency *ProofofefficiencySession) BridgeAddress() (common.Address, error) {
	return _Proofofefficiency.Contract.BridgeAddress(&_Proofofefficiency.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCallerSession) BridgeAddress() (common.Address, error) {
	return _Proofofefficiency.Contract.BridgeAddress(&_Proofofefficiency.CallOpts)
}

// CalculateRewardPerBatch is a free data retrieval call binding the contract method 0x99f5634e.
//
// Solidity: function calculateRewardPerBatch() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCaller) CalculateRewardPerBatch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "calculateRewardPerBatch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateRewardPerBatch is a free data retrieval call binding the contract method 0x99f5634e.
//
// Solidity: function calculateRewardPerBatch() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencySession) CalculateRewardPerBatch() (*big.Int, error) {
	return _Proofofefficiency.Contract.CalculateRewardPerBatch(&_Proofofefficiency.CallOpts)
}

// CalculateRewardPerBatch is a free data retrieval call binding the contract method 0x99f5634e.
//
// Solidity: function calculateRewardPerBatch() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCallerSession) CalculateRewardPerBatch() (*big.Int, error) {
	return _Proofofefficiency.Contract.CalculateRewardPerBatch(&_Proofofefficiency.CallOpts)
}

// ChainID is a free data retrieval call binding the contract method 0xadc879e9.
//
// Solidity: function chainID() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) ChainID(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "chainID")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ChainID is a free data retrieval call binding the contract method 0xadc879e9.
//
// Solidity: function chainID() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) ChainID() (uint64, error) {
	return _Proofofefficiency.Contract.ChainID(&_Proofofefficiency.CallOpts)
}

// ChainID is a free data retrieval call binding the contract method 0xadc879e9.
//
// Solidity: function chainID() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) ChainID() (uint64, error) {
	return _Proofofefficiency.Contract.ChainID(&_Proofofefficiency.CallOpts)
}

// ForceBatchAllowed is a free data retrieval call binding the contract method 0xd8f54db0.
//
// Solidity: function forceBatchAllowed() view returns(bool)
func (_Proofofefficiency *ProofofefficiencyCaller) ForceBatchAllowed(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "forceBatchAllowed")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ForceBatchAllowed is a free data retrieval call binding the contract method 0xd8f54db0.
//
// Solidity: function forceBatchAllowed() view returns(bool)
func (_Proofofefficiency *ProofofefficiencySession) ForceBatchAllowed() (bool, error) {
	return _Proofofefficiency.Contract.ForceBatchAllowed(&_Proofofefficiency.CallOpts)
}

// ForceBatchAllowed is a free data retrieval call binding the contract method 0xd8f54db0.
//
// Solidity: function forceBatchAllowed() view returns(bool)
func (_Proofofefficiency *ProofofefficiencyCallerSession) ForceBatchAllowed() (bool, error) {
	return _Proofofefficiency.Contract.ForceBatchAllowed(&_Proofofefficiency.CallOpts)
}

// ForcedBatches is a free data retrieval call binding the contract method 0x6b8616ce.
//
// Solidity: function forcedBatches(uint64 ) view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencyCaller) ForcedBatches(opts *bind.CallOpts, arg0 uint64) ([32]byte, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "forcedBatches", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ForcedBatches is a free data retrieval call binding the contract method 0x6b8616ce.
//
// Solidity: function forcedBatches(uint64 ) view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencySession) ForcedBatches(arg0 uint64) ([32]byte, error) {
	return _Proofofefficiency.Contract.ForcedBatches(&_Proofofefficiency.CallOpts, arg0)
}

// ForcedBatches is a free data retrieval call binding the contract method 0x6b8616ce.
//
// Solidity: function forcedBatches(uint64 ) view returns(bytes32)
func (_Proofofefficiency *ProofofefficiencyCallerSession) ForcedBatches(arg0 uint64) ([32]byte, error) {
	return _Proofofefficiency.Contract.ForcedBatches(&_Proofofefficiency.CallOpts, arg0)
}

// GetCurrentBatchFee is a free data retrieval call binding the contract method 0x9f0d039d.
//
// Solidity: function getCurrentBatchFee() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCaller) GetCurrentBatchFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "getCurrentBatchFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentBatchFee is a free data retrieval call binding the contract method 0x9f0d039d.
//
// Solidity: function getCurrentBatchFee() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencySession) GetCurrentBatchFee() (*big.Int, error) {
	return _Proofofefficiency.Contract.GetCurrentBatchFee(&_Proofofefficiency.CallOpts)
}

// GetCurrentBatchFee is a free data retrieval call binding the contract method 0x9f0d039d.
//
// Solidity: function getCurrentBatchFee() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCallerSession) GetCurrentBatchFee() (*big.Int, error) {
	return _Proofofefficiency.Contract.GetCurrentBatchFee(&_Proofofefficiency.CallOpts)
}

// GetInputSnarkBytes is a free data retrieval call binding the contract method 0x220d7899.
//
// Solidity: function getInputSnarkBytes(uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 oldStateRoot, bytes32 newStateRoot) view returns(bytes)
func (_Proofofefficiency *ProofofefficiencyCaller) GetInputSnarkBytes(opts *bind.CallOpts, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, oldStateRoot [32]byte, newStateRoot [32]byte) ([]byte, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "getInputSnarkBytes", initNumBatch, finalNewBatch, newLocalExitRoot, oldStateRoot, newStateRoot)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetInputSnarkBytes is a free data retrieval call binding the contract method 0x220d7899.
//
// Solidity: function getInputSnarkBytes(uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 oldStateRoot, bytes32 newStateRoot) view returns(bytes)
func (_Proofofefficiency *ProofofefficiencySession) GetInputSnarkBytes(initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, oldStateRoot [32]byte, newStateRoot [32]byte) ([]byte, error) {
	return _Proofofefficiency.Contract.GetInputSnarkBytes(&_Proofofefficiency.CallOpts, initNumBatch, finalNewBatch, newLocalExitRoot, oldStateRoot, newStateRoot)
}

// GetInputSnarkBytes is a free data retrieval call binding the contract method 0x220d7899.
//
// Solidity: function getInputSnarkBytes(uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 oldStateRoot, bytes32 newStateRoot) view returns(bytes)
func (_Proofofefficiency *ProofofefficiencyCallerSession) GetInputSnarkBytes(initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, oldStateRoot [32]byte, newStateRoot [32]byte) ([]byte, error) {
	return _Proofofefficiency.Contract.GetInputSnarkBytes(&_Proofofefficiency.CallOpts, initNumBatch, finalNewBatch, newLocalExitRoot, oldStateRoot, newStateRoot)
}

// GetLastVerifiedBatch is a free data retrieval call binding the contract method 0xc0ed84e0.
//
// Solidity: function getLastVerifiedBatch() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) GetLastVerifiedBatch(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "getLastVerifiedBatch")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetLastVerifiedBatch is a free data retrieval call binding the contract method 0xc0ed84e0.
//
// Solidity: function getLastVerifiedBatch() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) GetLastVerifiedBatch() (uint64, error) {
	return _Proofofefficiency.Contract.GetLastVerifiedBatch(&_Proofofefficiency.CallOpts)
}

// GetLastVerifiedBatch is a free data retrieval call binding the contract method 0xc0ed84e0.
//
// Solidity: function getLastVerifiedBatch() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) GetLastVerifiedBatch() (uint64, error) {
	return _Proofofefficiency.Contract.GetLastVerifiedBatch(&_Proofofefficiency.CallOpts)
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

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Proofofefficiency *ProofofefficiencyCaller) IsEmergencyState(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "isEmergencyState")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Proofofefficiency *ProofofefficiencySession) IsEmergencyState() (bool, error) {
	return _Proofofefficiency.Contract.IsEmergencyState(&_Proofofefficiency.CallOpts)
}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Proofofefficiency *ProofofefficiencyCallerSession) IsEmergencyState() (bool, error) {
	return _Proofofefficiency.Contract.IsEmergencyState(&_Proofofefficiency.CallOpts)
}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x383b3be8.
//
// Solidity: function isPendingStateConsolidable(uint64 pendingStateNum) view returns(bool)
func (_Proofofefficiency *ProofofefficiencyCaller) IsPendingStateConsolidable(opts *bind.CallOpts, pendingStateNum uint64) (bool, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "isPendingStateConsolidable", pendingStateNum)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x383b3be8.
//
// Solidity: function isPendingStateConsolidable(uint64 pendingStateNum) view returns(bool)
func (_Proofofefficiency *ProofofefficiencySession) IsPendingStateConsolidable(pendingStateNum uint64) (bool, error) {
	return _Proofofefficiency.Contract.IsPendingStateConsolidable(&_Proofofefficiency.CallOpts, pendingStateNum)
}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x383b3be8.
//
// Solidity: function isPendingStateConsolidable(uint64 pendingStateNum) view returns(bool)
func (_Proofofefficiency *ProofofefficiencyCallerSession) IsPendingStateConsolidable(pendingStateNum uint64) (bool, error) {
	return _Proofofefficiency.Contract.IsPendingStateConsolidable(&_Proofofefficiency.CallOpts, pendingStateNum)
}

// LastBatchSequenced is a free data retrieval call binding the contract method 0x423fa856.
//
// Solidity: function lastBatchSequenced() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) LastBatchSequenced(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "lastBatchSequenced")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastBatchSequenced is a free data retrieval call binding the contract method 0x423fa856.
//
// Solidity: function lastBatchSequenced() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) LastBatchSequenced() (uint64, error) {
	return _Proofofefficiency.Contract.LastBatchSequenced(&_Proofofefficiency.CallOpts)
}

// LastBatchSequenced is a free data retrieval call binding the contract method 0x423fa856.
//
// Solidity: function lastBatchSequenced() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) LastBatchSequenced() (uint64, error) {
	return _Proofofefficiency.Contract.LastBatchSequenced(&_Proofofefficiency.CallOpts)
}

// LastForceBatch is a free data retrieval call binding the contract method 0xe7a7ed02.
//
// Solidity: function lastForceBatch() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) LastForceBatch(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "lastForceBatch")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastForceBatch is a free data retrieval call binding the contract method 0xe7a7ed02.
//
// Solidity: function lastForceBatch() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) LastForceBatch() (uint64, error) {
	return _Proofofefficiency.Contract.LastForceBatch(&_Proofofefficiency.CallOpts)
}

// LastForceBatch is a free data retrieval call binding the contract method 0xe7a7ed02.
//
// Solidity: function lastForceBatch() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) LastForceBatch() (uint64, error) {
	return _Proofofefficiency.Contract.LastForceBatch(&_Proofofefficiency.CallOpts)
}

// LastForceBatchSequenced is a free data retrieval call binding the contract method 0x45605267.
//
// Solidity: function lastForceBatchSequenced() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) LastForceBatchSequenced(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "lastForceBatchSequenced")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastForceBatchSequenced is a free data retrieval call binding the contract method 0x45605267.
//
// Solidity: function lastForceBatchSequenced() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) LastForceBatchSequenced() (uint64, error) {
	return _Proofofefficiency.Contract.LastForceBatchSequenced(&_Proofofefficiency.CallOpts)
}

// LastForceBatchSequenced is a free data retrieval call binding the contract method 0x45605267.
//
// Solidity: function lastForceBatchSequenced() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) LastForceBatchSequenced() (uint64, error) {
	return _Proofofefficiency.Contract.LastForceBatchSequenced(&_Proofofefficiency.CallOpts)
}

// LastPendingState is a free data retrieval call binding the contract method 0x458c0477.
//
// Solidity: function lastPendingState() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) LastPendingState(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "lastPendingState")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastPendingState is a free data retrieval call binding the contract method 0x458c0477.
//
// Solidity: function lastPendingState() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) LastPendingState() (uint64, error) {
	return _Proofofefficiency.Contract.LastPendingState(&_Proofofefficiency.CallOpts)
}

// LastPendingState is a free data retrieval call binding the contract method 0x458c0477.
//
// Solidity: function lastPendingState() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) LastPendingState() (uint64, error) {
	return _Proofofefficiency.Contract.LastPendingState(&_Proofofefficiency.CallOpts)
}

// LastPendingStateConsolidated is a free data retrieval call binding the contract method 0x4a1a89a7.
//
// Solidity: function lastPendingStateConsolidated() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) LastPendingStateConsolidated(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "lastPendingStateConsolidated")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastPendingStateConsolidated is a free data retrieval call binding the contract method 0x4a1a89a7.
//
// Solidity: function lastPendingStateConsolidated() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) LastPendingStateConsolidated() (uint64, error) {
	return _Proofofefficiency.Contract.LastPendingStateConsolidated(&_Proofofefficiency.CallOpts)
}

// LastPendingStateConsolidated is a free data retrieval call binding the contract method 0x4a1a89a7.
//
// Solidity: function lastPendingStateConsolidated() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) LastPendingStateConsolidated() (uint64, error) {
	return _Proofofefficiency.Contract.LastPendingStateConsolidated(&_Proofofefficiency.CallOpts)
}

// LastTimestamp is a free data retrieval call binding the contract method 0x19d8ac61.
//
// Solidity: function lastTimestamp() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) LastTimestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "lastTimestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastTimestamp is a free data retrieval call binding the contract method 0x19d8ac61.
//
// Solidity: function lastTimestamp() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) LastTimestamp() (uint64, error) {
	return _Proofofefficiency.Contract.LastTimestamp(&_Proofofefficiency.CallOpts)
}

// LastTimestamp is a free data retrieval call binding the contract method 0x19d8ac61.
//
// Solidity: function lastTimestamp() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) LastTimestamp() (uint64, error) {
	return _Proofofefficiency.Contract.LastTimestamp(&_Proofofefficiency.CallOpts)
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

// NetworkName is a free data retrieval call binding the contract method 0x107bf28c.
//
// Solidity: function networkName() view returns(string)
func (_Proofofefficiency *ProofofefficiencyCaller) NetworkName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "networkName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NetworkName is a free data retrieval call binding the contract method 0x107bf28c.
//
// Solidity: function networkName() view returns(string)
func (_Proofofefficiency *ProofofefficiencySession) NetworkName() (string, error) {
	return _Proofofefficiency.Contract.NetworkName(&_Proofofefficiency.CallOpts)
}

// NetworkName is a free data retrieval call binding the contract method 0x107bf28c.
//
// Solidity: function networkName() view returns(string)
func (_Proofofefficiency *ProofofefficiencyCallerSession) NetworkName() (string, error) {
	return _Proofofefficiency.Contract.NetworkName(&_Proofofefficiency.CallOpts)
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

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) PendingStateTimeout(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "pendingStateTimeout")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) PendingStateTimeout() (uint64, error) {
	return _Proofofefficiency.Contract.PendingStateTimeout(&_Proofofefficiency.CallOpts)
}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) PendingStateTimeout() (uint64, error) {
	return _Proofofefficiency.Contract.PendingStateTimeout(&_Proofofefficiency.CallOpts)
}

// PendingStateTransitions is a free data retrieval call binding the contract method 0x837a4738.
//
// Solidity: function pendingStateTransitions(uint256 ) view returns(uint64 timestamp, uint64 lastVerifiedBatch, bytes32 exitRoot, bytes32 stateRoot)
func (_Proofofefficiency *ProofofefficiencyCaller) PendingStateTransitions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Timestamp         uint64
	LastVerifiedBatch uint64
	ExitRoot          [32]byte
	StateRoot         [32]byte
}, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "pendingStateTransitions", arg0)

	outstruct := new(struct {
		Timestamp         uint64
		LastVerifiedBatch uint64
		ExitRoot          [32]byte
		StateRoot         [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Timestamp = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.LastVerifiedBatch = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.ExitRoot = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.StateRoot = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// PendingStateTransitions is a free data retrieval call binding the contract method 0x837a4738.
//
// Solidity: function pendingStateTransitions(uint256 ) view returns(uint64 timestamp, uint64 lastVerifiedBatch, bytes32 exitRoot, bytes32 stateRoot)
func (_Proofofefficiency *ProofofefficiencySession) PendingStateTransitions(arg0 *big.Int) (struct {
	Timestamp         uint64
	LastVerifiedBatch uint64
	ExitRoot          [32]byte
	StateRoot         [32]byte
}, error) {
	return _Proofofefficiency.Contract.PendingStateTransitions(&_Proofofefficiency.CallOpts, arg0)
}

// PendingStateTransitions is a free data retrieval call binding the contract method 0x837a4738.
//
// Solidity: function pendingStateTransitions(uint256 ) view returns(uint64 timestamp, uint64 lastVerifiedBatch, bytes32 exitRoot, bytes32 stateRoot)
func (_Proofofefficiency *ProofofefficiencyCallerSession) PendingStateTransitions(arg0 *big.Int) (struct {
	Timestamp         uint64
	LastVerifiedBatch uint64
	ExitRoot          [32]byte
	StateRoot         [32]byte
}, error) {
	return _Proofofefficiency.Contract.PendingStateTransitions(&_Proofofefficiency.CallOpts, arg0)
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

// SequencedBatches is a free data retrieval call binding the contract method 0xb4d63f58.
//
// Solidity: function sequencedBatches(uint64 ) view returns(bytes32 accInputHash, uint64 sequencedTimestamp, uint64 previousLastBatchSequenced)
func (_Proofofefficiency *ProofofefficiencyCaller) SequencedBatches(opts *bind.CallOpts, arg0 uint64) (struct {
	AccInputHash               [32]byte
	SequencedTimestamp         uint64
	PreviousLastBatchSequenced uint64
}, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "sequencedBatches", arg0)

	outstruct := new(struct {
		AccInputHash               [32]byte
		SequencedTimestamp         uint64
		PreviousLastBatchSequenced uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AccInputHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.SequencedTimestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.PreviousLastBatchSequenced = *abi.ConvertType(out[2], new(uint64)).(*uint64)

	return *outstruct, err

}

// SequencedBatches is a free data retrieval call binding the contract method 0xb4d63f58.
//
// Solidity: function sequencedBatches(uint64 ) view returns(bytes32 accInputHash, uint64 sequencedTimestamp, uint64 previousLastBatchSequenced)
func (_Proofofefficiency *ProofofefficiencySession) SequencedBatches(arg0 uint64) (struct {
	AccInputHash               [32]byte
	SequencedTimestamp         uint64
	PreviousLastBatchSequenced uint64
}, error) {
	return _Proofofefficiency.Contract.SequencedBatches(&_Proofofefficiency.CallOpts, arg0)
}

// SequencedBatches is a free data retrieval call binding the contract method 0xb4d63f58.
//
// Solidity: function sequencedBatches(uint64 ) view returns(bytes32 accInputHash, uint64 sequencedTimestamp, uint64 previousLastBatchSequenced)
func (_Proofofefficiency *ProofofefficiencyCallerSession) SequencedBatches(arg0 uint64) (struct {
	AccInputHash               [32]byte
	SequencedTimestamp         uint64
	PreviousLastBatchSequenced uint64
}, error) {
	return _Proofofefficiency.Contract.SequencedBatches(&_Proofofefficiency.CallOpts, arg0)
}

// TrustedAggregator is a free data retrieval call binding the contract method 0x29878983.
//
// Solidity: function trustedAggregator() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCaller) TrustedAggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "trustedAggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TrustedAggregator is a free data retrieval call binding the contract method 0x29878983.
//
// Solidity: function trustedAggregator() view returns(address)
func (_Proofofefficiency *ProofofefficiencySession) TrustedAggregator() (common.Address, error) {
	return _Proofofefficiency.Contract.TrustedAggregator(&_Proofofefficiency.CallOpts)
}

// TrustedAggregator is a free data retrieval call binding the contract method 0x29878983.
//
// Solidity: function trustedAggregator() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCallerSession) TrustedAggregator() (common.Address, error) {
	return _Proofofefficiency.Contract.TrustedAggregator(&_Proofofefficiency.CallOpts)
}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCaller) TrustedAggregatorTimeout(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "trustedAggregatorTimeout")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencySession) TrustedAggregatorTimeout() (uint64, error) {
	return _Proofofefficiency.Contract.TrustedAggregatorTimeout(&_Proofofefficiency.CallOpts)
}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Proofofefficiency *ProofofefficiencyCallerSession) TrustedAggregatorTimeout() (uint64, error) {
	return _Proofofefficiency.Contract.TrustedAggregatorTimeout(&_Proofofefficiency.CallOpts)
}

// TrustedSequencer is a free data retrieval call binding the contract method 0xcfa8ed47.
//
// Solidity: function trustedSequencer() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCaller) TrustedSequencer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "trustedSequencer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TrustedSequencer is a free data retrieval call binding the contract method 0xcfa8ed47.
//
// Solidity: function trustedSequencer() view returns(address)
func (_Proofofefficiency *ProofofefficiencySession) TrustedSequencer() (common.Address, error) {
	return _Proofofefficiency.Contract.TrustedSequencer(&_Proofofefficiency.CallOpts)
}

// TrustedSequencer is a free data retrieval call binding the contract method 0xcfa8ed47.
//
// Solidity: function trustedSequencer() view returns(address)
func (_Proofofefficiency *ProofofefficiencyCallerSession) TrustedSequencer() (common.Address, error) {
	return _Proofofefficiency.Contract.TrustedSequencer(&_Proofofefficiency.CallOpts)
}

// TrustedSequencerURL is a free data retrieval call binding the contract method 0x542028d5.
//
// Solidity: function trustedSequencerURL() view returns(string)
func (_Proofofefficiency *ProofofefficiencyCaller) TrustedSequencerURL(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "trustedSequencerURL")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TrustedSequencerURL is a free data retrieval call binding the contract method 0x542028d5.
//
// Solidity: function trustedSequencerURL() view returns(string)
func (_Proofofefficiency *ProofofefficiencySession) TrustedSequencerURL() (string, error) {
	return _Proofofefficiency.Contract.TrustedSequencerURL(&_Proofofefficiency.CallOpts)
}

// TrustedSequencerURL is a free data retrieval call binding the contract method 0x542028d5.
//
// Solidity: function trustedSequencerURL() view returns(string)
func (_Proofofefficiency *ProofofefficiencyCallerSession) TrustedSequencerURL() (string, error) {
	return _Proofofefficiency.Contract.TrustedSequencerURL(&_Proofofefficiency.CallOpts)
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x7215541a.
//
// Solidity: function activateEmergencyState(uint64 sequencedBatchNum) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) ActivateEmergencyState(opts *bind.TransactOpts, sequencedBatchNum uint64) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "activateEmergencyState", sequencedBatchNum)
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x7215541a.
//
// Solidity: function activateEmergencyState(uint64 sequencedBatchNum) returns()
func (_Proofofefficiency *ProofofefficiencySession) ActivateEmergencyState(sequencedBatchNum uint64) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.ActivateEmergencyState(&_Proofofefficiency.TransactOpts, sequencedBatchNum)
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x7215541a.
//
// Solidity: function activateEmergencyState(uint64 sequencedBatchNum) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) ActivateEmergencyState(sequencedBatchNum uint64) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.ActivateEmergencyState(&_Proofofefficiency.TransactOpts, sequencedBatchNum)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x4a910e6a.
//
// Solidity: function consolidatePendingState(uint64 pendingStateNum) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) ConsolidatePendingState(opts *bind.TransactOpts, pendingStateNum uint64) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "consolidatePendingState", pendingStateNum)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x4a910e6a.
//
// Solidity: function consolidatePendingState(uint64 pendingStateNum) returns()
func (_Proofofefficiency *ProofofefficiencySession) ConsolidatePendingState(pendingStateNum uint64) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.ConsolidatePendingState(&_Proofofefficiency.TransactOpts, pendingStateNum)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x4a910e6a.
//
// Solidity: function consolidatePendingState(uint64 pendingStateNum) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) ConsolidatePendingState(pendingStateNum uint64) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.ConsolidatePendingState(&_Proofofefficiency.TransactOpts, pendingStateNum)
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) DeactivateEmergencyState(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "deactivateEmergencyState")
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Proofofefficiency *ProofofefficiencySession) DeactivateEmergencyState() (*types.Transaction, error) {
	return _Proofofefficiency.Contract.DeactivateEmergencyState(&_Proofofefficiency.TransactOpts)
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) DeactivateEmergencyState() (*types.Transaction, error) {
	return _Proofofefficiency.Contract.DeactivateEmergencyState(&_Proofofefficiency.TransactOpts)
}

// ForceBatch is a paid mutator transaction binding the contract method 0xeaeb077b.
//
// Solidity: function forceBatch(bytes transactions, uint256 maticAmount) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) ForceBatch(opts *bind.TransactOpts, transactions []byte, maticAmount *big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "forceBatch", transactions, maticAmount)
}

// ForceBatch is a paid mutator transaction binding the contract method 0xeaeb077b.
//
// Solidity: function forceBatch(bytes transactions, uint256 maticAmount) returns()
func (_Proofofefficiency *ProofofefficiencySession) ForceBatch(transactions []byte, maticAmount *big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.ForceBatch(&_Proofofefficiency.TransactOpts, transactions, maticAmount)
}

// ForceBatch is a paid mutator transaction binding the contract method 0xeaeb077b.
//
// Solidity: function forceBatch(bytes transactions, uint256 maticAmount) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) ForceBatch(transactions []byte, maticAmount *big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.ForceBatch(&_Proofofefficiency.TransactOpts, transactions, maticAmount)
}

// Initialize is a paid mutator transaction binding the contract method 0x60943d6a.
//
// Solidity: function initialize(address _globalExitRootManager, address _matic, address _rollupVerifier, address _bridgeAddress, (address,uint64,address,uint64,bool,address,uint64) initializePackedParameters, bytes32 genesisRoot, string _trustedSequencerURL, string _networkName) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) Initialize(opts *bind.TransactOpts, _globalExitRootManager common.Address, _matic common.Address, _rollupVerifier common.Address, _bridgeAddress common.Address, initializePackedParameters ProofOfEfficiencyInitializePackedParameters, genesisRoot [32]byte, _trustedSequencerURL string, _networkName string) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "initialize", _globalExitRootManager, _matic, _rollupVerifier, _bridgeAddress, initializePackedParameters, genesisRoot, _trustedSequencerURL, _networkName)
}

// Initialize is a paid mutator transaction binding the contract method 0x60943d6a.
//
// Solidity: function initialize(address _globalExitRootManager, address _matic, address _rollupVerifier, address _bridgeAddress, (address,uint64,address,uint64,bool,address,uint64) initializePackedParameters, bytes32 genesisRoot, string _trustedSequencerURL, string _networkName) returns()
func (_Proofofefficiency *ProofofefficiencySession) Initialize(_globalExitRootManager common.Address, _matic common.Address, _rollupVerifier common.Address, _bridgeAddress common.Address, initializePackedParameters ProofOfEfficiencyInitializePackedParameters, genesisRoot [32]byte, _trustedSequencerURL string, _networkName string) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.Initialize(&_Proofofefficiency.TransactOpts, _globalExitRootManager, _matic, _rollupVerifier, _bridgeAddress, initializePackedParameters, genesisRoot, _trustedSequencerURL, _networkName)
}

// Initialize is a paid mutator transaction binding the contract method 0x60943d6a.
//
// Solidity: function initialize(address _globalExitRootManager, address _matic, address _rollupVerifier, address _bridgeAddress, (address,uint64,address,uint64,bool,address,uint64) initializePackedParameters, bytes32 genesisRoot, string _trustedSequencerURL, string _networkName) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) Initialize(_globalExitRootManager common.Address, _matic common.Address, _rollupVerifier common.Address, _bridgeAddress common.Address, initializePackedParameters ProofOfEfficiencyInitializePackedParameters, genesisRoot [32]byte, _trustedSequencerURL string, _networkName string) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.Initialize(&_Proofofefficiency.TransactOpts, _globalExitRootManager, _matic, _rollupVerifier, _bridgeAddress, initializePackedParameters, genesisRoot, _trustedSequencerURL, _networkName)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0xe11f3f18.
//
// Solidity: function overridePendingState(uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) OverridePendingState(opts *bind.TransactOpts, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "overridePendingState", initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proofA, proofB, proofC)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0xe11f3f18.
//
// Solidity: function overridePendingState(uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencySession) OverridePendingState(initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.OverridePendingState(&_Proofofefficiency.TransactOpts, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proofA, proofB, proofC)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0xe11f3f18.
//
// Solidity: function overridePendingState(uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) OverridePendingState(initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.OverridePendingState(&_Proofofefficiency.TransactOpts, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proofA, proofB, proofC)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x75c508b3.
//
// Solidity: function proveNonDeterministicPendingState(uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) ProveNonDeterministicPendingState(opts *bind.TransactOpts, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "proveNonDeterministicPendingState", initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proofA, proofB, proofC)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x75c508b3.
//
// Solidity: function proveNonDeterministicPendingState(uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencySession) ProveNonDeterministicPendingState(initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.ProveNonDeterministicPendingState(&_Proofofefficiency.TransactOpts, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proofA, proofB, proofC)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x75c508b3.
//
// Solidity: function proveNonDeterministicPendingState(uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) ProveNonDeterministicPendingState(initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.ProveNonDeterministicPendingState(&_Proofofefficiency.TransactOpts, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proofA, proofB, proofC)
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

// SequenceBatches is a paid mutator transaction binding the contract method 0x3c158267.
//
// Solidity: function sequenceBatches((bytes,bytes32,uint64,uint64)[] batches) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) SequenceBatches(opts *bind.TransactOpts, batches []ProofOfEfficiencyBatchData) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "sequenceBatches", batches)
}

// SequenceBatches is a paid mutator transaction binding the contract method 0x3c158267.
//
// Solidity: function sequenceBatches((bytes,bytes32,uint64,uint64)[] batches) returns()
func (_Proofofefficiency *ProofofefficiencySession) SequenceBatches(batches []ProofOfEfficiencyBatchData) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SequenceBatches(&_Proofofefficiency.TransactOpts, batches)
}

// SequenceBatches is a paid mutator transaction binding the contract method 0x3c158267.
//
// Solidity: function sequenceBatches((bytes,bytes32,uint64,uint64)[] batches) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) SequenceBatches(batches []ProofOfEfficiencyBatchData) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SequenceBatches(&_Proofofefficiency.TransactOpts, batches)
}

// SequenceForceBatches is a paid mutator transaction binding the contract method 0xd8d1091b.
//
// Solidity: function sequenceForceBatches((bytes,bytes32,uint64)[] batches) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) SequenceForceBatches(opts *bind.TransactOpts, batches []ProofOfEfficiencyForcedBatchData) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "sequenceForceBatches", batches)
}

// SequenceForceBatches is a paid mutator transaction binding the contract method 0xd8d1091b.
//
// Solidity: function sequenceForceBatches((bytes,bytes32,uint64)[] batches) returns()
func (_Proofofefficiency *ProofofefficiencySession) SequenceForceBatches(batches []ProofOfEfficiencyForcedBatchData) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SequenceForceBatches(&_Proofofefficiency.TransactOpts, batches)
}

// SequenceForceBatches is a paid mutator transaction binding the contract method 0xd8d1091b.
//
// Solidity: function sequenceForceBatches((bytes,bytes32,uint64)[] batches) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) SequenceForceBatches(batches []ProofOfEfficiencyForcedBatchData) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SequenceForceBatches(&_Proofofefficiency.TransactOpts, batches)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address newAdmin) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) SetAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "setAdmin", newAdmin)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address newAdmin) returns()
func (_Proofofefficiency *ProofofefficiencySession) SetAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetAdmin(&_Proofofefficiency.TransactOpts, newAdmin)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x704b6c02.
//
// Solidity: function setAdmin(address newAdmin) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) SetAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetAdmin(&_Proofofefficiency.TransactOpts, newAdmin)
}

// SetForceBatchAllowed is a paid mutator transaction binding the contract method 0x8c4a0af7.
//
// Solidity: function setForceBatchAllowed(bool newForceBatchAllowed) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) SetForceBatchAllowed(opts *bind.TransactOpts, newForceBatchAllowed bool) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "setForceBatchAllowed", newForceBatchAllowed)
}

// SetForceBatchAllowed is a paid mutator transaction binding the contract method 0x8c4a0af7.
//
// Solidity: function setForceBatchAllowed(bool newForceBatchAllowed) returns()
func (_Proofofefficiency *ProofofefficiencySession) SetForceBatchAllowed(newForceBatchAllowed bool) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetForceBatchAllowed(&_Proofofefficiency.TransactOpts, newForceBatchAllowed)
}

// SetForceBatchAllowed is a paid mutator transaction binding the contract method 0x8c4a0af7.
//
// Solidity: function setForceBatchAllowed(bool newForceBatchAllowed) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) SetForceBatchAllowed(newForceBatchAllowed bool) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetForceBatchAllowed(&_Proofofefficiency.TransactOpts, newForceBatchAllowed)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) SetPendingStateTimeout(opts *bind.TransactOpts, newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "setPendingStateTimeout", newPendingStateTimeout)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Proofofefficiency *ProofofefficiencySession) SetPendingStateTimeout(newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetPendingStateTimeout(&_Proofofefficiency.TransactOpts, newPendingStateTimeout)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) SetPendingStateTimeout(newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetPendingStateTimeout(&_Proofofefficiency.TransactOpts, newPendingStateTimeout)
}

// SetTrustedAggregator is a paid mutator transaction binding the contract method 0xf14916d6.
//
// Solidity: function setTrustedAggregator(address newTrustedAggregator) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) SetTrustedAggregator(opts *bind.TransactOpts, newTrustedAggregator common.Address) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "setTrustedAggregator", newTrustedAggregator)
}

// SetTrustedAggregator is a paid mutator transaction binding the contract method 0xf14916d6.
//
// Solidity: function setTrustedAggregator(address newTrustedAggregator) returns()
func (_Proofofefficiency *ProofofefficiencySession) SetTrustedAggregator(newTrustedAggregator common.Address) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetTrustedAggregator(&_Proofofefficiency.TransactOpts, newTrustedAggregator)
}

// SetTrustedAggregator is a paid mutator transaction binding the contract method 0xf14916d6.
//
// Solidity: function setTrustedAggregator(address newTrustedAggregator) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) SetTrustedAggregator(newTrustedAggregator common.Address) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetTrustedAggregator(&_Proofofefficiency.TransactOpts, newTrustedAggregator)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) SetTrustedAggregatorTimeout(opts *bind.TransactOpts, newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "setTrustedAggregatorTimeout", newTrustedAggregatorTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Proofofefficiency *ProofofefficiencySession) SetTrustedAggregatorTimeout(newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetTrustedAggregatorTimeout(&_Proofofefficiency.TransactOpts, newTrustedAggregatorTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) SetTrustedAggregatorTimeout(newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetTrustedAggregatorTimeout(&_Proofofefficiency.TransactOpts, newTrustedAggregatorTimeout)
}

// SetTrustedSequencer is a paid mutator transaction binding the contract method 0x6ff512cc.
//
// Solidity: function setTrustedSequencer(address newTrustedSequencer) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) SetTrustedSequencer(opts *bind.TransactOpts, newTrustedSequencer common.Address) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "setTrustedSequencer", newTrustedSequencer)
}

// SetTrustedSequencer is a paid mutator transaction binding the contract method 0x6ff512cc.
//
// Solidity: function setTrustedSequencer(address newTrustedSequencer) returns()
func (_Proofofefficiency *ProofofefficiencySession) SetTrustedSequencer(newTrustedSequencer common.Address) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetTrustedSequencer(&_Proofofefficiency.TransactOpts, newTrustedSequencer)
}

// SetTrustedSequencer is a paid mutator transaction binding the contract method 0x6ff512cc.
//
// Solidity: function setTrustedSequencer(address newTrustedSequencer) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) SetTrustedSequencer(newTrustedSequencer common.Address) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetTrustedSequencer(&_Proofofefficiency.TransactOpts, newTrustedSequencer)
}

// SetTrustedSequencerURL is a paid mutator transaction binding the contract method 0xc89e42df.
//
// Solidity: function setTrustedSequencerURL(string newTrustedSequencerURL) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) SetTrustedSequencerURL(opts *bind.TransactOpts, newTrustedSequencerURL string) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "setTrustedSequencerURL", newTrustedSequencerURL)
}

// SetTrustedSequencerURL is a paid mutator transaction binding the contract method 0xc89e42df.
//
// Solidity: function setTrustedSequencerURL(string newTrustedSequencerURL) returns()
func (_Proofofefficiency *ProofofefficiencySession) SetTrustedSequencerURL(newTrustedSequencerURL string) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetTrustedSequencerURL(&_Proofofefficiency.TransactOpts, newTrustedSequencerURL)
}

// SetTrustedSequencerURL is a paid mutator transaction binding the contract method 0xc89e42df.
//
// Solidity: function setTrustedSequencerURL(string newTrustedSequencerURL) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) SetTrustedSequencerURL(newTrustedSequencerURL string) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.SetTrustedSequencerURL(&_Proofofefficiency.TransactOpts, newTrustedSequencerURL)
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

// TrustedVerifyBatches is a paid mutator transaction binding the contract method 0xedc41121.
//
// Solidity: function trustedVerifyBatches(uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) TrustedVerifyBatches(opts *bind.TransactOpts, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "trustedVerifyBatches", pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proofA, proofB, proofC)
}

// TrustedVerifyBatches is a paid mutator transaction binding the contract method 0xedc41121.
//
// Solidity: function trustedVerifyBatches(uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencySession) TrustedVerifyBatches(pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.TrustedVerifyBatches(&_Proofofefficiency.TransactOpts, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proofA, proofB, proofC)
}

// TrustedVerifyBatches is a paid mutator transaction binding the contract method 0xedc41121.
//
// Solidity: function trustedVerifyBatches(uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) TrustedVerifyBatches(pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.TrustedVerifyBatches(&_Proofofefficiency.TransactOpts, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proofA, proofB, proofC)
}

// VerifyBatches is a paid mutator transaction binding the contract method 0x4834a343.
//
// Solidity: function verifyBatches(uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactor) VerifyBatches(opts *bind.TransactOpts, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.contract.Transact(opts, "verifyBatches", pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proofA, proofB, proofC)
}

// VerifyBatches is a paid mutator transaction binding the contract method 0x4834a343.
//
// Solidity: function verifyBatches(uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencySession) VerifyBatches(pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.VerifyBatches(&_Proofofefficiency.TransactOpts, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proofA, proofB, proofC)
}

// VerifyBatches is a paid mutator transaction binding the contract method 0x4834a343.
//
// Solidity: function verifyBatches(uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, uint256[2] proofA, uint256[2][2] proofB, uint256[2] proofC) returns()
func (_Proofofefficiency *ProofofefficiencyTransactorSession) VerifyBatches(pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proofA [2]*big.Int, proofB [2][2]*big.Int, proofC [2]*big.Int) (*types.Transaction, error) {
	return _Proofofefficiency.Contract.VerifyBatches(&_Proofofefficiency.TransactOpts, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proofA, proofB, proofC)
}

// ProofofefficiencyConsolidatePendingStateIterator is returned from FilterConsolidatePendingState and is used to iterate over the raw logs and unpacked data for ConsolidatePendingState events raised by the Proofofefficiency contract.
type ProofofefficiencyConsolidatePendingStateIterator struct {
	Event *ProofofefficiencyConsolidatePendingState // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencyConsolidatePendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencyConsolidatePendingState)
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
		it.Event = new(ProofofefficiencyConsolidatePendingState)
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
func (it *ProofofefficiencyConsolidatePendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencyConsolidatePendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencyConsolidatePendingState represents a ConsolidatePendingState event raised by the Proofofefficiency contract.
type ProofofefficiencyConsolidatePendingState struct {
	NumBatch        uint64
	StateRoot       [32]byte
	PendingStateNum uint64
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConsolidatePendingState is a free log retrieval operation binding the contract event 0x328d3c6c0fd6f1be0515e422f2d87e59f25922cbc2233568515a0c4bc3f8510e.
//
// Solidity: event ConsolidatePendingState(uint64 indexed numBatch, bytes32 stateRoot, uint64 indexed pendingStateNum)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterConsolidatePendingState(opts *bind.FilterOpts, numBatch []uint64, pendingStateNum []uint64) (*ProofofefficiencyConsolidatePendingStateIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var pendingStateNumRule []interface{}
	for _, pendingStateNumItem := range pendingStateNum {
		pendingStateNumRule = append(pendingStateNumRule, pendingStateNumItem)
	}

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "ConsolidatePendingState", numBatchRule, pendingStateNumRule)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyConsolidatePendingStateIterator{contract: _Proofofefficiency.contract, event: "ConsolidatePendingState", logs: logs, sub: sub}, nil
}

// WatchConsolidatePendingState is a free log subscription operation binding the contract event 0x328d3c6c0fd6f1be0515e422f2d87e59f25922cbc2233568515a0c4bc3f8510e.
//
// Solidity: event ConsolidatePendingState(uint64 indexed numBatch, bytes32 stateRoot, uint64 indexed pendingStateNum)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchConsolidatePendingState(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyConsolidatePendingState, numBatch []uint64, pendingStateNum []uint64) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var pendingStateNumRule []interface{}
	for _, pendingStateNumItem := range pendingStateNum {
		pendingStateNumRule = append(pendingStateNumRule, pendingStateNumItem)
	}

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "ConsolidatePendingState", numBatchRule, pendingStateNumRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencyConsolidatePendingState)
				if err := _Proofofefficiency.contract.UnpackLog(event, "ConsolidatePendingState", log); err != nil {
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

// ParseConsolidatePendingState is a log parse operation binding the contract event 0x328d3c6c0fd6f1be0515e422f2d87e59f25922cbc2233568515a0c4bc3f8510e.
//
// Solidity: event ConsolidatePendingState(uint64 indexed numBatch, bytes32 stateRoot, uint64 indexed pendingStateNum)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseConsolidatePendingState(log types.Log) (*ProofofefficiencyConsolidatePendingState, error) {
	event := new(ProofofefficiencyConsolidatePendingState)
	if err := _Proofofefficiency.contract.UnpackLog(event, "ConsolidatePendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencyEmergencyStateActivatedIterator is returned from FilterEmergencyStateActivated and is used to iterate over the raw logs and unpacked data for EmergencyStateActivated events raised by the Proofofefficiency contract.
type ProofofefficiencyEmergencyStateActivatedIterator struct {
	Event *ProofofefficiencyEmergencyStateActivated // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencyEmergencyStateActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencyEmergencyStateActivated)
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
		it.Event = new(ProofofefficiencyEmergencyStateActivated)
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
func (it *ProofofefficiencyEmergencyStateActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencyEmergencyStateActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencyEmergencyStateActivated represents a EmergencyStateActivated event raised by the Proofofefficiency contract.
type ProofofefficiencyEmergencyStateActivated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEmergencyStateActivated is a free log retrieval operation binding the contract event 0x2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497.
//
// Solidity: event EmergencyStateActivated()
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterEmergencyStateActivated(opts *bind.FilterOpts) (*ProofofefficiencyEmergencyStateActivatedIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "EmergencyStateActivated")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyEmergencyStateActivatedIterator{contract: _Proofofefficiency.contract, event: "EmergencyStateActivated", logs: logs, sub: sub}, nil
}

// WatchEmergencyStateActivated is a free log subscription operation binding the contract event 0x2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497.
//
// Solidity: event EmergencyStateActivated()
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchEmergencyStateActivated(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyEmergencyStateActivated) (event.Subscription, error) {

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "EmergencyStateActivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencyEmergencyStateActivated)
				if err := _Proofofefficiency.contract.UnpackLog(event, "EmergencyStateActivated", log); err != nil {
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

// ParseEmergencyStateActivated is a log parse operation binding the contract event 0x2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497.
//
// Solidity: event EmergencyStateActivated()
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseEmergencyStateActivated(log types.Log) (*ProofofefficiencyEmergencyStateActivated, error) {
	event := new(ProofofefficiencyEmergencyStateActivated)
	if err := _Proofofefficiency.contract.UnpackLog(event, "EmergencyStateActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencyEmergencyStateDeactivatedIterator is returned from FilterEmergencyStateDeactivated and is used to iterate over the raw logs and unpacked data for EmergencyStateDeactivated events raised by the Proofofefficiency contract.
type ProofofefficiencyEmergencyStateDeactivatedIterator struct {
	Event *ProofofefficiencyEmergencyStateDeactivated // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencyEmergencyStateDeactivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencyEmergencyStateDeactivated)
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
		it.Event = new(ProofofefficiencyEmergencyStateDeactivated)
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
func (it *ProofofefficiencyEmergencyStateDeactivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencyEmergencyStateDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencyEmergencyStateDeactivated represents a EmergencyStateDeactivated event raised by the Proofofefficiency contract.
type ProofofefficiencyEmergencyStateDeactivated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEmergencyStateDeactivated is a free log retrieval operation binding the contract event 0x1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3.
//
// Solidity: event EmergencyStateDeactivated()
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterEmergencyStateDeactivated(opts *bind.FilterOpts) (*ProofofefficiencyEmergencyStateDeactivatedIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "EmergencyStateDeactivated")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyEmergencyStateDeactivatedIterator{contract: _Proofofefficiency.contract, event: "EmergencyStateDeactivated", logs: logs, sub: sub}, nil
}

// WatchEmergencyStateDeactivated is a free log subscription operation binding the contract event 0x1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3.
//
// Solidity: event EmergencyStateDeactivated()
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchEmergencyStateDeactivated(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyEmergencyStateDeactivated) (event.Subscription, error) {

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "EmergencyStateDeactivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencyEmergencyStateDeactivated)
				if err := _Proofofefficiency.contract.UnpackLog(event, "EmergencyStateDeactivated", log); err != nil {
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

// ParseEmergencyStateDeactivated is a log parse operation binding the contract event 0x1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3.
//
// Solidity: event EmergencyStateDeactivated()
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseEmergencyStateDeactivated(log types.Log) (*ProofofefficiencyEmergencyStateDeactivated, error) {
	event := new(ProofofefficiencyEmergencyStateDeactivated)
	if err := _Proofofefficiency.contract.UnpackLog(event, "EmergencyStateDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencyForceBatchIterator is returned from FilterForceBatch and is used to iterate over the raw logs and unpacked data for ForceBatch events raised by the Proofofefficiency contract.
type ProofofefficiencyForceBatchIterator struct {
	Event *ProofofefficiencyForceBatch // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencyForceBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencyForceBatch)
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
		it.Event = new(ProofofefficiencyForceBatch)
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
func (it *ProofofefficiencyForceBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencyForceBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencyForceBatch represents a ForceBatch event raised by the Proofofefficiency contract.
type ProofofefficiencyForceBatch struct {
	ForceBatchNum      uint64
	LastGlobalExitRoot [32]byte
	Sequencer          common.Address
	Transactions       []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterForceBatch is a free log retrieval operation binding the contract event 0xf94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc931.
//
// Solidity: event ForceBatch(uint64 indexed forceBatchNum, bytes32 lastGlobalExitRoot, address sequencer, bytes transactions)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterForceBatch(opts *bind.FilterOpts, forceBatchNum []uint64) (*ProofofefficiencyForceBatchIterator, error) {

	var forceBatchNumRule []interface{}
	for _, forceBatchNumItem := range forceBatchNum {
		forceBatchNumRule = append(forceBatchNumRule, forceBatchNumItem)
	}

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "ForceBatch", forceBatchNumRule)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyForceBatchIterator{contract: _Proofofefficiency.contract, event: "ForceBatch", logs: logs, sub: sub}, nil
}

// WatchForceBatch is a free log subscription operation binding the contract event 0xf94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc931.
//
// Solidity: event ForceBatch(uint64 indexed forceBatchNum, bytes32 lastGlobalExitRoot, address sequencer, bytes transactions)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchForceBatch(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyForceBatch, forceBatchNum []uint64) (event.Subscription, error) {

	var forceBatchNumRule []interface{}
	for _, forceBatchNumItem := range forceBatchNum {
		forceBatchNumRule = append(forceBatchNumRule, forceBatchNumItem)
	}

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "ForceBatch", forceBatchNumRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencyForceBatch)
				if err := _Proofofefficiency.contract.UnpackLog(event, "ForceBatch", log); err != nil {
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
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseForceBatch(log types.Log) (*ProofofefficiencyForceBatch, error) {
	event := new(ProofofefficiencyForceBatch)
	if err := _Proofofefficiency.contract.UnpackLog(event, "ForceBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencyInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Proofofefficiency contract.
type ProofofefficiencyInitializedIterator struct {
	Event *ProofofefficiencyInitialized // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencyInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencyInitialized)
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
		it.Event = new(ProofofefficiencyInitialized)
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
func (it *ProofofefficiencyInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencyInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencyInitialized represents a Initialized event raised by the Proofofefficiency contract.
type ProofofefficiencyInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterInitialized(opts *bind.FilterOpts) (*ProofofefficiencyInitializedIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyInitializedIterator{contract: _Proofofefficiency.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyInitialized) (event.Subscription, error) {

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencyInitialized)
				if err := _Proofofefficiency.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseInitialized(log types.Log) (*ProofofefficiencyInitialized, error) {
	event := new(ProofofefficiencyInitialized)
	if err := _Proofofefficiency.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencyOverridePendingStateIterator is returned from FilterOverridePendingState and is used to iterate over the raw logs and unpacked data for OverridePendingState events raised by the Proofofefficiency contract.
type ProofofefficiencyOverridePendingStateIterator struct {
	Event *ProofofefficiencyOverridePendingState // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencyOverridePendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencyOverridePendingState)
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
		it.Event = new(ProofofefficiencyOverridePendingState)
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
func (it *ProofofefficiencyOverridePendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencyOverridePendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencyOverridePendingState represents a OverridePendingState event raised by the Proofofefficiency contract.
type ProofofefficiencyOverridePendingState struct {
	NumBatch   uint64
	StateRoot  [32]byte
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterOverridePendingState is a free log retrieval operation binding the contract event 0xcc1b5520188bf1dd3e63f98164b577c4d75c11a619ddea692112f0d1aec4cf72.
//
// Solidity: event OverridePendingState(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterOverridePendingState(opts *bind.FilterOpts, numBatch []uint64, aggregator []common.Address) (*ProofofefficiencyOverridePendingStateIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "OverridePendingState", numBatchRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyOverridePendingStateIterator{contract: _Proofofefficiency.contract, event: "OverridePendingState", logs: logs, sub: sub}, nil
}

// WatchOverridePendingState is a free log subscription operation binding the contract event 0xcc1b5520188bf1dd3e63f98164b577c4d75c11a619ddea692112f0d1aec4cf72.
//
// Solidity: event OverridePendingState(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchOverridePendingState(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyOverridePendingState, numBatch []uint64, aggregator []common.Address) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "OverridePendingState", numBatchRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencyOverridePendingState)
				if err := _Proofofefficiency.contract.UnpackLog(event, "OverridePendingState", log); err != nil {
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

// ParseOverridePendingState is a log parse operation binding the contract event 0xcc1b5520188bf1dd3e63f98164b577c4d75c11a619ddea692112f0d1aec4cf72.
//
// Solidity: event OverridePendingState(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseOverridePendingState(log types.Log) (*ProofofefficiencyOverridePendingState, error) {
	event := new(ProofofefficiencyOverridePendingState)
	if err := _Proofofefficiency.contract.UnpackLog(event, "OverridePendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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

// ProofofefficiencyProveNonDeterministicPendingStateIterator is returned from FilterProveNonDeterministicPendingState and is used to iterate over the raw logs and unpacked data for ProveNonDeterministicPendingState events raised by the Proofofefficiency contract.
type ProofofefficiencyProveNonDeterministicPendingStateIterator struct {
	Event *ProofofefficiencyProveNonDeterministicPendingState // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencyProveNonDeterministicPendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencyProveNonDeterministicPendingState)
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
		it.Event = new(ProofofefficiencyProveNonDeterministicPendingState)
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
func (it *ProofofefficiencyProveNonDeterministicPendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencyProveNonDeterministicPendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencyProveNonDeterministicPendingState represents a ProveNonDeterministicPendingState event raised by the Proofofefficiency contract.
type ProofofefficiencyProveNonDeterministicPendingState struct {
	StoredStateRoot [32]byte
	ProvedStateRoot [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterProveNonDeterministicPendingState is a free log retrieval operation binding the contract event 0x1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010.
//
// Solidity: event ProveNonDeterministicPendingState(bytes32 storedStateRoot, bytes32 provedStateRoot)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterProveNonDeterministicPendingState(opts *bind.FilterOpts) (*ProofofefficiencyProveNonDeterministicPendingStateIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "ProveNonDeterministicPendingState")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyProveNonDeterministicPendingStateIterator{contract: _Proofofefficiency.contract, event: "ProveNonDeterministicPendingState", logs: logs, sub: sub}, nil
}

// WatchProveNonDeterministicPendingState is a free log subscription operation binding the contract event 0x1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010.
//
// Solidity: event ProveNonDeterministicPendingState(bytes32 storedStateRoot, bytes32 provedStateRoot)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchProveNonDeterministicPendingState(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyProveNonDeterministicPendingState) (event.Subscription, error) {

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "ProveNonDeterministicPendingState")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencyProveNonDeterministicPendingState)
				if err := _Proofofefficiency.contract.UnpackLog(event, "ProveNonDeterministicPendingState", log); err != nil {
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

// ParseProveNonDeterministicPendingState is a log parse operation binding the contract event 0x1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010.
//
// Solidity: event ProveNonDeterministicPendingState(bytes32 storedStateRoot, bytes32 provedStateRoot)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseProveNonDeterministicPendingState(log types.Log) (*ProofofefficiencyProveNonDeterministicPendingState, error) {
	event := new(ProofofefficiencyProveNonDeterministicPendingState)
	if err := _Proofofefficiency.contract.UnpackLog(event, "ProveNonDeterministicPendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencySequenceBatchesIterator is returned from FilterSequenceBatches and is used to iterate over the raw logs and unpacked data for SequenceBatches events raised by the Proofofefficiency contract.
type ProofofefficiencySequenceBatchesIterator struct {
	Event *ProofofefficiencySequenceBatches // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencySequenceBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencySequenceBatches)
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
		it.Event = new(ProofofefficiencySequenceBatches)
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
func (it *ProofofefficiencySequenceBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencySequenceBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencySequenceBatches represents a SequenceBatches event raised by the Proofofefficiency contract.
type ProofofefficiencySequenceBatches struct {
	NumBatch uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSequenceBatches is a free log retrieval operation binding the contract event 0x303446e6a8cb73c83dff421c0b1d5e5ce0719dab1bff13660fc254e58cc17fce.
//
// Solidity: event SequenceBatches(uint64 indexed numBatch)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSequenceBatches(opts *bind.FilterOpts, numBatch []uint64) (*ProofofefficiencySequenceBatchesIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "SequenceBatches", numBatchRule)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencySequenceBatchesIterator{contract: _Proofofefficiency.contract, event: "SequenceBatches", logs: logs, sub: sub}, nil
}

// WatchSequenceBatches is a free log subscription operation binding the contract event 0x303446e6a8cb73c83dff421c0b1d5e5ce0719dab1bff13660fc254e58cc17fce.
//
// Solidity: event SequenceBatches(uint64 indexed numBatch)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSequenceBatches(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySequenceBatches, numBatch []uint64) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "SequenceBatches", numBatchRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencySequenceBatches)
				if err := _Proofofefficiency.contract.UnpackLog(event, "SequenceBatches", log); err != nil {
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

// ParseSequenceBatches is a log parse operation binding the contract event 0x303446e6a8cb73c83dff421c0b1d5e5ce0719dab1bff13660fc254e58cc17fce.
//
// Solidity: event SequenceBatches(uint64 indexed numBatch)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseSequenceBatches(log types.Log) (*ProofofefficiencySequenceBatches, error) {
	event := new(ProofofefficiencySequenceBatches)
	if err := _Proofofefficiency.contract.UnpackLog(event, "SequenceBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencySequenceForceBatchesIterator is returned from FilterSequenceForceBatches and is used to iterate over the raw logs and unpacked data for SequenceForceBatches events raised by the Proofofefficiency contract.
type ProofofefficiencySequenceForceBatchesIterator struct {
	Event *ProofofefficiencySequenceForceBatches // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencySequenceForceBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencySequenceForceBatches)
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
		it.Event = new(ProofofefficiencySequenceForceBatches)
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
func (it *ProofofefficiencySequenceForceBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencySequenceForceBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencySequenceForceBatches represents a SequenceForceBatches event raised by the Proofofefficiency contract.
type ProofofefficiencySequenceForceBatches struct {
	NumBatch uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSequenceForceBatches is a free log retrieval operation binding the contract event 0x648a61dd2438f072f5a1960939abd30f37aea80d2e94c9792ad142d3e0a490a4.
//
// Solidity: event SequenceForceBatches(uint64 indexed numBatch)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSequenceForceBatches(opts *bind.FilterOpts, numBatch []uint64) (*ProofofefficiencySequenceForceBatchesIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "SequenceForceBatches", numBatchRule)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencySequenceForceBatchesIterator{contract: _Proofofefficiency.contract, event: "SequenceForceBatches", logs: logs, sub: sub}, nil
}

// WatchSequenceForceBatches is a free log subscription operation binding the contract event 0x648a61dd2438f072f5a1960939abd30f37aea80d2e94c9792ad142d3e0a490a4.
//
// Solidity: event SequenceForceBatches(uint64 indexed numBatch)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSequenceForceBatches(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySequenceForceBatches, numBatch []uint64) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "SequenceForceBatches", numBatchRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencySequenceForceBatches)
				if err := _Proofofefficiency.contract.UnpackLog(event, "SequenceForceBatches", log); err != nil {
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
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseSequenceForceBatches(log types.Log) (*ProofofefficiencySequenceForceBatches, error) {
	event := new(ProofofefficiencySequenceForceBatches)
	if err := _Proofofefficiency.contract.UnpackLog(event, "SequenceForceBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencySetAdminIterator is returned from FilterSetAdmin and is used to iterate over the raw logs and unpacked data for SetAdmin events raised by the Proofofefficiency contract.
type ProofofefficiencySetAdminIterator struct {
	Event *ProofofefficiencySetAdmin // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencySetAdminIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencySetAdmin)
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
		it.Event = new(ProofofefficiencySetAdmin)
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
func (it *ProofofefficiencySetAdminIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencySetAdminIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencySetAdmin represents a SetAdmin event raised by the Proofofefficiency contract.
type ProofofefficiencySetAdmin struct {
	NewAdmin common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSetAdmin is a free log retrieval operation binding the contract event 0x5a272403b402d892977df56625f4164ccaf70ca3863991c43ecfe76a6905b0a1.
//
// Solidity: event SetAdmin(address newAdmin)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSetAdmin(opts *bind.FilterOpts) (*ProofofefficiencySetAdminIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "SetAdmin")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencySetAdminIterator{contract: _Proofofefficiency.contract, event: "SetAdmin", logs: logs, sub: sub}, nil
}

// WatchSetAdmin is a free log subscription operation binding the contract event 0x5a272403b402d892977df56625f4164ccaf70ca3863991c43ecfe76a6905b0a1.
//
// Solidity: event SetAdmin(address newAdmin)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSetAdmin(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySetAdmin) (event.Subscription, error) {

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "SetAdmin")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencySetAdmin)
				if err := _Proofofefficiency.contract.UnpackLog(event, "SetAdmin", log); err != nil {
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

// ParseSetAdmin is a log parse operation binding the contract event 0x5a272403b402d892977df56625f4164ccaf70ca3863991c43ecfe76a6905b0a1.
//
// Solidity: event SetAdmin(address newAdmin)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseSetAdmin(log types.Log) (*ProofofefficiencySetAdmin, error) {
	event := new(ProofofefficiencySetAdmin)
	if err := _Proofofefficiency.contract.UnpackLog(event, "SetAdmin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencySetForceBatchAllowedIterator is returned from FilterSetForceBatchAllowed and is used to iterate over the raw logs and unpacked data for SetForceBatchAllowed events raised by the Proofofefficiency contract.
type ProofofefficiencySetForceBatchAllowedIterator struct {
	Event *ProofofefficiencySetForceBatchAllowed // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencySetForceBatchAllowedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencySetForceBatchAllowed)
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
		it.Event = new(ProofofefficiencySetForceBatchAllowed)
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
func (it *ProofofefficiencySetForceBatchAllowedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencySetForceBatchAllowedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencySetForceBatchAllowed represents a SetForceBatchAllowed event raised by the Proofofefficiency contract.
type ProofofefficiencySetForceBatchAllowed struct {
	NewForceBatchAllowed bool
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSetForceBatchAllowed is a free log retrieval operation binding the contract event 0xbacda50a4a8575be1d91a7ebe29ee45056f3a94f12a2281eb6b43afa33bcefe6.
//
// Solidity: event SetForceBatchAllowed(bool newForceBatchAllowed)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSetForceBatchAllowed(opts *bind.FilterOpts) (*ProofofefficiencySetForceBatchAllowedIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "SetForceBatchAllowed")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencySetForceBatchAllowedIterator{contract: _Proofofefficiency.contract, event: "SetForceBatchAllowed", logs: logs, sub: sub}, nil
}

// WatchSetForceBatchAllowed is a free log subscription operation binding the contract event 0xbacda50a4a8575be1d91a7ebe29ee45056f3a94f12a2281eb6b43afa33bcefe6.
//
// Solidity: event SetForceBatchAllowed(bool newForceBatchAllowed)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSetForceBatchAllowed(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySetForceBatchAllowed) (event.Subscription, error) {

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "SetForceBatchAllowed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencySetForceBatchAllowed)
				if err := _Proofofefficiency.contract.UnpackLog(event, "SetForceBatchAllowed", log); err != nil {
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

// ParseSetForceBatchAllowed is a log parse operation binding the contract event 0xbacda50a4a8575be1d91a7ebe29ee45056f3a94f12a2281eb6b43afa33bcefe6.
//
// Solidity: event SetForceBatchAllowed(bool newForceBatchAllowed)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseSetForceBatchAllowed(log types.Log) (*ProofofefficiencySetForceBatchAllowed, error) {
	event := new(ProofofefficiencySetForceBatchAllowed)
	if err := _Proofofefficiency.contract.UnpackLog(event, "SetForceBatchAllowed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencySetPendingStateTimeoutIterator is returned from FilterSetPendingStateTimeout and is used to iterate over the raw logs and unpacked data for SetPendingStateTimeout events raised by the Proofofefficiency contract.
type ProofofefficiencySetPendingStateTimeoutIterator struct {
	Event *ProofofefficiencySetPendingStateTimeout // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencySetPendingStateTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencySetPendingStateTimeout)
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
		it.Event = new(ProofofefficiencySetPendingStateTimeout)
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
func (it *ProofofefficiencySetPendingStateTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencySetPendingStateTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencySetPendingStateTimeout represents a SetPendingStateTimeout event raised by the Proofofefficiency contract.
type ProofofefficiencySetPendingStateTimeout struct {
	NewPendingStateTimeout uint64
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterSetPendingStateTimeout is a free log retrieval operation binding the contract event 0xc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c75.
//
// Solidity: event SetPendingStateTimeout(uint64 newPendingStateTimeout)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSetPendingStateTimeout(opts *bind.FilterOpts) (*ProofofefficiencySetPendingStateTimeoutIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "SetPendingStateTimeout")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencySetPendingStateTimeoutIterator{contract: _Proofofefficiency.contract, event: "SetPendingStateTimeout", logs: logs, sub: sub}, nil
}

// WatchSetPendingStateTimeout is a free log subscription operation binding the contract event 0xc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c75.
//
// Solidity: event SetPendingStateTimeout(uint64 newPendingStateTimeout)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSetPendingStateTimeout(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySetPendingStateTimeout) (event.Subscription, error) {

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "SetPendingStateTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencySetPendingStateTimeout)
				if err := _Proofofefficiency.contract.UnpackLog(event, "SetPendingStateTimeout", log); err != nil {
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

// ParseSetPendingStateTimeout is a log parse operation binding the contract event 0xc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c75.
//
// Solidity: event SetPendingStateTimeout(uint64 newPendingStateTimeout)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseSetPendingStateTimeout(log types.Log) (*ProofofefficiencySetPendingStateTimeout, error) {
	event := new(ProofofefficiencySetPendingStateTimeout)
	if err := _Proofofefficiency.contract.UnpackLog(event, "SetPendingStateTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencySetTrustedAggregatorIterator is returned from FilterSetTrustedAggregator and is used to iterate over the raw logs and unpacked data for SetTrustedAggregator events raised by the Proofofefficiency contract.
type ProofofefficiencySetTrustedAggregatorIterator struct {
	Event *ProofofefficiencySetTrustedAggregator // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencySetTrustedAggregatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencySetTrustedAggregator)
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
		it.Event = new(ProofofefficiencySetTrustedAggregator)
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
func (it *ProofofefficiencySetTrustedAggregatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencySetTrustedAggregatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencySetTrustedAggregator represents a SetTrustedAggregator event raised by the Proofofefficiency contract.
type ProofofefficiencySetTrustedAggregator struct {
	NewTrustedAggregator common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedAggregator is a free log retrieval operation binding the contract event 0x61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca.
//
// Solidity: event SetTrustedAggregator(address newTrustedAggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSetTrustedAggregator(opts *bind.FilterOpts) (*ProofofefficiencySetTrustedAggregatorIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "SetTrustedAggregator")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencySetTrustedAggregatorIterator{contract: _Proofofefficiency.contract, event: "SetTrustedAggregator", logs: logs, sub: sub}, nil
}

// WatchSetTrustedAggregator is a free log subscription operation binding the contract event 0x61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca.
//
// Solidity: event SetTrustedAggregator(address newTrustedAggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSetTrustedAggregator(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySetTrustedAggregator) (event.Subscription, error) {

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "SetTrustedAggregator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencySetTrustedAggregator)
				if err := _Proofofefficiency.contract.UnpackLog(event, "SetTrustedAggregator", log); err != nil {
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

// ParseSetTrustedAggregator is a log parse operation binding the contract event 0x61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca.
//
// Solidity: event SetTrustedAggregator(address newTrustedAggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseSetTrustedAggregator(log types.Log) (*ProofofefficiencySetTrustedAggregator, error) {
	event := new(ProofofefficiencySetTrustedAggregator)
	if err := _Proofofefficiency.contract.UnpackLog(event, "SetTrustedAggregator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencySetTrustedAggregatorTimeoutIterator is returned from FilterSetTrustedAggregatorTimeout and is used to iterate over the raw logs and unpacked data for SetTrustedAggregatorTimeout events raised by the Proofofefficiency contract.
type ProofofefficiencySetTrustedAggregatorTimeoutIterator struct {
	Event *ProofofefficiencySetTrustedAggregatorTimeout // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencySetTrustedAggregatorTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencySetTrustedAggregatorTimeout)
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
		it.Event = new(ProofofefficiencySetTrustedAggregatorTimeout)
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
func (it *ProofofefficiencySetTrustedAggregatorTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencySetTrustedAggregatorTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencySetTrustedAggregatorTimeout represents a SetTrustedAggregatorTimeout event raised by the Proofofefficiency contract.
type ProofofefficiencySetTrustedAggregatorTimeout struct {
	NewTrustedAggregatorTimeout uint64
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedAggregatorTimeout is a free log retrieval operation binding the contract event 0x1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1.
//
// Solidity: event SetTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSetTrustedAggregatorTimeout(opts *bind.FilterOpts) (*ProofofefficiencySetTrustedAggregatorTimeoutIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "SetTrustedAggregatorTimeout")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencySetTrustedAggregatorTimeoutIterator{contract: _Proofofefficiency.contract, event: "SetTrustedAggregatorTimeout", logs: logs, sub: sub}, nil
}

// WatchSetTrustedAggregatorTimeout is a free log subscription operation binding the contract event 0x1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1.
//
// Solidity: event SetTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSetTrustedAggregatorTimeout(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySetTrustedAggregatorTimeout) (event.Subscription, error) {

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "SetTrustedAggregatorTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencySetTrustedAggregatorTimeout)
				if err := _Proofofefficiency.contract.UnpackLog(event, "SetTrustedAggregatorTimeout", log); err != nil {
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

// ParseSetTrustedAggregatorTimeout is a log parse operation binding the contract event 0x1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1.
//
// Solidity: event SetTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseSetTrustedAggregatorTimeout(log types.Log) (*ProofofefficiencySetTrustedAggregatorTimeout, error) {
	event := new(ProofofefficiencySetTrustedAggregatorTimeout)
	if err := _Proofofefficiency.contract.UnpackLog(event, "SetTrustedAggregatorTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencySetTrustedSequencerIterator is returned from FilterSetTrustedSequencer and is used to iterate over the raw logs and unpacked data for SetTrustedSequencer events raised by the Proofofefficiency contract.
type ProofofefficiencySetTrustedSequencerIterator struct {
	Event *ProofofefficiencySetTrustedSequencer // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencySetTrustedSequencerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencySetTrustedSequencer)
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
		it.Event = new(ProofofefficiencySetTrustedSequencer)
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
func (it *ProofofefficiencySetTrustedSequencerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencySetTrustedSequencerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencySetTrustedSequencer represents a SetTrustedSequencer event raised by the Proofofefficiency contract.
type ProofofefficiencySetTrustedSequencer struct {
	NewTrustedSequencer common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedSequencer is a free log retrieval operation binding the contract event 0xf54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc0.
//
// Solidity: event SetTrustedSequencer(address newTrustedSequencer)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSetTrustedSequencer(opts *bind.FilterOpts) (*ProofofefficiencySetTrustedSequencerIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "SetTrustedSequencer")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencySetTrustedSequencerIterator{contract: _Proofofefficiency.contract, event: "SetTrustedSequencer", logs: logs, sub: sub}, nil
}

// WatchSetTrustedSequencer is a free log subscription operation binding the contract event 0xf54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc0.
//
// Solidity: event SetTrustedSequencer(address newTrustedSequencer)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSetTrustedSequencer(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySetTrustedSequencer) (event.Subscription, error) {

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "SetTrustedSequencer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencySetTrustedSequencer)
				if err := _Proofofefficiency.contract.UnpackLog(event, "SetTrustedSequencer", log); err != nil {
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
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseSetTrustedSequencer(log types.Log) (*ProofofefficiencySetTrustedSequencer, error) {
	event := new(ProofofefficiencySetTrustedSequencer)
	if err := _Proofofefficiency.contract.UnpackLog(event, "SetTrustedSequencer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencySetTrustedSequencerURLIterator is returned from FilterSetTrustedSequencerURL and is used to iterate over the raw logs and unpacked data for SetTrustedSequencerURL events raised by the Proofofefficiency contract.
type ProofofefficiencySetTrustedSequencerURLIterator struct {
	Event *ProofofefficiencySetTrustedSequencerURL // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencySetTrustedSequencerURLIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencySetTrustedSequencerURL)
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
		it.Event = new(ProofofefficiencySetTrustedSequencerURL)
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
func (it *ProofofefficiencySetTrustedSequencerURLIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencySetTrustedSequencerURLIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencySetTrustedSequencerURL represents a SetTrustedSequencerURL event raised by the Proofofefficiency contract.
type ProofofefficiencySetTrustedSequencerURL struct {
	NewTrustedSequencerURL string
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedSequencerURL is a free log retrieval operation binding the contract event 0x6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b20.
//
// Solidity: event SetTrustedSequencerURL(string newTrustedSequencerURL)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterSetTrustedSequencerURL(opts *bind.FilterOpts) (*ProofofefficiencySetTrustedSequencerURLIterator, error) {

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "SetTrustedSequencerURL")
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencySetTrustedSequencerURLIterator{contract: _Proofofefficiency.contract, event: "SetTrustedSequencerURL", logs: logs, sub: sub}, nil
}

// WatchSetTrustedSequencerURL is a free log subscription operation binding the contract event 0x6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b20.
//
// Solidity: event SetTrustedSequencerURL(string newTrustedSequencerURL)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchSetTrustedSequencerURL(opts *bind.WatchOpts, sink chan<- *ProofofefficiencySetTrustedSequencerURL) (event.Subscription, error) {

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "SetTrustedSequencerURL")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencySetTrustedSequencerURL)
				if err := _Proofofefficiency.contract.UnpackLog(event, "SetTrustedSequencerURL", log); err != nil {
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
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseSetTrustedSequencerURL(log types.Log) (*ProofofefficiencySetTrustedSequencerURL, error) {
	event := new(ProofofefficiencySetTrustedSequencerURL)
	if err := _Proofofefficiency.contract.UnpackLog(event, "SetTrustedSequencerURL", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencyTrustedVerifyBatchesIterator is returned from FilterTrustedVerifyBatches and is used to iterate over the raw logs and unpacked data for TrustedVerifyBatches events raised by the Proofofefficiency contract.
type ProofofefficiencyTrustedVerifyBatchesIterator struct {
	Event *ProofofefficiencyTrustedVerifyBatches // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencyTrustedVerifyBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencyTrustedVerifyBatches)
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
		it.Event = new(ProofofefficiencyTrustedVerifyBatches)
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
func (it *ProofofefficiencyTrustedVerifyBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencyTrustedVerifyBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencyTrustedVerifyBatches represents a TrustedVerifyBatches event raised by the Proofofefficiency contract.
type ProofofefficiencyTrustedVerifyBatches struct {
	NumBatch   uint64
	StateRoot  [32]byte
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTrustedVerifyBatches is a free log retrieval operation binding the contract event 0x0c0ce073a7d7b5850c04ccc4b20ee7d3179d5f57d0ac44399565792c0f72fce7.
//
// Solidity: event TrustedVerifyBatches(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterTrustedVerifyBatches(opts *bind.FilterOpts, numBatch []uint64, aggregator []common.Address) (*ProofofefficiencyTrustedVerifyBatchesIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "TrustedVerifyBatches", numBatchRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyTrustedVerifyBatchesIterator{contract: _Proofofefficiency.contract, event: "TrustedVerifyBatches", logs: logs, sub: sub}, nil
}

// WatchTrustedVerifyBatches is a free log subscription operation binding the contract event 0x0c0ce073a7d7b5850c04ccc4b20ee7d3179d5f57d0ac44399565792c0f72fce7.
//
// Solidity: event TrustedVerifyBatches(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchTrustedVerifyBatches(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyTrustedVerifyBatches, numBatch []uint64, aggregator []common.Address) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "TrustedVerifyBatches", numBatchRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencyTrustedVerifyBatches)
				if err := _Proofofefficiency.contract.UnpackLog(event, "TrustedVerifyBatches", log); err != nil {
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

// ParseTrustedVerifyBatches is a log parse operation binding the contract event 0x0c0ce073a7d7b5850c04ccc4b20ee7d3179d5f57d0ac44399565792c0f72fce7.
//
// Solidity: event TrustedVerifyBatches(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseTrustedVerifyBatches(log types.Log) (*ProofofefficiencyTrustedVerifyBatches, error) {
	event := new(ProofofefficiencyTrustedVerifyBatches)
	if err := _Proofofefficiency.contract.UnpackLog(event, "TrustedVerifyBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofofefficiencyVerifyBatchesIterator is returned from FilterVerifyBatches and is used to iterate over the raw logs and unpacked data for VerifyBatches events raised by the Proofofefficiency contract.
type ProofofefficiencyVerifyBatchesIterator struct {
	Event *ProofofefficiencyVerifyBatches // Event containing the contract specifics and raw log

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
func (it *ProofofefficiencyVerifyBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofofefficiencyVerifyBatches)
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
		it.Event = new(ProofofefficiencyVerifyBatches)
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
func (it *ProofofefficiencyVerifyBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofofefficiencyVerifyBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofofefficiencyVerifyBatches represents a VerifyBatches event raised by the Proofofefficiency contract.
type ProofofefficiencyVerifyBatches struct {
	NumBatch   uint64
	StateRoot  [32]byte
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVerifyBatches is a free log retrieval operation binding the contract event 0x9c72852172521097ba7e1482e6b44b351323df0155f97f4ea18fcec28e1f5966.
//
// Solidity: event VerifyBatches(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterVerifyBatches(opts *bind.FilterOpts, numBatch []uint64, aggregator []common.Address) (*ProofofefficiencyVerifyBatchesIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "VerifyBatches", numBatchRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyVerifyBatchesIterator{contract: _Proofofefficiency.contract, event: "VerifyBatches", logs: logs, sub: sub}, nil
}

// WatchVerifyBatches is a free log subscription operation binding the contract event 0x9c72852172521097ba7e1482e6b44b351323df0155f97f4ea18fcec28e1f5966.
//
// Solidity: event VerifyBatches(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchVerifyBatches(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyVerifyBatches, numBatch []uint64, aggregator []common.Address) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "VerifyBatches", numBatchRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofofefficiencyVerifyBatches)
				if err := _Proofofefficiency.contract.UnpackLog(event, "VerifyBatches", log); err != nil {
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
func (_Proofofefficiency *ProofofefficiencyFilterer) ParseVerifyBatches(log types.Log) (*ProofofefficiencyVerifyBatches, error) {
	event := new(ProofofefficiencyVerifyBatches)
	if err := _Proofofefficiency.contract.UnpackLog(event, "VerifyBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
