package vectors

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LoadStateTransitionTestCaseV2 loads the state-transition JSON file into a
// StateTransitionTestCaseV2 instance
func LoadStateTransitionTestCaseV2(path string) (StateTransitionTestCaseV2, error) {
	var testCase StateTransitionTestCaseV2

	jsonFile, err := os.Open(filepath.Clean(path))
	if err != nil {
		return testCase, err
	}
	defer func() { _ = jsonFile.Close() }()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return testCase, err
	}

	err = json.Unmarshal(bytes, &testCase)
	if err != nil {
		return testCase, err
	}
	if testCase.Description == "" {
		testCase.Description = strings.Replace(filepath.Base(path), ".json", "", 1)
	}

	return testCase, nil
}
