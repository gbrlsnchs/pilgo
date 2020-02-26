package fs

// FileSystem is a virtual file system that performs a set of
// operations in order for Pilgrim's core to work.
type FileSystem interface {
	Info(filename string) (FileInfo, error)
	ReadDir(dirname string) ([]string, error)
	ReadFile(filename string) ([]byte, error)
}
