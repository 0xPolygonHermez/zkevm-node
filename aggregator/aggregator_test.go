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
				m.stateMock.On("CleanupGeneratedProofs", mock.Anything, batchNumFinal, nil).Run(func(args mock.Arguments) {
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
				m.stateMock.On("CleanupGeneratedProofs", mock.Anything, batchNumFinal, nil).Run(func(args mock.Arguments) {
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
	proverCtx := context.WithValue(context.Background(), "owner", "prover")
	matchProverCtxFn := func(ctx context.Context) bool { return ctx.Value("owner") == "prover" }
	matchAggregatorCtxFn := func(ctx context.Context) bool { return ctx.Value("owner") == "aggregator" }
	batchNum := uint64(23)
	batchNumFinal := uint64(42)
	proof1 := state.Proof{
		Proof:       "proof1",
		BatchNumber: batchNum,
	}
	proof2 := state.Proof{
		Proof:            "proof2",
		BatchNumberFinal: batchNumFinal,
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
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(nil, nil, errBanana).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
			},
		},
		{
			name: "getAndLockProofsToAggregate returns ErrNotFound",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(nil, nil, state.ErrNotFound).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.NoError(err)
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
			},
		},
		{
			name: "AggregatedProof prover error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return("addr")
				dbTx := &mocks.DbTxMock{}
				lockProofsTxBegin := m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Once()
				lockProofsTxCommit := dbTx.On("Commit", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				proof1GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				proof2GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(nil, errBanana).Once()
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchAggregatorCtxFn)).Return(dbTx, nil).Once().NotBefore(lockProofsTxBegin)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.False(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once().
					NotBefore(proof1GeneratingTrueCall)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.False(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once().
					NotBefore(proof2GeneratingTrueCall)
				dbTx.On("Commit", mock.MatchedBy(matchAggregatorCtxFn)).Return(nil).Once().NotBefore(lockProofsTxCommit)
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
			},
		},
		{
			name: "WaitRecursiveProof prover error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return("addr")
				dbTx := &mocks.DbTxMock{}
				lockProofsTxBegin := m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Once()
				lockProofsTxCommit := dbTx.On("Commit", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				proof1GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				proof2GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return("", errBanana).Once()
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchAggregatorCtxFn)).Return(dbTx, nil).Once().NotBefore(lockProofsTxBegin)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.False(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once().
					NotBefore(proof1GeneratingTrueCall)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.False(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once().
					NotBefore(proof2GeneratingTrueCall)
				dbTx.On("Commit", mock.MatchedBy(matchAggregatorCtxFn)).Return(nil).Once().NotBefore(lockProofsTxCommit)
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
			},
		},
		{
			name: "unlockProofsToAggregate error after WaitRecursiveProof prover error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return(proverID)
				dbTx := &mocks.DbTxMock{}
				lockProofsTxBegin := m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Once()
				dbTx.On("Commit", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				proof1GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return("", errBanana).Once()
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchAggregatorCtxFn)).Return(dbTx, nil).Once().NotBefore(lockProofsTxBegin)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.False(args[1].(*state.Proof).Generating)
					}).
					Return(errBanana).
					Once().
					NotBefore(proof1GeneratingTrueCall)
				dbTx.On("Rollback", mock.MatchedBy(matchAggregatorCtxFn)).Return(nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
				proof, ok := a.proverProofs[proverID]
				if assert.True(ok) {
					assert.Equal(proofID, proof.ID)
					assert.Equal(batchNum, proof.batchNum)
					assert.Equal(batchNumFinal, proof.batchNumFinal)

				}
			},
		},
		{
			name: "rollback after DeleteGeneratedProofs error in db transaction",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return("addr")
				dbTx := &mocks.DbTxMock{}
				lockProofsTxBegin := m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Twice()
				lockProofsTxCommit := dbTx.On("Commit", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				proof1GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				proof2GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return(recursiveProof, nil).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.MatchedBy(matchProverCtxFn), proof1.BatchNumber, proof2.BatchNumberFinal, dbTx).Return(errBanana).Once()
				dbTx.On("Rollback", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchAggregatorCtxFn)).Return(dbTx, nil).Once().NotBefore(lockProofsTxBegin)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.False(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once().
					NotBefore(proof1GeneratingTrueCall)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.False(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once().
					NotBefore(proof2GeneratingTrueCall)
				dbTx.On("Commit", mock.MatchedBy(matchAggregatorCtxFn)).Return(nil).Once().NotBefore(lockProofsTxCommit)
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
			},
		},
		{
			name: "rollback after AddGeneratedProof error in db transaction",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return("addr")
				dbTx := &mocks.DbTxMock{}
				lockProofsTxBegin := m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Twice()
				lockProofsTxCommit := dbTx.On("Commit", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				proof1GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				proof2GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return(recursiveProof, nil).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.MatchedBy(matchProverCtxFn), proof1.BatchNumber, proof2.BatchNumberFinal, dbTx).Return(nil).Once()
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), mock.Anything, dbTx).Return(errBanana).Once()
				dbTx.On("Rollback", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchAggregatorCtxFn)).Return(dbTx, nil).Once().NotBefore(lockProofsTxBegin)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.False(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once().
					NotBefore(proof1GeneratingTrueCall)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.False(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once().
					NotBefore(proof2GeneratingTrueCall)
				dbTx.On("Commit", mock.MatchedBy(matchAggregatorCtxFn)).Return(nil).Once().NotBefore(lockProofsTxCommit)
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
			},
		},
		{
			name: "not time to send final ok",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				dbTx := &mocks.DbTxMock{}
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Twice()
				dbTx.On("Commit", mock.MatchedBy(matchProverCtxFn)).Return(nil).Twice()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.True(args[1].(*state.Proof).Generating)
					}).
					Return(nil).
					Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return(recursiveProof, nil).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.MatchedBy(matchProverCtxFn), proof1.BatchNumber, proof2.BatchNumberFinal, dbTx).Return(nil).Once()
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
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), &expectedGenProof, dbTx).Return(nil).Once()
				expectedUngenProof := state.Proof{
					BatchNumber:      proof1.BatchNumber,
					BatchNumberFinal: proof2.BatchNumberFinal,
					Prover:           &proverID,
					InputProver:      string(b),
					ProofID:          &proofID,
					Proof:            recursiveProof,
					Generating:       false,
				}
				m.stateMock.On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &expectedUngenProof, nil).Return(nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.True(result)
				assert.NoError(err)
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
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
			aggregatorCtx := context.WithValue(context.Background(), "owner", "aggregator")
			a.ctx, a.exit = context.WithCancel(aggregatorCtx)
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

			result, err := a.tryAggregateProofs(proverCtx, proverMock)

			if tc.asserts != nil {
				tc.asserts(result, &a, err)
			}
		})
	}
}

func TestTryGenerateBatchProof(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	cfg := Config{
		VerifyProofInterval:        configTypes.NewDuration(10000000),
		TxProfitabilityCheckerType: ProfitabilityAcceptAll,
	}
	pubAddr := common.BytesToAddress([]byte("pubAdddr"))
	lastVerifiedBatchNum := uint64(22)
	batchNum := uint64(23)
	lastVerifiedBatch := state.VerifiedBatch{
		BatchNumber: lastVerifiedBatchNum,
	}
	latestBatch := state.Batch{
		BatchNumber: lastVerifiedBatchNum,
	}
	batchToProve := state.Batch{
		BatchNumber: batchNum,
	}
	proofID := "proofId"
	proverID := "proverID"
	recursiveProof := "recursiveProof"
	errBanana := errors.New("banana")
	proverCtx := context.WithValue(context.Background(), "owner", "prover")
	matchProverCtxFn := func(ctx context.Context) bool { return ctx.Value("owner") == "prover" }
	matchAggregatorCtxFn := func(ctx context.Context) bool { return ctx.Value("owner") == "aggregator" }
	testCases := []struct {
		name    string
		setup   func(mox, *Aggregator)
		asserts func(bool, *Aggregator, error)
	}{
		{
			name: "getAndLockBatchToProve returns generic error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(nil, errBanana).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
			},
		},
		{
			name: "getAndLockBatchToProve returns ErrNotFound",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(nil, state.ErrNotFound).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.NoError(err)
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
			},
		},
		{
			name: "BatchProof prover error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&lastVerifiedBatch, nil).Once()
				m.stateMock.On("GetVirtualBatchToProve", mock.MatchedBy(matchProverCtxFn), lastVerifiedBatchNum, nil).Return(&batchToProve, nil).Once()
				expectedGenProof := state.Proof{
					BatchNumber:      batchToProve.BatchNumber,
					BatchNumberFinal: batchToProve.BatchNumber,
					Prover:           &proverID,
					Generating:       true,
				}
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), &expectedGenProof, nil).Return(nil).Once()
				m.stateMock.On("GetBatchByNumber", mock.Anything, lastVerifiedBatchNum, nil).Return(&latestBatch, nil).Twice()
				m.etherman.On("GetPublicAddress").Return(pubAddr, nil).Twice()
				expectedInputProver, err := a.buildInputProver(context.Background(), &batchToProve)
				require.NoError(err)
				m.proverMock.On("BatchProof", expectedInputProver).Return(nil, errBanana).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.MatchedBy(matchAggregatorCtxFn), expectedGenProof.BatchNumber, expectedGenProof.BatchNumberFinal, nil).Return(nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
			},
		},
		{
			name: "WaitRecursiveProof prover error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&lastVerifiedBatch, nil).Once()
				m.stateMock.On("GetVirtualBatchToProve", mock.MatchedBy(matchProverCtxFn), lastVerifiedBatchNum, nil).Return(&batchToProve, nil).Once()
				expectedGenProof := state.Proof{
					BatchNumber:      batchToProve.BatchNumber,
					BatchNumberFinal: batchToProve.BatchNumber,
					Prover:           &proverID,
					Generating:       true,
				}
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), &expectedGenProof, nil).Return(nil).Once()
				m.stateMock.On("GetBatchByNumber", mock.Anything, lastVerifiedBatchNum, nil).Return(&latestBatch, nil).Twice()
				m.etherman.On("GetPublicAddress").Return(pubAddr, nil).Twice()
				expectedInputProver, err := a.buildInputProver(context.Background(), &batchToProve)
				require.NoError(err)
				m.proverMock.On("BatchProof", expectedInputProver).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return("", errBanana).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.MatchedBy(matchAggregatorCtxFn), expectedGenProof.BatchNumber, expectedGenProof.BatchNumberFinal, nil).Return(nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
			},
		},
		{
			name: "DeleteGeneratedProofs error after WaitRecursiveProof prover error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return(proverID)
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&lastVerifiedBatch, nil).Once()
				m.stateMock.On("GetVirtualBatchToProve", mock.MatchedBy(matchProverCtxFn), lastVerifiedBatchNum, nil).Return(&batchToProve, nil).Once()
				expectedGenProof := state.Proof{
					BatchNumber:      batchToProve.BatchNumber,
					BatchNumberFinal: batchToProve.BatchNumber,
					Prover:           &proverID,
					Generating:       true,
				}
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), &expectedGenProof, nil).Return(nil).Once()
				m.stateMock.On("GetBatchByNumber", mock.Anything, lastVerifiedBatchNum, nil).Return(&latestBatch, nil).Twice()
				m.etherman.On("GetPublicAddress").Return(pubAddr, nil).Twice()
				expectedInputProver, err := a.buildInputProver(context.Background(), &batchToProve)
				require.NoError(err)
				m.proverMock.On("BatchProof", expectedInputProver).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return("", errBanana).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.MatchedBy(matchAggregatorCtxFn), expectedGenProof.BatchNumber, expectedGenProof.BatchNumberFinal, nil).Return(errBanana).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
				proof, ok := a.proverProofs[proverID]
				if assert.True(ok) {
					assert.Equal(proofID, proof.ID)
					assert.Equal(batchNum, proof.batchNum)
					assert.Equal(batchNum, proof.batchNumFinal)

				}
			},
		},
		{
			name: "not time to send final ok",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("ID").Return(proverID).Times(3)
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&lastVerifiedBatch, nil).Once()
				m.stateMock.On("GetVirtualBatchToProve", mock.MatchedBy(matchProverCtxFn), lastVerifiedBatchNum, nil).Return(&batchToProve, nil).Once()
				expectedGenProof := state.Proof{
					BatchNumber:      batchToProve.BatchNumber,
					BatchNumberFinal: batchToProve.BatchNumber,
					Prover:           &proverID,
					Generating:       true,
				}
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), &expectedGenProof, nil).Return(nil).Once()
				m.stateMock.On("GetBatchByNumber", mock.Anything, lastVerifiedBatchNum, nil).Return(&latestBatch, nil).Twice()
				m.etherman.On("GetPublicAddress").Return(pubAddr, nil).Twice()
				expectedInputProver, err := a.buildInputProver(context.Background(), &batchToProve)
				require.NoError(err)
				m.proverMock.On("BatchProof", expectedInputProver).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return(recursiveProof, nil).Once()
				b, err := json.Marshal(expectedInputProver)
				require.NoError(err)
				expectedUngenProof := state.Proof{
					BatchNumber:      batchToProve.BatchNumber,
					BatchNumberFinal: batchToProve.BatchNumber,
					Prover:           &proverID,
					InputProver:      string(b),
					ProofID:          &proofID,
					Proof:            recursiveProof,
					Generating:       false,
				}
				m.stateMock.On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &expectedUngenProof, nil).Return(nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.True(result)
				assert.NoError(err)
				_, ok := a.proverProofs[proverID]
				assert.False(ok)
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
			aggregatorCtx := context.WithValue(context.Background(), "owner", "aggregator")
			a.ctx, a.exit = context.WithCancel(aggregatorCtx)
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

			result, err := a.tryGenerateBatchProof(proverCtx, proverMock)

			if tc.asserts != nil {
				tc.asserts(result, &a, err)
			}
		})
	}
}
