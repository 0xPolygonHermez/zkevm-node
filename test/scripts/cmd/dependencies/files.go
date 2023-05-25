package dependencies

import (
	"io"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/spf13/afero"
)

func updateFiles(fs afero.Fs, sourceDir, targetDir string) error {
	const bufferSize = 20
	err := afero.Walk(fs, targetDir, func(wpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info == nil || info.IsDir() {
			return nil
		}
		relativePath := strings.Replace(wpath, targetDir, "", -1)
		sourcePath := path.Join(sourceDir, relativePath)

		sourceFile, err := fs.Open(sourcePath)
		if os.IsNotExist(err) {
			// we allow source files to not exist, for instance, test vectors that we
			// have in zkevm-node but are not present in the upstream repo
			return nil
		}
		if err != nil {
			return err
		}
		defer func() {
			if err := sourceFile.Close(); err != nil {
				log.Errorf("Could not close %s: %v", sourceFile.Name(), err)
			}
		}()
		targetFile, err := fs.OpenFile(wpath, os.O_RDWR|os.O_TRUNC, 0644) //nolint:gomnd
		if err != nil {
			return err
		}
		defer func() {
			if err := targetFile.Close(); err != nil {
				log.Errorf("Could not close %s: %v", targetFile.Name(), err)
			}
		}()
		buf := make([]byte, bufferSize)
		for {
			n, err := sourceFile.Read(buf)
			if err != nil && err != io.EOF {
				return err
			}
			if n == 0 {
				break
			}
			if _, err := targetFile.Write(buf[:n]); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func getTargetPath(targetPath string) string {
	if strings.HasPrefix(targetPath, "/") {
		return targetPath
	}
	_, filename, _, _ := runtime.Caller(1)

	return path.Join(path.Dir(filename), targetPath)
}
