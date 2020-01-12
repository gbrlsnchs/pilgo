package osfs

import (
	"io/ioutil"
	"os"

	"gsr.dev/pilgrim/linker"
)

// FileSystem is a OS file system that does real syscalls in order to work.
type FileSystem struct{}

var _ linker.FileSystem = new(FileSystem)

// Info returns real information about a file.
func (FileSystem) Info(filename string) (linker.FileInfo, error) {
	var info fileInfo
	fi, err := os.Lstat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return info, nil
	}
	info.exists = true
	info.isDir = fi.IsDir()
	if fi.Mode()&os.ModeSymlink != 0 {
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

type fileInfo struct {
	exists   bool
	isDir    bool
	linkname string
}

func (fi fileInfo) Exists() bool     { return fi.exists }
func (fi fileInfo) IsDir() bool      { return fi.isDir }
func (fi fileInfo) Linkname() string { return fi.linkname }
