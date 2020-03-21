package parser

import "path/filepath"

// File holds a file's metadata.
type File struct {
	BaseDir string
	Path    []string
}

// Dir returns the base directory of a file.
func (f File) Dir() string {
	pathc := len(f.Path)
	// Shortcut for the base directory itself.
	if pathc < 2 {
		return f.BaseDir
	}
	dir := append(
		make([]string, 0, len(f.Path)),
		f.BaseDir,
	)
	return filepath.Join(append(dir, f.Path[:pathc-1]...)...)
}

// FullPath joins base directory with all path segments.
func (f File) FullPath() string {
	fullPath := append(
		make([]string, 0, len(f.Path)+1),
		f.BaseDir,
	)
	return filepath.Join(append(fullPath, f.Path...)...)
}

func (f File) base() string { return f.Path[len(f.Path)-1] }
