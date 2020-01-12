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
	tgpath := n.Target.FullPath()
	target, err := ln.fs.Info(tgpath)
	if err != nil {
		return err
	}
	if !target.Exists() {
		n.Status = parser.StatusError
		return nil
	}
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
		return nil
	}
	if !target.IsDir() || !link.IsDir() {
		n.Status = parser.StatusConflict
		return nil
	}
	children, err := ln.fs.ReadDir(tgpath)
	if err != nil {
		return err
	}
	expand(n, children)
	n.Status = parser.StatusExpand
	return nil
}

func expand(n *parser.Node, children []string) {
	if len(children) <= 0 {
		return
	}
	n.Children = make([]*parser.Node, len(children))
	for i, c := range children {
		n.Children[i] = &parser.Node{
			Target: parser.File{
				BaseDir: n.Target.BaseDir,
				Path:    append(n.Target.Path, c),
			},
			Link: parser.File{
				BaseDir: n.Link.BaseDir,
				Path:    append(n.Link.Path, c),
			},
			Children: nil,
		}
	}
}
