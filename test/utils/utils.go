package utils

import (
	"fmt"
	"strings"

	"github.com/hermeznetwork/hermez-core/log"
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
