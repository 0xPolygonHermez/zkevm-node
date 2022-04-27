package instrumentation

// Trace contains a TX trace
type Trace struct {
	TraceAddress []uint64     `json:"traceAddress"`
	SubTraces    uint64       `json:"subtraces"`
	Action       TraceAction  `json:"action"`
	Result       *TraceResult `json:"result,omitempty"`
	Error        *string      `json:"error,omitempty"`
	Type         string       `json:"type"`
}

// TraceAction contains information about an action of a transaction
type TraceAction struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Value    uint64 `json:"value"`
	Gas      uint64 `json:"gas"`
	Input    []byte `json:"input"`
	CallType string `json:"callType"`
}

// TraceResult contains information about the execution result of a transaction
type TraceResult struct {
	GasUsed uint64 `json:"gasUsed"`
	Output  []byte `json:"output"`
}

// VMTrace contains a EVM trace
type VMTrace struct {
	// ParentStep int           `json:"parentStep"`
	Code       []byte        `json:"code"`
	Operations []VMOperation `json:"operations"`
	// Subs       []VMTrace     `json:"subs"`
}

// VMOperation contains one opcode execution metadata
type VMOperation struct {
	Pc          uint64              `json:"pc"`
	Instruction byte                `json:"instruction"`
	GasCost     uint64              `json:"gasCost"`
	Executed    VMExecutedOperation `json:"executed"`
	Sub         *VMTrace            `json:"sub"`
}

// VMExecutedOperation Contains information about data modified by an VMOperation
type VMExecutedOperation struct {
	GasUsed   uint64     `json:"gasUsed"`
	StackPush []uint64   `json:"stackPush"`
	MemDiff   MemoryDiff `json:"memDiff"`
	StoreDiff StoreDiff  `json:"storeDiff"`
}

// MemoryDiff contains modified memory data
type MemoryDiff struct {
	Offset uint64 `json:"offset"`
	Data   []byte `json:"data"`
}

// StoreDiff contains modified storage data
type StoreDiff struct {
	Location uint64 `json:"location"`
	Value    uint64 `json:"value"`
}
