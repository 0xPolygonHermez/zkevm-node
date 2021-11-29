package etherman

import (
	"io/ioutil"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func decryptKeystore(path, pw string) (*keystore.Key, error) {
	keystoreEncrypted, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keystoreEncrypted, pw)
	if err != nil {
		return nil, err
	}
	return key, nil
}
