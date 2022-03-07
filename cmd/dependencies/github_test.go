package dependencies

import (
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_cloneTargetRepo(t *testing.T) {
	var appFs = afero.NewMemMapFs()

	tmpdir, err := cloneTargetRepo(appFs, "https://github.com/git-fixtures/basic.git")
	require.NoError(t, err)

	expectedChangelog := "Initial changelog\n"
	actualChangelog, err := afero.ReadFile(appFs, path.Join(tmpdir, "CHANGELOG"))
	require.NoError(t, err)

	require.Equal(t, expectedChangelog, string(actualChangelog))
}
