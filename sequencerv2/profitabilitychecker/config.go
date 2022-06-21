package profitabilitychecker

// Config for profitability checker
type Config struct {
	// SendBatchesEvenWhenNotProfitable if true -> send unprofitable batch
	SendBatchesEvenWhenNotProfitable bool `mapstructure:"SendBatchesEvenWhenNotProfitable"`
}
