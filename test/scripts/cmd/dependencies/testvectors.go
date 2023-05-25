package dependencies

import (
	"os"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/spf13/afero"
)

// TVConfig is the configuration for the test vector updater.
type TVConfig struct {
	TargetDirPath string
	SourceRepo    string
}

type testVectorUpdater struct {
	fs afero.Fs

	gm *githubManager

	sourceRepo    string
	targetDirPath string
}

func newTestVectorUpdater(sourceRepo, targetDirPath string) *testVectorUpdater {
	aferoFs := afero.NewOsFs()

	gm := newGithubManager(aferoFs, os.Getenv("UPDATE_DEPS_SSH_PK"), os.Getenv("GITHUB_TOKEN"))

	return &testVectorUpdater{
		fs: aferoFs,

		gm: gm,

		sourceRepo:    sourceRepo,
		targetDirPath: targetDirPath,
	}
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
