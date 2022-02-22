package dependencies

import "github.com/spf13/afero"

const (
	defaultSourceRepo    = "https://github.com/hermeznetwork/test-vectors"
	defaultTargetDirPath = "../../test/vectors/src"
)

type testVectorUpdater struct {
	fs afero.Fs

	sourceRepo    string
	targetDirPath string
}

func init() {
	tv := &testVectorUpdater{
		fs: afero.NewOsFs(),

		sourceRepo:    defaultSourceRepo,
		targetDirPath: defaultTargetDirPath,
	}

	dependenciesList = append(dependenciesList, tv)
}

func (tu *testVectorUpdater) update() error {
	tmpdir, err := cloneTargetRepo(tu.fs, tu.sourceRepo)
	if err != nil {
		return nil
	}

	return updateFiles(tu.fs, tmpdir, tu.targetDirPath)
}
