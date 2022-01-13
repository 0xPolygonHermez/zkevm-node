package txselector

// TxSelectorType different types of tx selection logic
type TxSelectorType string

const (
	// AcceptAll strategy accepts all txs
	AcceptAll TxSelectorType = "acceptall"
	// Base strategy that have basic selection algorithm and can accept different sorting algorithms and profitability checkers
	Base = "base"
)

// Config for the tx selector configuration
type Config struct {
	TxSelectorType TxSelectorType `mapstructure:"TxSelectorType"`
	TxSorterType   TxSorterType   `mapstructure:"TxSorterType"`
}
