package dependencies

import (
	"io"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/hermeznetwork/hermez-core/log"
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
		if err != nil {
			return err
		}
		defer func() {
			if err := sourceFile.Close(); err != nil {
				log.Errorf("Could not close %s: %v", sourceFile.Name(), err)
			}
		}()
		destinationFile, err := fs.OpenFile(wpath, os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer func() {
			if err := destinationFile.Close(); err != nil {
				log.Errorf("Could not close %s: %v", destinationFile.Name(), err)
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
			if _, err := destinationFile.Write(buf[:n]); err != nil {
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
