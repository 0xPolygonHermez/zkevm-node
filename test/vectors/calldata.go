package vectors

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// LoadCallDataTestCases loads the calldata-test-vector.json
func LoadCallDataTestCases(path string) ([]CallDataTestCase, error) {
	var testCases []CallDataTestCase

	jsonFile, err := os.Open(filepath.Clean(path))
	if err != nil {
		return testCases, err
	}
	defer func() { _ = jsonFile.Close() }()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return testCases, err
	}

	err = json.Unmarshal(bytes, &testCases)
	if err != nil {
		return testCases, err
	}

	return testCases, nil
}
