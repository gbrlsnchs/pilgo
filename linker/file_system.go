package linker

import "gsr.dev/pilgrim/fs"

// FileSystem is an abstraction to read nodes.
type FileSystem interface {
	Info(string) (fs.FileInfo, error)
	ReadDir(string) ([]string, error)
}
