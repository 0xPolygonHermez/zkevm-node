package ethtxmanager

import (
	"math/big"
	"testing"

	ethman "github.com/0xPolygonHermez/zkevm-node/etherman"
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

func TestVerifyBatchWithROEthman(t *testing.T) {
	observedLogs, err := log.InitTestLogger(log.Config{Level: "INFO"})
	require.NoError(t, err)
	defer log.DeinitTestLogger()
	// add a dummy log to increase the number of logs and verify that the
	// filter works
	log.Infof("bla")
	ethManRO, _, _, _, _ := ethman.NewSimulatedEtherman(ethman.Config{}, nil)
	txMan := New(Config{MaxVerifyBatchTxRetries: 2}, ethManRO) // 3 executions in total

	txMan.VerifyBatch(42, nil)

	observedLogs = observedLogs.
		FilterLevelExact(zapcore.ErrorLevel).
		FilterMessageSnippet("etherman client in read-only mode, cannot send verifyBatch request")
	logs := observedLogs.All()
	assert.Len(t, logs, 3)
}
