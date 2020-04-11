package fstest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gbrlsnchs/pilgo/fs"
)

type createOpts int

const (
	mkdirOpt createOpts = 1 << iota
	overwriteOpt
)

var (
	// ErrNotExist means a file doesn't exist.
	ErrNotExist = errors.New("file doesn't exist")
	// ErrExist means a file already exists.
	ErrExist = errors.New("file already exist")
	// ErrNotDir means a file is not a directory.
	ErrNotDir = errors.New("file is not a directory")

	pathSep = string(filepath.Separator)
)

// InMemoryDriver is a synthetic file system that mimics
// simple behaviors of a real file system in memory.
type InMemoryDriver struct {
	Files map[string]File
}

// MkdirAll simulates the creation of a directory. It also creates the parents of
// such directory, if needed.
func (drv *InMemoryDriver) MkdirAll(dirname string) error {
	_, err := drv.mkdirAll(dirname)
	return err
}

// ReadDir simulates a directory read. If returns an error if it can't find dirname.
func (drv *InMemoryDriver) ReadDir(dirname string) ([]fs.FileInfo, error) {
	fstat, err := drv.find(dirname)
	if err != nil {
		return nil, err
	}
	if !fstat.IsDir() {
		return nil, ErrNotDir
	}
	files := fstat.File.Children
	list := make(sortedFiles, 0, len(files))
	for p, f := range files {
		list = append(list, FileStat{p, f})
	}
	sort.Sort(list)
	return list, nil
}

// ReadFile simulates a file read. It returns data associated with a file or an error
// if the file can't be found.
func (drv *InMemoryDriver) ReadFile(filename string) ([]byte, error) {
	fstat, err := drv.find(filename)
	if err != nil {
		return nil, err
	}
	return fstat.File.Data, nil
}

// Stat simulates reading information about filename. If filename doesn't exist, instead of
// returning an error, Stat returns an empty FileStat object.
func (drv *InMemoryDriver) Stat(filename string) (fs.FileInfo, error) {
	fstat, err := drv.find(filename)
	if errors.Is(err, ErrNotExist) {
		return FileStat{}, nil
	}
	return fstat, err
}

// Symlink simulates a symlink creation. It creates a symlink if newname doesn't exist,
// or return an error instead.
func (drv *InMemoryDriver) Symlink(oldname, newname string) error {
	f := File{
		Linkname: oldname,
		Perm:     os.ModePerm,
		Data:     nil,
		Children: nil,
	}
	return drv.create(newname, f, mkdirOpt)
}

// WriteFile simulates a file write by associating data and perm with filename.
// It overwrites files that are not directories, but still preserve the file perm.
func (drv *InMemoryDriver) WriteFile(filename string, data []byte, perm os.FileMode) error {
	f := File{
		Linkname: "",
		Perm:     perm,
		Data:     data,
		Children: nil,
	}
	return drv.create(filename, f, overwriteOpt)
}

func (drv *InMemoryDriver) create(filename string, f File, opts createOpts) error {
	fstat, err := drv.find(filename)
	if err != nil && !errors.Is(err, ErrNotExist) {
		return err
	}
	exists := fstat.Exists()
	if exists && (opts&overwriteOpt == 0 || fstat.IsDir()) {
		return ErrExist
	}
	if drv.Files == nil {
		drv.Files = make(map[string]File, 1)
	}
	parent := drv.Files
	dir, file := filepath.Split(filename)
	if dir != "" {
		dir = strings.TrimSuffix(dir, pathSep)
		fstatFn := drv.find
		if opts&mkdirOpt != 0 {
			fstatFn = drv.mkdirAll
		}
		fstat, err := fstatFn(dir)
		if err != nil {
			return err
		}
		parent = fstat.File.Children
	}
	if exists {
		f.Perm = fstat.Perm()
	}
	parent[file] = f
	return nil
}

func (drv *InMemoryDriver) find(filename string) (FileStat, error) {
	var (
		fstat FileStat
		files = drv.Files
		paths = strings.Split(filename, pathSep)
	)
	for _, p := range paths {
		if f, ok := files[p]; ok {
			files = f.Children
			fstat = FileStat{p, f}
			continue
		}
		return FileStat{}, fmt.Errorf("fstest: %s: %w", filename, ErrNotExist)
	}
	return fstat, nil
}

func (drv *InMemoryDriver) mkdirAll(dirname string) (FileStat, error) {
	if drv.Files == nil {
		drv.Files = make(map[string]File, 1)
	}
	var (
		fstat FileStat
		files = drv.Files
		paths = strings.Split(dirname, pathSep)
	)
	for _, p := range paths {
		if p == "" {
			continue
		}
		f, ok := files[p]
		if ok {
			fstat := FileStat{p, f}
			if !fstat.IsDir() {
				return FileStat{}, ErrExist
			}
		} else {
			f = File{
				Perm:     os.ModePerm,
				Linkname: "",
				Children: make(map[string]File, 0),
			}
			files[p] = f
		}
		files = f.Children
		fstat = FileStat{p, f}
	}
	return fstat, nil
}

// File is a representation of a file in a file system.
type File struct {
	Perm     os.FileMode
	Linkname string
	Data     []byte
	Children map[string]File
}

// FileStat is an abstraction of information about a file,
// which is basically the file itself plus its name.
type FileStat struct {
	Label string
	File  File
}

// Name returns a file's associated name.
func (f FileStat) Name() string { return f.Label }

// Exists returns whether a file exists.
func (f FileStat) Exists() bool { return f.Label != "" }

// IsDir returns whether a file is a directory.
func (f FileStat) IsDir() bool { return f.File.Children != nil }

// Linkname returns the name of a file a link is pointing to.
func (f FileStat) Linkname() string { return f.File.Linkname }

// Perm returns a file's associated permission.
func (f FileStat) Perm() os.FileMode { return f.File.Perm }

type sortedFiles []fs.FileInfo

func (sf sortedFiles) Len() int           { return len(sf) }
func (sf sortedFiles) Less(i, j int) bool { return sf[i].Name() < sf[j].Name() }
func (sf sortedFiles) Swap(i, j int)      { sf[i], sf[j] = sf[j], sf[i] }
