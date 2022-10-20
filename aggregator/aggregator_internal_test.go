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
	a := Aggregator{State: st}
	ctx := context.Background()

	verifiedBatch := &state.VerifiedBatch{BatchNumber: 1}

	st.On("GetLastVerifiedBatch", ctx, nil).Return(verifiedBatch, nil)
	st.On("GetVirtualBatchByNumber", ctx, verifiedBatch.BatchNumber+1, nil).Return(batchToVerify, nil)
	st.On("GetGeneratedProofByBatchNumber", ctx, verifiedBatch.BatchNumber+1, nil).Return(nil, state.ErrNotFound)

	res, err := a.getBatchToVerify(ctx)
	require.NoError(t, err)
	require.Equal(t, batchToVerify, res)
}

func TestBuildInputProver(t *testing.T) {
	st := new(aggrMocks.StateMock)
	ethTxManager := new(aggrMocks.EthTxManager)
	etherman := new(aggrMocks.Etherman)
	proverClient := new(aggrMocks.ProverClientMock)
	a := Aggregator{
		State:         st,
		EthTxManager:  ethTxManager,
		Ethman:        etherman,
		ProverClients: []proverClientInterface{proverClient},
	}
	var (
		oldStateRoot     = common.HexToHash("0xbdde84a5932a2f0a1a4c6c51f3b64ea265d4f1461749298cfdd09b31122ce0d6")
		newStateRoot     = common.HexToHash("0x613aabebf4fddf2ad0f034a8c73aa2f9c5a6fac3a07543023e0a6ee6f36e5795")
		oldLocalExitRoot = common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1")
		newLocalExitRoot = common.HexToHash("0x40a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0")
		seqAddress       = common.HexToAddress("0x123")
		batchL2Data      = []byte("data")
		previousBatch    = &state.Batch{
			BatchNumber:   1,
			StateRoot:     oldStateRoot,
			LocalExitRoot: oldLocalExitRoot,
		}
	)

	ctx := context.Background()

	verifiedBatch := &state.VerifiedBatch{BatchNumber: 1}

	st.On("GetLastVerifiedBatch", ctx, nil).Return(verifiedBatch, nil)
	etherman.On("GetLatestVerifiedBatchNum").Return(verifiedBatch.BatchNumber, nil)

	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	require.NoError(t, err)

	st.On("GetBatchByNumber", mock.Anything, lastVerifiedBatch.BatchNumber, nil).Return(previousBatch, nil)

	tx := *types.NewTransaction(1, common.HexToAddress("1"), big.NewInt(1), 0, big.NewInt(1), []byte("bbb"))
	batchToVerify := &state.Batch{
		BatchNumber:    2,
		Coinbase:       seqAddress,
		BatchL2Data:    batchL2Data,
		StateRoot:      newStateRoot,
		LocalExitRoot:  newLocalExitRoot,
		Timestamp:      time.Now(),
		Transactions:   []types.Transaction{tx},
		GlobalExitRoot: common.HexToHash("0xc1df82d9c4b87413eae2ef048f94b4d3554cea73d92b0f7af96e0271c691e2bb"),
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

	etherman.On("GetPublicAddress").Return(common.HexToAddress("0x123"))
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
		State:         st,
		EthTxManager:  ethTxManager,
		Ethman:        etherman,
		ProverClients: []proverClientInterface{proverClient},
	}
	var (
		newLocalExitRoot = common.HexToHash("0xbdde84a5932a2f0a1a4c6c51f3b64ea265d4f1461749298cfdd09b31122ce0d6")
		seqAddress       = common.HexToAddress("0x123")
		batchL2Data      = []byte("data")
	)
	ctx := context.Background()

	verifiedBatch := &state.VerifiedBatch{BatchNumber: 1}

	st.On("GetLastVerifiedBatch", ctx, nil).Return(verifiedBatch, nil)
	etherman.On("GetLatestVerifiedBatchNum").Return(verifiedBatch.BatchNumber, nil)

	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	require.NoError(t, err)

	st.On("GetBatchByNumber", mock.Anything, lastVerifiedBatch.BatchNumber, nil).Return(nil, errors.New("error"))

	tx := *types.NewTransaction(1, common.HexToAddress("1"), big.NewInt(1), 0, big.NewInt(1), []byte("bbb"))
	batchToVerify := &state.Batch{
		BatchNumber:    2,
		Coinbase:       seqAddress,
		BatchL2Data:    batchL2Data,
		LocalExitRoot:  newLocalExitRoot,
		Timestamp:      time.Now(),
		Transactions:   []types.Transaction{tx},
		GlobalExitRoot: common.HexToHash("0xc1df82d9c4b87413eae2ef048f94b4d3554cea73d92b0f7af96e0271c691e2bb"),
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
			IntervalToConsolidateState: cfgTypes.NewDuration(1 * time.Second),
			ChainID:                    1000,
		},
		State:                st,
		EthTxManager:         ethTxManager,
		Ethman:               etherman,
		ProverClients:        []proverClientInterface{proverClient},
		ProfitabilityChecker: NewTxProfitabilityCheckerAcceptAll(st, 1*time.Second),
	}
	var (
		oldStateRoot     = common.HexToHash("0xbdde84a5932a2f0a1a4c6c51f3b64ea265d4f1461749298cfdd09b31122ce0d6")
		newStateRoot     = common.HexToHash("0x613aabebf4fddf2ad0f034a8c73aa2f9c5a6fac3a07543023e0a6ee6f36e5795")
		oldLocalExitRoot = common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1")
		newLocalExitRoot = common.HexToHash("0x40a885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9a0")
		seqAddress       = common.HexToAddress("0x123")
		aggrAddress      = common.HexToAddress("0x123")
		verifiedBatch    = &state.VerifiedBatch{BatchNumber: 1}
		tx               = *types.NewTransaction(1, common.HexToAddress("1"), big.NewInt(1), 0, big.NewInt(1), []byte("bbb"))
		previousBatch    = &state.Batch{
			BatchNumber:   1,
			StateRoot:     oldStateRoot,
			LocalExitRoot: oldLocalExitRoot,
		}
		batchToVerify = &state.Batch{
			BatchNumber:    2,
			Coinbase:       seqAddress,
			StateRoot:      newStateRoot,
			LocalExitRoot:  newLocalExitRoot,
			Timestamp:      time.Now(),
			Transactions:   []types.Transaction{tx},
			GlobalExitRoot: common.HexToHash("0xc1df82d9c4b87413eae2ef048f94b4d3554cea73d92b0f7af96e0271c691e2bb"),
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
				AggregatorAddr:   aggrAddress.String(),
				ChainId:          a.cfg.ChainID,
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
	batchHashData := common.BytesToHash(keccak256.Hash(
		batchToVerify.BatchL2Data,
		batchToVerify.GlobalExitRoot[:],
		blockTimestampByte,
		batchToVerify.Coinbase[:],
	))
	expectedInputProver.PublicInputs.BatchHashData = batchHashData.String()
	// isSynced
	proverClient.On("IsIdle", mock.Anything).Return(true)
	proverClient.On("GetURI", mock.Anything).Return("mockProver:MockPort")
	st.On("GetLastVerifiedBatch", mock.Anything, nil).Return(verifiedBatch, nil)
	st.On("GetGeneratedProofByBatchNumber", mock.Anything, verifiedBatch.BatchNumber+1, nil).Return(nil, state.ErrNotFound)
	st.On("AddGeneratedProof", mock.Anything, mock.Anything, nil).Return(nil)
	st.On("UpdateGeneratedProof", mock.Anything, mock.Anything, nil).Return(nil)
	etherman.On("GetLatestVerifiedBatchNum").Return(uint64(1), nil)
	etherman.On("GetPublicAddress").Return(aggrAddress)

	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(context.Background(), nil)
	require.NoError(t, err)
	// get batch to verify
	st.On("GetVirtualBatchByNumber", mock.Anything, lastVerifiedBatch.BatchNumber+1, nil).Return(batchToVerify, nil)
	// build input prover
	st.On("GetBatchByNumber", mock.Anything, lastVerifiedBatch.BatchNumber, nil).Return(previousBatch, nil)
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
}
