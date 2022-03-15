package compilesc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/spf13/afero"
)

const (
	solcImage = "ethereum/solc:0.8.12-alpine"
)

// Contracts is a map from contract name to its compilation results, as obtained
// from solc combined json output.
type Contracts map[string]Result

// Output is an equivalent of solc's combined json output.
type Output struct {
	Contracts Contracts
}

// Result contains the actual compilation results of a smart contract.
type Result struct {
	Bin string
}

// Manager handles smart contract compilation.
type Manager struct {
	aferoFs afero.Fs
}

// NewManager is the Manager constructor.
func NewManager(aferoFs afero.Fs) *Manager {
	return &Manager{
		aferoFs: aferoFs,
	}
}

// Run executes the compilation of the smart contracts. The given input can be
// a directory or a single sol file. In case of a directory, for each .sol
// file found in it or any of its subdirectories it calls the compile command
// and stores a .bin and .abi file with the results of the compilation.
func (cm *Manager) Run(input string) error {
	file, err := cm.aferoFs.Open(input)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("Could not close file %q, %v", file.Name(), err)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		output, err := cm.Compile(file.Name())
		if err != nil {
			return err
		}

		base := filepath.Dir(file.Name())
		return cm.writeContracts(output.Contracts, base)
	}

	return afero.Walk(cm.aferoFs, file.Name(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil || info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".sol" {
			return nil
		}
		output, err := cm.Compile(path)
		if err != nil {
			return err
		}
		base := filepath.Dir(path)
		return cm.writeContracts(output.Contracts, base)
	})
}

// Compile invokes solc on the given sol file and returns the results and a
// potential error.
func (cm *Manager) Compile(file string) (*Output, error) {
	log.Infof("Compiling %q...", file)

	c := exec.Command("docker", "run", "-i", solcImage, "-", "--combined-json", "bin")
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	stdin, err := c.StdinPipe()
	if err != nil {
		return nil, err
	}
	inputFile, err := cm.aferoFs.Open(file)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(stdin, inputFile)
	if err != nil {
		return nil, err
	}
	err = stdin.Close()
	if err != nil {
		return nil, err
	}

	envPath := os.Getenv("PATH")
	c.Env = []string{fmt.Sprintf("PATH=%s", envPath)}

	if err := c.Run(); err != nil {
		return nil, err
	}

	o := &Output{}
	if err := json.Unmarshal(stdout.Bytes(), o); err != nil {
		return nil, err
	}

	return o, nil
}

func (cm *Manager) writeContracts(contracts Contracts, base string) error {
	for name, result := range contracts {
		items := strings.Split(name, ":")
		r := []rune(items[len(items)-1])
		r[0] = unicode.ToLower(r[0])
		contractName := string(r)
		targetFile := path.Join(base, contractName)
		if err := cm.write(targetFile, result); err != nil {
			return err
		}
	}
	return nil
}

func (cm *Manager) write(file string, result Result) error {
	baseName := strings.TrimSuffix(file, filepath.Ext(file))

	binFileName := fmt.Sprintf("%s.bin", baseName)
	return afero.WriteFile(cm.aferoFs, binFileName, []byte(result.Bin), 0644)
}

// ReadBytecode reads the bytecode of the given contract.
func ReadBytecode(contractPath string) (string, error) {
	const basePath = "../../../test/contracts"

	_, currentFilename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("Could not get name of current file")
	}
	fullBasePath := path.Join(path.Dir(currentFilename), basePath)

	binPath := strings.Replace(contractPath, ".sol", ".bin", -1)

	content, err := os.ReadFile(path.Join(fullBasePath, binPath))
	if err != nil {
		return "", err
	}
	return string(content), nil
}
