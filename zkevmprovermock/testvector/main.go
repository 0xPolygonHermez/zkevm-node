package testvector

import (
	"encoding/json"
	"path/filepath"

	"github.com/spf13/afero"
)

// StateDBRaw contains raw test vector for state db.
type StateDBRaw struct {
	Entries []*StateDBRawEntry
}

// StateDBRawEntry contains raw test vector for state db.
type StateDBRawEntry struct {
	Keys         []string
	Values       []string
	ExpectedRoot []string
}

// Container is a wrapper for test vectors.
type Container struct {
	StateDBRaw *StateDBRaw
}

// NewContainer is the Container constructor.
func NewContainer(testVectorPath string, aferoFs afero.Fs) (*Container, error) {
	stateDBRaw, err := getStateDBRaw(testVectorPath, aferoFs)
	if err != nil {
		return nil, err
	}

	return &Container{
		StateDBRaw: stateDBRaw,
	}, nil
}

func getStateDBRaw(testVectorPath string, aferoFs afero.Fs) (*StateDBRaw, error) {
	filePath := filepath.Join(testVectorPath, "merkle-tree/smt-raw.json")
	contents, err := afero.ReadFile(aferoFs, filePath)
	if err != nil {
		return nil, err
	}
	var stateDBRaw StateDBRaw
	if err := json.Unmarshal(contents, &stateDBRaw.Entries); err != nil {
		return nil, err
	}
	return &stateDBRaw, nil
}
