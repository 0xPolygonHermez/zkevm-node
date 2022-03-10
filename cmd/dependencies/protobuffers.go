package dependencies

import (
	"os"
	"os/exec"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/spf13/afero"
)

const (
	defaultPBSourceRepo    = "git@github.com:hermeznetwork/comms-protocol.git"
	defaultPBTargetDirPath = "../../pb/src"
)

type pbUpdater struct {
	fs afero.Fs

	gm *githubManager

	sourceRepo    string
	targetDirPath string
}

func init() {
	aferoFs := afero.NewOsFs()

	gm := newGithubManager(aferoFs, os.Getenv("UPDATE_DEPS_SSH_PK"), os.Getenv("GITHUB_TOKEN"))

	pb := &pbUpdater{
		fs: aferoFs,

		gm: gm,

		sourceRepo:    defaultPBSourceRepo,
		targetDirPath: defaultPBTargetDirPath,
	}

	dependenciesList = append(dependenciesList, pb)
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
