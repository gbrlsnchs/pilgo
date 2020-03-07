package fs

import (
	"io/ioutil"
	"os"

	"github.com/andybalholm/crlf"
	"golang.org/x/text/transform"
)

// NOTE(gbrlsnchs): Declaring this variable here might prevent concurrent file reads, since
// this is a stateful transformer used between reads and reset at the beginning of reads.
var normalize = new(crlf.Normalize)

// FileSystem is a concrete file system that implements a VFS contract.
type FileSystem struct{}

// Info returns real information about a file.
func (FileSystem) Info(filename string) (FileInfo, error) {
	fi, err := os.Lstat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return fileInfo{exists: false}, nil
	}
	mode := fi.Mode()
	info := fileInfo{
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

// ReadDir lists names of files from dirname.
func (FileSystem) ReadDir(dirname string) ([]string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(files))
	for i, f := range files {
		names[i] = f.Name()
	}
	return names, nil
}

// ReadFile returns the content of filename.
// It always transforms CRLF newlines into LF only.
func (FileSystem) ReadFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(transform.NewReader(f, normalize))
}

// WriteFile writes data to filename with permission perm.
func (FileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	// TODO(gbrlsnchs): use package "renameio"
	return ioutil.WriteFile(filename, data, perm)
}

type fileInfo struct {
	exists   bool
	isDir    bool
	linkname string
	perm     os.FileMode
}

func (fi fileInfo) Exists() bool      { return fi.exists }
func (fi fileInfo) IsDir() bool       { return fi.isDir }
func (fi fileInfo) Linkname() string  { return fi.linkname }
func (fi fileInfo) Perm() os.FileMode { return fi.perm }
