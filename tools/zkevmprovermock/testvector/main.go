package testvector

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/afero"
)

// E2E contains all the test vectors.
type E2E struct {
	Items []*E2EItem
}

// E2EItem contains an end-to-end test vector.
type E2EItem struct {
	StateTransition *StateTransition
	GenesisRaw      *GenesisRaw
}

// StateTransition contains the human-friendly genesis.
type StateTransition struct {
	Genesis []*GenesisItem
}

// GenesisItem is an element of the human-frinsdly genesis.
type GenesisItem struct {
	Address  string
	Nonce    string
	Balance  string
	PvtKey   string
	Bytecode string
	Storage  map[string]string
}

// GenesisRaw contains the genesis definition in raw mode as a set of keys, values,
// and expected roots as seen by the StateDB mock service.
type GenesisRaw struct {
	Keys          []string
	Values        []string
	ExpectedRoots []string
}

// Container is a wrapper for test vectors.
type Container struct {
	E2E *E2E
}

// NewContainer is the Container constructor.
func NewContainer(testVectorPath string, aferoFs afero.Fs) (*Container, error) {
	e2e, err := getE2E(testVectorPath, aferoFs)
	if err != nil {
		return nil, err
	}

	return &Container{
		E2E: e2e,
	}, nil
}

// FindE2EGenesisRaw searches for the given key on all the genesisRaw items
// present, checking also that the given root was the root returned by the
// previous item.
func (c *Container) FindE2EGenesisRaw(inputKey, oldRoot string) (value, newRoot string, err error) {
	for _, item := range c.E2E.Items {
		if item.GenesisRaw != nil {
			for index, key := range item.GenesisRaw.Keys {
				if key == inputKey &&
					(index > 0 && oldRoot == item.GenesisRaw.ExpectedRoots[index-1] ||
						index == 0 && oldRoot == "") {
					return item.GenesisRaw.Values[index], item.GenesisRaw.ExpectedRoots[index], nil
				}
			}
		}
	}
	return "", "", fmt.Errorf("key %q not found for oldRoot %q", inputKey, oldRoot)
}

func getE2E(testVectorPath string, aferoFs afero.Fs) (*E2E, error) {
	e2e := &E2E{}

	err := afero.Walk(aferoFs, testVectorPath, func(wpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil || info.IsDir() {
			return nil
		}
		e2eFile, err := getE2EFile(wpath, aferoFs)
		if err != nil {
			return err
		}
		e2e.Items = append(e2e.Items, e2eFile.Items...)

		return nil
	})

	return e2e, err
}

func getE2EFile(filePath string, aferoFs afero.Fs) (*E2E, error) {
	contents, err := afero.ReadFile(aferoFs, filePath)
	if err != nil {
		return nil, err
	}
	var e2e E2E
	if err := json.Unmarshal(contents, &e2e.Items); err != nil {
		return nil, err
	}
	return &e2e, nil
}
