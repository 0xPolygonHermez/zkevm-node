package ethtxmanager

import (
	"context"
	"math/big"
	"testing"

	ethman "github.com/0xPolygonHermez/zkevm-node/etherman"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestIncreaseGasLimit(t *testing.T) {
	actual := increaseGasLimit(100, 1)
	assert.Equal(t, uint64(101), actual)
}

func TestIncreaseGasPrice(t *testing.T) {
	actual := increaseGasPrice(big.NewInt(100), 1)
	assert.Equal(t, big.NewInt(101), actual)
}

func TestSequenceBatchesWithROEthman(t *testing.T) {
	observedLogs, err := log.InitTestLogger(log.Config{Level: "INFO"})
	require.NoError(t, err)
	defer log.DeinitTestLogger()
	// add a dummy log to increase the number of logs and verify that the
	// filter works
	log.Infof("bla")
	ethManRO, _, _, _, _ := ethman.NewSimulatedEtherman(ethman.Config{}, nil)
	txMan := New(Config{MaxSendBatchTxRetries: 2}, ethManRO) // 3 executions in total

	txMan.SequenceBatches(context.Background(), []ethmanTypes.Sequence{})

	observedLogs = observedLogs.
		FilterLevelExact(zapcore.ErrorLevel).
		FilterMessageSnippet(ethman.ErrIsReadOnlyMode.Error())
	logs := observedLogs.All()
	assert.Len(t, logs, 3)
}

func TestVerifyBatchWithROEthman(t *testing.T) {
	observedLogs, err := log.InitTestLogger(log.Config{Level: "INFO"})
	require.NoError(t, err)
	defer log.DeinitTestLogger()
	// add a dummy log to increase the number of logs and verify that the
	// filter works
	log.Infof("bla")
	ethManRO, _, _, _, _ := ethman.NewSimulatedEtherman(ethman.Config{}, nil)
	txMan := New(Config{MaxVerifyBatchTxRetries: 2}, ethManRO) // 3 executions in total

	txMan.VerifyBatch(context.Background(), 42, nil)

	observedLogs = observedLogs.
		FilterLevelExact(zapcore.ErrorLevel).
		FilterMessageSnippet(ethman.ErrIsReadOnlyMode.Error())
	logs := observedLogs.All()
	assert.Len(t, logs, 3)
}
