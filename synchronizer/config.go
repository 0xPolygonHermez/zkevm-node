package synchronizer

import "time"

// Duration is a wrapper type that parses time duration from text.
type Duration struct {
	time.Duration `validate:"required"`
}

// UnmarshalText unmarshalls time duration from text.
func (d *Duration) UnmarshalText(data []byte) error {
	duration, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	d.Duration = duration
	return nil
}

// Config represents the configuration of the synchronizer
type Config struct {
	// SyncInterval is the delay interval between reading new rollup information
	SyncInterval Duration `mapstructure:"SyncInterval"`

	// SyncChunkSize is the number of blocks to sync on each chunk
	SyncChunkSize uint64 `mapstructure:"SyncChunkSize"`
}
