package fs

import (
	"errors"
	"os"
)

// ErrNoDriver is the error for a nonfunctional file system.
var ErrNoDriver = errors.New("fs: nil driver")

// FileSystem is a concrete file system that implements a VFS contract.
type FileSystem struct {
	drv Driver
}

// Driver is the internal file system implementation.
type Driver interface {
	ReadDir(dirname string) ([]FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	Stat(filename string) (FileInfo, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// New creates a new FileSystem with drv as its engine.
func New(drv Driver) FileSystem {
	return FileSystem{drv}
}

// ReadDir lists names of files from dirname.
func (fs FileSystem) ReadDir(dirname string) ([]FileInfo, error) {
	fs.testDriver()
	return fs.drv.ReadDir(dirname)
}

// ReadFile returns the content of filename.
func (fs FileSystem) ReadFile(filename string) ([]byte, error) {
	fs.testDriver()
	return fs.drv.ReadFile(filename)
}

// Stat returns information about a file.
func (fs FileSystem) Stat(filename string) (FileInfo, error) {
	fs.testDriver()
	return fs.drv.Stat(filename)
}

// WriteFile writes data to filename with permission perm.
func (fs FileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	fs.testDriver()
	return fs.drv.WriteFile(filename, data, perm)
}

func (fs FileSystem) testDriver() {
	if fs.drv == nil {
		panic(ErrNoDriver)
	}
}
