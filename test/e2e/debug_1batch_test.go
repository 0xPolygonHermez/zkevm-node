package e2e

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/stretchr/testify/require"
)

func Test_DebugFirstBatch(t *testing.T) {
	// Send a request to the locally running zkevm-node to trace a block #1
	const zkevmAddr = "http://localhost:8123"

	debugOptions := map[string]interface{}{
		"tracer": "callTracer",
		"tracerConfig": map[string]interface{}{
			"onlyTopCall": false,
			"withLog":     true,
		},
	}

	response, err := client.JSONRPCCall(zkevmAddr, "debug_traceBlockByNumber", hex.EncodeBig(big.NewInt(1)), debugOptions)
	require.NoError(t, err)

	raw, err := json.MarshalIndent(response, "", "  ")
	require.NoError(t, err)

	t.Log(string(raw))
}
