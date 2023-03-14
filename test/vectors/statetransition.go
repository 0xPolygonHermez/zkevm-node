package vectors

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// LoadStateTransitionTestCases loads the state-transition.json into a
// StateTransitionVector instance
func LoadStateTransitionTestCases(path string) ([]StateTransitionTestCase, error) {
	var testCases []StateTransitionTestCase

	jsonFile, err := os.Open(filepath.Clean(path))
	if err != nil {
		return testCases, err
	}
	defer func() { _ = jsonFile.Close() }()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return testCases, err
	}

	err = json.Unmarshal(bytes, &testCases)
	if err != nil {
		return testCases, err
	}

	return testCases, nil
}
