package compilesc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	npm "github.com/aquasecurity/go-npm-version/pkg"
	"github.com/hermeznetwork/hermez-core/log"
	"github.com/spf13/afero"
)

const (
	solcImage = "ethereum/solc:%s-alpine"
)

var (
	validSolidityVersions = []string{
		"0.8.13", "0.8.12", "0.8.11", "0.8.10", "0.8.9", "0.8.8", "0.8.7", "0.8.6", "0.8.5", "0.8.4", "0.8.3", "0.8.2", "0.8.1", "0.8.0",
		"0.7.6", "0.7.5", "0.7.4", "0.7.3", "0.7.2", "0.7.1", "0.7.0",
		"0.6.12", "0.6.11", "0.6.10", "0.6.9", "0.6.8", "0.6.7", "0.6.6", "0.6.5", "0.6.4", "0.6.3", "0.6.2", "0.6.1", "0.6.0",
		"0.5.17", "0.5.16", "0.5.15", "0.5.14", "0.5.13", "0.5.12", "0.5.11", "0.5.10", "0.5.9", "0.5.8", "0.5.7", "0.5.6", "0.5.5", "0.5.4", "0.5.3", "0.5.2", "0.5.1", "0.5.0",
		"0.4.26", "0.4.25", "0.4.24", "0.4.23", "0.4.22", "0.4.21", "0.4.20", "0.4.19", "0.4.18", "0.4.17", "0.4.16", "0.4.15", "0.4.14", "0.4.13", "0.4.12", "0.4.11", "0.4.10", "0.4.9", "0.4.8", "0.4.7", "0.4.6", "0.4.5", "0.4.4", "0.4.3", "0.4.2", "0.4.1", "0.4.0",
		"0.3.6", "0.3.5", "0.3.4", "0.3.3", "0.3.2", "0.3.1", "0.3.0",
		"0.2.2", "0.2.1", "0.2.0",
		"0.1.7", "0.1.6", "0.1.5", "0.1.4", "0.1.3", "0.1.2",
	}
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

	solidityVersion, err := cm.getSolidityVersion(file)
	if err != nil {
		return nil, err
	}
	dockerImage := fmt.Sprintf(solcImage, solidityVersion)
	log.Debug("dockerImage ", dockerImage)

	fileDir := filepath.Dir(file)
	fileName := filepath.Base(file)

	c := exec.Command("docker", "run",
		"--rm",
		"-v", fmt.Sprintf("%s:/sources", fileDir),
		"-i", dockerImage,
		"-o", "/sources",
		"--abi",
		"--bin", fmt.Sprintf("/sources/%s", fileName))

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

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

func (cm *Manager) getSolidityVersion(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sanitizedText := scanner.Text()
		sanitizedText = strings.TrimSpace(sanitizedText)
		if !strings.HasPrefix(sanitizedText, "pragma solidity") {
			continue
		}

		versionConstraints := strings.ReplaceAll(sanitizedText, "pragma solidity", "")
		versionConstraints = strings.ReplaceAll(versionConstraints, ";", "")
		versionConstraints = strings.TrimSpace(versionConstraints)

		constraints, err := npm.NewConstraints(versionConstraints)
		if err != nil {
			return "", err
		}

		for _, version := range validSolidityVersions {
			v, _ := npm.NewVersion(version)
			if constraints.Check(v) {
				return v.String(), nil
			}
		}

		return "", fmt.Errorf("couldn't find a valid solidity version for the constraint: %s", constraints.String())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("pragma solidity not found")
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
	const basePath = "../../../test/contracts/bin"

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
