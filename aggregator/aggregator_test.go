package aggregator

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/mocks"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	configTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mox struct {
	stateMock    *mocks.StateMock
	ethTxManager *mocks.EthTxManager
	etherman     *mocks.Etherman
	proverMock   *mocks.ProverMock
}

func TestSendFinalProof(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	batchNum := uint64(23)
	batchNumFinal := uint64(42)
	currentNonce := uint64(1)
	estimatedGas := uint64(1)
	suggestedGasPrice := big.NewInt(1)
	var to *common.Address
	value := big.NewInt(0)
	var data []byte = nil
	finalBatch := state.Batch{
		LocalExitRoot: common.BytesToHash([]byte("localExitRoot")),
		StateRoot:     common.BytesToHash([]byte("stateRoot")),
	}
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    currentNonce,
		To:       to,
		Value:    value,
		Gas:      estimatedGas,
		GasPrice: suggestedGasPrice,
		Data:     data,
	})
	proofID := "proofId"
	proverID := "proverID"
	recursiveProof := &state.Proof{
		ProofID:          &proofID,
		BatchNumber:      batchNum,
		BatchNumberFinal: batchNumFinal,
	}
	finalProof := &pb.FinalProof{}
	cfg := Config{}

	testCases := []struct {
		name    string
		setup   func(mox, *Aggregator)
		asserts func(*Aggregator)
	}{
		{
			name: "GetBatchByNumber error",
			setup: func(m mox, a *Aggregator) {
				m.stateMock.On("GetBatchByNumber", mock.Anything, batchNumFinal, nil).Run(func(args mock.Arguments) {
					// test is done, stop the sendFinalProof method
					a.exit()
					assert.True(a.verifyingProof)
				}).Return(nil, errors.New("banana")).Once()
			},
			asserts: func(a *Aggregator) {
				proof, ok := a.proverProofs[proverID]
				if assert.True(ok) {
					assert.Equal(proofID, proof.ID)
					assert.Equal(batchNum, proof.batchNum)
					assert.Equal(batchNumFinal, proof.batchNumFinal)
				}
				assert.False(a.verifyingProof)
			},
		},
		{
			name: "VerifyBatches error",
			setup: func(m mox, a *Aggregator) {
				m.stateMock.On("GetBatchByNumber", mock.Anything, batchNumFinal, nil).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(&finalBatch, nil).Once()
				expectedInputs := ethmanTypes.FinalProofInputs{
					FinalProof:       finalProof,
					NewLocalExitRoot: finalBatch.LocalExitRoot.Bytes(),
					NewStateRoot:     finalBatch.StateRoot.Bytes(),
				}
				m.ethTxManager.On("VerifyBatches", mock.Anything, batchNum-1, batchNumFinal, &expectedInputs).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(nil, errors.New("banana")).Once()
				m.stateMock.On("UpdateGeneratedProof", mock.Anything, recursiveProof, nil).Run(func(args mock.Arguments) {
					// test is done, stop the sendFinalProof method
					a.exit()
				}).Return(nil).Once()
			},
			asserts: func(a *Aggregator) {
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
				assert.False(a.verifyingProof)
			},
		},
		{
			name: "UpdateGeneratedProof error after VerifyBatches error",
			setup: func(m mox, a *Aggregator) {
				m.stateMock.On("GetBatchByNumber", mock.Anything, batchNumFinal, nil).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(&finalBatch, nil).Once()
				expectedInputs := ethmanTypes.FinalProofInputs{
					FinalProof:       finalProof,
					NewLocalExitRoot: finalBatch.LocalExitRoot.Bytes(),
					NewStateRoot:     finalBatch.StateRoot.Bytes(),
				}
				m.ethTxManager.On("VerifyBatches", mock.Anything, batchNum-1, batchNumFinal, &expectedInputs).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(nil, errors.New("banana")).Once()
				m.stateMock.On("UpdateGeneratedProof", mock.Anything, recursiveProof, nil).Run(func(args mock.Arguments) {
					// test is done, stop the sendFinalProof method
					a.exit()
				}).Return(errors.New("banana")).Once()
			},
			asserts: func(a *Aggregator) {
				proof, ok := a.proverProofs[proverID]
				if assert.True(ok) {
					assert.Equal(proofID, proof.ID)
					assert.Equal(batchNum, proof.batchNum)
					assert.Equal(batchNumFinal, proof.batchNumFinal)
				}
				assert.False(a.verifyingProof)
			},
		},
		{
			name: "DeleteGeneratedProofs error",
			setup: func(m mox, a *Aggregator) {
				m.stateMock.On("GetBatchByNumber", mock.Anything, batchNumFinal, nil).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(&finalBatch, nil).Once()
				expectedInputs := ethmanTypes.FinalProofInputs{
					FinalProof:       finalProof,
					NewLocalExitRoot: finalBatch.LocalExitRoot.Bytes(),
					NewStateRoot:     finalBatch.StateRoot.Bytes(),
				}
				m.ethTxManager.On("VerifyBatches", mock.Anything, batchNum-1, batchNumFinal, &expectedInputs).Return(tx, nil).Once()
				verifiedBatch := state.VerifiedBatch{
					BatchNumber: batchNumFinal,
				}
				m.stateMock.On("GetLastVerifiedBatch", mock.Anything, nil).Return(&verifiedBatch, nil).Once()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(batchNumFinal, nil).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.Anything, batchNum, batchNumFinal, nil).Run(func(args mock.Arguments) {
					// test is done, stop the sendFinalProof method
					a.exit()
				}).Return(errors.New("banana")).Once()

			},
			asserts: func(a *Aggregator) {
				proof, ok := a.proverProofs[proverID]
				if assert.True(ok) {
					assert.Equal(proofID, proof.ID)
					assert.Equal(batchNum, proof.batchNum)
					assert.Equal(batchNumFinal, proof.batchNumFinal)
				}
				assert.False(a.verifyingProof)
			},
		},
		{
			name: "nominal case",
			setup: func(m mox, a *Aggregator) {
				m.stateMock.On("GetBatchByNumber", mock.Anything, batchNumFinal, nil).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(&finalBatch, nil).Once()
				expectedInputs := ethmanTypes.FinalProofInputs{
					FinalProof:       finalProof,
					NewLocalExitRoot: finalBatch.LocalExitRoot.Bytes(),
					NewStateRoot:     finalBatch.StateRoot.Bytes(),
				}
				m.ethTxManager.On("VerifyBatches", mock.Anything, batchNum-1, batchNumFinal, &expectedInputs).Return(tx, nil).Once()
				verifiedBatch := state.VerifiedBatch{
					BatchNumber: batchNumFinal,
				}
				m.stateMock.On("GetLastVerifiedBatch", mock.Anything, nil).Return(&verifiedBatch, nil).Once()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(batchNumFinal, nil).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.Anything, batchNum, batchNumFinal, nil).Run(func(args mock.Arguments) {
					// test is done, stop the sendFinalProof method
					a.exit()
				}).Return(nil).Once()

			},
			asserts: func(a *Aggregator) {
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
				assert.False(a.verifyingProof)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stateMock := mocks.NewStateMock(t)
			ethTxManager := mocks.NewEthTxManager(t)
			etherman := mocks.NewEtherman(t)
			a, err := New(cfg, stateMock, ethTxManager, etherman)
			require.NoError(err)
			a.ctx, a.exit = context.WithCancel(context.Background())
			// add an ongoing proof to the map
			a.proverProofs[proverID] = proverProof{
				ID:            proofID,
				batchNum:      batchNum,
				batchNumFinal: batchNumFinal,
			}
			m := mox{
				stateMock:    stateMock,
				ethTxManager: ethTxManager,
				etherman:     etherman,
			}
			if tc.setup != nil {
				tc.setup(m, &a)
			}
			// send a final proof over the channel
			go func() {
				finalMsg := finalProofMsg{
					proverID:       proverID,
					recursiveProof: recursiveProof,
					finalProof:     finalProof,
				}
				a.finalProof <- finalMsg
			}()

			a.sendFinalProof()

			if tc.asserts != nil {
				tc.asserts(&a)
			}
		})
	}
}

func TestTryAggregateProofs(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	cfg := Config{
		VerifyProofInterval: configTypes.NewDuration(10000000),
	}
	proofID := "proofId"
	proverID := "proverID"
	recursiveProof := "recursiveProof"
	errBanana := errors.New("banana")
	ctx := context.Background()
	proof1 := state.Proof{
		Proof:       "proof1",
		BatchNumber: uint64(23),
	}
	proof2 := state.Proof{
		Proof:            "proof2",
		BatchNumberFinal: uint64(42),
	}
	testCases := []struct {
		name    string
		setup   func(mox, *Aggregator)
		asserts func(bool, *Aggregator, error)
	}{
		{
			name: "getAndLockProofsToAggregate returns generic error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetProofsToAggregate", mock.Anything, nil).Return(nil, nil, errBanana).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "getAndLockProofsToAggregate returns ErrNotFound",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetProofsToAggregate", mock.Anything, nil).Return(nil, nil, state.ErrNotFound).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.NoError(err)
			},
		},
		{
			name: "AggregatedProof prover error",
			setup: func(m mox, a *Aggregator) {
				dbTx := &mocks.DbTxMock{}
				m.stateMock.On("BeginStateTransaction", mock.Anything).Return(dbTx, nil).Twice()
				dbTx.On("Commit", mock.Anything).Return(nil).Twice()
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetProofsToAggregate", mock.Anything, nil).Return(&proof1, &proof2, nil).Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(nil, errBanana).Once()
				proof1GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.Anything, &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				proof2GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.Anything, &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.Anything, &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.False(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once().
					NotBefore(proof1GeneratingTrueCall)
				m.stateMock.
					On("UpdateGeneratedProof", mock.Anything, &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.False(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once().
					NotBefore(proof2GeneratingTrueCall)
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "not time to send final ok",
			setup: func(m mox, a *Aggregator) {
				dbTx := &mocks.DbTxMock{}
				m.stateMock.On("BeginStateTransaction", mock.Anything).Return(dbTx, nil).Twice()
				dbTx.On("Commit", mock.Anything).Return(nil).Twice()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetProofsToAggregate", mock.Anything, nil).Return(&proof1, &proof2, nil).Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.Anything, &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.Anything, &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.Anything, proofID).Return(recursiveProof, nil).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.Anything, proof1.BatchNumber, proof2.BatchNumberFinal, dbTx).Return(nil).Once()
				expectedInputProver := map[string]interface{}{
					"recursive_proof_1": proof1.Proof,
					"recursive_proof_2": proof2.Proof,
				}
				b, err := json.Marshal(expectedInputProver)
				require.NoError(err)
				expectedGenProof := state.Proof{
					BatchNumber:      proof1.BatchNumber,
					BatchNumberFinal: proof2.BatchNumberFinal,
					Prover:           &proverID,
					InputProver:      string(b),
					ProofID:          &proofID,
					Proof:            recursiveProof,
					Generating:       true,
				}
				m.stateMock.On("AddGeneratedProof", mock.Anything, &expectedGenProof, dbTx).Return(nil).Once()
				expectedUngenProof := state.Proof{
					BatchNumber:      proof1.BatchNumber,
					BatchNumberFinal: proof2.BatchNumberFinal,
					Prover:           &proverID,
					InputProver:      string(b),
					ProofID:          &proofID,
					Proof:            recursiveProof,
					Generating:       false,
				}
				m.stateMock.On("UpdateGeneratedProof", mock.Anything, &expectedUngenProof, nil).Return(nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.True(result)
				assert.NoError(err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stateMock := mocks.NewStateMock(t)
			ethTxManager := mocks.NewEthTxManager(t)
			etherman := mocks.NewEtherman(t)
			proverMock := mocks.NewProverMock(t)
			a, err := New(cfg, stateMock, ethTxManager, etherman)
			require.NoError(err)
			a.ctx, a.exit = context.WithCancel(context.Background())
			m := mox{
				stateMock:    stateMock,
				ethTxManager: ethTxManager,
				etherman:     etherman,
				proverMock:   proverMock,
			}
			if tc.setup != nil {
				tc.setup(m, &a)
			}
			a.resetVerifyProofTime()

			result, err := a.tryAggregateProofs(ctx, proverMock)

			if tc.asserts != nil {
				tc.asserts(result, &a, err)
			}
		})
	}
}
