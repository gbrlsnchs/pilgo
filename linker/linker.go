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
	link, err := ln.fs.Info(n.Link.FullPath())
	if err != nil {
		return err
	}
	if !link.Exists() {
		n.Status = parser.StatusReady
	}
	return nil
}
