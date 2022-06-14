package instrumentation

type ExecutorTrace struct {
	Context Context `json:"context"`
	Steps   []Step  `json:"steps"`
}

type Context struct {
	Type         string `json:"type"`
	From         string `json:"from"`
	To           string `json:"to"`
	Input        string `json:"input"`
	Gas          string `json:"gas"`
	Value        string `json:"value"`
	Output       string `json:"output"`
	Nonce        uint64 `json:"nonce"`
	GasPrice     string `json:"gasPrice"`
	ChainID      uint64 `json:"chainId"`
	OldStateRoot string `json:"oldStateRoot"`
	Time         uint64 `json:"time"`
	GasUsed      string `json:"gasUsed"`
}

type Step struct {
	StateRoot  string   `json:"stateRoot"`
	Depth      int      `json:"depth"`
	Pc         uint64   `json:"pc"`
	Gas        string   `json:"gas"`
	OpCode     string   `json:"opcode"`
	Refund     string   `json:"refund"`
	Op         string   `json:"op"`
	Error      string   `json:"error"`
	Contract   Contract `json:"contract"`
	GasCost    string   `json:"gasCost"`
	Stack      []string `json:"stack"`
	Memory     []string `json:"memory"`
	ReturnData string   `json:"returnData"`
}

type Contract struct {
	Address string `json:"address"`
	Caller  string `json:"caller"`
	Value   string `json:"value"`
	Input   string `json:"input"`
	Gas     string `json:"gas"`
}

type Tracer struct {
	Code string `json:"tracer"`
}
