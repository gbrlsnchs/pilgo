package parser

import "path/filepath"

// File holds a file's metadata.
type File struct {
	BaseDir string
	Path    []string
}

// FullPath joins base directory with all path segments.
func (f File) FullPath() string {
	fullPath := append(
		make([]string, 0, len(f.Path)+1),
		f.BaseDir,
	)
	return filepath.Join(append(fullPath, f.Path...)...)
}
