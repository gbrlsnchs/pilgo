package linker

// FileSystem describes a file system able to retrieve
// a file's metadata and create a symlink.
type FileSystem interface {
	Info(filename string) (FileInfo, error)
}
