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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"ConsolidatePendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EmergencyStateActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EmergencyStateDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"forceBatchNum\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"lastGlobalExitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"}],\"name\":\"ForceBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"OverridePendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"storedStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"provedStateRoot\",\"type\":\"bytes32\"}],\"name\":\"ProveNonDeterministicPendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"}],\"name\":\"SequenceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"}],\"name\":\"SequenceForceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"SetAdmin\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"newForceBatchAllowed\",\"type\":\"bool\"}],\"name\":\"SetForceBatchAllowed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newPendingStateTimeout\",\"type\":\"uint64\"}],\"name\":\"SetPendingStateTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTrustedAggregator\",\"type\":\"address\"}],\"name\":\"SetTrustedAggregator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newTrustedAggregatorTimeout\",\"type\":\"uint64\"}],\"name\":\"SetTrustedAggregatorTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTrustedSequencer\",\"type\":\"address\"}],\"name\":\"SetTrustedSequencer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"newTrustedSequencerURL\",\"type\":\"string\"}],\"name\":\"SetTrustedSequencerURL\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"TrustedVerifyBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatches\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"FORCE_BATCH_TIMEOUT\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"HALT_AGGREGATION_TIMEOUT\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_BATCH_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_VERIFY_BATCHES\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"sequencedBatchNum\",\"type\":\"uint64\"}],\"name\":\"activateEmergencyState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"batchNumToStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"contractIBridge\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateBatchFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateRewardPerBatch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chainID\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"consolidatePendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deactivateEmergencyState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"maticAmount\",\"type\":\"uint256\"}],\"name\":\"forceBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"forceBatchAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"forcedBatches\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"getInputSnarkBytes\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastVerifiedBatch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"globalExitRootManager\",\"outputs\":[{\"internalType\":\"contractIGlobalExitRootManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIGlobalExitRootManager\",\"name\":\"_globalExitRootManager\",\"type\":\"address\"},{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"_matic\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"_rollupVerifier\",\"type\":\"address\"},{\"internalType\":\"contractIBridge\",\"name\":\"_bridgeAddress\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"trustedSequencer\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateTimeout\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"forceBatchAllowed\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"trustedAggregator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"trustedAggregatorTimeout\",\"type\":\"uint64\"}],\"internalType\":\"structProofOfEfficiency.InitializePackedParameters\",\"name\":\"initializePackedParameters\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"genesisRoot\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_trustedSequencerURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_networkName\",\"type\":\"string\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isEmergencyState\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"isPendingStateConsolidable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastBatchSequenced\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastForceBatch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastForceBatchSequenced\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastPendingState\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastPendingStateConsolidated\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastTimestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastVerifiedBatch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"matic\",\"outputs\":[{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"initPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofA\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"proofB\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofC\",\"type\":\"uint256[2]\"}],\"name\":\"overridePendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingStateTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"pendingStateTransitions\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"initPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofA\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"proofB\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofC\",\"type\":\"uint256[2]\"}],\"name\":\"proveNonDeterministicPendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupVerifier\",\"outputs\":[{\"internalType\":\"contractIVerifierRollup\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"globalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minForcedTimestamp\",\"type\":\"uint64\"}],\"internalType\":\"structProofOfEfficiency.BatchData[]\",\"name\":\"batches\",\"type\":\"tuple[]\"}],\"name\":\"sequenceBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"globalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"minForcedTimestamp\",\"type\":\"uint64\"}],\"internalType\":\"structProofOfEfficiency.ForcedBatchData[]\",\"name\":\"batches\",\"type\":\"tuple[]\"}],\"name\":\"sequenceForceBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"sequencedBatches\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"accInputHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequencedTimestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"previousLastBatchSequenced\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"newForceBatchAllowed\",\"type\":\"bool\"}],\"name\":\"setForceBatchAllowed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newPendingStateTimeout\",\"type\":\"uint64\"}],\"name\":\"setPendingStateTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newTrustedAggregator\",\"type\":\"address\"}],\"name\":\"setTrustedAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newTrustedAggregatorTimeout\",\"type\":\"uint64\"}],\"name\":\"setTrustedAggregatorTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newTrustedSequencer\",\"type\":\"address\"}],\"name\":\"setTrustedSequencer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newTrustedSequencerURL\",\"type\":\"string\"}],\"name\":\"setTrustedSequencerURL\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedAggregator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedAggregatorTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedSequencer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedSequencerURL\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofA\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"proofB\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofC\",\"type\":\"uint256[2]\"}],\"name\":\"trustedVerifyBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofA\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"proofB\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"proofC\",\"type\":\"uint256[2]\"}],\"name\":\"verifyBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50614e6a806100206000396000f3fe608060405234801561001057600080fd5b506004361061035d5760003560e01c80638b48931e116101d3578063d02103ca11610104578063e7a7ed02116100a2578063edc411211161007c578063edc411211461080c578063f14916d61461081f578063f2fde38b14610832578063f851a4401461084557600080fd5b8063e7a7ed02146107cc578063e8bf92ed146107e6578063eaeb077b146107f957600080fd5b8063d939b315116100de578063d939b3151461078e578063dbc16976146107a8578063e11f3f18146107b0578063e217cfd6146107c357600080fd5b8063d02103ca14610754578063d8d1091b14610767578063d8f54db01461077a57600080fd5b8063adc879e911610171578063b6b0b0971161014b578063b6b0b09714610707578063c0ed84e01461071f578063c89e42df14610727578063cfa8ed471461073a57600080fd5b8063adc879e914610684578063b02286c01461069e578063b4d63f58146106a757600080fd5b806399f5634e116101ad57806399f5634e146106565780639c9f3dfe1461065e578063a3c573eb14610671578063ab9fc5ef1461062857600080fd5b80638b48931e146106285780638c4a0af7146106325780638da5cb5b1461064557600080fd5b80634a910e6a116102ad578063715018a61161024b5780637abaf3e6116102255780637abaf3e6146105855780637fcb36531461058d578063837a4738146105a0578063841b24d71461060e57600080fd5b8063715018a6146105575780637215541a1461055f57806375c508b31461057257600080fd5b806360943d6a1161028757806360943d6a146104fe5780636b8616ce146105115780636ff512cc14610531578063704b6c021461054457600080fd5b80634a910e6a146104b55780635392c5e0146104c8578063542028d5146104f657600080fd5b8063394218e91161031a57806345605267116102f4578063456052671461045b578063458c0477146104755780634834a343146104885780634a1a89a71461049b57600080fd5b8063394218e9146104195780633c1582671461042e578063423fa8561461044157600080fd5b8063107bf28c1461036257806315064c961461038057806319d8ac611461039d578063220d7899146103c857806329878983146103db578063383b3be814610406575b600080fd5b61036a610858565b6040516103779190613f86565b60405180910390f35b60655461038d9060ff1681565b6040519015158152602001610377565b6068546103b0906001600160401b031681565b6040516001600160401b039091168152602001610377565b61036a6103d6366004613fbc565b6108e6565b606a546103ee906001600160a01b031681565b6040516001600160a01b039091168152602001610377565b61038d610414366004614009565b610aa9565b61042c610427366004614009565b610af0565b005b61042c61043c366004614146565b610c4d565b6068546103b090600160401b90046001600160401b031681565b6068546103b090600160801b90046001600160401b031681565b6072546103b0906001600160401b031681565b61042c610496366004614284565b6113c4565b6072546103b090600160401b90046001600160401b031681565b61042c6104c3366004614009565b6115d2565b6104e86104d6366004614009565b606d6020526000908152604090205481565b604051908152602001610377565b61036a611880565b61042c61050c366004614326565b61188d565b6104e861051f366004614009565b60666020526000908152604090205481565b61042c61053f3660046143fd565b611bb7565b61042c6105523660046143fd565b611c3b565b61042c611cb3565b61042c61056d366004614009565b611cc7565b61042c61058036600461441a565b611e81565b6104e8611f10565b6069546103b0906001600160401b031681565b6105e36105ae3660046144b8565b6071602052600090815260409020805460018201546002909201546001600160401b0380831693600160401b90930416919084565b604080516001600160401b039586168152949093166020850152918301526060820152608001610377565b6072546103b090600160c01b90046001600160401b031681565b6103b062093a8081565b61042c6106403660046144df565b611f66565b6033546001600160a01b03166103ee565b6104e8611fdd565b61042c61066c366004614009565b6120be565b6070546103ee906001600160a01b031681565b606c546103b090600160a81b90046001600160401b031681565b6104e861ea6081565b6106e26106b5366004614009565b606760205260009081526040902080546001909101546001600160401b0380821691600160401b90041683565b604080519384526001600160401b039283166020850152911690820152606001610377565b6065546103ee9061010090046001600160a01b031681565b6103b061221c565b61042c6107353660046144fc565b612269565b6069546103ee90600160401b90046001600160a01b031681565b606c546103ee906001600160a01b031681565b61042c610775366004614530565b6122cf565b606c5461038d90600160a01b900460ff1681565b6072546103b090600160801b90046001600160401b031681565b61042c612766565b61042c6107be36600461441a565b612822565b6103b06103e881565b6068546103b090600160c01b90046001600160401b031681565b606b546103ee906001600160a01b031681565b61042c610807366004614622565b612951565b61042c61081a366004614284565b612c83565b61042c61082d3660046143fd565b612da3565b61042c6108403660046143fd565b612e1b565b6073546103ee906001600160a01b031681565b606f805461086590614666565b80601f016020809104026020016040519081016040528092919081815260200182805461089190614666565b80156108de5780601f106108b3576101008083540402835291602001916108de565b820191906000526020600020905b8154815290600101906020018083116108c157829003601f168201915b505050505081565b6001600160401b0380861660008181526067602052604080822054938816825290205460609291158061091857508115155b61099d5760405162461bcd60e51b815260206004820152604560248201527f50726f6f664f66456666696369656e63793a3a676574496e707574536e61726b60448201527f42797465733a206f6c64416363496e7075744861736820646f6573206e6f7420606482015264195e1a5cdd60da1b608482015260a4015b60405180910390fd5b80610a1e5760405162461bcd60e51b815260206004820152604560248201527f50726f6f664f66456666696369656e63793a3a676574496e707574536e61726b60448201527f42797465733a206e6577416363496e7075744861736820646f6573206e6f7420606482015264195e1a5cdd60da1b608482015260a401610994565b606c54604080516bffffffffffffffffffffffff193360601b166020820152603481019790975260548701939093526001600160c01b031960c0998a1b81166074880152600160a81b909104891b8116607c870152608486019490945260a485015260c4840194909452509290931b90911660e4830152805180830360cc01815260ec909201905290565b6072546001600160401b0382811660009081526071602052604081205490924292610adf92600160801b909204811691166146b6565b6001600160401b0316111592915050565b6073546001600160a01b03163314610b1a5760405162461bcd60e51b8152600401610994906146e1565b62093a806001600160401b0382161115610b465760405162461bcd60e51b815260040161099490614729565b60655460ff16610bed576072546001600160401b03600160c01b909104811690821610610bed5760405162461bcd60e51b815260206004820152604960248201527f50726f6f664f66456666696369656e63793a3a7365745472757374656441676760448201527f72656761746f7254696d656f75743a206e65772074696d656f7574206d757374606482015268103132903637bbb2b960b91b608482015260a401610994565b607280546001600160c01b0316600160c01b6001600160401b038481168202929092179283905560405192041681527f1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1906020015b60405180910390a150565b60655460ff1615610c705760405162461bcd60e51b815260040161099490614799565b606954600160401b90046001600160a01b03163314610cf75760405162461bcd60e51b815260206004820152603f60248201527f50726f6f664f66456666696369656e63793a3a6f6e6c7954727573746564536560448201527f7175656e6365723a206f6e6c7920747275737465642073657175656e636572006064820152608401610994565b805180610d655760405162461bcd60e51b81526020600482015260426024820152600080516020614db583398151915260448201527f65733a204174206c65617374206d7573742073657175656e63652031206261746064820152610c6d60f31b608482015260a401610994565b6103e88110610d865760405162461bcd60e51b815260040161099490614801565b6068546001600160401b03600160401b82048116600081815260676020526040812054838516949293600160801b90930490921691905b858110156111f4576000878281518110610dd957610dd9614856565b60200260200101519050600081606001516001600160401b03161115610f795783610e038161486c565b945050600081600001518051906020012082602001518360600151604051602001610e3093929190614892565b60408051601f1981840301815291815281516020928301206001600160401b038816600090815260669093529120549091508114610ecf5760405162461bcd60e51b81526020600482015260426024820152600080516020614db583398151915260448201527f65733a20466f7263656420626174636865732064617461206d757374206d61746064820152610c6d60f31b608482015260a401610994565b81606001516001600160401b031682604001516001600160401b03161015610f735760405162461bcd60e51b815260206004820152605d6024820152600080516020614db583398151915260448201527f65733a20466f7263656420626174636865732074696d657374616d70206d757360648201527f7420626520626967676572206f7220657175616c207468616e206d696e000000608482015260a401610994565b506110ee565b602081015115806110025750606c5460208201516040516312bd9b1960e11b81526001600160a01b039092169163257b363291610fbc9160040190815260200190565b6020604051808303816000875af1158015610fdb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fff91906148b5565b15155b6110625760405162461bcd60e51b815260206004820152603f6024820152600080516020614db583398151915260448201527f65733a20476c6f62616c206578697420726f6f74206d757374206578697374006064820152608401610994565b80515161ea60116110ee5760405162461bcd60e51b815260206004820152604a60248201527f50726f6f664f664566666963696550656e64696e67537461746563793a3a736560448201527f7175656e6365426174636865733a205472616e73616374696f6e73206279746560648201526973206f766572666c6f7760b01b608482015260a401610994565b856001600160401b031681604001516001600160401b03161015801561112157504281604001516001600160401b031611155b61118c5760405162461bcd60e51b81526020600482015260426024820152600080516020614db583398151915260448201527f65733a2054696d657374616d70206d75737420626520696e736964652072616e606482015261676560f01b608482015260a401610994565b805180516020918201208183015160408085015190516111b39488949392913391016148ce565b60405160208183030381529060405280519060200120925084806111d69061486c565b955050806040015195505080806111ec90614912565b915050610dbd565b506068546001600160401b03600160c01b9091048116908316111561126f5760405162461bcd60e51b815260206004820152603a6024820152600080516020614db583398151915260448201527f65733a20466f7263652062617463686573206f766572666c6f770000000000006064820152608401610994565b60685460009061128f90600160801b90046001600160401b03168461492b565b6112a2906001600160401b031687614953565b60408051606081018252848152426001600160401b03908116602080840191825260688054600160401b9081900485168688019081528c861660008181526067909552979093209551865592516001909501805492519585166001600160801b03199384161795851684029590951790945583548b841691161793029290921767ffffffffffffffff60801b1916600160801b928716929092029190911790559050611375333083611352611f10565b61135c919061496a565b60655461010090046001600160a01b0316929190612e91565b61137d612efc565b606854604051600160401b9091046001600160401b0316907f303446e6a8cb73c83dff421c0b1d5e5ce0719dab1bff13660fc254e58cc17fce90600090a250505050505050565b60655460ff16156113e75760405162461bcd60e51b815260040161099490614799565b6072546001600160401b03878116600090815260676020526040902060010154429261141e92600160c01b909104811691166146b6565b6001600160401b0316111561149a5760405162461bcd60e51b81526020600482015260486024820152600080516020614e1583398151915260448201527f3a20747275737465642061676772656761746f722074696d656f7574206e6f7460648201526708195e1c1a5c995960c21b608482015260a401610994565b6103e86114a7888861492b565b6001600160401b0316106114cd5760405162461bcd60e51b815260040161099490614801565b6114dd8888888888888888612fa0565b6114e5612efc565b607280546001600160401b03169060006114fe8361486c565b82546001600160401b039182166101009390930a92830292820219169190911790915560408051608081018252428316815289831660208083018281528385018c8152606085018c8152607254881660009081526071909452928690209451855492518816600160401b026001600160801b0319909316971696909617178355935160018301559251600290910155513392507f9c72852172521097ba7e1482e6b44b351323df0155f97f4ea18fcec28e1f5966906115c09088815260200190565b60405180910390a35050505050505050565b6001600160401b038116158015906115ff57506072546001600160401b03600160401b9091048116908216115b801561161a57506072546001600160401b0390811690821611155b61169d5760405162461bcd60e51b815260206004820152604860248201527f50726f6f664f66456666696369656e63793a3a636f6e736f6c6964617465506560448201527f6e64696e6753746174653a2070656e64696e6753746174654e756d206d757374606482015267081a5b9d985b1a5960c21b608482015260a401610994565b606a546001600160a01b03163314611750576116b881610aa9565b6117505760405162461bcd60e51b815260206004820152605960248201527f50726f6f664f66456666696369656e63793a3a636f6e736f6c6964617465506560448201527f6e64696e6753746174653a2070656e64696e67207374617465206973206e6f7460648201527f20726561647920746f20626520636f6e736f6c69646174656400000000000000608482015260a401610994565b6001600160401b038181166000818152607160209081526040808320805460698054600160401b9283900490981667ffffffffffffffff19909816881790556002820154878652606d9094529382902092909255607280546fffffffffffffffff000000000000000019169390940292909217909255606c54600183015491516333d6247d60e01b815260048101929092529192916001600160a01b0316906333d6247d90602401600060405180830381600087803b15801561181257600080fd5b505af1158015611826573d6000803e3d6000fd5b50505050336001600160a01b0316816001600160401b03167f01f7d32e3b3278bace940a581067c87090c1aa09809730dd4ca002320c3a3cfa846002015460405161187391815260200190565b60405180910390a3505050565b606e805461086590614666565b600054610100900460ff16158080156118ad5750600054600160ff909116105b806118c75750303b1580156118c7575060005460ff166001145b61192a5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610994565b6000805460ff19166001179055801561194d576000805461ff0019166101001790555b606c80546001600160a01b03199081166001600160a01b038c81169190911790925560658054610100600160a81b0319166101008c851602179055606b805482168a8416179055607080549091169188169190911790556119b160208601866143fd565b607380546001600160a01b0319166001600160a01b03929092169190911790556119e160608601604087016143fd565b606980546001600160a01b0392909216600160401b02600160401b600160e01b0319909216919091179055611a1c60c0860160a087016143fd565b606a80546001600160a01b0319166001600160a01b039290921691909117905560008052606d6020527fda90043ba5b4096ba14704bc227ab0d3167da15b887e62ab2e76e37daa711356849055611a7960e0860160c08701614009565b607280546001600160401b0392909216600160c01b026001600160c01b03909216919091179055611ab06040860160208701614009565b606c80546001600160401b0392909216600160a81b0267ffffffffffffffff60a81b19909216919091179055611aec6080860160608701614009565b607280546001600160401b0392909216600160801b0267ffffffffffffffff60801b19909216919091179055611b2860a08601608087016144df565b606c8054911515600160a01b0260ff60a01b19909216919091179055606e611b5084826149d7565b50606f611b5d83826149d7565b50611b6661348a565b8015611bac576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050505050505050565b6073546001600160a01b03163314611be15760405162461bcd60e51b8152600401610994906146e1565b60698054600160401b600160e01b031916600160401b6001600160a01b038416908102919091179091556040519081527ff54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc090602001610c42565b6073546001600160a01b03163314611c655760405162461bcd60e51b8152600401610994906146e1565b607380546001600160a01b0319166001600160a01b0383169081179091556040519081527f5a272403b402d892977df56625f4164ccaf70ca3863991c43ecfe76a6905b0a190602001610c42565b611cbb6134fe565b611cc56000613558565b565b6033546001600160a01b03163314611e76576072546000906001600160401b031615611d1857506072546001600160401b03908116600090815260716020526040902054600160401b900416611d26565b506069546001600160401b03165b80826001600160401b031611611dae5760405162461bcd60e51b815260206004820152604160248201527f50726f6f664f66456666696369656e63793a3a6163746976617465456d65726760448201527f656e637953746174653a20426174636820616c726561647920766572696669656064820152601960fa1b608482015260a401610994565b6001600160401b038083166000908152606760205260409020600101544291611ddc9162093a8091166146b6565b6001600160401b03161115611e745760405162461bcd60e51b815260206004820152605260248201527f50726f6f664f66456666696369656e63793a3a6163746976617465456d65726760448201527f656e637953746174653a206167677265676174696f6e2068616c742074696d656064820152711bdd5d081a5cc81b9bdd08195e1c1a5c995960721b608482015260a401610994565b505b611e7e6135aa565b50565b60655460ff1615611ea45760405162461bcd60e51b815260040161099490614799565b611eb589898989898989898961361a565b6001600160401b0386166000908152606d60209081526040918290205482519081529081018690527f1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010910160405180910390a1611bac6135aa565b6068546000906001600160401b03600160801b8204811691611f3c91600160c01b9091041660016146b6565b611f46919061492b565b611f61906001600160401b0316670de0b6b3a764000061496a565b905090565b6073546001600160a01b03163314611f905760405162461bcd60e51b8152600401610994906146e1565b606c8054821515600160a01b0260ff60a01b199091161790556040517fbacda50a4a8575be1d91a7ebe29ee45056f3a94f12a2281eb6b43afa33bcefe690610c4290831515815260200190565b6065546040516370a0823160e01b815230600482015260009182916101009091046001600160a01b0316906370a0823190602401602060405180830381865afa15801561202e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061205291906148b5565b9050600061205e61221c565b6068546001600160401b03600160401b820481169161208e91600160801b8204811691600160c01b90041661492b565b61209891906146b6565b6120a2919061492b565b6001600160401b031690506120b78183614aac565b9250505090565b6073546001600160a01b031633146120e85760405162461bcd60e51b8152600401610994906146e1565b60725462093a80600160801b9091046001600160401b0316111561211e5760405162461bcd60e51b815260040161099490614729565b60655460ff166121c1576072546001600160401b03600160801b9091048116908216106121c15760405162461bcd60e51b8152602060048201526044602482018190527f50726f6f664f66456666696369656e63793a3a73657450656e64696e67537461908201527f746554696d656f75743a206e65772074696d656f7574206d757374206265206c60648201526337bbb2b960e11b608482015260a401610994565b6072805467ffffffffffffffff60801b1916600160801b6001600160401b038416908102919091179091556040519081527fc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c7590602001610c42565b6072546000906001600160401b03161561225957506072546001600160401b03908116600090815260716020526040902054600160401b90041690565b506069546001600160401b031690565b6073546001600160a01b031633146122935760405162461bcd60e51b8152600401610994906146e1565b606e61229f82826149d7565b507f6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b2081604051610c429190613f86565b60655460ff16156122f25760405162461bcd60e51b815260040161099490614799565b606c54600160a01b900460ff1615156001146123205760405162461bcd60e51b815260040161099490614ac0565b80518061238e5760405162461bcd60e51b81526020600482015260426024820152600080516020614df583398151915260448201527f42617463683a204d75737420666f726365206174206c656173742031206261746064820152610c6d60f31b608482015260a401610994565b6103e881106123af5760405162461bcd60e51b815260040161099490614801565b6068546001600160401b03600160c01b82048116916123d7918491600160801b900416614b2e565b11156124395760405162461bcd60e51b815260206004820152603a6024820152600080516020614df583398151915260448201527f42617463683a20466f72636520626174636820696e76616c69640000000000006064820152608401610994565b6068546001600160401b03600160401b820481166000818152606760205260408120549193600160801b9004909216915b8481101561266f57600086828151811061248657612486614856565b60200260200101519050838061249b9061486c565b9450506000816000015180519060200120826020015183604001516040516020016124c893929190614892565b60408051601f1981840301815291815281516020928301206001600160401b03881660009081526066909352912054909150811461256c5760405162461bcd60e51b81526020600482015260476024820152600080516020614df583398151915260448201527f426174636865733a20466f7263656420626174636865732064617461206d75736064820152660e840dac2e8c6d60cb1b608482015260a401610994565b612577600188614953565b8303612611574262093a80836040015161259191906146b6565b6001600160401b031611156126115760405162461bcd60e51b815260206004820152604c6024820152600080516020614df583398151915260448201527f42617463683a20466f72636564206261746368206973206e6f7420696e20746960648201526b1b595bdd5d081c195c9a5bd960a21b608482015260a401610994565b8151805160209182012081840151604051612634938893929142913391016148ce565b60405160208183030381529060405280519060200120935085806126579061486c565b9650505050808061266790614912565b91505061246a565b506068805467ffffffffffffffff1916426001600160401b03908116918217808455604080516060810182528681526020808201958652600160401b9384900485168284019081528a861660008181526067909352848320935184559651600193909301805491519387166001600160801b031990921691909117928616850292909217909155855477ffffffffffffffffffffffffffffffff0000000000000000191694830267ffffffffffffffff60801b191694909417600160801b88851602179485905551930416917f648a61dd2438f072f5a1960939abd30f37aea80d2e94c9792ad142d3e0a490a49190a25050505050565b60655460ff166127885760405162461bcd60e51b815260040161099490614b46565b6073546001600160a01b031633146127b25760405162461bcd60e51b8152600401610994906146e1565b607060009054906101000a90046001600160a01b03166001600160a01b031663dbc169766040518163ffffffff1660e01b8152600401600060405180830381600087803b15801561280257600080fd5b505af1158015612816573d6000803e3d6000fd5b50505050611cc5613be5565b606a546001600160a01b0316331461284c5760405162461bcd60e51b815260040161099490614ba3565b61285d89898989898989898961361a565b6069805467ffffffffffffffff19166001600160401b038881169182179092556000908152606d6020526040902085905560725416156128a857607280546001600160801b03191690555b606c546040516333d6247d60e01b8152600481018790526001600160a01b03909116906333d6247d90602401600060405180830381600087803b1580156128ee57600080fd5b505af1158015612902573d6000803e3d6000fd5b50506040518681523392506001600160401b03891691507fcc1b5520188bf1dd3e63f98164b577c4d75c11a619ddea692112f0d1aec4cf729060200160405180910390a3505050505050505050565b60655460ff16156129745760405162461bcd60e51b815260040161099490614799565b606c54600160a01b900460ff1615156001146129a25760405162461bcd60e51b815260040161099490614ac0565b60006129ac611f10565b905081811115612a165760405162461bcd60e51b815260206004820152602f60248201527f50726f6f664f66456666696369656e63793a3a666f72636542617463683a206e60448201526e6f7420656e6f756768206d6174696360881b6064820152608401610994565b61ea60835110612a8e5760405162461bcd60e51b815260206004820152603a60248201527f50726f6f664f66456666696369656e63793a3a666f72636542617463683a205460448201527f72616e73616374696f6e73206279746573206f766572666c6f770000000000006064820152608401610994565b606554612aab9061010090046001600160a01b0316333084612e91565b606c5460408051633ed691ef60e01b815290516000926001600160a01b031691633ed691ef9160048083019260209291908290030181865afa158015612af5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b1991906148b5565b60688054919250600160c01b9091046001600160401b0316906018612b3d8361486c565b91906101000a8154816001600160401b0302191690836001600160401b031602179055505083805190602001208142604051602001612b7e93929190614892565b60408051808303601f190181529181528151602092830120606854600160c01b90046001600160401b031660009081526066909352912055323303612c2257606854604080518381523360208201526060918101829052600091810191909152600160c01b9091046001600160401b0316907ff94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc9319060800160405180910390a2612c7d565b606860189054906101000a90046001600160401b03166001600160401b03167ff94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc931823387604051612c7493929190614c0a565b60405180910390a25b50505050565b606a546001600160a01b03163314612cad5760405162461bcd60e51b815260040161099490614ba3565b612cbd8888888888888888612fa0565b6069805467ffffffffffffffff19166001600160401b038881169182179092556000908152606d602052604090208590556072541615612d0857607280546001600160801b03191690555b606c546040516333d6247d60e01b8152600481018790526001600160a01b03909116906333d6247d90602401600060405180830381600087803b158015612d4e57600080fd5b505af1158015612d62573d6000803e3d6000fd5b50506040518681523392506001600160401b03891691507f0c0ce073a7d7b5850c04ccc4b20ee7d3179d5f57d0ac44399565792c0f72fce7906020016115c0565b6073546001600160a01b03163314612dcd5760405162461bcd60e51b8152600401610994906146e1565b606a80546001600160a01b0319166001600160a01b0383169081179091556040519081527f61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca90602001610c42565b612e236134fe565b6001600160a01b038116612e885760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610994565b611e7e81613558565b6040516001600160a01b0380851660248301528316604482015260648101829052612c7d9085906323b872dd60e01b906084015b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b031990931692909217909152613c3c565b6072546001600160401b03600160401b8204811691161115611cc557607254600090612f3990600160401b90046001600160401b031660016146b6565b9050612f4481610aa9565b15611e7e57607254600090600290612f669084906001600160401b031661492b565b612f709190614c3d565b612f7a90836146b6565b9050612f8581610aa9565b15612f9757612f93816115d2565b5050565b612f93826115d2565b600080612fab61221c565b90506001600160401b038a1615613107576072546001600160401b03908116908b1611156130555760405162461bcd60e51b815260206004820152605d6024820152600080516020614e1583398151915260448201527f3a2070656e64696e6753746174654e756d206d757374206265206c657373206f60648201527f7220657175616c207468616e206c61737450656e64696e675374617465000000608482015260a401610994565b6001600160401b03808b1660009081526071602052604090206002810154815490945090918b8116600160401b90920416146131015760405162461bcd60e51b81526020600482015260516024820152600080516020614e1583398151915260448201527f3a20696e69744e756d4261746368206d757374206d61746368207468652070656064820152700dcc8d2dcce40e6e8c2e8ca40c4c2e8c6d607b1b608482015260a401610994565b5061323b565b6001600160401b0389166000908152606d60205260409020549150816131945760405162461bcd60e51b81526020600482015260486024820152600080516020614e1583398151915260448201527f3a20696e69744e756d426174636820737461746520726f6f7420646f6573206e6064820152671bdd08195e1a5cdd60c21b608482015260a401610994565b806001600160401b0316896001600160401b0316111561323b5760405162461bcd60e51b81526020600482015260626024820152600080516020614e1583398151915260448201527f3a20696e69744e756d4261746368206d757374206265206c657373206f72206560648201527f7175616c207468616e2063757272656e744c61737456657269666965644261746084820152610c6d60f31b60a482015260c401610994565b806001600160401b0316886001600160401b0316116132d65760405162461bcd60e51b815260206004820152605c6024820152600080516020614e1583398151915260448201527f3a2066696e616c4e65774261746368206d75737420626520626967676572207460648201527f68616e2063757272656e744c6173745665726966696564426174636800000000608482015260a401610994565b60006132e58a8a8a868b6108e6565b905060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000160028360405161331a9190614c63565b602060405180830381855afa158015613337573d6000803e3d6000fd5b5050506040513d601f19601f8201168201806040525081019061335a91906148b5565b6133649190614c7f565b606b546040805160208101825283815290516343753b4d60e01b81529293506001600160a01b03909116916343753b4d916133a8918b918b918b9190600401614c93565b602060405180830381865afa1580156133c5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906133e99190614d0d565b61343b5760405162461bcd60e51b815260206004820152602f6024820152600080516020614e1583398151915260448201526e1d1024a72b20a624a22fa82927a7a360891b6064820152608401610994565b61347c33613449858d61492b565b6001600160401b031661345a611fdd565b613464919061496a565b60655461010090046001600160a01b03169190613d13565b505050505050505050505050565b600054610100900460ff166134f55760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b6064820152608401610994565b611cc533613558565b6033546001600160a01b03163314611cc55760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610994565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b607060009054906101000a90046001600160a01b03166001600160a01b0316632072f6c56040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156135fa57600080fd5b505af115801561360e573d6000803e3d6000fd5b50505050611cc5613d43565b60006001600160401b038a16156137a8576072546001600160401b03908116908b1611156136de5760405162461bcd60e51b81526020600482015260716024820152600080516020614dd583398151915260448201527f6d696e697374696350656e64696e6753746174653a2070656e64696e6753746160648201527f74654e756d206d757374206265206c657373206f7220657175616c207468616e608482015270206c61737450656e64696e67537461746560781b60a482015260c401610994565b506001600160401b03808a1660009081526071602052604090206002810154815490928a8116600160401b90920416146137a25760405162461bcd60e51b81526020600482015260656024820152600080516020614dd583398151915260448201527f6d696e697374696350656e64696e6753746174653a20696e69744e756d42617460648201527f6368206d757374206d61746368207468652070656e64696e67207374617465206084820152640c4c2e8c6d60db1b60a482015260c401610994565b50613901565b506001600160401b0387166000908152606d6020526040902054806138495760405162461bcd60e51b815260206004820152605c6024820152600080516020614dd583398151915260448201527f6d696e697374696350656e64696e6753746174653a20696e69744e756d42617460648201527f636820737461746520726f6f7420646f6573206e6f7420657869737400000000608482015260a401610994565b6069546001600160401b0390811690891611156139015760405162461bcd60e51b81526020600482015260766024820152600080516020614dd583398151915260448201527f6d696e697374696350656e64696e6753746174653a20696e69744e756d42617460648201527f6368206d757374206265206c657373206f7220657175616c207468616e2063756084820152750e4e4cadce898c2e6e8accae4d2ccd2cac884c2e8c6d60531b60a482015260c401610994565b6072546001600160401b03908116908a16118015906139315750896001600160401b0316896001600160401b0316115b801561395257506072546001600160401b03600160401b9091048116908a16115b61396e5760405162461bcd60e51b815260040161099490614d2a565b6001600160401b03898116600090815260716020526040902054600160401b90048116908816146139b15760405162461bcd60e51b815260040161099490614d2a565b60006139c0898989858a6108e6565b905060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f00000016002836040516139f59190614c63565b602060405180830381855afa158015613a12573d6000803e3d6000fd5b5050506040513d601f19601f82011682018060405250810190613a3591906148b5565b613a3f9190614c7f565b606b546040805160208101825283815290516343753b4d60e01b81529293506001600160a01b03909116916343753b4d91613a83918a918a918a9190600401614c93565b602060405180830381865afa158015613aa0573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190613ac49190614d0d565b613b305760405162461bcd60e51b81526020600482015260436024820152600080516020614dd583398151915260448201527f6d696e697374696350656e64696e6753746174653a20494e56414c49445f505260648201526227a7a360e91b608482015260a401610994565b6001600160401b038b1660009081526071602052604090206002015487900361347c5760405162461bcd60e51b81526020600482015260676024820152600080516020614dd583398151915260448201527f6d696e697374696350656e64696e6753746174653a2073746f72656420726f6f60648201527f74206d75737420626520646966666572656e74207468616e206e6577207374616084820152661d19481c9bdbdd60ca1b60a482015260c401610994565b60655460ff16613c075760405162461bcd60e51b815260040161099490614b46565b6065805460ff191690556040517f1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b390600090a1565b6000613c91826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b0316613d9e9092919063ffffffff16565b805190915015613d0e5780806020019051810190613caf9190614d0d565b613d0e5760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b6064820152608401610994565b505050565b6040516001600160a01b038316602482015260448101829052613d0e90849063a9059cbb60e01b90606401612ec5565b60655460ff1615613d665760405162461bcd60e51b815260040161099490614799565b6065805460ff191660011790556040517f2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a549790600090a1565b6060613dad8484600085613db5565b949350505050565b606082471015613e165760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b6064820152608401610994565b600080866001600160a01b03168587604051613e329190614c63565b60006040518083038185875af1925050503d8060008114613e6f576040519150601f19603f3d011682016040523d82523d6000602084013e613e74565b606091505b5091509150613e8587838387613e90565b979650505050505050565b60608315613eff578251600003613ef8576001600160a01b0385163b613ef85760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610994565b5081613dad565b613dad8383815115613f145781518083602001fd5b8060405162461bcd60e51b81526004016109949190613f86565b60005b83811015613f49578181015183820152602001613f31565b83811115612c7d5750506000910152565b60008151808452613f72816020860160208601613f2e565b601f01601f19169290920160200192915050565b602081526000613f996020830184613f5a565b9392505050565b80356001600160401b0381168114613fb757600080fd5b919050565b600080600080600060a08688031215613fd457600080fd5b613fdd86613fa0565b9450613feb60208701613fa0565b94979496505050506040830135926060810135926080909101359150565b60006020828403121561401b57600080fd5b613f9982613fa0565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b038111828210171561405c5761405c614024565b60405290565b604051606081016001600160401b038111828210171561405c5761405c614024565b604051601f8201601f191681016001600160401b03811182821017156140ac576140ac614024565b604052919050565b60006001600160401b038211156140cd576140cd614024565b5060051b60200190565b600082601f8301126140e857600080fd5b81356001600160401b0381111561410157614101614024565b614114601f8201601f1916602001614084565b81815284602083860101111561412957600080fd5b816020850160208301376000918101602001919091529392505050565b6000602080838503121561415957600080fd5b82356001600160401b038082111561417057600080fd5b818501915085601f83011261418457600080fd5b8135614197614192826140b4565b614084565b81815260059190911b830184019084810190888311156141b657600080fd5b8585015b8381101561424f578035858111156141d25760008081fd5b86016080818c03601f19018113156141ea5760008081fd5b6141f261403a565b89830135888111156142045760008081fd5b6142128e8c838701016140d7565b8252506040808401358b830152606061422c818601613fa0565b8284015261423b848601613fa0565b9083015250855250509186019186016141ba565b5098975050505050505050565b806040810183101561426d57600080fd5b92915050565b806080810183101561426d57600080fd5b6000806000806000806000806101a0898b0312156142a157600080fd5b6142aa89613fa0565b97506142b860208a01613fa0565b96506142c660408a01613fa0565b955060608901359450608089013593506142e38a60a08b0161425c565b92506142f28a60e08b01614273565b91506143028a6101608b0161425c565b90509295985092959890939650565b6001600160a01b0381168114611e7e57600080fd5b600080600080600080600080888a036101c081121561434457600080fd5b893561434f81614311565b985060208a013561435f81614311565b975060408a013561436f81614311565b965060608a013561437f81614311565b955060e0607f198201121561439357600080fd5b5060808901935061016089013592506101808901356001600160401b03808211156143bd57600080fd5b6143c98c838d016140d7565b93506101a08b01359150808211156143e057600080fd5b506143ed8b828c016140d7565b9150509295985092959890939650565b60006020828403121561440f57600080fd5b8135613f9981614311565b60008060008060008060008060006101c08a8c03121561443957600080fd5b6144428a613fa0565b985061445060208b01613fa0565b975061445e60408b01613fa0565b965061446c60608b01613fa0565b955060808a0135945060a08a013593506144898b60c08c0161425c565b92506144998b6101008c01614273565b91506144a98b6101808c0161425c565b90509295985092959850929598565b6000602082840312156144ca57600080fd5b5035919050565b8015158114611e7e57600080fd5b6000602082840312156144f157600080fd5b8135613f99816144d1565b60006020828403121561450e57600080fd5b81356001600160401b0381111561452457600080fd5b613dad848285016140d7565b6000602080838503121561454357600080fd5b82356001600160401b038082111561455a57600080fd5b818501915085601f83011261456e57600080fd5b813561457c614192826140b4565b81815260059190911b8301840190848101908883111561459b57600080fd5b8585015b8381101561424f578035858111156145b75760008081fd5b86016060818c03601f19018113156145cf5760008081fd5b6145d7614062565b89830135888111156145e95760008081fd5b6145f78e8c838701016140d7565b8252506040808401358b83015261460f838501613fa0565b908201528552505091860191860161459f565b6000806040838503121561463557600080fd5b82356001600160401b0381111561464b57600080fd5b614657858286016140d7565b95602094909401359450505050565b600181811c9082168061467a57607f821691505b60208210810361469a57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b60006001600160401b038083168185168083038211156146d8576146d86146a0565b01949350505050565b60208082526028908201527f50726f6f664f66456666696369656e63793a3a6f6e6c7941646d696e3a206f6e604082015267363c9030b236b4b760c11b606082015260800190565b6020808252604a908201527f50726f6f664f66456666696369656e63793a3a73657450656e64696e6753746160408201527f746554696d656f75743a206578636565642068616c74206167677265676174696060820152691bdb881d1a5b595bdd5d60b21b608082015260a00190565b60208082526042908201527f456d657267656e63794d616e616765723a3a69664e6f74456d657267656e637960408201527f53746174653a206f6e6c79206966206e6f7420656d657267656e637920737461606082015261746560f01b608082015260a00190565b6020808252604190820152600080516020614e1583398151915260408201527f3a2063616e6e6f74207665726966792074686174206d616e79206261746368656060820152607360f81b608082015260a00190565b634e487b7160e01b600052603260045260246000fd5b60006001600160401b03808316818103614888576148886146a0565b6001019392505050565b928352602083019190915260c01b6001600160c01b031916604082015260480190565b6000602082840312156148c757600080fd5b5051919050565b9485526020850193909352604084019190915260c01b6001600160c01b0319166060808401919091521b6bffffffffffffffffffffffff19166068820152607c0190565b600060018201614924576149246146a0565b5060010190565b60006001600160401b038381169083168181101561494b5761494b6146a0565b039392505050565b600082821015614965576149656146a0565b500390565b6000816000190483118215151615614984576149846146a0565b500290565b601f821115613d0e57600081815260208120601f850160051c810160208610156149b05750805b601f850160051c820191505b818110156149cf578281556001016149bc565b505050505050565b81516001600160401b038111156149f0576149f0614024565b614a04816149fe8454614666565b84614989565b602080601f831160018114614a395760008415614a215750858301515b600019600386901b1c1916600185901b1785556149cf565b600085815260208120601f198616915b82811015614a6857888601518255948401946001909101908401614a49565b5085821015614a865787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b634e487b7160e01b600052601260045260246000fd5b600082614abb57614abb614a96565b500490565b60208082526048908201527f50726f6f664f66456666696369656e63793a3a6973466f72636542617463684160408201527f6c6c6f7765643a206f6e6c7920696620666f72636520626174636820697320616060820152677661696c61626c6560c01b608082015260a00190565b60008219821115614b4157614b416146a0565b500190565b6020808252603b908201527f456d657267656e63794d616e616765723a3a6966456d657267656e637953746160408201527f74653a206f6e6c7920696620656d657267656e63792073746174650000000000606082015260800190565b60208082526041908201527f50726f6f664f66456666696369656e63793a3a6f6e6c7954727573746564416760408201527f6772656761746f723a206f6e6c7920747275737465642041676772656761746f6060820152603960f91b608082015260a00190565b8381526001600160a01b0383166020820152606060408201819052600090614c3490830184613f5a565b95945050505050565b60006001600160401b0380841680614c5757614c57614a96565b92169190910492915050565b60008251614c75818460208701613f2e565b9190910192915050565b600082614c8e57614c8e614a96565b500690565b61012081016040808784376000838201818152879190815b6002811015614ccb57848483379084018281529284019290600101614cab565b5050828760c0870137610100850181815286935091505b6001811015614d01578251825260209283019290910190600101614ce2565b50505095945050505050565b600060208284031215614d1f57600080fd5b8151613f99816144d1565b6020808252607090820152600080516020614dd583398151915260408201527f6d696e697374696350656e64696e6753746174653a2066696e616c4e6577426160608201527f746368206d75737420626520626967676572207468616e2063757272656e744c60808201526f0c2e6e8accae4d2ccd2cac884c2e8c6d60831b60a082015260c0019056fe50726f6f664f66456666696369656e63793a3a73657175656e6365426174636850726f6f664f66456666696369656e63793a3a70726f76654e6f6e446574657250726f6f664f66456666696369656e63793a3a73657175656e6365466f72636550726f6f664f66456666696369656e63793a3a76657269667942617463686573a2646970667358221220f125f483b2be59d5016f53c9eda62c838a3ee3605676c9d8df6b17fa319cfb0e64736f6c634300080f0033",
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

// MAXBATCHLENGTH is a free data retrieval call binding the contract method 0xb02286c0.
//
// Solidity: function MAX_BATCH_LENGTH() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCaller) MAXBATCHLENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "MAX_BATCH_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXBATCHLENGTH is a free data retrieval call binding the contract method 0xb02286c0.
//
// Solidity: function MAX_BATCH_LENGTH() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencySession) MAXBATCHLENGTH() (*big.Int, error) {
	return _Proofofefficiency.Contract.MAXBATCHLENGTH(&_Proofofefficiency.CallOpts)
}

// MAXBATCHLENGTH is a free data retrieval call binding the contract method 0xb02286c0.
//
// Solidity: function MAX_BATCH_LENGTH() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCallerSession) MAXBATCHLENGTH() (*big.Int, error) {
	return _Proofofefficiency.Contract.MAXBATCHLENGTH(&_Proofofefficiency.CallOpts)
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

// CalculateBatchFee is a free data retrieval call binding the contract method 0x7abaf3e6.
//
// Solidity: function calculateBatchFee() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCaller) CalculateBatchFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Proofofefficiency.contract.Call(opts, &out, "calculateBatchFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateBatchFee is a free data retrieval call binding the contract method 0x7abaf3e6.
//
// Solidity: function calculateBatchFee() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencySession) CalculateBatchFee() (*big.Int, error) {
	return _Proofofefficiency.Contract.CalculateBatchFee(&_Proofofefficiency.CallOpts)
}

// CalculateBatchFee is a free data retrieval call binding the contract method 0x7abaf3e6.
//
// Solidity: function calculateBatchFee() view returns(uint256)
func (_Proofofefficiency *ProofofefficiencyCallerSession) CalculateBatchFee() (*big.Int, error) {
	return _Proofofefficiency.Contract.CalculateBatchFee(&_Proofofefficiency.CallOpts)
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
	NumBatch   uint64
	StateRoot  [32]byte
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterConsolidatePendingState is a free log retrieval operation binding the contract event 0x01f7d32e3b3278bace940a581067c87090c1aa09809730dd4ca002320c3a3cfa.
//
// Solidity: event ConsolidatePendingState(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) FilterConsolidatePendingState(opts *bind.FilterOpts, numBatch []uint64, aggregator []common.Address) (*ProofofefficiencyConsolidatePendingStateIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Proofofefficiency.contract.FilterLogs(opts, "ConsolidatePendingState", numBatchRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &ProofofefficiencyConsolidatePendingStateIterator{contract: _Proofofefficiency.contract, event: "ConsolidatePendingState", logs: logs, sub: sub}, nil
}

// WatchConsolidatePendingState is a free log subscription operation binding the contract event 0x01f7d32e3b3278bace940a581067c87090c1aa09809730dd4ca002320c3a3cfa.
//
// Solidity: event ConsolidatePendingState(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Proofofefficiency *ProofofefficiencyFilterer) WatchConsolidatePendingState(opts *bind.WatchOpts, sink chan<- *ProofofefficiencyConsolidatePendingState, numBatch []uint64, aggregator []common.Address) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Proofofefficiency.contract.WatchLogs(opts, "ConsolidatePendingState", numBatchRule, aggregatorRule)
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

// ParseConsolidatePendingState is a log parse operation binding the contract event 0x01f7d32e3b3278bace940a581067c87090c1aa09809730dd4ca002320c3a3cfa.
//
// Solidity: event ConsolidatePendingState(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
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
