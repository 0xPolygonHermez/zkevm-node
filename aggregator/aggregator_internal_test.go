package aggregator

import (
	"context"
	"encoding/binary"
	"errors"
	"math/big"
	"testing"
	"time"

	aggrMocks "github.com/0xPolygonHermez/zkevm-node/aggregator/mocks"
	cfgTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iden3/go-iden3-crypto/keccak256"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestIsSynced(t *testing.T) {
	st := new(aggrMocks.StateMock)
	etherman := new(aggrMocks.Etherman)
	a := Aggregator{State: st, Ethman: etherman}
	ctx := context.Background()
	verifiedBatch := &state.VerifiedBatch{BatchNumber: 1}
	st.On("GetLastVerifiedBatch", ctx, nil).Return(verifiedBatch, nil)
	etherman.On("GetLatestVerifiedBatchNum").Return(uint64(1), nil)
	isSynced := a.isSynced(ctx)
	require.True(t, isSynced)
}

func TestIsSyncedNotSynced(t *testing.T) {
	st := new(aggrMocks.StateMock)
	etherman := new(aggrMocks.Etherman)
	a := Aggregator{State: st, Ethman: etherman}
	ctx := context.Background()
	verifiedBatch := &state.VerifiedBatch{BatchNumber: 1}
	st.On("GetLastVerifiedBatch", ctx, nil).Return(verifiedBatch, nil)
	etherman.On("GetLatestVerifiedBatchNum").Return(uint64(2), nil)
	isSynced := a.isSynced(ctx)
	require.False(t, isSynced)
}

func TestGetBatchToVerify(t *testing.T) {
	st := new(aggrMocks.StateMock)
	batchToVerify := &state.Batch{BatchNumber: 1}
	a := Aggregator{State: st, batchesSent: make(map[uint64]bool)}
	a.batchesSent[a.lastVerifiedBatchNum] = true
	ctx := context.Background()
	st.On("GetBatchByNumber", ctx, a.lastVerifiedBatchNum+1, nil).Return(batchToVerify, nil)
	res, err := a.getBatchToVerify(ctx)
	require.NoError(t, err)
	require.Equal(t, batchToVerify, res)
	require.False(t, a.batchesSent[a.lastVerifiedBatchNum])
}

func TestGetBatchToVerifyBatchAlreadySent(t *testing.T) {
	st := new(aggrMocks.StateMock)
	batchToVerify := &state.Batch{BatchNumber: 1}
	a := Aggregator{State: st, batchesSent: make(map[uint64]bool)}
	a.lastVerifiedBatchNum = 1
	a.batchesSent[a.lastVerifiedBatchNum] = true
	ctx := context.Background()
	st.On("GetBatchByNumber", ctx, a.lastVerifiedBatchNum+1, nil).Return(batchToVerify, nil)
	res, err := a.getBatchToVerify(ctx)
	require.NoError(t, err)
	require.Nil(t, res)
	require.False(t, a.batchesSent[a.lastVerifiedBatchNum])
}

func TestBuildInputProver(t *testing.T) {
	st := new(aggrMocks.StateMock)
	ethTxManager := new(aggrMocks.EthTxManager)
	etherman := new(aggrMocks.Etherman)
	proverClient := new(aggrMocks.ProverClientMock)
	a := Aggregator{
		State:                st,
		EthTxManager:         ethTxManager,
		Ethman:               etherman,
		ProverClient:         proverClient,
		lastVerifiedBatchNum: 1,
	}
	var (
		oldStateRoot     = common.HexToHash("oldroot")
		newStateRoot     = common.HexToHash("newroot")
		oldLocalExitRoot = common.HexToHash("oldleroot")
		newLocalExitRoot = common.HexToHash("newleroot")
		seqAddress       = common.HexToAddress("0x123")
		batchL2Data      = []byte("data")
	)
	st.On("GetStateRootByBatchNumber", mock.Anything, a.lastVerifiedBatchNum, nil).Return(oldStateRoot, nil)
	st.On("GetStateRootByBatchNumber", mock.Anything, a.lastVerifiedBatchNum+1, nil).Return(newStateRoot, nil)
	st.On("GetLocalExitRootByBatchNumber", mock.Anything, a.lastVerifiedBatchNum, nil).Return(oldLocalExitRoot, nil)

	ctx := context.Background()
	tx := *types.NewTransaction(1, common.HexToAddress("1"), big.NewInt(1), 0, big.NewInt(1), []byte("bbb"))
	batchToVerify := &state.Batch{
		BatchNumber:    2,
		Coinbase:       seqAddress,
		BatchL2Data:    batchL2Data,
		LocalExitRoot:  newLocalExitRoot,
		Timestamp:      time.Now(),
		Transactions:   []types.Transaction{tx},
		GlobalExitRoot: common.HexToHash("geroot"),
	}

	expectedInputProver := &pb.InputProver{
		PublicInputs: &pb.PublicInputs{
			OldStateRoot:     oldStateRoot.String(),
			OldLocalExitRoot: oldLocalExitRoot.String(),
			NewStateRoot:     newStateRoot.String(),
			NewLocalExitRoot: newLocalExitRoot.String(),
			SequencerAddr:    seqAddress.String(),
			BatchHashData:    hex.EncodeToString(batchL2Data),
			BatchNum:         uint32(batchToVerify.BatchNumber),
			EthTimestamp:     uint64(batchToVerify.Timestamp.Unix()),
		},
		GlobalExitRoot: batchToVerify.GlobalExitRoot.String(),
		BatchL2Data:    hex.EncodeToString(batchL2Data),
	}
	ip, err := a.buildInputProver(ctx, batchToVerify)
	require.NoError(t, err)
	require.NotNil(t, ip)
	require.Equal(t, expectedInputProver.PublicInputs.BatchNum, ip.PublicInputs.BatchNum)
	require.Equal(t, expectedInputProver.PublicInputs.OldStateRoot, ip.PublicInputs.OldStateRoot)
	require.Equal(t, expectedInputProver.PublicInputs.OldLocalExitRoot, ip.PublicInputs.OldLocalExitRoot)
	require.Equal(t, expectedInputProver.PublicInputs.NewStateRoot, ip.PublicInputs.NewStateRoot)
	require.Equal(t, expectedInputProver.PublicInputs.NewLocalExitRoot, ip.PublicInputs.NewLocalExitRoot)
	require.Equal(t, expectedInputProver.PublicInputs.SequencerAddr, ip.PublicInputs.SequencerAddr)
	require.Equal(t, expectedInputProver.PublicInputs.EthTimestamp, ip.PublicInputs.EthTimestamp)
	require.Equal(t, expectedInputProver.GlobalExitRoot, ip.GlobalExitRoot)
	require.Equal(t, expectedInputProver.BatchL2Data, ip.BatchL2Data)
}

func TestBuildInputProverError(t *testing.T) {
	st := new(aggrMocks.StateMock)
	ethTxManager := new(aggrMocks.EthTxManager)
	etherman := new(aggrMocks.Etherman)
	proverClient := new(aggrMocks.ProverClientMock)
	a := Aggregator{
		State:                st,
		EthTxManager:         ethTxManager,
		Ethman:               etherman,
		ProverClient:         proverClient,
		lastVerifiedBatchNum: 1,
	}
	var (
		newLocalExitRoot = common.HexToHash("newleroot")
		seqAddress       = common.HexToAddress("0x123")
		batchL2Data      = []byte("data")
	)
	st.On("GetStateRootByBatchNumber", mock.Anything, a.lastVerifiedBatchNum, nil).Return(nil, errors.New("error"))

	ctx := context.Background()
	tx := *types.NewTransaction(1, common.HexToAddress("1"), big.NewInt(1), 0, big.NewInt(1), []byte("bbb"))
	batchToVerify := &state.Batch{
		BatchNumber:    2,
		Coinbase:       seqAddress,
		BatchL2Data:    batchL2Data,
		LocalExitRoot:  newLocalExitRoot,
		Timestamp:      time.Now(),
		Transactions:   []types.Transaction{tx},
		GlobalExitRoot: common.HexToHash("geroot"),
	}

	ip, err := a.buildInputProver(ctx, batchToVerify)
	require.Error(t, err)
	require.Nil(t, ip)
}

func TestAggregatorFlow(t *testing.T) {
	st := new(aggrMocks.StateMock)
	ethTxManager := new(aggrMocks.EthTxManager)
	etherman := new(aggrMocks.Etherman)
	proverClient := new(aggrMocks.ProverClientMock)
	a := Aggregator{
		cfg: Config{
			IntervalToConsolidateState:                          cfgTypes.NewDuration(1 * time.Second),
			IntervalFrequencyToGetProofGenerationStateInSeconds: cfgTypes.Duration{},
			TxProfitabilityCheckerType:                          "",
			TxProfitabilityMinReward:                            TokenAmountWithDecimals{},
			IntervalAfterWhichBatchConsolidateAnyway:            cfgTypes.Duration{},
		},
		State:                st,
		EthTxManager:         ethTxManager,
		Ethman:               etherman,
		ProverClient:         proverClient,
		ProfitabilityChecker: NewTxProfitabilityCheckerAcceptAll(st, 1*time.Second),
		lastVerifiedBatchNum: 1,
		batchesSent:          make(map[uint64]bool),
	}
	var (
		oldStateRoot     = common.HexToHash("oldroot")
		newStateRoot     = common.HexToHash("newroot")
		oldLocalExitRoot = common.HexToHash("oldleroot")
		newLocalExitRoot = common.HexToHash("newleroot")
		seqAddress       = common.HexToAddress("0x123")
		verifiedBatch    = &state.VerifiedBatch{BatchNumber: 1}
		tx               = *types.NewTransaction(1, common.HexToAddress("1"), big.NewInt(1), 0, big.NewInt(1), []byte("bbb"))
		batchToVerify    = &state.Batch{
			BatchNumber:    2,
			Coinbase:       seqAddress,
			LocalExitRoot:  newLocalExitRoot,
			Timestamp:      time.Now(),
			Transactions:   []types.Transaction{tx},
			GlobalExitRoot: common.HexToHash("geroot"),
		}
		expectedInputProver = &pb.InputProver{
			PublicInputs: &pb.PublicInputs{
				OldStateRoot:     oldStateRoot.String(),
				OldLocalExitRoot: oldLocalExitRoot.String(),
				NewStateRoot:     newStateRoot.String(),
				NewLocalExitRoot: newLocalExitRoot.String(),
				SequencerAddr:    seqAddress.String(),
				BatchNum:         uint32(batchToVerify.BatchNumber),
				EthTimestamp:     uint64(batchToVerify.Timestamp.Unix()),
			},
			GlobalExitRoot:    batchToVerify.GlobalExitRoot.String(),
			Db:                map[string]string{},
			ContractsBytecode: map[string]string{},
		}
		getProofResponse = &pb.GetProofResponse{
			Id: "1",
			Proof: &pb.Proof{
				ProofA: []string{"1"},
				ProofB: []*pb.ProofB{},
				ProofC: []string{"1"},
			},
			Public: &pb.PublicInputsExtended{
				PublicInputs: expectedInputProver.PublicInputs,
				InputHash:    "0x2df6320fec1fa7dafa408f9c5f2b31603b2148cb02063935d2e4d105121c1967",
			},
			Result:       1,
			ResultString: "1",
		}
	)
	blockTimestampByte := make([]byte, 8) //nolint:gomnd
	binary.BigEndian.PutUint64(blockTimestampByte, uint64(batchToVerify.Timestamp.Unix()))
	rawTxs, err := state.EncodeTransactions(batchToVerify.Transactions)
	require.NoError(t, err)
	batchHashData := common.BytesToHash(keccak256.Hash(
		rawTxs,
		batchToVerify.GlobalExitRoot[:],
		blockTimestampByte,
		batchToVerify.Coinbase[:],
	))
	expectedInputProver.PublicInputs.BatchHashData = batchHashData.String()
	// isSynced
	st.On("GetLastVerifiedBatch", mock.Anything, nil).Return(verifiedBatch, nil)
	etherman.On("GetLatestVerifiedBatchNum").Return(uint64(1), nil)

	// get batch to verify
	st.On("GetBatchByNumber", mock.Anything, a.lastVerifiedBatchNum+1, nil).Return(batchToVerify, nil)
	// build input prover
	st.On("GetStateRootByBatchNumber", mock.Anything, a.lastVerifiedBatchNum, nil).Return(oldStateRoot, nil)
	st.On("GetStateRootByBatchNumber", mock.Anything, a.lastVerifiedBatchNum+1, nil).Return(newStateRoot, nil)
	st.On("GetLocalExitRootByBatchNumber", mock.Anything, a.lastVerifiedBatchNum, nil).Return(oldLocalExitRoot, nil)
	// gen proof id
	proverClient.On("GetGenProofID", mock.Anything, expectedInputProver).Return("1", nil)
	// get proof
	proverClient.On("GetResGetProof", mock.Anything, "1", batchToVerify.BatchNumber).Return(getProofResponse, nil)
	// send proof to the eth
	ethTxManager.On("VerifyBatch", batchToVerify.BatchNumber, getProofResponse).Return(nil)
	ticker := time.NewTicker(a.cfg.IntervalToConsolidateState.Duration)
	defer ticker.Stop()
	ctx := context.Background()
	a.tryVerifyBatch(ctx, ticker)
	require.True(t, a.batchesSent[batchToVerify.BatchNumber])
}
