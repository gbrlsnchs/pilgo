package linker

import "gsr.dev/pilgrim/parser"

// Linker is a file symlinker.
type Linker struct {
	fs FileSystem
}

// New creates a new linker with a given file system.
func New(fs FileSystem) *Linker { return &Linker{fs} }

// Resolve checks and resolves nodes in a parsed tree.
func (ln *Linker) Resolve(n *parser.Node) error {
	if len(n.Children) > 0 {
		n.Status = parser.StatusSkip
		return nil
	}
	lnpath := n.Link.FullPath()
	link, err := ln.fs.Info(lnpath)
	if err != nil {
		return err
	}
	if !link.Exists() {
		n.Status = parser.StatusReady
		return nil
	}
	if linkname := link.Linkname(); linkname != "" {
		if linkname == lnpath {
			n.Status = parser.StatusDone
			return nil
		}
		n.Status = parser.StatusConflict
	}
	if !link.IsDir() {
		n.Status = parser.StatusConflict
	}
	return nil
}
