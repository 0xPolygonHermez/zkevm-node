package metrics

const (
	//Endpoint the endpoint for exposing the metrics
	Endpoint = "/metrics"
	// ProfilingIndexEndpoint the endpoint for exposing the profiling metrics
	ProfilingIndexEndpoint = "/debug/pprof/"
	// ProfileEndpoint the endpoint for exposing the profile of the profiling metrics
	ProfileEndpoint = "/debug/pprof/profile"
	// ProfilingCmdEndpoint the endpoint for exposing the command-line of profiling metrics
	ProfilingCmdEndpoint = "/debug/pprof/cmdline"
	// ProfilingSymbolEndpoint the endpoint for exposing the symbol of profiling metrics
	ProfilingSymbolEndpoint = "/debug/pprof/symbol"
	// ProfilingTraceEndpoint the endpoint for exposing the trace of profiling metrics
	ProfilingTraceEndpoint = "/debug/pprof/trace"
)
