package dependencies

import (
	"testing"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_updateFiles(t *testing.T) {
	const (
		defaultTargetDir = "/a/b/src"
		defaultSourceDir = "/tmp/src"
	)

	var appFs = afero.NewMemMapFs()

	tcs := []struct {
		description         string
		initialSourceFiles  map[string]string
		initialTargetFiles  map[string]string
		expectedTargetFiles map[string]string
	}{
		{
			description: "single file matching file",
			initialSourceFiles: map[string]string{
				"/tmp/src/a": "new-a-content",
			},
			initialTargetFiles: map[string]string{
				"/a/b/src/a": "old-a-content",
			},
			expectedTargetFiles: map[string]string{
				"/a/b/src/a": "new-a-content",
			},
		},
		{
			description: "single file matching file with non-matching files",
			initialSourceFiles: map[string]string{
				"/tmp/src/a": "new-a-content",
				"/tmp/src/b": "new-b-content",
			},
			initialTargetFiles: map[string]string{
				"/a/b/src/a": "old-a-content",
			},
			expectedTargetFiles: map[string]string{
				"/a/b/src/a": "new-a-content",
			},
		},
		{
			description: "multiple matching files",
			initialSourceFiles: map[string]string{
				"/tmp/src/a.json":                 "new-a-content",
				"/tmp/src/subdir1/subdir2/b.json": "new-b-content",
			},
			initialTargetFiles: map[string]string{
				"/a/b/src/a.json":                 "old-a-content",
				"/a/b/src/subdir1/subdir2/b.json": "old-b-content",
			},
			expectedTargetFiles: map[string]string{
				"/a/b/src/a.json":                 "new-a-content",
				"/a/b/src/subdir1/subdir2/b.json": "new-b-content",
			},
		},
		{
			description: "multiple matching files with non matching files",
			initialSourceFiles: map[string]string{
				"/tmp/src/subdira1/a.json":          "new-a-content",
				"/tmp/src/subdirb1/subdirb2/b.json": "new-b-content",
				"/tmp/src/c.json":                   "new-c-content",
			},
			initialTargetFiles: map[string]string{
				"/a/b/src/subdira1/a.json":          "old-a-content",
				"/a/b/src/subdirb1/subdirb2/b.json": "old-b-content",
			},
			expectedTargetFiles: map[string]string{
				"/a/b/src/subdira1/a.json":          "new-a-content",
				"/a/b/src/subdirb1/subdirb2/b.json": "new-b-content",
			},
		},
		{
			description: "unexisting target file does not give error",
			initialSourceFiles: map[string]string{
				"/tmp/src/subdira1/a.json":          "new-a-content",
				"/tmp/src/subdirb1/subdirb2/b.json": "new-b-content",
				"/tmp/src/c.json":                   "new-c-content",
			},
			initialTargetFiles: map[string]string{
				"/a/b/src/subdira1/a.json":        "old-a-content",
				"/a/b/src/subdir1/subdir2/d.json": "old-d-content",
			},
			expectedTargetFiles: map[string]string{
				"/a/b/src/subdira1/a.json":        "new-a-content",
				"/a/b/src/subdir1/subdir2/d.json": "old-d-content",
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			require.NoError(t, createTestFiles(appFs, tc.initialSourceFiles))
			require.NoError(t, createTestFiles(appFs, tc.initialTargetFiles))

			require.NoError(t, updateFiles(appFs, defaultSourceDir, defaultTargetDir))
			a := afero.Afero{Fs: appFs}
			for path, expectedContent := range tc.expectedTargetFiles {
				actualContent, err := a.ReadFile(path)
				require.NoError(t, err)
				require.Equal(t, expectedContent, string(actualContent))
			}
			require.NoError(t, appFs.RemoveAll(defaultSourceDir))
			require.NoError(t, appFs.RemoveAll(defaultTargetDir))
		})
	}
}

func createTestFiles(appFs afero.Fs, files map[string]string) error {
	for path, content := range files {
		f, err := appFs.Create(path)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Errorf("Could not close %s: %v", f.Name(), err)
			}
		}()
		_, err = f.WriteString(content)

		if err != nil {
			return err
		}
	}
	return nil
}
