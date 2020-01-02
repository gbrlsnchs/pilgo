package linker

// FileInfo describes information about a file.
type FileInfo interface {
	Exists() bool
	Linkname() string
}
