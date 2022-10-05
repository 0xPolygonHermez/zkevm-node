package proverclient

// Config represents the configuration of the prover clients
type Config struct {
	// ProverURIs URIs to get access to the prover clients
	ProverURIs []string `mapstructure:"ProverURIs"`
}
