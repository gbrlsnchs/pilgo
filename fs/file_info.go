package fs

import "os"

// FileInfo describes information about a file.
type FileInfo interface {
	Name() string
	Exists() bool
	IsDir() bool
	Linkname() string
	Perm() os.FileMode
}
