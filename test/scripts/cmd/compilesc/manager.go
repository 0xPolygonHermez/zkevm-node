package compilesc

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

const (
	defaultAbigenImage = "ethereum/client-go:alltools-latest"
	defaultSolcImage   = "ethereum/solc"

	containerBase = "/contracts"
)

// CompileUnit represents a single contract to be compiled.
type CompileUnit struct {
	Name        string `yaml:"name"`
	SolcVersion string `yaml:"solcVersion"`
	InputPath   string `yaml:"inputPath"`
	OutputPath  string `yaml:"outputPath"`
}

// CompileUnits represents data for all the contracts to be compiled.
type CompileUnits struct {
	Parallel   []CompileUnit `yaml:"parallel"`
	Sequential []CompileUnit `yaml:"sequential"`
}

type compileIndex struct {
	CompileUnits CompileUnits `yaml:"compileUnits"`
}

// Manager handles smart contract compilation.
type Manager struct {
	basePath         string
	absoluteBasePath string
	currentUser      *user.User
}

// NewManager is the Manager constructor.
func NewManager(basePath string) (*Manager, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../../")
	absoluteBasePath := path.Join(dir, basePath)

	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	return &Manager{
		basePath:         basePath,
		absoluteBasePath: absoluteBasePath,
		currentUser:      currentUser,
	}, nil
}

// Run executes the compilation of smart contracts and generation of golang
// bindings.
func (cm *Manager) Run() error {
	yamlFile, err := os.ReadFile(path.Join(cm.absoluteBasePath, "index.yaml"))
	if err != nil {
		return err
	}
	ci := &compileIndex{}
	err = yaml.Unmarshal(yamlFile, ci)
	if err != nil {
		return err
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		for _, item := range ci.CompileUnits.Parallel {
			err = cm.parallelActions(item)
			if err != nil {
				return err
			}
		}
		return nil
	})

	g.Go(func() error {
		for _, item := range ci.CompileUnits.Sequential {
			err = cm.sequentialActions(item)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return g.Wait()
}

// parallelActions performs compile and generate actions in parallel for a given
// compile unit.
func (cm *Manager) parallelActions(item CompileUnit) error {
	entryPoint := path.Join(cm.basePath, item.InputPath)
	file, err := os.Open(entryPoint) // #nosec G304
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
		return cm.fileActions(path.Base(file.Name()), item.SolcVersion, item.InputPath, item.OutputPath)
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		return filepath.WalkDir(file.Name(), func(target string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if info == nil || info.IsDir() {
				return nil
			}
			if filepath.Ext(target) != ".sol" {
				return nil
			}
			fileName := strings.TrimSuffix(target, ".sol")

			g.Go(func() error {
				return cm.fileActions(path.Base(fileName), item.SolcVersion, item.InputPath, item.OutputPath)
			})
			return nil
		})
	})
	return g.Wait()
}

// sequentialActions performs compile and generate actions sequentially for a given
// compile unit.
func (cm *Manager) sequentialActions(item CompileUnit) error {
	entryPoint := path.Join(cm.basePath, item.InputPath, fmt.Sprintf("%s.sol", item.Name))
	file, err := os.Open(entryPoint) // #nosec G304
	if err != nil && err != os.ErrNotExist {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Errorf("Could not close file %q, %v", file.Name(), err)
		}
	}()

	return cm.fileActions(item.Name, item.SolcVersion, item.InputPath, item.OutputPath)
}

func (cm *Manager) fileActions(name, solcVersion, inputPath, outputPath string) error {
	err := cm.Compile(path.Base(name), solcVersion, inputPath, outputPath)
	if err != nil {
		return err
	}

	return cm.Abigen(path.Base(name), inputPath, outputPath)
}

// Compile invokes solc on the given sol file.
func (cm *Manager) Compile(name, solcVersion, inputPath, outputPath string) error {
	log.Infof("Compiling %s.sol with version %s", path.Join(cm.basePath, inputPath, name), solcVersion)

	solcImage := fmt.Sprintf("%s:%s", defaultSolcImage, solcVersion)

	c := exec.Command(
		"docker", "run", "--rm",
		"--user", fmt.Sprintf("%s:%s", cm.currentUser.Uid, cm.currentUser.Gid),
		"-v", fmt.Sprintf("%s:%s", cm.absoluteBasePath, containerBase),
		solcImage,
		"-",
		fmt.Sprintf("%s.sol", path.Join(containerBase, inputPath, name)),
		"-o", path.Join(containerBase, "bin", outputPath, name),
		"--abi", "--bin", "--overwrite", "--optimize") // #nosec G204

	envPath := os.Getenv("PATH")
	c.Env = []string{fmt.Sprintf("PATH=%s", envPath)}

	err := c.Run()
	if err != nil {
		return err
	}
	log.Infof("Compiler run successfully, artifacts can be found at %q", path.Join(cm.basePath, "bin", name))
	return nil
}

// Abigen generates bindings for the given file
func (cm *Manager) Abigen(name, inputPath, outputPath string) error {
	log.Infof("Generating go code for %q...", name)

	c := exec.Command(
		"docker", "run", "--rm",
		"--user", fmt.Sprintf("%s:%s", cm.currentUser.Uid, cm.currentUser.Gid),
		"-v", fmt.Sprintf("%s:%s", cm.absoluteBasePath, containerBase),
		defaultAbigenImage,
		"abigen",
		"--bin", path.Join(containerBase, "bin", outputPath, name, fmt.Sprintf("%s.bin", name)),
		"--abi", path.Join(containerBase, "bin", outputPath, name, fmt.Sprintf("%s.abi", name)),
		"--pkg", name,
		"--out", path.Join(containerBase, "bin", outputPath, name, fmt.Sprintf("%s.go", name))) // #nosec G204

	envPath := os.Getenv("PATH")
	c.Env = []string{fmt.Sprintf("PATH=%s", envPath)}

	err := c.Run()
	if err != nil {
		return err
	}
	log.Infof("Code generated at %q", path.Join(cm.basePath, "bin", outputPath, name, fmt.Sprintf("%s.go", name)))
	return nil
}
