package types

// KeystoreFileConfig has all the information needed to load a private key from a key store file
type KeystoreFileConfig struct {
	// Path is the file path for the key store file
	Path string `mapstructure:"Path"`

	// Password is the password to decrypt the key store file
	Password string `mapstructure:"Password"`
}
