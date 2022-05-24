package executor

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/hermeznetwork/hermez-core/state/runtime/executor/js"
	"github.com/stretchr/testify/require"
)

func Test_Trace(t *testing.T) {
	var (
		trace  Trace
		tracer Tracer
	)

	traceFile, err := os.Open("demo_trace.json")
	require.NoError(t, err)
	defer traceFile.Close()

	tracerFile, err := os.Open("tracer.json")
	require.NoError(t, err)
	defer tracerFile.Close()

	byteValue, err := ioutil.ReadAll(traceFile)
	require.NoError(t, err)

	err = json.Unmarshal(byteValue, &trace)
	require.NoError(t, err)

	byteCode, err := ioutil.ReadAll(tracerFile)
	require.NoError(t, err)

	err = json.Unmarshal(byteCode, &tracer)
	require.NoError(t, err)

	jsTracer, err := js.NewJsTracer(string(tracer.Code))
	require.NoError(t, err)
	result, err := jsTracer.GetResult()
	require.NoError(t, err)
	log.Debugf("%v", result)
}
