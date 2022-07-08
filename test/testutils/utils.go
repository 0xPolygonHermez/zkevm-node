package testutils

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/spf13/afero"
)

// CreateTestFiles creates the files in the map (path -> content) using the
// given afero file system.
func CreateTestFiles(appFs afero.Fs, files map[string]string) error {
	for path, content := range files {
		f, err := appFs.Create(path)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Errorf("Could not close %s: %v", f.Name(), err)
			}
		}()
		_, err = f.WriteString(content)

		if err != nil {
			return err
		}
	}
	return nil
}

// CheckError checks the given error taking into account if it was expected and
// potentially the message it should carry.
func CheckError(err error, expected bool, msg string) error {
	if !expected && err != nil {
		return fmt.Errorf("Unexpected error %v", err)
	}
	if expected {
		if err == nil {
			return fmt.Errorf("Expected error didn't happen")
		}
		if msg == "" {
			return fmt.Errorf("Expected error message not defined")
		}
		if !strings.HasPrefix(err.Error(), msg) {
			return fmt.Errorf("Wrong error, expected %q, got %q", msg, err.Error())
		}
	}
	return nil
}

// ReadBytecode reads the bytecode of the given contract.
func ReadBytecode(contractPath string) (string, error) {
	const basePath = "../../test/contracts/bin"

	_, currentFilename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("Could not get name of current file")
	}
	fullBasePath := path.Join(path.Dir(currentFilename), basePath)

	content, err := os.ReadFile(path.Join(fullBasePath, contractPath))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// GetEnv reads an environment variable, returning a given default value if not
// present.
func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
