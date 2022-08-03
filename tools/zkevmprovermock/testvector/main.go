package testvector

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/afero"
)

// E2E contains all the test vectors.
type E2E struct {
	Items []*E2EItem
}

// E2EItem contains an end-to-end test vector.
type E2EItem struct {
	GenesisRaw []*state.GenesisAction
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

// FindValue searches for the given key on all the genesisRaw items present,
// checking also that the given root was the root returned by the previous item.
// If both the value and the root of the previous item match it returns the
// associated value and new root.
func (c *Container) FindValue(inputKey, oldRoot string) (value, newRoot string, err error) {
	zero := common.HexToHash("").String()
	for _, item := range c.E2E.Items {
		for index, action := range item.GenesisRaw {
			if action.Key == inputKey &&
				(index > 0 && oldRoot == item.GenesisRaw[index-1].Root ||
					index == 0 && oldRoot == zero) {
				return item.GenesisRaw[index].Value, item.GenesisRaw[index].Root, nil
			}
		}
	}
	return "", "", fmt.Errorf("key %q not found for oldRoot %q", inputKey, oldRoot)
}

// FindBytecode searches for the given key on the value fields of all the
// genesisRaw items present and returns the associated bytecode field on match.
func (c *Container) FindBytecode(inputKey string) (bytecode string, err error) {
	for _, item := range c.E2E.Items {
		for index, action := range item.GenesisRaw {
			if action.Value == inputKey && action.Bytecode != "" {
				return item.GenesisRaw[index].Bytecode, nil
			}
		}
	}
	return "", fmt.Errorf("bytecode for key %q not found", inputKey)
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
