package dependencies

import (
	"os"
	"os/exec"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/spf13/afero"
)

// PBConfig is the configuration for the protobuffers updater.
type PBConfig struct {
	SourceRepo    string
	TargetDirPath string
}

type pbUpdater struct {
	fs afero.Fs

	gm *githubManager

	sourceRepo    string
	targetDirPath string
}

func newPBUpdater(sourceRepo, targetDirPath string) *pbUpdater {
	aferoFs := afero.NewOsFs()

	gm := newGithubManager(aferoFs, os.Getenv("UPDATE_DEPS_SSH_PK"), os.Getenv("GITHUB_TOKEN"))

	return &pbUpdater{
		fs: aferoFs,

		gm: gm,

		sourceRepo:    sourceRepo,
		targetDirPath: targetDirPath,
	}
}

func (pb *pbUpdater) update() error {
	log.Infof("Cloning %q...", pb.sourceRepo)
	tmpdir, err := pb.gm.cloneTargetRepo(pb.sourceRepo)
	if err != nil {
		return err
	}

	targetDirPath := getTargetPath(pb.targetDirPath)

	log.Infof("Updating files %q...", pb.sourceRepo)
	err = updateFiles(pb.fs, tmpdir, targetDirPath)
	if err != nil {
		return err
	}

	log.Infof("Generating stubs from proto files...")

	c := exec.Command("make", "generate-code-from-proto")
	c.Dir = "."
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
