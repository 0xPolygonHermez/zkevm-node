// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mockpolygonrollupmanager

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

// MockpolygonrollupmanagerMetaData contains all meta data concerning the Mockpolygonrollupmanager contract.
var MockpolygonrollupmanagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRoot\",\"name\":\"_globalExitRootManager\",\"type\":\"address\"},{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"_pol\",\"type\":\"address\"},{\"internalType\":\"contractIPolygonZkEVMBridge\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlOnlyCanRenounceRolesForSelf\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressDoNotHaveRequiredRole\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllzkEVMSequencedBatchesMustBeVerified\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BatchFeeOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ChainIDAlreadyExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExceedMaxVerifyBatches\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumBatchBelowLastVerifiedBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumBatchDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalPendingStateNumInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HaltTimeoutNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitBatchMustMatchCurrentForkID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitNumBatchAboveLastVerifiedBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitNumBatchDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeBatchTimeTarget\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeMultiplierBatchFee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSequenceSomeBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewAccInputHashDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewPendingStateTimeoutMustBeLower\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewStateRootNotInsidePrime\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewTrustedAggregatorTimeoutMustBeLower\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OldAccInputHashDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OldStateRootDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyEmergencyState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyNotEmergencyState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateNotConsolidable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupMustExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupTypeDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupTypeObsolete\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderMustBeRollup\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StoredRootMustBeDifferentThanNewRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TrustedAggregatorTimeoutNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpdateNotCompatible\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpdateToSameRollupTypeID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rollupAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatchBeforeUpgrade\",\"type\":\"uint64\"}],\"name\":\"AddExistingRollup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consensusImplementation\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"name\":\"AddNewRollupType\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"ConsolidatePendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rollupAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"gasTokenAddress\",\"type\":\"address\"}],\"name\":\"CreateNewRollup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EmergencyStateActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EmergencyStateDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"}],\"name\":\"ObsoleteRollupType\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"lastBatchSequenced\",\"type\":\"uint64\"}],\"name\":\"OnSequenceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"OverridePendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"storedStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"provedStateRoot\",\"type\":\"bytes32\"}],\"name\":\"ProveNonDeterministicPendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBatchFee\",\"type\":\"uint256\"}],\"name\":\"SetBatchFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"newMultiplierBatchFee\",\"type\":\"uint16\"}],\"name\":\"SetMultiplierBatchFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newPendingStateTimeout\",\"type\":\"uint64\"}],\"name\":\"SetPendingStateTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTrustedAggregator\",\"type\":\"address\"}],\"name\":\"SetTrustedAggregator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newTrustedAggregatorTimeout\",\"type\":\"uint64\"}],\"name\":\"SetTrustedAggregatorTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newVerifyBatchTimeTarget\",\"type\":\"uint64\"}],\"name\":\"SetVerifyBatchTimeTarget\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"newRollupTypeID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatchBeforeUpgrade\",\"type\":\"uint64\"}],\"name\":\"UpdateRollup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatchesTrustedAggregator\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activateEmergencyState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPolygonRollupBase\",\"name\":\"rollupAddress\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"}],\"name\":\"addExistingRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusImplementation\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"name\":\"addNewRollupType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMBridge\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateRewardPerBatch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"chainIDToRollupID\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"consolidatePendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"gasTokenAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"networkName\",\"type\":\"string\"}],\"name\":\"createNewRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deactivateEmergencyState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBatchFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getForcedBatchFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"getInputSnarkBytes\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"name\":\"getLastVerifiedBatch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"batchNum\",\"type\":\"uint64\"}],\"name\":\"getRollupBatchNumToStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"batchNum\",\"type\":\"uint64\"}],\"name\":\"getRollupPendingStateTransitions\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structLegacyZKEVMStateVariables.PendingState\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"batchNum\",\"type\":\"uint64\"}],\"name\":\"getRollupSequencedBatches\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accInputHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequencedTimestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"previousLastBatchSequenced\",\"type\":\"uint64\"}],\"internalType\":\"structLegacyZKEVMStateVariables.SequencedBatchData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"globalExitRootManager\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRoot\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"trustedAggregator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_pendingStateTimeout\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_trustedAggregatorTimeout\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"timelock\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"emergencyCouncil\",\"type\":\"address\"},{\"internalType\":\"contractPolygonZkEVMV2Existent\",\"name\":\"polygonZkEVM\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"zkEVMVerifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"zkEVMForkID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"zkEVMChainID\",\"type\":\"uint64\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"trustedAggregator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_pendingStateTimeout\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_trustedAggregatorTimeout\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"timelock\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"emergencyCouncil\",\"type\":\"address\"}],\"name\":\"initializeMock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isEmergencyState\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"isPendingStateConsolidable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastAggregationTimestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiplierBatchFee\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"}],\"name\":\"obsoleteRollupType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newSequencedBatches\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newAccInputHash\",\"type\":\"bytes32\"}],\"name\":\"onSequenceBatches\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"initPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"overridePendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingStateTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pol\",\"outputs\":[{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"localExitRoots\",\"type\":\"bytes32[]\"}],\"name\":\"prepareMockCalculateRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"initPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"proveNonDeterministicPendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rollupAddress\",\"type\":\"address\"}],\"name\":\"rollupAddressToID\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"name\":\"rollupIDToRollupData\",\"outputs\":[{\"internalType\":\"contractIPolygonRollupBase\",\"name\":\"rollupContract\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"lastLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"lastBatchSequenced\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastPendingState\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastPendingStateConsolidated\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatchBeforeUpgrade\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"rollupTypeID\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupTypeCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"}],\"name\":\"rollupTypeMap\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"consensusImplementation\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"obsolete\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newBatchFee\",\"type\":\"uint256\"}],\"name\":\"setBatchFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"newMultiplierBatchFee\",\"type\":\"uint16\"}],\"name\":\"setMultiplierBatchFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newPendingStateTimeout\",\"type\":\"uint64\"}],\"name\":\"setPendingStateTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newTrustedAggregatorTimeout\",\"type\":\"uint64\"}],\"name\":\"setTrustedAggregatorTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newVerifyBatchTimeTarget\",\"type\":\"uint64\"}],\"name\":\"setVerifyBatchTimeTarget\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSequencedBatches\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalVerifiedBatches\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedAggregatorTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractITransparentUpgradeableProxy\",\"name\":\"rollupContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"newRollupTypeID\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"upgradeData\",\"type\":\"bytes\"}],\"name\":\"updateRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifyBatchTimeTarget\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"verifyBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"verifyBatchesTrustedAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e06040523480156200001157600080fd5b50604051620061e8380380620061e883398101604081905262000034916200006b565b6001600160a01b0392831660805290821660c0521660a052620000bf565b6001600160a01b03811681146200006857600080fd5b50565b6000806000606084860312156200008157600080fd5b83516200008e8162000052565b6020850151909350620000a18162000052565b6040850151909250620000b48162000052565b809150509250925092565b60805160a05160c0516160c162000127600039600081816109bc015281816122cd0152613b4d01526000818161078201528181612e1c0152613e4701526000818161091601528181611313015281816114c301528181611fad0152613d3601526160c16000f3fe60806040523480156200001157600080fd5b5060043610620002aa5760003560e01c80630645af0914620002af578063066ec01214620002c8578063080b311114620002f45780630a0d9fbe146200031c5780630e36f582146200033757806311f6b287146200034e57806312b86e1914620003655780631489ed10146200037c57806315064c9614620003935780631608859c14620003a15780631796a1ae14620003b85780631816b7e514620003df5780632072f6c514620003f6578063248a9ca314620004005780632528016914620004265780632f2ff15d14620004db57806336568abe14620004f2578063394218e91462000509578063477fa270146200052057806355a71ee0146200052957806360469169146200056d57806365c0504d14620005775780637222020f1462000626578063727885e9146200063d5780637975fcfe14620006545780637fb6e76a146200067a578063841b24d714620006a357806387c20c0114620006be5780638bd4f07114620006d55780638f698ec514620006ec57806391d14854146200070357806399f5634e146200071a5780639a908e7314620007245780639c9f3dfe146200073b578063a066215c1462000752578063a217fddf1462000769578063a2967d991462000772578063a3c573eb146200077c578063afd23cbe14620007b3578063b99d0ad714620007dd578063c1acbc3414620008b5578063c4c928c214620008d0578063ceee281d14620008e7578063d02103ca1462000910578063d5073f6f1462000938578063d547741f146200094f578063d939b3151462000966578063dbc16976146200097a578063dde0ff771462000984578063e0bfd3d2146200099f578063e46761c414620009b6578063f34eb8eb14620009de578063f4e9267514620009f5578063f9c4c2ae1462000a06575b600080fd5b620002c6620002c0366004620047fe565b62000b1d565b005b608454620002dc906001600160401b031681565b604051620002eb9190620048d9565b60405180910390f35b6200030b6200030536600462004902565b62000f6f565b6040519015158152602001620002eb565b608554620002dc90600160401b90046001600160401b031681565b620002c6620003483660046200493a565b62000f99565b620002dc6200035f366004620049c5565b62001247565b620002c662000376366004620049f6565b62001267565b620002c66200038d36600462004a8d565b62001417565b606f546200030b9060ff1681565b620002c6620003b236600462004902565b620015a7565b607e54620003c99063ffffffff1681565b60405163ffffffff9091168152602001620002eb565b620002c6620003f036600462004b17565b6200163c565b620002c6620016e8565b620004176200041136600462004b44565b6200177e565b604051908152602001620002eb565b620004a76200043736600462004902565b60408051606080820183526000808352602080840182905292840181905263ffffffff959095168552608182528285206001600160401b03948516865260030182529382902082519485018352805485526001015480841691850191909152600160401b90049091169082015290565b60408051825181526020808401516001600160401b03908116918301919091529282015190921690820152606001620002eb565b620002c6620004ec36600462004b5e565b62001793565b620002c66200050336600462004b5e565b620017b5565b620002c66200051a36600462004b91565b620017ef565b60865462000417565b620004176200053a36600462004902565b63ffffffff821660009081526081602090815260408083206001600160401b038516845260020190915290205492915050565b620004176200189e565b620005dc62000588366004620049c5565b607f602052600090815260409020805460018201546002909201546001600160a01b0391821692918216916001600160401b03600160a01b8204169160ff600160e01b8304811692600160e81b9004169086565b604080516001600160a01b0397881681529690951660208701526001600160401b039093169385019390935260ff166060840152901515608083015260a082015260c001620002eb565b620002c662000637366004620049c5565b620018b6565b620002c66200064e36600462004c6e565b620019a1565b6200066b6200066536600462004d3b565b62001e09565b604051620002eb919062004df5565b620003c96200068b36600462004b91565b60836020526000908152604090205463ffffffff1681565b608454620002dc90600160c01b90046001600160401b031681565b620002c6620006cf36600462004a8d565b62001e3c565b620002c6620006e6366004620049f6565b62002160565b620002c6620006fd36600462004e0a565b62002216565b6200030b6200071436600462004b5e565b6200229d565b62000417620022c8565b620002dc6200073536600462004eba565b620023b4565b620002c66200074c36600462004b91565b62002581565b620002c66200076336600462004b91565b62002624565b62000417600081565b62000417620026c3565b620007a47f000000000000000000000000000000000000000000000000000000000000000081565b604051620002eb919062004ee7565b608554620007c990600160801b900461ffff1681565b60405161ffff9091168152602001620002eb565b62000873620007ee36600462004902565b604080516080808201835260008083526020808401829052838501829052606093840182905263ffffffff969096168152608186528381206001600160401b03958616825260040186528390208351918201845280548086168352600160401b9004909416948101949094526001830154918401919091526002909101549082015290565b604051620002eb919081516001600160401b03908116825260208084015190911690820152604082810151908201526060918201519181019190915260800190565b608454620002dc90600160801b90046001600160401b031681565b620002c6620008e136600462004efb565b62002a85565b620003c9620008f836600462004f93565b60826020526000908152604090205463ffffffff1681565b620007a47f000000000000000000000000000000000000000000000000000000000000000081565b620002c66200094936600462004b44565b62002d52565b620002c66200096036600462004b5e565b62002ddd565b608554620002dc906001600160401b031681565b620002c662002dff565b608454620002dc90600160401b90046001600160401b031681565b620002c6620009b036600462004fc5565b62002e9c565b620007a47f000000000000000000000000000000000000000000000000000000000000000081565b620002c6620009ef36600462005041565b62002f26565b608054620003c99063ffffffff1681565b62000a9d62000a17366004620049c5565b608160205260009081526040902080546001820154600583015460068401546007909401546001600160a01b0380851695600160a01b958690046001600160401b039081169692861695929092048216939282821692600160401b808404821693600160801b808204841694600160c01b90920484169380831693830416910460ff168c565b604080516001600160a01b039d8e1681526001600160401b039c8d1660208201529c909a16998c019990995296891660608b015260808a019590955292871660a089015290861660c0880152851660e0870152841661010086015283166101208501529190911661014083015260ff1661016082015261018001620002eb565b600054600290610100900460ff1615801562000b40575060005460ff8083169116105b62000b685760405162461bcd60e51b815260040162000b5f90620050d8565b60405180910390fd5b6000805461010060ff841661ffff199092169190911717905560858054608480546001600160c01b0316600160c01b6001600160401b038e8116919091029190911790915567016345785d8a00006086558c166001600160801b03199091161760e160431b1761ffff60801b19166101f560811b17905562000be962003111565b62000c046000805160206200606c8339815191528c6200317e565b62000c116000886200317e565b62000c2c60008051602062005f4c833981519152886200317e565b62000c4760008051602062005fcc833981519152886200317e565b62000c6260008051602062005eec833981519152886200317e565b62000c7d60008051602062005f2c833981519152896200317e565b62000c986000805160206200604c833981519152896200317e565b62000cb360008051602062005f6c833981519152896200317e565b62000cce60008051602062005fec833981519152896200317e565b62000cf86000805160206200606c83398151915260008051602062005ecc8339815191526200318a565b62000d1360008051602062005ecc833981519152896200317e565b62000d2e60008051602062005f0c833981519152896200317e565b62000d586000805160206200602c8339815191526000805160206200600c8339815191526200318a565b62000d736000805160206200602c833981519152876200317e565b62000d8e6000805160206200600c833981519152876200317e565b6073546074546001600160401b03600160401b9092048216911680821462000dc957604051632e4cc54360e11b815260040160405180910390fd5b600062000df1888888886000607460009054906101000a90046001600160401b0316620031df565b6001600160401b03838116600081815260756020908152604080832054600287018352818420558885168084526072808452828520600389018552948390208554815560018087018054919092018054918a166001600160401b03198084168217835593546001600160801b0319938416909117600160401b918290048c1682021790915560068b018054909216909417938802939093179092556007880180549092169095179055607a54606f54949092529154607354925163176b20e160e31b81529596506001600160a01b038e81169663bb5907089662000ef39694831695600160581b909504909216936076936077939092911690600401620051d2565b600060405180830381600087803b15801562000f0e57600080fd5b505af115801562000f23573d6000803e3d6000fd5b50506000805461ff0019169055505060405160ff8516815260008051602062005fac8339815191529350602001915062000f5a9050565b60405180910390a15050505050505050505050565b63ffffffff8216600090815260816020526040812062000f9090836200340d565b90505b92915050565b600054600290610100900460ff1615801562000fbc575060005460ff8083169116105b62000fdb5760405162461bcd60e51b815260040162000b5f90620050d8565b6000805461010060ff841661ffff199092169190911717905560858054608480546001600160c01b0316600160c01b6001600160401b038a8116919091029190911790915567016345785d8a000060865588166001600160801b03199091161760e160431b1761ffff60801b19166101f560811b1790556200105c62003111565b620010776000805160206200606c833981519152886200317e565b620010846000846200317e565b6200109f60008051602062005f4c833981519152846200317e565b620010ba60008051602062005fcc833981519152846200317e565b620010d560008051602062005eec833981519152846200317e565b620010f060008051602062005f2c833981519152856200317e565b6200110b6000805160206200604c833981519152856200317e565b6200112660008051602062005f6c833981519152856200317e565b6200114160008051602062005fec833981519152856200317e565b6200116b6000805160206200606c83398151915260008051602062005ecc8339815191526200318a565b6200118660008051602062005ecc833981519152856200317e565b620011a160008051602062005f0c833981519152856200317e565b620011cb6000805160206200602c8339815191526000805160206200600c8339815191526200318a565b620011e66000805160206200602c833981519152836200317e565b620012016000805160206200600c833981519152836200317e565b6200120e6000336200317e565b6000805461ff001916905560405160ff8216815260008051602062005fac8339815191529060200160405180910390a150505050505050565b63ffffffff8116600090815260816020526040812062000f939062003452565b6000805160206200606c8339815191526200128281620034c3565b63ffffffff89166000908152608160205260409020620012a9818a8a8a8a8a8a8a620034cf565b600681018054600160401b600160801b031916600160401b6001600160401b0389811691820292909217835560009081526002840160205260409020869055600583018790559054600160801b9004161562001311576006810180546001600160801b031690555b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166333d6247d6200134a620026c3565b6040518263ffffffff1660e01b81526004016200136991815260200190565b600060405180830381600087803b1580156200138457600080fd5b505af115801562001399573d6000803e3d6000fd5b5050608480546001600160c01b031661127560c71b1790555050604080516001600160401b03881681526020810186905290810186905233606082015263ffffffff8b16907f3182bd6e6f74fc1fdc88b60f3a4f4c7f79db6ae6f5b88a1b3f5a1e28ec210d5e9060800160405180910390a250505050505050505050565b6000805160206200606c8339815191526200143281620034c3565b63ffffffff8916600090815260816020526040902062001459818a8a8a8a8a8a8a62003857565b600681018054600160401b600160801b031916600160401b6001600160401b038a811691820292909217835560009081526002840160205260409020879055600583018890559054600160801b90041615620014c1576006810180546001600160801b031690555b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166333d6247d620014fa620026c3565b6040518263ffffffff1660e01b81526004016200151991815260200190565b600060405180830381600087803b1580156200153457600080fd5b505af115801562001549573d6000803e3d6000fd5b50505050336001600160a01b03168a63ffffffff167fd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d389888a604051620015939392919062005238565b60405180910390a350505050505050505050565b63ffffffff82166000908152608160205260409020620015d76000805160206200606c833981519152336200229d565b6200162b57606f5460ff16156200160157604051630bc011ff60e21b815260040160405180910390fd5b6200160d81836200340d565b6200162b57604051630674f25160e11b815260040160405180910390fd5b62001637818362003c53565b505050565b60008051602062005fec8339815191526200165781620034c3565b6103e88261ffff1610806200167157506103ff8261ffff16115b156200169057604051630984a67960e31b815260040160405180910390fd5b6085805461ffff60801b1916600160801b61ffff8516908102919091179091556040519081527f7019933d795eba185c180209e8ae8bffbaa25bcef293364687702c31f4d302c5906020015b60405180910390a15050565b620017036000805160206200602c833981519152336200229d565b6200177257608454600160801b90046001600160401b031615806200175357506084544290620017489062093a8090600160801b90046001600160401b03166200526f565b6001600160401b0316115b15620017725760405163692baaad60e11b815260040160405180910390fd5b6200177c62003e45565b565b60009081526034602052604090206001015490565b6200179e826200177e565b620017a981620034c3565b62001637838362003ec4565b6001600160a01b0381163314620017df57604051630b4ad1cd60e31b815260040160405180910390fd5b620017eb828262003f30565b5050565b60008051602062005fec8339815191526200180a81620034c3565b606f5460ff166200184c576084546001600160401b03600160c01b9091048116908316106200184c5760405163401636df60e01b815260040160405180910390fd5b608480546001600160c01b0316600160c01b6001600160401b038516021790556040517f1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a190620016dc908490620048d9565b60006086546064620018b1919062005299565b905090565b60008051602062005f2c833981519152620018d181620034c3565b63ffffffff82161580620018f05750607e5463ffffffff908116908316115b156200190f57604051637512e5cb60e01b815260040160405180910390fd5b63ffffffff82166000908152607f60205260409020600180820154600160e81b900460ff16151590036200195657604051633b8d3d9960e01b815260040160405180910390fd5b60018101805460ff60e81b1916600160e81b17905560405163ffffffff8416907f4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e4490600090a2505050565b6000805160206200604c833981519152620019bc81620034c3565b63ffffffff88161580620019db5750607e5463ffffffff908116908916115b15620019fa57604051637512e5cb60e01b815260040160405180910390fd5b63ffffffff88166000908152607f60205260409020600180820154600160e81b900460ff161515900362001a4157604051633b8d3d9960e01b815260040160405180910390fd5b6001600160401b03881660009081526083602052604090205463ffffffff161562001a7f576040516337c8fe0960e11b815260040160405180910390fd5b6080805460009190829062001a9a9063ffffffff16620052b3565b825463ffffffff8281166101009490940a9384029302191691909117909155825460408051600080825260208201928390529394506001600160a01b0390921691309162001ae890620047c2565b62001af693929190620052d9565b604051809103906000f08015801562001b13573d6000803e3d6000fd5b50905081608360008c6001600160401b03166001600160401b0316815260200190815260200160002060006101000a81548163ffffffff021916908363ffffffff1602179055508160826000836001600160a01b03166001600160a01b0316815260200190815260200160002060006101000a81548163ffffffff021916908363ffffffff1602179055506000608160008463ffffffff1663ffffffff1681526020019081526020016000209050818160000160006101000a8154816001600160a01b0302191690836001600160a01b031602179055508360010160149054906101000a90046001600160401b03168160010160146101000a8154816001600160401b0302191690836001600160401b031602179055508360010160009054906101000a90046001600160a01b03168160010160006101000a8154816001600160a01b0302191690836001600160a01b031602179055508a8160000160146101000a8154816001600160401b0302191690836001600160401b031602179055508360020154816002016000806001600160401b03168152602001908152602001600020819055508b63ffffffff168160070160086101000a8154816001600160401b0302191690836001600160401b0316021790555083600101601c9054906101000a900460ff168160070160106101000a81548160ff021916908360ff1602179055508263ffffffff167f194c983456df6701c6a50830b90fe80e72b823411d0d524970c9590dc277a6418d848e8c60405162001d87949392919063ffffffff9490941684526001600160a01b0392831660208501526001600160401b0391909116604084015216606082015260800190565b60405180910390a2604051633892b81160e11b81526001600160a01b0383169063712570229062001dc7908d908d9088908e908e908e9060040162005310565b600060405180830381600087803b15801562001de257600080fd5b505af115801562001df7573d6000803e3d6000fd5b50505050505050505050505050505050565b63ffffffff8616600090815260816020526040902060609062001e3190878787878762003f9a565b979650505050505050565b606f5460ff161562001e6157604051630bc011ff60e21b815260040160405180910390fd5b63ffffffff881660009081526081602090815260408083206084546001600160401b038a81168652600383019094529190932060010154429262001eb092600160c01b9004811691166200526f565b6001600160401b0316111562001ed957604051638a0704d360e01b815260040160405180910390fd5b6103e862001ee8888862005373565b6001600160401b0316111562001f1157604051635acfba9d60e11b815260040160405180910390fd5b62001f23818989898989898962003857565b62001f2f8187620040d5565b6085546001600160401b03166000036200203d57600681018054600160401b600160801b031916600160401b6001600160401b0389811691820292909217835560009081526002840160205260409020869055600583018790559054600160801b9004161562001fab576006810180546001600160801b031690555b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166333d6247d62001fe4620026c3565b6040518263ffffffff1660e01b81526004016200200391815260200190565b600060405180830381600087803b1580156200201e57600080fd5b505af115801562002033573d6000803e3d6000fd5b5050505062002107565b6200204881620042d2565b600681018054600160801b90046001600160401b03169060106200206c8362005396565b82546001600160401b039182166101009390930a92830292820219169190911790915560408051608081018252428316815289831660208083019182528284018b8152606084018b81526006890154600160801b90048716600090815260048a01909352949091209251835492518616600160401b026001600160801b03199093169516949094171781559151600183015551600290910155505b336001600160a01b03168963ffffffff167faac1e7a157b259544ebacd6e8a82ae5d6c8f174e12aa48696277bcc9a661f0b48887896040516200214d9392919062005238565b60405180910390a3505050505050505050565b606f5460ff16156200218557604051630bc011ff60e21b815260040160405180910390fd5b63ffffffff88166000908152608160205260409020620021ac8189898989898989620034cf565b6001600160401b03851660009081526002820160209081526040918290205482519081529081018590527f1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010910160405180910390a16200220b62003e45565b505050505050505050565b80516080805463ffffffff191663ffffffff90921691909117905560005b8151811015620017eb57818181518110620022535762002253620053bd565b6020026020010151608160008360016200226e9190620053d3565b63ffffffff168152602081019190915260400160002060050155806200229481620053e9565b91505062002234565b60009182526034602090815260408084206001600160a01b0393909316845291905290205460ff1690565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166370a08231306040518263ffffffff1660e01b815260040162002319919062004ee7565b602060405180830381865afa15801562002337573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200235d919062005405565b60845490915060009062002384906001600160401b03600160401b82048116911662005373565b6001600160401b0316905080600003620023a15760009250505090565b620023ad818362005435565b9250505090565b606f5460009060ff1615620023dc57604051630bc011ff60e21b815260040160405180910390fd5b3360009081526082602052604081205463ffffffff169081900362002414576040516371653c1560e01b815260040160405180910390fd5b836001600160401b03166000036200243f57604051632590ccf960e01b815260040160405180910390fd5b63ffffffff811660009081526081602052604081206084805491928792620024729084906001600160401b03166200526f565b82546101009290920a6001600160401b0381810219909316918316021790915560068301541690506000620024a887836200526f565b6006840180546001600160401b038084166001600160401b03199092168217909255604080516060810182528a81524284166020808301918252888616838501908152600095865260038b0190915292909320905181559151600192909201805491518416600160401b026001600160801b0319909216929093169190911717905590506200253783620042d2565b8363ffffffff167f1d9f30260051d51d70339da239ea7b080021adcaabfa71c9b0ea339a20cf9a25826040516200256f9190620048d9565b60405180910390a29695505050505050565b60008051602062005fec8339815191526200259c81620034c3565b606f5460ff16620025d7576085546001600160401b0390811690831610620025d75760405163048a05a960e41b815260040160405180910390fd5b608580546001600160401b0319166001600160401b0384161790556040517fc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c7590620016dc908490620048d9565b60008051602062005fec8339815191526200263f81620034c3565b62015180826001600160401b031611156200266d57604051631c0cfbfd60e31b815260040160405180910390fd5b60858054600160401b600160801b031916600160401b6001600160401b038516021790556040517f1b023231a1ab6b5d93992f168fb44498e1a7e64cef58daff6f1c216de6a68c2890620016dc908490620048d9565b60805460009063ffffffff16808203620026df57506000919050565b6000816001600160401b03811115620026fc57620026fc62004baf565b60405190808252806020026020018201604052801562002726578160200160208202803683370190505b50905060005b8281101562002799576081600062002746836001620053d3565b63ffffffff1663ffffffff16815260200190815260200160002060050154828281518110620027795762002779620053bd565b6020908102919091010152806200279081620053e9565b9150506200272c565b50600060205b83600114620029dd576000620027b76002866200544c565b620027c460028762005435565b620027d09190620053d3565b90506000816001600160401b03811115620027ef57620027ef62004baf565b60405190808252806020026020018201604052801562002819578160200160208202803683370190505b50905060005b8281101562002991576200283560018462005463565b811480156200285057506200284c6002886200544c565b6001145b15620028d057856200286482600262005299565b81518110620028775762002877620053bd565b6020026020010151856040516020016200289392919062005479565b60405160208183030381529060405280519060200120828281518110620028be57620028be620053bd565b6020026020010181815250506200297c565b85620028de82600262005299565b81518110620028f157620028f1620053bd565b60200260200101518682600262002909919062005299565b62002916906001620053d3565b81518110620029295762002929620053bd565b60200260200101516040516020016200294492919062005479565b604051602081830303815290604052805190602001208282815181106200296f576200296f620053bd565b6020026020010181815250505b806200298881620053e9565b9150506200281f565b508094508195508384604051602001620029ad92919062005479565b6040516020818303038152906040528051906020012093508280620029d29062005487565b93505050506200279f565b600083600081518110620029f557620029f5620053bd565b6020026020010151905060005b8281101562002a7b57818460405160200162002a2092919062005479565b604051602081830303815290604052805190602001209150838460405160200162002a4d92919062005479565b604051602081830303815290604052805190602001209350808062002a7290620053e9565b91505062002a02565b5095945050505050565b60008051602062005eec83398151915262002aa081620034c3565b63ffffffff8416158062002abf5750607e5463ffffffff908116908516115b1562002ade57604051637512e5cb60e01b815260040160405180910390fd5b6001600160a01b03851660009081526082602052604081205463ffffffff169081900362002b1f576040516374a086a360e01b815260040160405180910390fd5b63ffffffff8181166000908152608160205260409020600781015490918716600160401b9091046001600160401b03160362002b6e57604051634f61d51960e01b815260040160405180910390fd5b63ffffffff86166000908152607f60205260409020600180820154600160e81b900460ff161515900362002bb557604051633b8d3d9960e01b815260040160405180910390fd5b60018101546007830154600160801b900460ff908116600160e01b909204161462002bf357604051635aa0d5f160e11b815260040160405180910390fd5b6001808201805491840180546001600160a01b031981166001600160a01b03909416938417825591546001600160401b03600160a01b9182900416026001600160e01b0319909216909217179055600782018054600160401b63ffffffff8a1602600160401b600160801b0319909116179055600062002c738462001247565b6007840180546001600160401b0319166001600160401b038316179055825460405163278f794360e11b81529192506001600160a01b038b811692634f1ef2869262002cc89216908b908b90600401620054a1565b600060405180830381600087803b15801562002ce357600080fd5b505af115801562002cf8573d6000803e3d6000fd5b50506040805163ffffffff8c811682526001600160401b0386166020830152881693507ff585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d92500160405180910390a2505050505050505050565b60008051602062005f0c83398151915262002d6d81620034c3565b683635c9adc5dea0000082118062002d885750633b9aca0082105b1562002da757604051638586952560e01b815260040160405180910390fd5b60868290556040518281527ffb383653f53ee079978d0c9aff7aeff04a10166ce244cca9c9f9d8d96bed45b290602001620016dc565b62002de8826200177e565b62002df381620034c3565b62001637838362003f30565b60008051602062005f6c83398151915262002e1a81620034c3565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663dbc169766040518163ffffffff1660e01b8152600401600060405180830381600087803b15801562002e7657600080fd5b505af115801562002e8b573d6000803e3d6000fd5b5050505062002e996200439d565b50565b60008051602062005fcc83398151915262002eb781620034c3565b6001600160401b03841660009081526083602052604090205463ffffffff161562002ef5576040516337c8fe0960e11b815260040160405180910390fd5b600062002f0888888888876000620031df565b60008080526002909101602052604090209390935550505050505050565b60008051602062005f4c83398151915262002f4181620034c3565b607e805460009190829062002f5c9063ffffffff16620052b3565b91906101000a81548163ffffffff021916908363ffffffff160217905590506040518060c00160405280896001600160a01b03168152602001886001600160a01b03168152602001876001600160401b031681526020018660ff16815260200160001515815260200185815250607f60008363ffffffff1663ffffffff16815260200190815260200160002060008201518160000160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060208201518160010160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060408201518160010160146101000a8154816001600160401b0302191690836001600160401b03160217905550606082015181600101601c6101000a81548160ff021916908360ff160217905550608082015181600101601d6101000a81548160ff02191690831515021790555060a082015181600201559050508063ffffffff167fa2970448b3bd66ba7e524e7b2a5b9cf94fa29e32488fb942afdfe70dd4b77b52898989898989604051620030ff96959493929190620054e1565b60405180910390a25050505050505050565b600054610100900460ff166200177c5760405162461bcd60e51b815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201526a6e697469616c697a696e6760a81b606482015260840162000b5f565b620017eb828262003ec4565b600062003197836200177e565b600084815260346020526040808220600101859055519192508391839186917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9190a4505050565b6080805460009182918290620031fb9063ffffffff16620052b3565b91906101000a81548163ffffffff021916908363ffffffff160217905590508060836000876001600160401b03166001600160401b0316815260200190815260200160002060006101000a81548163ffffffff021916908363ffffffff16021790555080608260008a6001600160a01b03166001600160a01b0316815260200190815260200160002060006101000a81548163ffffffff021916908363ffffffff160217905550608160008263ffffffff1663ffffffff1681526020019081526020016000209150878260000160006101000a8154816001600160a01b0302191690836001600160a01b03160217905550858260010160146101000a8154816001600160401b0302191690836001600160401b03160217905550868260010160006101000a8154816001600160a01b0302191690836001600160a01b03160217905550848260000160146101000a8154816001600160401b0302191690836001600160401b03160217905550838260070160106101000a81548160ff021916908360ff1602179055508063ffffffff167fadfc7d56f7e39b08b321534f14bfb135ad27698f7d2f5ad0edc2356ea9a3f850878a888888604051620033fa9594939291906001600160401b0395861681526001600160a01b03949094166020850152918416604084015260ff166060830152909116608082015260a00190565b60405180910390a2509695505050505050565b6085546001600160401b03828116600090815260048501602052604081205490924292620034409291811691166200526f565b6001600160401b031611159392505050565b6006810154600090600160801b90046001600160401b031615620034a6575060068101546001600160401b03600160801b909104811660009081526004909201602052604090912054600160401b90041690565b5060060154600160401b90046001600160401b031690565b919050565b62002e998133620043f6565b60078801546000906001600160401b039081169087161015620035055760405163ead1340b60e01b815260040160405180910390fd5b6001600160401b03881615620035a65760068901546001600160401b03600160801b909104811690891611156200354f5760405163bb14c20560e01b815260040160405180910390fd5b506001600160401b03808816600090815260048a0160205260409020600281015481549092888116600160401b90920416146200359f57604051632bd2e3e760e01b815260040160405180910390fd5b506200361b565b506001600160401b038516600090815260028901602052604090205480620035e1576040516324cbdcc360e11b815260040160405180910390fd5b60068901546001600160401b03600160401b909104811690871611156200361b57604051630f2b74f160e11b815260040160405180910390fd5b60068901546001600160401b03600160801b90910481169088161180620036545750876001600160401b0316876001600160401b031611155b8062003678575060068901546001600160401b03600160c01b909104811690881611155b15620036975760405163bfa7079f60e01b815260040160405180910390fd5b6001600160401b03878116600090815260048b016020526040902054600160401b9004811690861614620036de576040516332a2a77f60e01b815260040160405180910390fd5b6000620036f08a888888868962003f9a565b9050600060008051602062005f8c8339815191526002836040516200371691906200552e565b602060405180830381855afa15801562003734573d6000803e3d6000fd5b5050506040513d601f19601f8201168201806040525081019062003759919062005405565b6200376591906200544c565b60018c0154604080516020810182528381529051634890ed4560e11b81529293506001600160a01b0390911691639121da8a91620037a9918891906004016200554c565b602060405180830381865afa158015620037c7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620037ed919062005589565b6200380b576040516309bde33960e01b815260040160405180910390fd5b6001600160401b038916600090815260048c0160205260409020600201548590036200384a5760405163a47276bd60e01b815260040160405180910390fd5b5050505050505050505050565b600080620038658a62003452565b60078b01549091506001600160401b0390811690891610156200389b5760405163ead1340b60e01b815260040160405180910390fd5b6001600160401b038916156200393e5760068a01546001600160401b03600160801b9091048116908a161115620038e55760405163bb14c20560e01b815260040160405180910390fd5b6001600160401b03808a16600090815260048c01602052604090206002810154815490945090918a8116600160401b90920416146200393757604051632bd2e3e760e01b815260040160405180910390fd5b50620039ae565b6001600160401b038816600090815260028b0160205260409020549150816200397a576040516324cbdcc360e11b815260040160405180910390fd5b806001600160401b0316886001600160401b03161115620039ae57604051630f2b74f160e11b815260040160405180910390fd5b806001600160401b0316876001600160401b031611620039e15760405163b9b18f5760e01b815260040160405180910390fd5b6000620039f38b8a8a8a878b62003f9a565b9050600060008051602062005f8c83398151915260028360405162003a1991906200552e565b602060405180830381855afa15801562003a37573d6000803e3d6000fd5b5050506040513d601f19601f8201168201806040525081019062003a5c919062005405565b62003a6891906200544c565b60018d0154604080516020810182528381529051634890ed4560e11b81529293506001600160a01b0390911691639121da8a9162003aac918991906004016200554c565b602060405180830381865afa15801562003aca573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062003af0919062005589565b62003b0e576040516309bde33960e01b815260040160405180910390fd5b600062003b1c848b62005373565b905062003b7587826001600160401b031662003b37620022c8565b62003b43919062005299565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016919062004420565b80608460088282829054906101000a90046001600160401b031662003b9b91906200526f565b82546101009290920a6001600160401b0381810219909316918316021790915560848054600160801b600160c01b031916600160801b428416021790558e546040516332c2d15360e01b8152918d166004830152602482018b90523360448301526001600160a01b031691506332c2d15390606401600060405180830381600087803b15801562003c2b57600080fd5b505af115801562003c40573d6000803e3d6000fd5b5050505050505050505050505050505050565b60068201546001600160401b03600160c01b909104811690821611158062003c92575060068201546001600160401b03600160801b9091048116908216115b1562003cb15760405163d086b70b60e01b815260040160405180910390fd5b6001600160401b03818116600081815260048501602090815260408083208054600689018054600160401b600160801b031916600160401b92839004909816918202979097178755600280830154828752908a0190945291909320919091556001820154600587015583546001600160c01b0316600160c01b909302929092179092557f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166333d6247d62003d6d620026c3565b6040518263ffffffff1660e01b815260040162003d8c91815260200190565b600060405180830381600087803b15801562003da757600080fd5b505af115801562003dbc573d6000803e3d6000fd5b505085546001600160a01b0316600090815260826020908152604091829020546002870154600188015484516001600160401b03898116825294810192909252818501529188166060830152915163ffffffff90921693507f581910eb7a27738945c2f00a91f2284b2d6de9d4e472b12f901c2b0df045e21b925081900360800190a250505050565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316632072f6c56040518163ffffffff1660e01b8152600401600060405180830381600087803b15801562003ea157600080fd5b505af115801562003eb6573d6000803e3d6000fd5b505050506200177c62004474565b62003ed082826200229d565b620017eb5760008281526034602090815260408083206001600160a01b0385168085529252808320805460ff1916600117905551339285917f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9190a45050565b62003f3c82826200229d565b15620017eb5760008281526034602090815260408083206001600160a01b0385168085529252808320805460ff1916905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b6001600160401b038086166000818152600389016020526040808220549388168252902054606092911580159062003fd0575081155b1562003fef5760405163340c614f60e11b815260040160405180910390fd5b806200400e576040516366385b5160e01b815260040160405180910390fd5b6200401984620044d1565b62004037576040516305dae44f60e21b815260040160405180910390fd5b885460018a01546040516001600160601b03193360601b16602082015260348101889052605481018590526001600160c01b031960c08c811b82166074840152600160a01b94859004811b8216607c84015293909204831b82166084820152608c810187905260ac810184905260cc81018990529189901b1660ec82015260f401604051602081830303815290604052925050509695505050505050565b6000620040e28362003452565b905081600080620040f4848462005373565b6085546001600160401b0391821692506000916200411b91600160401b9004164262005463565b90505b846001600160401b0316846001600160401b031614620041a5576001600160401b0380851660009081526003890160205260409020600181015490911682101562004180576001810154600160401b90046001600160401b031694506200419e565b6200418c868662005373565b6001600160401b0316935050620041a5565b506200411e565b6000620041b3848462005463565b9050838110156200421157808403600c8111620041d15780620041d4565b600c5b9050806103e80a81608560109054906101000a900461ffff1661ffff160a60865402816200420657620042066200541f565b046086555062004289565b838103600c811162004224578062004227565b600c5b90506000816103e80a82608560109054906101000a900461ffff1661ffff160a670de0b6b3a764000002816200426157620042616200541f565b04905080608654670de0b6b3a764000002816200428257620042826200541f565b0460865550505b683635c9adc5dea000006086541115620042b057683635c9adc5dea00000608655620042c8565b633b9aca006086541015620042c857633b9aca006086555b5050505050505050565b60068101546001600160401b03600160c01b82048116600160801b90920416111562002e995760068101546000906200431d90600160c01b90046001600160401b031660016200526f565b90506200432b82826200340d565b15620017eb57600682015460009060029062004359908490600160801b90046001600160401b031662005373565b620043659190620055ad565b6200437190836200526f565b90506200437f83826200340d565b15620043915762001637838262003c53565b62001637838362003c53565b606f5460ff16620043c157604051635386698160e01b815260040160405180910390fd5b606f805460ff191690556040517f1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b390600090a1565b6200440282826200229d565b620017eb57604051637615be1f60e11b815260040160405180910390fd5b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b031663a9059cbb60e01b1790526200163790849062004557565b606f5460ff16156200449957604051630bc011ff60e21b815260040160405180910390fd5b606f805460ff191660011790556040517f2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a549790600090a1565b600067ffffffff000000016001600160401b03831610801562004508575067ffffffff00000001604083901c6001600160401b0316105b801562004529575067ffffffff00000001608083901c6001600160401b0316105b801562004541575067ffffffff0000000160c083901c105b156200454f57506001919050565b506000919050565b6000620045ae826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b0316620046309092919063ffffffff16565b805190915015620016375780806020019051810190620045cf919062005589565b620016375760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840162000b5f565b606062004641848460008562004649565b949350505050565b606082471015620046ac5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840162000b5f565b600080866001600160a01b03168587604051620046ca91906200552e565b60006040518083038185875af1925050503d806000811462004709576040519150601f19603f3d011682016040523d82523d6000602084013e6200470e565b606091505b509150915062001e3187838387606083156200478f57825160000362004787576001600160a01b0385163b620047875760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640162000b5f565b508162004641565b620046418383815115620047a65781518083602001fd5b8060405162461bcd60e51b815260040162000b5f919062004df5565b6108f580620055d783390190565b6001600160a01b038116811462002e9957600080fd5b80356001600160401b0381168114620034be57600080fd5b6000806000806000806000806000806101408b8d0312156200481f57600080fd5b8a356200482c81620047d0565b99506200483c60208c01620047e6565b98506200484c60408c01620047e6565b975060608b01356200485e81620047d0565b965060808b01356200487081620047d0565b955060a08b01356200488281620047d0565b945060c08b01356200489481620047d0565b935060e08b0135620048a681620047d0565b9250620048b76101008c01620047e6565b9150620048c86101208c01620047e6565b90509295989b9194979a5092959850565b6001600160401b0391909116815260200190565b803563ffffffff81168114620034be57600080fd5b600080604083850312156200491657600080fd5b6200492183620048ed565b91506200493160208401620047e6565b90509250929050565b60008060008060008060c087890312156200495457600080fd5b86356200496181620047d0565b95506200497160208801620047e6565b94506200498160408801620047e6565b935060608701356200499381620047d0565b92506080870135620049a581620047d0565b915060a0870135620049b781620047d0565b809150509295509295509295565b600060208284031215620049d857600080fd5b62000f9082620048ed565b80610300810183101562000f9357600080fd5b6000806000806000806000806103e0898b03121562004a1457600080fd5b62004a1f89620048ed565b975062004a2f60208a01620047e6565b965062004a3f60408a01620047e6565b955062004a4f60608a01620047e6565b945062004a5f60808a01620047e6565b935060a0890135925060c0890135915062004a7e8a60e08b01620049e3565b90509295985092959890939650565b6000806000806000806000806103e0898b03121562004aab57600080fd5b62004ab689620048ed565b975062004ac660208a01620047e6565b965062004ad660408a01620047e6565b955062004ae660608a01620047e6565b94506080890135935060a0890135925060c089013562004b0681620047d0565b915062004a7e8a60e08b01620049e3565b60006020828403121562004b2a57600080fd5b813561ffff8116811462004b3d57600080fd5b9392505050565b60006020828403121562004b5757600080fd5b5035919050565b6000806040838503121562004b7257600080fd5b82359150602083013562004b8681620047d0565b809150509250929050565b60006020828403121562004ba457600080fd5b62000f9082620047e6565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b038111828210171562004bf05762004bf062004baf565b604052919050565b600082601f83011262004c0a57600080fd5b81356001600160401b0381111562004c265762004c2662004baf565b62004c3b601f8201601f191660200162004bc5565b81815284602083860101111562004c5157600080fd5b816020850160208301376000918101602001919091529392505050565b600080600080600080600060e0888a03121562004c8a57600080fd5b62004c9588620048ed565b965062004ca560208901620047e6565b9550604088013562004cb781620047d0565b9450606088013562004cc981620047d0565b9350608088013562004cdb81620047d0565b925060a08801356001600160401b038082111562004cf857600080fd5b62004d068b838c0162004bf8565b935060c08a013591508082111562004d1d57600080fd5b5062004d2c8a828b0162004bf8565b91505092959891949750929550565b60008060008060008060c0878903121562004d5557600080fd5b62004d6087620048ed565b955062004d7060208801620047e6565b945062004d8060408801620047e6565b9350606087013592506080870135915060a087013590509295509295509295565b60005b8381101562004dbe57818101518382015260200162004da4565b50506000910152565b6000815180845262004de181602086016020860162004da1565b601f01601f19169290920160200192915050565b60208152600062000f90602083018462004dc7565b6000602080838503121562004e1e57600080fd5b82356001600160401b038082111562004e3657600080fd5b818501915085601f83011262004e4b57600080fd5b81358181111562004e605762004e6062004baf565b8060051b915062004e7384830162004bc5565b818152918301840191848101908884111562004e8e57600080fd5b938501935b8385101562004eae5784358252938501939085019062004e93565b98975050505050505050565b6000806040838503121562004ece57600080fd5b62004ed983620047e6565b946020939093013593505050565b6001600160a01b0391909116815260200190565b6000806000806060858703121562004f1257600080fd5b843562004f1f81620047d0565b935062004f2f60208601620048ed565b925060408501356001600160401b038082111562004f4c57600080fd5b818701915087601f83011262004f6157600080fd5b81358181111562004f7157600080fd5b88602082850101111562004f8457600080fd5b95989497505060200194505050565b60006020828403121562004fa657600080fd5b813562004b3d81620047d0565b803560ff81168114620034be57600080fd5b60008060008060008060c0878903121562004fdf57600080fd5b863562004fec81620047d0565b9550602087013562004ffe81620047d0565b94506200500e60408801620047e6565b93506200501e60608801620047e6565b9250608087013591506200503560a0880162004fb3565b90509295509295509295565b60008060008060008060c087890312156200505b57600080fd5b86356200506881620047d0565b955060208701356200507a81620047d0565b94506200508a60408801620047e6565b93506200509a6060880162004fb3565b92506080870135915060a08701356001600160401b03811115620050bd57600080fd5b620050cb89828a0162004bf8565b9150509295509295509295565b6020808252602e908201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160408201526d191e481a5b9a5d1a585b1a5e995960921b606082015260800190565b8054600090600181811c90808316806200514157607f831692505b602080841082036200516357634e487b7160e01b600052602260045260246000fd5b838852602088018280156200518157600181146200519857620051c5565b60ff198716825285151560051b82019750620051c5565b60008981526020902060005b87811015620051bf57815484820152908601908401620051a4565b83019850505b5050505050505092915050565b6001600160a01b0387811682528616602082015260c060408201819052600090620052009083018762005126565b828103606084015262005214818762005126565b608084019590955250506001600160401b039190911660a090910152949350505050565b6001600160401b039390931683526020830191909152604082015260600190565b634e487b7160e01b600052601160045260246000fd5b6001600160401b0381811683821601908082111562005292576200529262005259565b5092915050565b808202811582820484141762000f935762000f9362005259565b600063ffffffff808316818103620052cf57620052cf62005259565b6001019392505050565b6001600160a01b03848116825283166020820152606060408201819052600090620053079083018462004dc7565b95945050505050565b6001600160a01b038781168252868116602083015263ffffffff861660408301528416606082015260c060808201819052600090620053529083018562004dc7565b82810360a084015262005366818562004dc7565b9998505050505050505050565b6001600160401b0382811682821603908082111562005292576200529262005259565b60006001600160401b038281166002600160401b03198101620052cf57620052cf62005259565b634e487b7160e01b600052603260045260246000fd5b8082018082111562000f935762000f9362005259565b600060018201620053fe57620053fe62005259565b5060010190565b6000602082840312156200541857600080fd5b5051919050565b634e487b7160e01b600052601260045260246000fd5b6000826200544757620054476200541f565b500490565b6000826200545e576200545e6200541f565b500690565b8181038181111562000f935762000f9362005259565b918252602082015260400190565b60008162005499576200549962005259565b506000190190565b6001600160a01b03841681526040602082018190528101829052818360608301376000818301606090810191909152601f909201601f1916010192915050565b6001600160a01b038781168252861660208201526001600160401b038516604082015260ff841660608201526080810183905260c060a0820181905260009062004eae9083018462004dc7565b600082516200554281846020870162004da1565b9190910192915050565b61032081016103008085843782018360005b60018110156200557f5781518352602092830192909101906001016200555e565b5050509392505050565b6000602082840312156200559c57600080fd5b8151801515811462004b3d57600080fd5b60006001600160401b0383811680620055ca57620055ca6200541f565b9216919091049291505056fe60a0604052604051620008f5380380620008f58339810160408190526100249161035b565b82816100308282610058565b50506001600160a01b03821660805261005061004b60805190565b6100b7565b505050610447565b61006182610126565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a28051156100ab576100a682826101a5565b505050565b6100b361021c565b5050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f6100f8600080516020620008d5833981519152546001600160a01b031690565b604080516001600160a01b03928316815291841660208301520160405180910390a16101238161023d565b50565b806001600160a01b03163b60000361016157604051634c9c8ce360e01b81526001600160a01b03821660048201526024015b60405180910390fd5b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5b80546001600160a01b0319166001600160a01b039290921691909117905550565b6060600080846001600160a01b0316846040516101c2919061042b565b600060405180830381855af49150503d80600081146101fd576040519150601f19603f3d011682016040523d82523d6000602084013e610202565b606091505b50909250905061021385838361027d565b95945050505050565b341561023b5760405163b398979f60e01b815260040160405180910390fd5b565b6001600160a01b03811661026757604051633173bdd160e11b815260006004820152602401610158565b80600080516020620008d5833981519152610184565b6060826102925761028d826102dc565b6102d5565b81511580156102a957506001600160a01b0384163b155b156102d257604051639996b31560e01b81526001600160a01b0385166004820152602401610158565b50805b9392505050565b8051156102ec5780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b80516001600160a01b038116811461031c57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60005b8381101561035257818101518382015260200161033a565b50506000910152565b60008060006060848603121561037057600080fd5b61037984610305565b925061038760208501610305565b60408501519092506001600160401b03808211156103a457600080fd5b818601915086601f8301126103b857600080fd5b8151818111156103ca576103ca610321565b604051601f8201601f19908116603f011681019083821181831017156103f2576103f2610321565b8160405282815289602084870101111561040b57600080fd5b61041c836020830160208801610337565b80955050505050509250925092565b6000825161043d818460208701610337565b9190910192915050565b608051610473620004626000396000601001526104736000f3fe608060405261000c61000e565b005b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316330361006a576000356001600160e01b03191663278f794360e11b146100625761006061006e565b565b61006061007e565b6100605b6100606100796100ad565b6100d3565b60008061008e36600481846102cb565b81019061009b919061030b565b915091506100a982826100f7565b5050565b60006100ce60008051602061041e833981519152546001600160a01b031690565b905090565b3660008037600080366000845af43d6000803e8080156100f2573d6000f35b3d6000fd5b61010082610152565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a280511561014a5761014582826101b7565b505050565b6100a961022d565b806001600160a01b03163b6000036101885780604051634c9c8ce360e01b815260040161017f91906103da565b60405180910390fd5b60008051602061041e83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080846001600160a01b0316846040516101d491906103ee565b600060405180830381855af49150503d806000811461020f576040519150601f19603f3d011682016040523d82523d6000602084013e610214565b606091505b509150915061022485838361024c565b95945050505050565b34156100605760405163b398979f60e01b815260040160405180910390fd5b6060826102615761025c826102a2565b61029b565b815115801561027857506001600160a01b0384163b155b156102985783604051639996b31560e01b815260040161017f91906103da565b50805b9392505050565b8051156102b25780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b600080858511156102db57600080fd5b838611156102e857600080fd5b5050820193919092039150565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561031e57600080fd5b82356001600160a01b038116811461033557600080fd5b915060208301356001600160401b038082111561035157600080fd5b818501915085601f83011261036557600080fd5b813581811115610377576103776102f5565b604051601f8201601f19908116603f0116810190838211818310171561039f5761039f6102f5565b816040528281528860208487010111156103b857600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b6001600160a01b0391909116815260200190565b6000825160005b8181101561040f57602081860181015185830152016103f5565b50600092019182525091905056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbca2646970667358221220a19e7a72432195d9a35c7ce8fa5f1284415aac66bb1ad08a4c2e1c252fd8690864736f6c63430008140033b53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d610373cb0569fdbea2544dae03fdb2fe10eda92a72a2e8cd2bd496e85b762505a3f066156603fe29d13f97c6f3e3dff4ef71919f9aa61c555be0182d954e94221aac8cf807f6970720f8e2c208c7c5037595982c7bd9ed93c380d09df743d0dcc3fbab66e11c4f712cd06ab11bf9339b48bef39e12d4a22eeef71d2860a0c90482bdac75d24dbb35ea80e25fab167da4dea46c1915260426570db84f184891f5f59062ba6ba2ffed8cfe316b583325ea41ac6e7ba9e5864d2bc6fabba7ac26d2f0f430644e72e131a029b85045b68181585d2833e84879b9709143e1f593f00000017f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024983dfe277d2a2c04b75fb2eb3743fa00005ae3678a20c299e65fdf4df76517f68ea5c5790f581d443ed43873ab47cfb8c5d66a6db268e58b5971bb33fc66e07db19b6f082d8d3644ae2f24a3c32e356d6f2d9b2844d9b26164fbc82663ff285951141f8f32ce6198eee741f695cec728bfd32d289f1acf73621fb303581000545ea0fab074aba36a6fa69f1a83ee86e5abfb8433966eb57efb13dc2fc2f24ddd08084e94f375e9d647f87f5b2ceffba1e062c70f6009fdbcf80291e803b5c9edd4a2646970667358221220f0c277530fc86ddd05a56918423fafab9c54bee2881e81568f1fab8a3f325a3d64736f6c63430008140033",
}

// MockpolygonrollupmanagerABI is the input ABI used to generate the binding from.
// Deprecated: Use MockpolygonrollupmanagerMetaData.ABI instead.
var MockpolygonrollupmanagerABI = MockpolygonrollupmanagerMetaData.ABI

// MockpolygonrollupmanagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockpolygonrollupmanagerMetaData.Bin instead.
var MockpolygonrollupmanagerBin = MockpolygonrollupmanagerMetaData.Bin

// DeployMockpolygonrollupmanager deploys a new Ethereum contract, binding an instance of Mockpolygonrollupmanager to it.
func DeployMockpolygonrollupmanager(auth *bind.TransactOpts, backend bind.ContractBackend, _globalExitRootManager common.Address, _pol common.Address, _bridgeAddress common.Address) (common.Address, *types.Transaction, *Mockpolygonrollupmanager, error) {
	parsed, err := MockpolygonrollupmanagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockpolygonrollupmanagerBin), backend, _globalExitRootManager, _pol, _bridgeAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Mockpolygonrollupmanager{MockpolygonrollupmanagerCaller: MockpolygonrollupmanagerCaller{contract: contract}, MockpolygonrollupmanagerTransactor: MockpolygonrollupmanagerTransactor{contract: contract}, MockpolygonrollupmanagerFilterer: MockpolygonrollupmanagerFilterer{contract: contract}}, nil
}

// Mockpolygonrollupmanager is an auto generated Go binding around an Ethereum contract.
type Mockpolygonrollupmanager struct {
	MockpolygonrollupmanagerCaller     // Read-only binding to the contract
	MockpolygonrollupmanagerTransactor // Write-only binding to the contract
	MockpolygonrollupmanagerFilterer   // Log filterer for contract events
}

// MockpolygonrollupmanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockpolygonrollupmanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockpolygonrollupmanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockpolygonrollupmanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockpolygonrollupmanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockpolygonrollupmanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockpolygonrollupmanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockpolygonrollupmanagerSession struct {
	Contract     *Mockpolygonrollupmanager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// MockpolygonrollupmanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockpolygonrollupmanagerCallerSession struct {
	Contract *MockpolygonrollupmanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// MockpolygonrollupmanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockpolygonrollupmanagerTransactorSession struct {
	Contract     *MockpolygonrollupmanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// MockpolygonrollupmanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockpolygonrollupmanagerRaw struct {
	Contract *Mockpolygonrollupmanager // Generic contract binding to access the raw methods on
}

// MockpolygonrollupmanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockpolygonrollupmanagerCallerRaw struct {
	Contract *MockpolygonrollupmanagerCaller // Generic read-only contract binding to access the raw methods on
}

// MockpolygonrollupmanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockpolygonrollupmanagerTransactorRaw struct {
	Contract *MockpolygonrollupmanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockpolygonrollupmanager creates a new instance of Mockpolygonrollupmanager, bound to a specific deployed contract.
func NewMockpolygonrollupmanager(address common.Address, backend bind.ContractBackend) (*Mockpolygonrollupmanager, error) {
	contract, err := bindMockpolygonrollupmanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Mockpolygonrollupmanager{MockpolygonrollupmanagerCaller: MockpolygonrollupmanagerCaller{contract: contract}, MockpolygonrollupmanagerTransactor: MockpolygonrollupmanagerTransactor{contract: contract}, MockpolygonrollupmanagerFilterer: MockpolygonrollupmanagerFilterer{contract: contract}}, nil
}

// NewMockpolygonrollupmanagerCaller creates a new read-only instance of Mockpolygonrollupmanager, bound to a specific deployed contract.
func NewMockpolygonrollupmanagerCaller(address common.Address, caller bind.ContractCaller) (*MockpolygonrollupmanagerCaller, error) {
	contract, err := bindMockpolygonrollupmanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerCaller{contract: contract}, nil
}

// NewMockpolygonrollupmanagerTransactor creates a new write-only instance of Mockpolygonrollupmanager, bound to a specific deployed contract.
func NewMockpolygonrollupmanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*MockpolygonrollupmanagerTransactor, error) {
	contract, err := bindMockpolygonrollupmanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerTransactor{contract: contract}, nil
}

// NewMockpolygonrollupmanagerFilterer creates a new log filterer instance of Mockpolygonrollupmanager, bound to a specific deployed contract.
func NewMockpolygonrollupmanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*MockpolygonrollupmanagerFilterer, error) {
	contract, err := bindMockpolygonrollupmanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerFilterer{contract: contract}, nil
}

// bindMockpolygonrollupmanager binds a generic wrapper to an already deployed contract.
func bindMockpolygonrollupmanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockpolygonrollupmanagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mockpolygonrollupmanager.Contract.MockpolygonrollupmanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.MockpolygonrollupmanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.MockpolygonrollupmanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mockpolygonrollupmanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Mockpolygonrollupmanager.Contract.DEFAULTADMINROLE(&_Mockpolygonrollupmanager.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Mockpolygonrollupmanager.Contract.DEFAULTADMINROLE(&_Mockpolygonrollupmanager.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) BridgeAddress() (common.Address, error) {
	return _Mockpolygonrollupmanager.Contract.BridgeAddress(&_Mockpolygonrollupmanager.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) BridgeAddress() (common.Address, error) {
	return _Mockpolygonrollupmanager.Contract.BridgeAddress(&_Mockpolygonrollupmanager.CallOpts)
}

// CalculateRewardPerBatch is a free data retrieval call binding the contract method 0x99f5634e.
//
// Solidity: function calculateRewardPerBatch() view returns(uint256)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) CalculateRewardPerBatch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "calculateRewardPerBatch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateRewardPerBatch is a free data retrieval call binding the contract method 0x99f5634e.
//
// Solidity: function calculateRewardPerBatch() view returns(uint256)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) CalculateRewardPerBatch() (*big.Int, error) {
	return _Mockpolygonrollupmanager.Contract.CalculateRewardPerBatch(&_Mockpolygonrollupmanager.CallOpts)
}

// CalculateRewardPerBatch is a free data retrieval call binding the contract method 0x99f5634e.
//
// Solidity: function calculateRewardPerBatch() view returns(uint256)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) CalculateRewardPerBatch() (*big.Int, error) {
	return _Mockpolygonrollupmanager.Contract.CalculateRewardPerBatch(&_Mockpolygonrollupmanager.CallOpts)
}

// ChainIDToRollupID is a free data retrieval call binding the contract method 0x7fb6e76a.
//
// Solidity: function chainIDToRollupID(uint64 chainID) view returns(uint32 rollupID)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) ChainIDToRollupID(opts *bind.CallOpts, chainID uint64) (uint32, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "chainIDToRollupID", chainID)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// ChainIDToRollupID is a free data retrieval call binding the contract method 0x7fb6e76a.
//
// Solidity: function chainIDToRollupID(uint64 chainID) view returns(uint32 rollupID)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) ChainIDToRollupID(chainID uint64) (uint32, error) {
	return _Mockpolygonrollupmanager.Contract.ChainIDToRollupID(&_Mockpolygonrollupmanager.CallOpts, chainID)
}

// ChainIDToRollupID is a free data retrieval call binding the contract method 0x7fb6e76a.
//
// Solidity: function chainIDToRollupID(uint64 chainID) view returns(uint32 rollupID)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) ChainIDToRollupID(chainID uint64) (uint32, error) {
	return _Mockpolygonrollupmanager.Contract.ChainIDToRollupID(&_Mockpolygonrollupmanager.CallOpts, chainID)
}

// GetBatchFee is a free data retrieval call binding the contract method 0x477fa270.
//
// Solidity: function getBatchFee() view returns(uint256)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) GetBatchFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "getBatchFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBatchFee is a free data retrieval call binding the contract method 0x477fa270.
//
// Solidity: function getBatchFee() view returns(uint256)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) GetBatchFee() (*big.Int, error) {
	return _Mockpolygonrollupmanager.Contract.GetBatchFee(&_Mockpolygonrollupmanager.CallOpts)
}

// GetBatchFee is a free data retrieval call binding the contract method 0x477fa270.
//
// Solidity: function getBatchFee() view returns(uint256)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) GetBatchFee() (*big.Int, error) {
	return _Mockpolygonrollupmanager.Contract.GetBatchFee(&_Mockpolygonrollupmanager.CallOpts)
}

// GetForcedBatchFee is a free data retrieval call binding the contract method 0x60469169.
//
// Solidity: function getForcedBatchFee() view returns(uint256)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) GetForcedBatchFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "getForcedBatchFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetForcedBatchFee is a free data retrieval call binding the contract method 0x60469169.
//
// Solidity: function getForcedBatchFee() view returns(uint256)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) GetForcedBatchFee() (*big.Int, error) {
	return _Mockpolygonrollupmanager.Contract.GetForcedBatchFee(&_Mockpolygonrollupmanager.CallOpts)
}

// GetForcedBatchFee is a free data retrieval call binding the contract method 0x60469169.
//
// Solidity: function getForcedBatchFee() view returns(uint256)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) GetForcedBatchFee() (*big.Int, error) {
	return _Mockpolygonrollupmanager.Contract.GetForcedBatchFee(&_Mockpolygonrollupmanager.CallOpts)
}

// GetInputSnarkBytes is a free data retrieval call binding the contract method 0x7975fcfe.
//
// Solidity: function getInputSnarkBytes(uint32 rollupID, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 oldStateRoot, bytes32 newStateRoot) view returns(bytes)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) GetInputSnarkBytes(opts *bind.CallOpts, rollupID uint32, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, oldStateRoot [32]byte, newStateRoot [32]byte) ([]byte, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "getInputSnarkBytes", rollupID, initNumBatch, finalNewBatch, newLocalExitRoot, oldStateRoot, newStateRoot)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetInputSnarkBytes is a free data retrieval call binding the contract method 0x7975fcfe.
//
// Solidity: function getInputSnarkBytes(uint32 rollupID, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 oldStateRoot, bytes32 newStateRoot) view returns(bytes)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) GetInputSnarkBytes(rollupID uint32, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, oldStateRoot [32]byte, newStateRoot [32]byte) ([]byte, error) {
	return _Mockpolygonrollupmanager.Contract.GetInputSnarkBytes(&_Mockpolygonrollupmanager.CallOpts, rollupID, initNumBatch, finalNewBatch, newLocalExitRoot, oldStateRoot, newStateRoot)
}

// GetInputSnarkBytes is a free data retrieval call binding the contract method 0x7975fcfe.
//
// Solidity: function getInputSnarkBytes(uint32 rollupID, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 oldStateRoot, bytes32 newStateRoot) view returns(bytes)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) GetInputSnarkBytes(rollupID uint32, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, oldStateRoot [32]byte, newStateRoot [32]byte) ([]byte, error) {
	return _Mockpolygonrollupmanager.Contract.GetInputSnarkBytes(&_Mockpolygonrollupmanager.CallOpts, rollupID, initNumBatch, finalNewBatch, newLocalExitRoot, oldStateRoot, newStateRoot)
}

// GetLastVerifiedBatch is a free data retrieval call binding the contract method 0x11f6b287.
//
// Solidity: function getLastVerifiedBatch(uint32 rollupID) view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) GetLastVerifiedBatch(opts *bind.CallOpts, rollupID uint32) (uint64, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "getLastVerifiedBatch", rollupID)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetLastVerifiedBatch is a free data retrieval call binding the contract method 0x11f6b287.
//
// Solidity: function getLastVerifiedBatch(uint32 rollupID) view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) GetLastVerifiedBatch(rollupID uint32) (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.GetLastVerifiedBatch(&_Mockpolygonrollupmanager.CallOpts, rollupID)
}

// GetLastVerifiedBatch is a free data retrieval call binding the contract method 0x11f6b287.
//
// Solidity: function getLastVerifiedBatch(uint32 rollupID) view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) GetLastVerifiedBatch(rollupID uint32) (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.GetLastVerifiedBatch(&_Mockpolygonrollupmanager.CallOpts, rollupID)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Mockpolygonrollupmanager.Contract.GetRoleAdmin(&_Mockpolygonrollupmanager.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Mockpolygonrollupmanager.Contract.GetRoleAdmin(&_Mockpolygonrollupmanager.CallOpts, role)
}

// GetRollupBatchNumToStateRoot is a free data retrieval call binding the contract method 0x55a71ee0.
//
// Solidity: function getRollupBatchNumToStateRoot(uint32 rollupID, uint64 batchNum) view returns(bytes32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) GetRollupBatchNumToStateRoot(opts *bind.CallOpts, rollupID uint32, batchNum uint64) ([32]byte, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "getRollupBatchNumToStateRoot", rollupID, batchNum)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRollupBatchNumToStateRoot is a free data retrieval call binding the contract method 0x55a71ee0.
//
// Solidity: function getRollupBatchNumToStateRoot(uint32 rollupID, uint64 batchNum) view returns(bytes32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) GetRollupBatchNumToStateRoot(rollupID uint32, batchNum uint64) ([32]byte, error) {
	return _Mockpolygonrollupmanager.Contract.GetRollupBatchNumToStateRoot(&_Mockpolygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupBatchNumToStateRoot is a free data retrieval call binding the contract method 0x55a71ee0.
//
// Solidity: function getRollupBatchNumToStateRoot(uint32 rollupID, uint64 batchNum) view returns(bytes32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) GetRollupBatchNumToStateRoot(rollupID uint32, batchNum uint64) ([32]byte, error) {
	return _Mockpolygonrollupmanager.Contract.GetRollupBatchNumToStateRoot(&_Mockpolygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupExitRoot is a free data retrieval call binding the contract method 0xa2967d99.
//
// Solidity: function getRollupExitRoot() view returns(bytes32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) GetRollupExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "getRollupExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRollupExitRoot is a free data retrieval call binding the contract method 0xa2967d99.
//
// Solidity: function getRollupExitRoot() view returns(bytes32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) GetRollupExitRoot() ([32]byte, error) {
	return _Mockpolygonrollupmanager.Contract.GetRollupExitRoot(&_Mockpolygonrollupmanager.CallOpts)
}

// GetRollupExitRoot is a free data retrieval call binding the contract method 0xa2967d99.
//
// Solidity: function getRollupExitRoot() view returns(bytes32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) GetRollupExitRoot() ([32]byte, error) {
	return _Mockpolygonrollupmanager.Contract.GetRollupExitRoot(&_Mockpolygonrollupmanager.CallOpts)
}

// GetRollupPendingStateTransitions is a free data retrieval call binding the contract method 0xb99d0ad7.
//
// Solidity: function getRollupPendingStateTransitions(uint32 rollupID, uint64 batchNum) view returns((uint64,uint64,bytes32,bytes32))
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) GetRollupPendingStateTransitions(opts *bind.CallOpts, rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesPendingState, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "getRollupPendingStateTransitions", rollupID, batchNum)

	if err != nil {
		return *new(LegacyZKEVMStateVariablesPendingState), err
	}

	out0 := *abi.ConvertType(out[0], new(LegacyZKEVMStateVariablesPendingState)).(*LegacyZKEVMStateVariablesPendingState)

	return out0, err

}

// GetRollupPendingStateTransitions is a free data retrieval call binding the contract method 0xb99d0ad7.
//
// Solidity: function getRollupPendingStateTransitions(uint32 rollupID, uint64 batchNum) view returns((uint64,uint64,bytes32,bytes32))
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) GetRollupPendingStateTransitions(rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesPendingState, error) {
	return _Mockpolygonrollupmanager.Contract.GetRollupPendingStateTransitions(&_Mockpolygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupPendingStateTransitions is a free data retrieval call binding the contract method 0xb99d0ad7.
//
// Solidity: function getRollupPendingStateTransitions(uint32 rollupID, uint64 batchNum) view returns((uint64,uint64,bytes32,bytes32))
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) GetRollupPendingStateTransitions(rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesPendingState, error) {
	return _Mockpolygonrollupmanager.Contract.GetRollupPendingStateTransitions(&_Mockpolygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupSequencedBatches is a free data retrieval call binding the contract method 0x25280169.
//
// Solidity: function getRollupSequencedBatches(uint32 rollupID, uint64 batchNum) view returns((bytes32,uint64,uint64))
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) GetRollupSequencedBatches(opts *bind.CallOpts, rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesSequencedBatchData, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "getRollupSequencedBatches", rollupID, batchNum)

	if err != nil {
		return *new(LegacyZKEVMStateVariablesSequencedBatchData), err
	}

	out0 := *abi.ConvertType(out[0], new(LegacyZKEVMStateVariablesSequencedBatchData)).(*LegacyZKEVMStateVariablesSequencedBatchData)

	return out0, err

}

// GetRollupSequencedBatches is a free data retrieval call binding the contract method 0x25280169.
//
// Solidity: function getRollupSequencedBatches(uint32 rollupID, uint64 batchNum) view returns((bytes32,uint64,uint64))
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) GetRollupSequencedBatches(rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesSequencedBatchData, error) {
	return _Mockpolygonrollupmanager.Contract.GetRollupSequencedBatches(&_Mockpolygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupSequencedBatches is a free data retrieval call binding the contract method 0x25280169.
//
// Solidity: function getRollupSequencedBatches(uint32 rollupID, uint64 batchNum) view returns((bytes32,uint64,uint64))
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) GetRollupSequencedBatches(rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesSequencedBatchData, error) {
	return _Mockpolygonrollupmanager.Contract.GetRollupSequencedBatches(&_Mockpolygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) GlobalExitRootManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "globalExitRootManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) GlobalExitRootManager() (common.Address, error) {
	return _Mockpolygonrollupmanager.Contract.GlobalExitRootManager(&_Mockpolygonrollupmanager.CallOpts)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) GlobalExitRootManager() (common.Address, error) {
	return _Mockpolygonrollupmanager.Contract.GlobalExitRootManager(&_Mockpolygonrollupmanager.CallOpts)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Mockpolygonrollupmanager.Contract.HasRole(&_Mockpolygonrollupmanager.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Mockpolygonrollupmanager.Contract.HasRole(&_Mockpolygonrollupmanager.CallOpts, role, account)
}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) IsEmergencyState(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "isEmergencyState")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) IsEmergencyState() (bool, error) {
	return _Mockpolygonrollupmanager.Contract.IsEmergencyState(&_Mockpolygonrollupmanager.CallOpts)
}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) IsEmergencyState() (bool, error) {
	return _Mockpolygonrollupmanager.Contract.IsEmergencyState(&_Mockpolygonrollupmanager.CallOpts)
}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x080b3111.
//
// Solidity: function isPendingStateConsolidable(uint32 rollupID, uint64 pendingStateNum) view returns(bool)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) IsPendingStateConsolidable(opts *bind.CallOpts, rollupID uint32, pendingStateNum uint64) (bool, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "isPendingStateConsolidable", rollupID, pendingStateNum)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x080b3111.
//
// Solidity: function isPendingStateConsolidable(uint32 rollupID, uint64 pendingStateNum) view returns(bool)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) IsPendingStateConsolidable(rollupID uint32, pendingStateNum uint64) (bool, error) {
	return _Mockpolygonrollupmanager.Contract.IsPendingStateConsolidable(&_Mockpolygonrollupmanager.CallOpts, rollupID, pendingStateNum)
}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x080b3111.
//
// Solidity: function isPendingStateConsolidable(uint32 rollupID, uint64 pendingStateNum) view returns(bool)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) IsPendingStateConsolidable(rollupID uint32, pendingStateNum uint64) (bool, error) {
	return _Mockpolygonrollupmanager.Contract.IsPendingStateConsolidable(&_Mockpolygonrollupmanager.CallOpts, rollupID, pendingStateNum)
}

// LastAggregationTimestamp is a free data retrieval call binding the contract method 0xc1acbc34.
//
// Solidity: function lastAggregationTimestamp() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) LastAggregationTimestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "lastAggregationTimestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastAggregationTimestamp is a free data retrieval call binding the contract method 0xc1acbc34.
//
// Solidity: function lastAggregationTimestamp() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) LastAggregationTimestamp() (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.LastAggregationTimestamp(&_Mockpolygonrollupmanager.CallOpts)
}

// LastAggregationTimestamp is a free data retrieval call binding the contract method 0xc1acbc34.
//
// Solidity: function lastAggregationTimestamp() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) LastAggregationTimestamp() (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.LastAggregationTimestamp(&_Mockpolygonrollupmanager.CallOpts)
}

// MultiplierBatchFee is a free data retrieval call binding the contract method 0xafd23cbe.
//
// Solidity: function multiplierBatchFee() view returns(uint16)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) MultiplierBatchFee(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "multiplierBatchFee")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// MultiplierBatchFee is a free data retrieval call binding the contract method 0xafd23cbe.
//
// Solidity: function multiplierBatchFee() view returns(uint16)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) MultiplierBatchFee() (uint16, error) {
	return _Mockpolygonrollupmanager.Contract.MultiplierBatchFee(&_Mockpolygonrollupmanager.CallOpts)
}

// MultiplierBatchFee is a free data retrieval call binding the contract method 0xafd23cbe.
//
// Solidity: function multiplierBatchFee() view returns(uint16)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) MultiplierBatchFee() (uint16, error) {
	return _Mockpolygonrollupmanager.Contract.MultiplierBatchFee(&_Mockpolygonrollupmanager.CallOpts)
}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) PendingStateTimeout(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "pendingStateTimeout")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) PendingStateTimeout() (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.PendingStateTimeout(&_Mockpolygonrollupmanager.CallOpts)
}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) PendingStateTimeout() (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.PendingStateTimeout(&_Mockpolygonrollupmanager.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) Pol(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "pol")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) Pol() (common.Address, error) {
	return _Mockpolygonrollupmanager.Contract.Pol(&_Mockpolygonrollupmanager.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) Pol() (common.Address, error) {
	return _Mockpolygonrollupmanager.Contract.Pol(&_Mockpolygonrollupmanager.CallOpts)
}

// RollupAddressToID is a free data retrieval call binding the contract method 0xceee281d.
//
// Solidity: function rollupAddressToID(address rollupAddress) view returns(uint32 rollupID)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) RollupAddressToID(opts *bind.CallOpts, rollupAddress common.Address) (uint32, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "rollupAddressToID", rollupAddress)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RollupAddressToID is a free data retrieval call binding the contract method 0xceee281d.
//
// Solidity: function rollupAddressToID(address rollupAddress) view returns(uint32 rollupID)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) RollupAddressToID(rollupAddress common.Address) (uint32, error) {
	return _Mockpolygonrollupmanager.Contract.RollupAddressToID(&_Mockpolygonrollupmanager.CallOpts, rollupAddress)
}

// RollupAddressToID is a free data retrieval call binding the contract method 0xceee281d.
//
// Solidity: function rollupAddressToID(address rollupAddress) view returns(uint32 rollupID)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) RollupAddressToID(rollupAddress common.Address) (uint32, error) {
	return _Mockpolygonrollupmanager.Contract.RollupAddressToID(&_Mockpolygonrollupmanager.CallOpts, rollupAddress)
}

// RollupCount is a free data retrieval call binding the contract method 0xf4e92675.
//
// Solidity: function rollupCount() view returns(uint32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) RollupCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "rollupCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RollupCount is a free data retrieval call binding the contract method 0xf4e92675.
//
// Solidity: function rollupCount() view returns(uint32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) RollupCount() (uint32, error) {
	return _Mockpolygonrollupmanager.Contract.RollupCount(&_Mockpolygonrollupmanager.CallOpts)
}

// RollupCount is a free data retrieval call binding the contract method 0xf4e92675.
//
// Solidity: function rollupCount() view returns(uint32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) RollupCount() (uint32, error) {
	return _Mockpolygonrollupmanager.Contract.RollupCount(&_Mockpolygonrollupmanager.CallOpts)
}

// RollupIDToRollupData is a free data retrieval call binding the contract method 0xf9c4c2ae.
//
// Solidity: function rollupIDToRollupData(uint32 rollupID) view returns(address rollupContract, uint64 chainID, address verifier, uint64 forkID, bytes32 lastLocalExitRoot, uint64 lastBatchSequenced, uint64 lastVerifiedBatch, uint64 lastPendingState, uint64 lastPendingStateConsolidated, uint64 lastVerifiedBatchBeforeUpgrade, uint64 rollupTypeID, uint8 rollupCompatibilityID)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) RollupIDToRollupData(opts *bind.CallOpts, rollupID uint32) (struct {
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
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "rollupIDToRollupData", rollupID)

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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) RollupIDToRollupData(rollupID uint32) (struct {
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
	return _Mockpolygonrollupmanager.Contract.RollupIDToRollupData(&_Mockpolygonrollupmanager.CallOpts, rollupID)
}

// RollupIDToRollupData is a free data retrieval call binding the contract method 0xf9c4c2ae.
//
// Solidity: function rollupIDToRollupData(uint32 rollupID) view returns(address rollupContract, uint64 chainID, address verifier, uint64 forkID, bytes32 lastLocalExitRoot, uint64 lastBatchSequenced, uint64 lastVerifiedBatch, uint64 lastPendingState, uint64 lastPendingStateConsolidated, uint64 lastVerifiedBatchBeforeUpgrade, uint64 rollupTypeID, uint8 rollupCompatibilityID)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) RollupIDToRollupData(rollupID uint32) (struct {
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
	return _Mockpolygonrollupmanager.Contract.RollupIDToRollupData(&_Mockpolygonrollupmanager.CallOpts, rollupID)
}

// RollupTypeCount is a free data retrieval call binding the contract method 0x1796a1ae.
//
// Solidity: function rollupTypeCount() view returns(uint32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) RollupTypeCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "rollupTypeCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RollupTypeCount is a free data retrieval call binding the contract method 0x1796a1ae.
//
// Solidity: function rollupTypeCount() view returns(uint32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) RollupTypeCount() (uint32, error) {
	return _Mockpolygonrollupmanager.Contract.RollupTypeCount(&_Mockpolygonrollupmanager.CallOpts)
}

// RollupTypeCount is a free data retrieval call binding the contract method 0x1796a1ae.
//
// Solidity: function rollupTypeCount() view returns(uint32)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) RollupTypeCount() (uint32, error) {
	return _Mockpolygonrollupmanager.Contract.RollupTypeCount(&_Mockpolygonrollupmanager.CallOpts)
}

// RollupTypeMap is a free data retrieval call binding the contract method 0x65c0504d.
//
// Solidity: function rollupTypeMap(uint32 rollupTypeID) view returns(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bool obsolete, bytes32 genesis)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) RollupTypeMap(opts *bind.CallOpts, rollupTypeID uint32) (struct {
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Obsolete                bool
	Genesis                 [32]byte
}, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "rollupTypeMap", rollupTypeID)

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

// RollupTypeMap is a free data retrieval call binding the contract method 0x65c0504d.
//
// Solidity: function rollupTypeMap(uint32 rollupTypeID) view returns(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bool obsolete, bytes32 genesis)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) RollupTypeMap(rollupTypeID uint32) (struct {
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Obsolete                bool
	Genesis                 [32]byte
}, error) {
	return _Mockpolygonrollupmanager.Contract.RollupTypeMap(&_Mockpolygonrollupmanager.CallOpts, rollupTypeID)
}

// RollupTypeMap is a free data retrieval call binding the contract method 0x65c0504d.
//
// Solidity: function rollupTypeMap(uint32 rollupTypeID) view returns(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bool obsolete, bytes32 genesis)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) RollupTypeMap(rollupTypeID uint32) (struct {
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Obsolete                bool
	Genesis                 [32]byte
}, error) {
	return _Mockpolygonrollupmanager.Contract.RollupTypeMap(&_Mockpolygonrollupmanager.CallOpts, rollupTypeID)
}

// TotalSequencedBatches is a free data retrieval call binding the contract method 0x066ec012.
//
// Solidity: function totalSequencedBatches() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) TotalSequencedBatches(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "totalSequencedBatches")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// TotalSequencedBatches is a free data retrieval call binding the contract method 0x066ec012.
//
// Solidity: function totalSequencedBatches() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) TotalSequencedBatches() (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.TotalSequencedBatches(&_Mockpolygonrollupmanager.CallOpts)
}

// TotalSequencedBatches is a free data retrieval call binding the contract method 0x066ec012.
//
// Solidity: function totalSequencedBatches() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) TotalSequencedBatches() (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.TotalSequencedBatches(&_Mockpolygonrollupmanager.CallOpts)
}

// TotalVerifiedBatches is a free data retrieval call binding the contract method 0xdde0ff77.
//
// Solidity: function totalVerifiedBatches() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) TotalVerifiedBatches(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "totalVerifiedBatches")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// TotalVerifiedBatches is a free data retrieval call binding the contract method 0xdde0ff77.
//
// Solidity: function totalVerifiedBatches() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) TotalVerifiedBatches() (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.TotalVerifiedBatches(&_Mockpolygonrollupmanager.CallOpts)
}

// TotalVerifiedBatches is a free data retrieval call binding the contract method 0xdde0ff77.
//
// Solidity: function totalVerifiedBatches() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) TotalVerifiedBatches() (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.TotalVerifiedBatches(&_Mockpolygonrollupmanager.CallOpts)
}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) TrustedAggregatorTimeout(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "trustedAggregatorTimeout")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) TrustedAggregatorTimeout() (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.TrustedAggregatorTimeout(&_Mockpolygonrollupmanager.CallOpts)
}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) TrustedAggregatorTimeout() (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.TrustedAggregatorTimeout(&_Mockpolygonrollupmanager.CallOpts)
}

// VerifyBatchTimeTarget is a free data retrieval call binding the contract method 0x0a0d9fbe.
//
// Solidity: function verifyBatchTimeTarget() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCaller) VerifyBatchTimeTarget(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Mockpolygonrollupmanager.contract.Call(opts, &out, "verifyBatchTimeTarget")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// VerifyBatchTimeTarget is a free data retrieval call binding the contract method 0x0a0d9fbe.
//
// Solidity: function verifyBatchTimeTarget() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) VerifyBatchTimeTarget() (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.VerifyBatchTimeTarget(&_Mockpolygonrollupmanager.CallOpts)
}

// VerifyBatchTimeTarget is a free data retrieval call binding the contract method 0x0a0d9fbe.
//
// Solidity: function verifyBatchTimeTarget() view returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerCallerSession) VerifyBatchTimeTarget() (uint64, error) {
	return _Mockpolygonrollupmanager.Contract.VerifyBatchTimeTarget(&_Mockpolygonrollupmanager.CallOpts)
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x2072f6c5.
//
// Solidity: function activateEmergencyState() returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) ActivateEmergencyState(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "activateEmergencyState")
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x2072f6c5.
//
// Solidity: function activateEmergencyState() returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) ActivateEmergencyState() (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.ActivateEmergencyState(&_Mockpolygonrollupmanager.TransactOpts)
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x2072f6c5.
//
// Solidity: function activateEmergencyState() returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) ActivateEmergencyState() (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.ActivateEmergencyState(&_Mockpolygonrollupmanager.TransactOpts)
}

// AddExistingRollup is a paid mutator transaction binding the contract method 0xe0bfd3d2.
//
// Solidity: function addExistingRollup(address rollupAddress, address verifier, uint64 forkID, uint64 chainID, bytes32 genesis, uint8 rollupCompatibilityID) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) AddExistingRollup(opts *bind.TransactOpts, rollupAddress common.Address, verifier common.Address, forkID uint64, chainID uint64, genesis [32]byte, rollupCompatibilityID uint8) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "addExistingRollup", rollupAddress, verifier, forkID, chainID, genesis, rollupCompatibilityID)
}

// AddExistingRollup is a paid mutator transaction binding the contract method 0xe0bfd3d2.
//
// Solidity: function addExistingRollup(address rollupAddress, address verifier, uint64 forkID, uint64 chainID, bytes32 genesis, uint8 rollupCompatibilityID) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) AddExistingRollup(rollupAddress common.Address, verifier common.Address, forkID uint64, chainID uint64, genesis [32]byte, rollupCompatibilityID uint8) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.AddExistingRollup(&_Mockpolygonrollupmanager.TransactOpts, rollupAddress, verifier, forkID, chainID, genesis, rollupCompatibilityID)
}

// AddExistingRollup is a paid mutator transaction binding the contract method 0xe0bfd3d2.
//
// Solidity: function addExistingRollup(address rollupAddress, address verifier, uint64 forkID, uint64 chainID, bytes32 genesis, uint8 rollupCompatibilityID) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) AddExistingRollup(rollupAddress common.Address, verifier common.Address, forkID uint64, chainID uint64, genesis [32]byte, rollupCompatibilityID uint8) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.AddExistingRollup(&_Mockpolygonrollupmanager.TransactOpts, rollupAddress, verifier, forkID, chainID, genesis, rollupCompatibilityID)
}

// AddNewRollupType is a paid mutator transaction binding the contract method 0xf34eb8eb.
//
// Solidity: function addNewRollupType(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) AddNewRollupType(opts *bind.TransactOpts, consensusImplementation common.Address, verifier common.Address, forkID uint64, rollupCompatibilityID uint8, genesis [32]byte, description string) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "addNewRollupType", consensusImplementation, verifier, forkID, rollupCompatibilityID, genesis, description)
}

// AddNewRollupType is a paid mutator transaction binding the contract method 0xf34eb8eb.
//
// Solidity: function addNewRollupType(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) AddNewRollupType(consensusImplementation common.Address, verifier common.Address, forkID uint64, rollupCompatibilityID uint8, genesis [32]byte, description string) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.AddNewRollupType(&_Mockpolygonrollupmanager.TransactOpts, consensusImplementation, verifier, forkID, rollupCompatibilityID, genesis, description)
}

// AddNewRollupType is a paid mutator transaction binding the contract method 0xf34eb8eb.
//
// Solidity: function addNewRollupType(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) AddNewRollupType(consensusImplementation common.Address, verifier common.Address, forkID uint64, rollupCompatibilityID uint8, genesis [32]byte, description string) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.AddNewRollupType(&_Mockpolygonrollupmanager.TransactOpts, consensusImplementation, verifier, forkID, rollupCompatibilityID, genesis, description)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x1608859c.
//
// Solidity: function consolidatePendingState(uint32 rollupID, uint64 pendingStateNum) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) ConsolidatePendingState(opts *bind.TransactOpts, rollupID uint32, pendingStateNum uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "consolidatePendingState", rollupID, pendingStateNum)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x1608859c.
//
// Solidity: function consolidatePendingState(uint32 rollupID, uint64 pendingStateNum) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) ConsolidatePendingState(rollupID uint32, pendingStateNum uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.ConsolidatePendingState(&_Mockpolygonrollupmanager.TransactOpts, rollupID, pendingStateNum)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x1608859c.
//
// Solidity: function consolidatePendingState(uint32 rollupID, uint64 pendingStateNum) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) ConsolidatePendingState(rollupID uint32, pendingStateNum uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.ConsolidatePendingState(&_Mockpolygonrollupmanager.TransactOpts, rollupID, pendingStateNum)
}

// CreateNewRollup is a paid mutator transaction binding the contract method 0x727885e9.
//
// Solidity: function createNewRollup(uint32 rollupTypeID, uint64 chainID, address admin, address sequencer, address gasTokenAddress, string sequencerURL, string networkName) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) CreateNewRollup(opts *bind.TransactOpts, rollupTypeID uint32, chainID uint64, admin common.Address, sequencer common.Address, gasTokenAddress common.Address, sequencerURL string, networkName string) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "createNewRollup", rollupTypeID, chainID, admin, sequencer, gasTokenAddress, sequencerURL, networkName)
}

// CreateNewRollup is a paid mutator transaction binding the contract method 0x727885e9.
//
// Solidity: function createNewRollup(uint32 rollupTypeID, uint64 chainID, address admin, address sequencer, address gasTokenAddress, string sequencerURL, string networkName) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) CreateNewRollup(rollupTypeID uint32, chainID uint64, admin common.Address, sequencer common.Address, gasTokenAddress common.Address, sequencerURL string, networkName string) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.CreateNewRollup(&_Mockpolygonrollupmanager.TransactOpts, rollupTypeID, chainID, admin, sequencer, gasTokenAddress, sequencerURL, networkName)
}

// CreateNewRollup is a paid mutator transaction binding the contract method 0x727885e9.
//
// Solidity: function createNewRollup(uint32 rollupTypeID, uint64 chainID, address admin, address sequencer, address gasTokenAddress, string sequencerURL, string networkName) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) CreateNewRollup(rollupTypeID uint32, chainID uint64, admin common.Address, sequencer common.Address, gasTokenAddress common.Address, sequencerURL string, networkName string) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.CreateNewRollup(&_Mockpolygonrollupmanager.TransactOpts, rollupTypeID, chainID, admin, sequencer, gasTokenAddress, sequencerURL, networkName)
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) DeactivateEmergencyState(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "deactivateEmergencyState")
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) DeactivateEmergencyState() (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.DeactivateEmergencyState(&_Mockpolygonrollupmanager.TransactOpts)
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) DeactivateEmergencyState() (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.DeactivateEmergencyState(&_Mockpolygonrollupmanager.TransactOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.GrantRole(&_Mockpolygonrollupmanager.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.GrantRole(&_Mockpolygonrollupmanager.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x0645af09.
//
// Solidity: function initialize(address trustedAggregator, uint64 _pendingStateTimeout, uint64 _trustedAggregatorTimeout, address admin, address timelock, address emergencyCouncil, address polygonZkEVM, address zkEVMVerifier, uint64 zkEVMForkID, uint64 zkEVMChainID) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) Initialize(opts *bind.TransactOpts, trustedAggregator common.Address, _pendingStateTimeout uint64, _trustedAggregatorTimeout uint64, admin common.Address, timelock common.Address, emergencyCouncil common.Address, polygonZkEVM common.Address, zkEVMVerifier common.Address, zkEVMForkID uint64, zkEVMChainID uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "initialize", trustedAggregator, _pendingStateTimeout, _trustedAggregatorTimeout, admin, timelock, emergencyCouncil, polygonZkEVM, zkEVMVerifier, zkEVMForkID, zkEVMChainID)
}

// Initialize is a paid mutator transaction binding the contract method 0x0645af09.
//
// Solidity: function initialize(address trustedAggregator, uint64 _pendingStateTimeout, uint64 _trustedAggregatorTimeout, address admin, address timelock, address emergencyCouncil, address polygonZkEVM, address zkEVMVerifier, uint64 zkEVMForkID, uint64 zkEVMChainID) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) Initialize(trustedAggregator common.Address, _pendingStateTimeout uint64, _trustedAggregatorTimeout uint64, admin common.Address, timelock common.Address, emergencyCouncil common.Address, polygonZkEVM common.Address, zkEVMVerifier common.Address, zkEVMForkID uint64, zkEVMChainID uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.Initialize(&_Mockpolygonrollupmanager.TransactOpts, trustedAggregator, _pendingStateTimeout, _trustedAggregatorTimeout, admin, timelock, emergencyCouncil, polygonZkEVM, zkEVMVerifier, zkEVMForkID, zkEVMChainID)
}

// Initialize is a paid mutator transaction binding the contract method 0x0645af09.
//
// Solidity: function initialize(address trustedAggregator, uint64 _pendingStateTimeout, uint64 _trustedAggregatorTimeout, address admin, address timelock, address emergencyCouncil, address polygonZkEVM, address zkEVMVerifier, uint64 zkEVMForkID, uint64 zkEVMChainID) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) Initialize(trustedAggregator common.Address, _pendingStateTimeout uint64, _trustedAggregatorTimeout uint64, admin common.Address, timelock common.Address, emergencyCouncil common.Address, polygonZkEVM common.Address, zkEVMVerifier common.Address, zkEVMForkID uint64, zkEVMChainID uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.Initialize(&_Mockpolygonrollupmanager.TransactOpts, trustedAggregator, _pendingStateTimeout, _trustedAggregatorTimeout, admin, timelock, emergencyCouncil, polygonZkEVM, zkEVMVerifier, zkEVMForkID, zkEVMChainID)
}

// InitializeMock is a paid mutator transaction binding the contract method 0x0e36f582.
//
// Solidity: function initializeMock(address trustedAggregator, uint64 _pendingStateTimeout, uint64 _trustedAggregatorTimeout, address admin, address timelock, address emergencyCouncil) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) InitializeMock(opts *bind.TransactOpts, trustedAggregator common.Address, _pendingStateTimeout uint64, _trustedAggregatorTimeout uint64, admin common.Address, timelock common.Address, emergencyCouncil common.Address) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "initializeMock", trustedAggregator, _pendingStateTimeout, _trustedAggregatorTimeout, admin, timelock, emergencyCouncil)
}

// InitializeMock is a paid mutator transaction binding the contract method 0x0e36f582.
//
// Solidity: function initializeMock(address trustedAggregator, uint64 _pendingStateTimeout, uint64 _trustedAggregatorTimeout, address admin, address timelock, address emergencyCouncil) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) InitializeMock(trustedAggregator common.Address, _pendingStateTimeout uint64, _trustedAggregatorTimeout uint64, admin common.Address, timelock common.Address, emergencyCouncil common.Address) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.InitializeMock(&_Mockpolygonrollupmanager.TransactOpts, trustedAggregator, _pendingStateTimeout, _trustedAggregatorTimeout, admin, timelock, emergencyCouncil)
}

// InitializeMock is a paid mutator transaction binding the contract method 0x0e36f582.
//
// Solidity: function initializeMock(address trustedAggregator, uint64 _pendingStateTimeout, uint64 _trustedAggregatorTimeout, address admin, address timelock, address emergencyCouncil) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) InitializeMock(trustedAggregator common.Address, _pendingStateTimeout uint64, _trustedAggregatorTimeout uint64, admin common.Address, timelock common.Address, emergencyCouncil common.Address) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.InitializeMock(&_Mockpolygonrollupmanager.TransactOpts, trustedAggregator, _pendingStateTimeout, _trustedAggregatorTimeout, admin, timelock, emergencyCouncil)
}

// ObsoleteRollupType is a paid mutator transaction binding the contract method 0x7222020f.
//
// Solidity: function obsoleteRollupType(uint32 rollupTypeID) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) ObsoleteRollupType(opts *bind.TransactOpts, rollupTypeID uint32) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "obsoleteRollupType", rollupTypeID)
}

// ObsoleteRollupType is a paid mutator transaction binding the contract method 0x7222020f.
//
// Solidity: function obsoleteRollupType(uint32 rollupTypeID) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) ObsoleteRollupType(rollupTypeID uint32) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.ObsoleteRollupType(&_Mockpolygonrollupmanager.TransactOpts, rollupTypeID)
}

// ObsoleteRollupType is a paid mutator transaction binding the contract method 0x7222020f.
//
// Solidity: function obsoleteRollupType(uint32 rollupTypeID) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) ObsoleteRollupType(rollupTypeID uint32) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.ObsoleteRollupType(&_Mockpolygonrollupmanager.TransactOpts, rollupTypeID)
}

// OnSequenceBatches is a paid mutator transaction binding the contract method 0x9a908e73.
//
// Solidity: function onSequenceBatches(uint64 newSequencedBatches, bytes32 newAccInputHash) returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) OnSequenceBatches(opts *bind.TransactOpts, newSequencedBatches uint64, newAccInputHash [32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "onSequenceBatches", newSequencedBatches, newAccInputHash)
}

// OnSequenceBatches is a paid mutator transaction binding the contract method 0x9a908e73.
//
// Solidity: function onSequenceBatches(uint64 newSequencedBatches, bytes32 newAccInputHash) returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) OnSequenceBatches(newSequencedBatches uint64, newAccInputHash [32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.OnSequenceBatches(&_Mockpolygonrollupmanager.TransactOpts, newSequencedBatches, newAccInputHash)
}

// OnSequenceBatches is a paid mutator transaction binding the contract method 0x9a908e73.
//
// Solidity: function onSequenceBatches(uint64 newSequencedBatches, bytes32 newAccInputHash) returns(uint64)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) OnSequenceBatches(newSequencedBatches uint64, newAccInputHash [32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.OnSequenceBatches(&_Mockpolygonrollupmanager.TransactOpts, newSequencedBatches, newAccInputHash)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0x12b86e19.
//
// Solidity: function overridePendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) OverridePendingState(opts *bind.TransactOpts, rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "overridePendingState", rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0x12b86e19.
//
// Solidity: function overridePendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) OverridePendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.OverridePendingState(&_Mockpolygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0x12b86e19.
//
// Solidity: function overridePendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) OverridePendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.OverridePendingState(&_Mockpolygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// PrepareMockCalculateRoot is a paid mutator transaction binding the contract method 0x8f698ec5.
//
// Solidity: function prepareMockCalculateRoot(bytes32[] localExitRoots) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) PrepareMockCalculateRoot(opts *bind.TransactOpts, localExitRoots [][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "prepareMockCalculateRoot", localExitRoots)
}

// PrepareMockCalculateRoot is a paid mutator transaction binding the contract method 0x8f698ec5.
//
// Solidity: function prepareMockCalculateRoot(bytes32[] localExitRoots) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) PrepareMockCalculateRoot(localExitRoots [][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.PrepareMockCalculateRoot(&_Mockpolygonrollupmanager.TransactOpts, localExitRoots)
}

// PrepareMockCalculateRoot is a paid mutator transaction binding the contract method 0x8f698ec5.
//
// Solidity: function prepareMockCalculateRoot(bytes32[] localExitRoots) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) PrepareMockCalculateRoot(localExitRoots [][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.PrepareMockCalculateRoot(&_Mockpolygonrollupmanager.TransactOpts, localExitRoots)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x8bd4f071.
//
// Solidity: function proveNonDeterministicPendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) ProveNonDeterministicPendingState(opts *bind.TransactOpts, rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "proveNonDeterministicPendingState", rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x8bd4f071.
//
// Solidity: function proveNonDeterministicPendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) ProveNonDeterministicPendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.ProveNonDeterministicPendingState(&_Mockpolygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x8bd4f071.
//
// Solidity: function proveNonDeterministicPendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) ProveNonDeterministicPendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.ProveNonDeterministicPendingState(&_Mockpolygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.RenounceRole(&_Mockpolygonrollupmanager.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.RenounceRole(&_Mockpolygonrollupmanager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.RevokeRole(&_Mockpolygonrollupmanager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.RevokeRole(&_Mockpolygonrollupmanager.TransactOpts, role, account)
}

// SetBatchFee is a paid mutator transaction binding the contract method 0xd5073f6f.
//
// Solidity: function setBatchFee(uint256 newBatchFee) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) SetBatchFee(opts *bind.TransactOpts, newBatchFee *big.Int) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "setBatchFee", newBatchFee)
}

// SetBatchFee is a paid mutator transaction binding the contract method 0xd5073f6f.
//
// Solidity: function setBatchFee(uint256 newBatchFee) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) SetBatchFee(newBatchFee *big.Int) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.SetBatchFee(&_Mockpolygonrollupmanager.TransactOpts, newBatchFee)
}

// SetBatchFee is a paid mutator transaction binding the contract method 0xd5073f6f.
//
// Solidity: function setBatchFee(uint256 newBatchFee) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) SetBatchFee(newBatchFee *big.Int) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.SetBatchFee(&_Mockpolygonrollupmanager.TransactOpts, newBatchFee)
}

// SetMultiplierBatchFee is a paid mutator transaction binding the contract method 0x1816b7e5.
//
// Solidity: function setMultiplierBatchFee(uint16 newMultiplierBatchFee) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) SetMultiplierBatchFee(opts *bind.TransactOpts, newMultiplierBatchFee uint16) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "setMultiplierBatchFee", newMultiplierBatchFee)
}

// SetMultiplierBatchFee is a paid mutator transaction binding the contract method 0x1816b7e5.
//
// Solidity: function setMultiplierBatchFee(uint16 newMultiplierBatchFee) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) SetMultiplierBatchFee(newMultiplierBatchFee uint16) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.SetMultiplierBatchFee(&_Mockpolygonrollupmanager.TransactOpts, newMultiplierBatchFee)
}

// SetMultiplierBatchFee is a paid mutator transaction binding the contract method 0x1816b7e5.
//
// Solidity: function setMultiplierBatchFee(uint16 newMultiplierBatchFee) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) SetMultiplierBatchFee(newMultiplierBatchFee uint16) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.SetMultiplierBatchFee(&_Mockpolygonrollupmanager.TransactOpts, newMultiplierBatchFee)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) SetPendingStateTimeout(opts *bind.TransactOpts, newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "setPendingStateTimeout", newPendingStateTimeout)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) SetPendingStateTimeout(newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.SetPendingStateTimeout(&_Mockpolygonrollupmanager.TransactOpts, newPendingStateTimeout)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) SetPendingStateTimeout(newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.SetPendingStateTimeout(&_Mockpolygonrollupmanager.TransactOpts, newPendingStateTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) SetTrustedAggregatorTimeout(opts *bind.TransactOpts, newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "setTrustedAggregatorTimeout", newTrustedAggregatorTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) SetTrustedAggregatorTimeout(newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.SetTrustedAggregatorTimeout(&_Mockpolygonrollupmanager.TransactOpts, newTrustedAggregatorTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) SetTrustedAggregatorTimeout(newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.SetTrustedAggregatorTimeout(&_Mockpolygonrollupmanager.TransactOpts, newTrustedAggregatorTimeout)
}

// SetVerifyBatchTimeTarget is a paid mutator transaction binding the contract method 0xa066215c.
//
// Solidity: function setVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) SetVerifyBatchTimeTarget(opts *bind.TransactOpts, newVerifyBatchTimeTarget uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "setVerifyBatchTimeTarget", newVerifyBatchTimeTarget)
}

// SetVerifyBatchTimeTarget is a paid mutator transaction binding the contract method 0xa066215c.
//
// Solidity: function setVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) SetVerifyBatchTimeTarget(newVerifyBatchTimeTarget uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.SetVerifyBatchTimeTarget(&_Mockpolygonrollupmanager.TransactOpts, newVerifyBatchTimeTarget)
}

// SetVerifyBatchTimeTarget is a paid mutator transaction binding the contract method 0xa066215c.
//
// Solidity: function setVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) SetVerifyBatchTimeTarget(newVerifyBatchTimeTarget uint64) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.SetVerifyBatchTimeTarget(&_Mockpolygonrollupmanager.TransactOpts, newVerifyBatchTimeTarget)
}

// UpdateRollup is a paid mutator transaction binding the contract method 0xc4c928c2.
//
// Solidity: function updateRollup(address rollupContract, uint32 newRollupTypeID, bytes upgradeData) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) UpdateRollup(opts *bind.TransactOpts, rollupContract common.Address, newRollupTypeID uint32, upgradeData []byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "updateRollup", rollupContract, newRollupTypeID, upgradeData)
}

// UpdateRollup is a paid mutator transaction binding the contract method 0xc4c928c2.
//
// Solidity: function updateRollup(address rollupContract, uint32 newRollupTypeID, bytes upgradeData) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) UpdateRollup(rollupContract common.Address, newRollupTypeID uint32, upgradeData []byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.UpdateRollup(&_Mockpolygonrollupmanager.TransactOpts, rollupContract, newRollupTypeID, upgradeData)
}

// UpdateRollup is a paid mutator transaction binding the contract method 0xc4c928c2.
//
// Solidity: function updateRollup(address rollupContract, uint32 newRollupTypeID, bytes upgradeData) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) UpdateRollup(rollupContract common.Address, newRollupTypeID uint32, upgradeData []byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.UpdateRollup(&_Mockpolygonrollupmanager.TransactOpts, rollupContract, newRollupTypeID, upgradeData)
}

// VerifyBatches is a paid mutator transaction binding the contract method 0x87c20c01.
//
// Solidity: function verifyBatches(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) VerifyBatches(opts *bind.TransactOpts, rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "verifyBatches", rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatches is a paid mutator transaction binding the contract method 0x87c20c01.
//
// Solidity: function verifyBatches(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) VerifyBatches(rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.VerifyBatches(&_Mockpolygonrollupmanager.TransactOpts, rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatches is a paid mutator transaction binding the contract method 0x87c20c01.
//
// Solidity: function verifyBatches(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) VerifyBatches(rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.VerifyBatches(&_Mockpolygonrollupmanager.TransactOpts, rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatchesTrustedAggregator is a paid mutator transaction binding the contract method 0x1489ed10.
//
// Solidity: function verifyBatchesTrustedAggregator(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactor) VerifyBatchesTrustedAggregator(opts *bind.TransactOpts, rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.contract.Transact(opts, "verifyBatchesTrustedAggregator", rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatchesTrustedAggregator is a paid mutator transaction binding the contract method 0x1489ed10.
//
// Solidity: function verifyBatchesTrustedAggregator(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerSession) VerifyBatchesTrustedAggregator(rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.VerifyBatchesTrustedAggregator(&_Mockpolygonrollupmanager.TransactOpts, rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatchesTrustedAggregator is a paid mutator transaction binding the contract method 0x1489ed10.
//
// Solidity: function verifyBatchesTrustedAggregator(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerTransactorSession) VerifyBatchesTrustedAggregator(rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockpolygonrollupmanager.Contract.VerifyBatchesTrustedAggregator(&_Mockpolygonrollupmanager.TransactOpts, rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// MockpolygonrollupmanagerAddExistingRollupIterator is returned from FilterAddExistingRollup and is used to iterate over the raw logs and unpacked data for AddExistingRollup events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerAddExistingRollupIterator struct {
	Event *MockpolygonrollupmanagerAddExistingRollup // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerAddExistingRollupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerAddExistingRollup)
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
		it.Event = new(MockpolygonrollupmanagerAddExistingRollup)
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
func (it *MockpolygonrollupmanagerAddExistingRollupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerAddExistingRollupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerAddExistingRollup represents a AddExistingRollup event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerAddExistingRollup struct {
	RollupID                       uint32
	ForkID                         uint64
	RollupAddress                  common.Address
	ChainID                        uint64
	RollupCompatibilityID          uint8
	LastVerifiedBatchBeforeUpgrade uint64
	Raw                            types.Log // Blockchain specific contextual infos
}

// FilterAddExistingRollup is a free log retrieval operation binding the contract event 0xadfc7d56f7e39b08b321534f14bfb135ad27698f7d2f5ad0edc2356ea9a3f850.
//
// Solidity: event AddExistingRollup(uint32 indexed rollupID, uint64 forkID, address rollupAddress, uint64 chainID, uint8 rollupCompatibilityID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterAddExistingRollup(opts *bind.FilterOpts, rollupID []uint32) (*MockpolygonrollupmanagerAddExistingRollupIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "AddExistingRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerAddExistingRollupIterator{contract: _Mockpolygonrollupmanager.contract, event: "AddExistingRollup", logs: logs, sub: sub}, nil
}

// WatchAddExistingRollup is a free log subscription operation binding the contract event 0xadfc7d56f7e39b08b321534f14bfb135ad27698f7d2f5ad0edc2356ea9a3f850.
//
// Solidity: event AddExistingRollup(uint32 indexed rollupID, uint64 forkID, address rollupAddress, uint64 chainID, uint8 rollupCompatibilityID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchAddExistingRollup(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerAddExistingRollup, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "AddExistingRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerAddExistingRollup)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "AddExistingRollup", log); err != nil {
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

// ParseAddExistingRollup is a log parse operation binding the contract event 0xadfc7d56f7e39b08b321534f14bfb135ad27698f7d2f5ad0edc2356ea9a3f850.
//
// Solidity: event AddExistingRollup(uint32 indexed rollupID, uint64 forkID, address rollupAddress, uint64 chainID, uint8 rollupCompatibilityID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseAddExistingRollup(log types.Log) (*MockpolygonrollupmanagerAddExistingRollup, error) {
	event := new(MockpolygonrollupmanagerAddExistingRollup)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "AddExistingRollup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerAddNewRollupTypeIterator is returned from FilterAddNewRollupType and is used to iterate over the raw logs and unpacked data for AddNewRollupType events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerAddNewRollupTypeIterator struct {
	Event *MockpolygonrollupmanagerAddNewRollupType // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerAddNewRollupTypeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerAddNewRollupType)
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
		it.Event = new(MockpolygonrollupmanagerAddNewRollupType)
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
func (it *MockpolygonrollupmanagerAddNewRollupTypeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerAddNewRollupTypeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerAddNewRollupType represents a AddNewRollupType event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerAddNewRollupType struct {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterAddNewRollupType(opts *bind.FilterOpts, rollupTypeID []uint32) (*MockpolygonrollupmanagerAddNewRollupTypeIterator, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "AddNewRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerAddNewRollupTypeIterator{contract: _Mockpolygonrollupmanager.contract, event: "AddNewRollupType", logs: logs, sub: sub}, nil
}

// WatchAddNewRollupType is a free log subscription operation binding the contract event 0xa2970448b3bd66ba7e524e7b2a5b9cf94fa29e32488fb942afdfe70dd4b77b52.
//
// Solidity: event AddNewRollupType(uint32 indexed rollupTypeID, address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchAddNewRollupType(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerAddNewRollupType, rollupTypeID []uint32) (event.Subscription, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "AddNewRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerAddNewRollupType)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "AddNewRollupType", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseAddNewRollupType(log types.Log) (*MockpolygonrollupmanagerAddNewRollupType, error) {
	event := new(MockpolygonrollupmanagerAddNewRollupType)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "AddNewRollupType", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerConsolidatePendingStateIterator is returned from FilterConsolidatePendingState and is used to iterate over the raw logs and unpacked data for ConsolidatePendingState events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerConsolidatePendingStateIterator struct {
	Event *MockpolygonrollupmanagerConsolidatePendingState // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerConsolidatePendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerConsolidatePendingState)
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
		it.Event = new(MockpolygonrollupmanagerConsolidatePendingState)
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
func (it *MockpolygonrollupmanagerConsolidatePendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerConsolidatePendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerConsolidatePendingState represents a ConsolidatePendingState event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerConsolidatePendingState struct {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterConsolidatePendingState(opts *bind.FilterOpts, rollupID []uint32) (*MockpolygonrollupmanagerConsolidatePendingStateIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "ConsolidatePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerConsolidatePendingStateIterator{contract: _Mockpolygonrollupmanager.contract, event: "ConsolidatePendingState", logs: logs, sub: sub}, nil
}

// WatchConsolidatePendingState is a free log subscription operation binding the contract event 0x581910eb7a27738945c2f00a91f2284b2d6de9d4e472b12f901c2b0df045e21b.
//
// Solidity: event ConsolidatePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, uint64 pendingStateNum)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchConsolidatePendingState(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerConsolidatePendingState, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "ConsolidatePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerConsolidatePendingState)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "ConsolidatePendingState", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseConsolidatePendingState(log types.Log) (*MockpolygonrollupmanagerConsolidatePendingState, error) {
	event := new(MockpolygonrollupmanagerConsolidatePendingState)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "ConsolidatePendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerCreateNewRollupIterator is returned from FilterCreateNewRollup and is used to iterate over the raw logs and unpacked data for CreateNewRollup events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerCreateNewRollupIterator struct {
	Event *MockpolygonrollupmanagerCreateNewRollup // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerCreateNewRollupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerCreateNewRollup)
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
		it.Event = new(MockpolygonrollupmanagerCreateNewRollup)
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
func (it *MockpolygonrollupmanagerCreateNewRollupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerCreateNewRollupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerCreateNewRollup represents a CreateNewRollup event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerCreateNewRollup struct {
	RollupID        uint32
	RollupTypeID    uint32
	RollupAddress   common.Address
	ChainID         uint64
	GasTokenAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCreateNewRollup is a free log retrieval operation binding the contract event 0x194c983456df6701c6a50830b90fe80e72b823411d0d524970c9590dc277a641.
//
// Solidity: event CreateNewRollup(uint32 indexed rollupID, uint32 rollupTypeID, address rollupAddress, uint64 chainID, address gasTokenAddress)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterCreateNewRollup(opts *bind.FilterOpts, rollupID []uint32) (*MockpolygonrollupmanagerCreateNewRollupIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "CreateNewRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerCreateNewRollupIterator{contract: _Mockpolygonrollupmanager.contract, event: "CreateNewRollup", logs: logs, sub: sub}, nil
}

// WatchCreateNewRollup is a free log subscription operation binding the contract event 0x194c983456df6701c6a50830b90fe80e72b823411d0d524970c9590dc277a641.
//
// Solidity: event CreateNewRollup(uint32 indexed rollupID, uint32 rollupTypeID, address rollupAddress, uint64 chainID, address gasTokenAddress)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchCreateNewRollup(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerCreateNewRollup, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "CreateNewRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerCreateNewRollup)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "CreateNewRollup", log); err != nil {
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

// ParseCreateNewRollup is a log parse operation binding the contract event 0x194c983456df6701c6a50830b90fe80e72b823411d0d524970c9590dc277a641.
//
// Solidity: event CreateNewRollup(uint32 indexed rollupID, uint32 rollupTypeID, address rollupAddress, uint64 chainID, address gasTokenAddress)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseCreateNewRollup(log types.Log) (*MockpolygonrollupmanagerCreateNewRollup, error) {
	event := new(MockpolygonrollupmanagerCreateNewRollup)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "CreateNewRollup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerEmergencyStateActivatedIterator is returned from FilterEmergencyStateActivated and is used to iterate over the raw logs and unpacked data for EmergencyStateActivated events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerEmergencyStateActivatedIterator struct {
	Event *MockpolygonrollupmanagerEmergencyStateActivated // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerEmergencyStateActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerEmergencyStateActivated)
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
		it.Event = new(MockpolygonrollupmanagerEmergencyStateActivated)
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
func (it *MockpolygonrollupmanagerEmergencyStateActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerEmergencyStateActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerEmergencyStateActivated represents a EmergencyStateActivated event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerEmergencyStateActivated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEmergencyStateActivated is a free log retrieval operation binding the contract event 0x2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497.
//
// Solidity: event EmergencyStateActivated()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterEmergencyStateActivated(opts *bind.FilterOpts) (*MockpolygonrollupmanagerEmergencyStateActivatedIterator, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "EmergencyStateActivated")
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerEmergencyStateActivatedIterator{contract: _Mockpolygonrollupmanager.contract, event: "EmergencyStateActivated", logs: logs, sub: sub}, nil
}

// WatchEmergencyStateActivated is a free log subscription operation binding the contract event 0x2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497.
//
// Solidity: event EmergencyStateActivated()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchEmergencyStateActivated(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerEmergencyStateActivated) (event.Subscription, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "EmergencyStateActivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerEmergencyStateActivated)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "EmergencyStateActivated", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseEmergencyStateActivated(log types.Log) (*MockpolygonrollupmanagerEmergencyStateActivated, error) {
	event := new(MockpolygonrollupmanagerEmergencyStateActivated)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "EmergencyStateActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerEmergencyStateDeactivatedIterator is returned from FilterEmergencyStateDeactivated and is used to iterate over the raw logs and unpacked data for EmergencyStateDeactivated events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerEmergencyStateDeactivatedIterator struct {
	Event *MockpolygonrollupmanagerEmergencyStateDeactivated // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerEmergencyStateDeactivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerEmergencyStateDeactivated)
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
		it.Event = new(MockpolygonrollupmanagerEmergencyStateDeactivated)
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
func (it *MockpolygonrollupmanagerEmergencyStateDeactivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerEmergencyStateDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerEmergencyStateDeactivated represents a EmergencyStateDeactivated event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerEmergencyStateDeactivated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEmergencyStateDeactivated is a free log retrieval operation binding the contract event 0x1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3.
//
// Solidity: event EmergencyStateDeactivated()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterEmergencyStateDeactivated(opts *bind.FilterOpts) (*MockpolygonrollupmanagerEmergencyStateDeactivatedIterator, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "EmergencyStateDeactivated")
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerEmergencyStateDeactivatedIterator{contract: _Mockpolygonrollupmanager.contract, event: "EmergencyStateDeactivated", logs: logs, sub: sub}, nil
}

// WatchEmergencyStateDeactivated is a free log subscription operation binding the contract event 0x1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3.
//
// Solidity: event EmergencyStateDeactivated()
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchEmergencyStateDeactivated(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerEmergencyStateDeactivated) (event.Subscription, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "EmergencyStateDeactivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerEmergencyStateDeactivated)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "EmergencyStateDeactivated", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseEmergencyStateDeactivated(log types.Log) (*MockpolygonrollupmanagerEmergencyStateDeactivated, error) {
	event := new(MockpolygonrollupmanagerEmergencyStateDeactivated)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "EmergencyStateDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerInitializedIterator struct {
	Event *MockpolygonrollupmanagerInitialized // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerInitialized)
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
		it.Event = new(MockpolygonrollupmanagerInitialized)
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
func (it *MockpolygonrollupmanagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerInitialized represents a Initialized event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*MockpolygonrollupmanagerInitializedIterator, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerInitializedIterator{contract: _Mockpolygonrollupmanager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerInitialized) (event.Subscription, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerInitialized)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseInitialized(log types.Log) (*MockpolygonrollupmanagerInitialized, error) {
	event := new(MockpolygonrollupmanagerInitialized)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerObsoleteRollupTypeIterator is returned from FilterObsoleteRollupType and is used to iterate over the raw logs and unpacked data for ObsoleteRollupType events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerObsoleteRollupTypeIterator struct {
	Event *MockpolygonrollupmanagerObsoleteRollupType // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerObsoleteRollupTypeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerObsoleteRollupType)
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
		it.Event = new(MockpolygonrollupmanagerObsoleteRollupType)
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
func (it *MockpolygonrollupmanagerObsoleteRollupTypeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerObsoleteRollupTypeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerObsoleteRollupType represents a ObsoleteRollupType event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerObsoleteRollupType struct {
	RollupTypeID uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterObsoleteRollupType is a free log retrieval operation binding the contract event 0x4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e44.
//
// Solidity: event ObsoleteRollupType(uint32 indexed rollupTypeID)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterObsoleteRollupType(opts *bind.FilterOpts, rollupTypeID []uint32) (*MockpolygonrollupmanagerObsoleteRollupTypeIterator, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "ObsoleteRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerObsoleteRollupTypeIterator{contract: _Mockpolygonrollupmanager.contract, event: "ObsoleteRollupType", logs: logs, sub: sub}, nil
}

// WatchObsoleteRollupType is a free log subscription operation binding the contract event 0x4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e44.
//
// Solidity: event ObsoleteRollupType(uint32 indexed rollupTypeID)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchObsoleteRollupType(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerObsoleteRollupType, rollupTypeID []uint32) (event.Subscription, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "ObsoleteRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerObsoleteRollupType)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "ObsoleteRollupType", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseObsoleteRollupType(log types.Log) (*MockpolygonrollupmanagerObsoleteRollupType, error) {
	event := new(MockpolygonrollupmanagerObsoleteRollupType)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "ObsoleteRollupType", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerOnSequenceBatchesIterator is returned from FilterOnSequenceBatches and is used to iterate over the raw logs and unpacked data for OnSequenceBatches events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerOnSequenceBatchesIterator struct {
	Event *MockpolygonrollupmanagerOnSequenceBatches // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerOnSequenceBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerOnSequenceBatches)
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
		it.Event = new(MockpolygonrollupmanagerOnSequenceBatches)
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
func (it *MockpolygonrollupmanagerOnSequenceBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerOnSequenceBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerOnSequenceBatches represents a OnSequenceBatches event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerOnSequenceBatches struct {
	RollupID           uint32
	LastBatchSequenced uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterOnSequenceBatches is a free log retrieval operation binding the contract event 0x1d9f30260051d51d70339da239ea7b080021adcaabfa71c9b0ea339a20cf9a25.
//
// Solidity: event OnSequenceBatches(uint32 indexed rollupID, uint64 lastBatchSequenced)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterOnSequenceBatches(opts *bind.FilterOpts, rollupID []uint32) (*MockpolygonrollupmanagerOnSequenceBatchesIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "OnSequenceBatches", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerOnSequenceBatchesIterator{contract: _Mockpolygonrollupmanager.contract, event: "OnSequenceBatches", logs: logs, sub: sub}, nil
}

// WatchOnSequenceBatches is a free log subscription operation binding the contract event 0x1d9f30260051d51d70339da239ea7b080021adcaabfa71c9b0ea339a20cf9a25.
//
// Solidity: event OnSequenceBatches(uint32 indexed rollupID, uint64 lastBatchSequenced)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchOnSequenceBatches(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerOnSequenceBatches, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "OnSequenceBatches", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerOnSequenceBatches)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "OnSequenceBatches", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseOnSequenceBatches(log types.Log) (*MockpolygonrollupmanagerOnSequenceBatches, error) {
	event := new(MockpolygonrollupmanagerOnSequenceBatches)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "OnSequenceBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerOverridePendingStateIterator is returned from FilterOverridePendingState and is used to iterate over the raw logs and unpacked data for OverridePendingState events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerOverridePendingStateIterator struct {
	Event *MockpolygonrollupmanagerOverridePendingState // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerOverridePendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerOverridePendingState)
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
		it.Event = new(MockpolygonrollupmanagerOverridePendingState)
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
func (it *MockpolygonrollupmanagerOverridePendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerOverridePendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerOverridePendingState represents a OverridePendingState event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerOverridePendingState struct {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterOverridePendingState(opts *bind.FilterOpts, rollupID []uint32) (*MockpolygonrollupmanagerOverridePendingStateIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "OverridePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerOverridePendingStateIterator{contract: _Mockpolygonrollupmanager.contract, event: "OverridePendingState", logs: logs, sub: sub}, nil
}

// WatchOverridePendingState is a free log subscription operation binding the contract event 0x3182bd6e6f74fc1fdc88b60f3a4f4c7f79db6ae6f5b88a1b3f5a1e28ec210d5e.
//
// Solidity: event OverridePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address aggregator)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchOverridePendingState(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerOverridePendingState, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "OverridePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerOverridePendingState)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "OverridePendingState", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseOverridePendingState(log types.Log) (*MockpolygonrollupmanagerOverridePendingState, error) {
	event := new(MockpolygonrollupmanagerOverridePendingState)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "OverridePendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerProveNonDeterministicPendingStateIterator is returned from FilterProveNonDeterministicPendingState and is used to iterate over the raw logs and unpacked data for ProveNonDeterministicPendingState events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerProveNonDeterministicPendingStateIterator struct {
	Event *MockpolygonrollupmanagerProveNonDeterministicPendingState // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerProveNonDeterministicPendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerProveNonDeterministicPendingState)
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
		it.Event = new(MockpolygonrollupmanagerProveNonDeterministicPendingState)
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
func (it *MockpolygonrollupmanagerProveNonDeterministicPendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerProveNonDeterministicPendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerProveNonDeterministicPendingState represents a ProveNonDeterministicPendingState event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerProveNonDeterministicPendingState struct {
	StoredStateRoot [32]byte
	ProvedStateRoot [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterProveNonDeterministicPendingState is a free log retrieval operation binding the contract event 0x1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010.
//
// Solidity: event ProveNonDeterministicPendingState(bytes32 storedStateRoot, bytes32 provedStateRoot)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterProveNonDeterministicPendingState(opts *bind.FilterOpts) (*MockpolygonrollupmanagerProveNonDeterministicPendingStateIterator, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "ProveNonDeterministicPendingState")
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerProveNonDeterministicPendingStateIterator{contract: _Mockpolygonrollupmanager.contract, event: "ProveNonDeterministicPendingState", logs: logs, sub: sub}, nil
}

// WatchProveNonDeterministicPendingState is a free log subscription operation binding the contract event 0x1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010.
//
// Solidity: event ProveNonDeterministicPendingState(bytes32 storedStateRoot, bytes32 provedStateRoot)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchProveNonDeterministicPendingState(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerProveNonDeterministicPendingState) (event.Subscription, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "ProveNonDeterministicPendingState")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerProveNonDeterministicPendingState)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "ProveNonDeterministicPendingState", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseProveNonDeterministicPendingState(log types.Log) (*MockpolygonrollupmanagerProveNonDeterministicPendingState, error) {
	event := new(MockpolygonrollupmanagerProveNonDeterministicPendingState)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "ProveNonDeterministicPendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerRoleAdminChangedIterator struct {
	Event *MockpolygonrollupmanagerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerRoleAdminChanged)
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
		it.Event = new(MockpolygonrollupmanagerRoleAdminChanged)
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
func (it *MockpolygonrollupmanagerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerRoleAdminChanged represents a RoleAdminChanged event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*MockpolygonrollupmanagerRoleAdminChangedIterator, error) {

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

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerRoleAdminChangedIterator{contract: _Mockpolygonrollupmanager.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerRoleAdminChanged)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseRoleAdminChanged(log types.Log) (*MockpolygonrollupmanagerRoleAdminChanged, error) {
	event := new(MockpolygonrollupmanagerRoleAdminChanged)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerRoleGrantedIterator struct {
	Event *MockpolygonrollupmanagerRoleGranted // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerRoleGranted)
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
		it.Event = new(MockpolygonrollupmanagerRoleGranted)
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
func (it *MockpolygonrollupmanagerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerRoleGranted represents a RoleGranted event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*MockpolygonrollupmanagerRoleGrantedIterator, error) {

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

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerRoleGrantedIterator{contract: _Mockpolygonrollupmanager.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerRoleGranted)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseRoleGranted(log types.Log) (*MockpolygonrollupmanagerRoleGranted, error) {
	event := new(MockpolygonrollupmanagerRoleGranted)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerRoleRevokedIterator struct {
	Event *MockpolygonrollupmanagerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerRoleRevoked)
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
		it.Event = new(MockpolygonrollupmanagerRoleRevoked)
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
func (it *MockpolygonrollupmanagerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerRoleRevoked represents a RoleRevoked event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*MockpolygonrollupmanagerRoleRevokedIterator, error) {

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

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerRoleRevokedIterator{contract: _Mockpolygonrollupmanager.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerRoleRevoked)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseRoleRevoked(log types.Log) (*MockpolygonrollupmanagerRoleRevoked, error) {
	event := new(MockpolygonrollupmanagerRoleRevoked)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerSetBatchFeeIterator is returned from FilterSetBatchFee and is used to iterate over the raw logs and unpacked data for SetBatchFee events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerSetBatchFeeIterator struct {
	Event *MockpolygonrollupmanagerSetBatchFee // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerSetBatchFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerSetBatchFee)
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
		it.Event = new(MockpolygonrollupmanagerSetBatchFee)
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
func (it *MockpolygonrollupmanagerSetBatchFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerSetBatchFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerSetBatchFee represents a SetBatchFee event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerSetBatchFee struct {
	NewBatchFee *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSetBatchFee is a free log retrieval operation binding the contract event 0xfb383653f53ee079978d0c9aff7aeff04a10166ce244cca9c9f9d8d96bed45b2.
//
// Solidity: event SetBatchFee(uint256 newBatchFee)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterSetBatchFee(opts *bind.FilterOpts) (*MockpolygonrollupmanagerSetBatchFeeIterator, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "SetBatchFee")
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerSetBatchFeeIterator{contract: _Mockpolygonrollupmanager.contract, event: "SetBatchFee", logs: logs, sub: sub}, nil
}

// WatchSetBatchFee is a free log subscription operation binding the contract event 0xfb383653f53ee079978d0c9aff7aeff04a10166ce244cca9c9f9d8d96bed45b2.
//
// Solidity: event SetBatchFee(uint256 newBatchFee)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchSetBatchFee(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerSetBatchFee) (event.Subscription, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "SetBatchFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerSetBatchFee)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "SetBatchFee", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseSetBatchFee(log types.Log) (*MockpolygonrollupmanagerSetBatchFee, error) {
	event := new(MockpolygonrollupmanagerSetBatchFee)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "SetBatchFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerSetMultiplierBatchFeeIterator is returned from FilterSetMultiplierBatchFee and is used to iterate over the raw logs and unpacked data for SetMultiplierBatchFee events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerSetMultiplierBatchFeeIterator struct {
	Event *MockpolygonrollupmanagerSetMultiplierBatchFee // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerSetMultiplierBatchFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerSetMultiplierBatchFee)
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
		it.Event = new(MockpolygonrollupmanagerSetMultiplierBatchFee)
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
func (it *MockpolygonrollupmanagerSetMultiplierBatchFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerSetMultiplierBatchFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerSetMultiplierBatchFee represents a SetMultiplierBatchFee event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerSetMultiplierBatchFee struct {
	NewMultiplierBatchFee uint16
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterSetMultiplierBatchFee is a free log retrieval operation binding the contract event 0x7019933d795eba185c180209e8ae8bffbaa25bcef293364687702c31f4d302c5.
//
// Solidity: event SetMultiplierBatchFee(uint16 newMultiplierBatchFee)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterSetMultiplierBatchFee(opts *bind.FilterOpts) (*MockpolygonrollupmanagerSetMultiplierBatchFeeIterator, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "SetMultiplierBatchFee")
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerSetMultiplierBatchFeeIterator{contract: _Mockpolygonrollupmanager.contract, event: "SetMultiplierBatchFee", logs: logs, sub: sub}, nil
}

// WatchSetMultiplierBatchFee is a free log subscription operation binding the contract event 0x7019933d795eba185c180209e8ae8bffbaa25bcef293364687702c31f4d302c5.
//
// Solidity: event SetMultiplierBatchFee(uint16 newMultiplierBatchFee)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchSetMultiplierBatchFee(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerSetMultiplierBatchFee) (event.Subscription, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "SetMultiplierBatchFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerSetMultiplierBatchFee)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "SetMultiplierBatchFee", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseSetMultiplierBatchFee(log types.Log) (*MockpolygonrollupmanagerSetMultiplierBatchFee, error) {
	event := new(MockpolygonrollupmanagerSetMultiplierBatchFee)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "SetMultiplierBatchFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerSetPendingStateTimeoutIterator is returned from FilterSetPendingStateTimeout and is used to iterate over the raw logs and unpacked data for SetPendingStateTimeout events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerSetPendingStateTimeoutIterator struct {
	Event *MockpolygonrollupmanagerSetPendingStateTimeout // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerSetPendingStateTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerSetPendingStateTimeout)
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
		it.Event = new(MockpolygonrollupmanagerSetPendingStateTimeout)
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
func (it *MockpolygonrollupmanagerSetPendingStateTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerSetPendingStateTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerSetPendingStateTimeout represents a SetPendingStateTimeout event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerSetPendingStateTimeout struct {
	NewPendingStateTimeout uint64
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterSetPendingStateTimeout is a free log retrieval operation binding the contract event 0xc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c75.
//
// Solidity: event SetPendingStateTimeout(uint64 newPendingStateTimeout)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterSetPendingStateTimeout(opts *bind.FilterOpts) (*MockpolygonrollupmanagerSetPendingStateTimeoutIterator, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "SetPendingStateTimeout")
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerSetPendingStateTimeoutIterator{contract: _Mockpolygonrollupmanager.contract, event: "SetPendingStateTimeout", logs: logs, sub: sub}, nil
}

// WatchSetPendingStateTimeout is a free log subscription operation binding the contract event 0xc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c75.
//
// Solidity: event SetPendingStateTimeout(uint64 newPendingStateTimeout)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchSetPendingStateTimeout(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerSetPendingStateTimeout) (event.Subscription, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "SetPendingStateTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerSetPendingStateTimeout)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "SetPendingStateTimeout", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseSetPendingStateTimeout(log types.Log) (*MockpolygonrollupmanagerSetPendingStateTimeout, error) {
	event := new(MockpolygonrollupmanagerSetPendingStateTimeout)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "SetPendingStateTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerSetTrustedAggregatorIterator is returned from FilterSetTrustedAggregator and is used to iterate over the raw logs and unpacked data for SetTrustedAggregator events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerSetTrustedAggregatorIterator struct {
	Event *MockpolygonrollupmanagerSetTrustedAggregator // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerSetTrustedAggregatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerSetTrustedAggregator)
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
		it.Event = new(MockpolygonrollupmanagerSetTrustedAggregator)
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
func (it *MockpolygonrollupmanagerSetTrustedAggregatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerSetTrustedAggregatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerSetTrustedAggregator represents a SetTrustedAggregator event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerSetTrustedAggregator struct {
	NewTrustedAggregator common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedAggregator is a free log retrieval operation binding the contract event 0x61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca.
//
// Solidity: event SetTrustedAggregator(address newTrustedAggregator)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterSetTrustedAggregator(opts *bind.FilterOpts) (*MockpolygonrollupmanagerSetTrustedAggregatorIterator, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "SetTrustedAggregator")
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerSetTrustedAggregatorIterator{contract: _Mockpolygonrollupmanager.contract, event: "SetTrustedAggregator", logs: logs, sub: sub}, nil
}

// WatchSetTrustedAggregator is a free log subscription operation binding the contract event 0x61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca.
//
// Solidity: event SetTrustedAggregator(address newTrustedAggregator)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchSetTrustedAggregator(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerSetTrustedAggregator) (event.Subscription, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "SetTrustedAggregator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerSetTrustedAggregator)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregator", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseSetTrustedAggregator(log types.Log) (*MockpolygonrollupmanagerSetTrustedAggregator, error) {
	event := new(MockpolygonrollupmanagerSetTrustedAggregator)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerSetTrustedAggregatorTimeoutIterator is returned from FilterSetTrustedAggregatorTimeout and is used to iterate over the raw logs and unpacked data for SetTrustedAggregatorTimeout events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerSetTrustedAggregatorTimeoutIterator struct {
	Event *MockpolygonrollupmanagerSetTrustedAggregatorTimeout // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerSetTrustedAggregatorTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerSetTrustedAggregatorTimeout)
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
		it.Event = new(MockpolygonrollupmanagerSetTrustedAggregatorTimeout)
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
func (it *MockpolygonrollupmanagerSetTrustedAggregatorTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerSetTrustedAggregatorTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerSetTrustedAggregatorTimeout represents a SetTrustedAggregatorTimeout event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerSetTrustedAggregatorTimeout struct {
	NewTrustedAggregatorTimeout uint64
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedAggregatorTimeout is a free log retrieval operation binding the contract event 0x1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1.
//
// Solidity: event SetTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterSetTrustedAggregatorTimeout(opts *bind.FilterOpts) (*MockpolygonrollupmanagerSetTrustedAggregatorTimeoutIterator, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "SetTrustedAggregatorTimeout")
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerSetTrustedAggregatorTimeoutIterator{contract: _Mockpolygonrollupmanager.contract, event: "SetTrustedAggregatorTimeout", logs: logs, sub: sub}, nil
}

// WatchSetTrustedAggregatorTimeout is a free log subscription operation binding the contract event 0x1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1.
//
// Solidity: event SetTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchSetTrustedAggregatorTimeout(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerSetTrustedAggregatorTimeout) (event.Subscription, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "SetTrustedAggregatorTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerSetTrustedAggregatorTimeout)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregatorTimeout", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseSetTrustedAggregatorTimeout(log types.Log) (*MockpolygonrollupmanagerSetTrustedAggregatorTimeout, error) {
	event := new(MockpolygonrollupmanagerSetTrustedAggregatorTimeout)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregatorTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerSetVerifyBatchTimeTargetIterator is returned from FilterSetVerifyBatchTimeTarget and is used to iterate over the raw logs and unpacked data for SetVerifyBatchTimeTarget events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerSetVerifyBatchTimeTargetIterator struct {
	Event *MockpolygonrollupmanagerSetVerifyBatchTimeTarget // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerSetVerifyBatchTimeTargetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerSetVerifyBatchTimeTarget)
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
		it.Event = new(MockpolygonrollupmanagerSetVerifyBatchTimeTarget)
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
func (it *MockpolygonrollupmanagerSetVerifyBatchTimeTargetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerSetVerifyBatchTimeTargetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerSetVerifyBatchTimeTarget represents a SetVerifyBatchTimeTarget event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerSetVerifyBatchTimeTarget struct {
	NewVerifyBatchTimeTarget uint64
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterSetVerifyBatchTimeTarget is a free log retrieval operation binding the contract event 0x1b023231a1ab6b5d93992f168fb44498e1a7e64cef58daff6f1c216de6a68c28.
//
// Solidity: event SetVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterSetVerifyBatchTimeTarget(opts *bind.FilterOpts) (*MockpolygonrollupmanagerSetVerifyBatchTimeTargetIterator, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "SetVerifyBatchTimeTarget")
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerSetVerifyBatchTimeTargetIterator{contract: _Mockpolygonrollupmanager.contract, event: "SetVerifyBatchTimeTarget", logs: logs, sub: sub}, nil
}

// WatchSetVerifyBatchTimeTarget is a free log subscription operation binding the contract event 0x1b023231a1ab6b5d93992f168fb44498e1a7e64cef58daff6f1c216de6a68c28.
//
// Solidity: event SetVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchSetVerifyBatchTimeTarget(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerSetVerifyBatchTimeTarget) (event.Subscription, error) {

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "SetVerifyBatchTimeTarget")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerSetVerifyBatchTimeTarget)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "SetVerifyBatchTimeTarget", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseSetVerifyBatchTimeTarget(log types.Log) (*MockpolygonrollupmanagerSetVerifyBatchTimeTarget, error) {
	event := new(MockpolygonrollupmanagerSetVerifyBatchTimeTarget)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "SetVerifyBatchTimeTarget", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerUpdateRollupIterator is returned from FilterUpdateRollup and is used to iterate over the raw logs and unpacked data for UpdateRollup events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerUpdateRollupIterator struct {
	Event *MockpolygonrollupmanagerUpdateRollup // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerUpdateRollupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerUpdateRollup)
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
		it.Event = new(MockpolygonrollupmanagerUpdateRollup)
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
func (it *MockpolygonrollupmanagerUpdateRollupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerUpdateRollupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerUpdateRollup represents a UpdateRollup event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerUpdateRollup struct {
	RollupID                       uint32
	NewRollupTypeID                uint32
	LastVerifiedBatchBeforeUpgrade uint64
	Raw                            types.Log // Blockchain specific contextual infos
}

// FilterUpdateRollup is a free log retrieval operation binding the contract event 0xf585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d.
//
// Solidity: event UpdateRollup(uint32 indexed rollupID, uint32 newRollupTypeID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterUpdateRollup(opts *bind.FilterOpts, rollupID []uint32) (*MockpolygonrollupmanagerUpdateRollupIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "UpdateRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerUpdateRollupIterator{contract: _Mockpolygonrollupmanager.contract, event: "UpdateRollup", logs: logs, sub: sub}, nil
}

// WatchUpdateRollup is a free log subscription operation binding the contract event 0xf585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d.
//
// Solidity: event UpdateRollup(uint32 indexed rollupID, uint32 newRollupTypeID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchUpdateRollup(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerUpdateRollup, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "UpdateRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerUpdateRollup)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "UpdateRollup", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseUpdateRollup(log types.Log) (*MockpolygonrollupmanagerUpdateRollup, error) {
	event := new(MockpolygonrollupmanagerUpdateRollup)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "UpdateRollup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerVerifyBatchesIterator is returned from FilterVerifyBatches and is used to iterate over the raw logs and unpacked data for VerifyBatches events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerVerifyBatchesIterator struct {
	Event *MockpolygonrollupmanagerVerifyBatches // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerVerifyBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerVerifyBatches)
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
		it.Event = new(MockpolygonrollupmanagerVerifyBatches)
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
func (it *MockpolygonrollupmanagerVerifyBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerVerifyBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerVerifyBatches represents a VerifyBatches event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerVerifyBatches struct {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterVerifyBatches(opts *bind.FilterOpts, rollupID []uint32, aggregator []common.Address) (*MockpolygonrollupmanagerVerifyBatchesIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "VerifyBatches", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerVerifyBatchesIterator{contract: _Mockpolygonrollupmanager.contract, event: "VerifyBatches", logs: logs, sub: sub}, nil
}

// WatchVerifyBatches is a free log subscription operation binding the contract event 0xaac1e7a157b259544ebacd6e8a82ae5d6c8f174e12aa48696277bcc9a661f0b4.
//
// Solidity: event VerifyBatches(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchVerifyBatches(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerVerifyBatches, rollupID []uint32, aggregator []common.Address) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "VerifyBatches", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerVerifyBatches)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "VerifyBatches", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseVerifyBatches(log types.Log) (*MockpolygonrollupmanagerVerifyBatches, error) {
	event := new(MockpolygonrollupmanagerVerifyBatches)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "VerifyBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockpolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator is returned from FilterVerifyBatchesTrustedAggregator and is used to iterate over the raw logs and unpacked data for VerifyBatchesTrustedAggregator events raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator struct {
	Event *MockpolygonrollupmanagerVerifyBatchesTrustedAggregator // Event containing the contract specifics and raw log

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
func (it *MockpolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockpolygonrollupmanagerVerifyBatchesTrustedAggregator)
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
		it.Event = new(MockpolygonrollupmanagerVerifyBatchesTrustedAggregator)
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
func (it *MockpolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockpolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockpolygonrollupmanagerVerifyBatchesTrustedAggregator represents a VerifyBatchesTrustedAggregator event raised by the Mockpolygonrollupmanager contract.
type MockpolygonrollupmanagerVerifyBatchesTrustedAggregator struct {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) FilterVerifyBatchesTrustedAggregator(opts *bind.FilterOpts, rollupID []uint32, aggregator []common.Address) (*MockpolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.FilterLogs(opts, "VerifyBatchesTrustedAggregator", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &MockpolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator{contract: _Mockpolygonrollupmanager.contract, event: "VerifyBatchesTrustedAggregator", logs: logs, sub: sub}, nil
}

// WatchVerifyBatchesTrustedAggregator is a free log subscription operation binding the contract event 0xd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d3.
//
// Solidity: event VerifyBatchesTrustedAggregator(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) WatchVerifyBatchesTrustedAggregator(opts *bind.WatchOpts, sink chan<- *MockpolygonrollupmanagerVerifyBatchesTrustedAggregator, rollupID []uint32, aggregator []common.Address) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Mockpolygonrollupmanager.contract.WatchLogs(opts, "VerifyBatchesTrustedAggregator", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockpolygonrollupmanagerVerifyBatchesTrustedAggregator)
				if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "VerifyBatchesTrustedAggregator", log); err != nil {
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
func (_Mockpolygonrollupmanager *MockpolygonrollupmanagerFilterer) ParseVerifyBatchesTrustedAggregator(log types.Log) (*MockpolygonrollupmanagerVerifyBatchesTrustedAggregator, error) {
	event := new(MockpolygonrollupmanagerVerifyBatchesTrustedAggregator)
	if err := _Mockpolygonrollupmanager.contract.UnpackLog(event, "VerifyBatchesTrustedAggregator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
