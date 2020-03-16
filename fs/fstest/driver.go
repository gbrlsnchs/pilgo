package fstest

import (
	"os"
	"reflect"

	"gsr.dev/pilgrim/fs"
)

// Driver is a stub and spy implementation of a file system's functionalities.
type Driver struct {
	// ReadDir
	ReadDirReturn map[string][]fs.FileInfo
	ReadDirErr    map[string]error

	// ReadFile
	ReadFileReturn map[string][]byte
	ReadFileErr    map[string]error

	// Stat
	StatReturn map[string]fs.FileInfo
	StatErr    map[string]error

	// WriteFile
	WriteFileErr map[string]error

	calls map[Method]Args
}

type Method interface{}
type Arg interface{}
type Args []Arg

func (drv *Driver) HasBeenCalled(fn Method) (bool, Args) {
	ptr := reflect.ValueOf(fn).Pointer()
	args, ok := drv.calls[ptr]
	return ok, args
}

func (drv *Driver) ReadDir(dirname string) ([]fs.FileInfo, error) {
	defer drv.setHasBeenCalled(drv.ReadDir, dirname)
	return drv.ReadDirReturn[dirname], drv.ReadDirErr[dirname]
}

func (drv *Driver) ReadFile(filename string) ([]byte, error) {
	defer drv.setHasBeenCalled(drv.ReadFile, filename)
	return drv.ReadFileReturn[filename], drv.ReadFileErr[filename]
}

func (drv *Driver) Stat(filename string) (fs.FileInfo, error) {
	defer drv.setHasBeenCalled(drv.Stat, filename)
	return drv.StatReturn[filename], drv.StatErr[filename]
}

func (drv *Driver) WriteFile(filename string, data []byte, perm os.FileMode) error {
	defer drv.setHasBeenCalled(drv.WriteFile, filename, data, perm)
	return drv.WriteFileErr[filename]
}

func (drv *Driver) setHasBeenCalled(fn Method, args ...Arg) {
	ptr := reflect.ValueOf(fn).Pointer()
	if drv.calls == nil {
		drv.calls = make(map[Method]Args, 5)
	}
	drv.calls[ptr] = args
}

type FileInfo struct {
	NameReturn     string
	ExistsReturn   bool
	IsDirReturn    bool
	LinknameReturn string
	PermReturn     os.FileMode
}

func (fi FileInfo) Name() string      { return fi.NameReturn }
func (fi FileInfo) Exists() bool      { return fi.ExistsReturn }
func (fi FileInfo) IsDir() bool       { return fi.IsDirReturn }
func (fi FileInfo) Linkname() string  { return fi.LinknameReturn }
func (fi FileInfo) Perm() os.FileMode { return fi.PermReturn }
