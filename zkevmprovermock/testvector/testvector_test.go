package testvector_test

import (
	"path/filepath"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/0xPolygonHermez/zkevm-node/zkevmprovermock/testvector"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestNewContainer(t *testing.T) {
	const defaultSourceDir = "/a/b/c"

	var appFs = afero.NewMemMapFs()

	tcs := []struct {
		description       string
		sourceFiles       map[string]string
		testVectorPath    string
		expectedContainer *testvector.Container
		expectedError     bool
		expectedErrorMsg  string
	}{
		{
			description: "happy path",
			sourceFiles: map[string]string{
				filepath.Join(defaultSourceDir, "merkle-tree/smt-raw.json"): `[
  {
    "keys": [
      "a","b","c"
    ],
    "values": [
      "1","2","3"
    ],
    "expectedRoot": "root"
  }
]`,
			},
			testVectorPath: defaultSourceDir,
			expectedContainer: &testvector.Container{
				StateDBRaw: &testvector.StateDBRaw{
					Entries: []*testvector.StateDBRawEntry{
						{
							Keys:         []string{"a", "b", "c"},
							Values:       []string{"1", "2", "3"},
							ExpectedRoot: "root",
						},
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			require.NoError(t, testutils.CreateTestFiles(appFs, tc.sourceFiles))

			actualContainer, err := testvector.NewContainer(tc.testVectorPath, appFs)

			require.NoError(t, testutils.CheckError(err, tc.expectedError, tc.expectedErrorMsg))

			require.Equal(t, tc.expectedContainer, actualContainer)
		})
	}
}
