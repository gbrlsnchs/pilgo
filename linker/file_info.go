package linker

// FileInfo describes information about a file.
type FileInfo interface {
	Exists() bool
	IsDir() bool
	Linkname() string
	List() ([]string, error)
}
