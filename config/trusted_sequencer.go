package config

// TrustedSequencerConfig is the configuration struct for Sequencer
type TrustedSequencerConfig struct {
	TrustedSequencerURL string `mapstructure:"TrustedSequencerURL"`
	IsTrustedSequencer  bool   `mapstructure:"IsTrustedSequencer"`
}
