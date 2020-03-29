package fstest

import (
	"os"
	"reflect"

	"github.com/gbrlsnchs/pilgo/fs"
)

// SpyDriver is a stub and spy implementation of a file system's functionalities.
type SpyDriver struct {
	// MkdirAll
	MkdirAllErr map[string]error

	// ReadDir
	ReadDirReturn map[string][]fs.FileInfo
	ReadDirErr    map[string]error

	// ReadFile
	ReadFileReturn map[string][]byte
	ReadFileErr    map[string]error

	// Stat
	StatReturn map[string]fs.FileInfo
	StatErr    map[string]error

	// Symlink
	SymlinkErr map[string]error

	// WriteFile
	WriteFileErr map[string]error

	calls map[Method]CallStack
}

type Method interface{}
type Arg interface{}
type Args []Arg
type CallStack []Args

func (drv *SpyDriver) HasBeenCalled(fn Method) (bool, CallStack) {
	ptr := reflect.ValueOf(fn).Pointer()
	args, ok := drv.calls[ptr]
	return ok, args
}

// MkdirAll returns a stub of directory creation.
func (drv *SpyDriver) MkdirAll(dirname string) error {
	defer drv.setHasBeenCalled(drv.MkdirAll, dirname)
	return drv.MkdirAllErr[dirname]
}

func (drv *SpyDriver) ReadDir(dirname string) ([]fs.FileInfo, error) {
	defer drv.setHasBeenCalled(drv.ReadDir, dirname)
	return drv.ReadDirReturn[dirname], drv.ReadDirErr[dirname]
}

func (drv *SpyDriver) ReadFile(filename string) ([]byte, error) {
	defer drv.setHasBeenCalled(drv.ReadFile, filename)
	return drv.ReadFileReturn[filename], drv.ReadFileErr[filename]
}

func (drv *SpyDriver) Stat(filename string) (fs.FileInfo, error) {
	defer drv.setHasBeenCalled(drv.Stat, filename)
	return drv.StatReturn[filename], drv.StatErr[filename]
}

// Symlink returns a stub of a symlink creation.
func (drv *SpyDriver) Symlink(oldname, newname string) error {
	defer drv.setHasBeenCalled(drv.Symlink, oldname, newname)
	return drv.SymlinkErr[oldname]
}

func (drv *SpyDriver) WriteFile(filename string, data []byte, perm os.FileMode) error {
	defer drv.setHasBeenCalled(drv.WriteFile, filename, data, perm)
	return drv.WriteFileErr[filename]
}

func (drv *SpyDriver) setHasBeenCalled(fn Method, args ...Arg) {
	ptr := reflect.ValueOf(fn).Pointer()
	if drv.calls == nil {
		drv.calls = make(map[Method]CallStack, 5)
	}
	drv.calls[ptr] = append(drv.calls[ptr], args)
}

type StubFile struct {
	NameReturn     string
	ExistsReturn   bool
	IsDirReturn    bool
	LinknameReturn string
	PermReturn     os.FileMode
}

func (fi StubFile) Name() string      { return fi.NameReturn }
func (fi StubFile) Exists() bool      { return fi.ExistsReturn }
func (fi StubFile) IsDir() bool       { return fi.IsDirReturn }
func (fi StubFile) Linkname() string  { return fi.LinknameReturn }
func (fi StubFile) Perm() os.FileMode { return fi.PermReturn }
