package ethtxmanager

import "github.com/0xPolygonHermez/zkevm-node/config/types"

// Config is configuration for ethereum transaction manager
type Config struct {
	// FrequencyToMonitorTxs frequency of the resending failed txs
	FrequencyToMonitorTxs types.Duration `mapstructure:"FrequencyToMonitorTxs"`
	// WaitTxToBeMined time to wait after transaction was sent to the ethereum
	WaitTxToBeMined types.Duration `mapstructure:"WaitTxToBeMined"`

	// PrivateKeys defines all the key store files that are going
	// to be read in order to provide the private keys to sign the L1 txs
	PrivateKeys []KeystoreFileConfig `mapstructure:"PrivateKeys"`
}

// KeystoreFileConfig has all the information needed to load a private key from a key store file
type KeystoreFileConfig struct {
	// Path is the file path for the key store file
	Path string `mapstructure:"Path"`

	// Password is the password to decrypt the key store file
	Password string `mapstructure:"Password"`
}
