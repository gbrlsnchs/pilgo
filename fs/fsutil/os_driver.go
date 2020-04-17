package fsutil

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/andybalholm/crlf"
	"github.com/gbrlsnchs/pilgo/fs"
	"golang.org/x/text/transform"
)

// NOTE(gbrlsnchs): Declaring this variable here might prevent concurrent file reads, since
// this is a stateful transformer used between reads and reset at the beginning of reads.
var normalize = new(crlf.Normalize)

// OSDriver is the driver for a concrete file system.
type OSDriver struct{}

// MkdirAll creates directories recursively or is a NOP when they already exist.
func (OSDriver) MkdirAll(dirname string) error {
	return os.MkdirAll(dirname, 0o755)
}

// ReadDir lists names of files from dirname.
func (OSDriver) ReadDir(dirname string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	names := make([]fs.FileInfo, len(files))
	for i, fi := range files {
		mode := fi.Mode()
		info := fileInfo{
			name:   fi.Name(),
			exists: true,
			isDir:  fi.IsDir(),
			perm:   mode.Perm(),
		}
		// TODO(gbrlsnchs): add test cases
		if mode&os.ModeSymlink != 0 {
			filename := filepath.Join(dirname, info.name)
			if info.linkname, err = os.Readlink(filename); err != nil {
				return nil, err
			}
		}
		names[i] = info
	}
	return names, nil
}

// ReadFile returns the content of filename.
// It always transforms CRLF newlines into LF only.
func (OSDriver) ReadFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(transform.NewReader(f, normalize))
}

// Stat returns real information about a file.
func (OSDriver) Stat(filename string) (fs.FileInfo, error) {
	fi, err := os.Lstat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return fileInfo{exists: false}, nil
	}
	mode := fi.Mode()
	info := fileInfo{
		name:   fi.Name(),
		exists: true,
		isDir:  fi.IsDir(),
		perm:   mode.Perm(),
	}
	if mode&os.ModeSymlink != 0 {
		if info.linkname, err = os.Readlink(filename); err != nil {
			return nil, err
		}
	}
	return info, nil
}

// Symlink creates a symbolic link newname of oldname.
func (OSDriver) Symlink(oldname, newname string) error {
	// TODO(gbrlsnchs): use renameio.Symlink to replace newname, if desired.
	return os.Symlink(oldname, newname)
}

type fileInfo struct {
	name     string
	exists   bool
	isDir    bool
	linkname string
	perm     os.FileMode
}

func (fi fileInfo) Name() string      { return fi.name }
func (fi fileInfo) Exists() bool      { return fi.exists }
func (fi fileInfo) IsDir() bool       { return fi.isDir }
func (fi fileInfo) Linkname() string  { return fi.linkname }
func (fi fileInfo) Perm() os.FileMode { return fi.perm }
