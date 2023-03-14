package vectors

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// LoadTxEventsSendBatchTestCases loads the calldata-test-vector.json
func LoadTxEventsSendBatchTestCases(path string) ([]TxEventsSendBatchTestCase, error) {
	var testCases []TxEventsSendBatchTestCase

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
