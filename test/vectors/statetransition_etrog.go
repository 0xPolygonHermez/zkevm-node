package vectors

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// LoadStateTransitionTestCasesEtrog loads the state-transition tests cases
func LoadStateTransitionTestCasesEtrog(path string) ([]StateTransitionTestCaseEtrog, error) {
	var testCases []StateTransitionTestCaseEtrog

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
