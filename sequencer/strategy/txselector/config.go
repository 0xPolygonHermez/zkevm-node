package txselector

// Type different types of tx selection logic
type Type string

const (
	// AcceptAllType strategy accepts all txs
	AcceptAllType Type = "acceptall"
	// BaseType strategy that have basic selection algorithm and can accept different sorting algorithms and profitability checkers
	BaseType Type = "base"
)

// Config for the tx selector configuration
type Config struct {
	Type         Type         `mapstructure:"Type"`
	TxSorterType TxSorterType `mapstructure:"TxSorterType"`
}
