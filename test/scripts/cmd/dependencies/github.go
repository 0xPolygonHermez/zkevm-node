package dependencies

import (
	"os"
	"path/filepath"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/helper/chroot"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/afero"
)

type githubManager struct {
	aferoFs afero.Fs

	sshKey string
	token  string
}

func newGithubManager(aferoFs afero.Fs, sshKey, token string) *githubManager {
	return &githubManager{
		aferoFs: aferoFs,
		token:   token,
		sshKey:  sshKey,
	}
}

func (gm *githubManager) cloneTargetRepo(repoURL string) (string, error) {
	tmpdir, err := afero.TempDir(gm.aferoFs, "", "zkevm-node-deps")
	if err != nil {
		return "", err
	}
	billyFS := newAdapter(gm.aferoFs)
	billyFS, err = billyFS.Chroot(tmpdir)
	if err != nil {
		return "", err
	}
	cloneOptions := &git.CloneOptions{
		URL: repoURL,
	}
	if gm.token != "" || gm.sshKey != "" {
		auth, err := gm.determineAuth()
		if err != nil {
			return "", err
		}
		if auth != nil {
			cloneOptions.Auth = auth
		}
	}
	storer := memory.NewStorage()
	_, err = git.Clone(storer, billyFS, cloneOptions)
	if err != nil {
		return "", err
	}
	return tmpdir, nil
}

func (gm *githubManager) determineAuth() (transport.AuthMethod, error) {
	if gm.token != "" {
		return &http.BasicAuth{
			Username: "int-bot", // this can be anything except an empty string
			Password: gm.token,
		}, nil
	}

	pvkFile, err := afero.TempFile(gm.aferoFs, "", "")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := gm.aferoFs.Remove(pvkFile.Name()); err != nil {
			log.Errorf("Could not remove temporary file %q: %v", pvkFile.Name(), err)
		}
	}()
	_, err = pvkFile.WriteString(gm.sshKey + "\n")
	if err != nil {
		return nil, err
	}
	const defaultUser = "git"
	auth, err := ssh.NewPublicKeysFromFile(defaultUser, pvkFile.Name(), "")
	if err != nil {
		return nil, err
	}

	return auth, nil
}

const (
	defaultDirectoryMode = 0755
	defaultCreateMode    = 0666
)

// AdapterFs holds an afero Fs interface for adaptation to billy.Filesystem.
type AdapterFs struct {
	fs afero.Fs
}

func newAdapter(fs afero.Fs) billy.Filesystem {
	return chroot.New(&AdapterFs{fs}, "/")
}

// Create creates a new file.
func (fs *AdapterFs) Create(filename string) (billy.File, error) {
	return fs.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, defaultCreateMode)
}

// OpenFile opens a file.
func (fs *AdapterFs) OpenFile(filename string, flag int, perm os.FileMode) (billy.File, error) {
	if flag&os.O_CREATE != 0 {
		if err := fs.createDir(filename); err != nil {
			return nil, err
		}
	}

	f, err := fs.fs.OpenFile(filename, flag, perm)
	if err != nil {
		return nil, err
	}

	mutexFile := &file{
		File: f,
	}

	return mutexFile, err
}

func (fs *AdapterFs) createDir(fullpath string) error {
	dir := filepath.Dir(fullpath)
	if dir != "." {
		if err := fs.fs.MkdirAll(dir, defaultDirectoryMode); err != nil {
			return err
		}
	}

	return nil
}

// ReadDir reads a directory.
func (fs *AdapterFs) ReadDir(path string) ([]os.FileInfo, error) {
	l, err := afero.ReadDir(fs.fs, path)
	if err != nil {
		return nil, err
	}

	var s = make([]os.FileInfo, len(l))
	copy(s, l)

	return s, nil
}

// Rename renames the given file.
func (fs *AdapterFs) Rename(from, to string) error {
	if err := fs.createDir(to); err != nil {
		return err
	}

	return os.Rename(from, to)
}

// MkdirAll creates directories recursively.
func (fs *AdapterFs) MkdirAll(path string, perm os.FileMode) error {
	return fs.fs.MkdirAll(path, defaultDirectoryMode)
}

// Open opens a file.
func (fs *AdapterFs) Open(filename string) (billy.File, error) {
	return fs.OpenFile(filename, os.O_RDONLY, 0)
}

// Stat returns information about a file.
func (fs *AdapterFs) Stat(filename string) (os.FileInfo, error) {
	return fs.fs.Stat(filename)
}

// Remove deletes a file.
func (fs *AdapterFs) Remove(filename string) error {
	return fs.fs.Remove(filename)
}

// TempFile creates a temporary file.
func (fs *AdapterFs) TempFile(dir, prefix string) (billy.File, error) {
	if err := fs.createDir(dir + string(os.PathSeparator)); err != nil {
		return nil, err
	}

	f, err := afero.TempFile(fs.fs, dir, prefix)
	if err != nil {
		return nil, err
	}
	return &file{File: f}, nil
}

// Join returns a string with joined paths.
func (fs *AdapterFs) Join(elem ...string) string {
	return filepath.Join(elem...)
}

// RemoveAll removes directories recursively.
func (fs *AdapterFs) RemoveAll(path string) error {
	return fs.fs.RemoveAll(filepath.Clean(path))
}

// Lstat returns information about a file.
func (fs *AdapterFs) Lstat(filename string) (os.FileInfo, error) {
	info, success := fs.fs.(afero.Lstater)
	if success {
		s, _, err := info.LstatIfPossible(filename)
		if err != nil {
			return nil, err
		}

		return s, nil
	}

	return fs.fs.Stat(filename)
}

// Symlink creates a symbolic link.
func (fs *AdapterFs) Symlink(target, link string) error {
	if err := fs.createDir(link); err != nil {
		return err
	}

	// TODO afero does not support symlinks
	return nil
}

// Readlink is not currently implemented.
func (fs *AdapterFs) Readlink(link string) (string, error) {
	// TODO afero does not support symlinks
	return "", nil
}

// Capabilities implements the Capable interface.
func (fs *AdapterFs) Capabilities() billy.Capability {
	return billy.DefaultCapabilities
}

// file is a wrapper for an os.File which adds support for file locking.
type file struct {
	afero.File
}

// Lock is not currently implemented.
func (f *file) Lock() error {
	return nil
}

// Unlock is not currently implemented.
func (f *file) Unlock() error {
	return nil
}
