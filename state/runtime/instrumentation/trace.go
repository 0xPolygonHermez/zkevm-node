package instrumentation

/*
type Trace struct {
	Output          string      `json:"output"`
	TransactionHash common.Hash `json:"transactionHash"`
	VmTrace         VMTrace     `json:"vmTrace"`
	Action          Action      `json:"action"`
}
*/

// VMTrace contains a EVM trace
type VMTrace struct {
	ParentStep int           `json:"parentStep"`
	Code       string        `json:"code"`
	Operations []VMOperation `json:"operations"`
	Subs       []VMTrace     `json:"subs"`
}

// VMOperation contains one opcode execution metadata
type VMOperation struct {
	Pc          uint64              `json:"pc"`
	Instruction byte                `json:"instruction"`
	GasCost     uint64              `json:"gasCost"`
	Executed    VMExecutedOperation `json:"executed"`
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
	Data   uint64 `json:"data"`
}

// StoreDiff contains modified storage data
type StoreDiff struct {
	Location uint64 `json:"location"`
	Value    uint   `json:"value"`
}

/*
type Action struct {
	CallType       string         `json:"calltype"`
	IncludeInTrace bool           `json:"includeInTrace"`
	IsPrecompiled  bool           `json:"isPrecompiled"`
	Type           string         `json:"type"`
	CreationMethod string         `json:"creationMethod"`
	From           common.Address `json:"from"`
	To             common.Address `json:"to"`
	Gas            uint64         `json:"gas"`
	Value          uint64         `json:"value"`
	Input          []byte         `json:"input"`
	SubTraces      []VMTrace      `json:"subTraces"`
	Result         Result         `json:"result"`
	Author         common.Address `json:"author"`
	RewardType     string         `json:"rewardType"`
	Error          string         `json:"error"`
}

type Result struct {
	GasUsed uint64         `json:"gasUsed"`
	Output  []byte         `json:"output"`
	Address common.Address `json:"address"`
	Code    []byte         `json:"code"`
}
*/
