package testvector_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/0xPolygonHermez/zkevm-node/zkevmprovermock/testvector"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestNewContainer(t *testing.T) {
	const defaultSourceDir = "/a/b/c"

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
    "expectedRoot": [
      "root1", "root2", "root3"
    ]
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
							ExpectedRoot: []string{"root1", "root2", "root3"},
						},
					},
				},
			},
		},
		{
			description: "invalid test vector causes error",
			sourceFiles: map[string]string{
				filepath.Join(defaultSourceDir, "merkle-tree/smt-raw.json"): "not a real json",
			},
			testVectorPath:   defaultSourceDir,
			expectedError:    true,
			expectedErrorMsg: "invalid character 'o' in literal null (expecting 'u')",
		},
		{
			description:      "unexisting test vector causes error",
			sourceFiles:      map[string]string{},
			testVectorPath:   defaultSourceDir,
			expectedError:    true,
			expectedErrorMsg: fmt.Sprintf("open %s/merkle-tree/smt-raw.json: file does not exist", defaultSourceDir),
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			var appFs = afero.NewMemMapFs()

			require.NoError(t, testutils.CreateTestFiles(appFs, tc.sourceFiles))

			actualContainer, err := testvector.NewContainer(tc.testVectorPath, appFs)

			require.NoError(t, testutils.CheckError(err, tc.expectedError, tc.expectedErrorMsg))

			require.Equal(t, tc.expectedContainer, actualContainer)
		})
	}
}
