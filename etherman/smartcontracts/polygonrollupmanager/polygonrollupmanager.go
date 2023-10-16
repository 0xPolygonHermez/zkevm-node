// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package polygonrollupmanager

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

// LegacyZKEVMStateVariablesPendingState is an auto generated low-level Go binding around an user-defined struct.
type LegacyZKEVMStateVariablesPendingState struct {
	Timestamp         uint64
	LastVerifiedBatch uint64
	ExitRoot          [32]byte
	StateRoot         [32]byte
}

// LegacyZKEVMStateVariablesSequencedBatchData is an auto generated low-level Go binding around an user-defined struct.
type LegacyZKEVMStateVariablesSequencedBatchData struct {
	AccInputHash               [32]byte
	SequencedTimestamp         uint64
	PreviousLastBatchSequenced uint64
}

// PolygonrollupmanagerMetaData contains all meta data concerning the Polygonrollupmanager contract.
var PolygonrollupmanagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRoot\",\"name\":\"_globalExitRootManager\",\"type\":\"address\"},{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"_pol\",\"type\":\"address\"},{\"internalType\":\"contractIPolygonZkEVMBridge\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlOnlyCanRenounceRolesForSelf\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressDoNotHaveRequiredRole\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllzkEVMSequencedBatchesMustBeVerified\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BatchAlreadyVerified\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BatchFeeOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BatchNotSequencedOrNotSequenceEnd\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ChainIDAlreadyExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConsensusAlreadyExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ConsensusDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExceedMaxVerifyBatches\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumBatchBelowLastVerifiedBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumBatchDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalPendingStateNumInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBatchNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBatchesAlreadyActive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HaltTimeoutNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitBatchMustMatchCurrentForkID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitNumBatchAboveLastVerifiedBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitNumBatchDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeBatchTimeTarget\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeForceBatchTimeout\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeMultiplierBatchFee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSequenceSomeBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewAccInputHashDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewPendingStateTimeoutMustBeLower\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewSequencedBatchMustBeBigger\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewStateRootNotInsidePrime\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewTrustedAggregatorTimeoutMustBeLower\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotSupportedCurrently\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OldAccInputHashDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OldStateRootDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyEmergencyState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyNotEmergencyState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPendingGovernance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyTrustedAggregator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateNotConsolidable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateTimeoutExceedHaltAggregationTimeout\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupMustExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupTypeDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupTypeObsolete\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderMustBeRollup\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StoredRootMustBeDifferentThanNewRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TrustedAggregatorTimeoutExceedHaltAggregationTimeout\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TrustedAggregatorTimeoutNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpdateNotCompatible\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpdateToSameRollupTypeID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VerifierAlreadyExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VerifierDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VerifierMustBeDifferent\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"VerifiersMustHaveSameGenesis\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rollupAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"}],\"name\":\"AddExistingRollup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consensusImplementation\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"name\":\"AddNewRollupType\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"ConsolidatePendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rollupAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"gasTokenAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasTokenNetwork\",\"type\":\"uint32\"}],\"name\":\"CreateNewRollup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EmergencyStateActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EmergencyStateDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"}],\"name\":\"ObsoleteRollupType\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"lastBatchSequenced\",\"type\":\"uint64\"}],\"name\":\"OnSequenceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"OverridePendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"storedStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"provedStateRoot\",\"type\":\"bytes32\"}],\"name\":\"ProveNonDeterministicPendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBatchFee\",\"type\":\"uint256\"}],\"name\":\"SetBatchFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"newMultiplierBatchFee\",\"type\":\"uint16\"}],\"name\":\"SetMultiplierBatchFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newPendingStateTimeout\",\"type\":\"uint64\"}],\"name\":\"SetPendingStateTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTrustedAggregator\",\"type\":\"address\"}],\"name\":\"SetTrustedAggregator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newTrustedAggregatorTimeout\",\"type\":\"uint64\"}],\"name\":\"SetTrustedAggregatorTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newVerifyBatchTimeTarget\",\"type\":\"uint64\"}],\"name\":\"SetVerifyBatchTimeTarget\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"newRollupTypeID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatchBeforeUpgrade\",\"type\":\"uint64\"}],\"name\":\"UpdateRollup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatchesTrustedAggregator\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activateEmergencyState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPolygonRollupBase\",\"name\":\"rollupAddress\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"}],\"name\":\"addExistingRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusImplementation\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"name\":\"addNewRollupType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMBridge\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateRewardPerBatch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"chainIDToRollupID\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"consolidatePendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"gasTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"gasTokenNetwork\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"networkName\",\"type\":\"string\"}],\"name\":\"createNewRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deactivateEmergencyState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBatchFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getForcedBatchFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"getInputSnarkBytes\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"name\":\"getLastVerifiedBatch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"batchNum\",\"type\":\"uint64\"}],\"name\":\"getRollupBatchNumToStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"batchNum\",\"type\":\"uint64\"}],\"name\":\"getRollupPendingStateTransitions\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structLegacyZKEVMStateVariables.PendingState\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"batchNum\",\"type\":\"uint64\"}],\"name\":\"getRollupSequencedBatches\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accInputHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequencedTimestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"previousLastBatchSequenced\",\"type\":\"uint64\"}],\"internalType\":\"structLegacyZKEVMStateVariables.SequencedBatchData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"globalExitRootManager\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRoot\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"trustedAggregator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_pendingStateTimeout\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_trustedAggregatorTimeout\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"timelock\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"emergencyCouncil\",\"type\":\"address\"},{\"internalType\":\"contractPolygonZkEVMV2Existent\",\"name\":\"polygonZkEVM\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"zkEVMVerifier\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isEmergencyState\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"isPendingStateConsolidable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastAggregationTimestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiplierBatchFee\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"}],\"name\":\"obsoleteRollupType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newSequencedBatches\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newAccInputHash\",\"type\":\"bytes32\"}],\"name\":\"onSequenceBatches\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"initPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"overridePendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingStateTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pol\",\"outputs\":[{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"initPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"proveNonDeterministicPendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rollupAddress\",\"type\":\"address\"}],\"name\":\"rollupAddressToID\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"name\":\"rollupIDToRollupData\",\"outputs\":[{\"internalType\":\"contractIPolygonRollupBase\",\"name\":\"rollupContract\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"lastLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"lastBatchSequenced\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastPendingState\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastPendingStateConsolidated\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatchBeforeUpgrade\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"rollupTypeID\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupTypeCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"rollupTypeID\",\"type\":\"uint64\"}],\"name\":\"rollupTypeMap\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"consensusImplementation\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"obsolete\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newBatchFee\",\"type\":\"uint256\"}],\"name\":\"setBatchFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"newMultiplierBatchFee\",\"type\":\"uint16\"}],\"name\":\"setMultiplierBatchFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newPendingStateTimeout\",\"type\":\"uint64\"}],\"name\":\"setPendingStateTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newTrustedAggregatorTimeout\",\"type\":\"uint64\"}],\"name\":\"setTrustedAggregatorTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newVerifyBatchTimeTarget\",\"type\":\"uint64\"}],\"name\":\"setVerifyBatchTimeTarget\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSequencedBatches\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalVerifiedBatches\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedAggregatorTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractITransparentUpgradeableProxy\",\"name\":\"rollupContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"newRollupTypeID\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"upgradeData\",\"type\":\"bytes\"}],\"name\":\"updateRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifyBatchTimeTarget\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"verifyBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"verifyBatchesTrustedAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b5060405162005d4638038062005d4683398101604081905262000034916200006b565b6001600160a01b0392831660805290821660c0521660a052620000bf565b6001600160a01b03811681146200006857600080fd5b50565b6000806000606084860312156200008157600080fd5b83516200008e8162000052565b6020850151909350620000a18162000052565b6040850151909250620000b48162000052565b809150509250925092565b60805160a05160c051615c1f6200012760003960008181610a0f01528181611c93015261371a0152600081816107010152818161281b01526139d801526000818161095201528181610c6601528181610e16015281816119d701526138c70152615c1f6000f3fe60806040523480156200001157600080fd5b5060043610620003095760003560e01c806391d14854116200019d578063ceee281d11620000eb578063dde0ff77116200009e578063dde0ff7714620009c0578063e0bfd3d214620009db578063e136edf714620009f2578063e46761c41462000a09578063f34eb8eb1462000a31578063f4e926751462000a48578063f9c4c2ae1462000a5957600080fd5b8063ceee281d1462000923578063d02103ca146200094c578063d5073f6f1462000974578063d547741f146200098b578063d939b31514620009a2578063dbc1697614620009b657600080fd5b8063a2967d991162000150578063a2967d9914620006f1578063a3c573eb14620006fb578063a4810bb2146200073c578063afd23cbe14620007eb578063b99d0ad71462000815578063c1acbc3414620008f1578063c4c928c2146200090c57600080fd5b806391d14854146200068257806399f5634e14620006995780639a908e7314620006a35780639c9f3dfe14620006ba578063a066215c14620006d1578063a217fddf14620006e857600080fd5b80632f2ff15d116200025b57806360469169116200020e5780636046916914620005c95780637222020f14620005d35780637975fcfe14620005ea5780637fb6e76a1462000610578063841b24d7146200063957806387c20c0114620006545780638bd4f071146200066b57600080fd5b80632f2ff15d146200052057806336568abe1462000537578063394218e9146200054e578063477fa27014620005655780634d4d3338146200056e57806355a71ee0146200058557600080fd5b806315064c9611620002c057806315064c9614620003c95780631608859c14620003d75780631796a1ae14620003ee5780631816b7e514620004155780632072f6c5146200042c578063248a9ca3146200043657806325280169146200046b57600080fd5b8063066ec012146200030e578063080b3111146200033f5780630a0d9fbe146200036757806311f6b287146200038257806312b86e1914620003995780631489ed1014620003b2575b600080fd5b60845462000322906001600160401b031681565b6040516001600160401b0390911681526020015b60405180910390f35b62000356620003503660046200467c565b62000b70565b604051901515815260200162000336565b6085546200032290600160401b90046001600160401b031681565b6200032262000393366004620046b4565b62000b9a565b620003b0620003aa366004620046e5565b62000bba565b005b620003b0620003c336600462004792565b62000d6a565b606f54620003569060ff1681565b620003b0620003e83660046200467c565b62000f00565b607e54620003ff9063ffffffff1681565b60405163ffffffff909116815260200162000336565b620003b0620004263660046200481c565b62000f95565b620003b062001041565b6200045c6200044736600462004849565b60009081526034602052604090206001015490565b60405190815260200162000336565b620004ec6200047c3660046200467c565b60408051606080820183526000808352602080840182905292840181905263ffffffff959095168552608182528285206001600160401b03948516865260030182529382902082519485018352805485526001015480841691850191909152600160401b90049091169082015290565b60408051825181526020808401516001600160401b0390811691830191909152928201519092169082015260600162000336565b620003b06200053136600462004863565b620010bb565b620003b06200054836600462004863565b620010e4565b620003b06200055f36600462004896565b6200111e565b6086546200045c565b620003b06200057f366004620048b4565b620011d2565b6200045c620005963660046200467c565b63ffffffff821660009081526081602090815260408083206001600160401b038516845260020190915290205492915050565b6200045c6200171f565b620003b0620005e4366004620046b4565b62001737565b62000601620005fb3660046200496a565b62001833565b60405162000336919062004a24565b620003ff6200062136600462004896565b60836020526000908152604090205463ffffffff1681565b6084546200032290600160c01b90046001600160401b031681565b620003b06200066536600462004792565b62001866565b620003b06200067c366004620046e5565b62001b90565b620003566200069336600462004863565b62001c46565b6200045c62001c71565b62000322620006b436600462004a39565b62001d58565b620003b0620006cb36600462004896565b62001f2a565b620003b0620006e236600462004896565b62001fd0565b6200045c600081565b6200045c62002074565b620007237f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200162000336565b620007a16200074d36600462004896565b607f602052600090815260409020805460018201546002909201546001600160a01b0391821692918216916001600160401b03600160a01b8204169160ff600160e01b8304811692600160e81b9004169086565b604080516001600160a01b0397881681529690951660208701526001600160401b039093169385019390935260ff166060840152901515608083015260a082015260c00162000336565b6085546200080190600160801b900461ffff1681565b60405161ffff909116815260200162000336565b620008ab620008263660046200467c565b604080516080808201835260008083526020808401829052838501829052606093840182905263ffffffff969096168152608186528381206001600160401b03958616825260040186528390208351918201845280548086168352600160401b9004909416948101949094526001830154918401919091526002909101549082015290565b60405162000336919060006080820190506001600160401b0380845116835280602085015116602084015250604083015160408301526060830151606083015292915050565b6084546200032290600160801b90046001600160401b031681565b620003b06200091d36600462004a66565b62002449565b620003ff6200093436600462004afe565b60826020526000908152604090205463ffffffff1681565b620007237f000000000000000000000000000000000000000000000000000000000000000081565b620003b06200098536600462004849565b62002728565b620003b06200099c36600462004863565b620027c4565b60855462000322906001600160401b031681565b620003b0620027ed565b6084546200032290600160401b90046001600160401b031681565b620003b0620009ec36600462004b30565b6200289b565b620003b062000a0336600462004c56565b62002934565b620007237f000000000000000000000000000000000000000000000000000000000000000081565b620003b062000a4236600462004d36565b62002db9565b608054620003ff9063ffffffff1681565b62000af062000a6a366004620046b4565b608160205260009081526040902080546001820154600583015460068401546007909401546001600160a01b0380851695600160a01b958690046001600160401b039081169692861695929092048216939282821692600160401b808404821693600160801b808204841694600160c01b90920484169380831693830416910460ff168c565b604080516001600160a01b039d8e1681526001600160401b039c8d1660208201529c909a16998c019990995296891660608b015260808a019590955292871660a089015290861660c0880152851660e0870152841661010086015283166101208501529190911661014083015260ff166101608201526101800162000336565b63ffffffff8216600090815260816020526040812062000b91908362002fb8565b90505b92915050565b63ffffffff8116600090815260816020526040812062000b949062002ffd565b60008051602062005bca83398151915262000bd5816200306e565b63ffffffff8916600090815260816020526040902062000bfc818a8a8a8a8a8a8a6200307a565b600681018054600160401b600160801b031916600160401b6001600160401b0389811691820292909217835560009081526002840160205260409020869055600583018790559054600160801b9004161562000c64576006810180546001600160801b031690555b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166333d6247d62000c9d62002074565b6040518263ffffffff1660e01b815260040162000cbc91815260200190565b600060405180830381600087803b15801562000cd757600080fd5b505af115801562000cec573d6000803e3d6000fd5b5050608480546001600160c01b031661127560c71b1790555050604080516001600160401b03881681526020810186905290810186905233606082015263ffffffff8b16907f3182bd6e6f74fc1fdc88b60f3a4f4c7f79db6ae6f5b88a1b3f5a1e28ec210d5e9060800160405180910390a250505050505050505050565b60008051602062005bca83398151915262000d85816200306e565b63ffffffff8916600090815260816020526040902062000dac818a8a8a8a8a8a8a62003413565b600681018054600160401b600160801b031916600160401b6001600160401b038a811691820292909217835560009081526002840160205260409020879055600583018890559054600160801b9004161562000e14576006810180546001600160801b031690555b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166333d6247d62000e4d62002074565b6040518263ffffffff1660e01b815260040162000e6c91815260200190565b600060405180830381600087803b15801562000e8757600080fd5b505af115801562000e9c573d6000803e3d6000fd5b5050604080516001600160401b038b1681526020810189905290810189905233925063ffffffff8d1691507fd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d39060600160405180910390a350505050505050505050565b63ffffffff8216600090815260816020526040902062000f3060008051602062005bca8339815191523362001c46565b62000f8457606f5460ff161562000f5a57604051630bc011ff60e21b815260040160405180910390fd5b62000f66818362002fb8565b62000f8457604051630674f25160e11b815260040160405180910390fd5b62000f908183620037e4565b505050565b60008051602062005b8a83398151915262000fb0816200306e565b6103e88261ffff16108062000fca57506103ff8261ffff16115b1562000fe957604051630984a67960e31b815260040160405180910390fd5b6085805461ffff60801b1916600160801b61ffff8516908102919091179091556040519081527f7019933d795eba185c180209e8ae8bffbaa25bcef293364687702c31f4d302c5906020015b60405180910390a15050565b6200105c60008051602062005baa8339815191523362001c46565b620010af576084544290620010869062093a8090600160801b90046001600160401b031662004de3565b6001600160401b03161115620010af5760405163692baaad60e11b815260040160405180910390fd5b620010b9620039d6565b565b600082815260346020526040902060010154620010d8816200306e565b62000f90838362003a55565b6001600160a01b03811633146200110e57604051630b4ad1cd60e31b815260040160405180910390fd5b6200111a828262003ac1565b5050565b60008051602062005b8a83398151915262001139816200306e565b606f5460ff166200117b576084546001600160401b03600160c01b9091048116908316106200117b5760405163401636df60e01b815260040160405180910390fd5b608480546001600160c01b0316600160c01b6001600160401b038516908102919091179091556040519081527f1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a19060200162001035565b600054600290610100900460ff16158015620011f5575060005460ff8083169116105b6200125e5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084015b60405180910390fd5b6000805461010060ff841661ffff199092169190911717905560858054608480546001600160c01b0316600160c01b6001600160401b038c8116919091029190911790915567016345785d8a00006086558a166001600160801b03199091161760e160431b1761ffff60801b19166101f560811b179055620012df62003b2b565b620012fa60008051602062005bca8339815191528a62003b98565b6200130760008662003b98565b620013337fac75d24dbb35ea80e25fab167da4dea46c1915260426570db84f184891f5f5908662003b98565b6200135f7f3dfe277d2a2c04b75fb2eb3743fa00005ae3678a20c299e65fdf4df76517f68e8662003b98565b6200138b7f66156603fe29d13f97c6f3e3dff4ef71919f9aa61c555be0182d954e94221aac8662003b98565b620013b77fab66e11c4f712cd06ab11bf9339b48bef39e12d4a22eeef71d2860a0c90482bd8762003b98565b620013e37fa0fab074aba36a6fa69f1a83ee86e5abfb8433966eb57efb13dc2fc2f24ddd088762003b98565b6200140f7f62ba6ba2ffed8cfe316b583325ea41ac6e7ba9e5864d2bc6fabba7ac26d2f0f48762003b98565b6200142a60008051602062005b8a8339815191528762003b98565b6200146560008051602062005bca8339815191527f73cb0569fdbea2544dae03fdb2fe10eda92a72a2e8cd2bd496e85b762505a3f062003ba4565b620014917f73cb0569fdbea2544dae03fdb2fe10eda92a72a2e8cd2bd496e85b762505a3f08762003b98565b620014bd7f8cf807f6970720f8e2c208c7c5037595982c7bd9ed93c380d09df743d0dcc3fb8762003b98565b620014f860008051602062005baa8339815191527f9b6f082d8d3644ae2f24a3c32e356d6f2d9b2844d9b26164fbc82663ff28595162003ba4565b6200151360008051602062005baa8339815191528562003b98565b6200153f7f9b6f082d8d3644ae2f24a3c32e356d6f2d9b2844d9b26164fbc82663ff2859518562003b98565b6000620015548484600561044d600062003bef565b6073546074549192506001600160401b03600160401b909104811691168082146200159257604051632e4cc54360e11b815260040160405180910390fd5b6001600160401b0381811660008181526075602090815260408083205460028901835281842055868516808452607280845282852060038b018552948390208554815560018087018054919092018054918a1667ffffffffffffffff198084168217835593546001600160801b0319938416909117600160401b918290048c1682021790915560068d0180549092169094179388029390931790925560078a0180549092169095179055607a54606f54949092529154607354925163176b20e160e31b81526001600160a01b038c81169663bb59070896620016949695831695600160581b909104909216936076936077939192919091169060040162004eb9565b600060405180830381600087803b158015620016af57600080fd5b505af1158015620016c4573d6000803e3d6000fd5b50506000805461ff0019169055505060405160ff851681527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498935060200191506200170c9050565b60405180910390a1505050505050505050565b6000608654606462001732919062004f1e565b905090565b7fab66e11c4f712cd06ab11bf9339b48bef39e12d4a22eeef71d2860a0c90482bd62001763816200306e565b63ffffffff82161580620017825750607e5463ffffffff908116908316115b15620017a157604051637512e5cb60e01b815260040160405180910390fd5b63ffffffff82166000908152607f60205260409020600180820154600160e81b900460ff1615159003620017e857604051633b8d3d9960e01b815260040160405180910390fd5b60018101805460ff60e81b1916600160e81b17905560405163ffffffff8416907f4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e4490600090a2505050565b63ffffffff861660009081526081602052604090206060906200185b90878787878762003e14565b979650505050505050565b606f5460ff16156200188b57604051630bc011ff60e21b815260040160405180910390fd5b63ffffffff881660009081526081602090815260408083206084546001600160401b038a811686526003830190945291909320600101544292620018da92600160c01b90048116911662004de3565b6001600160401b031611156200190357604051638a0704d360e01b815260040160405180910390fd5b6103e862001912888862004f38565b6001600160401b031611156200193b57604051635acfba9d60e11b815260040160405180910390fd5b6200194d818989898989898962003413565b62001959818762003f54565b6085546001600160401b031660000362001a6757600681018054600160401b600160801b031916600160401b6001600160401b0389811691820292909217835560009081526002840160205260409020869055600583018790559054600160801b90041615620019d5576006810180546001600160801b031690555b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166333d6247d62001a0e62002074565b6040518263ffffffff1660e01b815260040162001a2d91815260200190565b600060405180830381600087803b15801562001a4857600080fd5b505af115801562001a5d573d6000803e3d6000fd5b5050505062001b31565b62001a728162004151565b600681018054600160801b90046001600160401b031690601062001a968362004f5b565b82546001600160401b039182166101009390930a92830292820219169190911790915560408051608081018252428316815289831660208083019182528284018b8152606084018b81526006890154600160801b90048716600090815260048a01909352949091209251835492518616600160401b026001600160801b03199093169516949094171781559151600183015551600290910155505b604080516001600160401b038816815260208101869052908101869052339063ffffffff8b16907faac1e7a157b259544ebacd6e8a82ae5d6c8f174e12aa48696277bcc9a661f0b49060600160405180910390a3505050505050505050565b606f5460ff161562001bb557604051630bc011ff60e21b815260040160405180910390fd5b63ffffffff8816600090815260816020526040902062001bdc81898989898989896200307a565b6001600160401b03851660009081526002820160209081526040918290205482519081529081018590527f1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010910160405180910390a162001c3b620039d6565b505050505050505050565b60009182526034602090815260408084206001600160a01b0393909316845291905290205460ff1690565b6040516370a0823160e01b815230600482015260009081906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa15801562001cdb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001d01919062004f84565b60845490915060009062001d28906001600160401b03600160401b82048116911662004f38565b6001600160401b031690508060000362001d455760009250505090565b62001d51818362004fb4565b9250505090565b606f5460009060ff161562001d8057604051630bc011ff60e21b815260040160405180910390fd5b3360009081526082602052604081205463ffffffff169081900362001db8576040516371653c1560e01b815260040160405180910390fd5b836001600160401b031660000362001de357604051632590ccf960e01b815260040160405180910390fd5b63ffffffff81166000908152608160205260408120608480549192879262001e169084906001600160401b031662004de3565b82546101009290920a6001600160401b038181021990931691831602179091556006830154169050600062001e4c878362004de3565b6006840180546001600160401b0380841667ffffffffffffffff199092168217909255604080516060810182528a81524284166020808301918252888616838501908152600095865260038b0190915292909320905181559151600192909201805491518416600160401b026001600160801b03199092169290931691909117179055905062001edc8362004151565b6040516001600160401b038216815263ffffffff8516907f1d9f30260051d51d70339da239ea7b080021adcaabfa71c9b0ea339a20cf9a259060200160405180910390a29695505050505050565b60008051602062005b8a83398151915262001f45816200306e565b606f5460ff1662001f80576085546001600160401b039081169083161062001f805760405163048a05a960e41b815260040160405180910390fd5b6085805467ffffffffffffffff19166001600160401b0384169081179091556040519081527fc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c759060200162001035565b60008051602062005b8a83398151915262001feb816200306e565b62015180826001600160401b031611156200201957604051631c0cfbfd60e31b815260040160405180910390fd5b60858054600160401b600160801b031916600160401b6001600160401b038516908102919091179091556040519081527f1b023231a1ab6b5d93992f168fb44498e1a7e64cef58daff6f1c216de6a68c289060200162001035565b60805460009063ffffffff168082036200209057506000919050565b6000816001600160401b03811115620020ad57620020ad62004bac565b604051908082528060200260200182016040528015620020d7578160200160208202803683370190505b50905060005b828110156200214a5760816000620020f783600162004fcb565b63ffffffff1663ffffffff168152602001908152602001600020600501548282815181106200212a576200212a62004fe1565b602090810291909101015280620021418162004ff7565b915050620020dd565b50600060205b83600114620023a65760006200216860028662005013565b6200217560028762004fb4565b62002181919062004fcb565b90506000816001600160401b03811115620021a057620021a062004bac565b604051908082528060200260200182016040528015620021ca578160200160208202803683370190505b50905060005b828110156200235257620021e66001846200502a565b81148015620022015750620021fd60028862005013565b6001145b156200228957856200221582600262004f1e565b8151811062002228576200222862004fe1565b6020026020010151856040516020016200224c929190918252602082015260400190565b6040516020818303038152906040528051906020012082828151811062002277576200227762004fe1565b6020026020010181815250506200233d565b856200229782600262004f1e565b81518110620022aa57620022aa62004fe1565b602002602001015186826002620022c2919062004f1e565b620022cf90600162004fcb565b81518110620022e257620022e262004fe1565b602002602001015160405160200162002305929190918252602082015260400190565b6040516020818303038152906040528051906020012082828151811062002330576200233062004fe1565b6020026020010181815250505b80620023498162004ff7565b915050620021d0565b50809450819550838460405160200162002376929190918252602082015260400190565b60405160208183030381529060405280519060200120935082806200239b9062005040565b935050505062002150565b600083600081518110620023be57620023be62004fe1565b6020026020010151905060005b828110156200243f57604080516020810184905290810185905260600160408051601f19818403018152828252805160209182012090830187905290820186905292506060016040516020818303038152906040528051906020012093508080620024369062004ff7565b915050620023cb565b5095945050505050565b7f66156603fe29d13f97c6f3e3dff4ef71919f9aa61c555be0182d954e94221aac62002475816200306e565b63ffffffff84161580620024945750607e5463ffffffff908116908516115b15620024b357604051637512e5cb60e01b815260040160405180910390fd5b6001600160a01b03851660009081526082602052604081205463ffffffff1690819003620024f4576040516374a086a360e01b815260040160405180910390fd5b63ffffffff8181166000908152608160205260409020600781015490918716600160401b9091046001600160401b0316036200254357604051634f61d51960e01b815260040160405180910390fd5b63ffffffff86166000908152607f60205260409020600180820154600160e81b900460ff16151590036200258a57604051633b8d3d9960e01b815260040160405180910390fd5b60018101546007830154600160801b900460ff908116600160e01b9092041614620025c857604051635aa0d5f160e11b815260040160405180910390fd5b6001808201805491840180546001600160a01b031981166001600160a01b03909416938417825591546001600160401b03600160a01b9182900416026001600160e01b0319909216909217179055600782018054600160401b63ffffffff8a1602600160401b600160801b03199091161790556000620026488462000b9a565b60078401805467ffffffffffffffff19166001600160401b038316179055825460405163278f794360e11b81529192506001600160a01b038b811692634f1ef286926200269e9216908b908b906004016200505a565b600060405180830381600087803b158015620026b957600080fd5b505af1158015620026ce573d6000803e3d6000fd5b50506040805163ffffffff8c811682526001600160401b0386166020830152881693507ff585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d92500160405180910390a2505050505050505050565b7f8cf807f6970720f8e2c208c7c5037595982c7bd9ed93c380d09df743d0dcc3fb62002754816200306e565b683635c9adc5dea000008211806200276f5750633b9aca0082105b156200278e57604051638586952560e01b815260040160405180910390fd5b60868290556040518281527ffb383653f53ee079978d0c9aff7aeff04a10166ce244cca9c9f9d8d96bed45b29060200162001035565b600082815260346020526040902060010154620027e1816200306e565b62000f90838362003ac1565b7f62ba6ba2ffed8cfe316b583325ea41ac6e7ba9e5864d2bc6fabba7ac26d2f0f462002819816200306e565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663dbc169766040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156200287557600080fd5b505af11580156200288a573d6000803e3d6000fd5b50505050620028986200421c565b50565b7f3dfe277d2a2c04b75fb2eb3743fa00005ae3678a20c299e65fdf4df76517f68e620028c7816200306e565b6001600160401b03841660009081526083602052604090205463ffffffff161562002905576040516337c8fe0960e11b815260040160405180910390fd5b600062002916888888888762003bef565b60008080526002909101602052604090209390935550505050505050565b7fa0fab074aba36a6fa69f1a83ee86e5abfb8433966eb57efb13dc2fc2f24ddd0862002960816200306e565b63ffffffff891615806200297f5750607e5463ffffffff908116908a16115b156200299e57604051637512e5cb60e01b815260040160405180910390fd5b63ffffffff89166000908152607f60205260409020600180820154600160e81b900460ff1615159003620029e557604051633b8d3d9960e01b815260040160405180910390fd5b6001600160401b03891660009081526083602052604090205463ffffffff161562002a23576040516337c8fe0960e11b815260040160405180910390fd5b6080805460009190829062002a3e9063ffffffff166200509a565b825463ffffffff8281166101009490940a9384029302191691909117909155825460408051600080825260208201928390529394506001600160a01b0390921691309162002a8c9062004641565b62002a9a93929190620050b6565b604051809103906000f08015801562002ab7573d6000803e3d6000fd5b50905081608360008d6001600160401b03166001600160401b0316815260200190815260200160002060006101000a81548163ffffffff021916908363ffffffff1602179055508160826000836001600160a01b03166001600160a01b0316815260200190815260200160002060006101000a81548163ffffffff021916908363ffffffff1602179055506000608160008463ffffffff1663ffffffff1681526020019081526020016000209050818160000160006101000a8154816001600160a01b0302191690836001600160a01b031602179055508360010160149054906101000a90046001600160401b03168160010160146101000a8154816001600160401b0302191690836001600160401b031602179055508360010160009054906101000a90046001600160a01b03168160010160006101000a8154816001600160a01b0302191690836001600160a01b031602179055508b8160000160146101000a8154816001600160401b0302191690836001600160401b031602179055508360020154816002016000806001600160401b03168152602001908152602001600020819055508c63ffffffff168160070160086101000a8154816001600160401b0302191690836001600160401b0316021790555083600101601c9054906101000a900460ff168160070160106101000a81548160ff021916908360ff1602179055508263ffffffff167f9a84992ead4481f37097c7663e97ed883b31f216fae5685c6ee576a656483cd98e848f8d8d60405162002d3495949392919063ffffffff95861681526001600160a01b0394851660208201526001600160401b03939093166040840152921660608201529116608082015260a00190565b60405180910390a26040516308fd7e2160e31b81526001600160a01b038316906347ebf1089062002d76908e908e9088908f908f908f908f90600401620050ed565b600060405180830381600087803b15801562002d9157600080fd5b505af115801562002da6573d6000803e3d6000fd5b5050505050505050505050505050505050565b7fac75d24dbb35ea80e25fab167da4dea46c1915260426570db84f184891f5f59062002de5816200306e565b607e805460009190829062002e009063ffffffff166200509a565b91906101000a81548163ffffffff021916908363ffffffff160217905590506040518060c00160405280896001600160a01b03168152602001886001600160a01b03168152602001876001600160401b031681526020018660ff16815260200160001515815260200185815250607f60008363ffffffff166001600160401b0316815260200190815260200160002060008201518160000160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060208201518160010160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060408201518160010160146101000a8154816001600160401b0302191690836001600160401b03160217905550606082015181600101601c6101000a81548160ff021916908360ff160217905550608082015181600101601d6101000a81548160ff02191690831515021790555060a082015181600201559050508063ffffffff167fa2970448b3bd66ba7e524e7b2a5b9cf94fa29e32488fb942afdfe70dd4b77b5289898989898960405162002fa6969594939291906200515a565b60405180910390a25050505050505050565b6085546001600160401b0382811660009081526004850160205260408120549092429262002feb92918116911662004de3565b6001600160401b031611159392505050565b6006810154600090600160801b90046001600160401b03161562003051575060068101546001600160401b03600160801b909104811660009081526004909201602052604090912054600160401b90041690565b5060060154600160401b90046001600160401b031690565b919050565b62002898813362004275565b60078801546000906001600160401b039081169087161015620030b05760405163ead1340b60e01b815260040160405180910390fd5b6001600160401b03881615620031515760068901546001600160401b03600160801b90910481169089161115620030fa5760405163bb14c20560e01b815260040160405180910390fd5b506001600160401b03808816600090815260048a0160205260409020600281015481549092888116600160401b90920416146200314a57604051632bd2e3e760e01b815260040160405180910390fd5b50620031c6565b506001600160401b0385166000908152600289016020526040902054806200318c576040516324cbdcc360e11b815260040160405180910390fd5b60068901546001600160401b03600160401b90910481169087161115620031c657604051630f2b74f160e11b815260040160405180910390fd5b60068901546001600160401b03600160801b90910481169088161180620031ff5750876001600160401b0316876001600160401b031611155b8062003223575060068901546001600160401b03600160c01b909104811690881611155b15620032425760405163bfa7079f60e01b815260040160405180910390fd5b6001600160401b03878116600090815260048b016020526040902054600160401b900481169086161462003289576040516332a2a77f60e01b815260040160405180910390fd5b60006200329b8a888888868962003e14565b905060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001600283604051620032d29190620051b3565b602060405180830381855afa158015620032f0573d6000803e3d6000fd5b5050506040513d601f19601f8201168201806040525081019062003315919062004f84565b62003321919062005013565b60018c0154604080516020810182528381529051634890ed4560e11b81529293506001600160a01b0390911691639121da8a916200336591889190600401620051d1565b602060405180830381865afa15801562003383573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620033a991906200520e565b620033c7576040516309bde33960e01b815260040160405180910390fd5b6001600160401b038916600090815260048c016020526040902060020154859003620034065760405163a47276bd60e01b815260040160405180910390fd5b5050505050505050505050565b600080620034218a62002ffd565b60078b01549091506001600160401b039081169089161015620034575760405163ead1340b60e01b815260040160405180910390fd5b6001600160401b03891615620034fa5760068a01546001600160401b03600160801b9091048116908a161115620034a15760405163bb14c20560e01b815260040160405180910390fd5b6001600160401b03808a16600090815260048c01602052604090206002810154815490945090918a8116600160401b9092041614620034f357604051632bd2e3e760e01b815260040160405180910390fd5b506200356a565b6001600160401b038816600090815260028b01602052604090205491508162003536576040516324cbdcc360e11b815260040160405180910390fd5b806001600160401b0316886001600160401b031611156200356a57604051630f2b74f160e11b815260040160405180910390fd5b806001600160401b0316876001600160401b0316116200359d5760405163b9b18f5760e01b815260040160405180910390fd5b6000620035af8b8a8a8a878b62003e14565b905060007f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001600283604051620035e69190620051b3565b602060405180830381855afa15801562003604573d6000803e3d6000fd5b5050506040513d601f19601f8201168201806040525081019062003629919062004f84565b62003635919062005013565b60018d0154604080516020810182528381529051634890ed4560e11b81529293506001600160a01b0390911691639121da8a916200367991899190600401620051d1565b602060405180830381865afa15801562003697573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620036bd91906200520e565b620036db576040516309bde33960e01b815260040160405180910390fd5b6000620036e9848b62004f38565b90506200374287826001600160401b03166200370462001c71565b62003710919062004f1e565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001691906200429f565b80608460088282829054906101000a90046001600160401b031662003768919062004de3565b82546101009290920a6001600160401b038181021990931691831602179091556084805467ffffffffffffffff60801b1916600160801b428416021790558e546040516332c2d15360e01b8152918d166004830152602482018b90523360448301526001600160a01b031691506332c2d1539060640162002d76565b60068201546001600160401b03600160c01b909104811690821611158062003823575060068201546001600160401b03600160801b9091048116908216115b15620038425760405163d086b70b60e01b815260040160405180910390fd5b6001600160401b03818116600081815260048501602090815260408083208054600689018054600160401b600160801b031916600160401b92839004909816918202979097178755600280830154828752908a0190945291909320919091556001820154600587015583546001600160c01b0316600160c01b909302929092179092557f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166333d6247d620038fe62002074565b6040518263ffffffff1660e01b81526004016200391d91815260200190565b600060405180830381600087803b1580156200393857600080fd5b505af11580156200394d573d6000803e3d6000fd5b505085546001600160a01b0316600090815260826020908152604091829020546002870154600188015484516001600160401b03898116825294810192909252818501529188166060830152915163ffffffff90921693507f581910eb7a27738945c2f00a91f2284b2d6de9d4e472b12f901c2b0df045e21b925081900360800190a250505050565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316632072f6c56040518163ffffffff1660e01b8152600401600060405180830381600087803b15801562003a3257600080fd5b505af115801562003a47573d6000803e3d6000fd5b50505050620010b9620042f3565b62003a61828262001c46565b6200111a5760008281526034602090815260408083206001600160a01b0385168085529252808320805460ff1916600117905551339285917f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9190a45050565b62003acd828262001c46565b156200111a5760008281526034602090815260408083206001600160a01b0385168085529252808320805460ff1916905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b600054610100900460ff16620010b95760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b606482015260840162001255565b6200111a828262003a55565b600082815260346020526040808220600101805490849055905190918391839186917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9190a4505050565b608080546000918291829062003c0b9063ffffffff166200509a565b91906101000a81548163ffffffff021916908363ffffffff160217905590508060836000866001600160401b03166001600160401b0316815260200190815260200160002060006101000a81548163ffffffff021916908363ffffffff1602179055508060826000896001600160a01b03166001600160a01b0316815260200190815260200160002060006101000a81548163ffffffff021916908363ffffffff160217905550608160008263ffffffff1663ffffffff1681526020019081526020016000209150868260000160006101000a8154816001600160a01b0302191690836001600160a01b03160217905550848260010160146101000a8154816001600160401b0302191690836001600160401b03160217905550858260010160006101000a8154816001600160a01b0302191690836001600160a01b03160217905550838260000160146101000a8154816001600160401b0302191690836001600160401b03160217905550828260070160106101000a81548160ff021916908360ff1602179055508063ffffffff167f1c20c90adb9c53d1da3e79fb8531f2c13c52b73a31a2574202ec36fb449ab24a8689878760405162003e0294939291906001600160401b0394851681526001600160a01b039390931660208401529216604082015260ff91909116606082015260800190565b60405180910390a25095945050505050565b6001600160401b038086166000818152600389016020526040808220549388168252902054606092911580159062003e4a575081155b1562003e695760405163340c614f60e11b815260040160405180910390fd5b8062003e88576040516366385b5160e01b815260040160405180910390fd5b62003e938462004350565b62003eb1576040516305dae44f60e21b815260040160405180910390fd5b885460018a01546040516bffffffffffffffffffffffff193360601b16602082015260348101889052605481018590526001600160c01b031960c08c811b82166074840152600160a01b94859004811b8216607c84015293909204831b82166084820152608c810187905260ac810184905260cc81018990529189901b1660ec82015260f401604051602081830303815290604052925050509695505050505050565b600062003f618362002ffd565b90508160008062003f73848462004f38565b6085546001600160401b03918216925060009162003f9a91600160401b900416426200502a565b90505b846001600160401b0316846001600160401b03161462004024576001600160401b0380851660009081526003890160205260409020600181015490911682101562003fff576001810154600160401b90046001600160401b031694506200401d565b6200400b868662004f38565b6001600160401b031693505062004024565b5062003f9d565b60006200403284846200502a565b9050838110156200409057808403600c811162004050578062004053565b600c5b9050806103e80a81608560109054906101000a900461ffff1661ffff160a608654028162004085576200408562004f9e565b046086555062004108565b838103600c8111620040a35780620040a6565b600c5b90506000816103e80a82608560109054906101000a900461ffff1661ffff160a670de0b6b3a76400000281620040e057620040e062004f9e565b04905080608654670de0b6b3a7640000028162004101576200410162004f9e565b0460865550505b683635c9adc5dea0000060865411156200412f57683635c9adc5dea0000060865562004147565b633b9aca0060865410156200414757633b9aca006086555b5050505050505050565b60068101546001600160401b03600160c01b82048116600160801b909204161115620028985760068101546000906200419c90600160c01b90046001600160401b0316600162004de3565b9050620041aa828262002fb8565b156200111a576006820154600090600290620041d8908490600160801b90046001600160401b031662004f38565b620041e4919062005232565b620041f0908362004de3565b9050620041fe838262002fb8565b15620042105762000f908382620037e4565b62000f908383620037e4565b606f5460ff166200424057604051635386698160e01b815260040160405180910390fd5b606f805460ff191690556040517f1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b390600090a1565b62004281828262001c46565b6200111a57604051637615be1f60e11b815260040160405180910390fd5b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b031663a9059cbb60e01b17905262000f90908490620043d6565b606f5460ff16156200431857604051630bc011ff60e21b815260040160405180910390fd5b606f805460ff191660011790556040517f2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a549790600090a1565b600067ffffffff000000016001600160401b03831610801562004387575067ffffffff00000001604083901c6001600160401b0316105b8015620043a8575067ffffffff00000001608083901c6001600160401b0316105b8015620043c0575067ffffffff0000000160c083901c105b15620043ce57506001919050565b506000919050565b60006200442d826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b0316620044af9092919063ffffffff16565b80519091501562000f9057808060200190518101906200444e91906200520e565b62000f905760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840162001255565b6060620044c08484600085620044c8565b949350505050565b6060824710156200452b5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840162001255565b600080866001600160a01b03168587604051620045499190620051b3565b60006040518083038185875af1925050503d806000811462004588576040519150601f19603f3d011682016040523d82523d6000602084013e6200458d565b606091505b50915091506200185b87838387606083156200460e57825160000362004606576001600160a01b0385163b620046065760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640162001255565b5081620044c0565b620044c08383815115620046255781518083602001fd5b8060405162461bcd60e51b815260040162001255919062004a24565b61092e806200525c83390190565b803563ffffffff811681146200306957600080fd5b80356001600160401b03811681146200306957600080fd5b600080604083850312156200469057600080fd5b6200469b836200464f565b9150620046ab6020840162004664565b90509250929050565b600060208284031215620046c757600080fd5b62000b91826200464f565b80610300810183101562000b9457600080fd5b6000806000806000806000806103e0898b0312156200470357600080fd5b6200470e896200464f565b97506200471e60208a0162004664565b96506200472e60408a0162004664565b95506200473e60608a0162004664565b94506200474e60808a0162004664565b935060a0890135925060c089013591506200476d8a60e08b01620046d2565b90509295985092959890939650565b6001600160a01b03811681146200289857600080fd5b6000806000806000806000806103e0898b031215620047b057600080fd5b620047bb896200464f565b9750620047cb60208a0162004664565b9650620047db60408a0162004664565b9550620047eb60608a0162004664565b94506080890135935060a0890135925060c08901356200480b816200477c565b91506200476d8a60e08b01620046d2565b6000602082840312156200482f57600080fd5b813561ffff811681146200484257600080fd5b9392505050565b6000602082840312156200485c57600080fd5b5035919050565b600080604083850312156200487757600080fd5b8235915060208301356200488b816200477c565b809150509250929050565b600060208284031215620048a957600080fd5b62000b918262004664565b600080600080600080600080610100898b031215620048d257600080fd5b8835620048df816200477c565b9750620048ef60208a0162004664565b9650620048ff60408a0162004664565b9550606089013562004911816200477c565b9450608089013562004923816200477c565b935060a089013562004935816200477c565b925060c089013562004947816200477c565b915060e089013562004959816200477c565b809150509295985092959890939650565b60008060008060008060c087890312156200498457600080fd5b6200498f876200464f565b95506200499f6020880162004664565b9450620049af6040880162004664565b9350606087013592506080870135915060a087013590509295509295509295565b60005b83811015620049ed578181015183820152602001620049d3565b50506000910152565b6000815180845262004a10816020860160208601620049d0565b601f01601f19169290920160200192915050565b60208152600062000b916020830184620049f6565b6000806040838503121562004a4d57600080fd5b62004a588362004664565b946020939093013593505050565b6000806000806060858703121562004a7d57600080fd5b843562004a8a816200477c565b935062004a9a602086016200464f565b925060408501356001600160401b038082111562004ab757600080fd5b818701915087601f83011262004acc57600080fd5b81358181111562004adc57600080fd5b88602082850101111562004aef57600080fd5b95989497505060200194505050565b60006020828403121562004b1157600080fd5b813562004842816200477c565b803560ff811681146200306957600080fd5b60008060008060008060c0878903121562004b4a57600080fd5b863562004b57816200477c565b9550602087013562004b69816200477c565b945062004b796040880162004664565b935062004b896060880162004664565b92506080870135915062004ba060a0880162004b1e565b90509295509295509295565b634e487b7160e01b600052604160045260246000fd5b600082601f83011262004bd457600080fd5b81356001600160401b038082111562004bf15762004bf162004bac565b604051601f8301601f19908116603f0116810190828211818310171562004c1c5762004c1c62004bac565b8160405283815286602085880101111562004c3657600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600080600080600080610100898b03121562004c7457600080fd5b62004c7f896200464f565b975062004c8f60208a0162004664565b9650604089013562004ca1816200477c565b9550606089013562004cb3816200477c565b9450608089013562004cc5816200477c565b935062004cd560a08a016200464f565b925060c08901356001600160401b038082111562004cf257600080fd5b62004d008c838d0162004bc2565b935060e08b013591508082111562004d1757600080fd5b5062004d268b828c0162004bc2565b9150509295985092959890939650565b60008060008060008060c0878903121562004d5057600080fd5b863562004d5d816200477c565b9550602087013562004d6f816200477c565b945062004d7f6040880162004664565b935062004d8f6060880162004b1e565b92506080870135915060a08701356001600160401b0381111562004db257600080fd5b62004dc089828a0162004bc2565b9150509295509295509295565b634e487b7160e01b600052601160045260246000fd5b6001600160401b0381811683821601908082111562004e065762004e0662004dcd565b5092915050565b8054600090600181811c908083168062004e2857607f831692505b6020808410820362004e4a57634e487b7160e01b600052602260045260246000fd5b8388526020880182801562004e68576001811462004e7f5762004eac565b60ff198716825285151560051b8201975062004eac565b60008981526020902060005b8781101562004ea65781548482015290860190840162004e8b565b83019850505b5050505050505092915050565b6001600160a01b0387811682528616602082015260c06040820181905260009062004ee79083018762004e0d565b828103606084015262004efb818762004e0d565b9150508360808301526001600160401b03831660a0830152979650505050505050565b808202811582820484141762000b945762000b9462004dcd565b6001600160401b0382811682821603908082111562004e065762004e0662004dcd565b60006001600160401b0380831681810362004f7a5762004f7a62004dcd565b6001019392505050565b60006020828403121562004f9757600080fd5b5051919050565b634e487b7160e01b600052601260045260246000fd5b60008262004fc65762004fc662004f9e565b500490565b8082018082111562000b945762000b9462004dcd565b634e487b7160e01b600052603260045260246000fd5b6000600182016200500c576200500c62004dcd565b5060010190565b60008262005025576200502562004f9e565b500690565b8181038181111562000b945762000b9462004dcd565b60008162005052576200505262004dcd565b506000190190565b6001600160a01b03841681526040602082018190528101829052818360608301376000818301606090810191909152601f909201601f1916010192915050565b600063ffffffff80831681810362004f7a5762004f7a62004dcd565b6001600160a01b03848116825283166020820152606060408201819052600090620050e490830184620049f6565b95945050505050565b6001600160a01b038881168252878116602083015263ffffffff878116604084015290861660608301528416608082015260e060a082018190526000906200513890830185620049f6565b82810360c08401526200514c8185620049f6565b9a9950505050505050505050565b6001600160a01b038781168252861660208201526001600160401b038516604082015260ff841660608201526080810183905260c060a08201819052600090620051a790830184620049f6565b98975050505050505050565b60008251620051c7818460208701620049d0565b9190910192915050565b61032081016103008085843782018360005b600181101562005204578151835260209283019290910190600101620051e3565b5050509392505050565b6000602082840312156200522157600080fd5b815180151581146200484257600080fd5b60006001600160401b03808416806200524f576200524f62004f9e565b9216919091049291505056fe60a06040526040516200092e3803806200092e833981016040819052620000269162000383565b828162000034828262000060565b50506001600160a01b038216608052620000576200005160805190565b620000c6565b50505062000481565b6200006b8262000138565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a2805115620000b857620000b38282620001b8565b505050565b620000c262000235565b5050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f620001086000805160206200090e833981519152546001600160a01b031690565b604080516001600160a01b03928316815291841660208301520160405180910390a1620001358162000257565b50565b806001600160a01b03163b6000036200017457604051634c9c8ce360e01b81526001600160a01b03821660048201526024015b60405180910390fd5b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5b80546001600160a01b0319166001600160a01b039290921691909117905550565b6060600080846001600160a01b031684604051620001d7919062000463565b600060405180830381855af49150503d806000811462000214576040519150601f19603f3d011682016040523d82523d6000602084013e62000219565b606091505b5090925090506200022c8583836200029a565b95945050505050565b3415620002555760405163b398979f60e01b815260040160405180910390fd5b565b6001600160a01b0381166200028357604051633173bdd160e11b8152600060048201526024016200016b565b806000805160206200090e83398151915262000197565b606082620002b357620002ad8262000300565b620002f9565b8151158015620002cb57506001600160a01b0384163b155b15620002f657604051639996b31560e01b81526001600160a01b03851660048201526024016200016b565b50805b9392505050565b805115620003115780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b80516001600160a01b03811681146200034257600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60005b838110156200037a57818101518382015260200162000360565b50506000910152565b6000806000606084860312156200039957600080fd5b620003a4846200032a565b9250620003b4602085016200032a565b60408501519092506001600160401b0380821115620003d257600080fd5b818601915086601f830112620003e757600080fd5b815181811115620003fc57620003fc62000347565b604051601f8201601f19908116603f0116810190838211818310171562000427576200042762000347565b816040528281528960208487010111156200044157600080fd5b620004548360208301602088016200035d565b80955050505050509250925092565b60008251620004778184602087016200035d565b9190910192915050565b6080516104726200049c6000396000601001526104726000f3fe608060405261000c61000e565b005b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316330361006a576000356001600160e01b03191663278f794360e11b146100625761006061006e565b565b61006061007e565b6100605b6100606100796100ad565b6100e5565b60008061008e36600481846102fd565b81019061009b919061033d565b915091506100a98282610109565b5050565b60006100e07f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b905090565b3660008037600080366000845af43d6000803e808015610104573d6000f35b3d6000fd5b61011282610164565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a280511561015c5761015782826101e0565b505050565b6100a9610256565b806001600160a01b03163b60000361019f57604051634c9c8ce360e01b81526001600160a01b03821660048201526024015b60405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b0316846040516101fd919061040d565b600060405180830381855af49150503d8060008114610238576040519150601f19603f3d011682016040523d82523d6000602084013e61023d565b606091505b509150915061024d858383610275565b95945050505050565b34156100605760405163b398979f60e01b815260040160405180910390fd5b60608261028a57610285826102d4565b6102cd565b81511580156102a157506001600160a01b0384163b155b156102ca57604051639996b31560e01b81526001600160a01b0385166004820152602401610196565b50805b9392505050565b8051156102e45780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b6000808585111561030d57600080fd5b8386111561031a57600080fd5b5050820193919092039150565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561035057600080fd5b82356001600160a01b038116811461036757600080fd5b9150602083013567ffffffffffffffff8082111561038457600080fd5b818501915085601f83011261039857600080fd5b8135818111156103aa576103aa610327565b604051601f8201601f19908116603f011681019083821181831017156103d2576103d2610327565b816040528281528860208487010111156103eb57600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b6000825160005b8181101561042e5760208186018101518583015201610414565b50600092019182525091905056fea2646970667358221220d51c6fc430c553d1da2391464fd71dbfa43dc91a9531c4d842a3289bb581301664736f6c63430008140033b53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103a5c5790f581d443ed43873ab47cfb8c5d66a6db268e58b5971bb33fc66e07db1141f8f32ce6198eee741f695cec728bfd32d289f1acf73621fb303581000545e084e94f375e9d647f87f5b2ceffba1e062c70f6009fdbcf80291e803b5c9edd4a2646970667358221220f5ac7cd2c2bce84174660f9e38a04abba5f38d6b7992b6159a363918c4e4e4a664736f6c63430008140033",
}

// PolygonrollupmanagerABI is the input ABI used to generate the binding from.
// Deprecated: Use PolygonrollupmanagerMetaData.ABI instead.
var PolygonrollupmanagerABI = PolygonrollupmanagerMetaData.ABI

// PolygonrollupmanagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PolygonrollupmanagerMetaData.Bin instead.
var PolygonrollupmanagerBin = PolygonrollupmanagerMetaData.Bin

// DeployPolygonrollupmanager deploys a new Ethereum contract, binding an instance of Polygonrollupmanager to it.
func DeployPolygonrollupmanager(auth *bind.TransactOpts, backend bind.ContractBackend, _globalExitRootManager common.Address, _pol common.Address, _bridgeAddress common.Address) (common.Address, *types.Transaction, *Polygonrollupmanager, error) {
	parsed, err := PolygonrollupmanagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PolygonrollupmanagerBin), backend, _globalExitRootManager, _pol, _bridgeAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Polygonrollupmanager{PolygonrollupmanagerCaller: PolygonrollupmanagerCaller{contract: contract}, PolygonrollupmanagerTransactor: PolygonrollupmanagerTransactor{contract: contract}, PolygonrollupmanagerFilterer: PolygonrollupmanagerFilterer{contract: contract}}, nil
}

// Polygonrollupmanager is an auto generated Go binding around an Ethereum contract.
type Polygonrollupmanager struct {
	PolygonrollupmanagerCaller     // Read-only binding to the contract
	PolygonrollupmanagerTransactor // Write-only binding to the contract
	PolygonrollupmanagerFilterer   // Log filterer for contract events
}

// PolygonrollupmanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type PolygonrollupmanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolygonrollupmanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PolygonrollupmanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolygonrollupmanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PolygonrollupmanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolygonrollupmanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PolygonrollupmanagerSession struct {
	Contract     *Polygonrollupmanager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// PolygonrollupmanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PolygonrollupmanagerCallerSession struct {
	Contract *PolygonrollupmanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// PolygonrollupmanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PolygonrollupmanagerTransactorSession struct {
	Contract     *PolygonrollupmanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// PolygonrollupmanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type PolygonrollupmanagerRaw struct {
	Contract *Polygonrollupmanager // Generic contract binding to access the raw methods on
}

// PolygonrollupmanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PolygonrollupmanagerCallerRaw struct {
	Contract *PolygonrollupmanagerCaller // Generic read-only contract binding to access the raw methods on
}

// PolygonrollupmanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PolygonrollupmanagerTransactorRaw struct {
	Contract *PolygonrollupmanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPolygonrollupmanager creates a new instance of Polygonrollupmanager, bound to a specific deployed contract.
func NewPolygonrollupmanager(address common.Address, backend bind.ContractBackend) (*Polygonrollupmanager, error) {
	contract, err := bindPolygonrollupmanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Polygonrollupmanager{PolygonrollupmanagerCaller: PolygonrollupmanagerCaller{contract: contract}, PolygonrollupmanagerTransactor: PolygonrollupmanagerTransactor{contract: contract}, PolygonrollupmanagerFilterer: PolygonrollupmanagerFilterer{contract: contract}}, nil
}

// NewPolygonrollupmanagerCaller creates a new read-only instance of Polygonrollupmanager, bound to a specific deployed contract.
func NewPolygonrollupmanagerCaller(address common.Address, caller bind.ContractCaller) (*PolygonrollupmanagerCaller, error) {
	contract, err := bindPolygonrollupmanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerCaller{contract: contract}, nil
}

// NewPolygonrollupmanagerTransactor creates a new write-only instance of Polygonrollupmanager, bound to a specific deployed contract.
func NewPolygonrollupmanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*PolygonrollupmanagerTransactor, error) {
	contract, err := bindPolygonrollupmanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerTransactor{contract: contract}, nil
}

// NewPolygonrollupmanagerFilterer creates a new log filterer instance of Polygonrollupmanager, bound to a specific deployed contract.
func NewPolygonrollupmanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*PolygonrollupmanagerFilterer, error) {
	contract, err := bindPolygonrollupmanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerFilterer{contract: contract}, nil
}

// bindPolygonrollupmanager binds a generic wrapper to an already deployed contract.
func bindPolygonrollupmanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PolygonrollupmanagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Polygonrollupmanager *PolygonrollupmanagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Polygonrollupmanager.Contract.PolygonrollupmanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Polygonrollupmanager *PolygonrollupmanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.PolygonrollupmanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Polygonrollupmanager *PolygonrollupmanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.PolygonrollupmanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Polygonrollupmanager *PolygonrollupmanagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Polygonrollupmanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Polygonrollupmanager.Contract.DEFAULTADMINROLE(&_Polygonrollupmanager.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Polygonrollupmanager.Contract.DEFAULTADMINROLE(&_Polygonrollupmanager.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) BridgeAddress() (common.Address, error) {
	return _Polygonrollupmanager.Contract.BridgeAddress(&_Polygonrollupmanager.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) BridgeAddress() (common.Address, error) {
	return _Polygonrollupmanager.Contract.BridgeAddress(&_Polygonrollupmanager.CallOpts)
}

// CalculateRewardPerBatch is a free data retrieval call binding the contract method 0x99f5634e.
//
// Solidity: function calculateRewardPerBatch() view returns(uint256)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) CalculateRewardPerBatch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "calculateRewardPerBatch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateRewardPerBatch is a free data retrieval call binding the contract method 0x99f5634e.
//
// Solidity: function calculateRewardPerBatch() view returns(uint256)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) CalculateRewardPerBatch() (*big.Int, error) {
	return _Polygonrollupmanager.Contract.CalculateRewardPerBatch(&_Polygonrollupmanager.CallOpts)
}

// CalculateRewardPerBatch is a free data retrieval call binding the contract method 0x99f5634e.
//
// Solidity: function calculateRewardPerBatch() view returns(uint256)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) CalculateRewardPerBatch() (*big.Int, error) {
	return _Polygonrollupmanager.Contract.CalculateRewardPerBatch(&_Polygonrollupmanager.CallOpts)
}

// ChainIDToRollupID is a free data retrieval call binding the contract method 0x7fb6e76a.
//
// Solidity: function chainIDToRollupID(uint64 chainID) view returns(uint32 rollupID)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) ChainIDToRollupID(opts *bind.CallOpts, chainID uint64) (uint32, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "chainIDToRollupID", chainID)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// ChainIDToRollupID is a free data retrieval call binding the contract method 0x7fb6e76a.
//
// Solidity: function chainIDToRollupID(uint64 chainID) view returns(uint32 rollupID)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) ChainIDToRollupID(chainID uint64) (uint32, error) {
	return _Polygonrollupmanager.Contract.ChainIDToRollupID(&_Polygonrollupmanager.CallOpts, chainID)
}

// ChainIDToRollupID is a free data retrieval call binding the contract method 0x7fb6e76a.
//
// Solidity: function chainIDToRollupID(uint64 chainID) view returns(uint32 rollupID)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) ChainIDToRollupID(chainID uint64) (uint32, error) {
	return _Polygonrollupmanager.Contract.ChainIDToRollupID(&_Polygonrollupmanager.CallOpts, chainID)
}

// GetBatchFee is a free data retrieval call binding the contract method 0x477fa270.
//
// Solidity: function getBatchFee() view returns(uint256)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) GetBatchFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "getBatchFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBatchFee is a free data retrieval call binding the contract method 0x477fa270.
//
// Solidity: function getBatchFee() view returns(uint256)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) GetBatchFee() (*big.Int, error) {
	return _Polygonrollupmanager.Contract.GetBatchFee(&_Polygonrollupmanager.CallOpts)
}

// GetBatchFee is a free data retrieval call binding the contract method 0x477fa270.
//
// Solidity: function getBatchFee() view returns(uint256)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) GetBatchFee() (*big.Int, error) {
	return _Polygonrollupmanager.Contract.GetBatchFee(&_Polygonrollupmanager.CallOpts)
}

// GetForcedBatchFee is a free data retrieval call binding the contract method 0x60469169.
//
// Solidity: function getForcedBatchFee() view returns(uint256)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) GetForcedBatchFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "getForcedBatchFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetForcedBatchFee is a free data retrieval call binding the contract method 0x60469169.
//
// Solidity: function getForcedBatchFee() view returns(uint256)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) GetForcedBatchFee() (*big.Int, error) {
	return _Polygonrollupmanager.Contract.GetForcedBatchFee(&_Polygonrollupmanager.CallOpts)
}

// GetForcedBatchFee is a free data retrieval call binding the contract method 0x60469169.
//
// Solidity: function getForcedBatchFee() view returns(uint256)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) GetForcedBatchFee() (*big.Int, error) {
	return _Polygonrollupmanager.Contract.GetForcedBatchFee(&_Polygonrollupmanager.CallOpts)
}

// GetInputSnarkBytes is a free data retrieval call binding the contract method 0x7975fcfe.
//
// Solidity: function getInputSnarkBytes(uint32 rollupID, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 oldStateRoot, bytes32 newStateRoot) view returns(bytes)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) GetInputSnarkBytes(opts *bind.CallOpts, rollupID uint32, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, oldStateRoot [32]byte, newStateRoot [32]byte) ([]byte, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "getInputSnarkBytes", rollupID, initNumBatch, finalNewBatch, newLocalExitRoot, oldStateRoot, newStateRoot)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetInputSnarkBytes is a free data retrieval call binding the contract method 0x7975fcfe.
//
// Solidity: function getInputSnarkBytes(uint32 rollupID, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 oldStateRoot, bytes32 newStateRoot) view returns(bytes)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) GetInputSnarkBytes(rollupID uint32, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, oldStateRoot [32]byte, newStateRoot [32]byte) ([]byte, error) {
	return _Polygonrollupmanager.Contract.GetInputSnarkBytes(&_Polygonrollupmanager.CallOpts, rollupID, initNumBatch, finalNewBatch, newLocalExitRoot, oldStateRoot, newStateRoot)
}

// GetInputSnarkBytes is a free data retrieval call binding the contract method 0x7975fcfe.
//
// Solidity: function getInputSnarkBytes(uint32 rollupID, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 oldStateRoot, bytes32 newStateRoot) view returns(bytes)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) GetInputSnarkBytes(rollupID uint32, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, oldStateRoot [32]byte, newStateRoot [32]byte) ([]byte, error) {
	return _Polygonrollupmanager.Contract.GetInputSnarkBytes(&_Polygonrollupmanager.CallOpts, rollupID, initNumBatch, finalNewBatch, newLocalExitRoot, oldStateRoot, newStateRoot)
}

// GetLastVerifiedBatch is a free data retrieval call binding the contract method 0x11f6b287.
//
// Solidity: function getLastVerifiedBatch(uint32 rollupID) view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) GetLastVerifiedBatch(opts *bind.CallOpts, rollupID uint32) (uint64, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "getLastVerifiedBatch", rollupID)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetLastVerifiedBatch is a free data retrieval call binding the contract method 0x11f6b287.
//
// Solidity: function getLastVerifiedBatch(uint32 rollupID) view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) GetLastVerifiedBatch(rollupID uint32) (uint64, error) {
	return _Polygonrollupmanager.Contract.GetLastVerifiedBatch(&_Polygonrollupmanager.CallOpts, rollupID)
}

// GetLastVerifiedBatch is a free data retrieval call binding the contract method 0x11f6b287.
//
// Solidity: function getLastVerifiedBatch(uint32 rollupID) view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) GetLastVerifiedBatch(rollupID uint32) (uint64, error) {
	return _Polygonrollupmanager.Contract.GetLastVerifiedBatch(&_Polygonrollupmanager.CallOpts, rollupID)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Polygonrollupmanager.Contract.GetRoleAdmin(&_Polygonrollupmanager.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Polygonrollupmanager.Contract.GetRoleAdmin(&_Polygonrollupmanager.CallOpts, role)
}

// GetRollupBatchNumToStateRoot is a free data retrieval call binding the contract method 0x55a71ee0.
//
// Solidity: function getRollupBatchNumToStateRoot(uint32 rollupID, uint64 batchNum) view returns(bytes32)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) GetRollupBatchNumToStateRoot(opts *bind.CallOpts, rollupID uint32, batchNum uint64) ([32]byte, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "getRollupBatchNumToStateRoot", rollupID, batchNum)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRollupBatchNumToStateRoot is a free data retrieval call binding the contract method 0x55a71ee0.
//
// Solidity: function getRollupBatchNumToStateRoot(uint32 rollupID, uint64 batchNum) view returns(bytes32)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) GetRollupBatchNumToStateRoot(rollupID uint32, batchNum uint64) ([32]byte, error) {
	return _Polygonrollupmanager.Contract.GetRollupBatchNumToStateRoot(&_Polygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupBatchNumToStateRoot is a free data retrieval call binding the contract method 0x55a71ee0.
//
// Solidity: function getRollupBatchNumToStateRoot(uint32 rollupID, uint64 batchNum) view returns(bytes32)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) GetRollupBatchNumToStateRoot(rollupID uint32, batchNum uint64) ([32]byte, error) {
	return _Polygonrollupmanager.Contract.GetRollupBatchNumToStateRoot(&_Polygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupExitRoot is a free data retrieval call binding the contract method 0xa2967d99.
//
// Solidity: function getRollupExitRoot() view returns(bytes32)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) GetRollupExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "getRollupExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRollupExitRoot is a free data retrieval call binding the contract method 0xa2967d99.
//
// Solidity: function getRollupExitRoot() view returns(bytes32)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) GetRollupExitRoot() ([32]byte, error) {
	return _Polygonrollupmanager.Contract.GetRollupExitRoot(&_Polygonrollupmanager.CallOpts)
}

// GetRollupExitRoot is a free data retrieval call binding the contract method 0xa2967d99.
//
// Solidity: function getRollupExitRoot() view returns(bytes32)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) GetRollupExitRoot() ([32]byte, error) {
	return _Polygonrollupmanager.Contract.GetRollupExitRoot(&_Polygonrollupmanager.CallOpts)
}

// GetRollupPendingStateTransitions is a free data retrieval call binding the contract method 0xb99d0ad7.
//
// Solidity: function getRollupPendingStateTransitions(uint32 rollupID, uint64 batchNum) view returns((uint64,uint64,bytes32,bytes32))
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) GetRollupPendingStateTransitions(opts *bind.CallOpts, rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesPendingState, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "getRollupPendingStateTransitions", rollupID, batchNum)

	if err != nil {
		return *new(LegacyZKEVMStateVariablesPendingState), err
	}

	out0 := *abi.ConvertType(out[0], new(LegacyZKEVMStateVariablesPendingState)).(*LegacyZKEVMStateVariablesPendingState)

	return out0, err

}

// GetRollupPendingStateTransitions is a free data retrieval call binding the contract method 0xb99d0ad7.
//
// Solidity: function getRollupPendingStateTransitions(uint32 rollupID, uint64 batchNum) view returns((uint64,uint64,bytes32,bytes32))
func (_Polygonrollupmanager *PolygonrollupmanagerSession) GetRollupPendingStateTransitions(rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesPendingState, error) {
	return _Polygonrollupmanager.Contract.GetRollupPendingStateTransitions(&_Polygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupPendingStateTransitions is a free data retrieval call binding the contract method 0xb99d0ad7.
//
// Solidity: function getRollupPendingStateTransitions(uint32 rollupID, uint64 batchNum) view returns((uint64,uint64,bytes32,bytes32))
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) GetRollupPendingStateTransitions(rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesPendingState, error) {
	return _Polygonrollupmanager.Contract.GetRollupPendingStateTransitions(&_Polygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupSequencedBatches is a free data retrieval call binding the contract method 0x25280169.
//
// Solidity: function getRollupSequencedBatches(uint32 rollupID, uint64 batchNum) view returns((bytes32,uint64,uint64))
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) GetRollupSequencedBatches(opts *bind.CallOpts, rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesSequencedBatchData, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "getRollupSequencedBatches", rollupID, batchNum)

	if err != nil {
		return *new(LegacyZKEVMStateVariablesSequencedBatchData), err
	}

	out0 := *abi.ConvertType(out[0], new(LegacyZKEVMStateVariablesSequencedBatchData)).(*LegacyZKEVMStateVariablesSequencedBatchData)

	return out0, err

}

// GetRollupSequencedBatches is a free data retrieval call binding the contract method 0x25280169.
//
// Solidity: function getRollupSequencedBatches(uint32 rollupID, uint64 batchNum) view returns((bytes32,uint64,uint64))
func (_Polygonrollupmanager *PolygonrollupmanagerSession) GetRollupSequencedBatches(rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesSequencedBatchData, error) {
	return _Polygonrollupmanager.Contract.GetRollupSequencedBatches(&_Polygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupSequencedBatches is a free data retrieval call binding the contract method 0x25280169.
//
// Solidity: function getRollupSequencedBatches(uint32 rollupID, uint64 batchNum) view returns((bytes32,uint64,uint64))
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) GetRollupSequencedBatches(rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesSequencedBatchData, error) {
	return _Polygonrollupmanager.Contract.GetRollupSequencedBatches(&_Polygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) GlobalExitRootManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "globalExitRootManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) GlobalExitRootManager() (common.Address, error) {
	return _Polygonrollupmanager.Contract.GlobalExitRootManager(&_Polygonrollupmanager.CallOpts)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) GlobalExitRootManager() (common.Address, error) {
	return _Polygonrollupmanager.Contract.GlobalExitRootManager(&_Polygonrollupmanager.CallOpts)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Polygonrollupmanager.Contract.HasRole(&_Polygonrollupmanager.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Polygonrollupmanager.Contract.HasRole(&_Polygonrollupmanager.CallOpts, role, account)
}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) IsEmergencyState(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "isEmergencyState")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) IsEmergencyState() (bool, error) {
	return _Polygonrollupmanager.Contract.IsEmergencyState(&_Polygonrollupmanager.CallOpts)
}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) IsEmergencyState() (bool, error) {
	return _Polygonrollupmanager.Contract.IsEmergencyState(&_Polygonrollupmanager.CallOpts)
}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x080b3111.
//
// Solidity: function isPendingStateConsolidable(uint32 rollupID, uint64 pendingStateNum) view returns(bool)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) IsPendingStateConsolidable(opts *bind.CallOpts, rollupID uint32, pendingStateNum uint64) (bool, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "isPendingStateConsolidable", rollupID, pendingStateNum)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x080b3111.
//
// Solidity: function isPendingStateConsolidable(uint32 rollupID, uint64 pendingStateNum) view returns(bool)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) IsPendingStateConsolidable(rollupID uint32, pendingStateNum uint64) (bool, error) {
	return _Polygonrollupmanager.Contract.IsPendingStateConsolidable(&_Polygonrollupmanager.CallOpts, rollupID, pendingStateNum)
}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x080b3111.
//
// Solidity: function isPendingStateConsolidable(uint32 rollupID, uint64 pendingStateNum) view returns(bool)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) IsPendingStateConsolidable(rollupID uint32, pendingStateNum uint64) (bool, error) {
	return _Polygonrollupmanager.Contract.IsPendingStateConsolidable(&_Polygonrollupmanager.CallOpts, rollupID, pendingStateNum)
}

// LastAggregationTimestamp is a free data retrieval call binding the contract method 0xc1acbc34.
//
// Solidity: function lastAggregationTimestamp() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) LastAggregationTimestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "lastAggregationTimestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastAggregationTimestamp is a free data retrieval call binding the contract method 0xc1acbc34.
//
// Solidity: function lastAggregationTimestamp() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) LastAggregationTimestamp() (uint64, error) {
	return _Polygonrollupmanager.Contract.LastAggregationTimestamp(&_Polygonrollupmanager.CallOpts)
}

// LastAggregationTimestamp is a free data retrieval call binding the contract method 0xc1acbc34.
//
// Solidity: function lastAggregationTimestamp() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) LastAggregationTimestamp() (uint64, error) {
	return _Polygonrollupmanager.Contract.LastAggregationTimestamp(&_Polygonrollupmanager.CallOpts)
}

// MultiplierBatchFee is a free data retrieval call binding the contract method 0xafd23cbe.
//
// Solidity: function multiplierBatchFee() view returns(uint16)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) MultiplierBatchFee(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "multiplierBatchFee")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// MultiplierBatchFee is a free data retrieval call binding the contract method 0xafd23cbe.
//
// Solidity: function multiplierBatchFee() view returns(uint16)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) MultiplierBatchFee() (uint16, error) {
	return _Polygonrollupmanager.Contract.MultiplierBatchFee(&_Polygonrollupmanager.CallOpts)
}

// MultiplierBatchFee is a free data retrieval call binding the contract method 0xafd23cbe.
//
// Solidity: function multiplierBatchFee() view returns(uint16)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) MultiplierBatchFee() (uint16, error) {
	return _Polygonrollupmanager.Contract.MultiplierBatchFee(&_Polygonrollupmanager.CallOpts)
}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) PendingStateTimeout(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "pendingStateTimeout")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) PendingStateTimeout() (uint64, error) {
	return _Polygonrollupmanager.Contract.PendingStateTimeout(&_Polygonrollupmanager.CallOpts)
}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) PendingStateTimeout() (uint64, error) {
	return _Polygonrollupmanager.Contract.PendingStateTimeout(&_Polygonrollupmanager.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) Pol(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "pol")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) Pol() (common.Address, error) {
	return _Polygonrollupmanager.Contract.Pol(&_Polygonrollupmanager.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) Pol() (common.Address, error) {
	return _Polygonrollupmanager.Contract.Pol(&_Polygonrollupmanager.CallOpts)
}

// RollupAddressToID is a free data retrieval call binding the contract method 0xceee281d.
//
// Solidity: function rollupAddressToID(address rollupAddress) view returns(uint32 rollupID)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) RollupAddressToID(opts *bind.CallOpts, rollupAddress common.Address) (uint32, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "rollupAddressToID", rollupAddress)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RollupAddressToID is a free data retrieval call binding the contract method 0xceee281d.
//
// Solidity: function rollupAddressToID(address rollupAddress) view returns(uint32 rollupID)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) RollupAddressToID(rollupAddress common.Address) (uint32, error) {
	return _Polygonrollupmanager.Contract.RollupAddressToID(&_Polygonrollupmanager.CallOpts, rollupAddress)
}

// RollupAddressToID is a free data retrieval call binding the contract method 0xceee281d.
//
// Solidity: function rollupAddressToID(address rollupAddress) view returns(uint32 rollupID)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) RollupAddressToID(rollupAddress common.Address) (uint32, error) {
	return _Polygonrollupmanager.Contract.RollupAddressToID(&_Polygonrollupmanager.CallOpts, rollupAddress)
}

// RollupCount is a free data retrieval call binding the contract method 0xf4e92675.
//
// Solidity: function rollupCount() view returns(uint32)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) RollupCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "rollupCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RollupCount is a free data retrieval call binding the contract method 0xf4e92675.
//
// Solidity: function rollupCount() view returns(uint32)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) RollupCount() (uint32, error) {
	return _Polygonrollupmanager.Contract.RollupCount(&_Polygonrollupmanager.CallOpts)
}

// RollupCount is a free data retrieval call binding the contract method 0xf4e92675.
//
// Solidity: function rollupCount() view returns(uint32)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) RollupCount() (uint32, error) {
	return _Polygonrollupmanager.Contract.RollupCount(&_Polygonrollupmanager.CallOpts)
}

// RollupIDToRollupData is a free data retrieval call binding the contract method 0xf9c4c2ae.
//
// Solidity: function rollupIDToRollupData(uint32 rollupID) view returns(address rollupContract, uint64 chainID, address verifier, uint64 forkID, bytes32 lastLocalExitRoot, uint64 lastBatchSequenced, uint64 lastVerifiedBatch, uint64 lastPendingState, uint64 lastPendingStateConsolidated, uint64 lastVerifiedBatchBeforeUpgrade, uint64 rollupTypeID, uint8 rollupCompatibilityID)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) RollupIDToRollupData(opts *bind.CallOpts, rollupID uint32) (struct {
	RollupContract                 common.Address
	ChainID                        uint64
	Verifier                       common.Address
	ForkID                         uint64
	LastLocalExitRoot              [32]byte
	LastBatchSequenced             uint64
	LastVerifiedBatch              uint64
	LastPendingState               uint64
	LastPendingStateConsolidated   uint64
	LastVerifiedBatchBeforeUpgrade uint64
	RollupTypeID                   uint64
	RollupCompatibilityID          uint8
}, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "rollupIDToRollupData", rollupID)

	outstruct := new(struct {
		RollupContract                 common.Address
		ChainID                        uint64
		Verifier                       common.Address
		ForkID                         uint64
		LastLocalExitRoot              [32]byte
		LastBatchSequenced             uint64
		LastVerifiedBatch              uint64
		LastPendingState               uint64
		LastPendingStateConsolidated   uint64
		LastVerifiedBatchBeforeUpgrade uint64
		RollupTypeID                   uint64
		RollupCompatibilityID          uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RollupContract = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ChainID = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.Verifier = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.ForkID = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.LastLocalExitRoot = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)
	outstruct.LastBatchSequenced = *abi.ConvertType(out[5], new(uint64)).(*uint64)
	outstruct.LastVerifiedBatch = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.LastPendingState = *abi.ConvertType(out[7], new(uint64)).(*uint64)
	outstruct.LastPendingStateConsolidated = *abi.ConvertType(out[8], new(uint64)).(*uint64)
	outstruct.LastVerifiedBatchBeforeUpgrade = *abi.ConvertType(out[9], new(uint64)).(*uint64)
	outstruct.RollupTypeID = *abi.ConvertType(out[10], new(uint64)).(*uint64)
	outstruct.RollupCompatibilityID = *abi.ConvertType(out[11], new(uint8)).(*uint8)

	return *outstruct, err

}

// RollupIDToRollupData is a free data retrieval call binding the contract method 0xf9c4c2ae.
//
// Solidity: function rollupIDToRollupData(uint32 rollupID) view returns(address rollupContract, uint64 chainID, address verifier, uint64 forkID, bytes32 lastLocalExitRoot, uint64 lastBatchSequenced, uint64 lastVerifiedBatch, uint64 lastPendingState, uint64 lastPendingStateConsolidated, uint64 lastVerifiedBatchBeforeUpgrade, uint64 rollupTypeID, uint8 rollupCompatibilityID)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) RollupIDToRollupData(rollupID uint32) (struct {
	RollupContract                 common.Address
	ChainID                        uint64
	Verifier                       common.Address
	ForkID                         uint64
	LastLocalExitRoot              [32]byte
	LastBatchSequenced             uint64
	LastVerifiedBatch              uint64
	LastPendingState               uint64
	LastPendingStateConsolidated   uint64
	LastVerifiedBatchBeforeUpgrade uint64
	RollupTypeID                   uint64
	RollupCompatibilityID          uint8
}, error) {
	return _Polygonrollupmanager.Contract.RollupIDToRollupData(&_Polygonrollupmanager.CallOpts, rollupID)
}

// RollupIDToRollupData is a free data retrieval call binding the contract method 0xf9c4c2ae.
//
// Solidity: function rollupIDToRollupData(uint32 rollupID) view returns(address rollupContract, uint64 chainID, address verifier, uint64 forkID, bytes32 lastLocalExitRoot, uint64 lastBatchSequenced, uint64 lastVerifiedBatch, uint64 lastPendingState, uint64 lastPendingStateConsolidated, uint64 lastVerifiedBatchBeforeUpgrade, uint64 rollupTypeID, uint8 rollupCompatibilityID)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) RollupIDToRollupData(rollupID uint32) (struct {
	RollupContract                 common.Address
	ChainID                        uint64
	Verifier                       common.Address
	ForkID                         uint64
	LastLocalExitRoot              [32]byte
	LastBatchSequenced             uint64
	LastVerifiedBatch              uint64
	LastPendingState               uint64
	LastPendingStateConsolidated   uint64
	LastVerifiedBatchBeforeUpgrade uint64
	RollupTypeID                   uint64
	RollupCompatibilityID          uint8
}, error) {
	return _Polygonrollupmanager.Contract.RollupIDToRollupData(&_Polygonrollupmanager.CallOpts, rollupID)
}

// RollupTypeCount is a free data retrieval call binding the contract method 0x1796a1ae.
//
// Solidity: function rollupTypeCount() view returns(uint32)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) RollupTypeCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "rollupTypeCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RollupTypeCount is a free data retrieval call binding the contract method 0x1796a1ae.
//
// Solidity: function rollupTypeCount() view returns(uint32)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) RollupTypeCount() (uint32, error) {
	return _Polygonrollupmanager.Contract.RollupTypeCount(&_Polygonrollupmanager.CallOpts)
}

// RollupTypeCount is a free data retrieval call binding the contract method 0x1796a1ae.
//
// Solidity: function rollupTypeCount() view returns(uint32)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) RollupTypeCount() (uint32, error) {
	return _Polygonrollupmanager.Contract.RollupTypeCount(&_Polygonrollupmanager.CallOpts)
}

// RollupTypeMap is a free data retrieval call binding the contract method 0xa4810bb2.
//
// Solidity: function rollupTypeMap(uint64 rollupTypeID) view returns(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bool obsolete, bytes32 genesis)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) RollupTypeMap(opts *bind.CallOpts, rollupTypeID uint64) (struct {
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Obsolete                bool
	Genesis                 [32]byte
}, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "rollupTypeMap", rollupTypeID)

	outstruct := new(struct {
		ConsensusImplementation common.Address
		Verifier                common.Address
		ForkID                  uint64
		RollupCompatibilityID   uint8
		Obsolete                bool
		Genesis                 [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConsensusImplementation = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Verifier = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.ForkID = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.RollupCompatibilityID = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.Obsolete = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.Genesis = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// RollupTypeMap is a free data retrieval call binding the contract method 0xa4810bb2.
//
// Solidity: function rollupTypeMap(uint64 rollupTypeID) view returns(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bool obsolete, bytes32 genesis)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) RollupTypeMap(rollupTypeID uint64) (struct {
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Obsolete                bool
	Genesis                 [32]byte
}, error) {
	return _Polygonrollupmanager.Contract.RollupTypeMap(&_Polygonrollupmanager.CallOpts, rollupTypeID)
}

// RollupTypeMap is a free data retrieval call binding the contract method 0xa4810bb2.
//
// Solidity: function rollupTypeMap(uint64 rollupTypeID) view returns(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bool obsolete, bytes32 genesis)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) RollupTypeMap(rollupTypeID uint64) (struct {
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Obsolete                bool
	Genesis                 [32]byte
}, error) {
	return _Polygonrollupmanager.Contract.RollupTypeMap(&_Polygonrollupmanager.CallOpts, rollupTypeID)
}

// TotalSequencedBatches is a free data retrieval call binding the contract method 0x066ec012.
//
// Solidity: function totalSequencedBatches() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) TotalSequencedBatches(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "totalSequencedBatches")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// TotalSequencedBatches is a free data retrieval call binding the contract method 0x066ec012.
//
// Solidity: function totalSequencedBatches() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) TotalSequencedBatches() (uint64, error) {
	return _Polygonrollupmanager.Contract.TotalSequencedBatches(&_Polygonrollupmanager.CallOpts)
}

// TotalSequencedBatches is a free data retrieval call binding the contract method 0x066ec012.
//
// Solidity: function totalSequencedBatches() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) TotalSequencedBatches() (uint64, error) {
	return _Polygonrollupmanager.Contract.TotalSequencedBatches(&_Polygonrollupmanager.CallOpts)
}

// TotalVerifiedBatches is a free data retrieval call binding the contract method 0xdde0ff77.
//
// Solidity: function totalVerifiedBatches() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) TotalVerifiedBatches(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "totalVerifiedBatches")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// TotalVerifiedBatches is a free data retrieval call binding the contract method 0xdde0ff77.
//
// Solidity: function totalVerifiedBatches() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) TotalVerifiedBatches() (uint64, error) {
	return _Polygonrollupmanager.Contract.TotalVerifiedBatches(&_Polygonrollupmanager.CallOpts)
}

// TotalVerifiedBatches is a free data retrieval call binding the contract method 0xdde0ff77.
//
// Solidity: function totalVerifiedBatches() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) TotalVerifiedBatches() (uint64, error) {
	return _Polygonrollupmanager.Contract.TotalVerifiedBatches(&_Polygonrollupmanager.CallOpts)
}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) TrustedAggregatorTimeout(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "trustedAggregatorTimeout")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) TrustedAggregatorTimeout() (uint64, error) {
	return _Polygonrollupmanager.Contract.TrustedAggregatorTimeout(&_Polygonrollupmanager.CallOpts)
}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) TrustedAggregatorTimeout() (uint64, error) {
	return _Polygonrollupmanager.Contract.TrustedAggregatorTimeout(&_Polygonrollupmanager.CallOpts)
}

// VerifyBatchTimeTarget is a free data retrieval call binding the contract method 0x0a0d9fbe.
//
// Solidity: function verifyBatchTimeTarget() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCaller) VerifyBatchTimeTarget(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Polygonrollupmanager.contract.Call(opts, &out, "verifyBatchTimeTarget")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// VerifyBatchTimeTarget is a free data retrieval call binding the contract method 0x0a0d9fbe.
//
// Solidity: function verifyBatchTimeTarget() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) VerifyBatchTimeTarget() (uint64, error) {
	return _Polygonrollupmanager.Contract.VerifyBatchTimeTarget(&_Polygonrollupmanager.CallOpts)
}

// VerifyBatchTimeTarget is a free data retrieval call binding the contract method 0x0a0d9fbe.
//
// Solidity: function verifyBatchTimeTarget() view returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerCallerSession) VerifyBatchTimeTarget() (uint64, error) {
	return _Polygonrollupmanager.Contract.VerifyBatchTimeTarget(&_Polygonrollupmanager.CallOpts)
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x2072f6c5.
//
// Solidity: function activateEmergencyState() returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) ActivateEmergencyState(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "activateEmergencyState")
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x2072f6c5.
//
// Solidity: function activateEmergencyState() returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) ActivateEmergencyState() (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.ActivateEmergencyState(&_Polygonrollupmanager.TransactOpts)
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x2072f6c5.
//
// Solidity: function activateEmergencyState() returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) ActivateEmergencyState() (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.ActivateEmergencyState(&_Polygonrollupmanager.TransactOpts)
}

// AddExistingRollup is a paid mutator transaction binding the contract method 0xe0bfd3d2.
//
// Solidity: function addExistingRollup(address rollupAddress, address verifier, uint64 forkID, uint64 chainID, bytes32 genesis, uint8 rollupCompatibilityID) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) AddExistingRollup(opts *bind.TransactOpts, rollupAddress common.Address, verifier common.Address, forkID uint64, chainID uint64, genesis [32]byte, rollupCompatibilityID uint8) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "addExistingRollup", rollupAddress, verifier, forkID, chainID, genesis, rollupCompatibilityID)
}

// AddExistingRollup is a paid mutator transaction binding the contract method 0xe0bfd3d2.
//
// Solidity: function addExistingRollup(address rollupAddress, address verifier, uint64 forkID, uint64 chainID, bytes32 genesis, uint8 rollupCompatibilityID) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) AddExistingRollup(rollupAddress common.Address, verifier common.Address, forkID uint64, chainID uint64, genesis [32]byte, rollupCompatibilityID uint8) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.AddExistingRollup(&_Polygonrollupmanager.TransactOpts, rollupAddress, verifier, forkID, chainID, genesis, rollupCompatibilityID)
}

// AddExistingRollup is a paid mutator transaction binding the contract method 0xe0bfd3d2.
//
// Solidity: function addExistingRollup(address rollupAddress, address verifier, uint64 forkID, uint64 chainID, bytes32 genesis, uint8 rollupCompatibilityID) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) AddExistingRollup(rollupAddress common.Address, verifier common.Address, forkID uint64, chainID uint64, genesis [32]byte, rollupCompatibilityID uint8) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.AddExistingRollup(&_Polygonrollupmanager.TransactOpts, rollupAddress, verifier, forkID, chainID, genesis, rollupCompatibilityID)
}

// AddNewRollupType is a paid mutator transaction binding the contract method 0xf34eb8eb.
//
// Solidity: function addNewRollupType(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) AddNewRollupType(opts *bind.TransactOpts, consensusImplementation common.Address, verifier common.Address, forkID uint64, rollupCompatibilityID uint8, genesis [32]byte, description string) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "addNewRollupType", consensusImplementation, verifier, forkID, rollupCompatibilityID, genesis, description)
}

// AddNewRollupType is a paid mutator transaction binding the contract method 0xf34eb8eb.
//
// Solidity: function addNewRollupType(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) AddNewRollupType(consensusImplementation common.Address, verifier common.Address, forkID uint64, rollupCompatibilityID uint8, genesis [32]byte, description string) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.AddNewRollupType(&_Polygonrollupmanager.TransactOpts, consensusImplementation, verifier, forkID, rollupCompatibilityID, genesis, description)
}

// AddNewRollupType is a paid mutator transaction binding the contract method 0xf34eb8eb.
//
// Solidity: function addNewRollupType(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) AddNewRollupType(consensusImplementation common.Address, verifier common.Address, forkID uint64, rollupCompatibilityID uint8, genesis [32]byte, description string) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.AddNewRollupType(&_Polygonrollupmanager.TransactOpts, consensusImplementation, verifier, forkID, rollupCompatibilityID, genesis, description)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x1608859c.
//
// Solidity: function consolidatePendingState(uint32 rollupID, uint64 pendingStateNum) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) ConsolidatePendingState(opts *bind.TransactOpts, rollupID uint32, pendingStateNum uint64) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "consolidatePendingState", rollupID, pendingStateNum)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x1608859c.
//
// Solidity: function consolidatePendingState(uint32 rollupID, uint64 pendingStateNum) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) ConsolidatePendingState(rollupID uint32, pendingStateNum uint64) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.ConsolidatePendingState(&_Polygonrollupmanager.TransactOpts, rollupID, pendingStateNum)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x1608859c.
//
// Solidity: function consolidatePendingState(uint32 rollupID, uint64 pendingStateNum) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) ConsolidatePendingState(rollupID uint32, pendingStateNum uint64) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.ConsolidatePendingState(&_Polygonrollupmanager.TransactOpts, rollupID, pendingStateNum)
}

// CreateNewRollup is a paid mutator transaction binding the contract method 0xe136edf7.
//
// Solidity: function createNewRollup(uint32 rollupTypeID, uint64 chainID, address admin, address sequencer, address gasTokenAddress, uint32 gasTokenNetwork, string sequencerURL, string networkName) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) CreateNewRollup(opts *bind.TransactOpts, rollupTypeID uint32, chainID uint64, admin common.Address, sequencer common.Address, gasTokenAddress common.Address, gasTokenNetwork uint32, sequencerURL string, networkName string) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "createNewRollup", rollupTypeID, chainID, admin, sequencer, gasTokenAddress, gasTokenNetwork, sequencerURL, networkName)
}

// CreateNewRollup is a paid mutator transaction binding the contract method 0xe136edf7.
//
// Solidity: function createNewRollup(uint32 rollupTypeID, uint64 chainID, address admin, address sequencer, address gasTokenAddress, uint32 gasTokenNetwork, string sequencerURL, string networkName) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) CreateNewRollup(rollupTypeID uint32, chainID uint64, admin common.Address, sequencer common.Address, gasTokenAddress common.Address, gasTokenNetwork uint32, sequencerURL string, networkName string) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.CreateNewRollup(&_Polygonrollupmanager.TransactOpts, rollupTypeID, chainID, admin, sequencer, gasTokenAddress, gasTokenNetwork, sequencerURL, networkName)
}

// CreateNewRollup is a paid mutator transaction binding the contract method 0xe136edf7.
//
// Solidity: function createNewRollup(uint32 rollupTypeID, uint64 chainID, address admin, address sequencer, address gasTokenAddress, uint32 gasTokenNetwork, string sequencerURL, string networkName) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) CreateNewRollup(rollupTypeID uint32, chainID uint64, admin common.Address, sequencer common.Address, gasTokenAddress common.Address, gasTokenNetwork uint32, sequencerURL string, networkName string) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.CreateNewRollup(&_Polygonrollupmanager.TransactOpts, rollupTypeID, chainID, admin, sequencer, gasTokenAddress, gasTokenNetwork, sequencerURL, networkName)
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) DeactivateEmergencyState(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "deactivateEmergencyState")
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) DeactivateEmergencyState() (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.DeactivateEmergencyState(&_Polygonrollupmanager.TransactOpts)
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) DeactivateEmergencyState() (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.DeactivateEmergencyState(&_Polygonrollupmanager.TransactOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.GrantRole(&_Polygonrollupmanager.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.GrantRole(&_Polygonrollupmanager.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x4d4d3338.
//
// Solidity: function initialize(address trustedAggregator, uint64 _pendingStateTimeout, uint64 _trustedAggregatorTimeout, address admin, address timelock, address emergencyCouncil, address polygonZkEVM, address zkEVMVerifier) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) Initialize(opts *bind.TransactOpts, trustedAggregator common.Address, _pendingStateTimeout uint64, _trustedAggregatorTimeout uint64, admin common.Address, timelock common.Address, emergencyCouncil common.Address, polygonZkEVM common.Address, zkEVMVerifier common.Address) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "initialize", trustedAggregator, _pendingStateTimeout, _trustedAggregatorTimeout, admin, timelock, emergencyCouncil, polygonZkEVM, zkEVMVerifier)
}

// Initialize is a paid mutator transaction binding the contract method 0x4d4d3338.
//
// Solidity: function initialize(address trustedAggregator, uint64 _pendingStateTimeout, uint64 _trustedAggregatorTimeout, address admin, address timelock, address emergencyCouncil, address polygonZkEVM, address zkEVMVerifier) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) Initialize(trustedAggregator common.Address, _pendingStateTimeout uint64, _trustedAggregatorTimeout uint64, admin common.Address, timelock common.Address, emergencyCouncil common.Address, polygonZkEVM common.Address, zkEVMVerifier common.Address) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.Initialize(&_Polygonrollupmanager.TransactOpts, trustedAggregator, _pendingStateTimeout, _trustedAggregatorTimeout, admin, timelock, emergencyCouncil, polygonZkEVM, zkEVMVerifier)
}

// Initialize is a paid mutator transaction binding the contract method 0x4d4d3338.
//
// Solidity: function initialize(address trustedAggregator, uint64 _pendingStateTimeout, uint64 _trustedAggregatorTimeout, address admin, address timelock, address emergencyCouncil, address polygonZkEVM, address zkEVMVerifier) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) Initialize(trustedAggregator common.Address, _pendingStateTimeout uint64, _trustedAggregatorTimeout uint64, admin common.Address, timelock common.Address, emergencyCouncil common.Address, polygonZkEVM common.Address, zkEVMVerifier common.Address) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.Initialize(&_Polygonrollupmanager.TransactOpts, trustedAggregator, _pendingStateTimeout, _trustedAggregatorTimeout, admin, timelock, emergencyCouncil, polygonZkEVM, zkEVMVerifier)
}

// ObsoleteRollupType is a paid mutator transaction binding the contract method 0x7222020f.
//
// Solidity: function obsoleteRollupType(uint32 rollupTypeID) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) ObsoleteRollupType(opts *bind.TransactOpts, rollupTypeID uint32) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "obsoleteRollupType", rollupTypeID)
}

// ObsoleteRollupType is a paid mutator transaction binding the contract method 0x7222020f.
//
// Solidity: function obsoleteRollupType(uint32 rollupTypeID) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) ObsoleteRollupType(rollupTypeID uint32) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.ObsoleteRollupType(&_Polygonrollupmanager.TransactOpts, rollupTypeID)
}

// ObsoleteRollupType is a paid mutator transaction binding the contract method 0x7222020f.
//
// Solidity: function obsoleteRollupType(uint32 rollupTypeID) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) ObsoleteRollupType(rollupTypeID uint32) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.ObsoleteRollupType(&_Polygonrollupmanager.TransactOpts, rollupTypeID)
}

// OnSequenceBatches is a paid mutator transaction binding the contract method 0x9a908e73.
//
// Solidity: function onSequenceBatches(uint64 newSequencedBatches, bytes32 newAccInputHash) returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) OnSequenceBatches(opts *bind.TransactOpts, newSequencedBatches uint64, newAccInputHash [32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "onSequenceBatches", newSequencedBatches, newAccInputHash)
}

// OnSequenceBatches is a paid mutator transaction binding the contract method 0x9a908e73.
//
// Solidity: function onSequenceBatches(uint64 newSequencedBatches, bytes32 newAccInputHash) returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerSession) OnSequenceBatches(newSequencedBatches uint64, newAccInputHash [32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.OnSequenceBatches(&_Polygonrollupmanager.TransactOpts, newSequencedBatches, newAccInputHash)
}

// OnSequenceBatches is a paid mutator transaction binding the contract method 0x9a908e73.
//
// Solidity: function onSequenceBatches(uint64 newSequencedBatches, bytes32 newAccInputHash) returns(uint64)
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) OnSequenceBatches(newSequencedBatches uint64, newAccInputHash [32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.OnSequenceBatches(&_Polygonrollupmanager.TransactOpts, newSequencedBatches, newAccInputHash)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0x12b86e19.
//
// Solidity: function overridePendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) OverridePendingState(opts *bind.TransactOpts, rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "overridePendingState", rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0x12b86e19.
//
// Solidity: function overridePendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) OverridePendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.OverridePendingState(&_Polygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0x12b86e19.
//
// Solidity: function overridePendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) OverridePendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.OverridePendingState(&_Polygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x8bd4f071.
//
// Solidity: function proveNonDeterministicPendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) ProveNonDeterministicPendingState(opts *bind.TransactOpts, rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "proveNonDeterministicPendingState", rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x8bd4f071.
//
// Solidity: function proveNonDeterministicPendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) ProveNonDeterministicPendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.ProveNonDeterministicPendingState(&_Polygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x8bd4f071.
//
// Solidity: function proveNonDeterministicPendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) ProveNonDeterministicPendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.ProveNonDeterministicPendingState(&_Polygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.RenounceRole(&_Polygonrollupmanager.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.RenounceRole(&_Polygonrollupmanager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.RevokeRole(&_Polygonrollupmanager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.RevokeRole(&_Polygonrollupmanager.TransactOpts, role, account)
}

// SetBatchFee is a paid mutator transaction binding the contract method 0xd5073f6f.
//
// Solidity: function setBatchFee(uint256 newBatchFee) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) SetBatchFee(opts *bind.TransactOpts, newBatchFee *big.Int) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "setBatchFee", newBatchFee)
}

// SetBatchFee is a paid mutator transaction binding the contract method 0xd5073f6f.
//
// Solidity: function setBatchFee(uint256 newBatchFee) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) SetBatchFee(newBatchFee *big.Int) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.SetBatchFee(&_Polygonrollupmanager.TransactOpts, newBatchFee)
}

// SetBatchFee is a paid mutator transaction binding the contract method 0xd5073f6f.
//
// Solidity: function setBatchFee(uint256 newBatchFee) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) SetBatchFee(newBatchFee *big.Int) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.SetBatchFee(&_Polygonrollupmanager.TransactOpts, newBatchFee)
}

// SetMultiplierBatchFee is a paid mutator transaction binding the contract method 0x1816b7e5.
//
// Solidity: function setMultiplierBatchFee(uint16 newMultiplierBatchFee) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) SetMultiplierBatchFee(opts *bind.TransactOpts, newMultiplierBatchFee uint16) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "setMultiplierBatchFee", newMultiplierBatchFee)
}

// SetMultiplierBatchFee is a paid mutator transaction binding the contract method 0x1816b7e5.
//
// Solidity: function setMultiplierBatchFee(uint16 newMultiplierBatchFee) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) SetMultiplierBatchFee(newMultiplierBatchFee uint16) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.SetMultiplierBatchFee(&_Polygonrollupmanager.TransactOpts, newMultiplierBatchFee)
}

// SetMultiplierBatchFee is a paid mutator transaction binding the contract method 0x1816b7e5.
//
// Solidity: function setMultiplierBatchFee(uint16 newMultiplierBatchFee) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) SetMultiplierBatchFee(newMultiplierBatchFee uint16) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.SetMultiplierBatchFee(&_Polygonrollupmanager.TransactOpts, newMultiplierBatchFee)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) SetPendingStateTimeout(opts *bind.TransactOpts, newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "setPendingStateTimeout", newPendingStateTimeout)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) SetPendingStateTimeout(newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.SetPendingStateTimeout(&_Polygonrollupmanager.TransactOpts, newPendingStateTimeout)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) SetPendingStateTimeout(newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.SetPendingStateTimeout(&_Polygonrollupmanager.TransactOpts, newPendingStateTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) SetTrustedAggregatorTimeout(opts *bind.TransactOpts, newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "setTrustedAggregatorTimeout", newTrustedAggregatorTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) SetTrustedAggregatorTimeout(newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.SetTrustedAggregatorTimeout(&_Polygonrollupmanager.TransactOpts, newTrustedAggregatorTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) SetTrustedAggregatorTimeout(newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.SetTrustedAggregatorTimeout(&_Polygonrollupmanager.TransactOpts, newTrustedAggregatorTimeout)
}

// SetVerifyBatchTimeTarget is a paid mutator transaction binding the contract method 0xa066215c.
//
// Solidity: function setVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) SetVerifyBatchTimeTarget(opts *bind.TransactOpts, newVerifyBatchTimeTarget uint64) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "setVerifyBatchTimeTarget", newVerifyBatchTimeTarget)
}

// SetVerifyBatchTimeTarget is a paid mutator transaction binding the contract method 0xa066215c.
//
// Solidity: function setVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) SetVerifyBatchTimeTarget(newVerifyBatchTimeTarget uint64) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.SetVerifyBatchTimeTarget(&_Polygonrollupmanager.TransactOpts, newVerifyBatchTimeTarget)
}

// SetVerifyBatchTimeTarget is a paid mutator transaction binding the contract method 0xa066215c.
//
// Solidity: function setVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) SetVerifyBatchTimeTarget(newVerifyBatchTimeTarget uint64) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.SetVerifyBatchTimeTarget(&_Polygonrollupmanager.TransactOpts, newVerifyBatchTimeTarget)
}

// UpdateRollup is a paid mutator transaction binding the contract method 0xc4c928c2.
//
// Solidity: function updateRollup(address rollupContract, uint32 newRollupTypeID, bytes upgradeData) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) UpdateRollup(opts *bind.TransactOpts, rollupContract common.Address, newRollupTypeID uint32, upgradeData []byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "updateRollup", rollupContract, newRollupTypeID, upgradeData)
}

// UpdateRollup is a paid mutator transaction binding the contract method 0xc4c928c2.
//
// Solidity: function updateRollup(address rollupContract, uint32 newRollupTypeID, bytes upgradeData) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) UpdateRollup(rollupContract common.Address, newRollupTypeID uint32, upgradeData []byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.UpdateRollup(&_Polygonrollupmanager.TransactOpts, rollupContract, newRollupTypeID, upgradeData)
}

// UpdateRollup is a paid mutator transaction binding the contract method 0xc4c928c2.
//
// Solidity: function updateRollup(address rollupContract, uint32 newRollupTypeID, bytes upgradeData) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) UpdateRollup(rollupContract common.Address, newRollupTypeID uint32, upgradeData []byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.UpdateRollup(&_Polygonrollupmanager.TransactOpts, rollupContract, newRollupTypeID, upgradeData)
}

// VerifyBatches is a paid mutator transaction binding the contract method 0x87c20c01.
//
// Solidity: function verifyBatches(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) VerifyBatches(opts *bind.TransactOpts, rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "verifyBatches", rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatches is a paid mutator transaction binding the contract method 0x87c20c01.
//
// Solidity: function verifyBatches(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) VerifyBatches(rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.VerifyBatches(&_Polygonrollupmanager.TransactOpts, rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatches is a paid mutator transaction binding the contract method 0x87c20c01.
//
// Solidity: function verifyBatches(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) VerifyBatches(rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.VerifyBatches(&_Polygonrollupmanager.TransactOpts, rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatchesTrustedAggregator is a paid mutator transaction binding the contract method 0x1489ed10.
//
// Solidity: function verifyBatchesTrustedAggregator(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactor) VerifyBatchesTrustedAggregator(opts *bind.TransactOpts, rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.contract.Transact(opts, "verifyBatchesTrustedAggregator", rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatchesTrustedAggregator is a paid mutator transaction binding the contract method 0x1489ed10.
//
// Solidity: function verifyBatchesTrustedAggregator(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerSession) VerifyBatchesTrustedAggregator(rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.VerifyBatchesTrustedAggregator(&_Polygonrollupmanager.TransactOpts, rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatchesTrustedAggregator is a paid mutator transaction binding the contract method 0x1489ed10.
//
// Solidity: function verifyBatchesTrustedAggregator(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Polygonrollupmanager *PolygonrollupmanagerTransactorSession) VerifyBatchesTrustedAggregator(rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Polygonrollupmanager.Contract.VerifyBatchesTrustedAggregator(&_Polygonrollupmanager.TransactOpts, rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// PolygonrollupmanagerAddExistingRollupIterator is returned from FilterAddExistingRollup and is used to iterate over the raw logs and unpacked data for AddExistingRollup events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerAddExistingRollupIterator struct {
	Event *PolygonrollupmanagerAddExistingRollup // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerAddExistingRollupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerAddExistingRollup)
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
		it.Event = new(PolygonrollupmanagerAddExistingRollup)
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
func (it *PolygonrollupmanagerAddExistingRollupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerAddExistingRollupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerAddExistingRollup represents a AddExistingRollup event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerAddExistingRollup struct {
	RollupID              uint32
	ForkID                uint64
	RollupAddress         common.Address
	ChainID               uint64
	RollupCompatibilityID uint8
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterAddExistingRollup is a free log retrieval operation binding the contract event 0x1c20c90adb9c53d1da3e79fb8531f2c13c52b73a31a2574202ec36fb449ab24a.
//
// Solidity: event AddExistingRollup(uint32 indexed rollupID, uint64 forkID, address rollupAddress, uint64 chainID, uint8 rollupCompatibilityID)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterAddExistingRollup(opts *bind.FilterOpts, rollupID []uint32) (*PolygonrollupmanagerAddExistingRollupIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "AddExistingRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerAddExistingRollupIterator{contract: _Polygonrollupmanager.contract, event: "AddExistingRollup", logs: logs, sub: sub}, nil
}

// WatchAddExistingRollup is a free log subscription operation binding the contract event 0x1c20c90adb9c53d1da3e79fb8531f2c13c52b73a31a2574202ec36fb449ab24a.
//
// Solidity: event AddExistingRollup(uint32 indexed rollupID, uint64 forkID, address rollupAddress, uint64 chainID, uint8 rollupCompatibilityID)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchAddExistingRollup(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerAddExistingRollup, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "AddExistingRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerAddExistingRollup)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "AddExistingRollup", log); err != nil {
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

// ParseAddExistingRollup is a log parse operation binding the contract event 0x1c20c90adb9c53d1da3e79fb8531f2c13c52b73a31a2574202ec36fb449ab24a.
//
// Solidity: event AddExistingRollup(uint32 indexed rollupID, uint64 forkID, address rollupAddress, uint64 chainID, uint8 rollupCompatibilityID)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseAddExistingRollup(log types.Log) (*PolygonrollupmanagerAddExistingRollup, error) {
	event := new(PolygonrollupmanagerAddExistingRollup)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "AddExistingRollup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerAddNewRollupTypeIterator is returned from FilterAddNewRollupType and is used to iterate over the raw logs and unpacked data for AddNewRollupType events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerAddNewRollupTypeIterator struct {
	Event *PolygonrollupmanagerAddNewRollupType // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerAddNewRollupTypeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerAddNewRollupType)
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
		it.Event = new(PolygonrollupmanagerAddNewRollupType)
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
func (it *PolygonrollupmanagerAddNewRollupTypeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerAddNewRollupTypeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerAddNewRollupType represents a AddNewRollupType event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerAddNewRollupType struct {
	RollupTypeID            uint32
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Genesis                 [32]byte
	Description             string
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterAddNewRollupType is a free log retrieval operation binding the contract event 0xa2970448b3bd66ba7e524e7b2a5b9cf94fa29e32488fb942afdfe70dd4b77b52.
//
// Solidity: event AddNewRollupType(uint32 indexed rollupTypeID, address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterAddNewRollupType(opts *bind.FilterOpts, rollupTypeID []uint32) (*PolygonrollupmanagerAddNewRollupTypeIterator, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "AddNewRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerAddNewRollupTypeIterator{contract: _Polygonrollupmanager.contract, event: "AddNewRollupType", logs: logs, sub: sub}, nil
}

// WatchAddNewRollupType is a free log subscription operation binding the contract event 0xa2970448b3bd66ba7e524e7b2a5b9cf94fa29e32488fb942afdfe70dd4b77b52.
//
// Solidity: event AddNewRollupType(uint32 indexed rollupTypeID, address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchAddNewRollupType(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerAddNewRollupType, rollupTypeID []uint32) (event.Subscription, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "AddNewRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerAddNewRollupType)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "AddNewRollupType", log); err != nil {
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

// ParseAddNewRollupType is a log parse operation binding the contract event 0xa2970448b3bd66ba7e524e7b2a5b9cf94fa29e32488fb942afdfe70dd4b77b52.
//
// Solidity: event AddNewRollupType(uint32 indexed rollupTypeID, address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseAddNewRollupType(log types.Log) (*PolygonrollupmanagerAddNewRollupType, error) {
	event := new(PolygonrollupmanagerAddNewRollupType)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "AddNewRollupType", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerConsolidatePendingStateIterator is returned from FilterConsolidatePendingState and is used to iterate over the raw logs and unpacked data for ConsolidatePendingState events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerConsolidatePendingStateIterator struct {
	Event *PolygonrollupmanagerConsolidatePendingState // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerConsolidatePendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerConsolidatePendingState)
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
		it.Event = new(PolygonrollupmanagerConsolidatePendingState)
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
func (it *PolygonrollupmanagerConsolidatePendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerConsolidatePendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerConsolidatePendingState represents a ConsolidatePendingState event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerConsolidatePendingState struct {
	RollupID        uint32
	NumBatch        uint64
	StateRoot       [32]byte
	ExitRoot        [32]byte
	PendingStateNum uint64
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConsolidatePendingState is a free log retrieval operation binding the contract event 0x581910eb7a27738945c2f00a91f2284b2d6de9d4e472b12f901c2b0df045e21b.
//
// Solidity: event ConsolidatePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, uint64 pendingStateNum)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterConsolidatePendingState(opts *bind.FilterOpts, rollupID []uint32) (*PolygonrollupmanagerConsolidatePendingStateIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "ConsolidatePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerConsolidatePendingStateIterator{contract: _Polygonrollupmanager.contract, event: "ConsolidatePendingState", logs: logs, sub: sub}, nil
}

// WatchConsolidatePendingState is a free log subscription operation binding the contract event 0x581910eb7a27738945c2f00a91f2284b2d6de9d4e472b12f901c2b0df045e21b.
//
// Solidity: event ConsolidatePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, uint64 pendingStateNum)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchConsolidatePendingState(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerConsolidatePendingState, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "ConsolidatePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerConsolidatePendingState)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "ConsolidatePendingState", log); err != nil {
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

// ParseConsolidatePendingState is a log parse operation binding the contract event 0x581910eb7a27738945c2f00a91f2284b2d6de9d4e472b12f901c2b0df045e21b.
//
// Solidity: event ConsolidatePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, uint64 pendingStateNum)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseConsolidatePendingState(log types.Log) (*PolygonrollupmanagerConsolidatePendingState, error) {
	event := new(PolygonrollupmanagerConsolidatePendingState)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "ConsolidatePendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerCreateNewRollupIterator is returned from FilterCreateNewRollup and is used to iterate over the raw logs and unpacked data for CreateNewRollup events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerCreateNewRollupIterator struct {
	Event *PolygonrollupmanagerCreateNewRollup // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerCreateNewRollupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerCreateNewRollup)
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
		it.Event = new(PolygonrollupmanagerCreateNewRollup)
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
func (it *PolygonrollupmanagerCreateNewRollupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerCreateNewRollupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerCreateNewRollup represents a CreateNewRollup event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerCreateNewRollup struct {
	RollupID        uint32
	RollupTypeID    uint32
	RollupAddress   common.Address
	ChainID         uint64
	GasTokenAddress common.Address
	GasTokenNetwork uint32
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCreateNewRollup is a free log retrieval operation binding the contract event 0x9a84992ead4481f37097c7663e97ed883b31f216fae5685c6ee576a656483cd9.
//
// Solidity: event CreateNewRollup(uint32 indexed rollupID, uint32 rollupTypeID, address rollupAddress, uint64 chainID, address gasTokenAddress, uint32 gasTokenNetwork)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterCreateNewRollup(opts *bind.FilterOpts, rollupID []uint32) (*PolygonrollupmanagerCreateNewRollupIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "CreateNewRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerCreateNewRollupIterator{contract: _Polygonrollupmanager.contract, event: "CreateNewRollup", logs: logs, sub: sub}, nil
}

// WatchCreateNewRollup is a free log subscription operation binding the contract event 0x9a84992ead4481f37097c7663e97ed883b31f216fae5685c6ee576a656483cd9.
//
// Solidity: event CreateNewRollup(uint32 indexed rollupID, uint32 rollupTypeID, address rollupAddress, uint64 chainID, address gasTokenAddress, uint32 gasTokenNetwork)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchCreateNewRollup(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerCreateNewRollup, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "CreateNewRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerCreateNewRollup)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "CreateNewRollup", log); err != nil {
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

// ParseCreateNewRollup is a log parse operation binding the contract event 0x9a84992ead4481f37097c7663e97ed883b31f216fae5685c6ee576a656483cd9.
//
// Solidity: event CreateNewRollup(uint32 indexed rollupID, uint32 rollupTypeID, address rollupAddress, uint64 chainID, address gasTokenAddress, uint32 gasTokenNetwork)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseCreateNewRollup(log types.Log) (*PolygonrollupmanagerCreateNewRollup, error) {
	event := new(PolygonrollupmanagerCreateNewRollup)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "CreateNewRollup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerEmergencyStateActivatedIterator is returned from FilterEmergencyStateActivated and is used to iterate over the raw logs and unpacked data for EmergencyStateActivated events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerEmergencyStateActivatedIterator struct {
	Event *PolygonrollupmanagerEmergencyStateActivated // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerEmergencyStateActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerEmergencyStateActivated)
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
		it.Event = new(PolygonrollupmanagerEmergencyStateActivated)
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
func (it *PolygonrollupmanagerEmergencyStateActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerEmergencyStateActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerEmergencyStateActivated represents a EmergencyStateActivated event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerEmergencyStateActivated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEmergencyStateActivated is a free log retrieval operation binding the contract event 0x2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497.
//
// Solidity: event EmergencyStateActivated()
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterEmergencyStateActivated(opts *bind.FilterOpts) (*PolygonrollupmanagerEmergencyStateActivatedIterator, error) {

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "EmergencyStateActivated")
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerEmergencyStateActivatedIterator{contract: _Polygonrollupmanager.contract, event: "EmergencyStateActivated", logs: logs, sub: sub}, nil
}

// WatchEmergencyStateActivated is a free log subscription operation binding the contract event 0x2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497.
//
// Solidity: event EmergencyStateActivated()
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchEmergencyStateActivated(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerEmergencyStateActivated) (event.Subscription, error) {

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "EmergencyStateActivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerEmergencyStateActivated)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "EmergencyStateActivated", log); err != nil {
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
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseEmergencyStateActivated(log types.Log) (*PolygonrollupmanagerEmergencyStateActivated, error) {
	event := new(PolygonrollupmanagerEmergencyStateActivated)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "EmergencyStateActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerEmergencyStateDeactivatedIterator is returned from FilterEmergencyStateDeactivated and is used to iterate over the raw logs and unpacked data for EmergencyStateDeactivated events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerEmergencyStateDeactivatedIterator struct {
	Event *PolygonrollupmanagerEmergencyStateDeactivated // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerEmergencyStateDeactivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerEmergencyStateDeactivated)
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
		it.Event = new(PolygonrollupmanagerEmergencyStateDeactivated)
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
func (it *PolygonrollupmanagerEmergencyStateDeactivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerEmergencyStateDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerEmergencyStateDeactivated represents a EmergencyStateDeactivated event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerEmergencyStateDeactivated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEmergencyStateDeactivated is a free log retrieval operation binding the contract event 0x1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3.
//
// Solidity: event EmergencyStateDeactivated()
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterEmergencyStateDeactivated(opts *bind.FilterOpts) (*PolygonrollupmanagerEmergencyStateDeactivatedIterator, error) {

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "EmergencyStateDeactivated")
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerEmergencyStateDeactivatedIterator{contract: _Polygonrollupmanager.contract, event: "EmergencyStateDeactivated", logs: logs, sub: sub}, nil
}

// WatchEmergencyStateDeactivated is a free log subscription operation binding the contract event 0x1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3.
//
// Solidity: event EmergencyStateDeactivated()
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchEmergencyStateDeactivated(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerEmergencyStateDeactivated) (event.Subscription, error) {

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "EmergencyStateDeactivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerEmergencyStateDeactivated)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "EmergencyStateDeactivated", log); err != nil {
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
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseEmergencyStateDeactivated(log types.Log) (*PolygonrollupmanagerEmergencyStateDeactivated, error) {
	event := new(PolygonrollupmanagerEmergencyStateDeactivated)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "EmergencyStateDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerInitializedIterator struct {
	Event *PolygonrollupmanagerInitialized // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerInitialized)
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
		it.Event = new(PolygonrollupmanagerInitialized)
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
func (it *PolygonrollupmanagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerInitialized represents a Initialized event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*PolygonrollupmanagerInitializedIterator, error) {

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerInitializedIterator{contract: _Polygonrollupmanager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerInitialized) (event.Subscription, error) {

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerInitialized)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseInitialized(log types.Log) (*PolygonrollupmanagerInitialized, error) {
	event := new(PolygonrollupmanagerInitialized)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerObsoleteRollupTypeIterator is returned from FilterObsoleteRollupType and is used to iterate over the raw logs and unpacked data for ObsoleteRollupType events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerObsoleteRollupTypeIterator struct {
	Event *PolygonrollupmanagerObsoleteRollupType // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerObsoleteRollupTypeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerObsoleteRollupType)
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
		it.Event = new(PolygonrollupmanagerObsoleteRollupType)
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
func (it *PolygonrollupmanagerObsoleteRollupTypeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerObsoleteRollupTypeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerObsoleteRollupType represents a ObsoleteRollupType event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerObsoleteRollupType struct {
	RollupTypeID uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterObsoleteRollupType is a free log retrieval operation binding the contract event 0x4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e44.
//
// Solidity: event ObsoleteRollupType(uint32 indexed rollupTypeID)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterObsoleteRollupType(opts *bind.FilterOpts, rollupTypeID []uint32) (*PolygonrollupmanagerObsoleteRollupTypeIterator, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "ObsoleteRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerObsoleteRollupTypeIterator{contract: _Polygonrollupmanager.contract, event: "ObsoleteRollupType", logs: logs, sub: sub}, nil
}

// WatchObsoleteRollupType is a free log subscription operation binding the contract event 0x4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e44.
//
// Solidity: event ObsoleteRollupType(uint32 indexed rollupTypeID)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchObsoleteRollupType(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerObsoleteRollupType, rollupTypeID []uint32) (event.Subscription, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "ObsoleteRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerObsoleteRollupType)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "ObsoleteRollupType", log); err != nil {
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

// ParseObsoleteRollupType is a log parse operation binding the contract event 0x4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e44.
//
// Solidity: event ObsoleteRollupType(uint32 indexed rollupTypeID)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseObsoleteRollupType(log types.Log) (*PolygonrollupmanagerObsoleteRollupType, error) {
	event := new(PolygonrollupmanagerObsoleteRollupType)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "ObsoleteRollupType", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerOnSequenceBatchesIterator is returned from FilterOnSequenceBatches and is used to iterate over the raw logs and unpacked data for OnSequenceBatches events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerOnSequenceBatchesIterator struct {
	Event *PolygonrollupmanagerOnSequenceBatches // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerOnSequenceBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerOnSequenceBatches)
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
		it.Event = new(PolygonrollupmanagerOnSequenceBatches)
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
func (it *PolygonrollupmanagerOnSequenceBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerOnSequenceBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerOnSequenceBatches represents a OnSequenceBatches event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerOnSequenceBatches struct {
	RollupID           uint32
	LastBatchSequenced uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterOnSequenceBatches is a free log retrieval operation binding the contract event 0x1d9f30260051d51d70339da239ea7b080021adcaabfa71c9b0ea339a20cf9a25.
//
// Solidity: event OnSequenceBatches(uint32 indexed rollupID, uint64 lastBatchSequenced)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterOnSequenceBatches(opts *bind.FilterOpts, rollupID []uint32) (*PolygonrollupmanagerOnSequenceBatchesIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "OnSequenceBatches", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerOnSequenceBatchesIterator{contract: _Polygonrollupmanager.contract, event: "OnSequenceBatches", logs: logs, sub: sub}, nil
}

// WatchOnSequenceBatches is a free log subscription operation binding the contract event 0x1d9f30260051d51d70339da239ea7b080021adcaabfa71c9b0ea339a20cf9a25.
//
// Solidity: event OnSequenceBatches(uint32 indexed rollupID, uint64 lastBatchSequenced)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchOnSequenceBatches(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerOnSequenceBatches, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "OnSequenceBatches", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerOnSequenceBatches)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "OnSequenceBatches", log); err != nil {
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

// ParseOnSequenceBatches is a log parse operation binding the contract event 0x1d9f30260051d51d70339da239ea7b080021adcaabfa71c9b0ea339a20cf9a25.
//
// Solidity: event OnSequenceBatches(uint32 indexed rollupID, uint64 lastBatchSequenced)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseOnSequenceBatches(log types.Log) (*PolygonrollupmanagerOnSequenceBatches, error) {
	event := new(PolygonrollupmanagerOnSequenceBatches)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "OnSequenceBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerOverridePendingStateIterator is returned from FilterOverridePendingState and is used to iterate over the raw logs and unpacked data for OverridePendingState events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerOverridePendingStateIterator struct {
	Event *PolygonrollupmanagerOverridePendingState // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerOverridePendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerOverridePendingState)
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
		it.Event = new(PolygonrollupmanagerOverridePendingState)
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
func (it *PolygonrollupmanagerOverridePendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerOverridePendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerOverridePendingState represents a OverridePendingState event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerOverridePendingState struct {
	RollupID   uint32
	NumBatch   uint64
	StateRoot  [32]byte
	ExitRoot   [32]byte
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterOverridePendingState is a free log retrieval operation binding the contract event 0x3182bd6e6f74fc1fdc88b60f3a4f4c7f79db6ae6f5b88a1b3f5a1e28ec210d5e.
//
// Solidity: event OverridePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address aggregator)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterOverridePendingState(opts *bind.FilterOpts, rollupID []uint32) (*PolygonrollupmanagerOverridePendingStateIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "OverridePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerOverridePendingStateIterator{contract: _Polygonrollupmanager.contract, event: "OverridePendingState", logs: logs, sub: sub}, nil
}

// WatchOverridePendingState is a free log subscription operation binding the contract event 0x3182bd6e6f74fc1fdc88b60f3a4f4c7f79db6ae6f5b88a1b3f5a1e28ec210d5e.
//
// Solidity: event OverridePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address aggregator)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchOverridePendingState(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerOverridePendingState, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "OverridePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerOverridePendingState)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "OverridePendingState", log); err != nil {
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

// ParseOverridePendingState is a log parse operation binding the contract event 0x3182bd6e6f74fc1fdc88b60f3a4f4c7f79db6ae6f5b88a1b3f5a1e28ec210d5e.
//
// Solidity: event OverridePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address aggregator)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseOverridePendingState(log types.Log) (*PolygonrollupmanagerOverridePendingState, error) {
	event := new(PolygonrollupmanagerOverridePendingState)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "OverridePendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerProveNonDeterministicPendingStateIterator is returned from FilterProveNonDeterministicPendingState and is used to iterate over the raw logs and unpacked data for ProveNonDeterministicPendingState events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerProveNonDeterministicPendingStateIterator struct {
	Event *PolygonrollupmanagerProveNonDeterministicPendingState // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerProveNonDeterministicPendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerProveNonDeterministicPendingState)
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
		it.Event = new(PolygonrollupmanagerProveNonDeterministicPendingState)
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
func (it *PolygonrollupmanagerProveNonDeterministicPendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerProveNonDeterministicPendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerProveNonDeterministicPendingState represents a ProveNonDeterministicPendingState event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerProveNonDeterministicPendingState struct {
	StoredStateRoot [32]byte
	ProvedStateRoot [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterProveNonDeterministicPendingState is a free log retrieval operation binding the contract event 0x1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010.
//
// Solidity: event ProveNonDeterministicPendingState(bytes32 storedStateRoot, bytes32 provedStateRoot)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterProveNonDeterministicPendingState(opts *bind.FilterOpts) (*PolygonrollupmanagerProveNonDeterministicPendingStateIterator, error) {

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "ProveNonDeterministicPendingState")
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerProveNonDeterministicPendingStateIterator{contract: _Polygonrollupmanager.contract, event: "ProveNonDeterministicPendingState", logs: logs, sub: sub}, nil
}

// WatchProveNonDeterministicPendingState is a free log subscription operation binding the contract event 0x1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010.
//
// Solidity: event ProveNonDeterministicPendingState(bytes32 storedStateRoot, bytes32 provedStateRoot)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchProveNonDeterministicPendingState(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerProveNonDeterministicPendingState) (event.Subscription, error) {

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "ProveNonDeterministicPendingState")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerProveNonDeterministicPendingState)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "ProveNonDeterministicPendingState", log); err != nil {
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
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseProveNonDeterministicPendingState(log types.Log) (*PolygonrollupmanagerProveNonDeterministicPendingState, error) {
	event := new(PolygonrollupmanagerProveNonDeterministicPendingState)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "ProveNonDeterministicPendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerRoleAdminChangedIterator struct {
	Event *PolygonrollupmanagerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerRoleAdminChanged)
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
		it.Event = new(PolygonrollupmanagerRoleAdminChanged)
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
func (it *PolygonrollupmanagerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerRoleAdminChanged represents a RoleAdminChanged event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*PolygonrollupmanagerRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerRoleAdminChangedIterator{contract: _Polygonrollupmanager.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerRoleAdminChanged)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseRoleAdminChanged(log types.Log) (*PolygonrollupmanagerRoleAdminChanged, error) {
	event := new(PolygonrollupmanagerRoleAdminChanged)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerRoleGrantedIterator struct {
	Event *PolygonrollupmanagerRoleGranted // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerRoleGranted)
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
		it.Event = new(PolygonrollupmanagerRoleGranted)
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
func (it *PolygonrollupmanagerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerRoleGranted represents a RoleGranted event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*PolygonrollupmanagerRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerRoleGrantedIterator{contract: _Polygonrollupmanager.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerRoleGranted)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseRoleGranted(log types.Log) (*PolygonrollupmanagerRoleGranted, error) {
	event := new(PolygonrollupmanagerRoleGranted)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerRoleRevokedIterator struct {
	Event *PolygonrollupmanagerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerRoleRevoked)
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
		it.Event = new(PolygonrollupmanagerRoleRevoked)
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
func (it *PolygonrollupmanagerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerRoleRevoked represents a RoleRevoked event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*PolygonrollupmanagerRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerRoleRevokedIterator{contract: _Polygonrollupmanager.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerRoleRevoked)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseRoleRevoked(log types.Log) (*PolygonrollupmanagerRoleRevoked, error) {
	event := new(PolygonrollupmanagerRoleRevoked)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerSetBatchFeeIterator is returned from FilterSetBatchFee and is used to iterate over the raw logs and unpacked data for SetBatchFee events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerSetBatchFeeIterator struct {
	Event *PolygonrollupmanagerSetBatchFee // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerSetBatchFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerSetBatchFee)
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
		it.Event = new(PolygonrollupmanagerSetBatchFee)
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
func (it *PolygonrollupmanagerSetBatchFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerSetBatchFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerSetBatchFee represents a SetBatchFee event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerSetBatchFee struct {
	NewBatchFee *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSetBatchFee is a free log retrieval operation binding the contract event 0xfb383653f53ee079978d0c9aff7aeff04a10166ce244cca9c9f9d8d96bed45b2.
//
// Solidity: event SetBatchFee(uint256 newBatchFee)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterSetBatchFee(opts *bind.FilterOpts) (*PolygonrollupmanagerSetBatchFeeIterator, error) {

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "SetBatchFee")
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerSetBatchFeeIterator{contract: _Polygonrollupmanager.contract, event: "SetBatchFee", logs: logs, sub: sub}, nil
}

// WatchSetBatchFee is a free log subscription operation binding the contract event 0xfb383653f53ee079978d0c9aff7aeff04a10166ce244cca9c9f9d8d96bed45b2.
//
// Solidity: event SetBatchFee(uint256 newBatchFee)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchSetBatchFee(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerSetBatchFee) (event.Subscription, error) {

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "SetBatchFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerSetBatchFee)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "SetBatchFee", log); err != nil {
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

// ParseSetBatchFee is a log parse operation binding the contract event 0xfb383653f53ee079978d0c9aff7aeff04a10166ce244cca9c9f9d8d96bed45b2.
//
// Solidity: event SetBatchFee(uint256 newBatchFee)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseSetBatchFee(log types.Log) (*PolygonrollupmanagerSetBatchFee, error) {
	event := new(PolygonrollupmanagerSetBatchFee)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "SetBatchFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerSetMultiplierBatchFeeIterator is returned from FilterSetMultiplierBatchFee and is used to iterate over the raw logs and unpacked data for SetMultiplierBatchFee events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerSetMultiplierBatchFeeIterator struct {
	Event *PolygonrollupmanagerSetMultiplierBatchFee // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerSetMultiplierBatchFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerSetMultiplierBatchFee)
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
		it.Event = new(PolygonrollupmanagerSetMultiplierBatchFee)
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
func (it *PolygonrollupmanagerSetMultiplierBatchFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerSetMultiplierBatchFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerSetMultiplierBatchFee represents a SetMultiplierBatchFee event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerSetMultiplierBatchFee struct {
	NewMultiplierBatchFee uint16
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterSetMultiplierBatchFee is a free log retrieval operation binding the contract event 0x7019933d795eba185c180209e8ae8bffbaa25bcef293364687702c31f4d302c5.
//
// Solidity: event SetMultiplierBatchFee(uint16 newMultiplierBatchFee)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterSetMultiplierBatchFee(opts *bind.FilterOpts) (*PolygonrollupmanagerSetMultiplierBatchFeeIterator, error) {

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "SetMultiplierBatchFee")
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerSetMultiplierBatchFeeIterator{contract: _Polygonrollupmanager.contract, event: "SetMultiplierBatchFee", logs: logs, sub: sub}, nil
}

// WatchSetMultiplierBatchFee is a free log subscription operation binding the contract event 0x7019933d795eba185c180209e8ae8bffbaa25bcef293364687702c31f4d302c5.
//
// Solidity: event SetMultiplierBatchFee(uint16 newMultiplierBatchFee)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchSetMultiplierBatchFee(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerSetMultiplierBatchFee) (event.Subscription, error) {

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "SetMultiplierBatchFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerSetMultiplierBatchFee)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "SetMultiplierBatchFee", log); err != nil {
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

// ParseSetMultiplierBatchFee is a log parse operation binding the contract event 0x7019933d795eba185c180209e8ae8bffbaa25bcef293364687702c31f4d302c5.
//
// Solidity: event SetMultiplierBatchFee(uint16 newMultiplierBatchFee)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseSetMultiplierBatchFee(log types.Log) (*PolygonrollupmanagerSetMultiplierBatchFee, error) {
	event := new(PolygonrollupmanagerSetMultiplierBatchFee)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "SetMultiplierBatchFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerSetPendingStateTimeoutIterator is returned from FilterSetPendingStateTimeout and is used to iterate over the raw logs and unpacked data for SetPendingStateTimeout events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerSetPendingStateTimeoutIterator struct {
	Event *PolygonrollupmanagerSetPendingStateTimeout // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerSetPendingStateTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerSetPendingStateTimeout)
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
		it.Event = new(PolygonrollupmanagerSetPendingStateTimeout)
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
func (it *PolygonrollupmanagerSetPendingStateTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerSetPendingStateTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerSetPendingStateTimeout represents a SetPendingStateTimeout event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerSetPendingStateTimeout struct {
	NewPendingStateTimeout uint64
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterSetPendingStateTimeout is a free log retrieval operation binding the contract event 0xc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c75.
//
// Solidity: event SetPendingStateTimeout(uint64 newPendingStateTimeout)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterSetPendingStateTimeout(opts *bind.FilterOpts) (*PolygonrollupmanagerSetPendingStateTimeoutIterator, error) {

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "SetPendingStateTimeout")
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerSetPendingStateTimeoutIterator{contract: _Polygonrollupmanager.contract, event: "SetPendingStateTimeout", logs: logs, sub: sub}, nil
}

// WatchSetPendingStateTimeout is a free log subscription operation binding the contract event 0xc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c75.
//
// Solidity: event SetPendingStateTimeout(uint64 newPendingStateTimeout)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchSetPendingStateTimeout(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerSetPendingStateTimeout) (event.Subscription, error) {

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "SetPendingStateTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerSetPendingStateTimeout)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "SetPendingStateTimeout", log); err != nil {
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
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseSetPendingStateTimeout(log types.Log) (*PolygonrollupmanagerSetPendingStateTimeout, error) {
	event := new(PolygonrollupmanagerSetPendingStateTimeout)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "SetPendingStateTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerSetTrustedAggregatorIterator is returned from FilterSetTrustedAggregator and is used to iterate over the raw logs and unpacked data for SetTrustedAggregator events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerSetTrustedAggregatorIterator struct {
	Event *PolygonrollupmanagerSetTrustedAggregator // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerSetTrustedAggregatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerSetTrustedAggregator)
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
		it.Event = new(PolygonrollupmanagerSetTrustedAggregator)
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
func (it *PolygonrollupmanagerSetTrustedAggregatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerSetTrustedAggregatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerSetTrustedAggregator represents a SetTrustedAggregator event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerSetTrustedAggregator struct {
	NewTrustedAggregator common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedAggregator is a free log retrieval operation binding the contract event 0x61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca.
//
// Solidity: event SetTrustedAggregator(address newTrustedAggregator)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterSetTrustedAggregator(opts *bind.FilterOpts) (*PolygonrollupmanagerSetTrustedAggregatorIterator, error) {

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "SetTrustedAggregator")
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerSetTrustedAggregatorIterator{contract: _Polygonrollupmanager.contract, event: "SetTrustedAggregator", logs: logs, sub: sub}, nil
}

// WatchSetTrustedAggregator is a free log subscription operation binding the contract event 0x61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca.
//
// Solidity: event SetTrustedAggregator(address newTrustedAggregator)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchSetTrustedAggregator(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerSetTrustedAggregator) (event.Subscription, error) {

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "SetTrustedAggregator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerSetTrustedAggregator)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregator", log); err != nil {
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
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseSetTrustedAggregator(log types.Log) (*PolygonrollupmanagerSetTrustedAggregator, error) {
	event := new(PolygonrollupmanagerSetTrustedAggregator)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerSetTrustedAggregatorTimeoutIterator is returned from FilterSetTrustedAggregatorTimeout and is used to iterate over the raw logs and unpacked data for SetTrustedAggregatorTimeout events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerSetTrustedAggregatorTimeoutIterator struct {
	Event *PolygonrollupmanagerSetTrustedAggregatorTimeout // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerSetTrustedAggregatorTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerSetTrustedAggregatorTimeout)
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
		it.Event = new(PolygonrollupmanagerSetTrustedAggregatorTimeout)
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
func (it *PolygonrollupmanagerSetTrustedAggregatorTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerSetTrustedAggregatorTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerSetTrustedAggregatorTimeout represents a SetTrustedAggregatorTimeout event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerSetTrustedAggregatorTimeout struct {
	NewTrustedAggregatorTimeout uint64
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedAggregatorTimeout is a free log retrieval operation binding the contract event 0x1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1.
//
// Solidity: event SetTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterSetTrustedAggregatorTimeout(opts *bind.FilterOpts) (*PolygonrollupmanagerSetTrustedAggregatorTimeoutIterator, error) {

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "SetTrustedAggregatorTimeout")
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerSetTrustedAggregatorTimeoutIterator{contract: _Polygonrollupmanager.contract, event: "SetTrustedAggregatorTimeout", logs: logs, sub: sub}, nil
}

// WatchSetTrustedAggregatorTimeout is a free log subscription operation binding the contract event 0x1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1.
//
// Solidity: event SetTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchSetTrustedAggregatorTimeout(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerSetTrustedAggregatorTimeout) (event.Subscription, error) {

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "SetTrustedAggregatorTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerSetTrustedAggregatorTimeout)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregatorTimeout", log); err != nil {
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
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseSetTrustedAggregatorTimeout(log types.Log) (*PolygonrollupmanagerSetTrustedAggregatorTimeout, error) {
	event := new(PolygonrollupmanagerSetTrustedAggregatorTimeout)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregatorTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerSetVerifyBatchTimeTargetIterator is returned from FilterSetVerifyBatchTimeTarget and is used to iterate over the raw logs and unpacked data for SetVerifyBatchTimeTarget events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerSetVerifyBatchTimeTargetIterator struct {
	Event *PolygonrollupmanagerSetVerifyBatchTimeTarget // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerSetVerifyBatchTimeTargetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerSetVerifyBatchTimeTarget)
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
		it.Event = new(PolygonrollupmanagerSetVerifyBatchTimeTarget)
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
func (it *PolygonrollupmanagerSetVerifyBatchTimeTargetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerSetVerifyBatchTimeTargetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerSetVerifyBatchTimeTarget represents a SetVerifyBatchTimeTarget event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerSetVerifyBatchTimeTarget struct {
	NewVerifyBatchTimeTarget uint64
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterSetVerifyBatchTimeTarget is a free log retrieval operation binding the contract event 0x1b023231a1ab6b5d93992f168fb44498e1a7e64cef58daff6f1c216de6a68c28.
//
// Solidity: event SetVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterSetVerifyBatchTimeTarget(opts *bind.FilterOpts) (*PolygonrollupmanagerSetVerifyBatchTimeTargetIterator, error) {

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "SetVerifyBatchTimeTarget")
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerSetVerifyBatchTimeTargetIterator{contract: _Polygonrollupmanager.contract, event: "SetVerifyBatchTimeTarget", logs: logs, sub: sub}, nil
}

// WatchSetVerifyBatchTimeTarget is a free log subscription operation binding the contract event 0x1b023231a1ab6b5d93992f168fb44498e1a7e64cef58daff6f1c216de6a68c28.
//
// Solidity: event SetVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchSetVerifyBatchTimeTarget(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerSetVerifyBatchTimeTarget) (event.Subscription, error) {

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "SetVerifyBatchTimeTarget")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerSetVerifyBatchTimeTarget)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "SetVerifyBatchTimeTarget", log); err != nil {
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

// ParseSetVerifyBatchTimeTarget is a log parse operation binding the contract event 0x1b023231a1ab6b5d93992f168fb44498e1a7e64cef58daff6f1c216de6a68c28.
//
// Solidity: event SetVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseSetVerifyBatchTimeTarget(log types.Log) (*PolygonrollupmanagerSetVerifyBatchTimeTarget, error) {
	event := new(PolygonrollupmanagerSetVerifyBatchTimeTarget)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "SetVerifyBatchTimeTarget", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerUpdateRollupIterator is returned from FilterUpdateRollup and is used to iterate over the raw logs and unpacked data for UpdateRollup events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerUpdateRollupIterator struct {
	Event *PolygonrollupmanagerUpdateRollup // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerUpdateRollupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerUpdateRollup)
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
		it.Event = new(PolygonrollupmanagerUpdateRollup)
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
func (it *PolygonrollupmanagerUpdateRollupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerUpdateRollupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerUpdateRollup represents a UpdateRollup event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerUpdateRollup struct {
	RollupID                       uint32
	NewRollupTypeID                uint32
	LastVerifiedBatchBeforeUpgrade uint64
	Raw                            types.Log // Blockchain specific contextual infos
}

// FilterUpdateRollup is a free log retrieval operation binding the contract event 0xf585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d.
//
// Solidity: event UpdateRollup(uint32 indexed rollupID, uint32 newRollupTypeID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterUpdateRollup(opts *bind.FilterOpts, rollupID []uint32) (*PolygonrollupmanagerUpdateRollupIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "UpdateRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerUpdateRollupIterator{contract: _Polygonrollupmanager.contract, event: "UpdateRollup", logs: logs, sub: sub}, nil
}

// WatchUpdateRollup is a free log subscription operation binding the contract event 0xf585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d.
//
// Solidity: event UpdateRollup(uint32 indexed rollupID, uint32 newRollupTypeID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchUpdateRollup(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerUpdateRollup, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "UpdateRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerUpdateRollup)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "UpdateRollup", log); err != nil {
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

// ParseUpdateRollup is a log parse operation binding the contract event 0xf585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d.
//
// Solidity: event UpdateRollup(uint32 indexed rollupID, uint32 newRollupTypeID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseUpdateRollup(log types.Log) (*PolygonrollupmanagerUpdateRollup, error) {
	event := new(PolygonrollupmanagerUpdateRollup)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "UpdateRollup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerVerifyBatchesIterator is returned from FilterVerifyBatches and is used to iterate over the raw logs and unpacked data for VerifyBatches events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerVerifyBatchesIterator struct {
	Event *PolygonrollupmanagerVerifyBatches // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerVerifyBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerVerifyBatches)
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
		it.Event = new(PolygonrollupmanagerVerifyBatches)
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
func (it *PolygonrollupmanagerVerifyBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerVerifyBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerVerifyBatches represents a VerifyBatches event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerVerifyBatches struct {
	RollupID   uint32
	NumBatch   uint64
	StateRoot  [32]byte
	ExitRoot   [32]byte
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVerifyBatches is a free log retrieval operation binding the contract event 0xaac1e7a157b259544ebacd6e8a82ae5d6c8f174e12aa48696277bcc9a661f0b4.
//
// Solidity: event VerifyBatches(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterVerifyBatches(opts *bind.FilterOpts, rollupID []uint32, aggregator []common.Address) (*PolygonrollupmanagerVerifyBatchesIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "VerifyBatches", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerVerifyBatchesIterator{contract: _Polygonrollupmanager.contract, event: "VerifyBatches", logs: logs, sub: sub}, nil
}

// WatchVerifyBatches is a free log subscription operation binding the contract event 0xaac1e7a157b259544ebacd6e8a82ae5d6c8f174e12aa48696277bcc9a661f0b4.
//
// Solidity: event VerifyBatches(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchVerifyBatches(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerVerifyBatches, rollupID []uint32, aggregator []common.Address) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "VerifyBatches", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerVerifyBatches)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "VerifyBatches", log); err != nil {
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

// ParseVerifyBatches is a log parse operation binding the contract event 0xaac1e7a157b259544ebacd6e8a82ae5d6c8f174e12aa48696277bcc9a661f0b4.
//
// Solidity: event VerifyBatches(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseVerifyBatches(log types.Log) (*PolygonrollupmanagerVerifyBatches, error) {
	event := new(PolygonrollupmanagerVerifyBatches)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "VerifyBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator is returned from FilterVerifyBatchesTrustedAggregator and is used to iterate over the raw logs and unpacked data for VerifyBatchesTrustedAggregator events raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator struct {
	Event *PolygonrollupmanagerVerifyBatchesTrustedAggregator // Event containing the contract specifics and raw log

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
func (it *PolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygonrollupmanagerVerifyBatchesTrustedAggregator)
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
		it.Event = new(PolygonrollupmanagerVerifyBatchesTrustedAggregator)
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
func (it *PolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygonrollupmanagerVerifyBatchesTrustedAggregator represents a VerifyBatchesTrustedAggregator event raised by the Polygonrollupmanager contract.
type PolygonrollupmanagerVerifyBatchesTrustedAggregator struct {
	RollupID   uint32
	NumBatch   uint64
	StateRoot  [32]byte
	ExitRoot   [32]byte
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVerifyBatchesTrustedAggregator is a free log retrieval operation binding the contract event 0xd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d3.
//
// Solidity: event VerifyBatchesTrustedAggregator(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) FilterVerifyBatchesTrustedAggregator(opts *bind.FilterOpts, rollupID []uint32, aggregator []common.Address) (*PolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.FilterLogs(opts, "VerifyBatchesTrustedAggregator", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &PolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator{contract: _Polygonrollupmanager.contract, event: "VerifyBatchesTrustedAggregator", logs: logs, sub: sub}, nil
}

// WatchVerifyBatchesTrustedAggregator is a free log subscription operation binding the contract event 0xd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d3.
//
// Solidity: event VerifyBatchesTrustedAggregator(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) WatchVerifyBatchesTrustedAggregator(opts *bind.WatchOpts, sink chan<- *PolygonrollupmanagerVerifyBatchesTrustedAggregator, rollupID []uint32, aggregator []common.Address) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Polygonrollupmanager.contract.WatchLogs(opts, "VerifyBatchesTrustedAggregator", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygonrollupmanagerVerifyBatchesTrustedAggregator)
				if err := _Polygonrollupmanager.contract.UnpackLog(event, "VerifyBatchesTrustedAggregator", log); err != nil {
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

// ParseVerifyBatchesTrustedAggregator is a log parse operation binding the contract event 0xd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d3.
//
// Solidity: event VerifyBatchesTrustedAggregator(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Polygonrollupmanager *PolygonrollupmanagerFilterer) ParseVerifyBatchesTrustedAggregator(log types.Log) (*PolygonrollupmanagerVerifyBatchesTrustedAggregator, error) {
	event := new(PolygonrollupmanagerVerifyBatchesTrustedAggregator)
	if err := _Polygonrollupmanager.contract.UnpackLog(event, "VerifyBatchesTrustedAggregator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
