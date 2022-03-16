package dependencies

import (
	"os"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/spf13/afero"
)

const (
	defaultTUSourceRepo    = "git@github.com:hermeznetwork/test-vectors.git"
	defaultTUTargetDirPath = "../../test/vectors/src"
)

type testVectorUpdater struct {
	fs afero.Fs

	gm *githubManager

	sourceRepo    string
	targetDirPath string
}

func init() {
	aferoFs := afero.NewOsFs()

	gm := newGithubManager(aferoFs, os.Getenv("UPDATE_DEPS_SSH_PK"), os.Getenv("GITHUB_TOKEN"))
	tv := &testVectorUpdater{
		fs: aferoFs,

		gm: gm,

		sourceRepo:    defaultTUSourceRepo,
		targetDirPath: defaultTUTargetDirPath,
	}

	dependenciesList = append(dependenciesList, tv)
}

func (tu *testVectorUpdater) update() error {
	log.Infof("Cloning %q...", tu.sourceRepo)
	tmpdir, err := tu.gm.cloneTargetRepo(tu.sourceRepo)
	if err != nil {
		return err
	}

	targetDirPath := getTargetPath(tu.targetDirPath)

	log.Infof("Updating files %q...", tu.sourceRepo)
	return updateFiles(tu.fs, tmpdir, targetDirPath)
}
