package synchronizer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var forkID uint64 = 5

func TestGetBatchL2DataWithoutCommittee(t *testing.T) {
	m := mocks{
		State:       newStateMock(t),
		ZKEVMClient: newZkEVMClientMock(t),
		Etherman:    newEthermanMock(t),
	}

	ctx := context.Background()

	trustedSync := ClientSynchronizer{
		isTrustedSequencer:      true,
		state:                   m.State,
		zkEVMClient:             m.ZKEVMClient,
		etherMan:                m.Etherman,
		ctx:                     ctx,
		selectedCommitteeMember: -1,
	}

	permissionlessSync := ClientSynchronizer{
		isTrustedSequencer:      false,
		state:                   m.State,
		zkEVMClient:             m.ZKEVMClient,
		etherMan:                m.Etherman,
		ctx:                     ctx,
		selectedCommitteeMember: -1,
	}

	const batchNum uint64 = 5
	batchNumBig := big.NewInt(int64(batchNum))
	dataFromDB := []byte("i poli tis Kerkyras einai omorfi")
	errorHash := state.ZeroHash
	unexpectedErrFromDB := errors.New("unexpected db")
	errFromDBTemplate := "failed to get batch data from state for batch num " + strconv.Itoa(int(batchNum)) + ": %s"

	trustedResponse := &types.Batch{Transactions: []types.TransactionOrHash{
		{Tx: &types.Transaction{Nonce: 4}},
		{Tx: &types.Transaction{Nonce: 284}},
	}}
	txs := []ethTypes.Transaction{}
	effectivePercentages := []uint8{}
	for _, transaction := range trustedResponse.Transactions {
		tx := transaction.Tx.CoreTx()
		txs = append(txs, *tx)
		effectivePercentages = append(effectivePercentages, state.MaxEffectivePercentage)
	}
	dataFromTrusted, err := state.EncodeTransactions(txs, effectivePercentages, forkID)
	require.NoError(t, err)
	trustedResponse.BatchL2Data = dataFromTrusted
	trustedResponseEmpty := &types.Batch{Transactions: []types.TransactionOrHash{}}
	txs = []ethTypes.Transaction{}
	dataFromTrustedEmpty, err := state.EncodeTransactions(txs, effectivePercentages, forkID)
	require.NoError(t, err)

	unexpectedErrFromTrusted := errors.New("unexpected trusted")

	type testCase struct {
		Name           string
		ExpectedResult []byte
		ExpectedError  error
		Sync           *ClientSynchronizer
		SetupMocks     func(m *mocks)
	}

	testCases := []testCase{
		// Trusted sync cases
		{
			Name:           "Trusted sync fail if unexpected error from DB",
			ExpectedResult: nil,
			ExpectedError:  fmt.Errorf(errFromDBTemplate, unexpectedErrFromDB),
			Sync:           &trustedSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(nil, unexpectedErrFromDB).
					Once()
			},
		},
		{
			Name:           "Trusted sync fail if data not found on DB",
			ExpectedResult: nil,
			ExpectedError:  fmt.Errorf("data not found on the local DB nor on any data committee member"),
			Sync:           &trustedSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(nil, state.ErrNotFound).
					Once()
				m.Etherman.
					On("GetCurrentDataCommittee").
					Return(nil, nil).
					Once()
			},
		},
		{
			Name:           "Trusted sync fail if hash missmatch on DB",
			ExpectedResult: nil,
			ExpectedError:  fmt.Errorf("data not found on the local DB nor on any data committee member"),
			Sync:           &trustedSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(dataFromDB, nil).
					Once()
				m.Etherman.
					On("GetCurrentDataCommittee").
					Return(nil, nil).
					Once()
			},
		},
		{
			Name:           "Trusted sync succeeds if hash match on DB",
			ExpectedResult: dataFromDB,
			ExpectedError:  nil,
			Sync:           &trustedSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(dataFromDB, nil).
					Once()
			},
		},
		// Permissionless sync  cases
		{
			Name:           "Permissionless sync succeeds if hash match on DB",
			ExpectedResult: dataFromDB,
			ExpectedError:  nil,
			Sync:           &permissionlessSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(dataFromDB, nil).
					Once()
			},
		},
		{
			Name:           "Permissionless sync fail if unexpected error from DB",
			ExpectedResult: nil,
			ExpectedError:  fmt.Errorf(errFromDBTemplate, unexpectedErrFromDB),
			Sync:           &permissionlessSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(nil, unexpectedErrFromDB).
					Once()
			},
		},
		{
			Name:           "Permissionless sync fail if hash missmatch on the DB and error from trusted",
			ExpectedResult: nil,
			ExpectedError:  fmt.Errorf("data not found on the local DB, nor from the trusted sequencer nor on any data committee member"),
			Sync:           &permissionlessSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(dataFromDB, nil).
					Once()
				m.Etherman.
					On("GetCurrentDataCommittee").
					Return(nil, nil).
					Once()
				m.ZKEVMClient.
					On("BatchByNumber", ctx, batchNumBig).
					Return(nil, unexpectedErrFromTrusted).
					Once()
			},
		},
		{
			Name:           "Permissionless sync fail if hash missmatch on the DB and from trusted sequencer",
			ExpectedResult: nil,
			ExpectedError:  fmt.Errorf("data not found on the local DB, nor from the trusted sequencer nor on any data committee member"),
			Sync:           &permissionlessSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(dataFromDB, nil).
					Once()
				m.Etherman.
					On("GetCurrentDataCommittee").
					Return(nil, nil).
					Once()
				m.ZKEVMClient.
					On("BatchByNumber", ctx, batchNumBig).
					Return(trustedResponse, nil).
					Once()
			},
		},
		{
			Name:           "Permissionless sync succeeds if hash missmatch on the DB and match from trusted",
			ExpectedResult: dataFromTrusted,
			ExpectedError:  nil,
			Sync:           &permissionlessSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(dataFromDB, nil).
					Once()
				m.ZKEVMClient.
					On("BatchByNumber", ctx, batchNumBig).
					Return(trustedResponse, nil).
					Once()
			},
		},
		{
			Name:           "Permissionless sync fail if not found on the DB and error from trusted",
			ExpectedResult: nil,
			ExpectedError:  fmt.Errorf("data not found on the local DB, nor from the trusted sequencer nor on any data committee member"),
			Sync:           &permissionlessSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(nil, state.ErrNotFound).
					Once()
				m.Etherman.
					On("GetCurrentDataCommittee").
					Return(nil, nil).
					Once()
				m.ZKEVMClient.
					On("BatchByNumber", ctx, batchNumBig).
					Return(nil, unexpectedErrFromTrusted).
					Once()
			},
		},
		{
			Name:           "Permissionless sync fail fail if not found on the DB and hash missmatch trusted",
			ExpectedResult: nil,
			ExpectedError:  fmt.Errorf("data not found on the local DB, nor from the trusted sequencer nor on any data committee member"),
			Sync:           &permissionlessSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(nil, state.ErrNotFound).
					Once()
				m.Etherman.
					On("GetCurrentDataCommittee").
					Return(nil, nil).
					Once()
				m.ZKEVMClient.
					On("BatchByNumber", ctx, batchNumBig).
					Return(trustedResponse, nil).
					Once()
			},
		},
		{
			Name:           "Permissionless sync succeeds if not found on the DB and match from trusted",
			ExpectedResult: dataFromTrusted,
			ExpectedError:  nil,
			Sync:           &permissionlessSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(nil, state.ErrNotFound).
					Once()
				m.ZKEVMClient.
					On("BatchByNumber", ctx, batchNumBig).
					Return(trustedResponse, nil).
					Once()
			},
		},
		{
			Name:           "Permissionless sync succeeds if not found on the DB and match from trusted empty response",
			ExpectedResult: dataFromTrustedEmpty,
			ExpectedError:  nil,
			Sync:           &permissionlessSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(nil, state.ErrNotFound).
					Once()
				m.ZKEVMClient.
					On("BatchByNumber", ctx, batchNumBig).
					Return(trustedResponseEmpty, nil).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(&m)

			var expectedHash common.Hash
			if tc.ExpectedError == nil {
				expectedHash = crypto.Keccak256Hash(tc.ExpectedResult)
			} else {
				expectedHash = errorHash
			}

			res, err := tc.Sync.getBatchL2Data(batchNum, expectedHash)
			assert.Equal(t, tc.ExpectedResult, res)
			if tc.ExpectedError != nil {
				require.NotNil(t, err)
				assert.Equal(t, tc.ExpectedError.Error(), err.Error())
			}
		})
	}
}

func TestGetBatchL2DataWithCommittee(t *testing.T) {
	m := mocks{
		State:                      newStateMock(t),
		ZKEVMClient:                newZkEVMClientMock(t),
		Etherman:                   newEthermanMock(t),
		DataCommitteeClientFactory: newDataCommitteeClientFactoryMock(t),
	}

	ctx := context.Background()

	committeeMembers := []etherman.DataCommitteeMember{
		{
			URL:  "0",
			Addr: common.HexToAddress("0x0"),
		},
		{
			URL:  "1",
			Addr: common.HexToAddress("0x1"),
		},
		{
			URL:  "2",
			Addr: common.HexToAddress("0x2"),
		},
	}
	trustedSync := ClientSynchronizer{
		isTrustedSequencer:         true,
		state:                      m.State,
		zkEVMClient:                m.ZKEVMClient,
		etherMan:                   m.Etherman,
		ctx:                        ctx,
		selectedCommitteeMember:    0,
		committeeMembers:           committeeMembers,
		dataCommitteeClientFactory: m.DataCommitteeClientFactory,
	}

	permissionlessSync := ClientSynchronizer{
		isTrustedSequencer:         false,
		state:                      m.State,
		zkEVMClient:                m.ZKEVMClient,
		etherMan:                   m.Etherman,
		ctx:                        ctx,
		selectedCommitteeMember:    1,
		committeeMembers:           committeeMembers,
		dataCommitteeClientFactory: m.DataCommitteeClientFactory,
	}

	const batchNum uint64 = 5
	batchNumBig := big.NewInt(int64(batchNum))
	dataFromDB := []byte("i poli tis Kerkyras einai omorfi")
	errorHash := state.ZeroHash

	trustedResponse := &types.Batch{Transactions: []types.TransactionOrHash{
		{Tx: &types.Transaction{Nonce: 4}},
		{Tx: &types.Transaction{Nonce: 284}},
	}}
	txs := []ethTypes.Transaction{}
	effectivePercentages := []uint8{}
	for _, transaction := range trustedResponse.Transactions {
		tx := transaction.Tx.CoreTx()
		txs = append(txs, *tx)
		effectivePercentages = append(effectivePercentages, state.MaxEffectivePercentage)
	}
	dataFromTrusted, err := state.EncodeTransactions(txs, effectivePercentages, forkID)
	require.NoError(t, err)
	trustedResponse.BatchL2Data = dataFromTrusted

	type testCase struct {
		Name           string
		ExpectedResult []byte
		ExpectedError  error
		Sync           *ClientSynchronizer
		SetupMocks     func(m *mocks)
		Retry          bool
	}

	testCases := []testCase{
		// Trusted sync cases
		{
			Name:           "Trusted sync fail if all the members don't answer",
			ExpectedResult: nil,
			ExpectedError:  fmt.Errorf("data not found on the local DB nor on any data committee member"),
			Sync:           &trustedSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(nil, state.ErrNotFound).
					Once()
				DAClientMock := newDataCommitteeClientMock(t)
				m.DataCommitteeClientFactory.
					On("New", "0").
					Return(DAClientMock).
					Once()
				DAClientMock.
					On("GetOffChainData", trustedSync.ctx, state.ZeroHash).
					Return([]byte("not the correct data"), nil).
					Once()
				m.DataCommitteeClientFactory.
					On("New", "1").
					Return(DAClientMock).
					Once()
				DAClientMock.
					On("GetOffChainData", trustedSync.ctx, state.ZeroHash).
					Return(nil, errors.New("not today")).
					Once()
				m.DataCommitteeClientFactory.
					On("New", "2").
					Return(DAClientMock).
					Once()
				DAClientMock.
					On("GetOffChainData", trustedSync.ctx, state.ZeroHash).
					Return([]byte("not the correct data"), nil).
					Once()
				m.Etherman.
					On("GetCurrentDataCommittee").
					Return(nil, nil).
					Once()
			},
		},
		{
			Name:           "Trusted sync succeeds after 2nd committee member answers correctly",
			ExpectedResult: dataFromDB,
			ExpectedError:  nil,
			Sync:           &trustedSync,
			SetupMocks: func(m *mocks) {
				// Reset DAC
				trustedSync.committeeMembers = committeeMembers
				trustedSync.selectedCommitteeMember = 0
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(nil, state.ErrNotFound).
					Once()
				DAClientMock := newDataCommitteeClientMock(t)
				m.DataCommitteeClientFactory.
					On("New", "0").
					Return(DAClientMock).
					Once()
				DAClientMock.
					On("GetOffChainData", trustedSync.ctx, crypto.Keccak256Hash(dataFromDB)).
					Return([]byte("not the correct data"), nil).
					Once()
				m.DataCommitteeClientFactory.
					On("New", "1").
					Return(DAClientMock).
					Once()
				DAClientMock.
					On("GetOffChainData", trustedSync.ctx, crypto.Keccak256Hash(dataFromDB)).
					Return(nil, errors.New("not today")).
					Once()
				m.DataCommitteeClientFactory.
					On("New", "2").
					Return(DAClientMock).
					Once()
				DAClientMock.
					On("GetOffChainData", trustedSync.ctx, crypto.Keccak256Hash(dataFromDB)).
					Return(dataFromDB, nil).
					Once()
			},
		},
		// Permissionless sync  cases
		{
			Name:           "Permissionless sync succeeds after 2nd committee member answers correctly",
			ExpectedResult: dataFromDB,
			ExpectedError:  nil,
			Sync:           &permissionlessSync,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(nil, state.ErrNotFound).
					Once()
				m.ZKEVMClient.
					On("BatchByNumber", ctx, batchNumBig).
					Return(nil, errors.New("not today")).
					Once()
				DAClientMock := newDataCommitteeClientMock(t)
				m.DataCommitteeClientFactory.
					On("New", "1").
					Return(DAClientMock).
					Once()
				DAClientMock.
					On("GetOffChainData", trustedSync.ctx, crypto.Keccak256Hash(dataFromDB)).
					Return([]byte("not the correct data"), nil).
					Once()
				m.DataCommitteeClientFactory.
					On("New", "2").
					Return(DAClientMock).
					Once()
				DAClientMock.
					On("GetOffChainData", trustedSync.ctx, crypto.Keccak256Hash(dataFromDB)).
					Return(nil, errors.New("not today")).
					Once()
				m.DataCommitteeClientFactory.
					On("New", "0").
					Return(DAClientMock).
					Once()
				DAClientMock.
					On("GetOffChainData", trustedSync.ctx, crypto.Keccak256Hash(dataFromDB)).
					Return(dataFromDB, nil).
					Once()
			},
		},
		{
			Name:           "Permissionless sync succeeds after updating DAC",
			ExpectedResult: dataFromDB,
			ExpectedError:  nil,
			Sync:           &permissionlessSync,
			Retry:          true,
			SetupMocks: func(m *mocks) {
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(nil, state.ErrNotFound).
					Once()
				m.ZKEVMClient.
					On("BatchByNumber", ctx, batchNumBig).
					Return(nil, errors.New("not today")).
					Once()
				DAClientMock := newDataCommitteeClientMock(t)
				m.DataCommitteeClientFactory.
					On("New", "1").
					Return(DAClientMock).
					Once()
				DAClientMock.
					On("GetOffChainData", trustedSync.ctx, crypto.Keccak256Hash(dataFromDB)).
					Return([]byte("not the correct data"), nil).
					Once()
				m.DataCommitteeClientFactory.
					On("New", "2").
					Return(DAClientMock).
					Once()
				DAClientMock.
					On("GetOffChainData", trustedSync.ctx, crypto.Keccak256Hash(dataFromDB)).
					Return(nil, errors.New("not today")).
					Once()
				m.DataCommitteeClientFactory.
					On("New", "0").
					Return(DAClientMock).
					Once()
				DAClientMock.
					On("GetOffChainData", trustedSync.ctx, crypto.Keccak256Hash(dataFromDB)).
					Return(nil, errors.New("not today")).
					Once()
				const succesfullURL = "the time is now"
				m.Etherman.
					On("GetCurrentDataCommittee").
					Return(&etherman.DataCommittee{
						Members: []etherman.DataCommitteeMember{{
							URL:  succesfullURL,
							Addr: common.HexToAddress("0xff"),
						}},
					}, nil).
					Once()
				m.State.
					On("GetBatchL2DataByNumber", ctx, batchNum, nil).
					Return(nil, state.ErrNotFound).
					Once()
				m.ZKEVMClient.
					On("BatchByNumber", ctx, batchNumBig).
					Return(nil, errors.New("not today")).
					Once()
				m.DataCommitteeClientFactory.
					On("New", succesfullURL).
					Return(DAClientMock).
					Once()
				DAClientMock.
					On("GetOffChainData", trustedSync.ctx, crypto.Keccak256Hash(dataFromDB)).
					Return(dataFromDB, nil).
					Once()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tc := testCase
			tc.SetupMocks(&m)

			var expectedHash common.Hash
			if tc.ExpectedError == nil {
				expectedHash = crypto.Keccak256Hash(tc.ExpectedResult)
			} else {
				expectedHash = errorHash
			}

			res, err := tc.Sync.getBatchL2Data(batchNum, expectedHash)
			if tc.Retry {
				require.Error(t, err)
				res, err = tc.Sync.getBatchL2Data(batchNum, expectedHash)
			}
			assert.Equal(t, tc.ExpectedResult, res)
			if tc.ExpectedError != nil {
				require.NotNil(t, err)
				assert.Equal(t, tc.ExpectedError.Error(), err.Error())
			}
		})
	}
}
